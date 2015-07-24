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
	"github.com/golang/glog"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/api"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/client"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/config"
	flag "github.com/spf13/pflag"
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
	flagSet.StringVar(&nkm.mConfig.LogDir, "log_dir",
		"/var/log/nuagekubmon/", "Log Directory")
	flagSet.StringVar(&nkm.mConfig.LicenseFile, "license_file",
		"", "VSD License File Path")
}

func (nkm *NuageKubeMonitor) SetLogging(flagSet *flag.FlagSet) {
	if flagSet != nil {
		flagSet.Lookup("log_dir").Value.Set("/var/log/nuagekubemon")
		//flagSet.Lookup("v").Value.Set(10)
	}
}

func (nkm *NuageKubeMonitor) Run() {
	glog.Info("Starting NuageKubeMonitor...")
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
