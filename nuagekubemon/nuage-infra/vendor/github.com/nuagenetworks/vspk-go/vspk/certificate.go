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

// CertificateIdentity represents the Identity of the object
var CertificateIdentity = bambou.Identity{
	Name:     "certificate",
	Category: "certificates",
}

// CertificatesList represents a list of Certificates
type CertificatesList []*Certificate

// CertificatesAncestor is the interface that an ancestor of a Certificate must implement.
// An Ancestor is defined as an entity that has Certificate as a descendant.
// An Ancestor can get a list of its child Certificates, but not necessarily create one.
type CertificatesAncestor interface {
	Certificates(*bambou.FetchingInfo) (CertificatesList, *bambou.Error)
}

// CertificatesParent is the interface that a parent of a Certificate must implement.
// A Parent is defined as an entity that has Certificate as a child.
// A Parent is an Ancestor which can create a Certificate.
type CertificatesParent interface {
	CertificatesAncestor
	CreateCertificate(*Certificate) *bambou.Error
}

// Certificate represents the model of a certificate
type Certificate struct {
	ID           string `json:"ID,omitempty"`
	ParentID     string `json:"parentID,omitempty"`
	ParentType   string `json:"parentType,omitempty"`
	Owner        string `json:"owner,omitempty"`
	PemEncoded   string `json:"pemEncoded,omitempty"`
	SerialNumber int    `json:"serialNumber,omitempty"`
	EntityScope  string `json:"entityScope,omitempty"`
	IssuerDN     string `json:"issuerDN,omitempty"`
	SubjectDN    string `json:"subjectDN,omitempty"`
	PublicKey    string `json:"publicKey,omitempty"`
	ExternalID   string `json:"externalID,omitempty"`
}

// NewCertificate returns a new *Certificate
func NewCertificate() *Certificate {

	return &Certificate{}
}

// Identity returns the Identity of the object.
func (o *Certificate) Identity() bambou.Identity {

	return CertificateIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Certificate) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Certificate) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Certificate from the server
func (o *Certificate) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Certificate into the server
func (o *Certificate) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Certificate from the server
func (o *Certificate) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the Certificate
func (o *Certificate) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the Certificate
func (o *Certificate) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Certificate
func (o *Certificate) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Certificate
func (o *Certificate) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
