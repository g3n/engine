// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

// Vector3 is a 3D vector/point with X, Y and Z components.
type Vector3 struct {
	X float32
	Y float32
	Z float32
}

// NewVector3 creates and returns a pointer to a new Vector3 with
// the specified x, y and y components
func NewVector3(x, y, z float32) *Vector3 {

	return &Vector3{X: x, Y: y, Z: z}
}

// NewVec3 creates and returns a pointer to a new zero-ed Vector3.
func NewVec3() *Vector3 {

	return &Vector3{X: 0, Y: 0, Z: 0}
}

// Set sets this vector X, Y and Z components.
// Returns the pointer to this updated vector.
func (v *Vector3) Set(x, y, z float32) *Vector3 {

	v.X = x
	v.Y = y
	v.Z = z
	return v
}

// SetX sets this vector X component.
// Returns the pointer to this updated Vector.
func (v *Vector3) SetX(x float32) *Vector3 {

	v.X = x
	return v
}

// SetY sets this vector Y component.
// Returns the pointer to this updated vector.
func (v *Vector3) SetY(y float32) *Vector3 {

	v.Y = y
	return v
}

// SetZ sets this vector Z component.
// Returns the pointer to this updated vector.
func (v *Vector3) SetZ(z float32) *Vector3 {

	v.Z = z
	return v
}

// SetComponent sets this vector component value by its index: 0 for X, 1 for Y, 2 for Z.
// Returns the pointer to this updated vector
func (v *Vector3) SetComponent(index int, value float32) {

	switch index {
	case 0:
		v.X = value
	case 1:
		v.Y = value
	case 2:
		v.Z = value
	default:
		panic("index is out of range: ")
	}
}

// Component returns this vector component by its index: 0 for X, 1 for Y, 2 for Z.
func (v *Vector3) Component(index int) float32 {

	switch index {
	case 0:
		return v.X
	case 1:
		return v.Y
	case 2:
		return v.Z
	default:
		panic("index is out of range")
	}
}

// SetByName sets this vector component value by its case insensitive name: "x", "y", or "z".
func (v *Vector3) SetByName(name string, value float32) {

	switch name {
	case "x", "X":
		v.X = value
	case "y", "Y":
		v.Y = value
	case "z", "Z":
		v.Z = value
	default:
		panic("Invalid Vector3 component name: " + name)
	}
}

// Zero sets this vector X, Y and Z components to be zero.
// Returns the pointer to this updated vector.
func (v *Vector3) Zero() *Vector3 {

	v.X = 0
	v.Y = 0
	v.Z = 0
	return v
}

// Copy copies other vector to this one.
// It is equivalent to: *v = *other.
// Returns the pointer to this updated vector.
func (v *Vector3) Copy(other *Vector3) *Vector3 {

	*v = *other
	return v
}

// Add adds other vector to this one.
// Returns the pointer to this updated vector.
func (v *Vector3) Add(other *Vector3) *Vector3 {

	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	return v
}

// AddScalar adds scalar s to each component of this vector.
// Returns the pointer to this updated vector.
func (v *Vector3) AddScalar(s float32) *Vector3 {

	v.X += s
	v.Y += s
	v.Z += s
	return v
}

// AddVectors adds vectors a and b to this one.
// Returns the pointer to this updated vector.
func (v *Vector3) AddVectors(a, b *Vector3) *Vector3 {

	v.X = a.X + b.X
	v.Y = a.Y + b.Y
	v.Z = a.Z + b.Z
	return v
}

// Sub subtracts other vector from this one.
// Returns the pointer to this updated vector.
func (v *Vector3) Sub(other *Vector3) *Vector3 {

	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
	return v
}

// SubScalar subtracts scalar s from each component of this vector.
// Returns the pointer to this updated vector.
func (v *Vector3) SubScalar(s float32) *Vector3 {

	v.X -= s
	v.Y -= s
	v.Z -= s
	return v
}

// SubVectors sets this vector to a - b.
// Returns the pointer to this updated vector.
func (v *Vector3) SubVectors(a, b *Vector3) *Vector3 {

	v.X = a.X - b.X
	v.Y = a.Y - b.Y
	v.Z = a.Z - b.Z
	return v
}

// Multiply multiplies each component of this vector by the corresponding one from other vector.
// Returns the pointer to this updated vector.
func (v *Vector3) Multiply(other *Vector3) *Vector3 {

	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
	return v
}

// MultiplyScalar multiplies each component of this vector by the scalar s.
// Returns the pointer to this updated vector.
func (v *Vector3) MultiplyScalar(s float32) *Vector3 {

	v.X *= s
	v.Y *= s
	v.Z *= s
	return v
}

