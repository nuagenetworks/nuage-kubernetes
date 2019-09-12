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

// NetconfGatewayIdentity represents the Identity of the object
var NetconfGatewayIdentity = bambou.Identity{
	Name:     "netconfgateway",
	Category: "netconfgateways",
}

// NetconfGatewaysList represents a list of NetconfGateways
type NetconfGatewaysList []*NetconfGateway

// NetconfGatewaysAncestor is the interface that an ancestor of a NetconfGateway must implement.
// An Ancestor is defined as an entity that has NetconfGateway as a descendant.
// An Ancestor can get a list of its child NetconfGateways, but not necessarily create one.
type NetconfGatewaysAncestor interface {
	NetconfGateways(*bambou.FetchingInfo) (NetconfGatewaysList, *bambou.Error)
}

// NetconfGatewaysParent is the interface that a parent of a NetconfGateway must implement.
// A Parent is defined as an entity that has NetconfGateway as a child.
// A Parent is an Ancestor which can create a NetconfGateway.
type NetconfGatewaysParent interface {
	NetconfGatewaysAncestor
	CreateNetconfGateway(*NetconfGateway) *bambou.Error
}

// NetconfGateway represents the model of a netconfgateway
type NetconfGateway struct {
	ID                                 string        `json:"ID,omitempty"`
	ParentID                           string        `json:"parentID,omitempty"`
	ParentType                         string        `json:"parentType,omitempty"`
	Owner                              string        `json:"owner,omitempty"`
	MACAddress                         string        `json:"MACAddress,omitempty"`
	ZFBMatchAttribute                  string        `json:"ZFBMatchAttribute,omitempty"`
	ZFBMatchValue                      string        `json:"ZFBMatchValue,omitempty"`
	BIOSReleaseDate                    string        `json:"BIOSReleaseDate,omitempty"`
	BIOSVersion                        string        `json:"BIOSVersion,omitempty"`
	CPUType                            string        `json:"CPUType,omitempty"`
	UUID                               string        `json:"UUID,omitempty"`
	Name                               string        `json:"name,omitempty"`
	Family                             string        `json:"family,omitempty"`
	ManagementID                       string        `json:"managementID,omitempty"`
	LastUpdatedBy                      string        `json:"lastUpdatedBy,omitempty"`
	DatapathID                         string        `json:"datapathID,omitempty"`
	Patches                            string        `json:"patches,omitempty"`
	GatewayConnected                   bool          `json:"gatewayConnected"`
	GatewayModel                       string        `json:"gatewayModel,omitempty"`
	GatewayVersion                     string        `json:"gatewayVersion,omitempty"`
	RedundancyGroupID                  string        `json:"redundancyGroupID,omitempty"`
	Peer                               string        `json:"peer,omitempty"`
	TemplateID                         string        `json:"templateID,omitempty"`
	Pending                            bool          `json:"pending"`
	Vendor                             string        `json:"vendor,omitempty"`
	SerialNumber                       string        `json:"serialNumber,omitempty"`
	PermittedAction                    string        `json:"permittedAction,omitempty"`
	Personality                        string        `json:"personality,omitempty"`
	Description                        string        `json:"description,omitempty"`
	Libraries                          string        `json:"libraries,omitempty"`
	EmbeddedMetadata                   []interface{} `json:"embeddedMetadata,omitempty"`
	EnterpriseID                       string        `json:"enterpriseID,omitempty"`
	EntityScope                        string        `json:"entityScope,omitempty"`
	LocationID                         string        `json:"locationID,omitempty"`
	BootstrapID                        string        `json:"bootstrapID,omitempty"`
	BootstrapStatus                    string        `json:"bootstrapStatus,omitempty"`
	ProductName                        string        `json:"productName,omitempty"`
	UseGatewayVLANVNID                 bool          `json:"useGatewayVLANVNID"`
	AssociatedGatewaySecurityID        string        `json:"associatedGatewaySecurityID,omitempty"`
	AssociatedGatewaySecurityProfileID string        `json:"associatedGatewaySecurityProfileID,omitempty"`
	AssociatedNSGInfoID                string        `json:"associatedNSGInfoID,omitempty"`
	AssociatedNetconfProfileID         string        `json:"associatedNetconfProfileID,omitempty"`
	Vtep                               string        `json:"vtep,omitempty"`
	AutoDiscGatewayID                  string        `json:"autoDiscGatewayID,omitempty"`
	ExternalID                         string        `json:"externalID,omitempty"`
	SystemID                           string        `json:"systemID,omitempty"`
}

// NewNetconfGateway returns a new *NetconfGateway
func NewNetconfGateway() *NetconfGateway {

	return &NetconfGateway{
		ZFBMatchAttribute: "NONE",
		GatewayConnected:  false,
	}
}

