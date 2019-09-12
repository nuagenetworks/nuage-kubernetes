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

// VPNConnectionIdentity represents the Identity of the object
var VPNConnectionIdentity = bambou.Identity{
	Name:     "vpnconnection",
	Category: "vpnconnections",
}

// VPNConnectionsList represents a list of VPNConnections
type VPNConnectionsList []*VPNConnection

// VPNConnectionsAncestor is the interface that an ancestor of a VPNConnection must implement.
// An Ancestor is defined as an entity that has VPNConnection as a descendant.
// An Ancestor can get a list of its child VPNConnections, but not necessarily create one.
type VPNConnectionsAncestor interface {
	VPNConnections(*bambou.FetchingInfo) (VPNConnectionsList, *bambou.Error)
}

// VPNConnectionsParent is the interface that a parent of a VPNConnection must implement.
// A Parent is defined as an entity that has VPNConnection as a child.
// A Parent is an Ancestor which can create a VPNConnection.
type VPNConnectionsParent interface {
	VPNConnectionsAncestor
	CreateVPNConnection(*VPNConnection) *bambou.Error
}

// VPNConnection represents the model of a vpnconnection
type VPNConnection struct {
	ID                     string        `json:"ID,omitempty"`
	ParentID               string        `json:"parentID,omitempty"`
	ParentType             string        `json:"parentType,omitempty"`
	Owner                  string        `json:"owner,omitempty"`
	Name                   string        `json:"name,omitempty"`
	LastUpdatedBy          string        `json:"lastUpdatedBy,omitempty"`
	Description            string        `json:"description,omitempty"`
	EmbeddedMetadata       []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope            string        `json:"entityScope,omitempty"`
	AssociatedWANServiceID string        `json:"associatedWANServiceID,omitempty"`
	ExternalID             string        `json:"externalID,omitempty"`
}

// NewVPNConnection returns a new *VPNConnection
func NewVPNConnection() *VPNConnection {

	return &VPNConnection{}
}

// Identity returns the Identity of the object.
func (o *VPNConnection) Identity() bambou.Identity {

	return VPNConnectionIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VPNConnection) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VPNConnection) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VPNConnection from the server
func (o *VPNConnection) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VPNConnection into the server
func (o *VPNConnection) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VPNConnection from the server
func (o *VPNConnection) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the VPNConnection
func (o *VPNConnection) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VPNConnection
func (o *VPNConnection) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VPNConnection
func (o *VPNConnection) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VPNConnection
func (o *VPNConnection) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
