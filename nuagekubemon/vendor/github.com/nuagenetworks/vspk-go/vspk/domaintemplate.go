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

// DomainTemplateIdentity represents the Identity of the object
var DomainTemplateIdentity = bambou.Identity{
	Name:     "domaintemplate",
	Category: "domaintemplates",
}

// DomainTemplatesList represents a list of DomainTemplates
type DomainTemplatesList []*DomainTemplate

// DomainTemplatesAncestor is the interface that an ancestor of a DomainTemplate must implement.
// An Ancestor is defined as an entity that has DomainTemplate as a descendant.
// An Ancestor can get a list of its child DomainTemplates, but not necessarily create one.
type DomainTemplatesAncestor interface {
	DomainTemplates(*bambou.FetchingInfo) (DomainTemplatesList, *bambou.Error)
}

// DomainTemplatesParent is the interface that a parent of a DomainTemplate must implement.
// A Parent is defined as an entity that has DomainTemplate as a child.
// A Parent is an Ancestor which can create a DomainTemplate.
type DomainTemplatesParent interface {
	DomainTemplatesAncestor
	CreateDomainTemplate(*DomainTemplate) *bambou.Error
}

// DomainTemplate represents the model of a domaintemplate
type DomainTemplate struct {
	ID                              string        `json:"ID,omitempty"`
	ParentID                        string        `json:"parentID,omitempty"`
	ParentType                      string        `json:"parentType,omitempty"`
	Owner                           string        `json:"owner,omitempty"`
	DPI                             string        `json:"DPI,omitempty"`
	Name                            string        `json:"name,omitempty"`
	LastUpdatedBy                   string        `json:"lastUpdatedBy,omitempty"`
	Description                     string        `json:"description,omitempty"`
	EmbeddedMetadata                []interface{} `json:"embeddedMetadata,omitempty"`
	Encryption                      string        `json:"encryption,omitempty"`
	EntityScope                     string        `json:"entityScope,omitempty"`
	PolicyChangeStatus              string        `json:"policyChangeStatus,omitempty"`
	AssociatedBGPProfileID          string        `json:"associatedBGPProfileID,omitempty"`
	AssociatedMulticastChannelMapID string        `json:"associatedMulticastChannelMapID,omitempty"`
	AssociatedPATMapperID           string        `json:"associatedPATMapperID,omitempty"`
	Multicast                       string        `json:"multicast,omitempty"`
	ExternalID                      string        `json:"externalID,omitempty"`
}

// NewDomainTemplate returns a new *DomainTemplate
func NewDomainTemplate() *DomainTemplate {

	return &DomainTemplate{
		DPI: "DISABLED",
	}
}

// Identity returns the Identity of the object.
func (o *DomainTemplate) Identity() bambou.Identity {

	return DomainTemplateIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *DomainTemplate) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *DomainTemplate) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the DomainTemplate from the server
func (o *DomainTemplate) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the DomainTemplate into the server
func (o *DomainTemplate) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the DomainTemplate from the server
func (o *DomainTemplate) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// RedirectionTargetTemplates retrieves the list of child RedirectionTargetTemplates of the DomainTemplate
func (o *DomainTemplate) RedirectionTargetTemplates(info *bambou.FetchingInfo) (RedirectionTargetTemplatesList, *bambou.Error) {

	var list RedirectionTargetTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, RedirectionTargetTemplateIdentity, &list, info)
	return list, err
}

// CreateRedirectionTargetTemplate creates a new child RedirectionTargetTemplate under the DomainTemplate
func (o *DomainTemplate) CreateRedirectionTargetTemplate(child *RedirectionTargetTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Permissions retrieves the list of child Permissions of the DomainTemplate
func (o *DomainTemplate) Permissions(info *bambou.FetchingInfo) (PermissionsList, *bambou.Error) {

	var list PermissionsList
	err := bambou.CurrentSession().FetchChildren(o, PermissionIdentity, &list, info)
	return list, err
}

// CreatePermission creates a new child Permission under the DomainTemplate
func (o *DomainTemplate) CreatePermission(child *Permission) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the DomainTemplate
func (o *DomainTemplate) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the DomainTemplate
func (o *DomainTemplate) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PGExpressionTemplates retrieves the list of child PGExpressionTemplates of the DomainTemplate
func (o *DomainTemplate) PGExpressionTemplates(info *bambou.FetchingInfo) (PGExpressionTemplatesList, *bambou.Error) {

	var list PGExpressionTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, PGExpressionTemplateIdentity, &list, info)
	return list, err
}

// CreatePGExpressionTemplate creates a new child PGExpressionTemplate under the DomainTemplate
func (o *DomainTemplate) CreatePGExpressionTemplate(child *PGExpressionTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EgressACLTemplates retrieves the list of child EgressACLTemplates of the DomainTemplate
func (o *DomainTemplate) EgressACLTemplates(info *bambou.FetchingInfo) (EgressACLTemplatesList, *bambou.Error) {

	var list EgressACLTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, EgressACLTemplateIdentity, &list, info)
	return list, err
}

// CreateEgressACLTemplate creates a new child EgressACLTemplate under the DomainTemplate
func (o *DomainTemplate) CreateEgressACLTemplate(child *EgressACLTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EgressAdvFwdTemplates retrieves the list of child EgressAdvFwdTemplates of the DomainTemplate
func (o *DomainTemplate) EgressAdvFwdTemplates(info *bambou.FetchingInfo) (EgressAdvFwdTemplatesList, *bambou.Error) {

	var list EgressAdvFwdTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, EgressAdvFwdTemplateIdentity, &list, info)
	return list, err
}

