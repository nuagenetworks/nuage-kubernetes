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

// PatchIdentity represents the Identity of the object
var PatchIdentity = bambou.Identity{
	Name:     "patch",
	Category: "patches",
}

// PatchsList represents a list of Patchs
type PatchsList []*Patch

// PatchsAncestor is the interface that an ancestor of a Patch must implement.
// An Ancestor is defined as an entity that has Patch as a descendant.
// An Ancestor can get a list of its child Patchs, but not necessarily create one.
type PatchsAncestor interface {
	Patchs(*bambou.FetchingInfo) (PatchsList, *bambou.Error)
}

// PatchsParent is the interface that a parent of a Patch must implement.
// A Parent is defined as an entity that has Patch as a child.
// A Parent is an Ancestor which can create a Patch.
type PatchsParent interface {
	PatchsAncestor
	CreatePatch(*Patch) *bambou.Error
}

// Patch represents the model of a patch
type Patch struct {
	ID                          string        `json:"ID,omitempty"`
	ParentID                    string        `json:"parentID,omitempty"`
	ParentType                  string        `json:"parentType,omitempty"`
	Owner                       string        `json:"owner,omitempty"`
	Name                        string        `json:"name,omitempty"`
	LastUpdatedBy               string        `json:"lastUpdatedBy,omitempty"`
	PatchBuildNumber            int           `json:"patchBuildNumber,omitempty"`
	PatchSummary                string        `json:"patchSummary,omitempty"`
	PatchTag                    string        `json:"patchTag,omitempty"`
	PatchVersion                string        `json:"patchVersion,omitempty"`
	Description                 string        `json:"description,omitempty"`
	EmbeddedMetadata            []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope                 string        `json:"entityScope,omitempty"`
	SupportsDeletion            bool          `json:"supportsDeletion"`
	SupportsNetworkAcceleration bool          `json:"supportsNetworkAcceleration"`
	ExternalID                  string        `json:"externalID,omitempty"`
}

// NewPatch returns a new *Patch
func NewPatch() *Patch {

	return &Patch{
		SupportsDeletion:            false,
		SupportsNetworkAcceleration: false,
	}
}

// Identity returns the Identity of the object.
func (o *Patch) Identity() bambou.Identity {

	return PatchIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Patch) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Patch) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Patch from the server
func (o *Patch) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Patch into the server
func (o *Patch) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Patch from the server
func (o *Patch) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the Patch
func (o *Patch) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the Patch
func (o *Patch) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Patch
func (o *Patch) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Patch
func (o *Patch) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
