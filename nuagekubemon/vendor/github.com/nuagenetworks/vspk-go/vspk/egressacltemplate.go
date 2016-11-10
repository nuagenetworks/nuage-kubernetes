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

// EgressACLTemplateIdentity represents the Identity of the object
var EgressACLTemplateIdentity = bambou.Identity{
	Name:     "egressacltemplate",
	Category: "egressacltemplates",
}

// EgressACLTemplatesList represents a list of EgressACLTemplates
type EgressACLTemplatesList []*EgressACLTemplate

// EgressACLTemplatesAncestor is the interface of an ancestor of a EgressACLTemplate must implement.
type EgressACLTemplatesAncestor interface {
	EgressACLTemplates(*bambou.FetchingInfo) (EgressACLTemplatesList, *bambou.Error)
	CreateEgressACLTemplates(*EgressACLTemplate) *bambou.Error
}

// EgressACLTemplate represents the model of a egressacltemplate
type EgressACLTemplate struct {
	ID                             string `json:"ID,omitempty"`
	ParentID                       string `json:"parentID,omitempty"`
	ParentType                     string `json:"parentType,omitempty"`
	Owner                          string `json:"owner,omitempty"`
	Name                           string `json:"name,omitempty"`
	LastUpdatedBy                  string `json:"lastUpdatedBy,omitempty"`
	Active                         bool   `json:"active"`
	DefaultAllowIP                 bool   `json:"defaultAllowIP"`
	DefaultAllowNonIP              bool   `json:"defaultAllowNonIP"`
	DefaultInstallACLImplicitRules bool   `json:"defaultInstallACLImplicitRules"`
	Description                    string `json:"description,omitempty"`
	EntityScope                    string `json:"entityScope,omitempty"`
	PolicyState                    string `json:"policyState,omitempty"`
	Priority                       int    `json:"priority,omitempty"`
	PriorityType                   string `json:"priorityType,omitempty"`
	AssociatedLiveEntityID         string `json:"associatedLiveEntityID,omitempty"`
	ExternalID                     string `json:"externalID,omitempty"`
}

// NewEgressACLTemplate returns a new *EgressACLTemplate
func NewEgressACLTemplate() *EgressACLTemplate {

	return &EgressACLTemplate{}
}

// Identity returns the Identity of the object.
func (o *EgressACLTemplate) Identity() bambou.Identity {

	return EgressACLTemplateIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *EgressACLTemplate) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *EgressACLTemplate) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the EgressACLTemplate from the server
func (o *EgressACLTemplate) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the EgressACLTemplate into the server
func (o *EgressACLTemplate) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the EgressACLTemplate from the server
func (o *EgressACLTemplate) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the EgressACLTemplate
func (o *EgressACLTemplate) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the EgressACLTemplate
func (o *EgressACLTemplate) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EgressACLEntryTemplates retrieves the list of child EgressACLEntryTemplates of the EgressACLTemplate
func (o *EgressACLTemplate) EgressACLEntryTemplates(info *bambou.FetchingInfo) (EgressACLEntryTemplatesList, *bambou.Error) {

	var list EgressACLEntryTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, EgressACLEntryTemplateIdentity, &list, info)
	return list, err
}

// CreateEgressACLEntryTemplate creates a new child EgressACLEntryTemplate under the EgressACLTemplate
func (o *EgressACLTemplate) CreateEgressACLEntryTemplate(child *EgressACLEntryTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the EgressACLTemplate
func (o *EgressACLTemplate) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the EgressACLTemplate
func (o *EgressACLTemplate) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMs retrieves the list of child VMs of the EgressACLTemplate
func (o *EgressACLTemplate) VMs(info *bambou.FetchingInfo) (VMsList, *bambou.Error) {

	var list VMsList
	err := bambou.CurrentSession().FetchChildren(o, VMIdentity, &list, info)
	return list, err
}

// CreateVM creates a new child VM under the EgressACLTemplate
func (o *EgressACLTemplate) CreateVM(child *VM) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Jobs retrieves the list of child Jobs of the EgressACLTemplate
func (o *EgressACLTemplate) Jobs(info *bambou.FetchingInfo) (JobsList, *bambou.Error) {

	var list JobsList
	err := bambou.CurrentSession().FetchChildren(o, JobIdentity, &list, info)
	return list, err
}

// CreateJob creates a new child Job under the EgressACLTemplate
func (o *EgressACLTemplate) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Containers retrieves the list of child Containers of the EgressACLTemplate
func (o *EgressACLTemplate) Containers(info *bambou.FetchingInfo) (ContainersList, *bambou.Error) {

	var list ContainersList
	err := bambou.CurrentSession().FetchChildren(o, ContainerIdentity, &list, info)
	return list, err
}

// CreateContainer creates a new child Container under the EgressACLTemplate
func (o *EgressACLTemplate) CreateContainer(child *Container) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the EgressACLTemplate
func (o *EgressACLTemplate) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}

// CreateEventLog creates a new child EventLog under the EgressACLTemplate
func (o *EgressACLTemplate) CreateEventLog(child *EventLog) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
