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

// EnterpriseIdentity represents the Identity of the object
var EnterpriseIdentity = bambou.Identity{
	Name:     "enterprise",
	Category: "enterprises",
}

// EnterprisesList represents a list of Enterprises
type EnterprisesList []*Enterprise

// EnterprisesAncestor is the interface that an ancestor of a Enterprise must implement.
// An Ancestor is defined as an entity that has Enterprise as a descendant.
// An Ancestor can get a list of its child Enterprises, but not necessarily create one.
type EnterprisesAncestor interface {
	Enterprises(*bambou.FetchingInfo) (EnterprisesList, *bambou.Error)
}

// EnterprisesParent is the interface that a parent of a Enterprise must implement.
// A Parent is defined as an entity that has Enterprise as a child.
// A Parent is an Ancestor which can create a Enterprise.
type EnterprisesParent interface {
	EnterprisesAncestor
	CreateEnterprise(*Enterprise) *bambou.Error
}

// Enterprise represents the model of a enterprise
type Enterprise struct {
	ID                                     string        `json:"ID,omitempty"`
	ParentID                               string        `json:"parentID,omitempty"`
	ParentType                             string        `json:"parentType,omitempty"`
	Owner                                  string        `json:"owner,omitempty"`
	LDAPAuthorizationEnabled               bool          `json:"LDAPAuthorizationEnabled"`
	LDAPEnabled                            bool          `json:"LDAPEnabled"`
	BGPEnabled                             bool          `json:"BGPEnabled"`
	DHCPLeaseInterval                      int           `json:"DHCPLeaseInterval,omitempty"`
	VNFManagementEnabled                   bool          `json:"VNFManagementEnabled"`
	Name                                   string        `json:"name,omitempty"`
	LastUpdatedBy                          string        `json:"lastUpdatedBy,omitempty"`
	WebFilterEnabled                       bool          `json:"webFilterEnabled"`
	ReceiveMultiCastListID                 string        `json:"receiveMultiCastListID,omitempty"`
	SendMultiCastListID                    string        `json:"sendMultiCastListID,omitempty"`
	Description                            string        `json:"description,omitempty"`
	SharedEnterprise                       bool          `json:"sharedEnterprise"`
	DictionaryVersion                      int           `json:"dictionaryVersion,omitempty"`
	VirtualFirewallRulesEnabled            bool          `json:"virtualFirewallRulesEnabled"`
	AllowAdvancedQOSConfiguration          bool          `json:"allowAdvancedQOSConfiguration"`
	AllowGatewayManagement                 bool          `json:"allowGatewayManagement"`
	AllowTrustedForwardingClass            bool          `json:"allowTrustedForwardingClass"`
	AllowedForwardingClasses               []interface{} `json:"allowedForwardingClasses,omitempty"`
	AllowedForwardingMode                  string        `json:"allowedForwardingMode,omitempty"`
	FloatingIPsQuota                       int           `json:"floatingIPsQuota,omitempty"`
	FloatingIPsUsed                        int           `json:"floatingIPsUsed,omitempty"`
	FlowCollectionEnabled                  string        `json:"flowCollectionEnabled,omitempty"`
	EmbeddedMetadata                       []interface{} `json:"embeddedMetadata,omitempty"`
	EnableApplicationPerformanceManagement bool          `json:"enableApplicationPerformanceManagement"`
	EncryptionManagementMode               string        `json:"encryptionManagementMode,omitempty"`
	EnterpriseProfileID                    string        `json:"enterpriseProfileID,omitempty"`
	EntityScope                            string        `json:"entityScope,omitempty"`
	LocalAS                                int           `json:"localAS,omitempty"`
	ForwardingClass                        []interface{} `json:"forwardingClass,omitempty"`
	UseGlobalMAC                           bool          `json:"useGlobalMAC"`
	AssociatedEnterpriseSecurityID         string        `json:"associatedEnterpriseSecurityID,omitempty"`
	AssociatedGroupKeyEncryptionProfileID  string        `json:"associatedGroupKeyEncryptionProfileID,omitempty"`
	AssociatedKeyServerMonitorID           string        `json:"associatedKeyServerMonitorID,omitempty"`
	CustomerID                             int           `json:"customerID,omitempty"`
	AvatarData                             string        `json:"avatarData,omitempty"`
	AvatarType                             string        `json:"avatarType,omitempty"`
	ExternalID                             string        `json:"externalID,omitempty"`
}

// NewEnterprise returns a new *Enterprise
func NewEnterprise() *Enterprise {

	return &Enterprise{
		VNFManagementEnabled:                   false,
		WebFilterEnabled:                       false,
		DictionaryVersion:                      2,
		VirtualFirewallRulesEnabled:            false,
		FlowCollectionEnabled:                  "DISABLED",
		EnableApplicationPerformanceManagement: false,
		UseGlobalMAC:                           false,
	}
}

// Identity returns the Identity of the object.
func (o *Enterprise) Identity() bambou.Identity {

	return EnterpriseIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Enterprise) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Enterprise) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Enterprise from the server
