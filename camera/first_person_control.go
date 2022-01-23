// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package camera

import (
	"math"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

// FPEnabled specifies which control types are enabled.
type FPEnabled int

// The possible control types.
const (
	FPNone FPEnabled = 0x00
	FPRot  FPEnabled = 0x01
	FPZoom FPEnabled = 0x02
	FPPan  FPEnabled = 0x04
	FPKeys FPEnabled = 0x08
	FPAll  FPEnabled = 0xFF
)

// fpState bitmask
type fpState int

const (
	fpStateNone = fpState(iota)
	fpStateRotate
	fpStateZoom
	fpStatePan
)

// FirstPersonControl is a camera controller that allows looking around from a static point.
// It allows the user to rotate, zoom, and pan a 3D scene using the mouse or keyboard.
type FirstPersonControl struct {
	core.Dispatcher           // Embedded event dispatcher
	cam             *Camera   // Controlled camera
	enabled         FPEnabled // Which controls are enabled
	state           fpState   // Current control state

	// Public properties
	MinDistance     float32 // Minimum distance from target (default is 1)
	MaxDistance     float32 // Maximum distance from target (default is infinity)
	MinPolarAngle   float32 // Minimum polar angle in radians (default is 0)
	MaxPolarAngle   float32 // Maximum polar angle in radians (default is Pi)
	MinAzimuthAngle float32 // Minimum azimuthal angle in radians (default is negative infinity)
	MaxAzimuthAngle float32 // Maximum azimuthal angle in radians (default is infinity)
	RotSpeed        float32 // Rotation speed factor (default is 1)
	ZoomSpeed       float32 // Zoom speed factor (default is 0.1)
	KeyRotSpeed     float32 // Rotation delta in radians used on each rotation key event (default is the equivalent of 15 degrees)
	KeyZoomSpeed    float32 // Zoom delta used on each zoom key event (default is 2)
	KeyPanSpeed     float32 // Pan delta used on each pan key event (default is 35)

	// Internal
	rotStart  math32.Vector2
	panStart  math32.Vector2
	zoomStart float32
}

// NewFirstPersonControl creates and returns a pointer to a new first person control for the specified camera.
func NewFirstPersonControl(cam *Camera) *FirstPersonControl {

	fpc := new(FirstPersonControl)
	fpc.Dispatcher.Initialize()
	fpc.cam = cam
	fpc.enabled = FPAll

	fpc.MinDistance = 1.0
	fpc.MaxDistance = float32(math.Inf(1))
	fpc.MinPolarAngle = 0
	fpc.MaxPolarAngle = math32.Pi // 180 degrees as radians
	fpc.MinAzimuthAngle = float32(math.Inf(-1))
	fpc.MaxAzimuthAngle = float32(math.Inf(1))
	fpc.RotSpeed = 1.0
	fpc.ZoomSpeed = 0.1
	fpc.KeyRotSpeed = 15 * math32.Pi / 180 // 15 degrees as radians
	fpc.KeyZoomSpeed = 2.0
	fpc.KeyPanSpeed = 35.0

	// initialize Z axis rotation to zero for proper X/Y only view rotation
	fpc.cam.SetRotationZ(0)

	// Subscribe to events
	gui.Manager().SubscribeID(window.OnMouseUp, &fpc, fpc.onMouse)
	gui.Manager().SubscribeID(window.OnMouseDown, &fpc, fpc.onMouse)
	gui.Manager().SubscribeID(window.OnScroll, &fpc, fpc.onScroll)
	gui.Manager().SubscribeID(window.OnKeyDown, &fpc, fpc.onKey)
	gui.Manager().SubscribeID(window.OnKeyRepeat, &fpc, fpc.onKey)
	fpc.SubscribeID(window.OnCursor, &fpc, fpc.onCursor)

	return fpc
}

// Dispose unsubscribes from all events.
func (fpc *FirstPersonControl) Dispose() {

	gui.Manager().UnsubscribeID(window.OnMouseUp, &fpc)
	gui.Manager().UnsubscribeID(window.OnMouseDown, &fpc)
	gui.Manager().UnsubscribeID(window.OnScroll, &fpc)
	gui.Manager().UnsubscribeID(window.OnKeyDown, &fpc)
	gui.Manager().UnsubscribeID(window.OnKeyRepeat, &fpc)
	fpc.UnsubscribeID(window.OnCursor, &fpc)
}

// Enabled returns the current FPEnabled bitmask.
func (fpc *FirstPersonControl) Enabled() FPEnabled {

	return fpc.enabled
}

// SetEnabled sets the current FPEnabled bitmask.
func (fpc *FirstPersonControl) SetEnabled(bitmask FPEnabled) {

	fpc.enabled = bitmask
}

// Rotate rotates the camera in place by the specified angles.
func (fpc *FirstPersonControl) Rotate(thetaDelta, phiDelta float32) {

	// Rotate in place by rotation only in the X and Y axis
	r := fpc.cam.Rotation()
	fpc.cam.SetRotation(r.X+phiDelta, r.Y+thetaDelta, r.Z)
}

// Zoom moves the camera closer or farther from the target the specified amount
// and also updates the camera's orthographic size to match.
func (fpc *FirstPersonControl) Zoom(delta float32) {

	// TODO: rename as Move to allow movement on X/Z plane (or Y instead of Z?)
}

// Pan pans the camera and target the specified amount on the plane perpendicular to the viewing direction.
func (fpc *FirstPersonControl) Pan(deltaX, deltaY float32) {

	// TODO: same as Move?
}

// onMouse is called when an OnMouseDown/OnMouseUp event is received.
func (fpc *FirstPersonControl) onMouse(evname string, ev interface{}) {

	// If nothing enabled ignore event
	if fpc.enabled == FPNone {
		return
	}

	switch evname {
	case window.OnMouseDown:
		gui.Manager().SetCursorFocus(fpc)
		mev := ev.(*window.MouseEvent)
		switch mev.Button {
		case window.MouseButtonLeft: // Rotate
			if fpc.enabled&FPRot != 0 {
				fpc.state = fpStateRotate
				fpc.rotStart.Set(mev.Xpos, mev.Ypos)
			}
		case window.MouseButtonMiddle: // Zoom
			if fpc.enabled&FPZoom != 0 {
				fpc.state = fpStateZoom
				fpc.zoomStart = mev.Ypos
			}
		case window.MouseButtonRight: // Pan
			if fpc.enabled&FPPan != 0 {
				fpc.state = fpStatePan
				fpc.panStart.Set(mev.Xpos, mev.Ypos)
			}
		}
	case window.OnMouseUp:
		gui.Manager().SetCursorFocus(nil)
		fpc.state = fpStateNone
	}
}

// onCursor is called when an OnCursor event is received.
func (fpc *FirstPersonControl) onCursor(evname string, ev interface{}) {

	// If nothing enabled ignore event
	if fpc.enabled == FPNone || fpc.state == fpStateNone {
		return
	}

	mev := ev.(*window.CursorEvent)
	switch fpc.state {
	case fpStateRotate:
		c := -2 * math32.Pi * fpc.RotSpeed / fpc.winSize()
		fpc.Rotate(c*(mev.Xpos-fpc.rotStart.X),
			c*(mev.Ypos-fpc.rotStart.Y))
		fpc.rotStart.Set(mev.Xpos, mev.Ypos)
	case fpStateZoom:
		fpc.Zoom(fpc.ZoomSpeed * (mev.Ypos - fpc.zoomStart))
		fpc.zoomStart = mev.Ypos
	case fpStatePan:
		fpc.Pan(mev.Xpos-fpc.panStart.X,
			mev.Ypos-fpc.panStart.Y)
		fpc.panStart.Set(mev.Xpos, mev.Ypos)
	}
}

// onScroll is called when an OnScroll event is received.
func (fpc *FirstPersonControl) onScroll(evname string, ev interface{}) {

	if fpc.enabled&FPZoom != 0 {
		sev := ev.(*window.ScrollEvent)
		fpc.Zoom(-sev.Yoffset)
	}
}

// onKey is called when an OnKeyDown/OnKeyRepeat event is received.
func (fpc *FirstPersonControl) onKey(evname string, ev interface{}) {

	// If keyboard control is disabled ignore event
	if fpc.enabled&FPKeys == 0 {
		return
	}

	kev := ev.(*window.KeyEvent)
	if kev.Mods == 0 && fpc.enabled&FPRot != 0 {
		switch kev.Key {
		case window.KeyUp:
			fpc.Rotate(0, -fpc.KeyRotSpeed)
		case window.KeyDown:
			fpc.Rotate(0, fpc.KeyRotSpeed)
		case window.KeyLeft:
			fpc.Rotate(-fpc.KeyRotSpeed, 0)
		case window.KeyRight:
			fpc.Rotate(fpc.KeyRotSpeed, 0)
		}
	}
	if kev.Mods == window.ModControl && fpc.enabled&FPZoom != 0 {
		switch kev.Key {
		case window.KeyUp:
			fpc.Zoom(-fpc.KeyZoomSpeed)
		case window.KeyDown:
			fpc.Zoom(fpc.KeyZoomSpeed)
		}
	}
	if kev.Mods == window.ModShift && fpc.enabled&FPPan != 0 {
		switch kev.Key {
		case window.KeyUp:
			fpc.Pan(0, fpc.KeyPanSpeed)
		case window.KeyDown:
			fpc.Pan(0, -fpc.KeyPanSpeed)
		case window.KeyLeft:
			fpc.Pan(fpc.KeyPanSpeed, 0)
		case window.KeyRight:
			fpc.Pan(-fpc.KeyPanSpeed, 0)
		}
	}
}

// winSize returns the window height or width based on the camera reference axis.
func (fpc *FirstPersonControl) winSize() float32 {

	width, size := window.Get().GetSize()
	if fpc.cam.Axis() == Horizontal {
		size = width
	}
	return float32(size)
}
