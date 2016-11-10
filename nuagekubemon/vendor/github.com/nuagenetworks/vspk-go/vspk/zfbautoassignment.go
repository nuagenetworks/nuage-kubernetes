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

// ZFBAutoAssignmentIdentity represents the Identity of the object
var ZFBAutoAssignmentIdentity = bambou.Identity{
	Name:     "zfbautoassignment",
	Category: "zfbautoassignments",
}

// ZFBAutoAssignmentsList represents a list of ZFBAutoAssignments
type ZFBAutoAssignmentsList []*ZFBAutoAssignment

// ZFBAutoAssignmentsAncestor is the interface of an ancestor of a ZFBAutoAssignment must implement.
type ZFBAutoAssignmentsAncestor interface {
	ZFBAutoAssignments(*bambou.FetchingInfo) (ZFBAutoAssignmentsList, *bambou.Error)
	CreateZFBAutoAssignments(*ZFBAutoAssignment) *bambou.Error
}

// ZFBAutoAssignment represents the model of a zfbautoassignment
type ZFBAutoAssignment struct {
	ID                       string        `json:"ID,omitempty"`
	ParentID                 string        `json:"parentID,omitempty"`
	ParentType               string        `json:"parentType,omitempty"`
	Owner                    string        `json:"owner,omitempty"`
	ZFBMatchAttribute        string        `json:"ZFBMatchAttribute,omitempty"`
	ZFBMatchAttributeValues  []interface{} `json:"ZFBMatchAttributeValues,omitempty"`
	Name                     string        `json:"name,omitempty"`
	LastUpdatedBy            string        `json:"lastUpdatedBy,omitempty"`
	Description              string        `json:"description,omitempty"`
	EntityScope              string        `json:"entityScope,omitempty"`
	Priority                 int           `json:"priority,omitempty"`
	AssociatedEnterpriseID   string        `json:"associatedEnterpriseID,omitempty"`
	AssociatedEnterpriseName string        `json:"associatedEnterpriseName,omitempty"`
	ExternalID               string        `json:"externalID,omitempty"`
}

// NewZFBAutoAssignment returns a new *ZFBAutoAssignment
func NewZFBAutoAssignment() *ZFBAutoAssignment {

	return &ZFBAutoAssignment{}
}

// Identity returns the Identity of the object.
func (o *ZFBAutoAssignment) Identity() bambou.Identity {

	return ZFBAutoAssignmentIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *ZFBAutoAssignment) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *ZFBAutoAssignment) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the ZFBAutoAssignment from the server
func (o *ZFBAutoAssignment) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the ZFBAutoAssignment into the server
func (o *ZFBAutoAssignment) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the ZFBAutoAssignment from the server
func (o *ZFBAutoAssignment) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}
