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

// GroupIdentity represents the Identity of the object
var GroupIdentity = bambou.Identity{
	Name:     "group",
	Category: "groups",
}

// GroupsList represents a list of Groups
type GroupsList []*Group

// GroupsAncestor is the interface that an ancestor of a Group must implement.
// An Ancestor is defined as an entity that has Group as a descendant.
// An Ancestor can get a list of its child Groups, but not necessarily create one.
type GroupsAncestor interface {
	Groups(*bambou.FetchingInfo) (GroupsList, *bambou.Error)
}

// GroupsParent is the interface that a parent of a Group must implement.
// A Parent is defined as an entity that has Group as a child.
// A Parent is an Ancestor which can create a Group.
type GroupsParent interface {
	GroupsAncestor
	CreateGroup(*Group) *bambou.Error
}

// Group represents the model of a group
type Group struct {
	ID                  string  `json:"ID,omitempty"`
	ParentID            string  `json:"parentID,omitempty"`
	ParentType          string  `json:"parentType,omitempty"`
	Owner               string  `json:"owner,omitempty"`
	Name                string  `json:"name,omitempty"`
	ManagementMode      string  `json:"managementMode,omitempty"`
	LastUpdatedBy       string  `json:"lastUpdatedBy,omitempty"`
	AccountRestrictions bool    `json:"accountRestrictions"`
	Description         string  `json:"description,omitempty"`
	RestrictionDate     float64 `json:"restrictionDate,omitempty"`
	EntityScope         string  `json:"entityScope,omitempty"`
	Role                string  `json:"role,omitempty"`
	Private             bool    `json:"private"`
	ExternalID          string  `json:"externalID,omitempty"`
}

// NewGroup returns a new *Group
func NewGroup() *Group {

	return &Group{}
}

// Identity returns the Identity of the object.
func (o *Group) Identity() bambou.Identity {

	return GroupIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Group) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Group) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Group from the server
func (o *Group) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Group into the server
func (o *Group) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Group from the server
func (o *Group) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the Group
func (o *Group) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the Group
func (o *Group) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Group
func (o *Group) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Group
func (o *Group) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Users retrieves the list of child Users of the Group
func (o *Group) Users(info *bambou.FetchingInfo) (UsersList, *bambou.Error) {

	var list UsersList
	err := bambou.CurrentSession().FetchChildren(o, UserIdentity, &list, info)
	return list, err
}

// AssignUsers assigns the list of Users to the Group
func (o *Group) AssignUsers(children UsersList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, UserIdentity)
}

// EventLogs retrieves the list of child EventLogs of the Group
func (o *Group) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
