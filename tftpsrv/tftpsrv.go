package tftpsrv

import (
	"fmt"
	"io"
	"os"

	"github.com/coredhcp/coredhcp/logger"
	//"pack.ag/tftp"
	"github.com/pin/tftp/v3"
)

var log = logger.GetLogger("tftp")

func Run(port string) {
	s := tftp.NewServer(readHandler, writeHandler)
	s.EnableSinglePort()
	s.SetBlockSize(512)
	// readHandler := tftp.ReadHandlerFunc(proxyTFTP)

	// s.ReadHandler(readHandler)

	s.ListenAndServe(fmt.Sprintf(":%s", port))
	select {}

}

// func proxyTFTP(w tftp.ReadRequest) {
// 	log.Printf("[%s] GET %s\n", w.Addr().IP.String(), w.Name())

// 	file, err := os.Open("/var/lib/tftpboot/" + w.Name()) // For read access.
// 	if err != nil {
// 		log.Println(err)
// 		w.WriteError(tftp.ErrCodeFileNotFound, err.Error())
// 		return
// 	}
// 	defer file.Close()

// 	if _, err := io.Copy(w, file); err != nil {
// 		log.Println(err)
// 	}
// }

// readHandler is called when client starts file download from server
func readHandler(filename string, rf io.ReaderFrom) error {
	file, err := os.Open("/var/lib/tftpboot/" + filename)
	if err != nil {
		log.Println(err)

		return err
	}
	n, err := rf.ReadFrom(file)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Printf("%d bytes sent\n", n)
	log.Printf("sent %d bytes\n", n)
	return nil
}

// writeHandler is called when client starts file upload to server
func writeHandler(filename string, wt io.WriterTo) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	n, err := wt.WriteTo(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return err
	}
	fmt.Printf("%d bytes received\n", n)
	log.Printf("received %d bytes\n", n)
	return nil
}
func test() {}
