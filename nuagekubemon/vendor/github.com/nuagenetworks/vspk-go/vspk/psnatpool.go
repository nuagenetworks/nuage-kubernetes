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

// PSNATPoolIdentity represents the Identity of the object
var PSNATPoolIdentity = bambou.Identity{
	Name:     "psnatpool",
	Category: "psnatpools",
}

// PSNATPoolsList represents a list of PSNATPools
type PSNATPoolsList []*PSNATPool

// PSNATPoolsAncestor is the interface that an ancestor of a PSNATPool must implement.
// An Ancestor is defined as an entity that has PSNATPool as a descendant.
// An Ancestor can get a list of its child PSNATPools, but not necessarily create one.
type PSNATPoolsAncestor interface {
	PSNATPools(*bambou.FetchingInfo) (PSNATPoolsList, *bambou.Error)
}

// PSNATPoolsParent is the interface that a parent of a PSNATPool must implement.
// A Parent is defined as an entity that has PSNATPool as a child.
// A Parent is an Ancestor which can create a PSNATPool.
type PSNATPoolsParent interface {
	PSNATPoolsAncestor
	CreatePSNATPool(*PSNATPool) *bambou.Error
}

// PSNATPool represents the model of a psnatpool
type PSNATPool struct {
	ID               string        `json:"ID,omitempty"`
	ParentID         string        `json:"parentID,omitempty"`
	ParentType       string        `json:"parentType,omitempty"`
	Owner            string        `json:"owner,omitempty"`
	IPType           string        `json:"IPType,omitempty"`
	Name             string        `json:"name,omitempty"`
	LastUpdatedBy    string        `json:"lastUpdatedBy,omitempty"`
	Description      string        `json:"description,omitempty"`
	EmbeddedMetadata []interface{} `json:"embeddedMetadata,omitempty"`
	EndAddress       string        `json:"endAddress,omitempty"`
	EntityScope      string        `json:"entityScope,omitempty"`
	StartAddress     string        `json:"startAddress,omitempty"`
	ExternalID       string        `json:"externalID,omitempty"`
}

// NewPSNATPool returns a new *PSNATPool
func NewPSNATPool() *PSNATPool {

	return &PSNATPool{}
}

// Identity returns the Identity of the object.
func (o *PSNATPool) Identity() bambou.Identity {

	return PSNATPoolIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *PSNATPool) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *PSNATPool) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the PSNATPool from the server
func (o *PSNATPool) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the PSNATPool into the server
func (o *PSNATPool) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the PSNATPool from the server
func (o *PSNATPool) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the PSNATPool
func (o *PSNATPool) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the PSNATPool
func (o *PSNATPool) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the PSNATPool
func (o *PSNATPool) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the PSNATPool
func (o *PSNATPool) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PSPATMaps retrieves the list of child PSPATMaps of the PSNATPool
func (o *PSNATPool) PSPATMaps(info *bambou.FetchingInfo) (PSPATMapsList, *bambou.Error) {

	var list PSPATMapsList
	err := bambou.CurrentSession().FetchChildren(o, PSPATMapIdentity, &list, info)
	return list, err
}

// CreatePSPATMap creates a new child PSPATMap under the PSNATPool
func (o *PSNATPool) CreatePSPATMap(child *PSPATMap) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PTranslationMaps retrieves the list of child PTranslationMaps of the PSNATPool
func (o *PSNATPool) PTranslationMaps(info *bambou.FetchingInfo) (PTranslationMapsList, *bambou.Error) {

	var list PTranslationMapsList
	err := bambou.CurrentSession().FetchChildren(o, PTranslationMapIdentity, &list, info)
	return list, err
}

// CreatePTranslationMap creates a new child PTranslationMap under the PSNATPool
func (o *PSNATPool) CreatePTranslationMap(child *PTranslationMap) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
