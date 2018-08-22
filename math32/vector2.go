// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// Vector2 is a 2D vector/point with X and Y components.
type Vector2 struct {
	X float32
	Y float32
}

// NewVector2 creates and returns a pointer to a new Vector2 with
// the specified x and y components
func NewVector2(x, y float32) *Vector2 {

	return &Vector2{X: x, Y: y}
}

// NewVec2 creates and returns a pointer to a new zero-ed Vector2.
func NewVec2() *Vector2 {

	return &Vector2{X: 0, Y: 0}
}

// Set sets this vector X and Y components.
// Returns the pointer to this updated vector.
func (v *Vector2) Set(x, y float32) *Vector2 {

	v.X = x
	v.Y = y
	return v
}

// SetX sets this vector X component.
// Returns the pointer to this updated Vector.
func (v *Vector2) SetX(x float32) *Vector2 {

	v.X = x
	return v
}

// SetY sets this vector Y component.
// Returns the pointer to this updated vector.
func (v *Vector2) SetY(y float32) *Vector2 {

	v.Y = y
	return v
}

// SetComponent sets this vector component value by its index: 0 for X, 1 for Y.
// Returns the pointer to this updated vector
func (v *Vector2) SetComponent(index int, value float32) *Vector2 {

	switch index {
	case 0:
		v.X = value
	case 1:
		v.Y = value
	default:
		panic("index is out of range")
	}
	return v
}

// Component returns this vector component by its index: 0 for X, 1 for Y
func (v *Vector2) Component(index int) float32 {

	switch index {
	case 0:
		return v.X
	case 1:
		return v.Y
	default:
		panic("index is out of range")
	}
}

// SetByName sets this vector component value by its case insensitive name: "x" or "y".
func (v *Vector2) SetByName(name string, value float32) {

	switch name {
	case "x", "X":
		v.X = value
	case "y", "Y":
		v.Y = value
	default:
		panic("Invalid Vector2 component name: " + name)
	}
}

// Zero sets this vector X and Y components to be zero.
// Returns the pointer to this updated vector.
func (v *Vector2) Zero() *Vector2 {

	v.X = 0
	v.Y = 0
	return v
}

// Copy copies other vector to this one.
// It is equivalent to: *v = *other.
// Returns the pointer to this updated vector.
func (v *Vector2) Copy(other *Vector2) *Vector2 {

	v.X = other.X
	v.Y = other.Y
	return v
}

// Add adds other vector to this one.
// Returns the pointer to this updated vector.
func (v *Vector2) Add(other *Vector2) *Vector2 {

	v.X += other.X
	v.Y += other.Y
	return v
}

// AddScalar adds scalar s to each component of this vector.
// Returns the pointer to this updated vector.
func (v *Vector2) AddScalar(s float32) *Vector2 {

	v.X += s
	v.Y += s
	return v
}

// AddVectors adds vectors a and b to this one.
// Returns the pointer to this updated vector.
func (v *Vector2) AddVectors(a, b *Vector2) *Vector2 {

	v.X = a.X + b.X
	v.Y = a.Y + b.Y
	return v
}

// Sub subtracts other vector from this one.
// Returns the pointer to this updated vector.
func (v *Vector2) Sub(other *Vector2) *Vector2 {

	v.X -= other.X
	v.Y -= other.Y
	return v
}

// SubScalar subtracts scalar s from each component of this vector.
// Returns the pointer to this updated vector.
func (v *Vector2) SubScalar(s float32) *Vector2 {

	v.X -= s
	v.Y -= s
	return v
}

// SubVectors sets this vector to a - b.
// Returns the pointer to this updated vector.
func (v *Vector2) SubVectors(a, b *Vector2) *Vector2 {

	v.X = a.X - b.X
	v.Y = a.Y - b.Y
	return v
}

// Multiply multiplies each component of this vector by the corresponding one from other vector.
// Returns the pointer to this updated vector.
func (v *Vector2) Multiply(other *Vector2) *Vector2 {

	v.X *= other.X
	v.Y *= other.Y
	return v
}

// MultiplyScalar multiplies each component of this vector by the scalar s.
// Returns the pointer to this updated vector.
func (v *Vector2) MultiplyScalar(s float32) *Vector2 {

	v.X *= s
	v.Y *= s
	return v
}

// Divide divides each component of this vector by the corresponding one from other vector.
// Returns the pointer to this updated vector
func (v *Vector2) Divide(other *Vector2) *Vector2 {

	v.X /= other.X
	v.Y /= other.Y
	return v
}

// DivideScalar divides each component of this vector by the scalar s.
// If scalar is zero, sets this vector to zero.
// Returns the pointer to this updated vector.
func (v *Vector2) DivideScalar(scalar float32) *Vector2 {

	if scalar != 0 {
		invScalar := 1 / scalar
		v.X *= invScalar
		v.Y *= invScalar
	} else {
		v.X = 0
		v.Y = 0
	}
	return v
}

// Min sets this vector components to the minimum values of itself and other vector.
// Returns the pointer to this updated vector.
func (v *Vector2) Min(other *Vector2) *Vector2 {

	if v.X > other.X {
		v.X = other.X
	}
	if v.Y > other.Y {
		v.Y = other.Y
	}
	return v
}

