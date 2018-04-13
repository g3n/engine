// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// Triangle represents a triangle made of three vertices.
type Triangle struct {
	a Vector3
	b Vector3
	c Vector3
}

// NewTriangle returns a pointer to a new Triangle object.
func NewTriangle(a, b, c *Vector3) *Triangle {

	t := new(Triangle)
	if a != nil {
		t.a = *a
	}
	if b != nil {
		t.b = *b
	}
	if c != nil {
		t.c = *c
	}
	return t
}

// Normal returns the triangle's normal.
func Normal(a, b, c, optionalTarget *Vector3) *Vector3 {

	var v0 Vector3
	var result *Vector3
	if optionalTarget != nil {
		result = optionalTarget
	} else {
		result = NewVector3(0, 0, 0)
	}

	result.SubVectors(c, b)
	v0.SubVectors(a, b)
	result.Cross(&v0)

	resultLengthSq := result.LengthSq()
	if resultLengthSq > 0 {
		return result.MultiplyScalar(1 / Sqrt(resultLengthSq))
	}
	return result.Set(0, 0, 0)
}

// BarycoordFromPoint returns the barycentric coordinates for the specified point.
func BarycoordFromPoint(point, a, b, c, optionalTarget *Vector3) *Vector3 {

	var v0 Vector3
	var v1 Vector3
	var v2 Vector3

	v0.SubVectors(c, a)
	v1.SubVectors(b, a)
	v2.SubVectors(point, a)

	dot00 := v0.Dot(&v0)
	dot01 := v0.Dot(&v1)
	dot02 := v0.Dot(&v2)
	dot11 := v1.Dot(&v1)
	dot12 := v1.Dot(&v2)

	denom := dot00*dot11 - dot01*dot01

	var result *Vector3
	if optionalTarget != nil {
		result = optionalTarget
	} else {
		result = NewVector3(0, 0, 0)
	}

	// colinear or singular triangle
	if denom == 0 {
		// arbitrary location outside of triangle?
		// not sure if this is the best idea, maybe should be returning undefined
		return result.Set(-2, -1, -1)
	}

	invDenom := 1 / denom
	u := (dot11*dot02 - dot01*dot12) * invDenom
	v := (dot00*dot12 - dot01*dot02) * invDenom

	// barycoordinates must always sum to 1
	return result.Set(1-u-v, v, u)

}

// ContainsPoint returns whether a triangle contains a point.
func ContainsPoint(point, a, b, c *Vector3) bool {

	var v1 Vector3
	result := BarycoordFromPoint(point, a, b, c, &v1)

	return (result.X >= 0) && (result.Y >= 0) && ((result.X + result.Y) <= 1)
}

// Set sets the triangle's three vertices.
func (t *Triangle) Set(a, b, c *Vector3) *Triangle {

	t.a = *a
	t.b = *b
	t.c = *c
	return t
}

// SetFromPointsAndIndices sets the triangle's vertices based on the specified points and indices.
func (t *Triangle) SetFromPointsAndIndices(points []*Vector3, i0, i1, i2 int) *Triangle {

	t.a = *points[i0]
	t.b = *points[i1]
	t.c = *points[i2]
	return t
}

// Copy modifies the receiver triangle to match the provided triangle.
func (t *Triangle) Copy(triangle *Triangle) *Triangle {

	*t = *triangle
	return t
}

// Area returns the triangle's area.
func (t *Triangle) Area() float32 {

	var v0 Vector3
	var v1 Vector3

	v0.SubVectors(&t.c, &t.b)
	v1.SubVectors(&t.a, &t.b)
	return v0.Cross(&v1).Length() * 0.5
}

// Midpoint returns the triangle's midpoint.
func (t *Triangle) Midpoint(optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget != nil {
		result = optionalTarget
	} else {
		result = NewVector3(0, 0, 0)
	}
	return result.AddVectors(&t.a, &t.b).Add(&t.c).MultiplyScalar(1 / 3)
}

// Normal returns the triangle's normal.
func (t *Triangle) Normal(optionalTarget *Vector3) *Vector3 {

	return Normal(&t.a, &t.b, &t.c, optionalTarget)
}

// Plane returns a Plane object aligned with the triangle.
func (t *Triangle) Plane(optionalTarget *Plane) *Plane {

	var result *Plane
	if optionalTarget != nil {
		result = optionalTarget
	} else {
		result = NewPlane(nil, 0)
	}
	return result.SetFromCoplanarPoints(&t.a, &t.b, &t.c)
}

// BarycoordFromPoint returns the barycentric coordinates for the specified point.
func (t *Triangle) BarycoordFromPoint(point, optionalTarget *Vector3) *Vector3 {

	return BarycoordFromPoint(point, &t.a, &t.b, &t.c, optionalTarget)
}

// ContainsPoint returns whether the triangle contains a point.
func (t *Triangle) ContainsPoint(point *Vector3) bool {

	return ContainsPoint(point, &t.a, &t.b, &t.c)
}

// Equals returns whether the triangles are equal in all their vertices.
func (t *Triangle) Equals(triangle *Triangle) bool {

	return triangle.a.Equals(&t.a) && triangle.b.Equals(&t.b) && triangle.c.Equals(&t.c)
}

// Clone clones a triangle.
func (t *Triangle) Clone(triangle *Triangle) *Triangle {

	return NewTriangle(nil, nil, nil).Copy(t)
}
