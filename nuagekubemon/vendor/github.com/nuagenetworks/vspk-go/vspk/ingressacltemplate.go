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

// IngressACLTemplateIdentity represents the Identity of the object
var IngressACLTemplateIdentity = bambou.Identity{
	Name:     "ingressacltemplate",
	Category: "ingressacltemplates",
}

// IngressACLTemplatesList represents a list of IngressACLTemplates
type IngressACLTemplatesList []*IngressACLTemplate

// IngressACLTemplatesAncestor is the interface that an ancestor of a IngressACLTemplate must implement.
// An Ancestor is defined as an entity that has IngressACLTemplate as a descendant.
// An Ancestor can get a list of its child IngressACLTemplates, but not necessarily create one.
type IngressACLTemplatesAncestor interface {
	IngressACLTemplates(*bambou.FetchingInfo) (IngressACLTemplatesList, *bambou.Error)
}

// IngressACLTemplatesParent is the interface that a parent of a IngressACLTemplate must implement.
// A Parent is defined as an entity that has IngressACLTemplate as a child.
// A Parent is an Ancestor which can create a IngressACLTemplate.
type IngressACLTemplatesParent interface {
	IngressACLTemplatesAncestor
	CreateIngressACLTemplate(*IngressACLTemplate) *bambou.Error
}

// IngressACLTemplate represents the model of a ingressacltemplate
type IngressACLTemplate struct {
	ID                     string `json:"ID,omitempty"`
	ParentID               string `json:"parentID,omitempty"`
	ParentType             string `json:"parentType,omitempty"`
	Owner                  string `json:"owner,omitempty"`
	Name                   string `json:"name,omitempty"`
	LastUpdatedBy          string `json:"lastUpdatedBy,omitempty"`
	Active                 bool   `json:"active"`
	DefaultAllowIP         bool   `json:"defaultAllowIP"`
	DefaultAllowNonIP      bool   `json:"defaultAllowNonIP"`
	Description            string `json:"description,omitempty"`
	AllowAddressSpoof      bool   `json:"allowAddressSpoof"`
	AllowL2AddressSpoof    bool   `json:"allowL2AddressSpoof"`
	EntityScope            string `json:"entityScope,omitempty"`
	PolicyState            string `json:"policyState,omitempty"`
	Priority               int    `json:"priority,omitempty"`
	PriorityType           string `json:"priorityType,omitempty"`
	AssocAclTemplateId     string `json:"assocAclTemplateId,omitempty"`
	AssociatedLiveEntityID string `json:"associatedLiveEntityID,omitempty"`
	AutoGeneratePriority   bool   `json:"autoGeneratePriority"`
	ExternalID             string `json:"externalID,omitempty"`
}

// NewIngressACLTemplate returns a new *IngressACLTemplate
func NewIngressACLTemplate() *IngressACLTemplate {

	return &IngressACLTemplate{}
}

// Identity returns the Identity of the object.
func (o *IngressACLTemplate) Identity() bambou.Identity {

	return IngressACLTemplateIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *IngressACLTemplate) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *IngressACLTemplate) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the IngressACLTemplate from the server
func (o *IngressACLTemplate) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the IngressACLTemplate into the server
func (o *IngressACLTemplate) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the IngressACLTemplate from the server
func (o *IngressACLTemplate) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the IngressACLTemplate
func (o *IngressACLTemplate) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the IngressACLTemplate
func (o *IngressACLTemplate) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the IngressACLTemplate
func (o *IngressACLTemplate) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the IngressACLTemplate
func (o *IngressACLTemplate) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMs retrieves the list of child VMs of the IngressACLTemplate
func (o *IngressACLTemplate) VMs(info *bambou.FetchingInfo) (VMsList, *bambou.Error) {

	var list VMsList
	err := bambou.CurrentSession().FetchChildren(o, VMIdentity, &list, info)
	return list, err
}

// IngressACLEntryTemplates retrieves the list of child IngressACLEntryTemplates of the IngressACLTemplate
func (o *IngressACLTemplate) IngressACLEntryTemplates(info *bambou.FetchingInfo) (IngressACLEntryTemplatesList, *bambou.Error) {

	var list IngressACLEntryTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressACLEntryTemplateIdentity, &list, info)
	return list, err
}

// CreateIngressACLEntryTemplate creates a new child IngressACLEntryTemplate under the IngressACLTemplate
func (o *IngressACLTemplate) CreateIngressACLEntryTemplate(child *IngressACLEntryTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CreateJob creates a new child Job under the IngressACLTemplate
func (o *IngressACLTemplate) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Containers retrieves the list of child Containers of the IngressACLTemplate
func (o *IngressACLTemplate) Containers(info *bambou.FetchingInfo) (ContainersList, *bambou.Error) {

	var list ContainersList
	err := bambou.CurrentSession().FetchChildren(o, ContainerIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the IngressACLTemplate
func (o *IngressACLTemplate) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
