// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wasm

package window

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	_ "image/png"
	"syscall/js"
)

// Keycodes
const (
	KeyUnknown = Key(iota)
	KeySpace
	KeyApostrophe
	KeyComma
	KeyMinus
	KeyPeriod
	KeySlash
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeySemicolon
	KeyEqual
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyLeftBracket
	KeyBackslash
	KeyRightBracket
	KeyGraveAccent
	KeyWorld1
	KeyWorld2
	KeyEscape
	KeyEnter
	KeyTab
	KeyBackspace
	KeyInsert
	KeyDelete
	KeyRight
	KeyLeft
	KeyDown
	KeyUp
	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd
	KeyCapsLock
	KeyScrollLock
	KeyNumLock
	KeyPrintScreen
	KeyPause
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyF13
	KeyF14
	KeyF15
	KeyF16
	KeyF17
	KeyF18
	KeyF19
	KeyF20
	KeyF21
	KeyF22
	KeyF23
	KeyF24
	KeyF25
	KeyKP0
	KeyKP1
	KeyKP2
	KeyKP3
	KeyKP4
	KeyKP5
	KeyKP6
	KeyKP7
	KeyKP8
	KeyKP9
	KeyKPDecimal
	KeyKPDivide
	KeyKPMultiply
	KeyKPSubtract
	KeyKPAdd
	KeyKPEnter
	KeyKPEqual
	KeyLeftShift
	KeyLeftControl
	KeyLeftAlt
	KeyLeftSuper // Meta in Javascript
	KeyRightShift
	KeyRightControl
	KeyRightAlt
	KeyRightSuper
	KeyMenu
	KeyLast
)

var keyMap = map[string]Key{
	//"KeyUnknown":  KeyUnknown, TODO emit when key is not in map
	"Space":  KeySpace,
	"Quote":  KeyApostrophe,
	"Comma":  KeyComma,
	"Minus":  KeyMinus,
	"Period": KeyPeriod,
	"Slash":  KeySlash,

	"Digit0": Key0,
	"Digit1": Key1,
	"Digit2": Key2,
	"Digit3": Key3,
	"Digit4": Key4,
	"Digit5": Key5,
	"Digit6": Key6,
	"Digit7": Key7,
	"Digit8": Key8,
	"Digit9": Key9,

	"Semicolon": KeySemicolon,
	"Equal":     KeyEqual,

	"KeyA": KeyA,
	"KeyB": KeyB,
	"KeyC": KeyC,
	"KeyD": KeyD,
	"KeyE": KeyE,
	"KeyF": KeyF,
	"KeyG": KeyG,
	"KeyH": KeyH,
	"KeyI": KeyI,
	"KeyJ": KeyJ,
	"KeyK": KeyK,
	"KeyL": KeyL,
	"KeyM": KeyM,
	"KeyN": KeyN,
	"KeyO": KeyO,
	"KeyP": KeyP,
	"KeyQ": KeyQ,
	"KeyR": KeyR,
	"KeyS": KeyS,
	"KeyT": KeyT,
	"KeyU": KeyU,
	"KeyV": KeyV,
	"KeyW": KeyW,
	"KeyX": KeyX,
	"KeyY": KeyY,
	"KeyZ": KeyZ,

	"BracketLeft":  KeyLeftBracket,
	"Backslash":    KeyBackslash,
	"BracketRight": KeyRightBracket,
	"Backquote":    KeyGraveAccent,
	//"KeyWorld1": 	KeyWorld1,
	//"KeyWorld2": 	KeyWorld2,

	"Escape":      KeyEscape,
	"Enter":       KeyEnter,
	"Tab":         KeyTab,
	"Backspace":   KeyBackspace,
	"Insert":      KeyInsert,
	"Delete":      KeyDelete,
	"ArrowRight":  KeyRight,
	"ArrowLeft":   KeyLeft,
	"ArrowDown":   KeyDown,
	"ArrowUp":     KeyUp,
	"PageUp":      KeyPageUp,
	"PageDown":    KeyPageDown,
	"Home":        KeyHome,
	"End":         KeyEnd,
	"CapsLock":    KeyCapsLock,
	"ScrollLock":  KeyScrollLock,
	"NumLock":     KeyNumLock,
	"PrintScreen": KeyPrintScreen,
	"Pause":       KeyPause,

	"F1":  KeyF1,
	"F2":  KeyF2,
	"F3":  KeyF3,
	"F4":  KeyF4,
	"F5":  KeyF5,
	"F6":  KeyF6,
	"F7":  KeyF7,
	"F8":  KeyF8,
	"F9":  KeyF9,
	"F10": KeyF10,
	"F11": KeyF11,
	"F12": KeyF12,
	"F13": KeyF13,
	"F14": KeyF14,
	"F15": KeyF15,
	"F16": KeyF16,
	"F17": KeyF17,
	"F18": KeyF18,
	"F19": KeyF19,
	"F20": KeyF20,
	"F21": KeyF21,
	"F22": KeyF22,
	"F23": KeyF23,
	"F24": KeyF24,
	"F25": KeyF25,

	"Numpad0": KeyKP0,
	"Numpad1": KeyKP1,
	"Numpad2": KeyKP2,
	"Numpad3": KeyKP3,
	"Numpad4": KeyKP4,
	"Numpad5": KeyKP5,
	"Numpad6": KeyKP6,
	"Numpad7": KeyKP7,
	"Numpad8": KeyKP8,
	"Numpad9": KeyKP9,

	"NumpadDecimal":  KeyKPDecimal,
	"NumpadDivide":   KeyKPDivide,
	"NumpadMultiply": KeyKPMultiply,
	"NumpadSubtract": KeyKPSubtract,
	"NumpadAdd":      KeyKPAdd,
	"NumpadEnter":    KeyKPEnter,
	"NumpadEqual":    KeyKPEqual,

	"ShiftLeft":    KeyLeftShift,
	"ControlLeft":  KeyLeftControl,
	"AltLeft":      KeyLeftAlt,
	"MetaLeft":     KeyLeftSuper,
	"ShitRight":    KeyRightShift,
	"ControlRight": KeyRightControl,
	"AltRight":     KeyRightAlt,
	"MetaRight":    KeyRightSuper,
	"Menu":         KeyMenu,
}

