/*
###########################################################################
#
#   Filename:           types.go
#
#   Author:             Aniket Bhat
#   Created:            July 20, 2015
#
#   Description:        Nuage Monitor config types
#
###########################################################################
#
#              Copyright (c) 2015 Nuage Networks
#
###########################################################################

*/

package config

import (
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

type NuageKubeMonConfig struct {
	KubeConfigFile      string           `yaml:"kubeConfig"`
	MasterConfigFile    string           `yaml:"masterConfig"`
	NuageVsdApiUrl      string           `yaml:"vsdApiUrl"`
	NuageVspVersion     string           `yaml:"vspVersion"`
	LicenseFile         string           `yaml:"licenseFile"`
	EnterpriseName      string           `yaml:"enterpriseName"`
	DomainName          string           `yaml:"domainName"`
	StatsLogging        string           `yaml:"statsLogging"`
	RestServer          RestServerConfig `yaml:"nuageMonServer"`
	UserCertificateFile string           `yaml:"userCertificateFile"`
	UserKeyFile         string           `yaml:"userKeyFile"`
	PrivilegedProject   []string         `yaml:"privilegedProject"`
	PrivilegedNamespace []string         `yaml:"privilegedNamespace"`
	ConfigFile          string           `yaml:"-"` // yaml tag `-` denotes that this cannot be supplied in yaml.
	MasterConfig        MasterConfig     `yaml:"-"`
	EtcdClientConfig    EtcdConfig       `yaml:"etcdClientConfig"`
	AutoScaleSubnets    string           `yaml:"autoScaleSubnets"`
	UnderlaySupport     string           `yaml:"underlaySupport"`
	EncryptionEnabled   string           `yaml:"encryptionEnabled"`
}

type RestServerConfig struct {
	Url                   string `yaml:"URL"`
	CertificateDirectory  string `yaml:"certificateDirectory"`
	ClientCA              string `yaml:"clientCA"`
	ServerCertificate     string `yaml:"serverCertificate"`
	ServerKey             string `yaml:"serverKey"`
	ClientCAData          string `yaml:"clientCAData"`
	ServerCertificateData string `yaml:"serverCertificateData"`
	ServerKeyData         string `yaml:"serverKeyData"`
}

type networkConfig struct {
	ClusterNetworks []struct {
		CIDR         string `yaml:"cidr"`
		SubnetLength int    `yaml:"hostSubnetLength"`
	} `yaml:"clusterNetworks"`
	ServiceCIDR string `yaml:"serviceNetworkCIDR"`
}

/* Fields we care about in the openshift master-config.yaml */
type MasterConfig struct {
	NetworkConfig networkConfig `yaml:"networkConfig"`
}

type EtcdConfig struct {
	ServerCA          string   `yaml:"ca"`
	ClientCertificate string   `yaml:"certFile"`
	ClientKey         string   `yaml:"keyFile"`
	UrlList           []string `yaml:"urls"`
}

type NamespaceUpdateEvent int

const (
	AddSubnet NamespaceUpdateEvent = iota
	DeleteSubnet
)

type NamespaceUpdateRequest struct {
	NamespaceID string //Name of the namespace in the NamespaceData map
	Event       NamespaceUpdateEvent
}

type NamespaceMap map[string]bool

func DefaultEnterprise() string {
	programName := path.Base(os.Args[0])

	enterprise := "Openshift-Enterprise"
	if strings.ToLower(programName) == "nuagekubemon" {
		enterprise = "K8S-Enterprise"
	}

	return enterprise
}

func DefaultDomain() string {
	programName := path.Base(os.Args[0])

	domain := "Openshift-Domain"
	if strings.ToLower(programName) == "nuagekubemon" {
		domain = "K8S-Domain"
	}

	return domain
}

func (conf *NuageKubeMonConfig) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, conf); err != nil {
		return err
	}

	// TODO: Bounds checking and other validation on fields
	if conf.EnterpriseName == "" {
		conf.EnterpriseName = DefaultEnterprise()
	}

	if conf.DomainName == "" {
		conf.DomainName = DefaultDomain()
	}

	if conf.PrivilegedNamespace == nil {
		conf.PrivilegedNamespace = []string{"kube-system", "default"}
	}

	// To simplify execution, we'll use PrivilegedProject everywhere after
	// configuration is done.  If the system is nuagekubemon, we'll overwrite
	// the PrivilegedProject variable with the PrivilegedNamespace one.
	conf.PrivilegedProject = conf.PrivilegedNamespace

	return nil
}

func (conf *MasterConfig) Parse(data []byte) error {
	if err := yaml.Unmarshal(data, conf); err != nil {
		return err
	}
	// TODO: Bounds checking and other validation on fields
	return nil
}
