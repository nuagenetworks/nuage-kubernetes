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

// PSPATMapIdentity represents the Identity of the object
var PSPATMapIdentity = bambou.Identity{
	Name:     "pspatmap",
	Category: "pspatmaps",
}

// PSPATMapsList represents a list of PSPATMaps
type PSPATMapsList []*PSPATMap

// PSPATMapsAncestor is the interface that an ancestor of a PSPATMap must implement.
// An Ancestor is defined as an entity that has PSPATMap as a descendant.
// An Ancestor can get a list of its child PSPATMaps, but not necessarily create one.
type PSPATMapsAncestor interface {
	PSPATMaps(*bambou.FetchingInfo) (PSPATMapsList, *bambou.Error)
}

// PSPATMapsParent is the interface that a parent of a PSPATMap must implement.
// A Parent is defined as an entity that has PSPATMap as a child.
// A Parent is an Ancestor which can create a PSPATMap.
type PSPATMapsParent interface {
	PSPATMapsAncestor
	CreatePSPATMap(*PSPATMap) *bambou.Error
}

// PSPATMap represents the model of a pspatmap
type PSPATMap struct {
	ID                          string        `json:"ID,omitempty"`
	ParentID                    string        `json:"parentID,omitempty"`
	ParentType                  string        `json:"parentType,omitempty"`
	Owner                       string        `json:"owner,omitempty"`
	Name                        string        `json:"name,omitempty"`
	ReservedSPATIPs             []interface{} `json:"reservedSPATIPs,omitempty"`
	AssociatedSPATSourcesPoolID string        `json:"associatedSPATSourcesPoolID,omitempty"`
}

// NewPSPATMap returns a new *PSPATMap
func NewPSPATMap() *PSPATMap {

	return &PSPATMap{}
}

// Identity returns the Identity of the object.
func (o *PSPATMap) Identity() bambou.Identity {

	return PSPATMapIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *PSPATMap) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *PSPATMap) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the PSPATMap from the server
func (o *PSPATMap) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the PSPATMap into the server
func (o *PSPATMap) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the PSPATMap from the server
func (o *PSPATMap) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}
