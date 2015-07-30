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
	Children        string `json:"children"`
	ParentType      string `json:"parentType"`
	LastUpdatedBy   string `json:"lastUpdatedBy"`
	LastUpdatedDate int64  `json:"lastUpdateDate"`
	creationDate    int64
	avatarData      string
	avatarType      string
	Description     string `json:"description"`
	Name            string `json:"name"`
	owner           string
	ID              string
	parentID        string
	externalID      string
	customerID      string
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

// Generic VSD object. Most json objects returned by the VSD REST API will fit
// this "interface"
type VsdObject struct {
	ID          string
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Create a vsd object from a template.  Like VsdObject, but also contains the
// ID of the template.
type VsdObjectInstance struct {
	ID          string
	Name        string `json:"name"`
	Description string `json:"description"`
	TemplateID  string `json:"templateID"`
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
