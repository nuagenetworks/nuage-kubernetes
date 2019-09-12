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

// PolicyObjectGroupIdentity represents the Identity of the object
var PolicyObjectGroupIdentity = bambou.Identity{
	Name:     "policyobjectgroup",
	Category: "policyobjectgroups",
}

// PolicyObjectGroupsList represents a list of PolicyObjectGroups
type PolicyObjectGroupsList []*PolicyObjectGroup

// PolicyObjectGroupsAncestor is the interface that an ancestor of a PolicyObjectGroup must implement.
// An Ancestor is defined as an entity that has PolicyObjectGroup as a descendant.
// An Ancestor can get a list of its child PolicyObjectGroups, but not necessarily create one.
type PolicyObjectGroupsAncestor interface {
	PolicyObjectGroups(*bambou.FetchingInfo) (PolicyObjectGroupsList, *bambou.Error)
}

// PolicyObjectGroupsParent is the interface that a parent of a PolicyObjectGroup must implement.
// A Parent is defined as an entity that has PolicyObjectGroup as a child.
// A Parent is an Ancestor which can create a PolicyObjectGroup.
type PolicyObjectGroupsParent interface {
	PolicyObjectGroupsAncestor
	CreatePolicyObjectGroup(*PolicyObjectGroup) *bambou.Error
}

// PolicyObjectGroup represents the model of a policyobjectgroup
type PolicyObjectGroup struct {
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
	Type             string        `json:"type,omitempty"`
}

// NewPolicyObjectGroup returns a new *PolicyObjectGroup
func NewPolicyObjectGroup() *PolicyObjectGroup {

	return &PolicyObjectGroup{}
}

// Identity returns the Identity of the object.
func (o *PolicyObjectGroup) Identity() bambou.Identity {

	return PolicyObjectGroupIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *PolicyObjectGroup) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *PolicyObjectGroup) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the PolicyObjectGroup from the server
func (o *PolicyObjectGroup) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the PolicyObjectGroup into the server
func (o *PolicyObjectGroup) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the PolicyObjectGroup from the server
func (o *PolicyObjectGroup) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the PolicyObjectGroup
func (o *PolicyObjectGroup) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the PolicyObjectGroup
func (o *PolicyObjectGroup) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the PolicyObjectGroup
func (o *PolicyObjectGroup) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the PolicyObjectGroup
func (o *PolicyObjectGroup) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSGateways retrieves the list of child NSGateways of the PolicyObjectGroup
func (o *PolicyObjectGroup) NSGateways(info *bambou.FetchingInfo) (NSGatewaysList, *bambou.Error) {

	var list NSGatewaysList
	err := bambou.CurrentSession().FetchChildren(o, NSGatewayIdentity, &list, info)
	return list, err
}

// AssignNSGateways assigns the list of NSGateways to the PolicyObjectGroup
func (o *PolicyObjectGroup) AssignNSGateways(children NSGatewaysList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, NSGatewayIdentity)
}
