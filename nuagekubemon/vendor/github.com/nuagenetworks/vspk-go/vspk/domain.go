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

// DomainIdentity represents the Identity of the object
var DomainIdentity = bambou.Identity{
	Name:     "domain",
	Category: "domains",
}

// DomainsList represents a list of Domains
type DomainsList []*Domain

// DomainsAncestor is the interface that an ancestor of a Domain must implement.
// An Ancestor is defined as an entity that has Domain as a descendant.
// An Ancestor can get a list of its child Domains, but not necessarily create one.
type DomainsAncestor interface {
	Domains(*bambou.FetchingInfo) (DomainsList, *bambou.Error)
}

// DomainsParent is the interface that a parent of a Domain must implement.
// A Parent is defined as an entity that has Domain as a child.
// A Parent is an Ancestor which can create a Domain.
type DomainsParent interface {
	DomainsAncestor
	CreateDomain(*Domain) *bambou.Error
}

// Domain represents the model of a domain
type Domain struct {
	ID                              string        `json:"ID,omitempty"`
	ParentID                        string        `json:"parentID,omitempty"`
	ParentType                      string        `json:"parentType,omitempty"`
	Owner                           string        `json:"owner,omitempty"`
	PATEnabled                      string        `json:"PATEnabled,omitempty"`
	ECMPCount                       int           `json:"ECMPCount,omitempty"`
	BGPEnabled                      bool          `json:"BGPEnabled"`
	DHCPBehavior                    string        `json:"DHCPBehavior,omitempty"`
	DHCPServerAddress               string        `json:"DHCPServerAddress,omitempty"`
	DPI                             string        `json:"DPI,omitempty"`
	LabelID                         int           `json:"labelID,omitempty"`
	BackHaulRouteDistinguisher      string        `json:"backHaulRouteDistinguisher,omitempty"`
	BackHaulRouteTarget             string        `json:"backHaulRouteTarget,omitempty"`
	BackHaulSubnetIPAddress         string        `json:"backHaulSubnetIPAddress,omitempty"`
	BackHaulSubnetMask              string        `json:"backHaulSubnetMask,omitempty"`
	BackHaulVNID                    int           `json:"backHaulVNID,omitempty"`
	MaintenanceMode                 string        `json:"maintenanceMode,omitempty"`
	Name                            string        `json:"name,omitempty"`
	LastUpdatedBy                   string        `json:"lastUpdatedBy,omitempty"`
	LeakingEnabled                  bool          `json:"leakingEnabled"`
	SecondaryDHCPServerAddress      string        `json:"secondaryDHCPServerAddress,omitempty"`
	TemplateID                      string        `json:"templateID,omitempty"`
	PermittedAction                 string        `json:"permittedAction,omitempty"`
	ServiceID                       int           `json:"serviceID,omitempty"`
	Description                     string        `json:"description,omitempty"`
	DhcpServerAddresses             []interface{} `json:"dhcpServerAddresses,omitempty"`
	GlobalRoutingEnabled            bool          `json:"globalRoutingEnabled"`
	ImportRouteTarget               string        `json:"importRouteTarget,omitempty"`
	Encryption                      string        `json:"encryption,omitempty"`
	UnderlayEnabled                 string        `json:"underlayEnabled,omitempty"`
	EntityScope                     string        `json:"entityScope,omitempty"`
	PolicyChangeStatus              string        `json:"policyChangeStatus,omitempty"`
	DomainID                        int           `json:"domainID,omitempty"`
	DomainVLANID                    int           `json:"domainVLANID,omitempty"`
	RouteDistinguisher              string        `json:"routeDistinguisher,omitempty"`
	RouteTarget                     string        `json:"routeTarget,omitempty"`
	UplinkPreference                string        `json:"uplinkPreference,omitempty"`
	ApplicationDeploymentPolicy     string        `json:"applicationDeploymentPolicy,omitempty"`
	AssociatedBGPProfileID          string        `json:"associatedBGPProfileID,omitempty"`
	AssociatedMulticastChannelMapID string        `json:"associatedMulticastChannelMapID,omitempty"`
	AssociatedPATMapperID           string        `json:"associatedPATMapperID,omitempty"`
	Stretched                       bool          `json:"stretched"`
	Multicast                       string        `json:"multicast,omitempty"`
	TunnelType                      string        `json:"tunnelType,omitempty"`
	CustomerID                      int           `json:"customerID,omitempty"`
	ExportRouteTarget               string        `json:"exportRouteTarget,omitempty"`
	ExternalID                      string        `json:"externalID,omitempty"`
}

