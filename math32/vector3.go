// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

type Vector3 struct {
	X float32
	Y float32
	Z float32
}

// Common Vector3 values
var Vector3Zero = Vector3{0, 0, 0}
var Vector3Up = Vector3{0, 1, 0}
var Vector3Down = Vector3{0, -1, 0}
var Vector3Right = Vector3{1, 0, 0}
var Vector3Left = Vector3{-1, 0, 0}
var Vector3Front = Vector3{0, 0, 1}
var Vector3Back = Vector3{0, 0, -1}

// NewVector creates and returns a pointer to a Vector3 type
func NewVector3(x, y, z float32) *Vector3 {

	return &Vector3{X: x, Y: y, Z: z}
}

// Set sets the components of the vector
func (v *Vector3) Set(x, y, z float32) *Vector3 {

	v.X = x
	v.Y = y
	v.Z = z
	return v
}

// SetX sets the value of the X component of this vector
func (v *Vector3) SetX(x float32) *Vector3 {

	v.X = x
	return v
}

// SetY sets the value of the Y component of this vector
func (v *Vector3) SetY(y float32) *Vector3 {

	v.Y = y
	return v
}

// SetZ sets the value of the Z component of this vector
func (v *Vector3) SetZ(z float32) *Vector3 {

	v.Z = z
	return v
}

// SetComponent sets the value of this vector component
// specified by its index: X=0, Y=1, Z=2
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

// SetByName sets the value of this vector component
// specified by its name: "x|Z", "y|Y", or "z|Z".
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

// Copy copies the specified vector into this one.
func (v *Vector3) Copy(src *Vector3) *Vector3 {

	*v = *src
	return v
}

// Add adds the specified vector to this.
func (v *Vector3) Add(src *Vector3) *Vector3 {

	v.X += src.X
	v.Y += src.Y
	v.Z += src.Z
	return v
}

// AddScalar adds the specified value to each of the vector components
func (v *Vector3) AddScalar(s float32) *Vector3 {

	v.X += s
	v.Y += s
	v.Z += s
	return v
}

// AddVectors adds the specified vectors and saves the result in
// this vector: v = a * b
func (v *Vector3) AddVectors(a, b *Vector3) *Vector3 {

	v.X = a.X + b.X
	v.Y = a.Y + b.Y
	v.Z = a.Z + b.Z
	return v
}

// Sub subtracts the specified vector from this one:
// v = v - p
func (v *Vector3) Sub(p *Vector3) *Vector3 {

	v.X -= p.X
	v.Y -= p.Y
	v.Z -= p.Z
	return v
}

// SubScalar subtracts the specified scalar from this vector:
// v = v - s
func (v *Vector3) SubScalar(s float32) *Vector3 {

	v.X -= s
	v.Y -= s
	v.Z -= s
	return v
}

// SubVectors subtracts the specified vectors and stores
// the result in this vector: v = a - b
func (v *Vector3) SubVectors(a, b *Vector3) *Vector3 {

	v.X = a.X - b.X
	v.Y = a.Y - b.Y
	v.Z = a.Z - b.Z
	return v
}

// Multiply multiplies this vector by the specified vector:
// v = v * a
func (v *Vector3) Multiply(a *Vector3) *Vector3 {

	v.X *= a.X
	v.Y *= a.Y
	v.Z *= a.Z
	return v
}

// MultiplyScalar multiples this vector by the specified scalar:
// v = v * s
func (v *Vector3) MultiplyScalar(s float32) *Vector3 {

	v.X *= s
	v.Y *= s
	v.Z *= s
	return v
}

// MultiplyVectors multiply the specified vectors storing the
// result in this vector: v = a * b
func (v *Vector3) MultiplyVectors(a, b *Vector3) *Vector3 {

	v.X = a.X * b.X
	v.Y = a.Y * b.Y
	v.Z = a.Z * b.Z
	return v
}

//// ApplyEuler rotates this vector from the specified euler angles
//func (v *Vector3) ApplyEuler(euler *Euler) *Vector3 {
//
//	var quaternion Quaternion
//	v.ApplyQuaternion(quaternion.SetFromEuler2(euler))
//	return v
//}

// ApplyAxisAngle rotates the vertex around the specified axis by
// the specified angle
func (v *Vector3) ApplyAxisAngle(axis *Vector3, angle float32) *Vector3 {

	var quaternion Quaternion
	v.ApplyQuaternion(quaternion.SetFromAxisAngle(axis, angle))
	return v
}

