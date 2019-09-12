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

// VCenterIdentity represents the Identity of the object
var VCenterIdentity = bambou.Identity{
	Name:     "vcenter",
	Category: "vcenters",
}

// VCentersList represents a list of VCenters
type VCentersList []*VCenter

// VCentersAncestor is the interface that an ancestor of a VCenter must implement.
// An Ancestor is defined as an entity that has VCenter as a descendant.
// An Ancestor can get a list of its child VCenters, but not necessarily create one.
type VCentersAncestor interface {
	VCenters(*bambou.FetchingInfo) (VCentersList, *bambou.Error)
}

// VCentersParent is the interface that a parent of a VCenter must implement.
// A Parent is defined as an entity that has VCenter as a child.
// A Parent is an Ancestor which can create a VCenter.
type VCentersParent interface {
	VCentersAncestor
	CreateVCenter(*VCenter) *bambou.Error
}

// VCenter represents the model of a vcenter
type VCenter struct {
	ID                                     string        `json:"ID,omitempty"`
	ParentID                               string        `json:"parentID,omitempty"`
	ParentType                             string        `json:"parentType,omitempty"`
	Owner                                  string        `json:"owner,omitempty"`
	EAMExtensionName                       string        `json:"EAMExtensionName,omitempty"`
	ARPReply                               bool          `json:"ARPReply"`
	VRSConfigurationTimeLimit              int           `json:"VRSConfigurationTimeLimit,omitempty"`
	VRequireNuageMetadata                  bool          `json:"vRequireNuageMetadata"`
	Name                                   string        `json:"name,omitempty"`
	ManageVRSAvailability                  bool          `json:"manageVRSAvailability"`
	Password                               string        `json:"password,omitempty"`
	LastUpdatedBy                          string        `json:"lastUpdatedBy,omitempty"`
	DataDNS1                               string        `json:"dataDNS1,omitempty"`
	DataDNS2                               string        `json:"dataDNS2,omitempty"`
	DataGateway                            string        `json:"dataGateway,omitempty"`
	DataNetworkPortgroup                   string        `json:"dataNetworkPortgroup,omitempty"`
	DatapathSyncTimeout                    int           `json:"datapathSyncTimeout,omitempty"`
	SecondaryDataUplinkDHCPEnabled         bool          `json:"secondaryDataUplinkDHCPEnabled"`
	SecondaryDataUplinkEnabled             bool          `json:"secondaryDataUplinkEnabled"`
	SecondaryDataUplinkInterface           string        `json:"secondaryDataUplinkInterface,omitempty"`
	SecondaryDataUplinkMTU                 int           `json:"secondaryDataUplinkMTU,omitempty"`
	SecondaryDataUplinkPrimaryController   string        `json:"secondaryDataUplinkPrimaryController,omitempty"`
	SecondaryDataUplinkSecondaryController string        `json:"secondaryDataUplinkSecondaryController,omitempty"`
	SecondaryDataUplinkUnderlayID          int           `json:"secondaryDataUplinkUnderlayID,omitempty"`
	SecondaryDataUplinkVDFControlVLAN      int           `json:"secondaryDataUplinkVDFControlVLAN,omitempty"`
	SecondaryNuageController               string        `json:"secondaryNuageController,omitempty"`
	MemorySizeInGB                         string        `json:"memorySizeInGB,omitempty"`
	RemoteSyslogServerIP                   string        `json:"remoteSyslogServerIP,omitempty"`
	RemoteSyslogServerPort                 int           `json:"remoteSyslogServerPort,omitempty"`
	RemoteSyslogServerType                 string        `json:"remoteSyslogServerType,omitempty"`
	GenericSplitActivation                 bool          `json:"genericSplitActivation"`
	SeparateDataNetwork                    bool          `json:"separateDataNetwork"`
	Personality                            string        `json:"personality,omitempty"`
	Description                            string        `json:"description,omitempty"`
	DestinationMirrorPort                  string        `json:"destinationMirrorPort,omitempty"`
	MetadataServerIP                       string        `json:"metadataServerIP,omitempty"`
	MetadataServerListenPort               int           `json:"metadataServerListenPort,omitempty"`
	MetadataServerPort                     int           `json:"metadataServerPort,omitempty"`
	MetadataServiceEnabled                 bool          `json:"metadataServiceEnabled"`
	NetworkUplinkInterface                 string        `json:"networkUplinkInterface,omitempty"`
	NetworkUplinkInterfaceGateway          string        `json:"networkUplinkInterfaceGateway,omitempty"`
	NetworkUplinkInterfaceIp               string        `json:"networkUplinkInterfaceIp,omitempty"`
	NetworkUplinkInterfaceNetmask          string        `json:"networkUplinkInterfaceNetmask,omitempty"`
	RevertiveControllerEnabled             bool          `json:"revertiveControllerEnabled"`
	RevertiveTimer                         int           `json:"revertiveTimer,omitempty"`
	NfsLogServer                           string        `json:"nfsLogServer,omitempty"`
	NfsMountPath                           string        `json:"nfsMountPath,omitempty"`
	MgmtDNS1                               string        `json:"mgmtDNS1,omitempty"`
	MgmtDNS2                               string        `json:"mgmtDNS2,omitempty"`
	MgmtGateway                            string        `json:"mgmtGateway,omitempty"`
	MgmtNetworkPortgroup                   string        `json:"mgmtNetworkPortgroup,omitempty"`
	DhcpRelayServer                        string        `json:"dhcpRelayServer,omitempty"`
	MirrorNetworkPortgroup                 string        `json:"mirrorNetworkPortgroup,omitempty"`
	DisableGROOnDatapath                   bool          `json:"disableGROOnDatapath"`
	DisableLROOnDatapath                   bool          `json:"disableLROOnDatapath"`
	DisableNetworkDiscovery                bool          `json:"disableNetworkDiscovery"`
	SiteId                                 string        `json:"siteId,omitempty"`
	OldAgencyName                          string        `json:"oldAgencyName,omitempty"`
	AllowDataDHCP                          bool          `json:"allowDataDHCP"`
	AllowMgmtDHCP                          bool          `json:"allowMgmtDHCP"`
	FlowEvictionThreshold                  int           `json:"flowEvictionThreshold,omitempty"`
	VmNetworkPortgroup                     string        `json:"vmNetworkPortgroup,omitempty"`
	EmbeddedMetadata                       []interface{} `json:"embeddedMetadata,omitempty"`
	EnableVRSResourceReservation           bool          `json:"enableVRSResourceReservation"`
	EntityScope                            string        `json:"entityScope,omitempty"`
	ConfiguredMetricsPushInterval          int           `json:"configuredMetricsPushInterval,omitempty"`
	ConnectionStatus                       bool          `json:"connectionStatus"`
	PortgroupMetadata                      bool          `json:"portgroupMetadata"`
	HostLevelManagement                    bool          `json:"hostLevelManagement"`
	NovaClientVersion                      int           `json:"novaClientVersion,omitempty"`
	NovaIdentityURLVersion                 string        `json:"novaIdentityURLVersion,omitempty"`
	NovaMetadataServiceAuthUrl             string        `json:"novaMetadataServiceAuthUrl,omitempty"`
	NovaMetadataServiceEndpoint            string        `json:"novaMetadataServiceEndpoint,omitempty"`
	NovaMetadataServicePassword            string        `json:"novaMetadataServicePassword,omitempty"`
	NovaMetadataServiceTenant              string        `json:"novaMetadataServiceTenant,omitempty"`
	NovaMetadataServiceUsername            string        `json:"novaMetadataServiceUsername,omitempty"`
	NovaMetadataSharedSecret               string        `json:"novaMetadataSharedSecret,omitempty"`
	NovaOSKeystoneUsername                 string        `json:"novaOSKeystoneUsername,omitempty"`
	NovaProjectDomainName                  string        `json:"novaProjectDomainName,omitempty"`
	NovaProjectName                        string        `json:"novaProjectName,omitempty"`
	NovaRegionName                         string        `json:"novaRegionName,omitempty"`
	NovaUserDomainName                     string        `json:"novaUserDomainName,omitempty"`
	IpAddress                              string        `json:"ipAddress,omitempty"`
	UpgradePackagePassword                 string        `json:"upgradePackagePassword,omitempty"`
	UpgradePackageURL                      string        `json:"upgradePackageURL,omitempty"`
	UpgradePackageUsername                 string        `json:"upgradePackageUsername,omitempty"`
	UpgradeScriptTimeLimit                 int           `json:"upgradeScriptTimeLimit,omitempty"`
	CpuCount                               string        `json:"cpuCount,omitempty"`
	PrimaryDataUplinkUnderlayID            int           `json:"primaryDataUplinkUnderlayID,omitempty"`
	PrimaryDataUplinkVDFControlVLAN        int           `json:"primaryDataUplinkVDFControlVLAN,omitempty"`
	PrimaryNuageController                 string        `json:"primaryNuageController,omitempty"`
	VrsConfigID                            string        `json:"vrsConfigID,omitempty"`
	VrsPassword                            string        `json:"vrsPassword,omitempty"`
	VrsUserName                            string        `json:"vrsUserName,omitempty"`
	UserName                               string        `json:"userName,omitempty"`
	StaticRoute                            string        `json:"staticRoute,omitempty"`
	StaticRouteGateway                     string        `json:"staticRouteGateway,omitempty"`
	StaticRouteNetmask                     string        `json:"staticRouteNetmask,omitempty"`
	NtpServer1                             string        `json:"ntpServer1,omitempty"`
	NtpServer2                             string        `json:"ntpServer2,omitempty"`
	HttpPort                               int           `json:"httpPort,omitempty"`
	HttpsPort                              int           `json:"httpsPort,omitempty"`
	Mtu                                    int           `json:"mtu,omitempty"`
	MultiVMSsupport                        bool          `json:"multiVMSsupport"`
	MulticastReceiveInterface              string        `json:"multicastReceiveInterface,omitempty"`
	MulticastReceiveInterfaceIP            string        `json:"multicastReceiveInterfaceIP,omitempty"`
	MulticastReceiveInterfaceNetmask       string        `json:"multicastReceiveInterfaceNetmask,omitempty"`
	MulticastReceiveRange                  string        `json:"multicastReceiveRange,omitempty"`
	MulticastSendInterface                 string        `json:"multicastSendInterface,omitempty"`
	MulticastSendInterfaceIP               string        `json:"multicastSendInterfaceIP,omitempty"`
	MulticastSendInterfaceNetmask          string        `json:"multicastSendInterfaceNetmask,omitempty"`
	MulticastSourcePortgroup               string        `json:"multicastSourcePortgroup,omitempty"`
	CustomizedScriptURL                    string        `json:"customizedScriptURL,omitempty"`
	AutoResolveFrequency                   int           `json:"autoResolveFrequency,omitempty"`
	OvfURL                                 string        `json:"ovfURL,omitempty"`
	AvrsEnabled                            bool          `json:"avrsEnabled"`
	AvrsProfile                            string        `json:"avrsProfile,omitempty"`
	ExternalID                             string        `json:"externalID,omitempty"`
}

