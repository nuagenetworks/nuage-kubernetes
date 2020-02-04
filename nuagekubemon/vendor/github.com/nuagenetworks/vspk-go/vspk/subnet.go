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

// SubnetIdentity represents the Identity of the object
var SubnetIdentity = bambou.Identity{
	Name:     "subnet",
	Category: "subnets",
}

// SubnetsList represents a list of Subnets
type SubnetsList []*Subnet

// SubnetsAncestor is the interface that an ancestor of a Subnet must implement.
// An Ancestor is defined as an entity that has Subnet as a descendant.
// An Ancestor can get a list of its child Subnets, but not necessarily create one.
type SubnetsAncestor interface {
	Subnets(*bambou.FetchingInfo) (SubnetsList, *bambou.Error)
}

// SubnetsParent is the interface that a parent of a Subnet must implement.
// A Parent is defined as an entity that has Subnet as a child.
// A Parent is an Ancestor which can create a Subnet.
type SubnetsParent interface {
	SubnetsAncestor
	CreateSubnet(*Subnet) *bambou.Error
}

// Subnet represents the model of a subnet
type Subnet struct {
	ID                                string        `json:"ID,omitempty"`
	ParentID                          string        `json:"parentID,omitempty"`
	ParentType                        string        `json:"parentType,omitempty"`
	Owner                             string        `json:"owner,omitempty"`
	PATEnabled                        string        `json:"PATEnabled,omitempty"`
	DHCPRelayStatus                   string        `json:"DHCPRelayStatus,omitempty"`
	DPI                               string        `json:"DPI,omitempty"`
	IPType                            string        `json:"IPType,omitempty"`
	IPv6Address                       string        `json:"IPv6Address,omitempty"`
	IPv6Gateway                       string        `json:"IPv6Gateway,omitempty"`
	EVPNEnabled                       bool          `json:"EVPNEnabled"`
	MaintenanceMode                   string        `json:"maintenanceMode,omitempty"`
	Name                              string        `json:"name,omitempty"`
	LastUpdatedBy                     string        `json:"lastUpdatedBy,omitempty"`
	Gateway                           string        `json:"gateway,omitempty"`
	GatewayMACAddress                 string        `json:"gatewayMACAddress,omitempty"`
	AccessRestrictionEnabled          bool          `json:"accessRestrictionEnabled"`
	Address                           string        `json:"address,omitempty"`
	Advertise                         bool          `json:"advertise"`
	TemplateID                        string        `json:"templateID,omitempty"`
	ServiceID                         int           `json:"serviceID,omitempty"`
	Description                       string        `json:"description,omitempty"`
	ResourceType                      string        `json:"resourceType,omitempty"`
	Netmask                           string        `json:"netmask,omitempty"`
	LinkLocalAddress                  string        `json:"linkLocalAddress,omitempty"`
	EmbeddedMetadata                  []interface{} `json:"embeddedMetadata,omitempty"`
	VnId                              int           `json:"vnId,omitempty"`
	EnableDHCPv4                      bool          `json:"enableDHCPv4"`
	EnableDHCPv6                      bool          `json:"enableDHCPv6"`
	Encryption                        string        `json:"encryption,omitempty"`
	Underlay                          bool          `json:"underlay"`
	UnderlayEnabled                   string        `json:"underlayEnabled,omitempty"`
	IngressReplicationEnabled         bool          `json:"ingressReplicationEnabled"`
	EntityScope                       string        `json:"entityScope,omitempty"`
	EntityState                       string        `json:"entityState,omitempty"`
	PolicyGroupID                     int           `json:"policyGroupID,omitempty"`
	Color                             int           `json:"color,omitempty"`
	DomainServiceLabel                string        `json:"domainServiceLabel,omitempty"`
	RouteDistinguisher                string        `json:"routeDistinguisher,omitempty"`
	RouteTarget                       string        `json:"routeTarget,omitempty"`
	SplitSubnet                       bool          `json:"splitSubnet"`
	ProxyARP                          bool          `json:"proxyARP"`
	VrrpIPv6BackupAddress             string        `json:"vrrpIPv6BackupAddress,omitempty"`
	UseGlobalMAC                      string        `json:"useGlobalMAC,omitempty"`
	AssociatedMulticastChannelMapID   string        `json:"associatedMulticastChannelMapID,omitempty"`
	AssociatedSharedNetworkResourceID string        `json:"associatedSharedNetworkResourceID,omitempty"`
	DualStackDynamicIPAllocation      bool          `json:"dualStackDynamicIPAllocation"`
	Public                            bool          `json:"public"`
	SubnetVLANID                      int           `json:"subnetVLANID,omitempty"`
	MultiHomeEnabled                  bool          `json:"multiHomeEnabled"`
	Multicast                         string        `json:"multicast,omitempty"`
	CustomerID                        int           `json:"customerID,omitempty"`
	ExternalID                        string        `json:"externalID,omitempty"`
}

