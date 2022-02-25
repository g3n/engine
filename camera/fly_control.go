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

// FlyMovement are the ways the camera can move.
type FlyMovement int

const (
	Forward   FlyMovement = iota // positive translation in the direction the camera is looking.
	Backward                     // negative translation in the direction the camera is looking.
	Right                        // right ("positive") translation from the camera's point of view.
	Left                         // left ("negative") translation from the camera's point of view.
	Up                           // upward ("positive") translation from the camera's point of view.
	Down                         // downward ("negative") translation from the camera's point of view.
	YawRight                     // "right" rotation from the camera's point of view.
	YawLeft                      // "left" rotation from the camera's point of view.
	PitchUp                      // "up" rotation from the camera's point of view.
	PitchDown                    // "down" rotation from the camera's point of view.
	RollRight                    // rotation "right" (CW) along the direction the camera is looking.
	RollLeft                     // rotation "left" (CCW) along the direction the camera is looking.
	ZoomOut                      // increase in the field of view.
	ZoomIn                       // decrease in the field of view.
)

// MouseMotion are the ways a mouse can move.
type MouseMotion int

const (
	MouseNone        MouseMotion = iota
	MouseRight                   // +X
	MouseLeft                    // -X
	MouseDown                    // +Y
	MouseUp                      // -Y
	MouseScrollRight             // +X
	MouseScrollLeft              // -X
	MouseScrollDown              // +Y
	MouseScrollUp                // -Y
)

// MouseGuesture is an action performed with the mouse in a particular
// direction (eg up) and with optional mouse buttons pressed.
// Captured indicates if the MouseGuesture is active/valid only when the mouse is
// captured by the window (indicated by calling SetMouseIsCaptured(true)).
type MouseGuesture struct {
	Motion   MouseMotion          // motion of the guesture
	Buttons  []window.MouseButton // list of mouse buttons that must be pressed
	Captured bool                 // guesture happens when mouse is captured by window
	// Keys     []window.Key      // future use to specify keys (eg Ctrl+MouseUp)
}

const degrees = math32.Pi / 180.0
const radians = 1.0 / degrees
const twoPiOver64 = (2.0 * math32.Pi) / 64.0 // 5.625 deg

// FlyControl provides a camera controller that allows looking at the scene from a first
// person style view. It allows for translation in the Forward/Backward, Right/Left, and Up/Down
// directions, and rotation in the Yaw/Pitch/Roll directions from the camera's point
// of view. The camera's field of view can also be modified to allow ZoomOut/ZoomIn.
//
// For maximum flexibility, the FlyControl is very configurable, but two main modes
// of operation exist: "manual" and "automatic".
//
// "Manual" mode requires the user to directly make adjustments to the FlyControl's
// position and orientation by using the Forward(), Right(), Yaw(), Pitch() and other
// similar methods. Manual mode does not subscribe to key or mouse events, nor does it
// use the FlyControl's member maps "Speeds", "Keys", or "Mouse" ("Constraints" is
// used). A FlyControl created (eg with NewFlyControl()) with no options defaults
// to "manual" mode.
//
// "Automatic" mode subscribes to key and mouse events, and handles position and orientation
// adjustments according to parameters specified in the "Keys", "Mouse", "Speed", and "Constraint"
// maps. Generally, methods such as Forward() will not be used. Some preconfigured options
// are available with FPSStyle() and FlightSimStyle(), and can be further configured by
// the With*() options, and by modifying the FlyControl's Keys (etc) maps directly.
type FlyControl struct {
	core.Dispatcher         // Embedded event dispatcher
	cam             *Camera // Controlled camera

	position   math32.Vector3 // camera position in the world
	rotation   math32.Vector3 // camera yaw, pitch, and roll (radians)
	forward    math32.Vector3 // forward direction of camera
	up         math32.Vector3 // up direction of camera
	upWorld    math32.Vector3 // The world up direction
	useUpWorld bool           // true to use upWorld to do Up/Yaw/Pitch/Roll, false to use (camera) up

	mouseIsCaptured  bool           // indicates if the cursor is captured by the GL window
	mousePrevPos     math32.Vector2 // used to determine mouse cursor movements
	mouseBtnState    uint32         // bitfield. btn0 == 1, btn1 == 2, btn3 == 4 etc
	mouseSensitivity float32        // modifies sensitivity of mouse cursor position changes, in [0,1]

	isKeySubscribed   bool // indicates status of keyboard event subscription
	isMouseSubscribed bool // indicates status of mouse event subscription

	// Constraints map.
	// Specifies the maxium cumulative value in a particular movement.
	// Translation movements aren't used.
	// Angular rotation contraints are in radians.
	Constraints map[FlyMovement]float32

	// Speeds map.
	// Translation movements in units/event.
	// Angular rotation and zoom in radians/event.
	Speeds map[FlyMovement]float32

	// Keys map.
	// Maps a movment to a key press.
	Keys map[FlyMovement]window.Key

	// Mouse map.
	// Maps a movement to a mouse "guesture"
	// (a combination of mouse movement and optional mouse button).
	Mouse map[FlyMovement]MouseGuesture
}

