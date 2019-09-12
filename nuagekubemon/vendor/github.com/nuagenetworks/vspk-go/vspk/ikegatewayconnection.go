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

// IKEGatewayConnectionIdentity represents the Identity of the object
var IKEGatewayConnectionIdentity = bambou.Identity{
	Name:     "ikegatewayconnection",
	Category: "ikegatewayconnections",
}

// IKEGatewayConnectionsList represents a list of IKEGatewayConnections
type IKEGatewayConnectionsList []*IKEGatewayConnection

// IKEGatewayConnectionsAncestor is the interface that an ancestor of a IKEGatewayConnection must implement.
// An Ancestor is defined as an entity that has IKEGatewayConnection as a descendant.
// An Ancestor can get a list of its child IKEGatewayConnections, but not necessarily create one.
type IKEGatewayConnectionsAncestor interface {
	IKEGatewayConnections(*bambou.FetchingInfo) (IKEGatewayConnectionsList, *bambou.Error)
}

// IKEGatewayConnectionsParent is the interface that a parent of a IKEGatewayConnection must implement.
// A Parent is defined as an entity that has IKEGatewayConnection as a child.
// A Parent is an Ancestor which can create a IKEGatewayConnection.
type IKEGatewayConnectionsParent interface {
	IKEGatewayConnectionsAncestor
	CreateIKEGatewayConnection(*IKEGatewayConnection) *bambou.Error
}

// IKEGatewayConnection represents the model of a ikegatewayconnection
type IKEGatewayConnection struct {
	ID                               string        `json:"ID,omitempty"`
	ParentID                         string        `json:"parentID,omitempty"`
	ParentType                       string        `json:"parentType,omitempty"`
	Owner                            string        `json:"owner,omitempty"`
	NSGIdentifier                    string        `json:"NSGIdentifier,omitempty"`
	NSGIdentifierType                string        `json:"NSGIdentifierType,omitempty"`
	NSGRole                          string        `json:"NSGRole,omitempty"`
	Name                             string        `json:"name,omitempty"`
	Mark                             int           `json:"mark,omitempty"`
	LastUpdatedBy                    string        `json:"lastUpdatedBy,omitempty"`
	Sequence                         int           `json:"sequence,omitempty"`
	AllowAnySubnet                   bool          `json:"allowAnySubnet"`
	EmbeddedMetadata                 []interface{} `json:"embeddedMetadata,omitempty"`
	UnencryptedPSK                   string        `json:"unencryptedPSK,omitempty"`
	EntityScope                      string        `json:"entityScope,omitempty"`
	ConfigurationStatus              string        `json:"configurationStatus,omitempty"`
	PortVLANName                     string        `json:"portVLANName,omitempty"`
	Priority                         int           `json:"priority,omitempty"`
	AssociatedCloudID                string        `json:"associatedCloudID,omitempty"`
	AssociatedCloudType              string        `json:"associatedCloudType,omitempty"`
	AssociatedIKEAuthenticationID    string        `json:"associatedIKEAuthenticationID,omitempty"`
	AssociatedIKEAuthenticationType  string        `json:"associatedIKEAuthenticationType,omitempty"`
	AssociatedIKEEncryptionProfileID string        `json:"associatedIKEEncryptionProfileID,omitempty"`
	AssociatedIKEGatewayProfileID    string        `json:"associatedIKEGatewayProfileID,omitempty"`
	AssociatedVLANID                 string        `json:"associatedVLANID,omitempty"`
	ExternalID                       string        `json:"externalID,omitempty"`
}

// NewIKEGatewayConnection returns a new *IKEGatewayConnection
func NewIKEGatewayConnection() *IKEGatewayConnection {

	return &IKEGatewayConnection{
		NSGIdentifierType: "ID_KEY_ID",
		Mark:              1,
	}
}

// Identity returns the Identity of the object.
func (o *IKEGatewayConnection) Identity() bambou.Identity {

	return IKEGatewayConnectionIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *IKEGatewayConnection) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *IKEGatewayConnection) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the IKEGatewayConnection from the server
func (o *IKEGatewayConnection) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the IKEGatewayConnection into the server
func (o *IKEGatewayConnection) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the IKEGatewayConnection from the server
func (o *IKEGatewayConnection) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// PerformanceMonitors retrieves the list of child PerformanceMonitors of the IKEGatewayConnection
func (o *IKEGatewayConnection) PerformanceMonitors(info *bambou.FetchingInfo) (PerformanceMonitorsList, *bambou.Error) {

	var list PerformanceMonitorsList
	err := bambou.CurrentSession().FetchChildren(o, PerformanceMonitorIdentity, &list, info)
	return list, err
}

// AssignPerformanceMonitors assigns the list of PerformanceMonitors to the IKEGatewayConnection
func (o *IKEGatewayConnection) AssignPerformanceMonitors(children PerformanceMonitorsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, PerformanceMonitorIdentity)
}

// Metadatas retrieves the list of child Metadatas of the IKEGatewayConnection
func (o *IKEGatewayConnection) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the IKEGatewayConnection
func (o *IKEGatewayConnection) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Alarms retrieves the list of child Alarms of the IKEGatewayConnection
func (o *IKEGatewayConnection) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the IKEGatewayConnection
func (o *IKEGatewayConnection) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the IKEGatewayConnection
func (o *IKEGatewayConnection) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Jobs retrieves the list of child Jobs of the IKEGatewayConnection
func (o *IKEGatewayConnection) Jobs(info *bambou.FetchingInfo) (JobsList, *bambou.Error) {

	var list JobsList
	err := bambou.CurrentSession().FetchChildren(o, JobIdentity, &list, info)
	return list, err
}

// CreateJob creates a new child Job under the IKEGatewayConnection
func (o *IKEGatewayConnection) CreateJob(child *Job) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Subnets retrieves the list of child Subnets of the IKEGatewayConnection
func (o *IKEGatewayConnection) Subnets(info *bambou.FetchingInfo) (SubnetsList, *bambou.Error) {

	var list SubnetsList
	err := bambou.CurrentSession().FetchChildren(o, SubnetIdentity, &list, info)
	return list, err
}

// AssignSubnets assigns the list of Subnets to the IKEGatewayConnection
func (o *IKEGatewayConnection) AssignSubnets(children SubnetsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, SubnetIdentity)
}
