// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

type Line3 struct {
	start Vector3
	end   Vector3
}

func NewLine3(start, end *Vector3) *Line3 {

	this := new(Line3)
	this.Set(start, end)
	return this
}

func (this *Line3) Set(start, end *Vector3) *Line3 {

	if start != nil {
		this.start = *start
	}
	if end != nil {
		this.end = *end
	}
	return this
}

func (this *Line3) Copy(line *Line3) *Line3 {

	this.start = line.start
	this.end = line.end
	return this
}

func (this *Line3) Center(optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return result.AddVectors(&this.start, &this.end).MultiplyScalar(0.5)
}

func (this *Line3) Delta(optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return result.SubVectors(&this.end, &this.start).MultiplyScalar(0.5)
}

func (this *Line3) DistanceSq() float32 {

	return this.start.DistanceToSquared(&this.end)
}

func (this *Line3) Distance() float32 {

	return this.start.DistanceTo(&this.end)
}

func (this *Line3) At(t float32, optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return this.Delta(result).MultiplyScalar(t).Add(&this.start)
}

func (this *Line3) ClosestPointToPointParameter() func(*Vector3, bool) float32 {

	startP := NewVector3(0, 0, 0)
	startEnd := NewVector3(0, 0, 0)

	return func(point *Vector3, clampToLine bool) float32 {
		startP.SubVectors(point, &this.start)
		startEnd.SubVectors(&this.end, &this.start)

		startEnd2 := startEnd.Dot(startEnd)
		startEnd_startP := startEnd.Dot(startP)

		t := startEnd_startP / startEnd2
		if clampToLine {
			t = Clamp(t, 0, 1)
		}
		return t
	}
}

func (this *Line3) ClosestPointToPoint(point *Vector3, clampToLine bool, optionalTarget *Vector3) *Vector3 {

	t := this.ClosestPointToPointParameter()(point, clampToLine)
	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return this.Delta(result).MultiplyScalar(t).Add(&this.start)
}

func (this *Line3) ApplyMatrix4(matrix *Matrix4) *Line3 {

	this.start.ApplyMatrix4(matrix)
	this.end.ApplyMatrix4(matrix)

	return this
}

func (this *Line3) Equals(line *Line3) bool {

	return line.start.Equals(&this.start) && line.end.Equals(&this.end)
}

func (this *Line3) Clone() *Line3 {

	return NewLine3(&this.start, &this.end)
}
