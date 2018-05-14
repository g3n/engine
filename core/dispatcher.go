// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

// Dispatcher implements an event dispatcher
type Dispatcher struct {
	evmap  map[string][]subscription // maps event name to subcriptions list
	cancel bool                      // flag informing cancelled dispatch
}

// IDispatcher is the interface for dispatchers
type IDispatcher interface {
	Subscribe(evname string, cb Callback)
	SubscribeID(evname string, id interface{}, cb Callback)
	UnsubscribeID(evname string, id interface{}) int
	Dispatch(evname string, ev interface{}) bool
	ClearSubscriptions()
	CancelDispatch()
}

// Callback is the type for the Dispatcher callbacks functions
type Callback func(string, interface{})

type subscription struct {
	id interface{}
	cb func(string, interface{})
}

// NewDispatcher creates and returns a new Event Dispatcher
func NewDispatcher() *Dispatcher {

	d := new(Dispatcher)
	d.Initialize()
	return d
}

// Initialize initializes this event dispatcher.
// It is normally used by other types which embed an event dispatcher
func (d *Dispatcher) Initialize() {

	d.evmap = make(map[string][]subscription)
}

// Subscribe subscribes to receive events with the given name.
// If it is necessary to unsubscribe the event, the function SubscribeID should be used.
func (d *Dispatcher) Subscribe(evname string, cb Callback) {

	d.SubscribeID(evname, nil, cb)
}

// SubscribeID subscribes to receive events with the given name.
// The function accepts a unique id to be use to unsubscribe this event
func (d *Dispatcher) SubscribeID(evname string, id interface{}, cb Callback) {

	d.evmap[evname] = append(d.evmap[evname], subscription{id, cb})
}

// UnsubscribeID unsubscribes from the specified event and subscription id
// Returns the number of subscriptions found.
func (d *Dispatcher) UnsubscribeID(evname string, id interface{}) int {

	// Get list of subscribers for this event
	// If not found, nothing to do
	subs, ok := d.evmap[evname]
	if !ok {
		return 0
	}

	// Remove all subscribers with the specified id for this event
	found := 0
	pos := 0
	for pos < len(subs) {
		if subs[pos].id == id {
			copy(subs[pos:], subs[pos+1:])
			subs[len(subs)-1] = subscription{}
			subs = subs[:len(subs)-1]
			found++
		} else {
			pos++
		}
	}
	d.evmap[evname] = subs
	return found
}

// UnsubscribeAllID unsubscribes from all events with the specified subscription id.
// Returns the number of subscriptions found.
func (d *Dispatcher) UnsubscribeAllID(id interface{}) int {

	total := 0
	for evname := range d.evmap {
		found := d.UnsubscribeID(evname, id)
		total += found
	}
	return total
}

// Dispatch dispatch the specified event and data to all registered subscribers.
// The function returns true if the propagation was cancelled by a subscriber.
func (d *Dispatcher) Dispatch(evname string, ev interface{}) bool {

	// Get list of subscribers for this event
	subs := d.evmap[evname]
	if subs == nil {
		return false
	}

	// Dispatch to all subscribers
	d.cancel = false
	for i := 0; i < len(subs); i++ {
		subs[i].cb(evname, ev)
		if d.cancel {
			break
		}
	}
	return d.cancel
}

// ClearSubscriptions clear all subscriptions from this dispatcher
func (d *Dispatcher) ClearSubscriptions() {

	d.evmap = make(map[string][]subscription)
}

// CancelDispatch cancels the propagation of the current event.
// No more subscribers will be called for this event dispatch.
func (d *Dispatcher) CancelDispatch() {

	d.cancel = true
}
