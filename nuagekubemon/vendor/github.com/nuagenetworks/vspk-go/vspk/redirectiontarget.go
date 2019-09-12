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

// RedirectionTargetIdentity represents the Identity of the object
var RedirectionTargetIdentity = bambou.Identity{
	Name:     "redirectiontarget",
	Category: "redirectiontargets",
}

// RedirectionTargetsList represents a list of RedirectionTargets
type RedirectionTargetsList []*RedirectionTarget

// RedirectionTargetsAncestor is the interface that an ancestor of a RedirectionTarget must implement.
// An Ancestor is defined as an entity that has RedirectionTarget as a descendant.
// An Ancestor can get a list of its child RedirectionTargets, but not necessarily create one.
type RedirectionTargetsAncestor interface {
	RedirectionTargets(*bambou.FetchingInfo) (RedirectionTargetsList, *bambou.Error)
}

// RedirectionTargetsParent is the interface that a parent of a RedirectionTarget must implement.
// A Parent is defined as an entity that has RedirectionTarget as a child.
// A Parent is an Ancestor which can create a RedirectionTarget.
type RedirectionTargetsParent interface {
	RedirectionTargetsAncestor
	CreateRedirectionTarget(*RedirectionTarget) *bambou.Error
}

// RedirectionTarget represents the model of a redirectiontarget
type RedirectionTarget struct {
	ID                string        `json:"ID,omitempty"`
	ParentID          string        `json:"parentID,omitempty"`
	ParentType        string        `json:"parentType,omitempty"`
	Owner             string        `json:"owner,omitempty"`
	ESI               string        `json:"ESI,omitempty"`
	Name              string        `json:"name,omitempty"`
	LastUpdatedBy     string        `json:"lastUpdatedBy,omitempty"`
	RedundancyEnabled bool          `json:"redundancyEnabled"`
	TemplateID        string        `json:"templateID,omitempty"`
	Description       string        `json:"description,omitempty"`
	DestinationType   string        `json:"destinationType,omitempty"`
	VirtualNetworkID  string        `json:"virtualNetworkID,omitempty"`
	EmbeddedMetadata  []interface{} `json:"embeddedMetadata,omitempty"`
	EndPointType      string        `json:"endPointType,omitempty"`
	EntityScope       string        `json:"entityScope,omitempty"`
	TriggerType       string        `json:"triggerType,omitempty"`
	ExternalID        string        `json:"externalID,omitempty"`
}

// NewRedirectionTarget returns a new *RedirectionTarget
func NewRedirectionTarget() *RedirectionTarget {

	return &RedirectionTarget{
		EndPointType: "L3",
	}
}

// Identity returns the Identity of the object.
func (o *RedirectionTarget) Identity() bambou.Identity {

	return RedirectionTargetIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *RedirectionTarget) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *RedirectionTarget) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the RedirectionTarget from the server
func (o *RedirectionTarget) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the RedirectionTarget into the server
func (o *RedirectionTarget) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the RedirectionTarget from the server
func (o *RedirectionTarget) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the RedirectionTarget
func (o *RedirectionTarget) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the RedirectionTarget
func (o *RedirectionTarget) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VirtualIPs retrieves the list of child VirtualIPs of the RedirectionTarget
func (o *RedirectionTarget) VirtualIPs(info *bambou.FetchingInfo) (VirtualIPsList, *bambou.Error) {

	var list VirtualIPsList
	err := bambou.CurrentSession().FetchChildren(o, VirtualIPIdentity, &list, info)
	return list, err
}

// CreateVirtualIP creates a new child VirtualIP under the RedirectionTarget
func (o *RedirectionTarget) CreateVirtualIP(child *VirtualIP) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the RedirectionTarget
func (o *RedirectionTarget) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the RedirectionTarget
func (o *RedirectionTarget) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VPorts retrieves the list of child VPorts of the RedirectionTarget
func (o *RedirectionTarget) VPorts(info *bambou.FetchingInfo) (VPortsList, *bambou.Error) {

	var list VPortsList
	err := bambou.CurrentSession().FetchChildren(o, VPortIdentity, &list, info)
	return list, err
}

// AssignVPorts assigns the list of VPorts to the RedirectionTarget
func (o *RedirectionTarget) AssignVPorts(children VPortsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, VPortIdentity)
}

// EventLogs retrieves the list of child EventLogs of the RedirectionTarget
func (o *RedirectionTarget) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
