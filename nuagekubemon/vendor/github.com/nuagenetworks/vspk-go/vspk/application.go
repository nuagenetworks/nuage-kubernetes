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

// ApplicationIdentity represents the Identity of the object
var ApplicationIdentity = bambou.Identity{
	Name:     "application",
	Category: "applications",
}

// ApplicationsList represents a list of Applications
type ApplicationsList []*Application

// ApplicationsAncestor is the interface that an ancestor of a Application must implement.
// An Ancestor is defined as an entity that has Application as a descendant.
// An Ancestor can get a list of its child Applications, but not necessarily create one.
type ApplicationsAncestor interface {
	Applications(*bambou.FetchingInfo) (ApplicationsList, *bambou.Error)
}

// ApplicationsParent is the interface that a parent of a Application must implement.
// A Parent is defined as an entity that has Application as a child.
// A Parent is an Ancestor which can create a Application.
type ApplicationsParent interface {
	ApplicationsAncestor
	CreateApplication(*Application) *bambou.Error
}

// Application represents the model of a application
type Application struct {
	ID                                 string  `json:"ID,omitempty"`
	ParentID                           string  `json:"parentID,omitempty"`
	ParentType                         string  `json:"parentType,omitempty"`
	Owner                              string  `json:"owner,omitempty"`
	DSCP                               string  `json:"DSCP,omitempty"`
	Name                               string  `json:"name,omitempty"`
	ReadOnly                           bool    `json:"readOnly"`
	PerformanceMonitorType             string  `json:"performanceMonitorType,omitempty"`
	Description                        string  `json:"description,omitempty"`
	DestinationIP                      string  `json:"destinationIP,omitempty"`
	DestinationPort                    string  `json:"destinationPort,omitempty"`
	EnablePPS                          bool    `json:"enablePPS"`
	OneWayDelay                        int     `json:"oneWayDelay,omitempty"`
	OneWayJitter                       int     `json:"oneWayJitter,omitempty"`
	OneWayLoss                         float64 `json:"oneWayLoss,omitempty"`
	PostClassificationPath             string  `json:"postClassificationPath,omitempty"`
	SourceIP                           string  `json:"sourceIP,omitempty"`
	SourcePort                         string  `json:"sourcePort,omitempty"`
	OptimizePathSelection              string  `json:"optimizePathSelection,omitempty"`
	PreClassificationPath              string  `json:"preClassificationPath,omitempty"`
	Protocol                           string  `json:"protocol,omitempty"`
	AssociatedL7ApplicationSignatureID string  `json:"associatedL7ApplicationSignatureID,omitempty"`
	EtherType                          string  `json:"etherType,omitempty"`
	Symmetry                           bool    `json:"symmetry"`
}

// NewApplication returns a new *Application
func NewApplication() *Application {

	return &Application{
		ReadOnly:               false,
		PerformanceMonitorType: "FIRST_PACKET",
		EnablePPS:              false,
		PostClassificationPath: "ANY",
		PreClassificationPath:  "DEFAULT",
		Protocol:               "NONE",
	}
}

// Identity returns the Identity of the object.
func (o *Application) Identity() bambou.Identity {

	return ApplicationIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Application) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Application) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Application from the server
func (o *Application) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Application into the server
func (o *Application) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Application from the server
func (o *Application) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Monitorscopes retrieves the list of child Monitorscopes of the Application
func (o *Application) Monitorscopes(info *bambou.FetchingInfo) (MonitorscopesList, *bambou.Error) {

	var list MonitorscopesList
	err := bambou.CurrentSession().FetchChildren(o, MonitorscopeIdentity, &list, info)
	return list, err
}

// CreateMonitorscope creates a new child Monitorscope under the Application
func (o *Application) CreateMonitorscope(child *Monitorscope) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// ApplicationBindings retrieves the list of child ApplicationBindings of the Application
func (o *Application) ApplicationBindings(info *bambou.FetchingInfo) (ApplicationBindingsList, *bambou.Error) {

	var list ApplicationBindingsList
	err := bambou.CurrentSession().FetchChildren(o, ApplicationBindingIdentity, &list, info)
	return list, err
}
