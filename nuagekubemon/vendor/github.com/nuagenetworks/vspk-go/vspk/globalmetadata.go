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

// GlobalMetadataIdentity represents the Identity of the object
var GlobalMetadataIdentity = bambou.Identity{
	Name:     "globalmetadata",
	Category: "globalmetadatas",
}

// GlobalMetadatasList represents a list of GlobalMetadatas
type GlobalMetadatasList []*GlobalMetadata

// GlobalMetadatasAncestor is the interface of an ancestor of a GlobalMetadata must implement.
type GlobalMetadatasAncestor interface {
	GlobalMetadatas(*bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error)
	CreateGlobalMetadatas(*GlobalMetadata) *bambou.Error
}

// GlobalMetadata represents the model of a globalmetadata
type GlobalMetadata struct {
	ID                          string        `json:"ID,omitempty"`
	ParentID                    string        `json:"parentID,omitempty"`
	ParentType                  string        `json:"parentType,omitempty"`
	Owner                       string        `json:"owner,omitempty"`
	Name                        string        `json:"name,omitempty"`
	LastUpdatedBy               string        `json:"lastUpdatedBy,omitempty"`
	Description                 string        `json:"description,omitempty"`
	MetadataTagIDs              []interface{} `json:"metadataTagIDs,omitempty"`
	NetworkNotificationDisabled bool          `json:"networkNotificationDisabled"`
	Blob                        string        `json:"blob,omitempty"`
	GlobalMetadata              bool          `json:"globalMetadata"`
	EntityScope                 string        `json:"entityScope,omitempty"`
	ExternalID                  string        `json:"externalID,omitempty"`
}

// NewGlobalMetadata returns a new *GlobalMetadata
func NewGlobalMetadata() *GlobalMetadata {

	return &GlobalMetadata{}
}

// Identity returns the Identity of the object.
func (o *GlobalMetadata) Identity() bambou.Identity {

	return GlobalMetadataIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *GlobalMetadata) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *GlobalMetadata) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the GlobalMetadata from the server
func (o *GlobalMetadata) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the GlobalMetadata into the server
func (o *GlobalMetadata) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the GlobalMetadata from the server
func (o *GlobalMetadata) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the GlobalMetadata
func (o *GlobalMetadata) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the GlobalMetadata
func (o *GlobalMetadata) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// MetadataTags retrieves the list of child MetadataTags of the GlobalMetadata
func (o *GlobalMetadata) MetadataTags(info *bambou.FetchingInfo) (MetadataTagsList, *bambou.Error) {

	var list MetadataTagsList
	err := bambou.CurrentSession().FetchChildren(o, MetadataTagIdentity, &list, info)
	return list, err
}

// CreateMetadataTag creates a new child MetadataTag under the GlobalMetadata
func (o *GlobalMetadata) CreateMetadataTag(child *MetadataTag) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the GlobalMetadata
func (o *GlobalMetadata) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the GlobalMetadata
func (o *GlobalMetadata) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
