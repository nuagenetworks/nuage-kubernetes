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

// BootstrapIdentity represents the Identity of the object
var BootstrapIdentity = bambou.Identity{
	Name:     "bootstrap",
	Category: "bootstraps",
}

// BootstrapsList represents a list of Bootstraps
type BootstrapsList []*Bootstrap

// BootstrapsAncestor is the interface that an ancestor of a Bootstrap must implement.
// An Ancestor is defined as an entity that has Bootstrap as a descendant.
// An Ancestor can get a list of its child Bootstraps, but not necessarily create one.
type BootstrapsAncestor interface {
	Bootstraps(*bambou.FetchingInfo) (BootstrapsList, *bambou.Error)
}

// BootstrapsParent is the interface that a parent of a Bootstrap must implement.
// A Parent is defined as an entity that has Bootstrap as a child.
// A Parent is an Ancestor which can create a Bootstrap.
type BootstrapsParent interface {
	BootstrapsAncestor
	CreateBootstrap(*Bootstrap) *bambou.Error
}

// Bootstrap represents the model of a bootstrap
type Bootstrap struct {
	ID                string `json:"ID,omitempty"`
	ParentID          string `json:"parentID,omitempty"`
	ParentType        string `json:"parentType,omitempty"`
	Owner             string `json:"owner,omitempty"`
	ZFBInfo           string `json:"ZFBInfo,omitempty"`
	ZFBMatchAttribute string `json:"ZFBMatchAttribute,omitempty"`
	ZFBMatchValue     string `json:"ZFBMatchValue,omitempty"`
	LastUpdatedBy     string `json:"lastUpdatedBy,omitempty"`
	InstallerID       string `json:"installerID,omitempty"`
	EntityScope       string `json:"entityScope,omitempty"`
	Status            string `json:"status,omitempty"`
	ExternalID        string `json:"externalID,omitempty"`
}

// NewBootstrap returns a new *Bootstrap
func NewBootstrap() *Bootstrap {

	return &Bootstrap{
		Status: "INACTIVE",
	}
}

// Identity returns the Identity of the object.
func (o *Bootstrap) Identity() bambou.Identity {

	return BootstrapIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Bootstrap) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Bootstrap) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Bootstrap from the server
func (o *Bootstrap) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Bootstrap into the server
func (o *Bootstrap) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Bootstrap from the server
func (o *Bootstrap) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the Bootstrap
func (o *Bootstrap) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the Bootstrap
func (o *Bootstrap) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Bootstrap
func (o *Bootstrap) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Bootstrap
func (o *Bootstrap) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
