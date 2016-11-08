package implementer

import (
	"fmt"
	"github.com/nuagenetworks/openshift-integration/nuagekubemon/network-policy-engine/policies"
	"testing"
)

const (
	VSD_URL          = "https://192.168.103.200:8443"
	VSD_USERNAME     = "csproot"
	VSD_PASSWORD     = "csproot"
	VSD_ORGANIZATION = "csp"
	CLIENT_PG        = "ClientPG1"
	SERVER_PG        = "ServerPG"

	ENTERPRISE  = "nuage"
	DOMAIN      = "openshift"
	POLICY_NAME = "k8s-policy"
	POLICY_ID   = POLICY_NAME
)

var policyImplementer PolicyImplementer

func init() {
	var vsdCredentials VSDCredentials
	vsdCredentials.Username = VSD_USERNAME
	vsdCredentials.Password = VSD_PASSWORD
	vsdCredentials.Organization = VSD_ORGANIZATION
	vsdCredentials.URL = VSD_URL

	if err := policyImplementer.Init(&vsdCredentials); err != nil {
		fmt.Errorf("Unable to connect to VSD")
	}
}

func addPolicy() error {
	nuagePolicy := policies.NuagePolicy{
		Version:    policies.V1Alpha,
		Type:       policies.DEFAULT,
		Enterprise: ENTERPRISE,
		Domain:     DOMAIN,
		Name:       POLICY_NAME,
		ID:         POLICY_ID,
		Priority:   10000,
	}

	defaultPolicyElement := policies.DefaultPolicyElement{
		Name:   "Access Control",
		From:   policies.EndPoint{Name: CLIENT_PG, Type: policies.POLICY_GROUP},
		To:     policies.EndPoint{Name: SERVER_PG, Type: policies.POLICY_GROUP},
		Action: policies.ALLOW,
		NetworkParameters: policies.NetworkParameters{
			Protocol:             policies.TCP,
			DestinationPortRange: policies.PortRange{100, 200},
		},
	}

	nuagePolicy.PolicyElements = []policies.DefaultPolicyElement{defaultPolicyElement}
	err := policyImplementer.ImplementPolicy(&nuagePolicy)
	return err
}

func TestPolicyAdd(t *testing.T) {
	err := addPolicy()
	if err != nil {
		t.Fatalf("Failed to apply policy with error %+v", err)
	}
}

func TestPolicyRemove(t *testing.T) {
	err := addPolicy()
	if err != nil {
		t.Fatalf("Failed to apply policy with error %+v", err)
	}

	err = policyImplementer.DeletePolicy(POLICY_ID, ENTERPRISE, DOMAIN)
	if err != nil {
		t.Fatalf("Unable to delete the policy %+v", err)
	}
}
