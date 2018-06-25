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

// VNFThresholdPolicyIdentity represents the Identity of the object
var VNFThresholdPolicyIdentity = bambou.Identity{
	Name:     "vnfthresholdpolicy",
	Category: "vnfthresholdpolicies",
}

// VNFThresholdPoliciesList represents a list of VNFThresholdPolicies
type VNFThresholdPoliciesList []*VNFThresholdPolicy

// VNFThresholdPoliciesAncestor is the interface that an ancestor of a VNFThresholdPolicy must implement.
// An Ancestor is defined as an entity that has VNFThresholdPolicy as a descendant.
// An Ancestor can get a list of its child VNFThresholdPolicies, but not necessarily create one.
type VNFThresholdPoliciesAncestor interface {
	VNFThresholdPolicies(*bambou.FetchingInfo) (VNFThresholdPoliciesList, *bambou.Error)
}

// VNFThresholdPoliciesParent is the interface that a parent of a VNFThresholdPolicy must implement.
// A Parent is defined as an entity that has VNFThresholdPolicy as a child.
// A Parent is an Ancestor which can create a VNFThresholdPolicy.
type VNFThresholdPoliciesParent interface {
	VNFThresholdPoliciesAncestor
	CreateVNFThresholdPolicy(*VNFThresholdPolicy) *bambou.Error
}

// VNFThresholdPolicy represents the model of a vnfthresholdpolicy
type VNFThresholdPolicy struct {
	ID               string `json:"ID,omitempty"`
	ParentID         string `json:"parentID,omitempty"`
	ParentType       string `json:"parentType,omitempty"`
	Owner            string `json:"owner,omitempty"`
	CPUThreshold     int    `json:"CPUThreshold,omitempty"`
	Name             string `json:"name,omitempty"`
	Action           string `json:"action,omitempty"`
	MemoryThreshold  int    `json:"memoryThreshold,omitempty"`
	Description      string `json:"description,omitempty"`
	MinOccurrence    int    `json:"minOccurrence,omitempty"`
	MonitInterval    int    `json:"monitInterval,omitempty"`
	StorageThreshold int    `json:"storageThreshold,omitempty"`
}

// NewVNFThresholdPolicy returns a new *VNFThresholdPolicy
func NewVNFThresholdPolicy() *VNFThresholdPolicy {

	return &VNFThresholdPolicy{
		CPUThreshold:     80,
		Action:           "NONE",
		MemoryThreshold:  80,
		MinOccurrence:    5,
		MonitInterval:    10,
		StorageThreshold: 80,
	}
}

// Identity returns the Identity of the object.
func (o *VNFThresholdPolicy) Identity() bambou.Identity {

	return VNFThresholdPolicyIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VNFThresholdPolicy) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VNFThresholdPolicy) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VNFThresholdPolicy from the server
func (o *VNFThresholdPolicy) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VNFThresholdPolicy into the server
func (o *VNFThresholdPolicy) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VNFThresholdPolicy from the server
func (o *VNFThresholdPolicy) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}
