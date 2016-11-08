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

// KeyServerNotificationIdentity represents the Identity of the object
var KeyServerNotificationIdentity = bambou.Identity{
	Name:     "keyservernotification",
	Category: "keyservernotifications",
}

// KeyServerNotificationsList represents a list of KeyServerNotifications
type KeyServerNotificationsList []*KeyServerNotification

// KeyServerNotificationsAncestor is the interface of an ancestor of a KeyServerNotification must implement.
type KeyServerNotificationsAncestor interface {
	KeyServerNotifications(*bambou.FetchingInfo) (KeyServerNotificationsList, *bambou.Error)
	CreateKeyServerNotifications(*KeyServerNotification) *bambou.Error
}

// KeyServerNotification represents the model of a keyservernotification
type KeyServerNotification struct {
	ID               string      `json:"ID,omitempty"`
	ParentID         string      `json:"parentID,omitempty"`
	ParentType       string      `json:"parentType,omitempty"`
	Owner            string      `json:"owner,omitempty"`
	Base64JSONString string      `json:"base64JSONString,omitempty"`
	Message          interface{} `json:"message,omitempty"`
	EntityScope      string      `json:"entityScope,omitempty"`
	NotificationType string      `json:"notificationType,omitempty"`
	ExternalID       string      `json:"externalID,omitempty"`
}

// NewKeyServerNotification returns a new *KeyServerNotification
func NewKeyServerNotification() *KeyServerNotification {

	return &KeyServerNotification{}
}

// Identity returns the Identity of the object.
func (o *KeyServerNotification) Identity() bambou.Identity {

	return KeyServerNotificationIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *KeyServerNotification) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *KeyServerNotification) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the KeyServerNotification from the server
func (o *KeyServerNotification) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the KeyServerNotification into the server
func (o *KeyServerNotification) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the KeyServerNotification from the server
func (o *KeyServerNotification) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the KeyServerNotification
func (o *KeyServerNotification) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the KeyServerNotification
func (o *KeyServerNotification) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the KeyServerNotification
func (o *KeyServerNotification) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the KeyServerNotification
func (o *KeyServerNotification) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
