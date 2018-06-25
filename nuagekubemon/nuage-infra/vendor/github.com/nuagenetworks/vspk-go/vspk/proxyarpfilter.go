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

// ProxyARPFilterIdentity represents the Identity of the object
var ProxyARPFilterIdentity = bambou.Identity{
	Name:     "proxyarpfilter",
	Category: "proxyarpfilters",
}

// ProxyARPFiltersList represents a list of ProxyARPFilters
type ProxyARPFiltersList []*ProxyARPFilter

// ProxyARPFiltersAncestor is the interface that an ancestor of a ProxyARPFilter must implement.
// An Ancestor is defined as an entity that has ProxyARPFilter as a descendant.
// An Ancestor can get a list of its child ProxyARPFilters, but not necessarily create one.
type ProxyARPFiltersAncestor interface {
	ProxyARPFilters(*bambou.FetchingInfo) (ProxyARPFiltersList, *bambou.Error)
}

// ProxyARPFiltersParent is the interface that a parent of a ProxyARPFilter must implement.
// A Parent is defined as an entity that has ProxyARPFilter as a child.
// A Parent is an Ancestor which can create a ProxyARPFilter.
type ProxyARPFiltersParent interface {
	ProxyARPFiltersAncestor
	CreateProxyARPFilter(*ProxyARPFilter) *bambou.Error
}

// ProxyARPFilter represents the model of a proxyarpfilter
type ProxyARPFilter struct {
	ID            string `json:"ID,omitempty"`
	ParentID      string `json:"parentID,omitempty"`
	ParentType    string `json:"parentType,omitempty"`
	Owner         string `json:"owner,omitempty"`
	IPType        string `json:"IPType,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	MaxAddress    string `json:"maxAddress,omitempty"`
	MinAddress    string `json:"minAddress,omitempty"`
	EntityScope   string `json:"entityScope,omitempty"`
	ExternalID    string `json:"externalID,omitempty"`
}

// NewProxyARPFilter returns a new *ProxyARPFilter
func NewProxyARPFilter() *ProxyARPFilter {

	return &ProxyARPFilter{
		IPType: "IPV4",
	}
}

// Identity returns the Identity of the object.
func (o *ProxyARPFilter) Identity() bambou.Identity {

	return ProxyARPFilterIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *ProxyARPFilter) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *ProxyARPFilter) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the ProxyARPFilter from the server
func (o *ProxyARPFilter) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the ProxyARPFilter into the server
func (o *ProxyARPFilter) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the ProxyARPFilter from the server
func (o *ProxyARPFilter) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// EventLogs retrieves the list of child EventLogs of the ProxyARPFilter
func (o *ProxyARPFilter) EventLogs(info *bambou.FetchingInfo) (EventLogsList, *bambou.Error) {

	var list EventLogsList
	err := bambou.CurrentSession().FetchChildren(o, EventLogIdentity, &list, info)
	return list, err
}
