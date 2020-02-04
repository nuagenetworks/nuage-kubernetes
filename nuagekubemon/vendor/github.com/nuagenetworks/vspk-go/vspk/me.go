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

// MeIdentity represents the Identity of the object
var MeIdentity = bambou.Identity{
	Name:     "me",
	Category: "me",
}

// Me represents the model of a me
type Me struct {
	ID                    string        `json:"ID,omitempty"`
	ParentID              string        `json:"parentID,omitempty"`
	ParentType            string        `json:"parentType,omitempty"`
	Owner                 string        `json:"owner,omitempty"`
	AARFlowStatsInterval  int           `json:"AARFlowStatsInterval,omitempty"`
	AARProbeStatsInterval int           `json:"AARProbeStatsInterval,omitempty"`
	VSSStatsInterval      int           `json:"VSSStatsInterval,omitempty"`
	Password              string        `json:"password,omitempty"`
	LastName              string        `json:"lastName,omitempty"`
	LastUpdatedBy         string        `json:"lastUpdatedBy,omitempty"`
	FirstName             string        `json:"firstName,omitempty"`
	Disabled              bool          `json:"disabled"`
	ElasticSearchAddress  string        `json:"elasticSearchAddress,omitempty"`
	FlowCollectionEnabled bool          `json:"flowCollectionEnabled"`
	Email                 string        `json:"email,omitempty"`
	EmbeddedMetadata      []interface{} `json:"embeddedMetadata,omitempty"`
	EnterpriseID          string        `json:"enterpriseID,omitempty"`
	EnterpriseName        string        `json:"enterpriseName,omitempty"`
	EntityScope           string        `json:"entityScope,omitempty"`
	MobileNumber          string        `json:"mobileNumber,omitempty"`
	Role                  string        `json:"role,omitempty"`
	UserName              string        `json:"userName,omitempty"`
	StatisticsEnabled     bool          `json:"statisticsEnabled"`
	AvatarData            string        `json:"avatarData,omitempty"`
	AvatarType            string        `json:"avatarType,omitempty"`
	ExternalID            string        `json:"externalID,omitempty"`

	Token        string `json:"APIKey,omitempty"`
	Organization string `json:"enterprise,omitempty"`
}

// NewMe returns a new *Me
func NewMe() *Me {

	return &Me{
		AARFlowStatsInterval:  30,
		AARProbeStatsInterval: 30,
		VSSStatsInterval:      30,
	}
}

// Identity returns the Identity of the object.
func (o *Me) Identity() bambou.Identity {

	return MeIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Me) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Me) SetIdentifier(ID string) {

	o.ID = ID
}

// APIKey returns a the API Key
func (o *Me) APIKey() string {

	return o.Token
}

// SetAPIKey sets a the API Key
func (o *Me) SetAPIKey(key string) {

	o.Token = key
}

// Fetch retrieves the Me from the server
func (o *Me) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Me into the server
func (o *Me) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Me from the server
func (o *Me) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// L2Domains retrieves the list of child L2Domains of the Me
func (o *Me) L2Domains(info *bambou.FetchingInfo) (L2DomainsList, *bambou.Error) {

	var list L2DomainsList
	err := bambou.CurrentSession().FetchChildren(o, L2DomainIdentity, &list, info)
	return list, err
}

// L4Services retrieves the list of child L4Services of the Me
func (o *Me) L4Services(info *bambou.FetchingInfo) (L4ServicesList, *bambou.Error) {

	var list L4ServicesList
	err := bambou.CurrentSession().FetchChildren(o, L4ServiceIdentity, &list, info)
	return list, err
}

// CreateL4Service creates a new child L4Service under the Me
func (o *Me) CreateL4Service(child *L4Service) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// L7applicationsignatures retrieves the list of child L7applicationsignatures of the Me
func (o *Me) L7applicationsignatures(info *bambou.FetchingInfo) (L7applicationsignaturesList, *bambou.Error) {

	var list L7applicationsignaturesList
	err := bambou.CurrentSession().FetchChildren(o, L7applicationsignatureIdentity, &list, info)
	return list, err
}

// SaaSApplicationTypes retrieves the list of child SaaSApplicationTypes of the Me
func (o *Me) SaaSApplicationTypes(info *bambou.FetchingInfo) (SaaSApplicationTypesList, *bambou.Error) {

	var list SaaSApplicationTypesList
	err := bambou.CurrentSession().FetchChildren(o, SaaSApplicationTypeIdentity, &list, info)
	return list, err
}

// VCenterEAMConfigs retrieves the list of child VCenterEAMConfigs of the Me
func (o *Me) VCenterEAMConfigs(info *bambou.FetchingInfo) (VCenterEAMConfigsList, *bambou.Error) {

	var list VCenterEAMConfigsList
	err := bambou.CurrentSession().FetchChildren(o, VCenterEAMConfigIdentity, &list, info)
	return list, err
}

// RateLimiters retrieves the list of child RateLimiters of the Me
func (o *Me) RateLimiters(info *bambou.FetchingInfo) (RateLimitersList, *bambou.Error) {

	var list RateLimitersList
	err := bambou.CurrentSession().FetchChildren(o, RateLimiterIdentity, &list, info)
	return list, err
}

// CreateRateLimiter creates a new child RateLimiter under the Me
func (o *Me) CreateRateLimiter(child *RateLimiter) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Gateways retrieves the list of child Gateways of the Me
func (o *Me) Gateways(info *bambou.FetchingInfo) (GatewaysList, *bambou.Error) {

	var list GatewaysList
	err := bambou.CurrentSession().FetchChildren(o, GatewayIdentity, &list, info)
	return list, err
}

