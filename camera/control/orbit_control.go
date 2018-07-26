// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package control

import (
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/logger"
	"github.com/g3n/engine/window"
	"math"
)

// OrbitControl is a camera controller that allows orbiting a center point while looking at it.
type OrbitControl struct {
	Enabled         bool    // Control enabled state
	EnableRotate    bool    // Rotate enabled state
	EnableZoom      bool    // Zoom enabled state
	EnablePan       bool    // Pan enabled state
	EnableKeys      bool    // Enable keys state
	ZoomSpeed       float32 // Zoom speed factor. Default is 1.0
	RotateSpeed     float32 // Rotate speed factor. Default is 1.0
	MinDistance     float32 // Minimum distance from target. Default is 0.01
	MaxDistance     float32 // Maximum distance from target. Default is infinity
	MinPolarAngle   float32 // Minimum polar angle for rotatiom
	MaxPolarAngle   float32
	MinAzimuthAngle float32
	MaxAzimuthAngle float32
	KeyRotateSpeed  float32
	KeyPanSpeed     float32
	// Internal
	icam        camera.ICamera
	cam         *camera.Camera
	camPersp    *camera.Perspective
	camOrtho    *camera.Orthographic
	win         window.IWindow
	position0   math32.Vector3 // Initial camera position
	target0     math32.Vector3 // Initial camera target position
	state       int            // current active state
	phiDelta    float32        // rotation delta in the XZ plane
	thetaDelta  float32        // rotation delta in the YX plane
	rotateStart math32.Vector2
	rotateEnd   math32.Vector2
	rotateDelta math32.Vector2
	panStart    math32.Vector2 // initial pan screen coordinates
	panEnd      math32.Vector2 // final pan scren coordinates
	panDelta    math32.Vector2
	panOffset   math32.Vector2
	zoomStart   float32
	zoomEnd     float32
	zoomDelta   float32
	subsEvents  int // Address of this field is used as events subscription id
	subsPos     int // Address of this field is used as cursor pos events subscription id
}

const (
	stateNone = iota
	stateRotate
	stateZoom
	statePan
)

// Package logger
var log = logger.New("ORBIT", logger.Default)

// NewOrbitControl creates and returns a pointer to a new orbito control for
// the specified camera and window
func NewOrbitControl(icam camera.ICamera, win window.IWindow) *OrbitControl {

	oc := new(OrbitControl)
	oc.icam = icam
	oc.win = win

	oc.cam = icam.GetCamera()
	if persp, ok := icam.(*camera.Perspective); ok {
		oc.camPersp = persp
	} else if ortho, ok := icam.(*camera.Orthographic); ok {
		oc.camOrtho = ortho
	} else {
		panic("Invalid camera type")
	}

	// Set defaults
	oc.Enabled = true
	oc.EnableRotate = true
	oc.EnableZoom = true
	oc.EnablePan = true
	oc.EnableKeys = true
	oc.ZoomSpeed = 1.0
	oc.RotateSpeed = 1.0
	oc.MinDistance = 0.01
	oc.MaxDistance = float32(math.Inf(1))
	oc.MinPolarAngle = 0
	oc.MaxPolarAngle = math32.Pi
	oc.MinAzimuthAngle = float32(math.Inf(-1))
	oc.MaxAzimuthAngle = float32(math.Inf(1))
	oc.KeyPanSpeed = 5.0
	oc.KeyRotateSpeed = 0.02

	// Saves initial camera parameters
	oc.position0 = oc.cam.Position()
	oc.target0 = oc.cam.Target()

	// Subscribe to events
	oc.win.SubscribeID(window.OnMouseUp, &oc.subsEvents, oc.onMouse)
	oc.win.SubscribeID(window.OnMouseDown, &oc.subsEvents, oc.onMouse)
	oc.win.SubscribeID(window.OnScroll, &oc.subsEvents, oc.onScroll)
	oc.win.SubscribeID(window.OnKeyDown, &oc.subsEvents, oc.onKey)
	return oc
}

// Dispose unsubscribes from all events
func (oc *OrbitControl) Dispose() {

	// Unsubscribe to event handlers
	oc.win.UnsubscribeID(window.OnMouseUp, &oc.subsEvents)
	oc.win.UnsubscribeID(window.OnMouseDown, &oc.subsEvents)
	oc.win.UnsubscribeID(window.OnScroll, &oc.subsEvents)
	oc.win.UnsubscribeID(window.OnKeyDown, &oc.subsEvents)
	oc.win.UnsubscribeID(window.OnCursor, &oc.subsPos)
}

