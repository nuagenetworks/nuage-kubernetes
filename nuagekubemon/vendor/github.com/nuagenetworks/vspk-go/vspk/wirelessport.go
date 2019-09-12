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

// WirelessPortIdentity represents the Identity of the object
var WirelessPortIdentity = bambou.Identity{
	Name:     "wirelessport",
	Category: "wirelessports",
}

// WirelessPortsList represents a list of WirelessPorts
type WirelessPortsList []*WirelessPort

// WirelessPortsAncestor is the interface that an ancestor of a WirelessPort must implement.
// An Ancestor is defined as an entity that has WirelessPort as a descendant.
// An Ancestor can get a list of its child WirelessPorts, but not necessarily create one.
type WirelessPortsAncestor interface {
	WirelessPorts(*bambou.FetchingInfo) (WirelessPortsList, *bambou.Error)
}

// WirelessPortsParent is the interface that a parent of a WirelessPort must implement.
// A Parent is defined as an entity that has WirelessPort as a child.
// A Parent is an Ancestor which can create a WirelessPort.
type WirelessPortsParent interface {
	WirelessPortsAncestor
	CreateWirelessPort(*WirelessPort) *bambou.Error
}

// WirelessPort represents the model of a wirelessport
type WirelessPort struct {
	ID                          string        `json:"ID,omitempty"`
	ParentID                    string        `json:"parentID,omitempty"`
	ParentType                  string        `json:"parentType,omitempty"`
	Owner                       string        `json:"owner,omitempty"`
	VLANRange                   string        `json:"VLANRange,omitempty"`
	Name                        string        `json:"name,omitempty"`
	LastUpdatedBy               string        `json:"lastUpdatedBy,omitempty"`
	GenericConfig               string        `json:"genericConfig,omitempty"`
	PermittedAction             string        `json:"permittedAction,omitempty"`
	Description                 string        `json:"description,omitempty"`
	ChannelWidth                string        `json:"channelWidth,omitempty"`
	PhysicalName                string        `json:"physicalName,omitempty"`
	WifiFrequencyBand           string        `json:"wifiFrequencyBand,omitempty"`
	WifiMode                    string        `json:"wifiMode,omitempty"`
	EmbeddedMetadata            []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope                 string        `json:"entityScope,omitempty"`
	PortType                    string        `json:"portType,omitempty"`
	CountryCode                 string        `json:"countryCode,omitempty"`
	FrequencyChannel            string        `json:"frequencyChannel,omitempty"`
	UseUserMnemonic             bool          `json:"useUserMnemonic"`
	UserMnemonic                string        `json:"userMnemonic,omitempty"`
	AssociatedEgressQOSPolicyID string        `json:"associatedEgressQOSPolicyID,omitempty"`
	Status                      string        `json:"status,omitempty"`
	ExternalID                  string        `json:"externalID,omitempty"`
}

// NewWirelessPort returns a new *WirelessPort
func NewWirelessPort() *WirelessPort {

	return &WirelessPort{
		ChannelWidth:      "WIDTH_20_MHZ",
		WifiFrequencyBand: "FREQ_2_4_GHZ",
		WifiMode:          "WIFI_B_G_N",
		PortType:          "ACCESS",
		FrequencyChannel:  "CH_0",
	}
}

// Identity returns the Identity of the object.
func (o *WirelessPort) Identity() bambou.Identity {

	return WirelessPortIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *WirelessPort) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *WirelessPort) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the WirelessPort from the server
func (o *WirelessPort) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the WirelessPort into the server
func (o *WirelessPort) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the WirelessPort from the server
func (o *WirelessPort) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the WirelessPort
func (o *WirelessPort) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the WirelessPort
func (o *WirelessPort) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Alarms retrieves the list of child Alarms of the WirelessPort
func (o *WirelessPort) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the WirelessPort
func (o *WirelessPort) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the WirelessPort
func (o *WirelessPort) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// SSIDConnections retrieves the list of child SSIDConnections of the WirelessPort
func (o *WirelessPort) SSIDConnections(info *bambou.FetchingInfo) (SSIDConnectionsList, *bambou.Error) {

	var list SSIDConnectionsList
	err := bambou.CurrentSession().FetchChildren(o, SSIDConnectionIdentity, &list, info)
	return list, err
}

// CreateSSIDConnection creates a new child SSIDConnection under the WirelessPort
func (o *WirelessPort) CreateSSIDConnection(child *SSIDConnection) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the WirelessPort
func (o *WirelessPort) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
