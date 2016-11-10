// Copyright (c) 2015, Alcatel-Lucent Inc.
// All rights reserved.
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
// * Neither the name of bambou nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package bambou

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

var _currentSession Storer

// CurrentSession returns the current active and authenticated Session.
func CurrentSession() Storer {

	return _currentSession
}

// Storer is the interface that must be implemented by object that can
// perform CRUD operations on RemoteObjects.
type Storer interface {
	Start() *Error
	Reset()
	Root() Rootable

	FetchEntity(Identifiable) *Error
	SaveEntity(Identifiable) *Error
	DeleteEntity(Identifiable) *Error
	FetchChildren(Identifiable, Identity, interface{}, *FetchingInfo) *Error
	CreateChild(Identifiable, Identifiable) *Error
	AssignChildren(Identifiable, []Identifiable, Identity) *Error
	NextEvent(NotificationsChannel, string) *Error
}

// Session represents a user session. It provides the entire
// communication layer with the backend. It must implement the Operationable interface.
type Session struct {
	root         Rootable
	Certificate  string
	Username     string
	Password     string
	Organization string
	URL          string
	client       *http.Client
}

// NewSession returns a new *Session
// You need to provide a Rootable object that will be used to contain
// the results of the authentication process, like the api key for instance.
func NewSession(username, password, organization, url string, root Rootable) *Session {

	return &Session{
		Username:     username,
		Password:     password,
		Organization: organization,
		URL:          url,
		root:         root,
		client:       &http.Client{},
	}
}

// SetInsecureSkipVerify sets if the internal HTTP client should allow to connect
// to insecure API endpoints.
func (s *Session) SetInsecureSkipVerify(skip bool) *Error {

	if CurrentSession() != nil {
		return NewError(ErrorCodeSessionAlreadyStarted, "The session is already started. Stop it first")
	}

	s.client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: skip}}}

	return nil
}

func (s *Session) makeAuthorizationHeaders() (string, *Error) {

	if s.Username == "" {
		return "", NewError(ErrorCodeSessionUsernameNotSet, "Error while creating headers: User must be set")
	}

	if s.root == nil {
		return "", NewError(ErrorCodeSessionCannotForgeAuthToken, "Error while creating headers: no root user set")
	}

	key := s.root.APIKey()
	if s.Password == "" && key == "" {
		return "", NewError(ErrorCodeSessionCannotForgeAuthToken, "Error while creating headers: Password or APIKey must be set")
	}

	if key == "" {
		key = s.Password
	}

	return "XREST " + base64.StdEncoding.EncodeToString([]byte(s.Username+":"+key)), nil
}

func (s *Session) prepareHeaders(request *http.Request, info *FetchingInfo) *Error {

	authString, berr := s.makeAuthorizationHeaders()
	if berr != nil {
		return berr
	}

	request.Header.Set("Authorization", authString)
	request.Header.Set("X-Nuage-PageSize", "50")
	request.Header.Set("X-Nuage-Organization", s.Organization)
	request.Header.Set("Content-Type", "application/json")

	if info == nil {
		return nil
	}

	if info.Filter != "" {
		request.Header.Set("X-Nuage-Filter", info.Filter)
	}

	if info.OrderBy != "" {
		request.Header.Set("X-Nuage-OrderBy", info.OrderBy)
	}

	if info.Page != -1 {
		request.Header.Set("X-Nuage-Page", strconv.Itoa(info.Page))
	}

	if info.PageSize > 0 {
		request.Header.Set("X-Nuage-PageSize", strconv.Itoa(info.PageSize))
	}

	if len(info.GroupBy) > 0 {
		request.Header.Set("X-Nuage-GroupBy", "true")
		request.Header.Set("X-Nuage-Attributes", strings.Join(info.GroupBy, ", "))
	}

	return nil
}

