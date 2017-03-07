// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

type Sphere struct {
	Center Vector3 // center of the sphere
	Radius float32 // radius of the sphere
}

// NewSphere creates and returns a pointer to a new sphere with
// the specified center and radius
func NewSphere(center *Vector3, radius float32) *Sphere {

	s := new(Sphere)
	s.Center = *center
	s.Radius = radius
	return s
}

// Set sets the center and radius of the sphere
func (s *Sphere) Set(center *Vector3, radius float32) *Sphere {

	s.Center = *center
	s.Radius = radius
	return s
}

func (this *Sphere) SetFromPoints(points []Vector3, optionalCenter *Vector3) *Sphere {

	box := NewBox3(nil, nil)

	if optionalCenter != nil {
		this.Center.Copy(optionalCenter)
	} else {
		box.SetFromPoints(points).Center(&this.Center)
	}
	var maxRadiusSq float32 = 0.0
	for i := 0; i < len(points); i++ {
		maxRadiusSq = Max(maxRadiusSq, this.Center.DistanceToSquared(&points[i]))
	}
	this.Radius = Sqrt(maxRadiusSq)
	return this
}

func (this *Sphere) Copy(sphere *Sphere) *Sphere {

	*this = *sphere
	return this
}

// Empty checks if this sphere is empty (radius <= 0)
func (s *Sphere) Empty(sphere *Sphere) bool {

	if s.Radius <= 0 {
		return true
	}
	return false
}

// ContainsPoint checks if this sphere contains the specified point
func (s *Sphere) ContainsPoint(point *Vector3) bool {

	if point.DistanceToSquared(&s.Center) <= (s.Radius * s.Radius) {
		return true
	}
	return false
}

func (this *Sphere) DistanceToPoint(point *Vector3) float32 {

	return point.DistanceTo(&this.Center) - this.Radius
}

func (this *Sphere) IntersectSphere(sphere *Sphere) bool {

	radiusSum := this.Radius + sphere.Radius
	if sphere.Center.DistanceToSquared(&this.Center) <= (radiusSum * radiusSum) {
		return true
	}
	return false
}

func (this *Sphere) ClampPoint(point *Vector3, optionalTarget *Vector3) *Vector3 {

	deltaLengthSq := this.Center.DistanceToSquared(point)

	var result *Vector3
	if optionalTarget != nil {
		result = optionalTarget
	} else {
		result = NewVector3(0, 0, 0)
	}
	result.Copy(point)

	if deltaLengthSq > (this.Radius * this.Radius) {
		result.Sub(&this.Center).Normalize()
		result.MultiplyScalar(this.Radius).Add(&this.Center)
	}
	return result
}

func (s *Sphere) GetBoundingBox(optionalTarget *Box3) *Box3 {

	var box *Box3
	if optionalTarget != nil {
		box = optionalTarget
	} else {
		box = NewBox3(nil, nil)
	}

	box.Set(&s.Center, &s.Center)
	box.ExpandByScalar(s.Radius)
	return box
}

// ApplyMatrix4 applies the specified matrix transform to this sphere
func (s *Sphere) ApplyMatrix4(matrix *Matrix4) *Sphere {

	s.Center.ApplyMatrix4(matrix)
	s.Radius = s.Radius * matrix.GetMaxScaleOnAxis()
	return s
}

func (this *Sphere) Translate(offset *Vector3) *Sphere {

	this.Center.Add(offset)
	return this
}
