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

// VMInterfaceIdentity represents the Identity of the object
var VMInterfaceIdentity = bambou.Identity{
	Name:     "vminterface",
	Category: "vminterfaces",
}

// VMInterfacesList represents a list of VMInterfaces
type VMInterfacesList []*VMInterface

// VMInterfacesAncestor is the interface that an ancestor of a VMInterface must implement.
// An Ancestor is defined as an entity that has VMInterface as a descendant.
// An Ancestor can get a list of its child VMInterfaces, but not necessarily create one.
type VMInterfacesAncestor interface {
	VMInterfaces(*bambou.FetchingInfo) (VMInterfacesList, *bambou.Error)
}

// VMInterfacesParent is the interface that a parent of a VMInterface must implement.
// A Parent is defined as an entity that has VMInterface as a child.
// A Parent is an Ancestor which can create a VMInterface.
type VMInterfacesParent interface {
	VMInterfacesAncestor
	CreateVMInterface(*VMInterface) *bambou.Error
}

// VMInterface represents the model of a vminterface
type VMInterface struct {
	ID                  string        `json:"ID,omitempty"`
	ParentID            string        `json:"parentID,omitempty"`
	ParentType          string        `json:"parentType,omitempty"`
	Owner               string        `json:"owner,omitempty"`
	MAC                 string        `json:"MAC,omitempty"`
	VMUUID              string        `json:"VMUUID,omitempty"`
	IPAddress           string        `json:"IPAddress,omitempty"`
	VPortID             string        `json:"VPortID,omitempty"`
	VPortName           string        `json:"VPortName,omitempty"`
	IPv6Address         string        `json:"IPv6Address,omitempty"`
	IPv6Gateway         string        `json:"IPv6Gateway,omitempty"`
	Name                string        `json:"name,omitempty"`
	LastUpdatedBy       string        `json:"lastUpdatedBy,omitempty"`
	Gateway             string        `json:"gateway,omitempty"`
	Netmask             string        `json:"netmask,omitempty"`
	NetworkName         string        `json:"networkName,omitempty"`
	TierID              string        `json:"tierID,omitempty"`
	EmbeddedMetadata    []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope         string        `json:"entityScope,omitempty"`
	PolicyDecisionID    string        `json:"policyDecisionID,omitempty"`
	DomainID            string        `json:"domainID,omitempty"`
	DomainName          string        `json:"domainName,omitempty"`
	ZoneID              string        `json:"zoneID,omitempty"`
	ZoneName            string        `json:"zoneName,omitempty"`
	AttachedNetworkID   string        `json:"attachedNetworkID,omitempty"`
	AttachedNetworkType string        `json:"attachedNetworkType,omitempty"`
	MultiNICVPortName   string        `json:"multiNICVPortName,omitempty"`
	ExternalID          string        `json:"externalID,omitempty"`
}

// NewVMInterface returns a new *VMInterface
func NewVMInterface() *VMInterface {

	return &VMInterface{}
}

// Identity returns the Identity of the object.
func (o *VMInterface) Identity() bambou.Identity {

	return VMInterfaceIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VMInterface) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VMInterface) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VMInterface from the server
func (o *VMInterface) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VMInterface into the server
func (o *VMInterface) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VMInterface from the server
func (o *VMInterface) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// TCAs retrieves the list of child TCAs of the VMInterface
func (o *VMInterface) TCAs(info *bambou.FetchingInfo) (TCAsList, *bambou.Error) {

	var list TCAsList
	err := bambou.CurrentSession().FetchChildren(o, TCAIdentity, &list, info)
	return list, err
}

// RedirectionTargets retrieves the list of child RedirectionTargets of the VMInterface
func (o *VMInterface) RedirectionTargets(info *bambou.FetchingInfo) (RedirectionTargetsList, *bambou.Error) {

	var list RedirectionTargetsList
	err := bambou.CurrentSession().FetchChildren(o, RedirectionTargetIdentity, &list, info)
	return list, err
}

// Metadatas retrieves the list of child Metadatas of the VMInterface
func (o *VMInterface) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VMInterface
func (o *VMInterface) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DHCPOptions retrieves the list of child DHCPOptions of the VMInterface
func (o *VMInterface) DHCPOptions(info *bambou.FetchingInfo) (DHCPOptionsList, *bambou.Error) {

	var list DHCPOptionsList
	err := bambou.CurrentSession().FetchChildren(o, DHCPOptionIdentity, &list, info)
	return list, err
}

// DHCPv6Options retrieves the list of child DHCPv6Options of the VMInterface
func (o *VMInterface) DHCPv6Options(info *bambou.FetchingInfo) (DHCPv6OptionsList, *bambou.Error) {

	var list DHCPv6OptionsList
	err := bambou.CurrentSession().FetchChildren(o, DHCPv6OptionIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VMInterface
func (o *VMInterface) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VMInterface
func (o *VMInterface) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PolicyDecisions retrieves the list of child PolicyDecisions of the VMInterface
func (o *VMInterface) PolicyDecisions(info *bambou.FetchingInfo) (PolicyDecisionsList, *bambou.Error) {

	var list PolicyDecisionsList
	err := bambou.CurrentSession().FetchChildren(o, PolicyDecisionIdentity, &list, info)
	return list, err
}

// PolicyGroups retrieves the list of child PolicyGroups of the VMInterface
func (o *VMInterface) PolicyGroups(info *bambou.FetchingInfo) (PolicyGroupsList, *bambou.Error) {

	var list PolicyGroupsList
	err := bambou.CurrentSession().FetchChildren(o, PolicyGroupIdentity, &list, info)
	return list, err
}

// StaticRoutes retrieves the list of child StaticRoutes of the VMInterface
func (o *VMInterface) StaticRoutes(info *bambou.FetchingInfo) (StaticRoutesList, *bambou.Error) {

	var list StaticRoutesList
	err := bambou.CurrentSession().FetchChildren(o, StaticRouteIdentity, &list, info)
	return list, err
}

// Statistics retrieves the list of child Statistics of the VMInterface
func (o *VMInterface) Statistics(info *bambou.FetchingInfo) (StatisticsList, *bambou.Error) {

	var list StatisticsList
	err := bambou.CurrentSession().FetchChildren(o, StatisticsIdentity, &list, info)
	return list, err
}

// MultiCastChannelMaps retrieves the list of child MultiCastChannelMaps of the VMInterface
func (o *VMInterface) MultiCastChannelMaps(info *bambou.FetchingInfo) (MultiCastChannelMapsList, *bambou.Error) {

	var list MultiCastChannelMapsList
	err := bambou.CurrentSession().FetchChildren(o, MultiCastChannelMapIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the VMInterface
func (o *VMInterface) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
