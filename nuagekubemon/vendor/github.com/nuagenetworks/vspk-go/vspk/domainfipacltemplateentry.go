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

// DomainFIPAclTemplateEntryIdentity represents the Identity of the object
var DomainFIPAclTemplateEntryIdentity = bambou.Identity{
	Name:     "egressdomainfloatingipaclentrytemplate",
	Category: "egressdomainfloatingipaclentrytemplates",
}

// DomainFIPAclTemplateEntriesList represents a list of DomainFIPAclTemplateEntries
type DomainFIPAclTemplateEntriesList []*DomainFIPAclTemplateEntry

// DomainFIPAclTemplateEntriesAncestor is the interface that an ancestor of a DomainFIPAclTemplateEntry must implement.
// An Ancestor is defined as an entity that has DomainFIPAclTemplateEntry as a descendant.
// An Ancestor can get a list of its child DomainFIPAclTemplateEntries, but not necessarily create one.
type DomainFIPAclTemplateEntriesAncestor interface {
	DomainFIPAclTemplateEntries(*bambou.FetchingInfo) (DomainFIPAclTemplateEntriesList, *bambou.Error)
}

// DomainFIPAclTemplateEntriesParent is the interface that a parent of a DomainFIPAclTemplateEntry must implement.
// A Parent is defined as an entity that has DomainFIPAclTemplateEntry as a child.
// A Parent is an Ancestor which can create a DomainFIPAclTemplateEntry.
type DomainFIPAclTemplateEntriesParent interface {
	DomainFIPAclTemplateEntriesAncestor
	CreateDomainFIPAclTemplateEntry(*DomainFIPAclTemplateEntry) *bambou.Error
}

// DomainFIPAclTemplateEntry represents the model of a egressdomainfloatingipaclentrytemplate
type DomainFIPAclTemplateEntry struct {
	ID                              string      `json:"ID,omitempty"`
	ParentID                        string      `json:"parentID,omitempty"`
	ParentType                      string      `json:"parentType,omitempty"`
	Owner                           string      `json:"owner,omitempty"`
	ACLTemplateName                 string      `json:"ACLTemplateName,omitempty"`
	ICMPCode                        string      `json:"ICMPCode,omitempty"`
	ICMPType                        string      `json:"ICMPType,omitempty"`
	DSCP                            string      `json:"DSCP,omitempty"`
	LastUpdatedBy                   string      `json:"lastUpdatedBy,omitempty"`
	Action                          string      `json:"action,omitempty"`
	ActionDetails                   interface{} `json:"actionDetails,omitempty"`
	AddressOverride                 string      `json:"addressOverride,omitempty"`
	Reflexive                       bool        `json:"reflexive"`
	Description                     string      `json:"description,omitempty"`
	DestPgId                        string      `json:"destPgId,omitempty"`
	DestPgType                      string      `json:"destPgType,omitempty"`
	DestinationPort                 string      `json:"destinationPort,omitempty"`
	DestinationType                 string      `json:"destinationType,omitempty"`
	DestinationValue                string      `json:"destinationValue,omitempty"`
	NetworkID                       string      `json:"networkID,omitempty"`
	NetworkType                     string      `json:"networkType,omitempty"`
	MirrorDestinationID             string      `json:"mirrorDestinationID,omitempty"`
	FlowLoggingEnabled              bool        `json:"flowLoggingEnabled"`
	EnterpriseName                  string      `json:"enterpriseName,omitempty"`
	EntityScope                     string      `json:"entityScope,omitempty"`
	LocationID                      string      `json:"locationID,omitempty"`
	LocationType                    string      `json:"locationType,omitempty"`
	PolicyState                     string      `json:"policyState,omitempty"`
	DomainName                      string      `json:"domainName,omitempty"`
	SourcePgId                      string      `json:"sourcePgId,omitempty"`
	SourcePgType                    string      `json:"sourcePgType,omitempty"`
	SourcePort                      string      `json:"sourcePort,omitempty"`
	SourceType                      string      `json:"sourceType,omitempty"`
	SourceValue                     string      `json:"sourceValue,omitempty"`
	Priority                        int         `json:"priority,omitempty"`
	Protocol                        string      `json:"protocol,omitempty"`
	AssociatedApplicationID         string      `json:"associatedApplicationID,omitempty"`
	AssociatedApplicationObjectID   string      `json:"associatedApplicationObjectID,omitempty"`
	AssociatedApplicationObjectType string      `json:"associatedApplicationObjectType,omitempty"`
	AssociatedLiveEntityID          string      `json:"associatedLiveEntityID,omitempty"`
	Stateful                        bool        `json:"stateful"`
	StatsID                         string      `json:"statsID,omitempty"`
	StatsLoggingEnabled             bool        `json:"statsLoggingEnabled"`
	EtherType                       string      `json:"etherType,omitempty"`
	ExternalID                      string      `json:"externalID,omitempty"`
}

// NewDomainFIPAclTemplateEntry returns a new *DomainFIPAclTemplateEntry
func NewDomainFIPAclTemplateEntry() *DomainFIPAclTemplateEntry {

	return &DomainFIPAclTemplateEntry{}
}

// Identity returns the Identity of the object.
func (o *DomainFIPAclTemplateEntry) Identity() bambou.Identity {

	return DomainFIPAclTemplateEntryIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *DomainFIPAclTemplateEntry) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *DomainFIPAclTemplateEntry) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the DomainFIPAclTemplateEntry from the server
func (o *DomainFIPAclTemplateEntry) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the DomainFIPAclTemplateEntry into the server
func (o *DomainFIPAclTemplateEntry) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the DomainFIPAclTemplateEntry from the server
func (o *DomainFIPAclTemplateEntry) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the DomainFIPAclTemplateEntry
func (o *DomainFIPAclTemplateEntry) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the DomainFIPAclTemplateEntry
func (o *DomainFIPAclTemplateEntry) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the DomainFIPAclTemplateEntry
func (o *DomainFIPAclTemplateEntry) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the DomainFIPAclTemplateEntry
func (o *DomainFIPAclTemplateEntry) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
