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

// VPortMirrorIdentity represents the Identity of the object
var VPortMirrorIdentity = bambou.Identity{
	Name:     "vportmirror",
	Category: "vportmirrors",
}

// VPortMirrorsList represents a list of VPortMirrors
type VPortMirrorsList []*VPortMirror

// VPortMirrorsAncestor is the interface of an ancestor of a VPortMirror must implement.
type VPortMirrorsAncestor interface {
	VPortMirrors(*bambou.FetchingInfo) (VPortMirrorsList, *bambou.Error)
	CreateVPortMirrors(*VPortMirror) *bambou.Error
}

// VPortMirror represents the model of a vportmirror
type VPortMirror struct {
	ID                    string `json:"ID,omitempty"`
	ParentID              string `json:"parentID,omitempty"`
	ParentType            string `json:"parentType,omitempty"`
	Owner                 string `json:"owner,omitempty"`
	VPortName             string `json:"VPortName,omitempty"`
	LastUpdatedBy         string `json:"lastUpdatedBy,omitempty"`
	NetworkName           string `json:"networkName,omitempty"`
	MirrorDestinationID   string `json:"mirrorDestinationID,omitempty"`
	MirrorDestinationName string `json:"mirrorDestinationName,omitempty"`
	MirrorDirection       string `json:"mirrorDirection,omitempty"`
	EnterpiseName         string `json:"enterpiseName,omitempty"`
	EntityScope           string `json:"entityScope,omitempty"`
	DomainName            string `json:"domainName,omitempty"`
	VportId               string `json:"vportId,omitempty"`
	AttachedNetworkType   string `json:"attachedNetworkType,omitempty"`
	ExternalID            string `json:"externalID,omitempty"`
}

// NewVPortMirror returns a new *VPortMirror
func NewVPortMirror() *VPortMirror {

	return &VPortMirror{
		MirrorDirection: "BOTH",
	}
}

// Identity returns the Identity of the object.
func (o *VPortMirror) Identity() bambou.Identity {

	return VPortMirrorIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VPortMirror) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VPortMirror) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VPortMirror from the server
func (o *VPortMirror) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VPortMirror into the server
func (o *VPortMirror) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VPortMirror from the server
func (o *VPortMirror) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the VPortMirror
func (o *VPortMirror) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VPortMirror
func (o *VPortMirror) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VPortMirror
func (o *VPortMirror) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VPortMirror
func (o *VPortMirror) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
