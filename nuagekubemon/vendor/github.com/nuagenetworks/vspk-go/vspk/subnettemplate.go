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

// SubnetTemplateIdentity represents the Identity of the object
var SubnetTemplateIdentity = bambou.Identity{
	Name:     "subnettemplate",
	Category: "subnettemplates",
}

// SubnetTemplatesList represents a list of SubnetTemplates
type SubnetTemplatesList []*SubnetTemplate

// SubnetTemplatesAncestor is the interface that an ancestor of a SubnetTemplate must implement.
// An Ancestor is defined as an entity that has SubnetTemplate as a descendant.
// An Ancestor can get a list of its child SubnetTemplates, but not necessarily create one.
type SubnetTemplatesAncestor interface {
	SubnetTemplates(*bambou.FetchingInfo) (SubnetTemplatesList, *bambou.Error)
}

// SubnetTemplatesParent is the interface that a parent of a SubnetTemplate must implement.
// A Parent is defined as an entity that has SubnetTemplate as a child.
// A Parent is an Ancestor which can create a SubnetTemplate.
type SubnetTemplatesParent interface {
	SubnetTemplatesAncestor
	CreateSubnetTemplate(*SubnetTemplate) *bambou.Error
}

// SubnetTemplate represents the model of a subnettemplate
type SubnetTemplate struct {
	ID                              string        `json:"ID,omitempty"`
	ParentID                        string        `json:"parentID,omitempty"`
	ParentType                      string        `json:"parentType,omitempty"`
	Owner                           string        `json:"owner,omitempty"`
	DPI                             string        `json:"DPI,omitempty"`
	IPType                          string        `json:"IPType,omitempty"`
	IPv6Address                     string        `json:"IPv6Address,omitempty"`
	IPv6Gateway                     string        `json:"IPv6Gateway,omitempty"`
	Name                            string        `json:"name,omitempty"`
	LastUpdatedBy                   string        `json:"lastUpdatedBy,omitempty"`
	Gateway                         string        `json:"gateway,omitempty"`
	Address                         string        `json:"address,omitempty"`
	Description                     string        `json:"description,omitempty"`
	Netmask                         string        `json:"netmask,omitempty"`
	EmbeddedMetadata                []interface{} `json:"embeddedMetadata,omitempty"`
	EnableDHCPv4                    bool          `json:"enableDHCPv4"`
	EnableDHCPv6                    bool          `json:"enableDHCPv6"`
	Encryption                      string        `json:"encryption,omitempty"`
	EntityScope                     string        `json:"entityScope,omitempty"`
	SplitSubnet                     bool          `json:"splitSubnet"`
	ProxyARP                        bool          `json:"proxyARP"`
	UseGlobalMAC                    string        `json:"useGlobalMAC,omitempty"`
	AssociatedMulticastChannelMapID string        `json:"associatedMulticastChannelMapID,omitempty"`
	DualStackDynamicIPAllocation    bool          `json:"dualStackDynamicIPAllocation"`
	Multicast                       string        `json:"multicast,omitempty"`
	ExternalID                      string        `json:"externalID,omitempty"`
}

// NewSubnetTemplate returns a new *SubnetTemplate
func NewSubnetTemplate() *SubnetTemplate {

	return &SubnetTemplate{
		DPI:                          "INHERITED",
		IPType:                       "IPV4",
		EnableDHCPv4:                 true,
		EnableDHCPv6:                 false,
		UseGlobalMAC:                 "ENTERPRISE_DEFAULT",
		DualStackDynamicIPAllocation: false,
		Multicast:                    "INHERITED",
	}
}

// Identity returns the Identity of the object.
func (o *SubnetTemplate) Identity() bambou.Identity {

	return SubnetTemplateIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *SubnetTemplate) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *SubnetTemplate) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the SubnetTemplate from the server
func (o *SubnetTemplate) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the SubnetTemplate into the server
func (o *SubnetTemplate) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the SubnetTemplate from the server
func (o *SubnetTemplate) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// AddressRanges retrieves the list of child AddressRanges of the SubnetTemplate
func (o *SubnetTemplate) AddressRanges(info *bambou.FetchingInfo) (AddressRangesList, *bambou.Error) {

	var list AddressRangesList
	err := bambou.CurrentSession().FetchChildren(o, AddressRangeIdentity, &list, info)
	return list, err
}

// CreateAddressRange creates a new child AddressRange under the SubnetTemplate
func (o *SubnetTemplate) CreateAddressRange(child *AddressRange) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Metadatas retrieves the list of child Metadatas of the SubnetTemplate
func (o *SubnetTemplate) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the SubnetTemplate
func (o *SubnetTemplate) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the SubnetTemplate
func (o *SubnetTemplate) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the SubnetTemplate
func (o *SubnetTemplate) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// QOSs retrieves the list of child QOSs of the SubnetTemplate
func (o *SubnetTemplate) QOSs(info *bambou.FetchingInfo) (QOSsList, *bambou.Error) {

	var list QOSsList
	err := bambou.CurrentSession().FetchChildren(o, QOSIdentity, &list, info)
	return list, err
}

// CreateQOS creates a new child QOS under the SubnetTemplate
func (o *SubnetTemplate) CreateQOS(child *QOS) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Subnets retrieves the list of child Subnets of the SubnetTemplate
func (o *SubnetTemplate) Subnets(info *bambou.FetchingInfo) (SubnetsList, *bambou.Error) {

	var list SubnetsList
	err := bambou.CurrentSession().FetchChildren(o, SubnetIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the SubnetTemplate
func (o *SubnetTemplate) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
