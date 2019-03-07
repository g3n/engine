// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package window abstracts the OpenGL Window manager
// Currently only "glfw" is supported
package window

import (
	"github.com/g3n/engine/core"
)

// IWindowManager is the interface for all window managers
type IWindowManager interface {
	ScreenResolution(interface{}) (width, height int)
	CreateWindow(width, height int, title string, full bool) (IWindow, error)
	CreateCursor(imgFile string, xhot, yhot int) (int, error)
	DisposeCursor(key int)
	DisposeAllCursors()
	SetSwapInterval(interval int)
	PollEvents()
	Terminate()
}

// IWindow is the interface for all windows
type IWindow interface {
	core.IDispatcher
	Manager() IWindowManager
	MakeContextCurrent()
	FramebufferSize() (width int, height int)
	Scale() (x float64, y float64)
	Size() (width int, height int)
	SetSize(width int, height int)
	Pos() (xpos, ypos int)
	SetPos(xpos, ypos int)
	SetTitle(title string)
	SetStandardCursor(cursor StandardCursor)
	SetCustomCursor(int)
	SetInputMode(mode InputMode, state int)
	SetCursorPos(xpos, ypos float64)
	ShouldClose() bool
	SetShouldClose(bool)
	FullScreen() bool
	SetFullScreen(bool)
	SwapBuffers()
	Destroy()
}

// Key corresponds to a keyboard key.
type Key int

// ModifierKey corresponds to a modifier key.
type ModifierKey int

// MouseButton corresponds to a mouse button.
type MouseButton int

// StandardCursor corresponds to a g3n standard cursor icon.
type StandardCursor int

// Action corresponds to a key or button action.
type Action int

// InputMode corresponds to an input mode.
type InputMode int

// InputMode corresponds to an input mode.
type CursorMode int

//
// Window event names using for dispatch and subscribe
//
const (
	OnWindowPos  = "win.OnWindowPos"
	OnWindowSize = "win.OnWindowSize"
	OnKeyUp      = "win.OnKeyUp"
	OnKeyDown    = "win.OnKeyDown"
	OnKeyRepeat  = "win.OnKeyRepeat"
	OnChar       = "win.OnChar"
	OnCursor     = "win.OnCursor"
	OnMouseUp    = "win.OnMouseUp"
	OnMouseDown  = "win.OnMouseDown"
	OnScroll     = "win.OnScroll"
	OnFrame      = "win.OnFrame"
)

// PosEvent describes a windows position changed event
type PosEvent struct {
	W    IWindow
	Xpos int
	Ypos int
}

// SizeEvent describers a window size changed event
type SizeEvent struct {
	W      IWindow
	Width  int
	Height int
}

// KeyEvent describes a window key event
type KeyEvent struct {
	W        IWindow
	Keycode  Key
	Scancode int
	Action   Action
	Mods     ModifierKey
}

// CharEvent describes a window char event
type CharEvent struct {
	W    IWindow
	Char rune
	Mods ModifierKey
}

// MouseEvent describes a mouse event over the window
type MouseEvent struct {
	W      IWindow
	Xpos   float32
	Ypos   float32
	Button MouseButton
	Action Action
	Mods   ModifierKey
}

// CursorEvent describes a cursor position changed event
type CursorEvent struct {
	W    IWindow
	Xpos float32
	Ypos float32
}

// ScrollEvent describes a scroll event
type ScrollEvent struct {
	W       IWindow
	Xoffset float32
	Yoffset float32
}
