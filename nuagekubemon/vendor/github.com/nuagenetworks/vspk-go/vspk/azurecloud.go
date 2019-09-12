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

// AzureCloudIdentity represents the Identity of the object
var AzureCloudIdentity = bambou.Identity{
	Name:     "azurecloud",
	Category: "azureclouds",
}

// AzureCloudsList represents a list of AzureClouds
type AzureCloudsList []*AzureCloud

// AzureCloudsAncestor is the interface that an ancestor of a AzureCloud must implement.
// An Ancestor is defined as an entity that has AzureCloud as a descendant.
// An Ancestor can get a list of its child AzureClouds, but not necessarily create one.
type AzureCloudsAncestor interface {
	AzureClouds(*bambou.FetchingInfo) (AzureCloudsList, *bambou.Error)
}

// AzureCloudsParent is the interface that a parent of a AzureCloud must implement.
// A Parent is defined as an entity that has AzureCloud as a child.
// A Parent is an Ancestor which can create a AzureCloud.
type AzureCloudsParent interface {
	AzureCloudsAncestor
	CreateAzureCloud(*AzureCloud) *bambou.Error
}

// AzureCloud represents the model of a azurecloud
type AzureCloud struct {
	ID                               string        `json:"ID,omitempty"`
	ParentID                         string        `json:"parentID,omitempty"`
	ParentType                       string        `json:"parentType,omitempty"`
	Owner                            string        `json:"owner,omitempty"`
	Name                             string        `json:"name,omitempty"`
	LastUpdatedBy                    string        `json:"lastUpdatedBy,omitempty"`
	TenantID                         string        `json:"tenantID,omitempty"`
	ClientID                         string        `json:"clientID,omitempty"`
	ClientSecret                     string        `json:"clientSecret,omitempty"`
	EmbeddedMetadata                 []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope                      string        `json:"entityScope,omitempty"`
	AssociatedIKEEncryptionProfileID string        `json:"associatedIKEEncryptionProfileID,omitempty"`
	AssociatedIKEPSKID               string        `json:"associatedIKEPSKID,omitempty"`
	SubscriptionID                   string        `json:"subscriptionID,omitempty"`
	ExternalID                       string        `json:"externalID,omitempty"`
}

// NewAzureCloud returns a new *AzureCloud
func NewAzureCloud() *AzureCloud {

	return &AzureCloud{}
}

// Identity returns the Identity of the object.
func (o *AzureCloud) Identity() bambou.Identity {

	return AzureCloudIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *AzureCloud) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *AzureCloud) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the AzureCloud from the server
func (o *AzureCloud) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the AzureCloud into the server
func (o *AzureCloud) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the AzureCloud from the server
func (o *AzureCloud) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the AzureCloud
func (o *AzureCloud) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the AzureCloud
func (o *AzureCloud) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IKEGatewayProfiles retrieves the list of child IKEGatewayProfiles of the AzureCloud
func (o *AzureCloud) IKEGatewayProfiles(info *bambou.FetchingInfo) (IKEGatewayProfilesList, *bambou.Error) {

	var list IKEGatewayProfilesList
	err := bambou.CurrentSession().FetchChildren(o, IKEGatewayProfileIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the AzureCloud
func (o *AzureCloud) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the AzureCloud
func (o *AzureCloud) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Jobs retrieves the list of child Jobs of the AzureCloud
func (o *AzureCloud) Jobs(info *bambou.FetchingInfo) (JobsList, *bambou.Error) {

	var list JobsList
	err := bambou.CurrentSession().FetchChildren(o, JobIdentity, &list, info)
	return list, err
}

// CreateJob creates a new child Job under the AzureCloud
func (o *AzureCloud) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
