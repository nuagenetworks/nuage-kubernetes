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

// StaticRouteIdentity represents the Identity of the object
var StaticRouteIdentity = bambou.Identity{
	Name:     "staticroute",
	Category: "staticroutes",
}

// StaticRoutesList represents a list of StaticRoutes
type StaticRoutesList []*StaticRoute

// StaticRoutesAncestor is the interface that an ancestor of a StaticRoute must implement.
// An Ancestor is defined as an entity that has StaticRoute as a descendant.
// An Ancestor can get a list of its child StaticRoutes, but not necessarily create one.
type StaticRoutesAncestor interface {
	StaticRoutes(*bambou.FetchingInfo) (StaticRoutesList, *bambou.Error)
}

// StaticRoutesParent is the interface that a parent of a StaticRoute must implement.
// A Parent is defined as an entity that has StaticRoute as a child.
// A Parent is an Ancestor which can create a StaticRoute.
type StaticRoutesParent interface {
	StaticRoutesAncestor
	CreateStaticRoute(*StaticRoute) *bambou.Error
}

// StaticRoute represents the model of a staticroute
type StaticRoute struct {
	ID                   string        `json:"ID,omitempty"`
	ParentID             string        `json:"parentID,omitempty"`
	ParentType           string        `json:"parentType,omitempty"`
	Owner                string        `json:"owner,omitempty"`
	BFDEnabled           bool          `json:"BFDEnabled"`
	IPType               string        `json:"IPType,omitempty"`
	IPv6Address          string        `json:"IPv6Address,omitempty"`
	LastUpdatedBy        string        `json:"lastUpdatedBy,omitempty"`
	Address              string        `json:"address,omitempty"`
	Netmask              string        `json:"netmask,omitempty"`
	NextHopIp            string        `json:"nextHopIp,omitempty"`
	BlackHoleEnabled     bool          `json:"blackHoleEnabled"`
	EmbeddedMetadata     []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope          string        `json:"entityScope,omitempty"`
	RouteDistinguisher   string        `json:"routeDistinguisher,omitempty"`
	AssociatedGatewayIDs []interface{} `json:"associatedGatewayIDs,omitempty"`
	AssociatedSubnetID   string        `json:"associatedSubnetID,omitempty"`
	ExternalID           string        `json:"externalID,omitempty"`
	Type                 string        `json:"type,omitempty"`
}

// NewStaticRoute returns a new *StaticRoute
func NewStaticRoute() *StaticRoute {

	return &StaticRoute{
		BFDEnabled:       false,
		BlackHoleEnabled: false,
	}
}

// Identity returns the Identity of the object.
func (o *StaticRoute) Identity() bambou.Identity {

	return StaticRouteIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *StaticRoute) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *StaticRoute) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the StaticRoute from the server
func (o *StaticRoute) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the StaticRoute into the server
func (o *StaticRoute) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the StaticRoute from the server
func (o *StaticRoute) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// DeploymentFailures retrieves the list of child DeploymentFailures of the StaticRoute
func (o *StaticRoute) DeploymentFailures(info *bambou.FetchingInfo) (DeploymentFailuresList, *bambou.Error) {

	var list DeploymentFailuresList
	err := bambou.CurrentSession().FetchChildren(o, DeploymentFailureIdentity, &list, info)
	return list, err
}

// CreateDeploymentFailure creates a new child DeploymentFailure under the StaticRoute
func (o *StaticRoute) CreateDeploymentFailure(child *DeploymentFailure) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the StaticRoute
func (o *StaticRoute) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the StaticRoute
func (o *StaticRoute) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the StaticRoute
func (o *StaticRoute) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the StaticRoute
func (o *StaticRoute) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the StaticRoute
func (o *StaticRoute) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
