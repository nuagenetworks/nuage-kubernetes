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

// LDAPConfigurationIdentity represents the Identity of the object
var LDAPConfigurationIdentity = bambou.Identity{
	Name:     "ldapconfiguration",
	Category: "ldapconfigurations",
}

// LDAPConfigurationsList represents a list of LDAPConfigurations
type LDAPConfigurationsList []*LDAPConfiguration

// LDAPConfigurationsAncestor is the interface of an ancestor of a LDAPConfiguration must implement.
type LDAPConfigurationsAncestor interface {
	LDAPConfigurations(*bambou.FetchingInfo) (LDAPConfigurationsList, *bambou.Error)
	CreateLDAPConfigurations(*LDAPConfiguration) *bambou.Error
}

// LDAPConfiguration represents the model of a ldapconfiguration
type LDAPConfiguration struct {
	ID                    string `json:"ID,omitempty"`
	ParentID              string `json:"parentID,omitempty"`
	ParentType            string `json:"parentType,omitempty"`
	Owner                 string `json:"owner,omitempty"`
	SSLEnabled            bool   `json:"SSLEnabled"`
	Password              string `json:"password,omitempty"`
	LastUpdatedBy         string `json:"lastUpdatedBy,omitempty"`
	AcceptAllCertificates bool   `json:"acceptAllCertificates"`
	Certificate           string `json:"certificate,omitempty"`
	Server                string `json:"server,omitempty"`
	Enabled               bool   `json:"enabled"`
	EntityScope           string `json:"entityScope,omitempty"`
	Port                  string `json:"port,omitempty"`
	GroupDN               string `json:"groupDN,omitempty"`
	UserDNTemplate        string `json:"userDNTemplate,omitempty"`
	AuthorizationEnabled  bool   `json:"authorizationEnabled"`
	AuthorizingUserDN     string `json:"authorizingUserDN,omitempty"`
	ExternalID            string `json:"externalID,omitempty"`
}

// NewLDAPConfiguration returns a new *LDAPConfiguration
func NewLDAPConfiguration() *LDAPConfiguration {

	return &LDAPConfiguration{}
}

// Identity returns the Identity of the object.
func (o *LDAPConfiguration) Identity() bambou.Identity {

	return LDAPConfigurationIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *LDAPConfiguration) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *LDAPConfiguration) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the LDAPConfiguration from the server
func (o *LDAPConfiguration) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the LDAPConfiguration into the server
func (o *LDAPConfiguration) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the LDAPConfiguration from the server
func (o *LDAPConfiguration) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the LDAPConfiguration
func (o *LDAPConfiguration) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the LDAPConfiguration
func (o *LDAPConfiguration) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the LDAPConfiguration
func (o *LDAPConfiguration) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the LDAPConfiguration
func (o *LDAPConfiguration) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
