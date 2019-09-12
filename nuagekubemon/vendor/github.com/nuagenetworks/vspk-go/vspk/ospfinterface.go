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

// OSPFInterfaceIdentity represents the Identity of the object
var OSPFInterfaceIdentity = bambou.Identity{
	Name:     "ospfinterface",
	Category: "ospfinterfaces",
}

// OSPFInterfacesList represents a list of OSPFInterfaces
type OSPFInterfacesList []*OSPFInterface

// OSPFInterfacesAncestor is the interface that an ancestor of a OSPFInterface must implement.
// An Ancestor is defined as an entity that has OSPFInterface as a descendant.
// An Ancestor can get a list of its child OSPFInterfaces, but not necessarily create one.
type OSPFInterfacesAncestor interface {
	OSPFInterfaces(*bambou.FetchingInfo) (OSPFInterfacesList, *bambou.Error)
}

// OSPFInterfacesParent is the interface that a parent of a OSPFInterface must implement.
// A Parent is defined as an entity that has OSPFInterface as a child.
// A Parent is an Ancestor which can create a OSPFInterface.
type OSPFInterfacesParent interface {
	OSPFInterfacesAncestor
	CreateOSPFInterface(*OSPFInterface) *bambou.Error
}

// OSPFInterface represents the model of a ospfinterface
type OSPFInterface struct {
	ID                 string        `json:"ID,omitempty"`
	ParentID           string        `json:"parentID,omitempty"`
	ParentType         string        `json:"parentType,omitempty"`
	Owner              string        `json:"owner,omitempty"`
	Name               string        `json:"name,omitempty"`
	PassiveEnabled     bool          `json:"passiveEnabled"`
	LastUpdatedBy      string        `json:"lastUpdatedBy,omitempty"`
	AdminState         string        `json:"adminState,omitempty"`
	DeadInterval       int           `json:"deadInterval,omitempty"`
	HelloInterval      int           `json:"helloInterval,omitempty"`
	Description        string        `json:"description,omitempty"`
	MessageDigestKeys  []interface{} `json:"messageDigestKeys,omitempty"`
	Metric             int           `json:"metric,omitempty"`
	EmbeddedMetadata   []interface{} `json:"embeddedMetadata,omitempty"`
	InterfaceType      string        `json:"interfaceType,omitempty"`
	EntityScope        string        `json:"entityScope,omitempty"`
	Priority           int           `json:"priority,omitempty"`
	AssociatedSubnetID string        `json:"associatedSubnetID,omitempty"`
	Mtu                int           `json:"mtu,omitempty"`
	AuthenticationKey  string        `json:"authenticationKey,omitempty"`
	AuthenticationType string        `json:"authenticationType,omitempty"`
	ExternalID         string        `json:"externalID,omitempty"`
}

// NewOSPFInterface returns a new *OSPFInterface
func NewOSPFInterface() *OSPFInterface {

	return &OSPFInterface{
		PassiveEnabled:     false,
		AdminState:         "UP",
		DeadInterval:       40,
		HelloInterval:      10,
		InterfaceType:      "BROADCAST",
		Priority:           1,
		AuthenticationType: "NONE",
	}
}

// Identity returns the Identity of the object.
func (o *OSPFInterface) Identity() bambou.Identity {

	return OSPFInterfaceIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *OSPFInterface) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *OSPFInterface) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the OSPFInterface from the server
func (o *OSPFInterface) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the OSPFInterface into the server
func (o *OSPFInterface) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the OSPFInterface from the server
func (o *OSPFInterface) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the OSPFInterface
func (o *OSPFInterface) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the OSPFInterface
func (o *OSPFInterface) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the OSPFInterface
func (o *OSPFInterface) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the OSPFInterface
func (o *OSPFInterface) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
