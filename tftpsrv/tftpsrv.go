package tftpsrv

import (
        "pack.ag/tftp"
        "log"
        "os"
        "io"
        "fmt"
)

func Run(port string){
        s, err := tftp.NewServer(fmt.Sprintf(":%s",port), tftp.ServerSinglePort(true), tftp.ServerNet("udp4"))
        if err != nil {
                panic(err)
        }
        readHandler := tftp.ReadHandlerFunc(proxyTFTP)
        s.ReadHandler(readHandler)
        s. ListenAndServe()
        select{}

}

func proxyTFTP(w tftp.ReadRequest) {
        log.Printf("[%s] GET %s\n", w.Addr().IP.String(), w.Name() )
        file, err := os.Open("/var/lib/tftpboot/" + w.Name()) // For read access.
        if err != nil {
                log.Println(err)
                w.WriteError(tftp.ErrCodeFileNotFound, err.Error())
                return
        }
        defer file.Close()

        if _, err := io.Copy(w, file); err != nil {
                log.Println(err)
        }
}

func test(){}
