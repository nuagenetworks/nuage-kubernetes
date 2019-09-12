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

// NSGatewaysCountIdentity represents the Identity of the object
var NSGatewaysCountIdentity = bambou.Identity{
	Name:     "nsgatewayscount",
	Category: "nsgatewayscounts",
}

// NSGatewaysCountsList represents a list of NSGatewaysCounts
type NSGatewaysCountsList []*NSGatewaysCount

// NSGatewaysCountsAncestor is the interface that an ancestor of a NSGatewaysCount must implement.
// An Ancestor is defined as an entity that has NSGatewaysCount as a descendant.
// An Ancestor can get a list of its child NSGatewaysCounts, but not necessarily create one.
type NSGatewaysCountsAncestor interface {
	NSGatewaysCounts(*bambou.FetchingInfo) (NSGatewaysCountsList, *bambou.Error)
}

// NSGatewaysCountsParent is the interface that a parent of a NSGatewaysCount must implement.
// A Parent is defined as an entity that has NSGatewaysCount as a child.
// A Parent is an Ancestor which can create a NSGatewaysCount.
type NSGatewaysCountsParent interface {
	NSGatewaysCountsAncestor
	CreateNSGatewaysCount(*NSGatewaysCount) *bambou.Error
}

// NSGatewaysCount represents the model of a nsgatewayscount
type NSGatewaysCount struct {
	ID               string        `json:"ID,omitempty"`
	ParentID         string        `json:"parentID,omitempty"`
	ParentType       string        `json:"parentType,omitempty"`
	Owner            string        `json:"owner,omitempty"`
	ActiveNSGCount   int           `json:"activeNSGCount,omitempty"`
	AlarmedNSGCount  interface{}   `json:"alarmedNSGCount,omitempty"`
	EmbeddedMetadata []interface{} `json:"embeddedMetadata,omitempty"`
	InactiveNSGCount int           `json:"inactiveNSGCount,omitempty"`
	EntityScope      string        `json:"entityScope,omitempty"`
	ExternalID       string        `json:"externalID,omitempty"`
}

// NewNSGatewaysCount returns a new *NSGatewaysCount
func NewNSGatewaysCount() *NSGatewaysCount {

	return &NSGatewaysCount{}
}

// Identity returns the Identity of the object.
func (o *NSGatewaysCount) Identity() bambou.Identity {

	return NSGatewaysCountIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *NSGatewaysCount) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *NSGatewaysCount) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the NSGatewaysCount from the server
func (o *NSGatewaysCount) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the NSGatewaysCount into the server
func (o *NSGatewaysCount) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the NSGatewaysCount from the server
func (o *NSGatewaysCount) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the NSGatewaysCount
func (o *NSGatewaysCount) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the NSGatewaysCount
func (o *NSGatewaysCount) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the NSGatewaysCount
func (o *NSGatewaysCount) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the NSGatewaysCount
func (o *NSGatewaysCount) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
