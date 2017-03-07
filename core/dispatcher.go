// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import ()

type Dispatcher struct {
	evmap  map[string][]subscription // maps event name to subcriptions list
	cancel bool                      // flag informing cancelled dispatch
}

type IDispatcher interface {
	Subscribe(evname string, cb Callback)
	SubscribeID(evname string, id interface{}, cb Callback)
	UnsubscribeID(evname string, id interface{}) int
	Dispatch(evname string, ev interface{}) bool
	ClearSubscriptions()
	CancelDispatch()
}

type Callback func(string, interface{})

type subscription struct {
	id interface{}
	cb func(string, interface{})
}

// NewEventDispatcher creates and returns a pointer to an Event Dispatcher
func NewDispatcher() *Dispatcher {

	ed := new(Dispatcher)
	ed.Initialize()
	return ed
}

// Initialize initializes this event dispatcher.
// It is normally used by other types which embed an event dispatcher
func (ed *Dispatcher) Initialize() {

	ed.evmap = make(map[string][]subscription)
}

// Subscribe subscribes to receive events with the given name.
// If it is necessary to unsubscribe the event, the function SubscribeID
// should be used.
func (ed *Dispatcher) Subscribe(evname string, cb Callback) {

	ed.SubscribeID(evname, nil, cb)
}

// Subscribe subscribes to receive events with the given name.
// The function accepts a unique id to be use to unsubscribe this event
func (ed *Dispatcher) SubscribeID(evname string, id interface{}, cb Callback) {

	//log.Debug("Dispatcher(%p).SubscribeID:%s (%v)", ed, evname, id)
	ed.evmap[evname] = append(ed.evmap[evname], subscription{id, cb})
}

// Unsubscribe unsubscribes from the specified event and subscription id
// Returns the number of subscriptions found.
func (ed *Dispatcher) UnsubscribeID(evname string, id interface{}) int {

	// Get list of subscribers for this event
	// If not found, nothing to do
	subs, ok := ed.evmap[evname]
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
	//log.Debug("Dispatcher(%p).UnsubscribeID:%s (%p): %v",ed, evname, id, found)
	ed.evmap[evname] = subs
	return found
}

// Dispatch dispatch the specified event and data to all registered subscribers.
// The function returns true if the propagation was cancelled by a subscriber.
func (ed *Dispatcher) Dispatch(evname string, ev interface{}) bool {

	// Get list of subscribers for this event
	subs := ed.evmap[evname]
	if subs == nil {
		return false
	}

	// Dispatch to all subscribers
	//log.Debug("Dispatcher(%p).Dispatch:%s", ed, evname)
	ed.cancel = false
	for i := 0; i < len(subs); i++ {
		subs[i].cb(evname, ev)
		if ed.cancel {
			break
		}
	}
	return ed.cancel
}

// ClearSubscriptions clear all subscriptions from this dispatcher
func (ed *Dispatcher) ClearSubscriptions() {

	ed.evmap = make(map[string][]subscription)
	//log.Debug("Dispatcher(%p).ClearSubscriptions: %d", ed, len(ed.evmap))
}

// CancelDispatch cancels the propagation of the current event.
// No more subscribers will be called for this event dispatch.
func (ed *Dispatcher) CancelDispatch() {

	ed.cancel = true
}

//// LogSubscriptions is used for debugging to log the current
//// subscriptions of this dispatcher
//func (ed *Dispatcher) LogSubscriptions() {
//
//    for evname, subs := range ed.evmap {
//        log.Debug("event:%s", evname)
//        for _, sub := range subs {
//            log.Debug("   subscription:%v", sub)
//        }
//    }
//    log.Debug("")
//}
