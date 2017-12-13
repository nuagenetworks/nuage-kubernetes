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

// EnterpriseSecuredDataIdentity represents the Identity of the object
var EnterpriseSecuredDataIdentity = bambou.Identity{
	Name:     "enterprisesecureddata",
	Category: "enterprisesecureddatas",
}

// EnterpriseSecuredDatasList represents a list of EnterpriseSecuredDatas
type EnterpriseSecuredDatasList []*EnterpriseSecuredData

// EnterpriseSecuredDatasAncestor is the interface that an ancestor of a EnterpriseSecuredData must implement.
// An Ancestor is defined as an entity that has EnterpriseSecuredData as a descendant.
// An Ancestor can get a list of its child EnterpriseSecuredDatas, but not necessarily create one.
type EnterpriseSecuredDatasAncestor interface {
	EnterpriseSecuredDatas(*bambou.FetchingInfo) (EnterpriseSecuredDatasList, *bambou.Error)
}

// EnterpriseSecuredDatasParent is the interface that a parent of a EnterpriseSecuredData must implement.
// A Parent is defined as an entity that has EnterpriseSecuredData as a child.
// A Parent is an Ancestor which can create a EnterpriseSecuredData.
type EnterpriseSecuredDatasParent interface {
	EnterpriseSecuredDatasAncestor
	CreateEnterpriseSecuredData(*EnterpriseSecuredData) *bambou.Error
}

// EnterpriseSecuredData represents the model of a enterprisesecureddata
type EnterpriseSecuredData struct {
	ID                        string `json:"ID,omitempty"`
	ParentID                  string `json:"parentID,omitempty"`
	ParentType                string `json:"parentType,omitempty"`
	Owner                     string `json:"owner,omitempty"`
	Hash                      string `json:"hash,omitempty"`
	LastUpdatedBy             string `json:"lastUpdatedBy,omitempty"`
	Data                      string `json:"data,omitempty"`
	SekId                     int    `json:"sekId,omitempty"`
	KeyserverCertSerialNumber string `json:"keyserverCertSerialNumber,omitempty"`
	SignedHash                string `json:"signedHash,omitempty"`
	EntityScope               string `json:"entityScope,omitempty"`
	ExternalID                string `json:"externalID,omitempty"`
}

// NewEnterpriseSecuredData returns a new *EnterpriseSecuredData
func NewEnterpriseSecuredData() *EnterpriseSecuredData {

	return &EnterpriseSecuredData{}
}

// Identity returns the Identity of the object.
func (o *EnterpriseSecuredData) Identity() bambou.Identity {

	return EnterpriseSecuredDataIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *EnterpriseSecuredData) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *EnterpriseSecuredData) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the EnterpriseSecuredData from the server
func (o *EnterpriseSecuredData) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the EnterpriseSecuredData into the server
func (o *EnterpriseSecuredData) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the EnterpriseSecuredData from the server
func (o *EnterpriseSecuredData) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the EnterpriseSecuredData
func (o *EnterpriseSecuredData) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the EnterpriseSecuredData
func (o *EnterpriseSecuredData) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the EnterpriseSecuredData
func (o *EnterpriseSecuredData) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the EnterpriseSecuredData
func (o *EnterpriseSecuredData) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
