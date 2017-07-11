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

// VRSRedeploymentpolicyIdentity represents the Identity of the object
var VRSRedeploymentpolicyIdentity = bambou.Identity{
	Name:     "vrsredeploymentpolicy",
	Category: "vrsredeploymentpolicies",
}

// VRSRedeploymentpoliciesList represents a list of VRSRedeploymentpolicies
type VRSRedeploymentpoliciesList []*VRSRedeploymentpolicy

// VRSRedeploymentpoliciesAncestor is the interface of an ancestor of a VRSRedeploymentpolicy must implement.
type VRSRedeploymentpoliciesAncestor interface {
	VRSRedeploymentpolicies(*bambou.FetchingInfo) (VRSRedeploymentpoliciesList, *bambou.Error)
	CreateVRSRedeploymentpolicies(*VRSRedeploymentpolicy) *bambou.Error
}

// VRSRedeploymentpolicy represents the model of a vrsredeploymentpolicy
type VRSRedeploymentpolicy struct {
	ID                                   string  `json:"ID,omitempty"`
	ParentID                             string  `json:"parentID,omitempty"`
	ParentType                           string  `json:"parentType,omitempty"`
	Owner                                string  `json:"owner,omitempty"`
	ALUbr0StatusRedeploymentEnabled      bool    `json:"ALUbr0StatusRedeploymentEnabled"`
	CPUUtilizationRedeploymentEnabled    bool    `json:"CPUUtilizationRedeploymentEnabled"`
	CPUUtilizationThreshold              float64 `json:"CPUUtilizationThreshold,omitempty"`
	VRSCorrectiveActionDelay             int     `json:"VRSCorrectiveActionDelay,omitempty"`
	VRSProcessRedeploymentEnabled        bool    `json:"VRSProcessRedeploymentEnabled"`
	VRSVSCStatusRedeploymentEnabled      bool    `json:"VRSVSCStatusRedeploymentEnabled"`
	LastUpdatedBy                        string  `json:"lastUpdatedBy,omitempty"`
	RedeploymentDelay                    int     `json:"redeploymentDelay,omitempty"`
	MemoryUtilizationRedeploymentEnabled bool    `json:"memoryUtilizationRedeploymentEnabled"`
	MemoryUtilizationThreshold           float64 `json:"memoryUtilizationThreshold,omitempty"`
	DeploymentCountThreshold             int     `json:"deploymentCountThreshold,omitempty"`
	JesxmonProcessRedeploymentEnabled    bool    `json:"jesxmonProcessRedeploymentEnabled"`
	EntityScope                          string  `json:"entityScope,omitempty"`
	ExternalID                           string  `json:"externalID,omitempty"`
}

// NewVRSRedeploymentpolicy returns a new *VRSRedeploymentpolicy
func NewVRSRedeploymentpolicy() *VRSRedeploymentpolicy {

	return &VRSRedeploymentpolicy{}
}

// Identity returns the Identity of the object.
func (o *VRSRedeploymentpolicy) Identity() bambou.Identity {

	return VRSRedeploymentpolicyIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VRSRedeploymentpolicy) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VRSRedeploymentpolicy) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VRSRedeploymentpolicy from the server
func (o *VRSRedeploymentpolicy) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VRSRedeploymentpolicy into the server
func (o *VRSRedeploymentpolicy) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VRSRedeploymentpolicy from the server
func (o *VRSRedeploymentpolicy) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}