// NewDomain returns a new *Domain
func NewDomain() *Domain {

	return &Domain{
		PATEnabled:                  "INHERITED",
		DHCPBehavior:                "CONSUME",
		DPI:                         "DISABLED",
		MaintenanceMode:             "DISABLED",
		ApplicationDeploymentPolicy: "ZONE",
		TunnelType:                  "DC_DEFAULT",
	}
}

// Identity returns the Identity of the object.
func (o *Domain) Identity() bambou.Identity {

	return DomainIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Domain) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Domain) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Domain from the server
func (o *Domain) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Domain into the server
func (o *Domain) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Domain from the server
func (o *Domain) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// TCAs retrieves the list of child TCAs of the Domain
func (o *Domain) TCAs(info *bambou.FetchingInfo) (TCAsList, *bambou.Error) {

	var list TCAsList
	err := bambou.CurrentSession().FetchChildren(o, TCAIdentity, &list, info)
	return list, err
}

// CreateTCA creates a new child TCA under the Domain
func (o *Domain) CreateTCA(child *TCA) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// RedirectionTargets retrieves the list of child RedirectionTargets of the Domain
func (o *Domain) RedirectionTargets(info *bambou.FetchingInfo) (RedirectionTargetsList, *bambou.Error) {

	var list RedirectionTargetsList
	err := bambou.CurrentSession().FetchChildren(o, RedirectionTargetIdentity, &list, info)
	return list, err
}

// CreateRedirectionTarget creates a new child RedirectionTarget under the Domain
func (o *Domain) CreateRedirectionTarget(child *RedirectionTarget) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Permissions retrieves the list of child Permissions of the Domain
func (o *Domain) Permissions(info *bambou.FetchingInfo) (PermissionsList, *bambou.Error) {

	var list PermissionsList
	err := bambou.CurrentSession().FetchChildren(o, PermissionIdentity, &list, info)
	return list, err
}

// CreatePermission creates a new child Permission under the Domain
func (o *Domain) CreatePermission(child *Permission) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the Domain
func (o *Domain) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the Domain
func (o *Domain) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EgressACLEntryTemplates retrieves the list of child EgressACLEntryTemplates of the Domain
func (o *Domain) EgressACLEntryTemplates(info *bambou.FetchingInfo) (EgressACLEntryTemplatesList, *bambou.Error) {

	var list EgressACLEntryTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, EgressACLEntryTemplateIdentity, &list, info)
	return list, err
}

// EgressACLTemplates retrieves the list of child EgressACLTemplates of the Domain
func (o *Domain) EgressACLTemplates(info *bambou.FetchingInfo) (EgressACLTemplatesList, *bambou.Error) {

	var list EgressACLTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, EgressACLTemplateIdentity, &list, info)
	return list, err
}

// CreateEgressACLTemplate creates a new child EgressACLTemplate under the Domain
func (o *Domain) CreateEgressACLTemplate(child *EgressACLTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DomainFIPAclTemplates retrieves the list of child DomainFIPAclTemplates of the Domain
func (o *Domain) DomainFIPAclTemplates(info *bambou.FetchingInfo) (DomainFIPAclTemplatesList, *bambou.Error) {

	var list DomainFIPAclTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, DomainFIPAclTemplateIdentity, &list, info)
	return list, err
}

// CreateDomainFIPAclTemplate creates a new child DomainFIPAclTemplate under the Domain
func (o *Domain) CreateDomainFIPAclTemplate(child *DomainFIPAclTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// FloatingIPACLTemplates retrieves the list of child FloatingIPACLTemplates of the Domain
func (o *Domain) FloatingIPACLTemplates(info *bambou.FetchingInfo) (FloatingIPACLTemplatesList, *bambou.Error) {

	var list FloatingIPACLTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, FloatingIPACLTemplateIdentity, &list, info)
	return list, err
}

