// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

type Box3 struct {
	Min Vector3
	Max Vector3
}

func NewBox3(min, max *Vector3) *Box3 {

	this := new(Box3)
	this.Set(min, max)
	return this
}

func (this *Box3) Set(min, max *Vector3) *Box3 {

	if min != nil {
		this.Min = *min
	} else {
		this.Min.Set(Infinity, Infinity, Infinity)
	}
	if max != nil {
		this.Max = *max
	} else {
		this.Max.Set(-Infinity, -Infinity, -Infinity)
	}
	return this
}

func (this *Box3) SetFromPoints(points []Vector3) *Box3 {

	this.MakeEmpty()
	for i := 0; i < len(points); i++ {
		this.ExpandByPoint(&points[i])
	}
	return this
}

func (this *Box3) SetFromCenterAndSize(center, size *Vector3) *Box3 {

	v1 := NewVector3(0, 0, 0)
	halfSize := v1.Copy(size).MultiplyScalar(0.5)
	this.Min.Copy(center).Sub(halfSize)
	this.Max.Copy(center).Add(halfSize)
	return this
}

//func (this *Box3) SetFromObject(object *Object3D) *Box3 {
//
//	// TODO object.UpdateMatrixWorld(true)
//
//	return this
//}

func (this *Box3) Copy(box *Box3) *Box3 {

	this.Min = box.Min
	this.Max = box.Max
	return this
}

func (this *Box3) MakeEmpty() *Box3 {

	this.Min.X = Infinity
	this.Min.Y = Infinity
	this.Min.Z = Infinity
	this.Max.X = -Infinity
	this.Max.Y = -Infinity
	this.Max.Z = -Infinity
	return this
}

func (this *Box3) Empty() bool {

	return (this.Max.X < this.Min.X) || (this.Max.Y < this.Min.Y) || (this.Max.Z < this.Min.Z)
}

func (this *Box3) Center(optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return result.AddVectors(&this.Min, &this.Max).MultiplyScalar(0.5)
}

func (this *Box3) Size(optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return result.SubVectors(&this.Min, &this.Max)
}

func (this *Box3) ExpandByPoint(point *Vector3) *Box3 {

	this.Min.Min(point)
	this.Max.Max(point)
	return this
}

func (this *Box3) ExpandByVector(vector *Vector3) *Box3 {

	this.Min.Sub(vector)
	this.Max.Add(vector)
	return this
}

func (this *Box3) ExpandByScalar(scalar float32) *Box3 {

	this.Min.AddScalar(-scalar)
	this.Max.AddScalar(scalar)
	return this
}

func (this *Box3) ContainsPoint(point *Vector3) bool {

	if point.X < this.Min.X || point.X > this.Max.X ||
		point.Y < this.Min.Y || point.Y > this.Max.Y ||
		point.Z < this.Min.Z || point.Z > this.Max.Z {
		return false
	}
	return true
}

func (this *Box3) ContainsBox(box *Box3) bool {

	if (this.Min.X <= box.Max.X) && (box.Max.X <= this.Max.X) &&
		(this.Min.Y <= box.Min.Y) && (box.Max.Y <= this.Max.Y) &&
		(this.Min.Z <= box.Min.Z) && (box.Max.Z <= this.Max.Z) {
		return true

	}
	return false
}

func (this *Box3) GetParameter(point *Vector3, optionalTarget *Vector3) *Vector3 {

	// This can potentially have a divide by zero if the box
	// has a size dimension of 0.
	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return result.Set(
		(point.X-this.Min.X)/(this.Max.X-this.Min.X),
		(point.Y-this.Min.Y)/(this.Max.Y-this.Min.Y),
		(point.Z-this.Min.Z)/(this.Max.Z-this.Min.Z),
	)
}

func (this *Box3) IsIntersectionBox(box *Box3) bool {

	// using 6 splitting planes to rule out intersections.
	if box.Max.X < this.Min.X || box.Min.X > this.Max.X ||
		box.Max.Y < this.Min.Y || box.Min.Y > this.Max.Y ||
		box.Max.Z < this.Min.Z || box.Min.Z > this.Max.Z {
		return false
	}
	return true
}

func (this *Box3) ClampPoint(point *Vector3, optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return result.Copy(point).Clamp(&this.Min, &this.Max)
}

func (this *Box3) DistanceToPoint(point *Vector3) float32 {

	var v1 Vector3
	clampedPoint := v1.Copy(point).Clamp(&this.Min, &this.Max)
	return clampedPoint.Sub(point).Length()
}

func (this *Box3) GetBoundingSphere(optionalTarget *Sphere) *Sphere {

	var v1 Vector3
	var result *Sphere
	if optionalTarget == nil {
		result = NewSphere(nil, 0)
	} else {
		result = optionalTarget
	}

	result.Center = *this.Center(nil)
	result.Radius = this.Size(&v1).Length() * 0.5

	return result
}

func (this *Box3) Intersect(box *Box3) *Box3 {

	this.Min.Max(&box.Min)
	this.Max.Min(&box.Max)
	return this
}

func (this *Box3) Union(box *Box3) *Box3 {

	this.Min.Min(&box.Min)
	this.Max.Max(&box.Max)
	return this
}

func (this *Box3) ApplyMatrix4(matrix *Matrix4) *Box3 {

	points := []Vector3{
		Vector3{},
		Vector3{},
		Vector3{},
		Vector3{},
		Vector3{},
		Vector3{},
		Vector3{},
		Vector3{},
	}

	points[0].Set(this.Min.X, this.Min.Y, this.Min.Z).ApplyMatrix4(matrix) // 000
	points[1].Set(this.Min.X, this.Min.Y, this.Max.Z).ApplyMatrix4(matrix) // 001
	points[2].Set(this.Min.X, this.Max.Y, this.Min.Z).ApplyMatrix4(matrix) // 010
	points[3].Set(this.Min.X, this.Max.Y, this.Max.Z).ApplyMatrix4(matrix) // 011
	points[4].Set(this.Max.X, this.Min.Y, this.Min.Z).ApplyMatrix4(matrix) // 100
	points[5].Set(this.Max.X, this.Min.Y, this.Max.Z).ApplyMatrix4(matrix) // 101
	points[6].Set(this.Max.X, this.Max.Y, this.Min.Z).ApplyMatrix4(matrix) // 110
	points[7].Set(this.Max.X, this.Max.Y, this.Max.Z).ApplyMatrix4(matrix) // 111

	this.MakeEmpty()
	this.SetFromPoints(points)

	return this
}

func (this *Box3) Translate(offset *Vector3) *Box3 {

	this.Min.Add(offset)
	this.Max.Add(offset)
	return this
}

func (this *Box3) Equals(box *Box3) bool {

	return box.Min.Equals(&this.Min) && box.Max.Equals(&this.Max)
}

func (this *Box3) Clone() *Box3 {

	return NewBox3(&this.Min, &this.Max)
}
