package policies

import "fmt"

// EndPointType type of Network Endpoint
type EndPointType string

// Types of network endpoints
const (
	Zone         EndPointType = "ZONE"
	Subnet       EndPointType = "SUBNET"
	PolicyGroup  EndPointType = "POLICYGROUP"
	EndPointZone EndPointType = "ENDPOINT_ZONE"
	Invalid      EndPointType = "INVALID"
)

// EndPoint identifies a network endpoint
type EndPoint struct {
	Type EndPointType
	Name string
}

// ActionType identifies the policy action
type ActionType string

// Types of policy actions
const (
	Allow ActionType = "ALLOW"
	Deny  ActionType = "DENY"
)

// PortRange defines port range of a policy element
type PortRange struct {
	StartPort int `yaml:"start-port"`
	EndPort   int `yaml:"end-port"`
}

func (portRange PortRange) String() string {

	if portRange.StartPort == 0 && portRange.EndPort == 0 {
		return "*"
	} else if portRange.StartPort == portRange.EndPort {
		return fmt.Sprintf("%d", portRange.StartPort)
	} else {
		return fmt.Sprintf("%d-%d", portRange.StartPort, portRange.EndPort)
	}
}

// Protocol network protocol
type Protocol int

// Types of network protocols
const (
	TCP Protocol = 6
	UDP Protocol = 17
)

func (protocol Protocol) String() string {
	switch protocol {
	case TCP:
		return "6"
	case UDP:
		return "17"
	}
	return "-1"
}

// NetworkParameters network parameters for a policy element
type NetworkParameters struct {
	Protocol             Protocol
	SourcePortRange      PortRange `yaml:"source-port-range"`
	DestinationPortRange PortRange `yaml:"destination-port-range"`
}