// NewVCenter returns a new *VCenter
func NewVCenter() *VCenter {

	return &VCenter{
		ManageVRSAvailability:             false,
		SecondaryDataUplinkDHCPEnabled:    false,
		SecondaryDataUplinkEnabled:        false,
		SecondaryDataUplinkMTU:            1500,
		SecondaryDataUplinkUnderlayID:     1,
		SecondaryDataUplinkVDFControlVLAN: 0,
		MemorySizeInGB:                    "DEFAULT_4",
		RemoteSyslogServerPort:            514,
		RemoteSyslogServerType:            "NONE",
		Personality:                       "VRS",
		DestinationMirrorPort:             "no_mirror",
		RevertiveControllerEnabled:        false,
		RevertiveTimer:                    300,
		DisableGROOnDatapath:              false,
		DisableLROOnDatapath:              false,
		EnableVRSResourceReservation:      false,
		ConfiguredMetricsPushInterval:     60,
		CpuCount:                          "DEFAULT_2",
		PrimaryDataUplinkUnderlayID:       0,
		PrimaryDataUplinkVDFControlVLAN:   0,
		AvrsEnabled:                       false,
		AvrsProfile:                       "AVRS_25G",
	}
}

// Identity returns the Identity of the object.
func (o *VCenter) Identity() bambou.Identity {

	return VCenterIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VCenter) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VCenter) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VCenter from the server
