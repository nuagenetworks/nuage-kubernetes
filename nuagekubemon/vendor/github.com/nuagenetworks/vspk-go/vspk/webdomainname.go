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

// WebDomainNameIdentity represents the Identity of the object
var WebDomainNameIdentity = bambou.Identity{
	Name:     "webdomainname",
	Category: "webdomainnames",
}

// WebDomainNamesList represents a list of WebDomainNames
type WebDomainNamesList []*WebDomainName

// WebDomainNamesAncestor is the interface that an ancestor of a WebDomainName must implement.
// An Ancestor is defined as an entity that has WebDomainName as a descendant.
// An Ancestor can get a list of its child WebDomainNames, but not necessarily create one.
type WebDomainNamesAncestor interface {
	WebDomainNames(*bambou.FetchingInfo) (WebDomainNamesList, *bambou.Error)
}

// WebDomainNamesParent is the interface that a parent of a WebDomainName must implement.
// A Parent is defined as an entity that has WebDomainName as a child.
// A Parent is an Ancestor which can create a WebDomainName.
type WebDomainNamesParent interface {
	WebDomainNamesAncestor
	CreateWebDomainName(*WebDomainName) *bambou.Error
}

// WebDomainName represents the model of a webdomainname
type WebDomainName struct {
	ID                    string        `json:"ID,omitempty"`
	ParentID              string        `json:"parentID,omitempty"`
	ParentType            string        `json:"parentType,omitempty"`
	Owner                 string        `json:"owner,omitempty"`
	Name                  string        `json:"name,omitempty"`
	LastUpdatedBy         string        `json:"lastUpdatedBy,omitempty"`
	WebCategoryIdentifier int           `json:"webCategoryIdentifier,omitempty"`
	EmbeddedMetadata      []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope           string        `json:"entityScope,omitempty"`
	ExternalID            string        `json:"externalID,omitempty"`
}

// NewWebDomainName returns a new *WebDomainName
func NewWebDomainName() *WebDomainName {

	return &WebDomainName{}
}

// Identity returns the Identity of the object.
func (o *WebDomainName) Identity() bambou.Identity {

	return WebDomainNameIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *WebDomainName) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *WebDomainName) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the WebDomainName from the server
func (o *WebDomainName) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the WebDomainName into the server
func (o *WebDomainName) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the WebDomainName from the server
func (o *WebDomainName) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// WebCategories retrieves the list of child WebCategories of the WebDomainName
func (o *WebDomainName) WebCategories(info *bambou.FetchingInfo) (WebCategoriesList, *bambou.Error) {

	var list WebCategoriesList
	err := bambou.CurrentSession().FetchChildren(o, WebCategoryIdentity, &list, info)
	return list, err
}

// AssignWebCategories assigns the list of WebCategories to the WebDomainName
func (o *WebDomainName) AssignWebCategories(children WebCategoriesList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, WebCategoryIdentity)
}

// Metadatas retrieves the list of child Metadatas of the WebDomainName
func (o *WebDomainName) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the WebDomainName
func (o *WebDomainName) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the WebDomainName
func (o *WebDomainName) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the WebDomainName
func (o *WebDomainName) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
