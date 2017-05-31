/*
  Copyright (c) 2015, Alcatel-Lucent Inc
  All rights reserved.

  Redistribution and use in source and binary forms, with or without
  modification, are permitted provided that the following conditions are met:
      * Redistributions of source code must retain the above copyright
        notice, this list of conditions and the following disclaimer.
      * Redistributions in binary form must reproduce the above copyright
        notice, this list of conditions and the following disclaimer in the
        documentation and/or other materials provided with the distribution.
      * Neither the name of the copyright holder nor the names of its contributors
        may be used to endorse or promote products derived from this software without
        specific prior written permission.

  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
  ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
  WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
  DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY
  DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
  (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
  LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
  ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
  (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
  SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package vspk

import "github.com/nuagenetworks/go-bambou/bambou"

// EnterprisePermissionIdentity represents the Identity of the object
var EnterprisePermissionIdentity = bambou.Identity{
	Name:     "enterprisepermission",
	Category: "enterprisepermissions",
}

// EnterprisePermissionsList represents a list of EnterprisePermissions
type EnterprisePermissionsList []*EnterprisePermission

// EnterprisePermissionsAncestor is the interface that an ancestor of a EnterprisePermission must implement.
// An Ancestor is defined as an entity that has EnterprisePermission as a descendant.
// An Ancestor can get a list of its child EnterprisePermissions, but not necessarily create one.
type EnterprisePermissionsAncestor interface {
	EnterprisePermissions(*bambou.FetchingInfo) (EnterprisePermissionsList, *bambou.Error)
}

// EnterprisePermissionsParent is the interface that a parent of a EnterprisePermission must implement.
// A Parent is defined as an entity that has EnterprisePermission as a child.
// A Parent is an Ancestor which can create a EnterprisePermission.
type EnterprisePermissionsParent interface {
	EnterprisePermissionsAncestor
	CreateEnterprisePermission(*EnterprisePermission) *bambou.Error
}

// EnterprisePermission represents the model of a enterprisepermission
type EnterprisePermission struct {
	ID                         string `json:"ID,omitempty"`
	ParentID                   string `json:"parentID,omitempty"`
	ParentType                 string `json:"parentType,omitempty"`
	Owner                      string `json:"owner,omitempty"`
	Name                       string `json:"name,omitempty"`
	LastUpdatedBy              string `json:"lastUpdatedBy,omitempty"`
	PermittedAction            string `json:"permittedAction,omitempty"`
	PermittedEntityDescription string `json:"permittedEntityDescription,omitempty"`
	PermittedEntityID          string `json:"permittedEntityID,omitempty"`
	PermittedEntityName        string `json:"permittedEntityName,omitempty"`
	PermittedEntityType        string `json:"permittedEntityType,omitempty"`
	EntityScope                string `json:"entityScope,omitempty"`
	ExternalID                 string `json:"externalID,omitempty"`
}

// NewEnterprisePermission returns a new *EnterprisePermission
func NewEnterprisePermission() *EnterprisePermission {

	return &EnterprisePermission{
		PermittedAction: "USE",
	}
}

// Identity returns the Identity of the object.
func (o *EnterprisePermission) Identity() bambou.Identity {

	return EnterprisePermissionIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *EnterprisePermission) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *EnterprisePermission) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the EnterprisePermission from the server
func (o *EnterprisePermission) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the EnterprisePermission into the server
func (o *EnterprisePermission) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the EnterprisePermission from the server
func (o *EnterprisePermission) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the EnterprisePermission
func (o *EnterprisePermission) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the EnterprisePermission
func (o *EnterprisePermission) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the EnterprisePermission
func (o *EnterprisePermission) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the EnterprisePermission
func (o *EnterprisePermission) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