// NewSubnet returns a new *Subnet
func NewSubnet() *Subnet {

	return &Subnet{
		PATEnabled:                   "INHERITED",
		DPI:                          "INHERITED",
		IPType:                       "IPV4",
		EVPNEnabled:                  true,
		MaintenanceMode:              "DISABLED",
		AccessRestrictionEnabled:     false,
		Advertise:                    true,
		ResourceType:                 "STANDARD",
		LinkLocalAddress:             "fe80::1",
		EnableDHCPv4:                 true,
		EnableDHCPv6:                 false,
		Encryption:                   "INHERITED",
		UnderlayEnabled:              "INHERITED",
		IngressReplicationEnabled:    false,
		Color:                        0,
		VrrpIPv6BackupAddress:        "fe80::ffff:ffff:ffff:ffff",
		UseGlobalMAC:                 "ENTERPRISE_DEFAULT",
		DualStackDynamicIPAllocation: true,
		MultiHomeEnabled:             false,
		Multicast:                    "INHERITED",
	}
}

// Identity returns the Identity of the object.
func (o *Subnet) Identity() bambou.Identity {

	return SubnetIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Subnet) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Subnet) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Subnet from the server
func (o *Subnet) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Subnet into the server
func (o *Subnet) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Subnet from the server
func (o *Subnet) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// PATIPEntries retrieves the list of child PATIPEntries of the Subnet
func (o *Subnet) PATIPEntries(info *bambou.FetchingInfo) (PATIPEntriesList, *bambou.Error) {

	var list PATIPEntriesList
	err := bambou.CurrentSession().FetchChildren(o, PATIPEntryIdentity, &list, info)
	return list, err
}

// TCAs retrieves the list of child TCAs of the Subnet
func (o *Subnet) TCAs(info *bambou.FetchingInfo) (TCAsList, *bambou.Error) {

	var list TCAsList
	err := bambou.CurrentSession().FetchChildren(o, TCAIdentity, &list, info)
	return list, err
}

// CreateTCA creates a new child TCA under the Subnet
func (o *Subnet) CreateTCA(child *TCA) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// AddressRanges retrieves the list of child AddressRanges of the Subnet
func (o *Subnet) AddressRanges(info *bambou.FetchingInfo) (AddressRangesList, *bambou.Error) {

	var list AddressRangesList
	err := bambou.CurrentSession().FetchChildren(o, AddressRangeIdentity, &list, info)
	return list, err
}

// CreateAddressRange creates a new child AddressRange under the Subnet
func (o *Subnet) CreateAddressRange(child *AddressRange) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DefaultGateways retrieves the list of child DefaultGateways of the Subnet
func (o *Subnet) DefaultGateways(info *bambou.FetchingInfo) (DefaultGatewaysList, *bambou.Error) {

	var list DefaultGatewaysList
	err := bambou.CurrentSession().FetchChildren(o, DefaultGatewayIdentity, &list, info)
	return list, err
}

// DeploymentFailures retrieves the list of child DeploymentFailures of the Subnet
func (o *Subnet) DeploymentFailures(info *bambou.FetchingInfo) (DeploymentFailuresList, *bambou.Error) {

	var list DeploymentFailuresList
	err := bambou.CurrentSession().FetchChildren(o, DeploymentFailureIdentity, &list, info)
	return list, err
}

