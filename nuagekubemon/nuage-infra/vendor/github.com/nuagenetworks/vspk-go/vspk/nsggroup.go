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

// NSGGroupIdentity represents the Identity of the object
var NSGGroupIdentity = bambou.Identity{
	Name:     "nsggroup",
	Category: "nsggroups",
}

// NSGGroupsList represents a list of NSGGroups
type NSGGroupsList []*NSGGroup

// NSGGroupsAncestor is the interface that an ancestor of a NSGGroup must implement.
// An Ancestor is defined as an entity that has NSGGroup as a descendant.
// An Ancestor can get a list of its child NSGGroups, but not necessarily create one.
type NSGGroupsAncestor interface {
	NSGGroups(*bambou.FetchingInfo) (NSGGroupsList, *bambou.Error)
}

// NSGGroupsParent is the interface that a parent of a NSGGroup must implement.
// A Parent is defined as an entity that has NSGGroup as a child.
// A Parent is an Ancestor which can create a NSGGroup.
type NSGGroupsParent interface {
	NSGGroupsAncestor
	CreateNSGGroup(*NSGGroup) *bambou.Error
}

// NSGGroup represents the model of a nsggroup
type NSGGroup struct {
	ID          string `json:"ID,omitempty"`
	ParentID    string `json:"parentID,omitempty"`
	ParentType  string `json:"parentType,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// NewNSGGroup returns a new *NSGGroup
func NewNSGGroup() *NSGGroup {

	return &NSGGroup{}
}

// Identity returns the Identity of the object.
func (o *NSGGroup) Identity() bambou.Identity {

	return NSGGroupIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *NSGGroup) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *NSGGroup) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the NSGGroup from the server
func (o *NSGGroup) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the NSGGroup into the server
func (o *NSGGroup) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the NSGGroup from the server
func (o *NSGGroup) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// NSGateways retrieves the list of child NSGateways of the NSGGroup
func (o *NSGGroup) NSGateways(info *bambou.FetchingInfo) (NSGatewaysList, *bambou.Error) {

	var list NSGatewaysList
	err := bambou.CurrentSession().FetchChildren(o, NSGatewayIdentity, &list, info)
	return list, err
}

// AssignNSGateways assigns the list of NSGateways to the NSGGroup
func (o *NSGGroup) AssignNSGateways(children NSGatewaysList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, NSGatewayIdentity)
}

// DUCGroupBindings retrieves the list of child DUCGroupBindings of the NSGGroup
func (o *NSGGroup) DUCGroupBindings(info *bambou.FetchingInfo) (DUCGroupBindingsList, *bambou.Error) {

	var list DUCGroupBindingsList
	err := bambou.CurrentSession().FetchChildren(o, DUCGroupBindingIdentity, &list, info)
	return list, err
}

// CreateDUCGroupBinding creates a new child DUCGroupBinding under the NSGGroup
func (o *NSGGroup) CreateDUCGroupBinding(child *DUCGroupBinding) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