// Identity returns the Identity of the object.
func (o *NetconfGateway) Identity() bambou.Identity {

	return NetconfGatewayIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *NetconfGateway) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *NetconfGateway) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the NetconfGateway from the server
func (o *NetconfGateway) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the NetconfGateway into the server
func (o *NetconfGateway) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the NetconfGateway from the server
func (o *NetconfGateway) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// L2Domains retrieves the list of child L2Domains of the NetconfGateway
func (o *NetconfGateway) L2Domains(info *bambou.FetchingInfo) (L2DomainsList, *bambou.Error) {

	var list L2DomainsList
	err := bambou.CurrentSession().FetchChildren(o, L2DomainIdentity, &list, info)
	return list, err
}

// MACFilterProfiles retrieves the list of child MACFilterProfiles of the NetconfGateway
func (o *NetconfGateway) MACFilterProfiles(info *bambou.FetchingInfo) (MACFilterProfilesList, *bambou.Error) {

	var list MACFilterProfilesList
	err := bambou.CurrentSession().FetchChildren(o, MACFilterProfileIdentity, &list, info)
	return list, err
}

// SAPEgressQoSProfiles retrieves the list of child SAPEgressQoSProfiles of the NetconfGateway
func (o *NetconfGateway) SAPEgressQoSProfiles(info *bambou.FetchingInfo) (SAPEgressQoSProfilesList, *bambou.Error) {

	var list SAPEgressQoSProfilesList
	err := bambou.CurrentSession().FetchChildren(o, SAPEgressQoSProfileIdentity, &list, info)
	return list, err
}

// SAPIngressQoSProfiles retrieves the list of child SAPIngressQoSProfiles of the NetconfGateway
func (o *NetconfGateway) SAPIngressQoSProfiles(info *bambou.FetchingInfo) (SAPIngressQoSProfilesList, *bambou.Error) {

	var list SAPIngressQoSProfilesList
	err := bambou.CurrentSession().FetchChildren(o, SAPIngressQoSProfileIdentity, &list, info)
	return list, err
}

// GatewaySecurities retrieves the list of child GatewaySecurities of the NetconfGateway
func (o *NetconfGateway) GatewaySecurities(info *bambou.FetchingInfo) (GatewaySecuritiesList, *bambou.Error) {

	var list GatewaySecuritiesList
	err := bambou.CurrentSession().FetchChildren(o, GatewaySecurityIdentity, &list, info)
	return list, err
}

// PATNATPools retrieves the list of child PATNATPools of the NetconfGateway
func (o *NetconfGateway) PATNATPools(info *bambou.FetchingInfo) (PATNATPoolsList, *bambou.Error) {

	var list PATNATPoolsList
	err := bambou.CurrentSession().FetchChildren(o, PATNATPoolIdentity, &list, info)
	return list, err
}

// AssignPATNATPools assigns the list of PATNATPools to the NetconfGateway
func (o *NetconfGateway) AssignPATNATPools(children PATNATPoolsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, PATNATPoolIdentity)
}

// DeploymentFailures retrieves the list of child DeploymentFailures of the NetconfGateway
func (o *NetconfGateway) DeploymentFailures(info *bambou.FetchingInfo) (DeploymentFailuresList, *bambou.Error) {

	var list DeploymentFailuresList
	err := bambou.CurrentSession().FetchChildren(o, DeploymentFailureIdentity, &list, info)
	return list, err
}

// Permissions retrieves the list of child Permissions of the NetconfGateway
func (o *NetconfGateway) Permissions(info *bambou.FetchingInfo) (PermissionsList, *bambou.Error) {

	var list PermissionsList
	err := bambou.CurrentSession().FetchChildren(o, PermissionIdentity, &list, info)
	return list, err
}

// CreatePermission creates a new child Permission under the NetconfGateway
func (o *NetconfGateway) CreatePermission(child *Permission) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// WANServices retrieves the list of child WANServices of the NetconfGateway
func (o *NetconfGateway) WANServices(info *bambou.FetchingInfo) (WANServicesList, *bambou.Error) {

	var list WANServicesList
	err := bambou.CurrentSession().FetchChildren(o, WANServiceIdentity, &list, info)
	return list, err
}

// CreateWANService creates a new child WANService under the NetconfGateway
func (o *NetconfGateway) CreateWANService(child *WANService) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the NetconfGateway
func (o *NetconfGateway) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the NetconfGateway
func (o *NetconfGateway) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EgressProfiles retrieves the list of child EgressProfiles of the NetconfGateway
func (o *NetconfGateway) EgressProfiles(info *bambou.FetchingInfo) (EgressProfilesList, *bambou.Error) {

	var list EgressProfilesList
	err := bambou.CurrentSession().FetchChildren(o, EgressProfileIdentity, &list, info)
	return list, err
}

