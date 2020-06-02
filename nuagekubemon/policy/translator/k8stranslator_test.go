package translator

import (
	"fmt"
	"os"
	"testing"

	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/policyapi/implementer"
	gomega "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	kapi "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var appFs = afero.NewOsFs()

const (
	VsdURL       = "https://10.4.1.2:7443"
	UserCertFile = "/tmp/trans-usercert.pem"
	UserKeyFile  = "/tmp/trans-userkey.pem"
	ENTERPRISE   = "centos-operator"
	DOMAIN       = "centos-operator-1-14"
	PolicyName   = "nuagekubemon-test"

	PgServer1 = "ServerPG"
	PgClient1 = "ClientPG1"
	PgClient2 = "ClientPG2"
)

var policyImplementer implementer.PolicyImplementer

func init() {
	var vsdCredentials implementer.VSDCredentials
	// this is a fake certificate generated.
	userCert := `-----BEGIN CERTIFICATE-----
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
	userKey := `-----BEGIN PRIVATE KEY-----
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

	file, _ := appFs.OpenFile(UserCertFile, os.O_RDWR|os.O_CREATE, 0600)
	_, err := file.WriteString(userCert)
	if err != nil {
		er := fmt.Errorf("Cannot write to file %s", err)
		fmt.Println(er)
	}
	file2, _ := appFs.OpenFile(UserKeyFile, os.O_RDWR|os.O_CREATE, 0600)
	_, err = file2.WriteString(userKey)
	if err != nil {
		er := fmt.Errorf("Cannot write to file %s", err)
		fmt.Println(er)
	}
	vsdCredentials.UserCertFile = UserCertFile
	vsdCredentials.UserKeyFile = UserKeyFile
	vsdCredentials.URL = VsdURL

	if err := policyImplementer.Init(&vsdCredentials); err != nil {
		er := fmt.Errorf("Unable to connect to VSD %s", err)
		fmt.Println(er)
	}
}

func TestNuageKubemonK8SPolicyCreation(t *testing.T) {
	pgMap := make(map[string]api.PgInfo)

	var tcp kapi.Protocol = kapi.ProtocolTCP
	var udp kapi.Protocol = kapi.ProtocolTCP
	var port1 = intstr.IntOrString{Type: intstr.Int, IntVal: 1000}
	var port2 = intstr.IntOrString{Type: intstr.Int, IntVal: 2000}

	networkPolicyPort1 := networkingV1.NetworkPolicyPort{Protocol: &tcp, Port: &port1}
	networkPolicyPort2 := networkingV1.NetworkPolicyPort{Protocol: &udp, Port: &port2}
	networkPolicyPorts := []networkingV1.NetworkPolicyPort{networkPolicyPort1, networkPolicyPort2}

	var podSelector1 metav1.LabelSelector
	podSelector1.MatchLabels = make(map[string]string)
	podSelector1.MatchLabels["openstack.io/client"] = "client1"
	pod1SelectorKey, err := metav1.LabelSelectorAsSelector(&podSelector1)
	if err != nil {
		t.Fatalf("Unable to create pod1 selector key")
	}
	client1PgInfo := api.PgInfo{PgName: PgClient1, Selector: podSelector1}
	pgMap[pod1SelectorKey.String()] = client1PgInfo

	var podSelector2 metav1.LabelSelector
	podSelector2.MatchLabels = make(map[string]string)
	podSelector2.MatchLabels["openstack.io/client"] = "client2"
	pod2SelectorKey, err := metav1.LabelSelectorAsSelector(&podSelector2)
	if err != nil {
		t.Fatalf("Unable to create pod2 selector key")
	}
	client2PgInfo := api.PgInfo{PgName: PgClient2, Selector: podSelector2}
	pgMap[pod2SelectorKey.String()] = client2PgInfo

	networkPolicyPeer1 := networkingV1.NetworkPolicyPeer{PodSelector: &podSelector1}
	networkPolicyPeer2 := networkingV1.NetworkPolicyPeer{PodSelector: &podSelector2}
	networkPolicyPeers := []networkingV1.NetworkPolicyPeer{networkPolicyPeer1, networkPolicyPeer2}

	ingressRule := networkingV1.NetworkPolicyIngressRule{Ports: networkPolicyPorts,
		From: networkPolicyPeers}

	var targetLabel metav1.LabelSelector
	targetLabel.MatchLabels = make(map[string]string)
	targetLabel.MatchLabels["openstack.io/server"] = "server"
	server1PgInfo := api.PgInfo{PgName: PgServer1, Selector: targetLabel}

	targetKey, err := metav1.LabelSelectorAsSelector(&targetLabel)
	if err != nil {
		t.Fatalf("Unable to create target selector key")
	}
	pgMap[targetKey.String()] = server1PgInfo
	var label = make(map[string]string)
	label["openstack.io/server"] = "server"
	label["nuage.io/priority"] = "1234"
	networkPolicySpec := networkingV1.NetworkPolicySpec{PodSelector: targetLabel}
	networkPolicySpec.Ingress = []networkingV1.NetworkPolicyIngressRule{ingressRule}
	networkPolicyEvent := api.NetworkPolicyEvent{Type: api.Added, Name: "test-policy-2", Namespace: "default", Policy: networkPolicySpec, Labels: label}
	metadata := make(map[string]string)
	metadata["enterpriseName"] = ENTERPRISE
	metadata["domainName"] = DOMAIN
	nuagePolicy, err := CreateNuagePGPolicy(&networkPolicyEvent, pgMap, metadata)

	if err != nil {
		t.Fatalf("Error creating the nuage policy %+v", err)
	}

	d, err := yaml.Marshal(&nuagePolicy)
	if err != nil {
		t.Fatalf("Error while marshalling %+v", err)
	}
	t.Logf("Marshalled YAML %s", string(d))

	err = policyImplementer.ImplementPolicy(nuagePolicy)
	g := gomega.NewGomegaWithT(t)
	g.Expect(err).To(gomega.HaveOccurred())
}
