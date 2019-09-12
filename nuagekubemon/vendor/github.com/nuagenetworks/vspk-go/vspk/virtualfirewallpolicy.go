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

// VirtualFirewallPolicyIdentity represents the Identity of the object
var VirtualFirewallPolicyIdentity = bambou.Identity{
	Name:     "virtualfirewallpolicy",
	Category: "virtualfirewallpolicies",
}

// VirtualFirewallPoliciesList represents a list of VirtualFirewallPolicies
type VirtualFirewallPoliciesList []*VirtualFirewallPolicy

// VirtualFirewallPoliciesAncestor is the interface that an ancestor of a VirtualFirewallPolicy must implement.
// An Ancestor is defined as an entity that has VirtualFirewallPolicy as a descendant.
// An Ancestor can get a list of its child VirtualFirewallPolicies, but not necessarily create one.
type VirtualFirewallPoliciesAncestor interface {
	VirtualFirewallPolicies(*bambou.FetchingInfo) (VirtualFirewallPoliciesList, *bambou.Error)
}

// VirtualFirewallPoliciesParent is the interface that a parent of a VirtualFirewallPolicy must implement.
// A Parent is defined as an entity that has VirtualFirewallPolicy as a child.
// A Parent is an Ancestor which can create a VirtualFirewallPolicy.
type VirtualFirewallPoliciesParent interface {
	VirtualFirewallPoliciesAncestor
	CreateVirtualFirewallPolicy(*VirtualFirewallPolicy) *bambou.Error
}

// VirtualFirewallPolicy represents the model of a virtualfirewallpolicy
type VirtualFirewallPolicy struct {
	ID                             string        `json:"ID,omitempty"`
	ParentID                       string        `json:"parentID,omitempty"`
	ParentType                     string        `json:"parentType,omitempty"`
	Owner                          string        `json:"owner,omitempty"`
	Name                           string        `json:"name,omitempty"`
	LastUpdatedBy                  string        `json:"lastUpdatedBy,omitempty"`
	Active                         bool          `json:"active"`
	DefaultAllowIP                 bool          `json:"defaultAllowIP"`
	DefaultAllowNonIP              bool          `json:"defaultAllowNonIP"`
	DefaultInstallACLImplicitRules bool          `json:"defaultInstallACLImplicitRules"`
	Description                    string        `json:"description,omitempty"`
	AllowAddressSpoof              bool          `json:"allowAddressSpoof"`
	EmbeddedMetadata               []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope                    string        `json:"entityScope,omitempty"`
	PolicyState                    string        `json:"policyState,omitempty"`
	Priority                       int           `json:"priority,omitempty"`
	PriorityType                   string        `json:"priorityType,omitempty"`
	AssociatedEgressTemplateID     string        `json:"associatedEgressTemplateID,omitempty"`
	AssociatedIngressTemplateID    string        `json:"associatedIngressTemplateID,omitempty"`
	AssociatedLiveEntityID         string        `json:"associatedLiveEntityID,omitempty"`
	AutoGeneratePriority           bool          `json:"autoGeneratePriority"`
	ExternalID                     string        `json:"externalID,omitempty"`
}

// NewVirtualFirewallPolicy returns a new *VirtualFirewallPolicy
func NewVirtualFirewallPolicy() *VirtualFirewallPolicy {

	return &VirtualFirewallPolicy{
		Active:                         false,
		DefaultAllowIP:                 false,
		DefaultAllowNonIP:              false,
		DefaultInstallACLImplicitRules: false,
		AllowAddressSpoof:              false,
		AutoGeneratePriority:           false,
	}
}

// Identity returns the Identity of the object.
func (o *VirtualFirewallPolicy) Identity() bambou.Identity {

	return VirtualFirewallPolicyIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VirtualFirewallPolicy) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VirtualFirewallPolicy) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VirtualFirewallPolicy from the server
func (o *VirtualFirewallPolicy) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VirtualFirewallPolicy into the server
func (o *VirtualFirewallPolicy) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VirtualFirewallPolicy from the server
func (o *VirtualFirewallPolicy) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the VirtualFirewallPolicy
func (o *VirtualFirewallPolicy) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VirtualFirewallPolicy
func (o *VirtualFirewallPolicy) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VirtualFirewallRules retrieves the list of child VirtualFirewallRules of the VirtualFirewallPolicy
func (o *VirtualFirewallPolicy) VirtualFirewallRules(info *bambou.FetchingInfo) (VirtualFirewallRulesList, *bambou.Error) {

	var list VirtualFirewallRulesList
	err := bambou.CurrentSession().FetchChildren(o, VirtualFirewallRuleIdentity, &list, info)
	return list, err
}

// CreateVirtualFirewallRule creates a new child VirtualFirewallRule under the VirtualFirewallPolicy
func (o *VirtualFirewallPolicy) CreateVirtualFirewallRule(child *VirtualFirewallRule) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VirtualFirewallPolicy
func (o *VirtualFirewallPolicy) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VirtualFirewallPolicy
func (o *VirtualFirewallPolicy) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
