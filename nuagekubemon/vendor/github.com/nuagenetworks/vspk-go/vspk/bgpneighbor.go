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

// BGPNeighborIdentity represents the Identity of the object
var BGPNeighborIdentity = bambou.Identity{
	Name:     "bgpneighbor",
	Category: "bgpneighbors",
}

// BGPNeighborsList represents a list of BGPNeighbors
type BGPNeighborsList []*BGPNeighbor

// BGPNeighborsAncestor is the interface that an ancestor of a BGPNeighbor must implement.
// An Ancestor is defined as an entity that has BGPNeighbor as a descendant.
// An Ancestor can get a list of its child BGPNeighbors, but not necessarily create one.
type BGPNeighborsAncestor interface {
	BGPNeighbors(*bambou.FetchingInfo) (BGPNeighborsList, *bambou.Error)
}

// BGPNeighborsParent is the interface that a parent of a BGPNeighbor must implement.
// A Parent is defined as an entity that has BGPNeighbor as a child.
// A Parent is an Ancestor which can create a BGPNeighbor.
type BGPNeighborsParent interface {
	BGPNeighborsAncestor
	CreateBGPNeighbor(*BGPNeighbor) *bambou.Error
}

// BGPNeighbor represents the model of a bgpneighbor
type BGPNeighbor struct {
	ID                              string `json:"ID,omitempty"`
	ParentID                        string `json:"parentID,omitempty"`
	ParentType                      string `json:"parentType,omitempty"`
	Owner                           string `json:"owner,omitempty"`
	Name                            string `json:"name,omitempty"`
	DampeningEnabled                bool   `json:"dampeningEnabled"`
	PeerAS                          int    `json:"peerAS,omitempty"`
	PeerIP                          string `json:"peerIP,omitempty"`
	Description                     string `json:"description,omitempty"`
	Session                         string `json:"session,omitempty"`
	EntityScope                     string `json:"entityScope,omitempty"`
	AssociatedExportRoutingPolicyID string `json:"associatedExportRoutingPolicyID,omitempty"`
	AssociatedImportRoutingPolicyID string `json:"associatedImportRoutingPolicyID,omitempty"`
	ExternalID                      string `json:"externalID,omitempty"`
}

// NewBGPNeighbor returns a new *BGPNeighbor
func NewBGPNeighbor() *BGPNeighbor {

	return &BGPNeighbor{}
}

// Identity returns the Identity of the object.
func (o *BGPNeighbor) Identity() bambou.Identity {

	return BGPNeighborIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *BGPNeighbor) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *BGPNeighbor) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the BGPNeighbor from the server
func (o *BGPNeighbor) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the BGPNeighbor into the server
func (o *BGPNeighbor) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the BGPNeighbor from the server
func (o *BGPNeighbor) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the BGPNeighbor
func (o *BGPNeighbor) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the BGPNeighbor
func (o *BGPNeighbor) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the BGPNeighbor
func (o *BGPNeighbor) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the BGPNeighbor
func (o *BGPNeighbor) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
