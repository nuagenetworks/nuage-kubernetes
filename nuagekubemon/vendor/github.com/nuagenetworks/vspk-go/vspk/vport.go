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

// VPortIdentity represents the Identity of the object
var VPortIdentity = bambou.Identity{
	Name:     "vport",
	Category: "vports",
}

// VPortsList represents a list of VPorts
type VPortsList []*VPort

// VPortsAncestor is the interface that an ancestor of a VPort must implement.
// An Ancestor is defined as an entity that has VPort as a descendant.
// An Ancestor can get a list of its child VPorts, but not necessarily create one.
type VPortsAncestor interface {
	VPorts(*bambou.FetchingInfo) (VPortsList, *bambou.Error)
}

// VPortsParent is the interface that a parent of a VPort must implement.
// A Parent is defined as an entity that has VPort as a child.
// A Parent is an Ancestor which can create a VPort.
type VPortsParent interface {
	VPortsAncestor
	CreateVPort(*VPort) *bambou.Error
}

// VPort represents the model of a vport
type VPort struct {
	ID                                  string `json:"ID,omitempty"`
	ParentID                            string `json:"parentID,omitempty"`
	ParentType                          string `json:"parentType,omitempty"`
	Owner                               string `json:"owner,omitempty"`
	VLANID                              string `json:"VLANID,omitempty"`
	DPI                                 string `json:"DPI,omitempty"`
	Name                                string `json:"name,omitempty"`
	HasAttachedInterfaces               bool   `json:"hasAttachedInterfaces"`
	LastUpdatedBy                       string `json:"lastUpdatedBy,omitempty"`
	Active                              bool   `json:"active"`
	AddressSpoofing                     string `json:"addressSpoofing,omitempty"`
	Description                         string `json:"description,omitempty"`
	EntityScope                         string `json:"entityScope,omitempty"`
	DomainID                            string `json:"domainID,omitempty"`
	ZoneID                              string `json:"zoneID,omitempty"`
	OperationalState                    string `json:"operationalState,omitempty"`
	AssociatedFloatingIPID              string `json:"associatedFloatingIPID,omitempty"`
	AssociatedMulticastChannelMapID     string `json:"associatedMulticastChannelMapID,omitempty"`
	AssociatedSendMulticastChannelMapID string `json:"associatedSendMulticastChannelMapID,omitempty"`
	MultiNICVPortID                     string `json:"multiNICVPortID,omitempty"`
	Multicast                           string `json:"multicast,omitempty"`
	ExternalID                          string `json:"externalID,omitempty"`
	Type                                string `json:"type,omitempty"`
	SystemType                          string `json:"systemType,omitempty"`
}

// NewVPort returns a new *VPort
func NewVPort() *VPort {

	return &VPort{
		DPI:              "INHERITED",
		AddressSpoofing:  "INHERITED",
		OperationalState: "INIT",
		Multicast:        "INHERITED",
		Type:             "VM",
	}
}

// Identity returns the Identity of the object.
func (o *VPort) Identity() bambou.Identity {

	return VPortIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VPort) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VPort) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VPort from the server
func (o *VPort) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VPort into the server
func (o *VPort) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VPort from the server
func (o *VPort) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// TCAs retrieves the list of child TCAs of the VPort
func (o *VPort) TCAs(info *bambou.FetchingInfo) (TCAsList, *bambou.Error) {

	var list TCAsList
	err := bambou.CurrentSession().FetchChildren(o, TCAIdentity, &list, info)
	return list, err
}

// CreateTCA creates a new child TCA under the VPort
func (o *VPort) CreateTCA(child *TCA) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// RedirectionTargets retrieves the list of child RedirectionTargets of the VPort
func (o *VPort) RedirectionTargets(info *bambou.FetchingInfo) (RedirectionTargetsList, *bambou.Error) {

	var list RedirectionTargetsList
	err := bambou.CurrentSession().FetchChildren(o, RedirectionTargetIdentity, &list, info)
	return list, err
}

// AssignRedirectionTargets assigns the list of RedirectionTargets to the VPort
func (o *VPort) AssignRedirectionTargets(children RedirectionTargetsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, RedirectionTargetIdentity)
}

// Metadatas retrieves the list of child Metadatas of the VPort
func (o *VPort) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VPort
func (o *VPort) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// AggregateMetadatas retrieves the list of child AggregateMetadatas of the VPort
func (o *VPort) AggregateMetadatas(info *bambou.FetchingInfo) (AggregateMetadatasList, *bambou.Error) {

	var list AggregateMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, AggregateMetadataIdentity, &list, info)
	return list, err
}

// DHCPOptions retrieves the list of child DHCPOptions of the VPort
func (o *VPort) DHCPOptions(info *bambou.FetchingInfo) (DHCPOptionsList, *bambou.Error) {

	var list DHCPOptionsList
	err := bambou.CurrentSession().FetchChildren(o, DHCPOptionIdentity, &list, info)
	return list, err
}

// CreateDHCPOption creates a new child DHCPOption under the VPort
func (o *VPort) CreateDHCPOption(child *DHCPOption) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VirtualIPs retrieves the list of child VirtualIPs of the VPort
func (o *VPort) VirtualIPs(info *bambou.FetchingInfo) (VirtualIPsList, *bambou.Error) {

	var list VirtualIPsList
	err := bambou.CurrentSession().FetchChildren(o, VirtualIPIdentity, &list, info)
	return list, err
}

