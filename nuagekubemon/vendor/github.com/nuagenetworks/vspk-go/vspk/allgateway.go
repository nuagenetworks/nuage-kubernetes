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

// AllGatewayIdentity represents the Identity of the object
var AllGatewayIdentity = bambou.Identity{
	Name:     "allgateway",
	Category: "allgateways",
}

// AllGatewaysList represents a list of AllGateways
type AllGatewaysList []*AllGateway

// AllGatewaysAncestor is the interface that an ancestor of a AllGateway must implement.
// An Ancestor is defined as an entity that has AllGateway as a descendant.
// An Ancestor can get a list of its child AllGateways, but not necessarily create one.
type AllGatewaysAncestor interface {
	AllGateways(*bambou.FetchingInfo) (AllGatewaysList, *bambou.Error)
}

// AllGatewaysParent is the interface that a parent of a AllGateway must implement.
// A Parent is defined as an entity that has AllGateway as a child.
// A Parent is an Ancestor which can create a AllGateway.
type AllGatewaysParent interface {
	AllGatewaysAncestor
	CreateAllGateway(*AllGateway) *bambou.Error
}

// AllGateway represents the model of a allgateway
type AllGateway struct {
	ID                                 string        `json:"ID,omitempty"`
	ParentID                           string        `json:"parentID,omitempty"`
	ParentType                         string        `json:"parentType,omitempty"`
	Owner                              string        `json:"owner,omitempty"`
	MACAddress                         string        `json:"MACAddress,omitempty"`
	ZFBMatchAttribute                  string        `json:"ZFBMatchAttribute,omitempty"`
	ZFBMatchValue                      string        `json:"ZFBMatchValue,omitempty"`
	BIOSReleaseDate                    string        `json:"BIOSReleaseDate,omitempty"`
	BIOSVersion                        string        `json:"BIOSVersion,omitempty"`
	CPUType                            string        `json:"CPUType,omitempty"`
	UUID                               string        `json:"UUID,omitempty"`
	Name                               string        `json:"name,omitempty"`
	Family                             string        `json:"family,omitempty"`
	ManagementID                       string        `json:"managementID,omitempty"`
	LastUpdatedBy                      string        `json:"lastUpdatedBy,omitempty"`
	DatapathID                         string        `json:"datapathID,omitempty"`
	Patches                            string        `json:"patches,omitempty"`
	GatewayConnected                   bool          `json:"gatewayConnected"`
	GatewayVersion                     string        `json:"gatewayVersion,omitempty"`
	RedundancyGroupID                  string        `json:"redundancyGroupID,omitempty"`
	Peer                               string        `json:"peer,omitempty"`
	TemplateID                         string        `json:"templateID,omitempty"`
	Pending                            bool          `json:"pending"`
	SerialNumber                       string        `json:"serialNumber,omitempty"`
	PermittedAction                    string        `json:"permittedAction,omitempty"`
	Personality                        string        `json:"personality,omitempty"`
	Description                        string        `json:"description,omitempty"`
	Libraries                          string        `json:"libraries,omitempty"`
	EmbeddedMetadata                   []interface{} `json:"embeddedMetadata,omitempty"`
	EnterpriseID                       string        `json:"enterpriseID,omitempty"`
	EntityScope                        string        `json:"entityScope,omitempty"`
	LocationID                         string        `json:"locationID,omitempty"`
	BootstrapID                        string        `json:"bootstrapID,omitempty"`
	BootstrapStatus                    string        `json:"bootstrapStatus,omitempty"`
	ProductName                        string        `json:"productName,omitempty"`
	UseGatewayVLANVNID                 bool          `json:"useGatewayVLANVNID"`
	AssociatedGatewaySecurityID        string        `json:"associatedGatewaySecurityID,omitempty"`
	AssociatedGatewaySecurityProfileID string        `json:"associatedGatewaySecurityProfileID,omitempty"`
	AssociatedNSGInfoID                string        `json:"associatedNSGInfoID,omitempty"`
	AssociatedNetconfProfileID         string        `json:"associatedNetconfProfileID,omitempty"`
	Vtep                               string        `json:"vtep,omitempty"`
	AutoDiscGatewayID                  string        `json:"autoDiscGatewayID,omitempty"`
	ExternalID                         string        `json:"externalID,omitempty"`
	SystemID                           string        `json:"systemID,omitempty"`
}

// NewAllGateway returns a new *AllGateway
func NewAllGateway() *AllGateway {

	return &AllGateway{}
}

// Identity returns the Identity of the object.
func (o *AllGateway) Identity() bambou.Identity {

	return AllGatewayIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *AllGateway) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *AllGateway) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the AllGateway from the server
func (o *AllGateway) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the AllGateway into the server
func (o *AllGateway) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the AllGateway from the server
func (o *AllGateway) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the AllGateway
func (o *AllGateway) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the AllGateway
func (o *AllGateway) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the AllGateway
func (o *AllGateway) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the AllGateway
func (o *AllGateway) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
