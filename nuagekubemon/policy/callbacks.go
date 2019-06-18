package policy

import "github.com/nuagenetworks/nuage-kubernetes/nuagekubemon/api"

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
