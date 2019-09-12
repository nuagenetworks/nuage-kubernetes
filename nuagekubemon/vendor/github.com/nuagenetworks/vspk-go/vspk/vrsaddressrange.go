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

// VRSAddressRangeIdentity represents the Identity of the object
var VRSAddressRangeIdentity = bambou.Identity{
	Name:     "vrsaddressrange",
	Category: "vrsaddressranges",
}

// VRSAddressRangesList represents a list of VRSAddressRanges
type VRSAddressRangesList []*VRSAddressRange

// VRSAddressRangesAncestor is the interface that an ancestor of a VRSAddressRange must implement.
// An Ancestor is defined as an entity that has VRSAddressRange as a descendant.
// An Ancestor can get a list of its child VRSAddressRanges, but not necessarily create one.
type VRSAddressRangesAncestor interface {
	VRSAddressRanges(*bambou.FetchingInfo) (VRSAddressRangesList, *bambou.Error)
}

// VRSAddressRangesParent is the interface that a parent of a VRSAddressRange must implement.
// A Parent is defined as an entity that has VRSAddressRange as a child.
// A Parent is an Ancestor which can create a VRSAddressRange.
type VRSAddressRangesParent interface {
	VRSAddressRangesAncestor
	CreateVRSAddressRange(*VRSAddressRange) *bambou.Error
}

// VRSAddressRange represents the model of a vrsaddressrange
type VRSAddressRange struct {
	ID               string        `json:"ID,omitempty"`
	ParentID         string        `json:"parentID,omitempty"`
	ParentType       string        `json:"parentType,omitempty"`
	Owner            string        `json:"owner,omitempty"`
	IPType           string        `json:"IPType,omitempty"`
	LastUpdatedBy    string        `json:"lastUpdatedBy,omitempty"`
	MaxAddress       string        `json:"maxAddress,omitempty"`
	MinAddress       string        `json:"minAddress,omitempty"`
	EmbeddedMetadata []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope      string        `json:"entityScope,omitempty"`
	ExternalID       string        `json:"externalID,omitempty"`
}

// NewVRSAddressRange returns a new *VRSAddressRange
func NewVRSAddressRange() *VRSAddressRange {

	return &VRSAddressRange{}
}

// Identity returns the Identity of the object.
func (o *VRSAddressRange) Identity() bambou.Identity {

	return VRSAddressRangeIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *VRSAddressRange) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *VRSAddressRange) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the VRSAddressRange from the server
func (o *VRSAddressRange) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the VRSAddressRange into the server
func (o *VRSAddressRange) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the VRSAddressRange from the server
func (o *VRSAddressRange) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the VRSAddressRange
func (o *VRSAddressRange) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the VRSAddressRange
func (o *VRSAddressRange) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the VRSAddressRange
func (o *VRSAddressRange) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the VRSAddressRange
func (o *VRSAddressRange) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
