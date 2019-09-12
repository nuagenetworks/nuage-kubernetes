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

// GatewayIdentity represents the Identity of the object
var GatewayIdentity = bambou.Identity{
	Name:     "gateway",
	Category: "gateways",
}

// GatewaysList represents a list of Gateways
type GatewaysList []*Gateway

// GatewaysAncestor is the interface that an ancestor of a Gateway must implement.
// An Ancestor is defined as an entity that has Gateway as a descendant.
// An Ancestor can get a list of its child Gateways, but not necessarily create one.
type GatewaysAncestor interface {
	Gateways(*bambou.FetchingInfo) (GatewaysList, *bambou.Error)
}

// GatewaysParent is the interface that a parent of a Gateway must implement.
// A Parent is defined as an entity that has Gateway as a child.
// A Parent is an Ancestor which can create a Gateway.
type GatewaysParent interface {
	GatewaysAncestor
	CreateGateway(*Gateway) *bambou.Error
}

// Gateway represents the model of a gateway
type Gateway struct {
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
	GatewayConfigRawVersion            string        `json:"gatewayConfigRawVersion,omitempty"`
	GatewayConfigVersion               string        `json:"gatewayConfigVersion,omitempty"`
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

// NewGateway returns a new *Gateway
func NewGateway() *Gateway {

	return &Gateway{
		ZFBMatchAttribute: "NONE",
		GatewayConnected:  false,
		Personality:       "VRSG",
	}
}

// Identity returns the Identity of the object.
func (o *Gateway) Identity() bambou.Identity {

	return GatewayIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Gateway) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Gateway) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Gateway from the server
func (o *Gateway) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Gateway into the server
func (o *Gateway) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Gateway from the server
func (o *Gateway) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// L2Domains retrieves the list of child L2Domains of the Gateway
func (o *Gateway) L2Domains(info *bambou.FetchingInfo) (L2DomainsList, *bambou.Error) {

	var list L2DomainsList
	err := bambou.CurrentSession().FetchChildren(o, L2DomainIdentity, &list, info)
	return list, err
}

// MACFilterProfiles retrieves the list of child MACFilterProfiles of the Gateway
func (o *Gateway) MACFilterProfiles(info *bambou.FetchingInfo) (MACFilterProfilesList, *bambou.Error) {

	var list MACFilterProfilesList
	err := bambou.CurrentSession().FetchChildren(o, MACFilterProfileIdentity, &list, info)
	return list, err
}

// SAPEgressQoSProfiles retrieves the list of child SAPEgressQoSProfiles of the Gateway
func (o *Gateway) SAPEgressQoSProfiles(info *bambou.FetchingInfo) (SAPEgressQoSProfilesList, *bambou.Error) {

	var list SAPEgressQoSProfilesList
	err := bambou.CurrentSession().FetchChildren(o, SAPEgressQoSProfileIdentity, &list, info)
	return list, err
}

// SAPIngressQoSProfiles retrieves the list of child SAPIngressQoSProfiles of the Gateway
func (o *Gateway) SAPIngressQoSProfiles(info *bambou.FetchingInfo) (SAPIngressQoSProfilesList, *bambou.Error) {

	var list SAPIngressQoSProfilesList
	err := bambou.CurrentSession().FetchChildren(o, SAPIngressQoSProfileIdentity, &list, info)
	return list, err
}

// GatewaySecurities retrieves the list of child GatewaySecurities of the Gateway
func (o *Gateway) GatewaySecurities(info *bambou.FetchingInfo) (GatewaySecuritiesList, *bambou.Error) {

	var list GatewaySecuritiesList
	err := bambou.CurrentSession().FetchChildren(o, GatewaySecurityIdentity, &list, info)
	return list, err
}

// PATNATPools retrieves the list of child PATNATPools of the Gateway
func (o *Gateway) PATNATPools(info *bambou.FetchingInfo) (PATNATPoolsList, *bambou.Error) {

	var list PATNATPoolsList
	err := bambou.CurrentSession().FetchChildren(o, PATNATPoolIdentity, &list, info)
	return list, err
}

// AssignPATNATPools assigns the list of PATNATPools to the Gateway
func (o *Gateway) AssignPATNATPools(children PATNATPoolsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, PATNATPoolIdentity)
}

// DeploymentFailures retrieves the list of child DeploymentFailures of the Gateway
func (o *Gateway) DeploymentFailures(info *bambou.FetchingInfo) (DeploymentFailuresList, *bambou.Error) {

	var list DeploymentFailuresList
	err := bambou.CurrentSession().FetchChildren(o, DeploymentFailureIdentity, &list, info)
	return list, err
}

// Permissions retrieves the list of child Permissions of the Gateway
func (o *Gateway) Permissions(info *bambou.FetchingInfo) (PermissionsList, *bambou.Error) {

	var list PermissionsList
	err := bambou.CurrentSession().FetchChildren(o, PermissionIdentity, &list, info)
	return list, err
}

// CreatePermission creates a new child Permission under the Gateway
func (o *Gateway) CreatePermission(child *Permission) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// WANServices retrieves the list of child WANServices of the Gateway
func (o *Gateway) WANServices(info *bambou.FetchingInfo) (WANServicesList, *bambou.Error) {

	var list WANServicesList
	err := bambou.CurrentSession().FetchChildren(o, WANServiceIdentity, &list, info)
	return list, err
}

