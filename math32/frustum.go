// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// Frustum represents a frustum
type Frustum struct {
	planes []Plane
}

// NewFrustumFromMatrix creates and returns a Frustum based on the provided matrix
func NewFrustumFromMatrix(m *Matrix4) *Frustum {
	f := new(Frustum)
	f.planes = make([]Plane, 6)
	f.SetFromMatrix(m)
	return f
}

// NewFrustum returns a pointer to a new Frustum object
func NewFrustum(p0, p1, p2, p3, p4, p5 *Plane) *Frustum {

	f := new(Frustum)
	f.planes = make([]Plane, 6)
	f.Set(p0, p1, p2, p3, p4, p5)
	return f
}

// Set sets the frustum's planes
func (f *Frustum) Set(p0, p1, p2, p3, p4, p5 *Plane) *Frustum {

	if p0 != nil {
		f.planes[0] = *p0
	}
	if p1 != nil {
		f.planes[1] = *p1
	}
	if p2 != nil {
		f.planes[2] = *p2
	}
	if p3 != nil {
		f.planes[3] = *p3
	}
	if p4 != nil {
		f.planes[4] = *p4
	}
	if p5 != nil {
		f.planes[5] = *p5
	}
	return f
}

// Copy modifies the receiver frustum to match the provided frustum
func (f *Frustum) Copy(frustum *Frustum) *Frustum {

	for i := 0; i < 6; i++ {
		f.planes[i] = frustum.planes[i]
	}
	return f
}

// SetFromMatrix sets the frustum's planes based on the specified Matrix4
func (f *Frustum) SetFromMatrix(m *Matrix4) *Frustum {

	planes := f.planes
	me0 := m[0]
	me1 := m[1]
	me2 := m[2]
	me3 := m[3]
	me4 := m[4]
	me5 := m[5]
	me6 := m[6]
	me7 := m[7]
	me8 := m[8]
	me9 := m[9]
	me10 := m[10]
	me11 := m[11]
	me12 := m[12]
	me13 := m[13]
	me14 := m[14]
	me15 := m[15]

	planes[0].SetComponents(me3-me0, me7-me4, me11-me8, me15-me12).Normalize()
	planes[1].SetComponents(me3+me0, me7+me4, me11+me8, me15+me12).Normalize()
	planes[2].SetComponents(me3+me1, me7+me5, me11+me9, me15+me13).Normalize()
	planes[3].SetComponents(me3-me1, me7-me5, me11-me9, me15-me13).Normalize()
	planes[4].SetComponents(me3-me2, me7-me6, me11-me10, me15-me14).Normalize()
	planes[5].SetComponents(me3+me2, me7+me6, me11+me10, me15+me14).Normalize()

	return f
}

/**
SHOULD NOT DEPEND on core package (Move to core ?)
func (this *Frustum) IntersectsObject(geometry *core.Geometry) bool {


    return false
}
*/

// IntersectsSphere determines whether the specified sphere is intersecting the frustum
func (f *Frustum) IntersectsSphere(sphere *Sphere) bool {

	planes := f.planes
	negRadius := -sphere.Radius

	for i := 0; i < 6; i++ {
		distance := planes[i].DistanceToPoint(&sphere.Center)

		if distance < negRadius {
			return false
		}
	}
	return true
}

// IntersectsBox determines whether the specified box is intersecting the frustum
func (f *Frustum) IntersectsBox(box *Box3) bool {

	var p1 Vector3
	var p2 Vector3

	for i := 0; i < 6; i++ {
		plane := &f.planes[i]
		if plane.normal.X > 0 {
			p1.X = box.Min.X
		} else {
			p1.X = box.Max.X
		}
		if plane.normal.X > 0 {
			p2.X = box.Max.X
		} else {
			p2.X = box.Min.X
		}
		if plane.normal.Y > 0 {
			p1.Y = box.Min.Y
		} else {
			p1.Y = box.Max.Y
		}
		if plane.normal.Y > 0 {
			p2.Y = box.Max.Y
		} else {
			p2.Y = box.Min.Y
		}
		if plane.normal.Z > 0 {
			p1.Z = box.Min.Z
		} else {
			p1.Z = box.Max.Z
		}
		if plane.normal.Z > 0 {
			p2.Z = box.Max.Z
		} else {
			p2.Z = box.Min.Z
		}

		d1 := plane.DistanceToPoint(&p1)
		d2 := plane.DistanceToPoint(&p2)

		// if both outside plane, no intersection

		if d1 < 0 && d2 < 0 {
			return false
		}

	}
	return true
}

// ContainsPoint determines whether the frustum contains the specified point
func (f *Frustum) ContainsPoint(point *Vector3) bool {

	for i := 0; i < 6; i++ {
		if f.planes[i].DistanceToPoint(point) < 0 {
			return false
		}
	}
	return true
}

// Clone returns a pointer to a new Frustum object with the same planes as the original
func (f *Frustum) Clone() *Frustum {

	return NewFrustum(nil, nil, nil, nil, nil, nil).Copy(f)
}
