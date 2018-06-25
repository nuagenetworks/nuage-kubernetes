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
	RedirectOption                   string        `json:"redirectOption,omitempty"`
	RedirectURL                      string        `json:"redirectURL,omitempty"`
	GenericConfig                    string        `json:"genericConfig,omitempty"`
	Description                      string        `json:"description,omitempty"`
	WhiteList                        []interface{} `json:"whiteList,omitempty"`
	BlackList                        []interface{} `json:"blackList,omitempty"`
	InterfaceName                    string        `json:"interfaceName,omitempty"`
	VportID                          string        `json:"vportID,omitempty"`
	BroadcastSSID                    bool          `json:"broadcastSSID"`
	AssociatedCaptivePortalProfileID string        `json:"associatedCaptivePortalProfileID,omitempty"`
	AssociatedEgressQOSPolicyID      string        `json:"associatedEgressQOSPolicyID,omitempty"`
	AuthenticationMode               string        `json:"authenticationMode,omitempty"`
}

// NewSSIDConnection returns a new *SSIDConnection
func NewSSIDConnection() *SSIDConnection {

	return &SSIDConnection{
		RedirectOption:     "ORIGINAL_REQUEST",
		BroadcastSSID:      true,
		AuthenticationMode: "OPEN",
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

// CaptivePortalProfiles retrieves the list of child CaptivePortalProfiles of the SSIDConnection
func (o *SSIDConnection) CaptivePortalProfiles(info *bambou.FetchingInfo) (CaptivePortalProfilesList, *bambou.Error) {

	var list CaptivePortalProfilesList
	err := bambou.CurrentSession().FetchChildren(o, CaptivePortalProfileIdentity, &list, info)
	return list, err
}

// AssignCaptivePortalProfiles assigns the list of CaptivePortalProfiles to the SSIDConnection
func (o *SSIDConnection) AssignCaptivePortalProfiles(children CaptivePortalProfilesList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, CaptivePortalProfileIdentity)
}

// Alarms retrieves the list of child Alarms of the SSIDConnection
func (o *SSIDConnection) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the SSIDConnection
func (o *SSIDConnection) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
