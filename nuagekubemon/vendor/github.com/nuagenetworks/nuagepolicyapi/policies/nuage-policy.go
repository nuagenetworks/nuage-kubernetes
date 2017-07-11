package policies

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

// PolicyType identifies the type of policy
type PolicyType string

// Types of policies supported
const (
	Default PolicyType = "default"
)

// Version policy version
type Version string

// Supported versions of policies
const (
	V1Alpha Version = "v1-alpha"
)

// PolicyUpdateOperation defines the types of update ops
type PolicyUpdateOperation int

// Supported policy updates
const (
	UpdateAdd    PolicyUpdateOperation = 1
	UpdateRemove PolicyUpdateOperation = 2
)

// NuagePolicy idenfies a Nuage policy
type NuagePolicy struct {
	Version    Version
	Type       PolicyType
	Enterprise string
	Domain     string

	Name           string
	ID             string
	Priority       int
	PolicyElements interface{} `yaml:"policy-elements"`
}

// LoadPolicyFromYAML creates nuage policy object from a YAML string
func LoadPolicyFromYAML(policyYAML string) (*NuagePolicy, error) {
	var nuagePolicy NuagePolicy
	err := yaml.Unmarshal([]byte(policyYAML), &nuagePolicy)

	if err != nil {
		return nil, err
	}

	if nuagePolicy.Version != V1Alpha {
		return nil, fmt.Errorf("Invalid policy version %+v", nuagePolicy.Version)
	}

	policyElementsYAML, err := yaml.Marshal(nuagePolicy.PolicyElements)
	if err != nil {
		return nil, fmt.Errorf("Error extracting the policy elements %+v", err)
	}

	switch nuagePolicy.Type {
	case Default:
		var defaultPolicyElements []DefaultPolicyElement
		err = yaml.Unmarshal(policyElementsYAML, &defaultPolicyElements)
		if err != nil {
			return nil, fmt.Errorf("Error extracting the default policy elements %+v", err)
		}
		nuagePolicy.PolicyElements = defaultPolicyElements
	default:
		return nil, fmt.Errorf("Invalid policy type %+v", nuagePolicy.Type)
	}

	return &nuagePolicy, nil
}

// ConvertPolicyActionToNuageAction converts policy actions to nuage action
// types
func ConvertPolicyActionToNuageAction(action ActionType) string {
	switch action {
	case Allow:
		return "FORWARD"
	case Deny:
		return "DROP"
	default:
		panic(fmt.Sprintf("Invalid action %s", action))
	}
}

// ConvertPolicyEndPointStringToEndPointType converts policy endpoint types to
// Nuage endpoint types
func ConvertPolicyEndPointStringToEndPointType(endPointTypeString string) (EndPointType, error) {
	endptTypeStr := strings.ToUpper(strings.Replace(endPointTypeString, "-", "", -1))
	switch endptTypeStr {
	case "ZONE":
		return Zone, nil
	case "SUBNET":
		return Subnet, nil
	case "POLICYGROUP":
		return PolicyGroup, nil
	case "ENDPOINT_ZONE":
		return EndPointZone, nil
	case "ENDPOINTZONE":
		return EndPointZone, nil
	}

	return Invalid, fmt.Errorf(fmt.Sprintf("Invalid endpoint type %s %s", endPointTypeString, endptTypeStr))
}
