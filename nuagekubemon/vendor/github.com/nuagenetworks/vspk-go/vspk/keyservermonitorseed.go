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

// KeyServerMonitorSeedIdentity represents the Identity of the object
var KeyServerMonitorSeedIdentity = bambou.Identity{
	Name:     "keyservermonitorseed",
	Category: "keyservermonitorseeds",
}

// KeyServerMonitorSeedsList represents a list of KeyServerMonitorSeeds
type KeyServerMonitorSeedsList []*KeyServerMonitorSeed

// KeyServerMonitorSeedsAncestor is the interface of an ancestor of a KeyServerMonitorSeed must implement.
type KeyServerMonitorSeedsAncestor interface {
	KeyServerMonitorSeeds(*bambou.FetchingInfo) (KeyServerMonitorSeedsList, *bambou.Error)
	CreateKeyServerMonitorSeeds(*KeyServerMonitorSeed) *bambou.Error
}

// KeyServerMonitorSeed represents the model of a keyservermonitorseed
type KeyServerMonitorSeed struct {
	ID                                 string `json:"ID,omitempty"`
	ParentID                           string `json:"parentID,omitempty"`
	ParentType                         string `json:"parentType,omitempty"`
	Owner                              string `json:"owner,omitempty"`
	LastUpdatedBy                      string `json:"lastUpdatedBy,omitempty"`
	SeedTrafficAuthenticationAlgorithm string `json:"seedTrafficAuthenticationAlgorithm,omitempty"`
	SeedTrafficEncryptionAlgorithm     string `json:"seedTrafficEncryptionAlgorithm,omitempty"`
	SeedTrafficEncryptionKeyLifetime   int    `json:"seedTrafficEncryptionKeyLifetime,omitempty"`
	Lifetime                           int    `json:"lifetime,omitempty"`
	EntityScope                        string `json:"entityScope,omitempty"`
	CreationTime                       int    `json:"creationTime,omitempty"`
	StartTime                          int    `json:"startTime,omitempty"`
	ExternalID                         string `json:"externalID,omitempty"`
}

// NewKeyServerMonitorSeed returns a new *KeyServerMonitorSeed
func NewKeyServerMonitorSeed() *KeyServerMonitorSeed {

	return &KeyServerMonitorSeed{}
}

// Identity returns the Identity of the object.
func (o *KeyServerMonitorSeed) Identity() bambou.Identity {

	return KeyServerMonitorSeedIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *KeyServerMonitorSeed) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *KeyServerMonitorSeed) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the KeyServerMonitorSeed from the server
func (o *KeyServerMonitorSeed) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the KeyServerMonitorSeed into the server
func (o *KeyServerMonitorSeed) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the KeyServerMonitorSeed from the server
func (o *KeyServerMonitorSeed) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the KeyServerMonitorSeed
func (o *KeyServerMonitorSeed) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the KeyServerMonitorSeed
func (o *KeyServerMonitorSeed) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// KeyServerMonitorEncryptedSeeds retrieves the list of child KeyServerMonitorEncryptedSeeds of the KeyServerMonitorSeed
func (o *KeyServerMonitorSeed) KeyServerMonitorEncryptedSeeds(info *bambou.FetchingInfo) (KeyServerMonitorEncryptedSeedsList, *bambou.Error) {

	var list KeyServerMonitorEncryptedSeedsList
	err := bambou.CurrentSession().FetchChildren(o, KeyServerMonitorEncryptedSeedIdentity, &list, info)
	return list, err
}

// CreateKeyServerMonitorEncryptedSeed creates a new child KeyServerMonitorEncryptedSeed under the KeyServerMonitorSeed
func (o *KeyServerMonitorSeed) CreateKeyServerMonitorEncryptedSeed(child *KeyServerMonitorEncryptedSeed) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the KeyServerMonitorSeed
func (o *KeyServerMonitorSeed) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the KeyServerMonitorSeed
func (o *KeyServerMonitorSeed) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