// CreateGateway creates a new child Gateway under the Me
func (o *Me) CreateGateway(child *Gateway) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GatewayTemplates retrieves the list of child GatewayTemplates of the Me
func (o *Me) GatewayTemplates(info *bambou.FetchingInfo) (GatewayTemplatesList, *bambou.Error) {

	var list GatewayTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, GatewayTemplateIdentity, &list, info)
	return list, err
}

// CreateGatewayTemplate creates a new child GatewayTemplate under the Me
func (o *Me) CreateGatewayTemplate(child *GatewayTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PATMappers retrieves the list of child PATMappers of the Me
func (o *Me) PATMappers(info *bambou.FetchingInfo) (PATMappersList, *bambou.Error) {

	var list PATMappersList
	err := bambou.CurrentSession().FetchChildren(o, PATMapperIdentity, &list, info)
	return list, err
}

// CreatePATMapper creates a new child PATMapper under the Me
func (o *Me) CreatePATMapper(child *PATMapper) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PATNATPools retrieves the list of child PATNATPools of the Me
func (o *Me) PATNATPools(info *bambou.FetchingInfo) (PATNATPoolsList, *bambou.Error) {

	var list PATNATPoolsList
	err := bambou.CurrentSession().FetchChildren(o, PATNATPoolIdentity, &list, info)
	return list, err
}

// CreatePATNATPool creates a new child PATNATPool under the Me
func (o *Me) CreatePATNATPool(child *PATNATPool) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// TCAs retrieves the list of child TCAs of the Me
func (o *Me) TCAs(info *bambou.FetchingInfo) (TCAsList, *bambou.Error) {

	var list TCAsList
	err := bambou.CurrentSession().FetchChildren(o, TCAIdentity, &list, info)
	return list, err
}

// VCenters retrieves the list of child VCenters of the Me
func (o *Me) VCenters(info *bambou.FetchingInfo) (VCentersList, *bambou.Error) {

	var list VCentersList
	err := bambou.CurrentSession().FetchChildren(o, VCenterIdentity, &list, info)
	return list, err
}

// CreateVCenter creates a new child VCenter under the Me
func (o *Me) CreateVCenter(child *VCenter) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VCenterHypervisors retrieves the list of child VCenterHypervisors of the Me
func (o *Me) VCenterHypervisors(info *bambou.FetchingInfo) (VCenterHypervisorsList, *bambou.Error) {

	var list VCenterHypervisorsList
	err := bambou.CurrentSession().FetchChildren(o, VCenterHypervisorIdentity, &list, info)
	return list, err
}

// RedirectionTargets retrieves the list of child RedirectionTargets of the Me
func (o *Me) RedirectionTargets(info *bambou.FetchingInfo) (RedirectionTargetsList, *bambou.Error) {

	var list RedirectionTargetsList
	err := bambou.CurrentSession().FetchChildren(o, RedirectionTargetIdentity, &list, info)
	return list, err
}

// RedundancyGroups retrieves the list of child RedundancyGroups of the Me
func (o *Me) RedundancyGroups(info *bambou.FetchingInfo) (RedundancyGroupsList, *bambou.Error) {

	var list RedundancyGroupsList
	err := bambou.CurrentSession().FetchChildren(o, RedundancyGroupIdentity, &list, info)
	return list, err
}

// CreateRedundancyGroup creates a new child RedundancyGroup under the Me
func (o *Me) CreateRedundancyGroup(child *RedundancyGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// RemoteVrsInfos retrieves the list of child RemoteVrsInfos of the Me
func (o *Me) RemoteVrsInfos(info *bambou.FetchingInfo) (RemoteVrsInfosList, *bambou.Error) {

	var list RemoteVrsInfosList
	err := bambou.CurrentSession().FetchChildren(o, RemoteVrsInfoIdentity, &list, info)
	return list, err
}

// CreateRemoteVrsInfo creates a new child RemoteVrsInfo under the Me
func (o *Me) CreateRemoteVrsInfo(child *RemoteVrsInfo) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PerformanceMonitors retrieves the list of child PerformanceMonitors of the Me
func (o *Me) PerformanceMonitors(info *bambou.FetchingInfo) (PerformanceMonitorsList, *bambou.Error) {

	var list PerformanceMonitorsList
	err := bambou.CurrentSession().FetchChildren(o, PerformanceMonitorIdentity, &list, info)
	return list, err
}

// CreatePerformanceMonitor creates a new child PerformanceMonitor under the Me
func (o *Me) CreatePerformanceMonitor(child *PerformanceMonitor) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CreateCertificate creates a new child Certificate under the Me
func (o *Me) CreateCertificate(child *Certificate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// TestDefinitions retrieves the list of child TestDefinitions of the Me
func (o *Me) TestDefinitions(info *bambou.FetchingInfo) (TestDefinitionsList, *bambou.Error) {

	var list TestDefinitionsList
	err := bambou.CurrentSession().FetchChildren(o, TestDefinitionIdentity, &list, info)
	return list, err
}

// CreateTestDefinition creates a new child TestDefinition under the Me
func (o *Me) CreateTestDefinition(child *TestDefinition) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the Me
func (o *Me) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// NetconfProfiles retrieves the list of child NetconfProfiles of the Me
func (o *Me) NetconfProfiles(info *bambou.FetchingInfo) (NetconfProfilesList, *bambou.Error) {

	var list NetconfProfilesList
	err := bambou.CurrentSession().FetchChildren(o, NetconfProfileIdentity, &list, info)
	return list, err
}

// CreateNetconfProfile creates a new child NetconfProfile under the Me
func (o *Me) CreateNetconfProfile(child *NetconfProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NetworkLayouts retrieves the list of child NetworkLayouts of the Me
func (o *Me) NetworkLayouts(info *bambou.FetchingInfo) (NetworkLayoutsList, *bambou.Error) {

	var list NetworkLayoutsList
	err := bambou.CurrentSession().FetchChildren(o, NetworkLayoutIdentity, &list, info)
	return list, err
}

// NetworkPerformanceMeasurements retrieves the list of child NetworkPerformanceMeasurements of the Me
func (o *Me) NetworkPerformanceMeasurements(info *bambou.FetchingInfo) (NetworkPerformanceMeasurementsList, *bambou.Error) {

	var list NetworkPerformanceMeasurementsList
	err := bambou.CurrentSession().FetchChildren(o, NetworkPerformanceMeasurementIdentity, &list, info)
	return list, err
}

// KeyServerMembers retrieves the list of child KeyServerMembers of the Me
func (o *Me) KeyServerMembers(info *bambou.FetchingInfo) (KeyServerMembersList, *bambou.Error) {

	var list KeyServerMembersList
	err := bambou.CurrentSession().FetchChildren(o, KeyServerMemberIdentity, &list, info)
	return list, err
}

// CreateKeyServerMember creates a new child KeyServerMember under the Me
func (o *Me) CreateKeyServerMember(child *KeyServerMember) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// ZFBAutoAssignments retrieves the list of child ZFBAutoAssignments of the Me
func (o *Me) ZFBAutoAssignments(info *bambou.FetchingInfo) (ZFBAutoAssignmentsList, *bambou.Error) {

	var list ZFBAutoAssignmentsList
	err := bambou.CurrentSession().FetchChildren(o, ZFBAutoAssignmentIdentity, &list, info)
	return list, err
}

// CreateZFBAutoAssignment creates a new child ZFBAutoAssignment under the Me
func (o *Me) CreateZFBAutoAssignment(child *ZFBAutoAssignment) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// ZFBRequests retrieves the list of child ZFBRequests of the Me
func (o *Me) ZFBRequests(info *bambou.FetchingInfo) (ZFBRequestsList, *bambou.Error) {

	var list ZFBRequestsList
	err := bambou.CurrentSession().FetchChildren(o, ZFBRequestIdentity, &list, info)
	return list, err
}

// CreateZFBRequest creates a new child ZFBRequest under the Me
func (o *Me) CreateZFBRequest(child *ZFBRequest) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// BGPNeighbors retrieves the list of child BGPNeighbors of the Me
func (o *Me) BGPNeighbors(info *bambou.FetchingInfo) (BGPNeighborsList, *bambou.Error) {

	var list BGPNeighborsList
	err := bambou.CurrentSession().FetchChildren(o, BGPNeighborIdentity, &list, info)
	return list, err
}

// BGPProfiles retrieves the list of child BGPProfiles of the Me
func (o *Me) BGPProfiles(info *bambou.FetchingInfo) (BGPProfilesList, *bambou.Error) {

	var list BGPProfilesList
	err := bambou.CurrentSession().FetchChildren(o, BGPProfileIdentity, &list, info)
	return list, err
}

// EgressACLEntryTemplates retrieves the list of child EgressACLEntryTemplates of the Me
func (o *Me) EgressACLEntryTemplates(info *bambou.FetchingInfo) (EgressACLEntryTemplatesList, *bambou.Error) {

	var list EgressACLEntryTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, EgressACLEntryTemplateIdentity, &list, info)
	return list, err
}

// EgressACLTemplates retrieves the list of child EgressACLTemplates of the Me
func (o *Me) EgressACLTemplates(info *bambou.FetchingInfo) (EgressACLTemplatesList, *bambou.Error) {

	var list EgressACLTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, EgressACLTemplateIdentity, &list, info)
	return list, err
}

// EgressAdvFwdEntryTemplates retrieves the list of child EgressAdvFwdEntryTemplates of the Me
func (o *Me) EgressAdvFwdEntryTemplates(info *bambou.FetchingInfo) (EgressAdvFwdEntryTemplatesList, *bambou.Error) {

	var list EgressAdvFwdEntryTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, EgressAdvFwdEntryTemplateIdentity, &list, info)
	return list, err
}

// DomainFIPAclTemplates retrieves the list of child DomainFIPAclTemplates of the Me
func (o *Me) DomainFIPAclTemplates(info *bambou.FetchingInfo) (DomainFIPAclTemplatesList, *bambou.Error) {

	var list DomainFIPAclTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, DomainFIPAclTemplateIdentity, &list, info)
	return list, err
}

