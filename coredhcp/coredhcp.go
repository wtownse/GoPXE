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
	"io/ioutil"
	"os"
	"time"

	"github.com/coredhcp/coredhcp/config"
	"github.com/coredhcp/coredhcp/logger"
	"github.com/coredhcp/coredhcp/server"
	"gopkg.in/yaml.v3"

	"github.com/coredhcp/coredhcp/plugins"
	pl_dns "github.com/coredhcp/coredhcp/plugins/dns"
	pl_file "github.com/coredhcp/coredhcp/plugins/file"
	pl_leasetime "github.com/coredhcp/coredhcp/plugins/leasetime"
	pl_mtu "github.com/coredhcp/coredhcp/plugins/mtu"
	pl_nbp "github.com/coredhcp/coredhcp/plugins/nbp"
	pl_netmask "github.com/coredhcp/coredhcp/plugins/netmask"
	pl_prefix "github.com/coredhcp/coredhcp/plugins/prefix"
	pl_range "github.com/coredhcp/coredhcp/plugins/range"
	pl_router "github.com/coredhcp/coredhcp/plugins/router"
	pl_searchdomains "github.com/coredhcp/coredhcp/plugins/searchdomains"
	pl_serverid "github.com/coredhcp/coredhcp/plugins/serverid"
	pl_sleep "github.com/coredhcp/coredhcp/plugins/sleep"
	pl_staticroute "github.com/coredhcp/coredhcp/plugins/staticroute"

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var (
	flagLogFile       = flag.StringP("logfile", "l", "", "Name of the log file to append to. Default: stdout/stderr only")
	flagLogNoStdout   = flag.BoolP("nostdout", "N", false, "Disable logging to stdout/stderr")
	flagLogLevel      = flag.StringP("loglevel", "L", "info", fmt.Sprintf("Log level. One of %v", getLogLevels()))
	flagConfig        = flag.StringP("conf", "c", "", "Use this configuration file instead of the default location")
	flagDynamicConfig = flag.BoolP("usedynamicconfig", "U", false, "use configuration file.  Overwrites file with command options if false.  Defaults to false")
	flagPlugins       = flag.BoolP("plugins", "P", false, "list plugins")
	dns               = flag.StringP("dns", "D", "8.8.8.8 8.8.4.4", "Override default dns configuration.  Servers are space separated.")
	nbpFile           = flag.StringP("nbp", "B", "tftp://172.16.130.107/pxelinux.0", "Override default tftp boot file url pxelinux.0.")
	dhcpRange         = flag.StringP("range", "R", "172.16.130.100 172.16.130.120", "Set the dhcp server range.")
	dhcpNetmask       = flag.StringP("netmask", "M", "255.255.255.0", "Set the dhcp server mask.")
	dhcpRouter        = flag.StringP("router", "G", "172.16.130.1", "Set the dhcp server default gateway.")
	dhcpServerID      = flag.StringP("sid", "S", "172.16.130.107", "Set the server ID.")
	dhcpLeaseDuration = flag.StringP("lease", "d", "300s", "Set the lease duration")
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

var desiredPlugins = []*plugins.Plugin{
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

func dynamic(log logrus.Entry) {
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

	if *flagPlugins {
		for _, p := range desiredPlugins {
			fmt.Println(p.Name)
		}
		os.Exit(0)
	}

}

func Run() {
	flag.Parse()
	log := logger.GetLogger("main")
	if *flagDynamicConfig {
		dynamic(*log)

	}

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
	config, err := config.Load(*flagConfig)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// register plugins
	for _, plugin := range desiredPlugins {
		if err := plugins.RegisterPlugin(plugin); err != nil {
			log.Fatalf("Failed to register plugin '%s': %v", plugin.Name, err)
		}
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
