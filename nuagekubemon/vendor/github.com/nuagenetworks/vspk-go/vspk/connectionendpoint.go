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

// ConnectionendpointIdentity represents the Identity of the object
var ConnectionendpointIdentity = bambou.Identity{
	Name:     "connectionendpoint",
	Category: "connectionendpoints",
}

// ConnectionendpointsList represents a list of Connectionendpoints
type ConnectionendpointsList []*Connectionendpoint

// ConnectionendpointsAncestor is the interface that an ancestor of a Connectionendpoint must implement.
// An Ancestor is defined as an entity that has Connectionendpoint as a descendant.
// An Ancestor can get a list of its child Connectionendpoints, but not necessarily create one.
type ConnectionendpointsAncestor interface {
	Connectionendpoints(*bambou.FetchingInfo) (ConnectionendpointsList, *bambou.Error)
}

// ConnectionendpointsParent is the interface that a parent of a Connectionendpoint must implement.
// A Parent is defined as an entity that has Connectionendpoint as a child.
// A Parent is an Ancestor which can create a Connectionendpoint.
type ConnectionendpointsParent interface {
	ConnectionendpointsAncestor
	CreateConnectionendpoint(*Connectionendpoint) *bambou.Error
}

// Connectionendpoint represents the model of a connectionendpoint
type Connectionendpoint struct {
	ID               string        `json:"ID,omitempty"`
	ParentID         string        `json:"parentID,omitempty"`
	ParentType       string        `json:"parentType,omitempty"`
	Owner            string        `json:"owner,omitempty"`
	IPAddress        string        `json:"IPAddress,omitempty"`
	IPType           string        `json:"IPType,omitempty"`
	IPv6Address      string        `json:"IPv6Address,omitempty"`
	Name             string        `json:"name,omitempty"`
	LastUpdatedBy    string        `json:"lastUpdatedBy,omitempty"`
	Description      string        `json:"description,omitempty"`
	EmbeddedMetadata []interface{} `json:"embeddedMetadata,omitempty"`
	EndPointType     string        `json:"endPointType,omitempty"`
	EntityScope      string        `json:"entityScope,omitempty"`
	ExternalID       string        `json:"externalID,omitempty"`
}

// NewConnectionendpoint returns a new *Connectionendpoint
func NewConnectionendpoint() *Connectionendpoint {

	return &Connectionendpoint{
		IPType:       "IPV4",
		EndPointType: "SOURCE",
	}
}

// Identity returns the Identity of the object.
func (o *Connectionendpoint) Identity() bambou.Identity {

	return ConnectionendpointIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Connectionendpoint) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Connectionendpoint) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Connectionendpoint from the server
func (o *Connectionendpoint) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Connectionendpoint into the server
func (o *Connectionendpoint) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Connectionendpoint from the server
func (o *Connectionendpoint) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the Connectionendpoint
func (o *Connectionendpoint) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the Connectionendpoint
func (o *Connectionendpoint) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Connectionendpoint
func (o *Connectionendpoint) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Connectionendpoint
func (o *Connectionendpoint) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