func (s *Session) readHeaders(response *http.Response, info *FetchingInfo) {

	if info == nil {
		return
	}

	info.Filter = response.Header.Get("X-Nuage-Filter")
	info.FilterType = response.Header.Get("X-Nuage-FilterType")
	info.OrderBy = response.Header.Get("X-Nuage-OrderBy")
	info.Page, _ = strconv.Atoi(response.Header.Get("X-Nuage-Page"))
	info.PageSize, _ = strconv.Atoi(response.Header.Get("X-Nuage-PageSize"))
	info.TotalCount, _ = strconv.Atoi(response.Header.Get("X-Nuage-Count"))

	// info.GroupBy = response.Header.Get("X-Nuage-GroupBy")
}

func (s *Session) send(request *http.Request, info *FetchingInfo) (*http.Response, *Error) {

	s.prepareHeaders(request, info)

	response, err := s.client.Do(request)

	if err != nil {
		return response, NewError(ErrorCodeSessionCannotProcessRequest, err.Error())
	}

	switch response.StatusCode {

	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		s.readHeaders(response, info)
		return response, nil

	case http.StatusMultipleChoices:
		newURL := request.URL.String() + "?responseChoice=1"
		request, _ = http.NewRequest(request.Method, newURL, request.Body)
		return s.send(request, info)

	case http.StatusConflict:
		berr := NewError(response.StatusCode, "")
		if err := json.NewDecoder(response.Body).Decode(&berr); err != nil {
			return nil, NewError(ErrorCodeJSONCannotDecode, err.Error())
		}
		return nil, berr

	default:
		return nil, NewError(response.StatusCode, response.Status)
	}
}

func (s *Session) getGeneralURL(o Identifiable) string {

	return s.URL + "/" + o.Identity().Category
}

func (s *Session) getPersonalURL(o Identifiable) (string, *Error) {

	if _, ok := o.(Rootable); ok {
		return s.URL + "/" + o.Identity().Name, nil
	}

	if o.Identifier() == "" {
		return "", NewError(ErrorCodeSessionIDNotSet, "Cannot GetPersonalURL of an object with no ID set")
	}

	return s.getGeneralURL(o) + "/" + o.Identifier(), nil
}

func (s *Session) getURLForChildrenIdentity(o Identifiable, childrenIdentity Identity) (string, *Error) {

	if _, ok := o.(Rootable); ok {
		return s.URL + "/" + childrenIdentity.Category, nil
	}

	url, err := s.getPersonalURL(o)
	if err != nil {
		return "", err
	}

	return url + "/" + childrenIdentity.Category, nil
}

// Root returns the Root API object.
func (s *Session) Root() Rootable {

	return s.root
}

// Start starts the session.
// At that point the authentication will be done.
func (s *Session) Start() *Error {

	_currentSession = s

	err := s.FetchEntity(s.root)

	if err != nil {
		return err
	}

	return nil
}

// Reset resets the session.
func (s *Session) Reset() {

	s.root.SetAPIKey("")

	_currentSession = nil
}

// FetchEntity fetchs the given Identifiable from the server.
func (s *Session) FetchEntity(object Identifiable) *Error {

	url, berr := s.getPersonalURL(object)
	if berr != nil {
		return berr
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NewError(http.StatusBadRequest, err.Error())
	}

	response, berr := s.send(request, nil)
	if berr != nil {
		return berr
	}

	defer response.Body.Close()

	arr := IdentifiablesList{object} // trick for weird api..
	if err := json.NewDecoder(response.Body).Decode(&arr); err != nil {
		return NewError(ErrorCodeJSONCannotDecode, err.Error())
	}

	return nil
}

// SaveEntity saves the given Identifiable into the server.
func (s *Session) SaveEntity(object Identifiable) *Error {

	url, berr := s.getPersonalURL(object)
	if berr != nil {
		return berr
	}

	buffer := &bytes.Buffer{}
	if err := json.NewEncoder(buffer).Encode(object); err != nil {
		return NewError(ErrorCodeJSONCannotEncode, err.Error())
	}

	request, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		return NewError(http.StatusBadRequest, err.Error())
	}

	response, berr := s.send(request, nil)
	if berr != nil {
		return berr
	}

	defer response.Body.Close()

	dest := IdentifiablesList{object}
	if err := json.NewDecoder(response.Body).Decode(&dest); err != nil {
		return NewError(ErrorCodeJSONCannotDecode, err.Error())
	}

	return nil
}

