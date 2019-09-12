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

// EnterpriseProfileIdentity represents the Identity of the object
var EnterpriseProfileIdentity = bambou.Identity{
	Name:     "enterpriseprofile",
	Category: "enterpriseprofiles",
}

// EnterpriseProfilesList represents a list of EnterpriseProfiles
type EnterpriseProfilesList []*EnterpriseProfile

// EnterpriseProfilesAncestor is the interface that an ancestor of a EnterpriseProfile must implement.
// An Ancestor is defined as an entity that has EnterpriseProfile as a descendant.
// An Ancestor can get a list of its child EnterpriseProfiles, but not necessarily create one.
type EnterpriseProfilesAncestor interface {
	EnterpriseProfiles(*bambou.FetchingInfo) (EnterpriseProfilesList, *bambou.Error)
}

// EnterpriseProfilesParent is the interface that a parent of a EnterpriseProfile must implement.
// A Parent is defined as an entity that has EnterpriseProfile as a child.
// A Parent is an Ancestor which can create a EnterpriseProfile.
type EnterpriseProfilesParent interface {
	EnterpriseProfilesAncestor
	CreateEnterpriseProfile(*EnterpriseProfile) *bambou.Error
}

// EnterpriseProfile represents the model of a enterpriseprofile
type EnterpriseProfile struct {
	ID                                     string        `json:"ID,omitempty"`
	ParentID                               string        `json:"parentID,omitempty"`
	ParentType                             string        `json:"parentType,omitempty"`
	Owner                                  string        `json:"owner,omitempty"`
	BGPEnabled                             bool          `json:"BGPEnabled"`
	DHCPLeaseInterval                      int           `json:"DHCPLeaseInterval,omitempty"`
	VNFManagementEnabled                   bool          `json:"VNFManagementEnabled"`
	Name                                   string        `json:"name,omitempty"`
	LastUpdatedBy                          string        `json:"lastUpdatedBy,omitempty"`
	WebFilterEnabled                       bool          `json:"webFilterEnabled"`
	ReceiveMultiCastListID                 string        `json:"receiveMultiCastListID,omitempty"`
	SendMultiCastListID                    string        `json:"sendMultiCastListID,omitempty"`
	Description                            string        `json:"description,omitempty"`
	AllowAdvancedQOSConfiguration          bool          `json:"allowAdvancedQOSConfiguration"`
	AllowGatewayManagement                 bool          `json:"allowGatewayManagement"`
	AllowTrustedForwardingClass            bool          `json:"allowTrustedForwardingClass"`
	AllowedForwardingClasses               []interface{} `json:"allowedForwardingClasses,omitempty"`
	FloatingIPsQuota                       int           `json:"floatingIPsQuota,omitempty"`
	EmbeddedMetadata                       []interface{} `json:"embeddedMetadata,omitempty"`
	EnableApplicationPerformanceManagement bool          `json:"enableApplicationPerformanceManagement"`
	EncryptionManagementMode               string        `json:"encryptionManagementMode,omitempty"`
	EntityScope                            string        `json:"entityScope,omitempty"`
	ForwardingClass                        []interface{} `json:"forwardingClass,omitempty"`
	ExternalID                             string        `json:"externalID,omitempty"`
}

// NewEnterpriseProfile returns a new *EnterpriseProfile
func NewEnterpriseProfile() *EnterpriseProfile {

	return &EnterpriseProfile{
		DHCPLeaseInterval:                      24,
		VNFManagementEnabled:                   false,
		WebFilterEnabled:                       false,
		FloatingIPsQuota:                       100,
		EnableApplicationPerformanceManagement: false,
	}
}

// Identity returns the Identity of the object.
func (o *EnterpriseProfile) Identity() bambou.Identity {

	return EnterpriseProfileIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *EnterpriseProfile) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *EnterpriseProfile) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the EnterpriseProfile from the server
func (o *EnterpriseProfile) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the EnterpriseProfile into the server
func (o *EnterpriseProfile) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the EnterpriseProfile from the server
func (o *EnterpriseProfile) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the EnterpriseProfile
func (o *EnterpriseProfile) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the EnterpriseProfile
func (o *EnterpriseProfile) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the EnterpriseProfile
func (o *EnterpriseProfile) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the EnterpriseProfile
func (o *EnterpriseProfile) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// Enterprises retrieves the list of child Enterprises of the EnterpriseProfile
func (o *EnterpriseProfile) Enterprises(info *bambou.FetchingInfo) (EnterprisesList, *bambou.Error) {

	var list EnterprisesList
	err := bambou.CurrentSession().FetchChildren(o, EnterpriseIdentity, &list, info)
	return list, err
}

// MultiCastLists retrieves the list of child MultiCastLists of the EnterpriseProfile
func (o *EnterpriseProfile) MultiCastLists(info *bambou.FetchingInfo) (MultiCastListsList, *bambou.Error) {

	var list MultiCastListsList
	err := bambou.CurrentSession().FetchChildren(o, MultiCastListIdentity, &list, info)
	return list, err
}

// EventLogs retrieves the list of child EventLogs of the EnterpriseProfile
func (o *EnterpriseProfile) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