// CreateFloatingIPACLTemplate creates a new child FloatingIPACLTemplate under the Domain
func (o *Domain) CreateFloatingIPACLTemplate(child *FloatingIPACLTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DHCPOptions retrieves the list of child DHCPOptions of the Domain
func (o *Domain) DHCPOptions(info *bambou.FetchingInfo) (DHCPOptionsList, *bambou.Error) {

	var list DHCPOptionsList
	err := bambou.CurrentSession().FetchChildren(o, DHCPOptionIdentity, &list, info)
	return list, err
}

// CreateDHCPOption creates a new child DHCPOption under the Domain
func (o *Domain) CreateDHCPOption(child *DHCPOption) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Links retrieves the list of child Links of the Domain
func (o *Domain) Links(info *bambou.FetchingInfo) (LinksList, *bambou.Error) {

	var list LinksList
	err := bambou.CurrentSession().FetchChildren(o, LinkIdentity, &list, info)
	return list, err
}

// CreateLink creates a new child Link under the Domain
func (o *Domain) CreateLink(child *Link) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// FirewallAcls retrieves the list of child FirewallAcls of the Domain
func (o *Domain) FirewallAcls(info *bambou.FetchingInfo) (FirewallAclsList, *bambou.Error) {

	var list FirewallAclsList
	err := bambou.CurrentSession().FetchChildren(o, FirewallAclIdentity, &list, info)
	return list, err
}

// FloatingIps retrieves the list of child FloatingIps of the Domain
func (o *Domain) FloatingIps(info *bambou.FetchingInfo) (FloatingIpsList, *bambou.Error) {

	var list FloatingIpsList
	err := bambou.CurrentSession().FetchChildren(o, FloatingIpIdentity, &list, info)
	return list, err
}

// CreateFloatingIp creates a new child FloatingIp under the Domain
func (o *Domain) CreateFloatingIp(child *FloatingIp) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Domain
func (o *Domain) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Domain
func (o *Domain) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMs retrieves the list of child VMs of the Domain
func (o *Domain) VMs(info *bambou.FetchingInfo) (VMsList, *bambou.Error) {

	var list VMsList
	err := bambou.CurrentSession().FetchChildren(o, VMIdentity, &list, info)
	return list, err
}

// VMInterfaces retrieves the list of child VMInterfaces of the Domain
func (o *Domain) VMInterfaces(info *bambou.FetchingInfo) (VMInterfacesList, *bambou.Error) {

	var list VMInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, VMInterfaceIdentity, &list, info)
	return list, err
}

// IngressACLEntryTemplates retrieves the list of child IngressACLEntryTemplates of the Domain
func (o *Domain) IngressACLEntryTemplates(info *bambou.FetchingInfo) (IngressACLEntryTemplatesList, *bambou.Error) {

	var list IngressACLEntryTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressACLEntryTemplateIdentity, &list, info)
	return list, err
}

// IngressACLTemplates retrieves the list of child IngressACLTemplates of the Domain
func (o *Domain) IngressACLTemplates(info *bambou.FetchingInfo) (IngressACLTemplatesList, *bambou.Error) {

	var list IngressACLTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressACLTemplateIdentity, &list, info)
	return list, err
}

// CreateIngressACLTemplate creates a new child IngressACLTemplate under the Domain
func (o *Domain) CreateIngressACLTemplate(child *IngressACLTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IngressAdvFwdTemplates retrieves the list of child IngressAdvFwdTemplates of the Domain
func (o *Domain) IngressAdvFwdTemplates(info *bambou.FetchingInfo) (IngressAdvFwdTemplatesList, *bambou.Error) {

	var list IngressAdvFwdTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressAdvFwdTemplateIdentity, &list, info)
	return list, err
}

// CreateIngressAdvFwdTemplate creates a new child IngressAdvFwdTemplate under the Domain
func (o *Domain) CreateIngressAdvFwdTemplate(child *IngressAdvFwdTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IngressExternalServiceTemplates retrieves the list of child IngressExternalServiceTemplates of the Domain
func (o *Domain) IngressExternalServiceTemplates(info *bambou.FetchingInfo) (IngressExternalServiceTemplatesList, *bambou.Error) {

	var list IngressExternalServiceTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressExternalServiceTemplateIdentity, &list, info)
	return list, err
}

