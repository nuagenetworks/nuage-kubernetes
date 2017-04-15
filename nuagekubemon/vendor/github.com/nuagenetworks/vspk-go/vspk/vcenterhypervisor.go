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

// VCenterHypervisorIdentity represents the Identity of the object
var VCenterHypervisorIdentity = bambou.Identity{
	Name:     "vcenterhypervisor",
	Category: "vcenterhypervisors",
}

// VCenterHypervisorsList represents a list of VCenterHypervisors
type VCenterHypervisorsList []*VCenterHypervisor

// VCenterHypervisorsAncestor is the interface that an ancestor of a VCenterHypervisor must implement.
// An Ancestor is defined as an entity that has VCenterHypervisor as a descendant.
// An Ancestor can get a list of its child VCenterHypervisors, but not necessarily create one.
type VCenterHypervisorsAncestor interface {
	VCenterHypervisors(*bambou.FetchingInfo) (VCenterHypervisorsList, *bambou.Error)
}

// VCenterHypervisorsParent is the interface that a parent of a VCenterHypervisor must implement.
// A Parent is defined as an entity that has VCenterHypervisor as a child.
// A Parent is an Ancestor which can create a VCenterHypervisor.
type VCenterHypervisorsParent interface {
	VCenterHypervisorsAncestor
	CreateVCenterHypervisor(*VCenterHypervisor) *bambou.Error
}

