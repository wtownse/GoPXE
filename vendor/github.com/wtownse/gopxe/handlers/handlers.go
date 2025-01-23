package handlers

import (
	"encoding/json"
	//"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/coredhcp/coredhcp/logger"
	"github.com/gorilla/mux"

	db "github.com/wtownse/gopxe/bbolt"
	//db "github.com/coreos/bbolt"
	flag "github.com/spf13/pflag"
)

var (
	conn      db.BoltDB
	Templates map[string]*template.Template
)

const (
	tftpROOT = "/var/lib/tftpboot/"
)

var log = logger.GetLogger("gopxe")

type PXEBOOTTYPE struct {
	BootAction string `json:"bootaction"`
	KSFile     string `json:"ksfile"`
	OS         string `json:"os"`
	Version    string `json:"version"`
	Hostname   string `json:"hostname"`
	IP         string `json:"ip"`
	MASK       string `json:"mask"`
	NS1        string `json:"ns1"`
	NS2        string `json:"ns2"`
	GW         string `json:"gw"`
	UUID       string `json:"uuid"`
}

type ACTIONTYPE struct {
	Default     string `json:"default"`
	Label       string `json:"label"`
	Menu        string `json:"menu"`
	Kernel      string `json:"kernel"`
	KSDevice    string `json:"ksdevice"`
	IP          string `json:"ip"`
	LoadRamdisk string `json:"load_ramdisk"`
	Initrd      string `json:"initrd"`
}

func LoadTemplates() {
	var baseTemplate = "public/layouts/base.html"
	Templates = make(map[string]*template.Template)
	Templates["index"] = template.Must(template.ParseFiles(baseTemplate, "public/pages/index.html"))
	Templates["actions"] = template.Must(template.ParseFiles(baseTemplate, "public/pages/bootactions.html"))
	Templates["pxeboot"] = template.Must(template.ParseFiles(baseTemplate, "public/pages/pxeboot.html"))
	Templates["pxeinfo"] = template.Must(template.ParseFiles(baseTemplate, "public/pages/pxeinfo.html"))
	Templates["createbootaction"] = template.Must(template.ParseFiles(baseTemplate, "public/pages/createbootaction.html"))
}

func getBucket() string {
	return flag.Lookup("bucket").Value.String()
}

func getEFIBucket() string {
	return flag.Lookup("efibucket").Value.String()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func isExists(f string) bool {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return false
	}
	return true
}

func KsGenerate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	checkError(err)

	profile := r.Form.Get("name")

	tmplFile := "ksTempl/" + profile + ".tmpl"
	tp, err := ioutil.ReadFile(tmplFile)
	checkError(err)

	data := make(map[string]string)
	for i, j := range r.Form {
		data[i] = j[0]
	}

	t, err := template.New("index").Parse(string(tp))
	checkError(err)

	err = t.Execute(w, data)
	checkError(err)
}
func ShowTest(w http.ResponseWriter, r *http.Request) {
	var action ACTIONTYPE
	vars2 := mux.Vars(r)
	key := vars2["key"]
	r.ParseForm()

	vars := r.Form
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &action); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	log.Printf("Key %s\n", key)
	log.Printf("Action %s\n", action.Default)
	log.Println(vars)
	log.Println(vars["default"][0])
	value := fmt.Sprintf("default %s\n label %s\n MENU LABEL %s\n KERNEL %s\n APPEND ksdevice=%s ip=%s load_ramdisk=%s initrd=%s", vars["default"][0], vars["label"][0], vars["menu"][0], vars["kernel"][0], vars["ksdevice"][0], vars["ip"][0], vars["load_ramdisk"][0], vars["initrd"][0])
	err1 := conn.PutBootAction(getBucket(), key, value)
	if err1 != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"Status": "Couldn't store bootaction"}`)
		log.Printf("Error %s\n", err1)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"Status":"bootaction recorded"}`)
		log.Printf("Bootaction %s recorded\n", key)
	}
}
func GetAllBA(w http.ResponseWriter, r *http.Request) {
	var actions map[string]string

	err, actions := conn.GetAllBootActions(getBucket())
	getEFIBucket()
	if err != nil {
		log.Printf("Couldn't retrieve bootaction %v \n", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"Status": "Couldn't retrieve bootaction"}`)
		fmt.Fprintln(w, err)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		for i, j := range actions {
			fmt.Fprintln(w, "Bootaction name:", i)
			io.WriteString(w, j)
			//fmt.Fprintln(w, "\n")
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}

func GetBA(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	err, v := conn.GetBootAction(getBucket(), key)
	getEFIBucket()
	if err != nil {
		log.Printf("Couldn't retrieve bootaction %v \n", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"Status": "Couldn't retrieve bootaction"}`)
		fmt.Fprintln(w, err)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, v)
		log.Printf("Getting this for v %s", v)
		return
	}
}

