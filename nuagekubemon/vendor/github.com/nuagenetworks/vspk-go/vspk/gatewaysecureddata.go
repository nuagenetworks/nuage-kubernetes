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

// GatewaySecuredDataIdentity represents the Identity of the object
var GatewaySecuredDataIdentity = bambou.Identity{
	Name:     "gatewaysecureddata",
	Category: "gatewaysecureddatas",
}

// GatewaySecuredDatasList represents a list of GatewaySecuredDatas
type GatewaySecuredDatasList []*GatewaySecuredData

// GatewaySecuredDatasAncestor is the interface that an ancestor of a GatewaySecuredData must implement.
// An Ancestor is defined as an entity that has GatewaySecuredData as a descendant.
// An Ancestor can get a list of its child GatewaySecuredDatas, but not necessarily create one.
type GatewaySecuredDatasAncestor interface {
	GatewaySecuredDatas(*bambou.FetchingInfo) (GatewaySecuredDatasList, *bambou.Error)
}

// GatewaySecuredDatasParent is the interface that a parent of a GatewaySecuredData must implement.
// A Parent is defined as an entity that has GatewaySecuredData as a child.
// A Parent is an Ancestor which can create a GatewaySecuredData.
type GatewaySecuredDatasParent interface {
	GatewaySecuredDatasAncestor
	CreateGatewaySecuredData(*GatewaySecuredData) *bambou.Error
}

// GatewaySecuredData represents the model of a gatewaysecureddata
type GatewaySecuredData struct {
	ID                        string        `json:"ID,omitempty"`
	ParentID                  string        `json:"parentID,omitempty"`
	ParentType                string        `json:"parentType,omitempty"`
	Owner                     string        `json:"owner,omitempty"`
	LastUpdatedBy             string        `json:"lastUpdatedBy,omitempty"`
	Data                      string        `json:"data,omitempty"`
	GatewayCertSerialNumber   string        `json:"gatewayCertSerialNumber,omitempty"`
	KeyserverCertSerialNumber string        `json:"keyserverCertSerialNumber,omitempty"`
	SignedData                string        `json:"signedData,omitempty"`
	EmbeddedMetadata          []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope               string        `json:"entityScope,omitempty"`
	AssociatedEnterpriseID    string        `json:"associatedEnterpriseID,omitempty"`
	ExternalID                string        `json:"externalID,omitempty"`
}

// NewGatewaySecuredData returns a new *GatewaySecuredData
func NewGatewaySecuredData() *GatewaySecuredData {

	return &GatewaySecuredData{}
}

// Identity returns the Identity of the object.
func (o *GatewaySecuredData) Identity() bambou.Identity {

	return GatewaySecuredDataIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *GatewaySecuredData) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *GatewaySecuredData) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the GatewaySecuredData from the server
func (o *GatewaySecuredData) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the GatewaySecuredData into the server
func (o *GatewaySecuredData) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the GatewaySecuredData from the server
func (o *GatewaySecuredData) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the GatewaySecuredData
func (o *GatewaySecuredData) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the GatewaySecuredData
func (o *GatewaySecuredData) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the GatewaySecuredData
func (o *GatewaySecuredData) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the GatewaySecuredData
func (o *GatewaySecuredData) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
