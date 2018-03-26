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

// EventLogIdentity represents the Identity of the object
var EventLogIdentity = bambou.Identity{
	Name:     "eventlog",
	Category: "eventlogs",
}

// EventLogsList represents a list of EventLogs
type EventLogsList []*EventLog

// EventLogsAncestor is the interface that an ancestor of a EventLog must implement.
// An Ancestor is defined as an entity that has EventLog as a descendant.
// An Ancestor can get a list of its child EventLogs, but not necessarily create one.
type EventLogsAncestor interface {
	EventLogs(*bambou.FetchingInfo) (EventLogsList, *bambou.Error)
}

// EventLogsParent is the interface that a parent of a EventLog must implement.
// A Parent is defined as an entity that has EventLog as a child.
// A Parent is an Ancestor which can create a EventLog.
type EventLogsParent interface {
	EventLogsAncestor
	CreateEventLog(*EventLog) *bambou.Error
}

// EventLog represents the model of a eventlog
type EventLog struct {
	ID                string        `json:"ID,omitempty"`
	ParentID          string        `json:"parentID,omitempty"`
	ParentType        string        `json:"parentType,omitempty"`
	Owner             string        `json:"owner,omitempty"`
	Diff              interface{}   `json:"diff,omitempty"`
	Enterprise        string        `json:"enterprise,omitempty"`
	Entities          []interface{} `json:"entities,omitempty"`
	EntityID          string        `json:"entityID,omitempty"`
	EntityParentID    string        `json:"entityParentID,omitempty"`
	EntityParentType  string        `json:"entityParentType,omitempty"`
	EntityScope       string        `json:"entityScope,omitempty"`
	EntityType        string        `json:"entityType,omitempty"`
	User              string        `json:"user,omitempty"`
	EventReceivedTime float64       `json:"eventReceivedTime,omitempty"`
	ExternalID        string        `json:"externalID,omitempty"`
	Type              string        `json:"type,omitempty"`
}

// NewEventLog returns a new *EventLog
func NewEventLog() *EventLog {

	return &EventLog{}
}

// Identity returns the Identity of the object.
func (o *EventLog) Identity() bambou.Identity {

	return EventLogIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *EventLog) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *EventLog) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the EventLog from the server
func (o *EventLog) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the EventLog into the server
func (o *EventLog) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the EventLog from the server
func (o *EventLog) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the EventLog
func (o *EventLog) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the EventLog
func (o *EventLog) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the EventLog
func (o *EventLog) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the EventLog
func (o *EventLog) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
