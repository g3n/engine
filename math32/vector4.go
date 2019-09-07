// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// Vector4 is a vector/point in homogeneous coordinates with X, Y, Z and W components.
type Vector4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

// NewVector4 creates and returns a pointer to a new Vector4
func NewVector4(x, y, z, w float32) *Vector4 {

	return &Vector4{X: x, Y: y, Z: z, W: w}
}

// NewVec4 creates and returns a pointer to a new zero-ed Vector4 (with W=1).
func NewVec4() *Vector4 {

	return &Vector4{X: 0, Y: 0, Z: 0, W: 1}
}

// Set sets this vector X, Y, Z and W components.
// Returns the pointer to this updated vector.
func (v *Vector4) Set(x, y, z, w float32) *Vector4 {

	v.X = x
	v.Y = y
	v.Z = z
	v.W = w
	return v
}

// SetVector3 sets this vector from another Vector3 and W
func (v *Vector4) SetVector3(other *Vector3, w float32) *Vector4 {

	v.X = other.X
	v.Y = other.Y
	v.Z = other.Z
	v.W = w
	return v
}

// SetX sets this vector X component.
// Returns the pointer to this updated Vector.
func (v *Vector4) SetX(x float32) *Vector4 {

	v.X = x
	return v
}

// SetY sets this vector Y component.
// Returns the pointer to this updated vector.
func (v *Vector4) SetY(y float32) *Vector4 {

	v.Y = y
	return v
}

// SetZ sets this vector Z component.
// Returns the pointer to this updated vector.
func (v *Vector4) SetZ(z float32) *Vector4 {

	v.Z = z
	return v
}

// SetW sets this vector W component.
// Returns the pointer to this updated vector.
func (v *Vector4) SetW(w float32) *Vector4 {

	v.W = w
	return v
}

// SetComponent sets this vector component value by its index: 0 for X, 1 for Y, 2 for Z, 3 for W.
// Returns the pointer to this updated vector
func (v *Vector4) SetComponent(index int, value float32) *Vector4 {

	switch index {
	case 0:
		v.X = value
	case 1:
		v.Y = value
	case 2:
		v.Z = value
	case 3:
		v.W = value
	default:
		panic("index is out of range")
	}
	return v
}

// Component returns this vector component by its index: 0 for X, 1 for Y, 2 for Z, 3 for W.
func (v *Vector4) Component(index int) float32 {

	switch index {
	case 0:
		return v.X
	case 1:
		return v.Y
	case 2:
		return v.Z
	case 3:
		return v.W
	default:
		panic("index is out of range")
	}
}

// SetByName sets this vector component value by its case insensitive name: "x", "y", "z" or "w".
func (v *Vector4) SetByName(name string, value float32) {

	switch name {
	case "x", "X":
		v.X = value
	case "y", "Y":
		v.Y = value
	case "z", "Z":
		v.Z = value
	case "w", "W":
		v.W = value
	default:
		panic("Invalid Vector4 component name: " + name)
	}
}

// Zero sets this vector X, Y and Z components to be zero and W to be one.
// Returns the pointer to this updated vector.
func (v *Vector4) Zero() *Vector4 {

	v.X = 0
	v.Y = 0
	v.Z = 0
	v.W = 1
	return v
}

// Copy copies other vector to this one.
// Returns the pointer to this updated vector.
func (v *Vector4) Copy(other *Vector4) *Vector4 {

	*v = *other
	return v
}

// Add adds other vector to this one.
// Returns the pointer to this updated vector.
func (v *Vector4) Add(other *Vector4) *Vector4 {

	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	v.W += other.W
	return v
}

// AddScalar adds scalar s to each component of this vector.
// Returns the pointer to this updated vector.
func (v *Vector4) AddScalar(s float32) *Vector4 {

	v.X += s
	v.Y += s
	v.Z += s
	v.W += s
	return v
}

// AddVectors adds vectors a and b to this one.
// Returns the pointer to this updated vector.
func (v *Vector4) AddVectors(a, b *Vector4) *Vector4 {

	v.X = a.X + b.X
	v.Y = a.Y + b.Y
	v.Z = a.Z + b.Z
	v.W = a.W + b.W
	return v
}

