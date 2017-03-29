/*
###########################################################################
#
#   Filename:           api.go
#
#   Author:             Aniket Bhat
#   Created:            July 20, 2015
#
#   Description:        NuageKubeMon event API
#
###########################################################################
#
#              Copyright (c) 2015 Nuage Networks
#
###########################################################################

*/

package api

import (
	"fmt"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

type PgInfo struct {
	PgName   string
	PgId     string
	Selector unversioned.LabelSelector
}

type EventType string

const (
	Added    EventType = "ADDED"
	Deleted  EventType = "DELETED"
	Modified EventType = "MODIFIED"
)

const (
	PATEnabled   = "ENABLED"
	PATInherited = "INHERITED"
	PATDisabled  = "DISABLED"
)

const (
	IngressAclTemplateName     = "Auto-generated Ingress Policies"
	EgressAclTemplateName      = "Auto-generated Egress Policies"
	ZoneAnnotationTemplateName = "Namespace Annotations"
)

type Namespace string

type NamespaceEvent struct {
	Type        EventType
	Name        string
	Annotations map[string]string `json:"annotations,omitempty"`
}

type ServiceEvent struct {
	Type        EventType
	Name        string
	ClusterIP   string
	Namespace   string
	NuageLabels map[string]string `json:"labels,omitempty"`
}

type NetworkPolicyEvent struct {
	Type      EventType
	Name      string
	Namespace string
	Policy    extensions.NetworkPolicySpec
	Labels    map[string]string `json:"labels,omitempty"`
}

type PodEvent struct {
	Type      EventType
	Name      string
	Namespace string
	Labels    map[string]string `json:"labels,omitempty"`
}

type GetPod func(string, string) (*PodEvent, error)
type FilterPods func(*kapi.ListOptions, string) (*[]*PodEvent, error)
type FilterNamespaces func(*kapi.ListOptions) (*[]*NamespaceEvent, error)
type FilterServices func(*kapi.ListOptions) (*[]*ServiceEvent, error)
type FilterNetworkPolicies func(*kapi.ListOptions) (*[]*NetworkPolicyEvent, error)

type ClusterClientCallBacks struct {
	FilterPods FilterPods
	GetPod     GetPod
}

type RESTError struct {
	Errors []struct {
		Property     string `json:"property"`
		Descriptions []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		} `json:"descriptions"`
	} `json:"errors"`
	InternalErrorCode int `json:"internalErrorCode"`
}

func (restErr RESTError) String() string {
	outString := fmt.Sprintf("InternalErrorCode: %d\n",
		restErr.InternalErrorCode)
	for _, vsdErr := range restErr.Errors {
		outString += "\tProperty: " + vsdErr.Property + "\n"
		for _, description := range vsdErr.Descriptions {
			outString += "\t\tTitle: " + description.Title +
				"\n\t\tDescription: " + description.Description + "\n"
		}
	}
	return outString
}

type VsdEnterprise struct {
	Description string `json:"description"`
	Name        string `json:"name"`
	ID          string
}

type VsdUser struct {
	ID        string
	UserName  string `json:"userName"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type VsdGroup struct {
	ID   string
	Role string `json:"role"`
}

type VsdLicense struct {
	ID        string
	License   string `json:"license"`
	LicenseId int    `json:"licenseID"`
}

type VsdSubnet struct {
	ID              string
	IPType          string
	Name            string `json:"name"`
	Address         string `json:"address"`
	Netmask         string `json:"netmask"`
	Description     string `json:"description"`
	PATEnabled      string
	UnderlayEnabled string `json:"underlayEnabled,omitempty"`
}

// Generic VSD object. Most json objects returned by the VSD REST API will fit
// this "interface"
type VsdObject struct {
	ID          string
	Name        string `json:"name"`
	Description string `json:"description"`
}

type VsdDomain struct {
	ID              string
	Name            string `json:"name"`
	Description     string `json:"description"`
	TemplateID      string `json:"templateID"`
	PATEnabled      string
	UnderlayEnabled string `json:"underlayEnabled,omitempty"`
}

type VsdAuthToken struct {
	APIKey       string
	APIKeyExpiry int64
	ID           string
	Email        string `json:"email"`
	EnterpriseID string `json:"enterpriseID"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Role         string `json:"role"`
	UserName     string `json:"userName"`
}