// CreateDomainFIPAclTemplate creates a new child DomainFIPAclTemplate under the Me
func (o *Me) CreateDomainFIPAclTemplate(child *DomainFIPAclTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EgressQOSPolicies retrieves the list of child EgressQOSPolicies of the Me
func (o *Me) EgressQOSPolicies(info *bambou.FetchingInfo) (EgressQOSPoliciesList, *bambou.Error) {

	var list EgressQOSPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, EgressQOSPolicyIdentity, &list, info)
	return list, err
}

// CreateEgressQOSPolicy creates a new child EgressQOSPolicy under the Me
func (o *Me) CreateEgressQOSPolicy(child *EgressQOSPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// SharedNetworkResources retrieves the list of child SharedNetworkResources of the Me
func (o *Me) SharedNetworkResources(info *bambou.FetchingInfo) (SharedNetworkResourcesList, *bambou.Error) {

	var list SharedNetworkResourcesList
	err := bambou.CurrentSession().FetchChildren(o, SharedNetworkResourceIdentity, &list, info)
	return list, err
}

// CreateSharedNetworkResource creates a new child SharedNetworkResource under the Me
func (o *Me) CreateSharedNetworkResource(child *SharedNetworkResource) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Licenses retrieves the list of child Licenses of the Me
func (o *Me) Licenses(info *bambou.FetchingInfo) (LicensesList, *bambou.Error) {

	var list LicensesList
	err := bambou.CurrentSession().FetchChildren(o, LicenseIdentity, &list, info)
	return list, err
}

// CreateLicense creates a new child License under the Me
func (o *Me) CreateLicense(child *License) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// LicenseStatus retrieves the list of child LicenseStatus of the Me
func (o *Me) LicenseStatus(info *bambou.FetchingInfo) (LicenseStatusList, *bambou.Error) {

	var list LicenseStatusList
	err := bambou.CurrentSession().FetchChildren(o, LicenseStatusIdentity, &list, info)
	return list, err
}

// MirrorDestinations retrieves the list of child MirrorDestinations of the Me
func (o *Me) MirrorDestinations(info *bambou.FetchingInfo) (MirrorDestinationsList, *bambou.Error) {

	var list MirrorDestinationsList
	err := bambou.CurrentSession().FetchChildren(o, MirrorDestinationIdentity, &list, info)
	return list, err
}

// CreateMirrorDestination creates a new child MirrorDestination under the Me
func (o *Me) CreateMirrorDestination(child *MirrorDestination) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VirtualFirewallPolicies retrieves the list of child VirtualFirewallPolicies of the Me
func (o *Me) VirtualFirewallPolicies(info *bambou.FetchingInfo) (VirtualFirewallPoliciesList, *bambou.Error) {

	var list VirtualFirewallPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, VirtualFirewallPolicyIdentity, &list, info)
	return list, err
}

