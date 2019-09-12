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

// VMIdentity represents the Identity of the object
var VMIdentity = bambou.Identity{
	Name:     "vm",
	Category: "vms",
}

// VMsList represents a list of VMs
type VMsList []*VM

// VMsAncestor is the interface that an ancestor of a VM must implement.
// An Ancestor is defined as an entity that has VM as a descendant.
// An Ancestor can get a list of its child VMs, but not necessarily create one.
type VMsAncestor interface {
	VMs(*bambou.FetchingInfo) (VMsList, *bambou.Error)
}

// VMsParent is the interface that a parent of a VM must implement.
// A Parent is defined as an entity that has VM as a child.
// A Parent is an Ancestor which can create a VM.
type VMsParent interface {
	VMsAncestor
	CreateVM(*VM) *bambou.Error
}

// VM represents the model of a vm
type VM struct {
	ID                 string        `json:"ID,omitempty"`
	ParentID           string        `json:"parentID,omitempty"`
	ParentType         string        `json:"parentType,omitempty"`
	Owner              string        `json:"owner,omitempty"`
	L2DomainIDs        []interface{} `json:"l2DomainIDs,omitempty"`
	VRSID              string        `json:"VRSID,omitempty"`
	UUID               string        `json:"UUID,omitempty"`
	Name               string        `json:"name,omitempty"`
	LastUpdatedBy      string        `json:"lastUpdatedBy,omitempty"`
	ReasonType         string        `json:"reasonType,omitempty"`
	DeleteExpiry       int           `json:"deleteExpiry,omitempty"`
	DeleteMode         string        `json:"deleteMode,omitempty"`
	ResyncInfo         interface{}   `json:"resyncInfo,omitempty"`
	SiteIdentifier     string        `json:"siteIdentifier,omitempty"`
	EmbeddedMetadata   []interface{} `json:"embeddedMetadata,omitempty"`
	Interfaces         []interface{} `json:"interfaces,omitempty"`
	EnterpriseID       string        `json:"enterpriseID,omitempty"`
	EnterpriseName     string        `json:"enterpriseName,omitempty"`
	EntityScope        string        `json:"entityScope,omitempty"`
	DomainIDs          []interface{} `json:"domainIDs,omitempty"`
	ComputeProvisioned bool          `json:"computeProvisioned"`
	ZoneIDs            []interface{} `json:"zoneIDs,omitempty"`
	OrchestrationID    string        `json:"orchestrationID,omitempty"`
	VrsRawVersion      string        `json:"vrsRawVersion,omitempty"`
	VrsVersion         string        `json:"vrsVersion,omitempty"`
	UserID             string        `json:"userID,omitempty"`
	UserName           string        `json:"userName,omitempty"`
	Status             string        `json:"status,omitempty"`
	SubnetIDs          []interface{} `json:"subnetIDs,omitempty"`
	ExternalID         string        `json:"externalID,omitempty"`
	HypervisorIP       string        `json:"hypervisorIP,omitempty"`
}

// NewVM returns a new *VM
func NewVM() *VM {

	return &VM{}
}

// Identity returns the Identity of the object.
func (o *VM) Identity() bambou.Identity {

	return VMIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VM) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VM) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VM from the server
func (o *VM) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VM into the server
func (o *VM) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VM from the server
func (o *VM) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// VMResyncs retrieves the list of child VMResyncs of the VM
func (o *VM) VMResyncs(info *bambou.FetchingInfo) (VMResyncsList, *bambou.Error) {

	var list VMResyncsList
	err := bambou.CurrentSession().FetchChildren(o, VMResyncIdentity, &list, info)
	return list, err
}

// CreateVMResync creates a new child VMResync under the VM
func (o *VM) CreateVMResync(child *VMResync) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the VM
func (o *VM) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VM
func (o *VM) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Alarms retrieves the list of child Alarms of the VM
func (o *VM) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VM
func (o *VM) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VM
func (o *VM) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMInterfaces retrieves the list of child VMInterfaces of the VM
func (o *VM) VMInterfaces(info *bambou.FetchingInfo) (VMInterfacesList, *bambou.Error) {

	var list VMInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, VMInterfaceIdentity, &list, info)
	return list, err
}

// CreateVMInterface creates a new child VMInterface under the VM
func (o *VM) CreateVMInterface(child *VMInterface) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VRSs retrieves the list of child VRSs of the VM
func (o *VM) VRSs(info *bambou.FetchingInfo) (VRSsList, *bambou.Error) {

	var list VRSsList
	err := bambou.CurrentSession().FetchChildren(o, VRSIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the VM
func (o *VM) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
