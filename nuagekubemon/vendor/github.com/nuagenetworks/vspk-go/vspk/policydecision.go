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

// PolicyDecisionIdentity represents the Identity of the object
var PolicyDecisionIdentity = bambou.Identity{
	Name:     "policydecision",
	Category: "policydecisions",
}

// PolicyDecisionsList represents a list of PolicyDecisions
type PolicyDecisionsList []*PolicyDecision

// PolicyDecisionsAncestor is the interface that an ancestor of a PolicyDecision must implement.
// An Ancestor is defined as an entity that has PolicyDecision as a descendant.
// An Ancestor can get a list of its child PolicyDecisions, but not necessarily create one.
type PolicyDecisionsAncestor interface {
	PolicyDecisions(*bambou.FetchingInfo) (PolicyDecisionsList, *bambou.Error)
}

// PolicyDecisionsParent is the interface that a parent of a PolicyDecision must implement.
// A Parent is defined as an entity that has PolicyDecision as a child.
// A Parent is an Ancestor which can create a PolicyDecision.
type PolicyDecisionsParent interface {
	PolicyDecisionsAncestor
	CreatePolicyDecision(*PolicyDecision) *bambou.Error
}

// PolicyDecision represents the model of a policydecision
type PolicyDecision struct {
	ID                         string        `json:"ID,omitempty"`
	ParentID                   string        `json:"parentID,omitempty"`
	ParentType                 string        `json:"parentType,omitempty"`
	Owner                      string        `json:"owner,omitempty"`
	LastUpdatedBy              string        `json:"lastUpdatedBy,omitempty"`
	EgressACLs                 []interface{} `json:"egressACLs,omitempty"`
	EgressQos                  interface{}   `json:"egressQos,omitempty"`
	FipACLs                    []interface{} `json:"fipACLs,omitempty"`
	IngressACLs                []interface{} `json:"ingressACLs,omitempty"`
	IngressAdvFwd              []interface{} `json:"ingressAdvFwd,omitempty"`
	IngressExternalServiceACLs []interface{} `json:"ingressExternalServiceACLs,omitempty"`
	EntityScope                string        `json:"entityScope,omitempty"`
	Qos                        interface{}   `json:"qos,omitempty"`
	Stats                      interface{}   `json:"stats,omitempty"`
	ExternalID                 string        `json:"externalID,omitempty"`
}

// NewPolicyDecision returns a new *PolicyDecision
func NewPolicyDecision() *PolicyDecision {

	return &PolicyDecision{}
}

// Identity returns the Identity of the object.
func (o *PolicyDecision) Identity() bambou.Identity {

	return PolicyDecisionIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *PolicyDecision) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *PolicyDecision) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the PolicyDecision from the server
func (o *PolicyDecision) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the PolicyDecision into the server
func (o *PolicyDecision) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the PolicyDecision from the server
func (o *PolicyDecision) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the PolicyDecision
func (o *PolicyDecision) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the PolicyDecision
func (o *PolicyDecision) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the PolicyDecision
func (o *PolicyDecision) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the PolicyDecision
func (o *PolicyDecision) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// QOSs retrieves the list of child QOSs of the PolicyDecision
func (o *PolicyDecision) QOSs(info *bambou.FetchingInfo) (QOSsList, *bambou.Error) {

	var list QOSsList
	err := bambou.CurrentSession().FetchChildren(o, QOSIdentity, &list, info)
	return list, err
}
