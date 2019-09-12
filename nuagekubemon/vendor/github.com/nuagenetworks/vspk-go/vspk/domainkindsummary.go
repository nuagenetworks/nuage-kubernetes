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

// DomainKindSummaryIdentity represents the Identity of the object
var DomainKindSummaryIdentity = bambou.Identity{
	Name:     "domainkindsummary",
	Category: "domainkindsummaries",
}

// DomainKindSummariesList represents a list of DomainKindSummaries
type DomainKindSummariesList []*DomainKindSummary

// DomainKindSummariesAncestor is the interface that an ancestor of a DomainKindSummary must implement.
// An Ancestor is defined as an entity that has DomainKindSummary as a descendant.
// An Ancestor can get a list of its child DomainKindSummaries, but not necessarily create one.
type DomainKindSummariesAncestor interface {
	DomainKindSummaries(*bambou.FetchingInfo) (DomainKindSummariesList, *bambou.Error)
}

// DomainKindSummariesParent is the interface that a parent of a DomainKindSummary must implement.
// A Parent is defined as an entity that has DomainKindSummary as a child.
// A Parent is an Ancestor which can create a DomainKindSummary.
type DomainKindSummariesParent interface {
	DomainKindSummariesAncestor
	CreateDomainKindSummary(*DomainKindSummary) *bambou.Error
}

// DomainKindSummary represents the model of a domainkindsummary
type DomainKindSummary struct {
	ID                    string        `json:"ID,omitempty"`
	ParentID              string        `json:"parentID,omitempty"`
	ParentType            string        `json:"parentType,omitempty"`
	Owner                 string        `json:"owner,omitempty"`
	MajorAlarmsCount      int           `json:"majorAlarmsCount,omitempty"`
	LastUpdatedBy         string        `json:"lastUpdatedBy,omitempty"`
	GatewayCount          int           `json:"gatewayCount,omitempty"`
	MeshGroupCount        int           `json:"meshGroupCount,omitempty"`
	MinorAlarmsCount      int           `json:"minorAlarmsCount,omitempty"`
	EmbeddedMetadata      []interface{} `json:"embeddedMetadata,omitempty"`
	InfoAlarmsCount       int           `json:"infoAlarmsCount,omitempty"`
	EntityScope           string        `json:"entityScope,omitempty"`
	DomainKindDescription string        `json:"domainKindDescription,omitempty"`
	DomainKindName        string        `json:"domainKindName,omitempty"`
	ZoneCount             int           `json:"zoneCount,omitempty"`
	TrafficVolume         int           `json:"trafficVolume,omitempty"`
	CriticalAlarmsCount   int           `json:"criticalAlarmsCount,omitempty"`
	NsgCount              int           `json:"nsgCount,omitempty"`
	SubNetworkCount       int           `json:"subNetworkCount,omitempty"`
	ExternalID            string        `json:"externalID,omitempty"`
}

// NewDomainKindSummary returns a new *DomainKindSummary
func NewDomainKindSummary() *DomainKindSummary {

	return &DomainKindSummary{}
}

// Identity returns the Identity of the object.
func (o *DomainKindSummary) Identity() bambou.Identity {

	return DomainKindSummaryIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *DomainKindSummary) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *DomainKindSummary) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the DomainKindSummary from the server
func (o *DomainKindSummary) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the DomainKindSummary into the server
func (o *DomainKindSummary) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the DomainKindSummary from the server
func (o *DomainKindSummary) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the DomainKindSummary
func (o *DomainKindSummary) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the DomainKindSummary
func (o *DomainKindSummary) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the DomainKindSummary
func (o *DomainKindSummary) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the DomainKindSummary
func (o *DomainKindSummary) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
