// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package camera contains common camera types used for rendering 3D scenes.
package camera

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

// ICamera is interface for all camera types.
type ICamera interface {
	GetCamera() *Camera
	SetAspect(float32)
	ViewMatrix(*math32.Matrix4)
	ProjMatrix(*math32.Matrix4)
	Project(*math32.Vector3) (*math32.Vector3, error)
	Unproject(*math32.Vector3) (*math32.Vector3, error)
	SetRaycaster(rc *core.Raycaster, x, y float32) error
}

// Camera is the base camera which is normally embedded in other camera types.
type Camera struct {
	core.Node                 // Embedded Node
	target     math32.Vector3 // Camera target in world coordinates
	up         math32.Vector3 // Camera Up vector
	viewMatrix math32.Matrix4 // Last calculated view matrix
}

// Initialize initializes the base camera.
// Normally used by other camera types which embed this base camera.
func (cam *Camera) Initialize() {

	cam.Node.Init()
	cam.target.Set(0, 0, 0)
	cam.up.Set(0, 1, 0)
	cam.SetDirection(0, 0, -1)
}

// LookAt rotates the camera to look at the specified target position.
// This method does not support objects with rotated and/or translated parent(s).
// TODO: maybe move method to Node, or create similar in Node.
func (cam *Camera) LookAt(target *math32.Vector3) {

	cam.target = *target

	var rotMat math32.Matrix4
	pos := cam.Position()
	rotMat.LookAt(&pos, &cam.target, &cam.up)

	var q math32.Quaternion
	q.SetFromRotationMatrix(&rotMat)
	cam.SetQuaternionQuat(&q)
}

// GetCamera satisfies the ICamera interface.
func (cam *Camera) GetCamera() *Camera {

	return cam
}

// Target get the current target position.
func (cam *Camera) Target() math32.Vector3 {

	return cam.target
}

// Up get the current camera up vector.
func (cam *Camera) Up() math32.Vector3 {

	return cam.up
}

// SetUp sets the camera up vector.
func (cam *Camera) SetUp(up *math32.Vector3) {

	cam.up = *up
	cam.LookAt(&cam.target) // TODO Maybe remove and let user call LookAt explicitly
}

// ViewMatrix returns the current view matrix of this camera.
func (cam *Camera) ViewMatrix(m *math32.Matrix4) {

	cam.UpdateMatrixWorld()
	matrixWorld := cam.MatrixWorld()
	err := m.GetInverse(&matrixWorld)
	if err != nil {
		panic("Camera.ViewMatrix: Couldn't invert matrix")
	}
}
