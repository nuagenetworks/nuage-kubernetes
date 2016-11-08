package policies

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"testing"
)

const (
	CLIENT_PG = "ClientPG"
	SERVER_PG = "ServerPG"
)

func init() {
	fmt.Println("Initing test bundle")
}

func TestDefaultPolicyMarshalling(t *testing.T) {
	nuagePolicy := NuagePolicy{
		Version:    V1Alpha,
		Type:       DEFAULT,
		Enterprise: "nuage",
		Domain:     "openshift",
		Name:       "k8s allow traffic",
		ID:         "k8s allow traffic",
		Priority:   10000,
	}

	defaultPolicyElement := DefaultPolicyElement{
		Name:   "Access Control",
		From:   EndPoint{Name: CLIENT_PG, Type: POLICY_GROUP},
		To:     EndPoint{Name: SERVER_PG, Type: POLICY_GROUP},
		Action: ALLOW,
		NetworkParameters: NetworkParameters{
			Protocol:             TCP,
			DestinationPortRange: PortRange{100, 200},
		},
	}

	nuagePolicy.PolicyElements = []DefaultPolicyElement{defaultPolicyElement}
	d, err := yaml.Marshal(&nuagePolicy)
	if err != nil {
		t.Fatalf("Error while marshalling %+v", err)
	}
	t.Logf("Marshalled YAML %s", string(d))
}

const testYaml = `
--- 
version: v1
type: default
enterprise: nuage
domain: openshift
id: "k8s allow web traffic"
name: "k8s allow web traffic"
policy-elements: 
    - name: "Access control"
      from: 
        name: busybox
        type: policy-group
      to: 
        name: nginx
        type: policy-group
      action: ALLOW
      network-parameters:
        protocol: 6
        destination-port-range: 
          start-port: 80
          end-port: 80
        source-port-range:
          start-port: 0
          end-port: 65535
          
`

func TestDefaultPolicyUnMarshalling(t *testing.T) {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(testYaml), &m)

	if err != nil {
		t.Fatalf("Unable to unmarshal the default policy %+v", err)
	}
	t.Logf("Unmarshalled data %+v", m)
	t.Fatalf("Incomplete implementation")
}
