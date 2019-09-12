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

// ZFBRequestIdentity represents the Identity of the object
var ZFBRequestIdentity = bambou.Identity{
	Name:     "zfbrequest",
	Category: "zfbrequests",
}

// ZFBRequestsList represents a list of ZFBRequests
type ZFBRequestsList []*ZFBRequest

// ZFBRequestsAncestor is the interface that an ancestor of a ZFBRequest must implement.
// An Ancestor is defined as an entity that has ZFBRequest as a descendant.
// An Ancestor can get a list of its child ZFBRequests, but not necessarily create one.
type ZFBRequestsAncestor interface {
	ZFBRequests(*bambou.FetchingInfo) (ZFBRequestsList, *bambou.Error)
}

// ZFBRequestsParent is the interface that a parent of a ZFBRequest must implement.
// A Parent is defined as an entity that has ZFBRequest as a child.
// A Parent is an Ancestor which can create a ZFBRequest.
type ZFBRequestsParent interface {
	ZFBRequestsAncestor
	CreateZFBRequest(*ZFBRequest) *bambou.Error
}

// ZFBRequest represents the model of a zfbrequest
type ZFBRequest struct {
	ID                       string        `json:"ID,omitempty"`
	ParentID                 string        `json:"parentID,omitempty"`
	ParentType               string        `json:"parentType,omitempty"`
	Owner                    string        `json:"owner,omitempty"`
	MACAddress               string        `json:"MACAddress,omitempty"`
	ZFBApprovalStatus        string        `json:"ZFBApprovalStatus,omitempty"`
	ZFBBootstrapEnabled      bool          `json:"ZFBBootstrapEnabled"`
	ZFBInfo                  string        `json:"ZFBInfo,omitempty"`
	ZFBRequestRetryTimer     int           `json:"ZFBRequestRetryTimer,omitempty"`
	SKU                      string        `json:"SKU,omitempty"`
	IPAddress                string        `json:"IPAddress,omitempty"`
	CPUType                  string        `json:"CPUType,omitempty"`
	NSGVersion               string        `json:"NSGVersion,omitempty"`
	UUID                     string        `json:"UUID,omitempty"`
	Family                   string        `json:"family,omitempty"`
	LastConnectedTime        float64       `json:"lastConnectedTime,omitempty"`
	LastUpdatedBy            string        `json:"lastUpdatedBy,omitempty"`
	RegistrationURL          string        `json:"registrationURL,omitempty"`
	SerialNumber             string        `json:"serialNumber,omitempty"`
	EmbeddedMetadata         []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope              string        `json:"entityScope,omitempty"`
	Hostname                 string        `json:"hostname,omitempty"`
	AssociatedEnterpriseID   string        `json:"associatedEnterpriseID,omitempty"`
	AssociatedEnterpriseName string        `json:"associatedEnterpriseName,omitempty"`
	AssociatedEntityType     string        `json:"associatedEntityType,omitempty"`
	AssociatedGatewayID      string        `json:"associatedGatewayID,omitempty"`
	AssociatedGatewayName    string        `json:"associatedGatewayName,omitempty"`
	AssociatedNSGatewayID    string        `json:"associatedNSGatewayID,omitempty"`
	AssociatedNSGatewayName  string        `json:"associatedNSGatewayName,omitempty"`
	StatusString             string        `json:"statusString,omitempty"`
	ExternalID               string        `json:"externalID,omitempty"`
}

// NewZFBRequest returns a new *ZFBRequest
func NewZFBRequest() *ZFBRequest {

	return &ZFBRequest{}
}

// Identity returns the Identity of the object.
func (o *ZFBRequest) Identity() bambou.Identity {

	return ZFBRequestIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *ZFBRequest) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *ZFBRequest) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the ZFBRequest from the server
func (o *ZFBRequest) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the ZFBRequest into the server
func (o *ZFBRequest) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the ZFBRequest from the server
func (o *ZFBRequest) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the ZFBRequest
func (o *ZFBRequest) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the ZFBRequest
func (o *ZFBRequest) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the ZFBRequest
func (o *ZFBRequest) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the ZFBRequest
func (o *ZFBRequest) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CreateJob creates a new child Job under the ZFBRequest
func (o *ZFBRequest) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
