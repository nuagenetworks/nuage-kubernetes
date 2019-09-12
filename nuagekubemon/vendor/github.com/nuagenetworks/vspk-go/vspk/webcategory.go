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

// WebCategoryIdentity represents the Identity of the object
var WebCategoryIdentity = bambou.Identity{
	Name:     "webcategory",
	Category: "webcategories",
}

// WebCategoriesList represents a list of WebCategories
type WebCategoriesList []*WebCategory

// WebCategoriesAncestor is the interface that an ancestor of a WebCategory must implement.
// An Ancestor is defined as an entity that has WebCategory as a descendant.
// An Ancestor can get a list of its child WebCategories, but not necessarily create one.
type WebCategoriesAncestor interface {
	WebCategories(*bambou.FetchingInfo) (WebCategoriesList, *bambou.Error)
}

// WebCategoriesParent is the interface that a parent of a WebCategory must implement.
// A Parent is defined as an entity that has WebCategory as a child.
// A Parent is an Ancestor which can create a WebCategory.
type WebCategoriesParent interface {
	WebCategoriesAncestor
	CreateWebCategory(*WebCategory) *bambou.Error
}

// WebCategory represents the model of a webcategory
type WebCategory struct {
	ID                    string        `json:"ID,omitempty"`
	ParentID              string        `json:"parentID,omitempty"`
	ParentType            string        `json:"parentType,omitempty"`
	Owner                 string        `json:"owner,omitempty"`
	Name                  string        `json:"name,omitempty"`
	LastUpdatedBy         string        `json:"lastUpdatedBy,omitempty"`
	WebCategoryIdentifier int           `json:"webCategoryIdentifier,omitempty"`
	DefaultCategory       bool          `json:"defaultCategory"`
	Description           string        `json:"description,omitempty"`
	EmbeddedMetadata      []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope           string        `json:"entityScope,omitempty"`
	ExternalID            string        `json:"externalID,omitempty"`
	Type                  string        `json:"type,omitempty"`
}

// NewWebCategory returns a new *WebCategory
func NewWebCategory() *WebCategory {

	return &WebCategory{
		DefaultCategory: false,
		Type:            "WEB_DOMAIN_NAME",
	}
}

// Identity returns the Identity of the object.
func (o *WebCategory) Identity() bambou.Identity {

	return WebCategoryIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *WebCategory) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *WebCategory) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the WebCategory from the server
func (o *WebCategory) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the WebCategory into the server
func (o *WebCategory) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the WebCategory from the server
func (o *WebCategory) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// WebDomainNames retrieves the list of child WebDomainNames of the WebCategory
func (o *WebCategory) WebDomainNames(info *bambou.FetchingInfo) (WebDomainNamesList, *bambou.Error) {

	var list WebDomainNamesList
	err := bambou.CurrentSession().FetchChildren(o, WebDomainNameIdentity, &list, info)
	return list, err
}

// AssignWebDomainNames assigns the list of WebDomainNames to the WebCategory
func (o *WebCategory) AssignWebDomainNames(children WebDomainNamesList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, WebDomainNameIdentity)
}

// Metadatas retrieves the list of child Metadatas of the WebCategory
func (o *WebCategory) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the WebCategory
func (o *WebCategory) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the WebCategory
func (o *WebCategory) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the WebCategory
func (o *WebCategory) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
