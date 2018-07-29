// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shape

import "github.com/g3n/engine/math32"

// IShape is the interface for all collision shapes.
// Shapes in this package satisfy this interface and also geometry.Geometry.
type IShape interface {
	BoundingBox() math32.Box3
	BoundingSphere() math32.Sphere
	Area() float32
	Volume() float32
	RotationalInertia(mass float32) math32.Matrix3
	ProjectOntoAxis(localAxis *math32.Vector3) (float32, float32)
}

// Shape is a collision shape.
// It can be an analytical geometry such as a sphere, plane, etc.. or it can be defined by a polygonal Geometry.
type Shape struct {

	// TODO
	//material

	// Collision filtering
	colFilterGroup int
	colFilterMask  int
}

func (s *Shape) initialize() {

	// Collision filtering
	s.colFilterGroup = 1
	s.colFilterMask = -1
}
