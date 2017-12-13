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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFetchingInfo_NewFetchingInfo(t *testing.T) {

	Convey("Given I create a FetchingInfo", t, func() {
		f := NewFetchingInfo()

		Convey("Then Page should be -1", func() {
			So(f.Page, ShouldEqual, -1)
		})

		Convey("Then PageSize should -1", func() {
			So(f.PageSize, ShouldEqual, -1)
		})
	})
}

func TestFetchingInfo_String(t *testing.T) {

	Convey("Given I create a FetchingInfo", t, func() {
		f := NewFetchingInfo()

		Convey("When I set some values", func() {

			f.Filter = "filer"
			f.Page = 2
			f.PageSize = 50

			Convey("Then string representation should <FetchingInfo page: 2, pagesize: 50, totalcount: 0>", func() {
				So(f.String(), ShouldEqual, "<FetchingInfo page: 2, pagesize: 50, totalcount: 0>")
			})
		})
	})
}