// FlyControlOption is a functional option when creating a new FlyControl.
type FlyControlOption func(fc *FlyControl)

// WithSpeeds provides a map of movement speeds for the FlyControl.
func WithSpeeds(speeds map[FlyMovement]float32) FlyControlOption {
	return func(fc *FlyControl) {
		fc.Speeds = speeds
	}
}

// WithKeys provides a map of movements to keys for the FlyControl.
func WithKeys(keys map[FlyMovement]window.Key) FlyControlOption {
	return func(fc *FlyControl) {
		fc.Keys = keys
		fc.Subscribe(true, false)
	}
}

// WithMouse provides a map of movements to mouse guestures for the FlyControl.
func WithMouse(mouse map[FlyMovement]MouseGuesture) FlyControlOption {
	return func(fc *FlyControl) {
		fc.Mouse = mouse
		fc.Subscribe(false, true)
	}
}

// WithConstraints provides a map of rotation (and FoV) constraints for the FlyControl.
func WithConstraints(constraints map[FlyMovement]float32) FlyControlOption {
	return func(fc *FlyControl) {
		fc.Constraints = constraints
	}
}

// FPSStyle configures the FlyControl with some default settings to make it act
// similarly to most FPS style games.
//
// Yaw causes the camera to rotate around world up instead of camera up (IsUsingWorldUp()==true).
// Forward/Backward/Right/Left is mapped to the WADS keys while Yaw and Pitch are mapped to the
// arrow keys and also the mouse. Mouse look requires the user to left-click and hold. Pitch is
// constrained between ±85 degrees, and FoV is constrained between 0 and 100 degrees.
// The control subscribes to keyboard and mouse events.
//
// These settings can be altered by providing maps with WithKeys() etc functional
// options in NewFlyControl(), or by altering the maps directly in the FlyControl
// struct. This option should be listed before any of the With*() options.
func FPSStyle() FlyControlOption {
	return func(fc *FlyControl) {
		fc.UseWorldUp(true)

		// WADS style keys
		WithKeys(map[FlyMovement]window.Key{
			Forward:  window.KeyW,
			Backward: window.KeyS,
			Right:    window.KeyD,
			Left:     window.KeyA,
			// Up:        window.KeyPageUp,
			// Down:      window.KeyPageDown,
			YawRight:  window.KeyRight,
			YawLeft:   window.KeyLeft,
			PitchUp:   window.KeyUp,
			PitchDown: window.KeyDown,
			// RollRight: window.KeyE,
			// RollLeft:  window.KeyQ,
			ZoomOut: window.KeyMinus,
			ZoomIn:  window.KeyEqual,
		})(fc)

		// click+hold to look with mouse
		WithMouse(map[FlyMovement]MouseGuesture{
			YawRight:  {Motion: MouseRight, Buttons: []window.MouseButton{window.MouseButtonLeft}},
			YawLeft:   {Motion: MouseLeft, Buttons: []window.MouseButton{window.MouseButtonLeft}},
			PitchUp:   {Motion: MouseUp, Buttons: []window.MouseButton{window.MouseButtonLeft}},
			PitchDown: {Motion: MouseDown, Buttons: []window.MouseButton{window.MouseButtonLeft}},
			ZoomIn:    {Motion: MouseScrollUp},
			ZoomOut:   {Motion: MouseScrollDown},
		})(fc)

		WithSpeeds(map[FlyMovement]float32{
			Forward:  1,
			Backward: -1,
			Right:    1,
			Left:     -1,
			// Up:        1,
			// Down:      -1,
			YawRight:  twoPiOver64,
			YawLeft:   -twoPiOver64,
			PitchUp:   twoPiOver64,
			PitchDown: -twoPiOver64,
			// RollRight: twoPiOver64,
			// RollLeft:  -twoPiOver64,
			ZoomOut: twoPiOver64,
			ZoomIn:  -twoPiOver64,
		})(fc)

		// constrain pitch due to using world up
		WithConstraints(map[FlyMovement]float32{
			// YawRight:  45 * degrees,
			// YawLeft:   -45 * degrees,
			PitchUp:   85 * degrees,
			PitchDown: -85 * degrees,
			// RollRight: 45 * degrees,
			// RollLeft:  -45 * degrees,
			ZoomOut: 100.0 * degrees,
			ZoomIn:  1.0 * degrees,
		})(fc)
	}
}

