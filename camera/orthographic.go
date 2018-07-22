// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package camera

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

// Orthographic is an orthographic camera.
type Orthographic struct {
	Camera              // Embedded camera
	left        float32 // left plane x coordinate
	right       float32 // right plane x coordinate
	top         float32 // top plane y coordinate
	bottom      float32 // bottom plane y coordinate
	near        float32 // near plane z coordinate
	far         float32 // far plane z coordinate
	zoom        float32
	projChanged bool           // camera projection parameters changed (needs to recalculates projection matrix)
	projMatrix  math32.Matrix4 // last calculated projection matrix
}

// NewOrthographic creates and returns a pointer to a new orthographic camera with the specified parameters.
func NewOrthographic(left, right, top, bottom, near, far float32) *Orthographic {

	cam := new(Orthographic)
	cam.Camera.Initialize()
	cam.left = left
	cam.right = right
	cam.top = top
	cam.bottom = bottom
	cam.near = near
	cam.far = far
	cam.zoom = 1.0
	cam.projChanged = true
	return cam
}

// SetAspect sets the camera aspect ratio (width/height).
func (cam *Orthographic) SetAspect(aspect float32) {

	height := cam.top - cam.bottom
	halfwidth := height * aspect * 0.5
	center := (cam.left + cam.right) * 0.5
	cam.left = center - halfwidth
	cam.right = center + halfwidth
	cam.projChanged = true
}

// SetZoom sets the zoom factor of the camera.
func (cam *Orthographic) SetZoom(zoom float32) {

	cam.zoom = math32.Abs(zoom)
	cam.projChanged = true
}

// Zoom returns the zoom factor of the camera.
func (cam *Orthographic) Zoom() float32 {

	return cam.zoom
}

// Planes returns the coordinates of the camera planes.
func (cam *Orthographic) Planes() (left, right, top, bottom, near, far float32) {

	return cam.left, cam.right, cam.top, cam.bottom, cam.near, cam.far
}

// ProjMatrix satisfies the ICamera interface.
func (cam *Orthographic) ProjMatrix(m *math32.Matrix4) {

	if cam.projChanged {
		cam.projMatrix.MakeOrthographic(cam.left/cam.zoom, cam.right/cam.zoom, cam.top/cam.zoom, cam.bottom/cam.zoom, cam.near, cam.far)
		cam.projChanged = false
	}
	*m = cam.projMatrix
}

// Project satisfies the ICamera interface and must
// be implemented for specific camera types.
func (cam *Camera) Project(v *math32.Vector3) (*math32.Vector3, error) {

	panic("Not implemented")
}

// Unproject satisfies the ICamera interface and must
// be implemented for specific camera types.
func (cam *Camera) Unproject(v *math32.Vector3) (*math32.Vector3, error) {

	panic("Not implemented")
}

// SetRaycaster satisfies the ICamera interface and must
// be implemented for specific camera types.
func (cam *Camera) SetRaycaster(rc *core.Raycaster, x, y float32) error {

	panic("Not implemented")
}
