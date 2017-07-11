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

// GatewayTemplateIdentity represents the Identity of the object
var GatewayTemplateIdentity = bambou.Identity{
	Name:     "gatewaytemplate",
	Category: "gatewaytemplates",
}

// GatewayTemplatesList represents a list of GatewayTemplates
type GatewayTemplatesList []*GatewayTemplate

// GatewayTemplatesAncestor is the interface of an ancestor of a GatewayTemplate must implement.
type GatewayTemplatesAncestor interface {
	GatewayTemplates(*bambou.FetchingInfo) (GatewayTemplatesList, *bambou.Error)
	CreateGatewayTemplates(*GatewayTemplate) *bambou.Error
}

// GatewayTemplate represents the model of a gatewaytemplate
type GatewayTemplate struct {
	ID            string `json:"ID,omitempty"`
	ParentID      string `json:"parentID,omitempty"`
	ParentType    string `json:"parentType,omitempty"`
	Owner         string `json:"owner,omitempty"`
	Name          string `json:"name,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	Personality   string `json:"personality,omitempty"`
	Description   string `json:"description,omitempty"`
	EnterpriseID  string `json:"enterpriseID,omitempty"`
	EntityScope   string `json:"entityScope,omitempty"`
	ExternalID    string `json:"externalID,omitempty"`
}

// NewGatewayTemplate returns a new *GatewayTemplate
func NewGatewayTemplate() *GatewayTemplate {

	return &GatewayTemplate{
		Personality: "VRSG",
	}
}

// Identity returns the Identity of the object.
func (o *GatewayTemplate) Identity() bambou.Identity {

	return GatewayTemplateIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *GatewayTemplate) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *GatewayTemplate) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the GatewayTemplate from the server
func (o *GatewayTemplate) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the GatewayTemplate into the server
func (o *GatewayTemplate) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the GatewayTemplate from the server
func (o *GatewayTemplate) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the GatewayTemplate
func (o *GatewayTemplate) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the GatewayTemplate
func (o *GatewayTemplate) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the GatewayTemplate
func (o *GatewayTemplate) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the GatewayTemplate
func (o *GatewayTemplate) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PortTemplates retrieves the list of child PortTemplates of the GatewayTemplate
func (o *GatewayTemplate) PortTemplates(info *bambou.FetchingInfo) (PortTemplatesList, *bambou.Error) {

	var list PortTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, PortTemplateIdentity, &list, info)
	return list, err
}

// CreatePortTemplate creates a new child PortTemplate under the GatewayTemplate
func (o *GatewayTemplate) CreatePortTemplate(child *PortTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
