// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window

import (
	"runtime"

	"bytes"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui/assets"
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

// Standard cursors for g3n. The diagonal cursors are not standard for GLFW.
const (
	ArrowCursor       = StandardCursor(glfw.ArrowCursor)
	IBeamCursor       = StandardCursor(glfw.IBeamCursor)
	CrosshairCursor   = StandardCursor(glfw.CrosshairCursor)
	HandCursor        = StandardCursor(glfw.HandCursor)
	HResizeCursor     = StandardCursor(glfw.HResizeCursor)
	VResizeCursor     = StandardCursor(glfw.VResizeCursor)
	DiagResize1Cursor = StandardCursor(VResizeCursor + 1)
	DiagResize2Cursor = StandardCursor(VResizeCursor + 2)
)

// Actions
const (
	// Release indicates that key or mouse button was released
	Release = Action(glfw.Release)
	// Press indicates that key or mouse button was pressed
	Press = Action(glfw.Press)
	// Repeat indicates that key was held down until it repeated
	Repeat = Action(glfw.Repeat)
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

// glfwManager contains data shared by all windows
type glfwManager struct {
	arrowCursor     *glfw.Cursor // Preallocated standard arrow cursor
	ibeamCursor     *glfw.Cursor // Preallocated standard ibeam cursor
	crosshairCursor *glfw.Cursor // Preallocated standard cross hair cursor
	handCursor      *glfw.Cursor // Preallocated standard hand cursor
	hresizeCursor   *glfw.Cursor // Preallocated standard horizontal resize cursor
	vresizeCursor   *glfw.Cursor // Preallocated standard vertical resize cursor

	// Non GLFW standard cursors (but g3n standard)
	diag1Cursor *glfw.Cursor // Preallocated diagonal resize cursor (/)
	diag2Cursor *glfw.Cursor // Preallocated diagonal resize cursor (\)

	// User-created custom cursors
	customCursors map[int]*glfw.Cursor
	lastCursorKey int
}

// glfwWindow describes one glfw window
type glfwWindow struct {
	core.Dispatcher              // Embedded event dispatcher
	gls             *gls.GLS     // Associated OpenGL State
	win             *glfw.Window // Pointer to native glfw window
	mgr             *glfwManager // Pointer to window manager
	fullScreen      bool
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
}

// glfw manager singleton
var manager *glfwManager

// Manager returns the glfw window manager
func Manager() (IWindowManager, error) {

	if manager != nil {
		return manager, nil
	}

	// Initialize glfw
	err := glfw.Init()
	if err != nil {
		return nil, err
	}

	// Sets window hints
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Samples, 8)
	// Sets OpenGL forward compatible context only for OSX because it is required for OSX.
	// When this is set, glLineWidth(width) only accepts width=1.0 and generates an error
	// for any other values although the spec says it should ignore unsupported widths
	// and generate an error only when width <= 0.
	if runtime.GOOS == "darwin" {
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	}

	manager = new(glfwManager)

	// Preallocate GLFW standard cursors
	manager.arrowCursor = glfw.CreateStandardCursor(glfw.ArrowCursor)
	manager.ibeamCursor = glfw.CreateStandardCursor(glfw.IBeamCursor)
	manager.crosshairCursor = glfw.CreateStandardCursor(glfw.CrosshairCursor)
	manager.handCursor = glfw.CreateStandardCursor(glfw.HandCursor)
	manager.hresizeCursor = glfw.CreateStandardCursor(glfw.HResizeCursor)
	manager.vresizeCursor = glfw.CreateStandardCursor(glfw.VResizeCursor)

	// Preallocate g3n cursors (diagonal cursors)
	cursorDiag1Png := assets.MustAsset("cursors/diag1.png")
	cursorDiag2Png := assets.MustAsset("cursors/diag2.png")

	diag1Img, _, err := image.Decode(bytes.NewReader(cursorDiag1Png))
	diag2Img, _, err := image.Decode(bytes.NewReader(cursorDiag2Png))
	if err != nil {
		return nil, err
	}
	manager.diag1Cursor = glfw.CreateCursor(diag1Img, 8, 8)
	manager.diag2Cursor = glfw.CreateCursor(diag2Img, 8, 8)

	// Create map for custom cursors
	manager.customCursors = make(map[int]*glfw.Cursor)

	return manager, nil
}

// ScreenResolution returns the screen resolution
func (m *glfwManager) ScreenResolution(p interface{}) (width, height int) {

	mon := glfw.GetPrimaryMonitor()
	vmode := mon.GetVideoMode()
	return vmode.Width, vmode.Height
}

// PollEvents process events in the event queue
func (m *glfwManager) PollEvents() {

	glfw.PollEvents()
}

// SetSwapInterval sets the number of screen updates to wait from the time SwapBuffer()
// is called before swapping the buffers and returning.
func (m *glfwManager) SetSwapInterval(interval int) {

	glfw.SwapInterval(interval)
}

// Terminate destroys any remaining window, cursors and other related objects.
func (m *glfwManager) Terminate() {

	glfw.Terminate()
	manager = nil
}

// CreateCursor creates a new custom cursor and returns an int handle.
func (m *glfwManager) CreateCursor(imgFile string, xhot, yhot int) (int, error) {

	// Open image file
	file, err := os.Open(imgFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Decodes image
	img, _, err := image.Decode(file)
	if err != nil {
		return 0, err
	}

	cursor := glfw.CreateCursor(img, xhot, yhot)
	if err != nil {
		return 0, err
	}
	m.lastCursorKey += 1
	m.customCursors[m.lastCursorKey] = cursor

	return m.lastCursorKey, nil
}

// DisposeCursor deletes the existing custom cursor with the provided int handle.
func (m *glfwManager) DisposeCursor(key int) {

	delete(m.customCursors, key)
}

// DisposeAllCursors deletes all existing custom cursors.
func (m *glfwManager) DisposeAllCursors() {

	m.customCursors = make(map[int]*glfw.Cursor)
	m.lastCursorKey = 0
}

// CreateWindow creates and returns a new window with the specified width and height in screen coordinates
func (m *glfwManager) CreateWindow(width, height int, title string, fullscreen bool) (IWindow, error) {

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
	w := new(glfwWindow)
	w.win = win
	w.mgr = m
	w.Dispatcher.Initialize()

	fbw, fbh := w.FramebufferSize()
	w.scaleX = float64(fbw) / float64(width)
	w.scaleY = float64(fbh) / float64(height)

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
		w.mouseEv.Xpos = float32(xpos * w.scaleX)
		w.mouseEv.Ypos = float32(ypos * w.scaleY)

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

		fbw, fbh := x.GetFramebufferSize()
		w.sizeEv.W = w
		w.sizeEv.Width = width
		w.sizeEv.Height = height
		w.scaleX = float64(fbw) / float64(width)
		w.scaleY = float64(fbh) / float64(height)
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
		w.cursorEv.Xpos = float32(xpos * w.scaleX)
		w.cursorEv.Ypos = float32(ypos * w.scaleY)
		w.Dispatch(OnCursor, &w.cursorEv)
	})

	// Set mouse wheel scroll event callback to dispatch event
	win.SetScrollCallback(func(x *glfw.Window, xoff float64, yoff float64) {

		w.scrollEv.W = w
		w.scrollEv.Xoffset = float32(xoff)
		w.scrollEv.Yoffset = float32(yoff)
		w.Dispatch(OnScroll, &w.scrollEv)
	})

	// Sets full screen if requested
	if fullscreen {
		w.SetFullScreen(true)
	}

	// Create OpenGL state
	gl, err := gls.New()
	if err != nil {
		return nil, err
	}
	w.gls = gl

	return w, nil
}

