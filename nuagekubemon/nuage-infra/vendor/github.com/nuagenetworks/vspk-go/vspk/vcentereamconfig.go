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

// VCenterEAMConfigIdentity represents the Identity of the object
var VCenterEAMConfigIdentity = bambou.Identity{
	Name:     "eamconfig",
	Category: "eamconfigs",
}

// VCenterEAMConfigsList represents a list of VCenterEAMConfigs
type VCenterEAMConfigsList []*VCenterEAMConfig

// VCenterEAMConfigsAncestor is the interface that an ancestor of a VCenterEAMConfig must implement.
// An Ancestor is defined as an entity that has VCenterEAMConfig as a descendant.
// An Ancestor can get a list of its child VCenterEAMConfigs, but not necessarily create one.
type VCenterEAMConfigsAncestor interface {
	VCenterEAMConfigs(*bambou.FetchingInfo) (VCenterEAMConfigsList, *bambou.Error)
}

// VCenterEAMConfigsParent is the interface that a parent of a VCenterEAMConfig must implement.
// A Parent is defined as an entity that has VCenterEAMConfig as a child.
// A Parent is an Ancestor which can create a VCenterEAMConfig.
type VCenterEAMConfigsParent interface {
	VCenterEAMConfigsAncestor
	CreateVCenterEAMConfig(*VCenterEAMConfig) *bambou.Error
}

// VCenterEAMConfig represents the model of a eamconfig
type VCenterEAMConfig struct {
	ID                  string `json:"ID,omitempty"`
	ParentID            string `json:"parentID,omitempty"`
	ParentType          string `json:"parentType,omitempty"`
	Owner               string `json:"owner,omitempty"`
	EamServerIP         string `json:"eamServerIP,omitempty"`
	EamServerPortNumber int    `json:"eamServerPortNumber,omitempty"`
	EamServerPortType   string `json:"eamServerPortType,omitempty"`
	LastUpdatedBy       string `json:"lastUpdatedBy,omitempty"`
	VibURL              string `json:"vibURL,omitempty"`
	EntityScope         string `json:"entityScope,omitempty"`
	OvfURL              string `json:"ovfURL,omitempty"`
	ExtensionKey        string `json:"extensionKey,omitempty"`
	ExternalID          string `json:"externalID,omitempty"`
}

// NewVCenterEAMConfig returns a new *VCenterEAMConfig
func NewVCenterEAMConfig() *VCenterEAMConfig {

	return &VCenterEAMConfig{}
}

// Identity returns the Identity of the object.
func (o *VCenterEAMConfig) Identity() bambou.Identity {

	return VCenterEAMConfigIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VCenterEAMConfig) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VCenterEAMConfig) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VCenterEAMConfig from the server
func (o *VCenterEAMConfig) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VCenterEAMConfig into the server
func (o *VCenterEAMConfig) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VCenterEAMConfig from the server
func (o *VCenterEAMConfig) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the VCenterEAMConfig
func (o *VCenterEAMConfig) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VCenterEAMConfig
func (o *VCenterEAMConfig) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VCenterEAMConfig
func (o *VCenterEAMConfig) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VCenterEAMConfig
func (o *VCenterEAMConfig) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
