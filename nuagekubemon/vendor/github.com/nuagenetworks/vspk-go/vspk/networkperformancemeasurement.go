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

// NetworkPerformanceMeasurementIdentity represents the Identity of the object
var NetworkPerformanceMeasurementIdentity = bambou.Identity{
	Name:     "networkperformancemeasurement",
	Category: "networkperformancemeasurements",
}

// NetworkPerformanceMeasurementsList represents a list of NetworkPerformanceMeasurements
type NetworkPerformanceMeasurementsList []*NetworkPerformanceMeasurement

// NetworkPerformanceMeasurementsAncestor is the interface that an ancestor of a NetworkPerformanceMeasurement must implement.
// An Ancestor is defined as an entity that has NetworkPerformanceMeasurement as a descendant.
// An Ancestor can get a list of its child NetworkPerformanceMeasurements, but not necessarily create one.
type NetworkPerformanceMeasurementsAncestor interface {
	NetworkPerformanceMeasurements(*bambou.FetchingInfo) (NetworkPerformanceMeasurementsList, *bambou.Error)
}

// NetworkPerformanceMeasurementsParent is the interface that a parent of a NetworkPerformanceMeasurement must implement.
// A Parent is defined as an entity that has NetworkPerformanceMeasurement as a child.
// A Parent is an Ancestor which can create a NetworkPerformanceMeasurement.
type NetworkPerformanceMeasurementsParent interface {
	NetworkPerformanceMeasurementsAncestor
	CreateNetworkPerformanceMeasurement(*NetworkPerformanceMeasurement) *bambou.Error
}

// NetworkPerformanceMeasurement represents the model of a networkperformancemeasurement
type NetworkPerformanceMeasurement struct {
	ID                             string `json:"ID,omitempty"`
	ParentID                       string `json:"parentID,omitempty"`
	ParentType                     string `json:"parentType,omitempty"`
	Owner                          string `json:"owner,omitempty"`
	Name                           string `json:"name,omitempty"`
	ReadOnly                       bool   `json:"readOnly"`
	Description                    string `json:"description,omitempty"`
	AssociatedPerformanceMonitorID string `json:"associatedPerformanceMonitorID,omitempty"`
}

// NewNetworkPerformanceMeasurement returns a new *NetworkPerformanceMeasurement
func NewNetworkPerformanceMeasurement() *NetworkPerformanceMeasurement {

	return &NetworkPerformanceMeasurement{
		ReadOnly: false,
	}
}

// Identity returns the Identity of the object.
func (o *NetworkPerformanceMeasurement) Identity() bambou.Identity {

	return NetworkPerformanceMeasurementIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *NetworkPerformanceMeasurement) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *NetworkPerformanceMeasurement) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the NetworkPerformanceMeasurement from the server
func (o *NetworkPerformanceMeasurement) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the NetworkPerformanceMeasurement into the server
func (o *NetworkPerformanceMeasurement) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the NetworkPerformanceMeasurement from the server
func (o *NetworkPerformanceMeasurement) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// NetworkPerformanceBindings retrieves the list of child NetworkPerformanceBindings of the NetworkPerformanceMeasurement
func (o *NetworkPerformanceMeasurement) NetworkPerformanceBindings(info *bambou.FetchingInfo) (NetworkPerformanceBindingsList, *bambou.Error) {

	var list NetworkPerformanceBindingsList
	err := bambou.CurrentSession().FetchChildren(o, NetworkPerformanceBindingIdentity, &list, info)
	return list, err
}

// AssignNetworkPerformanceBindings assigns the list of NetworkPerformanceBindings to the NetworkPerformanceMeasurement
func (o *NetworkPerformanceMeasurement) AssignNetworkPerformanceBindings(children NetworkPerformanceBindingsList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, NetworkPerformanceBindingIdentity)
}

// Monitorscopes retrieves the list of child Monitorscopes of the NetworkPerformanceMeasurement
func (o *NetworkPerformanceMeasurement) Monitorscopes(info *bambou.FetchingInfo) (MonitorscopesList, *bambou.Error) {

	var list MonitorscopesList
	err := bambou.CurrentSession().FetchChildren(o, MonitorscopeIdentity, &list, info)
	return list, err
}

// CreateMonitorscope creates a new child Monitorscope under the NetworkPerformanceMeasurement
func (o *NetworkPerformanceMeasurement) CreateMonitorscope(child *Monitorscope) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
