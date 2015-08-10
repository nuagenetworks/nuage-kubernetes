package client

import (
	"flag"
	"fmt"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/config"
	"os"
	"testing"
)

var kubemonConfig *config.NuageKubeMonConfig
var vsdClient *NuageVsdClient

func TestMain(m *testing.M) {
	kubemonConfig = &config.NuageKubeMonConfig{}
	addArgs(kubemonConfig, flag.CommandLine)
	flag.CommandLine.Parse(os.Args[1:])
	vsdClient = new(NuageVsdClient)
	vsdClient.domains = make(map[string]string)
	vsdClient.version = kubemonConfig.NuageVspVersion
	vsdClient.url = kubemonConfig.NuageVsdApiUrl + "/nuage/api/" + vsdClient.version + "/"
	vsdClient.CreateSession()
	if err := vsdClient.GetAuthorizationToken(); err != nil {
		fmt.Printf("Error getting VSD auth token: %s\n", err)
		vsdClient = nil
	}
	os.Exit(m.Run())
}

func addArgs(myConfig *config.NuageKubeMonConfig, flagSet *flag.FlagSet) {
	flagSet.StringVar(&myConfig.OsClusterAdmin, "osusername",
		"system:admin", "User name of the cluster administrator")
	flagSet.StringVar(&myConfig.KubeConfigFile, "kubeconfig",
		"", "kubeconfig File for Openshift User")
	flagSet.StringVar(&myConfig.OsMasterConfigFile, "osmasterconfig",
		"", "Path to master-config.yaml for the cluster master")
	flagSet.StringVar(&myConfig.NuageVsdApiUrl, "nuagevsdurl",
		"", "Nuage VSD URL")
	flagSet.StringVar(&myConfig.NuageVspVersion, "nuagevspversion",
		"", "Nuage VSP Version")
	flagSet.StringVar(&myConfig.LicenseFile, "license_file",
		"", "VSD License File Path")
}
