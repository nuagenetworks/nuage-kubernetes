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

// VNFMetadataIdentity represents the Identity of the object
var VNFMetadataIdentity = bambou.Identity{
	Name:     "vnfmetadata",
	Category: "vnfmetadatas",
}

// VNFMetadatasList represents a list of VNFMetadatas
type VNFMetadatasList []*VNFMetadata

// VNFMetadatasAncestor is the interface that an ancestor of a VNFMetadata must implement.
// An Ancestor is defined as an entity that has VNFMetadata as a descendant.
// An Ancestor can get a list of its child VNFMetadatas, but not necessarily create one.
type VNFMetadatasAncestor interface {
	VNFMetadatas(*bambou.FetchingInfo) (VNFMetadatasList, *bambou.Error)
}

// VNFMetadatasParent is the interface that a parent of a VNFMetadata must implement.
// A Parent is defined as an entity that has VNFMetadata as a child.
// A Parent is an Ancestor which can create a VNFMetadata.
type VNFMetadatasParent interface {
	VNFMetadatasAncestor
	CreateVNFMetadata(*VNFMetadata) *bambou.Error
}

// VNFMetadata represents the model of a vnfmetadata
type VNFMetadata struct {
	ID               string        `json:"ID,omitempty"`
	ParentID         string        `json:"parentID,omitempty"`
	ParentType       string        `json:"parentType,omitempty"`
	Owner            string        `json:"owner,omitempty"`
	Name             string        `json:"name,omitempty"`
	LastUpdatedBy    string        `json:"lastUpdatedBy,omitempty"`
	Description      string        `json:"description,omitempty"`
	Blob             string        `json:"blob,omitempty"`
	EmbeddedMetadata []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope      string        `json:"entityScope,omitempty"`
	AssocEntityType  string        `json:"assocEntityType,omitempty"`
	ExternalID       string        `json:"externalID,omitempty"`
}

// NewVNFMetadata returns a new *VNFMetadata
func NewVNFMetadata() *VNFMetadata {

	return &VNFMetadata{}
}

// Identity returns the Identity of the object.
func (o *VNFMetadata) Identity() bambou.Identity {

	return VNFMetadataIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VNFMetadata) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VNFMetadata) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VNFMetadata from the server
func (o *VNFMetadata) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VNFMetadata into the server
func (o *VNFMetadata) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VNFMetadata from the server
func (o *VNFMetadata) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the VNFMetadata
func (o *VNFMetadata) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VNFMetadata
func (o *VNFMetadata) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VNFMetadata
func (o *VNFMetadata) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VNFMetadata
func (o *VNFMetadata) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
