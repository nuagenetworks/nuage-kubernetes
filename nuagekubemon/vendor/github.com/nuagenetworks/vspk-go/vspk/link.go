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

// LinkIdentity represents the Identity of the object
var LinkIdentity = bambou.Identity{
	Name:     "link",
	Category: "links",
}

// LinksList represents a list of Links
type LinksList []*Link

// LinksAncestor is the interface that an ancestor of a Link must implement.
// An Ancestor is defined as an entity that has Link as a descendant.
// An Ancestor can get a list of its child Links, but not necessarily create one.
type LinksAncestor interface {
	Links(*bambou.FetchingInfo) (LinksList, *bambou.Error)
}

// LinksParent is the interface that a parent of a Link must implement.
// A Parent is defined as an entity that has Link as a child.
// A Parent is an Ancestor which can create a Link.
type LinksParent interface {
	LinksAncestor
	CreateLink(*Link) *bambou.Error
}

// Link represents the model of a link
type Link struct {
	ID                        string        `json:"ID,omitempty"`
	ParentID                  string        `json:"parentID,omitempty"`
	ParentType                string        `json:"parentType,omitempty"`
	Owner                     string        `json:"owner,omitempty"`
	LastUpdatedBy             string        `json:"lastUpdatedBy,omitempty"`
	AcceptanceCriteria        string        `json:"acceptanceCriteria,omitempty"`
	ReadOnly                  bool          `json:"readOnly"`
	EmbeddedMetadata          []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope               string        `json:"entityScope,omitempty"`
	AssociatedDestinationID   string        `json:"associatedDestinationID,omitempty"`
	AssociatedDestinationName string        `json:"associatedDestinationName,omitempty"`
	AssociatedDestinationType string        `json:"associatedDestinationType,omitempty"`
	AssociatedSourceID        string        `json:"associatedSourceID,omitempty"`
	AssociatedSourceName      string        `json:"associatedSourceName,omitempty"`
	AssociatedSourceType      string        `json:"associatedSourceType,omitempty"`
	ExternalID                string        `json:"externalID,omitempty"`
	Type                      string        `json:"type,omitempty"`
}

// NewLink returns a new *Link
func NewLink() *Link {

	return &Link{
		AcceptanceCriteria: "ALL",
	}
}

// Identity returns the Identity of the object.
func (o *Link) Identity() bambou.Identity {

	return LinkIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Link) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Link) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Link from the server
func (o *Link) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Link into the server
func (o *Link) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Link from the server
func (o *Link) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// DemarcationServices retrieves the list of child DemarcationServices of the Link
func (o *Link) DemarcationServices(info *bambou.FetchingInfo) (DemarcationServicesList, *bambou.Error) {

	var list DemarcationServicesList
	err := bambou.CurrentSession().FetchChildren(o, DemarcationServiceIdentity, &list, info)
	return list, err
}

// CreateDemarcationService creates a new child DemarcationService under the Link
func (o *Link) CreateDemarcationService(child *DemarcationService) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the Link
func (o *Link) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the Link
func (o *Link) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NextHops retrieves the list of child NextHops of the Link
func (o *Link) NextHops(info *bambou.FetchingInfo) (NextHopsList, *bambou.Error) {

	var list NextHopsList
	err := bambou.CurrentSession().FetchChildren(o, NextHopIdentity, &list, info)
	return list, err
}

// CreateNextHop creates a new child NextHop under the Link
func (o *Link) CreateNextHop(child *NextHop) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Link
func (o *Link) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Link
func (o *Link) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PolicyStatements retrieves the list of child PolicyStatements of the Link
func (o *Link) PolicyStatements(info *bambou.FetchingInfo) (PolicyStatementsList, *bambou.Error) {

	var list PolicyStatementsList
	err := bambou.CurrentSession().FetchChildren(o, PolicyStatementIdentity, &list, info)
	return list, err
}

// CreatePolicyStatement creates a new child PolicyStatement under the Link
func (o *Link) CreatePolicyStatement(child *PolicyStatement) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CSNATPools retrieves the list of child CSNATPools of the Link
func (o *Link) CSNATPools(info *bambou.FetchingInfo) (CSNATPoolsList, *bambou.Error) {

	var list CSNATPoolsList
	err := bambou.CurrentSession().FetchChildren(o, CSNATPoolIdentity, &list, info)
	return list, err
}

// CreateCSNATPool creates a new child CSNATPool under the Link
func (o *Link) CreateCSNATPool(child *CSNATPool) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PSNATPools retrieves the list of child PSNATPools of the Link
func (o *Link) PSNATPools(info *bambou.FetchingInfo) (PSNATPoolsList, *bambou.Error) {

	var list PSNATPoolsList
	err := bambou.CurrentSession().FetchChildren(o, PSNATPoolIdentity, &list, info)
	return list, err
}

// CreatePSNATPool creates a new child PSNATPool under the Link
func (o *Link) CreatePSNATPool(child *PSNATPool) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// OverlayAddressPools retrieves the list of child OverlayAddressPools of the Link
func (o *Link) OverlayAddressPools(info *bambou.FetchingInfo) (OverlayAddressPoolsList, *bambou.Error) {

	var list OverlayAddressPoolsList
	err := bambou.CurrentSession().FetchChildren(o, OverlayAddressPoolIdentity, &list, info)
	return list, err
}

// CreateOverlayAddressPool creates a new child OverlayAddressPool under the Link
func (o *Link) CreateOverlayAddressPool(child *OverlayAddressPool) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
