package application

import (
	"github.com/g3n/engine/window"
)

// KeyState keeps track of the state of all keys.
type KeyState struct {
	win    window.IWindow
	states map[window.Key]bool
}

// NewKeyState returns a new KeyState object.
func NewKeyState(win window.IWindow) *KeyState {

	ks := new(KeyState)
	ks.win = win
	ks.states = map[window.Key]bool{
		window.KeyUnknown      : false,
		window.KeySpace        : false,
		window.KeyApostrophe   : false,
		window.KeyComma        : false,
		window.KeyMinus        : false,
		window.KeyPeriod       : false,
		window.KeySlash        : false,
		window.Key0            : false,
		window.Key1            : false,
		window.Key2            : false,
		window.Key3            : false,
		window.Key4            : false,
		window.Key5            : false,
		window.Key6            : false,
		window.Key7            : false,
		window.Key8            : false,
		window.Key9            : false,
		window.KeySemicolon    : false,
		window.KeyEqual        : false,
		window.KeyA            : false,
		window.KeyB            : false,
		window.KeyC            : false,
		window.KeyD            : false,
		window.KeyE            : false,
		window.KeyF            : false,
		window.KeyG            : false,
		window.KeyH            : false,
		window.KeyI            : false,
		window.KeyJ            : false,
		window.KeyK            : false,
		window.KeyL            : false,
		window.KeyM            : false,
		window.KeyN            : false,
		window.KeyO            : false,
		window.KeyP            : false,
		window.KeyQ            : false,
		window.KeyR            : false,
		window.KeyS            : false,
		window.KeyT            : false,
		window.KeyU            : false,
		window.KeyV            : false,
		window.KeyW            : false,
		window.KeyX            : false,
		window.KeyY            : false,
		window.KeyZ            : false,
		window.KeyLeftBracket  : false,
		window.KeyBackslash    : false,
		window.KeyRightBracket : false,
		window.KeyGraveAccent  : false,
		window.KeyWorld1       : false,
		window.KeyWorld2       : false,
		window.KeyEscape       : false,
		window.KeyEnter        : false,
		window.KeyTab          : false,
		window.KeyBackspace    : false,
		window.KeyInsert       : false,
		window.KeyDelete       : false,
		window.KeyRight        : false,
		window.KeyLeft         : false,
		window.KeyDown         : false,
		window.KeyUp           : false,
		window.KeyPageUp       : false,
		window.KeyPageDown     : false,
		window.KeyHome         : false,
		window.KeyEnd          : false,
		window.KeyCapsLock     : false,
		window.KeyScrollLock   : false,
		window.KeyNumLock      : false,
		window.KeyPrintScreen  : false,
		window.KeyPause        : false,
		window.KeyF1           : false,
		window.KeyF2           : false,
		window.KeyF3           : false,
		window.KeyF4           : false,
		window.KeyF5           : false,
		window.KeyF6           : false,
		window.KeyF7           : false,
		window.KeyF8           : false,
		window.KeyF9           : false,
		window.KeyF10          : false,
		window.KeyF11          : false,
		window.KeyF12          : false,
		window.KeyF13          : false,
		window.KeyF14          : false,
		window.KeyF15          : false,
		window.KeyF16          : false,
		window.KeyF17          : false,
		window.KeyF18          : false,
		window.KeyF19          : false,
		window.KeyF20          : false,
		window.KeyF21          : false,
		window.KeyF22          : false,
		window.KeyF23          : false,
		window.KeyF24          : false,
		window.KeyF25          : false,
		window.KeyKP0          : false,
		window.KeyKP1          : false,
		window.KeyKP2          : false,
		window.KeyKP3          : false,
		window.KeyKP4          : false,
		window.KeyKP5          : false,
		window.KeyKP6          : false,
		window.KeyKP7          : false,
		window.KeyKP8          : false,
		window.KeyKP9          : false,
		window.KeyKPDecimal    : false,
		window.KeyKPDivide     : false,
		window.KeyKPMultiply   : false,
		window.KeyKPSubtract   : false,
		window.KeyKPAdd        : false,
		window.KeyKPEnter      : false,
		window.KeyKPEqual      : false,
		window.KeyLeftShift    : false,
		window.KeyLeftControl  : false,
		window.KeyLeftAlt      : false,
		window.KeyLeftSuper    : false,
		window.KeyRightShift   : false,
		window.KeyRightControl : false,
		window.KeyRightAlt     : false,
		window.KeyRightSuper   : false,
		window.KeyMenu         : false,
	}

	// Subscribe to window key events
	ks.win.SubscribeID(window.OnKeyUp, &ks, ks.onKey)
	ks.win.SubscribeID(window.OnKeyDown, &ks, ks.onKey)

	return ks
}

// Dispose unsubscribes from the window events.
func (ks *KeyState) Dispose() {

	ks.win.UnsubscribeID(window.OnKeyUp, &ks)
	ks.win.UnsubscribeID(window.OnKeyDown, &ks)
}

// Pressed returns whether the specified key is currently pressed.
func (ks *KeyState) Pressed(k window.Key) bool {

	return ks.states[k]
}

// onKey receives key events and updates the internal map of states.
func (ks *KeyState) onKey(evname string, ev interface{}) {

	kev := ev.(*window.KeyEvent)
	switch evname {
	case window.OnKeyUp:
		ks.states[kev.Keycode] = false
	case window.OnKeyDown:
		ks.states[kev.Keycode] = true
	}
}
