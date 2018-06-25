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

// UserIdentity represents the Identity of the object
var UserIdentity = bambou.Identity{
	Name:     "user",
	Category: "users",
}

// UsersList represents a list of Users
type UsersList []*User

// UsersAncestor is the interface that an ancestor of a User must implement.
// An Ancestor is defined as an entity that has User as a descendant.
// An Ancestor can get a list of its child Users, but not necessarily create one.
type UsersAncestor interface {
	Users(*bambou.FetchingInfo) (UsersList, *bambou.Error)
}

// UsersParent is the interface that a parent of a User must implement.
// A Parent is defined as an entity that has User as a child.
// A Parent is an Ancestor which can create a User.
type UsersParent interface {
	UsersAncestor
	CreateUser(*User) *bambou.Error
}

// User represents the model of a user
type User struct {
	ID             string `json:"ID,omitempty"`
	ParentID       string `json:"parentID,omitempty"`
	ParentType     string `json:"parentType,omitempty"`
	Owner          string `json:"owner,omitempty"`
	LDAPUserDN     string `json:"LDAPUserDN,omitempty"`
	ManagementMode string `json:"managementMode,omitempty"`
	Password       string `json:"password,omitempty"`
	LastName       string `json:"lastName,omitempty"`
	LastUpdatedBy  string `json:"lastUpdatedBy,omitempty"`
	FirstName      string `json:"firstName,omitempty"`
	Disabled       bool   `json:"disabled"`
	Email          string `json:"email,omitempty"`
	EntityScope    string `json:"entityScope,omitempty"`
	MobileNumber   string `json:"mobileNumber,omitempty"`
	UserName       string `json:"userName,omitempty"`
	AvatarData     string `json:"avatarData,omitempty"`
	AvatarType     string `json:"avatarType,omitempty"`
	ExternalID     string `json:"externalID,omitempty"`
}

// NewUser returns a new *User
func NewUser() *User {

	return &User{}
}

// Identity returns the Identity of the object.
func (o *User) Identity() bambou.Identity {

	return UserIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *User) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *User) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the User from the server
func (o *User) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the User into the server
func (o *User) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the User from the server
func (o *User) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the User
func (o *User) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the User
func (o *User) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the User
func (o *User) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the User
func (o *User) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMs retrieves the list of child VMs of the User
func (o *User) VMs(info *bambou.FetchingInfo) (VMsList, *bambou.Error) {

	var list VMsList
	err := bambou.CurrentSession().FetchChildren(o, VMIdentity, &list, info)
	return list, err
}

// Containers retrieves the list of child Containers of the User
func (o *User) Containers(info *bambou.FetchingInfo) (ContainersList, *bambou.Error) {

	var list ContainersList
	err := bambou.CurrentSession().FetchChildren(o, ContainerIdentity, &list, info)
	return list, err
}

// Groups retrieves the list of child Groups of the User
func (o *User) Groups(info *bambou.FetchingInfo) (GroupsList, *bambou.Error) {

	var list GroupsList
	err := bambou.CurrentSession().FetchChildren(o, GroupIdentity, &list, info)
	return list, err
}

// Avatars retrieves the list of child Avatars of the User
func (o *User) Avatars(info *bambou.FetchingInfo) (AvatarsList, *bambou.Error) {

	var list AvatarsList
	err := bambou.CurrentSession().FetchChildren(o, AvatarIdentity, &list, info)
	return list, err
}

// CreateAvatar creates a new child Avatar under the User
func (o *User) CreateAvatar(child *Avatar) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the User
func (o *User) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
