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
	ID                      string `json:"ID,omitempty"`
	ParentID                string `json:"parentID,omitempty"`
	ParentType              string `json:"parentType,omitempty"`
	Owner                   string `json:"owner,omitempty"`
	DNSAddress              string `json:"DNSAddress,omitempty"`
	Password                string `json:"password,omitempty"`
	Gateway                 string `json:"gateway,omitempty"`
	Address                 string `json:"address,omitempty"`
	AdvertisementCriteria   string `json:"advertisementCriteria,omitempty"`
	Netmask                 string `json:"netmask,omitempty"`
	InterfaceConnectionType string `json:"interfaceConnectionType,omitempty"`
	Mode                    string `json:"mode,omitempty"`
	Role                    string `json:"role,omitempty"`
	UplinkID                string `json:"uplinkID,omitempty"`
	Username                string `json:"username,omitempty"`
	AssocUnderlayID         string `json:"assocUnderlayID,omitempty"`
	AssociatedUnderlayName  string `json:"associatedUnderlayName,omitempty"`
	AssociatedVSCProfileID  string `json:"associatedVSCProfileID,omitempty"`
	AuxiliaryLink           bool   `json:"auxiliaryLink"`
}

// NewUplinkConnection returns a new *UplinkConnection
func NewUplinkConnection() *UplinkConnection {

	return &UplinkConnection{
		Address:                 "IPv4",
		InterfaceConnectionType: "AUTOMATIC",
		Mode:          "Dynamic",
		Role:          "PRIMARY",
		AuxiliaryLink: false,
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

// Underlays retrieves the list of child Underlays of the UplinkConnection
func (o *UplinkConnection) Underlays(info *bambou.FetchingInfo) (UnderlaysList, *bambou.Error) {

	var list UnderlaysList
	err := bambou.CurrentSession().FetchChildren(o, UnderlayIdentity, &list, info)
	return list, err
}

// CreateUnderlay creates a new child Underlay under the UplinkConnection
func (o *UplinkConnection) CreateUnderlay(child *Underlay) *bambou.Error {

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