// Sub subtracts other vector from this one.
// Returns the pointer to this updated vector.
func (v *Vector4) Sub(other *Vector4) *Vector4 {

	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
	v.W -= other.W
	return v
}

// SubScalar subtracts scalar s from each component of this vector.
// Returns the pointer to this updated vector.
func (v *Vector4) SubScalar(s float32) *Vector4 {

	v.X -= s
	v.Y -= s
	v.Z -= s
	v.W -= s
	return v
}

// SubVectors sets this vector to a - b.
// Returns the pointer to this updated vector.
func (v *Vector4) SubVectors(a, b *Vector4) *Vector4 {

	v.X = a.X - b.X
	v.Y = a.Y - b.Y
	v.Z = a.Y - b.Z
	v.W = a.Y - b.W
	return v
}

// Multiply multiplies each component of this vector by the corresponding one from other vector.
// Returns the pointer to this updated vector.
func (v *Vector4) Multiply(other *Vector4) *Vector4 {

	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
	v.W *= other.W
	return v
}

// MultiplyScalar multiplies each component of this vector by the scalar s.
// Returns the pointer to this updated vector.
func (v *Vector4) MultiplyScalar(scalar float32) *Vector4 {

	v.X *= scalar
	v.Y *= scalar
	v.Z *= scalar
	v.W *= scalar
	return v
}

// Divide divides each component of this vector by the corresponding one from other vector.
// Returns the pointer to this updated vector
func (v *Vector4) Divide(other *Vector4) *Vector4 {

	v.X /= other.X
	v.Y /= other.Y
	v.Z /= other.Z
	v.W /= other.W
	return v
}

// DivideScalar divides each component of this vector by the scalar s.
// If scalar is zero, sets this vector to zero.
// Returns the pointer to this updated vector.
func (v *Vector4) DivideScalar(scalar float32) *Vector4 {

	if scalar != 0 {
		invScalar := 1 / scalar
		v.X *= invScalar
		v.Y *= invScalar
		v.Z *= invScalar
		v.W *= invScalar
	} else {
		v.X = 0
		v.Y = 0
		v.Z = 0
		v.W = 0
	}
	return v
}

// Min sets this vector components to the minimum values of itself and other vector.
// Returns the pointer to this updated vector.
func (v *Vector4) Min(other *Vector4) *Vector4 {

	if v.X > other.X {
		v.X = other.X
	}
	if v.Y > other.Y {
		v.Y = other.Y
	}
	if v.Z > other.Z {
		v.Z = other.Z
	}
	if v.W > other.W {
		v.W = other.W
	}
	return v
}

// Max sets this vector components to the maximum value of itself and other vector.
// Returns the pointer to this updated vector.
func (v *Vector4) Max(other *Vector4) *Vector4 {

	if v.X < other.X {
		v.X = other.X
	}
	if v.Y < other.Y {
		v.Y = other.Y
	}
	if v.Z < other.Z {
		v.Z = other.Z
	}
	if v.W < other.W {
		v.W = other.W
	}
	return v
}

// Clamp sets this vector components to be no less than the corresponding components of min
// and not greater than the corresponding component of max.
// Assumes min < max, if this assumption isn't true it will not operate correctly.
// Returns the pointer to this updated vector.
func (v *Vector4) Clamp(min, max *Vector4) *Vector4 {

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

	if v.Z < min.Z {
		v.Z = min.Z
	} else if v.Z > max.Z {
		v.Z = max.Z
	}

	if v.W < min.W {
		v.W = min.W
	} else if v.W > max.W {
		v.W = max.W
	}
	return v
}

// ClampScalar sets this vector components to be no less than minVal and not greater than maxVal.
// Returns the pointer to this updated vector.
func (v *Vector4) ClampScalar(minVal, maxVal float32) *Vector4 {

	min := NewVector4(minVal, minVal, minVal, minVal)
	max := NewVector4(maxVal, maxVal, maxVal, maxVal)
	return v.Clamp(min, max)
}

