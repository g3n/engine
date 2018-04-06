// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// Box2 represents a 2D bounding box defined by two points:
// the point with minimum coordinates and the point with maximum coordinates.
type Box2 struct {
	min Vector2
	max Vector2
}

// NewBox2 creates and returns a pointer to a new Box2 defined
// by its minimum and maximum coordinates.
func NewBox2(min, max *Vector2) *Box2 {

	b := new(Box2)
	b.Set(min, max)
	return b
}

// Set sets this bounding box minimum and maximum coordinates.
// Returns pointer to this updated bounding box.
func (b *Box2) Set(min, max *Vector2) *Box2 {

	if min != nil {
		b.min = *min
	} else {
		b.min.Set(Infinity, Infinity)
	}
	if max != nil {
		b.max = *max
	} else {
		b.max.Set(-Infinity, -Infinity)
	}
	return b
}

// SetFromPoints set this bounding box from the specified array of points.
// Returns pointer to this updated bounding box.
func (b *Box2) SetFromPoints(points []*Vector2) *Box2 {

	b.MakeEmpty()
	for i := 0; i < len(points); i++ {
		b.ExpandByPoint(points[i])
	}
	return b
}

// SetFromCenterAndSize set this bounding box from a center point and size.
// Size is a vector from the minimum point to the maximum point.
// Returns pointer to this updated bounding box.
func (b *Box2) SetFromCenterAndSize(center, size *Vector2) *Box2 {

	var v1 Vector2
	halfSize := v1.Copy(size).MultiplyScalar(0.5)
	b.min.Copy(center).Sub(halfSize)
	b.max.Copy(center).Add(halfSize)
	return b
}

// Copy copy other to this bounding box.
// Returns pointer to this updated bounding box.
func (b *Box2) Copy(box *Box2) *Box2 {

	b.min = box.min
	b.max = box.max
	return b
}

// MakeEmpty set this bounding box to empty.
// Returns pointer to this updated bounding box.
func (b *Box2) MakeEmpty() *Box2 {

	b.min.X = Infinity
	b.min.Y = Infinity
	b.max.X = -Infinity
	b.max.Y = -Infinity
	return b
}

// Empty returns if this bounding box is empty.
func (b *Box2) Empty() bool {

	return (b.max.X < b.min.X) || (b.max.Y < b.min.Y)
}

// Center calculates the center point of this bounding box and
// stores its pointer to optionalTarget, if not nil, and also returns it.
func (b *Box2) Center(optionalTarget *Vector2) *Vector2 {

	var result *Vector2
	if optionalTarget == nil {
		result = NewVector2(0, 0)
	} else {
		result = optionalTarget
	}
	return result.AddVectors(&b.min, &b.max).MultiplyScalar(0.5)
}

// Size calculates the size of this bounding box: the vector  from
// its minimum point to its maximum point.
// Store pointer to the calculated size into optionalTarget, if not nil,
// and also returns it.
func (b *Box2) Size(optionalTarget *Vector2) *Vector2 {

	var result *Vector2
	if optionalTarget == nil {
		result = NewVector2(0, 0)
	} else {
		result = optionalTarget
	}
	return result.SubVectors(&b.min, &b.max)
}

// ExpandByPoint may expand this bounding box to include the specified point.
// Returns pointer to this updated bounding box.
func (b *Box2) ExpandByPoint(point *Vector2) *Box2 {

	b.min.Min(point)
	b.max.Max(point)
	return b
}

// ExpandByVector expands this bounding box by the specified vector.
// Returns pointer to this updated bounding box.
func (b *Box2) ExpandByVector(vector *Vector2) *Box2 {

	b.min.Sub(vector)
	b.max.Add(vector)
	return b
}

// ExpandByScalar expands this bounding box by the specified scalar.
// Returns pointer to this updated bounding box.
func (b *Box2) ExpandByScalar(scalar float32) *Box2 {

	b.min.AddScalar(-scalar)
	b.max.AddScalar(scalar)
	return b
}

// ContainsPoint returns if this bounding box contains the specified point.
func (b *Box2) ContainsPoint(point *Vector2) bool {

	if point.X < b.min.X || point.X > b.max.X ||
		point.Y < b.min.Y || point.Y > b.max.Y {
		return false
	}
	return true
}

// ContainsBox returns if this bounding box contains other box.
func (b *Box2) ContainsBox(other *Box2) bool {

	if (b.min.X <= other.min.X) && (other.max.X <= b.max.X) &&
		(b.min.Y <= other.min.Y) && (other.max.Y <= b.max.Y) {
		return true

	}
	return false
}

// IsIntersectionBox returns if other box intersects this one.
func (b *Box2) IsIntersectionBox(other *Box2) bool {

	// using 6 splitting planes to rule out intersections.
	if other.max.X < b.min.X || other.min.X > b.max.X ||
		other.max.Y < b.min.Y || other.min.Y > b.max.Y {
		return false
	}
	return true
}

// ClampPoint calculates a new point which is the specified point clamped inside this box.
// Stores the pointer to this new point into optionaTarget, if not nil, and also returns it.
func (b *Box2) ClampPoint(point *Vector2, optionalTarget *Vector2) *Vector2 {

	var result *Vector2
	if optionalTarget == nil {
		result = NewVector2(0, 0)
	} else {
		result = optionalTarget
	}
	return result.Copy(point).Clamp(&b.min, &b.max)
}

// DistanceToPoint returns the distance from this box to the specified point.
func (b *Box2) DistanceToPoint(point *Vector2) float32 {

	v1 := NewVector2(0, 0)
	clampedPoint := v1.Copy(point).Clamp(&b.min, &b.max)
	return clampedPoint.Sub(point).Length()
}

// Intersect sets this box to the intersection with other box.
// Returns pointer to this updated bounding box.
func (b *Box2) Intersect(other *Box2) *Box2 {

	b.min.Max(&other.min)
	b.max.Min(&other.max)
	return b
}

// Union set this box to the union with other box.
// Returns pointer to this updated bounding box.
func (b *Box2) Union(other *Box2) *Box2 {

	b.min.Min(&other.min)
	b.max.Max(&other.max)
	return b
}

// Translate translates the position of this box by offset.
// Returns pointer to this updated box.
func (b *Box2) Translate(offset *Vector2) *Box2 {

	b.min.Add(offset)
	b.max.Add(offset)
	return b
}

// Equals returns if this box is equal to other
func (b *Box2) Equals(other *Box2) bool {

	return other.min.Equals(&b.min) && other.max.Equals(&b.max)
}
