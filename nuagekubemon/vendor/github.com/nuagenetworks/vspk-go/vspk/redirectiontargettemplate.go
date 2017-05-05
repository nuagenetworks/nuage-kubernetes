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

// RedirectionTargetTemplateIdentity represents the Identity of the object
var RedirectionTargetTemplateIdentity = bambou.Identity{
	Name:     "redirectiontargettemplate",
	Category: "redirectiontargettemplates",
}

// RedirectionTargetTemplatesList represents a list of RedirectionTargetTemplates
type RedirectionTargetTemplatesList []*RedirectionTargetTemplate

// RedirectionTargetTemplatesAncestor is the interface that an ancestor of a RedirectionTargetTemplate must implement.
// An Ancestor is defined as an entity that has RedirectionTargetTemplate as a descendant.
// An Ancestor can get a list of its child RedirectionTargetTemplates, but not necessarily create one.
type RedirectionTargetTemplatesAncestor interface {
	RedirectionTargetTemplates(*bambou.FetchingInfo) (RedirectionTargetTemplatesList, *bambou.Error)
}

// RedirectionTargetTemplatesParent is the interface that a parent of a RedirectionTargetTemplate must implement.
// A Parent is defined as an entity that has RedirectionTargetTemplate as a child.
// A Parent is an Ancestor which can create a RedirectionTargetTemplate.
type RedirectionTargetTemplatesParent interface {
	RedirectionTargetTemplatesAncestor
	CreateRedirectionTargetTemplate(*RedirectionTargetTemplate) *bambou.Error
}

// RedirectionTargetTemplate represents the model of a redirectiontargettemplate
type RedirectionTargetTemplate struct {
	ID                string `json:"ID,omitempty"`
	ParentID          string `json:"parentID,omitempty"`
	ParentType        string `json:"parentType,omitempty"`
	Owner             string `json:"owner,omitempty"`
	Name              string `json:"name,omitempty"`
	LastUpdatedBy     string `json:"lastUpdatedBy,omitempty"`
	RedundancyEnabled bool   `json:"redundancyEnabled"`
	Description       string `json:"description,omitempty"`
	EndPointType      string `json:"endPointType,omitempty"`
	EntityScope       string `json:"entityScope,omitempty"`
	TriggerType       string `json:"triggerType,omitempty"`
	ExternalID        string `json:"externalID,omitempty"`
}

// NewRedirectionTargetTemplate returns a new *RedirectionTargetTemplate
func NewRedirectionTargetTemplate() *RedirectionTargetTemplate {

	return &RedirectionTargetTemplate{
		EndPointType: "L3",
	}
}

// Identity returns the Identity of the object.
func (o *RedirectionTargetTemplate) Identity() bambou.Identity {

	return RedirectionTargetTemplateIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *RedirectionTargetTemplate) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *RedirectionTargetTemplate) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the RedirectionTargetTemplate from the server
func (o *RedirectionTargetTemplate) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the RedirectionTargetTemplate into the server
func (o *RedirectionTargetTemplate) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the RedirectionTargetTemplate from the server
func (o *RedirectionTargetTemplate) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the RedirectionTargetTemplate
func (o *RedirectionTargetTemplate) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the RedirectionTargetTemplate
func (o *RedirectionTargetTemplate) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the RedirectionTargetTemplate
func (o *RedirectionTargetTemplate) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the RedirectionTargetTemplate
func (o *RedirectionTargetTemplate) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CreateJob creates a new child Job under the RedirectionTargetTemplate
func (o *RedirectionTargetTemplate) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the RedirectionTargetTemplate
func (o *RedirectionTargetTemplate) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
