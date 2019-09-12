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

// DHCPv6OptionIdentity represents the Identity of the object
var DHCPv6OptionIdentity = bambou.Identity{
	Name:     "dhcpv6option",
	Category: "dhcpv6options",
}

// DHCPv6OptionsList represents a list of DHCPv6Options
type DHCPv6OptionsList []*DHCPv6Option

// DHCPv6OptionsAncestor is the interface that an ancestor of a DHCPv6Option must implement.
// An Ancestor is defined as an entity that has DHCPv6Option as a descendant.
// An Ancestor can get a list of its child DHCPv6Options, but not necessarily create one.
type DHCPv6OptionsAncestor interface {
	DHCPv6Options(*bambou.FetchingInfo) (DHCPv6OptionsList, *bambou.Error)
}

// DHCPv6OptionsParent is the interface that a parent of a DHCPv6Option must implement.
// A Parent is defined as an entity that has DHCPv6Option as a child.
// A Parent is an Ancestor which can create a DHCPv6Option.
type DHCPv6OptionsParent interface {
	DHCPv6OptionsAncestor
	CreateDHCPv6Option(*DHCPv6Option) *bambou.Error
}

// DHCPv6Option represents the model of a dhcpv6option
type DHCPv6Option struct {
	ID               string        `json:"ID,omitempty"`
	ParentID         string        `json:"parentID,omitempty"`
	ParentType       string        `json:"parentType,omitempty"`
	Owner            string        `json:"owner,omitempty"`
	Value            string        `json:"value,omitempty"`
	LastUpdatedBy    string        `json:"lastUpdatedBy,omitempty"`
	ActualType       int           `json:"actualType,omitempty"`
	ActualValues     []interface{} `json:"actualValues,omitempty"`
	Length           string        `json:"length,omitempty"`
	EmbeddedMetadata []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope      string        `json:"entityScope,omitempty"`
	ExternalID       string        `json:"externalID,omitempty"`
	Type             string        `json:"type,omitempty"`
}

// NewDHCPv6Option returns a new *DHCPv6Option
func NewDHCPv6Option() *DHCPv6Option {

	return &DHCPv6Option{}
}

// Identity returns the Identity of the object.
func (o *DHCPv6Option) Identity() bambou.Identity {

	return DHCPv6OptionIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *DHCPv6Option) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *DHCPv6Option) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the DHCPv6Option from the server
func (o *DHCPv6Option) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the DHCPv6Option into the server
func (o *DHCPv6Option) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the DHCPv6Option from the server
func (o *DHCPv6Option) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the DHCPv6Option
func (o *DHCPv6Option) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the DHCPv6Option
func (o *DHCPv6Option) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the DHCPv6Option
func (o *DHCPv6Option) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the DHCPv6Option
func (o *DHCPv6Option) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// EventLogs retrieves the list of child EventLogs of the DHCPv6Option
func (o *DHCPv6Option) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