// FlightSimStyle configures the FlyControl with some default settings to make it act
// similarly to a airplane or spacecraft.
//
// Yaw causes the camera to rotate around camera up instead of world up (IsUsingWorldUp()==false).
// Forward/Backward/Right/Left is mapped to the arrow keys while Yaw/Pitch/Roll are mapped to the
// WADS and also the mouse. Mouse look requires the user to left-click and hold. There are no rotational
// constraints, allowing the camera to perform manuvers such as loops and aileron rolls.
// FoV is constrained between 0 and 100 degrees. The control subscribes to keyboard and mouse events.
//
// These settings can be altered by providing maps with WithKeys() etc functional
// options in NewFlyControl(), or by altering the maps directly in the FlyControl
// struct. This option should be listed before any of the With*() options.
func FlightSimStyle() FlyControlOption {
	return func(fc *FlyControl) {
		fc.UseWorldUp(false)

		// arrow keys for translation, WADS keys for rotation
		WithKeys(map[FlyMovement]window.Key{
			Forward:   window.KeyUp,
			Backward:  window.KeyDown,
			Right:     window.KeyRight,
			Left:      window.KeyLeft,
			Up:        window.KeyPageUp,
			Down:      window.KeyPageDown,
			YawRight:  window.KeyD,
			YawLeft:   window.KeyA,
			PitchUp:   window.KeyW,
			PitchDown: window.KeyS,
			RollRight: window.KeyE,
			RollLeft:  window.KeyQ,
			ZoomOut:   window.KeyMinus,
			ZoomIn:    window.KeyEqual,
		})(fc)

		// click+hold to look with mouse
		WithMouse(map[FlyMovement]MouseGuesture{
			YawRight:  {Motion: MouseRight, Buttons: []window.MouseButton{window.MouseButtonLeft}},
			YawLeft:   {Motion: MouseLeft, Buttons: []window.MouseButton{window.MouseButtonLeft}},
			PitchUp:   {Motion: MouseUp, Buttons: []window.MouseButton{window.MouseButtonLeft}},
			PitchDown: {Motion: MouseDown, Buttons: []window.MouseButton{window.MouseButtonLeft}},
			ZoomIn:    {Motion: MouseScrollUp},
			ZoomOut:   {Motion: MouseScrollDown},
		})(fc)

		WithSpeeds(map[FlyMovement]float32{
			Forward:   1,
			Backward:  -1,
			Right:     1,
			Left:      -1,
			Up:        1,
			Down:      -1,
			YawRight:  twoPiOver64,
			YawLeft:   -twoPiOver64,
			PitchUp:   twoPiOver64,
			PitchDown: -twoPiOver64,
			RollRight: twoPiOver64,
			RollLeft:  -twoPiOver64,
			ZoomOut:   twoPiOver64,
			ZoomIn:    -twoPiOver64,
		})(fc)

		// no need to constrain rotation for flight sim style
		WithConstraints(map[FlyMovement]float32{
			// YawRight:  45 * degrees,
			// YawLeft:   -45 * degrees,
			// PitchUp:   85 * degrees,
			// PitchDown: -85 * degrees,
			// RollRight: 45 * degrees,
			// RollLeft:  -45 * degrees,
			ZoomOut: 100.0 * degrees,
			ZoomIn:  1.0 * degrees,
		})(fc)
	}
}

