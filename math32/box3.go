// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// Box3 represents a 3D bounding box defined by two points:
// the point with minimum coordinates and the point with maximum coordinates.
type Box3 struct {
	Min Vector3
	Max Vector3
}

// NewBox3 creates and returns a pointer to a new Box3 defined
// by its minimum and maximum coordinates.
func NewBox3(min, max *Vector3) *Box3 {

	b := new(Box3)
	b.Set(min, max)
	return b
}

// Set sets this bounding box minimum and maximum coordinates.
// Returns pointer to this updated bounding box.
func (b *Box3) Set(min, max *Vector3) *Box3 {

	if min != nil {
		b.Min = *min
	} else {
		b.Min.Set(Infinity, Infinity, Infinity)
	}
	if max != nil {
		b.Max = *max
	} else {
		b.Max.Set(-Infinity, -Infinity, -Infinity)
	}
	return b
}

// SetFromPoints set this bounding box from the specified array of points.
// Returns pointer to this updated bounding box.
func (b *Box3) SetFromPoints(points []Vector3) *Box3 {

	b.MakeEmpty()
	for i := 0; i < len(points); i++ {
		b.ExpandByPoint(&points[i])
	}
	return b
}

// SetFromCenterAndSize set this bounding box from a center point and size.
// Size is a vector from the minimum point to the maximum point.
// Returns pointer to this updated bounding box.
func (b *Box3) SetFromCenterAndSize(center, size *Vector3) *Box3 {

	v1 := NewVector3(0, 0, 0)
	halfSize := v1.Copy(size).MultiplyScalar(0.5)
	b.Min.Copy(center).Sub(halfSize)
	b.Max.Copy(center).Add(halfSize)
	return b
}

// Copy copy other to this bounding box.
// Returns pointer to this updated bounding box.
func (b *Box3) Copy(other *Box3) *Box3 {

	b.Min = other.Min
	b.Max = other.Max
	return b
}

// MakeEmpty set this bounding box to empty.
// Returns pointer to this updated bounding box.
func (b *Box3) MakeEmpty() *Box3 {

	b.Min.X = Infinity
	b.Min.Y = Infinity
	b.Min.Z = Infinity
	b.Max.X = -Infinity
	b.Max.Y = -Infinity
	b.Max.Z = -Infinity
	return b
}

// Empty returns if this bounding box is empty.
func (b *Box3) Empty() bool {

	return (b.Max.X < b.Min.X) || (b.Max.Y < b.Min.Y) || (b.Max.Z < b.Min.Z)
}

// Center calculates the center point of this bounding box and
// stores its pointer to optionalTarget, if not nil, and also returns it.
func (b *Box3) Center(optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return result.AddVectors(&b.Min, &b.Max).MultiplyScalar(0.5)
}

// Size calculates the size of this bounding box: the vector  from
// its minimum point to its maximum point.
// Store pointer to the calculated size into optionalTarget, if not nil,
// and also returns it.
func (b *Box3) Size(optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return result.SubVectors(&b.Min, &b.Max)
}

// ExpandByPoint may expand this bounding box to include the specified point.
// Returns pointer to this updated bounding box.
func (b *Box3) ExpandByPoint(point *Vector3) *Box3 {

	b.Min.Min(point)
	b.Max.Max(point)
	return b
}

// ExpandByVector expands this bounding box by the specified vector.
// Returns pointer to this updated bounding box.
func (b *Box3) ExpandByVector(vector *Vector3) *Box3 {

	b.Min.Sub(vector)
	b.Max.Add(vector)
	return b
}

// ExpandByScalar expands this bounding box by the specified scalar.
// Returns pointer to this updated bounding box.
func (b *Box3) ExpandByScalar(scalar float32) *Box3 {

	b.Min.AddScalar(-scalar)
	b.Max.AddScalar(scalar)
	return b
}

// ContainsPoint returns if this bounding box contains the specified point.
func (b *Box3) ContainsPoint(point *Vector3) bool {

	if point.X < b.Min.X || point.X > b.Max.X ||
		point.Y < b.Min.Y || point.Y > b.Max.Y ||
		point.Z < b.Min.Z || point.Z > b.Max.Z {
		return false
	}
	return true
}

// ContainsBox returns if this bounding box contains other box.
func (b *Box3) ContainsBox(box *Box3) bool {

	if (b.Min.X <= box.Max.X) && (box.Max.X <= b.Max.X) &&
		(b.Min.Y <= box.Min.Y) && (box.Max.Y <= b.Max.Y) &&
		(b.Min.Z <= box.Min.Z) && (box.Max.Z <= b.Max.Z) {
		return true

	}
	return false
}