// CreateDeploymentFailure creates a new child DeploymentFailure under the Subnet
func (o *Subnet) CreateDeploymentFailure(child *DeploymentFailure) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMResyncs retrieves the list of child VMResyncs of the Subnet
func (o *Subnet) VMResyncs(info *bambou.FetchingInfo) (VMResyncsList, *bambou.Error) {

	var list VMResyncsList
	err := bambou.CurrentSession().FetchChildren(o, VMResyncIdentity, &list, info)
	return list, err
}

// CreateVMResync creates a new child VMResync under the Subnet
func (o *Subnet) CreateVMResync(child *VMResync) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the Subnet
func (o *Subnet) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the Subnet
func (o *Subnet) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// BGPNeighbors retrieves the list of child BGPNeighbors of the Subnet
func (o *Subnet) BGPNeighbors(info *bambou.FetchingInfo) (BGPNeighborsList, *bambou.Error) {

	var list BGPNeighborsList
	err := bambou.CurrentSession().FetchChildren(o, BGPNeighborIdentity, &list, info)
	return list, err
}

// CreateBGPNeighbor creates a new child BGPNeighbor under the Subnet
func (o *Subnet) CreateBGPNeighbor(child *BGPNeighbor) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DHCPOptions retrieves the list of child DHCPOptions of the Subnet
func (o *Subnet) DHCPOptions(info *bambou.FetchingInfo) (DHCPOptionsList, *bambou.Error) {

	var list DHCPOptionsList
	err := bambou.CurrentSession().FetchChildren(o, DHCPOptionIdentity, &list, info)
	return list, err
}

// CreateDHCPOption creates a new child DHCPOption under the Subnet
func (o *Subnet) CreateDHCPOption(child *DHCPOption) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DHCPv6Options retrieves the list of child DHCPv6Options of the Subnet
func (o *Subnet) DHCPv6Options(info *bambou.FetchingInfo) (DHCPv6OptionsList, *bambou.Error) {

	var list DHCPv6OptionsList
	err := bambou.CurrentSession().FetchChildren(o, DHCPv6OptionIdentity, &list, info)
	return list, err
}

// CreateDHCPv6Option creates a new child DHCPv6Option under the Subnet
func (o *Subnet) CreateDHCPv6Option(child *DHCPv6Option) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VirtualIPs retrieves the list of child VirtualIPs of the Subnet
func (o *Subnet) VirtualIPs(info *bambou.FetchingInfo) (VirtualIPsList, *bambou.Error) {

	var list VirtualIPsList
	err := bambou.CurrentSession().FetchChildren(o, VirtualIPIdentity, &list, info)
	return list, err
}

// IKEGatewayConnections retrieves the list of child IKEGatewayConnections of the Subnet
func (o *Subnet) IKEGatewayConnections(info *bambou.FetchingInfo) (IKEGatewayConnectionsList, *bambou.Error) {

	var list IKEGatewayConnectionsList
	err := bambou.CurrentSession().FetchChildren(o, IKEGatewayConnectionIdentity, &list, info)
	return list, err
}

// AssignIKEGatewayConnections assigns the list of IKEGatewayConnections to the Subnet
func (o *Subnet) AssignIKEGatewayConnections(children IKEGatewayConnectionsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, IKEGatewayConnectionIdentity)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Subnet
func (o *Subnet) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Subnet
func (o *Subnet) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMs retrieves the list of child VMs of the Subnet
func (o *Subnet) VMs(info *bambou.FetchingInfo) (VMsList, *bambou.Error) {

	var list VMsList
	err := bambou.CurrentSession().FetchChildren(o, VMIdentity, &list, info)
	return list, err
}

// VMInterfaces retrieves the list of child VMInterfaces of the Subnet
func (o *Subnet) VMInterfaces(info *bambou.FetchingInfo) (VMInterfacesList, *bambou.Error) {

	var list VMInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, VMInterfaceIdentity, &list, info)
	return list, err
}

// EnterprisePermissions retrieves the list of child EnterprisePermissions of the Subnet
func (o *Subnet) EnterprisePermissions(info *bambou.FetchingInfo) (EnterprisePermissionsList, *bambou.Error) {

	var list EnterprisePermissionsList
	err := bambou.CurrentSession().FetchChildren(o, EnterprisePermissionIdentity, &list, info)
	return list, err
}