func (o *VCenter) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VCenter into the server
func (o *VCenter) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VCenter from the server
func (o *VCenter) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// VCenterDataCenters retrieves the list of child VCenterDataCenters of the VCenter
func (o *VCenter) VCenterDataCenters(info *bambou.FetchingInfo) (VCenterDataCentersList, *bambou.Error) {

	var list VCenterDataCentersList
	err := bambou.CurrentSession().FetchChildren(o, VCenterDataCenterIdentity, &list, info)
	return list, err
}

// CreateVCenterDataCenter creates a new child VCenterDataCenter under the VCenter
func (o *VCenter) CreateVCenterDataCenter(child *VCenterDataCenter) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the VCenter
func (o *VCenter) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VCenter
func (o *VCenter) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VCenter
func (o *VCenter) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VCenter
func (o *VCenter) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CreateJob creates a new child Job under the VCenter
func (o *VCenter) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VRSAddressRanges retrieves the list of child VRSAddressRanges of the VCenter
func (o *VCenter) VRSAddressRanges(info *bambou.FetchingInfo) (VRSAddressRangesList, *bambou.Error) {

	var list VRSAddressRangesList
	err := bambou.CurrentSession().FetchChildren(o, VRSAddressRangeIdentity, &list, info)
	return list, err
}

// CreateVRSAddressRange creates a new child VRSAddressRange under the VCenter
func (o *VCenter) CreateVRSAddressRange(child *VRSAddressRange) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VRSRedeploymentpolicies retrieves the list of child VRSRedeploymentpolicies of the VCenter
func (o *VCenter) VRSRedeploymentpolicies(info *bambou.FetchingInfo) (VRSRedeploymentpoliciesList, *bambou.Error) {

	var list VRSRedeploymentpoliciesList
	err := bambou.CurrentSession().FetchChildren(o, VRSRedeploymentpolicyIdentity, &list, info)
	return list, err
}

// CreateVRSRedeploymentpolicy creates a new child VRSRedeploymentpolicy under the VCenter
func (o *VCenter) CreateVRSRedeploymentpolicy(child *VRSRedeploymentpolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Autodiscovereddatacenters retrieves the list of child Autodiscovereddatacenters of the VCenter
func (o *VCenter) Autodiscovereddatacenters(info *bambou.FetchingInfo) (AutodiscovereddatacentersList, *bambou.Error) {

	var list AutodiscovereddatacentersList
	err := bambou.CurrentSession().FetchChildren(o, AutodiscovereddatacenterIdentity, &list, info)
	return list, err
}
