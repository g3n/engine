// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

type Triangle struct {
	a Vector3
	b Vector3
	c Vector3
}

func NewTriangle(a, b, c *Vector3) *Triangle {

	this := new(Triangle)
	if a != nil {
		this.a = *a
	}
	if b != nil {
		this.b = *b
	}
	if c != nil {
		this.c = *c
	}
	return this
}

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

func ContainsPoint(point, a, b, c *Vector3) bool {

	var v1 Vector3
	result := BarycoordFromPoint(point, a, b, c, &v1)

	return (result.X >= 0) && (result.Y >= 0) && ((result.X + result.Y) <= 1)
}

func (this *Triangle) Set(a, b, c *Vector3) *Triangle {

	this.a = *a
	this.b = *b
	this.c = *c
	return this
}

func (this *Triangle) SetFromPointsAndIndices(points []*Vector3, i0, i1, i2 int) *Triangle {

	this.a = *points[i0]
	this.b = *points[i1]
	this.c = *points[i2]
	return this
}

func (this *Triangle) Copy(triangle *Triangle) *Triangle {

	*this = *triangle
	return this
}

func (this *Triangle) Area() float32 {

	var v0 Vector3
	var v1 Vector3

	v0.SubVectors(&this.c, &this.b)
	v1.SubVectors(&this.a, &this.b)
	return v0.Cross(&v1).Length() * 0.5
}

func (this *Triangle) Midpoint(optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget != nil {
		result = optionalTarget
	} else {
		result = NewVector3(0, 0, 0)
	}
	return result.AddVectors(&this.a, &this.b).Add(&this.c).MultiplyScalar(1 / 3)
}

func (this *Triangle) Normal(optionalTarget *Vector3) *Vector3 {

	return Normal(&this.a, &this.b, &this.c, optionalTarget)
}

func (this *Triangle) Plane(optionalTarget *Plane) *Plane {

	var result *Plane
	if optionalTarget != nil {
		result = optionalTarget
	} else {
		result = NewPlane(nil, 0)
	}
	return result.SetFromCoplanarPoints(&this.a, &this.b, &this.c)
}

func (this *Triangle) BarycoordFromPoint(point, optionalTarget *Vector3) *Vector3 {

	return BarycoordFromPoint(point, &this.a, &this.b, &this.c, optionalTarget)
}

func (this *Triangle) ContainsPoint(point *Vector3) bool {

	return ContainsPoint(point, &this.a, &this.b, &this.c)
}

func (this *Triangle) Equals(triangle *Triangle) bool {

	return triangle.a.Equals(&this.a) && triangle.b.Equals(&this.b) && triangle.c.Equals(&this.c)
}

func (this *Triangle) Clone(triangle *Triangle) *Triangle {

	return NewTriangle(nil, nil, nil).Copy(this)
}
