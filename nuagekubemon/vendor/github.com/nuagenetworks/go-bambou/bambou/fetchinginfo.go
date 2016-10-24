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

import "fmt"

// FetchingInfo is a structure that contains differents values about
// the fetching of children, either for the user to set, of for the server
// to fill.
// This structure will be used to pass and get back information
// during the fetching of some children.
type FetchingInfo struct {
	Filter     string
	FilterType string
	GroupBy    []string
	OrderBy    string
	Page       int
	PageSize   int
	TotalCount int
}

// NewFetchingInfo returns a new *FetchingInfo
func NewFetchingInfo() *FetchingInfo {

	return &FetchingInfo{
		Page:     -1,
		PageSize: -1,
	}
}

// String returns the string representation of the FetchingInfo.
func (f *FetchingInfo) String() string {

	return fmt.Sprintf("<FetchingInfo page: %d, pagesize: %d, totalcount: %d>", f.Page, f.PageSize, f.TotalCount)
}
