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
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNotification_NewNotification(t *testing.T) {

	Convey("Given I create a new Notification", t, func() {
		n := NewNotification()

		Convey("Then Events should not be nil", func() {
			So(n.Events, ShouldNotBeNil)
		})

		Convey("Then UUID should not be nil", func() {
			So(n.UUID, ShouldNotBeNil)
		})
	})
}

func TestNotification_FromJSON(t *testing.T) {

	Convey("Given I create a new notification", t, func() {
		n := NewNotification()

		Convey("When I unmarshal son json data", func() {
			d := `{"uuid": "007", "events": [{"entityType": "cat", "type": "UPDATE", "updateMechanism": "useless", "entities":[{"name": "hello"}]}]}`
			json.NewDecoder(bytes.NewBuffer([]byte(d))).Decode(n)

			Convey("Then UUI should be '007'", func() {
				So(n.UUID, ShouldEqual, "007")
			})

			Convey("Then lenght of Events should be 1", func() {
				So(len(n.Events), ShouldEqual, 1)
			})

			Convey("When I retrieve the Events", func() {
				e := n.Events[0]

				Convey("Then EntityType should be cat", func() {
					So(e.EntityType, ShouldEqual, "cat")
				})

				Convey("Then Type should UPDATE", func() {
					So(e.Type, ShouldEqual, "UPDATE")
				})

				Convey("Then UpdateMechanism should useless", func() {
					So(e.UpdateMechanism, ShouldEqual, "useless")
				})

				Convey("Then the lenght of DataMap should be 1", func() {
					So(len(e.DataMap), ShouldEqual, 1)
				})

				Convey("Then the value of item 0 of DataMap should hello", func() {
					So(e.DataMap[0]["name"], ShouldEqual, "hello")
				})
			})
		})
	})
}
