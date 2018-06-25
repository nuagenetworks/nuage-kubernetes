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

// PGExpressionTemplateIdentity represents the Identity of the object
var PGExpressionTemplateIdentity = bambou.Identity{
	Name:     "pgexpressiontemplate",
	Category: "pgexpressiontemplates",
}

// PGExpressionTemplatesList represents a list of PGExpressionTemplates
type PGExpressionTemplatesList []*PGExpressionTemplate

// PGExpressionTemplatesAncestor is the interface that an ancestor of a PGExpressionTemplate must implement.
// An Ancestor is defined as an entity that has PGExpressionTemplate as a descendant.
// An Ancestor can get a list of its child PGExpressionTemplates, but not necessarily create one.
type PGExpressionTemplatesAncestor interface {
	PGExpressionTemplates(*bambou.FetchingInfo) (PGExpressionTemplatesList, *bambou.Error)
}

// PGExpressionTemplatesParent is the interface that a parent of a PGExpressionTemplate must implement.
// A Parent is defined as an entity that has PGExpressionTemplate as a child.
// A Parent is an Ancestor which can create a PGExpressionTemplate.
type PGExpressionTemplatesParent interface {
	PGExpressionTemplatesAncestor
	CreatePGExpressionTemplate(*PGExpressionTemplate) *bambou.Error
}

// PGExpressionTemplate represents the model of a pgexpressiontemplate
type PGExpressionTemplate struct {
	ID            string `json:"ID,omitempty"`
	ParentID      string `json:"parentID,omitempty"`
	ParentType    string `json:"parentType,omitempty"`
	Owner         string `json:"owner,omitempty"`
	Name          string `json:"name,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	Description   string `json:"description,omitempty"`
	EntityScope   string `json:"entityScope,omitempty"`
	Expression    string `json:"expression,omitempty"`
	ExternalID    string `json:"externalID,omitempty"`
}

// NewPGExpressionTemplate returns a new *PGExpressionTemplate
func NewPGExpressionTemplate() *PGExpressionTemplate {

	return &PGExpressionTemplate{}
}

// Identity returns the Identity of the object.
func (o *PGExpressionTemplate) Identity() bambou.Identity {

	return PGExpressionTemplateIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *PGExpressionTemplate) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *PGExpressionTemplate) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the PGExpressionTemplate from the server
func (o *PGExpressionTemplate) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the PGExpressionTemplate into the server
func (o *PGExpressionTemplate) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the PGExpressionTemplate from the server
func (o *PGExpressionTemplate) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}
