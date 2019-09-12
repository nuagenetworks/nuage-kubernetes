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

// PerformanceMonitorIdentity represents the Identity of the object
var PerformanceMonitorIdentity = bambou.Identity{
	Name:     "performancemonitor",
	Category: "performancemonitors",
}

// PerformanceMonitorsList represents a list of PerformanceMonitors
type PerformanceMonitorsList []*PerformanceMonitor

// PerformanceMonitorsAncestor is the interface that an ancestor of a PerformanceMonitor must implement.
// An Ancestor is defined as an entity that has PerformanceMonitor as a descendant.
// An Ancestor can get a list of its child PerformanceMonitors, but not necessarily create one.
type PerformanceMonitorsAncestor interface {
	PerformanceMonitors(*bambou.FetchingInfo) (PerformanceMonitorsList, *bambou.Error)
}

// PerformanceMonitorsParent is the interface that a parent of a PerformanceMonitor must implement.
// A Parent is defined as an entity that has PerformanceMonitor as a child.
// A Parent is an Ancestor which can create a PerformanceMonitor.
type PerformanceMonitorsParent interface {
	PerformanceMonitorsAncestor
	CreatePerformanceMonitor(*PerformanceMonitor) *bambou.Error
}

// PerformanceMonitor represents the model of a performancemonitor
type PerformanceMonitor struct {
	ID              string `json:"ID,omitempty"`
	ParentID        string `json:"parentID,omitempty"`
	ParentType      string `json:"parentType,omitempty"`
	Owner           string `json:"owner,omitempty"`
	Name            string `json:"name,omitempty"`
	LastUpdatedBy   string `json:"lastUpdatedBy,omitempty"`
	PayloadSize     int    `json:"payloadSize,omitempty"`
	ReadOnly        bool   `json:"readOnly"`
	ServiceClass    string `json:"serviceClass,omitempty"`
	Description     string `json:"description,omitempty"`
	Interval        int    `json:"interval,omitempty"`
	EntityScope     string `json:"entityScope,omitempty"`
	HoldDownTimer   int    `json:"holdDownTimer,omitempty"`
	ProbeType       string `json:"probeType,omitempty"`
	NumberOfPackets int    `json:"numberOfPackets,omitempty"`
	ExternalID      string `json:"externalID,omitempty"`
}

// NewPerformanceMonitor returns a new *PerformanceMonitor
func NewPerformanceMonitor() *PerformanceMonitor {

	return &PerformanceMonitor{
		PayloadSize:     137,
		ReadOnly:        false,
		ServiceClass:    "H",
		Interval:        180,
		HoldDownTimer:   1000,
		ProbeType:       "ONEWAY",
		NumberOfPackets: 1,
	}
}

// Identity returns the Identity of the object.
func (o *PerformanceMonitor) Identity() bambou.Identity {

	return PerformanceMonitorIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *PerformanceMonitor) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *PerformanceMonitor) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the PerformanceMonitor from the server
func (o *PerformanceMonitor) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the PerformanceMonitor into the server
func (o *PerformanceMonitor) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the PerformanceMonitor from the server
func (o *PerformanceMonitor) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Tiers retrieves the list of child Tiers of the PerformanceMonitor
func (o *PerformanceMonitor) Tiers(info *bambou.FetchingInfo) (TiersList, *bambou.Error) {

	var list TiersList
	err := bambou.CurrentSession().FetchChildren(o, TierIdentity, &list, info)
	return list, err
}

// Applicationperformancemanagements retrieves the list of child Applicationperformancemanagements of the PerformanceMonitor
func (o *PerformanceMonitor) Applicationperformancemanagements(info *bambou.FetchingInfo) (ApplicationperformancemanagementsList, *bambou.Error) {

	var list ApplicationperformancemanagementsList
	err := bambou.CurrentSession().FetchChildren(o, ApplicationperformancemanagementIdentity, &list, info)
	return list, err
}
