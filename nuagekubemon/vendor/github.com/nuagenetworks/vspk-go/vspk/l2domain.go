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

// L2DomainIdentity represents the Identity of the object
var L2DomainIdentity = bambou.Identity{
	Name:     "l2domain",
	Category: "l2domains",
}

// L2DomainsList represents a list of L2Domains
type L2DomainsList []*L2Domain

// L2DomainsAncestor is the interface that an ancestor of a L2Domain must implement.
// An Ancestor is defined as an entity that has L2Domain as a descendant.
// An Ancestor can get a list of its child L2Domains, but not necessarily create one.
type L2DomainsAncestor interface {
	L2Domains(*bambou.FetchingInfo) (L2DomainsList, *bambou.Error)
}

// L2DomainsParent is the interface that a parent of a L2Domain must implement.
// A Parent is defined as an entity that has L2Domain as a child.
// A Parent is an Ancestor which can create a L2Domain.
type L2DomainsParent interface {
	L2DomainsAncestor
	CreateL2Domain(*L2Domain) *bambou.Error
}

// L2Domain represents the model of a l2domain
type L2Domain struct {
	ID                                string `json:"ID,omitempty"`
	ParentID                          string `json:"parentID,omitempty"`
	ParentType                        string `json:"parentType,omitempty"`
	Owner                             string `json:"owner,omitempty"`
	DHCPManaged                       bool   `json:"DHCPManaged"`
	DPI                               string `json:"DPI,omitempty"`
	IPType                            string `json:"IPType,omitempty"`
	MaintenanceMode                   string `json:"maintenanceMode,omitempty"`
	Name                              string `json:"name,omitempty"`
	LastUpdatedBy                     string `json:"lastUpdatedBy,omitempty"`
	Gateway                           string `json:"gateway,omitempty"`
	GatewayMACAddress                 string `json:"gatewayMACAddress,omitempty"`
	Address                           string `json:"address,omitempty"`
	TemplateID                        string `json:"templateID,omitempty"`
	ServiceID                         int    `json:"serviceID,omitempty"`
	Description                       string `json:"description,omitempty"`
	Netmask                           string `json:"netmask,omitempty"`
	VnId                              int    `json:"vnId,omitempty"`
	Encryption                        string `json:"encryption,omitempty"`
	EntityScope                       string `json:"entityScope,omitempty"`
	PolicyChangeStatus                string `json:"policyChangeStatus,omitempty"`
	RouteDistinguisher                string `json:"routeDistinguisher,omitempty"`
	RouteTarget                       string `json:"routeTarget,omitempty"`
	UplinkPreference                  string `json:"uplinkPreference,omitempty"`
	AssociatedMulticastChannelMapID   string `json:"associatedMulticastChannelMapID,omitempty"`
	AssociatedSharedNetworkResourceID string `json:"associatedSharedNetworkResourceID,omitempty"`
	Stretched                         bool   `json:"stretched"`
	Multicast                         string `json:"multicast,omitempty"`
	ExternalID                        string `json:"externalID,omitempty"`
}

// NewL2Domain returns a new *L2Domain
func NewL2Domain() *L2Domain {

	return &L2Domain{
		DPI:             "DISABLED",
		MaintenanceMode: "DISABLED",
	}
}

// Identity returns the Identity of the object.
func (o *L2Domain) Identity() bambou.Identity {

	return L2DomainIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *L2Domain) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *L2Domain) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the L2Domain from the server
func (o *L2Domain) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the L2Domain into the server
func (o *L2Domain) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the L2Domain from the server
func (o *L2Domain) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// TCAs retrieves the list of child TCAs of the L2Domain
func (o *L2Domain) TCAs(info *bambou.FetchingInfo) (TCAsList, *bambou.Error) {

	var list TCAsList
	err := bambou.CurrentSession().FetchChildren(o, TCAIdentity, &list, info)
	return list, err
}

