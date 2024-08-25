package conf

import (
	"flag"
)

var (
	tftpPath     string
	ksURL        string
	wsHOST       string
	wsPORT       string
	port         string
	bucket       string
	efiBucket    string
	dbName       string
	acFILE       string
	bootFILEPath string
	webFILEPath  string
	staticLease  string
)

const (
	tftpROOT = "/var/lib/tftpboot/"
)

func Setup() {
	// Define flags
	flag.StringVar(&tftpPath, "tftpPath", tftpROOT+"pxelinux.cfg/", "tftp conf path e.g /var/lib/tftpboot/pxelinux.cfg/")
	flag.StringVar(&port, "port", "9090", "tcp port")
	flag.StringVar(&bucket, "bucket", "bootactions", "db bucket")
	flag.StringVar(&efiBucket, "efibucket", "EFIbootactions", "EFI db bucket")
	flag.StringVar(&dbName, "dbName", "gopxe.db", "database file name")
	flag.StringVar(&ksURL, "ksURL", "localhost", "kickstart url")
	flag.StringVar(&wsHOST, "wsHOST", "localhost", "external webserver host ip")
	flag.StringVar(&wsPORT, "wsPORT", "80", "external webserver port")
	flag.StringVar(&acFILE, "acFILE", tftpROOT+"coreos/configs/agent-config.yaml", "agent config file path relative to tftpROOT /var/lib/tftpboot/")
	flag.StringVar(&bootFILEPath, "bootFILEPath", "/", " path to pxe boot files")
	flag.StringVar(&webFILEPath, "webFILEPath", "/", "path to web root folder")
	flag.StringVar(&staticLease, "staticLease", "leases.txt", "static lease file")
	// Parsing flags
	flag.Parse()
}