type VsdAclTemplate struct {
	ID                string
	Name              string `json:"name"`
	DefaultAllowIP    bool   `json:"defaultAllowIP"`
	DefaultAllowNonIP bool   `json:"defaultAllowNonIP"`
	Active            bool   `json:"active"`
	Priority          int    `json:"priority"`
}

type VsdAclEntry struct {
	DSCP         string `json:"DSCP,omitempty"`
	ID           string
	Action       string `json:"action"`
	Description  string `json:"description"`
	EntityScope  string `json:"entityScope"`
	EtherType    string `json:"etherType"`
	LocationID   string `json:"locationID"`
	LocationType string `json:"locationType"`
	NetworkID    string `json:"networkID"`
	NetworkType  string `json:"networkType"`
	PolicyState  string `json:"policyState"`
	Priority     int    `json:"priority"`
	Protocol     string `json:"protocol"`
	Reflexive    bool   `json:"reflexive"`
}

const MAX_VSD_ACL_PRIORITY = 1000000000 //the maximum priority allowed in VSD is 1 billion.

type VsdNetworkMacro struct {
	ID      string
	Name    string `json:"name"`
	IPType  string `json:"IPType"`
	Address string `json:"address"`
	Netmask string `json:"netmask"`
}

type VsdJobType string

const (
	BeginPolicyChanges   VsdJobType = "BEGIN_POLICY_CHANGES"
	ApplyPolicyChanges   VsdJobType = "APPLY_POLICY_CHANGES"
	DiscardPolicyChanges VsdJobType = "DISCARD_POLICY_CHANGES"
)

func (jobType VsdJobType) String() string {
	strType := "INVALID"

	switch jobType {
	case BeginPolicyChanges:
		strType = "BEGIN_POLICY_CHANGES"
	case ApplyPolicyChanges:
		strType = "APPLY_POLICY_CHANGES"
	case DiscardPolicyChanges:
		strType = "DISCARD_POLICY_CHANGES"
	}

	return strType
}

type VsdJob struct {
	ID              string      `json:"ID,omitempty"`
	ParentID        string      `json:"parentID,omitempty"`
	ParentType      string      `json:"parentType,omitempty"`
	Owner           string      `json:"owner,omitempty"`
	Parameters      interface{} `json:"parameters,omitempty"`
	LastUpdatedBy   string      `json:"lastUpdatedBy,omitempty"`
	Result          interface{} `json:"result,omitempty"`
	EntityScope     string      `json:"entityScope,omitempty"`
	Command         string      `json:"command,omitempty"`
	Progress        float64     `json:"progress,omitempty"`
	AssocEntityType string      `json:"assocEntityType,omitempty"`
	Status          string      `json:"status,omitempty"`
	ExternalID      string      `json:"externalID,omitempty"`
}

func (lhs *VsdNetworkMacro) IsEqual(rhs *VsdNetworkMacro) bool {
	if lhs.Name != rhs.Name {
		return false
	}
	if lhs.IPType != rhs.IPType {
		return false
	}
	if lhs.Address != rhs.Address {
		return false
	}
	if lhs.Netmask != rhs.Netmask {
		return false
	}
	return true
}

func (acl *VsdAclEntry) TryNextAclPriority() {
	if acl.Priority == MAX_VSD_ACL_PRIORITY {
		acl.Priority = acl.Priority - 1
	} else {
		acl.Priority = acl.Priority + 1
	}
}

func (lhs *VsdAclEntry) IsEqual(rhs *VsdAclEntry) bool {
	if lhs.DSCP != "" && lhs.DSCP != rhs.DSCP {
		return false
	}
	if lhs.Action != "" && lhs.Action != rhs.Action {
		return false
	}
	if lhs.EntityScope != "" && lhs.EntityScope != rhs.EntityScope {
		return false
	}
	if lhs.EtherType != "" && lhs.EtherType != rhs.EtherType {
		return false
	}
	if lhs.LocationID != "" && lhs.LocationID != rhs.LocationID {
		return false
	}
	if lhs.LocationType != "" && lhs.LocationType != rhs.LocationType {
		return false
	}
	if lhs.NetworkID != "" && lhs.NetworkID != rhs.NetworkID {
		return false
	}
	if lhs.NetworkType != "" && lhs.NetworkType != rhs.NetworkType {
		return false
	}
	if lhs.PolicyState != "" && lhs.PolicyState != rhs.PolicyState {
		return false
	}
	if lhs.Protocol != "" && lhs.Protocol != rhs.Protocol {
		return false
	}
	return true
}

