// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

type Box2 struct {
	min Vector2
	max Vector2
}

func NewBox2(min, max *Vector2) *Box2 {

	this := new(Box2)
	this.Set(min, max)
	return this
}

func (this *Box2) Set(min, max *Vector2) *Box2 {

	if min != nil {
		this.min = *min
	} else {
		this.min.Set(Infinity, Infinity)
	}
	if max != nil {
		this.max = *max
	} else {
		this.max.Set(-Infinity, -Infinity)
	}
	return this
}

func (this *Box2) SetFromPoints(points []*Vector2) *Box2 {

	this.MakeEmpty()
	for i := 0; i < len(points); i++ {
		this.ExpandByPoint(points[i])
	}
	return this
}

func (this *Box2) SetFromCenterAndSize(center, size *Vector2) *Box2 {

	var v1 Vector2
	halfSize := v1.Copy(size).MultiplyScalar(0.5)
	this.min.Copy(center).Sub(halfSize)
	this.max.Copy(center).Add(halfSize)
	return this
}

func (this *Box2) Copy(box *Box2) *Box2 {

	this.min = box.min
	this.max = box.max
	return this
}

func (this *Box2) MakeEmpty() *Box2 {

	this.min.X = Infinity
	this.min.Y = Infinity
	this.max.X = -Infinity
	this.max.Y = -Infinity
	return this
}

func (this *Box2) Empty() bool {

	return (this.max.X < this.min.X) || (this.max.Y < this.min.Y)
}

func (this *Box2) Center(optionalTarget *Vector2) *Vector2 {

	var result *Vector2
	if optionalTarget == nil {
		result = NewVector2(0, 0)
	} else {
		result = optionalTarget
	}
	return result.AddVectors(&this.min, &this.max).MultiplyScalar(0.5)
}

func (this *Box2) Size(optionalTarget *Vector2) *Vector2 {

	var result *Vector2
	if optionalTarget == nil {
		result = NewVector2(0, 0)
	} else {
		result = optionalTarget
	}
	return result.SubVectors(&this.min, &this.max)
}

func (this *Box2) ExpandByPoint(point *Vector2) *Box2 {

	this.min.Min(point)
	this.max.Max(point)
	return this
}

func (this *Box2) ExpandByVector(vector *Vector2) *Box2 {

	this.min.Sub(vector)
	this.max.Add(vector)
	return this
}

func (this *Box2) ExpandByScalar(scalar float32) *Box2 {

	this.min.AddScalar(-scalar)
	this.max.AddScalar(scalar)
	return this
}

func (this *Box2) ContainsPoint(point *Vector2) bool {

	if point.X < this.min.X || point.X > this.max.X ||
		point.Y < this.min.Y || point.Y > this.max.Y {
		return false
	}
	return true
}

func (this *Box2) ContainsBox(box *Box2) bool {

	if (this.min.X <= box.min.X) && (box.max.X <= this.max.X) &&
		(this.min.Y <= box.min.Y) && (box.max.Y <= this.max.Y) {
		return true

	}
	return false
}

func (this *Box2) GetParameter(point *Vector2, optionalTarget *Vector2) *Vector2 {

	var result *Vector2
	if optionalTarget == nil {
		result = NewVector2(0, 0)
	} else {
		result = optionalTarget
	}
	return result.Set(
		(point.X-this.min.X)/(this.max.X-this.min.X),
		(point.Y-this.min.Y)/(this.max.Y-this.min.Y),
	)
}

func (this *Box2) IsIntersectionBox(box *Box2) bool {

	// using 6 splitting planes to rule out intersections.
	if box.max.X < this.min.X || box.min.X > this.max.X ||
		box.max.Y < this.min.Y || box.min.Y > this.max.Y {
		return false
	}
	return true
}

func (this *Box2) ClampPoint(point *Vector2, optionalTarget *Vector2) *Vector2 {

	var result *Vector2
	if optionalTarget == nil {
		result = NewVector2(0, 0)
	} else {
		result = optionalTarget
	}
	return result.Copy(point).Clamp(&this.min, &this.max)
}

func (this *Box2) DistanceToPoint(point *Vector2) float32 {

	v1 := NewVector2(0, 0)
	clampedPoint := v1.Copy(point).Clamp(&this.min, &this.max)
	return clampedPoint.Sub(point).Length()
}

func (this *Box2) Intersect(box *Box2) *Box2 {

	this.min.Max(&box.min)
	this.max.Min(&box.max)
	return this
}

func (this *Box2) Union(box *Box2) *Box2 {

	this.min.Min(&box.min)
	this.max.Max(&box.max)
	return this
}

func (this *Box2) Translate(offset *Vector2) *Box2 {

	this.min.Add(offset)
	this.max.Add(offset)
	return this
}

func (this *Box2) Equals(box *Box2) bool {

	return box.min.Equals(&this.min) && box.max.Equals(&this.max)
}
