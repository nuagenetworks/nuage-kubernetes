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

type NuageKubeMonConfig struct {
	KubeConfigFile     string
	OsClusterAdmin     string
	OsMasterConfigFile string
	NuageVsdApiUrl     string
	NuageVspVersion    string
	LogDir             string
	LicenseFile        string
}

type NamespaceMap map[string]bool
