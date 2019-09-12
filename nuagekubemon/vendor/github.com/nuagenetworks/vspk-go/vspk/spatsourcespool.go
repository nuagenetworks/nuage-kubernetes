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

// SPATSourcesPoolIdentity represents the Identity of the object
var SPATSourcesPoolIdentity = bambou.Identity{
	Name:     "spatsourcespool",
	Category: "spatsourcespools",
}

// SPATSourcesPoolsList represents a list of SPATSourcesPools
type SPATSourcesPoolsList []*SPATSourcesPool

// SPATSourcesPoolsAncestor is the interface that an ancestor of a SPATSourcesPool must implement.
// An Ancestor is defined as an entity that has SPATSourcesPool as a descendant.
// An Ancestor can get a list of its child SPATSourcesPools, but not necessarily create one.
type SPATSourcesPoolsAncestor interface {
	SPATSourcesPools(*bambou.FetchingInfo) (SPATSourcesPoolsList, *bambou.Error)
}

// SPATSourcesPoolsParent is the interface that a parent of a SPATSourcesPool must implement.
// A Parent is defined as an entity that has SPATSourcesPool as a child.
// A Parent is an Ancestor which can create a SPATSourcesPool.
type SPATSourcesPoolsParent interface {
	SPATSourcesPoolsAncestor
	CreateSPATSourcesPool(*SPATSourcesPool) *bambou.Error
}

// SPATSourcesPool represents the model of a spatsourcespool
type SPATSourcesPool struct {
	ID               string        `json:"ID,omitempty"`
	ParentID         string        `json:"parentID,omitempty"`
	ParentType       string        `json:"parentType,omitempty"`
	Owner            string        `json:"owner,omitempty"`
	Name             string        `json:"name,omitempty"`
	Family           string        `json:"family,omitempty"`
	LastUpdatedBy    string        `json:"lastUpdatedBy,omitempty"`
	AddressList      []interface{} `json:"addressList,omitempty"`
	EmbeddedMetadata []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope      string        `json:"entityScope,omitempty"`
	ExternalID       string        `json:"externalID,omitempty"`
}

// NewSPATSourcesPool returns a new *SPATSourcesPool
func NewSPATSourcesPool() *SPATSourcesPool {

	return &SPATSourcesPool{
		Name: "IPV4",
	}
}

// Identity returns the Identity of the object.
func (o *SPATSourcesPool) Identity() bambou.Identity {

	return SPATSourcesPoolIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *SPATSourcesPool) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *SPATSourcesPool) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the SPATSourcesPool from the server
func (o *SPATSourcesPool) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the SPATSourcesPool into the server
func (o *SPATSourcesPool) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the SPATSourcesPool from the server
func (o *SPATSourcesPool) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the SPATSourcesPool
func (o *SPATSourcesPool) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the SPATSourcesPool
func (o *SPATSourcesPool) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the SPATSourcesPool
func (o *SPATSourcesPool) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the SPATSourcesPool
func (o *SPATSourcesPool) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
