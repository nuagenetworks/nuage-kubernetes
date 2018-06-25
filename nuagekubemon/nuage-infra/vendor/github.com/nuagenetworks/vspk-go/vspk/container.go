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

// ContainerIdentity represents the Identity of the object
var ContainerIdentity = bambou.Identity{
	Name:     "container",
	Category: "containers",
}

// ContainersList represents a list of Containers
type ContainersList []*Container

// ContainersAncestor is the interface that an ancestor of a Container must implement.
// An Ancestor is defined as an entity that has Container as a descendant.
// An Ancestor can get a list of its child Containers, but not necessarily create one.
type ContainersAncestor interface {
	Containers(*bambou.FetchingInfo) (ContainersList, *bambou.Error)
}

// ContainersParent is the interface that a parent of a Container must implement.
// A Parent is defined as an entity that has Container as a child.
// A Parent is an Ancestor which can create a Container.
type ContainersParent interface {
	ContainersAncestor
	CreateContainer(*Container) *bambou.Error
}

// Container represents the model of a container
type Container struct {
	ID              string        `json:"ID,omitempty"`
	ParentID        string        `json:"parentID,omitempty"`
	ParentType      string        `json:"parentType,omitempty"`
	Owner           string        `json:"owner,omitempty"`
	L2DomainIDs     []interface{} `json:"l2DomainIDs,omitempty"`
	VRSID           string        `json:"VRSID,omitempty"`
	UUID            string        `json:"UUID,omitempty"`
	Name            string        `json:"name,omitempty"`
	LastUpdatedBy   string        `json:"lastUpdatedBy,omitempty"`
	ReasonType      string        `json:"reasonType,omitempty"`
	DeleteExpiry    int           `json:"deleteExpiry,omitempty"`
	DeleteMode      string        `json:"deleteMode,omitempty"`
	ResyncInfo      interface{}   `json:"resyncInfo,omitempty"`
	SiteIdentifier  string        `json:"siteIdentifier,omitempty"`
	ImageID         string        `json:"imageID,omitempty"`
	ImageName       string        `json:"imageName,omitempty"`
	Interfaces      []interface{} `json:"interfaces,omitempty"`
	EnterpriseID    string        `json:"enterpriseID,omitempty"`
	EnterpriseName  string        `json:"enterpriseName,omitempty"`
	EntityScope     string        `json:"entityScope,omitempty"`
	DomainIDs       []interface{} `json:"domainIDs,omitempty"`
	ZoneIDs         []interface{} `json:"zoneIDs,omitempty"`
	OrchestrationID string        `json:"orchestrationID,omitempty"`
	UserID          string        `json:"userID,omitempty"`
	UserName        string        `json:"userName,omitempty"`
	Status          string        `json:"status,omitempty"`
	SubnetIDs       []interface{} `json:"subnetIDs,omitempty"`
	ExternalID      string        `json:"externalID,omitempty"`
	HypervisorIP    string        `json:"hypervisorIP,omitempty"`
}

// NewContainer returns a new *Container
func NewContainer() *Container {

	return &Container{}
}

// Identity returns the Identity of the object.
func (o *Container) Identity() bambou.Identity {

	return ContainerIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Container) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Container) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Container from the server
func (o *Container) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Container into the server
func (o *Container) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Container from the server
func (o *Container) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the Container
func (o *Container) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the Container
func (o *Container) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Alarms retrieves the list of child Alarms of the Container
func (o *Container) Alarms(info *bambou.FetchingInfo) (AlarmsList, *bambou.Error) {

	var list AlarmsList
	err := bambou.CurrentSession().FetchChildren(o, AlarmIdentity, &list, info)
	return list, err
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Container
func (o *Container) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Container
func (o *Container) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// ContainerInterfaces retrieves the list of child ContainerInterfaces of the Container
func (o *Container) ContainerInterfaces(info *bambou.FetchingInfo) (ContainerInterfacesList, *bambou.Error) {

	var list ContainerInterfacesList
	err := bambou.CurrentSession().FetchChildren(o, ContainerInterfaceIdentity, &list, info)
	return list, err
}

// CreateContainerInterface creates a new child ContainerInterface under the Container
func (o *Container) CreateContainerInterface(child *ContainerInterface) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// ContainerResyncs retrieves the list of child ContainerResyncs of the Container
func (o *Container) ContainerResyncs(info *bambou.FetchingInfo) (ContainerResyncsList, *bambou.Error) {

	var list ContainerResyncsList
	err := bambou.CurrentSession().FetchChildren(o, ContainerResyncIdentity, &list, info)
	return list, err
}

// CreateContainerResync creates a new child ContainerResync under the Container
func (o *Container) CreateContainerResync(child *ContainerResync) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VRSs retrieves the list of child VRSs of the Container
func (o *Container) VRSs(info *bambou.FetchingInfo) (VRSsList, *bambou.Error) {

	var list VRSsList
	err := bambou.CurrentSession().FetchChildren(o, VRSIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the Container
func (o *Container) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
