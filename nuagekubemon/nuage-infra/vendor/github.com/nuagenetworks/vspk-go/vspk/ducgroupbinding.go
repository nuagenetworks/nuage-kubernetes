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

// DUCGroupBindingIdentity represents the Identity of the object
var DUCGroupBindingIdentity = bambou.Identity{
	Name:     "ducgroupbinding",
	Category: "ducgroupbindings",
}

// DUCGroupBindingsList represents a list of DUCGroupBindings
type DUCGroupBindingsList []*DUCGroupBinding

// DUCGroupBindingsAncestor is the interface that an ancestor of a DUCGroupBinding must implement.
// An Ancestor is defined as an entity that has DUCGroupBinding as a descendant.
// An Ancestor can get a list of its child DUCGroupBindings, but not necessarily create one.
type DUCGroupBindingsAncestor interface {
	DUCGroupBindings(*bambou.FetchingInfo) (DUCGroupBindingsList, *bambou.Error)
}

// DUCGroupBindingsParent is the interface that a parent of a DUCGroupBinding must implement.
// A Parent is defined as an entity that has DUCGroupBinding as a child.
// A Parent is an Ancestor which can create a DUCGroupBinding.
type DUCGroupBindingsParent interface {
	DUCGroupBindingsAncestor
	CreateDUCGroupBinding(*DUCGroupBinding) *bambou.Error
}

// DUCGroupBinding represents the model of a ducgroupbinding
type DUCGroupBinding struct {
	ID                   string `json:"ID,omitempty"`
	ParentID             string `json:"parentID,omitempty"`
	ParentType           string `json:"parentType,omitempty"`
	Owner                string `json:"owner,omitempty"`
	OneWayDelay          int    `json:"oneWayDelay,omitempty"`
	Priority             int    `json:"priority,omitempty"`
	AssociatedDUCGroupID string `json:"associatedDUCGroupID,omitempty"`
}

// NewDUCGroupBinding returns a new *DUCGroupBinding
func NewDUCGroupBinding() *DUCGroupBinding {

	return &DUCGroupBinding{
		OneWayDelay: 50,
	}
}

// Identity returns the Identity of the object.
func (o *DUCGroupBinding) Identity() bambou.Identity {

	return DUCGroupBindingIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *DUCGroupBinding) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *DUCGroupBinding) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the DUCGroupBinding from the server
func (o *DUCGroupBinding) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the DUCGroupBinding into the server
func (o *DUCGroupBinding) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the DUCGroupBinding from the server
func (o *DUCGroupBinding) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}
