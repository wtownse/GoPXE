package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/coredhcp/coredhcp/logger"
	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"

	//"github.com/wtownse/gopxe/conf"

	"github.com/wtownse/gopxe/coredhcpsrv"
	"github.com/wtownse/gopxe/routes"
	"github.com/wtownse/gopxe/tftpsrv"
)

const (
	tftpROOT = "/var/lib/tftpboot/"
)

var (
	flagLogFile       = flag.StringP("logfile", "l", "", "Name of the log file to append to. Default: stdout/stderr only")
	flagLogNoStdout   = flag.BoolP("nostdout", "N", false, "Disable logging to stdout/stderr")
	flagLogLevel      = flag.StringP("loglevel", "L", "info", fmt.Sprintf("Log level. One of %v", getLogLevels()))
	flagConfig        = flag.StringP("conf", "c", "", "Use this configuration file instead of the default location")
	flagDynamicConfig = flag.BoolP("usedynamicconfig", "U", false, "use configuration file.  Overwrites file with command options if false.  Defaults to false")
	flagPlugins       = flag.BoolP("plugins", "P", false, "list plugins")
	dns               = flag.StringP("dns", "D", "8.8.8.8 8.8.4.4", "Override default dns configuration.  Servers are space separated.")
	nbpFile           = flag.StringP("nbp", "B", "tftp://172.16.130.107/pxelinux.0 tftp://172.16.130.107/bootx64.efi", "Override default tftp boot file url pxelinux.0.")
	dhcpRange         = flag.StringP("range", "R", "172.16.130.100 172.16.130.120", "Set the dhcp server range.")
	dhcpNetmask       = flag.StringP("netmask", "M", "255.255.255.0", "Set the dhcp server mask.")
	dhcpRouter        = flag.StringP("router", "G", "172.16.130.1", "Set the dhcp server default gateway.")
	dhcpServerID      = flag.StringP("sid", "S", "172.16.130.107", "Set the server ID.")
	dhcpLeaseDuration = flag.StringP("lease", "d", "300s", "Set the lease duration")
	port              = flag.StringP("port", "p", "9090", "Set gopxe port")
	tftpPath          = flag.StringP("tftpPath", "T", tftpROOT+"pxelinux.cfg/", "tftp conf path e.g /var/lib/tftpboot/pxelinux.cfg/")
	ksURL             = flag.StringP("ksURL", "K", "localhost", "kickstart url")
	wsHOST            = flag.StringP("wsHOST", "s", "localhost", "external webserver host ip")
	wsPORT            = flag.StringP("wsPORT", "r", "80", "external webserver port")
	bucket            = flag.StringP("bucket", "", "bootactions", "db bucket")
	efibucket         = flag.StringP("efibucket", "", "efibootactions", "db bucket")
	dbName            = flag.StringP("dbName", "", "gopxe.db", "database file name")
	acFILE            = flag.StringP("acFILE", "a", tftpROOT+"coreos/configs/agent-config.yaml", "agent config file path relative to tftpROOT /var/lib/tftpboot/")
	bootFILEPath      = flag.StringP("bootFILEPath", "b", "/", " path to pxe boot files")
	webFILEPath       = flag.StringP("webFILEPath", "w", "/", "path to web root folder")
	staticLease       = flag.StringP("staticLease", "m", "leases.txt", "static lease file")
)

var logLevels = map[string]func(*logrus.Logger){
	"none":    func(l *logrus.Logger) { l.SetOutput(ioutil.Discard) },
	"debug":   func(l *logrus.Logger) { l.SetLevel(logrus.DebugLevel) },
	"info":    func(l *logrus.Logger) { l.SetLevel(logrus.InfoLevel) },
	"warning": func(l *logrus.Logger) { l.SetLevel(logrus.WarnLevel) },
	"error":   func(l *logrus.Logger) { l.SetLevel(logrus.ErrorLevel) },
	"fatal":   func(l *logrus.Logger) { l.SetLevel(logrus.FatalLevel) },
}

func getLogLevels() []string {
	var levels []string
	for k := range logLevels {
		levels = append(levels, k)
	}
	return levels
}

func getDhcpFlags() map[string]interface{} {
	flag.Parse()
	DNS, NBP, RANGE, NETMASK, ROUTER, SID, LEASE := *dns, *nbpFile, *dhcpRange, *dhcpNetmask, *dhcpRouter, *dhcpServerID, *dhcpLeaseDuration
	ovr := map[string]interface{}{
		//		"server6": map[string]interface{}{
		//			"plugins": []map[string]string{
		//				{
		//					"file": "leases.txt",
		//				},
		//			},
		//		},
		"server4": map[string]interface{}{
			"plugins": []map[string]string{
				{
					"lease_time": LEASE,
				},
				{
					"server_id": SID,
				},
				{
					"dns": DNS,
				},
				{
					"router": ROUTER,
				},
				{
					"netmask": NETMASK,
				},
				{
					"range": "leases.txt " + RANGE + " 60s",
				},
				{
					"nbp": NBP,
				},
			},
		},
	}
	return ovr
}

// This is the main package
// Output is webserver om port
func main() {

	log := logger.GetLogger("main")

	fn, ok := logLevels[*flagLogLevel]
	if !ok {
		log.Fatalf("Invalid log level '%s'. Valid log levels are %v", *flagLogLevel, getLogLevels())
	}
	fn(log.Logger)
	log.Infof("Setting log level to '%s'", *flagLogLevel)
	if *flagLogFile != "" {
		log.Infof("Logging to file %s", *flagLogFile)
		logger.WithFile(log, *flagLogFile)
	}
	if *flagLogNoStdout {
		log.Infof("Disabling logging to stdout/stderr")
		logger.WithNoStdOutErr(log)
	}

	//conf.Setup()
	port := flag.Lookup("port").Value
	wg := new(sync.WaitGroup)
	wg.Add(4)
	go func() {
		tftpsrv.Run("69")
		log.Printf("tftp listening on port: 69")
	}()
	go func() {
		routes := routes.New()
		log.Printf("Serving on port: %s", port)
		if err := http.ListenAndServe(":"+port.String(), routes); err != nil {
			log.Fatal("ListenAndServe: ", err.Error())
		}
	}()
	go func() {
		requestURL := fmt.Sprintf("http://localhost:%s/acparse", port)
		time.Sleep(3 * time.Second)
		res, err := http.Get(requestURL)
		if err != nil {
			log.Printf("error making http request: %s\n", err)
		}
		log.Printf("client: got response!\n")
		log.Printf("client: status code: %d\n", res.StatusCode)
	}()
	go func() {
		coredhcpsrv.Run(*log, *flagConfig, *flagDynamicConfig, getDhcpFlags())
	}()

	wg.Wait()

}
