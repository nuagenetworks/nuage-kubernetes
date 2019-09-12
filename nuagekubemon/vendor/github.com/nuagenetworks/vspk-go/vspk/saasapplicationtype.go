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

// SaaSApplicationTypeIdentity represents the Identity of the object
var SaaSApplicationTypeIdentity = bambou.Identity{
	Name:     "saasapplicationtype",
	Category: "saasapplicationtypes",
}

// SaaSApplicationTypesList represents a list of SaaSApplicationTypes
type SaaSApplicationTypesList []*SaaSApplicationType

// SaaSApplicationTypesAncestor is the interface that an ancestor of a SaaSApplicationType must implement.
// An Ancestor is defined as an entity that has SaaSApplicationType as a descendant.
// An Ancestor can get a list of its child SaaSApplicationTypes, but not necessarily create one.
type SaaSApplicationTypesAncestor interface {
	SaaSApplicationTypes(*bambou.FetchingInfo) (SaaSApplicationTypesList, *bambou.Error)
}

// SaaSApplicationTypesParent is the interface that a parent of a SaaSApplicationType must implement.
// A Parent is defined as an entity that has SaaSApplicationType as a child.
// A Parent is an Ancestor which can create a SaaSApplicationType.
type SaaSApplicationTypesParent interface {
	SaaSApplicationTypesAncestor
	CreateSaaSApplicationType(*SaaSApplicationType) *bambou.Error
}

// SaaSApplicationType represents the model of a saasapplicationtype
type SaaSApplicationType struct {
	ID               string        `json:"ID,omitempty"`
	ParentID         string        `json:"parentID,omitempty"`
	ParentType       string        `json:"parentType,omitempty"`
	Owner            string        `json:"owner,omitempty"`
	Name             string        `json:"name,omitempty"`
	LastUpdatedBy    string        `json:"lastUpdatedBy,omitempty"`
	ReadOnly         bool          `json:"readOnly"`
	Description      string        `json:"description,omitempty"`
	EmbeddedMetadata []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope      string        `json:"entityScope,omitempty"`
	ExternalID       string        `json:"externalID,omitempty"`
}

// NewSaaSApplicationType returns a new *SaaSApplicationType
func NewSaaSApplicationType() *SaaSApplicationType {

	return &SaaSApplicationType{
		ReadOnly: false,
	}
}

// Identity returns the Identity of the object.
func (o *SaaSApplicationType) Identity() bambou.Identity {

	return SaaSApplicationTypeIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *SaaSApplicationType) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *SaaSApplicationType) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the SaaSApplicationType from the server
func (o *SaaSApplicationType) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the SaaSApplicationType into the server
func (o *SaaSApplicationType) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the SaaSApplicationType from the server
func (o *SaaSApplicationType) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the SaaSApplicationType
func (o *SaaSApplicationType) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the SaaSApplicationType
func (o *SaaSApplicationType) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the SaaSApplicationType
func (o *SaaSApplicationType) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the SaaSApplicationType
func (o *SaaSApplicationType) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EnterpriseNetworks retrieves the list of child EnterpriseNetworks of the SaaSApplicationType
func (o *SaaSApplicationType) EnterpriseNetworks(info *bambou.FetchingInfo) (EnterpriseNetworksList, *bambou.Error) {

	var list EnterpriseNetworksList
	err := bambou.CurrentSession().FetchChildren(o, EnterpriseNetworkIdentity, &list, info)
	return list, err
}

// CreateEnterpriseNetwork creates a new child EnterpriseNetwork under the SaaSApplicationType
func (o *SaaSApplicationType) CreateEnterpriseNetwork(child *EnterpriseNetwork) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