// CreateTCA creates a new child TCA under the L2Domain
func (o *L2Domain) CreateTCA(child *TCA) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// AddressRanges retrieves the list of child AddressRanges of the L2Domain
func (o *L2Domain) AddressRanges(info *bambou.FetchingInfo) (AddressRangesList, *bambou.Error) {

	var list AddressRangesList
	err := bambou.CurrentSession().FetchChildren(o, AddressRangeIdentity, &list, info)
	return list, err
}

// RedirectionTargets retrieves the list of child RedirectionTargets of the L2Domain
func (o *L2Domain) RedirectionTargets(info *bambou.FetchingInfo) (RedirectionTargetsList, *bambou.Error) {

	var list RedirectionTargetsList
	err := bambou.CurrentSession().FetchChildren(o, RedirectionTargetIdentity, &list, info)
	return list, err
}

// CreateRedirectionTarget creates a new child RedirectionTarget under the L2Domain
func (o *L2Domain) CreateRedirectionTarget(child *RedirectionTarget) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Permissions retrieves the list of child Permissions of the L2Domain
func (o *L2Domain) Permissions(info *bambou.FetchingInfo) (PermissionsList, *bambou.Error) {

	var list PermissionsList
	err := bambou.CurrentSession().FetchChildren(o, PermissionIdentity, &list, info)
	return list, err
}

// CreatePermission creates a new child Permission under the L2Domain
func (o *L2Domain) CreatePermission(child *Permission) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the L2Domain
func (o *L2Domain) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the L2Domain
func (o *L2Domain) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EgressACLEntryTemplates retrieves the list of child EgressACLEntryTemplates of the L2Domain
func (o *L2Domain) EgressACLEntryTemplates(info *bambou.FetchingInfo) (EgressACLEntryTemplatesList, *bambou.Error) {

	var list EgressACLEntryTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, EgressACLEntryTemplateIdentity, &list, info)
	return list, err
}

// EgressACLTemplates retrieves the list of child EgressACLTemplates of the L2Domain
func (o *L2Domain) EgressACLTemplates(info *bambou.FetchingInfo) (EgressACLTemplatesList, *bambou.Error) {

	var list EgressACLTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, EgressACLTemplateIdentity, &list, info)
	return list, err
}

// CreateEgressACLTemplate creates a new child EgressACLTemplate under the L2Domain
func (o *L2Domain) CreateEgressACLTemplate(child *EgressACLTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DHCPOptions retrieves the list of child DHCPOptions of the L2Domain
func (o *L2Domain) DHCPOptions(info *bambou.FetchingInfo) (DHCPOptionsList, *bambou.Error) {

	var list DHCPOptionsList
	err := bambou.CurrentSession().FetchChildren(o, DHCPOptionIdentity, &list, info)
	return list, err
}

// CreateDHCPOption creates a new child DHCPOption under the L2Domain
func (o *L2Domain) CreateDHCPOption(child *DHCPOption) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the L2Domain
func (o *L2Domain) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the L2Domain
func (o *L2Domain) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMs retrieves the list of child VMs of the L2Domain
func (o *L2Domain) VMs(info *bambou.FetchingInfo) (VMsList, *bambou.Error) {

	var list VMsList
	err := bambou.CurrentSession().FetchChildren(o, VMIdentity, &list, info)
	return list, err
}

// VMInterfaces retrieves the list of child VMInterfaces of the L2Domain
func (o *L2Domain) VMInterfaces(info *bambou.FetchingInfo) (VMInterfacesList, *bambou.Error) {

	var list VMInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, VMInterfaceIdentity, &list, info)
	return list, err
}

// IngressACLEntryTemplates retrieves the list of child IngressACLEntryTemplates of the L2Domain
func (o *L2Domain) IngressACLEntryTemplates(info *bambou.FetchingInfo) (IngressACLEntryTemplatesList, *bambou.Error) {

	var list IngressACLEntryTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressACLEntryTemplateIdentity, &list, info)
	return list, err
}