// Divide divides each component of this vector by the corresponding one from other vector.
// Returns the pointer to this updated vector
func (v *Vector3) Divide(other *Vector3) *Vector3 {

	v.X /= other.X
	v.Y /= other.Y
	v.Z /= other.Z
	return v
}

// DivideScalar divides each component of this vector by the scalar s.
// If scalar is zero, sets this vector to zero.
// Returns the pointer to this updated vector.
func (v *Vector3) DivideScalar(scalar float32) *Vector3 {

	if scalar != 0 {
		invScalar := 1 / scalar
		v.X *= invScalar
		v.Y *= invScalar
		v.Z *= invScalar
	} else {
		v.X = 0
		v.Y = 0
		v.Z = 0
	}
	return v
}

// Min sets this vector components to the minimum values of itself and other vector.
// Returns the pointer to this updated vector.
func (v *Vector3) Min(other *Vector3) *Vector3 {

	if v.X > other.X {
		v.X = other.X
	}
	if v.Y > other.Y {
		v.Y = other.Y
	}
	if v.Z > other.Z {
		v.Z = other.Z
	}
	return v
}

// Max sets this vector components to the maximum value of itself and other vector.
// Returns the pointer to this updated vector.
func (v *Vector3) Max(other *Vector3) *Vector3 {

	if v.X < other.X {
		v.X = other.X
	}
	if v.Y < other.Y {
		v.Y = other.Y
	}
	if v.Z < other.Z {
		v.Z = other.Z
	}
	return v
}

// Clamp sets this vector components to be no less than the corresponding components of min
// and not greater than the corresponding component of max.
// Assumes min < max, if this assumption isn't true it will not operate correctly.
// Returns the pointer to this updated vector.
func (v *Vector3) Clamp(min, max *Vector3) *Vector3 {

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
	return v
}

// ClampScalar sets this vector components to be no less than minVal and not greater than maxVal.
// Returns the pointer to this updated vector.
func (v *Vector3) ClampScalar(minVal, maxVal float32) *Vector3 {

	min := NewVector3(minVal, minVal, minVal)
	max := NewVector3(maxVal, maxVal, maxVal)
	return v.Clamp(min, max)
}

// Floor applies math32.Floor() to each of this vector's components.
// Returns the pointer to this updated vector.
func (v *Vector3) Floor() *Vector3 {

	v.X = Floor(v.X)
	v.Y = Floor(v.Y)
	v.Z = Floor(v.Z)
	return v
}

// Ceil applies math32.Ceil() to each of this vector's components.
// Returns the pointer to this updated vector.
func (v *Vector3) Ceil() *Vector3 {

	v.X = Ceil(v.X)
	v.Y = Ceil(v.Y)
	v.Z = Ceil(v.Z)
	return v
}

// Round rounds each of this vector's components.
// Returns the pointer to this updated vector.
func (v *Vector3) Round() *Vector3 {

	v.X = Floor(v.X + 0.5)
	v.Y = Floor(v.Y + 0.5)
	v.Z = Floor(v.Z + 0.5)
	return v
}

// Negate negates each of this vector's components.
// Returns the pointer to this updated vector.
func (v *Vector3) Negate() *Vector3 {

	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	return v
}

