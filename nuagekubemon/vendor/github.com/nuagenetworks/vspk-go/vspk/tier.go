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

// TierIdentity represents the Identity of the object
var TierIdentity = bambou.Identity{
	Name:     "tier",
	Category: "tiers",
}

// TiersList represents a list of Tiers
type TiersList []*Tier

// TiersAncestor is the interface that an ancestor of a Tier must implement.
// An Ancestor is defined as an entity that has Tier as a descendant.
// An Ancestor can get a list of its child Tiers, but not necessarily create one.
type TiersAncestor interface {
	Tiers(*bambou.FetchingInfo) (TiersList, *bambou.Error)
}

// TiersParent is the interface that a parent of a Tier must implement.
// A Parent is defined as an entity that has Tier as a child.
// A Parent is an Ancestor which can create a Tier.
type TiersParent interface {
	TiersAncestor
	CreateTier(*Tier) *bambou.Error
}

// Tier represents the model of a tier
type Tier struct {
	ID                          string `json:"ID,omitempty"`
	ParentID                    string `json:"parentID,omitempty"`
	ParentType                  string `json:"parentType,omitempty"`
	Owner                       string `json:"owner,omitempty"`
	Name                        string `json:"name,omitempty"`
	LastUpdatedBy               string `json:"lastUpdatedBy,omitempty"`
	Gateway                     string `json:"gateway,omitempty"`
	Address                     string `json:"address,omitempty"`
	Description                 string `json:"description,omitempty"`
	Metadata                    string `json:"metadata,omitempty"`
	Netmask                     string `json:"netmask,omitempty"`
	EntityScope                 string `json:"entityScope,omitempty"`
	AssociatedApplicationID     string `json:"associatedApplicationID,omitempty"`
	AssociatedFloatingIPPoolID  string `json:"associatedFloatingIPPoolID,omitempty"`
	AssociatedNetworkMacroID    string `json:"associatedNetworkMacroID,omitempty"`
	AssociatedNetworkObjectID   string `json:"associatedNetworkObjectID,omitempty"`
	AssociatedNetworkObjectType string `json:"associatedNetworkObjectType,omitempty"`
	ExternalID                  string `json:"externalID,omitempty"`
	Type                        string `json:"type,omitempty"`
}

// NewTier returns a new *Tier
func NewTier() *Tier {

	return &Tier{
		Type: "STANDARD",
	}
}

// Identity returns the Identity of the object.
func (o *Tier) Identity() bambou.Identity {

	return TierIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Tier) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Tier) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Tier from the server
func (o *Tier) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Tier into the server
func (o *Tier) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Tier from the server
func (o *Tier) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// TCAs retrieves the list of child TCAs of the Tier
func (o *Tier) TCAs(info *bambou.FetchingInfo) (TCAsList, *bambou.Error) {

	var list TCAsList
	err := bambou.CurrentSession().FetchChildren(o, TCAIdentity, &list, info)
	return list, err
}

// Metadatas retrieves the list of child Metadatas of the Tier
func (o *Tier) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the Tier
func (o *Tier) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Tier
func (o *Tier) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Tier
func (o *Tier) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMs retrieves the list of child VMs of the Tier
func (o *Tier) VMs(info *bambou.FetchingInfo) (VMsList, *bambou.Error) {

	var list VMsList
	err := bambou.CurrentSession().FetchChildren(o, VMIdentity, &list, info)
	return list, err
}

// Containers retrieves the list of child Containers of the Tier
func (o *Tier) Containers(info *bambou.FetchingInfo) (ContainersList, *bambou.Error) {

	var list ContainersList
	err := bambou.CurrentSession().FetchChildren(o, ContainerIdentity, &list, info)
	return list, err
}

// VPorts retrieves the list of child VPorts of the Tier
func (o *Tier) VPorts(info *bambou.FetchingInfo) (VPortsList, *bambou.Error) {

	var list VPortsList
	err := bambou.CurrentSession().FetchChildren(o, VPortIdentity, &list, info)
	return list, err
}

// Statistics retrieves the list of child Statistics of the Tier
func (o *Tier) Statistics(info *bambou.FetchingInfo) (StatisticsList, *bambou.Error) {

	var list StatisticsList
	err := bambou.CurrentSession().FetchChildren(o, StatisticsIdentity, &list, info)
	return list, err
}

// StatisticsPolicies retrieves the list of child StatisticsPolicies of the Tier
func (o *Tier) StatisticsPolicies(info *bambou.FetchingInfo) (StatisticsPoliciesList, *bambou.Error) {

	var list StatisticsPoliciesList
	err := bambou.CurrentSession().FetchChildren(o, StatisticsPolicyIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the Tier
func (o *Tier) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