// Max sets this vector components to the maximum value of itself and other vector.
// Returns the pointer to this updated vector.
func (v *Vector2) Max(other *Vector2) *Vector2 {

	if v.X < other.X {
		v.X = other.X
	}
	if v.Y < other.Y {
		v.Y = other.Y
	}
	return v
}

// Clamp sets this vector components to be no less than the corresponding components of min
// and not greater than the corresponding components of max.
// Assumes min < max, if this assumption isn't true it will not operate correctly.
// Returns the pointer to this updated vector.
func (v *Vector2) Clamp(min, max *Vector2) *Vector2 {

	if v.X < min.X {
		v.X = min.X
	} else if v.X > max.X {
		v.X = max.X
	}

	if v.Y < min.Y {
		v.Y = min.Y
	} else if v.Y > max.Y {
		v.Y = max.Y
	}
	return v
}

// ClampScalar sets this vector components to be no less than minVal and not greater than maxVal.
// Returns the pointer to this updated vector.
func (v *Vector2) ClampScalar(minVal, maxVal float32) *Vector2 {

	if v.X < minVal {
		v.X = minVal
	} else if v.X > maxVal {
		v.X = maxVal
	}

	if v.Y < minVal {
		v.Y = minVal
	} else if v.Y > maxVal {
		v.Y = maxVal
	}
	return v
}

// Floor applies math32.Floor() to each of this vector's components.
// Returns the pointer to this updated vector.
func (v *Vector2) Floor() *Vector2 {

	v.X = Floor(v.X)
	v.Y = Floor(v.Y)
	return v
}

// Ceil applies math32.Ceil() to each of this vector's components.
// Returns the pointer to this updated vector.
func (v *Vector2) Ceil() *Vector2 {

	v.X = Ceil(v.X)
	v.Y = Ceil(v.Y)
	return v
}

// Round rounds each of this vector's components.
// Returns the pointer to this updated vector.
func (v *Vector2) Round() *Vector2 {

	v.X = Floor(v.X + 0.5)
	v.Y = Floor(v.Y + 0.5)
	return v
}

// Negate negates each of this vector's components.
// Returns the pointer to this updated vector.
func (v *Vector2) Negate() *Vector2 {

	v.X = -v.X
	v.Y = -v.Y
	return v
}

// Dot returns the dot product of this vector with other.
// None of the vectors are changed.
func (v *Vector2) Dot(other *Vector2) float32 {

	return v.X*other.X + v.Y*other.Y
}

// LengthSq returns the length squared of this vector.
// LengthSq can be used to compare vectors' lengths without the need to perform a square root.
func (v *Vector2) LengthSq() float32 {

	return v.X*v.X + v.Y*v.Y
}

// Length returns the length of this vector.
func (v *Vector2) Length() float32 {

	return Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normalize normalizes this vector so its length will be 1.
// Returns the pointer to this updated vector.
func (v *Vector2) Normalize() *Vector2 {

	return v.DivideScalar(v.Length())
}

// DistanceTo returns the distance of this point to other.
func (v *Vector2) DistanceTo(other *Vector2) float32 {

	return Sqrt(v.DistanceToSquared(other))
}

// DistanceToSquared returns the distance squared of this point to other.
func (v *Vector2) DistanceToSquared(other *Vector2) float32 {

	dx := v.X - other.X
	dy := v.Y - other.Y
	return dx*dx + dy*dy
}

// SetLength sets this vector to have the specified length.
// Returns the pointer to this updated vector.
func (v *Vector2) SetLength(l float32) *Vector2 {

	oldLength := v.Length()
	if oldLength != 0 && l != oldLength {
		v.MultiplyScalar(l / oldLength)
	}
	return v
}

// Lerp sets each of this vector's components to the linear interpolated value of
// alpha between ifself and the corresponding other component.
// Returns the pointer to this updated vector.
func (v *Vector2) Lerp(other *Vector2, alpha float32) *Vector2 {

	v.X += (other.X - v.X) * alpha
	v.Y += (other.Y - v.Y) * alpha
	return v
}

// Equals returns if this vector is equal to other.
func (v *Vector2) Equals(other *Vector2) bool {

	return (other.X == v.X) && (other.Y == v.Y)
}

// FromArray sets this vector's components from the specified array and offset
// Returns the pointer to this updated vector.
func (v *Vector2) FromArray(array []float32, offset int) *Vector2 {

	v.X = array[offset]
	v.Y = array[offset+1]
	return v
}

// ToArray copies this vector's components to array starting at offset.
// Returns the array.
func (v *Vector2) ToArray(array []float32, offset int) []float32 {

	array[offset] = v.X
	array[offset+1] = v.Y
	return array
}

// InTriangle returns whether the vector is inside the specified triangle.
func (v *Vector2) InTriangle(p0, p1, p2 *Vector2) bool {

	A := 0.5 * (-p1.Y*p2.X + p0.Y*(-p1.X+p2.X) + p0.X*(p1.Y-p2.Y) + p1.X*p2.Y)
	sign := float32(1)
	if A < 0 {
		sign = float32(-1)
	}
	s := (p0.Y*p2.X - p0.X*p2.Y + (p2.Y-p0.Y)*v.X + (p0.X-p2.X)*v.Y) * sign
	t := (p0.X*p1.Y - p0.Y*p1.X + (p0.Y-p1.Y)*v.X + (p1.X-p0.X)*v.Y) * sign

	return s >= 0 && t >= 0 && (s+t) < 2*A*sign
}