// Reset to initial camera position
func (oc *OrbitControl) Reset() {

	oc.state = stateNone
	oc.cam.SetPositionVec(&oc.position0)
	oc.cam.LookAt(&oc.target0)
}

// Pan the camera and target by the specified deltas
func (oc *OrbitControl) Pan(deltaX, deltaY float32) {

	width, height := oc.win.Size()
	oc.pan(deltaX, deltaY, width, height)
	oc.updatePan()
}

// Zoom in or out
func (oc *OrbitControl) Zoom(delta float32) {

	oc.zoomDelta = delta
	oc.updateZoom()
}

// RotateLeft rotates the camera left by specified angle
func (oc *OrbitControl) RotateLeft(angle float32) {

	oc.thetaDelta -= angle
	oc.updateRotate()
}

// RotateUp rotates the camera up by specified angle
func (oc *OrbitControl) RotateUp(angle float32) {

	oc.phiDelta -= angle
	oc.updateRotate()
}

// Updates the camera rotation from thetaDelta and phiDelta
func (oc *OrbitControl) updateRotate() {

	const EPS = 0.01

	// Get camera parameters
	position := oc.cam.Position()
	target := oc.cam.Target()
	up := oc.cam.Up()

	// Camera UP is the orbit axis
	var quat math32.Quaternion
	quat.SetFromUnitVectors(&up, &math32.Vector3{0, 1, 0})
	quatInverse := quat
	quatInverse.Inverse()

	// Calculates direction vector from camera position to target
	vdir := position
	vdir.Sub(&target)
	vdir.ApplyQuaternion(&quat)

	// Calculate angles from current camera position
	radius := vdir.Length()
	theta := math32.Atan2(vdir.X, vdir.Z)
	phi := math32.Acos(vdir.Y / radius)

	// Add deltas to the angles
	theta += oc.thetaDelta
	phi += oc.phiDelta

	// Restrict phi (elevation) to be between desired limits
	phi = math32.Max(oc.MinPolarAngle, math32.Min(oc.MaxPolarAngle, phi))
	phi = math32.Max(EPS, math32.Min(math32.Pi-EPS, phi))
	// Restrict theta to be between desired limits
	theta = math32.Max(oc.MinAzimuthAngle, math32.Min(oc.MaxAzimuthAngle, theta))

	// Calculate new cartesian coordinates
	vdir.X = radius * math32.Sin(phi) * math32.Sin(theta)
	vdir.Y = radius * math32.Cos(phi)
	vdir.Z = radius * math32.Sin(phi) * math32.Cos(theta)

	// Rotate offset back to "camera-up-vector-is-up" space
	vdir.ApplyQuaternion(&quatInverse)

	position = target
	position.Add(&vdir)
	oc.cam.SetPositionVec(&position)
	oc.cam.LookAt(&target)

	// Reset deltas
	oc.thetaDelta = 0
	oc.phiDelta = 0
}

// Updates camera rotation from tethaDelta and phiDelta
// ALTERNATIVE rotation algorithm
func (oc *OrbitControl) updateRotate2() {

	const EPS = 0.01

	// Get camera parameters
	position := oc.cam.Position()
	target := oc.cam.Target()
	up := oc.cam.Up()

	// Calculates direction vector from target to camera
	vdir := position
	vdir.Sub(&target)

	// Calculates right and up vectors
	var vright math32.Vector3
	vright.CrossVectors(&up, &vdir)
	vright.Normalize()
	var vup math32.Vector3
	vup.CrossVectors(&vdir, &vright)
	vup.Normalize()

	phi := vdir.AngleTo(&math32.Vector3{0, 1, 0})
	newphi := phi + oc.phiDelta
	if newphi < EPS || newphi > math32.Pi-EPS {
		oc.phiDelta = 0
	} else if newphi < oc.MinPolarAngle || newphi > oc.MaxPolarAngle {
		oc.phiDelta = 0
	}

	// Rotates position around the two vectors
	vdir.ApplyAxisAngle(&vup, oc.thetaDelta)
	vdir.ApplyAxisAngle(&vright, oc.phiDelta)

	// Adds target back get final position
	position = target
	position.Add(&vdir)
	log.Debug("orbit set position")
	oc.cam.SetPositionVec(&position)
	oc.cam.LookAt(&target)

	// Reset deltas
	oc.thetaDelta = 0
	oc.phiDelta = 0
}