// NewFlyControl initalizes a FlyControl to manipulate cam. It starts positioned
// at the cam's position, and oriented looking at target with camera's up aligned
// in the direction of worldUp. Configuration options can be provided. The FlyControl
// will default to "manual" mode if no options are given. See documentation for
// FlyControl for more information about common usage patterns.
func NewFlyControl(cam *Camera, target, worldUp *math32.Vector3,
	options ...FlyControlOption) *FlyControl {

	fc := new(FlyControl)

	// init camera position and orientation
	fc.cam = cam
	fc.position = cam.Position()
	fc.Reorient(target, worldUp)

	// init mouse state
	fc.resetMousePrevPosition()
	fc.mouseIsCaptured = false
	fc.mouseBtnState = 0
	fc.mouseSensitivity = 0.5

	// apply options
	for _, option := range options {
		option(fc)
	}

	return fc
}

// Subscribe to input events. A value of false for either key or mouse
// does not unsubscribe from that type of event. Multiple calls will subscribe
// only once.
func (fc *FlyControl) Subscribe(key, mouse bool) {
	if key && !fc.isKeySubscribed {
		gui.Manager().SubscribeID(window.OnKeyDown, fc, fc.onKey)
		gui.Manager().SubscribeID(window.OnKeyRepeat, fc, fc.onKey)
		// gui.Manager().SubscribeID(window.OnKeyUp, fc, fc.onKeyUp)
		fc.isKeySubscribed = true
	}

	if mouse && !fc.isMouseSubscribed {
		gui.Manager().SubscribeID(window.OnMouseUp, fc, fc.onMouse)
		gui.Manager().SubscribeID(window.OnMouseDown, fc, fc.onMouse)
		gui.Manager().SubscribeID(window.OnScroll, fc, fc.onScroll)
		gui.Manager().SubscribeID(window.OnCursor, fc, fc.onCursor)
		fc.isMouseSubscribed = true
	}
}

// Unsubscribe from input events.
func (fc *FlyControl) Unsubscribe(key, mouse bool) {
	if key && fc.isKeySubscribed {
		gui.Manager().UnsubscribeID(window.OnKeyDown, fc)
		gui.Manager().UnsubscribeID(window.OnKeyRepeat, fc)
		// gui.Manager().UnsubscribeID(window.OnKeyUp, fc)
		fc.isKeySubscribed = false
	}

	if mouse && fc.isMouseSubscribed {
		gui.Manager().UnsubscribeID(window.OnMouseUp, fc)
		gui.Manager().UnsubscribeID(window.OnMouseDown, fc)
		gui.Manager().UnsubscribeID(window.OnScroll, fc)
		gui.Manager().UnsubscribeID(window.OnCursor, fc)
		fc.isMouseSubscribed = false
	}
}

// Dispose unsubscribes from all events.
func (fc *FlyControl) Dispose() {
	fc.Unsubscribe(true, true)
}

// Reposition the camera to the new position.
func (fc *FlyControl) Reposition(position *math32.Vector3) {
	fc.position.Copy(position)
	fc.cam.SetPositionVec(position)
}

// Reorient the camera to look at target and use the given world up direction.
func (fc *FlyControl) Reorient(target, worldUp *math32.Vector3) {
	fc.rotation.Set(0, 0, 0) // reset the total rotation
	fc.upWorld.Copy(worldUp) // worldUp might have changed
	fc.forward.Copy(target.Clone().Sub(&fc.position).Normalize())
	right := fc.forward.Clone().Cross(worldUp).Normalize()
	fc.up.Copy(right.Clone().Cross(&fc.forward).Normalize())
}

// "getters"

// GetPosition returns a copy of the camera's position.
func (fc *FlyControl) GetPosition() (position math32.Vector3) {
	return fc.position
}

