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

// TestSuiteRunIdentity represents the Identity of the object
var TestSuiteRunIdentity = bambou.Identity{
	Name:     "testsuiterun",
	Category: "testsuiteruns",
}

// TestSuiteRunsList represents a list of TestSuiteRuns
type TestSuiteRunsList []*TestSuiteRun

// TestSuiteRunsAncestor is the interface that an ancestor of a TestSuiteRun must implement.
// An Ancestor is defined as an entity that has TestSuiteRun as a descendant.
// An Ancestor can get a list of its child TestSuiteRuns, but not necessarily create one.
type TestSuiteRunsAncestor interface {
	TestSuiteRuns(*bambou.FetchingInfo) (TestSuiteRunsList, *bambou.Error)
}

// TestSuiteRunsParent is the interface that a parent of a TestSuiteRun must implement.
// A Parent is defined as an entity that has TestSuiteRun as a child.
// A Parent is an Ancestor which can create a TestSuiteRun.
type TestSuiteRunsParent interface {
	TestSuiteRunsAncestor
	CreateTestSuiteRun(*TestSuiteRun) *bambou.Error
}

// TestSuiteRun represents the model of a testsuiterun
type TestSuiteRun struct {
	ID                      string        `json:"ID,omitempty"`
	ParentID                string        `json:"parentID,omitempty"`
	ParentType              string        `json:"parentType,omitempty"`
	Owner                   string        `json:"owner,omitempty"`
	VPortName               string        `json:"VPortName,omitempty"`
	NSGatewayName           string        `json:"NSGatewayName,omitempty"`
	LastUpdatedBy           string        `json:"lastUpdatedBy,omitempty"`
	DatapathID              string        `json:"datapathID,omitempty"`
	Destination             string        `json:"destination,omitempty"`
	EmbeddedMetadata        []interface{} `json:"embeddedMetadata,omitempty"`
	EntityScope             string        `json:"entityScope,omitempty"`
	DomainName              string        `json:"domainName,omitempty"`
	ZoneName                string        `json:"zoneName,omitempty"`
	OperationStatus         string        `json:"operationStatus,omitempty"`
	AssociatedEntityType    string        `json:"associatedEntityType,omitempty"`
	AssociatedTestSuiteID   string        `json:"associatedTestSuiteID,omitempty"`
	AssociatedTestSuiteName string        `json:"associatedTestSuiteName,omitempty"`
	SubnetName              string        `json:"subnetName,omitempty"`
	ExternalID              string        `json:"externalID,omitempty"`
}

// NewTestSuiteRun returns a new *TestSuiteRun
func NewTestSuiteRun() *TestSuiteRun {

	return &TestSuiteRun{}
}

// Identity returns the Identity of the object.
func (o *TestSuiteRun) Identity() bambou.Identity {

	return TestSuiteRunIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *TestSuiteRun) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *TestSuiteRun) SetIdentifier(ID string) {

	o.ID = ID
}

// Fetch retrieves the TestSuiteRun from the server
func (o *TestSuiteRun) Fetch() *bambou.Error {

	return bambou.CurrentSession().FetchEntity(o)
}

// Save saves the TestSuiteRun into the server
func (o *TestSuiteRun) Save() *bambou.Error {

	return bambou.CurrentSession().SaveEntity(o)
}

// Delete deletes the TestSuiteRun from the server
func (o *TestSuiteRun) Delete() *bambou.Error {

	return bambou.CurrentSession().DeleteEntity(o)
}

// TestRuns retrieves the list of child TestRuns of the TestSuiteRun
func (o *TestSuiteRun) TestRuns(info *bambou.FetchingInfo) (TestRunsList, *bambou.Error) {

	var list TestRunsList
	err := bambou.CurrentSession().FetchChildren(o, TestRunIdentity, &list, info)
	return list, err
}

// Metadatas retrieves the list of child Metadatas of the TestSuiteRun
func (o *TestSuiteRun) Metadatas(info *bambou.FetchingInfo) (MetadatasList, *bambou.Error) {

	var list MetadatasList
	err := bambou.CurrentSession().FetchChildren(o, MetadataIdentity, &list, info)
	return list, err
}

// CreateMetadata creates a new child Metadata under the TestSuiteRun
func (o *TestSuiteRun) CreateMetadata(child *Metadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}

// GlobalMetadatas retrieves the list of child GlobalMetadatas of the TestSuiteRun
func (o *TestSuiteRun) GlobalMetadatas(info *bambou.FetchingInfo) (GlobalMetadatasList, *bambou.Error) {

	var list GlobalMetadatasList
	err := bambou.CurrentSession().FetchChildren(o, GlobalMetadataIdentity, &list, info)
	return list, err
}

// CreateGlobalMetadata creates a new child GlobalMetadata under the TestSuiteRun
func (o *TestSuiteRun) CreateGlobalMetadata(child *GlobalMetadata) *bambou.Error {

	return bambou.CurrentSession().CreateChild(o, child)
}
