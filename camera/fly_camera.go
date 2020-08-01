package camera

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"math"
)


// FlyControl is a camera controller that allows flying through a 3D scene
// It allows the user to rotate camera using mouse and move camera using keyboard.
type FlyControl struct {
	core.Dispatcher                // Embedded event dispatcher
	cam             *Camera        // Controlled camera
	enabled         bool   // If fly control is enabled

	// Public properties
	MinPolarAngle   float32 // Minimum polar angle in radians (default is 0)
	MaxPolarAngle   float32 // Maximum polar angle in radians (default is Pi)
	MinAzimuthAngle float32 // Minimum azimuthal angle in radians (default is negative infinity)
	MaxAzimuthAngle float32 // Maximum azimuthal angle in radians (default is infinity)

	RotSpeed        float32 // Rotation speed factor (default is 0.6)
	FlySpeed     	float32 // Fly speed factor (default is 0.3)

	// Internal
	rotStart math32.Vector2
}

// NewFlyControl creates and returns a pointer to a new fly control for the specified camera
func NewFlyControl(cam *Camera) *FlyControl {

	oc := new(FlyControl)
	oc.Dispatcher.Initialize()
	oc.cam = cam
	oc.enabled = true

	oc.MinPolarAngle = -math32.Pi/2
	oc.MaxPolarAngle = math32.Pi/2 // 90 degrees as radians
	oc.MinAzimuthAngle = float32(math.Inf(-1))
	oc.MaxAzimuthAngle = float32(math.Inf(1))
	oc.RotSpeed = 0.6
	oc.FlySpeed = 0.3


	gui.Manager().SetCursorFocus(oc)

	// Subscribe to events
	gui.Manager().SubscribeID(window.OnKeyDown, &oc, oc.onKey)
	gui.Manager().SubscribeID(window.OnKeyRepeat, &oc, oc.onKey)
	oc.SubscribeID(window.OnCursor, &oc, oc.onCursor)

	return oc
}

// Dispose unsubscribes from all events.
func (oc *FlyControl) Dispose() {
	gui.Manager().UnsubscribeID(window.OnKeyDown, &oc)
	gui.Manager().UnsubscribeID(window.OnKeyRepeat, &oc)
	oc.UnsubscribeID(window.OnCursor, &oc)
}

// Enabled returns the current enabled state
func (oc *FlyControl) Enabled() bool {
	return oc.enabled
}

// SetEnabled sets the current enabled state.
func (oc *FlyControl) SetEnabled(enabled bool) {
	if enabled == false {
		gui.Manager().SetCursorFocus(nil)
	}
	oc.enabled = enabled
}

// Rotate rotates the camera by the specified angles.
func (oc *FlyControl) Rotate(thetaDelta, phiDelta float32) {
	rot := oc.cam.Rotation()

	phi := math32.Clamp(rot.X+(phiDelta*oc.RotSpeed), oc.MinPolarAngle, oc.MaxPolarAngle)
	oc.cam.SetRotationX(phi)

	oc.cam.SetRotationY(rot.Y+(thetaDelta*oc.RotSpeed))
}

// Move moves the camera the specified amount through a 3D scene perpendicular to the viewing direction.
func (oc *FlyControl) Move(deltaX, deltaZ float32) {
	oc.cam.TranslateX(deltaX*oc.FlySpeed)
	oc.cam.TranslateZ(deltaZ*oc.FlySpeed)
}

// onCursor is called when an OnCursor event is received.
func (oc *FlyControl) onCursor(evname string, ev interface{}) {
	if oc.enabled == false {
		return
	}

	mev := ev.(*window.CursorEvent)
	c := -2 * math32.Pi * oc.RotSpeed / oc.winSize()
	oc.Rotate(c*(mev.Xpos-oc.rotStart.X),
		c*(mev.Ypos-oc.rotStart.Y))
	oc.rotStart.Set(mev.Xpos, mev.Ypos)
}

// onKey is called when an OnKeyDown/OnKeyRepeat event is received.
func (oc *FlyControl) onKey(evname string, ev interface{}) {

	if oc.enabled == false {
		return
	}

	deltas := [2]float32{0,0}

	kev := ev.(*window.KeyEvent)
	switch kev.Key {
	case window.KeyUp, window.KeyW:
		deltas[1] = -1
	case window.KeyDown, window.KeyS:
		deltas[1] = 1
	case window.KeyLeft, window.KeyA:
		deltas[0] = -1
	case window.KeyRight, window.KeyD:
		deltas[0] = 1
	}

	if deltas[0] != 0 || deltas[1] != 0{
		oc.Move(deltas[0], deltas[1])
	}
}

// winSize returns the window height or width based on the camera reference axis.
func (oc *FlyControl) winSize() float32 {

	width, size := window.Get().GetSize()
	if oc.cam.Axis() == Horizontal {
		size = width
	}
	return float32(size)
}
