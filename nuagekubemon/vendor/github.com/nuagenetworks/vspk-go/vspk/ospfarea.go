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

// OSPFAreaIdentity represents the Identity of the object
var OSPFAreaIdentity = bambou.Identity{
	Name:     "ospfarea",
	Category: "ospfareas",
}

// OSPFAreasList represents a list of OSPFAreas
type OSPFAreasList []*OSPFArea

// OSPFAreasAncestor is the interface that an ancestor of a OSPFArea must implement.
// An Ancestor is defined as an entity that has OSPFArea as a descendant.
// An Ancestor can get a list of its child OSPFAreas, but not necessarily create one.
type OSPFAreasAncestor interface {
	OSPFAreas(*bambou.FetchingInfo) (OSPFAreasList, *bambou.Error)
}

// OSPFAreasParent is the interface that a parent of a OSPFArea must implement.
// A Parent is defined as an entity that has OSPFArea as a child.
// A Parent is an Ancestor which can create a OSPFArea.
type OSPFAreasParent interface {
	OSPFAreasAncestor
	CreateOSPFArea(*OSPFArea) *bambou.Error
}

// OSPFArea represents the model of a ospfarea
type OSPFArea struct {
	ID                          string        `json:"ID,omitempty"`
	ParentID                    string        `json:"parentID,omitempty"`
	ParentType                  string        `json:"parentType,omitempty"`
	Owner                       string        `json:"owner,omitempty"`
	LastUpdatedBy               string        `json:"lastUpdatedBy,omitempty"`
	RedistributeExternalEnabled bool          `json:"redistributeExternalEnabled"`
	DefaultMetric               int           `json:"defaultMetric,omitempty"`
	DefaultOriginateOption      string        `json:"defaultOriginateOption,omitempty"`
	Description                 string        `json:"description,omitempty"`
	AggregateAreaRange          []interface{} `json:"aggregateAreaRange,omitempty"`
	AggregateAreaRangeNSSA      []interface{} `json:"aggregateAreaRangeNSSA,omitempty"`
	EmbeddedMetadata            []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope                 string        `json:"entityScope,omitempty"`
	AreaID                      int           `json:"areaID,omitempty"`
	AreaType                    string        `json:"areaType,omitempty"`
	SummariesEnabled            bool          `json:"summariesEnabled"`
	SuppressAreaRange           []interface{} `json:"suppressAreaRange,omitempty"`
	SuppressAreaRangeNSSA       []interface{} `json:"suppressAreaRangeNSSA,omitempty"`
	ExternalID                  string        `json:"externalID,omitempty"`
}

// NewOSPFArea returns a new *OSPFArea
func NewOSPFArea() *OSPFArea {

	return &OSPFArea{
		AreaType: "NORMAL",
	}
}

// Identity returns the Identity of the object.
func (o *OSPFArea) Identity() bambou.Identity {

	return OSPFAreaIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *OSPFArea) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *OSPFArea) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the OSPFArea from the server
func (o *OSPFArea) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the OSPFArea into the server
func (o *OSPFArea) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the OSPFArea from the server
func (o *OSPFArea) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the OSPFArea
func (o *OSPFArea) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the OSPFArea
func (o *OSPFArea) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the OSPFArea
func (o *OSPFArea) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the OSPFArea
func (o *OSPFArea) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// OSPFInterfaces retrieves the list of child OSPFInterfaces of the OSPFArea
func (o *OSPFArea) OSPFInterfaces(info *bambou.FetchingInfo) (OSPFInterfacesList, *bambou.Error) {

	var list OSPFInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, OSPFInterfaceIdentity, &list, info)
	return list, err
}

// CreateOSPFInterface creates a new child OSPFInterface under the OSPFArea
func (o *OSPFArea) CreateOSPFInterface(child *OSPFInterface) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
