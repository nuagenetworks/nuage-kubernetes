/*
###########################################################################
#
#   Filename:           client_test.go
#
#   Author:             Ryan Fredette
#   Created:            August 3, 2015
#
#   Description:        Common functions and test setup for the client
#                       package
#
###########################################################################
#
#              Copyright (c) 2015 Nuage Networks
#
###########################################################################

*/

package client

import (
	"flag"
	"fmt"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/config"
	"os"
	"os/exec"
	"testing"
)

var kubemonConfig *config.NuageKubeMonConfig
var vsdClient *NuageVsdClient
var isTargetSystem bool

func TestMain(m *testing.M) {
	kubemonConfig = &config.NuageKubeMonConfig{}
	addArgs(kubemonConfig, flag.CommandLine)
	flag.CommandLine.Parse(os.Args[1:])
	if testing.Short() {
		isTargetSystem = false
	} else {
		vsdClient = new(NuageVsdClient)
		vsdClient.namespaces = make(map[string]NamespaceData)
		vsdClient.version = kubemonConfig.NuageVspVersion
		vsdClient.url = kubemonConfig.NuageVsdApiUrl + "/nuage/api/" + vsdClient.version + "/"
		vsdClient.CreateSession()
		if err := vsdClient.GetAuthorizationToken(); err != nil {
			fmt.Printf("Error getting VSD auth token: %s\n", err)
			vsdClient = nil
		}
		// Check if we have `oc`.  If it's not present, this isn't the
		// target system for nuagekubemon, so some tests cannot be run.
		_, err := exec.Command("oc", "whoami").CombinedOutput()
		if err != nil {
			isTargetSystem = false
		} else {
			isTargetSystem = true
		}
	}
	os.Exit(m.Run())
}

func addArgs(myConfig *config.NuageKubeMonConfig, flagSet *flag.FlagSet) {
	flagSet.StringVar(&myConfig.KubeConfigFile, "kubeconfig",
		"", "kubeconfig File for Openshift User")
	flagSet.StringVar(&myConfig.MasterConfigFile, "masterconfig",
		"", "Path to master-config.yaml for the cluster master")
	flagSet.StringVar(&myConfig.NuageVsdApiUrl, "nuagevsdurl",
		"", "Nuage VSD URL")
	flagSet.StringVar(&myConfig.NuageVspVersion, "nuagevspversion",
		"", "Nuage VSP Version")
	flagSet.StringVar(&myConfig.LicenseFile, "license_file",
		"", "VSD License File Path")
}
