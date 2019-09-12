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

// ShuntLinkIdentity represents the Identity of the object
var ShuntLinkIdentity = bambou.Identity{
	Name:     "shuntlink",
	Category: "shuntlinks",
}

// ShuntLinksList represents a list of ShuntLinks
type ShuntLinksList []*ShuntLink

// ShuntLinksAncestor is the interface that an ancestor of a ShuntLink must implement.
// An Ancestor is defined as an entity that has ShuntLink as a descendant.
// An Ancestor can get a list of its child ShuntLinks, but not necessarily create one.
type ShuntLinksAncestor interface {
	ShuntLinks(*bambou.FetchingInfo) (ShuntLinksList, *bambou.Error)
}

// ShuntLinksParent is the interface that a parent of a ShuntLink must implement.
// A Parent is defined as an entity that has ShuntLink as a child.
// A Parent is an Ancestor which can create a ShuntLink.
type ShuntLinksParent interface {
	ShuntLinksAncestor
	CreateShuntLink(*ShuntLink) *bambou.Error
}

// ShuntLink represents the model of a shuntlink
type ShuntLink struct {
	ID               string        `json:"ID,omitempty"`
	ParentID         string        `json:"parentID,omitempty"`
	ParentType       string        `json:"parentType,omitempty"`
	Owner            string        `json:"owner,omitempty"`
	VLANPeer1ID      string        `json:"VLANPeer1ID,omitempty"`
	VLANPeer2ID      string        `json:"VLANPeer2ID,omitempty"`
	Name             string        `json:"name,omitempty"`
	LastUpdatedBy    string        `json:"lastUpdatedBy,omitempty"`
	GatewayPeer1ID   string        `json:"gatewayPeer1ID,omitempty"`
	GatewayPeer2ID   string        `json:"gatewayPeer2ID,omitempty"`
	PermittedAction  string        `json:"permittedAction,omitempty"`
	Description      string        `json:"description,omitempty"`
	EmbeddedMetadata []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope      string        `json:"entityScope,omitempty"`
	ExternalID       string        `json:"externalID,omitempty"`
}

// NewShuntLink returns a new *ShuntLink
func NewShuntLink() *ShuntLink {

	return &ShuntLink{}
}

// Identity returns the Identity of the object.
func (o *ShuntLink) Identity() bambou.Identity {

	return ShuntLinkIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *ShuntLink) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *ShuntLink) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the ShuntLink from the server
func (o *ShuntLink) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the ShuntLink into the server
func (o *ShuntLink) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the ShuntLink from the server
func (o *ShuntLink) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the ShuntLink
func (o *ShuntLink) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the ShuntLink
func (o *ShuntLink) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Alarms retrieves the list of child Alarms of the ShuntLink
func (o *ShuntLink) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the ShuntLink
func (o *ShuntLink) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the ShuntLink
func (o *ShuntLink) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
