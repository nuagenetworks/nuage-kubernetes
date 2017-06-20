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

type VsdErrorList struct {
	VsdErrors    []VsdError `json:"errors"`
	VsdErrorCode int        `json:"internalErrorCode"`
}

type VsdError struct {
	Property     string  `json:"property"`
	Descriptions []Error `json:"descriptions"` // XXX -- note
}

// Errors at this level can be of two types: 1) Connection logic / setup errors (e.g. invalid credentials) 2) VSD error response.
// We use Bambou Error for both, even if they are conceptually different -- hence the note above

type Error struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewBambouError(title, description string) *Error {
	return &Error{
		Title:       title,
		Description: description,
	}
}

func NewError(code int, description string) *Error {
	return &Error{
		Title:       fmt.Sprintf("Error code: %d", code),
		Description: description,
	}
}

// Error returns the string representation of a Bambou Error (making it an "error")
// Valid JSON formatted
func (be *Error) Error() string {
	// return fmt.Sprintf("%+v", *be)
	return fmt.Sprintf("{\"title\": \"%s\", \"description\": \"%s\"}", be.Title, be.Description)
}