// CreateIngressExternalServiceTemplate creates a new child IngressExternalServiceTemplate under the Domain
func (o *Domain) CreateIngressExternalServiceTemplate(child *IngressExternalServiceTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CreateJob creates a new child Job under the Domain
func (o *Domain) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PolicyGroups retrieves the list of child PolicyGroups of the Domain
func (o *Domain) PolicyGroups(info *bambou.FetchingInfo) (PolicyGroupsList, *bambou.Error) {

	var list PolicyGroupsList
	err := bambou.CurrentSession().FetchChildren(o, PolicyGroupIdentity, &list, info)
	return list, err
}

// CreatePolicyGroup creates a new child PolicyGroup under the Domain
func (o *Domain) CreatePolicyGroup(child *PolicyGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Domains retrieves the list of child Domains of the Domain
func (o *Domain) Domains(info *bambou.FetchingInfo) (DomainsList, *bambou.Error) {

	var list DomainsList
	err := bambou.CurrentSession().FetchChildren(o, DomainIdentity, &list, info)
	return list, err
}

// AssignDomains assigns the list of Domains to the Domain
func (o *Domain) AssignDomains(children DomainsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, DomainIdentity)
}

// DomainTemplates retrieves the list of child DomainTemplates of the Domain
func (o *Domain) DomainTemplates(info *bambou.FetchingInfo) (DomainTemplatesList, *bambou.Error) {

	var list DomainTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, DomainTemplateIdentity, &list, info)
	return list, err
}

// Zones retrieves the list of child Zones of the Domain
func (o *Domain) Zones(info *bambou.FetchingInfo) (ZonesList, *bambou.Error) {

	var list ZonesList
	err := bambou.CurrentSession().FetchChildren(o, ZoneIdentity, &list, info)
	return list, err
}

// CreateZone creates a new child Zone under the Domain
func (o *Domain) CreateZone(child *Zone) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Containers retrieves the list of child Containers of the Domain
func (o *Domain) Containers(info *bambou.FetchingInfo) (ContainersList, *bambou.Error) {

	var list ContainersList
	err := bambou.CurrentSession().FetchChildren(o, ContainerIdentity, &list, info)
	return list, err
}

// ContainerInterfaces retrieves the list of child ContainerInterfaces of the Domain
func (o *Domain) ContainerInterfaces(info *bambou.FetchingInfo) (ContainerInterfacesList, *bambou.Error) {

	var list ContainerInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, ContainerInterfaceIdentity, &list, info)
	return list, err
}

// QOSs retrieves the list of child QOSs of the Domain
func (o *Domain) QOSs(info *bambou.FetchingInfo) (QOSsList, *bambou.Error) {

	var list QOSsList
	err := bambou.CurrentSession().FetchChildren(o, QOSIdentity, &list, info)
	return list, err
}

// CreateQOS creates a new child QOS under the Domain
func (o *Domain) CreateQOS(child *QOS) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// HostInterfaces retrieves the list of child HostInterfaces of the Domain
func (o *Domain) HostInterfaces(info *bambou.FetchingInfo) (HostInterfacesList, *bambou.Error) {

	var list HostInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, HostInterfaceIdentity, &list, info)
	return list, err
}

// RoutingPolicies retrieves the list of child RoutingPolicies of the Domain
func (o *Domain) RoutingPolicies(info *bambou.FetchingInfo) (RoutingPoliciesList, *bambou.Error) {

	var list RoutingPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, RoutingPolicyIdentity, &list, info)
	return list, err
}

