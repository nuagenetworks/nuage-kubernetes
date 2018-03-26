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

// VLANIdentity represents the Identity of the object
var VLANIdentity = bambou.Identity{
	Name:     "vlan",
	Category: "vlans",
}

// VLANsList represents a list of VLANs
type VLANsList []*VLAN

// VLANsAncestor is the interface that an ancestor of a VLAN must implement.
// An Ancestor is defined as an entity that has VLAN as a descendant.
// An Ancestor can get a list of its child VLANs, but not necessarily create one.
type VLANsAncestor interface {
	VLANs(*bambou.FetchingInfo) (VLANsList, *bambou.Error)
}

// VLANsParent is the interface that a parent of a VLAN must implement.
// A Parent is defined as an entity that has VLAN as a child.
// A Parent is an Ancestor which can create a VLAN.
type VLANsParent interface {
	VLANsAncestor
	CreateVLAN(*VLAN) *bambou.Error
}

// VLAN represents the model of a vlan
type VLAN struct {
	ID                           string `json:"ID,omitempty"`
	ParentID                     string `json:"parentID,omitempty"`
	ParentType                   string `json:"parentType,omitempty"`
	Owner                        string `json:"owner,omitempty"`
	Value                        int    `json:"value,omitempty"`
	LastUpdatedBy                string `json:"lastUpdatedBy,omitempty"`
	GatewayID                    string `json:"gatewayID,omitempty"`
	Readonly                     bool   `json:"readonly"`
	TemplateID                   string `json:"templateID,omitempty"`
	PermittedAction              string `json:"permittedAction,omitempty"`
	Description                  string `json:"description,omitempty"`
	Restricted                   bool   `json:"restricted"`
	EntityScope                  string `json:"entityScope,omitempty"`
	VportID                      string `json:"vportID,omitempty"`
	UseUserMnemonic              bool   `json:"useUserMnemonic"`
	UserMnemonic                 string `json:"userMnemonic,omitempty"`
	AssociatedBGPProfileID       string `json:"associatedBGPProfileID,omitempty"`
	AssociatedEgressQOSPolicyID  string `json:"associatedEgressQOSPolicyID,omitempty"`
	AssociatedUplinkConnectionID string `json:"associatedUplinkConnectionID,omitempty"`
	AssociatedVSCProfileID       string `json:"associatedVSCProfileID,omitempty"`
	Status                       string `json:"status,omitempty"`
	DucVlan                      bool   `json:"ducVlan"`
	ExternalID                   string `json:"externalID,omitempty"`
}

// NewVLAN returns a new *VLAN
func NewVLAN() *VLAN {

	return &VLAN{
		DucVlan: false,
	}
}

// Identity returns the Identity of the object.
func (o *VLAN) Identity() bambou.Identity {

	return VLANIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VLAN) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VLAN) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VLAN from the server
func (o *VLAN) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VLAN into the server
func (o *VLAN) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VLAN from the server
func (o *VLAN) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// PATNATPools retrieves the list of child PATNATPools of the VLAN
func (o *VLAN) PATNATPools(info *bambou.FetchingInfo) (PATNATPoolsList, *bambou.Error) {

	var list PATNATPoolsList
	err := bambou.CurrentSession().FetchChildren(o, PATNATPoolIdentity, &list, info)
	return list, err
}

// AssignPATNATPools assigns the list of PATNATPools to the VLAN
func (o *VLAN) AssignPATNATPools(children PATNATPoolsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, PATNATPoolIdentity)
}

// Permissions retrieves the list of child Permissions of the VLAN
func (o *VLAN) Permissions(info *bambou.FetchingInfo) (PermissionsList, *bambou.Error) {

	var list PermissionsList
	err := bambou.CurrentSession().FetchChildren(o, PermissionIdentity, &list, info)
	return list, err
}

// CreatePermission creates a new child Permission under the VLAN
func (o *VLAN) CreatePermission(child *Permission) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the VLAN
func (o *VLAN) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VLAN
func (o *VLAN) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// BGPNeighbors retrieves the list of child BGPNeighbors of the VLAN
func (o *VLAN) BGPNeighbors(info *bambou.FetchingInfo) (BGPNeighborsList, *bambou.Error) {

	var list BGPNeighborsList
	err := bambou.CurrentSession().FetchChildren(o, BGPNeighborIdentity, &list, info)
	return list, err
}

// CreateBGPNeighbor creates a new child BGPNeighbor under the VLAN
func (o *VLAN) CreateBGPNeighbor(child *BGPNeighbor) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// IKEGatewayConnections retrieves the list of child IKEGatewayConnections of the VLAN
func (o *VLAN) IKEGatewayConnections(info *bambou.FetchingInfo) (IKEGatewayConnectionsList, *bambou.Error) {

	var list IKEGatewayConnectionsList
	err := bambou.CurrentSession().FetchChildren(o, IKEGatewayConnectionIdentity, &list, info)
	return list, err
}

// CreateIKEGatewayConnection creates a new child IKEGatewayConnection under the VLAN
func (o *VLAN) CreateIKEGatewayConnection(child *IKEGatewayConnection) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Alarms retrieves the list of child Alarms of the VLAN
func (o *VLAN) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VLAN
func (o *VLAN) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VLAN
func (o *VLAN) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EnterprisePermissions retrieves the list of child EnterprisePermissions of the VLAN
func (o *VLAN) EnterprisePermissions(info *bambou.FetchingInfo) (EnterprisePermissionsList, *bambou.Error) {

	var list EnterprisePermissionsList
	err := bambou.CurrentSession().FetchChildren(o, EnterprisePermissionIdentity, &list, info)
	return list, err
}

// CreateEnterprisePermission creates a new child EnterprisePermission under the VLAN
func (o *VLAN) CreateEnterprisePermission(child *EnterprisePermission) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// UplinkConnections retrieves the list of child UplinkConnections of the VLAN
func (o *VLAN) UplinkConnections(info *bambou.FetchingInfo) (UplinkConnectionsList, *bambou.Error) {

	var list UplinkConnectionsList
	err := bambou.CurrentSession().FetchChildren(o, UplinkConnectionIdentity, &list, info)
	return list, err
}

// CreateUplinkConnection creates a new child UplinkConnection under the VLAN
func (o *VLAN) CreateUplinkConnection(child *UplinkConnection) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// BRConnections retrieves the list of child BRConnections of the VLAN
func (o *VLAN) BRConnections(info *bambou.FetchingInfo) (BRConnectionsList, *bambou.Error) {

	var list BRConnectionsList
	err := bambou.CurrentSession().FetchChildren(o, BRConnectionIdentity, &list, info)
	return list, err
}

// CreateBRConnection creates a new child BRConnection under the VLAN
func (o *VLAN) CreateBRConnection(child *BRConnection) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Ltestatistics retrieves the list of child Ltestatistics of the VLAN
func (o *VLAN) Ltestatistics(info *bambou.FetchingInfo) (LtestatisticsList, *bambou.Error) {

	var list LtestatisticsList
	err := bambou.CurrentSession().FetchChildren(o, LtestatisticsIdentity, &list, info)
	return list, err
}

// CreateLtestatistics creates a new child Ltestatistics under the VLAN
func (o *VLAN) CreateLtestatistics(child *Ltestatistics) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the VLAN
func (o *VLAN) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
