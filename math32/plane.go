// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

import ()

type Plane struct {
	normal   Vector3
	constant float32
}

func NewPlane(normal *Vector3, constant float32) *Plane {

	this := new(Plane)
	if normal != nil {
		this.normal = *normal
	}
	this.constant = constant
	return this
}

func (this *Plane) Set(normal *Vector3, constant float32) *Plane {

	this.normal = *normal
	this.constant = constant
	return this
}

func (this *Plane) SetComponents(x, y, z, w float32) *Plane {

	this.normal.Set(x, y, z)
	this.constant = w
	return this
}

func (this *Plane) SetFromNormalAndCoplanarPoint(normal *Vector3, point *Vector3) *Plane {

	this.normal.Copy(normal)
	this.constant = -point.Dot(&this.normal)
	return this
}

func (this *Plane) SetFromCoplanarPoints(a, b, c *Vector3) *Plane {

	var v1 Vector3
	var v2 Vector3

	normal := v1.SubVectors(c, b).Cross(v2.SubVectors(a, b)).Normalize()
	// Q: should an error be thrown if normal is zero (e.g. degenerate plane)?
	this.SetFromNormalAndCoplanarPoint(normal, a)
	return this
}

func (this *Plane) Copy(plane *Plane) *Plane {

	this.normal.Copy(&plane.normal)
	this.constant = plane.constant
	return this
}

func (this *Plane) Normalize() *Plane {

	// Note: will lead to a divide by zero if the plane is invalid.
	inverseNormalLength := 1.0 / this.normal.Length()
	this.normal.MultiplyScalar(inverseNormalLength)
	this.constant *= inverseNormalLength
	return this
}

func (this *Plane) Negate() *Plane {

	this.constant *= -1
	this.normal.Negate()
	return this
}

func (this *Plane) DistanceToPoint(point *Vector3) float32 {

	return this.normal.Dot(point) + this.constant
}

func (this *Plane) DistanceToSphere(sphere *Sphere) float32 {

	return this.DistanceToPoint(&sphere.Center) - sphere.Radius
}

func (this *Plane) ProjectPoint(point *Vector3, optionalTarget *Vector3) *Vector3 {

	return this.OrthoPoint(point, optionalTarget).Sub(point).Negate()
}

func (this *Plane) OrthoPoint(point *Vector3, optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	perpendicularMagnitude := this.DistanceToPoint(point)
	return result.Copy(&this.normal).MultiplyScalar(perpendicularMagnitude)
}

func (this *Plane) IsIntersectionLine(line *Line3) bool {

	// Note: this tests if a line intersects the plane, not whether it (or its end-points) are coplanar with it.
	startSign := this.DistanceToPoint(&line.start)
	endSign := this.DistanceToPoint(&line.end)

	return (startSign < 0 && endSign > 0) || (endSign < 0 && startSign > 0)

}

func (this *Plane) IntersectLine(line *Line3, optionalTarget *Vector3) *Vector3 {

	var v1 Vector3
	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}

	direction := line.Delta(&v1)
	denominator := this.normal.Dot(direction)
	if denominator == 0 {
		// line is coplanar, return origin
		if this.DistanceToPoint(&line.start) == 0 {
			return result.Copy(&line.start)
		}
		// Unsure if this is the correct method to handle this case.
		return nil
	}

	var t = -(line.start.Dot(&this.normal) + this.constant) / denominator
	if t < 0 || t > 1 {
		return nil
	}
	return result.Copy(direction).MultiplyScalar(t).Add(&line.start)
}

func (this *Plane) CoplanarPoint(optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return result.Copy(&this.normal).MultiplyScalar(-this.constant)
}

func (this *Plane) ApplyMatrix4(matrix *Matrix4, optionalNormalMatrix *Matrix3) *Plane {
	// compute new normal based on theory here:
	// http://www.songho.ca/opengl/gl_normaltransform.html

	var v1 Vector3
	var v2 Vector3
	m1 := NewMatrix3()

	var normalMatrix *Matrix3
	if optionalNormalMatrix != nil {
		normalMatrix = optionalNormalMatrix
	} else {
		normalMatrix = m1.GetNormalMatrix(matrix)
	}

	newNormal := v1.Copy(&this.normal).ApplyMatrix3(normalMatrix)

	newCoplanarPoint := this.CoplanarPoint(&v2)
	newCoplanarPoint.ApplyMatrix4(matrix)

	this.SetFromNormalAndCoplanarPoint(newNormal, newCoplanarPoint)
	return this
}

func (this *Plane) Translate(offset *Vector3) *Plane {

	this.constant = this.constant - offset.Dot(&this.normal)
	return this
}

func (this *Plane) Equals(plane *Plane) bool {

	return plane.normal.Equals(&this.normal) && (plane.constant == this.constant)
}

func (this *Plane) Clone(plane *Plane) *Plane {

	return NewPlane(&plane.normal, plane.constant)
}