// VCenterHypervisor represents the model of a vcenterhypervisor
type VCenterHypervisor struct {
	ID                               string        `json:"ID,omitempty"`
	ParentID                         string        `json:"parentID,omitempty"`
	ParentType                       string        `json:"parentType,omitempty"`
	Owner                            string        `json:"owner,omitempty"`
	VCenterIP                        string        `json:"vCenterIP,omitempty"`
	VCenterPassword                  string        `json:"vCenterPassword,omitempty"`
	VCenterUser                      string        `json:"vCenterUser,omitempty"`
	VRSConfigurationTimeLimit        int           `json:"VRSConfigurationTimeLimit,omitempty"`
	VRSMetricsID                     string        `json:"VRSMetricsID,omitempty"`
	VRSState                         string        `json:"VRSState,omitempty"`
	VRequireNuageMetadata            bool          `json:"vRequireNuageMetadata"`
	Name                             string        `json:"name,omitempty"`
	ManagedObjectID                  string        `json:"managedObjectID,omitempty"`
	LastUpdatedBy                    string        `json:"lastUpdatedBy,omitempty"`
	LastVRSDeployedDate              float64       `json:"lastVRSDeployedDate,omitempty"`
	DataDNS1                         string        `json:"dataDNS1,omitempty"`
	DataDNS2                         string        `json:"dataDNS2,omitempty"`
	DataGateway                      string        `json:"dataGateway,omitempty"`
	DataIPAddress                    string        `json:"dataIPAddress,omitempty"`
	DataNetmask                      string        `json:"dataNetmask,omitempty"`
	DataNetworkPortgroup             string        `json:"dataNetworkPortgroup,omitempty"`
	DatapathSyncTimeout              int           `json:"datapathSyncTimeout,omitempty"`
	Scope                            bool          `json:"scope"`
	SecondaryNuageController         string        `json:"secondaryNuageController,omitempty"`
	RemovedFromVCenterInventory      bool          `json:"removedFromVCenterInventory"`
	GenericSplitActivation           bool          `json:"genericSplitActivation"`
	SeparateDataNetwork              bool          `json:"separateDataNetwork"`
	DeploymentCount                  int           `json:"deploymentCount,omitempty"`
	Personality                      string        `json:"personality,omitempty"`
	Description                      string        `json:"description,omitempty"`
	DestinationMirrorPort            string        `json:"destinationMirrorPort,omitempty"`
	MetadataServerIP                 string        `json:"metadataServerIP,omitempty"`
	MetadataServerListenPort         int           `json:"metadataServerListenPort,omitempty"`
	MetadataServerPort               int           `json:"metadataServerPort,omitempty"`
	MetadataServiceEnabled           bool          `json:"metadataServiceEnabled"`
	NetworkUplinkInterface           string        `json:"networkUplinkInterface,omitempty"`
	NetworkUplinkInterfaceGateway    string        `json:"networkUplinkInterfaceGateway,omitempty"`
	NetworkUplinkInterfaceIp         string        `json:"networkUplinkInterfaceIp,omitempty"`
	NetworkUplinkInterfaceNetmask    string        `json:"networkUplinkInterfaceNetmask,omitempty"`
	NfsLogServer                     string        `json:"nfsLogServer,omitempty"`
	NfsMountPath                     string        `json:"nfsMountPath,omitempty"`
	MgmtDNS1                         string        `json:"mgmtDNS1,omitempty"`
	MgmtDNS2                         string        `json:"mgmtDNS2,omitempty"`
	MgmtGateway                      string        `json:"mgmtGateway,omitempty"`
	MgmtIPAddress                    string        `json:"mgmtIPAddress,omitempty"`
	MgmtNetmask                      string        `json:"mgmtNetmask,omitempty"`
	MgmtNetworkPortgroup             string        `json:"mgmtNetworkPortgroup,omitempty"`
	DhcpRelayServer                  string        `json:"dhcpRelayServer,omitempty"`
	MirrorNetworkPortgroup           string        `json:"mirrorNetworkPortgroup,omitempty"`
	SiteId                           string        `json:"siteId,omitempty"`
	AllowDataDHCP                    bool          `json:"allowDataDHCP"`
	AllowMgmtDHCP                    bool          `json:"allowMgmtDHCP"`
	FlowEvictionThreshold            int           `json:"flowEvictionThreshold,omitempty"`
	VmNetworkPortgroup               string        `json:"vmNetworkPortgroup,omitempty"`
	EntityScope                      string        `json:"entityScope,omitempty"`
	ToolboxDeploymentMode            bool          `json:"toolboxDeploymentMode"`
	ToolboxGroup                     string        `json:"toolboxGroup,omitempty"`
	ToolboxIP                        string        `json:"toolboxIP,omitempty"`
	ToolboxPassword                  string        `json:"toolboxPassword,omitempty"`
	ToolboxUserName                  string        `json:"toolboxUserName,omitempty"`
	PortgroupMetadata                bool          `json:"portgroupMetadata"`
	NovaClientVersion                int           `json:"novaClientVersion,omitempty"`
	NovaMetadataServiceAuthUrl       string        `json:"novaMetadataServiceAuthUrl,omitempty"`
	NovaMetadataServiceEndpoint      string        `json:"novaMetadataServiceEndpoint,omitempty"`
	NovaMetadataServicePassword      string        `json:"novaMetadataServicePassword,omitempty"`
	NovaMetadataServiceTenant        string        `json:"novaMetadataServiceTenant,omitempty"`
	NovaMetadataServiceUsername      string        `json:"novaMetadataServiceUsername,omitempty"`
	NovaMetadataSharedSecret         string        `json:"novaMetadataSharedSecret,omitempty"`
	NovaRegionName                   string        `json:"novaRegionName,omitempty"`
	UpgradePackagePassword           string        `json:"upgradePackagePassword,omitempty"`
	UpgradePackageURL                string        `json:"upgradePackageURL,omitempty"`
	UpgradePackageUsername           string        `json:"upgradePackageUsername,omitempty"`
	UpgradeScriptTimeLimit           int           `json:"upgradeScriptTimeLimit,omitempty"`
	UpgradeStatus                    string        `json:"upgradeStatus,omitempty"`
	UpgradeTimedout                  bool          `json:"upgradeTimedout"`
	PrimaryNuageController           string        `json:"primaryNuageController,omitempty"`
	VrsId                            string        `json:"vrsId,omitempty"`
	VrsPassword                      string        `json:"vrsPassword,omitempty"`
	VrsUserName                      string        `json:"vrsUserName,omitempty"`
	StaticRoute                      string        `json:"staticRoute,omitempty"`
	StaticRouteGateway               string        `json:"staticRouteGateway,omitempty"`
	StaticRouteNetmask               string        `json:"staticRouteNetmask,omitempty"`
	NtpServer1                       string        `json:"ntpServer1,omitempty"`
	NtpServer2                       string        `json:"ntpServer2,omitempty"`
	Mtu                              int           `json:"mtu,omitempty"`
	MultiVMSsupport                  bool          `json:"multiVMSsupport"`
	MulticastReceiveInterface        string        `json:"multicastReceiveInterface,omitempty"`
	MulticastReceiveInterfaceIP      string        `json:"multicastReceiveInterfaceIP,omitempty"`
	MulticastReceiveInterfaceNetmask string        `json:"multicastReceiveInterfaceNetmask,omitempty"`
	MulticastReceiveRange            string        `json:"multicastReceiveRange,omitempty"`
	MulticastSendInterface           string        `json:"multicastSendInterface,omitempty"`
	MulticastSendInterfaceIP         string        `json:"multicastSendInterfaceIP,omitempty"`
	MulticastSendInterfaceNetmask    string        `json:"multicastSendInterfaceNetmask,omitempty"`
	MulticastSourcePortgroup         string        `json:"multicastSourcePortgroup,omitempty"`
	CustomizedScriptURL              string        `json:"customizedScriptURL,omitempty"`
	AvailableNetworks                []interface{} `json:"availableNetworks,omitempty"`
	OvfURL                           string        `json:"ovfURL,omitempty"`
	ExternalID                       string        `json:"externalID,omitempty"`
	HypervisorIP                     string        `json:"hypervisorIP,omitempty"`
	HypervisorPassword               string        `json:"hypervisorPassword,omitempty"`
	HypervisorUser                   string        `json:"hypervisorUser,omitempty"`
}

