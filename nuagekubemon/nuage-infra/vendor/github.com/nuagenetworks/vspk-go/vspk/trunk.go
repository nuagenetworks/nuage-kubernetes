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

// TrunkIdentity represents the Identity of the object
var TrunkIdentity = bambou.Identity{
	Name:     "trunk",
	Category: "trunks",
}

// TrunksList represents a list of Trunks
type TrunksList []*Trunk

// TrunksAncestor is the interface that an ancestor of a Trunk must implement.
// An Ancestor is defined as an entity that has Trunk as a descendant.
// An Ancestor can get a list of its child Trunks, but not necessarily create one.
type TrunksAncestor interface {
	Trunks(*bambou.FetchingInfo) (TrunksList, *bambou.Error)
}

// TrunksParent is the interface that a parent of a Trunk must implement.
// A Parent is defined as an entity that has Trunk as a child.
// A Parent is an Ancestor which can create a Trunk.
type TrunksParent interface {
	TrunksAncestor
	CreateTrunk(*Trunk) *bambou.Error
}

// Trunk represents the model of a trunk
type Trunk struct {
	ID                string `json:"ID,omitempty"`
	ParentID          string `json:"parentID,omitempty"`
	ParentType        string `json:"parentType,omitempty"`
	Owner             string `json:"owner,omitempty"`
	Name              string `json:"name,omitempty"`
	AssociatedVPortID string `json:"associatedVPortID,omitempty"`
}

// NewTrunk returns a new *Trunk
func NewTrunk() *Trunk {

	return &Trunk{}
}

// Identity returns the Identity of the object.
func (o *Trunk) Identity() bambou.Identity {

	return TrunkIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Trunk) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Trunk) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Trunk from the server
func (o *Trunk) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Trunk into the server
func (o *Trunk) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Trunk from the server
func (o *Trunk) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// VPorts retrieves the list of child VPorts of the Trunk
func (o *Trunk) VPorts(info *bambou.FetchingInfo) (VPortsList, *bambou.Error) {

	var list VPortsList
	err := bambou.CurrentSession().FetchChildren(o, VPortIdentity, &list, info)
	return list, err
}