// Gls returns the associated OpenGL state
func (w *glfwWindow) Gls() *gls.GLS {

	return w.gls
}

// Manager returns the window manager and satisfies the IWindow interface
func (w *glfwWindow) Manager() IWindowManager {

	return w.mgr
}

// MakeContextCurrent makes the OpenGL context of this window current on the calling thread
func (w *glfwWindow) MakeContextCurrent() {

	w.win.MakeContextCurrent()
}

// FullScreen returns this window full screen state for the primary monitor
func (w *glfwWindow) FullScreen() bool {

	return w.fullScreen
}

// SetFullScreen sets this window full screen state for the primary monitor
func (w *glfwWindow) SetFullScreen(full bool) {

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

// Destroy destroys this window and its context
func (w *glfwWindow) Destroy() {

	w.win.Destroy()
	w.win = nil
}

// SwapBuffers swaps the front and back buffers of this window.
// If the swap interval is greater than zero,
// the GPU driver waits the specified number of screen updates before swapping the buffers.
func (w *glfwWindow) SwapBuffers() {

	w.win.SwapBuffers()
}

// FramebufferSize returns framebuffer size of this window
func (w *glfwWindow) FramebufferSize() (width int, height int) {

	return w.win.GetFramebufferSize()
}

// Scale returns this window's DPI scale factor (FramebufferSize / Size)
func (w *glfwWindow) Scale() (x float64, y float64) {

	return w.scaleX, w.scaleY
}

// Size returns this window's size in screen coordinates
func (w *glfwWindow) Size() (width int, height int) {

	return w.win.GetSize()
}

// SetSize sets the size, in screen coordinates, of the client area of this window
func (w *glfwWindow) SetSize(width int, height int) {

	w.win.SetSize(width, height)
}

// Pos returns the position, in screen coordinates, of the upper-left corner of the client area of this window
func (w *glfwWindow) Pos() (xpos, ypos int) {

	return w.win.GetPos()
}

// SetPos sets the position, in screen coordinates, of the upper-left corner of the client area of this window.
// If the window is a full screen window, this function does nothing.
func (w *glfwWindow) SetPos(xpos, ypos int) {

	w.win.SetPos(xpos, ypos)
}

// SetTitle sets this window title, encoded as UTF-8
func (w *glfwWindow) SetTitle(title string) {

	w.win.SetTitle(title)
}

// ShouldClose returns the current state of this window  should close flag
func (w *glfwWindow) ShouldClose() bool {

	return w.win.ShouldClose()
}

// SetShouldClose sets the state of this windows should close flag
func (w *glfwWindow) SetShouldClose(v bool) {

	w.win.SetShouldClose(v)
}

// SetStandardCursor sets the window's cursor to a standard one
func (w *glfwWindow) SetStandardCursor(cursor StandardCursor) {

	switch cursor {
	case ArrowCursor:
		w.win.SetCursor(w.mgr.arrowCursor)
	case IBeamCursor:
		w.win.SetCursor(w.mgr.ibeamCursor)
	case CrosshairCursor:
		w.win.SetCursor(w.mgr.crosshairCursor)
	case HandCursor:
		w.win.SetCursor(w.mgr.handCursor)
	case HResizeCursor:
		w.win.SetCursor(w.mgr.hresizeCursor)
	case VResizeCursor:
		w.win.SetCursor(w.mgr.vresizeCursor)
	// Non-GLFW cursors (but standard cursors for g3n)
	case DiagResize1Cursor:
		w.win.SetCursor(w.mgr.diag1Cursor)
	case DiagResize2Cursor:
		w.win.SetCursor(w.mgr.diag2Cursor)
	default:
		panic("Invalid cursor")
	}
}

// SetStandardCursor sets this window's cursor to a custom, user-created one
func (w *glfwWindow) SetCustomCursor(key int) {

	w.win.SetCursor(w.mgr.customCursors[key])
}

// SetInputMode changes specified input to specified state
// Reference: http://www.glfw.org/docs/latest/group__input.html#gaa92336e173da9c8834558b54ee80563b
func (w *glfwWindow) SetInputMode(mode InputMode, state int) {

	w.win.SetInputMode(glfw.InputMode(mode), state)
}

// SetCursorPos sets cursor position in window coordinates
// Reference: http://www.glfw.org/docs/latest/group__input.html#ga04b03af936d906ca123c8f4ee08b39e7
func (w *glfwWindow) SetCursorPos(xpos, ypos float64) {

	w.win.SetCursorPos(xpos, ypos)
}
