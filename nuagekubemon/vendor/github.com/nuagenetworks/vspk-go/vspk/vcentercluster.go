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

// VCenterClusterIdentity represents the Identity of the object
var VCenterClusterIdentity = bambou.Identity{
	Name:     "vcentercluster",
	Category: "vcenterclusters",
}

// VCenterClustersList represents a list of VCenterClusters
type VCenterClustersList []*VCenterCluster

// VCenterClustersAncestor is the interface that an ancestor of a VCenterCluster must implement.
// An Ancestor is defined as an entity that has VCenterCluster as a descendant.
// An Ancestor can get a list of its child VCenterClusters, but not necessarily create one.
type VCenterClustersAncestor interface {
	VCenterClusters(*bambou.FetchingInfo) (VCenterClustersList, *bambou.Error)
}

// VCenterClustersParent is the interface that a parent of a VCenterCluster must implement.
// A Parent is defined as an entity that has VCenterCluster as a child.
// A Parent is an Ancestor which can create a VCenterCluster.
type VCenterClustersParent interface {
	VCenterClustersAncestor
	CreateVCenterCluster(*VCenterCluster) *bambou.Error
}

// VCenterCluster represents the model of a vcentercluster
type VCenterCluster struct {
	ID                               string `json:"ID,omitempty"`
	ParentID                         string `json:"parentID,omitempty"`
	ParentType                       string `json:"parentType,omitempty"`
	Owner                            string `json:"owner,omitempty"`
	VRSConfigurationTimeLimit        int    `json:"VRSConfigurationTimeLimit,omitempty"`
	VRequireNuageMetadata            bool   `json:"vRequireNuageMetadata"`
	Name                             string `json:"name,omitempty"`
	ManagedObjectID                  string `json:"managedObjectID,omitempty"`
	LastUpdatedBy                    string `json:"lastUpdatedBy,omitempty"`
	DataDNS1                         string `json:"dataDNS1,omitempty"`
	DataDNS2                         string `json:"dataDNS2,omitempty"`
	DataGateway                      string `json:"dataGateway,omitempty"`
	DataNetworkPortgroup             string `json:"dataNetworkPortgroup,omitempty"`
	DatapathSyncTimeout              int    `json:"datapathSyncTimeout,omitempty"`
	Scope                            bool   `json:"scope"`
	SecondaryNuageController         string `json:"secondaryNuageController,omitempty"`
	DeletedFromVCenterDataCenter     bool   `json:"deletedFromVCenterDataCenter"`
	GenericSplitActivation           bool   `json:"genericSplitActivation"`
	SeparateDataNetwork              bool   `json:"separateDataNetwork"`
	Personality                      string `json:"personality,omitempty"`
	Description                      string `json:"description,omitempty"`
	DestinationMirrorPort            string `json:"destinationMirrorPort,omitempty"`
	MetadataServerIP                 string `json:"metadataServerIP,omitempty"`
	MetadataServerListenPort         int    `json:"metadataServerListenPort,omitempty"`
	MetadataServerPort               int    `json:"metadataServerPort,omitempty"`
	MetadataServiceEnabled           bool   `json:"metadataServiceEnabled"`
	NetworkUplinkInterface           string `json:"networkUplinkInterface,omitempty"`
	NetworkUplinkInterfaceGateway    string `json:"networkUplinkInterfaceGateway,omitempty"`
	NetworkUplinkInterfaceIp         string `json:"networkUplinkInterfaceIp,omitempty"`
	NetworkUplinkInterfaceNetmask    string `json:"networkUplinkInterfaceNetmask,omitempty"`
	NfsLogServer                     string `json:"nfsLogServer,omitempty"`
	NfsMountPath                     string `json:"nfsMountPath,omitempty"`
	MgmtDNS1                         string `json:"mgmtDNS1,omitempty"`
	MgmtDNS2                         string `json:"mgmtDNS2,omitempty"`
	MgmtGateway                      string `json:"mgmtGateway,omitempty"`
	MgmtNetworkPortgroup             string `json:"mgmtNetworkPortgroup,omitempty"`
	DhcpRelayServer                  string `json:"dhcpRelayServer,omitempty"`
	MirrorNetworkPortgroup           string `json:"mirrorNetworkPortgroup,omitempty"`
	SiteId                           string `json:"siteId,omitempty"`
	AllowDataDHCP                    bool   `json:"allowDataDHCP"`
	AllowMgmtDHCP                    bool   `json:"allowMgmtDHCP"`
	FlowEvictionThreshold            int    `json:"flowEvictionThreshold,omitempty"`
	VmNetworkPortgroup               string `json:"vmNetworkPortgroup,omitempty"`
	EntityScope                      string `json:"entityScope,omitempty"`
	PortgroupMetadata                bool   `json:"portgroupMetadata"`
	NovaClientVersion                int    `json:"novaClientVersion,omitempty"`
	NovaMetadataServiceAuthUrl       string `json:"novaMetadataServiceAuthUrl,omitempty"`
	NovaMetadataServiceEndpoint      string `json:"novaMetadataServiceEndpoint,omitempty"`
	NovaMetadataServicePassword      string `json:"novaMetadataServicePassword,omitempty"`
	NovaMetadataServiceTenant        string `json:"novaMetadataServiceTenant,omitempty"`
	NovaMetadataServiceUsername      string `json:"novaMetadataServiceUsername,omitempty"`
	NovaMetadataSharedSecret         string `json:"novaMetadataSharedSecret,omitempty"`
	NovaRegionName                   string `json:"novaRegionName,omitempty"`
	UpgradePackagePassword           string `json:"upgradePackagePassword,omitempty"`
	UpgradePackageURL                string `json:"upgradePackageURL,omitempty"`
	UpgradePackageUsername           string `json:"upgradePackageUsername,omitempty"`
	UpgradeScriptTimeLimit           int    `json:"upgradeScriptTimeLimit,omitempty"`
	PrimaryNuageController           string `json:"primaryNuageController,omitempty"`
	VrsPassword                      string `json:"vrsPassword,omitempty"`
	VrsUserName                      string `json:"vrsUserName,omitempty"`
	AssocVCenterDataCenterID         string `json:"assocVCenterDataCenterID,omitempty"`
	AssocVCenterID                   string `json:"assocVCenterID,omitempty"`
	StaticRoute                      string `json:"staticRoute,omitempty"`
	StaticRouteGateway               string `json:"staticRouteGateway,omitempty"`
	StaticRouteNetmask               string `json:"staticRouteNetmask,omitempty"`
	NtpServer1                       string `json:"ntpServer1,omitempty"`
	NtpServer2                       string `json:"ntpServer2,omitempty"`
	Mtu                              int    `json:"mtu,omitempty"`
	MultiVMSsupport                  bool   `json:"multiVMSsupport"`
	MulticastReceiveInterface        string `json:"multicastReceiveInterface,omitempty"`
	MulticastReceiveInterfaceIP      string `json:"multicastReceiveInterfaceIP,omitempty"`
	MulticastReceiveInterfaceNetmask string `json:"multicastReceiveInterfaceNetmask,omitempty"`
	MulticastReceiveRange            string `json:"multicastReceiveRange,omitempty"`
	MulticastSendInterface           string `json:"multicastSendInterface,omitempty"`
	MulticastSendInterfaceIP         string `json:"multicastSendInterfaceIP,omitempty"`
	MulticastSendInterfaceNetmask    string `json:"multicastSendInterfaceNetmask,omitempty"`
	MulticastSourcePortgroup         string `json:"multicastSourcePortgroup,omitempty"`
	CustomizedScriptURL              string `json:"customizedScriptURL,omitempty"`
	OvfURL                           string `json:"ovfURL,omitempty"`
	ExternalID                       string `json:"externalID,omitempty"`
}

