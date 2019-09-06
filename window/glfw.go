// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !wasm

package window

import (
	"bytes"
	"fmt"
	"github.com/g3n/engine/gui/assets"
	"runtime"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/go-gl/glfw/v3.2/glfw"
	"image"
	_ "image/png"
	"os"
)

// Keycodes
const (
	KeyUnknown      = Key(glfw.KeyUnknown)
	KeySpace        = Key(glfw.KeySpace)
	KeyApostrophe   = Key(glfw.KeyApostrophe)
	KeyComma        = Key(glfw.KeyComma)
	KeyMinus        = Key(glfw.KeyMinus)
	KeyPeriod       = Key(glfw.KeyPeriod)
	KeySlash        = Key(glfw.KeySlash)
	Key0            = Key(glfw.Key0)
	Key1            = Key(glfw.Key1)
	Key2            = Key(glfw.Key2)
	Key3            = Key(glfw.Key3)
	Key4            = Key(glfw.Key4)
	Key5            = Key(glfw.Key5)
	Key6            = Key(glfw.Key6)
	Key7            = Key(glfw.Key7)
	Key8            = Key(glfw.Key8)
	Key9            = Key(glfw.Key9)
	KeySemicolon    = Key(glfw.KeySemicolon)
	KeyEqual        = Key(glfw.KeyEqual)
	KeyA            = Key(glfw.KeyA)
	KeyB            = Key(glfw.KeyB)
	KeyC            = Key(glfw.KeyC)
	KeyD            = Key(glfw.KeyD)
	KeyE            = Key(glfw.KeyE)
	KeyF            = Key(glfw.KeyF)
	KeyG            = Key(glfw.KeyG)
	KeyH            = Key(glfw.KeyH)
	KeyI            = Key(glfw.KeyI)
	KeyJ            = Key(glfw.KeyJ)
	KeyK            = Key(glfw.KeyK)
	KeyL            = Key(glfw.KeyL)
	KeyM            = Key(glfw.KeyM)
	KeyN            = Key(glfw.KeyN)
	KeyO            = Key(glfw.KeyO)
	KeyP            = Key(glfw.KeyP)
	KeyQ            = Key(glfw.KeyQ)
	KeyR            = Key(glfw.KeyR)
	KeyS            = Key(glfw.KeyS)
	KeyT            = Key(glfw.KeyT)
	KeyU            = Key(glfw.KeyU)
	KeyV            = Key(glfw.KeyV)
	KeyW            = Key(glfw.KeyW)
	KeyX            = Key(glfw.KeyX)
	KeyY            = Key(glfw.KeyY)
	KeyZ            = Key(glfw.KeyZ)
	KeyLeftBracket  = Key(glfw.KeyLeftBracket)
	KeyBackslash    = Key(glfw.KeyBackslash)
	KeyRightBracket = Key(glfw.KeyRightBracket)
	KeyGraveAccent  = Key(glfw.KeyGraveAccent)
	KeyWorld1       = Key(glfw.KeyWorld1)
	KeyWorld2       = Key(glfw.KeyWorld2)
	KeyEscape       = Key(glfw.KeyEscape)
	KeyEnter        = Key(glfw.KeyEnter)
	KeyTab          = Key(glfw.KeyTab)
	KeyBackspace    = Key(glfw.KeyBackspace)
	KeyInsert       = Key(glfw.KeyInsert)
	KeyDelete       = Key(glfw.KeyDelete)
	KeyRight        = Key(glfw.KeyRight)
	KeyLeft         = Key(glfw.KeyLeft)
	KeyDown         = Key(glfw.KeyDown)
	KeyUp           = Key(glfw.KeyUp)
	KeyPageUp       = Key(glfw.KeyPageUp)
	KeyPageDown     = Key(glfw.KeyPageDown)
	KeyHome         = Key(glfw.KeyHome)
	KeyEnd          = Key(glfw.KeyEnd)
	KeyCapsLock     = Key(glfw.KeyCapsLock)
	KeyScrollLock   = Key(glfw.KeyScrollLock)
	KeyNumLock      = Key(glfw.KeyNumLock)
	KeyPrintScreen  = Key(glfw.KeyPrintScreen)
	KeyPause        = Key(glfw.KeyPause)
	KeyF1           = Key(glfw.KeyF1)
	KeyF2           = Key(glfw.KeyF2)
	KeyF3           = Key(glfw.KeyF3)
	KeyF4           = Key(glfw.KeyF4)
	KeyF5           = Key(glfw.KeyF5)
	KeyF6           = Key(glfw.KeyF6)
	KeyF7           = Key(glfw.KeyF7)
	KeyF8           = Key(glfw.KeyF8)
	KeyF9           = Key(glfw.KeyF9)
	KeyF10          = Key(glfw.KeyF10)
	KeyF11          = Key(glfw.KeyF11)
	KeyF12          = Key(glfw.KeyF12)
	KeyF13          = Key(glfw.KeyF13)
	KeyF14          = Key(glfw.KeyF14)
	KeyF15          = Key(glfw.KeyF15)
	KeyF16          = Key(glfw.KeyF16)
	KeyF17          = Key(glfw.KeyF17)
	KeyF18          = Key(glfw.KeyF18)
	KeyF19          = Key(glfw.KeyF19)
	KeyF20          = Key(glfw.KeyF20)
	KeyF21          = Key(glfw.KeyF21)
	KeyF22          = Key(glfw.KeyF22)
	KeyF23          = Key(glfw.KeyF23)
	KeyF24          = Key(glfw.KeyF24)
	KeyF25          = Key(glfw.KeyF25)
	KeyKP0          = Key(glfw.KeyKP0)
	KeyKP1          = Key(glfw.KeyKP1)
	KeyKP2          = Key(glfw.KeyKP2)
	KeyKP3          = Key(glfw.KeyKP3)
	KeyKP4          = Key(glfw.KeyKP4)
	KeyKP5          = Key(glfw.KeyKP5)
	KeyKP6          = Key(glfw.KeyKP6)
	KeyKP7          = Key(glfw.KeyKP7)
	KeyKP8          = Key(glfw.KeyKP8)
	KeyKP9          = Key(glfw.KeyKP9)
	KeyKPDecimal    = Key(glfw.KeyKPDecimal)
	KeyKPDivide     = Key(glfw.KeyKPDivide)
	KeyKPMultiply   = Key(glfw.KeyKPMultiply)
	KeyKPSubtract   = Key(glfw.KeyKPSubtract)
	KeyKPAdd        = Key(glfw.KeyKPAdd)
	KeyKPEnter      = Key(glfw.KeyKPEnter)
	KeyKPEqual      = Key(glfw.KeyKPEqual)
	KeyLeftShift    = Key(glfw.KeyLeftShift)
	KeyLeftControl  = Key(glfw.KeyLeftControl)
	KeyLeftAlt      = Key(glfw.KeyLeftAlt)
	KeyLeftSuper    = Key(glfw.KeyLeftSuper)
	KeyRightShift   = Key(glfw.KeyRightShift)
	KeyRightControl = Key(glfw.KeyRightControl)
	KeyRightAlt     = Key(glfw.KeyRightAlt)
	KeyRightSuper   = Key(glfw.KeyRightSuper)
	KeyMenu         = Key(glfw.KeyMenu)
	KeyLast         = Key(glfw.KeyLast)
)

