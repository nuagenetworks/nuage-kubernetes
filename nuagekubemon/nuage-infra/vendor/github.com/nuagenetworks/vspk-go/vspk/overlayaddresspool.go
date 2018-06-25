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

// OverlayAddressPoolIdentity represents the Identity of the object
var OverlayAddressPoolIdentity = bambou.Identity{
	Name:     "overlayaddresspool",
	Category: "overlayaddresspools",
}

// OverlayAddressPoolsList represents a list of OverlayAddressPools
type OverlayAddressPoolsList []*OverlayAddressPool

// OverlayAddressPoolsAncestor is the interface that an ancestor of a OverlayAddressPool must implement.
// An Ancestor is defined as an entity that has OverlayAddressPool as a descendant.
// An Ancestor can get a list of its child OverlayAddressPools, but not necessarily create one.
type OverlayAddressPoolsAncestor interface {
	OverlayAddressPools(*bambou.FetchingInfo) (OverlayAddressPoolsList, *bambou.Error)
}

// OverlayAddressPoolsParent is the interface that a parent of a OverlayAddressPool must implement.
// A Parent is defined as an entity that has OverlayAddressPool as a child.
// A Parent is an Ancestor which can create a OverlayAddressPool.
type OverlayAddressPoolsParent interface {
	OverlayAddressPoolsAncestor
	CreateOverlayAddressPool(*OverlayAddressPool) *bambou.Error
}

// OverlayAddressPool represents the model of a overlayaddresspool
type OverlayAddressPool struct {
	ID                 string `json:"ID,omitempty"`
	ParentID           string `json:"parentID,omitempty"`
	ParentType         string `json:"parentType,omitempty"`
	Owner              string `json:"owner,omitempty"`
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	EndAddressRange    string `json:"endAddressRange,omitempty"`
	AssociatedDomainID string `json:"associatedDomainID,omitempty"`
	StartAddressRange  string `json:"startAddressRange,omitempty"`
}

// NewOverlayAddressPool returns a new *OverlayAddressPool
func NewOverlayAddressPool() *OverlayAddressPool {

	return &OverlayAddressPool{}
}

// Identity returns the Identity of the object.
func (o *OverlayAddressPool) Identity() bambou.Identity {

	return OverlayAddressPoolIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *OverlayAddressPool) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *OverlayAddressPool) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the OverlayAddressPool from the server
func (o *OverlayAddressPool) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the OverlayAddressPool into the server
func (o *OverlayAddressPool) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the OverlayAddressPool from the server
func (o *OverlayAddressPool) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// OverlayPATNATEntries retrieves the list of child OverlayPATNATEntries of the OverlayAddressPool
func (o *OverlayAddressPool) OverlayPATNATEntries(info *bambou.FetchingInfo) (OverlayPATNATEntriesList, *bambou.Error) {

	var list OverlayPATNATEntriesList
	err := bambou.CurrentSession().FetchChildren(o, OverlayPATNATEntryIdentity, &list, info)
	return list, err
}

// CreateOverlayPATNATEntry creates a new child OverlayPATNATEntry under the OverlayAddressPool
func (o *OverlayAddressPool) CreateOverlayPATNATEntry(child *OverlayPATNATEntry) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