// VirtualFirewallRules retrieves the list of child VirtualFirewallRules of the Me
func (o *Me) VirtualFirewallRules(info *bambou.FetchingInfo) (VirtualFirewallRulesList, *bambou.Error) {

	var list VirtualFirewallRulesList
	err := bambou.CurrentSession().FetchChildren(o, VirtualFirewallRuleIdentity, &list, info)
	return list, err
}

// SiteInfos retrieves the list of child SiteInfos of the Me
func (o *Me) SiteInfos(info *bambou.FetchingInfo) (SiteInfosList, *bambou.Error) {

	var list SiteInfosList
	err := bambou.CurrentSession().FetchChildren(o, SiteInfoIdentity, &list, info)
	return list, err
}

// CreateSiteInfo creates a new child SiteInfo under the Me
func (o *Me) CreateSiteInfo(child *SiteInfo) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// AllGateways retrieves the list of child AllGateways of the Me
func (o *Me) AllGateways(info *bambou.FetchingInfo) (AllGatewaysList, *bambou.Error) {

	var list AllGatewaysList
	err := bambou.CurrentSession().FetchChildren(o, AllGatewayIdentity, &list, info)
	return list, err
}

// AllRedundancyGroups retrieves the list of child AllRedundancyGroups of the Me
func (o *Me) AllRedundancyGroups(info *bambou.FetchingInfo) (AllRedundancyGroupsList, *bambou.Error) {

	var list AllRedundancyGroupsList
	err := bambou.CurrentSession().FetchChildren(o, AllRedundancyGroupIdentity, &list, info)
	return list, err
}

