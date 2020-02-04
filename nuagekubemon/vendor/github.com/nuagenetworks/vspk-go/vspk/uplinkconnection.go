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

// UplinkConnectionIdentity represents the Identity of the object
var UplinkConnectionIdentity = bambou.Identity{
	Name:     "uplinkconnection",
	Category: "uplinkconnections",
}

// UplinkConnectionsList represents a list of UplinkConnections
type UplinkConnectionsList []*UplinkConnection

// UplinkConnectionsAncestor is the interface that an ancestor of a UplinkConnection must implement.
// An Ancestor is defined as an entity that has UplinkConnection as a descendant.
// An Ancestor can get a list of its child UplinkConnections, but not necessarily create one.
type UplinkConnectionsAncestor interface {
	UplinkConnections(*bambou.FetchingInfo) (UplinkConnectionsList, *bambou.Error)
}

// UplinkConnectionsParent is the interface that a parent of a UplinkConnection must implement.
// A Parent is defined as an entity that has UplinkConnection as a child.
// A Parent is an Ancestor which can create a UplinkConnection.
type UplinkConnectionsParent interface {
	UplinkConnectionsAncestor
	CreateUplinkConnection(*UplinkConnection) *bambou.Error
}

// UplinkConnection represents the model of a uplinkconnection
type UplinkConnection struct {
	ID                      string        `json:"ID,omitempty"`
	ParentID                string        `json:"parentID,omitempty"`
	ParentType              string        `json:"parentType,omitempty"`
	Owner                   string        `json:"owner,omitempty"`
	PATEnabled              bool          `json:"PATEnabled"`
	DNSAddress              string        `json:"DNSAddress,omitempty"`
	DNSAddressV6            string        `json:"DNSAddressV6,omitempty"`
	Password                string        `json:"password,omitempty"`
	LastUpdatedBy           string        `json:"lastUpdatedBy,omitempty"`
	Gateway                 string        `json:"gateway,omitempty"`
	GatewayV6               string        `json:"gatewayV6,omitempty"`
	Address                 string        `json:"address,omitempty"`
	AddressFamily           string        `json:"addressFamily,omitempty"`
	AddressV6               string        `json:"addressV6,omitempty"`
	AdvertisementCriteria   string        `json:"advertisementCriteria,omitempty"`
	SecondaryAddress        string        `json:"secondaryAddress,omitempty"`
	Netmask                 string        `json:"netmask,omitempty"`
	Vlan                    int           `json:"vlan,omitempty"`
	EmbeddedMetadata        []interface{} `json:"embeddedMetadata,omitempty"`
	UnderlayEnabled         bool          `json:"underlayEnabled"`
	UnderlayID              int           `json:"underlayID,omitempty"`
	Inherited               bool          `json:"inherited"`
	InstallerManaged        bool          `json:"installerManaged"`
	InterfaceConnectionType string        `json:"interfaceConnectionType,omitempty"`
	EntityScope             string        `json:"entityScope,omitempty"`
	Mode                    string        `json:"mode,omitempty"`
	Role                    string        `json:"role,omitempty"`
	RoleOrder               int           `json:"roleOrder,omitempty"`
	PortName                string        `json:"portName,omitempty"`
	DownloadRateLimit       float64       `json:"downloadRateLimit,omitempty"`
	UplinkID                int           `json:"uplinkID,omitempty"`
	Username                string        `json:"username,omitempty"`
	AssocUnderlayID         string        `json:"assocUnderlayID,omitempty"`
	AssociatedBGPNeighborID string        `json:"associatedBGPNeighborID,omitempty"`
	AssociatedUnderlayName  string        `json:"associatedUnderlayName,omitempty"`
	AuxMode                 string        `json:"auxMode,omitempty"`
	AuxiliaryLink           bool          `json:"auxiliaryLink"`
	ExternalID              string        `json:"externalID,omitempty"`
}

// NewUplinkConnection returns a new *UplinkConnection
func NewUplinkConnection() *UplinkConnection {

	return &UplinkConnection{
		PATEnabled:              true,
		AddressFamily:           "IPV4",
		UnderlayEnabled:         true,
		Inherited:               false,
		InstallerManaged:        false,
		InterfaceConnectionType: "AUTOMATIC",
		Mode:                    "Dynamic",
		Role:                    "PRIMARY",
		DownloadRateLimit:       8.0,
		AuxMode:                 "NONE",
		AuxiliaryLink:           false,
	}
}

// Identity returns the Identity of the object.
func (o *UplinkConnection) Identity() bambou.Identity {

	return UplinkConnectionIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *UplinkConnection) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *UplinkConnection) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the UplinkConnection from the server
func (o *UplinkConnection) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the UplinkConnection into the server
func (o *UplinkConnection) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the UplinkConnection from the server
func (o *UplinkConnection) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the UplinkConnection
func (o *UplinkConnection) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the UplinkConnection
func (o *UplinkConnection) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// BFDSessions retrieves the list of child BFDSessions of the UplinkConnection
func (o *UplinkConnection) BFDSessions(info *bambou.FetchingInfo) (BFDSessionsList, *bambou.Error) {

	var list BFDSessionsList
	err := bambou.CurrentSession().FetchChildren(o, BFDSessionIdentity, &list, info)
	return list, err
}

// CreateBFDSession creates a new child BFDSession under the UplinkConnection
func (o *UplinkConnection) CreateBFDSession(child *BFDSession) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the UplinkConnection
func (o *UplinkConnection) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the UplinkConnection
func (o *UplinkConnection) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CustomProperties retrieves the list of child CustomProperties of the UplinkConnection
func (o *UplinkConnection) CustomProperties(info *bambou.FetchingInfo) (CustomPropertiesList, *bambou.Error) {

	var list CustomPropertiesList
	err := bambou.CurrentSession().FetchChildren(o, CustomPropertyIdentity, &list, info)
	return list, err
}

// CreateCustomProperty creates a new child CustomProperty under the UplinkConnection
func (o *UplinkConnection) CreateCustomProperty(child *CustomProperty) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
