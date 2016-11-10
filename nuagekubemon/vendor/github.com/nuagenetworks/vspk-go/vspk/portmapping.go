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

// PortMappingIdentity represents the Identity of the object
var PortMappingIdentity = bambou.Identity{
	Name:     "portmapping",
	Category: "portmappings",
}

// PortMappingsList represents a list of PortMappings
type PortMappingsList []*PortMapping

// PortMappingsAncestor is the interface of an ancestor of a PortMapping must implement.
type PortMappingsAncestor interface {
	PortMappings(*bambou.FetchingInfo) (PortMappingsList, *bambou.Error)
	CreatePortMappings(*PortMapping) *bambou.Error
}

// PortMapping represents the model of a portmapping
type PortMapping struct {
	ID            string `json:"ID,omitempty"`
	ParentID      string `json:"parentID,omitempty"`
	ParentType    string `json:"parentType,omitempty"`
	Owner         string `json:"owner,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	EntityScope   string `json:"entityScope,omitempty"`
	PrivatePort   string `json:"privatePort,omitempty"`
	PublicPort    string `json:"publicPort,omitempty"`
	ExternalID    string `json:"externalID,omitempty"`
}

// NewPortMapping returns a new *PortMapping
func NewPortMapping() *PortMapping {

	return &PortMapping{}
}

// Identity returns the Identity of the object.
func (o *PortMapping) Identity() bambou.Identity {

	return PortMappingIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *PortMapping) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *PortMapping) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the PortMapping from the server
func (o *PortMapping) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the PortMapping into the server
func (o *PortMapping) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the PortMapping from the server
func (o *PortMapping) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}
