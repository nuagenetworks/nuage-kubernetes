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

// VSPIdentity represents the Identity of the object
var VSPIdentity = bambou.Identity{
	Name:     "vsp",
	Category: "vsps",
}

// VSPsList represents a list of VSPs
type VSPsList []*VSP

// VSPsAncestor is the interface of an ancestor of a VSP must implement.
type VSPsAncestor interface {
	VSPs(*bambou.FetchingInfo) (VSPsList, *bambou.Error)
	CreateVSPs(*VSP) *bambou.Error
}

// VSP represents the model of a vsp
type VSP struct {
	ID             string `json:"ID,omitempty"`
	ParentID       string `json:"parentID,omitempty"`
	ParentType     string `json:"parentType,omitempty"`
	Owner          string `json:"owner,omitempty"`
	Name           string `json:"name,omitempty"`
	LastUpdatedBy  string `json:"lastUpdatedBy,omitempty"`
	Description    string `json:"description,omitempty"`
	EntityScope    string `json:"entityScope,omitempty"`
	Location       string `json:"location,omitempty"`
	ProductVersion string `json:"productVersion,omitempty"`
	ExternalID     string `json:"externalID,omitempty"`
}

// NewVSP returns a new *VSP
func NewVSP() *VSP {

	return &VSP{}
}

// Identity returns the Identity of the object.
func (o *VSP) Identity() bambou.Identity {

	return VSPIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VSP) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VSP) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VSP from the server
func (o *VSP) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VSP into the server
func (o *VSP) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VSP from the server
func (o *VSP) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the VSP
func (o *VSP) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VSP
func (o *VSP) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VSP
func (o *VSP) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VSP
func (o *VSP) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// HSCs retrieves the list of child HSCs of the VSP
func (o *VSP) HSCs(info *bambou.FetchingInfo) (HSCsList, *bambou.Error) {

	var list HSCsList
	err := bambou.CurrentSession().FetchChildren(o, HSCIdentity, &list, info)
	return list, err
}

// CreateHSC creates a new child HSC under the VSP
func (o *VSP) CreateHSC(child *HSC) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VSCs retrieves the list of child VSCs of the VSP
func (o *VSP) VSCs(info *bambou.FetchingInfo) (VSCsList, *bambou.Error) {

	var list VSCsList
	err := bambou.CurrentSession().FetchChildren(o, VSCIdentity, &list, info)
	return list, err
}

// CreateVSC creates a new child VSC under the VSP
func (o *VSP) CreateVSC(child *VSC) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VSDs retrieves the list of child VSDs of the VSP
func (o *VSP) VSDs(info *bambou.FetchingInfo) (VSDsList, *bambou.Error) {

	var list VSDsList
	err := bambou.CurrentSession().FetchChildren(o, VSDIdentity, &list, info)
	return list, err
}

// CreateVSD creates a new child VSD under the VSP
func (o *VSP) CreateVSD(child *VSD) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the VSP
func (o *VSP) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}

// CreateEventLog creates a new child EventLog under the VSP
func (o *VSP) CreateEventLog(child *EventLog) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
