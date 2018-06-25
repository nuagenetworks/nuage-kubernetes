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

// IngressQOSPolicyIdentity represents the Identity of the object
var IngressQOSPolicyIdentity = bambou.Identity{
	Name:     "ingressqospolicy",
	Category: "ingressqospolicies",
}

// IngressQOSPoliciesList represents a list of IngressQOSPolicies
type IngressQOSPoliciesList []*IngressQOSPolicy

// IngressQOSPoliciesAncestor is the interface that an ancestor of a IngressQOSPolicy must implement.
// An Ancestor is defined as an entity that has IngressQOSPolicy as a descendant.
// An Ancestor can get a list of its child IngressQOSPolicies, but not necessarily create one.
type IngressQOSPoliciesAncestor interface {
	IngressQOSPolicies(*bambou.FetchingInfo) (IngressQOSPoliciesList, *bambou.Error)
}

// IngressQOSPoliciesParent is the interface that a parent of a IngressQOSPolicy must implement.
// A Parent is defined as an entity that has IngressQOSPolicy as a child.
// A Parent is an Ancestor which can create a IngressQOSPolicy.
type IngressQOSPoliciesParent interface {
	IngressQOSPoliciesAncestor
	CreateIngressQOSPolicy(*IngressQOSPolicy) *bambou.Error
}

// IngressQOSPolicy represents the model of a ingressqospolicy
type IngressQOSPolicy struct {
	ID                                 string        `json:"ID,omitempty"`
	ParentID                           string        `json:"parentID,omitempty"`
	ParentType                         string        `json:"parentType,omitempty"`
	Owner                              string        `json:"owner,omitempty"`
	Name                               string        `json:"name,omitempty"`
	ParentQueueAssociatedRateLimiterID string        `json:"parentQueueAssociatedRateLimiterID,omitempty"`
	LastUpdatedBy                      string        `json:"lastUpdatedBy,omitempty"`
	Description                        string        `json:"description,omitempty"`
	EntityScope                        string        `json:"entityScope,omitempty"`
	AssocEgressQosId                   string        `json:"assocEgressQosId,omitempty"`
	Queue1AssociatedRateLimiterID      string        `json:"queue1AssociatedRateLimiterID,omitempty"`
	Queue1ForwardingClasses            []interface{} `json:"queue1ForwardingClasses,omitempty"`
	Queue2AssociatedRateLimiterID      string        `json:"queue2AssociatedRateLimiterID,omitempty"`
	Queue2ForwardingClasses            []interface{} `json:"queue2ForwardingClasses,omitempty"`
	Queue3AssociatedRateLimiterID      string        `json:"queue3AssociatedRateLimiterID,omitempty"`
	Queue3ForwardingClasses            []interface{} `json:"queue3ForwardingClasses,omitempty"`
	Queue4AssociatedRateLimiterID      string        `json:"queue4AssociatedRateLimiterID,omitempty"`
	Queue4ForwardingClasses            []interface{} `json:"queue4ForwardingClasses,omitempty"`
	ExternalID                         string        `json:"externalID,omitempty"`
}

// NewIngressQOSPolicy returns a new *IngressQOSPolicy
func NewIngressQOSPolicy() *IngressQOSPolicy {

	return &IngressQOSPolicy{}
}

// Identity returns the Identity of the object.
func (o *IngressQOSPolicy) Identity() bambou.Identity {

	return IngressQOSPolicyIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *IngressQOSPolicy) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *IngressQOSPolicy) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the IngressQOSPolicy from the server
func (o *IngressQOSPolicy) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the IngressQOSPolicy into the server
func (o *IngressQOSPolicy) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the IngressQOSPolicy from the server
func (o *IngressQOSPolicy) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the IngressQOSPolicy
func (o *IngressQOSPolicy) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the IngressQOSPolicy
func (o *IngressQOSPolicy) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the IngressQOSPolicy
func (o *IngressQOSPolicy) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the IngressQOSPolicy
func (o *IngressQOSPolicy) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
