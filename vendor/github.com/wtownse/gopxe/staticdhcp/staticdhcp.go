package staticdhcp

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/coredhcp/coredhcp/logger"
	mux "github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	pl_range "github.com/wtownse/gopxe/coredhcp/plugins/range"
	"github.com/wtownse/gopxe/handlers"
)

var log = logger.GetLogger("staticdhcp")

// StaticRecords holds a MAC -> IP address mapping
var StaticRecords map[string]net.IP

type macip struct {
	Mac string `json:"mac"`
	Ip  string `json:"ip"`
}

type dhclientInfo struct {
	Arch     string
	Ip       string
	Hostname string
}

var recLock sync.RWMutex

func mkStaticDhcpEntry(path string, append string) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Printf("Cannot create file %v", err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(append)
	if err != nil {
		log.Printf("Error %v", err)
		return err
	}
	defer file.Close()

	return err
}
func clearfile(path string) error {
	file, err := os.OpenFile(path, os.O_TRUNC, 0644)

	if err != nil {
		log.Printf("Cannot create file %v", err)
		return err
	}
	defer file.Close()

	return err
}
func LoadDHCPv4Records(filename string) (map[string]net.IP, error) {
	log.Infof("reading leases from %s", filename)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	records := make(map[string]net.IP)
	for _, lineBytes := range bytes.Split(data, []byte{'\n'}) {
		line := string(lineBytes)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		tokens := strings.Fields(line)
		if len(tokens) != 2 {
			return nil, fmt.Errorf("malformed line, want 2 fields, got %d: %s", len(tokens), line)
		}
		hwaddr, err := net.ParseMAC(tokens[0])
		if err != nil {
			return nil, fmt.Errorf("malformed hardware address: %s", tokens[0])
		}
		ipaddr := net.ParseIP(tokens[1])
		if ipaddr.To4() == nil {
			return nil, fmt.Errorf("expected an IPv4 address, got: %v", ipaddr)
		}
		records[hwaddr.String()] = ipaddr
	}
	recLock.Lock()
	defer recLock.Unlock()

	StaticRecords = records

	return records, nil
}
func GetDHCPEntry(rw http.ResponseWriter, req *http.Request) {
	p := pl_range.PS
	vars := mux.Vars(req)
	mac := vars["mac"]
	ip := p.Recordsv4[mac]
	arch := p.Recordsv4[mac].Arch
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.WriteHeader(http.StatusOK)
	io.WriteString(rw, ip.IP.String())
	log.Printf("MAC: %v IP: %v Arch: %v", mac, ip, arch)

}
func GetDHCPEntriesJSON(rw http.ResponseWriter, req *http.Request) {
	p := pl_range.PS
	macips := make(map[string]dhclientInfo)
	for k, v := range p.Recordsv4 {
		macips[k] = dhclientInfo{
			Arch:     v.Arch,
			Ip:       v.IP.String(),
			Hostname: v.Hostname,
		}
	}
	jsonData, err := json.Marshal(macips)
	if err == nil {
		rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
		rw.WriteHeader(http.StatusOK)
		rw.Write(jsonData)
		log.Printf("DHCP to IP mappings: %v", (string(jsonData)))
	}
}
func GetDHCPEntries(rw http.ResponseWriter, req *http.Request) {
	p := pl_range.PS
	macips := make(map[string]dhclientInfo)
	for k, v := range p.Recordsv4 {
		macips[k] = dhclientInfo{
			Arch:     v.Arch,
			Ip:       v.IP.String(),
			Hostname: v.Hostname,
		}
	}
	if err := handlers.Templates["pxeinfo"].Execute(rw, macips); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
func StaticDHCP(rw http.ResponseWriter, req *http.Request) {
	d := json.NewDecoder(req.Body)
	var macip macip
	err := d.Decode(&macip)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("StaticRecord: %v", macip)
	log.Printf("%v %v", macip.Mac, macip.Ip)
	// validate mac address and return if invalid
	if _, err := net.ParseMAC(macip.Mac); err != nil {
		log.Printf("%v is an invalid mac address", macip.Mac)
		return
	}
	// validate ip address and return if invalid
	if err := net.ParseIP(macip.Ip); err != nil {
		log.Printf("%v is an invalid ip address", macip.Ip)
		return
	}
	p := pl_range.PS

	db := p.Leasedb
	if err != nil {
		log.Printf("Failed to load database: %v", err)
	}
	records, err := pl_range.LoadRecords(db)
	if err != nil {
		log.Printf("Failed to load records: %v", err)
	}
	//log.Printf("Loaded %d DHCPv4 leases from %s", len(records), "leases.db")
	p.Recordsv4 = records

	pl_Record := pl_range.Record{
		IP:       net.ParseIP(macip.Ip),
		Expires:  int(time.Now().Add(p.LeaseTime).Unix()),
		Hostname: "",
		Arch:     "",
	}
	addr, _ := net.ParseMAC(macip.Mac)
	_, ok := p.Recordsv4[addr.String()]
	if ok {
		log.Printf("MAC address already exists")
	} else {
		log.Printf("MAC address does not exist")
	}

	p.Recordsv4[macip.Mac] = &pl_Record

	err = p.SaveIPAddress(addr, &pl_Record)
	p.UpdateRecordsv4(pl_Record, addr)
	if err != nil {
		log.Printf("Failed to save IP address: %v", err)
	}

	log.Printf("StaticRecord: %v", macip)
	log.Printf("%v %v", macip.Mac, macip.Ip)
	for _, m := range p.Recordsv4 {
		log.Printf("leases", m)
	}
	LoadDHCPv4Records("leases.txt")
	if StaticRecords[macip.Mac] != nil && StaticRecords[macip.Mac].String() == macip.Ip {
		log.Printf("MAC address already exists")
		return
	} else {
		log.Printf("MAC address does not exist")
		StaticRecords[macip.Mac] = net.ParseIP(macip.Ip)
	}

	Updatemacaddresses()

}
func Updatemacaddresses() {
	clearfile("leases.txt")
	for mac, ip := range StaticRecords {
		mkStaticDhcpEntry("leases.txt", mac+" "+ip.String()+"\n")
	}
}
func loadDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", path))
	if err != nil {
		return nil, fmt.Errorf("failed to open database (%T): %w", err, err)
	}
	return db, nil
}