// Modifier keys
const (
	ModShift   = ModifierKey(glfw.ModShift)
	ModControl = ModifierKey(glfw.ModControl)
	ModAlt     = ModifierKey(glfw.ModAlt)
	ModSuper   = ModifierKey(glfw.ModSuper)
)

// Mouse buttons
const (
	MouseButton1      = MouseButton(glfw.MouseButton1)
	MouseButton2      = MouseButton(glfw.MouseButton2)
	MouseButton3      = MouseButton(glfw.MouseButton3)
	MouseButton4      = MouseButton(glfw.MouseButton4)
	MouseButton5      = MouseButton(glfw.MouseButton5)
	MouseButton6      = MouseButton(glfw.MouseButton6)
	MouseButton7      = MouseButton(glfw.MouseButton7)
	MouseButton8      = MouseButton(glfw.MouseButton8)
	MouseButtonLast   = MouseButton(glfw.MouseButtonLast)
	MouseButtonLeft   = MouseButton(glfw.MouseButtonLeft)
	MouseButtonRight  = MouseButton(glfw.MouseButtonRight)
	MouseButtonMiddle = MouseButton(glfw.MouseButtonMiddle)
)

// Input modes
const (
	CursorInputMode             = InputMode(glfw.CursorMode)             // See Cursor mode values
	StickyKeysInputMode         = InputMode(glfw.StickyKeysMode)         // Value can be either 1 or 0
	StickyMouseButtonsInputMode = InputMode(glfw.StickyMouseButtonsMode) // Value can be either 1 or 0
)

