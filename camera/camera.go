// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package camera contain common camera types used for rendering 3D scenes.
package camera

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

// ICamera is interface for all camera types
type ICamera interface {
	GetCamera() *Camera
	ViewMatrix(*math32.Matrix4)
	ProjMatrix(*math32.Matrix4)
	Project(*math32.Vector3) *math32.Vector3
	Unproject(*math32.Vector3) *math32.Vector3
	SetRaycaster(rc *core.Raycaster, x, y float32)
}

// Camera is the base camera which is normally embedded in other camera types
type Camera struct {
	core.Node                 // Embedded Node
	target     math32.Vector3 // camera target in world coordinates
	up         math32.Vector3 // camera Up vector
	viewMatrix math32.Matrix4 // last calculated view matrix
}

// Initialize initializes the base camera.
// It is used by other camera types which embed this base camera
func (cam *Camera) Initialize() {

	cam.Node.Init()
	cam.target.Set(0, 0, 0)
	cam.up.Set(0, 1, 0)
	cam.SetDirection(0, 0, -1)
	cam.updateQuaternion()
}

// WorldDirection updates the specified vector with the
// current world direction the camera is pointed
func (cam *Camera) WorldDirection(result *math32.Vector3) {

	var wpos math32.Vector3
	cam.WorldPosition(&wpos)
	*result = cam.target
	result.Sub(&wpos).Normalize()
}

// LookAt sets the camera target position
func (cam *Camera) LookAt(target *math32.Vector3) {

	cam.target = *target
	cam.updateQuaternion()
}

// GetCamera satisfies the ICamera interface
func (cam *Camera) GetCamera() *Camera {

	return cam
}

// Target get the current target position
func (cam *Camera) Target() math32.Vector3 {

	return cam.target
}

// Up get the current up position
func (cam *Camera) Up() math32.Vector3 {

	return cam.up
}

// SetUp sets the camera up vector
func (cam *Camera) SetUp(up *math32.Vector3) {

	cam.up = *up
}

// SetPosition sets this camera world position
// This method overrides the Node method to update
// the camera quaternion, because changing the camera position
// may change its rotation
func (cam *Camera) SetPosition(x, y, z float32) {

	cam.Node.SetPosition(x, y, z)
	cam.updateQuaternion()
}

// SetPositionX sets this camera world position
// This method overrides the Node method to update
// the camera quaternion, because changing the camera position
// may change its rotation
func (cam *Camera) SetPositionX(x float32) {
	cam.Node.SetPositionX(x)
	cam.updateQuaternion()
}

// SetPositionY sets this camera world position
// This method overrides the Node method to update
// the camera quaternion, because changing the camera position
// may change its rotation
func (cam *Camera) SetPositionY(y float32) {
	cam.Node.SetPositionY(y)
	cam.updateQuaternion()
}

// SetPositionZ sets this camera world position
// This method overrides the Node method to update
// the camera quaternion, because changing the camera position
// may change its rotation
func (cam *Camera) SetPositionZ(z float32) {
	cam.Node.SetPositionZ(z)
	cam.updateQuaternion()
}

// SetPositionVec sets this node position from the specified vector pointer
// This method overrides the Node method to update
// the camera quaternion, because changing the camera position
// may change its rotation
func (cam *Camera) SetPositionVec(vpos *math32.Vector3) {

	cam.Node.SetPositionVec(vpos)
	cam.updateQuaternion()
}

// ViewMatrix returns the current view matrix of this camera
func (cam *Camera) ViewMatrix(m *math32.Matrix4) {

	var wpos math32.Vector3
	cam.WorldPosition(&wpos)
	cam.viewMatrix.LookAt(&wpos, &cam.target, &cam.up)
	*m = cam.viewMatrix
}

// updateQuaternion must be called when the camera position or target
// is changed to update its quaternion.
// This is important if the camera has children, such as an audio listener
func (cam *Camera) updateQuaternion() {

	var wdir math32.Vector3
	cam.WorldDirection(&wdir)
	var q math32.Quaternion
	q.SetFromUnitVectors(&math32.Vector3{0, 0, -1}, &wdir)
	cam.SetQuaternionQuat(&q)
}

// Project satisfies the ICamera interface and must
// be implemented for specific camera types.
func (cam *Camera) Project(v *math32.Vector3) *math32.Vector3 {

	panic("Not implemented")
}

// Unproject satisfies the ICamera interface and must
// be implemented for specific camera types.
func (cam *Camera) Unproject(v *math32.Vector3) *math32.Vector3 {

	panic("Not implemented")
}

// SetRaycaster satisfies the ICamera interface and must
// be implemented for specific camera types.
func (cam *Camera) SetRaycaster(rc *core.Raycaster, x, y float32) {

	panic("Not implemented")
}