// IngressACLTemplates retrieves the list of child IngressACLTemplates of the L2Domain
func (o *L2Domain) IngressACLTemplates(info *bambou.FetchingInfo) (IngressACLTemplatesList, *bambou.Error) {

	var list IngressACLTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressACLTemplateIdentity, &list, info)
	return list, err
}

// CreateIngressACLTemplate creates a new child IngressACLTemplate under the L2Domain
func (o *L2Domain) CreateIngressACLTemplate(child *IngressACLTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IngressAdvFwdTemplates retrieves the list of child IngressAdvFwdTemplates of the L2Domain
func (o *L2Domain) IngressAdvFwdTemplates(info *bambou.FetchingInfo) (IngressAdvFwdTemplatesList, *bambou.Error) {

	var list IngressAdvFwdTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressAdvFwdTemplateIdentity, &list, info)
	return list, err
}

// CreateIngressAdvFwdTemplate creates a new child IngressAdvFwdTemplate under the L2Domain
func (o *L2Domain) CreateIngressAdvFwdTemplate(child *IngressAdvFwdTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IngressExternalServiceTemplates retrieves the list of child IngressExternalServiceTemplates of the L2Domain
func (o *L2Domain) IngressExternalServiceTemplates(info *bambou.FetchingInfo) (IngressExternalServiceTemplatesList, *bambou.Error) {

	var list IngressExternalServiceTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressExternalServiceTemplateIdentity, &list, info)
	return list, err
}

// CreateIngressExternalServiceTemplate creates a new child IngressExternalServiceTemplate under the L2Domain
func (o *L2Domain) CreateIngressExternalServiceTemplate(child *IngressExternalServiceTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CreateJob creates a new child Job under the L2Domain
func (o *L2Domain) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PolicyGroups retrieves the list of child PolicyGroups of the L2Domain
func (o *L2Domain) PolicyGroups(info *bambou.FetchingInfo) (PolicyGroupsList, *bambou.Error) {

	var list PolicyGroupsList
	err := bambou.CurrentSession().FetchChildren(o, PolicyGroupIdentity, &list, info)
	return list, err
}

// CreatePolicyGroup creates a new child PolicyGroup under the L2Domain
func (o *L2Domain) CreatePolicyGroup(child *PolicyGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Containers retrieves the list of child Containers of the L2Domain
func (o *L2Domain) Containers(info *bambou.FetchingInfo) (ContainersList, *bambou.Error) {

	var list ContainersList
	err := bambou.CurrentSession().FetchChildren(o, ContainerIdentity, &list, info)
	return list, err
}

// ContainerInterfaces retrieves the list of child ContainerInterfaces of the L2Domain
func (o *L2Domain) ContainerInterfaces(info *bambou.FetchingInfo) (ContainerInterfacesList, *bambou.Error) {

	var list ContainerInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, ContainerInterfaceIdentity, &list, info)
	return list, err
}

// QOSs retrieves the list of child QOSs of the L2Domain
func (o *L2Domain) QOSs(info *bambou.FetchingInfo) (QOSsList, *bambou.Error) {

	var list QOSsList
	err := bambou.CurrentSession().FetchChildren(o, QOSIdentity, &list, info)
	return list, err
}

// CreateQOS creates a new child QOS under the L2Domain
func (o *L2Domain) CreateQOS(child *QOS) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// HostInterfaces retrieves the list of child HostInterfaces of the L2Domain
func (o *L2Domain) HostInterfaces(info *bambou.FetchingInfo) (HostInterfacesList, *bambou.Error) {

	var list HostInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, HostInterfaceIdentity, &list, info)
	return list, err
}

