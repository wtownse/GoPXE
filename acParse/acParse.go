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
	//set bootoption
        webServerHost := flag.Lookup("wsHOST").Value.(flag.Getter).Get().(string)
        webServerPort := flag.Lookup("wsPORT").Value.(flag.Getter).Get().(string)
        agentFilePath := flag.Lookup("acFILE").Value.(flag.Getter).Get().(string)
        bootFilePath := flag.Lookup("bootFILEPath").Value.(flag.Getter).Get().(string)
        webFilePath := flag.Lookup("webFILEPath").Value.(flag.Getter).Get().(string)
        requestURL := fmt.Sprintf("http://localhost:%d/bootaction/myfirstbootaction", serverPort)
        kernelString := fmt.Sprintf(bootFilePath+"agent.x86_64-vmlinuz coreos.live.rootfs_url=http://%s:%s"+webFilePath+"agent.x86_64-rootfs.img ignition.firstboot ignition.platform.id=metal",webServerHost,webServerPort)
        fmt.Printf("webserverhost: " + webServerHost + " webserverport: " + webServerPort + " agent file path: " + agentFilePath + "\n\n")
        data := map[string]interface{}{  //create boot menu
        "default": "coreos",
        "label": "coreos",
        "menu": "coreos",
        "kernel": kernelString,
        "ksdevice": "link",
        "load_ramdisk": "1",
        "initrd": bootFilePath+"/agent.x86_64-initrd.img",
        }
           time.Sleep(100 * time.Millisecond)
           jsonBMData, err := json.Marshal(data)
           if err != nil {
             fmt.Printf("could not marshal json: %s\n", err)
           }else {  // create boot menu entry if it doesn't exist
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
	//check if boot option is set
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
        //check if agent config file exists
        if os.IsNotExist(err){
            fmt.Printf(agentFilePath + " file doesn't exist\n")
        } else {
            //create host boot configuration files
            fmt.Printf("found agent config file\n")
            agentConfigData, err := ioutil.ReadFile(agentFilePath) //Read agent-config file
            if err != nil {
                fmt.Printf("can't read file")
            }

            var aconfig map[string]interface{}  //define map to store agent-config file data
            err2 := yaml.Unmarshal(agentConfigData, &aconfig) //parse yaml into agentConfigData
            if err2 != nil {
                fmt.Printf("%s",err2)
            }

            hosts := aconfig["hosts"].([]interface{}) //populate list of hosts from agent-config
            for _, v := range hosts {  //loop hosts
                hs := v.(map[string]interface{})  //get list of host interfaces
                ifaces := hs["interfaces"].([]interface{})
                    iface := ifaces[0].(map[string]interface{})  //get first host interface
                    uuid := "01-" + strings.ReplaceAll(string(iface["macAddress"].(interface{}).(string)),":","-")
                    data := map[string]interface{}{  //create host config map
                    "bootaction": "myfirstbootaction",
                    "ksfile": "default",
                    "os": "coreos",
                    "version": "4.13",
                    "hostname": "test-myvm.local",
                    "uuid": uuid,
                    }

           jsonData, err := json.Marshal(data) //convert host config map to json

           if err != nil {
               fmt.Printf("could not marshal json: %s\n", err)
           }else {
           time.Sleep(100 * time.Millisecond)
           jsonBody := []byte(jsonData)  //convert json data to byte array
           bodyReader := bytes.NewReader(jsonBody)  
           requestURL2 := fmt.Sprintf("http://localhost:%d/pxeboot", serverPort) // pxeboot config url
           req, err := http.NewRequest(http.MethodPost, requestURL2, bodyReader) // create html request

              if err != nil {
                 fmt.Printf("client: could not create request: %s\n", err)
                 os.Exit(1)
                 }

                 req.Header.Set("Content-Type", "application/json") //set html request header
                 client := http.Client{  //create html client
                 Timeout: 30 * time.Second,
                 }

                 res, err := client.Do(req)  //perform post operation
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

