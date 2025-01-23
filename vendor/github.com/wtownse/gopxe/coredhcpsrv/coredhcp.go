// Copyright 2018-present the CoreDHCP Authors. All rights reserved
// This source code is licensed under the MIT license found in the
// LICENSE file in the root directory of this source tree.

// This is a generated file, edits should be made in the corresponding source file
// And this file regenerated using `coredhcp-generator --from core-plugins.txt`
package coredhcpsrv

import (
	"bytes"
	"errors"
	"fmt"

	"os"
	"time"

	"github.com/coredhcp/coredhcp/config"

	"github.com/coredhcp/coredhcp/server"
	"gopkg.in/yaml.v3"

	"github.com/coredhcp/coredhcp/plugins"
	pl_dns "github.com/coredhcp/coredhcp/plugins/dns"
	//pl_file "github.com/coredhcp/coredhcp/plugins/file"
	pl_file "github.com/wtownse/gopxe/coredhcp/plugins/file"
	pl_leasetime "github.com/coredhcp/coredhcp/plugins/leasetime"
	pl_mtu "github.com/coredhcp/coredhcp/plugins/mtu"
	//pl_nbp "github.com/coredhcp/coredhcp/plugins/nbp"
	pl_nbp "github.com/wtownse/gopxe/coredhcp/plugins/nbp"
	pl_netmask "github.com/coredhcp/coredhcp/plugins/netmask"
	pl_prefix "github.com/coredhcp/coredhcp/plugins/prefix"
	//pl_range "github.com/coredhcp/coredhcp/plugins/range"
	pl_range "github.com/wtownse/gopxe/coredhcp/plugins/range"
	pl_router "github.com/coredhcp/coredhcp/plugins/router"
	pl_searchdomains "github.com/coredhcp/coredhcp/plugins/searchdomains"
	pl_serverid "github.com/coredhcp/coredhcp/plugins/serverid"
	pl_sleep "github.com/coredhcp/coredhcp/plugins/sleep"
	pl_staticroute "github.com/coredhcp/coredhcp/plugins/staticroute"

	"github.com/sirupsen/logrus"
)

var DesiredPlugins = []*plugins.Plugin{
	&pl_dns.Plugin,
	&pl_file.Plugin,
	&pl_leasetime.Plugin,
	&pl_mtu.Plugin,
	&pl_nbp.Plugin,
	&pl_netmask.Plugin,
	&pl_prefix.Plugin,
	&pl_range.Plugin,
	&pl_router.Plugin,
	&pl_searchdomains.Plugin,
	&pl_serverid.Plugin,
	&pl_sleep.Plugin,
	&pl_staticroute.Plugin,
}

func dynamic(log logrus.Entry, ovr map[string]interface{}) {

	//config.Server4.Plugins = append(config.Server4.Plugins, ovr)
	b, err := yaml.Marshal(ovr)
	var bb = *bytes.NewBuffer(b)
	yamlEncoder := yaml.NewEncoder(&bb)
	yamlEncoder.SetIndent(2)

	fmt.Printf("%s", bb.Bytes())
	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}

	path := "/coredhcp/"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	fileName := path + "config.yaml"
	err = os.WriteFile(fileName, bb.Bytes(), 0644)

	if err != nil {
		panic("Unable to write data into the file")
	}

}

func Run(log logrus.Entry, flagConfig string, flagDynamicConfig bool, dyn map[string]interface{}) {

	// register plugins
	for _, plugin := range DesiredPlugins {
		if err := plugins.RegisterPlugin(plugin); err != nil {
			log.Fatalf("Failed to register plugin '%s': %v", plugin.Name, err)
		}
	}
	if flagDynamicConfig {
		dynamic(log, dyn)

	}
	config, err := config.Load(flagConfig)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	// start server
	srv, err := server.Start(config)
	//srv, err := server.Start(b)
	if err != nil {
		log.Fatal(err)
	}
	if err := srv.Wait(); err != nil {
		log.Print(err)
	}
	
	time.Sleep(time.Second)
}
