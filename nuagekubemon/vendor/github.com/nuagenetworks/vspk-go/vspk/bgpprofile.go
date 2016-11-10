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

// BGPProfileIdentity represents the Identity of the object
var BGPProfileIdentity = bambou.Identity{
	Name:     "bgpprofile",
	Category: "bgpprofiles",
}

// BGPProfilesList represents a list of BGPProfiles
type BGPProfilesList []*BGPProfile

// BGPProfilesAncestor is the interface of an ancestor of a BGPProfile must implement.
type BGPProfilesAncestor interface {
	BGPProfiles(*bambou.FetchingInfo) (BGPProfilesList, *bambou.Error)
	CreateBGPProfiles(*BGPProfile) *bambou.Error
}

// BGPProfile represents the model of a bgpprofile
type BGPProfile struct {
	ID                              string `json:"ID,omitempty"`
	ParentID                        string `json:"parentID,omitempty"`
	ParentType                      string `json:"parentType,omitempty"`
	Owner                           string `json:"owner,omitempty"`
	Name                            string `json:"name,omitempty"`
	DampeningHalfLife               int    `json:"dampeningHalfLife,omitempty"`
	DampeningMaxSuppress            int    `json:"dampeningMaxSuppress,omitempty"`
	DampeningName                   string `json:"dampeningName,omitempty"`
	DampeningReuse                  int    `json:"dampeningReuse,omitempty"`
	DampeningSuppress               int    `json:"dampeningSuppress,omitempty"`
	Description                     string `json:"description,omitempty"`
	EntityScope                     string `json:"entityScope,omitempty"`
	AssociatedExportRoutingPolicyID string `json:"associatedExportRoutingPolicyID,omitempty"`
	AssociatedImportRoutingPolicyID string `json:"associatedImportRoutingPolicyID,omitempty"`
	ExternalID                      string `json:"externalID,omitempty"`
}

// NewBGPProfile returns a new *BGPProfile
func NewBGPProfile() *BGPProfile {

	return &BGPProfile{}
}

// Identity returns the Identity of the object.
func (o *BGPProfile) Identity() bambou.Identity {

	return BGPProfileIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *BGPProfile) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *BGPProfile) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the BGPProfile from the server
func (o *BGPProfile) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the BGPProfile into the server
func (o *BGPProfile) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the BGPProfile from the server
func (o *BGPProfile) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the BGPProfile
func (o *BGPProfile) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the BGPProfile
func (o *BGPProfile) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the BGPProfile
func (o *BGPProfile) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the BGPProfile
func (o *BGPProfile) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
