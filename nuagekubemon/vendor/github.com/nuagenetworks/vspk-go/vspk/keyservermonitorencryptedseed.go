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

// KeyServerMonitorEncryptedSeedIdentity represents the Identity of the object
var KeyServerMonitorEncryptedSeedIdentity = bambou.Identity{
	Name:     "keyservermonitorencryptedseed",
	Category: "keyservermonitorencryptedseeds",
}

// KeyServerMonitorEncryptedSeedsList represents a list of KeyServerMonitorEncryptedSeeds
type KeyServerMonitorEncryptedSeedsList []*KeyServerMonitorEncryptedSeed

// KeyServerMonitorEncryptedSeedsAncestor is the interface that an ancestor of a KeyServerMonitorEncryptedSeed must implement.
// An Ancestor is defined as an entity that has KeyServerMonitorEncryptedSeed as a descendant.
// An Ancestor can get a list of its child KeyServerMonitorEncryptedSeeds, but not necessarily create one.
type KeyServerMonitorEncryptedSeedsAncestor interface {
	KeyServerMonitorEncryptedSeeds(*bambou.FetchingInfo) (KeyServerMonitorEncryptedSeedsList, *bambou.Error)
}

// KeyServerMonitorEncryptedSeedsParent is the interface that a parent of a KeyServerMonitorEncryptedSeed must implement.
// A Parent is defined as an entity that has KeyServerMonitorEncryptedSeed as a child.
// A Parent is an Ancestor which can create a KeyServerMonitorEncryptedSeed.
type KeyServerMonitorEncryptedSeedsParent interface {
	KeyServerMonitorEncryptedSeedsAncestor
	CreateKeyServerMonitorEncryptedSeed(*KeyServerMonitorEncryptedSeed) *bambou.Error
}

// KeyServerMonitorEncryptedSeed represents the model of a keyservermonitorencryptedseed
type KeyServerMonitorEncryptedSeed struct {
	ID                                         string        `json:"ID,omitempty"`
	ParentID                                   string        `json:"parentID,omitempty"`
	ParentType                                 string        `json:"parentType,omitempty"`
	Owner                                      string        `json:"owner,omitempty"`
	SEKCreationTime                            int           `json:"SEKCreationTime,omitempty"`
	LastUpdatedBy                              string        `json:"lastUpdatedBy,omitempty"`
	SeedType                                   string        `json:"seedType,omitempty"`
	KeyServerCertificateSerialNumber           int           `json:"keyServerCertificateSerialNumber,omitempty"`
	EmbeddedMetadata                           []interface{} `json:"embeddedMetadata,omitempty"`
	EnterpriseSecuredDataID                    string        `json:"enterpriseSecuredDataID,omitempty"`
	EntityScope                                string        `json:"entityScope,omitempty"`
	AssociatedKeyServerMonitorSEKCreationTime  int           `json:"associatedKeyServerMonitorSEKCreationTime,omitempty"`
	AssociatedKeyServerMonitorSEKID            string        `json:"associatedKeyServerMonitorSEKID,omitempty"`
	AssociatedKeyServerMonitorSeedCreationTime int           `json:"associatedKeyServerMonitorSeedCreationTime,omitempty"`
	AssociatedKeyServerMonitorSeedID           string        `json:"associatedKeyServerMonitorSeedID,omitempty"`
	ExternalID                                 string        `json:"externalID,omitempty"`
}

// NewKeyServerMonitorEncryptedSeed returns a new *KeyServerMonitorEncryptedSeed
func NewKeyServerMonitorEncryptedSeed() *KeyServerMonitorEncryptedSeed {

	return &KeyServerMonitorEncryptedSeed{}
}

// Identity returns the Identity of the object.
func (o *KeyServerMonitorEncryptedSeed) Identity() bambou.Identity {

	return KeyServerMonitorEncryptedSeedIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *KeyServerMonitorEncryptedSeed) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *KeyServerMonitorEncryptedSeed) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the KeyServerMonitorEncryptedSeed from the server
func (o *KeyServerMonitorEncryptedSeed) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the KeyServerMonitorEncryptedSeed into the server
func (o *KeyServerMonitorEncryptedSeed) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the KeyServerMonitorEncryptedSeed from the server
func (o *KeyServerMonitorEncryptedSeed) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the KeyServerMonitorEncryptedSeed
func (o *KeyServerMonitorEncryptedSeed) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the KeyServerMonitorEncryptedSeed
func (o *KeyServerMonitorEncryptedSeed) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the KeyServerMonitorEncryptedSeed
func (o *KeyServerMonitorEncryptedSeed) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the KeyServerMonitorEncryptedSeed
func (o *KeyServerMonitorEncryptedSeed) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