// Floor applies math32.Floor() to each of this vector's components.
// Returns the pointer to this updated vector.
func (v *Vector4) Floor() *Vector4 {

	v.X = Floor(v.X)
	v.Y = Floor(v.Y)
	v.Z = Floor(v.Z)
	v.W = Floor(v.W)
	return v
}

// Ceil applies math32.Ceil() to each of this vector's components.
// Returns the pointer to this updated vector.
func (v *Vector4) Ceil() *Vector4 {

	v.X = Ceil(v.X)
	v.Y = Ceil(v.Y)
	v.Z = Ceil(v.Z)
	v.W = Ceil(v.W)
	return v
}

// Round rounds each of this vector's components.
// Returns the pointer to this updated vector.
func (v *Vector4) Round() *Vector4 {

	v.X = Floor(v.X + 0.5)
	v.Y = Floor(v.Y + 0.5)
	v.Z = Floor(v.Z + 0.5)
	v.W = Floor(v.W + 0.5)
	return v
}

// Negate negates each of this vector's components.
// Returns the pointer to this updated vector.
func (v *Vector4) Negate() *Vector4 {

	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	v.W = -v.W
	return v
}

// Dot returns the dot product of this vector with other.
// None of the vectors are changed.
func (v *Vector4) Dot(other *Vector4) float32 {

	return v.X*other.X + v.Y*other.Y + v.Z*other.Z + v.W*other.W
}

// LengthSq returns the length squared of this vector.
// LengthSq can be used to compare vectors' lengths without the need to perform a square root.
func (v *Vector4) LengthSq() float32 {

	return v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W
}

// Length returns the length of this vector.
func (v *Vector4) Length() float32 {

	return Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W)
}

// Normalize normalizes this vector so its length will be 1.
// Returns the pointer to this updated vector.
func (v *Vector4) Normalize() *Vector4 {

	return v.DivideScalar(v.Length())
}

// SetLength sets this vector to have the specified length.
// If the current length is zero, does nothing.
// Returns the pointer to this updated vector.
func (v *Vector4) SetLength(l float32) *Vector4 {

	oldLength := v.Length()
	if oldLength != 0 && l != oldLength {
		v.MultiplyScalar(l / oldLength)
	}
	return v
}

// Lerp sets each of this vector's components to the linear interpolated value of
// alpha between ifself and the corresponding other component.
// Returns the pointer to this updated vector.
func (v *Vector4) Lerp(other *Vector4, alpha float32) *Vector4 {

	v.X += (other.X - v.X) * alpha
	v.Y += (other.Y - v.Y) * alpha
	v.Z += (other.Z - v.Z) * alpha
	v.W += (other.W - v.W) * alpha
	return v
}

// Equals returns if this vector is equal to other.
func (v *Vector4) Equals(other *Vector4) bool {

	return (other.X == v.X) && (other.Y == v.Y) && (other.Z == v.Z) && (other.W == v.W)
}

// FromArray sets this vector's components from the specified array and offset
// Returns the pointer to this updated vector.
func (v *Vector4) FromArray(array []float32, offset int) *Vector4 {

	v.X = array[offset]
	v.Y = array[offset+1]
	v.Z = array[offset+2]
	v.W = array[offset+3]
	return v
}

// ToArray copies this vector's components to array starting at offset.
// Returns the array.
func (v *Vector4) ToArray(array []float32, offset int) []float32 {

	array[offset] = v.X
	array[offset+1] = v.Y
	array[offset+2] = v.Z
	array[offset+3] = v.W
	return array
}

// ApplyMatrix4 multiplies the specified 4x4 matrix by this vector.
// Returns the pointer to this updated vector.
func (v *Vector4) ApplyMatrix4(m *Matrix4) *Vector4 {

	x := v.X
	y := v.Y
	z := v.Z
	w := v.W

	v.X = m[0]*x + m[4]*y + m[8]*z + m[12]*w
	v.Y = m[1]*x + m[5]*y + m[9]*z + m[13]*w
	v.Z = m[2]*x + m[6]*y + m[10]*z + m[14]*w
	v.W = m[3]*x + m[7]*y + m[11]*z + m[15]*w

	return v
}

