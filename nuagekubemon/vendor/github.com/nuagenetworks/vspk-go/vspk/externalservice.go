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

// ExternalServiceIdentity represents the Identity of the object
var ExternalServiceIdentity = bambou.Identity{
	Name:     "externalservice",
	Category: "externalservices",
}

// ExternalServicesList represents a list of ExternalServices
type ExternalServicesList []*ExternalService

// ExternalServicesAncestor is the interface of an ancestor of a ExternalService must implement.
type ExternalServicesAncestor interface {
	ExternalServices(*bambou.FetchingInfo) (ExternalServicesList, *bambou.Error)
	CreateExternalServices(*ExternalService) *bambou.Error
}

// ExternalService represents the model of a externalservice
type ExternalService struct {
	ID            string `json:"ID,omitempty"`
	ParentID      string `json:"parentID,omitempty"`
	ParentType    string `json:"parentType,omitempty"`
	Owner         string `json:"owner,omitempty"`
	Name          string `json:"name,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	ServiceType   string `json:"serviceType,omitempty"`
	Description   string `json:"description,omitempty"`
	Direction     string `json:"direction,omitempty"`
	EntityScope   string `json:"entityScope,omitempty"`
	Stage         string `json:"stage,omitempty"`
	ExternalID    string `json:"externalID,omitempty"`
}

// NewExternalService returns a new *ExternalService
func NewExternalService() *ExternalService {

	return &ExternalService{}
}

// Identity returns the Identity of the object.
func (o *ExternalService) Identity() bambou.Identity {

	return ExternalServiceIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *ExternalService) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *ExternalService) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the ExternalService from the server
func (o *ExternalService) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the ExternalService into the server
func (o *ExternalService) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the ExternalService from the server
func (o *ExternalService) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the ExternalService
func (o *ExternalService) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the ExternalService
func (o *ExternalService) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// MetadataTags retrieves the list of child MetadataTags of the ExternalService
func (o *ExternalService) MetadataTags(info *bambou.FetchingInfo) (MetadataTagsList, *bambou.Error) {

	var list MetadataTagsList
	err := bambou.CurrentSession().FetchChildren(o, MetadataTagIdentity, &list, info)
	return list, err
}

// CreateMetadataTag creates a new child MetadataTag under the ExternalService
func (o *ExternalService) CreateMetadataTag(child *MetadataTag) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the ExternalService
func (o *ExternalService) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the ExternalService
func (o *ExternalService) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EndPoints retrieves the list of child EndPoints of the ExternalService
func (o *ExternalService) EndPoints(info *bambou.FetchingInfo) (EndPointsList, *bambou.Error) {

	var list EndPointsList
	err := bambou.CurrentSession().FetchChildren(o, EndPointIdentity, &list, info)
	return list, err
}

// CreateEndPoint creates a new child EndPoint under the ExternalService
func (o *ExternalService) CreateEndPoint(child *EndPoint) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the ExternalService
func (o *ExternalService) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}

// CreateEventLog creates a new child EventLog under the ExternalService
func (o *ExternalService) CreateEventLog(child *EventLog) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