// GetRotation returns a copy of the camera's cumulative rotation. Yaw, Pitch, and Roll
// are in the X, Y, and Z coordinates.
func (fc *FlyControl) GetRotation() (rotation math32.Vector3) {
	return fc.rotation
}

// GetDirections returns copies of the camera's current "forward" and "up" direction.
// The up direction is from the camera's perspective regardless of the status
// of IsUsingWorldUp().
func (fc *FlyControl) GetDirections() (forward, up math32.Vector3) {
	return fc.forward, fc.up
}

// movement changes

// Forward and backward translation along the camera's forward axis.
func (fc *FlyControl) Forward(delta float32) {
	deltaDir := fc.forward.Clone().MultiplyScalar(delta)
	fc.position.Add(deltaDir)
	fc.cam.SetPositionVec(&fc.position)
}

// Right and left translation along the camera's right axis.
func (fc *FlyControl) Right(delta float32) {
	up := fc.whichUp()
	// right direction from forward cross up
	deltaDir := fc.forward.Clone().Cross(up).MultiplyScalar(delta)
	fc.position.Add(deltaDir)
	fc.cam.SetPositionVec(&fc.position)
}

// Up and down translation along the current "up" axis.
func (fc *FlyControl) Up(delta float32) {
	up := fc.whichUp()
	deltaDir := up.Clone().MultiplyScalar(delta)
	fc.position.Add(deltaDir)
	fc.cam.SetPositionVec(&fc.position)
}

// Yaw adjustment in radians. Yaw rotates the camera about the current "up" axis.
// Positive yaw is "right" from the camera's point of view.
func (fc *FlyControl) Yaw(delta float32) {
	yaw := fc.rotation.X + delta
	if fc.constraintOk(yaw, YawLeft, YawRight) {
		fc.rotation.X = yaw
		// rotation about up axis
		// affects only "forward" direction
		// because of the right hand coord system, positive rotation about up-axis
		// makes camera appear to yaw left instead of right. thus, delta
		// must be inverted below.
		up := fc.whichUp()
		fc.forward.ApplyAxisAngle(up, -delta)
		fc.cam.LookAt(fc.position.Clone().Add(&fc.forward), up)
	}
}

// Pitch adjustment in radians. Pitch rotates the camera about its right axis.
//
// Caution: If using world up, take care not to allow the camera's forward
// direction to become parallel to the world up direction by setting constraints
// on PitchUp and PitchDown to be in the interval (-π/2, π/2). If camera forward
// and world up become parallel, NaNs will happen.
func (fc *FlyControl) Pitch(delta float32) {
	pitch := fc.rotation.Y + delta
	if fc.constraintOk(pitch, PitchDown, PitchUp) {
		fc.rotation.Y = pitch
		// rotation about right axis
		// affects both "forward" and "up" directions
		up := fc.whichUp()
		right := fc.forward.Clone().Cross(up)
		fc.forward.ApplyAxisAngle(right, delta)
		fc.up.ApplyAxisAngle(right, delta)
		fc.cam.LookAt(fc.position.Clone().Add(&fc.forward), up)
	}
}

// Roll adjustment in radians. Roll rotates the camera about its forward axis.
func (fc *FlyControl) Roll(delta float32) {
	roll := fc.rotation.Z + delta
	if fc.constraintOk(roll, RollLeft, RollRight) {
		fc.rotation.Z = roll
		// rotation about forward axis
		// affects only "up" direction
		up := fc.whichUp()
		fc.up.ApplyAxisAngle(&fc.forward, delta)
		fc.cam.LookAt(fc.position.Clone().Add(&fc.forward), up)
	}
}

// Field of view changes

// Zoom modifies the fov based on the delta change in radians.
func (fc *FlyControl) Zoom(delta float32) {
	newFovRad := fc.cam.Fov()*degrees + delta
	if fc.constraintOk(newFovRad, ZoomIn, ZoomOut) {
		fc.cam.SetFov(newFovRad * radians)
	}
}

