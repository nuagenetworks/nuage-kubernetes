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

// RateLimiterIdentity represents the Identity of the object
var RateLimiterIdentity = bambou.Identity{
	Name:     "ratelimiter",
	Category: "ratelimiters",
}

// RateLimitersList represents a list of RateLimiters
type RateLimitersList []*RateLimiter

// RateLimitersAncestor is the interface that an ancestor of a RateLimiter must implement.
// An Ancestor is defined as an entity that has RateLimiter as a descendant.
// An Ancestor can get a list of its child RateLimiters, but not necessarily create one.
type RateLimitersAncestor interface {
	RateLimiters(*bambou.FetchingInfo) (RateLimitersList, *bambou.Error)
}

// RateLimitersParent is the interface that a parent of a RateLimiter must implement.
// A Parent is defined as an entity that has RateLimiter as a child.
// A Parent is an Ancestor which can create a RateLimiter.
type RateLimitersParent interface {
	RateLimitersAncestor
	CreateRateLimiter(*RateLimiter) *bambou.Error
}

// RateLimiter represents the model of a ratelimiter
type RateLimiter struct {
	ID                       string `json:"ID,omitempty"`
	ParentID                 string `json:"parentID,omitempty"`
	ParentType               string `json:"parentType,omitempty"`
	Owner                    string `json:"owner,omitempty"`
	Name                     string `json:"name,omitempty"`
	LastUpdatedBy            string `json:"lastUpdatedBy,omitempty"`
	PeakBurstSize            string `json:"peakBurstSize,omitempty"`
	PeakInformationRate      string `json:"peakInformationRate,omitempty"`
	Description              string `json:"description,omitempty"`
	EntityScope              string `json:"entityScope,omitempty"`
	CommittedInformationRate string `json:"committedInformationRate,omitempty"`
	ExternalID               string `json:"externalID,omitempty"`
}

// NewRateLimiter returns a new *RateLimiter
func NewRateLimiter() *RateLimiter {

	return &RateLimiter{}
}

// Identity returns the Identity of the object.
func (o *RateLimiter) Identity() bambou.Identity {

	return RateLimiterIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *RateLimiter) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *RateLimiter) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the RateLimiter from the server
func (o *RateLimiter) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the RateLimiter into the server
func (o *RateLimiter) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the RateLimiter from the server
func (o *RateLimiter) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the RateLimiter
func (o *RateLimiter) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the RateLimiter
func (o *RateLimiter) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the RateLimiter
func (o *RateLimiter) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the RateLimiter
func (o *RateLimiter) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
