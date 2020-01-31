package translate

import (
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"
	"github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/pkg/policyapi/policies"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//GetEnterpriseFunc get enterprise from VSD
type GetEnterpriseFunc func(string) (string, error)

//GetNetworkMacroFunc fetch network macro from VSD
type GetNetworkMacroFunc func(string, string) (*api.VsdNetworkMacro, error)

//CreatePgFunc creates pg on VSD
type CreatePgFunc func(string, string) (string, string, error)

//DeletePgFunc deletes pg on VSD
type DeletePgFunc func(string) error

//CreateNetworkMacroFunc creates network macro on VSD
type CreateNetworkMacroFunc func(*api.VsdNetworkMacro) (string, error)

//DeleteNetworkMacroFunc delete network macro on VSD
type DeleteNetworkMacroFunc func(string) error

//AddPortsToPgFunc add vports to a policy group
type AddPortsToPgFunc func(string, []string) error

//DeletePortsFromPgFunc delete ports from a policy group
type DeletePortsFromPgFunc func(string) error

//VSDMetaData vsd metadata map
type VSDMetaData map[string]string

//NamespaceLabelsMap label to list of namespaces mapping
type NamespaceLabelsMap map[string][]string

//NWMacroMap map of network macros on VSD per zone
type NWMacroMap map[string]NWMacroInfo

//NWMacroExceptMap map of network macros for except CIDR on VSD per zone
type NWMacroExceptMap map[string]NWMacroInfo

//SelectorMap map of selectors to policy groups
type SelectorMap map[string]*PgInfo

//PgMap map of namespace to policies in a namespace
type PgMap map[string]SelectorMap

// CallBacks has callbacks to VSD API
type CallBacks struct {
	GetEnterprise      GetEnterpriseFunc
	GetNetworkMacro    GetNetworkMacroFunc
	AddPg              CreatePgFunc
	DeletePg           DeletePgFunc
	AddPortsToPg       AddPortsToPgFunc
	DeletePortsFromPg  DeletePortsFromPgFunc
	AddNetworkMacro    CreateNetworkMacroFunc
	DeleteNetworkMacro DeleteNetworkMacroFunc
}

// PgInfo policy group info for a selector
type PgInfo struct {
	PgName     string
	PgID       string
	PolicyName string
	Selector   metav1.LabelSelector
	RefCount   int
}

// NWMacroInfo network macro info for an ip block
type NWMacroInfo struct {
	Name     string
	ID       string
	CIDR     string
	RefCount int
}

// NWMacroGroup network macro group
type NWMacroGroup struct {
	Name     string
	ID       string
	NWMacros []NWMacroInfo
}

// VSDObjsMap map of network policy peers to vsd objects
type VSDObjsMap struct {
	PGMap       PgMap
	NWMacroMap  NWMacroMap
	NSLabelsMap NamespaceLabelsMap
}

// PolicyData contains metadata for policy element
type PolicyData struct {
	Name       string
	SourceName string
	TargetName string
	SourceType policies.EndPointType
	TargetType policies.EndPointType
	Action     policies.ActionType
}

// InitVSDObjsMap initializes vsd objects map
func InitVSDObjsMap() *VSDObjsMap {
	m := &VSDObjsMap{}
	m.PGMap = make(PgMap)
	m.NWMacroMap = make(NWMacroMap)
	m.NSLabelsMap = make(NamespaceLabelsMap)
	return m
}