func PutBA(w http.ResponseWriter, r *http.Request) {
	var action ACTIONTYPE
	vars := mux.Vars(r)
	key := vars["key"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &action); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	tftpPath := flag.Lookup("tftpPath").Value.String()

	if !isExists(tftpPath) {
		http.Error(w, "Couldnt store bootaction", http.StatusNotFound)
		io.WriteString(w, `{"Status": "Couldn't store bootaction, tftpd path doesn't exist"}`)
		return
	}
	value := fmt.Sprintf("default %s\n label %s\n MENU LABEL %s\n KERNEL %s\n APPEND ksdevice=%s ip=%s load_ramdisk=%s initrd=%s", action.Default, action.Label, action.Menu, action.Kernel, action.KSDevice, action.IP, action.LoadRamdisk, action.Initrd)
	valueEFI := fmt.Sprintf(" set default=0\n set timeout=60\n menuentry %s {\n linuxeefi %s\n initrdefi %s\n}", action.Label, action.Kernel, action.Initrd)
	log.Printf(valueEFI)
	log.Printf(getEFIBucket())
	err1 := conn.PutBootAction(getBucket(), key, value)
	err1 = conn.PutBootAction(getEFIBucket(), key, valueEFI)
	if err1 != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"Status": "Couldn't store bootaction"}`)
		log.Printf("Error %s\n", err1)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"Status":"bootaction recorded"}`)
		log.Printf("Bootaction %s recorded\n", key)
	}

}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"Status": alive}`)
}

func Index(w http.ResponseWriter, r *http.Request) {

	if err := Templates["index"].Execute(w, ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func CreateBootAction(w http.ResponseWriter, r *http.Request) {

	if err := Templates["createbootaction"].Execute(w, ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Print("Called create boot action")
}
func BootactionHandler(w http.ResponseWriter, r *http.Request) {
	_, v := conn.GetAllBootActions("bootactions")

	if err := Templates["actions"].Execute(w, v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func PxebootHandler(w http.ResponseWriter, r *http.Request) {
	_, v := conn.GetAllBootActions("pxe")

	if err := Templates["pxeboot"].Execute(w, v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func mkBootEntry(path string, append string) error {
	file, err := os.Create(path)
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

func PXEBOOT(w http.ResponseWriter, r *http.Request) {
	var pxe PXEBOOTTYPE

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &pxe); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	tftpPath := flag.Lookup("tftpPath").Value.String()
	ksURL := flag.Lookup("ksURL").Value.String()
	ksPort := flag.Lookup("port").Value.String()
	filePath := tftpPath + pxe.UUID
	EFIfilePath := tftpROOT + "grub.cfg-" + pxe.UUID
	var kickstart string

	if pxe.IP != "" && pxe.MASK != "" && pxe.NS1 != "" && pxe.NS2 != "" && pxe.GW != "" {
		kickstart = "http://" + ksURL + ":" + ksPort + "/kickstart/" + "?name=" + pxe.KSFile + "&os=" + pxe.OS + "&version=" + pxe.Version + "&fqdn=" + pxe.Hostname + "&ip=" + pxe.IP + "&mask=" + pxe.MASK + "&gw=" + pxe.MASK + "&ns1=" + pxe.NS1 + "&ns2=" + pxe.NS2
	} else {
		kickstart = "http://" + ksURL + ":" + ksPort + "/kickstart/" + "?name=" + pxe.KSFile + "&os=" + pxe.OS + "&version=" + pxe.Version + "&fqdn=" + pxe.Hostname
	}

	err, results := conn.GetBootAction("bootactions", pxe.BootAction)
	if err != nil {
		panic(err)
	}
	err, EFIresults := conn.GetBootAction("efibootactions", pxe.BootAction)
	if err != nil {
		panic(err)
	}

	if results == "" {
		io.WriteString(w, `{"bootaction not found, make sure it exist"}`)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bootAction := results + " " + "ks=" + kickstart

	// Record the pxeboot action in the db
	err = conn.PutBootAction("pxe", pxe.UUID, bootAction)
	checkError(err)
	err = conn.PutBootAction("efipxe", pxe.UUID, EFIresults)
	checkError(err)

	err = mkBootEntry(filePath, bootAction)
	mkBootEntry(EFIfilePath, EFIresults)
	
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"Status": "Failed"}`)
		log.Printf("Error %v \n", err)
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		io.WriteString(w, `{"Status": "success"}`)
		log.Printf("Host %s was pxebooted using profile %v\n", pxe.Hostname, pxe.KSFile)
		return
	}
}
