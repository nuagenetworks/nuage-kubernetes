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

// VRSMetricsIdentity represents the Identity of the object
var VRSMetricsIdentity = bambou.Identity{
	Name:     "vrsmetrics",
	Category: "vrsmetrics",
}

// VRSMetricsList represents a list of VRSMetrics
type VRSMetricsList []*VRSMetrics

// VRSMetricsAncestor is the interface of an ancestor of a VRSMetrics must implement.
type VRSMetricsAncestor interface {
	VRSMetrics(*bambou.FetchingInfo) (VRSMetricsList, *bambou.Error)
	CreateVRSMetrics(*VRSMetrics) *bambou.Error
}

// VRSMetrics represents the model of a vrsmetrics
type VRSMetrics struct {
	ID                            string  `json:"ID,omitempty"`
	ParentID                      string  `json:"parentID,omitempty"`
	ParentType                    string  `json:"parentType,omitempty"`
	Owner                         string  `json:"owner,omitempty"`
	ALUbr0Status                  bool    `json:"ALUbr0Status"`
	CPUUtilization                float64 `json:"CPUUtilization,omitempty"`
	VRSProcess                    bool    `json:"VRSProcess"`
	VRSVSCStatus                  bool    `json:"VRSVSCStatus"`
	LastUpdatedBy                 string  `json:"lastUpdatedBy,omitempty"`
	ReDeploy                      bool    `json:"reDeploy"`
	ReceivingMetrics              bool    `json:"receivingMetrics"`
	MemoryUtilization             float64 `json:"memoryUtilization,omitempty"`
	JesxmonProcess                bool    `json:"jesxmonProcess"`
	AgentName                     string  `json:"agentName,omitempty"`
	EntityScope                   string  `json:"entityScope,omitempty"`
	AssociatedVCenterHypervisorID string  `json:"associatedVCenterHypervisorID,omitempty"`
	ExternalID                    string  `json:"externalID,omitempty"`
}

// NewVRSMetrics returns a new *VRSMetrics
func NewVRSMetrics() *VRSMetrics {

	return &VRSMetrics{}
}

// Identity returns the Identity of the object.
func (o *VRSMetrics) Identity() bambou.Identity {

	return VRSMetricsIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VRSMetrics) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VRSMetrics) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VRSMetrics from the server
func (o *VRSMetrics) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VRSMetrics into the server
func (o *VRSMetrics) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VRSMetrics from the server
func (o *VRSMetrics) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}