// CreateEgressProfile creates a new child EgressProfile under the NetconfGateway
func (o *NetconfGateway) CreateEgressProfile(child *EgressProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Alarms retrieves the list of child Alarms of the NetconfGateway
func (o *NetconfGateway) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the NetconfGateway
func (o *NetconfGateway) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the NetconfGateway
func (o *NetconfGateway) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// InfrastructureConfigs retrieves the list of child InfrastructureConfigs of the NetconfGateway
func (o *NetconfGateway) InfrastructureConfigs(info *bambou.FetchingInfo) (InfrastructureConfigsList, *bambou.Error) {

	var list InfrastructureConfigsList
	err := bambou.CurrentSession().FetchChildren(o, InfrastructureConfigIdentity, &list, info)
	return list, err
}

// IngressProfiles retrieves the list of child IngressProfiles of the NetconfGateway
func (o *NetconfGateway) IngressProfiles(info *bambou.FetchingInfo) (IngressProfilesList, *bambou.Error) {

	var list IngressProfilesList
	err := bambou.CurrentSession().FetchChildren(o, IngressProfileIdentity, &list, info)
	return list, err
}

// CreateIngressProfile creates a new child IngressProfile under the NetconfGateway
func (o *NetconfGateway) CreateIngressProfile(child *IngressProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EnterprisePermissions retrieves the list of child EnterprisePermissions of the NetconfGateway
func (o *NetconfGateway) EnterprisePermissions(info *bambou.FetchingInfo) (EnterprisePermissionsList, *bambou.Error) {

	var list EnterprisePermissionsList
	err := bambou.CurrentSession().FetchChildren(o, EnterprisePermissionIdentity, &list, info)
	return list, err
}

// CreateEnterprisePermission creates a new child EnterprisePermission under the NetconfGateway
func (o *NetconfGateway) CreateEnterprisePermission(child *EnterprisePermission) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Jobs retrieves the list of child Jobs of the NetconfGateway
func (o *NetconfGateway) Jobs(info *bambou.FetchingInfo) (JobsList, *bambou.Error) {

	var list JobsList
	err := bambou.CurrentSession().FetchChildren(o, JobIdentity, &list, info)
	return list, err
}

// CreateJob creates a new child Job under the NetconfGateway
func (o *NetconfGateway) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Locations retrieves the list of child Locations of the NetconfGateway
func (o *NetconfGateway) Locations(info *bambou.FetchingInfo) (LocationsList, *bambou.Error) {

	var list LocationsList
	err := bambou.CurrentSession().FetchChildren(o, LocationIdentity, &list, info)
	return list, err
}

// Domains retrieves the list of child Domains of the NetconfGateway
func (o *NetconfGateway) Domains(info *bambou.FetchingInfo) (DomainsList, *bambou.Error) {

	var list DomainsList
	err := bambou.CurrentSession().FetchChildren(o, DomainIdentity, &list, info)
	return list, err
}

// Bootstraps retrieves the list of child Bootstraps of the NetconfGateway
func (o *NetconfGateway) Bootstraps(info *bambou.FetchingInfo) (BootstrapsList, *bambou.Error) {

	var list BootstrapsList
	err := bambou.CurrentSession().FetchChildren(o, BootstrapIdentity, &list, info)
	return list, err
}

// CreateBootstrapActivation creates a new child BootstrapActivation under the NetconfGateway
func (o *NetconfGateway) CreateBootstrapActivation(child *BootstrapActivation) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Ports retrieves the list of child Ports of the NetconfGateway
func (o *NetconfGateway) Ports(info *bambou.FetchingInfo) (PortsList, *bambou.Error) {

	var list PortsList
	err := bambou.CurrentSession().FetchChildren(o, PortIdentity, &list, info)
	return list, err
}

// CreatePort creates a new child Port under the NetconfGateway
func (o *NetconfGateway) CreatePort(child *Port) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IPFilterProfiles retrieves the list of child IPFilterProfiles of the NetconfGateway
func (o *NetconfGateway) IPFilterProfiles(info *bambou.FetchingInfo) (IPFilterProfilesList, *bambou.Error) {

	var list IPFilterProfilesList
	err := bambou.CurrentSession().FetchChildren(o, IPFilterProfileIdentity, &list, info)
	return list, err
}

// IPv6FilterProfiles retrieves the list of child IPv6FilterProfiles of the NetconfGateway
func (o *NetconfGateway) IPv6FilterProfiles(info *bambou.FetchingInfo) (IPv6FilterProfilesList, *bambou.Error) {

	var list IPv6FilterProfilesList
	err := bambou.CurrentSession().FetchChildren(o, IPv6FilterProfileIdentity, &list, info)
	return list, err
}

// Subnets retrieves the list of child Subnets of the NetconfGateway
func (o *NetconfGateway) Subnets(info *bambou.FetchingInfo) (SubnetsList, *bambou.Error) {

	var list SubnetsList
	err := bambou.CurrentSession().FetchChildren(o, SubnetIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the NetconfGateway
func (o *NetconfGateway) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
