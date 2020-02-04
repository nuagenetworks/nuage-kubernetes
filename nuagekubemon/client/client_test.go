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
	"os"
	"os/exec"
	"testing"

	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/config"
)

var kubemonConfig *config.NuageKubeMonConfig
var vsdClient *NuageVsdClient
var isTargetSystem bool

func TestMain(m *testing.M) {
	kubemonConfig = &config.NuageKubeMonConfig{}
	addArgs(kubemonConfig, flag.CommandLine)
	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		isTargetSystem = false
	}

	// this is a fake certificate generated.
	userCertFile := `-----BEGIN CERTIFICATE-----
MIIDBzCCAe+gAwIBAgIJAOLz9k9vdwloMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNV
BAMMD3d3dy5leGFtcGxlLmNvbTAeFw0yMDAyMjgyMzE5MThaFw0zMDAyMjUyMzE5
MThaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBALSEqYCu8ghfuZnNKUVZ30WeEmW6RWVNankcXlbnxZjq
9bhdrxMDMma1KCGFCPvIE8scFq539vtBbjYs/ZVAPhq/dTubz7lr8PvgRI50VpOj
5PbJg1w9E0Is6HH33hnnfIjRUW0eJTk7UJzD9h5S4MutfNhF6OJ7JcxPVNosEzm5
UW+Ydpo0L1ExjSalcZFoLhjGN8xZRFmpSfSZHntZqP52DhboPHOJpQs84JJGCghs
wIifeQCEHZH7dYXE0AYRZqjwUenH5SG/R151sKEKANex9MXR77nVlXWU3udJ2JRt
+vn8JaagiBy5vrdrU4EWgkU8zDpFG4uMYSf44GbLoZsCAwEAAaNQME4wHQYDVR0O
BBYEFD1EMpXxuRhGgmBQjX4jorOkck4EMB8GA1UdIwQYMBaAFD1EMpXxuRhGgmBQ
jX4jorOkck4EMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBABmZPPu9
SVlo2psSX09Eyku//1pAcnqHf5VATQQIwr1M9f7ag3v4FmAcl+/LMrgjNkwaWLbg
iK6KS31fF0+ortoYPtxgiZ4uGEqqFw/PghPlq3OfYT4XkfjWJ4RS/r9zltW9llVe
/EeUAtOXA8AYGLUT5BZDX/AFAhPKoSQHuhJ98SnqLffdgwkIqqpr23vJTPT2zdlx
nSUx8nCOCS6cw/ZPS4wsmUF6j9QatprMG50wo6NL1SbyNzkGuFgxAV7WrPPeTGvV
WeOSFkV7TFj6i8ofNFT0mddBU5OEVq9qcv8mpOclzx3I/iyjCk1FIPZ0DV8WLQOt
6ic0UXJ5xo4CzI0=
-----END CERTIFICATE-----`

	// this is a fake private key generated.
	userKeyFile := `-----BEGIN PRIVATE KEY-----
MIIEpAIBAAKCAQEAtISpgK7yCF+5mc0pRVnfRZ4SZbpFZU1qeRxeVufFmOr1uF2v
EwMyZrUoIYUI+8gTyxwWrnf2+0FuNiz9lUA+Gr91O5vPuWvw++BEjnRWk6Pk9smD
XD0TQizocffeGed8iNFRbR4lOTtQnMP2HlLgy6182EXo4nslzE9U2iwTOblRb5h2
mjQvUTGNJqVxkWguGMY3zFlEWalJ9Jkee1mo/nYOFug8c4mlCzzgkkYKCGzAiJ95
AIQdkft1hcTQBhFmqPBR6cflIb9HXnWwoQoA17H0xdHvudWVdZTe50nYlG36+fwl
pqCIHLm+t2tTgRaCRTzMOkUbi4xhJ/jgZsuhmwIDAQABAoIBAQCEkwV1d4ZTdhH2
DYGo6CcclsnGIjYC/wcaKSZzxsYM10pc+5ivauKiIZt2eqCtYTSAL4HM4lfmERii
+wnFiifSNxgfDgBRmh+irANNZ82Joo1uXXJ21HgHWrnfsX1RIvwH80pMzB3kWVaL
uzNO8+kaTLBqmXU+l9ibowubK1F3SxHsOUiRfxLdE/UP449FEq7W5aesAOeECT/t
w/DQTahkNaKwAknMKl3DnHqlHGpMYvHNiqBt3LoS4FQyOdPkZbVZy80Sp0ZkhdiB
ZisrhUQJl3c+yP9Zi13SG1iG+PGja6r/Ki1+asMwAhaJJ0uYc/LSP+bZD2CzsnFS
NNL8xXE5AoGBANeeIlpYhrxk/DYSZjw0ThaDO0p3/wmBS4iE85YDsgvvxF+NlWVv
Qn47/pTdz8mS1cbDqNVEn+Wglf+5cpTFunboVasGlMKbUA3Ks8CKU6fnfFCfndB3
uU/iVTi/2OM88d4L7qUnRQVnKl1qOyMrHxOX01jx8JNkDkOMNduSH1QXAoGBANZT
q5Tgu70ZHbfyOtHQNkbpyOjXnMwjvuolccndxKlDyU76/f7DBpuuFwEW5EI7Rzwl
SPFwdDBJNTZ4JpsbGW+3ERDANfFlc4MAjOag0YH1Micmfq45pPakNEIRdMC0bJiG
Lj23cWpjMlyZMJFaNUMaxgou00UYHthLpMfF250dAoGBAMP58EFruzMbGn5PJOtN
ozglGUvrWzyZbzzrkrbkLv1YdYVgG8zxXl98Sj2mikktk+6wQhFt6WN+HTgsp29/
dKbFL7BeL/Hd1tpiRhUX5Ud0SHLDUV58o0tvbYRCI3EPIMtwzvz/f2WUylXTy2KA
vCND2Q48AS0GQUy18PHck2sLAoGAC2gomaPcWhQcIM4jk0chnGSU7M+M6NB+OLgF
dlj3Por9C9cP7Z8zmtWJI+W0AFJnWCwj1bXGeUtsKZn7dAXdNLTpk5qnRFHB9Bbz
aNLmU6RZJvxFgcBPp1DV9y42qIrxvKxniaFZx++/nm4Ix7OlYgzqvWAAnozKF3jv
LDK7nYECgYBhfOmpSh8VSbQyNHGPSWUFPPQEvW+TX2NRV5iY718Xvg9g5ctKwNFB
qslkZ3GGpfUmuIzPM3tFJJi9lnGOqTGixoTkvsT7hY9QjwoMWeWMEJyM0pWvyYjH
G3PY7QEvUYkh3lD36FAQxssSDuZZb0kHmTGEeR/oAhXqrwOJmh1HWA==
-----END PRIVATE KEY-----`

	if testing.Short() {
		isTargetSystem = false
	} else {
		vsdClient = new(NuageVsdClient)
		vsdClient.namespaces = make(map[string]NamespaceData)
		vsdClient.version = kubemonConfig.NuageVspVersion
		vsdClient.url = kubemonConfig.NuageVsdApiUrl + "/nuage/api/" + vsdClient.version + "/"
		vsdClient.CreateSession(userCertFile, userKeyFile)
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
