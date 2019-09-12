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

// TestSuiteIdentity represents the Identity of the object
var TestSuiteIdentity = bambou.Identity{
	Name:     "testsuite",
	Category: "testsuites",
}

// TestSuitesList represents a list of TestSuites
type TestSuitesList []*TestSuite

// TestSuitesAncestor is the interface that an ancestor of a TestSuite must implement.
// An Ancestor is defined as an entity that has TestSuite as a descendant.
// An Ancestor can get a list of its child TestSuites, but not necessarily create one.
type TestSuitesAncestor interface {
	TestSuites(*bambou.FetchingInfo) (TestSuitesList, *bambou.Error)
}

// TestSuitesParent is the interface that a parent of a TestSuite must implement.
// A Parent is defined as an entity that has TestSuite as a child.
// A Parent is an Ancestor which can create a TestSuite.
type TestSuitesParent interface {
	TestSuitesAncestor
	CreateTestSuite(*TestSuite) *bambou.Error
}

// TestSuite represents the model of a testsuite
type TestSuite struct {
	ID               string        `json:"ID,omitempty"`
	ParentID         string        `json:"parentID,omitempty"`
	ParentType       string        `json:"parentType,omitempty"`
	Owner            string        `json:"owner,omitempty"`
	Name             string        `json:"name,omitempty"`
	LastUpdatedBy    string        `json:"lastUpdatedBy,omitempty"`
	Description      string        `json:"description,omitempty"`
	EmbeddedMetadata []interface{} `json:"embeddedMetadata,omitempty"`
	EnterpriseID     string        `json:"enterpriseID,omitempty"`
	EntityScope      string        `json:"entityScope,omitempty"`
	ExternalID       string        `json:"externalID,omitempty"`
}

// NewTestSuite returns a new *TestSuite
func NewTestSuite() *TestSuite {

	return &TestSuite{}
}

// Identity returns the Identity of the object.
func (o *TestSuite) Identity() bambou.Identity {

	return TestSuiteIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *TestSuite) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *TestSuite) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the TestSuite from the server
func (o *TestSuite) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the TestSuite into the server
func (o *TestSuite) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the TestSuite from the server
func (o *TestSuite) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// Tests retrieves the list of child Tests of the TestSuite
func (o *TestSuite) Tests(info *bambou.FetchingInfo) (TestsList, *bambou.Error) {

	var list TestsList
	err := bambou.CurrentSession().FetchChildren(o, TestIdentity, &list, info)
	return list, err
}

// CreateTest creates a new child Test under the TestSuite
func (o *TestSuite) CreateTest(child *Test) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// TestSuiteRuns retrieves the list of child TestSuiteRuns of the TestSuite
func (o *TestSuite) TestSuiteRuns(info *bambou.FetchingInfo) (TestSuiteRunsList, *bambou.Error) {

	var list TestSuiteRunsList
	err := bambou.CurrentSession().FetchChildren(o, TestSuiteRunIdentity, &list, info)
	return list, err
}

// Metadatas retrieves the list of child Metadatas of the TestSuite
func (o *TestSuite) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the TestSuite
func (o *TestSuite) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the TestSuite
func (o *TestSuite) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the TestSuite
func (o *TestSuite) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
