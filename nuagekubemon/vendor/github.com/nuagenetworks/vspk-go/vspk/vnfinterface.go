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

// VNFInterfaceIdentity represents the Identity of the object
var VNFInterfaceIdentity = bambou.Identity{
	Name:     "vnfinterface",
	Category: "vnfinterfaces",
}

// VNFInterfacesList represents a list of VNFInterfaces
type VNFInterfacesList []*VNFInterface

// VNFInterfacesAncestor is the interface that an ancestor of a VNFInterface must implement.
// An Ancestor is defined as an entity that has VNFInterface as a descendant.
// An Ancestor can get a list of its child VNFInterfaces, but not necessarily create one.
type VNFInterfacesAncestor interface {
	VNFInterfaces(*bambou.FetchingInfo) (VNFInterfacesList, *bambou.Error)
}

// VNFInterfacesParent is the interface that a parent of a VNFInterface must implement.
// A Parent is defined as an entity that has VNFInterface as a child.
// A Parent is an Ancestor which can create a VNFInterface.
type VNFInterfacesParent interface {
	VNFInterfacesAncestor
	CreateVNFInterface(*VNFInterface) *bambou.Error
}

// VNFInterface represents the model of a vnfinterface
type VNFInterface struct {
	ID                  string        `json:"ID,omitempty"`
	ParentID            string        `json:"parentID,omitempty"`
	ParentType          string        `json:"parentType,omitempty"`
	Owner               string        `json:"owner,omitempty"`
	MAC                 string        `json:"MAC,omitempty"`
	VNFUUID             string        `json:"VNFUUID,omitempty"`
	IPAddress           string        `json:"IPAddress,omitempty"`
	VPortID             string        `json:"VPortID,omitempty"`
	VPortName           string        `json:"VPortName,omitempty"`
	IPv6Address         string        `json:"IPv6Address,omitempty"`
	IPv6Gateway         string        `json:"IPv6Gateway,omitempty"`
	Name                string        `json:"name,omitempty"`
	LastUpdatedBy       string        `json:"lastUpdatedBy,omitempty"`
	Gateway             string        `json:"gateway,omitempty"`
	Netmask             string        `json:"netmask,omitempty"`
	NetworkName         string        `json:"networkName,omitempty"`
	EmbeddedMetadata    []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope         string        `json:"entityScope,omitempty"`
	PolicyDecisionID    string        `json:"policyDecisionID,omitempty"`
	DomainID            string        `json:"domainID,omitempty"`
	DomainName          string        `json:"domainName,omitempty"`
	ZoneID              string        `json:"zoneID,omitempty"`
	ZoneName            string        `json:"zoneName,omitempty"`
	AttachedNetworkID   string        `json:"attachedNetworkID,omitempty"`
	AttachedNetworkType string        `json:"attachedNetworkType,omitempty"`
	ExternalID          string        `json:"externalID,omitempty"`
	Type                string        `json:"type,omitempty"`
}

// NewVNFInterface returns a new *VNFInterface
func NewVNFInterface() *VNFInterface {

	return &VNFInterface{}
}

// Identity returns the Identity of the object.
func (o *VNFInterface) Identity() bambou.Identity {

	return VNFInterfaceIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VNFInterface) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VNFInterface) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VNFInterface from the server
func (o *VNFInterface) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VNFInterface into the server
func (o *VNFInterface) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VNFInterface from the server
func (o *VNFInterface) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the VNFInterface
func (o *VNFInterface) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VNFInterface
func (o *VNFInterface) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VNFInterface
func (o *VNFInterface) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VNFInterface
func (o *VNFInterface) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
