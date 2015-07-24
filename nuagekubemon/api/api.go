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

type RESTErrorResponse struct {
	Message string
	Errors  []struct {
		Resource string
		Field    string
		Code     string
	}
}

type VsdEnterpriseResponse struct {
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

type VsdAdminCreateResponse struct {
	UserName string `json:"userName"`
	ID       string
}

type VsdAdminGroupResponse struct {
	Role string `json:"role"`
	ID   string
}

type VsdCreateLicenseResponse struct {
	License   string `json:"license"`
	LicenseId int    `json:"licenseID"`
	ID        string
}