// IsIntersectionBox returns if other box intersects this one.
func (b *Box3) IsIntersectionBox(other *Box3) bool {

	// using 6 splitting planes to rule out intersections.
	if other.Max.X < b.Min.X || other.Min.X > b.Max.X ||
		other.Max.Y < b.Min.Y || other.Min.Y > b.Max.Y ||
		other.Max.Z < b.Min.Z || other.Min.Z > b.Max.Z {
		return false
	}
	return true
}

// ClampPoint calculates a new point which is the specified point clamped inside this box.
// Stores the pointer to this new point into optionaTarget, if not nil, and also returns it.
func (b *Box3) ClampPoint(point *Vector3, optionalTarget *Vector3) *Vector3 {

	var result *Vector3
	if optionalTarget == nil {
		result = NewVector3(0, 0, 0)
	} else {
		result = optionalTarget
	}
	return result.Copy(point).Clamp(&b.Min, &b.Max)
}

// DistanceToPoint returns the distance from this box to the specified point.
func (b *Box3) DistanceToPoint(point *Vector3) float32 {

	var v1 Vector3
	clampedPoint := v1.Copy(point).Clamp(&b.Min, &b.Max)
	return clampedPoint.Sub(point).Length()
}

// GetBoundingSphere creates a bounding sphere to this bounding box.
// Store its pointer into optionalTarget, if not nil, and also returns it.
func (b *Box3) GetBoundingSphere(optionalTarget *Sphere) *Sphere {

	var v1 Vector3
	var result *Sphere
	if optionalTarget == nil {
		result = NewSphere(nil, 0)
	} else {
		result = optionalTarget
	}

	result.Center = *b.Center(nil)
	result.Radius = b.Size(&v1).Length() * 0.5

	return result
}

// Intersect sets this box to the intersection with other box.
// Returns pointer to this updated bounding box.
func (b *Box3) Intersect(other *Box3) *Box3 {

	b.Min.Max(&other.Min)
	b.Max.Min(&other.Max)
	return b
}

// Union set this box to the union with other box.
// Returns pointer to this updated bounding box.
func (b *Box3) Union(other *Box3) *Box3 {

	b.Min.Min(&other.Min)
	b.Max.Max(&other.Max)
	return b
}

// ApplyMatrix4 applies the specified matrix to the vertices of this bounding box.
// Returns pointer to this updated bounding box.
func (b *Box3) ApplyMatrix4(m *Matrix4) *Box3 {

	xax := m[0] * b.Min.X
	xay := m[1] * b.Min.X
	xaz := m[2] * b.Min.X
	xbx := m[0] * b.Max.X
	xby := m[1] * b.Max.X
	xbz := m[2] * b.Max.X
	yax := m[4] * b.Min.Y
	yay := m[5] * b.Min.Y
	yaz := m[6] * b.Min.Y
	ybx := m[4] * b.Max.Y
	yby := m[5] * b.Max.Y
	ybz := m[6] * b.Max.Y
	zax := m[8] * b.Min.Z
	zay := m[9] * b.Min.Z
	zaz := m[10] * b.Min.Z
	zbx := m[8] * b.Max.Z
	zby := m[9] * b.Max.Z
	zbz := m[10] * b.Max.Z

	b.Min.X = Min(xax, xbx) + Min(yax, ybx) + Min(zax, zbx) + m[12]
	b.Min.Y = Min(xay, xby) + Min(yay, yby) + Min(zay, zby) + m[13]
	b.Min.Z = Min(xaz, xbz) + Min(yaz, ybz) + Min(zaz, zbz) + m[14]
	b.Max.X = Max(xax, xbx) + Max(yax, ybx) + Max(zax, zbx) + m[12]
	b.Max.Y = Max(xay, xby) + Max(yay, yby) + Max(zay, zby) + m[13]
	b.Max.Z = Max(xaz, xbz) + Max(yaz, ybz) + Max(zaz, zbz) + m[14]

	return b
}

// Translate translates the position of this box by offset.
// Returns pointer to this updated box.
func (b *Box3) Translate(offset *Vector3) *Box3 {

	b.Min.Add(offset)
	b.Max.Add(offset)
	return b
}

// Equals returns if this box is equal to other
func (b *Box3) Equals(other *Box3) bool {

	return other.Min.Equals(&b.Min) && other.Max.Equals(&b.Max)
}

// Clone creates and returns a pointer to copy of this bounding box
func (b *Box3) Clone() *Box3 {

	return NewBox3(&b.Min, &b.Max)
}
