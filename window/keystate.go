// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window

import "github.com/g3n/engine/core"

// KeyState keeps track of the state of all keys.
type KeyState struct {
	win    core.IDispatcher
	states map[Key]bool
}

// NewKeyState returns a new KeyState object.
func NewKeyState(win core.IDispatcher) *KeyState {

	ks := new(KeyState)
	ks.win = win
	ks.states = map[Key]bool{
		KeyUnknown:      false,
		KeySpace:        false,
		KeyApostrophe:   false,
		KeyComma:        false,
		KeyMinus:        false,
		KeyPeriod:       false,
		KeySlash:        false,
		Key0:            false,
		Key1:            false,
		Key2:            false,
		Key3:            false,
		Key4:            false,
		Key5:            false,
		Key6:            false,
		Key7:            false,
		Key8:            false,
		Key9:            false,
		KeySemicolon:    false,
		KeyEqual:        false,
		KeyA:            false,
		KeyB:            false,
		KeyC:            false,
		KeyD:            false,
		KeyE:            false,
		KeyF:            false,
		KeyG:            false,
		KeyH:            false,
		KeyI:            false,
		KeyJ:            false,
		KeyK:            false,
		KeyL:            false,
		KeyM:            false,
		KeyN:            false,
		KeyO:            false,
		KeyP:            false,
		KeyQ:            false,
		KeyR:            false,
		KeyS:            false,
		KeyT:            false,
		KeyU:            false,
		KeyV:            false,
		KeyW:            false,
		KeyX:            false,
		KeyY:            false,
		KeyZ:            false,
		KeyLeftBracket:  false,
		KeyBackslash:    false,
		KeyRightBracket: false,
		KeyGraveAccent:  false,
		KeyWorld1:       false,
		KeyWorld2:       false,
		KeyEscape:       false,
		KeyEnter:        false,
		KeyTab:          false,
		KeyBackspace:    false,
		KeyInsert:       false,
		KeyDelete:       false,
		KeyRight:        false,
		KeyLeft:         false,
		KeyDown:         false,
		KeyUp:           false,
		KeyPageUp:       false,
		KeyPageDown:     false,
		KeyHome:         false,
		KeyEnd:          false,
		KeyCapsLock:     false,
		KeyScrollLock:   false,
		KeyNumLock:      false,
		KeyPrintScreen:  false,
		KeyPause:        false,
		KeyF1:           false,
		KeyF2:           false,
		KeyF3:           false,
		KeyF4:           false,
		KeyF5:           false,
		KeyF6:           false,
		KeyF7:           false,
		KeyF8:           false,
		KeyF9:           false,
		KeyF10:          false,
		KeyF11:          false,
		KeyF12:          false,
		KeyF13:          false,
		KeyF14:          false,
		KeyF15:          false,
		KeyF16:          false,
		KeyF17:          false,
		KeyF18:          false,
		KeyF19:          false,
		KeyF20:          false,
		KeyF21:          false,
		KeyF22:          false,
		KeyF23:          false,
		KeyF24:          false,
		KeyF25:          false,
		KeyKP0:          false,
		KeyKP1:          false,
		KeyKP2:          false,
		KeyKP3:          false,
		KeyKP4:          false,
		KeyKP5:          false,
		KeyKP6:          false,
		KeyKP7:          false,
		KeyKP8:          false,
		KeyKP9:          false,
		KeyKPDecimal:    false,
		KeyKPDivide:     false,
		KeyKPMultiply:   false,
		KeyKPSubtract:   false,
		KeyKPAdd:        false,
		KeyKPEnter:      false,
		KeyKPEqual:      false,
		KeyLeftShift:    false,
		KeyLeftControl:  false,
		KeyLeftAlt:      false,
		KeyLeftSuper:    false,
		KeyRightShift:   false,
		KeyRightControl: false,
		KeyRightAlt:     false,
		KeyRightSuper:   false,
		KeyMenu:         false,
	}

	// Subscribe to window key events
	ks.win.SubscribeID(OnKeyUp, &ks, ks.onKey)
	ks.win.SubscribeID(OnKeyDown, &ks, ks.onKey)

	return ks
}

// Dispose unsubscribes from the window events.
func (ks *KeyState) Dispose() {

	ks.win.UnsubscribeID(OnKeyUp, &ks)
	ks.win.UnsubscribeID(OnKeyDown, &ks)
}

// Pressed returns whether the specified key is currently pressed.
func (ks *KeyState) Pressed(k Key) bool {

	return ks.states[k]
}

// onKey receives key events and updates the internal map of states.
func (ks *KeyState) onKey(evname string, ev interface{}) {

	kev := ev.(*KeyEvent)
	switch evname {
	case OnKeyUp:
		ks.states[kev.Key] = false
	case OnKeyDown:
		ks.states[kev.Key] = true
	}
}