// DeleteEntity deletes the given Identifiable from the server.
func (s *Session) DeleteEntity(object Identifiable) *Error {

	url, berr := s.getPersonalURL(object)
	if berr != nil {
		return berr
	}

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return NewError(http.StatusBadRequest, err.Error())
	}

	_, berr = s.send(request, nil)
	if berr != nil {
		return berr
	}

	return nil
}

// FetchChildren fetches the children with of given parent identified by the given Identity.
func (s *Session) FetchChildren(parent Identifiable, identity Identity, dest interface{}, info *FetchingInfo) *Error {

	url, berr := s.getURLForChildrenIdentity(parent, identity)
	if berr != nil {
		return berr
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NewError(http.StatusBadRequest, err.Error())
	}

	response, berr := s.send(request, info)
	if berr != nil {
		return berr
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent || response.ContentLength == 0 {
		return nil
	}

	if err := json.NewDecoder(response.Body).Decode(&dest); err != nil {
		return NewError(ErrorCodeJSONCannotDecode, err.Error())
	}

	return nil
}

// CreateChild creates a new child Identifiable under the given parent Identifiable in the server.
func (s *Session) CreateChild(parent Identifiable, child Identifiable) *Error {

	url, berr := s.getURLForChildrenIdentity(parent, child.Identity())
	if berr != nil {
		return berr
	}

	buffer := &bytes.Buffer{}
	if err := json.NewEncoder(buffer).Encode(child); err != nil {
		return NewError(ErrorCodeJSONCannotEncode, err.Error())
	}

	request, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return NewError(http.StatusBadRequest, err.Error())
	}

	response, berr := s.send(request, nil)
	if berr != nil {
		return berr
	}

	defer response.Body.Close()

	dest := IdentifiablesList{child}
	if err := json.NewDecoder(response.Body).Decode(&dest); err != nil {
		return NewError(ErrorCodeJSONCannotDecode, err.Error())
	}

	return nil
}

// AssignChildren assigns the list of given child Identifiables to the given Identifiable parent in the server.
func (s *Session) AssignChildren(parent Identifiable, children []Identifiable, identity Identity) *Error {

	url, berr := s.getURLForChildrenIdentity(parent, identity)
	if berr != nil {
		return berr
	}

	var ids []string
	for _, c := range children {

		if i := c.Identifier(); i != "" {
			ids = append(ids, c.Identifier())
		} else {
			return NewError(ErrorCodeSessionIDNotSet, "One of the object to assign has no ID")
		}
	}

	buffer := &bytes.Buffer{}
	json.NewEncoder(buffer).Encode(ids)

	request, err := http.NewRequest("PUT", url, buffer)
	if err != nil {
		return NewError(http.StatusBadRequest, err.Error())
	}

	_, berr = s.send(request, nil)
	if berr != nil {
		return berr
	}

	return nil
}

// NextEvent will return the next notification from the backend as it occurs and will
// send it to the correct channel.
func (s *Session) NextEvent(channel NotificationsChannel, lastEventID string) *Error {

	currentURL := s.URL + "/events"
	if lastEventID != "" {
		currentURL += "?uuid=" + lastEventID
	}

	request, err := http.NewRequest("GET", currentURL, nil)
	if err != nil {
		return NewError(http.StatusBadRequest, err.Error())
	}

	response, berr := s.send(request, nil)
	if berr != nil {
		return berr
	}

	notification := NewNotification()
	if err := json.NewDecoder(response.Body).Decode(notification); err != nil {
		return NewError(ErrorCodeJSONCannotDecode, err.Error())
	}

	if len(notification.Events) > 0 {
		channel <- notification
	}

	return nil
}
