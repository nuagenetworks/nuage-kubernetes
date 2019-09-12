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

// BridgeInterfaceIdentity represents the Identity of the object
var BridgeInterfaceIdentity = bambou.Identity{
	Name:     "bridgeinterface",
	Category: "bridgeinterfaces",
}

// BridgeInterfacesList represents a list of BridgeInterfaces
type BridgeInterfacesList []*BridgeInterface

// BridgeInterfacesAncestor is the interface that an ancestor of a BridgeInterface must implement.
// An Ancestor is defined as an entity that has BridgeInterface as a descendant.
// An Ancestor can get a list of its child BridgeInterfaces, but not necessarily create one.
type BridgeInterfacesAncestor interface {
	BridgeInterfaces(*bambou.FetchingInfo) (BridgeInterfacesList, *bambou.Error)
}

// BridgeInterfacesParent is the interface that a parent of a BridgeInterface must implement.
// A Parent is defined as an entity that has BridgeInterface as a child.
// A Parent is an Ancestor which can create a BridgeInterface.
type BridgeInterfacesParent interface {
	BridgeInterfacesAncestor
	CreateBridgeInterface(*BridgeInterface) *bambou.Error
}

// BridgeInterface represents the model of a bridgeinterface
type BridgeInterface struct {
	ID                          string        `json:"ID,omitempty"`
	ParentID                    string        `json:"parentID,omitempty"`
	ParentType                  string        `json:"parentType,omitempty"`
	Owner                       string        `json:"owner,omitempty"`
	VPortID                     string        `json:"VPortID,omitempty"`
	VPortName                   string        `json:"VPortName,omitempty"`
	IPv6Gateway                 string        `json:"IPv6Gateway,omitempty"`
	Name                        string        `json:"name,omitempty"`
	LastUpdatedBy               string        `json:"lastUpdatedBy,omitempty"`
	Gateway                     string        `json:"gateway,omitempty"`
	Netmask                     string        `json:"netmask,omitempty"`
	NetworkName                 string        `json:"networkName,omitempty"`
	TierID                      string        `json:"tierID,omitempty"`
	EmbeddedMetadata            []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope                 string        `json:"entityScope,omitempty"`
	PolicyDecisionID            string        `json:"policyDecisionID,omitempty"`
	DomainID                    string        `json:"domainID,omitempty"`
	DomainName                  string        `json:"domainName,omitempty"`
	ZoneID                      string        `json:"zoneID,omitempty"`
	ZoneName                    string        `json:"zoneName,omitempty"`
	AssociatedFloatingIPAddress string        `json:"associatedFloatingIPAddress,omitempty"`
	AttachedNetworkID           string        `json:"attachedNetworkID,omitempty"`
	AttachedNetworkType         string        `json:"attachedNetworkType,omitempty"`
	ExternalID                  string        `json:"externalID,omitempty"`
}

// NewBridgeInterface returns a new *BridgeInterface
func NewBridgeInterface() *BridgeInterface {

	return &BridgeInterface{}
}

// Identity returns the Identity of the object.
func (o *BridgeInterface) Identity() bambou.Identity {

	return BridgeInterfaceIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *BridgeInterface) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *BridgeInterface) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the BridgeInterface from the server
func (o *BridgeInterface) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the BridgeInterface into the server
func (o *BridgeInterface) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the BridgeInterface from the server
func (o *BridgeInterface) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// TCAs retrieves the list of child TCAs of the BridgeInterface
func (o *BridgeInterface) TCAs(info *bambou.FetchingInfo) (TCAsList, *bambou.Error) {

	var list TCAsList
	err := bambou.CurrentSession().FetchChildren(o, TCAIdentity, &list, info)
	return list, err
}

// RedirectionTargets retrieves the list of child RedirectionTargets of the BridgeInterface
func (o *BridgeInterface) RedirectionTargets(info *bambou.FetchingInfo) (RedirectionTargetsList, *bambou.Error) {

	var list RedirectionTargetsList
	err := bambou.CurrentSession().FetchChildren(o, RedirectionTargetIdentity, &list, info)
	return list, err
}

// DeploymentFailures retrieves the list of child DeploymentFailures of the BridgeInterface
func (o *BridgeInterface) DeploymentFailures(info *bambou.FetchingInfo) (DeploymentFailuresList, *bambou.Error) {

	var list DeploymentFailuresList
	err := bambou.CurrentSession().FetchChildren(o, DeploymentFailureIdentity, &list, info)
	return list, err
}

// CreateDeploymentFailure creates a new child DeploymentFailure under the BridgeInterface
func (o *BridgeInterface) CreateDeploymentFailure(child *DeploymentFailure) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the BridgeInterface
func (o *BridgeInterface) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the BridgeInterface
func (o *BridgeInterface) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// DHCPOptions retrieves the list of child DHCPOptions of the BridgeInterface
func (o *BridgeInterface) DHCPOptions(info *bambou.FetchingInfo) (DHCPOptionsList, *bambou.Error) {

	var list DHCPOptionsList
	err := bambou.CurrentSession().FetchChildren(o, DHCPOptionIdentity, &list, info)
	return list, err
}

// DHCPv6Options retrieves the list of child DHCPv6Options of the BridgeInterface
func (o *BridgeInterface) DHCPv6Options(info *bambou.FetchingInfo) (DHCPv6OptionsList, *bambou.Error) {

	var list DHCPv6OptionsList
	err := bambou.CurrentSession().FetchChildren(o, DHCPv6OptionIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the BridgeInterface
func (o *BridgeInterface) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the BridgeInterface
func (o *BridgeInterface) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// PolicyDecisions retrieves the list of child PolicyDecisions of the BridgeInterface
func (o *BridgeInterface) PolicyDecisions(info *bambou.FetchingInfo) (PolicyDecisionsList, *bambou.Error) {

	var list PolicyDecisionsList
	err := bambou.CurrentSession().FetchChildren(o, PolicyDecisionIdentity, &list, info)
	return list, err
}

// PolicyGroups retrieves the list of child PolicyGroups of the BridgeInterface
func (o *BridgeInterface) PolicyGroups(info *bambou.FetchingInfo) (PolicyGroupsList, *bambou.Error) {

	var list PolicyGroupsList
	err := bambou.CurrentSession().FetchChildren(o, PolicyGroupIdentity, &list, info)
	return list, err
}

// QOSs retrieves the list of child QOSs of the BridgeInterface
func (o *BridgeInterface) QOSs(info *bambou.FetchingInfo) (QOSsList, *bambou.Error) {

	var list QOSsList
	err := bambou.CurrentSession().FetchChildren(o, QOSIdentity, &list, info)
	return list, err
}

// CreateQOS creates a new child QOS under the BridgeInterface
func (o *BridgeInterface) CreateQOS(child *QOS) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Statistics retrieves the list of child Statistics of the BridgeInterface
func (o *BridgeInterface) Statistics(info *bambou.FetchingInfo) (StatisticsList, *bambou.Error) {

	var list StatisticsList
	err := bambou.CurrentSession().FetchChildren(o, StatisticsIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the BridgeInterface
func (o *BridgeInterface) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
