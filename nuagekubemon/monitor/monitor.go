/*
###########################################################################
#
#   Filename:           nuagekubemon.go
#
#   Author:             Aniket Bhat
#   Created:            July 20, 2015
#
#   Description:        NuageKubeMon monitor interface
#
###########################################################################
#
#              Copyright (c) 2015 Nuage Networks
#
###########################################################################

*/

/* package monitor implements a kubernetes/openshift monitor for integration with Nuage VSP */

package monitor

import (
	"flag"
	"github.com/golang/glog"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/api"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/client"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/config"
	"io/ioutil"
)

type NuageKubeMonitor struct {
	mConfig    config.NuageKubeMonConfig
	mVsdClient *client.NuageVsdClient
	mOsClient  *client.NuageOsClient
	//mOsNodeClient nuageosnodeclient.NuageOsNodeClient
}

func NewNuageKubeMonitor() *NuageKubeMonitor {
	nkm := new(NuageKubeMonitor)
	return nkm
}

func (nkm *NuageKubeMonitor) ParseArgs(flagSet *flag.FlagSet) {
	flagSet.StringVar(&nkm.mConfig.OsClusterAdmin, "osusername",
		"system:admin", "User name of the cluster administrator")
	flagSet.StringVar(&nkm.mConfig.KubeConfigFile, "kubeconfig",
		"", "kubeconfig File for Openshift User")
	flagSet.StringVar(&nkm.mConfig.OsMasterConfigFile, "osmasterconfig",
		"", "Path to master-config.yaml for the cluster master")
	flagSet.StringVar(&nkm.mConfig.NuageVsdApiUrl, "nuagevsdurl",
		"", "Nuage VSD URL")
	flagSet.StringVar(&nkm.mConfig.NuageVspVersion, "nuagevspversion",
		"", "Nuage VSP Version")
	flagSet.StringVar(&nkm.mConfig.LicenseFile, "license_file",
		"", "VSD License File Path")
	flagSet.StringVar(&nkm.mConfig.ConfigFile, "config",
		"", "Configuration file for nuagekubemon.  If this argument is specified, all other commandline arguments will be ignored.")
	// Set the values for log_dir and logtostderr.  Because this happens before
	// flag.Parse(), cli arguments will override these.  Also set the DefValue
	// parameter so -help shows the new defaults.
	log_dir := flagSet.Lookup("log_dir")
	log_dir.Value.Set("/var/log/nuagekubemon")
	log_dir.DefValue = "/var/log/nuagekubemon"
	logtostderr := flagSet.Lookup("logtostderr")
	logtostderr.Value.Set("false")
	logtostderr.DefValue = "false"
}

func (nkm *NuageKubeMonitor) LoadConfig() error {
	if nkm.mConfig.ConfigFile == "" {
		// If there was no config file specified, don't try to read the nothing.
		return nil
	}
	configData, err := ioutil.ReadFile(nkm.mConfig.ConfigFile)
	if err != nil {
		return err
	}
	return nkm.mConfig.Parse(configData)
}

func (nkm *NuageKubeMonitor) Run() {
	glog.Info("Starting NuageKubeMonitor...")
	// Read the config file if it was specified.  If there was an error reading
	// it, don't continue.
	if err := nkm.LoadConfig(); err != nil {
		glog.Fatalf("Error reading config file %s! Error: %v\n",
			nkm.mConfig.ConfigFile, err)
	}
	nkm.mOsClient = client.NewNuageOsClient(&(nkm.mConfig))
	nkm.mVsdClient = client.NewNuageVsdClient(&(nkm.mConfig))
	//nkm.mOsNodeClient = client.NuageOsNodeClient(nkm.mConfig)
	stop := make(chan bool)
	nsEventChannel := make(chan *api.NamespaceEvent)
	go nkm.mOsClient.Run(nsEventChannel, stop)
	go nkm.mVsdClient.Run(nsEventChannel, stop)
	//go nkm.mOsNodeClient.Run()
	select {}
}