// Dot returns the dot product of this vector with other.
// None of the vectors are changed.
func (v *Vector3) Dot(other *Vector3) float32 {

	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// LengthSq returns the length squared of this vector.
// LengthSq can be used to compare vectors' lengths without the need to perform a square root.
func (v *Vector3) LengthSq() float32 {

	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Length returns the length of this vector.
func (v *Vector3) Length() float32 {

	return Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize normalizes this vector so its length will be 1.
// Returns the pointer to this updated vector.
func (v *Vector3) Normalize() *Vector3 {

	return v.DivideScalar(v.Length())
}

// DistanceTo returns the distance of this point to other.
func (v *Vector3) DistanceTo(other *Vector3) float32 {

	return Sqrt(v.DistanceToSquared(other))
}

// DistanceToSquared returns the distance squared of this point to other.
func (v *Vector3) DistanceToSquared(other *Vector3) float32 {

	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return dx*dx + dy*dy + dz*dz
}

// SetLength sets this vector to have the specified length.
// If the current length is zero, does nothing.
// Returns the pointer to this updated vector.
func (v *Vector3) SetLength(l float32) *Vector3 {

	oldLength := v.Length()
	if oldLength != 0 && l != oldLength {
		v.MultiplyScalar(l / oldLength)
	}
	return v
}

// Lerp sets each of this vector's components to the linear interpolated value of
// alpha between ifself and the corresponding other component.
// Returns the pointer to this updated vector.
func (v *Vector3) Lerp(other *Vector3, alpha float32) *Vector3 {

	v.X += (other.X - v.X) * alpha
	v.Y += (other.Y - v.Y) * alpha
	v.Z += (other.Z - v.Z) * alpha
	return v
}

// Equals returns if this vector is equal to other.
func (v *Vector3) Equals(other *Vector3) bool {

	return (other.X == v.X) && (other.Y == v.Y) && (other.Z == v.Z)
}

// FromArray sets this vector's components from the specified array and offset
// Returns the pointer to this updated vector.
func (v *Vector3) FromArray(array []float32, offset int) *Vector3 {

	v.X = array[offset]
	v.Y = array[offset+1]
	v.Z = array[offset+2]
	return v
}

// ToArray copies this vector's components to array starting at offset.
// Returns the array.
func (v *Vector3) ToArray(array []float32, offset int) []float32 {

	array[offset] = v.X
	array[offset+1] = v.Y
	array[offset+2] = v.Z
	return array
}

// MultiplyVectors multiply vectors a and b storing the result in this vector.
// Returns the pointer to this updated vector.
func (v *Vector3) MultiplyVectors(a, b *Vector3) *Vector3 {

	v.X = a.X * b.X
	v.Y = a.Y * b.Y
	v.Z = a.Z * b.Z
	return v
}

// ApplyAxisAngle rotates the vector around axis by angle.
// Returns the pointer to this updated vector.
func (v *Vector3) ApplyAxisAngle(axis *Vector3, angle float32) *Vector3 {

	var quaternion Quaternion
	v.ApplyQuaternion(quaternion.SetFromAxisAngle(axis, angle))
	return v
}

// ApplyMatrix3 multiplies the specified 3x3 matrix by this vector.
// Returns the pointer to this updated vector.
func (v *Vector3) ApplyMatrix3(m *Matrix3) *Vector3 {

	x := v.X
	y := v.Y
	z := v.Z
	v.X = m[0]*x + m[3]*y + m[6]*z
	v.Y = m[1]*x + m[4]*y + m[7]*z
	v.Z = m[2]*x + m[5]*y + m[8]*z
	return v
}

// ApplyMatrix4 multiplies the specified 4x4 matrix by this vector.
// Returns the pointer to this updated vector.
func (v *Vector3) ApplyMatrix4(m *Matrix4) *Vector3 {

	x := v.X
	y := v.Y
	z := v.Z
	v.X = m[0]*x + m[4]*y + m[8]*z + m[12]
	v.Y = m[1]*x + m[5]*y + m[9]*z + m[13]
	v.Z = m[2]*x + m[6]*y + m[10]*z + m[14]
	return v
}

// ApplyProjection applies the projection matrix m to this vector
// Returns the pointer to this updated vector.
func (v *Vector3) ApplyProjection(m *Matrix4) *Vector3 {

	x := v.X
	y := v.Y
	z := v.Z
	d := 1 / (m[3]*x + m[7]*y + m[11]*z + m[15]) // perspective divide
	v.X = (m[0]*x + m[4]*y + m[8]*z + m[12]) * d
	v.Y = (m[1]*x + m[5]*y + m[9]*z + m[13]) * d
	v.Z = (m[2]*x + m[6]*y + m[10]*z + m[14]) * d
	return v
}

// ApplyQuaternion transforms this vector by multiplying it by
// the specified quaternion and then by the quaternion inverse.
// It basically applies the rotation encoded in the quaternion to this vector.
// Returns the pointer to this updated vector.
func (v *Vector3) ApplyQuaternion(q *Quaternion) *Vector3 {

	x := v.X
	y := v.Y
	z := v.Z

	qx := q.X
	qy := q.Y
	qz := q.Z
	qw := q.W

	// calculate quat * vector
	ix := qw*x + qy*z - qz*y
	iy := qw*y + qz*x - qx*z
	iz := qw*z + qx*y - qy*x
	iw := -qx*x - qy*y - qz*z
	// calculate result * inverse quat
	v.X = ix*qw + iw*-qx + iy*-qz - iz*-qy
	v.Y = iy*qw + iw*-qy + iz*-qx - ix*-qz
	v.Z = iz*qw + iw*-qz + ix*-qy - iy*-qx
	return v
}

// Cross calculates the cross product of this vector with other and returns the result vector.
func (v *Vector3) Cross(other *Vector3) *Vector3 {

	cx := v.Y*other.Z - v.Z*other.Y
	cy := v.Z*other.X - v.X*other.Z
	cz := v.X*other.Y - v.Y*other.X
	v.X = cx
	v.Y = cy
	v.Z = cz
	return v
}

// CrossVectors calculates the cross product of a and b storing the result in this vector.
// Returns the pointer to this updated vector.
func (v *Vector3) CrossVectors(a, b *Vector3) *Vector3 {

	cx := a.Y*b.Z - a.Z*b.Y
	cy := a.Z*b.X - a.X*b.Z
	cz := a.X*b.Y - a.Y*b.X
	v.X = cx
	v.Y = cy
	v.Z = cz
	return v
}

// ProjectOnVector sets this vector to its projection on other vector.
// Returns the pointer to this updated vector.
func (v *Vector3) ProjectOnVector(other *Vector3) *Vector3 {

	var on Vector3
	on.Copy(other).Normalize()
	dot := v.Dot(&on)
	return v.Copy(&on).MultiplyScalar(dot)
}

// ProjectOnPlane sets this vector to its projection on the plane
// specified by its normal vector.
// Returns the pointer to this updated vector.
func (v *Vector3) ProjectOnPlane(planeNormal *Vector3) *Vector3 {

	var tmp Vector3
	tmp.Copy(v).ProjectOnVector(planeNormal)
	return v.Sub(&tmp)
}

// Reflect sets this vector to its reflection relative to the normal vector.
// The normal vector is assumed to be normalized.
// Returns the pointer to this updated vector.
func (v *Vector3) Reflect(normal *Vector3) *Vector3 {

	var tmp Vector3
	return v.Sub(tmp.Copy(normal).MultiplyScalar(2 * v.Dot(normal)))
}

// AngleTo returns the angle between this vector and other
func (v *Vector3) AngleTo(other *Vector3) float32 {

	theta := v.Dot(other) / (v.Length() * other.Length())
	// clamp, to handle numerical problems
	return Acos(Clamp(theta, -1, 1))
}

// SetFromMatrixPosition set this vector from the translation coordinates
// in the specified transformation matrix.
func (v *Vector3) SetFromMatrixPosition(m *Matrix4) *Vector3 {

	v.X = m[12]
	v.Y = m[13]
	v.Z = m[14]
	return v
}

// SetFromMatrixColumn set this vector with the column at index of the m matrix.
// Returns the pointer to this updated vector.
func (v *Vector3) SetFromMatrixColumn(index int, m *Matrix4) *Vector3 {

	offset := index * 4
	v.X = m[offset]
	v.Y = m[offset+1]
	v.Z = m[offset+2]
	return v
}

// Clone returns a copy of this vector
func (v *Vector3) Clone() *Vector3 {

	return NewVector3(v.X, v.Y, v.Z)
}

// SetFromRotationMatrix sets this vector components to the Euler angles
// from the specified pure rotation matrix.
// Returns the pointer to this updated vector.
func (v *Vector3) SetFromRotationMatrix(m *Matrix4) *Vector3 {

	m11 := m[0]
	m12 := m[4]
	m13 := m[8]
	m22 := m[5]
	m23 := m[9]
	m32 := m[6]
	m33 := m[10]

	v.Y = Asin(Clamp(m13, -1, 1))
	if Abs(m13) < 0.99999 {
		v.X = Atan2(-m23, m33)
		v.Z = Atan2(-m12, m11)
	} else {
		v.X = Atan2(m32, m22)
		v.Z = 0
	}
	return v
}

// SetFromQuaternion sets this vector components to the Euler angles
// from the specified quaternion
// Returns the pointer to this updated vector.
func (v *Vector3) SetFromQuaternion(q *Quaternion) *Vector3 {

	matrix := NewMatrix4()
	matrix.MakeRotationFromQuaternion(q)
	v.SetFromRotationMatrix(matrix)
	return v
}

// RandomTangents computes and returns two arbitrary tangents to the vector.
func (v *Vector3) RandomTangents() (*Vector3, *Vector3) {

	t1 := NewVector3(0,0,0)
	t2 := NewVector3(0,0,0)
	length := v.Length()
	if length > 0 {
		n := NewVector3(v.X/length, v.Y/length, v.Z/length)
		randVec := NewVector3(0,0,0)
		if Abs(n.X) < 0.9 {
			randVec.SetX(1)
			t1.CrossVectors(n, randVec)
		} else if Abs(n.Y) < 0.9 {
			randVec.SetY(1)
			t1.CrossVectors(n, randVec)
		} else {
			randVec.SetZ(1)
			t1.CrossVectors(n, randVec)
		}
		t2.CrossVectors(n, t1)
	} else {
		t1.SetX(1)
		t2.SetY(1)
	}

	return t1, t2
}

// TODO: implement similar methods for Vector2 and Vector4
// AlmostEquals returns whether the vector is almost equal to another vector within the specified tolerance.
func (v *Vector3) AlmostEquals(other *Vector3, tolerance float32) bool {

	if (Abs(v.X - other.X) < tolerance) &&
		(Abs(v.Y - other.Y) < tolerance) &&
		(Abs(v.Z - other.Z) < tolerance) {
			return true
	}
	return false
}
