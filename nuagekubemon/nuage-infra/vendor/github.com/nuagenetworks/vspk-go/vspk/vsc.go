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

// VSCIdentity represents the Identity of the object
var VSCIdentity = bambou.Identity{
	Name:     "vsc",
	Category: "vscs",
}

// VSCsList represents a list of VSCs
type VSCsList []*VSC

// VSCsAncestor is the interface that an ancestor of a VSC must implement.
// An Ancestor is defined as an entity that has VSC as a descendant.
// An Ancestor can get a list of its child VSCs, but not necessarily create one.
type VSCsAncestor interface {
	VSCs(*bambou.FetchingInfo) (VSCsList, *bambou.Error)
}

// VSCsParent is the interface that a parent of a VSC must implement.
// A Parent is defined as an entity that has VSC as a child.
// A Parent is an Ancestor which can create a VSC.
type VSCsParent interface {
	VSCsAncestor
	CreateVSC(*VSC) *bambou.Error
}

// VSC represents the model of a vsc
type VSC struct {
	ID                          string        `json:"ID,omitempty"`
	ParentID                    string        `json:"parentID,omitempty"`
	ParentType                  string        `json:"parentType,omitempty"`
	Owner                       string        `json:"owner,omitempty"`
	Name                        string        `json:"name,omitempty"`
	ManagementIP                string        `json:"managementIP,omitempty"`
	LastStateChange             int           `json:"lastStateChange,omitempty"`
	LastUpdatedBy               string        `json:"lastUpdatedBy,omitempty"`
	Address                     string        `json:"address,omitempty"`
	PeakCPUUsage                float64       `json:"peakCPUUsage,omitempty"`
	PeakMemoryUsage             float64       `json:"peakMemoryUsage,omitempty"`
	Description                 string        `json:"description,omitempty"`
	Messages                    []interface{} `json:"messages,omitempty"`
	Disks                       []interface{} `json:"disks,omitempty"`
	AlreadyMarkedForUnavailable bool          `json:"alreadyMarkedForUnavailable"`
	UnavailableTimestamp        int           `json:"unavailableTimestamp,omitempty"`
	EntityScope                 string        `json:"entityScope,omitempty"`
	Location                    string        `json:"location,omitempty"`
	ProductVersion              string        `json:"productVersion,omitempty"`
	Vsds                        []interface{} `json:"vsds,omitempty"`
	Status                      string        `json:"status,omitempty"`
	CurrentCPUUsage             float64       `json:"currentCPUUsage,omitempty"`
	CurrentMemoryUsage          float64       `json:"currentMemoryUsage,omitempty"`
	AverageCPUUsage             float64       `json:"averageCPUUsage,omitempty"`
	AverageMemoryUsage          float64       `json:"averageMemoryUsage,omitempty"`
	ExternalID                  string        `json:"externalID,omitempty"`
}

// NewVSC returns a new *VSC
func NewVSC() *VSC {

	return &VSC{}
}

// Identity returns the Identity of the object.
func (o *VSC) Identity() bambou.Identity {

	return VSCIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VSC) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VSC) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VSC from the server
func (o *VSC) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VSC into the server
func (o *VSC) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VSC from the server
func (o *VSC) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the VSC
func (o *VSC) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VSC
func (o *VSC) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// BGPPeers retrieves the list of child BGPPeers of the VSC
func (o *VSC) BGPPeers(info *bambou.FetchingInfo) (BGPPeersList, *bambou.Error) {

	var list BGPPeersList
	err := bambou.CurrentSession().FetchChildren(o, BGPPeerIdentity, &list, info)
	return list, err
}

// Alarms retrieves the list of child Alarms of the VSC
func (o *VSC) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VSC
func (o *VSC) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VSC
func (o *VSC) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CreateJob creates a new child Job under the VSC
func (o *VSC) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// MonitoringPorts retrieves the list of child MonitoringPorts of the VSC
func (o *VSC) MonitoringPorts(info *bambou.FetchingInfo) (MonitoringPortsList, *bambou.Error) {

	var list MonitoringPortsList
	err := bambou.CurrentSession().FetchChildren(o, MonitoringPortIdentity, &list, info)
	return list, err
}

// VRSs retrieves the list of child VRSs of the VSC
func (o *VSC) VRSs(info *bambou.FetchingInfo) (VRSsList, *bambou.Error) {

	var list VRSsList
	err := bambou.CurrentSession().FetchChildren(o, VRSIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the VSC
func (o *VSC) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
