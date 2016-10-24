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

// ContainerInterfaceIdentity represents the Identity of the object
var ContainerInterfaceIdentity = bambou.Identity{
	Name:     "containerinterface",
	Category: "containerinterfaces",
}

// ContainerInterfacesList represents a list of ContainerInterfaces
type ContainerInterfacesList []*ContainerInterface

// ContainerInterfacesAncestor is the interface of an ancestor of a ContainerInterface must implement.
type ContainerInterfacesAncestor interface {
	ContainerInterfaces(*bambou.FetchingInfo) (ContainerInterfacesList, *bambou.Error)
	CreateContainerInterfaces(*ContainerInterface) *bambou.Error
}

// ContainerInterface represents the model of a containerinterface
type ContainerInterface struct {
	ID                          string `json:"ID,omitempty"`
	ParentID                    string `json:"parentID,omitempty"`
	ParentType                  string `json:"parentType,omitempty"`
	Owner                       string `json:"owner,omitempty"`
	MAC                         string `json:"MAC,omitempty"`
	IPAddress                   string `json:"IPAddress,omitempty"`
	VPortID                     string `json:"VPortID,omitempty"`
	VPortName                   string `json:"VPortName,omitempty"`
	Name                        string `json:"name,omitempty"`
	LastUpdatedBy               string `json:"lastUpdatedBy,omitempty"`
	Gateway                     string `json:"gateway,omitempty"`
	Netmask                     string `json:"netmask,omitempty"`
	NetworkID                   string `json:"networkID,omitempty"`
	NetworkName                 string `json:"networkName,omitempty"`
	TierID                      string `json:"tierID,omitempty"`
	EndpointID                  string `json:"endpointID,omitempty"`
	EntityScope                 string `json:"entityScope,omitempty"`
	PolicyDecisionID            string `json:"policyDecisionID,omitempty"`
	DomainID                    string `json:"domainID,omitempty"`
	DomainName                  string `json:"domainName,omitempty"`
	ZoneID                      string `json:"zoneID,omitempty"`
	ZoneName                    string `json:"zoneName,omitempty"`
	ContainerUUID               string `json:"containerUUID,omitempty"`
	AssociatedFloatingIPAddress string `json:"associatedFloatingIPAddress,omitempty"`
	AttachedNetworkID           string `json:"attachedNetworkID,omitempty"`
	AttachedNetworkType         string `json:"attachedNetworkType,omitempty"`
	MultiNICVPortName           string `json:"multiNICVPortName,omitempty"`
	ExternalID                  string `json:"externalID,omitempty"`
}

// NewContainerInterface returns a new *ContainerInterface
func NewContainerInterface() *ContainerInterface {

	return &ContainerInterface{}
}

// Identity returns the Identity of the object.
func (o *ContainerInterface) Identity() bambou.Identity {

	return ContainerInterfaceIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *ContainerInterface) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *ContainerInterface) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the ContainerInterface from the server
func (o *ContainerInterface) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the ContainerInterface into the server
func (o *ContainerInterface) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the ContainerInterface from the server
func (o *ContainerInterface) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// TCAs retrieves the list of child TCAs of the ContainerInterface
func (o *ContainerInterface) TCAs(info *bambou.FetchingInfo) (TCAsList, *bambou.Error) {

	var list TCAsList
	err := bambou.CurrentSession().FetchChildren(o, TCAIdentity, &list, info)
	return list, err
}

// CreateTCA creates a new child TCA under the ContainerInterface
func (o *ContainerInterface) CreateTCA(child *TCA) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// RedirectionTargets retrieves the list of child RedirectionTargets of the ContainerInterface
func (o *ContainerInterface) RedirectionTargets(info *bambou.FetchingInfo) (RedirectionTargetsList, *bambou.Error) {

	var list RedirectionTargetsList
	err := bambou.CurrentSession().FetchChildren(o, RedirectionTargetIdentity, &list, info)
	return list, err
}

