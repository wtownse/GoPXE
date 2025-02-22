package routes

import (
	//	"flag"
	"net/http"

	flag "github.com/spf13/pflag"
	ac "github.com/wtownse/gopxe/acParse"
	h "github.com/wtownse/gopxe/handlers"
	static "github.com/wtownse/gopxe/staticdhcp"

	//External dependencies
	"github.com/gorilla/mux"
)

func New() http.Handler {

	const localrepo string = "/opt/localrepo"
	tftpPath := flag.Lookup("tftpPath").Value.String()

	router := mux.NewRouter()
	router = mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/localrepo").Handler(http.StripPrefix("/localrepo/", http.FileServer(http.Dir(localrepo))))
	router.PathPrefix("/pxelinux").Handler(http.StripPrefix("/pxelinux/", http.FileServer(http.Dir(tftpPath))))
	router.HandleFunc("/", h.Index)
	router.HandleFunc("/viewbootaction", h.BootactionHandler)
	router.HandleFunc("/viewpxeboot", h.PxebootHandler)
	router.HandleFunc("/health", h.StatusHandler)
	router.HandleFunc("/bootaction/{key}", h.GetBA).Methods("GET")
	router.HandleFunc("/bootaction/{key}", h.PutBA).Methods("POST")
	router.HandleFunc("/bootaction", h.GetAllBA).Methods("GET")
	router.HandleFunc("/kickstart/", h.KsGenerate)
	router.HandleFunc("/pxeboot", h.PXEBOOT).Methods("POST")
	router.HandleFunc("/pxeboot2", h.PXEBOOT2).Methods("POST")
	router.HandleFunc("/createbootaction", h.CreateBootAction).Methods("GET")
	router.HandleFunc("/createbootrecord", h.CreatePxeBootRecord).Methods("GET")
	router.HandleFunc("/createbootaction/{key}", h.CreateBA).Methods("POST")
	router.HandleFunc("/acparse", ac.Create)
	router.HandleFunc("/staticdhcp", static.StaticDHCP)
	router.HandleFunc("/getdhcpentry/{mac}", static.GetDHCPEntry)
	router.HandleFunc("/getdhcpentries", static.GetDHCPEntries).Methods("GET")
	router.HandleFunc("/getdhcpentriesjson", static.GetDHCPEntriesJSON).Methods("GET")
	h.LoadTemplates()
	return router
}
