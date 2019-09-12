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

// DUCGroupIdentity represents the Identity of the object
var DUCGroupIdentity = bambou.Identity{
	Name:     "ducgroup",
	Category: "ducgroups",
}

// DUCGroupsList represents a list of DUCGroups
type DUCGroupsList []*DUCGroup

// DUCGroupsAncestor is the interface that an ancestor of a DUCGroup must implement.
// An Ancestor is defined as an entity that has DUCGroup as a descendant.
// An Ancestor can get a list of its child DUCGroups, but not necessarily create one.
type DUCGroupsAncestor interface {
	DUCGroups(*bambou.FetchingInfo) (DUCGroupsList, *bambou.Error)
}

// DUCGroupsParent is the interface that a parent of a DUCGroup must implement.
// A Parent is defined as an entity that has DUCGroup as a child.
// A Parent is an Ancestor which can create a DUCGroup.
type DUCGroupsParent interface {
	DUCGroupsAncestor
	CreateDUCGroup(*DUCGroup) *bambou.Error
}

// DUCGroup represents the model of a ducgroup
type DUCGroup struct {
	ID                             string        `json:"ID,omitempty"`
	ParentID                       string        `json:"parentID,omitempty"`
	ParentType                     string        `json:"parentType,omitempty"`
	Owner                          string        `json:"owner,omitempty"`
	Name                           string        `json:"name,omitempty"`
	LastUpdatedBy                  string        `json:"lastUpdatedBy,omitempty"`
	Description                    string        `json:"description,omitempty"`
	EmbeddedMetadata               []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope                    string        `json:"entityScope,omitempty"`
	AssociatedPerformanceMonitorID string        `json:"associatedPerformanceMonitorID,omitempty"`
	DucMeshGroupID                 int           `json:"ducMeshGroupID,omitempty"`
	Function                       string        `json:"function,omitempty"`
	ExternalID                     string        `json:"externalID,omitempty"`
}

// NewDUCGroup returns a new *DUCGroup
func NewDUCGroup() *DUCGroup {

	return &DUCGroup{
		Function: "UBR",
	}
}

// Identity returns the Identity of the object.
func (o *DUCGroup) Identity() bambou.Identity {

	return DUCGroupIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *DUCGroup) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *DUCGroup) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the DUCGroup from the server
func (o *DUCGroup) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the DUCGroup into the server
func (o *DUCGroup) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the DUCGroup from the server
func (o *DUCGroup) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the DUCGroup
func (o *DUCGroup) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the DUCGroup
func (o *DUCGroup) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the DUCGroup
func (o *DUCGroup) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the DUCGroup
func (o *DUCGroup) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSGateways retrieves the list of child NSGateways of the DUCGroup
func (o *DUCGroup) NSGateways(info *bambou.FetchingInfo) (NSGatewaysList, *bambou.Error) {

	var list NSGatewaysList
	err := bambou.CurrentSession().FetchChildren(o, NSGatewayIdentity, &list, info)
	return list, err
}

// AssignNSGateways assigns the list of NSGateways to the DUCGroup
func (o *DUCGroup) AssignNSGateways(children NSGatewaysList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, NSGatewayIdentity)
}