// Modifier keys
const (
	ModShift = ModifierKey(1 << iota) // Bitmask
	ModControl
	ModAlt
	ModSuper // Meta in Javascript
)

// Mouse buttons
const (
	//MouseButton1      = MouseButton(0)
	//MouseButton2      = MouseButton(0)
	//MouseButton3      = MouseButton(0)
	//MouseButton4      = MouseButton(0)
	//MouseButton5      = MouseButton(0)
	//MouseButton6      = MouseButton(0)
	//MouseButton7      = MouseButton(0)
	//MouseButton8      = MouseButton(0)
	//MouseButtonLast   = MouseButton(0)
	MouseButtonLeft   = MouseButton(0)
	MouseButtonRight  = MouseButton(2)
	MouseButtonMiddle = MouseButton(1)
)

// Input modes
const (
	CursorInputMode             = InputMode(iota) // See Cursor mode values
	StickyKeysInputMode                           // Value can be either 1 or 0
	StickyMouseButtonsInputMode                   // Value can be either 1 or 0
)

// Cursor mode values
const (
	CursorNormal = CursorMode(iota)
	CursorHidden
	CursorDisabled
)

// WebGlCanvas is a browser-based WebGL canvas.
type WebGlCanvas struct {
	core.Dispatcher          // Embedded event dispatcher
	canvas          js.Value // Associated WebGL canvas
	gls             *gls.GLS // Associated WebGL state

	// Events
	keyEv    KeyEvent
	charEv   CharEvent
	mouseEv  MouseEvent
	posEv    PosEvent
	sizeEv   SizeEvent
	cursorEv CursorEvent
	scrollEv ScrollEvent

	// Callbacks
	onCtxMenu  js.Func
	keyDown    js.Func
	keyUp      js.Func
	mouseDown  js.Func
	mouseUp    js.Func
	mouseMove  js.Func
	mouseWheel js.Func
	winResize  js.Func
}

