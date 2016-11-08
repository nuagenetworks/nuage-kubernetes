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

// ApplicationperformancemanagementIdentity represents the Identity of the object
var ApplicationperformancemanagementIdentity = bambou.Identity{
	Name:     "applicationperformancemanagement",
	Category: "applicationperformancemanagements",
}

// ApplicationperformancemanagementsList represents a list of Applicationperformancemanagements
type ApplicationperformancemanagementsList []*Applicationperformancemanagement

// ApplicationperformancemanagementsAncestor is the interface of an ancestor of a Applicationperformancemanagement must implement.
type ApplicationperformancemanagementsAncestor interface {
	Applicationperformancemanagements(*bambou.FetchingInfo) (ApplicationperformancemanagementsList, *bambou.Error)
	CreateApplicationperformancemanagements(*Applicationperformancemanagement) *bambou.Error
}

// Applicationperformancemanagement represents the model of a applicationperformancemanagement
type Applicationperformancemanagement struct {
	ID                             string `json:"ID,omitempty"`
	ParentID                       string `json:"parentID,omitempty"`
	ParentType                     string `json:"parentType,omitempty"`
	Owner                          string `json:"owner,omitempty"`
	Name                           string `json:"name,omitempty"`
	ReadOnly                       bool   `json:"readOnly"`
	Description                    string `json:"description,omitempty"`
	ApplicationGroupType           string `json:"applicationGroupType,omitempty"`
	AssociatedPerformanceMonitorID string `json:"associatedPerformanceMonitorID,omitempty"`
}

// NewApplicationperformancemanagement returns a new *Applicationperformancemanagement
func NewApplicationperformancemanagement() *Applicationperformancemanagement {

	return &Applicationperformancemanagement{}
}

// Identity returns the Identity of the object.
func (o *Applicationperformancemanagement) Identity() bambou.Identity {

	return ApplicationperformancemanagementIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Applicationperformancemanagement) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Applicationperformancemanagement) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Applicationperformancemanagement from the server
func (o *Applicationperformancemanagement) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Applicationperformancemanagement into the server
func (o *Applicationperformancemanagement) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Applicationperformancemanagement from the server
func (o *Applicationperformancemanagement) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// ApplicationBindings retrieves the list of child ApplicationBindings of the Applicationperformancemanagement
func (o *Applicationperformancemanagement) ApplicationBindings(info *bambou.FetchingInfo) (ApplicationBindingsList, *bambou.Error) {

	var list ApplicationBindingsList
	err := bambou.CurrentSession().FetchChildren(o, ApplicationBindingIdentity, &list, info)
	return list, err
}

// CreateApplicationBinding creates a new child ApplicationBinding under the Applicationperformancemanagement
func (o *Applicationperformancemanagement) CreateApplicationBinding(child *ApplicationBinding) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
