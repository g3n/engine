// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// Line3 represents a 3D line segment defined by a start and an end point.
type Line3 struct {
	start Vector3
	end   Vector3
}

// NewLine3 creates and returns a pointer to a new Line3 with the
// specified start and end points.
func NewLine3(start, end *Vector3) *Line3 {

	l := new(Line3)
	l.Set(start, end)
	return l
}

// Set sets this line segment start and end points.
// Returns pointer to this updated line segment.
func (l *Line3) Set(start, end *Vector3) *Line3 {

	if start != nil {
		l.start = *start
	}
	if end != nil {
		l.end = *end
	}
	return l
}

// Copy copy other line segment to this one.
// Returns pointer to this updated line segment.
func (l *Line3) Copy(other *Line3) *Line3 {

	*l = *other
	return l
}

// Center calculates this line segment center point.
// Store its pointer into optionalTarget, if not nil, and also returns it.
func (l *Line3) Center(optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return result.AddVectors(&l.start, &l.end).MultiplyScalar(0.5)
}

// Delta calculates the vector from the start to end point of this line segment.
// Store its pointer in optionalTarget, if not nil, and also returns it.
func (l *Line3) Delta(optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return result.SubVectors(&l.end, &l.start)
}

// DistanceSq returns the square of the distance from the start point to the end point.
func (l *Line3) DistanceSq() float32 {

	return l.start.DistanceToSquared(&l.end)
}

// Distance returns the distance from the start point to the end point.
func (l *Line3) Distance() float32 {

	return l.start.DistanceTo(&l.end)
}

// ApplyMatrix4 applies the specified matrix to this line segment start and end points.
// Returns pointer to this updated line segment.
func (l *Line3) ApplyMatrix4(matrix *Matrix4) *Line3 {

	l.start.ApplyMatrix4(matrix)
	l.end.ApplyMatrix4(matrix)
	return l
}

// Equals returns if this line segement is equal to other.
func (l *Line3) Equals(other *Line3) bool {

	return other.start.Equals(&l.start) && other.end.Equals(&l.end)
}

// Clone creates and returns a pointer to a copy of this line segment.
func (l *Line3) Clone() *Line3 {

	return NewLine3(&l.start, &l.end)
}
