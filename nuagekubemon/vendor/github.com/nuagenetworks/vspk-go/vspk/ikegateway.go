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

// IKEGatewayIdentity represents the Identity of the object
var IKEGatewayIdentity = bambou.Identity{
	Name:     "ikegateway",
	Category: "ikegateways",
}

// IKEGatewaysList represents a list of IKEGateways
type IKEGatewaysList []*IKEGateway

// IKEGatewaysAncestor is the interface that an ancestor of a IKEGateway must implement.
// An Ancestor is defined as an entity that has IKEGateway as a descendant.
// An Ancestor can get a list of its child IKEGateways, but not necessarily create one.
type IKEGatewaysAncestor interface {
	IKEGateways(*bambou.FetchingInfo) (IKEGatewaysList, *bambou.Error)
}

// IKEGatewaysParent is the interface that a parent of a IKEGateway must implement.
// A Parent is defined as an entity that has IKEGateway as a child.
// A Parent is an Ancestor which can create a IKEGateway.
type IKEGatewaysParent interface {
	IKEGatewaysAncestor
	CreateIKEGateway(*IKEGateway) *bambou.Error
}

// IKEGateway represents the model of a ikegateway
type IKEGateway struct {
	ID                     string `json:"ID,omitempty"`
	ParentID               string `json:"parentID,omitempty"`
	ParentType             string `json:"parentType,omitempty"`
	Owner                  string `json:"owner,omitempty"`
	IKEVersion             string `json:"IKEVersion,omitempty"`
	IKEv1Mode              string `json:"IKEv1Mode,omitempty"`
	IPAddress              string `json:"IPAddress,omitempty"`
	Name                   string `json:"name,omitempty"`
	LastUpdatedBy          string `json:"lastUpdatedBy,omitempty"`
	Description            string `json:"description,omitempty"`
	EntityScope            string `json:"entityScope,omitempty"`
	AssociatedEnterpriseID string `json:"associatedEnterpriseID,omitempty"`
	ExternalID             string `json:"externalID,omitempty"`
}

// NewIKEGateway returns a new *IKEGateway
func NewIKEGateway() *IKEGateway {

	return &IKEGateway{
		IKEVersion: "V2",
		IKEv1Mode:  "NONE",
	}
}

// Identity returns the Identity of the object.
func (o *IKEGateway) Identity() bambou.Identity {

	return IKEGatewayIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *IKEGateway) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *IKEGateway) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the IKEGateway from the server
func (o *IKEGateway) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the IKEGateway into the server
func (o *IKEGateway) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the IKEGateway from the server
func (o *IKEGateway) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the IKEGateway
func (o *IKEGateway) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the IKEGateway
func (o *IKEGateway) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IKEGatewayConfigs retrieves the list of child IKEGatewayConfigs of the IKEGateway
func (o *IKEGateway) IKEGatewayConfigs(info *bambou.FetchingInfo) (IKEGatewayConfigsList, *bambou.Error) {

	var list IKEGatewayConfigsList
	err := bambou.CurrentSession().FetchChildren(o, IKEGatewayConfigIdentity, &list, info)
	return list, err
}

// AssignIKEGatewayConfigs assigns the list of IKEGatewayConfigs to the IKEGateway
func (o *IKEGateway) AssignIKEGatewayConfigs(children IKEGatewayConfigsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, IKEGatewayConfigIdentity)
}

// IKESubnets retrieves the list of child IKESubnets of the IKEGateway
func (o *IKEGateway) IKESubnets(info *bambou.FetchingInfo) (IKESubnetsList, *bambou.Error) {

	var list IKESubnetsList
	err := bambou.CurrentSession().FetchChildren(o, IKESubnetIdentity, &list, info)
	return list, err
}

// CreateIKESubnet creates a new child IKESubnet under the IKEGateway
func (o *IKEGateway) CreateIKESubnet(child *IKESubnet) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the IKEGateway
func (o *IKEGateway) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the IKEGateway
func (o *IKEGateway) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