// CreateEgressAdvFwdTemplate creates a new child EgressAdvFwdTemplate under the DomainTemplate
func (o *DomainTemplate) CreateEgressAdvFwdTemplate(child *EgressAdvFwdTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DomainFIPAclTemplates retrieves the list of child DomainFIPAclTemplates of the DomainTemplate
func (o *DomainTemplate) DomainFIPAclTemplates(info *bambou.FetchingInfo) (DomainFIPAclTemplatesList, *bambou.Error) {

	var list DomainFIPAclTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, DomainFIPAclTemplateIdentity, &list, info)
	return list, err
}

// CreateDomainFIPAclTemplate creates a new child DomainFIPAclTemplate under the DomainTemplate
func (o *DomainTemplate) CreateDomainFIPAclTemplate(child *DomainFIPAclTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VirtualFirewallPolicies retrieves the list of child VirtualFirewallPolicies of the DomainTemplate
func (o *DomainTemplate) VirtualFirewallPolicies(info *bambou.FetchingInfo) (VirtualFirewallPoliciesList, *bambou.Error) {

	var list VirtualFirewallPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, VirtualFirewallPolicyIdentity, &list, info)
	return list, err
}

// CreateVirtualFirewallPolicy creates a new child VirtualFirewallPolicy under the DomainTemplate
func (o *DomainTemplate) CreateVirtualFirewallPolicy(child *VirtualFirewallPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the DomainTemplate
func (o *DomainTemplate) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the DomainTemplate
func (o *DomainTemplate) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IngressACLTemplates retrieves the list of child IngressACLTemplates of the DomainTemplate
func (o *DomainTemplate) IngressACLTemplates(info *bambou.FetchingInfo) (IngressACLTemplatesList, *bambou.Error) {

	var list IngressACLTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressACLTemplateIdentity, &list, info)
	return list, err
}

// CreateIngressACLTemplate creates a new child IngressACLTemplate under the DomainTemplate
func (o *DomainTemplate) CreateIngressACLTemplate(child *IngressACLTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IngressAdvFwdTemplates retrieves the list of child IngressAdvFwdTemplates of the DomainTemplate
func (o *DomainTemplate) IngressAdvFwdTemplates(info *bambou.FetchingInfo) (IngressAdvFwdTemplatesList, *bambou.Error) {

	var list IngressAdvFwdTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressAdvFwdTemplateIdentity, &list, info)
	return list, err
}

// CreateIngressAdvFwdTemplate creates a new child IngressAdvFwdTemplate under the DomainTemplate
func (o *DomainTemplate) CreateIngressAdvFwdTemplate(child *IngressAdvFwdTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CreateJob creates a new child Job under the DomainTemplate
func (o *DomainTemplate) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PolicyGroupTemplates retrieves the list of child PolicyGroupTemplates of the DomainTemplate
func (o *DomainTemplate) PolicyGroupTemplates(info *bambou.FetchingInfo) (PolicyGroupTemplatesList, *bambou.Error) {

	var list PolicyGroupTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, PolicyGroupTemplateIdentity, &list, info)
	return list, err
}

// CreatePolicyGroupTemplate creates a new child PolicyGroupTemplate under the DomainTemplate
func (o *DomainTemplate) CreatePolicyGroupTemplate(child *PolicyGroupTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Domains retrieves the list of child Domains of the DomainTemplate
func (o *DomainTemplate) Domains(info *bambou.FetchingInfo) (DomainsList, *bambou.Error) {

	var list DomainsList
	err := bambou.CurrentSession().FetchChildren(o, DomainIdentity, &list, info)
	return list, err
}

// AssignDomains assigns the list of Domains to the DomainTemplate
func (o *DomainTemplate) AssignDomains(children DomainsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, DomainIdentity)
}

// ZoneTemplates retrieves the list of child ZoneTemplates of the DomainTemplate
func (o *DomainTemplate) ZoneTemplates(info *bambou.FetchingInfo) (ZoneTemplatesList, *bambou.Error) {

	var list ZoneTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, ZoneTemplateIdentity, &list, info)
	return list, err
}

// CreateZoneTemplate creates a new child ZoneTemplate under the DomainTemplate
func (o *DomainTemplate) CreateZoneTemplate(child *ZoneTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// QOSs retrieves the list of child QOSs of the DomainTemplate
func (o *DomainTemplate) QOSs(info *bambou.FetchingInfo) (QOSsList, *bambou.Error) {

	var list QOSsList
	err := bambou.CurrentSession().FetchChildren(o, QOSIdentity, &list, info)
	return list, err
}

// CreateQOS creates a new child QOS under the DomainTemplate
func (o *DomainTemplate) CreateQOS(child *QOS) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Groups retrieves the list of child Groups of the DomainTemplate
func (o *DomainTemplate) Groups(info *bambou.FetchingInfo) (GroupsList, *bambou.Error) {

	var list GroupsList
	err := bambou.CurrentSession().FetchChildren(o, GroupIdentity, &list, info)
	return list, err
}

// SubnetTemplates retrieves the list of child SubnetTemplates of the DomainTemplate
func (o *DomainTemplate) SubnetTemplates(info *bambou.FetchingInfo) (SubnetTemplatesList, *bambou.Error) {

	var list SubnetTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, SubnetTemplateIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the DomainTemplate
func (o *DomainTemplate) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
