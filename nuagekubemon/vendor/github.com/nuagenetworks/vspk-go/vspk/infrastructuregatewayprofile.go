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

// InfrastructureGatewayProfileIdentity represents the Identity of the object
var InfrastructureGatewayProfileIdentity = bambou.Identity{
	Name:     "infrastructuregatewayprofile",
	Category: "infrastructuregatewayprofiles",
}

// InfrastructureGatewayProfilesList represents a list of InfrastructureGatewayProfiles
type InfrastructureGatewayProfilesList []*InfrastructureGatewayProfile

// InfrastructureGatewayProfilesAncestor is the interface of an ancestor of a InfrastructureGatewayProfile must implement.
type InfrastructureGatewayProfilesAncestor interface {
	InfrastructureGatewayProfiles(*bambou.FetchingInfo) (InfrastructureGatewayProfilesList, *bambou.Error)
	CreateInfrastructureGatewayProfiles(*InfrastructureGatewayProfile) *bambou.Error
}

// InfrastructureGatewayProfile represents the model of a infrastructuregatewayprofile
type InfrastructureGatewayProfile struct {
	ID                           string `json:"ID,omitempty"`
	ParentID                     string `json:"parentID,omitempty"`
	ParentType                   string `json:"parentType,omitempty"`
	Owner                        string `json:"owner,omitempty"`
	NTPServerKey                 string `json:"NTPServerKey,omitempty"`
	NTPServerKeyID               int    `json:"NTPServerKeyID,omitempty"`
	Name                         string `json:"name,omitempty"`
	LastUpdatedBy                string `json:"lastUpdatedBy,omitempty"`
	DatapathSyncTimeout          int    `json:"datapathSyncTimeout,omitempty"`
	DeadTimer                    string `json:"deadTimer,omitempty"`
	DeadTimerEnabled             bool   `json:"deadTimerEnabled"`
	RemoteLogMode                string `json:"remoteLogMode,omitempty"`
	RemoteLogServerAddress       string `json:"remoteLogServerAddress,omitempty"`
	RemoteLogServerPort          int    `json:"remoteLogServerPort,omitempty"`
	Description                  string `json:"description,omitempty"`
	MetadataUpgradePath          string `json:"metadataUpgradePath,omitempty"`
	EnterpriseID                 string `json:"enterpriseID,omitempty"`
	EntityScope                  string `json:"entityScope,omitempty"`
	ControllerLessDuration       string `json:"controllerLessDuration,omitempty"`
	ControllerLessEnabled        bool   `json:"controllerLessEnabled"`
	ControllerLessForwardingMode string `json:"controllerLessForwardingMode,omitempty"`
	ControllerLessRemoteDuration string `json:"controllerLessRemoteDuration,omitempty"`
	ForceImmediateSystemSync     bool   `json:"forceImmediateSystemSync"`
	UpgradeAction                string `json:"upgradeAction,omitempty"`
	ProxyDNSName                 string `json:"proxyDNSName,omitempty"`
	UseTwoFactor                 bool   `json:"useTwoFactor"`
	StatsCollectorPort           int    `json:"statsCollectorPort,omitempty"`
	ExternalID                   string `json:"externalID,omitempty"`
	SystemSyncScheduler          string `json:"systemSyncScheduler,omitempty"`
}

// NewInfrastructureGatewayProfile returns a new *InfrastructureGatewayProfile
func NewInfrastructureGatewayProfile() *InfrastructureGatewayProfile {

	return &InfrastructureGatewayProfile{
		UpgradeAction:       "NONE",
		StatsCollectorPort:  29090,
		SystemSyncScheduler: "0 0 * * 0",
		DeadTimer:           "ONE_HOUR",
		UseTwoFactor:        true,
		RemoteLogMode:       "DISABLED",
		DatapathSyncTimeout: 1000,
	}
}

// Identity returns the Identity of the object.
func (o *InfrastructureGatewayProfile) Identity() bambou.Identity {

	return InfrastructureGatewayProfileIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *InfrastructureGatewayProfile) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *InfrastructureGatewayProfile) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the InfrastructureGatewayProfile from the server
func (o *InfrastructureGatewayProfile) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the InfrastructureGatewayProfile into the server
func (o *InfrastructureGatewayProfile) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the InfrastructureGatewayProfile from the server
func (o *InfrastructureGatewayProfile) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the InfrastructureGatewayProfile
func (o *InfrastructureGatewayProfile) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the InfrastructureGatewayProfile
func (o *InfrastructureGatewayProfile) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the InfrastructureGatewayProfile
func (o *InfrastructureGatewayProfile) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the InfrastructureGatewayProfile
func (o *InfrastructureGatewayProfile) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