// CreateAllRedundancyGroup creates a new child AllRedundancyGroup under the Me
func (o *Me) CreateAllRedundancyGroup(child *AllRedundancyGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// FloatingIps retrieves the list of child FloatingIps of the Me
func (o *Me) FloatingIps(info *bambou.FetchingInfo) (FloatingIpsList, *bambou.Error) {

	var list FloatingIpsList
	err := bambou.CurrentSession().FetchChildren(o, FloatingIpIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Me
func (o *Me) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Me
func (o *Me) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMs retrieves the list of child VMs of the Me
func (o *Me) VMs(info *bambou.FetchingInfo) (VMsList, *bambou.Error) {

	var list VMsList
	err := bambou.CurrentSession().FetchChildren(o, VMIdentity, &list, info)
	return list, err
}

// CreateVM creates a new child VM under the Me
func (o *Me) CreateVM(child *VM) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMInterfaces retrieves the list of child VMInterfaces of the Me
func (o *Me) VMInterfaces(info *bambou.FetchingInfo) (VMInterfacesList, *bambou.Error) {

	var list VMInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, VMInterfaceIdentity, &list, info)
	return list, err
}

// CloudMgmtSystems retrieves the list of child CloudMgmtSystems of the Me
func (o *Me) CloudMgmtSystems(info *bambou.FetchingInfo) (CloudMgmtSystemsList, *bambou.Error) {

	var list CloudMgmtSystemsList
	err := bambou.CurrentSession().FetchChildren(o, CloudMgmtSystemIdentity, &list, info)
	return list, err
}

// CreateCloudMgmtSystem creates a new child CloudMgmtSystem under the Me
func (o *Me) CreateCloudMgmtSystem(child *CloudMgmtSystem) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Underlays retrieves the list of child Underlays of the Me
func (o *Me) Underlays(info *bambou.FetchingInfo) (UnderlaysList, *bambou.Error) {

	var list UnderlaysList
	err := bambou.CurrentSession().FetchChildren(o, UnderlayIdentity, &list, info)
	return list, err
}

// CreateUnderlay creates a new child Underlay under the Me
func (o *Me) CreateUnderlay(child *Underlay) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VNFCatalogs retrieves the list of child VNFCatalogs of the Me
func (o *Me) VNFCatalogs(info *bambou.FetchingInfo) (VNFCatalogsList, *bambou.Error) {

	var list VNFCatalogsList
	err := bambou.CurrentSession().FetchChildren(o, VNFCatalogIdentity, &list, info)
	return list, err
}

// VNFMetadatas retrieves the list of child VNFMetadatas of the Me
func (o *Me) VNFMetadatas(info *bambou.FetchingInfo) (VNFMetadatasList, *bambou.Error) {

	var list VNFMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, VNFMetadataIdentity, &list, info)
	return list, err
}

// CreateVNFMetadata creates a new child VNFMetadata under the Me
func (o *Me) CreateVNFMetadata(child *VNFMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// InfrastructureAccessProfiles retrieves the list of child InfrastructureAccessProfiles of the Me
func (o *Me) InfrastructureAccessProfiles(info *bambou.FetchingInfo) (InfrastructureAccessProfilesList, *bambou.Error) {

	var list InfrastructureAccessProfilesList
	err := bambou.CurrentSession().FetchChildren(o, InfrastructureAccessProfileIdentity, &list, info)
	return list, err
}

// CreateInfrastructureAccessProfile creates a new child InfrastructureAccessProfile under the Me
func (o *Me) CreateInfrastructureAccessProfile(child *InfrastructureAccessProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// InfrastructureEVDFProfiles retrieves the list of child InfrastructureEVDFProfiles of the Me
func (o *Me) InfrastructureEVDFProfiles(info *bambou.FetchingInfo) (InfrastructureEVDFProfilesList, *bambou.Error) {

	var list InfrastructureEVDFProfilesList
	err := bambou.CurrentSession().FetchChildren(o, InfrastructureEVDFProfileIdentity, &list, info)
	return list, err
}

// CreateInfrastructureEVDFProfile creates a new child InfrastructureEVDFProfile under the Me
func (o *Me) CreateInfrastructureEVDFProfile(child *InfrastructureEVDFProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// InfrastructureGatewayProfiles retrieves the list of child InfrastructureGatewayProfiles of the Me
func (o *Me) InfrastructureGatewayProfiles(info *bambou.FetchingInfo) (InfrastructureGatewayProfilesList, *bambou.Error) {

	var list InfrastructureGatewayProfilesList
	err := bambou.CurrentSession().FetchChildren(o, InfrastructureGatewayProfileIdentity, &list, info)
	return list, err
}

// CreateInfrastructureGatewayProfile creates a new child InfrastructureGatewayProfile under the Me
func (o *Me) CreateInfrastructureGatewayProfile(child *InfrastructureGatewayProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// InfrastructureVscProfiles retrieves the list of child InfrastructureVscProfiles of the Me
func (o *Me) InfrastructureVscProfiles(info *bambou.FetchingInfo) (InfrastructureVscProfilesList, *bambou.Error) {

	var list InfrastructureVscProfilesList
	err := bambou.CurrentSession().FetchChildren(o, InfrastructureVscProfileIdentity, &list, info)
	return list, err
}

// CreateInfrastructureVscProfile creates a new child InfrastructureVscProfile under the Me
func (o *Me) CreateInfrastructureVscProfile(child *InfrastructureVscProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VNFThresholdPolicies retrieves the list of child VNFThresholdPolicies of the Me
func (o *Me) VNFThresholdPolicies(info *bambou.FetchingInfo) (VNFThresholdPoliciesList, *bambou.Error) {

	var list VNFThresholdPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, VNFThresholdPolicyIdentity, &list, info)
	return list, err
}

// CreateVNFThresholdPolicy creates a new child VNFThresholdPolicy under the Me
func (o *Me) CreateVNFThresholdPolicy(child *VNFThresholdPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IngressACLEntryTemplates retrieves the list of child IngressACLEntryTemplates of the Me
func (o *Me) IngressACLEntryTemplates(info *bambou.FetchingInfo) (IngressACLEntryTemplatesList, *bambou.Error) {

	var list IngressACLEntryTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressACLEntryTemplateIdentity, &list, info)
	return list, err
}

// IngressACLTemplates retrieves the list of child IngressACLTemplates of the Me
func (o *Me) IngressACLTemplates(info *bambou.FetchingInfo) (IngressACLTemplatesList, *bambou.Error) {

	var list IngressACLTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressACLTemplateIdentity, &list, info)
	return list, err
}

// IngressAdvFwdEntryTemplates retrieves the list of child IngressAdvFwdEntryTemplates of the Me
func (o *Me) IngressAdvFwdEntryTemplates(info *bambou.FetchingInfo) (IngressAdvFwdEntryTemplatesList, *bambou.Error) {

	var list IngressAdvFwdEntryTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, IngressAdvFwdEntryTemplateIdentity, &list, info)
	return list, err
}

// IngressQOSPolicies retrieves the list of child IngressQOSPolicies of the Me
func (o *Me) IngressQOSPolicies(info *bambou.FetchingInfo) (IngressQOSPoliciesList, *bambou.Error) {

	var list IngressQOSPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, IngressQOSPolicyIdentity, &list, info)
	return list, err
}

// CreateIngressQOSPolicy creates a new child IngressQOSPolicy under the Me
func (o *Me) CreateIngressQOSPolicy(child *IngressQOSPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Enterprises retrieves the list of child Enterprises of the Me
func (o *Me) Enterprises(info *bambou.FetchingInfo) (EnterprisesList, *bambou.Error) {

	var list EnterprisesList
	err := bambou.CurrentSession().FetchChildren(o, EnterpriseIdentity, &list, info)
	return list, err
}

// CreateEnterprise creates a new child Enterprise under the Me
func (o *Me) CreateEnterprise(child *Enterprise) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EnterpriseProfiles retrieves the list of child EnterpriseProfiles of the Me
func (o *Me) EnterpriseProfiles(info *bambou.FetchingInfo) (EnterpriseProfilesList, *bambou.Error) {

	var list EnterpriseProfilesList
	err := bambou.CurrentSession().FetchChildren(o, EnterpriseProfileIdentity, &list, info)
	return list, err
}

// CreateEnterpriseProfile creates a new child EnterpriseProfile under the Me
func (o *Me) CreateEnterpriseProfile(child *EnterpriseProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Jobs retrieves the list of child Jobs of the Me
func (o *Me) Jobs(info *bambou.FetchingInfo) (JobsList, *bambou.Error) {

	var list JobsList
	err := bambou.CurrentSession().FetchChildren(o, JobIdentity, &list, info)
	return list, err
}

// CreateJob creates a new child Job under the Me
func (o *Me) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PolicyGroups retrieves the list of child PolicyGroups of the Me
func (o *Me) PolicyGroups(info *bambou.FetchingInfo) (PolicyGroupsList, *bambou.Error) {

	var list PolicyGroupsList
	err := bambou.CurrentSession().FetchChildren(o, PolicyGroupIdentity, &list, info)
	return list, err
}

// PolicyObjectGroups retrieves the list of child PolicyObjectGroups of the Me
func (o *Me) PolicyObjectGroups(info *bambou.FetchingInfo) (PolicyObjectGroupsList, *bambou.Error) {

	var list PolicyObjectGroupsList
	err := bambou.CurrentSession().FetchChildren(o, PolicyObjectGroupIdentity, &list, info)
	return list, err
}

// Domains retrieves the list of child Domains of the Me
func (o *Me) Domains(info *bambou.FetchingInfo) (DomainsList, *bambou.Error) {

	var list DomainsList
	err := bambou.CurrentSession().FetchChildren(o, DomainIdentity, &list, info)
	return list, err
}

// Commands retrieves the list of child Commands of the Me
func (o *Me) Commands(info *bambou.FetchingInfo) (CommandsList, *bambou.Error) {

	var list CommandsList
	err := bambou.CurrentSession().FetchChildren(o, CommandIdentity, &list, info)
	return list, err
}

// CreateCommand creates a new child Command under the Me
func (o *Me) CreateCommand(child *Command) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Zones retrieves the list of child Zones of the Me
func (o *Me) Zones(info *bambou.FetchingInfo) (ZonesList, *bambou.Error) {

	var list ZonesList
	err := bambou.CurrentSession().FetchChildren(o, ZoneIdentity, &list, info)
	return list, err
}

// Containers retrieves the list of child Containers of the Me
func (o *Me) Containers(info *bambou.FetchingInfo) (ContainersList, *bambou.Error) {

	var list ContainersList
	err := bambou.CurrentSession().FetchChildren(o, ContainerIdentity, &list, info)
	return list, err
}

// CreateContainer creates a new child Container under the Me
func (o *Me) CreateContainer(child *Container) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// ContainerInterfaces retrieves the list of child ContainerInterfaces of the Me
func (o *Me) ContainerInterfaces(info *bambou.FetchingInfo) (ContainerInterfacesList, *bambou.Error) {

	var list ContainerInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, ContainerInterfaceIdentity, &list, info)
	return list, err
}

// QosPolicers retrieves the list of child QosPolicers of the Me
func (o *Me) QosPolicers(info *bambou.FetchingInfo) (QosPolicersList, *bambou.Error) {

	var list QosPolicersList
	err := bambou.CurrentSession().FetchChildren(o, QosPolicerIdentity, &list, info)
	return list, err
}

// CreateQosPolicer creates a new child QosPolicer under the Me
func (o *Me) CreateQosPolicer(child *QosPolicer) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// COSRemarkingPolicyTables retrieves the list of child COSRemarkingPolicyTables of the Me
func (o *Me) COSRemarkingPolicyTables(info *bambou.FetchingInfo) (COSRemarkingPolicyTablesList, *bambou.Error) {

	var list COSRemarkingPolicyTablesList
	err := bambou.CurrentSession().FetchChildren(o, COSRemarkingPolicyTableIdentity, &list, info)
	return list, err
}

// CreateCOSRemarkingPolicyTable creates a new child COSRemarkingPolicyTable under the Me
func (o *Me) CreateCOSRemarkingPolicyTable(child *COSRemarkingPolicyTable) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// HostInterfaces retrieves the list of child HostInterfaces of the Me
func (o *Me) HostInterfaces(info *bambou.FetchingInfo) (HostInterfacesList, *bambou.Error) {

	var list HostInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, HostInterfaceIdentity, &list, info)
	return list, err
}

// RoutingPolicies retrieves the list of child RoutingPolicies of the Me
func (o *Me) RoutingPolicies(info *bambou.FetchingInfo) (RoutingPoliciesList, *bambou.Error) {

	var list RoutingPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, RoutingPolicyIdentity, &list, info)
	return list, err
}

// UplinkRDs retrieves the list of child UplinkRDs of the Me
func (o *Me) UplinkRDs(info *bambou.FetchingInfo) (UplinkRDsList, *bambou.Error) {

	var list UplinkRDsList
	err := bambou.CurrentSession().FetchChildren(o, UplinkRDIdentity, &list, info)
	return list, err
}

// Applications retrieves the list of child Applications of the Me
func (o *Me) Applications(info *bambou.FetchingInfo) (ApplicationsList, *bambou.Error) {

	var list ApplicationsList
	err := bambou.CurrentSession().FetchChildren(o, ApplicationIdentity, &list, info)
	return list, err
}

// Applicationperformancemanagements retrieves the list of child Applicationperformancemanagements of the Me
func (o *Me) Applicationperformancemanagements(info *bambou.FetchingInfo) (ApplicationperformancemanagementsList, *bambou.Error) {

	var list ApplicationperformancemanagementsList
	err := bambou.CurrentSession().FetchChildren(o, ApplicationperformancemanagementIdentity, &list, info)
	return list, err
}

// VRSs retrieves the list of child VRSs of the Me
func (o *Me) VRSs(info *bambou.FetchingInfo) (VRSsList, *bambou.Error) {

	var list VRSsList
	err := bambou.CurrentSession().FetchChildren(o, VRSIdentity, &list, info)
	return list, err
}

// VCenterVRSConfigs retrieves the list of child VCenterVRSConfigs of the Me
func (o *Me) VCenterVRSConfigs(info *bambou.FetchingInfo) (VCenterVRSConfigsList, *bambou.Error) {

	var list VCenterVRSConfigsList
	err := bambou.CurrentSession().FetchChildren(o, VCenterVRSConfigIdentity, &list, info)
	return list, err
}

// VrsInfos retrieves the list of child vrsInfos of the Me
func (o *Me) VrsInfos(info *bambou.FetchingInfo) (vrsInfosList, *bambou.Error) {

	var list vrsInfosList
	err := bambou.CurrentSession().FetchChildren(o, vrsInfoIdentity, &list, info)
	return list, err
}

// DSCPRemarkingPolicyTables retrieves the list of child DSCPRemarkingPolicyTables of the Me
func (o *Me) DSCPRemarkingPolicyTables(info *bambou.FetchingInfo) (DSCPRemarkingPolicyTablesList, *bambou.Error) {

	var list DSCPRemarkingPolicyTablesList
	err := bambou.CurrentSession().FetchChildren(o, DSCPRemarkingPolicyTableIdentity, &list, info)
	return list, err
}

// CreateDSCPRemarkingPolicyTable creates a new child DSCPRemarkingPolicyTable under the Me
func (o *Me) CreateDSCPRemarkingPolicyTable(child *DSCPRemarkingPolicyTable) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Users retrieves the list of child Users of the Me
func (o *Me) Users(info *bambou.FetchingInfo) (UsersList, *bambou.Error) {

	var list UsersList
	err := bambou.CurrentSession().FetchChildren(o, UserIdentity, &list, info)
	return list, err
}

// CreateUser creates a new child User under the Me
func (o *Me) CreateUser(child *User) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// UserContexts retrieves the list of child UserContexts of the Me
func (o *Me) UserContexts(info *bambou.FetchingInfo) (UserContextsList, *bambou.Error) {

	var list UserContextsList
	err := bambou.CurrentSession().FetchChildren(o, UserContextIdentity, &list, info)
	return list, err
}

// NSGateways retrieves the list of child NSGateways of the Me
func (o *Me) NSGateways(info *bambou.FetchingInfo) (NSGatewaysList, *bambou.Error) {

	var list NSGatewaysList
	err := bambou.CurrentSession().FetchChildren(o, NSGatewayIdentity, &list, info)
	return list, err
}

// CreateNSGateway creates a new child NSGateway under the Me
func (o *Me) CreateNSGateway(child *NSGateway) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSGatewayTemplates retrieves the list of child NSGatewayTemplates of the Me
func (o *Me) NSGatewayTemplates(info *bambou.FetchingInfo) (NSGatewayTemplatesList, *bambou.Error) {

	var list NSGatewayTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, NSGatewayTemplateIdentity, &list, info)
	return list, err
}

