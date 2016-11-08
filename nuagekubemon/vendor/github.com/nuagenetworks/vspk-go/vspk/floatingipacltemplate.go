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

// FloatingIPACLTemplateIdentity represents the Identity of the object
var FloatingIPACLTemplateIdentity = bambou.Identity{
	Name:     "egressfloatingipacltemplate",
	Category: "egressfloatingipacltemplates",
}

// FloatingIPACLTemplatesList represents a list of FloatingIPACLTemplates
type FloatingIPACLTemplatesList []*FloatingIPACLTemplate

// FloatingIPACLTemplatesAncestor is the interface of an ancestor of a FloatingIPACLTemplate must implement.
type FloatingIPACLTemplatesAncestor interface {
	FloatingIPACLTemplates(*bambou.FetchingInfo) (FloatingIPACLTemplatesList, *bambou.Error)
	CreateFloatingIPACLTemplates(*FloatingIPACLTemplate) *bambou.Error
}

// FloatingIPACLTemplate represents the model of a egressfloatingipacltemplate
type FloatingIPACLTemplate struct {
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
	EntityScope            string `json:"entityScope,omitempty"`
	PolicyState            string `json:"policyState,omitempty"`
	Priority               int    `json:"priority,omitempty"`
	PriorityType           string `json:"priorityType,omitempty"`
	AssociatedLiveEntityID string `json:"associatedLiveEntityID,omitempty"`
	ExternalID             string `json:"externalID,omitempty"`
}

// NewFloatingIPACLTemplate returns a new *FloatingIPACLTemplate
func NewFloatingIPACLTemplate() *FloatingIPACLTemplate {

	return &FloatingIPACLTemplate{}
}

// Identity returns the Identity of the object.
func (o *FloatingIPACLTemplate) Identity() bambou.Identity {

	return FloatingIPACLTemplateIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *FloatingIPACLTemplate) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *FloatingIPACLTemplate) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the FloatingIPACLTemplate from the server
func (o *FloatingIPACLTemplate) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the FloatingIPACLTemplate into the server
func (o *FloatingIPACLTemplate) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the FloatingIPACLTemplate from the server
func (o *FloatingIPACLTemplate) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the FloatingIPACLTemplate
func (o *FloatingIPACLTemplate) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the FloatingIPACLTemplate
func (o *FloatingIPACLTemplate) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// FloatingIPACLTemplateEntries retrieves the list of child FloatingIPACLTemplateEntries of the FloatingIPACLTemplate
func (o *FloatingIPACLTemplate) FloatingIPACLTemplateEntries(info *bambou.FetchingInfo) (FloatingIPACLTemplateEntriesList, *bambou.Error) {

	var list FloatingIPACLTemplateEntriesList
	err := bambou.CurrentSession().FetchChildren(o, FloatingIPACLTemplateEntryIdentity, &list, info)
	return list, err
}

// CreateFloatingIPACLTemplateEntry creates a new child FloatingIPACLTemplateEntry under the FloatingIPACLTemplate
func (o *FloatingIPACLTemplate) CreateFloatingIPACLTemplateEntry(child *FloatingIPACLTemplateEntry) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the FloatingIPACLTemplate
func (o *FloatingIPACLTemplate) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the FloatingIPACLTemplate
func (o *FloatingIPACLTemplate) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
