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

// EnterpriseSecurityIdentity represents the Identity of the object
var EnterpriseSecurityIdentity = bambou.Identity{
	Name:     "enterprisesecurity",
	Category: "enterprisesecurities",
}

// EnterpriseSecuritiesList represents a list of EnterpriseSecurities
type EnterpriseSecuritiesList []*EnterpriseSecurity

// EnterpriseSecuritiesAncestor is the interface that an ancestor of a EnterpriseSecurity must implement.
// An Ancestor is defined as an entity that has EnterpriseSecurity as a descendant.
// An Ancestor can get a list of its child EnterpriseSecurities, but not necessarily create one.
type EnterpriseSecuritiesAncestor interface {
	EnterpriseSecurities(*bambou.FetchingInfo) (EnterpriseSecuritiesList, *bambou.Error)
}

// EnterpriseSecuritiesParent is the interface that a parent of a EnterpriseSecurity must implement.
// A Parent is defined as an entity that has EnterpriseSecurity as a child.
// A Parent is an Ancestor which can create a EnterpriseSecurity.
type EnterpriseSecuritiesParent interface {
	EnterpriseSecuritiesAncestor
	CreateEnterpriseSecurity(*EnterpriseSecurity) *bambou.Error
}

// EnterpriseSecurity represents the model of a enterprisesecurity
type EnterpriseSecurity struct {
	ID                      string `json:"ID,omitempty"`
	ParentID                string `json:"parentID,omitempty"`
	ParentType              string `json:"parentType,omitempty"`
	Owner                   string `json:"owner,omitempty"`
	LastUpdatedBy           string `json:"lastUpdatedBy,omitempty"`
	GatewaySecurityRevision int    `json:"gatewaySecurityRevision,omitempty"`
	Revision                int    `json:"revision,omitempty"`
	EnterpriseID            string `json:"enterpriseID,omitempty"`
	EntityScope             string `json:"entityScope,omitempty"`
	ExternalID              string `json:"externalID,omitempty"`
}

// NewEnterpriseSecurity returns a new *EnterpriseSecurity
func NewEnterpriseSecurity() *EnterpriseSecurity {

	return &EnterpriseSecurity{}
}

// Identity returns the Identity of the object.
func (o *EnterpriseSecurity) Identity() bambou.Identity {

	return EnterpriseSecurityIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *EnterpriseSecurity) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *EnterpriseSecurity) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the EnterpriseSecurity from the server
func (o *EnterpriseSecurity) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the EnterpriseSecurity into the server
func (o *EnterpriseSecurity) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the EnterpriseSecurity from the server
func (o *EnterpriseSecurity) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the EnterpriseSecurity
func (o *EnterpriseSecurity) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the EnterpriseSecurity
func (o *EnterpriseSecurity) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the EnterpriseSecurity
func (o *EnterpriseSecurity) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the EnterpriseSecurity
func (o *EnterpriseSecurity) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EnterpriseSecuredDatas retrieves the list of child EnterpriseSecuredDatas of the EnterpriseSecurity
func (o *EnterpriseSecurity) EnterpriseSecuredDatas(info *bambou.FetchingInfo) (EnterpriseSecuredDatasList, *bambou.Error) {

	var list EnterpriseSecuredDatasList
	err := bambou.CurrentSession().FetchChildren(o, EnterpriseSecuredDataIdentity, &list, info)
	return list, err
}

// CreateEnterpriseSecuredData creates a new child EnterpriseSecuredData under the EnterpriseSecurity
func (o *EnterpriseSecurity) CreateEnterpriseSecuredData(child *EnterpriseSecuredData) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
