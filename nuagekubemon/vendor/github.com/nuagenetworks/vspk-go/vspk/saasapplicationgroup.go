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

// SaaSApplicationGroupIdentity represents the Identity of the object
var SaaSApplicationGroupIdentity = bambou.Identity{
	Name:     "saasapplicationgroup",
	Category: "saasapplicationgroups",
}

// SaaSApplicationGroupsList represents a list of SaaSApplicationGroups
type SaaSApplicationGroupsList []*SaaSApplicationGroup

// SaaSApplicationGroupsAncestor is the interface that an ancestor of a SaaSApplicationGroup must implement.
// An Ancestor is defined as an entity that has SaaSApplicationGroup as a descendant.
// An Ancestor can get a list of its child SaaSApplicationGroups, but not necessarily create one.
type SaaSApplicationGroupsAncestor interface {
	SaaSApplicationGroups(*bambou.FetchingInfo) (SaaSApplicationGroupsList, *bambou.Error)
}

// SaaSApplicationGroupsParent is the interface that a parent of a SaaSApplicationGroup must implement.
// A Parent is defined as an entity that has SaaSApplicationGroup as a child.
// A Parent is an Ancestor which can create a SaaSApplicationGroup.
type SaaSApplicationGroupsParent interface {
	SaaSApplicationGroupsAncestor
	CreateSaaSApplicationGroup(*SaaSApplicationGroup) *bambou.Error
}

// SaaSApplicationGroup represents the model of a saasapplicationgroup
type SaaSApplicationGroup struct {
	ID               string        `json:"ID,omitempty"`
	ParentID         string        `json:"parentID,omitempty"`
	ParentType       string        `json:"parentType,omitempty"`
	Owner            string        `json:"owner,omitempty"`
	Name             string        `json:"name,omitempty"`
	LastUpdatedBy    string        `json:"lastUpdatedBy,omitempty"`
	Description      string        `json:"description,omitempty"`
	EmbeddedMetadata []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope      string        `json:"entityScope,omitempty"`
	ExternalID       string        `json:"externalID,omitempty"`
}

// NewSaaSApplicationGroup returns a new *SaaSApplicationGroup
func NewSaaSApplicationGroup() *SaaSApplicationGroup {

	return &SaaSApplicationGroup{}
}

// Identity returns the Identity of the object.
func (o *SaaSApplicationGroup) Identity() bambou.Identity {

	return SaaSApplicationGroupIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *SaaSApplicationGroup) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *SaaSApplicationGroup) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the SaaSApplicationGroup from the server
func (o *SaaSApplicationGroup) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the SaaSApplicationGroup into the server
func (o *SaaSApplicationGroup) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the SaaSApplicationGroup from the server
func (o *SaaSApplicationGroup) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// SaaSApplicationTypes retrieves the list of child SaaSApplicationTypes of the SaaSApplicationGroup
func (o *SaaSApplicationGroup) SaaSApplicationTypes(info *bambou.FetchingInfo) (SaaSApplicationTypesList, *bambou.Error) {

	var list SaaSApplicationTypesList
	err := bambou.CurrentSession().FetchChildren(o, SaaSApplicationTypeIdentity, &list, info)
	return list, err
}

// AssignSaaSApplicationTypes assigns the list of SaaSApplicationTypes to the SaaSApplicationGroup
func (o *SaaSApplicationGroup) AssignSaaSApplicationTypes(children SaaSApplicationTypesList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, SaaSApplicationTypeIdentity)
}

// Metadatas retrieves the list of child Metadatas of the SaaSApplicationGroup
func (o *SaaSApplicationGroup) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the SaaSApplicationGroup
func (o *SaaSApplicationGroup) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the SaaSApplicationGroup
func (o *SaaSApplicationGroup) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the SaaSApplicationGroup
func (o *SaaSApplicationGroup) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
