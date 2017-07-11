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

// StatisticsPolicyIdentity represents the Identity of the object
var StatisticsPolicyIdentity = bambou.Identity{
	Name:     "statisticspolicy",
	Category: "statisticspolicies",
}

// StatisticsPoliciesList represents a list of StatisticsPolicies
type StatisticsPoliciesList []*StatisticsPolicy

// StatisticsPoliciesAncestor is the interface of an ancestor of a StatisticsPolicy must implement.
type StatisticsPoliciesAncestor interface {
	StatisticsPolicies(*bambou.FetchingInfo) (StatisticsPoliciesList, *bambou.Error)
	CreateStatisticsPolicies(*StatisticsPolicy) *bambou.Error
}

// StatisticsPolicy represents the model of a statisticspolicy
type StatisticsPolicy struct {
	ID                      string `json:"ID,omitempty"`
	ParentID                string `json:"parentID,omitempty"`
	ParentType              string `json:"parentType,omitempty"`
	Owner                   string `json:"owner,omitempty"`
	Name                    string `json:"name,omitempty"`
	LastUpdatedBy           string `json:"lastUpdatedBy,omitempty"`
	DataCollectionFrequency int    `json:"dataCollectionFrequency,omitempty"`
	Description             string `json:"description,omitempty"`
	EntityScope             string `json:"entityScope,omitempty"`
	ExternalID              string `json:"externalID,omitempty"`
}

// NewStatisticsPolicy returns a new *StatisticsPolicy
func NewStatisticsPolicy() *StatisticsPolicy {

	return &StatisticsPolicy{}
}

// Identity returns the Identity of the object.
func (o *StatisticsPolicy) Identity() bambou.Identity {

	return StatisticsPolicyIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *StatisticsPolicy) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *StatisticsPolicy) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the StatisticsPolicy from the server
func (o *StatisticsPolicy) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the StatisticsPolicy into the server
func (o *StatisticsPolicy) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the StatisticsPolicy from the server
func (o *StatisticsPolicy) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Metadatas retrieves the list of child Metadatas of the StatisticsPolicy
func (o *StatisticsPolicy) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the StatisticsPolicy
func (o *StatisticsPolicy) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the StatisticsPolicy
func (o *StatisticsPolicy) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the StatisticsPolicy
func (o *StatisticsPolicy) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