// CreateEnterprisePermission creates a new child EnterprisePermission under the Subnet
func (o *Subnet) CreateEnterprisePermission(child *EnterprisePermission) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Containers retrieves the list of child Containers of the Subnet
func (o *Subnet) Containers(info *bambou.FetchingInfo) (ContainersList, *bambou.Error) {

	var list ContainersList
	err := bambou.CurrentSession().FetchChildren(o, ContainerIdentity, &list, info)
	return list, err
}

// ContainerInterfaces retrieves the list of child ContainerInterfaces of the Subnet
func (o *Subnet) ContainerInterfaces(info *bambou.FetchingInfo) (ContainerInterfacesList, *bambou.Error) {

	var list ContainerInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, ContainerInterfaceIdentity, &list, info)
	return list, err
}

// ContainerResyncs retrieves the list of child ContainerResyncs of the Subnet
func (o *Subnet) ContainerResyncs(info *bambou.FetchingInfo) (ContainerResyncsList, *bambou.Error) {

	var list ContainerResyncsList
	err := bambou.CurrentSession().FetchChildren(o, ContainerResyncIdentity, &list, info)
	return list, err
}

// CreateContainerResync creates a new child ContainerResync under the Subnet
func (o *Subnet) CreateContainerResync(child *ContainerResync) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// QOSs retrieves the list of child QOSs of the Subnet
func (o *Subnet) QOSs(info *bambou.FetchingInfo) (QOSsList, *bambou.Error) {

	var list QOSsList
	err := bambou.CurrentSession().FetchChildren(o, QOSIdentity, &list, info)
	return list, err
}

// CreateQOS creates a new child QOS under the Subnet
func (o *Subnet) CreateQOS(child *QOS) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VPorts retrieves the list of child VPorts of the Subnet
func (o *Subnet) VPorts(info *bambou.FetchingInfo) (VPortsList, *bambou.Error) {

	var list VPortsList
	err := bambou.CurrentSession().FetchChildren(o, VPortIdentity, &list, info)
	return list, err
}

// CreateVPort creates a new child VPort under the Subnet
func (o *Subnet) CreateVPort(child *VPort) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IPReservations retrieves the list of child IPReservations of the Subnet
func (o *Subnet) IPReservations(info *bambou.FetchingInfo) (IPReservationsList, *bambou.Error) {

	var list IPReservationsList
	err := bambou.CurrentSession().FetchChildren(o, IPReservationIdentity, &list, info)
	return list, err
}

// CreateIPReservation creates a new child IPReservation under the Subnet
func (o *Subnet) CreateIPReservation(child *IPReservation) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// ProxyARPFilters retrieves the list of child ProxyARPFilters of the Subnet
func (o *Subnet) ProxyARPFilters(info *bambou.FetchingInfo) (ProxyARPFiltersList, *bambou.Error) {

	var list ProxyARPFiltersList
	err := bambou.CurrentSession().FetchChildren(o, ProxyARPFilterIdentity, &list, info)
	return list, err
}

// CreateProxyARPFilter creates a new child ProxyARPFilter under the Subnet
func (o *Subnet) CreateProxyARPFilter(child *ProxyARPFilter) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Statistics retrieves the list of child Statistics of the Subnet
func (o *Subnet) Statistics(info *bambou.FetchingInfo) (StatisticsList, *bambou.Error) {

	var list StatisticsList
	err := bambou.CurrentSession().FetchChildren(o, StatisticsIdentity, &list, info)
	return list, err
}

// StatisticsPolicies retrieves the list of child StatisticsPolicies of the Subnet
func (o *Subnet) StatisticsPolicies(info *bambou.FetchingInfo) (StatisticsPoliciesList, *bambou.Error) {

	var list StatisticsPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, StatisticsPolicyIdentity, &list, info)
	return list, err
}

// CreateStatisticsPolicy creates a new child StatisticsPolicy under the Subnet
func (o *Subnet) CreateStatisticsPolicy(child *StatisticsPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the Subnet
func (o *Subnet) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