// ScaleZoom modifies the fov based on scaling the current fov.
func (fc *FlyControl) ScaleZoom(fovScale float32) {
	newFovRad := (fc.cam.Fov() * degrees) * fovScale
	if fc.constraintOk(newFovRad, ZoomIn, ZoomOut) {
		fc.cam.SetFov(newFovRad * radians)
	}
}

// constraintOk checks if newValue is allowed according to the given low and high
// movement constraints.
func (fc *FlyControl) constraintOk(newValue float32, low, high FlyMovement) bool {
	min, minDefined := fc.Constraints[low]
	max, maxDefined := fc.Constraints[high]
	if minDefined && newValue < min {
		return false
	}
	if maxDefined && newValue > max {
		return false
	}
	return true
}

// whichUp returns the currently used up (world or camera).
func (fc *FlyControl) whichUp() *math32.Vector3 {
	if fc.useUpWorld {
		return &fc.upWorld
	}
	return &fc.up
}

// UseWorldUp sets the FlyControl to use world up instead of camera up
// when calculating movements Up/Down, Yaw, Pitch, and Roll. Setting this to
// false will use the "up" direction relative to the camera's point of view.
//
// Caution: If using world up, take care not to allow the camera's forward
// direction to become parallel to the world up direction by setting constraints
// on PitchUp and PitchDown to be in the interval (-π/2, π/2). If camera forward
// and world up become parallel, NaNs will happen.
func (fc *FlyControl) UseWorldUp(use bool) {
	changed := fc.useUpWorld != use
	fc.useUpWorld = use
	if changed {
		fc.Reorient(fc.position.Clone().Add(&fc.forward), &fc.upWorld)
	}
}

// IsUsingWorldUp returns true if the "up" direction used for movement and rotation
// is the world up, and false if it is the camera up.
func (fc *FlyControl) IsUsingWorldUp() bool { return fc.useUpWorld }

// apply maps FlyMovements to the appropriate method.
func (fc *FlyControl) apply(movement FlyMovement, delta float32) {
	switch movement {
	case Backward, Forward:
		fc.Forward(delta)
	case Left, Right:
		fc.Right(delta)
	case Up, Down:
		fc.Up(delta)
	case YawLeft, YawRight:
		fc.Yaw(delta)
	case PitchUp, PitchDown:
		fc.Pitch(delta)
	case RollLeft, RollRight:
		fc.Roll(delta)
	case ZoomIn, ZoomOut:
		fc.Zoom(delta)
	}
}

// clamp and smoothly transition in [0,1]
func smoothstep(x float32) float32 {
	// https://en.wikipedia.org/wiki/Smoothstep
	if x <= 0 {
		return 0
	}
	if x >= 1 {
		return 1
	}
	return 3*x*x - 2*x*x*x
}

// GetMouseSensitivity gets the current mouse sensitivity in [0,1].
func (fc *FlyControl) GetMouseSensitivity() float32 { return fc.mouseSensitivity }

// SetMouseSensitivity to value in [0,1]. Values outside of [0,1] are clamped
// to that range.
func (fc *FlyControl) SetMouseSensitivity(value float32) {
	fc.mouseSensitivity = smoothstep(value)
}

// SetMouseIsCaptured is used to inform the FlyControl about the cursor
// being captured by the GL window so the FlyControl can respond to mouse events
// appropriately. Mouse/cursor capture is not modifed by the FlyControl.
func (fc *FlyControl) SetMouseIsCaptured(captured bool) {
	fc.mouseIsCaptured = captured
	fc.resetMousePrevPosition()
}

// sets/unsets the appropriate bit for button.
func (fc *FlyControl) toggleMouseButton(button window.MouseButton) {
	fc.mouseBtnState ^= 1 << button // toggle button state bit
}

// determines if all buttons are pressed.
func (fc *FlyControl) mouseButtonsPressed(buttons ...window.MouseButton) bool {
	if fc.mouseBtnState == 0 && len(buttons) == 0 {
		return true
	}
	var bits uint32
	for _, b := range buttons {
		bits ^= uint32(1 << b)
	}
	return fc.mouseBtnState == bits
}

func (fc *FlyControl) resetMousePrevPosition() {
	// uses NaN to indicate "unset" since numbers are valid values
	fc.mousePrevPos.X = math32.NaN()
	fc.mousePrevPos.Y = math32.NaN()
}

