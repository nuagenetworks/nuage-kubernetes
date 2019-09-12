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

// MirrorDestinationGroupIdentity represents the Identity of the object
var MirrorDestinationGroupIdentity = bambou.Identity{
	Name:     "mirrordestinationgroup",
	Category: "mirrordestinationgroups",
}

// MirrorDestinationGroupsList represents a list of MirrorDestinationGroups
type MirrorDestinationGroupsList []*MirrorDestinationGroup

// MirrorDestinationGroupsAncestor is the interface that an ancestor of a MirrorDestinationGroup must implement.
// An Ancestor is defined as an entity that has MirrorDestinationGroup as a descendant.
// An Ancestor can get a list of its child MirrorDestinationGroups, but not necessarily create one.
type MirrorDestinationGroupsAncestor interface {
	MirrorDestinationGroups(*bambou.FetchingInfo) (MirrorDestinationGroupsList, *bambou.Error)
}

// MirrorDestinationGroupsParent is the interface that a parent of a MirrorDestinationGroup must implement.
// A Parent is defined as an entity that has MirrorDestinationGroup as a child.
// A Parent is an Ancestor which can create a MirrorDestinationGroup.
type MirrorDestinationGroupsParent interface {
	MirrorDestinationGroupsAncestor
	CreateMirrorDestinationGroup(*MirrorDestinationGroup) *bambou.Error
}

// MirrorDestinationGroup represents the model of a mirrordestinationgroup
type MirrorDestinationGroup struct {
	ID               string        `json:"ID,omitempty"`
	ParentID         string        `json:"parentID,omitempty"`
	ParentType       string        `json:"parentType,omitempty"`
	Owner            string        `json:"owner,omitempty"`
	Name             string        `json:"name,omitempty"`
	LastUpdatedBy    string        `json:"lastUpdatedBy,omitempty"`
	Description      string        `json:"description,omitempty"`
	EmbeddedMetadata []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope      string        `json:"entityScope,omitempty"`
	ExternalID       string        `json:"externalID,omitempty"`
}

// NewMirrorDestinationGroup returns a new *MirrorDestinationGroup
func NewMirrorDestinationGroup() *MirrorDestinationGroup {

	return &MirrorDestinationGroup{}
}

// Identity returns the Identity of the object.
func (o *MirrorDestinationGroup) Identity() bambou.Identity {

	return MirrorDestinationGroupIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *MirrorDestinationGroup) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *MirrorDestinationGroup) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the MirrorDestinationGroup from the server
func (o *MirrorDestinationGroup) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the MirrorDestinationGroup into the server
func (o *MirrorDestinationGroup) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the MirrorDestinationGroup from the server
func (o *MirrorDestinationGroup) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the MirrorDestinationGroup
func (o *MirrorDestinationGroup) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the MirrorDestinationGroup
func (o *MirrorDestinationGroup) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// MirrorDestinations retrieves the list of child MirrorDestinations of the MirrorDestinationGroup
func (o *MirrorDestinationGroup) MirrorDestinations(info *bambou.FetchingInfo) (MirrorDestinationsList, *bambou.Error) {

	var list MirrorDestinationsList
	err := bambou.CurrentSession().FetchChildren(o, MirrorDestinationIdentity, &list, info)
	return list, err
}

// AssignMirrorDestinations assigns the list of MirrorDestinations to the MirrorDestinationGroup
func (o *MirrorDestinationGroup) AssignMirrorDestinations(children MirrorDestinationsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, MirrorDestinationIdentity)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the MirrorDestinationGroup
func (o *MirrorDestinationGroup) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the MirrorDestinationGroup
func (o *MirrorDestinationGroup) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// OverlayMirrorDestinations retrieves the list of child OverlayMirrorDestinations of the MirrorDestinationGroup
func (o *MirrorDestinationGroup) OverlayMirrorDestinations(info *bambou.FetchingInfo) (OverlayMirrorDestinationsList, *bambou.Error) {

	var list OverlayMirrorDestinationsList
	err := bambou.CurrentSession().FetchChildren(o, OverlayMirrorDestinationIdentity, &list, info)
	return list, err
}

// AssignOverlayMirrorDestinations assigns the list of OverlayMirrorDestinations to the MirrorDestinationGroup
func (o *MirrorDestinationGroup) AssignOverlayMirrorDestinations(children OverlayMirrorDestinationsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, OverlayMirrorDestinationIdentity)
}
