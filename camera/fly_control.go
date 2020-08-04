package camera

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"math"
	"time"
)


// FlyControl is a camera controller that allows flying through a 3D scene
// It allows the user to rotate camera using mouse and move camera using keyboard.
type FlyControl struct {
	core.Dispatcher                // Embedded event dispatcher
	cam             *Camera        // Controlled camera
	win 			window.IWindow
	keyState 		*window.KeyState
	enabled         bool   // If fly control is enabled

	// Public properties
	MinPolarAngle   float32 // Minimum polar angle in radians (default is -Pi/2)
	MaxPolarAngle   float32 // Maximum polar angle in radians (default is Pi/2)
	MinAzimuthAngle float32 // Minimum azimuthal angle in radians (default is negative infinity)
	MaxAzimuthAngle float32 // Maximum azimuthal angle in radians (default is infinity)

	Sensitivity     float32 // Mouse sensitivity (default is 1)
	Speed     		float32 // Fly speed factor (default is 1)

	// Internal
	cursorPosition [2]float32
}

// NewFlyControl creates and returns a pointer to a new fly control for the specified camera
func NewFlyControl(cam *Camera) *FlyControl {

	fc := new(FlyControl)
	fc.Dispatcher.Initialize()
	fc.cam = cam
	fc.win = window.Get()
	fc.keyState = window.NewKeyState(fc.win)
	fc.enabled = true

	fc.MinPolarAngle = -math32.Pi/2
	fc.MaxPolarAngle = math32.Pi/2 // 90 degrees as radians
	fc.MinAzimuthAngle = float32(math.Inf(-1))
	fc.MaxAzimuthAngle = float32(math.Inf(1))
	fc.Sensitivity = 1
	fc.Speed = 1

	fc.cursorPosition = [2]float32{0,0}

	fc.win.SetCursorMode(window.CursorDisabled)

	return fc
}

// Dispose unsubscribes from all events.
func (fc *FlyControl) Dispose() {
	fc.keyState.Dispose()
}

// Enabled returns the current enabled state
func (fc *FlyControl) Enabled() bool {
	return fc.enabled
}

// SetEnabled sets the current enabled state.
func (fc *FlyControl) SetEnabled(enabled bool) {
	if enabled == false {
		fc.win.SetCursorMode(window.CursorNormal)
		gui.Manager().SetCursorFocus(nil)
	} else {
		cursorX, cursorY := fc.win.CursorPosition()
		fc.cursorPosition[0] = float32(cursorX)
		fc.cursorPosition[1] = float32(cursorY)

		fc.win.SetCursorMode(window.CursorDisabled)

		gui.Manager().SetCursorFocus(fc)
	}
	fc.enabled = enabled
}

// Rotate rotates the camera by the specified angles.
func (fc *FlyControl) rotate(thetaDelta, phiDelta float32) {
	rot := fc.cam.Rotation()

	phi := math32.Clamp(rot.X-(phiDelta*fc.Sensitivity), fc.MinPolarAngle, fc.MaxPolarAngle)
	fc.cam.SetRotationX(phi)

	fc.cam.SetRotationY(rot.Y - (thetaDelta * fc.Sensitivity))
}

// Move moves the camera the specified amount through a 3D scene perpendicular to the viewing direction.
func (fc *FlyControl) move(deltaTime, x, z float32) {
	if x == 0 && z == 0 {
		return
	}

	camPos := fc.cam.Position()
	camQuaternion := fc.cam.Quaternion()

	deltaX := x * 3 * deltaTime * fc.Speed
	v := math32.NewVector3(1, 0, 0)
	v.ApplyQuaternion(&camQuaternion)
	v.MultiplyScalar(deltaX)
	camPos.Add(v)

	deltaZ := z * 3 * deltaTime * fc.Speed
	v = math32.NewVector3(0, 0, 1)
	v.ApplyQuaternion(&camQuaternion)
	v.MultiplyScalar(deltaZ)
	camPos.Add(v)

	fc.cam.SetPositionVec(&camPos)
}

// Update should be called in application update loop
func (fc *FlyControl) Update(deltaTime time.Duration) {
	if fc.enabled == false {
		return
	}
	dt := float32(deltaTime.Seconds())

	if fc.win.CursorMode() == window.CursorDisabled {
		cursorX, cursorY := fc.win.CursorPosition()
		//println(fmt.Sprintf("Cursor x %v, y %v", cursorX, cursorY))

		deltaX := (cursorX - fc.cursorPosition[0]) * 0.1 * dt
		deltaY := (cursorY - fc.cursorPosition[1]) * 0.1 * dt
		fc.rotate(deltaX, deltaY)
		fc.cursorPosition[0] = cursorX
		fc.cursorPosition[1] = cursorY
	}

	deltas := [2]float32{0, 0}
	for key, _ := range fc.keyState.PressedKeys() {
		switch key {
		case window.KeyUp, window.KeyW:
			deltas[1] = -1
		case window.KeyDown, window.KeyS:
			deltas[1] = 0.8
		case window.KeyLeft, window.KeyA:
			deltas[0] = -0.66
		case window.KeyRight, window.KeyD:
			deltas[0] = 0.66
		}
	}

	if deltas[0] != 0 || deltas[1] != 0 {
		fc.move(dt, deltas[0], deltas[1])
	}
}

// winSize returns the window height or width based on the camera reference axis.
func (fc *FlyControl) winSize() float32 {

	width, size := window.Get().GetSize()
	if fc.cam.Axis() == Horizontal {
		size = width
	}
	return float32(size)
}