// ApplyMatrix3 multiplies the specified 3x3 matrix by this vector:
// v = m * v
func (v *Vector3) ApplyMatrix3(m *Matrix3) *Vector3 {

	x := v.X
	y := v.Y
	z := v.Z
	v.X = m[0]*x + m[3]*y + m[6]*z
	v.Y = m[1]*x + m[4]*y + m[7]*z
	v.Z = m[2]*x + m[5]*y + m[8]*z
	return v
}

// ApplyMatrix4 multiplies the specified 4x4 matrix by this vector:
// v = m * v
func (v *Vector3) ApplyMatrix4(m *Matrix4) *Vector3 {

	x := v.X
	y := v.Y
	z := v.Z
	v.X = m[0]*x + m[4]*y + m[8]*z + m[12]
	v.Y = m[1]*x + m[5]*y + m[9]*z + m[13]
	v.Z = m[2]*x + m[6]*y + m[10]*z + m[14]
	return v
}

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
// the specified quaternion and then by the quaternion
// inverse.
// It basically applies the rotation encoded in the quaternion
// to this vector.
// v = q * v * inverse(q)
func (v *Vector3) ApplyQuaternion(q *Quaternion) *Vector3 {

	x := v.X
	y := v.Y
	z := v.Z

	qx := q.x
	qy := q.y
	qz := q.z
	qw := q.w

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

func (v *Vector3) TransformDirection(m *Matrix4) *Vector3 {

	x := v.X
	y := v.Y
	z := v.Z
	v.X = m[0]*x + m[4]*y + m[8]*z
	v.Y = m[1]*x + m[5]*y + m[9]*z
	v.Z = m[2]*x + m[6]*y + m[10]*z
	v.Normalize()
	return v
}

func (v *Vector3) Divide(other *Vector3) *Vector3 {

	v.X /= other.X
	v.Y /= other.Y
	v.Z /= other.Z
	return v
}

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

// Max sets this vector with maximum components from itself and the other vector
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

func (this *Vector3) Clamp(min, max *Vector3) *Vector3 {

	// This function assumes min < max, if this assumption isn't true it will not operate correctly
	if this.X < min.X {
		this.X = min.X
	} else if this.X > max.X {
		this.X = max.X
	}

	if this.Y < min.Y {
		this.Y = min.Y
	} else if this.Y > max.Y {
		this.Y = max.Y
	}

	if this.Z < min.Z {
		this.Z = min.Z
	} else if this.Z > max.Z {
		this.Z = max.Z
	}
	return this
}

func (this *Vector3) ClampScalar(minVal, maxVal float32) *Vector3 {

	min := NewVector3(minVal, minVal, minVal)
	max := NewVector3(maxVal, maxVal, maxVal)
	return this.Clamp(min, max)
}

func (this *Vector3) Floor() *Vector3 {

	this.X = Floor(this.X)
	this.Y = Floor(this.Y)
	this.Z = Floor(this.Z)
	return this
}

func (this *Vector3) Ceil() *Vector3 {

	this.X = Ceil(this.X)
	this.Y = Ceil(this.Y)
	this.Z = Ceil(this.Z)
	return this
}

func (this *Vector3) Round() *Vector3 {

	this.X = Floor(this.X + 0.5)
	this.Y = Floor(this.Y + 0.5)
	this.Z = Floor(this.Z + 0.5)
	return this
}

func (this *Vector3) RoundToZero() *Vector3 {

	if this.X < 0 {
		this.X = Ceil(this.X)
	} else {
		this.X = Floor(this.X)
	}

	if this.Y < 0 {
		this.Y = Ceil(this.Y)
	} else {
		this.Y = Floor(this.Y)
	}

	if this.Z < 0 {
		this.Z = Ceil(this.Z)
	} else {
		this.Z = Floor(this.Z)
	}
	return this
}

func (this *Vector3) Negate() *Vector3 {

	this.X = -this.X
	this.Y = -this.Y
	this.Z = -this.Z
	return this
}

func (this *Vector3) Dot(v *Vector3) float32 {

	return this.X*v.X + this.Y*v.Y + this.Z*v.Z
}

func (this *Vector3) LengthSq() float32 {

	return this.X*this.X + this.Y*this.Y + this.Z*this.Z
}

func (this *Vector3) Length() float32 {

	return Sqrt(this.X*this.X + this.Y*this.Y + this.Z*this.Z)
}

func (this *Vector3) LengthManhattan(v *Vector3) float32 {

	return Abs(this.X + Abs(this.Y+Abs(this.Z)))
}

func (this *Vector3) Normalize() *Vector3 {

	return this.DivideScalar(this.Length())
}

func (this *Vector3) SetLength(l float32) *Vector3 {

	oldLength := this.Length()
	if oldLength != 0 && l != oldLength {
		this.MultiplyScalar(l / oldLength)
	}
	return this
}

func (this *Vector3) Lerp(v *Vector3, alpha float32) *Vector3 {

	this.X += (v.X - this.X) * alpha
	this.Y += (v.Y - this.Y) * alpha
	this.Z += (v.Z - this.Z) * alpha
	return this
}

func (this *Vector3) LerpVectors(v1, v2 *Vector3, alpha float32) *Vector3 {

	this.SubVectors(v2, v2).MultiplyScalar(alpha).Add(v1)
	return this
}

func (this *Vector3) Cross(v *Vector3) *Vector3 {

	cx := this.Y*v.Z - this.Z*v.Y
	cy := this.Z*v.X - this.X*v.Z
	cz := this.X*v.Y - this.Y*v.X
	this.X = cx
	this.Y = cy
	this.Z = cz
	return this
}

func (this *Vector3) CrossVectors(a, b *Vector3) *Vector3 {

	cx := a.Y*b.Z - a.Z*b.Y
	cy := a.Z*b.X - a.X*b.Z
	cz := a.X*b.Y - a.Y*b.X
	this.X = cx
	this.Y = cy
	this.Z = cz
	return this
}

func (this *Vector3) ProjectOnVector(vector *Vector3) *Vector3 {

	var v1 Vector3
	v1.Copy(vector).Normalize()
	dot := this.Dot(&v1)
	return this.Copy(&v1).MultiplyScalar(dot)
}

func (this *Vector3) ProjectOnPlane(planeNormal *Vector3) *Vector3 {

	var v1 Vector3
	v1.Copy(this).ProjectOnVector(planeNormal)
	return this.Sub(&v1)
}

func (this *Vector3) Reflect(normal *Vector3) *Vector3 {
	// reflect incident vector off plane orthogonal to normal
	// normal is assumed to have unit length

	var v1 Vector3
	return this.Sub(v1.Copy(normal).MultiplyScalar(2 * this.Dot(normal)))
}

func (this *Vector3) AngleTo(v *Vector3) float32 {

	theta := this.Dot(v) / (this.Length() * v.Length())
	// clamp, to handle numerical problems
	return Acos(Clamp(theta, -1, 1))
}

func (this *Vector3) DistanceTo(v *Vector3) float32 {

	return Sqrt(this.DistanceToSquared(v))
}

func (this *Vector3) DistanceToSquared(v *Vector3) float32 {

	dx := this.X - v.X
	dy := this.Y - v.Y
	dz := this.Z - v.Z
	return dx*dx + dy*dy + dz*dz
}

// SetFromMatrixPosition set this vector from translation coordinates
// in the specified transformation matrix.
func (v *Vector3) SetFromMatrixPosition(m *Matrix4) *Vector3 {

	v.X = m[12]
	v.Y = m[13]
	v.Z = m[14]
	return v
}

func (this *Vector3) SetFromMatrixScale(m *Matrix4) *Vector3 {

	sx := this.Set(m[0], m[1], m[2]).Length()
	sy := this.Set(m[4], m[5], m[6]).Length()
	sz := this.Set(m[8], m[9], m[10]).Length()

	this.X = sx
	this.Y = sy
	this.Z = sz

	return this
}

// SetFromMatrixColumn set this vector with the column at index 'index'
// of the 'm' matrix.
func (v *Vector3) SetFromMatrixColumn(index int, m *Matrix4) *Vector3 {

	offset := index * 4
	v.X = m[offset]
	v.Y = m[offset+1]
	v.Z = m[offset+2]
	return v
}

func (this *Vector3) Equals(v *Vector3) bool {

	return (v.X == this.X) && (v.Y == this.Y) && (v.Z == this.Z)
}

func (v *Vector3) FromArray(array []float32, offset int) *Vector3 {

	v.X = array[offset]
	v.Y = array[offset+1]
	v.Z = array[offset+2]
	return v
}

func (v *Vector3) ToArray(array []float32, offset int) []float32 {

	array[offset] = v.X
	array[offset+1] = v.Y
	array[offset+2] = v.Z
	return array
}

func (this *Vector3) Clone() *Vector3 {

	return NewVector3(this.X, this.Y, this.Z)
}

// SetFromRotationMatrix sets this vector components to the Euler angles
// from the specified pure rotation matrix.
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
func (v *Vector3) SetFromQuaternion(q *Quaternion) *Vector3 {

	matrix := NewMatrix4()
	matrix.MakeRotationFromQuaternion(q)
	v.SetFromRotationMatrix(matrix)
	return v
}
