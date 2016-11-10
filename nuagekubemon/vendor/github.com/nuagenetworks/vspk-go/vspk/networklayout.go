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

// NetworkLayoutIdentity represents the Identity of the object
var NetworkLayoutIdentity = bambou.Identity{
	Name:     "networklayout",
	Category: "networklayout",
}

// NetworkLayoutsList represents a list of NetworkLayouts
type NetworkLayoutsList []*NetworkLayout

// NetworkLayoutsAncestor is the interface of an ancestor of a NetworkLayout must implement.
type NetworkLayoutsAncestor interface {
	NetworkLayouts(*bambou.FetchingInfo) (NetworkLayoutsList, *bambou.Error)
	CreateNetworkLayouts(*NetworkLayout) *bambou.Error
}

// NetworkLayout represents the model of a networklayout
type NetworkLayout struct {
	ID                  string `json:"ID,omitempty"`
	ParentID            string `json:"parentID,omitempty"`
	ParentType          string `json:"parentType,omitempty"`
	Owner               string `json:"owner,omitempty"`
	LastUpdatedBy       string `json:"lastUpdatedBy,omitempty"`
	ServiceType         string `json:"serviceType,omitempty"`
	EntityScope         string `json:"entityScope,omitempty"`
	RouteReflectorIP    string `json:"routeReflectorIP,omitempty"`
	AutonomousSystemNum int    `json:"autonomousSystemNum,omitempty"`
	ExternalID          string `json:"externalID,omitempty"`
}

// NewNetworkLayout returns a new *NetworkLayout
func NewNetworkLayout() *NetworkLayout {

	return &NetworkLayout{}
}

// Identity returns the Identity of the object.
func (o *NetworkLayout) Identity() bambou.Identity {

	return NetworkLayoutIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *NetworkLayout) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *NetworkLayout) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the NetworkLayout from the server
func (o *NetworkLayout) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the NetworkLayout into the server
func (o *NetworkLayout) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the NetworkLayout from the server
func (o *NetworkLayout) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the NetworkLayout
func (o *NetworkLayout) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the NetworkLayout
func (o *NetworkLayout) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the NetworkLayout
func (o *NetworkLayout) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the NetworkLayout
func (o *NetworkLayout) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
