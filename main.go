package main

import (
	"flag"
	"log"
        "sync"
	"net/http"

	"github.com/wtownse/gopxe/conf"
	"github.com/wtownse/gopxe/routes"
        "github.com/wtownse/gopxe/tftpsrv"
)

// This is the main package
// Output is webserver om port
func main() {
        conf.Setup()
        port := flag.Lookup("port").Value.(flag.Getter).Get().(string)
        wg := new(sync.WaitGroup)
        wg.Add(2)
        go func(){
        tftpsrv.Run("69")
        log.Printf("tftp listening on port: 69")
        }()
        go func(){
	routes := routes.New()
	log.Printf("Serving on port: %s", port)
	if err := http.ListenAndServe(":"+port, routes); err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
        }()
        wg.Wait()

}
