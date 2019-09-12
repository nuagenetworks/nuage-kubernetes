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

// MetadataIdentity represents the Identity of the object
var MetadataIdentity = bambou.Identity{
	Name:     "metadata",
	Category: "metadatas",
}

// MetadatasList represents a list of Metadatas
type MetadatasList []*Metadata

// MetadatasAncestor is the interface that an ancestor of a Metadata must implement.
// An Ancestor is defined as an entity that has Metadata as a descendant.
// An Ancestor can get a list of its child Metadatas, but not necessarily create one.
type MetadatasAncestor interface {
	Metadatas(*bambou.FetchingInfo) (MetadatasList, *bambou.Error)
}

// MetadatasParent is the interface that a parent of a Metadata must implement.
// A Parent is defined as an entity that has Metadata as a child.
// A Parent is an Ancestor which can create a Metadata.
type MetadatasParent interface {
	MetadatasAncestor
	CreateMetadata(*Metadata) *bambou.Error
}

// Metadata represents the model of a metadata
type Metadata struct {
	ID                          string        `json:"ID,omitempty"`
	ParentID                    string        `json:"parentID,omitempty"`
	ParentType                  string        `json:"parentType,omitempty"`
	Owner                       string        `json:"owner,omitempty"`
	Name                        string        `json:"name,omitempty"`
	LastUpdatedBy               string        `json:"lastUpdatedBy,omitempty"`
	Description                 string        `json:"description,omitempty"`
	MetadataTagIDs              []interface{} `json:"metadataTagIDs,omitempty"`
	NetworkNotificationDisabled bool          `json:"networkNotificationDisabled"`
	Blob                        string        `json:"blob,omitempty"`
	GlobalMetadata              bool          `json:"globalMetadata"`
	EntityScope                 string        `json:"entityScope,omitempty"`
	AssocEntityID               string        `json:"assocEntityID,omitempty"`
	AssocEntityType             string        `json:"assocEntityType,omitempty"`
	ExternalID                  string        `json:"externalID,omitempty"`
}

// NewMetadata returns a new *Metadata
func NewMetadata() *Metadata {

	return &Metadata{}
}

// Identity returns the Identity of the object.
func (o *Metadata) Identity() bambou.Identity {

	return MetadataIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Metadata) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Metadata) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Metadata from the server
func (o *Metadata) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Metadata into the server
func (o *Metadata) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Metadata from the server
func (o *Metadata) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// EventLogs retrieves the list of child EventLogs of the Metadata
func (o *Metadata) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