// Updates camera pan from panOffset
func (oc *OrbitControl) updatePan() {

	// Get camera parameters
	position := oc.cam.Position()
	target := oc.cam.Target()
	up := oc.cam.Up()

	// Calculates direction vector from camera position to target
	vdir := target
	vdir.Sub(&position)
	vdir.Normalize()

	// Calculates vector perpendicular to direction and up (side vector)
	var vpanx math32.Vector3
	vpanx.CrossVectors(&up, &vdir)
	vpanx.Normalize()

	// Calculates vector perpendicular to direction and vpanx
	var vpany math32.Vector3
	vpany.CrossVectors(&vdir, &vpanx)
	vpany.Normalize()

	// Adds pan offsets
	vpanx.MultiplyScalar(oc.panOffset.X)
	vpany.MultiplyScalar(oc.panOffset.Y)
	var vpan math32.Vector3
	vpan.AddVectors(&vpanx, &vpany)

	// Adds offsets to camera position and target
	position.Add(&vpan)
	target.Add(&vpan)

	// Sets new camera parameters
	oc.cam.SetPositionVec(&position)
	oc.cam.LookAt(&target)

	// Reset deltas
	oc.panOffset.Set(0, 0)
}

// Updates camera zoom from zoomDelta
func (oc *OrbitControl) updateZoom() {

	if oc.camOrtho != nil {
		zoom := oc.camOrtho.Zoom() - 0.01*oc.zoomDelta
		oc.camOrtho.SetZoom(zoom)
		// Reset delta
		oc.zoomDelta = 0
		return
	}

	// Get camera and target positions
	position := oc.cam.Position()
	target := oc.cam.Target()

	// Calculates direction vector from target to camera position
	vdir := position
	vdir.Sub(&target)

	// Calculates new distance from target and applies limits
	dist := vdir.Length() * (1.0 + oc.zoomDelta*oc.ZoomSpeed/10.0)
	dist = math32.Max(oc.MinDistance, math32.Min(oc.MaxDistance, dist))
	vdir.SetLength(dist)

	// Adds new distance to target to get new camera position
	target.Add(&vdir)
	oc.cam.SetPositionVec(&target)

	// Reset delta
	oc.zoomDelta = 0
}

// Called when mouse button event is received
func (oc *OrbitControl) onMouse(evname string, ev interface{}) {

	// If control not enabled ignore event
	if !oc.Enabled {
		return
	}

	mev := ev.(*window.MouseEvent)
	// Mouse button pressed
	if mev.Action == window.Press {
		// Left button pressed sets Rotate state
		if mev.Button == window.MouseButtonLeft {
			if !oc.EnableRotate {
				return
			}
			oc.state = stateRotate
			oc.rotateStart.Set(float32(mev.Xpos), float32(mev.Ypos))
		} else
		// Middle button pressed sets Zoom state
		if mev.Button == window.MouseButtonMiddle {
			if !oc.EnableZoom {
				return
			}
			oc.state = stateZoom
			oc.zoomStart = float32(mev.Ypos)
		} else
		// Right button pressed sets Pan state
		if mev.Button == window.MouseButtonRight {
			if !oc.EnablePan {
				return
			}
			oc.state = statePan
			oc.panStart.Set(float32(mev.Xpos), float32(mev.Ypos))
		}
		// If a valid state is set requests mouse position events
		if oc.state != stateNone {
			oc.win.SubscribeID(window.OnCursor, &oc.subsPos, oc.onCursorPos)
		}
		return
	}

	// Mouse button released
	if mev.Action == window.Release {
		oc.win.UnsubscribeID(window.OnCursor, &oc.subsPos)
		oc.state = stateNone
	}
}