// SetAxisAngleFromQuaternion set this vector to be the axis (x, y, z) and angle (w) of a rotation specified the quaternion q.
// Assumes q is normalized.
func (v *Vector4) SetAxisAngleFromQuaternion(q *Quaternion) *Vector4 {

	// http://www.euclideanspace.com/maths/geometry/rotations/conversions/quaternionToAngle/index.htm
	v.W = 2 * Acos(q.W)
	s := Sqrt(1 - q.W*q.W)
	if s < 0.0001 {
		v.X = 1
		v.Y = 0
		v.Z = 0
	} else {
		v.X = q.X / s
		v.Y = q.Y / s
		v.Z = q.Z / s
	}
	return v
}

// SetAxisFromRotationMatrix this vector to be the axis (x, y, z) and angle (w) of a rotation specified the matrix m.
// Assumes the upper 3x3 of m is a pure rotation matrix (i.e, unscaled).
func (v *Vector4) SetAxisFromRotationMatrix(m *Matrix4) *Vector4 {

	// http://www.euclideanspace.com/maths/geometry/rotations/conversions/matrixToAngle/index.htm
	var angle, x, y, z float32 // variables for result
	var epsilon float32 = 0.01 // margin to allow for rounding errors
	var epsilon2 float32 = 0.1 // margin to distinguish between 0 and 180 degrees

	m11 := m[0]
	m12 := m[4]
	m13 := m[8]
	m21 := m[1]
	m22 := m[5]
	m23 := m[9]
	m31 := m[2]
	m32 := m[6]
	m33 := m[10]

	if (Abs(m12-m21) < epsilon) && (Abs(m13-m31) < epsilon) && (Abs(m23-m32) < epsilon) {

		// singularity found
		// first check for identity matrix which must have +1 for all terms
		// in leading diagonal and zero in other terms

		if (Abs(m12+m21) < epsilon2) && (Abs(m13+m31) < epsilon2) && (Abs(m23+m32) < epsilon2) && (Abs(m11+m22+m33-3) < epsilon2) {

			// v singularity is identity matrix so angle = 0

			v.Set(1, 0, 0, 0)

			return v // zero angle, arbitrary axis
		}

		// otherwise this singularity is angle = 180

		angle = Pi

		var xx = (m11 + 1) / 2
		var yy = (m22 + 1) / 2
		var zz = (m33 + 1) / 2
		var xy = (m12 + m21) / 4
		var xz = (m13 + m31) / 4
		var yz = (m23 + m32) / 4

		if (xx > yy) && (xx > zz) { // m11 is the largest diagonal term

			if xx < epsilon {

				x = 0
				y = 0.707106781
				z = 0.707106781

			} else {

				x = Sqrt(xx)
				y = xy / x
				z = xz / x

			}

		} else if yy > zz { // m22 is the largest diagonal term

			if yy < epsilon {

				x = 0.707106781
				y = 0
				z = 0.707106781

			} else {

				y = Sqrt(yy)
				x = xy / y
				z = yz / y

			}

		} else { // m33 is the largest diagonal term so base result on this

			if zz < epsilon {

				x = 0.707106781
				y = 0.707106781
				z = 0

			} else {

				z = Sqrt(zz)
				x = xz / z
				y = yz / z

			}

		}

		v.Set(x, y, z, angle)

		return v // return 180 deg rotation
	}

	// as we have reached here there are no singularities so we can handle normally

	s := Sqrt((m32-m23)*(m32-m23) + (m13-m31)*(m13-m31) + (m21-m12)*(m21-m12)) // used to normalize

	if Abs(s) < 0.001 {
		s = 1
	}

	// prevent divide by zero, should not happen if matrix is orthogonal and should be
	// caught by singularity test above, but I've left it in just in case

	v.X = (m32 - m23) / s
	v.Y = (m13 - m31) / s
	v.Z = (m21 - m12) / s
	v.W = Acos((m11 + m22 + m33 - 1) / 2)

	return v
}

// Clone returns a copy of this vector
func (v *Vector4) Clone() *Vector4 {

	return NewVector4(v.X, v.Y, v.Z, v.W)
}
