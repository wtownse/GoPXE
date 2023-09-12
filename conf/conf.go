package conf

import (
        "flag"
)

var (
        tftpPath   string
        ksURL      string
        wsHOST     string
        wsPORT     string
        port       string
        bucket     string
        dbName     string
        acfile     string
)

const (
        tftpROOT = "/var/lib/tftpboot/"
)

func Setup() {
        // Define flags 
        flag.StringVar(&tftpPath, "tftpPath", tftpROOT+"pxelinux.cfg/", "tftp conf path e.g /var/lib/tftpboot/pxelinux.cfg/")
        flag.StringVar(&port, "port", "9090", "tcp port")
        flag.StringVar(&bucket, "bucket", "bootactions", "db bucket")
        flag.StringVar(&dbName, "dbName", "gopxe.db", "database file name")
        flag.StringVar(&ksURL, "ksURL", "localhost", "kickstart url")
        flag.StringVar(&wsHOST, "wsHOST", "localhost", "external webserver host ip")
        flag.StringVar(&wsPORT, "wsPORT", "80", "external webserver port")
        flag.StringVar(&acfile, "acFILE", tftpROOT+"configs/agent-installer.yaml", "agent config file path relative to tftpROOT /var/lib/tftpboot/")
        // Parsing flags
        flag.Parse()
}

