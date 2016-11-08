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

// PATMapperIdentity represents the Identity of the object
var PATMapperIdentity = bambou.Identity{
	Name:     "patmapper",
	Category: "patmappers",
}

// PATMappersList represents a list of PATMappers
type PATMappersList []*PATMapper

// PATMappersAncestor is the interface of an ancestor of a PATMapper must implement.
type PATMappersAncestor interface {
	PATMappers(*bambou.FetchingInfo) (PATMappersList, *bambou.Error)
	CreatePATMappers(*PATMapper) *bambou.Error
}

// PATMapper represents the model of a patmapper
type PATMapper struct {
	ID            string `json:"ID,omitempty"`
	ParentID      string `json:"parentID,omitempty"`
	ParentType    string `json:"parentType,omitempty"`
	Owner         string `json:"owner,omitempty"`
	Name          string `json:"name,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	Description   string `json:"description,omitempty"`
	EntityScope   string `json:"entityScope,omitempty"`
	ExternalID    string `json:"externalID,omitempty"`
}

// NewPATMapper returns a new *PATMapper
func NewPATMapper() *PATMapper {

	return &PATMapper{}
}

// Identity returns the Identity of the object.
func (o *PATMapper) Identity() bambou.Identity {

	return PATMapperIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *PATMapper) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *PATMapper) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the PATMapper from the server
func (o *PATMapper) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the PATMapper into the server
func (o *PATMapper) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the PATMapper from the server
func (o *PATMapper) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// SharedNetworkResources retrieves the list of child SharedNetworkResources of the PATMapper
func (o *PATMapper) SharedNetworkResources(info *bambou.FetchingInfo) (SharedNetworkResourcesList, *bambou.Error) {

	var list SharedNetworkResourcesList
	err := bambou.CurrentSession().FetchChildren(o, SharedNetworkResourceIdentity, &list, info)
	return list, err
}

// AssignSharedNetworkResources assigns the list of SharedNetworkResources to the PATMapper
func (o *PATMapper) AssignSharedNetworkResources(children SharedNetworkResourcesList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, SharedNetworkResourceIdentity)
}
