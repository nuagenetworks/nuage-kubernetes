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

// SystemConfigIdentity represents the Identity of the object
var SystemConfigIdentity = bambou.Identity{
	Name:     "systemconfig",
	Category: "systemconfigs",
}

// SystemConfigsList represents a list of SystemConfigs
type SystemConfigsList []*SystemConfig

// SystemConfigsAncestor is the interface that an ancestor of a SystemConfig must implement.
// An Ancestor is defined as an entity that has SystemConfig as a descendant.
// An Ancestor can get a list of its child SystemConfigs, but not necessarily create one.
type SystemConfigsAncestor interface {
	SystemConfigs(*bambou.FetchingInfo) (SystemConfigsList, *bambou.Error)
}

// SystemConfigsParent is the interface that a parent of a SystemConfig must implement.
// A Parent is defined as an entity that has SystemConfig as a child.
// A Parent is an Ancestor which can create a SystemConfig.
type SystemConfigsParent interface {
	SystemConfigsAncestor
	CreateSystemConfig(*SystemConfig) *bambou.Error
}

// SystemConfig represents the model of a systemconfig
type SystemConfig struct {
	ID                                                string        `json:"ID,omitempty"`
	ParentID                                          string        `json:"parentID,omitempty"`
	ParentType                                        string        `json:"parentType,omitempty"`
	Owner                                             string        `json:"owner,omitempty"`
	AARFlowStatsInterval                              int           `json:"AARFlowStatsInterval,omitempty"`
	AARProbeStatsInterval                             int           `json:"AARProbeStatsInterval,omitempty"`
	ACLAllowOrigin                                    string        `json:"ACLAllowOrigin,omitempty"`
	ECMPCount                                         int           `json:"ECMPCount,omitempty"`
	LDAPSyncInterval                                  int           `json:"LDAPSyncInterval,omitempty"`
	LDAPTrustStoreCertifcate                          string        `json:"LDAPTrustStoreCertifcate,omitempty"`
	LDAPTrustStorePassword                            string        `json:"LDAPTrustStorePassword,omitempty"`
	ADGatewayPurgeTime                                int           `json:"ADGatewayPurgeTime,omitempty"`
	RDLowerLimit                                      int           `json:"RDLowerLimit,omitempty"`
	RDPublicNetworkLowerLimit                         int           `json:"RDPublicNetworkLowerLimit,omitempty"`
	RDPublicNetworkUpperLimit                         int           `json:"RDPublicNetworkUpperLimit,omitempty"`
	RDUpperLimit                                      int           `json:"RDUpperLimit,omitempty"`
	ZFBBootstrapEnabled                               bool          `json:"ZFBBootstrapEnabled"`
	ZFBRequestRetryTimer                              int           `json:"ZFBRequestRetryTimer,omitempty"`
	ZFBSchedulerStaleRequestTimeout                   int           `json:"ZFBSchedulerStaleRequestTimeout,omitempty"`
	PGIDLowerLimit                                    int           `json:"PGIDLowerLimit,omitempty"`
	PGIDUpperLimit                                    int           `json:"PGIDUpperLimit,omitempty"`
	DHCPOptionSize                                    int           `json:"DHCPOptionSize,omitempty"`
	VLANIDLowerLimit                                  int           `json:"VLANIDLowerLimit,omitempty"`
	VLANIDUpperLimit                                  int           `json:"VLANIDUpperLimit,omitempty"`
	VMCacheSize                                       int           `json:"VMCacheSize,omitempty"`
	VMPurgeTime                                       int           `json:"VMPurgeTime,omitempty"`
	VMResyncDeletionWaitTime                          int           `json:"VMResyncDeletionWaitTime,omitempty"`
	VMResyncOutstandingInterval                       int           `json:"VMResyncOutstandingInterval,omitempty"`
	VMUnreachableCleanupTime                          int           `json:"VMUnreachableCleanupTime,omitempty"`
	VMUnreachableTime                                 int           `json:"VMUnreachableTime,omitempty"`
	VNFTaskTimeout                                    int           `json:"VNFTaskTimeout,omitempty"`
	VNIDLowerLimit                                    int           `json:"VNIDLowerLimit,omitempty"`
	VNIDPublicNetworkLowerLimit                       int           `json:"VNIDPublicNetworkLowerLimit,omitempty"`
	VNIDPublicNetworkUpperLimit                       int           `json:"VNIDPublicNetworkUpperLimit,omitempty"`
	VNIDUpperLimit                                    int           `json:"VNIDUpperLimit,omitempty"`
	APIKeyRenewalInterval                             int           `json:"APIKeyRenewalInterval,omitempty"`
	APIKeyValidity                                    int           `json:"APIKeyValidity,omitempty"`
	VPortInitStatefulTimer                            int           `json:"VPortInitStatefulTimer,omitempty"`
	LRUCacheSizePerSubnet                             int           `json:"LRUCacheSizePerSubnet,omitempty"`
	VSCOnSameVersionAsVSD                             bool          `json:"VSCOnSameVersionAsVSD"`
	VSDAARApplicationVersion                          string        `json:"VSDAARApplicationVersion,omitempty"`
	VSDAARApplicationVersionPublishDate               string        `json:"VSDAARApplicationVersionPublishDate,omitempty"`
	VSDReadOnlyMode                                   bool          `json:"VSDReadOnlyMode"`
	VSDUpgradeIsComplete                              bool          `json:"VSDUpgradeIsComplete"`
	NSGUplinkHoldDownTimer                            int           `json:"NSGUplinkHoldDownTimer,omitempty"`
	ASNumber                                          int           `json:"ASNumber,omitempty"`
	VSSStatsInterval                                  int           `json:"VSSStatsInterval,omitempty"`
	RTLowerLimit                                      int           `json:"RTLowerLimit,omitempty"`
	RTPublicNetworkLowerLimit                         int           `json:"RTPublicNetworkLowerLimit,omitempty"`
	RTPublicNetworkUpperLimit                         int           `json:"RTPublicNetworkUpperLimit,omitempty"`
	RTUpperLimit                                      int           `json:"RTUpperLimit,omitempty"`
	EVPNBGPCommunityTagASNumber                       int           `json:"EVPNBGPCommunityTagASNumber,omitempty"`
	EVPNBGPCommunityTagLowerLimit                     int           `json:"EVPNBGPCommunityTagLowerLimit,omitempty"`
	EVPNBGPCommunityTagUpperLimit                     int           `json:"EVPNBGPCommunityTagUpperLimit,omitempty"`
	SaaSApplicationsPublishDate                       string        `json:"SaaSApplicationsPublishDate,omitempty"`
	PageMaxSize                                       int           `json:"pageMaxSize,omitempty"`
	PageSize                                          int           `json:"pageSize,omitempty"`
	LastExecutedMigrationPhase                        string        `json:"lastExecutedMigrationPhase,omitempty"`
	LastUpdatedBy                                     string        `json:"lastUpdatedBy,omitempty"`
	GatewayProbeInterval                              int           `json:"gatewayProbeInterval,omitempty"`
	GatewayProbeWindow                                int           `json:"gatewayProbeWindow,omitempty"`
	GatewayRebalancingInterval                        int           `json:"gatewayRebalancingInterval,omitempty"`
	MaxFailedLogins                                   int           `json:"maxFailedLogins,omitempty"`
	MaxResponse                                       int           `json:"maxResponse,omitempty"`
	AccumulateLicensesEnabled                         bool          `json:"accumulateLicensesEnabled"`
	VcinLoadBalancerIP                                string        `json:"vcinLoadBalancerIP,omitempty"`
	SecondaryASNumber                                 int           `json:"secondaryASNumber,omitempty"`
	SecondaryRTLowerLimit                             int           `json:"secondaryRTLowerLimit,omitempty"`
	SecondaryRTUpperLimit                             int           `json:"secondaryRTUpperLimit,omitempty"`
	PerDomainVlanIdEnabled                            bool          `json:"perDomainVlanIdEnabled"`
	ServiceIDUpperLimit                               int           `json:"serviceIDUpperLimit,omitempty"`
	KeyServerMonitorEnabled                           bool          `json:"keyServerMonitorEnabled"`
	KeyServerVSDDataSynchronizationInterval           int           `json:"keyServerVSDDataSynchronizationInterval,omitempty"`
	OffsetCustomerID                                  int           `json:"offsetCustomerID,omitempty"`
	OffsetServiceID                                   int           `json:"offsetServiceID,omitempty"`
	VirtualFirewallRulesEnabled                       bool          `json:"virtualFirewallRulesEnabled"`
	EjbcaNSGCertificateProfile                        string        `json:"ejbcaNSGCertificateProfile,omitempty"`
	EjbcaNSGEndEntityProfile                          string        `json:"ejbcaNSGEndEntityProfile,omitempty"`
	EjbcaOCSPResponderCN                              string        `json:"ejbcaOCSPResponderCN,omitempty"`
	EjbcaOCSPResponderURI                             string        `json:"ejbcaOCSPResponderURI,omitempty"`
	EjbcaVspRootCa                                    string        `json:"ejbcaVspRootCa,omitempty"`
	AlarmsMaxPerObject                                int           `json:"alarmsMaxPerObject,omitempty"`
	ElasticClusterName                                string        `json:"elasticClusterName,omitempty"`
	AllowEnterpriseAvatarOnNSG                        bool          `json:"allowEnterpriseAvatarOnNSG"`
	GlobalMACAddress                                  string        `json:"globalMACAddress,omitempty"`
	FlowCollectionEnabled                             bool          `json:"flowCollectionEnabled"`
	EmbeddedMetadata                                  []interface{} `json:"embeddedMetadata,omitempty"`
	EmbeddedMetadataSize                              int           `json:"embeddedMetadataSize,omitempty"`
	ImportedSaaSApplicationsVersion                   string        `json:"importedSaaSApplicationsVersion,omitempty"`
	InactiveTimeout                                   int           `json:"inactiveTimeout,omitempty"`
	InfrastructureBGPASNumber                         int           `json:"infrastructureBGPASNumber,omitempty"`
	EntityScope                                       string        `json:"entityScope,omitempty"`
	DomainTunnelType                                  string        `json:"domainTunnelType,omitempty"`
	GoogleMapsAPIKey                                  string        `json:"googleMapsAPIKey,omitempty"`
	PostProcessorThreadsCount                         int           `json:"postProcessorThreadsCount,omitempty"`
	GroupKeyDefaultSEKGenerationInterval              int           `json:"groupKeyDefaultSEKGenerationInterval,omitempty"`
	GroupKeyDefaultSEKLifetime                        int           `json:"groupKeyDefaultSEKLifetime,omitempty"`
	GroupKeyDefaultSEKPayloadEncryptionAlgorithm      string        `json:"groupKeyDefaultSEKPayloadEncryptionAlgorithm,omitempty"`
	GroupKeyDefaultSEKPayloadSigningAlgorithm         string        `json:"groupKeyDefaultSEKPayloadSigningAlgorithm,omitempty"`
	GroupKeyDefaultSeedGenerationInterval             int           `json:"groupKeyDefaultSeedGenerationInterval,omitempty"`
	GroupKeyDefaultSeedLifetime                       int           `json:"groupKeyDefaultSeedLifetime,omitempty"`
	GroupKeyDefaultSeedPayloadAuthenticationAlgorithm string        `json:"groupKeyDefaultSeedPayloadAuthenticationAlgorithm,omitempty"`
	GroupKeyDefaultSeedPayloadEncryptionAlgorithm     string        `json:"groupKeyDefaultSeedPayloadEncryptionAlgorithm,omitempty"`
	GroupKeyDefaultSeedPayloadSigningAlgorithm        string        `json:"groupKeyDefaultSeedPayloadSigningAlgorithm,omitempty"`
	GroupKeyDefaultTrafficAuthenticationAlgorithm     string        `json:"groupKeyDefaultTrafficAuthenticationAlgorithm,omitempty"`
	GroupKeyDefaultTrafficEncryptionAlgorithm         string        `json:"groupKeyDefaultTrafficEncryptionAlgorithm,omitempty"`
	GroupKeyDefaultTrafficEncryptionKeyLifetime       int           `json:"groupKeyDefaultTrafficEncryptionKeyLifetime,omitempty"`
	GroupKeyGenerationIntervalOnForcedReKey           int           `json:"groupKeyGenerationIntervalOnForcedReKey,omitempty"`
	GroupKeyGenerationIntervalOnRevoke                int           `json:"groupKeyGenerationIntervalOnRevoke,omitempty"`
	GroupKeyMinimumSEKGenerationInterval              int           `json:"groupKeyMinimumSEKGenerationInterval,omitempty"`
	GroupKeyMinimumSEKLifetime                        int           `json:"groupKeyMinimumSEKLifetime,omitempty"`
	GroupKeyMinimumSeedGenerationInterval             int           `json:"groupKeyMinimumSeedGenerationInterval,omitempty"`
	GroupKeyMinimumSeedLifetime                       int           `json:"groupKeyMinimumSeedLifetime,omitempty"`
	GroupKeyMinimumTrafficEncryptionKeyLifetime       int           `json:"groupKeyMinimumTrafficEncryptionKeyLifetime,omitempty"`
	NsgBootstrapEndpoint                              string        `json:"nsgBootstrapEndpoint,omitempty"`
	NsgConfigEndpoint                                 string        `json:"nsgConfigEndpoint,omitempty"`
	NsgLocalUiUrl                                     string        `json:"nsgLocalUiUrl,omitempty"`
	EsiID                                             int           `json:"esiID,omitempty"`
	CsprootAuthenticationMethod                       string        `json:"csprootAuthenticationMethod,omitempty"`
	StackTraceEnabled                                 bool          `json:"stackTraceEnabled"`
	StatefulACLNonTCPTimeout                          int           `json:"statefulACLNonTCPTimeout,omitempty"`
	StatefulACLTCPTimeout                             int           `json:"statefulACLTCPTimeout,omitempty"`
	StaticWANServicePurgeTime                         int           `json:"staticWANServicePurgeTime,omitempty"`
	StatisticsEnabled                                 bool          `json:"statisticsEnabled"`
	StatsCollectorAddress                             string        `json:"statsCollectorAddress,omitempty"`
	StatsCollectorPort                                string        `json:"statsCollectorPort,omitempty"`
	StatsCollectorProtoBufPort                        string        `json:"statsCollectorProtoBufPort,omitempty"`
	StatsDatabaseProxy                                string        `json:"statsDatabaseProxy,omitempty"`
	StatsMaxDataPoints                                int           `json:"statsMaxDataPoints,omitempty"`
	StatsMinDuration                                  int           `json:"statsMinDuration,omitempty"`
	StatsNumberOfDataPoints                           int           `json:"statsNumberOfDataPoints,omitempty"`
	StatsTSDBServerAddress                            string        `json:"statsTSDBServerAddress,omitempty"`
	StickyECMPIdleTimeout                             int           `json:"stickyECMPIdleTimeout,omitempty"`
	AttachProbeToIPsecNPM                             bool          `json:"attachProbeToIPsecNPM"`
	AttachProbeToVXLANNPM                             bool          `json:"attachProbeToVXLANNPM"`
	SubnetResyncInterval                              int           `json:"subnetResyncInterval,omitempty"`
	SubnetResyncOutstandingInterval                   int           `json:"subnetResyncOutstandingInterval,omitempty"`
	CustomerIDUpperLimit                              int           `json:"customerIDUpperLimit,omitempty"`
	CustomerKey                                       string        `json:"customerKey,omitempty"`
	AvatarBasePath                                    string        `json:"avatarBasePath,omitempty"`
	AvatarBaseURL                                     string        `json:"avatarBaseURL,omitempty"`
	EventLogCleanupInterval                           int           `json:"eventLogCleanupInterval,omitempty"`
	EventLogEntryMaxAge                               int           `json:"eventLogEntryMaxAge,omitempty"`
	EventProcessorInterval                            int           `json:"eventProcessorInterval,omitempty"`
	EventProcessorMaxEventsCount                      int           `json:"eventProcessorMaxEventsCount,omitempty"`
	EventProcessorTimeout                             int           `json:"eventProcessorTimeout,omitempty"`
	TwoFactorCodeExpiry                               int           `json:"twoFactorCodeExpiry,omitempty"`
	TwoFactorCodeLength                               int           `json:"twoFactorCodeLength,omitempty"`
	TwoFactorCodeSeedLength                           int           `json:"twoFactorCodeSeedLength,omitempty"`
	ExternalID                                        string        `json:"externalID,omitempty"`
	DynamicWANServiceDiffTime                         int           `json:"dynamicWANServiceDiffTime,omitempty"`
	SyslogDestinationHost                             string        `json:"syslogDestinationHost,omitempty"`
	SyslogDestinationPort                             int           `json:"syslogDestinationPort,omitempty"`
	SysmonCleanupTaskInterval                         int           `json:"sysmonCleanupTaskInterval,omitempty"`
	SysmonNodePresenceTimeout                         int           `json:"sysmonNodePresenceTimeout,omitempty"`
	SysmonProbeResponseTimeout                        int           `json:"sysmonProbeResponseTimeout,omitempty"`
	SysmonPurgeInterval                               int           `json:"sysmonPurgeInterval,omitempty"`
	SystemAvatarData                                  string        `json:"systemAvatarData,omitempty"`
	SystemAvatarType                                  string        `json:"systemAvatarType,omitempty"`
}

