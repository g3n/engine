// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window

import (
	"github.com/g3n/engine/core"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type GLFW struct {
	core.Dispatcher
	win             *glfw.Window
	keyEv           KeyEvent
	charEv          CharEvent
	mouseEv         MouseEvent
	posEv           PosEvent
	sizeEv          SizeEvent
	cursorEv        CursorEvent
	scrollEv        ScrollEvent
	arrowCursor     *glfw.Cursor
	ibeamCursor     *glfw.Cursor
	crosshairCursor *glfw.Cursor
	handCursor      *glfw.Cursor
	hresizeCursor   *glfw.Cursor
	vresizeCursor   *glfw.Cursor
	fullScreen      bool
	lastX           int
	lastY           int
	lastWidth       int
	lastHeight      int
}

// Global GLFW initialization flag
// is initialized when the first window is created
var initialized bool = false

func newGLFW(width, height int, title string, full bool) (*GLFW, error) {

	// Initialize GLFW once before the first window is created
	if !initialized {
		err := glfw.Init()
		if err != nil {
			return nil, err
		}
		// Sets window hints
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.Samples, 8)
		initialized = true
	}

	// Creates window and sets it as the current context.
	// The window is created always as not full screen because if it is
	// created as full screen it not possible to revert it to windowed mode.
	// At the end of this function, the window will be set to full screen if requested.
	win, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}
	win.MakeContextCurrent()

	// Create wrapper window with dispacher
	w := new(GLFW)
	w.win = win
	w.Dispatcher.Initialize()

	// Set key callback to dispatch event
	win.SetKeyCallback(func(x *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

		w.keyEv.W = w
		w.keyEv.Keycode = Key(key)
		w.keyEv.Scancode = scancode
		w.keyEv.Action = Action(action)
		w.keyEv.Mods = ModifierKey(mods)
		if action == glfw.Press {
			w.Dispatch(OnKeyDown, &w.keyEv)
			return
		}
		if action == glfw.Release {
			w.Dispatch(OnKeyUp, &w.keyEv)
			return
		}
		if action == glfw.Repeat {
			w.Dispatch(OnKeyRepeat, &w.keyEv)
			return
		}
	})

	// Set char callback
	win.SetCharModsCallback(func(x *glfw.Window, char rune, mods glfw.ModifierKey) {

		w.charEv.W = w
		w.charEv.Char = char
		w.charEv.Mods = ModifierKey(mods)
		w.Dispatch(OnChar, &w.charEv)
	})

	// Set mouse button callback to dispatch event
	win.SetMouseButtonCallback(func(x *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {

		xpos, ypos := x.GetCursorPos()
		w.mouseEv.W = w
		w.mouseEv.Button = MouseButton(button)
		w.mouseEv.Action = Action(action)
		w.mouseEv.Mods = ModifierKey(mods)
		w.mouseEv.Xpos = float32(xpos)
		w.mouseEv.Ypos = float32(ypos)
		if action == glfw.Press {
			w.Dispatch(OnMouseDown, &w.mouseEv)
			return
		}
		if action == glfw.Release {
			w.Dispatch(OnMouseUp, &w.mouseEv)
			return
		}
	})

	// Set window size callback to dispatch event
	win.SetSizeCallback(func(x *glfw.Window, width int, height int) {

		w.sizeEv.W = w
		w.sizeEv.Width = width
		w.sizeEv.Height = height
		w.Dispatch(OnWindowSize, &w.sizeEv)
	})

	// Set window position event callback to dispatch event
	win.SetPosCallback(func(x *glfw.Window, xpos int, ypos int) {

		w.posEv.W = w
		w.posEv.Xpos = xpos
		w.posEv.Ypos = ypos
		w.Dispatch(OnWindowPos, &w.posEv)
	})

	// Set window cursor position event callback to dispatch event
	win.SetCursorPosCallback(func(x *glfw.Window, xpos float64, ypos float64) {

		w.cursorEv.W = w
		w.cursorEv.Xpos = float32(xpos)
		w.cursorEv.Ypos = float32(ypos)
		w.Dispatch(OnCursor, &w.cursorEv)
	})

	// Set mouse wheel scroll event callback to dispatch event
	win.SetScrollCallback(func(x *glfw.Window, xoff float64, yoff float64) {

		w.scrollEv.W = w
		w.scrollEv.Xoffset = float32(xoff)
		w.scrollEv.Yoffset = float32(yoff)
		w.Dispatch(OnScroll, &w.scrollEv)
	})

	// Preallocate standard cursors
	w.arrowCursor = glfw.CreateStandardCursor(glfw.ArrowCursor)
	w.ibeamCursor = glfw.CreateStandardCursor(glfw.IBeamCursor)
	w.crosshairCursor = glfw.CreateStandardCursor(glfw.CrosshairCursor)
	w.handCursor = glfw.CreateStandardCursor(glfw.HandCursor)
	w.hresizeCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
	w.vresizeCursor = glfw.CreateStandardCursor(glfw.VResizeCursor)

	// Sets full screen if requested
	if full {
		w.SetFullScreen(true)
	}
	return w, nil
}

