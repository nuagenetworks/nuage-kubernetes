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

// PGExpressionIdentity represents the Identity of the object
var PGExpressionIdentity = bambou.Identity{
	Name:     "pgexpression",
	Category: "pgexpressions",
}

// PGExpressionsList represents a list of PGExpressions
type PGExpressionsList []*PGExpression

// PGExpressionsAncestor is the interface that an ancestor of a PGExpression must implement.
// An Ancestor is defined as an entity that has PGExpression as a descendant.
// An Ancestor can get a list of its child PGExpressions, but not necessarily create one.
type PGExpressionsAncestor interface {
	PGExpressions(*bambou.FetchingInfo) (PGExpressionsList, *bambou.Error)
}

// PGExpressionsParent is the interface that a parent of a PGExpression must implement.
// A Parent is defined as an entity that has PGExpression as a child.
// A Parent is an Ancestor which can create a PGExpression.
type PGExpressionsParent interface {
	PGExpressionsAncestor
	CreatePGExpression(*PGExpression) *bambou.Error
}

// PGExpression represents the model of a pgexpression
type PGExpression struct {
	ID            string `json:"ID,omitempty"`
	ParentID      string `json:"parentID,omitempty"`
	ParentType    string `json:"parentType,omitempty"`
	Owner         string `json:"owner,omitempty"`
	Name          string `json:"name,omitempty"`
	LastUpdatedBy string `json:"lastUpdatedBy,omitempty"`
	TemplateID    string `json:"templateID,omitempty"`
	Description   string `json:"description,omitempty"`
	EntityScope   string `json:"entityScope,omitempty"`
	Expression    string `json:"expression,omitempty"`
	ExternalID    string `json:"externalID,omitempty"`
}

// NewPGExpression returns a new *PGExpression
func NewPGExpression() *PGExpression {

	return &PGExpression{}
}

// Identity returns the Identity of the object.
func (o *PGExpression) Identity() bambou.Identity {

	return PGExpressionIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *PGExpression) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *PGExpression) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the PGExpression from the server
func (o *PGExpression) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the PGExpression into the server
func (o *PGExpression) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the PGExpression from the server
func (o *PGExpression) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}
