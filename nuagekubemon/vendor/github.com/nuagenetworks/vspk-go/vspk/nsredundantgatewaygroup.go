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

// NSRedundantGatewayGroupIdentity represents the Identity of the object
var NSRedundantGatewayGroupIdentity = bambou.Identity{
	Name:     "nsgredundancygroup",
	Category: "nsgredundancygroups",
}

// NSRedundantGatewayGroupsList represents a list of NSRedundantGatewayGroups
type NSRedundantGatewayGroupsList []*NSRedundantGatewayGroup

// NSRedundantGatewayGroupsAncestor is the interface that an ancestor of a NSRedundantGatewayGroup must implement.
// An Ancestor is defined as an entity that has NSRedundantGatewayGroup as a descendant.
// An Ancestor can get a list of its child NSRedundantGatewayGroups, but not necessarily create one.
type NSRedundantGatewayGroupsAncestor interface {
	NSRedundantGatewayGroups(*bambou.FetchingInfo) (NSRedundantGatewayGroupsList, *bambou.Error)
}

// NSRedundantGatewayGroupsParent is the interface that a parent of a NSRedundantGatewayGroup must implement.
// A Parent is defined as an entity that has NSRedundantGatewayGroup as a child.
// A Parent is an Ancestor which can create a NSRedundantGatewayGroup.
type NSRedundantGatewayGroupsParent interface {
	NSRedundantGatewayGroupsAncestor
	CreateNSRedundantGatewayGroup(*NSRedundantGatewayGroup) *bambou.Error
}

// NSRedundantGatewayGroup represents the model of a nsgredundancygroup
type NSRedundantGatewayGroup struct {
	ID                                  string        `json:"ID,omitempty"`
	ParentID                            string        `json:"parentID,omitempty"`
	ParentType                          string        `json:"parentType,omitempty"`
	Owner                               string        `json:"owner,omitempty"`
	Name                                string        `json:"name,omitempty"`
	LastUpdatedBy                       string        `json:"lastUpdatedBy,omitempty"`
	GatewayPeer1AutodiscoveredGatewayID string        `json:"gatewayPeer1AutodiscoveredGatewayID,omitempty"`
	GatewayPeer1Connected               bool          `json:"gatewayPeer1Connected"`
	GatewayPeer1ID                      string        `json:"gatewayPeer1ID,omitempty"`
	GatewayPeer1Name                    string        `json:"gatewayPeer1Name,omitempty"`
	GatewayPeer2AutodiscoveredGatewayID string        `json:"gatewayPeer2AutodiscoveredGatewayID,omitempty"`
	GatewayPeer2Connected               bool          `json:"gatewayPeer2Connected"`
	GatewayPeer2ID                      string        `json:"gatewayPeer2ID,omitempty"`
	GatewayPeer2Name                    string        `json:"gatewayPeer2Name,omitempty"`
	HeartbeatInterval                   int           `json:"heartbeatInterval,omitempty"`
	HeartbeatVLANID                     int           `json:"heartbeatVLANID,omitempty"`
	RedundancyPortIDs                   []interface{} `json:"redundancyPortIDs,omitempty"`
	RedundantGatewayStatus              string        `json:"redundantGatewayStatus,omitempty"`
	PermittedAction                     string        `json:"permittedAction,omitempty"`
	Personality                         string        `json:"personality,omitempty"`
	Description                         string        `json:"description,omitempty"`
	EmbeddedMetadata                    []interface{} `json:"embeddedMetadata,omitempty"`
	EnterpriseID                        string        `json:"enterpriseID,omitempty"`
	EntityScope                         string        `json:"entityScope,omitempty"`
	ConsecutiveFailuresCount            int           `json:"consecutiveFailuresCount,omitempty"`
	ExternalID                          string        `json:"externalID,omitempty"`
}

// NewNSRedundantGatewayGroup returns a new *NSRedundantGatewayGroup
func NewNSRedundantGatewayGroup() *NSRedundantGatewayGroup {

	return &NSRedundantGatewayGroup{
		GatewayPeer1Connected:    false,
		GatewayPeer2Connected:    false,
		HeartbeatInterval:        500,
		HeartbeatVLANID:          4094,
		ConsecutiveFailuresCount: 3,
	}
}

// Identity returns the Identity of the object.
func (o *NSRedundantGatewayGroup) Identity() bambou.Identity {

	return NSRedundantGatewayGroupIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *NSRedundantGatewayGroup) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *NSRedundantGatewayGroup) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the NSRedundantGatewayGroup from the server
func (o *NSRedundantGatewayGroup) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the NSRedundantGatewayGroup into the server
func (o *NSRedundantGatewayGroup) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the NSRedundantGatewayGroup from the server
func (o *NSRedundantGatewayGroup) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Permissions retrieves the list of child Permissions of the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) Permissions(info *bambou.FetchingInfo) (PermissionsList, *bambou.Error) {

	var list PermissionsList
	err := bambou.CurrentSession().FetchChildren(o, PermissionIdentity, &list, info)
	return list, err
}

// CreatePermission creates a new child Permission under the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) CreatePermission(child *Permission) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// ShuntLinks retrieves the list of child ShuntLinks of the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) ShuntLinks(info *bambou.FetchingInfo) (ShuntLinksList, *bambou.Error) {

	var list ShuntLinksList
	err := bambou.CurrentSession().FetchChildren(o, ShuntLinkIdentity, &list, info)
	return list, err
}

// CreateShuntLink creates a new child ShuntLink under the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) CreateShuntLink(child *ShuntLink) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Alarms retrieves the list of child Alarms of the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EnterprisePermissions retrieves the list of child EnterprisePermissions of the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) EnterprisePermissions(info *bambou.FetchingInfo) (EnterprisePermissionsList, *bambou.Error) {

	var list EnterprisePermissionsList
	err := bambou.CurrentSession().FetchChildren(o, EnterprisePermissionIdentity, &list, info)
	return list, err
}

// CreateEnterprisePermission creates a new child EnterprisePermission under the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) CreateEnterprisePermission(child *EnterprisePermission) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// NSGateways retrieves the list of child NSGateways of the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) NSGateways(info *bambou.FetchingInfo) (NSGatewaysList, *bambou.Error) {

	var list NSGatewaysList
	err := bambou.CurrentSession().FetchChildren(o, NSGatewayIdentity, &list, info)
	return list, err
}

// RedundantPorts retrieves the list of child RedundantPorts of the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) RedundantPorts(info *bambou.FetchingInfo) (RedundantPortsList, *bambou.Error) {

	var list RedundantPortsList
	err := bambou.CurrentSession().FetchChildren(o, RedundantPortIdentity, &list, info)
	return list, err
}

// CreateRedundantPort creates a new child RedundantPort under the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) CreateRedundantPort(child *RedundantPort) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the NSRedundantGatewayGroup
func (o *NSRedundantGatewayGroup) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