// GetScreenResolution returns the resolution of the primary screen in pixels.
// The parameter is currently ignored
func (w *GLFW) GetScreenResolution(p interface{}) (width, height int) {

	mon := glfw.GetPrimaryMonitor()
	vmode := mon.GetVideoMode()
	return vmode.Width, vmode.Height
}

func (w *GLFW) SwapInterval(interval int) {

	glfw.SwapInterval(interval)
}

func (w *GLFW) MakeContextCurrent() {

	w.win.MakeContextCurrent()
}

func (w *GLFW) GetSize() (width int, height int) {

	return w.win.GetSize()
}

func (w *GLFW) SetSize(width int, height int) {

	w.win.SetSize(width, height)
}

func (w *GLFW) GetPos() (xpos, ypos int) {

	return w.win.GetPos()
}

func (w *GLFW) SetPos(xpos, ypos int) {

	w.win.SetPos(xpos, ypos)
}

func (w *GLFW) SetTitle(title string) {

	w.win.SetTitle(title)
}

func (w *GLFW) SetStandardCursor(cursor StandardCursor) {

	switch cursor {
	case ArrowCursor:
		w.win.SetCursor(w.arrowCursor)
	case IBeamCursor:
		w.win.SetCursor(w.ibeamCursor)
	case CrosshairCursor:
		w.win.SetCursor(w.crosshairCursor)
	case HandCursor:
		w.win.SetCursor(w.handCursor)
	case HResizeCursor:
		w.win.SetCursor(w.hresizeCursor)
	case VResizeCursor:
		w.win.SetCursor(w.vresizeCursor)
	default:
		panic("Invalid cursor")
	}
}

// FullScreen returns this window full screen state for the primary monitor
func (w *GLFW) FullScreen() bool {

	return w.fullScreen
}

// SetFullScreen sets this window full screen state for the primary monitor
func (w *GLFW) SetFullScreen(full bool) {

	// If already in the desired state, nothing to do
	if w.fullScreen == full {
		return
	}
	// Sets this window full screen for the primary monitor
	if full {
		// Get primary monitor
		mon := glfw.GetPrimaryMonitor()
		vmode := mon.GetVideoMode()
		width := vmode.Width
		height := vmode.Height
		// Saves current position and size of the window
		w.lastX, w.lastY = w.win.GetPos()
		w.lastWidth, w.lastHeight = w.win.GetSize()
		// Sets monitor for full screen
		w.win.SetMonitor(mon, 0, 0, width, height, vmode.RefreshRate)
		w.fullScreen = true
	} else {
		// Restore window to previous position and size
		w.win.SetMonitor(nil, w.lastX, w.lastY, w.lastWidth, w.lastHeight, glfw.DontCare)
		w.fullScreen = false
	}
}

// ShouldClose returns the current state of this window  should close flag
func (w *GLFW) ShouldClose() bool {

	return w.win.ShouldClose()
}

// SetShouldClose sets the state of this windows should close flag
func (w *GLFW) SetShouldClose(v bool) {

	w.win.SetShouldClose(v)
}

func (w *GLFW) SwapBuffers() {

	w.win.SwapBuffers()
}

func (w *GLFW) Destroy() {

	w.win.Destroy()
	w.win = nil
}

func (w *GLFW) PollEvents() {

	glfw.PollEvents()
}

func (w *GLFW) GetTime() float64 {

	return glfw.GetTime()
}
