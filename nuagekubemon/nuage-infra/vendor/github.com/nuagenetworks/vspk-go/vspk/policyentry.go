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

// PolicyEntryIdentity represents the Identity of the object
var PolicyEntryIdentity = bambou.Identity{
	Name:     "policyentry",
	Category: "policyentries",
}

// PolicyEntriesList represents a list of PolicyEntries
type PolicyEntriesList []*PolicyEntry

// PolicyEntriesAncestor is the interface that an ancestor of a PolicyEntry must implement.
// An Ancestor is defined as an entity that has PolicyEntry as a descendant.
// An Ancestor can get a list of its child PolicyEntries, but not necessarily create one.
type PolicyEntriesAncestor interface {
	PolicyEntries(*bambou.FetchingInfo) (PolicyEntriesList, *bambou.Error)
}

// PolicyEntriesParent is the interface that a parent of a PolicyEntry must implement.
// A Parent is defined as an entity that has PolicyEntry as a child.
// A Parent is an Ancestor which can create a PolicyEntry.
type PolicyEntriesParent interface {
	PolicyEntriesAncestor
	CreatePolicyEntry(*PolicyEntry) *bambou.Error
}

// PolicyEntry represents the model of a policyentry
type PolicyEntry struct {
	ID            string      `json:"ID,omitempty"`
	ParentID      string      `json:"parentID,omitempty"`
	ParentType    string      `json:"parentType,omitempty"`
	Owner         string      `json:"owner,omitempty"`
	Name          string      `json:"name,omitempty"`
	MatchCriteria interface{} `json:"matchCriteria,omitempty"`
	Actions       interface{} `json:"actions,omitempty"`
	Description   string      `json:"description,omitempty"`
}

// NewPolicyEntry returns a new *PolicyEntry
func NewPolicyEntry() *PolicyEntry {

	return &PolicyEntry{}
}

// Identity returns the Identity of the object.
func (o *PolicyEntry) Identity() bambou.Identity {

	return PolicyEntryIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *PolicyEntry) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *PolicyEntry) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the PolicyEntry from the server
func (o *PolicyEntry) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the PolicyEntry into the server
func (o *PolicyEntry) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the PolicyEntry from the server
func (o *PolicyEntry) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}
