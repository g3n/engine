// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package camera

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

// FPEnabled specifies which control types are enabled.
type FPEnabled int

// The possible control types.
const (
	FPNone   FPEnabled = 0x00
	FPRot    FPEnabled = 0x01
	FPMoveXZ FPEnabled = 0x02
	FPMoveY  FPEnabled = 0x04
	FPKeys   FPEnabled = 0x08
	FPAll    FPEnabled = 0xFF
)

// fpState bitmask
type fpState int

const (
	fpStateNone = fpState(iota)
	fpStateRotate
	fpStateMoveHorizontal
	fpStateMoveVertical
	fpStateZoom
)

// FirstPersonControl is a camera controller that allows looking around from a static point.
// It allows the user to rotate, zoom, and pan a 3D scene using the mouse or keyboard.
type FirstPersonControl struct {
	core.Dispatcher           // Embedded event dispatcher
	cam             *Camera   // Controlled camera
	enabled         FPEnabled // Which controls are enabled
	state           fpState   // Current control state

	// Public properties
	RotSpeed     float32 // Rotation speed factor (default is 1)
	MoveSpeed    float32 // Move speed factor (default is 0.1)
	KeyRotSpeed  float32 // Rotation delta in radians used on each rotation key event (default is the equivalent of 15 degrees)
	KeyMoveSpeed float32 // Move delta used on each move key event (default is 0.5)

	// Internal
	rotStart  math32.Vector2
	moveStart math32.Vector3
}

// NewFirstPersonControl creates and returns a pointer to a new first person control for the specified camera.
func NewFirstPersonControl(cam *Camera) *FirstPersonControl {

	fpc := new(FirstPersonControl)
	fpc.Dispatcher.Initialize()
	fpc.cam = cam
	fpc.enabled = FPAll

	fpc.RotSpeed = 1.0
	fpc.MoveSpeed = 0.1
	fpc.KeyRotSpeed = 15 * math32.Pi / 180 // 15 degrees as radians
	fpc.KeyMoveSpeed = 0.5

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

// Move moves the camera the specified amount relative from the camera position based on its view direction.
func (fpc *FirstPersonControl) Move(deltaX, deltaY, deltaZ float32) {

	// TODO: The Translate calls handle relaitive camera movement well, however it "flies" in the direction
	// pointed at rather than "walks" along the XZ plane as desired for this purpose. A new function
	// for "flying" could be useful but still need to figure out movement without flying.

	// Translate camera by deltas
	fpc.cam.TranslateX(deltaX)
	fpc.cam.TranslateY(deltaY)
	fpc.cam.TranslateZ(deltaZ)
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
		case window.MouseButtonMiddle: // Zoom movement only
			if fpc.enabled&FPMoveXZ != 0 {
				fpc.state = fpStateZoom
				fpc.moveStart = math32.Vector3{X: 0, Y: 0, Z: mev.Ypos}
			}
		case window.MouseButtonRight: // Move horizontal only
			if fpc.enabled&FPMoveXZ != 0 {
				fpc.state = fpStateMoveHorizontal
				fpc.moveStart = math32.Vector3{X: mev.Xpos, Y: 0, Z: mev.Ypos}
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
		fpc.Move(0, 0, fpc.MoveSpeed*(mev.Ypos-fpc.moveStart.Z))
		fpc.moveStart = math32.Vector3{X: 0, Y: 0, Z: mev.Ypos}
	case fpStateMoveHorizontal:
		fpc.Move(fpc.MoveSpeed*(mev.Xpos-fpc.moveStart.X), 0, fpc.MoveSpeed*(mev.Ypos-fpc.moveStart.Z))
		fpc.moveStart = math32.Vector3{X: mev.Xpos, Y: 0, Z: mev.Ypos}
	}
}

// onScroll is called when an OnScroll event is received.
func (fpc *FirstPersonControl) onScroll(evname string, ev interface{}) {

	if fpc.enabled&FPMoveXZ != 0 {
		sev := ev.(*window.ScrollEvent)
		fpc.Move(0, 0, fpc.KeyMoveSpeed*(-sev.Yoffset))
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
			fpc.Rotate(0, fpc.KeyRotSpeed)
		case window.KeyDown:
			fpc.Rotate(0, -fpc.KeyRotSpeed)
		case window.KeyLeft, window.KeyQ:
			fpc.Rotate(fpc.KeyRotSpeed, 0)
		case window.KeyRight, window.KeyE:
			fpc.Rotate(-fpc.KeyRotSpeed, 0)
		}
	}
	if kev.Mods == 0 && fpc.enabled&FPMoveY != 0 {
		switch kev.Key {
		case window.KeyR:
			fpc.Move(0, fpc.KeyMoveSpeed, 0)
		case window.KeyF:
			fpc.Move(0, -fpc.KeyMoveSpeed, 0)
		}
	}
	if kev.Mods == 0 && fpc.enabled&FPMoveXZ != 0 {
		switch kev.Key {
		case window.KeyW:
			fpc.Move(0, 0, -fpc.KeyMoveSpeed)
		case window.KeyS:
			fpc.Move(0, 0, fpc.KeyMoveSpeed)
		case window.KeyA:
			fpc.Move(-fpc.KeyMoveSpeed, 0, 0)
		case window.KeyD:
			fpc.Move(fpc.KeyMoveSpeed, 0, 0)
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
