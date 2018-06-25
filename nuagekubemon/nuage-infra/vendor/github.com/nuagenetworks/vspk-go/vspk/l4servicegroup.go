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

// L4ServiceGroupIdentity represents the Identity of the object
var L4ServiceGroupIdentity = bambou.Identity{
	Name:     "l4servicegroup",
	Category: "l4servicegroups",
}

// L4ServiceGroupsList represents a list of L4ServiceGroups
type L4ServiceGroupsList []*L4ServiceGroup

// L4ServiceGroupsAncestor is the interface that an ancestor of a L4ServiceGroup must implement.
// An Ancestor is defined as an entity that has L4ServiceGroup as a descendant.
// An Ancestor can get a list of its child L4ServiceGroups, but not necessarily create one.
type L4ServiceGroupsAncestor interface {
	L4ServiceGroups(*bambou.FetchingInfo) (L4ServiceGroupsList, *bambou.Error)
}

// L4ServiceGroupsParent is the interface that a parent of a L4ServiceGroup must implement.
// A Parent is defined as an entity that has L4ServiceGroup as a child.
// A Parent is an Ancestor which can create a L4ServiceGroup.
type L4ServiceGroupsParent interface {
	L4ServiceGroupsAncestor
	CreateL4ServiceGroup(*L4ServiceGroup) *bambou.Error
}

// L4ServiceGroup represents the model of a l4servicegroup
type L4ServiceGroup struct {
	ID            string `json:"ID,omitempty"`
	ParentID      string `json:"parentID,omitempty"`
	ParentType    string `json:"parentType,omitempty"`
	Owner         string `json:"owner,omitempty"`
	Name          string `json:"name,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	Description   string `json:"description,omitempty"`
	EntityScope   string `json:"entityScope,omitempty"`
	ExternalID    string `json:"externalID,omitempty"`
}

// NewL4ServiceGroup returns a new *L4ServiceGroup
func NewL4ServiceGroup() *L4ServiceGroup {

	return &L4ServiceGroup{}
}

// Identity returns the Identity of the object.
func (o *L4ServiceGroup) Identity() bambou.Identity {

	return L4ServiceGroupIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *L4ServiceGroup) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *L4ServiceGroup) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the L4ServiceGroup from the server
func (o *L4ServiceGroup) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the L4ServiceGroup into the server
func (o *L4ServiceGroup) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the L4ServiceGroup from the server
func (o *L4ServiceGroup) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// L4Services retrieves the list of child L4Services of the L4ServiceGroup
func (o *L4ServiceGroup) L4Services(info *bambou.FetchingInfo) (L4ServicesList, *bambou.Error) {

	var list L4ServicesList
	err := bambou.CurrentSession().FetchChildren(o, L4ServiceIdentity, &list, info)
	return list, err
}

// AssignL4Services assigns the list of L4Services to the L4ServiceGroup
func (o *L4ServiceGroup) AssignL4Services(children L4ServicesList) *bambou.Error {

	list := []bambou.Identifiable{}
	for _, c := range children {
		list = append(list, c)
	}

	return bambou.CurrentSession().AssignChildren(o, list, L4ServiceIdentity)
}
