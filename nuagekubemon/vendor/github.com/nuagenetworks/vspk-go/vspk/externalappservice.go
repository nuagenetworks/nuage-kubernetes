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

// ExternalAppServiceIdentity represents the Identity of the object
var ExternalAppServiceIdentity = bambou.Identity{
	Name:     "externalappservice",
	Category: "externalappservices",
}

// ExternalAppServicesList represents a list of ExternalAppServices
type ExternalAppServicesList []*ExternalAppService

// ExternalAppServicesAncestor is the interface that an ancestor of a ExternalAppService must implement.
// An Ancestor is defined as an entity that has ExternalAppService as a descendant.
// An Ancestor can get a list of its child ExternalAppServices, but not necessarily create one.
type ExternalAppServicesAncestor interface {
	ExternalAppServices(*bambou.FetchingInfo) (ExternalAppServicesList, *bambou.Error)
}

// ExternalAppServicesParent is the interface that a parent of a ExternalAppService must implement.
// A Parent is defined as an entity that has ExternalAppService as a child.
// A Parent is an Ancestor which can create a ExternalAppService.
type ExternalAppServicesParent interface {
	ExternalAppServicesAncestor
	CreateExternalAppService(*ExternalAppService) *bambou.Error
}

// ExternalAppService represents the model of a externalappservice
type ExternalAppService struct {
	ID                                 string `json:"ID,omitempty"`
	ParentID                           string `json:"parentID,omitempty"`
	ParentType                         string `json:"parentType,omitempty"`
	Owner                              string `json:"owner,omitempty"`
	Name                               string `json:"name,omitempty"`
	LastUpdatedBy                      string `json:"lastUpdatedBy,omitempty"`
	Description                        string `json:"description,omitempty"`
	DestinationNATAddress              string `json:"destinationNATAddress,omitempty"`
	DestinationNATEnabled              bool   `json:"destinationNATEnabled"`
	DestinationNATMask                 string `json:"destinationNATMask,omitempty"`
	Metadata                           string `json:"metadata,omitempty"`
	EgressType                         string `json:"egressType,omitempty"`
	VirtualIP                          string `json:"virtualIP,omitempty"`
	VirtualIPRequired                  bool   `json:"virtualIPRequired"`
	IngressType                        string `json:"ingressType,omitempty"`
	EntityScope                        string `json:"entityScope,omitempty"`
	SourceNATAddress                   string `json:"sourceNATAddress,omitempty"`
	SourceNATEnabled                   bool   `json:"sourceNATEnabled"`
	AssociatedServiceEgressGroupID     string `json:"associatedServiceEgressGroupID,omitempty"`
	AssociatedServiceEgressRedirectID  string `json:"associatedServiceEgressRedirectID,omitempty"`
	AssociatedServiceIngressGroupID    string `json:"associatedServiceIngressGroupID,omitempty"`
	AssociatedServiceIngressRedirectID string `json:"associatedServiceIngressRedirectID,omitempty"`
	ExternalID                         string `json:"externalID,omitempty"`
}

// NewExternalAppService returns a new *ExternalAppService
func NewExternalAppService() *ExternalAppService {

	return &ExternalAppService{}
}

// Identity returns the Identity of the object.
func (o *ExternalAppService) Identity() bambou.Identity {

	return ExternalAppServiceIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *ExternalAppService) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *ExternalAppService) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the ExternalAppService from the server
func (o *ExternalAppService) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the ExternalAppService into the server
func (o *ExternalAppService) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the ExternalAppService from the server
func (o *ExternalAppService) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the ExternalAppService
func (o *ExternalAppService) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the ExternalAppService
func (o *ExternalAppService) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the ExternalAppService
func (o *ExternalAppService) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the ExternalAppService
func (o *ExternalAppService) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