// CreateWANService creates a new child WANService under the Gateway
func (o *Gateway) CreateWANService(child *WANService) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the Gateway
func (o *Gateway) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the Gateway
func (o *Gateway) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EgressProfiles retrieves the list of child EgressProfiles of the Gateway
func (o *Gateway) EgressProfiles(info *bambou.FetchingInfo) (EgressProfilesList, *bambou.Error) {

	var list EgressProfilesList
	err := bambou.CurrentSession().FetchChildren(o, EgressProfileIdentity, &list, info)
	return list, err
}

// CreateEgressProfile creates a new child EgressProfile under the Gateway
func (o *Gateway) CreateEgressProfile(child *EgressProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Alarms retrieves the list of child Alarms of the Gateway
func (o *Gateway) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Gateway
func (o *Gateway) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Gateway
func (o *Gateway) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// InfrastructureConfigs retrieves the list of child InfrastructureConfigs of the Gateway
func (o *Gateway) InfrastructureConfigs(info *bambou.FetchingInfo) (InfrastructureConfigsList, *bambou.Error) {

	var list InfrastructureConfigsList
	err := bambou.CurrentSession().FetchChildren(o, InfrastructureConfigIdentity, &list, info)
	return list, err
}

// IngressProfiles retrieves the list of child IngressProfiles of the Gateway
func (o *Gateway) IngressProfiles(info *bambou.FetchingInfo) (IngressProfilesList, *bambou.Error) {

	var list IngressProfilesList
	err := bambou.CurrentSession().FetchChildren(o, IngressProfileIdentity, &list, info)
	return list, err
}

// CreateIngressProfile creates a new child IngressProfile under the Gateway
func (o *Gateway) CreateIngressProfile(child *IngressProfile) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EnterprisePermissions retrieves the list of child EnterprisePermissions of the Gateway
func (o *Gateway) EnterprisePermissions(info *bambou.FetchingInfo) (EnterprisePermissionsList, *bambou.Error) {

	var list EnterprisePermissionsList
	err := bambou.CurrentSession().FetchChildren(o, EnterprisePermissionIdentity, &list, info)
	return list, err
}

// CreateEnterprisePermission creates a new child EnterprisePermission under the Gateway
func (o *Gateway) CreateEnterprisePermission(child *EnterprisePermission) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Jobs retrieves the list of child Jobs of the Gateway
func (o *Gateway) Jobs(info *bambou.FetchingInfo) (JobsList, *bambou.Error) {

	var list JobsList
	err := bambou.CurrentSession().FetchChildren(o, JobIdentity, &list, info)
	return list, err
}

// CreateJob creates a new child Job under the Gateway
func (o *Gateway) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Locations retrieves the list of child Locations of the Gateway
func (o *Gateway) Locations(info *bambou.FetchingInfo) (LocationsList, *bambou.Error) {

	var list LocationsList
	err := bambou.CurrentSession().FetchChildren(o, LocationIdentity, &list, info)
	return list, err
}

// Domains retrieves the list of child Domains of the Gateway
func (o *Gateway) Domains(info *bambou.FetchingInfo) (DomainsList, *bambou.Error) {

	var list DomainsList
	err := bambou.CurrentSession().FetchChildren(o, DomainIdentity, &list, info)
	return list, err
}

// Bootstraps retrieves the list of child Bootstraps of the Gateway
func (o *Gateway) Bootstraps(info *bambou.FetchingInfo) (BootstrapsList, *bambou.Error) {

	var list BootstrapsList
	err := bambou.CurrentSession().FetchChildren(o, BootstrapIdentity, &list, info)
	return list, err
}

// CreateBootstrapActivation creates a new child BootstrapActivation under the Gateway
func (o *Gateway) CreateBootstrapActivation(child *BootstrapActivation) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Ports retrieves the list of child Ports of the Gateway
func (o *Gateway) Ports(info *bambou.FetchingInfo) (PortsList, *bambou.Error) {

	var list PortsList
	err := bambou.CurrentSession().FetchChildren(o, PortIdentity, &list, info)
	return list, err
}

// CreatePort creates a new child Port under the Gateway
func (o *Gateway) CreatePort(child *Port) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IPFilterProfiles retrieves the list of child IPFilterProfiles of the Gateway
func (o *Gateway) IPFilterProfiles(info *bambou.FetchingInfo) (IPFilterProfilesList, *bambou.Error) {

	var list IPFilterProfilesList
	err := bambou.CurrentSession().FetchChildren(o, IPFilterProfileIdentity, &list, info)
	return list, err
}

// IPv6FilterProfiles retrieves the list of child IPv6FilterProfiles of the Gateway
func (o *Gateway) IPv6FilterProfiles(info *bambou.FetchingInfo) (IPv6FilterProfilesList, *bambou.Error) {

	var list IPv6FilterProfilesList
	err := bambou.CurrentSession().FetchChildren(o, IPv6FilterProfileIdentity, &list, info)
	return list, err
}

// Subnets retrieves the list of child Subnets of the Gateway
func (o *Gateway) Subnets(info *bambou.FetchingInfo) (SubnetsList, *bambou.Error) {

	var list SubnetsList
	err := bambou.CurrentSession().FetchChildren(o, SubnetIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the Gateway
func (o *Gateway) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