func (fc *FlyControl) isMousePrevPositionUnset() bool {
	return math32.IsNaN(fc.mousePrevPos.X)
}

// event listeners

// onMouse is called when an OnMouseDown/OnMouseUp event is received.
func (fc *FlyControl) onMouse(evname string, ev interface{}) {
	mev := ev.(*window.MouseEvent)
	fc.toggleMouseButton(mev.Button)
}

// onCursor is called when an OnCursor event is received.
func (fc *FlyControl) onCursor(evname string, ev interface{}) {

	mev := ev.(*window.CursorEvent)
	if fc.isMousePrevPositionUnset() {
		fc.mousePrevPos.X = mev.Xpos
		fc.mousePrevPos.Y = mev.Ypos
	}

	const moderator = 0.25 // 'magic' number to help mouse sensitivity be realistic in [0,1]
	dx := (mev.Xpos - fc.mousePrevPos.X) * moderator * fc.mouseSensitivity
	dy := (mev.Ypos - fc.mousePrevPos.Y) * moderator * fc.mouseSensitivity
	fc.mousePrevPos.X = mev.Xpos
	fc.mousePrevPos.Y = mev.Ypos

	// determine which MouseMotion happened
	mouseX := MouseNone
	mouseY := MouseNone
	if dx != 0 {
		if dx > 0 {
			mouseX = MouseRight
		} else {
			mouseX = MouseLeft
		}
	}
	if dy != 0 {
		if dy > 0 {
			mouseY = MouseDown
		} else {
			mouseY = MouseUp
		}
	}

	// find then find and apply the appropriate MouseGuesture
	for m, g := range fc.Mouse {
		pressed := fc.mouseButtonsPressed(g.Buttons...)
		captured := g.Captured == fc.mouseIsCaptured
		if pressed && captured {
			if g.Motion == mouseX {
				speed := fc.Speeds[m]
				fc.apply(m, math32.Abs(dx)*speed)
			}
			if g.Motion == mouseY {
				speed := fc.Speeds[m]
				fc.apply(m, math32.Abs(dy)*speed)
			}
		}
	}
}

// onScroll is called when an OnScroll event is received.
func (fc *FlyControl) onScroll(evname string, ev interface{}) {
	// x/y offset appears to always be +-1
	sev := ev.(*window.ScrollEvent)

	// determine which MouseMotion happened
	scrollX := MouseNone
	scrollY := MouseNone
	if sev.Xoffset != 0 {
		if sev.Xoffset > 0 {
			scrollX = MouseScrollRight
		} else {
			scrollX = MouseScrollLeft
		}
	}
	if sev.Yoffset != 0 {
		if sev.Yoffset > 0 {
			scrollY = MouseScrollUp
		} else {
			scrollY = MouseScrollDown
		}
	}

	// find then find and apply the appropriate MouseGuesture
	for m, g := range fc.Mouse {
		pressed := fc.mouseButtonsPressed(g.Buttons...)
		captured := g.Captured == fc.mouseIsCaptured
		if pressed && captured {
			if g.Motion == scrollX {
				speed := fc.Speeds[m]
				fc.apply(m, math32.Abs(sev.Xoffset)*speed)
			}
			if g.Motion == scrollY {
				speed := fc.Speeds[m]
				fc.apply(m, math32.Abs(sev.Yoffset)*speed)
			}
		}
	}
}

// onKey is called when an OnKeyDown/OnKeyRepeat event is received.
func (fc *FlyControl) onKey(evname string, ev interface{}) {

	kev := ev.(*window.KeyEvent)

	// find which movement the key corresponds to
	var movement FlyMovement = -1
	for m, k := range fc.Keys {
		if k == kev.Key {
			movement = m
			break
		}
	}
	if movement == -1 {
		// the pressed key is not mapped to a camera movement
		return
	}

	delta := fc.Speeds[movement]
	fc.apply(movement, delta)
}

// onKeyUp is called when an OnKeyUp event is received.
// func (fc *FlyControl) onKeyUp(evname string, ev interface{}) {
// }
