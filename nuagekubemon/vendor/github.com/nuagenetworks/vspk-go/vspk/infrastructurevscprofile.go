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

// InfrastructureVscProfileIdentity represents the Identity of the object
var InfrastructureVscProfileIdentity = bambou.Identity{
	Name:     "infrastructurevscprofile",
	Category: "infrastructurevscprofiles",
}

// InfrastructureVscProfilesList represents a list of InfrastructureVscProfiles
type InfrastructureVscProfilesList []*InfrastructureVscProfile

// InfrastructureVscProfilesAncestor is the interface that an ancestor of a InfrastructureVscProfile must implement.
// An Ancestor is defined as an entity that has InfrastructureVscProfile as a descendant.
// An Ancestor can get a list of its child InfrastructureVscProfiles, but not necessarily create one.
type InfrastructureVscProfilesAncestor interface {
	InfrastructureVscProfiles(*bambou.FetchingInfo) (InfrastructureVscProfilesList, *bambou.Error)
}

// InfrastructureVscProfilesParent is the interface that a parent of a InfrastructureVscProfile must implement.
// A Parent is defined as an entity that has InfrastructureVscProfile as a child.
// A Parent is an Ancestor which can create a InfrastructureVscProfile.
type InfrastructureVscProfilesParent interface {
	InfrastructureVscProfilesAncestor
	CreateInfrastructureVscProfile(*InfrastructureVscProfile) *bambou.Error
}

// InfrastructureVscProfile represents the model of a infrastructurevscprofile
type InfrastructureVscProfile struct {
	ID               string `json:"ID,omitempty"`
	ParentID         string `json:"parentID,omitempty"`
	ParentType       string `json:"parentType,omitempty"`
	Owner            string `json:"owner,omitempty"`
	Name             string `json:"name,omitempty"`
	LastUpdatedBy    string `json:"lastUpdatedBy,omitempty"`
	SecondController string `json:"secondController,omitempty"`
	Description      string `json:"description,omitempty"`
	FirstController  string `json:"firstController,omitempty"`
	EnterpriseID     string `json:"enterpriseID,omitempty"`
	EntityScope      string `json:"entityScope,omitempty"`
	ProbeInterval    int    `json:"probeInterval,omitempty"`
	ExternalID       string `json:"externalID,omitempty"`
}

// NewInfrastructureVscProfile returns a new *InfrastructureVscProfile
func NewInfrastructureVscProfile() *InfrastructureVscProfile {

	return &InfrastructureVscProfile{}
}

// Identity returns the Identity of the object.
func (o *InfrastructureVscProfile) Identity() bambou.Identity {

	return InfrastructureVscProfileIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *InfrastructureVscProfile) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *InfrastructureVscProfile) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the InfrastructureVscProfile from the server
func (o *InfrastructureVscProfile) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the InfrastructureVscProfile into the server
func (o *InfrastructureVscProfile) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the InfrastructureVscProfile from the server
func (o *InfrastructureVscProfile) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the InfrastructureVscProfile
func (o *InfrastructureVscProfile) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the InfrastructureVscProfile
func (o *InfrastructureVscProfile) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the InfrastructureVscProfile
func (o *InfrastructureVscProfile) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the InfrastructureVscProfile
func (o *InfrastructureVscProfile) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
