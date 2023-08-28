package acParse

import (
        "fmt"
        "net/http"
        "os"
        "time"
        "flag"
        "io/ioutil"
        "gopkg.in/yaml.v3"
        "strings"
        "bytes"
        "encoding/json"
)

const (
        serverPort = 9090
        agentFilePath = "/opt/storage/agent-config.yaml"
)


type myhosts struct {
  Hosts []host  `yaml: "hosts"`
}

type host struct {
  Hostname string `yaml: "hostname"`
  Interfaces []interfaces `yaml: "interfaces"`
}

type interfaces struct {
  Name string `yaml: "name"`
  Macaddress string `yaml:"macAddress"`
}

func Create(http.ResponseWriter, *http.Request) {
        webServerHost := flag.Lookup("wsHOST").Value.(flag.Getter).Get().(string)
        webServerPort := flag.Lookup("wsPORT").Value.(flag.Getter).Get().(string)
        requestURL := fmt.Sprintf("http://localhost:%d/bootaction/myfirstbootaction", serverPort)
        kernelString := fmt.Sprintf("coreos/coreos/agent.x86_64-vmlinuz coreos.live.rootfs_url=http://%s:%s/coreos/agent.x86_64-rootfs.img ignition.firstboot ignition.platform.id=metal",webServerHost,webServerPort)
        data := map[string]interface{}{
        "default": "coreos",
        "label": "coreos",
        "menu": "coreos",
        "kernel": kernelString,
        "ksdevice": "link",
        "load_ramdisk": "1",
        "initrd": "coreos/coreos/agent.x86_64-initrd.img",
        }
           time.Sleep(100 * time.Millisecond)
           jsonBMData, err := json.Marshal(data)
           if err != nil {
     fmt.Printf("could not marshal json: %s\n", err)
           }else {

           BootmenuJsonBody := []byte(jsonBMData)
           bodyReader := bytes.NewReader(BootmenuJsonBody)
           BMreq, err := http.NewRequest(http.MethodPost, requestURL, bodyReader                                                               )

           if err != nil {
                 fmt.Printf("client: could not create request: %s\n", err)
                 os.Exit(1)
          }
          BMreq.Header.Set("Content-Type", "application/json")

         client := http.Client{
          Timeout: 30 * time.Second,
          }

          res, err := client.Do(BMreq)
          if err != nil {
          fmt.Printf("client: error making http request: %s\n", err)
          os.Exit(1)
          }

                fmt.Printf("%s",jsonBMData)
                fmt.Printf("%s",res)
         }

        fmt.Printf("%s",data)
        time.Sleep(100 * time.Millisecond)
        req, err := http.NewRequest(http.MethodGet, requestURL, nil)
        if err != nil {
                fmt.Printf("error making http request: %s\n", err)
                os.Exit(1)
        }
        res, err := http.DefaultClient.Do(req)
        resBody, err := ioutil.ReadAll(res.Body)
        fmt.Printf("client: got response!\n")
        fmt.Printf("client: status code: %d\n", res.StatusCode)
        bootstring := string(resBody)
        if bootstring != "" {
                fmt.Printf("bootaction is set\n")
        } else {
                fmt.Printf("no bootaction is set\n")
        }

        _, err = os.Stat(agentFilePath)

        if os.IsNotExist(err){
        fmt.Printf("file doesn't exist\n")
        } else {
        fmt.Printf("found agent config file\n")
        agentConfigData, err := ioutil.ReadFile(agentFilePath)
            if err != nil {
            fmt.Printf("can't read file")
            }
             //Acg := myhosts{}
             var aconfig map[string]interface{}
             err2 := yaml.Unmarshal(agentConfigData, &aconfig)

              if err2 != nil {

               fmt.Printf("%s",err2)
               }
            hosts := aconfig["hosts"].([]interface{})
            for _, v := range hosts {
            //fmt.Printf("key:",k,"value:",v)
            //fmt.Printf("\n")
            hs := v.(map[string]interface{})
            ifaces := hs["interfaces"].([]interface{})
            iface := ifaces[0].(map[string]interface{})
            //fmt.Printf("%s\n",strings.ReplaceAll(string(iface["macAddress"].(i                                                               nterface{}).(string)),":","-"))
            //fmt.Printf("%s\n",hs["hostname"])
            //fmt.Printf("\n\n")
                  uuid := "01-" + strings.ReplaceAll(string(iface["macAddress"].                                                               (interface{}).(string)),":","-")
                  data := map[string]interface{}{
                  "bootaction": "myfirstbootaction",
                  "ksfile": "default",
                  "os": "coreos",
                  "version": "4.13",
                  "hostname": "test-myvm.local",
                  "uuid": uuid,
                   }
           jsonData, err := json.Marshal(data)
           if err != nil {
     fmt.Printf("could not marshal json: %s\n", err)
           }else {
           time.Sleep(100 * time.Millisecond)

           jsonBody := []byte(jsonData)
           bodyReader := bytes.NewReader(jsonBody)
           requestURL2 := fmt.Sprintf("http://localhost:%d/pxeboot", serverPort)
           req, err := http.NewRequest(http.MethodPost, requestURL2, bodyReader)

           if err != nil {
                 fmt.Printf("client: could not create request: %s\n", err)
                 os.Exit(1)
          }
          req.Header.Set("Content-Type", "application/json")

         client := http.Client{
          Timeout: 30 * time.Second,
          }

          res, err := client.Do(req)
          if err != nil {
          fmt.Printf("client: error making http request: %s\n", err)
          os.Exit(1)
          }

                fmt.Printf("%s",jsonData)
                fmt.Printf("%s",res)
           }
         }
}
}

