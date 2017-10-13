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
	"errors"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/client"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/config"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type NuageKubeMonitor struct {
	mConfig     config.NuageKubeMonConfig
	mVsdClient  *client.NuageVsdClient
	mOsClient   *client.NuageClusterClient
	metcdClient *client.NuageEtcdClient
	//mOsNodeClient nuageosnodeclient.NuageOsNodeClient
}

func NewNuageKubeMonitor() *NuageKubeMonitor {
	nkm := new(NuageKubeMonitor)
	return nkm
}

func (nkm *NuageKubeMonitor) ParseArgs(flagSet *flag.FlagSet) {
	programName := path.Base(os.Args[0])
	flagSet.StringVar(&nkm.mConfig.KubeConfigFile, "kubeconfig",
		"", "kubeconfig File for Openshift User")
	flagSet.StringVar(&nkm.mConfig.MasterConfigFile, "masterconfig",
		"", "Path to master-config.yaml for the cluster master")
	flagSet.StringVar(&nkm.mConfig.NuageVsdApiUrl, "nuagevsdurl",
		"", "Nuage VSD URL")
	flagSet.StringVar(&nkm.mConfig.NuageVspVersion, "nuagevspversion",
		"", "Nuage VSP Version")
	flagSet.StringVar(&nkm.mConfig.LicenseFile, "license_file",
		"", "VSD License File Path")
	flagSet.StringVar(&nkm.mConfig.ConfigFile, "config",
		"", "Configuration file for nuagekubemon.  If this argument is specified, all other commandline arguments will be ignored.")
	flagSet.StringVar(&nkm.mConfig.EnterpriseName, "enterprise",
		config.DefaultEnterprise(), "Enterprise in which the containers will reside")
	flagSet.StringVar(&nkm.mConfig.DomainName, "domain",
		config.DefaultDomain(), "Domain in which the containers will reside")
	// Set the values for log_dir and logtostderr.  Because this happens before
	// flag.Parse(), cli arguments will override these.  Also set the DefValue
	// parameter so -help shows the new defaults.
	log_dir := flagSet.Lookup("log_dir")
	log_dir.Value.Set(fmt.Sprintf("/var/log/%s", programName))
	log_dir.DefValue = fmt.Sprintf("/var/log/%s", programName)
	logtostderr := flagSet.Lookup("logtostderr")
	logtostderr.Value.Set("false")
	logtostderr.DefValue = "false"
	stderrlogthreshold := flagSet.Lookup("stderrthreshold")
	stderrlogthreshold.Value.Set("4")
	stderrlogthreshold.DefValue = "4"
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func (nkm *NuageKubeMonitor) LoadConfig() error {
	if nkm.mConfig.ConfigFile != "" {
		// If there was a config file specified, overwrite the cli arguments
		configData, err := ioutil.ReadFile(nkm.mConfig.ConfigFile)
		if err != nil {
			return err
		}
		if err := nkm.mConfig.Parse(configData); err != nil {
			return err
		}
	}
	if nkm.mConfig.MasterConfigFile == "" {
		return errors.New("No master config file specified")
	}
	osMasterData, err := ioutil.ReadFile(nkm.mConfig.MasterConfigFile)
	if err != nil {
		return err
	}
	return nkm.mConfig.MasterConfig.Parse(osMasterData)
}

func (nkm *NuageKubeMonitor) Run() {
	var err error
	programName := path.Base(os.Args[0])
	glog.Infof("Starting %s...", programName)
	defer glog.Flush() //Flush logs when the monitor exits
	// Read the config file if it was specified.  If there was an error reading
	// it, don't continue.
	if err := nkm.LoadConfig(); err != nil {
		glog.Fatalf("Error reading config file %s! Error: %v\n",
			nkm.mConfig.ConfigFile, err)
	}
	if nkm.mConfig.KubeConfigFile == "" {
		glog.Error(fmt.Sprintf("No valid kubeconfig file specified...%s cannot continue.", programName))
		glog.Error(fmt.Sprintf("Please restart %s after specifying a valid kubeconfig path either in the config file or as a command line parameter",
			programName))
		return
	}
	if nkm.metcdClient, err = client.NewNuageEtcdClient(&nkm.mConfig); err != nil {
		glog.Errorf("Creating etcd client failed with error: %v", err)
		return
	}

	etcdChannel := make(chan *api.EtcdEvent)

	nkm.mOsClient = client.NewNuageOsClient(&(nkm.mConfig))
	nkm.mVsdClient = client.NewNuageVsdClient(&(nkm.mConfig), nkm.mOsClient.GetClusterClientCallBacks(), etcdChannel)
	stop := make(chan bool)
	nsEventChannel := make(chan *api.NamespaceEvent)
	serviceEventChannel := make(chan *api.ServiceEvent)
	policyEventChannel := make(chan *api.NetworkPolicyEvent)
	go nkm.metcdClient.Run(etcdChannel)
	go nkm.mVsdClient.Run(nsEventChannel, serviceEventChannel, policyEventChannel, stop)
	nkm.mOsClient.GetExistingEvents(nsEventChannel, serviceEventChannel, policyEventChannel)
	go nkm.mOsClient.RunNamespaceWatcher(nsEventChannel, stop)
	go nkm.mOsClient.RunServiceWatcher(serviceEventChannel, stop)
	go nkm.mOsClient.RunNetworkPolicyWatcher(policyEventChannel, stop)
	select {}
}