func (lhs *VsdAclEntry) BuildFilter() string {
	filter := ""
	if lhs.DSCP != "" {
		dscpClause := `DSCP == "` + lhs.DSCP + `"`
		filter = dscpClause
	}
	if lhs.Action != "" {
		actionClause := `action == "` + lhs.Action + `"`
		if filter != "" {
			filter = filter + ` and ` + actionClause
		} else {
			filter = actionClause
		}
	}
	// The entity scope is not filterable currently in VSD

	// if lhs.EntityScope != "" {
	// 	entityScopeClause := `entityScope == "` + lhs.EntityScope + `"`
	// 	if filter != "" {
	// 		filter = filter + ` and ` + entityScopeClause
	// 	} else {
	// 		filter = entityScopeClause
	// 	}
	// }
	if lhs.EtherType != "" {
		etherTypeClause := `etherType == "` + lhs.EtherType + `"`
		if filter != "" {
			filter = filter + ` and ` + etherTypeClause
		} else {
			filter = etherTypeClause
		}
	}
	if lhs.LocationID != "" {
		locationIDClause := `locationID == "` + lhs.LocationID + `"`
		if filter != "" {
			filter = filter + ` and ` + locationIDClause
		} else {
			filter = locationIDClause
		}
	}
	if lhs.LocationType != "" {
		locationTypeClause := `locationType == "` + lhs.LocationType + `"`
		if filter != "" {
			filter = filter + ` and ` + locationTypeClause
		} else {
			filter = locationTypeClause
		}
	}
	if lhs.NetworkID != "" {
		networkIDClause := `networkID == "` + lhs.NetworkID + `"`
		if filter != "" {
			filter = filter + ` and ` + networkIDClause
		} else {
			filter = networkIDClause
		}
	}
	if lhs.NetworkType != "" {
		networkTypeClause := `networkType == "` + lhs.NetworkType + `"`
		if filter != "" {
			filter = filter + ` and ` + networkTypeClause
		} else {
			filter = networkTypeClause
		}
	}
	if lhs.PolicyState != "" {
		policyStateClause := `policyState == "` + lhs.PolicyState + `"`
		if filter != "" {
			filter = filter + ` and ` + policyStateClause
		} else {
			filter = policyStateClause
		}
	}
	if lhs.Protocol != "" {
		protocolClause := `protocol == "` + lhs.Protocol + `"`
		if filter != "" {
			filter = filter + ` and ` + protocolClause
		} else {
			filter = protocolClause
		}
	}
	return filter
}

func (lhs *VsdAclEntry) String() string {
	return fmt.Sprintf(`\nACL Entry: ID: %v, Description: %v,\n`+
		`Priority: %v, Action: %v,\n`+
		`DSCP: %v, EntityScope: %v, EtherType: %v, Protocol: %v\n`+
		`LocationID: %v, LocationType: %v\n`+
		`NetworkID: %v, NetworkType: %v, PolicyState: %v, Reflexive %v`,
		lhs.ID, lhs.Description, lhs.Priority, lhs.Action, lhs.DSCP,
		lhs.EntityScope, lhs.EtherType, lhs.Protocol, lhs.LocationID,
		lhs.LocationType, lhs.NetworkID, lhs.NetworkType, lhs.PolicyState,
		lhs.Reflexive)
}

func (svc *ServiceEvent) String() string {
	return fmt.Sprintf(`\nService: Name: %v, Namespace %v,\n`+
		`ClusterIP: %v, Labels: %v,\n`+
		`EventType: %v`, svc.Name, svc.Namespace, svc.ClusterIP,
		svc.NuageLabels, svc.Type)
}

func (lhs *VsdNetworkMacro) String() string {
	return fmt.Sprintf(`\nNetwork Macro: Name: %v, ID: %v,\n`+
		`IPType: %v, Address: %v,\n`+
		`Netmask: %v`, lhs.Name, lhs.ID, lhs.IPType, lhs.Address, lhs.Netmask)
}