// NewVCenterHypervisor returns a new *VCenterHypervisor
func NewVCenterHypervisor() *VCenterHypervisor {

	return &VCenterHypervisor{
		VRSState:              "NOT_DEPLOYED",
		DestinationMirrorPort: "no_mirror",
	}
}

// Identity returns the Identity of the object.
func (o *VCenterHypervisor) Identity() bambou.Identity {

	return VCenterHypervisorIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VCenterHypervisor) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VCenterHypervisor) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VCenterHypervisor from the server
func (o *VCenterHypervisor) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VCenterHypervisor into the server
func (o *VCenterHypervisor) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VCenterHypervisor from the server
func (o *VCenterHypervisor) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the VCenterHypervisor
func (o *VCenterHypervisor) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VCenterHypervisor
func (o *VCenterHypervisor) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VCenterHypervisor
func (o *VCenterHypervisor) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VCenterHypervisor
func (o *VCenterHypervisor) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CreateJob creates a new child Job under the VCenterHypervisor
func (o *VCenterHypervisor) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VRSAddressRanges retrieves the list of child VRSAddressRanges of the VCenterHypervisor
func (o *VCenterHypervisor) VRSAddressRanges(info *bambou.FetchingInfo) (VRSAddressRangesList, *bambou.Error) {

	var list VRSAddressRangesList
	err := bambou.CurrentSession().FetchChildren(o, VRSAddressRangeIdentity, &list, info)
	return list, err
}

// CreateVRSAddressRange creates a new child VRSAddressRange under the VCenterHypervisor
func (o *VCenterHypervisor) CreateVRSAddressRange(child *VRSAddressRange) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VRSMetrics retrieves the list of child VRSMetrics of the VCenterHypervisor
func (o *VCenterHypervisor) VRSMetrics(info *bambou.FetchingInfo) (VRSMetricsList, *bambou.Error) {

	var list VRSMetricsList
	err := bambou.CurrentSession().FetchChildren(o, VRSMetricsIdentity, &list, info)
	return list, err
}

// VRSRedeploymentpolicies retrieves the list of child VRSRedeploymentpolicies of the VCenterHypervisor
func (o *VCenterHypervisor) VRSRedeploymentpolicies(info *bambou.FetchingInfo) (VRSRedeploymentpoliciesList, *bambou.Error) {

	var list VRSRedeploymentpoliciesList
	err := bambou.CurrentSession().FetchChildren(o, VRSRedeploymentpolicyIdentity, &list, info)
	return list, err
}

// CreateVRSRedeploymentpolicy creates a new child VRSRedeploymentpolicy under the VCenterHypervisor
func (o *VCenterHypervisor) CreateVRSRedeploymentpolicy(child *VRSRedeploymentpolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
