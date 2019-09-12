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

// DSCPRemarkingPolicyTableIdentity represents the Identity of the object
var DSCPRemarkingPolicyTableIdentity = bambou.Identity{
	Name:     "dscpremarkingpolicytable",
	Category: "dscpremarkingpolicytables",
}

// DSCPRemarkingPolicyTablesList represents a list of DSCPRemarkingPolicyTables
type DSCPRemarkingPolicyTablesList []*DSCPRemarkingPolicyTable

// DSCPRemarkingPolicyTablesAncestor is the interface that an ancestor of a DSCPRemarkingPolicyTable must implement.
// An Ancestor is defined as an entity that has DSCPRemarkingPolicyTable as a descendant.
// An Ancestor can get a list of its child DSCPRemarkingPolicyTables, but not necessarily create one.
type DSCPRemarkingPolicyTablesAncestor interface {
	DSCPRemarkingPolicyTables(*bambou.FetchingInfo) (DSCPRemarkingPolicyTablesList, *bambou.Error)
}

// DSCPRemarkingPolicyTablesParent is the interface that a parent of a DSCPRemarkingPolicyTable must implement.
// A Parent is defined as an entity that has DSCPRemarkingPolicyTable as a child.
// A Parent is an Ancestor which can create a DSCPRemarkingPolicyTable.
type DSCPRemarkingPolicyTablesParent interface {
	DSCPRemarkingPolicyTablesAncestor
	CreateDSCPRemarkingPolicyTable(*DSCPRemarkingPolicyTable) *bambou.Error
}

// DSCPRemarkingPolicyTable represents the model of a dscpremarkingpolicytable
type DSCPRemarkingPolicyTable struct {
	ID               string        `json:"ID,omitempty"`
	ParentID         string        `json:"parentID,omitempty"`
	ParentType       string        `json:"parentType,omitempty"`
	Owner            string        `json:"owner,omitempty"`
	Name             string        `json:"name,omitempty"`
	LastUpdatedBy    string        `json:"lastUpdatedBy,omitempty"`
	Description      string        `json:"description,omitempty"`
	EmbeddedMetadata []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope      string        `json:"entityScope,omitempty"`
	ExternalID       string        `json:"externalID,omitempty"`
}

// NewDSCPRemarkingPolicyTable returns a new *DSCPRemarkingPolicyTable
func NewDSCPRemarkingPolicyTable() *DSCPRemarkingPolicyTable {

	return &DSCPRemarkingPolicyTable{}
}

// Identity returns the Identity of the object.
func (o *DSCPRemarkingPolicyTable) Identity() bambou.Identity {

	return DSCPRemarkingPolicyTableIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *DSCPRemarkingPolicyTable) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *DSCPRemarkingPolicyTable) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the DSCPRemarkingPolicyTable from the server
func (o *DSCPRemarkingPolicyTable) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the DSCPRemarkingPolicyTable into the server
func (o *DSCPRemarkingPolicyTable) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the DSCPRemarkingPolicyTable from the server
func (o *DSCPRemarkingPolicyTable) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the DSCPRemarkingPolicyTable
func (o *DSCPRemarkingPolicyTable) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the DSCPRemarkingPolicyTable
func (o *DSCPRemarkingPolicyTable) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the DSCPRemarkingPolicyTable
func (o *DSCPRemarkingPolicyTable) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the DSCPRemarkingPolicyTable
func (o *DSCPRemarkingPolicyTable) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DSCPRemarkingPolicies retrieves the list of child DSCPRemarkingPolicies of the DSCPRemarkingPolicyTable
func (o *DSCPRemarkingPolicyTable) DSCPRemarkingPolicies(info *bambou.FetchingInfo) (DSCPRemarkingPoliciesList, *bambou.Error) {

	var list DSCPRemarkingPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, DSCPRemarkingPolicyIdentity, &list, info)
	return list, err
}

// CreateDSCPRemarkingPolicy creates a new child DSCPRemarkingPolicy under the DSCPRemarkingPolicyTable
func (o *DSCPRemarkingPolicyTable) CreateDSCPRemarkingPolicy(child *DSCPRemarkingPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
