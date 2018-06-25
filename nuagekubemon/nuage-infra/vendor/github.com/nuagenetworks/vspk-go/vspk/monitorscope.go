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

// MonitorscopeIdentity represents the Identity of the object
var MonitorscopeIdentity = bambou.Identity{
	Name:     "monitorscope",
	Category: "monitorscopes",
}

// MonitorscopesList represents a list of Monitorscopes
type MonitorscopesList []*Monitorscope

// MonitorscopesAncestor is the interface that an ancestor of a Monitorscope must implement.
// An Ancestor is defined as an entity that has Monitorscope as a descendant.
// An Ancestor can get a list of its child Monitorscopes, but not necessarily create one.
type MonitorscopesAncestor interface {
	Monitorscopes(*bambou.FetchingInfo) (MonitorscopesList, *bambou.Error)
}

// MonitorscopesParent is the interface that a parent of a Monitorscope must implement.
// A Parent is defined as an entity that has Monitorscope as a child.
// A Parent is an Ancestor which can create a Monitorscope.
type MonitorscopesParent interface {
	MonitorscopesAncestor
	CreateMonitorscope(*Monitorscope) *bambou.Error
}

// Monitorscope represents the model of a monitorscope
type Monitorscope struct {
	ID                      string        `json:"ID,omitempty"`
	ParentID                string        `json:"parentID,omitempty"`
	ParentType              string        `json:"parentType,omitempty"`
	Owner                   string        `json:"owner,omitempty"`
	Name                    string        `json:"name,omitempty"`
	ReadOnly                bool          `json:"readOnly"`
	DestinationNSGs         []interface{} `json:"destinationNSGs,omitempty"`
	AllowAllDestinationNSGs bool          `json:"allowAllDestinationNSGs"`
	AllowAllSourceNSGs      bool          `json:"allowAllSourceNSGs"`
	SourceNSGs              []interface{} `json:"sourceNSGs,omitempty"`
}

// NewMonitorscope returns a new *Monitorscope
func NewMonitorscope() *Monitorscope {

	return &Monitorscope{
		ReadOnly: false,
	}
}

// Identity returns the Identity of the object.
func (o *Monitorscope) Identity() bambou.Identity {

	return MonitorscopeIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Monitorscope) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Monitorscope) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Monitorscope from the server
func (o *Monitorscope) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Monitorscope into the server
func (o *Monitorscope) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Monitorscope from the server
func (o *Monitorscope) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}