// UplinkRDs retrieves the list of child UplinkRDs of the L2Domain
func (o *L2Domain) UplinkRDs(info *bambou.FetchingInfo) (UplinkRDsList, *bambou.Error) {

	var list UplinkRDsList
	err := bambou.CurrentSession().FetchChildren(o, UplinkRDIdentity, &list, info)
	return list, err
}

// VPNConnections retrieves the list of child VPNConnections of the L2Domain
func (o *L2Domain) VPNConnections(info *bambou.FetchingInfo) (VPNConnectionsList, *bambou.Error) {

	var list VPNConnectionsList
	err := bambou.CurrentSession().FetchChildren(o, VPNConnectionIdentity, &list, info)
	return list, err
}

// CreateVPNConnection creates a new child VPNConnection under the L2Domain
func (o *L2Domain) CreateVPNConnection(child *VPNConnection) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VPorts retrieves the list of child VPorts of the L2Domain
func (o *L2Domain) VPorts(info *bambou.FetchingInfo) (VPortsList, *bambou.Error) {

	var list VPortsList
	err := bambou.CurrentSession().FetchChildren(o, VPortIdentity, &list, info)
	return list, err
}

// CreateVPort creates a new child VPort under the L2Domain
func (o *L2Domain) CreateVPort(child *VPort) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Applicationperformancemanagementbindings retrieves the list of child Applicationperformancemanagementbindings of the L2Domain
func (o *L2Domain) Applicationperformancemanagementbindings(info *bambou.FetchingInfo) (ApplicationperformancemanagementbindingsList, *bambou.Error) {

	var list ApplicationperformancemanagementbindingsList
	err := bambou.CurrentSession().FetchChildren(o, ApplicationperformancemanagementbindingIdentity, &list, info)
	return list, err
}

// CreateApplicationperformancemanagementbinding creates a new child Applicationperformancemanagementbinding under the L2Domain
func (o *L2Domain) CreateApplicationperformancemanagementbinding(child *Applicationperformancemanagementbinding) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// BridgeInterfaces retrieves the list of child BridgeInterfaces of the L2Domain
func (o *L2Domain) BridgeInterfaces(info *bambou.FetchingInfo) (BridgeInterfacesList, *bambou.Error) {

	var list BridgeInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, BridgeInterfaceIdentity, &list, info)
	return list, err
}

// Groups retrieves the list of child Groups of the L2Domain
func (o *L2Domain) Groups(info *bambou.FetchingInfo) (GroupsList, *bambou.Error) {

	var list GroupsList
	err := bambou.CurrentSession().FetchChildren(o, GroupIdentity, &list, info)
	return list, err
}

// StaticRoutes retrieves the list of child StaticRoutes of the L2Domain
func (o *L2Domain) StaticRoutes(info *bambou.FetchingInfo) (StaticRoutesList, *bambou.Error) {

	var list StaticRoutesList
	err := bambou.CurrentSession().FetchChildren(o, StaticRouteIdentity, &list, info)
	return list, err
}

// CreateStaticRoute creates a new child StaticRoute under the L2Domain
func (o *L2Domain) CreateStaticRoute(child *StaticRoute) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Statistics retrieves the list of child Statistics of the L2Domain
func (o *L2Domain) Statistics(info *bambou.FetchingInfo) (StatisticsList, *bambou.Error) {

	var list StatisticsList
	err := bambou.CurrentSession().FetchChildren(o, StatisticsIdentity, &list, info)
	return list, err
}

// StatisticsPolicies retrieves the list of child StatisticsPolicies of the L2Domain
func (o *L2Domain) StatisticsPolicies(info *bambou.FetchingInfo) (StatisticsPoliciesList, *bambou.Error) {

	var list StatisticsPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, StatisticsPolicyIdentity, &list, info)
	return list, err
}

// CreateStatisticsPolicy creates a new child StatisticsPolicy under the L2Domain
func (o *L2Domain) CreateStatisticsPolicy(child *StatisticsPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the L2Domain
func (o *L2Domain) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
