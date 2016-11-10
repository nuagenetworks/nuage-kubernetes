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

// IKECertificateIdentity represents the Identity of the object
var IKECertificateIdentity = bambou.Identity{
	Name:     "ikecertificate",
	Category: "ikecertificates",
}

// IKECertificatesList represents a list of IKECertificates
type IKECertificatesList []*IKECertificate

// IKECertificatesAncestor is the interface of an ancestor of a IKECertificate must implement.
type IKECertificatesAncestor interface {
	IKECertificates(*bambou.FetchingInfo) (IKECertificatesList, *bambou.Error)
	CreateIKECertificates(*IKECertificate) *bambou.Error
}

// IKECertificate represents the model of a ikecertificate
type IKECertificate struct {
	ID                     string  `json:"ID,omitempty"`
	ParentID               string  `json:"parentID,omitempty"`
	ParentType             string  `json:"parentType,omitempty"`
	Owner                  string  `json:"owner,omitempty"`
	PEMEncoded             string  `json:"PEMEncoded,omitempty"`
	Name                   string  `json:"name,omitempty"`
	LastUpdatedBy          string  `json:"lastUpdatedBy,omitempty"`
	SerialNumber           int     `json:"serialNumber,omitempty"`
	Description            string  `json:"description,omitempty"`
	EntityScope            string  `json:"entityScope,omitempty"`
	NotAfter               float64 `json:"notAfter,omitempty"`
	NotBefore              float64 `json:"notBefore,omitempty"`
	AssociatedEnterpriseID string  `json:"associatedEnterpriseID,omitempty"`
	IssuerDN               string  `json:"issuerDN,omitempty"`
	SubjectDN              string  `json:"subjectDN,omitempty"`
	ExternalID             string  `json:"externalID,omitempty"`
}

// NewIKECertificate returns a new *IKECertificate
func NewIKECertificate() *IKECertificate {

	return &IKECertificate{}
}

// Identity returns the Identity of the object.
func (o *IKECertificate) Identity() bambou.Identity {

	return IKECertificateIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *IKECertificate) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *IKECertificate) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the IKECertificate from the server
func (o *IKECertificate) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the IKECertificate into the server
func (o *IKECertificate) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the IKECertificate from the server
func (o *IKECertificate) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the IKECertificate
func (o *IKECertificate) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the IKECertificate
func (o *IKECertificate) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the IKECertificate
func (o *IKECertificate) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the IKECertificate
func (o *IKECertificate) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
