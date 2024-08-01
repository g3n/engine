// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package camera contains virtual cameras and associated controls.
package camera

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/logger"
)

// Package logger
var log = logger.New("CAMERA", logger.Default)

// Axis represents a camera axis.
type Axis int

// The two possible camera axes.
const (
	Vertical = Axis(iota)
	Horizontal
)

// Projection represents a camera projection.
type Projection int

// The possible camera projections.
const (
	Perspective = Projection(iota)
	Orthographic
)

// ICamera is the interface for all cameras.
type ICamera interface {
	ViewMatrix(m *math32.Matrix4)
	ProjMatrix(m *math32.Matrix4)
}

// Camera represents a virtual camera, which specifies how to project a 3D scene onto an image.
type Camera struct {
	core.Node                  // Embedded Node
	aspect      float32        // Aspect ratio (width/height)
	near        float32        // Near plane depth
	far         float32        // Far plane depth
	axis        Axis           // The reference axis
	proj        Projection     // Projection method
	fov         float32        // Perspective field-of-view along reference axis
	size        float32        // Orthographic size along reference axis
	projChanged bool           // Flag indicating that the projection matrix needs to be recalculated
	projMatrix  math32.Matrix4 // Last calculated projection matrix
}

// New creates and returns a new perspective camera with the specified aspect ratio and default parameters.
func New(aspect float32) *Camera {
	return NewPerspective(aspect, 0.3, 1000, 60, Vertical)
}

// NewPerspective creates and returns a new perspective camera with the specified parameters.
func NewPerspective(aspect, near, far, fov float32, axis Axis) *Camera {

	c := new(Camera)
	c.Node.Init(c)
	c.SetDirection(0, 0, -1)
	c.aspect = aspect
	c.near = near
	c.far = far
	c.axis = axis
	c.proj = Perspective
	c.fov = fov
	c.size = 8
	c.projChanged = true
	return c
}

// NewOrthographic creates and returns a new orthographic camera with the specified parameters.
func NewOrthographic(aspect, near, far, size float32, axis Axis) *Camera {

	c := new(Camera)
	c.Node.Init(c)
	c.SetDirection(0, 0, -1)
	c.aspect = aspect
	c.near = near
	c.far = far
	c.axis = axis
	c.proj = Orthographic
	c.fov = 60
	c.size = size
	c.projChanged = true
	return c
}

// Aspect returns the camera aspect ratio.
func (c *Camera) Aspect() float32 {

	return c.aspect
}

// SetAspect sets the camera aspect ratio.
func (c *Camera) SetAspect(aspect float32) {

	if aspect == c.aspect {
		return
	}
	c.aspect = aspect
	c.projChanged = true
}

// Near returns the camera near plane Z coordinate.
func (c *Camera) Near() float32 {

	return c.near
}

// SetNear sets the camera near plane Z coordinate.
func (c *Camera) SetNear(near float32) {

	if near == c.near {
		return
	}
	c.near = near
	c.projChanged = true
}

// Far returns the camera far plane Z coordinate.
func (c *Camera) Far() float32 {

	return c.far
}

// SetFar sets the camera far plane Z coordinate.
func (c *Camera) SetFar(far float32) {

	if far == c.far {
		return
	}
	c.far = far
	c.projChanged = true
}

// Axis returns the reference axis associated with the camera size/fov.
func (c *Camera) Axis() Axis {

	return c.axis
}

// SetAxis sets the reference axis associated with the camera size/fov.
func (c *Camera) SetAxis(axis Axis) {

	if axis == c.axis {
		return
	}
	c.axis = axis
	c.projChanged = true
}

// Projection returns the projection method used by the camera.
func (c *Camera) Projection() Projection {

	return c.proj
}

// SetProjection sets the projection method used by the camera.
func (c *Camera) SetProjection(proj Projection) {

	if proj == c.proj {
		return
	}
	c.proj = proj
	c.projChanged = true
}