// Called when cursor position event is received
func (oc *OrbitControl) onCursorPos(evname string, ev interface{}) {

	// If control not enabled ignore event
	if !oc.Enabled {
		return
	}

	mev := ev.(*window.CursorEvent)
	// Rotation
	if oc.state == stateRotate {
		oc.rotateEnd.Set(float32(mev.Xpos), float32(mev.Ypos))
		oc.rotateDelta.SubVectors(&oc.rotateEnd, &oc.rotateStart)
		oc.rotateStart = oc.rotateEnd
		// rotating across whole screen goes 360 degrees around
		width, height := oc.win.Size()
		oc.RotateLeft(2 * math32.Pi * oc.rotateDelta.X / float32(width) * oc.RotateSpeed)
		// rotating up and down along whole screen attempts to go 360, but limited to 180
		oc.RotateUp(2 * math32.Pi * oc.rotateDelta.Y / float32(height) * oc.RotateSpeed)
		return
	}

	// Panning
	if oc.state == statePan {
		oc.panEnd.Set(float32(mev.Xpos), float32(mev.Ypos))
		oc.panDelta.SubVectors(&oc.panEnd, &oc.panStart)
		oc.panStart = oc.panEnd
		oc.Pan(oc.panDelta.X, oc.panDelta.Y)
		return
	}

	// Zooming
	if oc.state == stateZoom {
		oc.zoomEnd = float32(mev.Ypos)
		oc.zoomDelta = oc.zoomEnd - oc.zoomStart
		oc.zoomStart = oc.zoomEnd
		oc.Zoom(oc.zoomDelta)
	}
}

// Called when mouse button scroll event is received
func (oc *OrbitControl) onScroll(evname string, ev interface{}) {

	if !oc.Enabled || !oc.EnableZoom || oc.state != stateNone {
		return
	}
	sev := ev.(*window.ScrollEvent)
	oc.Zoom(float32(-sev.Yoffset))
}

// Called when key is pressed, released or repeats.
func (oc *OrbitControl) onKey(evname string, ev interface{}) {

	if !oc.Enabled || !oc.EnableKeys {
		return
	}

	kev := ev.(*window.KeyEvent)
	if kev.Action == window.Release {
		return
	}

	if oc.EnablePan && kev.Mods == 0 {
		switch kev.Keycode {
		case window.KeyUp:
			oc.Pan(0, oc.KeyPanSpeed)
		case window.KeyDown:
			oc.Pan(0, -oc.KeyPanSpeed)
		case window.KeyLeft:
			oc.Pan(oc.KeyPanSpeed, 0)
		case window.KeyRight:
			oc.Pan(-oc.KeyPanSpeed, 0)
		}
	}

	if oc.EnableRotate && kev.Mods == window.ModShift {
		switch kev.Keycode {
		case window.KeyUp:
			oc.RotateUp(oc.KeyRotateSpeed)
		case window.KeyDown:
			oc.RotateUp(-oc.KeyRotateSpeed)
		case window.KeyLeft:
			oc.RotateLeft(-oc.KeyRotateSpeed)
		case window.KeyRight:
			oc.RotateLeft(oc.KeyRotateSpeed)
		}
	}

	if oc.EnableZoom && kev.Mods == window.ModControl {
		switch kev.Keycode {
		case window.KeyUp:
			oc.Zoom(-1.0)
		case window.KeyDown:
			oc.Zoom(1.0)
		}
	}
}

func (oc *OrbitControl) pan(deltaX, deltaY float32, swidth, sheight int) {

	// Perspective camera
	if oc.camPersp != nil {
		position := oc.cam.Position()
		target := oc.cam.Target()
		offset := position.Clone().Sub(&target)
		targetDistance := offset.Length()
		// Half the FOV is center to top of screen
		targetDistance += math32.Tan((oc.camPersp.Fov() / 2.0) * math32.Pi / 180.0)
		// we actually don't use screenWidth, since perspective camera is fixed to screen height
		oc.panLeft(2 * deltaX * targetDistance / float32(sheight))
		oc.panUp(2 * deltaY * targetDistance / float32(sheight))
		return
	}
	// Orthographic camera
	left, right, top, bottom, _, _ := oc.camOrtho.Planes()
	oc.panLeft(deltaX * (right - left) / float32(swidth))
	oc.panUp(deltaY * (top - bottom) / float32(sheight))
}

func (oc *OrbitControl) panLeft(distance float32) {

	oc.panOffset.X += distance
}

func (oc *OrbitControl) panUp(distance float32) {

	oc.panOffset.Y += distance
}
