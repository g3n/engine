// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package camera

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

// Perspective is a perspective camera.
type Perspective struct {
	Camera                     // Embedded camera
	fov         float32        // field of view in degrees
	aspect      float32        // aspect ratio (width/height)
	near        float32        // near plane z coordinate
	far         float32        // far plane z coordinate
	projChanged bool           // camera projection parameters changed (needs to recalculates projection matrix)
	projMatrix  math32.Matrix4 // last calculated projection matrix
}

// NewPerspective creates and returns a pointer to a new perspective camera with the
// specified parameters.
func NewPerspective(fov, aspect, near, far float32) *Perspective {

	cam := new(Perspective)
	cam.Camera.Initialize()
	cam.fov = fov
	cam.aspect = aspect
	cam.near = near
	cam.far = far
	cam.projChanged = true
	return cam
}

// SetFov sets the camera field of view in degrees.
func (cam *Perspective) SetFov(fov float32) {

	cam.fov = fov
	cam.projChanged = true
}

// SetAspect sets the camera aspect ratio (width/height).
func (cam *Perspective) SetAspect(aspect float32) {

	cam.aspect = aspect
	cam.projChanged = true
}

// Fov returns the current camera FOV (field of view) in degrees.
func (cam *Perspective) Fov() float32 {

	return cam.fov
}

// Aspect returns the current camera aspect ratio.
func (cam Perspective) Aspect() float32 {

	return cam.aspect
}

// Near returns the current camera near plane Z coordinate.
func (cam *Perspective) Near() float32 {

	return cam.near
}

// Far returns the current camera far plane Z coordinate.
func (cam *Perspective) Far() float32 {

	return cam.far
}

// ProjMatrix satisfies the ICamera interface.
func (cam *Perspective) ProjMatrix(m *math32.Matrix4) {

	cam.updateProjMatrix()
	*m = cam.projMatrix
}

// Project transforms the specified position from world coordinates to this camera projected coordinates.
func (cam *Perspective) Project(v *math32.Vector3) (*math32.Vector3, error) {

	// Get camera view matrix
	var matrix math32.Matrix4
	matrixWorld := cam.MatrixWorld()
	err := matrix.GetInverse(&matrixWorld)
	if err != nil {
		return nil, err
	}

	// Update camera projection matrix
	cam.updateProjMatrix()

	// Multiply viewMatrix by projMatrix and apply the resulting projection matrix to the provided vector
	matrix.MultiplyMatrices(&cam.projMatrix, &matrix)
	v.ApplyProjection(&matrix)
	return v, nil
}

// Unproject transforms the specified position from camera projected coordinates to world coordinates.
func (cam *Perspective) Unproject(v *math32.Vector3) (*math32.Vector3, error) {

	// Get inverted camera view matrix
	invertedViewMatrix := cam.MatrixWorld()

	// Get inverted camera projection matrix
	cam.updateProjMatrix()
	var invertedProjMatrix math32.Matrix4
	err := invertedProjMatrix.GetInverse(&cam.projMatrix)
	if err != nil {
		return nil, err
	}

	// Multiply invertedViewMatrix by invertedProjMatrix
	// to get transformation from camera projected coordinates to world coordinates
	// and project vector using this transformation
	var matrix math32.Matrix4
	matrix.MultiplyMatrices(&invertedViewMatrix, &invertedProjMatrix)
	v.ApplyProjection(&matrix)
	return v, nil
}

// SetRaycaster sets the specified raycaster with this camera position in world coordinates
// pointing to the direction defined by the specified coordinates unprojected using this camera.
func (cam *Perspective) SetRaycaster(rc *core.Raycaster, sx, sy float32) error {

	var origin, direction math32.Vector3
	matrixWorld := cam.MatrixWorld()
	origin.SetFromMatrixPosition(&matrixWorld)
	direction.Set(sx, sy, 0.5)
	unproj, err := cam.Unproject(&direction)
	if err != nil {
		return err
	}
	unproj.Sub(&origin).Normalize()
	rc.Set(&origin, &direction)
	// Updates the view matrix of the raycaster
	cam.ViewMatrix(&rc.ViewMatrix)
	return nil
}

// updateProjMatrix updates the projection matrix if necessary.
func (cam *Perspective) updateProjMatrix() {

	if cam.projChanged {
		cam.projMatrix.MakePerspective(cam.fov, cam.aspect, cam.near, cam.far)
		cam.projChanged = false
	}
}
