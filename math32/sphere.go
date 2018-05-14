// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// Sphere represents a 3D sphere defined by its center point and a radius
type Sphere struct {
	Center Vector3 // center of the sphere
	Radius float32 // radius of the sphere
}

// NewSphere creates and returns a pointer to a new sphere with
// the specified center and radius.
func NewSphere(center *Vector3, radius float32) *Sphere {

	s := new(Sphere)
	s.Center = *center
	s.Radius = radius
	return s
}

// Set sets the center and radius of this sphere.
// Returns pointer to this update sphere.
func (s *Sphere) Set(center *Vector3, radius float32) *Sphere {

	s.Center = *center
	s.Radius = radius
	return s
}

// SetFromPoints sets this sphere from the specified points array and optional center.
// Returns pointer to this update sphere.
func (s *Sphere) SetFromPoints(points []Vector3, optionalCenter *Vector3) *Sphere {

	box := NewBox3(nil, nil)

	if optionalCenter != nil {
		s.Center.Copy(optionalCenter)
	} else {
		box.SetFromPoints(points).Center(&s.Center)
	}
	var maxRadiusSq float32
	for i := 0; i < len(points); i++ {
		maxRadiusSq = Max(maxRadiusSq, s.Center.DistanceToSquared(&points[i]))
	}
	s.Radius = Sqrt(maxRadiusSq)
	return s
}

// Copy copy other sphere to this one.
// Returns pointer to this update sphere.
func (s *Sphere) Copy(other *Sphere) *Sphere {

	*s = *other
	return s
}

// Empty checks if this sphere is empty (radius <= 0)
func (s *Sphere) Empty(sphere *Sphere) bool {

	if s.Radius <= 0 {
		return true
	}
	return false
}

// ContainsPoint returns if this sphere contains the specified point.
func (s *Sphere) ContainsPoint(point *Vector3) bool {

	if point.DistanceToSquared(&s.Center) <= (s.Radius * s.Radius) {
		return true
	}
	return false
}

// DistanceToPoint returns the distance from the sphere surface to the specified point.
func (s *Sphere) DistanceToPoint(point *Vector3) float32 {

	return point.DistanceTo(&s.Center) - s.Radius
}

// IntersectSphere returns if other sphere intersects this one.
func (s *Sphere) IntersectSphere(other *Sphere) bool {

	radiusSum := s.Radius + other.Radius
	if other.Center.DistanceToSquared(&s.Center) <= (radiusSum * radiusSum) {
		return true
	}
	return false
}

// ClampPoint clamps the specified point inside the sphere.
// If the specified point is inside the sphere, it is the clamped point.
// Otherwise the clamped point is the the point in the sphere surface in the
// nearest of the specified point.
// The clamped point is stored in optionalTarget, if not nil, and returned.
func (s *Sphere) ClampPoint(point *Vector3, optionalTarget *Vector3) *Vector3 {

	deltaLengthSq := s.Center.DistanceToSquared(point)

	var result *Vector3
	if optionalTarget != nil {
		result = optionalTarget
	} else {
		result = NewVector3(0, 0, 0)
	}
	result.Copy(point)

	if deltaLengthSq > (s.Radius * s.Radius) {
		result.Sub(&s.Center).Normalize()
		result.MultiplyScalar(s.Radius).Add(&s.Center)
	}
	return result
}

// GetBoundingBox calculates a Box3 which bounds this sphere.
// Update optionalTarget with the calculated Box3, if not nil, and also returns it.
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

// ApplyMatrix4 applies the specified matrix transform to this sphere.
// Returns pointer to this updated sphere.
func (s *Sphere) ApplyMatrix4(matrix *Matrix4) *Sphere {

	s.Center.ApplyMatrix4(matrix)
	s.Radius = s.Radius * matrix.GetMaxScaleOnAxis()
	return s
}

// Translate translates this sphere by the specified offset.
// Returns pointer to this updated sphere.
func (s *Sphere) Translate(offset *Vector3) *Sphere {

	s.Center.Add(offset)
	return s
}
