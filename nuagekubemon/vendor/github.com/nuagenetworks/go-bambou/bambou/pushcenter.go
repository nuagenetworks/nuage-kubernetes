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
	"errors"
)

// NotificationsChannel is used to received notification from the session
type NotificationsChannel chan *Notification

// EventHandler is prototype of a Push Center Handler.
type EventHandler func(*Event)

// eventHandlers represents a map of EventHandler based on the identity.
type eventHandlers map[string]EventHandler

// PushCenter is a structure that allows the user to deal with notifications.
// You can register multiple handlers for several Identity. When a notification
// is sent by the server and the Identity of its content matches one of the
// registered handler, this handler will be called.
type PushCenter struct {
	isRunning bool
	Channel   NotificationsChannel

	handlers      eventHandlers
	defaultHander EventHandler
	stop          chan bool
	session       *Session
}

// NewPushCenter creates a new PushCenter.
func NewPushCenter(session *Session) *PushCenter {

	return &PushCenter{
		Channel:  make(NotificationsChannel),
		stop:     make(chan bool),
		handlers: eventHandlers{},
		session:  session,
	}
}

// RegisterHandlerForIdentity registers the given EventHandler for the given Entity Identity.
// You can pass the bambou.AllIdentity as identity to register the handler
// for all events. If you pass a handler for an Identity that is already registered
// the previous handler will be silently overwriten.
func (p *PushCenter) RegisterHandlerForIdentity(handler EventHandler, identity Identity) {

	if identity.Name == AllIdentity.Name {
		p.defaultHander = handler
		return
	}

	p.handlers[identity.Name] = handler
}

// UnregisterHandlerForIdentity unegisters the given EventHandler for the given Entity Identity.
func (p *PushCenter) UnregisterHandlerForIdentity(identity Identity) {

	if identity.Name == AllIdentity.Name {
		p.defaultHander = nil
		return
	}

	if _, exists := p.handlers[identity.Name]; exists {
		delete(p.handlers, identity.Name)
	}
}

// HasHandlerForIdentity verifies if the given identity has a registered handler.
func (p *PushCenter) HasHandlerForIdentity(identity Identity) bool {

	if identity.Name == AllIdentity.Name {
		return p.defaultHander != nil
	}
	_, exists := p.handlers[identity.Name]
	return exists
}

// Start starts the Push Center.
func (p *PushCenter) Start() error {

	if p.isRunning {
		return errors.New("the push center is already started")
	}

	p.isRunning = true

	go func() {
		lastEventID := ""
		for {
			go p.session.NextEvent(p.Channel, lastEventID)
			select {
			case notification := <-p.Channel:
				for _, event := range notification.Events {

					buffer := &bytes.Buffer{}
					if err := json.NewEncoder(buffer).Encode(event.DataMap[0]); err != nil {
						continue
					}
					event.Data = buffer.Bytes()

					lastEventID = notification.UUID
					if p.defaultHander != nil {
						p.defaultHander(event)
					}

					if handler, exists := p.handlers[event.EntityType]; exists {
						handler(event)
					}
				}
			case <-p.stop:
				return
			}
		}
	}()

	return nil
}

// Stop stops a running PushCenter.
func (p *PushCenter) Stop() error {

	if !p.isRunning {
		return errors.New("the push center is not started")
	}

	p.stop <- true
	p.isRunning = false

	return nil
}