// NewVCenterCluster returns a new *VCenterCluster
func NewVCenterCluster() *VCenterCluster {

	return &VCenterCluster{
		DestinationMirrorPort: "no_mirror",
	}
}

// Identity returns the Identity of the object.
func (o *VCenterCluster) Identity() bambou.Identity {

	return VCenterClusterIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VCenterCluster) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VCenterCluster) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VCenterCluster from the server
func (o *VCenterCluster) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VCenterCluster into the server
func (o *VCenterCluster) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VCenterCluster from the server
func (o *VCenterCluster) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// VCenterHypervisors retrieves the list of child VCenterHypervisors of the VCenterCluster
func (o *VCenterCluster) VCenterHypervisors(info *bambou.FetchingInfo) (VCenterHypervisorsList, *bambou.Error) {

	var list VCenterHypervisorsList
	err := bambou.CurrentSession().FetchChildren(o, VCenterHypervisorIdentity, &list, info)
	return list, err
}

// CreateVCenterHypervisor creates a new child VCenterHypervisor under the VCenterCluster
func (o *VCenterCluster) CreateVCenterHypervisor(child *VCenterHypervisor) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the VCenterCluster
func (o *VCenterCluster) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VCenterCluster
func (o *VCenterCluster) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VCenterCluster
func (o *VCenterCluster) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VCenterCluster
func (o *VCenterCluster) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// CreateJob creates a new child Job under the VCenterCluster
func (o *VCenterCluster) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VRSAddressRanges retrieves the list of child VRSAddressRanges of the VCenterCluster
func (o *VCenterCluster) VRSAddressRanges(info *bambou.FetchingInfo) (VRSAddressRangesList, *bambou.Error) {

	var list VRSAddressRangesList
	err := bambou.CurrentSession().FetchChildren(o, VRSAddressRangeIdentity, &list, info)
	return list, err
}

// CreateVRSAddressRange creates a new child VRSAddressRange under the VCenterCluster
func (o *VCenterCluster) CreateVRSAddressRange(child *VRSAddressRange) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VRSRedeploymentpolicies retrieves the list of child VRSRedeploymentpolicies of the VCenterCluster
func (o *VCenterCluster) VRSRedeploymentpolicies(info *bambou.FetchingInfo) (VRSRedeploymentpoliciesList, *bambou.Error) {

	var list VRSRedeploymentpoliciesList
	err := bambou.CurrentSession().FetchChildren(o, VRSRedeploymentpolicyIdentity, &list, info)
	return list, err
}

// CreateVRSRedeploymentpolicy creates a new child VRSRedeploymentpolicy under the VCenterCluster
func (o *VCenterCluster) CreateVRSRedeploymentpolicy(child *VRSRedeploymentpolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// AutoDiscoverHypervisorFromClusters retrieves the list of child AutoDiscoverHypervisorFromClusters of the VCenterCluster
func (o *VCenterCluster) AutoDiscoverHypervisorFromClusters(info *bambou.FetchingInfo) (AutoDiscoverHypervisorFromClustersList, *bambou.Error) {

	var list AutoDiscoverHypervisorFromClustersList
	err := bambou.CurrentSession().FetchChildren(o, AutoDiscoverHypervisorFromClusterIdentity, &list, info)
	return list, err
}
