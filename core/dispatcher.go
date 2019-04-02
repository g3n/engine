// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

// IDispatcher is the interface for event dispatchers.
type IDispatcher interface {
	Subscribe(evname string, cb Callback)
	SubscribeID(evname string, id interface{}, cb Callback)
	UnsubscribeID(evname string, id interface{}) int
	UnsubscribeAllID(id interface{}) int
	Dispatch(evname string, ev interface{}) int
}

// Dispatcher implements an event dispatcher.
type Dispatcher struct {
	evmap map[string][]subscription // Map of event names to subscription lists
}

// Callback is the type for Dispatcher callback functions.
type Callback func(string, interface{})

// subscription links a Callback with a user-provided unique id.
type subscription struct {
	id interface{}
	cb Callback
}

// NewDispatcher creates and returns a new event dispatcher.
func NewDispatcher() *Dispatcher {

	d := new(Dispatcher)
	d.Initialize()
	return d
}

// Initialize initializes the event dispatcher.
// It is normally used by other types which embed a dispatcher.
func (d *Dispatcher) Initialize() {

	d.evmap = make(map[string][]subscription)
}

// Subscribe subscribes a callback to events with the given name.
// If it is necessary to unsubscribe later, SubscribeID should be used instead.
func (d *Dispatcher) Subscribe(evname string, cb Callback) {

	d.evmap[evname] = append(d.evmap[evname], subscription{nil, cb})
}

// SubscribeID subscribes a callback to events events with the given name.
// The user-provided unique id can be used to unsubscribe via UnsubscribeID.
func (d *Dispatcher) SubscribeID(evname string, id interface{}, cb Callback) {

	d.evmap[evname] = append(d.evmap[evname], subscription{id, cb})
}

// UnsubscribeID removes all subscribed callbacks with the specified unique id from the specified event.
// Returns the number of subscriptions removed.
func (d *Dispatcher) UnsubscribeID(evname string, id interface{}) int {

	// Get list of subscribers for this event
	subs := d.evmap[evname]
	if len(subs) == 0 {
		return 0
	}

	// Remove all subscribers of the specified event with the specified id, counting how many were removed
	rm := 0
	i := 0
	for _, s := range subs {
		if s.id == id {
			rm++
		} else {
			subs[i] = s
			i++
		}
	}
	d.evmap[evname] = subs[:i]
	return rm
}

// UnsubscribeAllID removes all subscribed callbacks with the specified unique id from all events.
// Returns the number of subscriptions removed.
func (d *Dispatcher) UnsubscribeAllID(id interface{}) int {

	// Remove all subscribers with the specified id (for all events), counting how many were removed
	total := 0
	for evname := range d.evmap {
		total += d.UnsubscribeID(evname, id)
	}
	return total
}

// Dispatch dispatches the specified event to all registered subscribers.
// The function returns the number of subscribers to which the event was dispatched.
func (d *Dispatcher) Dispatch(evname string, ev interface{}) int {

	// Get list of subscribers for this event
	subs := d.evmap[evname]
	nsubs := len(subs)
	if nsubs == 0 {
		return 0
	}

	// Dispatch event to all subscribers
	for _, s := range subs {
		s.cb(evname, ev)
	}
	return nsubs
}
