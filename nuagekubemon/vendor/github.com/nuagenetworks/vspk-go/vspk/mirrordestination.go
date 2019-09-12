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

// MirrorDestinationIdentity represents the Identity of the object
var MirrorDestinationIdentity = bambou.Identity{
	Name:     "mirrordestination",
	Category: "mirrordestinations",
}

// MirrorDestinationsList represents a list of MirrorDestinations
type MirrorDestinationsList []*MirrorDestination

// MirrorDestinationsAncestor is the interface that an ancestor of a MirrorDestination must implement.
// An Ancestor is defined as an entity that has MirrorDestination as a descendant.
// An Ancestor can get a list of its child MirrorDestinations, but not necessarily create one.
type MirrorDestinationsAncestor interface {
	MirrorDestinations(*bambou.FetchingInfo) (MirrorDestinationsList, *bambou.Error)
}

// MirrorDestinationsParent is the interface that a parent of a MirrorDestination must implement.
// A Parent is defined as an entity that has MirrorDestination as a child.
// A Parent is an Ancestor which can create a MirrorDestination.
type MirrorDestinationsParent interface {
	MirrorDestinationsAncestor
	CreateMirrorDestination(*MirrorDestination) *bambou.Error
}

// MirrorDestination represents the model of a mirrordestination
type MirrorDestination struct {
	ID               string        `json:"ID,omitempty"`
	ParentID         string        `json:"parentID,omitempty"`
	ParentType       string        `json:"parentType,omitempty"`
	Owner            string        `json:"owner,omitempty"`
	Name             string        `json:"name,omitempty"`
	LastUpdatedBy    string        `json:"lastUpdatedBy,omitempty"`
	ServiceId        int           `json:"serviceId,omitempty"`
	DestinationIp    string        `json:"destinationIp,omitempty"`
	EmbeddedMetadata []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope      string        `json:"entityScope,omitempty"`
	ExternalID       string        `json:"externalID,omitempty"`
}

// NewMirrorDestination returns a new *MirrorDestination
func NewMirrorDestination() *MirrorDestination {

	return &MirrorDestination{}
}

// Identity returns the Identity of the object.
func (o *MirrorDestination) Identity() bambou.Identity {

	return MirrorDestinationIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *MirrorDestination) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *MirrorDestination) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the MirrorDestination from the server
func (o *MirrorDestination) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the MirrorDestination into the server
func (o *MirrorDestination) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the MirrorDestination from the server
func (o *MirrorDestination) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the MirrorDestination
func (o *MirrorDestination) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the MirrorDestination
func (o *MirrorDestination) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EgressACLEntryTemplates retrieves the list of child EgressACLEntryTemplates of the MirrorDestination
func (o *MirrorDestination) EgressACLEntryTemplates(info *bambou.FetchingInfo) (EgressACLEntryTemplatesList, *bambou.Error) {

	var list EgressACLEntryTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, EgressACLEntryTemplateIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the MirrorDestination
func (o *MirrorDestination) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the MirrorDestination
func (o *MirrorDestination) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IngressACLEntryTemplates retrieves the list of child IngressACLEntryTemplates of the MirrorDestination
func (o *MirrorDestination) IngressACLEntryTemplates(info *bambou.FetchingInfo) (IngressACLEntryTemplatesList, *bambou.Error) {

	var list IngressACLEntryTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressACLEntryTemplateIdentity, &list, info)
	return list, err
}

// IngressAdvFwdEntryTemplates retrieves the list of child IngressAdvFwdEntryTemplates of the MirrorDestination
func (o *MirrorDestination) IngressAdvFwdEntryTemplates(info *bambou.FetchingInfo) (IngressAdvFwdEntryTemplatesList, *bambou.Error) {

	var list IngressAdvFwdEntryTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressAdvFwdEntryTemplateIdentity, &list, info)
	return list, err
}

// VPortMirrors retrieves the list of child VPortMirrors of the MirrorDestination
func (o *MirrorDestination) VPortMirrors(info *bambou.FetchingInfo) (VPortMirrorsList, *bambou.Error) {

	var list VPortMirrorsList
	err := bambou.CurrentSession().FetchChildren(o, VPortMirrorIdentity, &list, info)
	return list, err
}