// NewSystemConfig returns a new *SystemConfig
func NewSystemConfig() *SystemConfig {

	return &SystemConfig{
		AARFlowStatsInterval:                        30,
		AARProbeStatsInterval:                       30,
		ACLAllowOrigin:                              "*",
		ECMPCount:                                   1,
		LDAPSyncInterval:                            600,
		LDAPTrustStoreCertifcate:                    "/usr/lib/jvm/java-8-oracle/jre/lib/security/cacerts",
		LDAPTrustStorePassword:                      "changeit",
		ADGatewayPurgeTime:                          7200,
		RDLowerLimit:                                0,
		RDPublicNetworkLowerLimit:                   0,
		RDPublicNetworkUpperLimit:                   65535,
		RDUpperLimit:                                65535,
		ZFBRequestRetryTimer:                        30,
		ZFBSchedulerStaleRequestTimeout:             1440,
		PGIDLowerLimit:                              65536,
		PGIDUpperLimit:                              2147483647,
		DHCPOptionSize:                              16,
		VLANIDLowerLimit:                            0,
		VLANIDUpperLimit:                            0,
		VMCacheSize:                                 5000,
		VMPurgeTime:                                 60,
		VMResyncDeletionWaitTime:                    2,
		VMResyncOutstandingInterval:                 1000,
		VMUnreachableCleanupTime:                    7200,
		VMUnreachableTime:                           3600,
		VNFTaskTimeout:                              3600,
		VNIDLowerLimit:                              1,
		VNIDPublicNetworkLowerLimit:                 1,
		VNIDPublicNetworkUpperLimit:                 16777215,
		VNIDUpperLimit:                              1048575,
		APIKeyRenewalInterval:                       300,
		APIKeyValidity:                              86400,
		VPortInitStatefulTimer:                      300,
		LRUCacheSizePerSubnet:                       32,
		VSDReadOnlyMode:                             false,
		NSGUplinkHoldDownTimer:                      5,
		ASNumber:                                    65534,
		VSSStatsInterval:                            30,
		RTLowerLimit:                                0,
		RTPublicNetworkLowerLimit:                   0,
		RTPublicNetworkUpperLimit:                   65535,
		RTUpperLimit:                                65535,
		EVPNBGPCommunityTagASNumber:                 65534,
		EVPNBGPCommunityTagLowerLimit:               0,
		EVPNBGPCommunityTagUpperLimit:               65535,
		PageMaxSize:                                 500,
		PageSize:                                    50,
		LastExecutedMigrationPhase:                  "REDUCE",
		GatewayProbeInterval:                        5,
		GatewayProbeWindow:                          120,
		GatewayRebalancingInterval:                  600,
		MaxFailedLogins:                             0,
		MaxResponse:                                 500,
		AccumulateLicensesEnabled:                   false,
		SecondaryASNumber:                           65533,
		SecondaryRTLowerLimit:                       0,
		SecondaryRTUpperLimit:                       65533,
		PerDomainVlanIdEnabled:                      false,
		ServiceIDUpperLimit:                         2147483648,
		OffsetCustomerID:                            10000,
		OffsetServiceID:                             20001,
		VirtualFirewallRulesEnabled:                 false,
		EjbcaNSGCertificateProfile:                  "VSPClient",
		EjbcaNSGEndEntityProfile:                    "NSG",
		EjbcaOCSPResponderCN:                        "ocspsigner",
		EjbcaOCSPResponderURI:                       "http://localhost:7080/ejbca/publicweb/status/ocsp",
		EjbcaVspRootCa:                              "VSPCA",
		AlarmsMaxPerObject:                          10,
		ElasticClusterName:                          "nuage_elasticsearch",
		AllowEnterpriseAvatarOnNSG:                  true,
		EmbeddedMetadataSize:                        10,
		ImportedSaaSApplicationsVersion:             "1.0",
		InactiveTimeout:                             600000,
		InfrastructureBGPASNumber:                   65500,
		PostProcessorThreadsCount:                   20,
		GroupKeyDefaultSEKGenerationInterval:        1200,
		GroupKeyDefaultSEKLifetime:                  86400,
		GroupKeyDefaultSeedGenerationInterval:       1200,
		GroupKeyDefaultSeedLifetime:                 14400,
		GroupKeyDefaultTrafficEncryptionKeyLifetime: 600,
		GroupKeyGenerationIntervalOnForcedReKey:     60,
		GroupKeyGenerationIntervalOnRevoke:          60,
		GroupKeyMinimumSEKGenerationInterval:        20,
		GroupKeyMinimumSEKLifetime:                  20,
		GroupKeyMinimumSeedGenerationInterval:       20,
		GroupKeyMinimumSeedLifetime:                 30,
		GroupKeyMinimumTrafficEncryptionKeyLifetime: 30,
		NsgBootstrapEndpoint:                        "https://proxy-bootstrap:12443/nuage/api",
		NsgConfigEndpoint:                           "https://{proxyDNSName}:11443/nuage/api",
		NsgLocalUiUrl:                               "http://registration.nsg",
		EsiID:                                       10000,
		CsprootAuthenticationMethod:                 "LOCAL",
		StackTraceEnabled:                           false,
		StatefulACLNonTCPTimeout:                    180,
		StatefulACLTCPTimeout:                       3600,
		StaticWANServicePurgeTime:                   3600,
		StatsCollectorAddress:                       "localhost",
		StatsCollectorPort:                          "29090",
		StatsCollectorProtoBufPort:                  "39090",
		StatsMaxDataPoints:                          10000,
		StatsMinDuration:                            2592000,
		StatsNumberOfDataPoints:                     30,
		StatsTSDBServerAddress:                      "localhost:9300",
		StickyECMPIdleTimeout:                       0,
		AttachProbeToIPsecNPM:                       true,
		AttachProbeToVXLANNPM:                       true,
		SubnetResyncInterval:                        10,
		SubnetResyncOutstandingInterval:             20,
		CustomerIDUpperLimit:                        2147483647,
		AvatarBasePath:                              "/opt/vsd/jboss/standalone/deployments/CloudMgmt-web.war",
		AvatarBaseURL:                               "https://localhost:8443",
		EventLogCleanupInterval:                     3600,
		EventLogEntryMaxAge:                         7,
		EventProcessorInterval:                      250,
		EventProcessorMaxEventsCount:                100,
		EventProcessorTimeout:                       25000,
		TwoFactorCodeExpiry:                         300,
		TwoFactorCodeLength:                         6,
		TwoFactorCodeSeedLength:                     96,
		DynamicWANServiceDiffTime:                   1,
		SyslogDestinationHost:                       "http://localhost",
		SysmonCleanupTaskInterval:                   20,
		SysmonNodePresenceTimeout:                   3600,
		SysmonProbeResponseTimeout:                  30,
		SysmonPurgeInterval:                         86400,
	}
}

// Identity returns the Identity of the object.
func (o *SystemConfig) Identity() bambou.Identity {

	return SystemConfigIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *SystemConfig) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *SystemConfig) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the SystemConfig from the server
func (o *SystemConfig) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the SystemConfig into the server
func (o *SystemConfig) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the SystemConfig from the server
func (o *SystemConfig) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the SystemConfig
func (o *SystemConfig) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the SystemConfig
func (o *SystemConfig) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the SystemConfig
func (o *SystemConfig) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the SystemConfig
func (o *SystemConfig) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
