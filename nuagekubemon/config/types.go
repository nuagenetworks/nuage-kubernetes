/*
###########################################################################
#
#   Filename:           types.go
#
#   Author:             Aniket Bhat
#   Created:            July 20, 2015
#
#   Description:        NuageKubeMon config types
#
###########################################################################
#
#              Copyright (c) 2015 Nuage Networks
#
###########################################################################

*/

package config

import (
	"gopkg.in/yaml.v2"
)

type NuageKubeMonConfig struct {
	KubeConfigFile     string       `yaml:"kubeConfig"`
	OsClusterAdmin     string       `yaml:"openshiftAdmin"`
	OsMasterConfigFile string       `yaml:"openshiftMasterConfig"`
	NuageVsdApiUrl     string       `yaml:"vsdApiUrl"`
	NuageVspVersion    string       `yaml:"vspVersion"`
	LicenseFile        string       `yaml:"licenseFile"`
	EnterpriseName     string       `yaml:"enterpriseName"`
	DomainName         string       `yaml:"domainName"`
	ConfigFile         string       `yaml:"-"` // yaml tag `-` denotes that this cannot be supplied in yaml.
	OsMasterConfig     MasterConfig `yaml:"-"`
}

type openshiftNetworkConfig struct {
	ClusterCIDR  string `yaml:"clusterNetworkCIDR"`
	SubnetLength int    `yaml:"hostSubnetLength"`
}

/* Fields we care about in the openshift master-config.yaml */
type MasterConfig struct {
	NetworkConfig openshiftNetworkConfig `yaml:"networkConfig"`
}

type NamespaceMap map[string]bool

func (conf *NuageKubeMonConfig) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, conf); err != nil {
		return err
	}
	// TODO: Bounds checking and other validation on fields
	if conf.EnterpriseName == "" {
		conf.EnterpriseName = "K8S-Enterprise"
	}
	if conf.DomainName == "" {
		conf.DomainName = "K8S-Domain"
	}
	return nil
}

func (conf *MasterConfig) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, conf); err != nil {
		return err
	}
	// TODO: Bounds checking and other validation on fields
	return nil
}
