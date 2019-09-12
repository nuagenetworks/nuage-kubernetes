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

// IngressProfileIdentity represents the Identity of the object
var IngressProfileIdentity = bambou.Identity{
	Name:     "ingressprofile",
	Category: "ingressprofiles",
}

// IngressProfilesList represents a list of IngressProfiles
type IngressProfilesList []*IngressProfile

// IngressProfilesAncestor is the interface that an ancestor of a IngressProfile must implement.
// An Ancestor is defined as an entity that has IngressProfile as a descendant.
// An Ancestor can get a list of its child IngressProfiles, but not necessarily create one.
type IngressProfilesAncestor interface {
	IngressProfiles(*bambou.FetchingInfo) (IngressProfilesList, *bambou.Error)
}

// IngressProfilesParent is the interface that a parent of a IngressProfile must implement.
// A Parent is defined as an entity that has IngressProfile as a child.
// A Parent is an Ancestor which can create a IngressProfile.
type IngressProfilesParent interface {
	IngressProfilesAncestor
	CreateIngressProfile(*IngressProfile) *bambou.Error
}

// IngressProfile represents the model of a ingressprofile
type IngressProfile struct {
	ID                                 string        `json:"ID,omitempty"`
	ParentID                           string        `json:"parentID,omitempty"`
	ParentType                         string        `json:"parentType,omitempty"`
	Owner                              string        `json:"owner,omitempty"`
	Name                               string        `json:"name,omitempty"`
	LastUpdatedBy                      string        `json:"lastUpdatedBy,omitempty"`
	Description                        string        `json:"description,omitempty"`
	EmbeddedMetadata                   []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope                        string        `json:"entityScope,omitempty"`
	AssocEntityType                    string        `json:"assocEntityType,omitempty"`
	AssociatedIPFilterProfileID        string        `json:"associatedIPFilterProfileID,omitempty"`
	AssociatedIPFilterProfileName      string        `json:"associatedIPFilterProfileName,omitempty"`
	AssociatedIPv6FilterProfileID      string        `json:"associatedIPv6FilterProfileID,omitempty"`
	AssociatedIPv6FilterProfileName    string        `json:"associatedIPv6FilterProfileName,omitempty"`
	AssociatedMACFilterProfileID       string        `json:"associatedMACFilterProfileID,omitempty"`
	AssociatedMACFilterProfileName     string        `json:"associatedMACFilterProfileName,omitempty"`
	AssociatedSAPIngressQoSProfileID   string        `json:"associatedSAPIngressQoSProfileID,omitempty"`
	AssociatedSAPIngressQoSProfileName string        `json:"associatedSAPIngressQoSProfileName,omitempty"`
	ExternalID                         string        `json:"externalID,omitempty"`
}

// NewIngressProfile returns a new *IngressProfile
func NewIngressProfile() *IngressProfile {

	return &IngressProfile{}
}

// Identity returns the Identity of the object.
func (o *IngressProfile) Identity() bambou.Identity {

	return IngressProfileIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *IngressProfile) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *IngressProfile) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the IngressProfile from the server
func (o *IngressProfile) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the IngressProfile into the server
func (o *IngressProfile) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the IngressProfile from the server
func (o *IngressProfile) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// DeploymentFailures retrieves the list of child DeploymentFailures of the IngressProfile
func (o *IngressProfile) DeploymentFailures(info *bambou.FetchingInfo) (DeploymentFailuresList, *bambou.Error) {

	var list DeploymentFailuresList
	err := bambou.CurrentSession().FetchChildren(o, DeploymentFailureIdentity, &list, info)
	return list, err
}

// CreateDeploymentFailure creates a new child DeploymentFailure under the IngressProfile
func (o *IngressProfile) CreateDeploymentFailure(child *DeploymentFailure) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the IngressProfile
func (o *IngressProfile) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the IngressProfile
func (o *IngressProfile) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the IngressProfile
func (o *IngressProfile) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the IngressProfile
func (o *IngressProfile) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VPorts retrieves the list of child VPorts of the IngressProfile
func (o *IngressProfile) VPorts(info *bambou.FetchingInfo) (VPortsList, *bambou.Error) {

	var list VPortsList
	err := bambou.CurrentSession().FetchChildren(o, VPortIdentity, &list, info)
	return list, err
}
