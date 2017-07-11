package translator

import (
	"fmt"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuagepolicyapi/implementer"
	"gopkg.in/yaml.v2"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/util/intstr"
	"testing"
)

const (
	VSD_URL          = "https://192.168.103.200:8443"
	VSD_USERNAME     = "csproot"
	VSD_PASSWORD     = "csproot"
	VSD_ORGANIZATION = "csp"

	ENTERPRISE  = "nuage"
	DOMAIN      = "openshift"
	POLICY_NAME = "nuagekubemon-test"

	PG_SERVER1 = "ServerPG"
	PG_CLIENT1 = "ClientPG1"
	PG_CLIENT2 = "ClientPG2"
)

var policyImplementer implementer.PolicyImplementer

func init() {
	var vsdCredentials implementer.VSDCredentials
	vsdCredentials.Username = VSD_USERNAME
	vsdCredentials.Password = VSD_PASSWORD
	vsdCredentials.Organization = VSD_ORGANIZATION
	vsdCredentials.URL = VSD_URL

	if err := policyImplementer.Init(&vsdCredentials); err != nil {
		fmt.Errorf("Unable to connect to VSD")
	}
}

func TestNuageKubemonK8SPolicyCreation(t *testing.T) {
	pgMap := make(map[string]api.PgInfo)

	var tcp kapi.Protocol = kapi.ProtocolTCP
	var udp kapi.Protocol = kapi.ProtocolTCP
	var port1 = intstr.IntOrString{Type: intstr.Int, IntVal: 1000}
	var port2 = intstr.IntOrString{Type: intstr.Int, IntVal: 2000}

	networkPolicyPort1 := extensions.NetworkPolicyPort{Protocol: &tcp, Port: &port1}
	networkPolicyPort2 := extensions.NetworkPolicyPort{Protocol: &udp, Port: &port2}
	networkPolicyPorts := []extensions.NetworkPolicyPort{networkPolicyPort1, networkPolicyPort2}

	var podSelector1 unversioned.LabelSelector
	podSelector1.MatchLabels = make(map[string]string)
	podSelector1.MatchLabels["openstack.io/client"] = "client1"
	pod1SelectorKey, err := unversioned.LabelSelectorAsSelector(&podSelector1)
	if err != nil {
		t.Fatalf("Unable to create pod1 selector key")
	}
	client1PgInfo := api.PgInfo{PgName: PG_CLIENT1, Selector: podSelector1}
	pgMap[pod1SelectorKey.String()] = client1PgInfo

	var podSelector2 unversioned.LabelSelector
	podSelector2.MatchLabels = make(map[string]string)
	podSelector2.MatchLabels["openstack.io/client"] = "client2"
	pod2SelectorKey, err := unversioned.LabelSelectorAsSelector(&podSelector2)
	if err != nil {
		t.Fatalf("Unable to create pod2 selector key")
	}
	client2PgInfo := api.PgInfo{PgName: PG_CLIENT2, Selector: podSelector2}
	pgMap[pod2SelectorKey.String()] = client2PgInfo

	networkPolicyPeer1 := extensions.NetworkPolicyPeer{PodSelector: &podSelector1}
	networkPolicyPeer2 := extensions.NetworkPolicyPeer{PodSelector: &podSelector2}
	networkPolicyPeers := []extensions.NetworkPolicyPeer{networkPolicyPeer1, networkPolicyPeer2}

	ingressRule := extensions.NetworkPolicyIngressRule{Ports: networkPolicyPorts,
		From: networkPolicyPeers}

	var targetLabel unversioned.LabelSelector
	targetLabel.MatchLabels = make(map[string]string)
	targetLabel.MatchLabels["openstack.io/server"] = "server"
	server1PgInfo := api.PgInfo{PgName: PG_SERVER1, Selector: targetLabel}

	targetKey, err := unversioned.LabelSelectorAsSelector(&targetLabel)
	if err != nil {
		t.Fatalf("Unable to create target selector key")
	}
	pgMap[targetKey.String()] = server1PgInfo

	networkPolicySpec := extensions.NetworkPolicySpec{PodSelector: targetLabel}
	networkPolicySpec.Ingress = []extensions.NetworkPolicyIngressRule{ingressRule}

	metadata := make(map[string]string)
	metadata["enterpriseName"] = ENTERPRISE
	metadata["domainName"] = DOMAIN
	nuagePolicy, err := CreateNuagePGPolicy(&networkPolicySpec, "test-policy-2", pgMap, metadata)

	if err != nil {
		t.Fatalf("Error creating the nuage policy %+v", err)
	}

	d, err := yaml.Marshal(&nuagePolicy)
	if err != nil {
		t.Fatalf("Error while marshalling %+v", err)
	}
	t.Logf("Marshalled YAML %s", string(d))

	err = policyImplementer.ImplementPolicy(nuagePolicy)
	if err != nil {
		t.Fatalf("Failed to apply policy with error %+v", err)
	}
}
