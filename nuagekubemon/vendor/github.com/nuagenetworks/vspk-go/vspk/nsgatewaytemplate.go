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

// NSGatewayTemplateIdentity represents the Identity of the object
var NSGatewayTemplateIdentity = bambou.Identity{
	Name:     "nsgatewaytemplate",
	Category: "nsgatewaytemplates",
}

// NSGatewayTemplatesList represents a list of NSGatewayTemplates
type NSGatewayTemplatesList []*NSGatewayTemplate

// NSGatewayTemplatesAncestor is the interface that an ancestor of a NSGatewayTemplate must implement.
// An Ancestor is defined as an entity that has NSGatewayTemplate as a descendant.
// An Ancestor can get a list of its child NSGatewayTemplates, but not necessarily create one.
type NSGatewayTemplatesAncestor interface {
	NSGatewayTemplates(*bambou.FetchingInfo) (NSGatewayTemplatesList, *bambou.Error)
}

// NSGatewayTemplatesParent is the interface that a parent of a NSGatewayTemplate must implement.
// A Parent is defined as an entity that has NSGatewayTemplate as a child.
// A Parent is an Ancestor which can create a NSGatewayTemplate.
type NSGatewayTemplatesParent interface {
	NSGatewayTemplatesAncestor
	CreateNSGatewayTemplate(*NSGatewayTemplate) *bambou.Error
}

// NSGatewayTemplate represents the model of a nsgatewaytemplate
type NSGatewayTemplate struct {
	ID                            string `json:"ID,omitempty"`
	ParentID                      string `json:"parentID,omitempty"`
	ParentType                    string `json:"parentType,omitempty"`
	Owner                         string `json:"owner,omitempty"`
	SSHService                    string `json:"SSHService,omitempty"`
	Name                          string `json:"name,omitempty"`
	LastUpdatedBy                 string `json:"lastUpdatedBy,omitempty"`
	Personality                   string `json:"personality,omitempty"`
	Description                   string `json:"description,omitempty"`
	InfrastructureAccessProfileID string `json:"infrastructureAccessProfileID,omitempty"`
	InfrastructureProfileID       string `json:"infrastructureProfileID,omitempty"`
	InstanceSSHOverride           string `json:"instanceSSHOverride,omitempty"`
	EnterpriseID                  string `json:"enterpriseID,omitempty"`
	EntityScope                   string `json:"entityScope,omitempty"`
	ExternalID                    string `json:"externalID,omitempty"`
}

// NewNSGatewayTemplate returns a new *NSGatewayTemplate
func NewNSGatewayTemplate() *NSGatewayTemplate {

	return &NSGatewayTemplate{
		SSHService:          "ENABLED",
		InstanceSSHOverride: "DISALLOWED",
	}
}

// Identity returns the Identity of the object.
func (o *NSGatewayTemplate) Identity() bambou.Identity {

	return NSGatewayTemplateIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *NSGatewayTemplate) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *NSGatewayTemplate) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the NSGatewayTemplate from the server
func (o *NSGatewayTemplate) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the NSGatewayTemplate into the server
func (o *NSGatewayTemplate) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the NSGatewayTemplate from the server
func (o *NSGatewayTemplate) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the NSGatewayTemplate
func (o *NSGatewayTemplate) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the NSGatewayTemplate
func (o *NSGatewayTemplate) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the NSGatewayTemplate
func (o *NSGatewayTemplate) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the NSGatewayTemplate
func (o *NSGatewayTemplate) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSPortTemplates retrieves the list of child NSPortTemplates of the NSGatewayTemplate
func (o *NSGatewayTemplate) NSPortTemplates(info *bambou.FetchingInfo) (NSPortTemplatesList, *bambou.Error) {

	var list NSPortTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, NSPortTemplateIdentity, &list, info)
	return list, err
}

// CreateNSPortTemplate creates a new child NSPortTemplate under the NSGatewayTemplate
func (o *NSGatewayTemplate) CreateNSPortTemplate(child *NSPortTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
