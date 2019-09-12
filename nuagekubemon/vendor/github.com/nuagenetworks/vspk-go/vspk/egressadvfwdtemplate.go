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

// EgressAdvFwdTemplateIdentity represents the Identity of the object
var EgressAdvFwdTemplateIdentity = bambou.Identity{
	Name:     "egressadvfwdtemplate",
	Category: "egressadvfwdtemplates",
}

// EgressAdvFwdTemplatesList represents a list of EgressAdvFwdTemplates
type EgressAdvFwdTemplatesList []*EgressAdvFwdTemplate

// EgressAdvFwdTemplatesAncestor is the interface that an ancestor of a EgressAdvFwdTemplate must implement.
// An Ancestor is defined as an entity that has EgressAdvFwdTemplate as a descendant.
// An Ancestor can get a list of its child EgressAdvFwdTemplates, but not necessarily create one.
type EgressAdvFwdTemplatesAncestor interface {
	EgressAdvFwdTemplates(*bambou.FetchingInfo) (EgressAdvFwdTemplatesList, *bambou.Error)
}

// EgressAdvFwdTemplatesParent is the interface that a parent of a EgressAdvFwdTemplate must implement.
// A Parent is defined as an entity that has EgressAdvFwdTemplate as a child.
// A Parent is an Ancestor which can create a EgressAdvFwdTemplate.
type EgressAdvFwdTemplatesParent interface {
	EgressAdvFwdTemplatesAncestor
	CreateEgressAdvFwdTemplate(*EgressAdvFwdTemplate) *bambou.Error
}

// EgressAdvFwdTemplate represents the model of a egressadvfwdtemplate
type EgressAdvFwdTemplate struct {
	ID                     string        `json:"ID,omitempty"`
	ParentID               string        `json:"parentID,omitempty"`
	ParentType             string        `json:"parentType,omitempty"`
	Owner                  string        `json:"owner,omitempty"`
	Name                   string        `json:"name,omitempty"`
	LastUpdatedBy          string        `json:"lastUpdatedBy,omitempty"`
	Active                 bool          `json:"active"`
	DefaultAllowIP         bool          `json:"defaultAllowIP"`
	DefaultAllowNonIP      bool          `json:"defaultAllowNonIP"`
	Description            string        `json:"description,omitempty"`
	EmbeddedMetadata       []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope            string        `json:"entityScope,omitempty"`
	PolicyState            string        `json:"policyState,omitempty"`
	Priority               int           `json:"priority,omitempty"`
	PriorityType           string        `json:"priorityType,omitempty"`
	AssociatedLiveEntityID string        `json:"associatedLiveEntityID,omitempty"`
	AutoGeneratePriority   bool          `json:"autoGeneratePriority"`
	ExternalID             string        `json:"externalID,omitempty"`
}

// NewEgressAdvFwdTemplate returns a new *EgressAdvFwdTemplate
func NewEgressAdvFwdTemplate() *EgressAdvFwdTemplate {

	return &EgressAdvFwdTemplate{}
}

// Identity returns the Identity of the object.
func (o *EgressAdvFwdTemplate) Identity() bambou.Identity {

	return EgressAdvFwdTemplateIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *EgressAdvFwdTemplate) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *EgressAdvFwdTemplate) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the EgressAdvFwdTemplate from the server
func (o *EgressAdvFwdTemplate) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the EgressAdvFwdTemplate into the server
func (o *EgressAdvFwdTemplate) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the EgressAdvFwdTemplate from the server
func (o *EgressAdvFwdTemplate) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the EgressAdvFwdTemplate
func (o *EgressAdvFwdTemplate) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the EgressAdvFwdTemplate
func (o *EgressAdvFwdTemplate) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EgressAdvFwdEntryTemplates retrieves the list of child EgressAdvFwdEntryTemplates of the EgressAdvFwdTemplate
func (o *EgressAdvFwdTemplate) EgressAdvFwdEntryTemplates(info *bambou.FetchingInfo) (EgressAdvFwdEntryTemplatesList, *bambou.Error) {

	var list EgressAdvFwdEntryTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, EgressAdvFwdEntryTemplateIdentity, &list, info)
	return list, err
}

// CreateEgressAdvFwdEntryTemplate creates a new child EgressAdvFwdEntryTemplate under the EgressAdvFwdTemplate
func (o *EgressAdvFwdTemplate) CreateEgressAdvFwdEntryTemplate(child *EgressAdvFwdEntryTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the EgressAdvFwdTemplate
func (o *EgressAdvFwdTemplate) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the EgressAdvFwdTemplate
func (o *EgressAdvFwdTemplate) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
