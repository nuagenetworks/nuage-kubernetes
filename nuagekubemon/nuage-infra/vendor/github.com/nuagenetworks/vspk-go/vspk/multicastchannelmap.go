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

// MultiCastChannelMapIdentity represents the Identity of the object
var MultiCastChannelMapIdentity = bambou.Identity{
	Name:     "multicastchannelmap",
	Category: "multicastchannelmaps",
}

// MultiCastChannelMapsList represents a list of MultiCastChannelMaps
type MultiCastChannelMapsList []*MultiCastChannelMap

// MultiCastChannelMapsAncestor is the interface that an ancestor of a MultiCastChannelMap must implement.
// An Ancestor is defined as an entity that has MultiCastChannelMap as a descendant.
// An Ancestor can get a list of its child MultiCastChannelMaps, but not necessarily create one.
type MultiCastChannelMapsAncestor interface {
	MultiCastChannelMaps(*bambou.FetchingInfo) (MultiCastChannelMapsList, *bambou.Error)
}

// MultiCastChannelMapsParent is the interface that a parent of a MultiCastChannelMap must implement.
// A Parent is defined as an entity that has MultiCastChannelMap as a child.
// A Parent is an Ancestor which can create a MultiCastChannelMap.
type MultiCastChannelMapsParent interface {
	MultiCastChannelMapsAncestor
	CreateMultiCastChannelMap(*MultiCastChannelMap) *bambou.Error
}

// MultiCastChannelMap represents the model of a multicastchannelmap
type MultiCastChannelMap struct {
	ID            string `json:"ID,omitempty"`
	ParentID      string `json:"parentID,omitempty"`
	ParentType    string `json:"parentType,omitempty"`
	Owner         string `json:"owner,omitempty"`
	Name          string `json:"name,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	Description   string `json:"description,omitempty"`
	EntityScope   string `json:"entityScope,omitempty"`
	ExternalID    string `json:"externalID,omitempty"`
}

// NewMultiCastChannelMap returns a new *MultiCastChannelMap
func NewMultiCastChannelMap() *MultiCastChannelMap {

	return &MultiCastChannelMap{}
}

// Identity returns the Identity of the object.
func (o *MultiCastChannelMap) Identity() bambou.Identity {

	return MultiCastChannelMapIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *MultiCastChannelMap) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *MultiCastChannelMap) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the MultiCastChannelMap from the server
func (o *MultiCastChannelMap) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the MultiCastChannelMap into the server
func (o *MultiCastChannelMap) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the MultiCastChannelMap from the server
func (o *MultiCastChannelMap) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the MultiCastChannelMap
func (o *MultiCastChannelMap) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the MultiCastChannelMap
func (o *MultiCastChannelMap) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the MultiCastChannelMap
func (o *MultiCastChannelMap) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the MultiCastChannelMap
func (o *MultiCastChannelMap) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// MultiCastRanges retrieves the list of child MultiCastRanges of the MultiCastChannelMap
func (o *MultiCastChannelMap) MultiCastRanges(info *bambou.FetchingInfo) (MultiCastRangesList, *bambou.Error) {

	var list MultiCastRangesList
	err := bambou.CurrentSession().FetchChildren(o, MultiCastRangeIdentity, &list, info)
	return list, err
}

// CreateMultiCastRange creates a new child MultiCastRange under the MultiCastChannelMap
func (o *MultiCastChannelMap) CreateMultiCastRange(child *MultiCastRange) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the MultiCastChannelMap
func (o *MultiCastChannelMap) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
