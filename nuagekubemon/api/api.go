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

type EventType string

const (
	Added   EventType = "ADDED"
	Deleted EventType = "DELETED"
)

const (
	PATEnabled   = "ENABLED"
	PATInherited = "INHERITED"
	PATDisabled  = "DISABLED"
)

type Namespace string

type NamespaceEvent struct {
	Type EventType
	Name string
}

type RESTError struct {
	Message string
	Errors  []struct {
		Resource string
		Field    string
		Code     string
	}
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
	ID          string
	IPType      string
	Name        string `json:"name"`
	Address     string `json:"address"`
	Netmask     string `json:"netmask"`
	Description string `json:"description"`
	PATEnabled  string
}

// Generic VSD object. Most json objects returned by the VSD REST API will fit
// this "interface"
type VsdObject struct {
	ID          string
	Name        string `json:"name"`
	Description string `json:"description"`
}

type VsdDomain struct {
	ID          string
	Name        string `json:"name"`
	Description string `json:"description"`
	TemplateID  string `json:"templateID"`
	PATEnabled  string
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
}