// Fov returns the perspective field-of-view in degrees along the reference axis.
func (c *Camera) Fov() float32 {

	return c.fov
}

// SetFov sets the perspective field-of-view in degrees along the reference axis.
func (c *Camera) SetFov(fov float32) {

	if fov == c.fov {
		return
	}
	c.fov = fov
	if c.proj == Perspective {
		c.projChanged = true
	}
}

// UpdateFov updates the field-of-view such that the frustum matches
// the orthographic size at the depth specified by targetDist.
func (c *Camera) UpdateFov(targetDist float32) {

	c.fov = 2 * math32.Atan(c.size/(2*targetDist)) * 180 / math32.Pi
	if c.proj == Perspective {
		c.projChanged = true
	}
}

// Size returns the orthographic view size along the camera's reference axis.
func (c *Camera) Size() float32 {

	return c.size
}

// SetSize sets the orthographic view size along the camera's reference axis.
func (c *Camera) SetSize(size float32) {

	if size == c.size {
		return
	}
	c.size = size
	if c.proj == Orthographic {
		c.projChanged = true
	}
}

// UpdateSize updates the orthographic size to match the current
// field-of-view frustum at the depth specified by targetDist.
func (c *Camera) UpdateSize(targetDist float32) {

	c.size = 2 * targetDist * math32.Tan(math32.Pi/180*c.fov/2)
	if c.proj == Orthographic {
		c.projChanged = true
	}
}

// ViewMatrix returns the view matrix of the camera.
func (c *Camera) ViewMatrix(m *math32.Matrix4) {

	c.UpdateMatrixWorld()
	matrixWorld := c.MatrixWorld()
	err := m.GetInverse(&matrixWorld)
	if err != nil {
		panic("Camera.ViewMatrix: Couldn't invert matrix")
	}
}

// ProjMatrix returns the projection matrix of the camera.
func (c *Camera) ProjMatrix(m *math32.Matrix4) {

	if c.projChanged {
		switch c.proj {
		case Perspective:
			t := c.near * math32.Tan(math32.DegToRad(c.fov*0.5))
			ymax := t
			ymin := -t
			xmax := t
			xmin := -t
			switch c.axis {
			case Vertical:
				xmax *= c.aspect
				xmin *= c.aspect
			case Horizontal:
				ymax /= c.aspect
				ymin /= c.aspect
			}
			c.projMatrix.MakeFrustum(xmin, xmax, ymin, ymax, c.near, c.far)
		case Orthographic:
			s := c.size / 2
			var h, w float32
			switch c.axis {
			case Vertical:
				h = s
				w = s * c.aspect
			case Horizontal:
				h = s / c.aspect
				w = s
			}
			c.projMatrix.MakeOrthographic(-w, w, h, -h, c.near, c.far)
		}
		c.projChanged = false
	}
	*m = c.projMatrix
}

// Project transforms the specified position from world coordinates to this camera projected coordinates.
func (c *Camera) Project(v *math32.Vector3) *math32.Vector3 {

	// Get camera view matrix
	var viewMat, projMat math32.Matrix4
	c.ViewMatrix(&viewMat)
	c.ProjMatrix(&projMat)

	// Apply projMat * viewMat to the provided vector
	v.ApplyProjection(projMat.Multiply(&viewMat))
	return v
}

// Unproject transforms the specified position from camera projected coordinates to world coordinates.
func (c *Camera) Unproject(v *math32.Vector3) *math32.Vector3 {

	// Get inverted camera view matrix
	invViewMat := c.MatrixWorld()

	// Get inverted camera projection matrix
	var invProjMat math32.Matrix4
	c.ProjMatrix(&invProjMat)
	err := invProjMat.GetInverse(&invProjMat)
	if err != nil {
		panic("Camera.Unproject: Couldn't invert matrix")
	}

	// Apply invViewMat * invProjMat to the provided vector
	v.ApplyProjection(invViewMat.Multiply(&invProjMat))
	return v
}
