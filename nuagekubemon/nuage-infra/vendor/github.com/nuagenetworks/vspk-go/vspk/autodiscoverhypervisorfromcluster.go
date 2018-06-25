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

// AutoDiscoverHypervisorFromClusterIdentity represents the Identity of the object
var AutoDiscoverHypervisorFromClusterIdentity = bambou.Identity{
	Name:     "autodiscoveredhypervisor",
	Category: "autodiscoveredhypervisors",
}

// AutoDiscoverHypervisorFromClustersList represents a list of AutoDiscoverHypervisorFromClusters
type AutoDiscoverHypervisorFromClustersList []*AutoDiscoverHypervisorFromCluster

// AutoDiscoverHypervisorFromClustersAncestor is the interface that an ancestor of a AutoDiscoverHypervisorFromCluster must implement.
// An Ancestor is defined as an entity that has AutoDiscoverHypervisorFromCluster as a descendant.
// An Ancestor can get a list of its child AutoDiscoverHypervisorFromClusters, but not necessarily create one.
type AutoDiscoverHypervisorFromClustersAncestor interface {
	AutoDiscoverHypervisorFromClusters(*bambou.FetchingInfo) (AutoDiscoverHypervisorFromClustersList, *bambou.Error)
}

// AutoDiscoverHypervisorFromClustersParent is the interface that a parent of a AutoDiscoverHypervisorFromCluster must implement.
// A Parent is defined as an entity that has AutoDiscoverHypervisorFromCluster as a child.
// A Parent is an Ancestor which can create a AutoDiscoverHypervisorFromCluster.
type AutoDiscoverHypervisorFromClustersParent interface {
	AutoDiscoverHypervisorFromClustersAncestor
	CreateAutoDiscoverHypervisorFromCluster(*AutoDiscoverHypervisorFromCluster) *bambou.Error
}

// AutoDiscoverHypervisorFromCluster represents the model of a autodiscoveredhypervisor
type AutoDiscoverHypervisorFromCluster struct {
	ID            string        `json:"ID,omitempty"`
	ParentID      string        `json:"parentID,omitempty"`
	ParentType    string        `json:"parentType,omitempty"`
	Owner         string        `json:"owner,omitempty"`
	LastUpdatedBy string        `json:"lastUpdatedBy,omitempty"`
	NetworkList   []interface{} `json:"networkList,omitempty"`
	EntityScope   string        `json:"entityScope,omitempty"`
	AssocEntityID string        `json:"assocEntityID,omitempty"`
	ExternalID    string        `json:"externalID,omitempty"`
	HypervisorIP  string        `json:"hypervisorIP,omitempty"`
}

// NewAutoDiscoverHypervisorFromCluster returns a new *AutoDiscoverHypervisorFromCluster
func NewAutoDiscoverHypervisorFromCluster() *AutoDiscoverHypervisorFromCluster {

	return &AutoDiscoverHypervisorFromCluster{}
}

// Identity returns the Identity of the object.
func (o *AutoDiscoverHypervisorFromCluster) Identity() bambou.Identity {

	return AutoDiscoverHypervisorFromClusterIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *AutoDiscoverHypervisorFromCluster) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *AutoDiscoverHypervisorFromCluster) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the AutoDiscoverHypervisorFromCluster from the server
func (o *AutoDiscoverHypervisorFromCluster) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the AutoDiscoverHypervisorFromCluster into the server
func (o *AutoDiscoverHypervisorFromCluster) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the AutoDiscoverHypervisorFromCluster from the server
func (o *AutoDiscoverHypervisorFromCluster) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}