// Init initializes the WebGlCanvas singleton.
// If canvasId is provided, the pre-existing WebGlCanvas with that id is used.
// If canvasId is the empty string then it creates a new WebGL canvas.
func Init(canvasId string) error {

	// Panic if already created
	if win != nil {
		panic(fmt.Errorf("can only call window.Init() once"))
	}

	// Create wrapper window with dispatcher
	w := new(WebGlCanvas)
	w.Dispatcher.Initialize()

	// Create or get WebGlCanvas
	doc := js.Global().Get("document")
	if canvasId == "" {
		w.canvas = doc.Call("createElement", "WebGlCanvas")
	} else {
		w.canvas = doc.Call("getElementById", canvasId)
		if w.canvas == js.Null() {
			panic(fmt.Sprintf("Cannot find canvas with provided id: %s", canvasId))
		}
	}

	// Get reference to WebGL context
	webglCtx := w.canvas.Call("getContext", "webgl2")
	if webglCtx == js.Undefined() {
		return fmt.Errorf("Browser doesn't support WebGL2")
	}

	// Create WebGL state
	gl, err := gls.New(webglCtx)
	if err != nil {
		return err
	}
	w.gls = gl

	// Disable right-click context menu on the canvas
	w.onCtxMenu = js.FuncOf(func(this js.Value, args []js.Value) interface{} { return false })
	w.canvas.Set("oncontextmenu", w.onCtxMenu)

	// TODO scaling/hidpi (device pixel ratio)

	// Set up key down callback to dispatch event
	w.keyDown = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		eventCode := event.Get("code").String()
		w.keyEv.Key = Key(keyMap[eventCode])
		w.keyEv.Mods = getModifiers(event)
		w.Dispatch(OnKeyDown, &w.keyEv)
		return nil
	})
	js.Global().Call("addEventListener", "keydown", w.keyDown)

	// Set up key up callback to dispatch event
	w.keyUp = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		eventCode := event.Get("code").String()
		w.keyEv.Key = Key(keyMap[eventCode])
		w.keyEv.Mods = getModifiers(event)
		w.Dispatch(OnKeyUp, &w.keyEv)
		return nil
	})
	js.Global().Call("addEventListener", "keyup", w.keyUp)

	// Set up mouse down callback to dispatch event
	w.mouseDown = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		w.mouseEv.Button = MouseButton(event.Get("button").Int())
		w.mouseEv.Xpos = float32(event.Get("offsetX").Int()) //* float32(w.scaleX) TODO
		w.mouseEv.Ypos = float32(event.Get("offsetY").Int()) //* float32(w.scaleY)
		w.mouseEv.Mods = getModifiers(event)
		w.Dispatch(OnMouseDown, &w.mouseEv)
		return nil
	})
	w.canvas.Call("addEventListener", "mousedown", w.mouseDown)

	// Set up mouse down callback to dispatch event
	w.mouseUp = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		w.mouseEv.Button = MouseButton(event.Get("button").Int())
		w.mouseEv.Xpos = float32(event.Get("offsetX").Float()) //* float32(w.scaleX) TODO
		w.mouseEv.Ypos = float32(event.Get("offsetY").Float()) //* float32(w.scaleY)
		w.mouseEv.Mods = getModifiers(event)
		w.Dispatch(OnMouseUp, &w.mouseEv)
		return nil
	})
	w.canvas.Call("addEventListener", "mouseup", w.mouseUp)

	// Set up mouse move callback to dispatch event
	w.mouseMove = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		w.cursorEv.Xpos = float32(event.Get("offsetX").Float()) //* float32(w.scaleX) TODO
		w.cursorEv.Ypos = float32(event.Get("offsetY").Float()) //* float32(w.scaleY)
		w.cursorEv.Mods = getModifiers(event)
		w.Dispatch(OnCursor, &w.cursorEv)
		return nil
	})
	w.canvas.Call("addEventListener", "mousemove", w.mouseMove)

	// Set up mouse wheel callback to dispatch event
	w.mouseWheel = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		event.Call("preventDefault")
		w.scrollEv.Xoffset = -float32(event.Get("deltaX").Float()) / 100.0
		w.scrollEv.Yoffset = -float32(event.Get("deltaY").Float()) / 100.0
		w.scrollEv.Mods = getModifiers(event)
		w.Dispatch(OnScroll, &w.scrollEv)
		return nil
	})
	w.canvas.Call("addEventListener", "wheel", w.mouseWheel)

	// Set up window resize callback to dispatch event
	w.winResize = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.sizeEv.Width = w.canvas.Get("width").Int()
		w.sizeEv.Height = w.canvas.Get("height").Int()
		// TODO device pixel ratio
		//fbw, fbh := x.GetFramebufferSize()
		//w.scaleX = float64(fbw) / float64(width)
		//w.scaleY = float64(fbh) / float64(height)
		w.Dispatch(OnWindowSize, &w.sizeEv)
		return nil
	})
	js.Global().Get("window").Call("addEventListener", "resize", w.winResize)

	//// Set up char callback to dispatch event TODO
	//w.SetCharModsCallback(func(x *glfw.Window, char rune, mods glfw.ModifierKey) {	//
	//	w.charEv.Char = char
	//	w.charEv.Mods = ModifierKey(mods)
	//	w.Dispatch(OnChar, &w.charEv)
	//})

	win = w // Set singleton
	return nil
}

