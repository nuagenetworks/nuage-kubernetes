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

// FirewallAclIdentity represents the Identity of the object
var FirewallAclIdentity = bambou.Identity{
	Name:     "firewallacl",
	Category: "firewallacls",
}

// FirewallAclsList represents a list of FirewallAcls
type FirewallAclsList []*FirewallAcl

// FirewallAclsAncestor is the interface that an ancestor of a FirewallAcl must implement.
// An Ancestor is defined as an entity that has FirewallAcl as a descendant.
// An Ancestor can get a list of its child FirewallAcls, but not necessarily create one.
type FirewallAclsAncestor interface {
	FirewallAcls(*bambou.FetchingInfo) (FirewallAclsList, *bambou.Error)
}

// FirewallAclsParent is the interface that a parent of a FirewallAcl must implement.
// A Parent is defined as an entity that has FirewallAcl as a child.
// A Parent is an Ancestor which can create a FirewallAcl.
type FirewallAclsParent interface {
	FirewallAclsAncestor
	CreateFirewallAcl(*FirewallAcl) *bambou.Error
}

// FirewallAcl represents the model of a firewallacl
type FirewallAcl struct {
	ID                   string        `json:"ID,omitempty"`
	ParentID             string        `json:"parentID,omitempty"`
	ParentType           string        `json:"parentType,omitempty"`
	Owner                string        `json:"owner,omitempty"`
	Name                 string        `json:"name,omitempty"`
	LastUpdatedBy        string        `json:"lastUpdatedBy,omitempty"`
	Active               bool          `json:"active"`
	DefaultAllowIP       bool          `json:"defaultAllowIP"`
	DefaultAllowNonIP    bool          `json:"defaultAllowNonIP"`
	Description          string        `json:"description,omitempty"`
	EmbeddedMetadata     []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope          string        `json:"entityScope,omitempty"`
	RuleIds              []interface{} `json:"ruleIds,omitempty"`
	AutoGeneratePriority bool          `json:"autoGeneratePriority"`
	ExternalID           string        `json:"externalID,omitempty"`
}

// NewFirewallAcl returns a new *FirewallAcl
func NewFirewallAcl() *FirewallAcl {

	return &FirewallAcl{}
}

// Identity returns the Identity of the object.
func (o *FirewallAcl) Identity() bambou.Identity {

	return FirewallAclIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *FirewallAcl) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *FirewallAcl) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the FirewallAcl from the server
func (o *FirewallAcl) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the FirewallAcl into the server
func (o *FirewallAcl) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the FirewallAcl from the server
func (o *FirewallAcl) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the FirewallAcl
func (o *FirewallAcl) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the FirewallAcl
func (o *FirewallAcl) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// FirewallRules retrieves the list of child FirewallRules of the FirewallAcl
func (o *FirewallAcl) FirewallRules(info *bambou.FetchingInfo) (FirewallRulesList, *bambou.Error) {

	var list FirewallRulesList
	err := bambou.CurrentSession().FetchChildren(o, FirewallRuleIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the FirewallAcl
func (o *FirewallAcl) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the FirewallAcl
func (o *FirewallAcl) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Domains retrieves the list of child Domains of the FirewallAcl
func (o *FirewallAcl) Domains(info *bambou.FetchingInfo) (DomainsList, *bambou.Error) {

	var list DomainsList
	err := bambou.CurrentSession().FetchChildren(o, DomainIdentity, &list, info)
	return list, err
}