// CreateRedirectionTarget creates a new child RedirectionTarget under the ContainerInterface
func (o *ContainerInterface) CreateRedirectionTarget(child *RedirectionTarget) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the ContainerInterface
func (o *ContainerInterface) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the ContainerInterface
func (o *ContainerInterface) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DHCPOptions retrieves the list of child DHCPOptions of the ContainerInterface
func (o *ContainerInterface) DHCPOptions(info *bambou.FetchingInfo) (DHCPOptionsList, *bambou.Error) {

	var list DHCPOptionsList
	err := bambou.CurrentSession().FetchChildren(o, DHCPOptionIdentity, &list, info)
	return list, err
}

// CreateDHCPOption creates a new child DHCPOption under the ContainerInterface
func (o *ContainerInterface) CreateDHCPOption(child *DHCPOption) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the ContainerInterface
func (o *ContainerInterface) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the ContainerInterface
func (o *ContainerInterface) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PolicyDecisions retrieves the list of child PolicyDecisions of the ContainerInterface
func (o *ContainerInterface) PolicyDecisions(info *bambou.FetchingInfo) (PolicyDecisionsList, *bambou.Error) {

	var list PolicyDecisionsList
	err := bambou.CurrentSession().FetchChildren(o, PolicyDecisionIdentity, &list, info)
	return list, err
}

// CreatePolicyDecision creates a new child PolicyDecision under the ContainerInterface
func (o *ContainerInterface) CreatePolicyDecision(child *PolicyDecision) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PolicyGroups retrieves the list of child PolicyGroups of the ContainerInterface
func (o *ContainerInterface) PolicyGroups(info *bambou.FetchingInfo) (PolicyGroupsList, *bambou.Error) {

	var list PolicyGroupsList
	err := bambou.CurrentSession().FetchChildren(o, PolicyGroupIdentity, &list, info)
	return list, err
}

// CreatePolicyGroup creates a new child PolicyGroup under the ContainerInterface
func (o *ContainerInterface) CreatePolicyGroup(child *PolicyGroup) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// StaticRoutes retrieves the list of child StaticRoutes of the ContainerInterface
func (o *ContainerInterface) StaticRoutes(info *bambou.FetchingInfo) (StaticRoutesList, *bambou.Error) {

	var list StaticRoutesList
	err := bambou.CurrentSession().FetchChildren(o, StaticRouteIdentity, &list, info)
	return list, err
}

// CreateStaticRoute creates a new child StaticRoute under the ContainerInterface
func (o *ContainerInterface) CreateStaticRoute(child *StaticRoute) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Statistics retrieves the list of child Statistics of the ContainerInterface
func (o *ContainerInterface) Statistics(info *bambou.FetchingInfo) (StatisticsList, *bambou.Error) {

	var list StatisticsList
	err := bambou.CurrentSession().FetchChildren(o, StatisticsIdentity, &list, info)
	return list, err
}

// CreateStatistics creates a new child Statistics under the ContainerInterface
func (o *ContainerInterface) CreateStatistics(child *Statistics) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// MultiCastChannelMaps retrieves the list of child MultiCastChannelMaps of the ContainerInterface
func (o *ContainerInterface) MultiCastChannelMaps(info *bambou.FetchingInfo) (MultiCastChannelMapsList, *bambou.Error) {

	var list MultiCastChannelMapsList
	err := bambou.CurrentSession().FetchChildren(o, MultiCastChannelMapIdentity, &list, info)
	return list, err
}

// CreateMultiCastChannelMap creates a new child MultiCastChannelMap under the ContainerInterface
func (o *ContainerInterface) CreateMultiCastChannelMap(child *MultiCastChannelMap) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the ContainerInterface
func (o *ContainerInterface) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}

// CreateEventLog creates a new child EventLog under the ContainerInterface
func (o *ContainerInterface) CreateEventLog(child *EventLog) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