// getModifiers extracts a ModifierKey bitmask from a Javascript event object.
func getModifiers(event js.Value) ModifierKey {

	shiftKey := event.Get("shiftKey").Bool()
	ctrlKey := event.Get("ctrlKey").Bool()
	altKey := event.Get("altKey").Bool()
	metaKey := event.Get("metaKey").Bool()
	var mods ModifierKey
	if shiftKey {
		mods = mods | ModShift
	}
	if ctrlKey {
		mods = mods | ModControl
	}
	if altKey {
		mods = mods | ModAlt
	}
	if metaKey {
		mods = mods | ModSuper
	}
	return mods
}

// Canvas returns the associated WebGL WebGlCanvas.
func (w *WebGlCanvas) Canvas() js.Value {

	return w.canvas
}

// Gls returns the associated OpenGL state
func (w *WebGlCanvas) Gls() *gls.GLS {

	return w.gls
}

// FullScreen returns whether this canvas is fullscreen
func (w *WebGlCanvas) FullScreen() bool {

	// TODO
	return false
}

// SetFullScreen sets this window full screen state for the primary monitor
func (w *WebGlCanvas) SetFullScreen(full bool) {

	// TODO
	// Make it so that the first user interaction (e.g. click) should set the canvas as fullscreen.
}

// Destroy destroys the WebGL canvas and removes all event listeners.
func (w *WebGlCanvas) Destroy() {

	// Remove event listeners
	w.canvas.Set("oncontextmenu", js.Null())
	js.Global().Call("removeEventListener", "keydown", w.keyDown)
	js.Global().Call("removeEventListener", "keyup", w.keyUp)
	w.canvas.Call("removeEventListener", "mousedown", w.mouseDown)
	w.canvas.Call("removeEventListener", "mouseup", w.mouseUp)
	w.canvas.Call("removeEventListener", "mousemove", w.mouseMove)
	w.canvas.Call("removeEventListener", "wheel", w.mouseWheel)
	js.Global().Get("window").Call("removeEventListener", "resize", w.winResize)

	// Release callbacks
	w.onCtxMenu.Release()
	w.keyDown.Release()
	w.keyUp.Release()
	w.mouseDown.Release()
	w.mouseUp.Release()
	w.mouseMove.Release()
	w.mouseWheel.Release()
	w.winResize.Release()
}

// GetFramebufferSize returns the framebuffer size.
func (w *WebGlCanvas) GetFramebufferSize() (width int, height int) {

	// TODO device pixel ratio
	return w.canvas.Get("width").Int(), w.canvas.Get("height").Int()
}

// GetSize returns this window's size in screen coordinates.
func (w *WebGlCanvas) GetSize() (width int, height int) {

	return w.canvas.Get("width").Int(), w.canvas.Get("height").Int()
}

// SetSize sets the size, in screen coordinates, of the canvas.
func (w *WebGlCanvas) SetSize(width int, height int) {

	w.canvas.Set("width", width)
	w.canvas.Set("height", height)
}

// Scale returns this window's DPI scale factor (FramebufferSize / Size)
func (w *WebGlCanvas) GetScale() (x float64, y float64) {

	// TODO device pixel ratio
	return 1, 1
}

// CreateCursor creates a new custom cursor and returns an int handle.
func (w *WebGlCanvas) CreateCursor(imgFile string, xhot, yhot int) (Cursor, error) {

	// TODO
	return 0, nil
}

// SetCursor sets the window's cursor to a standard one
func (w *WebGlCanvas) SetCursor(cursor Cursor) {

	// TODO
}

// DisposeAllCursors deletes all existing custom cursors.
func (w *WebGlCanvas) DisposeAllCustomCursors() {

	// TODO
}

// SetInputMode changes specified input to specified state
//func (w *WebGlCanvas) SetInputMode(mode InputMode, state int) {
//
//	// TODO
//	// Hide cursor etc
//}
