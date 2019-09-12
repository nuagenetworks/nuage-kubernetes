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

// IPReservationIdentity represents the Identity of the object
var IPReservationIdentity = bambou.Identity{
	Name:     "ipreservation",
	Category: "ipreservations",
}

// IPReservationsList represents a list of IPReservations
type IPReservationsList []*IPReservation

// IPReservationsAncestor is the interface that an ancestor of a IPReservation must implement.
// An Ancestor is defined as an entity that has IPReservation as a descendant.
// An Ancestor can get a list of its child IPReservations, but not necessarily create one.
type IPReservationsAncestor interface {
	IPReservations(*bambou.FetchingInfo) (IPReservationsList, *bambou.Error)
}

// IPReservationsParent is the interface that a parent of a IPReservation must implement.
// A Parent is defined as an entity that has IPReservation as a child.
// A Parent is an Ancestor which can create a IPReservation.
type IPReservationsParent interface {
	IPReservationsAncestor
	CreateIPReservation(*IPReservation) *bambou.Error
}

// IPReservation represents the model of a ipreservation
type IPReservation struct {
	ID                       string        `json:"ID,omitempty"`
	ParentID                 string        `json:"parentID,omitempty"`
	ParentType               string        `json:"parentType,omitempty"`
	Owner                    string        `json:"owner,omitempty"`
	MAC                      string        `json:"MAC,omitempty"`
	IPAddress                string        `json:"IPAddress,omitempty"`
	LastUpdatedBy            string        `json:"lastUpdatedBy,omitempty"`
	EmbeddedMetadata         []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope              string        `json:"entityScope,omitempty"`
	ExternalID               string        `json:"externalID,omitempty"`
	DynamicAllocationEnabled bool          `json:"dynamicAllocationEnabled"`
}

// NewIPReservation returns a new *IPReservation
func NewIPReservation() *IPReservation {

	return &IPReservation{}
}

// Identity returns the Identity of the object.
func (o *IPReservation) Identity() bambou.Identity {

	return IPReservationIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *IPReservation) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *IPReservation) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the IPReservation from the server
func (o *IPReservation) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the IPReservation into the server
func (o *IPReservation) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the IPReservation from the server
func (o *IPReservation) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the IPReservation
func (o *IPReservation) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the IPReservation
func (o *IPReservation) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the IPReservation
func (o *IPReservation) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the IPReservation
func (o *IPReservation) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the IPReservation
func (o *IPReservation) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
