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

// DestinationurlIdentity represents the Identity of the object
var DestinationurlIdentity = bambou.Identity{
	Name:     "destinationurl",
	Category: "destinationurls",
}

// DestinationurlsList represents a list of Destinationurls
type DestinationurlsList []*Destinationurl

// DestinationurlsAncestor is the interface that an ancestor of a Destinationurl must implement.
// An Ancestor is defined as an entity that has Destinationurl as a descendant.
// An Ancestor can get a list of its child Destinationurls, but not necessarily create one.
type DestinationurlsAncestor interface {
	Destinationurls(*bambou.FetchingInfo) (DestinationurlsList, *bambou.Error)
}

// DestinationurlsParent is the interface that a parent of a Destinationurl must implement.
// A Parent is defined as an entity that has Destinationurl as a child.
// A Parent is an Ancestor which can create a Destinationurl.
type DestinationurlsParent interface {
	DestinationurlsAncestor
	CreateDestinationurl(*Destinationurl) *bambou.Error
}

// Destinationurl represents the model of a destinationurl
type Destinationurl struct {
	ID                 string `json:"ID,omitempty"`
	ParentID           string `json:"parentID,omitempty"`
	ParentType         string `json:"parentType,omitempty"`
	Owner              string `json:"owner,omitempty"`
	URL                string `json:"URL,omitempty"`
	HTTPMethod         string `json:"HTTPMethod,omitempty"`
	PacketCount        int    `json:"packetCount,omitempty"`
	LastUpdatedBy      string `json:"lastUpdatedBy,omitempty"`
	PercentageWeight   int    `json:"percentageWeight,omitempty"`
	Timeout            int    `json:"timeout,omitempty"`
	EntityScope        string `json:"entityScope,omitempty"`
	DownThresholdCount int    `json:"downThresholdCount,omitempty"`
	ProbeInterval      int    `json:"probeInterval,omitempty"`
	ExternalID         string `json:"externalID,omitempty"`
}

// NewDestinationurl returns a new *Destinationurl
func NewDestinationurl() *Destinationurl {

	return &Destinationurl{
		HTTPMethod:         "HEAD",
		PacketCount:        1,
		Timeout:            3000,
		DownThresholdCount: 3,
		ProbeInterval:      10,
	}
}

// Identity returns the Identity of the object.
func (o *Destinationurl) Identity() bambou.Identity {

	return DestinationurlIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Destinationurl) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Destinationurl) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the Destinationurl from the server
func (o *Destinationurl) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the Destinationurl into the server
func (o *Destinationurl) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the Destinationurl from the server
func (o *Destinationurl) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the Destinationurl
func (o *Destinationurl) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the Destinationurl
func (o *Destinationurl) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the Destinationurl
func (o *Destinationurl) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the Destinationurl
func (o *Destinationurl) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
