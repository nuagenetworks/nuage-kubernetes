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

// GroupKeyEncryptionProfileIdentity represents the Identity of the object
var GroupKeyEncryptionProfileIdentity = bambou.Identity{
	Name:     "groupkeyencryptionprofile",
	Category: "groupkeyencryptionprofiles",
}

// GroupKeyEncryptionProfilesList represents a list of GroupKeyEncryptionProfiles
type GroupKeyEncryptionProfilesList []*GroupKeyEncryptionProfile

// GroupKeyEncryptionProfilesAncestor is the interface of an ancestor of a GroupKeyEncryptionProfile must implement.
type GroupKeyEncryptionProfilesAncestor interface {
	GroupKeyEncryptionProfiles(*bambou.FetchingInfo) (GroupKeyEncryptionProfilesList, *bambou.Error)
	CreateGroupKeyEncryptionProfiles(*GroupKeyEncryptionProfile) *bambou.Error
}

// GroupKeyEncryptionProfile represents the model of a groupkeyencryptionprofile
type GroupKeyEncryptionProfile struct {
	ID                                   string `json:"ID,omitempty"`
	ParentID                             string `json:"parentID,omitempty"`
	ParentType                           string `json:"parentType,omitempty"`
	Owner                                string `json:"owner,omitempty"`
	SEKGenerationInterval                int    `json:"SEKGenerationInterval,omitempty"`
	SEKLifetime                          int    `json:"SEKLifetime,omitempty"`
	SEKPayloadEncryptionAlgorithm        string `json:"SEKPayloadEncryptionAlgorithm,omitempty"`
	SEKPayloadEncryptionBCAlgorithm      string `json:"SEKPayloadEncryptionBCAlgorithm,omitempty"`
	SEKPayloadEncryptionKeyLength        int    `json:"SEKPayloadEncryptionKeyLength,omitempty"`
	SEKPayloadSigningAlgorithm           string `json:"SEKPayloadSigningAlgorithm,omitempty"`
	Name                                 string `json:"name,omitempty"`
	LastUpdatedBy                        string `json:"lastUpdatedBy,omitempty"`
	SeedGenerationInterval               int    `json:"seedGenerationInterval,omitempty"`
	SeedLifetime                         int    `json:"seedLifetime,omitempty"`
	SeedPayloadAuthenticationAlgorithm   string `json:"seedPayloadAuthenticationAlgorithm,omitempty"`
	SeedPayloadAuthenticationBCAlgorithm string `json:"seedPayloadAuthenticationBCAlgorithm,omitempty"`
	SeedPayloadAuthenticationKeyLength   int    `json:"seedPayloadAuthenticationKeyLength,omitempty"`
	SeedPayloadEncryptionAlgorithm       string `json:"seedPayloadEncryptionAlgorithm,omitempty"`
	SeedPayloadEncryptionBCAlgorithm     string `json:"seedPayloadEncryptionBCAlgorithm,omitempty"`
	SeedPayloadEncryptionKeyLength       int    `json:"seedPayloadEncryptionKeyLength,omitempty"`
	SeedPayloadSigningAlgorithm          string `json:"seedPayloadSigningAlgorithm,omitempty"`
	Description                          string `json:"description,omitempty"`
	EntityScope                          string `json:"entityScope,omitempty"`
	TrafficAuthenticationAlgorithm       string `json:"trafficAuthenticationAlgorithm,omitempty"`
	TrafficEncryptionAlgorithm           string `json:"trafficEncryptionAlgorithm,omitempty"`
	TrafficEncryptionKeyLifetime         int    `json:"trafficEncryptionKeyLifetime,omitempty"`
	AssociatedEnterpriseID               string `json:"associatedEnterpriseID,omitempty"`
	ExternalID                           string `json:"externalID,omitempty"`
}

// NewGroupKeyEncryptionProfile returns a new *GroupKeyEncryptionProfile
func NewGroupKeyEncryptionProfile() *GroupKeyEncryptionProfile {

	return &GroupKeyEncryptionProfile{}
}

// Identity returns the Identity of the object.
func (o *GroupKeyEncryptionProfile) Identity() bambou.Identity {

	return GroupKeyEncryptionProfileIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *GroupKeyEncryptionProfile) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *GroupKeyEncryptionProfile) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the GroupKeyEncryptionProfile from the server
func (o *GroupKeyEncryptionProfile) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the GroupKeyEncryptionProfile into the server
func (o *GroupKeyEncryptionProfile) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the GroupKeyEncryptionProfile from the server
func (o *GroupKeyEncryptionProfile) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the GroupKeyEncryptionProfile
func (o *GroupKeyEncryptionProfile) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the GroupKeyEncryptionProfile
func (o *GroupKeyEncryptionProfile) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the GroupKeyEncryptionProfile
func (o *GroupKeyEncryptionProfile) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the GroupKeyEncryptionProfile
func (o *GroupKeyEncryptionProfile) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