// CreateNSGatewayTemplate creates a new child NSGatewayTemplate under the Me
func (o *Me) CreateNSGatewayTemplate(child *NSGatewayTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSGGroups retrieves the list of child NSGGroups of the Me
func (o *Me) NSGGroups(info *bambou.FetchingInfo) (NSGGroupsList, *bambou.Error) {

	var list NSGGroupsList
	err := bambou.CurrentSession().FetchChildren(o, NSGGroupIdentity, &list, info)
	return list, err
}

// CreateNSGGroup creates a new child NSGGroup under the Me
func (o *Me) CreateNSGGroup(child *NSGGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSGInfos retrieves the list of child NSGInfos of the Me
func (o *Me) NSGInfos(info *bambou.FetchingInfo) (NSGInfosList, *bambou.Error) {

	var list NSGInfosList
	err := bambou.CurrentSession().FetchChildren(o, NSGInfoIdentity, &list, info)
	return list, err
}

// NSGPatchProfiles retrieves the list of child NSGPatchProfiles of the Me
func (o *Me) NSGPatchProfiles(info *bambou.FetchingInfo) (NSGPatchProfilesList, *bambou.Error) {

	var list NSGPatchProfilesList
	err := bambou.CurrentSession().FetchChildren(o, NSGPatchProfileIdentity, &list, info)
	return list, err
}

// CreateNSGPatchProfile creates a new child NSGPatchProfile under the Me
func (o *Me) CreateNSGPatchProfile(child *NSGPatchProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSRedundantGatewayGroups retrieves the list of child NSRedundantGatewayGroups of the Me
func (o *Me) NSRedundantGatewayGroups(info *bambou.FetchingInfo) (NSRedundantGatewayGroupsList, *bambou.Error) {

	var list NSRedundantGatewayGroupsList
	err := bambou.CurrentSession().FetchChildren(o, NSRedundantGatewayGroupIdentity, &list, info)
	return list, err
}

// CreateNSRedundantGatewayGroup creates a new child NSRedundantGatewayGroup under the Me
func (o *Me) CreateNSRedundantGatewayGroup(child *NSRedundantGatewayGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSGUpgradeProfiles retrieves the list of child NSGUpgradeProfiles of the Me
func (o *Me) NSGUpgradeProfiles(info *bambou.FetchingInfo) (NSGUpgradeProfilesList, *bambou.Error) {

	var list NSGUpgradeProfilesList
	err := bambou.CurrentSession().FetchChildren(o, NSGUpgradeProfileIdentity, &list, info)
	return list, err
}

// CreateNSGUpgradeProfile creates a new child NSGUpgradeProfile under the Me
func (o *Me) CreateNSGUpgradeProfile(child *NSGUpgradeProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VSPs retrieves the list of child VSPs of the Me
func (o *Me) VSPs(info *bambou.FetchingInfo) (VSPsList, *bambou.Error) {

	var list VSPsList
	err := bambou.CurrentSession().FetchChildren(o, VSPIdentity, &list, info)
	return list, err
}

// StaticRoutes retrieves the list of child StaticRoutes of the Me
func (o *Me) StaticRoutes(info *bambou.FetchingInfo) (StaticRoutesList, *bambou.Error) {

	var list StaticRoutesList
	err := bambou.CurrentSession().FetchChildren(o, StaticRouteIdentity, &list, info)
	return list, err
}

// StatsCollectorInfos retrieves the list of child StatsCollectorInfos of the Me
func (o *Me) StatsCollectorInfos(info *bambou.FetchingInfo) (StatsCollectorInfosList, *bambou.Error) {

	var list StatsCollectorInfosList
	err := bambou.CurrentSession().FetchChildren(o, StatsCollectorInfoIdentity, &list, info)
	return list, err
}

// Subnets retrieves the list of child Subnets of the Me
func (o *Me) Subnets(info *bambou.FetchingInfo) (SubnetsList, *bambou.Error) {

	var list SubnetsList
	err := bambou.CurrentSession().FetchChildren(o, SubnetIdentity, &list, info)
	return list, err
}

// DUCGroups retrieves the list of child DUCGroups of the Me
func (o *Me) DUCGroups(info *bambou.FetchingInfo) (DUCGroupsList, *bambou.Error) {

	var list DUCGroupsList
	err := bambou.CurrentSession().FetchChildren(o, DUCGroupIdentity, &list, info)
	return list, err
}

// CreateDUCGroup creates a new child DUCGroup under the Me
func (o *Me) CreateDUCGroup(child *DUCGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// MultiCastChannelMaps retrieves the list of child MultiCastChannelMaps of the Me
func (o *Me) MultiCastChannelMaps(info *bambou.FetchingInfo) (MultiCastChannelMapsList, *bambou.Error) {

	var list MultiCastChannelMapsList
	err := bambou.CurrentSession().FetchChildren(o, MultiCastChannelMapIdentity, &list, info)
	return list, err
}

// CreateMultiCastChannelMap creates a new child MultiCastChannelMap under the Me
func (o *Me) CreateMultiCastChannelMap(child *MultiCastChannelMap) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// AutoDiscoveredGateways retrieves the list of child AutoDiscoveredGateways of the Me
func (o *Me) AutoDiscoveredGateways(info *bambou.FetchingInfo) (AutoDiscoveredGatewaysList, *bambou.Error) {

	var list AutoDiscoveredGatewaysList
	err := bambou.CurrentSession().FetchChildren(o, AutoDiscoveredGatewayIdentity, &list, info)
	return list, err
}

// OverlayMirrorDestinations retrieves the list of child OverlayMirrorDestinations of the Me
func (o *Me) OverlayMirrorDestinations(info *bambou.FetchingInfo) (OverlayMirrorDestinationsList, *bambou.Error) {

	var list OverlayMirrorDestinationsList
	err := bambou.CurrentSession().FetchChildren(o, OverlayMirrorDestinationIdentity, &list, info)
	return list, err
}

// SystemConfigs retrieves the list of child SystemConfigs of the Me
func (o *Me) SystemConfigs(info *bambou.FetchingInfo) (SystemConfigsList, *bambou.Error) {

	var list SystemConfigsList
	err := bambou.CurrentSession().FetchChildren(o, SystemConfigIdentity, &list, info)
	return list, err
}
