package config

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestReadMasterConfig(t *testing.T) {
	const masterConfigFile = "../testfiles/master-config.yaml"
	const expectedClusterCIDR = "172.30.0.0/16"
	masterConfigData, err := ioutil.ReadFile(masterConfigFile)
	if err != nil {
		t.Fatalf("Failed to read %s! Could not complete this test.",
			masterConfigFile)
	}
	//fmt.Println("masterConfigData:\n", string(masterConfigData))
	myConfig := &MasterConfig{}
	if err := myConfig.Parse(masterConfigData); err != nil {
		t.Fatalf("Failed to parse %s! Error: %s\n", masterConfigFile, err)
	}
	if myConfig.NetworkConfig.ClusterCIDR != expectedClusterCIDR {
		t.Fatalf("clusterNetworkCIDR mismatch! Expected: %q, Got: %q\n",
			expectedClusterCIDR, myConfig.NetworkConfig.ClusterCIDR)
	}
}

func TestReadKubemonConfig(t *testing.T) {
	const kubemonConfigFile = "../testfiles/nuagekubemon-config.yaml"
	const (
		kubeConfig            = "admin.kubeconfig"
		openshiftAdmin        = "system:admin"
		openshiftMasterConfig = "master-config.yaml"
		vsdApiUrl             = "https://xmpp.example.com:8443"
		vspVersion            = "v3_2"
		licenseFile           = "base_vsp_license.txt"
	)
	kubemonConfigData, err := ioutil.ReadFile(kubemonConfigFile)
	if err != nil {
		t.Fatalf("Failed to read %s! Could not complete this test.",
			kubemonConfigFile)
	}
	fmt.Printf("kubemonConfigData:\n%s", string(kubemonConfigData))
	myConfig := &NuageKubeMonConfig{}
	if err := myConfig.Parse(kubemonConfigData); err != nil {
		t.Fatalf("Failed to parse %s! Error: %s\n", kubemonConfigFile, err)
	}
	if myConfig.KubeConfigFile != kubeConfig {
		t.Fatalf("kubeConfig mismatch! Expected: %q, Got: %q",
			kubeConfig, myConfig.KubeConfigFile)
	}

	if myConfig.MasterConfigFile != openshiftMasterConfig {
		t.Fatalf("openshiftMasterConfig mismatch! Expected: %q, Got: %q",
			openshiftMasterConfig, myConfig.MasterConfigFile)
	}
	if myConfig.NuageVsdApiUrl != vsdApiUrl {
		t.Fatalf("vsdApiUrl mismatch! Expected: %q, Got: %q",
			vsdApiUrl, myConfig.NuageVsdApiUrl)
	}
	if myConfig.NuageVspVersion != vspVersion {
		t.Fatalf("vspVersion mismatch! Expected: %q, Got: %q",
			vspVersion, myConfig.NuageVspVersion)
	}
	if myConfig.LicenseFile != licenseFile {
		t.Fatalf("licenseFile mismatch! Expected: %q, Got: %q",
			licenseFile, myConfig.LicenseFile)
	}
}
