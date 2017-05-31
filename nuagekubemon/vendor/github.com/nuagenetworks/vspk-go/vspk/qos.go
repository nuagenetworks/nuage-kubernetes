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

// QOSIdentity represents the Identity of the object
var QOSIdentity = bambou.Identity{
	Name:     "qos",
	Category: "qos",
}

// QOSsList represents a list of QOSs
type QOSsList []*QOS

// QOSsAncestor is the interface that an ancestor of a QOS must implement.
// An Ancestor is defined as an entity that has QOS as a descendant.
// An Ancestor can get a list of its child QOSs, but not necessarily create one.
type QOSsAncestor interface {
	QOSs(*bambou.FetchingInfo) (QOSsList, *bambou.Error)
}

// QOSsParent is the interface that a parent of a QOS must implement.
// A Parent is defined as an entity that has QOS as a child.
// A Parent is an Ancestor which can create a QOS.
type QOSsParent interface {
	QOSsAncestor
	CreateQOS(*QOS) *bambou.Error
}

// QOS represents the model of a qos
type QOS struct {
	ID                                     string `json:"ID,omitempty"`
	ParentID                               string `json:"parentID,omitempty"`
	ParentType                             string `json:"parentType,omitempty"`
	Owner                                  string `json:"owner,omitempty"`
	FIPCommittedBurstSize                  string `json:"FIPCommittedBurstSize,omitempty"`
	FIPCommittedInformationRate            string `json:"FIPCommittedInformationRate,omitempty"`
	FIPPeakBurstSize                       string `json:"FIPPeakBurstSize,omitempty"`
	FIPPeakInformationRate                 string `json:"FIPPeakInformationRate,omitempty"`
	FIPRateLimitingActive                  bool   `json:"FIPRateLimitingActive"`
	BUMCommittedBurstSize                  string `json:"BUMCommittedBurstSize,omitempty"`
	BUMCommittedInformationRate            string `json:"BUMCommittedInformationRate,omitempty"`
	BUMPeakBurstSize                       string `json:"BUMPeakBurstSize,omitempty"`
	BUMPeakInformationRate                 string `json:"BUMPeakInformationRate,omitempty"`
	BUMRateLimitingActive                  bool   `json:"BUMRateLimitingActive"`
	Name                                   string `json:"name,omitempty"`
	LastUpdatedBy                          string `json:"lastUpdatedBy,omitempty"`
	RateLimitingActive                     bool   `json:"rateLimitingActive"`
	Active                                 bool   `json:"active"`
	Peak                                   string `json:"peak,omitempty"`
	ServiceClass                           string `json:"serviceClass,omitempty"`
	Description                            string `json:"description,omitempty"`
	RewriteForwardingClass                 bool   `json:"rewriteForwardingClass"`
	EgressFIPCommittedBurstSize            string `json:"EgressFIPCommittedBurstSize,omitempty"`
	EgressFIPCommittedInformationRate      string `json:"EgressFIPCommittedInformationRate,omitempty"`
	EgressFIPPeakBurstSize                 string `json:"EgressFIPPeakBurstSize,omitempty"`
	EgressFIPPeakInformationRate           string `json:"EgressFIPPeakInformationRate,omitempty"`
	EntityScope                            string `json:"entityScope,omitempty"`
	CommittedBurstSize                     string `json:"committedBurstSize,omitempty"`
	CommittedInformationRate               string `json:"committedInformationRate,omitempty"`
	TrustedForwardingClass                 bool   `json:"trustedForwardingClass"`
	AssocQosId                             string `json:"assocQosId,omitempty"`
	AssociatedDSCPForwardingClassTableID   string `json:"associatedDSCPForwardingClassTableID,omitempty"`
	AssociatedDSCPForwardingClassTableName string `json:"associatedDSCPForwardingClassTableName,omitempty"`
	Burst                                  string `json:"burst,omitempty"`
	ExternalID                             string `json:"externalID,omitempty"`
}

// NewQOS returns a new *QOS
func NewQOS() *QOS {

	return &QOS{}
}

// Identity returns the Identity of the object.
func (o *QOS) Identity() bambou.Identity {

	return QOSIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *QOS) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *QOS) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the QOS from the server
func (o *QOS) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the QOS into the server
func (o *QOS) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the QOS from the server
func (o *QOS) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the QOS
func (o *QOS) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the QOS
func (o *QOS) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the QOS
func (o *QOS) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the QOS
func (o *QOS) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VMs retrieves the list of child VMs of the QOS
func (o *QOS) VMs(info *bambou.FetchingInfo) (VMsList, *bambou.Error) {

	var list VMsList
	err := bambou.CurrentSession().FetchChildren(o, VMIdentity, &list, info)
	return list, err
}

// Containers retrieves the list of child Containers of the QOS
func (o *QOS) Containers(info *bambou.FetchingInfo) (ContainersList, *bambou.Error) {

	var list ContainersList
	err := bambou.CurrentSession().FetchChildren(o, ContainerIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the QOS
func (o *QOS) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
