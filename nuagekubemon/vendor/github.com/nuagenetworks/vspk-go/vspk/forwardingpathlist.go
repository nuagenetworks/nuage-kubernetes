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

// ForwardingPathListIdentity represents the Identity of the object
var ForwardingPathListIdentity = bambou.Identity{
	Name:     "forwardingpathlist",
	Category: "forwardingpathlists",
}

// ForwardingPathListsList represents a list of ForwardingPathLists
type ForwardingPathListsList []*ForwardingPathList

// ForwardingPathListsAncestor is the interface that an ancestor of a ForwardingPathList must implement.
// An Ancestor is defined as an entity that has ForwardingPathList as a descendant.
// An Ancestor can get a list of its child ForwardingPathLists, but not necessarily create one.
type ForwardingPathListsAncestor interface {
	ForwardingPathLists(*bambou.FetchingInfo) (ForwardingPathListsList, *bambou.Error)
}

// ForwardingPathListsParent is the interface that a parent of a ForwardingPathList must implement.
// A Parent is defined as an entity that has ForwardingPathList as a child.
// A Parent is an Ancestor which can create a ForwardingPathList.
type ForwardingPathListsParent interface {
	ForwardingPathListsAncestor
	CreateForwardingPathList(*ForwardingPathList) *bambou.Error
}

// ForwardingPathList represents the model of a forwardingpathlist
type ForwardingPathList struct {
	ID                   string        `json:"ID,omitempty"`
	ParentID             string        `json:"parentID,omitempty"`
	ParentType           string        `json:"parentType,omitempty"`
	Owner                string        `json:"owner,omitempty"`
	Name                 string        `json:"name,omitempty"`
	LastUpdatedBy        string        `json:"lastUpdatedBy,omitempty"`
	Description          string        `json:"description,omitempty"`
	EmbeddedMetadata     []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope          string        `json:"entityScope,omitempty"`
	ForwardingPathListID int           `json:"forwardingPathListID,omitempty"`
	ExternalID           string        `json:"externalID,omitempty"`
}

// NewForwardingPathList returns a new *ForwardingPathList
func NewForwardingPathList() *ForwardingPathList {

	return &ForwardingPathList{}
}

// Identity returns the Identity of the object.
func (o *ForwardingPathList) Identity() bambou.Identity {

	return ForwardingPathListIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *ForwardingPathList) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *ForwardingPathList) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the ForwardingPathList from the server
func (o *ForwardingPathList) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the ForwardingPathList into the server
func (o *ForwardingPathList) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the ForwardingPathList from the server
func (o *ForwardingPathList) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the ForwardingPathList
func (o *ForwardingPathList) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the ForwardingPathList
func (o *ForwardingPathList) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the ForwardingPathList
func (o *ForwardingPathList) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the ForwardingPathList
func (o *ForwardingPathList) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// ForwardingPathListEntries retrieves the list of child ForwardingPathListEntries of the ForwardingPathList
func (o *ForwardingPathList) ForwardingPathListEntries(info *bambou.FetchingInfo) (ForwardingPathListEntriesList, *bambou.Error) {

	var list ForwardingPathListEntriesList
	err := bambou.CurrentSession().FetchChildren(o, ForwardingPathListEntryIdentity, &list, info)
	return list, err
}

// CreateForwardingPathListEntry creates a new child ForwardingPathListEntry under the ForwardingPathList
func (o *ForwardingPathList) CreateForwardingPathListEntry(child *ForwardingPathListEntry) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