// CreateVirtualIP creates a new child VirtualIP under the VPort
func (o *VPort) CreateVirtualIP(child *VirtualIP) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Alarms retrieves the list of child Alarms of the VPort
func (o *VPort) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VPort
func (o *VPort) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VPort
func (o *VPort) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMs retrieves the list of child VMs of the VPort
func (o *VPort) VMs(info *bambou.FetchingInfo) (VMsList, *bambou.Error) {

	var list VMsList
	err := bambou.CurrentSession().FetchChildren(o, VMIdentity, &list, info)
	return list, err
}

// VMInterfaces retrieves the list of child VMInterfaces of the VPort
func (o *VPort) VMInterfaces(info *bambou.FetchingInfo) (VMInterfacesList, *bambou.Error) {

	var list VMInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, VMInterfaceIdentity, &list, info)
	return list, err
}

// PolicyGroups retrieves the list of child PolicyGroups of the VPort
func (o *VPort) PolicyGroups(info *bambou.FetchingInfo) (PolicyGroupsList, *bambou.Error) {

	var list PolicyGroupsList
	err := bambou.CurrentSession().FetchChildren(o, PolicyGroupIdentity, &list, info)
	return list, err
}

// AssignPolicyGroups assigns the list of PolicyGroups to the VPort
func (o *VPort) AssignPolicyGroups(children PolicyGroupsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, PolicyGroupIdentity)
}

// Containers retrieves the list of child Containers of the VPort
func (o *VPort) Containers(info *bambou.FetchingInfo) (ContainersList, *bambou.Error) {

	var list ContainersList
	err := bambou.CurrentSession().FetchChildren(o, ContainerIdentity, &list, info)
	return list, err
}

// ContainerInterfaces retrieves the list of child ContainerInterfaces of the VPort
func (o *VPort) ContainerInterfaces(info *bambou.FetchingInfo) (ContainerInterfacesList, *bambou.Error) {

	var list ContainerInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, ContainerInterfaceIdentity, &list, info)
	return list, err
}

// PortMappings retrieves the list of child PortMappings of the VPort
func (o *VPort) PortMappings(info *bambou.FetchingInfo) (PortMappingsList, *bambou.Error) {

	var list PortMappingsList
	err := bambou.CurrentSession().FetchChildren(o, PortMappingIdentity, &list, info)
	return list, err
}

// QOSs retrieves the list of child QOSs of the VPort
func (o *VPort) QOSs(info *bambou.FetchingInfo) (QOSsList, *bambou.Error) {

	var list QOSsList
	err := bambou.CurrentSession().FetchChildren(o, QOSIdentity, &list, info)
	return list, err
}

// CreateQOS creates a new child QOS under the VPort
func (o *VPort) CreateQOS(child *QOS) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// HostInterfaces retrieves the list of child HostInterfaces of the VPort
func (o *VPort) HostInterfaces(info *bambou.FetchingInfo) (HostInterfacesList, *bambou.Error) {

	var list HostInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, HostInterfaceIdentity, &list, info)
	return list, err
}

// CreateHostInterface creates a new child HostInterface under the VPort
func (o *VPort) CreateHostInterface(child *HostInterface) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VPortMirrors retrieves the list of child VPortMirrors of the VPort
func (o *VPort) VPortMirrors(info *bambou.FetchingInfo) (VPortMirrorsList, *bambou.Error) {

	var list VPortMirrorsList
	err := bambou.CurrentSession().FetchChildren(o, VPortMirrorIdentity, &list, info)
	return list, err
}

// CreateVPortMirror creates a new child VPortMirror under the VPort
func (o *VPort) CreateVPortMirror(child *VPortMirror) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Applicationperformancemanagements retrieves the list of child Applicationperformancemanagements of the VPort
func (o *VPort) Applicationperformancemanagements(info *bambou.FetchingInfo) (ApplicationperformancemanagementsList, *bambou.Error) {

	var list ApplicationperformancemanagementsList
	err := bambou.CurrentSession().FetchChildren(o, ApplicationperformancemanagementIdentity, &list, info)
	return list, err
}

// AssignApplicationperformancemanagements assigns the list of Applicationperformancemanagements to the VPort
func (o *VPort) AssignApplicationperformancemanagements(children ApplicationperformancemanagementsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, ApplicationperformancemanagementIdentity)
}

// BridgeInterfaces retrieves the list of child BridgeInterfaces of the VPort
func (o *VPort) BridgeInterfaces(info *bambou.FetchingInfo) (BridgeInterfacesList, *bambou.Error) {

	var list BridgeInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, BridgeInterfaceIdentity, &list, info)
	return list, err
}

// CreateBridgeInterface creates a new child BridgeInterface under the VPort
func (o *VPort) CreateBridgeInterface(child *BridgeInterface) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VRSs retrieves the list of child VRSs of the VPort
func (o *VPort) VRSs(info *bambou.FetchingInfo) (VRSsList, *bambou.Error) {

	var list VRSsList
	err := bambou.CurrentSession().FetchChildren(o, VRSIdentity, &list, info)
	return list, err
}

// Statistics retrieves the list of child Statistics of the VPort
func (o *VPort) Statistics(info *bambou.FetchingInfo) (StatisticsList, *bambou.Error) {

	var list StatisticsList
	err := bambou.CurrentSession().FetchChildren(o, StatisticsIdentity, &list, info)
	return list, err
}

// StatisticsPolicies retrieves the list of child StatisticsPolicies of the VPort
func (o *VPort) StatisticsPolicies(info *bambou.FetchingInfo) (StatisticsPoliciesList, *bambou.Error) {

	var list StatisticsPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, StatisticsPolicyIdentity, &list, info)
	return list, err
}

// CreateStatisticsPolicy creates a new child StatisticsPolicy under the VPort
func (o *VPort) CreateStatisticsPolicy(child *StatisticsPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the VPort
func (o *VPort) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