// CreateRoutingPolicy creates a new child RoutingPolicy under the Domain
func (o *Domain) CreateRoutingPolicy(child *RoutingPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// UplinkRDs retrieves the list of child UplinkRDs of the Domain
func (o *Domain) UplinkRDs(info *bambou.FetchingInfo) (UplinkRDsList, *bambou.Error) {

	var list UplinkRDsList
	err := bambou.CurrentSession().FetchChildren(o, UplinkRDIdentity, &list, info)
	return list, err
}

// VPNConnections retrieves the list of child VPNConnections of the Domain
func (o *Domain) VPNConnections(info *bambou.FetchingInfo) (VPNConnectionsList, *bambou.Error) {

	var list VPNConnectionsList
	err := bambou.CurrentSession().FetchChildren(o, VPNConnectionIdentity, &list, info)
	return list, err
}

// CreateVPNConnection creates a new child VPNConnection under the Domain
func (o *Domain) CreateVPNConnection(child *VPNConnection) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VPorts retrieves the list of child VPorts of the Domain
func (o *Domain) VPorts(info *bambou.FetchingInfo) (VPortsList, *bambou.Error) {

	var list VPortsList
	err := bambou.CurrentSession().FetchChildren(o, VPortIdentity, &list, info)
	return list, err
}

// Applicationperformancemanagementbindings retrieves the list of child Applicationperformancemanagementbindings of the Domain
func (o *Domain) Applicationperformancemanagementbindings(info *bambou.FetchingInfo) (ApplicationperformancemanagementbindingsList, *bambou.Error) {

	var list ApplicationperformancemanagementbindingsList
	err := bambou.CurrentSession().FetchChildren(o, ApplicationperformancemanagementbindingIdentity, &list, info)
	return list, err
}

// CreateApplicationperformancemanagementbinding creates a new child Applicationperformancemanagementbinding under the Domain
func (o *Domain) CreateApplicationperformancemanagementbinding(child *Applicationperformancemanagementbinding) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// BridgeInterfaces retrieves the list of child BridgeInterfaces of the Domain
func (o *Domain) BridgeInterfaces(info *bambou.FetchingInfo) (BridgeInterfacesList, *bambou.Error) {

	var list BridgeInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, BridgeInterfaceIdentity, &list, info)
	return list, err
}

// Groups retrieves the list of child Groups of the Domain
func (o *Domain) Groups(info *bambou.FetchingInfo) (GroupsList, *bambou.Error) {

	var list GroupsList
	err := bambou.CurrentSession().FetchChildren(o, GroupIdentity, &list, info)
	return list, err
}

// StaticRoutes retrieves the list of child StaticRoutes of the Domain
func (o *Domain) StaticRoutes(info *bambou.FetchingInfo) (StaticRoutesList, *bambou.Error) {

	var list StaticRoutesList
	err := bambou.CurrentSession().FetchChildren(o, StaticRouteIdentity, &list, info)
	return list, err
}

// CreateStaticRoute creates a new child StaticRoute under the Domain
func (o *Domain) CreateStaticRoute(child *StaticRoute) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Statistics retrieves the list of child Statistics of the Domain
func (o *Domain) Statistics(info *bambou.FetchingInfo) (StatisticsList, *bambou.Error) {

	var list StatisticsList
	err := bambou.CurrentSession().FetchChildren(o, StatisticsIdentity, &list, info)
	return list, err
}

// StatisticsPolicies retrieves the list of child StatisticsPolicies of the Domain
func (o *Domain) StatisticsPolicies(info *bambou.FetchingInfo) (StatisticsPoliciesList, *bambou.Error) {

	var list StatisticsPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, StatisticsPolicyIdentity, &list, info)
	return list, err
}

// CreateStatisticsPolicy creates a new child StatisticsPolicy under the Domain
func (o *Domain) CreateStatisticsPolicy(child *StatisticsPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Subnets retrieves the list of child Subnets of the Domain
func (o *Domain) Subnets(info *bambou.FetchingInfo) (SubnetsList, *bambou.Error) {

	var list SubnetsList
	err := bambou.CurrentSession().FetchChildren(o, SubnetIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the Domain
func (o *Domain) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}

// ExternalAppServices retrieves the list of child ExternalAppServices of the Domain
func (o *Domain) ExternalAppServices(info *bambou.FetchingInfo) (ExternalAppServicesList, *bambou.Error) {

	var list ExternalAppServicesList
	err := bambou.CurrentSession().FetchChildren(o, ExternalAppServiceIdentity, &list, info)
	return list, err
}

// CreateExternalAppService creates a new child ExternalAppService under the Domain
func (o *Domain) CreateExternalAppService(child *ExternalAppService) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
