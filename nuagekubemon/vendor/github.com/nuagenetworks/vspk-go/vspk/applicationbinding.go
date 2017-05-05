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

// ApplicationBindingIdentity represents the Identity of the object
var ApplicationBindingIdentity = bambou.Identity{
	Name:     "applicationbinding",
	Category: "applicationbindings",
}

// ApplicationBindingsList represents a list of ApplicationBindings
type ApplicationBindingsList []*ApplicationBinding

// ApplicationBindingsAncestor is the interface that an ancestor of a ApplicationBinding must implement.
// An Ancestor is defined as an entity that has ApplicationBinding as a descendant.
// An Ancestor can get a list of its child ApplicationBindings, but not necessarily create one.
type ApplicationBindingsAncestor interface {
	ApplicationBindings(*bambou.FetchingInfo) (ApplicationBindingsList, *bambou.Error)
}

// ApplicationBindingsParent is the interface that a parent of a ApplicationBinding must implement.
// A Parent is defined as an entity that has ApplicationBinding as a child.
// A Parent is an Ancestor which can create a ApplicationBinding.
type ApplicationBindingsParent interface {
	ApplicationBindingsAncestor
	CreateApplicationBinding(*ApplicationBinding) *bambou.Error
}

// ApplicationBinding represents the model of a applicationbinding
type ApplicationBinding struct {
	ID                      string `json:"ID,omitempty"`
	ParentID                string `json:"parentID,omitempty"`
	ParentType              string `json:"parentType,omitempty"`
	Owner                   string `json:"owner,omitempty"`
	ReadOnly                bool   `json:"readOnly"`
	Priority                int    `json:"priority,omitempty"`
	AssociatedApplicationID string `json:"associatedApplicationID,omitempty"`
}

// NewApplicationBinding returns a new *ApplicationBinding
func NewApplicationBinding() *ApplicationBinding {

	return &ApplicationBinding{
		ReadOnly: false,
	}
}

// Identity returns the Identity of the object.
func (o *ApplicationBinding) Identity() bambou.Identity {

	return ApplicationBindingIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *ApplicationBinding) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *ApplicationBinding) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the ApplicationBinding from the server
func (o *ApplicationBinding) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the ApplicationBinding into the server
func (o *ApplicationBinding) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the ApplicationBinding from the server
func (o *ApplicationBinding) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}