func (o *Enterprise) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Enterprise into the server
func (o *Enterprise) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Enterprise from the server
func (o *Enterprise) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// L2Domains retrieves the list of child L2Domains of the Enterprise
func (o *Enterprise) L2Domains(info *bambou.FetchingInfo) (L2DomainsList, *bambou.Error) {

	var list L2DomainsList
	err := bambou.CurrentSession().FetchChildren(o, L2DomainIdentity, &list, info)
	return list, err
}

// CreateL2Domain creates a new child L2Domain under the Enterprise
func (o *Enterprise) CreateL2Domain(child *L2Domain) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// L2DomainTemplates retrieves the list of child L2DomainTemplates of the Enterprise
func (o *Enterprise) L2DomainTemplates(info *bambou.FetchingInfo) (L2DomainTemplatesList, *bambou.Error) {

	var list L2DomainTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, L2DomainTemplateIdentity, &list, info)
	return list, err
}

// CreateL2DomainTemplate creates a new child L2DomainTemplate under the Enterprise
func (o *Enterprise) CreateL2DomainTemplate(child *L2DomainTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// L4Services retrieves the list of child L4Services of the Enterprise
func (o *Enterprise) L4Services(info *bambou.FetchingInfo) (L4ServicesList, *bambou.Error) {

	var list L4ServicesList
	err := bambou.CurrentSession().FetchChildren(o, L4ServiceIdentity, &list, info)
	return list, err
}

// CreateL4Service creates a new child L4Service under the Enterprise
func (o *Enterprise) CreateL4Service(child *L4Service) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// L4ServiceGroups retrieves the list of child L4ServiceGroups of the Enterprise
func (o *Enterprise) L4ServiceGroups(info *bambou.FetchingInfo) (L4ServiceGroupsList, *bambou.Error) {

	var list L4ServiceGroupsList
	err := bambou.CurrentSession().FetchChildren(o, L4ServiceGroupIdentity, &list, info)
	return list, err
}

// CreateL4ServiceGroup creates a new child L4ServiceGroup under the Enterprise
func (o *Enterprise) CreateL4ServiceGroup(child *L4ServiceGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// L7applicationsignatures retrieves the list of child L7applicationsignatures of the Enterprise
func (o *Enterprise) L7applicationsignatures(info *bambou.FetchingInfo) (L7applicationsignaturesList, *bambou.Error) {

	var list L7applicationsignaturesList
	err := bambou.CurrentSession().FetchChildren(o, L7applicationsignatureIdentity, &list, info)
	return list, err
}

// SaaSApplicationGroups retrieves the list of child SaaSApplicationGroups of the Enterprise
func (o *Enterprise) SaaSApplicationGroups(info *bambou.FetchingInfo) (SaaSApplicationGroupsList, *bambou.Error) {

	var list SaaSApplicationGroupsList
	err := bambou.CurrentSession().FetchChildren(o, SaaSApplicationGroupIdentity, &list, info)
	return list, err
}

// CreateSaaSApplicationGroup creates a new child SaaSApplicationGroup under the Enterprise
func (o *Enterprise) CreateSaaSApplicationGroup(child *SaaSApplicationGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// SaaSApplicationTypes retrieves the list of child SaaSApplicationTypes of the Enterprise
func (o *Enterprise) SaaSApplicationTypes(info *bambou.FetchingInfo) (SaaSApplicationTypesList, *bambou.Error) {

	var list SaaSApplicationTypesList
	err := bambou.CurrentSession().FetchChildren(o, SaaSApplicationTypeIdentity, &list, info)
	return list, err
}

// CreateSaaSApplicationType creates a new child SaaSApplicationType under the Enterprise
func (o *Enterprise) CreateSaaSApplicationType(child *SaaSApplicationType) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CaptivePortalProfiles retrieves the list of child CaptivePortalProfiles of the Enterprise
func (o *Enterprise) CaptivePortalProfiles(info *bambou.FetchingInfo) (CaptivePortalProfilesList, *bambou.Error) {

	var list CaptivePortalProfilesList
	err := bambou.CurrentSession().FetchChildren(o, CaptivePortalProfileIdentity, &list, info)
	return list, err
}

// CreateCaptivePortalProfile creates a new child CaptivePortalProfile under the Enterprise
func (o *Enterprise) CreateCaptivePortalProfile(child *CaptivePortalProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// RateLimiters retrieves the list of child RateLimiters of the Enterprise
func (o *Enterprise) RateLimiters(info *bambou.FetchingInfo) (RateLimitersList, *bambou.Error) {

	var list RateLimitersList
	err := bambou.CurrentSession().FetchChildren(o, RateLimiterIdentity, &list, info)
	return list, err
}

// CreateRateLimiter creates a new child RateLimiter under the Enterprise
func (o *Enterprise) CreateRateLimiter(child *RateLimiter) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Gateways retrieves the list of child Gateways of the Enterprise
func (o *Enterprise) Gateways(info *bambou.FetchingInfo) (GatewaysList, *bambou.Error) {

	var list GatewaysList
	err := bambou.CurrentSession().FetchChildren(o, GatewayIdentity, &list, info)
	return list, err
}

// CreateGateway creates a new child Gateway under the Enterprise
func (o *Enterprise) CreateGateway(child *Gateway) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GatewaysLocations retrieves the list of child GatewaysLocations of the Enterprise
func (o *Enterprise) GatewaysLocations(info *bambou.FetchingInfo) (GatewaysLocationsList, *bambou.Error) {

	var list GatewaysLocationsList
	err := bambou.CurrentSession().FetchChildren(o, GatewaysLocationIdentity, &list, info)
	return list, err
}

// GatewayTemplates retrieves the list of child GatewayTemplates of the Enterprise
func (o *Enterprise) GatewayTemplates(info *bambou.FetchingInfo) (GatewayTemplatesList, *bambou.Error) {

	var list GatewayTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, GatewayTemplateIdentity, &list, info)
	return list, err
}

// CreateGatewayTemplate creates a new child GatewayTemplate under the Enterprise
func (o *Enterprise) CreateGatewayTemplate(child *GatewayTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PATNATPools retrieves the list of child PATNATPools of the Enterprise
func (o *Enterprise) PATNATPools(info *bambou.FetchingInfo) (PATNATPoolsList, *bambou.Error) {

	var list PATNATPoolsList
	err := bambou.CurrentSession().FetchChildren(o, PATNATPoolIdentity, &list, info)
	return list, err
}

// LDAPConfigurations retrieves the list of child LDAPConfigurations of the Enterprise
func (o *Enterprise) LDAPConfigurations(info *bambou.FetchingInfo) (LDAPConfigurationsList, *bambou.Error) {

	var list LDAPConfigurationsList
	err := bambou.CurrentSession().FetchChildren(o, LDAPConfigurationIdentity, &list, info)
	return list, err
}

// WebCategories retrieves the list of child WebCategories of the Enterprise
func (o *Enterprise) WebCategories(info *bambou.FetchingInfo) (WebCategoriesList, *bambou.Error) {

	var list WebCategoriesList
	err := bambou.CurrentSession().FetchChildren(o, WebCategoryIdentity, &list, info)
	return list, err
}

// CreateWebCategory creates a new child WebCategory under the Enterprise
func (o *Enterprise) CreateWebCategory(child *WebCategory) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// WebDomainNames retrieves the list of child WebDomainNames of the Enterprise
func (o *Enterprise) WebDomainNames(info *bambou.FetchingInfo) (WebDomainNamesList, *bambou.Error) {

	var list WebDomainNamesList
	err := bambou.CurrentSession().FetchChildren(o, WebDomainNameIdentity, &list, info)
	return list, err
}

// CreateWebDomainName creates a new child WebDomainName under the Enterprise
func (o *Enterprise) CreateWebDomainName(child *WebDomainName) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// RedundancyGroups retrieves the list of child RedundancyGroups of the Enterprise
func (o *Enterprise) RedundancyGroups(info *bambou.FetchingInfo) (RedundancyGroupsList, *bambou.Error) {

	var list RedundancyGroupsList
	err := bambou.CurrentSession().FetchChildren(o, RedundancyGroupIdentity, &list, info)
	return list, err
}

// CreateRedundancyGroup creates a new child RedundancyGroup under the Enterprise
func (o *Enterprise) CreateRedundancyGroup(child *RedundancyGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DeploymentFailures retrieves the list of child DeploymentFailures of the Enterprise
func (o *Enterprise) DeploymentFailures(info *bambou.FetchingInfo) (DeploymentFailuresList, *bambou.Error) {

	var list DeploymentFailuresList
	err := bambou.CurrentSession().FetchChildren(o, DeploymentFailureIdentity, &list, info)
	return list, err
}

// PerformanceMonitors retrieves the list of child PerformanceMonitors of the Enterprise
func (o *Enterprise) PerformanceMonitors(info *bambou.FetchingInfo) (PerformanceMonitorsList, *bambou.Error) {

	var list PerformanceMonitorsList
	err := bambou.CurrentSession().FetchChildren(o, PerformanceMonitorIdentity, &list, info)
	return list, err
}

// CreatePerformanceMonitor creates a new child PerformanceMonitor under the Enterprise
func (o *Enterprise) CreatePerformanceMonitor(child *PerformanceMonitor) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// TestDefinitions retrieves the list of child TestDefinitions of the Enterprise
func (o *Enterprise) TestDefinitions(info *bambou.FetchingInfo) (TestDefinitionsList, *bambou.Error) {

	var list TestDefinitionsList
	err := bambou.CurrentSession().FetchChildren(o, TestDefinitionIdentity, &list, info)
	return list, err
}

// CreateTestDefinition creates a new child TestDefinition under the Enterprise
func (o *Enterprise) CreateTestDefinition(child *TestDefinition) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// TestSuites retrieves the list of child TestSuites of the Enterprise
func (o *Enterprise) TestSuites(info *bambou.FetchingInfo) (TestSuitesList, *bambou.Error) {

	var list TestSuitesList
	err := bambou.CurrentSession().FetchChildren(o, TestSuiteIdentity, &list, info)
	return list, err
}

// CreateTestSuite creates a new child TestSuite under the Enterprise
func (o *Enterprise) CreateTestSuite(child *TestSuite) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the Enterprise
func (o *Enterprise) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the Enterprise
func (o *Enterprise) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NetconfProfiles retrieves the list of child NetconfProfiles of the Enterprise
func (o *Enterprise) NetconfProfiles(info *bambou.FetchingInfo) (NetconfProfilesList, *bambou.Error) {

	var list NetconfProfilesList
	err := bambou.CurrentSession().FetchChildren(o, NetconfProfileIdentity, &list, info)
	return list, err
}

// CreateNetconfProfile creates a new child NetconfProfile under the Enterprise
func (o *Enterprise) CreateNetconfProfile(child *NetconfProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NetworkMacroGroups retrieves the list of child NetworkMacroGroups of the Enterprise
func (o *Enterprise) NetworkMacroGroups(info *bambou.FetchingInfo) (NetworkMacroGroupsList, *bambou.Error) {

	var list NetworkMacroGroupsList
	err := bambou.CurrentSession().FetchChildren(o, NetworkMacroGroupIdentity, &list, info)
	return list, err
}

// CreateNetworkMacroGroup creates a new child NetworkMacroGroup under the Enterprise
func (o *Enterprise) CreateNetworkMacroGroup(child *NetworkMacroGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NetworkPerformanceMeasurements retrieves the list of child NetworkPerformanceMeasurements of the Enterprise
func (o *Enterprise) NetworkPerformanceMeasurements(info *bambou.FetchingInfo) (NetworkPerformanceMeasurementsList, *bambou.Error) {

	var list NetworkPerformanceMeasurementsList
	err := bambou.CurrentSession().FetchChildren(o, NetworkPerformanceMeasurementIdentity, &list, info)
	return list, err
}

// CreateNetworkPerformanceMeasurement creates a new child NetworkPerformanceMeasurement under the Enterprise
func (o *Enterprise) CreateNetworkPerformanceMeasurement(child *NetworkPerformanceMeasurement) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// KeyServerMonitors retrieves the list of child KeyServerMonitors of the Enterprise
func (o *Enterprise) KeyServerMonitors(info *bambou.FetchingInfo) (KeyServerMonitorsList, *bambou.Error) {

	var list KeyServerMonitorsList
	err := bambou.CurrentSession().FetchChildren(o, KeyServerMonitorIdentity, &list, info)
	return list, err
}

// ZFBRequests retrieves the list of child ZFBRequests of the Enterprise
func (o *Enterprise) ZFBRequests(info *bambou.FetchingInfo) (ZFBRequestsList, *bambou.Error) {

	var list ZFBRequestsList
	err := bambou.CurrentSession().FetchChildren(o, ZFBRequestIdentity, &list, info)
	return list, err
}

// CreateZFBRequest creates a new child ZFBRequest under the Enterprise
func (o *Enterprise) CreateZFBRequest(child *ZFBRequest) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// BGPProfiles retrieves the list of child BGPProfiles of the Enterprise
func (o *Enterprise) BGPProfiles(info *bambou.FetchingInfo) (BGPProfilesList, *bambou.Error) {

	var list BGPProfilesList
	err := bambou.CurrentSession().FetchChildren(o, BGPProfileIdentity, &list, info)
	return list, err
}

// CreateBGPProfile creates a new child BGPProfile under the Enterprise
func (o *Enterprise) CreateBGPProfile(child *BGPProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EgressQOSPolicies retrieves the list of child EgressQOSPolicies of the Enterprise
func (o *Enterprise) EgressQOSPolicies(info *bambou.FetchingInfo) (EgressQOSPoliciesList, *bambou.Error) {

	var list EgressQOSPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, EgressQOSPolicyIdentity, &list, info)
	return list, err
}

// CreateEgressQOSPolicy creates a new child EgressQOSPolicy under the Enterprise
func (o *Enterprise) CreateEgressQOSPolicy(child *EgressQOSPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// SharedNetworkResources retrieves the list of child SharedNetworkResources of the Enterprise
func (o *Enterprise) SharedNetworkResources(info *bambou.FetchingInfo) (SharedNetworkResourcesList, *bambou.Error) {

	var list SharedNetworkResourcesList
	err := bambou.CurrentSession().FetchChildren(o, SharedNetworkResourceIdentity, &list, info)
	return list, err
}

// FirewallAcls retrieves the list of child FirewallAcls of the Enterprise
func (o *Enterprise) FirewallAcls(info *bambou.FetchingInfo) (FirewallAclsList, *bambou.Error) {

	var list FirewallAclsList
	err := bambou.CurrentSession().FetchChildren(o, FirewallAclIdentity, &list, info)
	return list, err
}

// CreateFirewallAcl creates a new child FirewallAcl under the Enterprise
func (o *Enterprise) CreateFirewallAcl(child *FirewallAcl) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// FirewallRules retrieves the list of child FirewallRules of the Enterprise
func (o *Enterprise) FirewallRules(info *bambou.FetchingInfo) (FirewallRulesList, *bambou.Error) {

	var list FirewallRulesList
	err := bambou.CurrentSession().FetchChildren(o, FirewallRuleIdentity, &list, info)
	return list, err
}

// CreateFirewallRule creates a new child FirewallRule under the Enterprise
func (o *Enterprise) CreateFirewallRule(child *FirewallRule) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IKECertificates retrieves the list of child IKECertificates of the Enterprise
func (o *Enterprise) IKECertificates(info *bambou.FetchingInfo) (IKECertificatesList, *bambou.Error) {

	var list IKECertificatesList
	err := bambou.CurrentSession().FetchChildren(o, IKECertificateIdentity, &list, info)
	return list, err
}

// CreateIKECertificate creates a new child IKECertificate under the Enterprise
func (o *Enterprise) CreateIKECertificate(child *IKECertificate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IKEEncryptionprofiles retrieves the list of child IKEEncryptionprofiles of the Enterprise
func (o *Enterprise) IKEEncryptionprofiles(info *bambou.FetchingInfo) (IKEEncryptionprofilesList, *bambou.Error) {

	var list IKEEncryptionprofilesList
	err := bambou.CurrentSession().FetchChildren(o, IKEEncryptionprofileIdentity, &list, info)
	return list, err
}

// CreateIKEEncryptionprofile creates a new child IKEEncryptionprofile under the Enterprise
func (o *Enterprise) CreateIKEEncryptionprofile(child *IKEEncryptionprofile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IKEGateways retrieves the list of child IKEGateways of the Enterprise
func (o *Enterprise) IKEGateways(info *bambou.FetchingInfo) (IKEGatewaysList, *bambou.Error) {

	var list IKEGatewaysList
	err := bambou.CurrentSession().FetchChildren(o, IKEGatewayIdentity, &list, info)
	return list, err
}

// CreateIKEGateway creates a new child IKEGateway under the Enterprise
func (o *Enterprise) CreateIKEGateway(child *IKEGateway) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IKEGatewayProfiles retrieves the list of child IKEGatewayProfiles of the Enterprise
func (o *Enterprise) IKEGatewayProfiles(info *bambou.FetchingInfo) (IKEGatewayProfilesList, *bambou.Error) {

	var list IKEGatewayProfilesList
	err := bambou.CurrentSession().FetchChildren(o, IKEGatewayProfileIdentity, &list, info)
	return list, err
}

// CreateIKEGatewayProfile creates a new child IKEGatewayProfile under the Enterprise
func (o *Enterprise) CreateIKEGatewayProfile(child *IKEGatewayProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IKEPSKs retrieves the list of child IKEPSKs of the Enterprise
func (o *Enterprise) IKEPSKs(info *bambou.FetchingInfo) (IKEPSKsList, *bambou.Error) {

	var list IKEPSKsList
	err := bambou.CurrentSession().FetchChildren(o, IKEPSKIdentity, &list, info)
	return list, err
}

// CreateIKEPSK creates a new child IKEPSK under the Enterprise
func (o *Enterprise) CreateIKEPSK(child *IKEPSK) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Alarms retrieves the list of child Alarms of the Enterprise
func (o *Enterprise) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// AllAlarms retrieves the list of child AllAlarms of the Enterprise
func (o *Enterprise) AllAlarms(info *bambou.FetchingInfo) (AllAlarmsList, *bambou.Error) {

	var list AllAlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AllAlarmIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Enterprise
func (o *Enterprise) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Enterprise
func (o *Enterprise) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMs retrieves the list of child VMs of the Enterprise
func (o *Enterprise) VMs(info *bambou.FetchingInfo) (VMsList, *bambou.Error) {

	var list VMsList
	err := bambou.CurrentSession().FetchChildren(o, VMIdentity, &list, info)
	return list, err
}

// VNFs retrieves the list of child VNFs of the Enterprise
func (o *Enterprise) VNFs(info *bambou.FetchingInfo) (VNFsList, *bambou.Error) {

	var list VNFsList
	err := bambou.CurrentSession().FetchChildren(o, VNFIdentity, &list, info)
	return list, err
}

// CreateVNF creates a new child VNF under the Enterprise
func (o *Enterprise) CreateVNF(child *VNF) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VNFMetadatas retrieves the list of child VNFMetadatas of the Enterprise
func (o *Enterprise) VNFMetadatas(info *bambou.FetchingInfo) (VNFMetadatasList, *bambou.Error) {

	var list VNFMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, VNFMetadataIdentity, &list, info)
	return list, err
}

// CreateVNFMetadata creates a new child VNFMetadata under the Enterprise
func (o *Enterprise) CreateVNFMetadata(child *VNFMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VNFThresholdPolicies retrieves the list of child VNFThresholdPolicies of the Enterprise
func (o *Enterprise) VNFThresholdPolicies(info *bambou.FetchingInfo) (VNFThresholdPoliciesList, *bambou.Error) {

	var list VNFThresholdPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, VNFThresholdPolicyIdentity, &list, info)
	return list, err
}

// CreateVNFThresholdPolicy creates a new child VNFThresholdPolicy under the Enterprise
func (o *Enterprise) CreateVNFThresholdPolicy(child *VNFThresholdPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IngressQOSPolicies retrieves the list of child IngressQOSPolicies of the Enterprise
func (o *Enterprise) IngressQOSPolicies(info *bambou.FetchingInfo) (IngressQOSPoliciesList, *bambou.Error) {

	var list IngressQOSPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, IngressQOSPolicyIdentity, &list, info)
	return list, err
}

// CreateIngressQOSPolicy creates a new child IngressQOSPolicy under the Enterprise
func (o *Enterprise) CreateIngressQOSPolicy(child *IngressQOSPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EnterpriseNetworks retrieves the list of child EnterpriseNetworks of the Enterprise
func (o *Enterprise) EnterpriseNetworks(info *bambou.FetchingInfo) (EnterpriseNetworksList, *bambou.Error) {

	var list EnterpriseNetworksList
	err := bambou.CurrentSession().FetchChildren(o, EnterpriseNetworkIdentity, &list, info)
	return list, err
}

// CreateEnterpriseNetwork creates a new child EnterpriseNetwork under the Enterprise
func (o *Enterprise) CreateEnterpriseNetwork(child *EnterpriseNetwork) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EnterpriseSecurities retrieves the list of child EnterpriseSecurities of the Enterprise
func (o *Enterprise) EnterpriseSecurities(info *bambou.FetchingInfo) (EnterpriseSecuritiesList, *bambou.Error) {

	var list EnterpriseSecuritiesList
	err := bambou.CurrentSession().FetchChildren(o, EnterpriseSecurityIdentity, &list, info)
	return list, err
}

// Jobs retrieves the list of child Jobs of the Enterprise
func (o *Enterprise) Jobs(info *bambou.FetchingInfo) (JobsList, *bambou.Error) {

	var list JobsList
	err := bambou.CurrentSession().FetchChildren(o, JobIdentity, &list, info)
	return list, err
}

// CreateJob creates a new child Job under the Enterprise
func (o *Enterprise) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PolicyGroupCategories retrieves the list of child PolicyGroupCategories of the Enterprise
func (o *Enterprise) PolicyGroupCategories(info *bambou.FetchingInfo) (PolicyGroupCategoriesList, *bambou.Error) {

	var list PolicyGroupCategoriesList
	err := bambou.CurrentSession().FetchChildren(o, PolicyGroupCategoryIdentity, &list, info)
	return list, err
}

// CreatePolicyGroupCategory creates a new child PolicyGroupCategory under the Enterprise
func (o *Enterprise) CreatePolicyGroupCategory(child *PolicyGroupCategory) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PolicyObjectGroups retrieves the list of child PolicyObjectGroups of the Enterprise
func (o *Enterprise) PolicyObjectGroups(info *bambou.FetchingInfo) (PolicyObjectGroupsList, *bambou.Error) {

	var list PolicyObjectGroupsList
	err := bambou.CurrentSession().FetchChildren(o, PolicyObjectGroupIdentity, &list, info)
	return list, err
}

// CreatePolicyObjectGroup creates a new child PolicyObjectGroup under the Enterprise
func (o *Enterprise) CreatePolicyObjectGroup(child *PolicyObjectGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Domains retrieves the list of child Domains of the Enterprise
func (o *Enterprise) Domains(info *bambou.FetchingInfo) (DomainsList, *bambou.Error) {

	var list DomainsList
	err := bambou.CurrentSession().FetchChildren(o, DomainIdentity, &list, info)
	return list, err
}

// CreateDomain creates a new child Domain under the Enterprise
func (o *Enterprise) CreateDomain(child *Domain) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DomainKindSummaries retrieves the list of child DomainKindSummaries of the Enterprise
func (o *Enterprise) DomainKindSummaries(info *bambou.FetchingInfo) (DomainKindSummariesList, *bambou.Error) {

	var list DomainKindSummariesList
	err := bambou.CurrentSession().FetchChildren(o, DomainKindSummaryIdentity, &list, info)
	return list, err
}

// DomainTemplates retrieves the list of child DomainTemplates of the Enterprise
func (o *Enterprise) DomainTemplates(info *bambou.FetchingInfo) (DomainTemplatesList, *bambou.Error) {

	var list DomainTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, DomainTemplateIdentity, &list, info)
	return list, err
}

// CreateDomainTemplate creates a new child DomainTemplate under the Enterprise
func (o *Enterprise) CreateDomainTemplate(child *DomainTemplate) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Containers retrieves the list of child Containers of the Enterprise
func (o *Enterprise) Containers(info *bambou.FetchingInfo) (ContainersList, *bambou.Error) {

	var list ContainersList
	err := bambou.CurrentSession().FetchChildren(o, ContainerIdentity, &list, info)
	return list, err
}

// COSRemarkingPolicyTables retrieves the list of child COSRemarkingPolicyTables of the Enterprise
func (o *Enterprise) COSRemarkingPolicyTables(info *bambou.FetchingInfo) (COSRemarkingPolicyTablesList, *bambou.Error) {

	var list COSRemarkingPolicyTablesList
	err := bambou.CurrentSession().FetchChildren(o, COSRemarkingPolicyTableIdentity, &list, info)
	return list, err
}

// CreateCOSRemarkingPolicyTable creates a new child COSRemarkingPolicyTable under the Enterprise
func (o *Enterprise) CreateCOSRemarkingPolicyTable(child *COSRemarkingPolicyTable) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// RoutingPolicies retrieves the list of child RoutingPolicies of the Enterprise
func (o *Enterprise) RoutingPolicies(info *bambou.FetchingInfo) (RoutingPoliciesList, *bambou.Error) {

	var list RoutingPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, RoutingPolicyIdentity, &list, info)
	return list, err
}

// CreateRoutingPolicy creates a new child RoutingPolicy under the Enterprise
func (o *Enterprise) CreateRoutingPolicy(child *RoutingPolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Applications retrieves the list of child Applications of the Enterprise
func (o *Enterprise) Applications(info *bambou.FetchingInfo) (ApplicationsList, *bambou.Error) {

	var list ApplicationsList
	err := bambou.CurrentSession().FetchChildren(o, ApplicationIdentity, &list, info)
	return list, err
}

// CreateApplication creates a new child Application under the Enterprise
func (o *Enterprise) CreateApplication(child *Application) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Applicationperformancemanagements retrieves the list of child Applicationperformancemanagements of the Enterprise
func (o *Enterprise) Applicationperformancemanagements(info *bambou.FetchingInfo) (ApplicationperformancemanagementsList, *bambou.Error) {

	var list ApplicationperformancemanagementsList
	err := bambou.CurrentSession().FetchChildren(o, ApplicationperformancemanagementIdentity, &list, info)
	return list, err
}

// CreateApplicationperformancemanagement creates a new child Applicationperformancemanagement under the Enterprise
func (o *Enterprise) CreateApplicationperformancemanagement(child *Applicationperformancemanagement) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Groups retrieves the list of child Groups of the Enterprise
func (o *Enterprise) Groups(info *bambou.FetchingInfo) (GroupsList, *bambou.Error) {

	var list GroupsList
	err := bambou.CurrentSession().FetchChildren(o, GroupIdentity, &list, info)
	return list, err
}

// CreateGroup creates a new child Group under the Enterprise
func (o *Enterprise) CreateGroup(child *Group) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GroupKeyEncryptionProfiles retrieves the list of child GroupKeyEncryptionProfiles of the Enterprise
func (o *Enterprise) GroupKeyEncryptionProfiles(info *bambou.FetchingInfo) (GroupKeyEncryptionProfilesList, *bambou.Error) {

	var list GroupKeyEncryptionProfilesList
	err := bambou.CurrentSession().FetchChildren(o, GroupKeyEncryptionProfileIdentity, &list, info)
	return list, err
}

// Trunks retrieves the list of child Trunks of the Enterprise
func (o *Enterprise) Trunks(info *bambou.FetchingInfo) (TrunksList, *bambou.Error) {

	var list TrunksList
	err := bambou.CurrentSession().FetchChildren(o, TrunkIdentity, &list, info)
	return list, err
}

// CreateTrunk creates a new child Trunk under the Enterprise
func (o *Enterprise) CreateTrunk(child *Trunk) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DSCPForwardingClassTables retrieves the list of child DSCPForwardingClassTables of the Enterprise
func (o *Enterprise) DSCPForwardingClassTables(info *bambou.FetchingInfo) (DSCPForwardingClassTablesList, *bambou.Error) {

	var list DSCPForwardingClassTablesList
	err := bambou.CurrentSession().FetchChildren(o, DSCPForwardingClassTableIdentity, &list, info)
	return list, err
}

// CreateDSCPForwardingClassTable creates a new child DSCPForwardingClassTable under the Enterprise
func (o *Enterprise) CreateDSCPForwardingClassTable(child *DSCPForwardingClassTable) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DSCPRemarkingPolicyTables retrieves the list of child DSCPRemarkingPolicyTables of the Enterprise
func (o *Enterprise) DSCPRemarkingPolicyTables(info *bambou.FetchingInfo) (DSCPRemarkingPolicyTablesList, *bambou.Error) {

	var list DSCPRemarkingPolicyTablesList
	err := bambou.CurrentSession().FetchChildren(o, DSCPRemarkingPolicyTableIdentity, &list, info)
	return list, err
}

// CreateDSCPRemarkingPolicyTable creates a new child DSCPRemarkingPolicyTable under the Enterprise
func (o *Enterprise) CreateDSCPRemarkingPolicyTable(child *DSCPRemarkingPolicyTable) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Users retrieves the list of child Users of the Enterprise
func (o *Enterprise) Users(info *bambou.FetchingInfo) (UsersList, *bambou.Error) {

	var list UsersList
	err := bambou.CurrentSession().FetchChildren(o, UserIdentity, &list, info)
	return list, err
}

// CreateUser creates a new child User under the Enterprise
func (o *Enterprise) CreateUser(child *User) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSGateways retrieves the list of child NSGateways of the Enterprise
func (o *Enterprise) NSGateways(info *bambou.FetchingInfo) (NSGatewaysList, *bambou.Error) {

	var list NSGatewaysList
	err := bambou.CurrentSession().FetchChildren(o, NSGatewayIdentity, &list, info)
	return list, err
}

// CreateNSGateway creates a new child NSGateway under the Enterprise
func (o *Enterprise) CreateNSGateway(child *NSGateway) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSGatewaysCounts retrieves the list of child NSGatewaysCounts of the Enterprise
func (o *Enterprise) NSGatewaysCounts(info *bambou.FetchingInfo) (NSGatewaysCountsList, *bambou.Error) {

	var list NSGatewaysCountsList
	err := bambou.CurrentSession().FetchChildren(o, NSGatewaysCountIdentity, &list, info)
	return list, err
}

// NSGatewaySummaries retrieves the list of child NSGatewaySummaries of the Enterprise
func (o *Enterprise) NSGatewaySummaries(info *bambou.FetchingInfo) (NSGatewaySummariesList, *bambou.Error) {

	var list NSGatewaySummariesList
	err := bambou.CurrentSession().FetchChildren(o, NSGatewaySummaryIdentity, &list, info)
	return list, err
}

// NSGatewayTemplates retrieves the list of child NSGatewayTemplates of the Enterprise
func (o *Enterprise) NSGatewayTemplates(info *bambou.FetchingInfo) (NSGatewayTemplatesList, *bambou.Error) {

	var list NSGatewayTemplatesList
	err := bambou.CurrentSession().FetchChildren(o, NSGatewayTemplateIdentity, &list, info)
	return list, err
}

// NSGGroups retrieves the list of child NSGGroups of the Enterprise
func (o *Enterprise) NSGGroups(info *bambou.FetchingInfo) (NSGGroupsList, *bambou.Error) {

	var list NSGGroupsList
	err := bambou.CurrentSession().FetchChildren(o, NSGGroupIdentity, &list, info)
	return list, err
}

// CreateNSGGroup creates a new child NSGGroup under the Enterprise
func (o *Enterprise) CreateNSGGroup(child *NSGGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSRedundantGatewayGroups retrieves the list of child NSRedundantGatewayGroups of the Enterprise
func (o *Enterprise) NSRedundantGatewayGroups(info *bambou.FetchingInfo) (NSRedundantGatewayGroupsList, *bambou.Error) {

	var list NSRedundantGatewayGroupsList
	err := bambou.CurrentSession().FetchChildren(o, NSRedundantGatewayGroupIdentity, &list, info)
	return list, err
}

// CreateNSRedundantGatewayGroup creates a new child NSRedundantGatewayGroup under the Enterprise
func (o *Enterprise) CreateNSRedundantGatewayGroup(child *NSRedundantGatewayGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PublicNetworkMacros retrieves the list of child PublicNetworkMacros of the Enterprise
func (o *Enterprise) PublicNetworkMacros(info *bambou.FetchingInfo) (PublicNetworkMacrosList, *bambou.Error) {

	var list PublicNetworkMacrosList
	err := bambou.CurrentSession().FetchChildren(o, PublicNetworkMacroIdentity, &list, info)
	return list, err
}

// CreatePublicNetworkMacro creates a new child PublicNetworkMacro under the Enterprise
func (o *Enterprise) CreatePublicNetworkMacro(child *PublicNetworkMacro) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// MultiCastLists retrieves the list of child MultiCastLists of the Enterprise
func (o *Enterprise) MultiCastLists(info *bambou.FetchingInfo) (MultiCastListsList, *bambou.Error) {

	var list MultiCastListsList
	err := bambou.CurrentSession().FetchChildren(o, MultiCastListIdentity, &list, info)
	return list, err
}

// Avatars retrieves the list of child Avatars of the Enterprise
func (o *Enterprise) Avatars(info *bambou.FetchingInfo) (AvatarsList, *bambou.Error) {

	var list AvatarsList
	err := bambou.CurrentSession().FetchChildren(o, AvatarIdentity, &list, info)
	return list, err
}

// CreateAvatar creates a new child Avatar under the Enterprise
func (o *Enterprise) CreateAvatar(child *Avatar) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the Enterprise
func (o *Enterprise) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}

// OverlayManagementProfiles retrieves the list of child OverlayManagementProfiles of the Enterprise
func (o *Enterprise) OverlayManagementProfiles(info *bambou.FetchingInfo) (OverlayManagementProfilesList, *bambou.Error) {

	var list OverlayManagementProfilesList
	err := bambou.CurrentSession().FetchChildren(o, OverlayManagementProfileIdentity, &list, info)
	return list, err
}

// CreateOverlayManagementProfile creates a new child OverlayManagementProfile under the Enterprise
func (o *Enterprise) CreateOverlayManagementProfile(child *OverlayManagementProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// SyslogDestinations retrieves the list of child SyslogDestinations of the Enterprise
func (o *Enterprise) SyslogDestinations(info *bambou.FetchingInfo) (SyslogDestinationsList, *bambou.Error) {

	var list SyslogDestinationsList
	err := bambou.CurrentSession().FetchChildren(o, SyslogDestinationIdentity, &list, info)
	return list, err
}

// CreateSyslogDestination creates a new child SyslogDestination under the Enterprise
func (o *Enterprise) CreateSyslogDestination(child *SyslogDestination) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// AzureClouds retrieves the list of child AzureClouds of the Enterprise
func (o *Enterprise) AzureClouds(info *bambou.FetchingInfo) (AzureCloudsList, *bambou.Error) {

	var list AzureCloudsList
	err := bambou.CurrentSession().FetchChildren(o, AzureCloudIdentity, &list, info)
	return list, err
}

// CreateAzureCloud creates a new child AzureCloud under the Enterprise
func (o *Enterprise) CreateAzureCloud(child *AzureCloud) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
