// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shape

import "github.com/g3n/engine/math32"

// Sphere is an analytical collision sphere.
type Sphere struct {
	radius float32
}

// NewSphere creates and returns a pointer to a new analytical collision sphere.
func NewSphere(radius float32) *Sphere {

	s := new(Sphere)
	s.radius = radius
	return s
}

// SetRadius sets the radius of the analytical collision sphere.
func (s *Sphere) SetRadius(radius float32) {

	s.radius = radius
}

// Radius returns the radius of the analytical collision sphere.
func (s *Sphere) Radius() float32 {

	return s.radius
}

// IShape =============================================================

// BoundingBox computes and returns the bounding box of the analytical collision sphere.
func (s *Sphere) BoundingBox() math32.Box3 {

	return math32.Box3{math32.Vector3{-s.radius, -s.radius, -s.radius}, math32.Vector3{s.radius, s.radius, s.radius}}
}

// BoundingSphere computes and returns the bounding sphere of the analytical collision sphere.
func (s *Sphere) BoundingSphere() math32.Sphere {

	return *math32.NewSphere(math32.NewVec3(), s.radius)
}

// Area computes and returns the surface area of the analytical collision sphere.
func (s *Sphere) Area() float32 {

	return 4 * math32.Pi * s.radius * s.radius
}

// Volume computes and returns the volume of the analytical collision sphere.
func (s *Sphere) Volume() float32 {

	return (4/3) * math32.Pi * s.radius * s.radius * s.radius
}

// RotationalInertia computes and returns the rotational inertia of the analytical collision sphere.
func (s *Sphere) RotationalInertia(mass float32) math32.Matrix3 {

	v := (2/5) * mass * s.radius * s.radius
	return *math32.NewMatrix3().Set(
		v, 0, 0,
		0, v, 0,
		0, 0, v,
	)
}

// ProjectOntoAxis computes and returns the minimum and maximum distances of the analytical collision sphere projected onto the specified local axis.
func (s *Sphere) ProjectOntoAxis(localAxis *math32.Vector3) (float32, float32) {

	return -s.radius, s.radius
}
