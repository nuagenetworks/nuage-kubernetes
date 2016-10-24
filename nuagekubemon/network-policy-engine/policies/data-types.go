package policies

import "fmt"

type EndPointType string

const (
	ZONE         EndPointType = "ZONE"
	SUBNET       EndPointType = "SUBNET"
	POLICY_GROUP EndPointType = "POLICYGROUP"
)

type EndPoint struct {
	Type EndPointType
	Name string
}

type ActionType string

const (
	ALLOW ActionType = "FORWARD"
	DENY  ActionType = "DROP"
)

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

	return "*"
}

type Protocol int

const (
	ANY Protocol = 0
	TCP Protocol = 6
	UDP Protocol = 17
)

func (protocol Protocol) String() string {
	switch protocol {
	case ANY:
		return "ANY"
	case TCP:
		return "6"
	case UDP:
		return "17"
	}
	return "-1"
}

type NetworkParameters struct {
	Protocol             Protocol
	SourcePortRange      PortRange `yaml:"source-port-range"`
	DestinationPortRange PortRange `yaml:"destination-port-range"`
}
