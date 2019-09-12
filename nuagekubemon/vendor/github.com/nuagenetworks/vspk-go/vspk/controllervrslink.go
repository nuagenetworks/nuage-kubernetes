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

// ControllerVRSLinkIdentity represents the Identity of the object
var ControllerVRSLinkIdentity = bambou.Identity{
	Name:     "controllervrslink",
	Category: "controllervrslinks",
}

// ControllerVRSLinksList represents a list of ControllerVRSLinks
type ControllerVRSLinksList []*ControllerVRSLink

// ControllerVRSLinksAncestor is the interface that an ancestor of a ControllerVRSLink must implement.
// An Ancestor is defined as an entity that has ControllerVRSLink as a descendant.
// An Ancestor can get a list of its child ControllerVRSLinks, but not necessarily create one.
type ControllerVRSLinksAncestor interface {
	ControllerVRSLinks(*bambou.FetchingInfo) (ControllerVRSLinksList, *bambou.Error)
}

// ControllerVRSLinksParent is the interface that a parent of a ControllerVRSLink must implement.
// A Parent is defined as an entity that has ControllerVRSLink as a child.
// A Parent is an Ancestor which can create a ControllerVRSLink.
type ControllerVRSLinksParent interface {
	ControllerVRSLinksAncestor
	CreateControllerVRSLink(*ControllerVRSLink) *bambou.Error
}

// ControllerVRSLink represents the model of a controllervrslink
type ControllerVRSLink struct {
	ID                     string        `json:"ID,omitempty"`
	ParentID               string        `json:"parentID,omitempty"`
	ParentType             string        `json:"parentType,omitempty"`
	Owner                  string        `json:"owner,omitempty"`
	VRSID                  string        `json:"VRSID,omitempty"`
	VRSPersonality         string        `json:"VRSPersonality,omitempty"`
	VSCConfigState         string        `json:"VSCConfigState,omitempty"`
	VSCCurrentState        string        `json:"VSCCurrentState,omitempty"`
	JSONRPCConnectionState string        `json:"JSONRPCConnectionState,omitempty"`
	Name                   string        `json:"name,omitempty"`
	LastUpdatedBy          string        `json:"lastUpdatedBy,omitempty"`
	Peer                   string        `json:"peer,omitempty"`
	ClusterNodeRole        string        `json:"clusterNodeRole,omitempty"`
	EmbeddedMetadata       []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope            string        `json:"entityScope,omitempty"`
	Role                   string        `json:"role,omitempty"`
	Connections            []interface{} `json:"connections,omitempty"`
	ControllerID           string        `json:"controllerID,omitempty"`
	ControllerType         string        `json:"controllerType,omitempty"`
	Status                 string        `json:"status,omitempty"`
	ExternalID             string        `json:"externalID,omitempty"`
	Dynamic                bool          `json:"dynamic"`
}

// NewControllerVRSLink returns a new *ControllerVRSLink
func NewControllerVRSLink() *ControllerVRSLink {

	return &ControllerVRSLink{}
}

// Identity returns the Identity of the object.
func (o *ControllerVRSLink) Identity() bambou.Identity {

	return ControllerVRSLinkIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *ControllerVRSLink) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *ControllerVRSLink) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the ControllerVRSLink from the server
func (o *ControllerVRSLink) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the ControllerVRSLink into the server
func (o *ControllerVRSLink) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the ControllerVRSLink from the server
func (o *ControllerVRSLink) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the ControllerVRSLink
func (o *ControllerVRSLink) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the ControllerVRSLink
func (o *ControllerVRSLink) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the ControllerVRSLink
func (o *ControllerVRSLink) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the ControllerVRSLink
func (o *ControllerVRSLink) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// VRSs retrieves the list of child VRSs of the ControllerVRSLink
func (o *ControllerVRSLink) VRSs(info *bambou.FetchingInfo) (VRSsList, *bambou.Error) {

	var list VRSsList
	err := bambou.CurrentSession().FetchChildren(o, VRSIdentity, &list, info)
	return list, err
}

// HSCs retrieves the list of child HSCs of the ControllerVRSLink
func (o *ControllerVRSLink) HSCs(info *bambou.FetchingInfo) (HSCsList, *bambou.Error) {

	var list HSCsList
	err := bambou.CurrentSession().FetchChildren(o, HSCIdentity, &list, info)
	return list, err
}

// VSCs retrieves the list of child VSCs of the ControllerVRSLink
func (o *ControllerVRSLink) VSCs(info *bambou.FetchingInfo) (VSCsList, *bambou.Error) {

	var list VSCsList
	err := bambou.CurrentSession().FetchChildren(o, VSCIdentity, &list, info)
	return list, err
}
