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

// SSIDConnectionIdentity represents the Identity of the object
var SSIDConnectionIdentity = bambou.Identity{
	Name:     "ssidconnection",
	Category: "ssidconnections",
}

// SSIDConnectionsList represents a list of SSIDConnections
type SSIDConnectionsList []*SSIDConnection

// SSIDConnectionsAncestor is the interface that an ancestor of a SSIDConnection must implement.
// An Ancestor is defined as an entity that has SSIDConnection as a descendant.
// An Ancestor can get a list of its child SSIDConnections, but not necessarily create one.
type SSIDConnectionsAncestor interface {
	SSIDConnections(*bambou.FetchingInfo) (SSIDConnectionsList, *bambou.Error)
}

// SSIDConnectionsParent is the interface that a parent of a SSIDConnection must implement.
// A Parent is defined as an entity that has SSIDConnection as a child.
// A Parent is an Ancestor which can create a SSIDConnection.
type SSIDConnectionsParent interface {
	SSIDConnectionsAncestor
	CreateSSIDConnection(*SSIDConnection) *bambou.Error
}

// SSIDConnection represents the model of a ssidconnection
type SSIDConnection struct {
	ID                               string        `json:"ID,omitempty"`
	ParentID                         string        `json:"parentID,omitempty"`
	ParentType                       string        `json:"parentType,omitempty"`
	Owner                            string        `json:"owner,omitempty"`
	Name                             string        `json:"name,omitempty"`
	Passphrase                       string        `json:"passphrase,omitempty"`
	LastUpdatedBy                    string        `json:"lastUpdatedBy,omitempty"`
	GatewayID                        string        `json:"gatewayID,omitempty"`
	Readonly                         bool          `json:"readonly"`
	RedirectOption                   string        `json:"redirectOption,omitempty"`
	RedirectURL                      string        `json:"redirectURL,omitempty"`
	GenericConfig                    string        `json:"genericConfig,omitempty"`
	PermittedAction                  string        `json:"permittedAction,omitempty"`
	Description                      string        `json:"description,omitempty"`
	Restricted                       bool          `json:"restricted"`
	WhiteList                        []interface{} `json:"whiteList,omitempty"`
	BlackList                        []interface{} `json:"blackList,omitempty"`
	VlanID                           int           `json:"vlanID,omitempty"`
	EmbeddedMetadata                 []interface{} `json:"embeddedMetadata,omitempty"`
	InterfaceName                    string        `json:"interfaceName,omitempty"`
	EntityScope                      string        `json:"entityScope,omitempty"`
	VportID                          string        `json:"vportID,omitempty"`
	BroadcastSSID                    bool          `json:"broadcastSSID"`
	AssociatedCaptivePortalProfileID string        `json:"associatedCaptivePortalProfileID,omitempty"`
	AssociatedEgressQOSPolicyID      string        `json:"associatedEgressQOSPolicyID,omitempty"`
	Status                           string        `json:"status,omitempty"`
	AuthenticationMode               string        `json:"authenticationMode,omitempty"`
	ExternalID                       string        `json:"externalID,omitempty"`
}

// NewSSIDConnection returns a new *SSIDConnection
func NewSSIDConnection() *SSIDConnection {

	return &SSIDConnection{
		RedirectOption:     "ORIGINAL_REQUEST",
		BroadcastSSID:      true,
		AuthenticationMode: "WPA2",
	}
}

// Identity returns the Identity of the object.
func (o *SSIDConnection) Identity() bambou.Identity {

	return SSIDConnectionIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *SSIDConnection) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *SSIDConnection) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the SSIDConnection from the server
func (o *SSIDConnection) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the SSIDConnection into the server
func (o *SSIDConnection) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the SSIDConnection from the server
func (o *SSIDConnection) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the SSIDConnection
func (o *SSIDConnection) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the SSIDConnection
func (o *SSIDConnection) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Alarms retrieves the list of child Alarms of the SSIDConnection
func (o *SSIDConnection) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the SSIDConnection
func (o *SSIDConnection) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the SSIDConnection
func (o *SSIDConnection) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the SSIDConnection
func (o *SSIDConnection) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