// Cursor mode values
const (
	CursorNormal   = CursorMode(glfw.CursorNormal)
	CursorHidden   = CursorMode(glfw.CursorHidden)
	CursorDisabled = CursorMode(glfw.CursorDisabled)
)

// GlfwWindow describes one glfw window
type GlfwWindow struct {
	*glfw.Window             // Embedded GLFW window
	core.Dispatcher          // Embedded event dispatcher
	gls             *gls.GLS // Associated OpenGL State
	fullscreen      bool
	lastX           int
	lastY           int
	lastWidth       int
	lastHeight      int
	scaleX          float64
	scaleY          float64

	// Events
	keyEv    KeyEvent
	charEv   CharEvent
	mouseEv  MouseEvent
	posEv    PosEvent
	sizeEv   SizeEvent
	cursorEv CursorEvent
	scrollEv ScrollEvent

	mods ModifierKey // Current modifier keys

	// Cursors
	cursors       map[Cursor]*glfw.Cursor
	lastCursorKey Cursor
}

// Init initializes the GlfwWindow singleton with the specified width, height, and title.
func Init(width, height int, title string) error {

	// Panic if already created
	if win != nil {
		panic(fmt.Errorf("can only call window.Init() once"))
	}

	// OpenGL functions must be executed in the same thread where
	// the context was created (by wmgr.CreateWindow())
	runtime.LockOSThread()

	// Create wrapper window with dispatcher
	w := new(GlfwWindow)
	w.Dispatcher.Initialize()
	var err error

	// Initialize GLFW
	err = glfw.Init()
	if err != nil {
		return err
	}

	// Set window hints
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Samples, 8)
	// Set OpenGL forward compatible context only for OSX because it is required for OSX.
	// When this is set, glLineWidth(width) only accepts width=1.0 and generates an error
	// for any other values although the spec says it should ignore unsupported widths
	// and generate an error only when width <= 0.
	if runtime.GOOS == "darwin" {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	}

	// Create window and set it as the current context.
	// The window is created always as not full screen because if it is
	// created as full screen it not possible to revert it to windowed mode.
	// At the end of this function, the window will be set to full screen if requested.
	w.Window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return err
	}
	w.MakeContextCurrent()

	// Create OpenGL state
	w.gls, err = gls.New()
	if err != nil {
		return err
	}

	// Compute and store scale
	fbw, fbh := w.GetFramebufferSize()
	w.scaleX = float64(fbw) / float64(width)
	w.scaleY = float64(fbh) / float64(height)

	// Create map for cursors
	w.cursors = make(map[Cursor]*glfw.Cursor)
	w.lastCursorKey = CursorLast

	// Preallocate GLFW standard cursors
	w.cursors[ArrowCursor] = glfw.CreateStandardCursor(glfw.ArrowCursor)
	w.cursors[IBeamCursor] = glfw.CreateStandardCursor(glfw.IBeamCursor)
	w.cursors[CrosshairCursor] = glfw.CreateStandardCursor(glfw.CrosshairCursor)
	w.cursors[HandCursor] = glfw.CreateStandardCursor(glfw.HandCursor)
	w.cursors[HResizeCursor] = glfw.CreateStandardCursor(glfw.HResizeCursor)
	w.cursors[VResizeCursor] = glfw.CreateStandardCursor(glfw.VResizeCursor)

	// Preallocate extra G3N standard cursors (diagonal resize cursors)
	cursorDiag1Png := assets.MustAsset("cursors/diag1.png") // [/]
	cursorDiag2Png := assets.MustAsset("cursors/diag2.png") // [\]
	diag1Img, _, err := image.Decode(bytes.NewReader(cursorDiag1Png))
	diag2Img, _, err := image.Decode(bytes.NewReader(cursorDiag2Png))
	if err != nil {
		return err
	}
	w.cursors[DiagResize1Cursor] = glfw.CreateCursor(diag1Img, 8, 8) // [/]
	w.cursors[DiagResize2Cursor] = glfw.CreateCursor(diag2Img, 8, 8) // [\]

	// Set up key callback to dispatch event
	w.SetKeyCallback(func(x *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		w.keyEv.Key = Key(key)
		w.keyEv.Mods = ModifierKey(mods)
		w.mods = w.keyEv.Mods
		if action == glfw.Press {
			w.Dispatch(OnKeyDown, &w.keyEv)
		} else if action == glfw.Release {
			w.Dispatch(OnKeyUp, &w.keyEv)
		} else if action == glfw.Repeat {
			w.Dispatch(OnKeyRepeat, &w.keyEv)
		}
	})

	// Set up char callback to dispatch event
	w.SetCharModsCallback(func(x *glfw.Window, char rune, mods glfw.ModifierKey) {
		w.charEv.Char = char
		w.charEv.Mods = ModifierKey(mods)
		w.Dispatch(OnChar, &w.charEv)
	})

	// Set up mouse button callback to dispatch event
	w.SetMouseButtonCallback(func(x *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		xpos, ypos := x.GetCursorPos()
		w.mouseEv.Button = MouseButton(button)
		w.mouseEv.Mods = ModifierKey(mods)
		w.mouseEv.Xpos = float32(xpos * w.scaleX)
		w.mouseEv.Ypos = float32(ypos * w.scaleY)
		if action == glfw.Press {
			w.Dispatch(OnMouseDown, &w.mouseEv)
		} else if action == glfw.Release {
			w.Dispatch(OnMouseUp, &w.mouseEv)
		}
	})

	// Set up window size callback to dispatch event
	w.SetSizeCallback(func(x *glfw.Window, width int, height int) {
		fbw, fbh := x.GetFramebufferSize()
		w.sizeEv.Width = width
		w.sizeEv.Height = height
		w.scaleX = float64(fbw) / float64(width)
		w.scaleY = float64(fbh) / float64(height)
		w.Dispatch(OnWindowSize, &w.sizeEv)
	})

	// Set up window position callback to dispatch event
	w.SetPosCallback(func(x *glfw.Window, xpos int, ypos int) {
		w.posEv.Xpos = xpos
		w.posEv.Ypos = ypos
		w.Dispatch(OnWindowPos, &w.posEv)
	})

	// Set up window cursor position callback to dispatch event
	w.SetCursorPosCallback(func(x *glfw.Window, xpos float64, ypos float64) {
		w.cursorEv.Xpos = float32(xpos * w.scaleX)
		w.cursorEv.Ypos = float32(ypos * w.scaleY)
		w.cursorEv.Mods = w.mods
		w.Dispatch(OnCursor, &w.cursorEv)
	})

	// Set up mouse wheel scroll callback to dispatch event
	w.SetScrollCallback(func(x *glfw.Window, xoff float64, yoff float64) {
		w.scrollEv.Xoffset = float32(xoff)
		w.scrollEv.Yoffset = float32(yoff)
		w.scrollEv.Mods = w.mods
		w.Dispatch(OnScroll, &w.scrollEv)
	})

	win = w // Set singleton
	return nil
}

