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

// IngressAdvFwdEntryTemplateIdentity represents the Identity of the object
var IngressAdvFwdEntryTemplateIdentity = bambou.Identity{
	Name:     "ingressadvfwdentrytemplate",
	Category: "ingressadvfwdentrytemplates",
}

// IngressAdvFwdEntryTemplatesList represents a list of IngressAdvFwdEntryTemplates
type IngressAdvFwdEntryTemplatesList []*IngressAdvFwdEntryTemplate

// IngressAdvFwdEntryTemplatesAncestor is the interface that an ancestor of a IngressAdvFwdEntryTemplate must implement.
// An Ancestor is defined as an entity that has IngressAdvFwdEntryTemplate as a descendant.
// An Ancestor can get a list of its child IngressAdvFwdEntryTemplates, but not necessarily create one.
type IngressAdvFwdEntryTemplatesAncestor interface {
	IngressAdvFwdEntryTemplates(*bambou.FetchingInfo) (IngressAdvFwdEntryTemplatesList, *bambou.Error)
}

// IngressAdvFwdEntryTemplatesParent is the interface that a parent of a IngressAdvFwdEntryTemplate must implement.
// A Parent is defined as an entity that has IngressAdvFwdEntryTemplate as a child.
// A Parent is an Ancestor which can create a IngressAdvFwdEntryTemplate.
type IngressAdvFwdEntryTemplatesParent interface {
	IngressAdvFwdEntryTemplatesAncestor
	CreateIngressAdvFwdEntryTemplate(*IngressAdvFwdEntryTemplate) *bambou.Error
}

// IngressAdvFwdEntryTemplate represents the model of a ingressadvfwdentrytemplate
type IngressAdvFwdEntryTemplate struct {
	ID                              string `json:"ID,omitempty"`
	ParentID                        string `json:"parentID,omitempty"`
	ParentType                      string `json:"parentType,omitempty"`
	Owner                           string `json:"owner,omitempty"`
	ACLTemplateName                 string `json:"ACLTemplateName,omitempty"`
	ICMPCode                        string `json:"ICMPCode,omitempty"`
	ICMPType                        string `json:"ICMPType,omitempty"`
	FCOverride                      string `json:"FCOverride,omitempty"`
	IPv6AddressOverride             string `json:"IPv6AddressOverride,omitempty"`
	DSCP                            string `json:"DSCP,omitempty"`
	Name                            string `json:"name,omitempty"`
	LastUpdatedBy                   string `json:"lastUpdatedBy,omitempty"`
	Action                          string `json:"action,omitempty"`
	AddressOverride                 string `json:"addressOverride,omitempty"`
	RedirectVPortTagID              string `json:"redirectVPortTagID,omitempty"`
	Description                     string `json:"description,omitempty"`
	DestinationPort                 string `json:"destinationPort,omitempty"`
	NetworkID                       string `json:"networkID,omitempty"`
	NetworkType                     string `json:"networkType,omitempty"`
	MirrorDestinationID             string `json:"mirrorDestinationID,omitempty"`
	FlowLoggingEnabled              bool   `json:"flowLoggingEnabled"`
	EnterpriseName                  string `json:"enterpriseName,omitempty"`
	EntityScope                     string `json:"entityScope,omitempty"`
	LocationID                      string `json:"locationID,omitempty"`
	LocationType                    string `json:"locationType,omitempty"`
	PolicyState                     string `json:"policyState,omitempty"`
	DomainName                      string `json:"domainName,omitempty"`
	SourcePort                      string `json:"sourcePort,omitempty"`
	UplinkPreference                string `json:"uplinkPreference,omitempty"`
	Priority                        int    `json:"priority,omitempty"`
	Protocol                        string `json:"protocol,omitempty"`
	AssociatedApplicationID         string `json:"associatedApplicationID,omitempty"`
	AssociatedApplicationObjectID   string `json:"associatedApplicationObjectID,omitempty"`
	AssociatedApplicationObjectType string `json:"associatedApplicationObjectType,omitempty"`
	AssociatedLiveEntityID          string `json:"associatedLiveEntityID,omitempty"`
	StatsID                         string `json:"statsID,omitempty"`
	StatsLoggingEnabled             bool   `json:"statsLoggingEnabled"`
	EtherType                       string `json:"etherType,omitempty"`
	ExternalID                      string `json:"externalID,omitempty"`
}

// NewIngressAdvFwdEntryTemplate returns a new *IngressAdvFwdEntryTemplate
func NewIngressAdvFwdEntryTemplate() *IngressAdvFwdEntryTemplate {

	return &IngressAdvFwdEntryTemplate{
		DSCP:         "*",
		Action:       "FORWARD",
		NetworkType:  "ANY",
		LocationType: "ANY",
		Protocol:     "6",
		EtherType:    "0x0800",
	}
}

// Identity returns the Identity of the object.
func (o *IngressAdvFwdEntryTemplate) Identity() bambou.Identity {

	return IngressAdvFwdEntryTemplateIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *IngressAdvFwdEntryTemplate) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *IngressAdvFwdEntryTemplate) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the IngressAdvFwdEntryTemplate from the server
func (o *IngressAdvFwdEntryTemplate) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the IngressAdvFwdEntryTemplate into the server
func (o *IngressAdvFwdEntryTemplate) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the IngressAdvFwdEntryTemplate from the server
func (o *IngressAdvFwdEntryTemplate) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the IngressAdvFwdEntryTemplate
func (o *IngressAdvFwdEntryTemplate) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the IngressAdvFwdEntryTemplate
func (o *IngressAdvFwdEntryTemplate) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the IngressAdvFwdEntryTemplate
func (o *IngressAdvFwdEntryTemplate) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the IngressAdvFwdEntryTemplate
func (o *IngressAdvFwdEntryTemplate) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CreateJob creates a new child Job under the IngressAdvFwdEntryTemplate
func (o *IngressAdvFwdEntryTemplate) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Statistics retrieves the list of child Statistics of the IngressAdvFwdEntryTemplate
func (o *IngressAdvFwdEntryTemplate) Statistics(info *bambou.FetchingInfo) (StatisticsList, *bambou.Error) {

	var list StatisticsList
	err := bambou.CurrentSession().FetchChildren(o, StatisticsIdentity, &list, info)
	return list, err
}
