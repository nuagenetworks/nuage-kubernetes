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
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPushCenter_NewPushCenter(t *testing.T) {

	Convey("Given I create a new PushCenter", t, func() {

		p := NewPushCenter(nil)

		Convey("Then the Channel should not be nil", func() {
			So(p.Channel, ShouldNotBeNil)
			So(p.Channel, ShouldHaveSameTypeAs, make(NotificationsChannel))
		})

		Convey("Then the stop channel should not be nil", func() {
			So(p.stop, ShouldNotBeNil)
			So(p.stop, ShouldHaveSameTypeAs, make(chan bool))
		})

		Convey("Then the handlers list should not be nil", func() {
			So(p.handlers, ShouldNotBeNil)
			So(p.handlers, ShouldHaveSameTypeAs, make(eventHandlers))
		})
	})
}

func TestPushCenter_HandlersRegistration(t *testing.T) {

	Convey("Given I create a new PushCenter and a handler", t, func() {

		p := NewPushCenter(nil)
		h := func(*Event) {}

		Convey("When I register the handler for an identity", func() {
			p.RegisterHandlerForIdentity(h, FakeIdentity)

			Convey("Then it should be registered in the list for that identity", func() {
				So(p.HasHandlerForIdentity(FakeIdentity), ShouldBeTrue)
			})

			Convey("Then the default handler should be nil", func() {
				So(p.defaultHander, ShouldBeNil)
			})

			Convey("When I unregister the handler for that identity", func() {

				p.UnregisterHandlerForIdentity(FakeIdentity)

				Convey("Then it should not be registered in the list anymore", func() {
					So(p.HasHandlerForIdentity(FakeIdentity), ShouldBeFalse)
				})

				Convey("Then the default handler should be nil", func() {
					So(p.defaultHander, ShouldBeNil)
				})
			})
		})

		Convey("When I register handler for the all identity", func() {
			p.RegisterHandlerForIdentity(h, AllIdentity)

			Convey("Then it should not be registered in the list", func() {
				So(p.HasHandlerForIdentity(AllIdentity), ShouldBeTrue)
			})

			Convey("Then it should be set as the defaultHandler", func() {
				So(p.defaultHander, ShouldEqual, h)
			})

			Convey("When I unregister the handler for the all identity", func() {
				p.UnregisterHandlerForIdentity(AllIdentity)

				Convey("Then it should not be registered in the lista anymore", func() {
					So(p.HasHandlerForIdentity(AllIdentity), ShouldBeFalse)
				})

				Convey("Then it should not be set as the defaultHandler anymore", func() {
					So(p.defaultHander, ShouldBeNil)
				})
			})
		})
	})
}

func TestPushCenter_Start(t *testing.T) {

	Convey("Given I create a new PushCenter and resgister a handler", t, func() {

		cont := make(chan bool)

		type notifList struct {
			notifications EventsList
			lock          sync.Mutex
		}

		n := notifList{
			notifications: make(EventsList, 0),
		}
		c := 0

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if c == 0 {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"uuid": "x", "events": [{"type": "CREATE", "entityType": "fake", "updateMechanism": "DEFAULT", "entities": [{"ID": "x"}]}]}`)
			} else if c == 1 {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"uuid": "y", "events": [{"type": "CREATE", "entityType": "notfake", "updateMechanism": "DEFAULT", "entities": [{"ID": "y"}]}]}`)
			} else {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, `{"uuid": "z", "events": [{"type": "CREATE", "entityType": "fake", "updateMechanism": "DEFAULT", "entities": [{"ID": "z"}]}]}`)
				cont <- true
			}
			c++
		}))
		defer ts.Close()

		r := NewFakeRootObject()

		session := NewSession("username", "password", "organization", ts.URL, r)

		p := NewPushCenter(session)
		h1 := func(e *Event) {
			n.lock.Lock()
			n.notifications = append(n.notifications, e)
			n.lock.Unlock()
		}
		h2 := func(e *Event) {
			n.lock.Lock()
			n.notifications = append(n.notifications, e)
			n.lock.Unlock()
		}
		p.RegisterHandlerForIdentity(h1, AllIdentity)
		p.RegisterHandlerForIdentity(h2, FakeIdentity)

		Convey("When I start the push center and receive the notifications", func() {

			p.Start()
			<-cont
			n.lock.Lock()

			Convey("Then the number of notifications should be 3", func() {
				So(len(n.notifications), ShouldEqual, 3)
			})

			Convey("Then events Data should not be empty ", func() {
				So(len(n.notifications[0].Data), ShouldBeGreaterThan, 0)
				So(len(n.notifications[1].Data), ShouldBeGreaterThan, 0)
				So(len(n.notifications[2].Data), ShouldBeGreaterThan, 0)
			})

			Convey("When I decode the data map", func() {

				o1 := NewFakeObject("")
				o2 := NewFakeObject("")
				o3 := NewFakeObject("")
				json.NewDecoder(bytes.NewBuffer([]byte(n.notifications[0].Data))).Decode(o1)
				json.NewDecoder(bytes.NewBuffer([]byte(n.notifications[1].Data))).Decode(o2)
				json.NewDecoder(bytes.NewBuffer([]byte(n.notifications[2].Data))).Decode(o3)

				Convey("Then the ID of the decoded object should be correct ", func() {
					So(o1.ID, ShouldEqual, "x")
					So(o2.ID, ShouldEqual, "x")
					So(o3.ID, ShouldEqual, "y")
				})
			})

			n.lock.Unlock()
		})
	})
}

func TestPushCenter_Stop(t *testing.T) {

	Convey("Given I have a started Push Center", t, func() {

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
		}))
		defer ts.Close()

		r := NewFakeRootObject()

		session := NewSession("username", "password", "organization", ts.URL, r)

		p := NewPushCenter(session)
		err := p.Start()

		Convey("Then the error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("When I start it again", func() {

			err = p.Start()

			Convey("Then the error should not be nil", func() {
				So(err, ShouldNotBeNil)
			})
		})

		Convey("When I stop the push center", func() {

			p.Stop()

			Convey("Then the error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("When I stop it again", func() {

				err = p.Stop()

				Convey("Then the error should not be nil", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})
	})
}