// Gls returns the associated OpenGL state.
func (w *GlfwWindow) Gls() *gls.GLS {

	return w.gls
}

// Fullscreen returns whether this windows is currently fullscreen.
func (w *GlfwWindow) Fullscreen() bool {

	return w.fullscreen
}

// SetFullscreen sets this window as fullscreen on the primary monitor
// TODO allow for fullscreen with resolutions different than the monitor's
func (w *GlfwWindow) SetFullscreen(full bool) {

	// If already in the desired state, nothing to do
	if w.fullscreen == full {
		return
	}
	// Set window fullscreen on the primary monitor
	if full {
		// Get size of primary monitor
		mon := glfw.GetPrimaryMonitor()
		vmode := mon.GetVideoMode()
		width := vmode.Width
		height := vmode.Height
		// Set as fullscreen on the primary monitor
		w.SetMonitor(mon, 0, 0, width, height, vmode.RefreshRate)
		w.fullscreen = true
		// Save current position and size of the window
		w.lastX, w.lastY = w.GetPos()
		w.lastWidth, w.lastHeight = w.GetSize()
	} else {
		// Restore window to previous position and size
		w.SetMonitor(nil, w.lastX, w.lastY, w.lastWidth, w.lastHeight, glfw.DontCare)
		w.fullscreen = false
	}
}

