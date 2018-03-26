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

// VCenterVRSConfigIdentity represents the Identity of the object
var VCenterVRSConfigIdentity = bambou.Identity{
	Name:     "vrsconfig",
	Category: "vrsconfigs",
}

// VCenterVRSConfigsList represents a list of VCenterVRSConfigs
type VCenterVRSConfigsList []*VCenterVRSConfig

// VCenterVRSConfigsAncestor is the interface that an ancestor of a VCenterVRSConfig must implement.
// An Ancestor is defined as an entity that has VCenterVRSConfig as a descendant.
// An Ancestor can get a list of its child VCenterVRSConfigs, but not necessarily create one.
type VCenterVRSConfigsAncestor interface {
	VCenterVRSConfigs(*bambou.FetchingInfo) (VCenterVRSConfigsList, *bambou.Error)
}

// VCenterVRSConfigsParent is the interface that a parent of a VCenterVRSConfig must implement.
// A Parent is defined as an entity that has VCenterVRSConfig as a child.
// A Parent is an Ancestor which can create a VCenterVRSConfig.
type VCenterVRSConfigsParent interface {
	VCenterVRSConfigsAncestor
	CreateVCenterVRSConfig(*VCenterVRSConfig) *bambou.Error
}

// VCenterVRSConfig represents the model of a vrsconfig
type VCenterVRSConfig struct {
	ID                               string `json:"ID,omitempty"`
	ParentID                         string `json:"parentID,omitempty"`
	ParentType                       string `json:"parentType,omitempty"`
	Owner                            string `json:"owner,omitempty"`
	VRequireNuageMetadata            bool   `json:"vRequireNuageMetadata"`
	LastUpdatedBy                    string `json:"lastUpdatedBy,omitempty"`
	DataDNS1                         string `json:"dataDNS1,omitempty"`
	DataDNS2                         string `json:"dataDNS2,omitempty"`
	DataGateway                      string `json:"dataGateway,omitempty"`
	DataNetworkPortgroup             string `json:"dataNetworkPortgroup,omitempty"`
	DatapathSyncTimeout              int    `json:"datapathSyncTimeout,omitempty"`
	SecondaryNuageController         string `json:"secondaryNuageController,omitempty"`
	GenericSplitActivation           bool   `json:"genericSplitActivation"`
	SeparateDataNetwork              bool   `json:"separateDataNetwork"`
	Personality                      string `json:"personality,omitempty"`
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
	PrimaryNuageController           string `json:"primaryNuageController,omitempty"`
	VrsPassword                      string `json:"vrsPassword,omitempty"`
	VrsUserName                      string `json:"vrsUserName,omitempty"`
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
	ExternalID                       string `json:"externalID,omitempty"`
}

// NewVCenterVRSConfig returns a new *VCenterVRSConfig
func NewVCenterVRSConfig() *VCenterVRSConfig {

	return &VCenterVRSConfig{}
}

// Identity returns the Identity of the object.
func (o *VCenterVRSConfig) Identity() bambou.Identity {

	return VCenterVRSConfigIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VCenterVRSConfig) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VCenterVRSConfig) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VCenterVRSConfig from the server
func (o *VCenterVRSConfig) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VCenterVRSConfig into the server
func (o *VCenterVRSConfig) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VCenterVRSConfig from the server
func (o *VCenterVRSConfig) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the VCenterVRSConfig
func (o *VCenterVRSConfig) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VCenterVRSConfig
func (o *VCenterVRSConfig) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VCenterVRSConfig
func (o *VCenterVRSConfig) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VCenterVRSConfig
func (o *VCenterVRSConfig) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VRSAddressRanges retrieves the list of child VRSAddressRanges of the VCenterVRSConfig
func (o *VCenterVRSConfig) VRSAddressRanges(info *bambou.FetchingInfo) (VRSAddressRangesList, *bambou.Error) {

	var list VRSAddressRangesList
	err := bambou.CurrentSession().FetchChildren(o, VRSAddressRangeIdentity, &list, info)
	return list, err
}

// CreateVRSAddressRange creates a new child VRSAddressRange under the VCenterVRSConfig
func (o *VCenterVRSConfig) CreateVRSAddressRange(child *VRSAddressRange) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VRSRedeploymentpolicies retrieves the list of child VRSRedeploymentpolicies of the VCenterVRSConfig
func (o *VCenterVRSConfig) VRSRedeploymentpolicies(info *bambou.FetchingInfo) (VRSRedeploymentpoliciesList, *bambou.Error) {

	var list VRSRedeploymentpoliciesList
	err := bambou.CurrentSession().FetchChildren(o, VRSRedeploymentpolicyIdentity, &list, info)
	return list, err
}

// CreateVRSRedeploymentpolicy creates a new child VRSRedeploymentpolicy under the VCenterVRSConfig
func (o *VCenterVRSConfig) CreateVRSRedeploymentpolicy(child *VRSRedeploymentpolicy) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
