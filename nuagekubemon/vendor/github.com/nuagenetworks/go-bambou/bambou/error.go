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
	"fmt"
)

const (
	// ErrorCodeJSONCannotDecode is the code that means it is impossible to
	// decode json data.
	ErrorCodeJSONCannotDecode = 10001

	// ErrorCodeJSONCannotEncode is the code that means it is impossible to
	// encode json data.
	ErrorCodeJSONCannotEncode = 10002

	// ErrorCodeSessionAlreadyStarted is the error code that means a session
	// is already stared.
	ErrorCodeSessionAlreadyStarted = 11001

	// ErrorCodeSessionCannotForgetAuthToken is the code that means no password
	// or token has been given to the session.
	ErrorCodeSessionCannotForgeAuthToken = 11002

	// ErrorCodeSessionCannotProcessRequest is the code that means that it was
	// impossible to process a request.
	ErrorCodeSessionCannotProcessRequest = 11003

	// ErrorCodeSessionIDNotSet is the code that means the Identifiable is
	// missing a required ID.
	ErrorCodeSessionIDNotSet = 11004

	// ErrorCodeSessionUsernameNotSet is the code that means the username
	// is missing.
	ErrorCodeSessionUsernameNotSet = 11005
)

// ErrorDescriptionsList represents a list of *ErrorDescriptions.
type ErrorDescriptionsList []*ErrorDescription

// ErrorDescription represents an entry in an Error.
type ErrorDescription struct {
	Description string `json:"description"`
	Title       string `json:"title"`
}

// Error represent an connection error.
type Error struct {
	Code         int                   `json:"-"`
	Property     string                `json:"property"`
	Message      string                `json:"type"`
	Descriptions ErrorDescriptionsList `json:"descriptions"`
}

// NewError returns a new *Error.
// You can give a message, that will be used, if no additional
// information is given by the server. Otherwhise Message will be
// overwritten.
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Error returns the string representation of the Error.
func (e *Error) Error() string {

	return fmt.Sprintf("<Error: %d, message: %s>", e.Code, e.Message)
}