// Destroy destroys this window and its context
func (w *GlfwWindow) Destroy() {

	w.Window.Destroy()
	glfw.Terminate()
	runtime.UnlockOSThread() // Important when using the execution tracer
}

// Scale returns this window's DPI scale factor (FramebufferSize / Size)
func (w *GlfwWindow) GetScale() (x float64, y float64) {

	return w.scaleX, w.scaleY
}

// ScreenResolution returns the screen resolution
func (w *GlfwWindow) ScreenResolution(p interface{}) (width, height int) {

	mon := glfw.GetPrimaryMonitor()
	vmode := mon.GetVideoMode()
	return vmode.Width, vmode.Height
}

// PollEvents process events in the event queue
func (w *GlfwWindow) PollEvents() {

	glfw.PollEvents()
}

// SetSwapInterval sets the number of screen updates to wait from the time SwapBuffer()
// is called before swapping the buffers and returning.
func (w *GlfwWindow) SetSwapInterval(interval int) {

	glfw.SwapInterval(interval)
}

// SetCursor sets the window's cursor.
func (w *GlfwWindow) SetCursor(cursor Cursor) {

	cur, ok := w.cursors[cursor]
	if !ok {
		panic("Invalid cursor")
	}
	w.Window.SetCursor(cur)
}

// CreateCursor creates a new custom cursor and returns an int handle.
func (w *GlfwWindow) CreateCursor(imgFile string, xhot, yhot int) (Cursor, error) {

	// Open image file
	file, err := os.Open(imgFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return 0, err
	}
	// Create and store cursor
	w.lastCursorKey += 1
	w.cursors[Cursor(w.lastCursorKey)] = glfw.CreateCursor(img, xhot, yhot)

	return w.lastCursorKey, nil
}

// DisposeCursor deletes the existing custom cursor with the provided int handle.
func (w *GlfwWindow) DisposeCursor(cursor Cursor) {

	if cursor <= CursorLast {
		panic("Can't dispose standard cursor")
	}
	w.cursors[cursor].Destroy()
	delete(w.cursors, cursor)
}

// DisposeAllCursors deletes all existing custom cursors.
func (w *GlfwWindow) DisposeAllCustomCursors() {

	// Destroy and delete all custom cursors
	for key := range w.cursors {
		if key > CursorLast {
			w.cursors[key].Destroy()
			delete(w.cursors, key)
		}
	}
	// Set the next cursor key as the last standard cursor key + 1
	w.lastCursorKey = CursorLast
}

// Center centers the window on the screen.
//func (w *GlfwWindow) Center() {
//
//	// TODO
//}
