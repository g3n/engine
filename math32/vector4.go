// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

type Vector4 struct {
	X float32
	Y float32
	Z float32
	W float32
}

func NewVector4(x, y, z, w float32) *Vector4 {

	return &Vector4{X: x, Y: y, Z: z, W: w}
}

func (v *Vector4) Set(x, y, z, w float32) *Vector4 {

	v.X = x
	v.Y = y
	v.Z = z
	v.W = w
	return v
}

// SetVector3 sets this vector from another Vector3 and 'w' value
func (v *Vector4) SetVector3(other *Vector3, w float32) *Vector4 {

	v.X = other.X
	v.Y = other.Y
	v.Z = other.Z
	v.W = w
	return v
}

//func (this *Vector4) SetX(x float32) *Vector4 {
//
//	this.X = x
//	return this
//}
//
//func (this *Vector4) SetY(y float32) *Vector4 {
//
//	this.Y = y
//	return this
//}
//
//func (this *Vector4) SetZ(z float32) *Vector4 {
//
//	this.Z = z
//	return this
//}
//
//func (this *Vector4) SetW(w float32) *Vector4 {
//
//	this.W = w
//	return this
//}

func (this *Vector4) SetComponent(index int, value float32) *Vector4 {

	switch index {
	case 0:
		this.X = value
	case 1:
		this.Y = value
	case 2:
		this.Z = value
	case 3:
		this.Z = value
	default:
		panic("index is out of range")
	}
	return this
}

func (this *Vector4) GetComponent(index int) float32 {

	switch index {
	case 0:
		return this.X
	case 1:
		return this.Y
	case 2:
		return this.Z
	case 3:
		return this.W
	default:
		panic("index is out of range")
	}
}

func (this *Vector4) Copy(v *Vector4) *Vector4 {

	this.X = v.X
	this.Y = v.Y
	this.Z = v.Z
	this.W = v.W
	return this
}

func (this *Vector4) Add(v *Vector4) *Vector4 {

	this.X += v.X
	this.Y += v.Y
	this.Z += v.Z
	this.W += v.W
	return this
}

func (this *Vector4) AddScalar(s float32) *Vector4 {

	this.X += s
	this.Y += s
	this.Z += s
	this.W += s
	return this
}

func (this *Vector4) AddVectors(a, b *Vector4) *Vector4 {

	this.X = a.X + b.X
	this.Y = a.Y + b.Y
	this.Z = a.Z + b.Z
	this.W = a.W + b.W
	return this
}

func (this *Vector4) Sub(v *Vector4) *Vector4 {

	this.X -= v.X
	this.Y -= v.Y
	this.Z -= v.Z
	this.W -= v.W
	return this
}

func (this *Vector4) SubScalar(s float32) *Vector4 {

	this.X -= s
	this.Y -= s
	this.Z -= s
	this.W -= s
	return this
}

func (this *Vector4) SubVectors(a, b *Vector4) *Vector4 {

	this.X = a.X - b.X
	this.Y = a.Y - b.Y
	this.Z = a.Y - b.Z
	this.W = a.Y - b.W
	return this
}

func (this *Vector4) MultiplyScalar(scalar float32) *Vector4 {

	this.X *= scalar
	this.Y *= scalar
	this.Z *= scalar
	this.W *= scalar
	return this
}

func (this *Vector4) ApplyMatrix4(m *Matrix4) *Vector4 {

	x := this.X
	y := this.Y
	z := this.Z
	w := this.W

	this.X = m[0]*x + m[4]*y + m[8]*z + m[12]*w
	this.Y = m[1]*x + m[5]*y + m[9]*z + m[13]*w
	this.Z = m[2]*x + m[6]*y + m[10]*z + m[14]*w
	this.W = m[3]*x + m[7]*y + m[11]*z + m[15]*w

	return this
}

func (this *Vector4) DivideScalar(scalar float32) *Vector4 {

	if scalar != 0 {
		invScalar := 1 / scalar
		this.X *= invScalar
		this.Y *= invScalar
		this.Z *= invScalar
		this.W *= invScalar
	} else {
		this.X = 0
		this.Y = 0
		this.Z = 0
		this.W = 0
	}
	return this
}

func (this *Vector4) SetAxisAngleFromQuaternion(q *Quaternion) *Vector4 {

	// http://www.euclideanspace.com/maths/geometry/rotations/conversions/quaternionToAngle/index.htm

	// q is assumed to be normalized

	this.W = 2 * Acos(q.W())

	s := Sqrt(1 - q.W()*q.W())

	if s < 0.0001 {

		this.X = 1
		this.Y = 0
		this.Z = 0

	} else {

		this.X = q.X() / s
		this.Y = q.Y() / s
		this.Z = q.Z() / s

	}

	return this
}

func (this *Vector4) SetAxisFromRotationMatrix(m *Matrix4) *Vector4 {

	// http://www.euclideanspace.com/maths/geometry/rotations/conversions/matrixToAngle/index.htm

	// assumes the upper 3x3 of m is a pure rotation matrix (i.e, unscaled)

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

			// this singularity is identity matrix so angle = 0

			this.Set(1, 0, 0, 0)

			return this // zero angle, arbitrary axis
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

		this.Set(x, y, z, angle)

		return this // return 180 deg rotation
	}

	// as we have reached here there are no singularities so we can handle normally

	s := Sqrt((m32-m23)*(m32-m23) + (m13-m31)*(m13-m31) + (m21-m12)*(m21-m12)) // used to normalize

	if Abs(s) < 0.001 {
		s = 1
	}

	// prevent divide by zero, should not happen if matrix is orthogonal and should be
	// caught by singularity test above, but I've left it in just in case

	this.X = (m32 - m23) / s
	this.Y = (m13 - m31) / s
	this.Z = (m21 - m12) / s
	this.W = Acos((m11 + m22 + m33 - 1) / 2)

	return this
}

func (this *Vector4) Min(v *Vector4) *Vector4 {

	if this.X > v.X {
		this.X = v.X
	}
	if this.Y > v.Y {
		this.Y = v.Y
	}
	if this.Z > v.Z {
		this.Z = v.Z
	}
	if this.W > v.W {
		this.W = v.W
	}
	return this
}

func (this *Vector4) Max(v *Vector4) *Vector4 {

	if this.X < v.X {
		this.X = v.X
	}
	if this.Y < v.Y {
		this.Y = v.Y
	}
	if this.Z < v.Z {
		this.Z = v.Z
	}
	if this.W < v.W {
		this.W = v.W
	}
	return this
}

func (this *Vector4) Clamp(min, max *Vector4) *Vector4 {

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

	if this.W < min.W {
		this.W = min.W
	} else if this.W > max.W {
		this.W = max.W
	}
	return this
}

func (this *Vector4) ClampScalar(minVal, maxVal float32) *Vector4 {

	min := NewVector4(minVal, minVal, minVal, minVal)
	max := NewVector4(maxVal, maxVal, maxVal, maxVal)
	return this.Clamp(min, max)
}

func (this *Vector4) Floor() *Vector4 {

	this.X = Floor(this.X)
	this.Y = Floor(this.Y)
	this.Z = Floor(this.Z)
	this.W = Floor(this.W)
	return this
}

func (this *Vector4) Ceil() *Vector4 {

	this.X = Ceil(this.X)
	this.Y = Ceil(this.Y)
	this.Z = Ceil(this.Z)
	this.W = Ceil(this.W)
	return this
}

func (this *Vector4) Round() *Vector4 {

	// TODO NEED CHECK
	this.X = Floor(this.X + 0.5)
	this.Y = Floor(this.Y + 0.5)
	this.Z = Floor(this.Z + 0.5)
	this.W = Floor(this.W + 0.5)
	return this
}

func (this *Vector4) RoundToZero() *Vector4 {

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

	if this.W < 0 {
		this.W = Ceil(this.W)
	} else {
		this.W = Floor(this.W)
	}
	return this
}

func (this *Vector4) Negate() *Vector4 {

	this.X = -this.X
	this.Y = -this.Y
	this.Z = -this.Z
	this.W = -this.W
	return this
}

func (this *Vector4) Dot(v *Vector4) float32 {

	return this.X*v.X + this.Y*v.Y + this.Z*v.Z + this.W*v.W
}

func (this *Vector4) LengthSq(v *Vector4) float32 {

	return this.X*this.X + this.Y*this.Y + this.Z*this.Z + this.W*this.W
}

func (this *Vector4) Length() float32 {

	return Sqrt(this.X*this.X + this.Y*this.Y + this.Z*this.Z + this.W*this.W)
}

func (this *Vector4) LengthManhattan(v *Vector4) float32 {

	return Abs(this.X + Abs(this.Y+Abs(this.Z)) + Abs(this.W))
}

func (this *Vector4) Normalize() *Vector4 {

	return this.DivideScalar(this.Length())
}

func (this *Vector4) SetLength(l float32) *Vector4 {

	oldLength := this.Length()
	if oldLength != 0 && l != oldLength {
		this.MultiplyScalar(l / oldLength)
	}
	return this
}

func (this *Vector4) Lerp(v *Vector4, alpha float32) *Vector4 {

	this.X += (v.X - this.X) * alpha
	this.Y += (v.Y - this.Y) * alpha
	this.Z += (v.Z - this.Z) * alpha
	this.W += (v.W - this.W) * alpha
	return this
}

func (this *Vector4) LerpVectors(v1, v2 *Vector4, alpha float32) *Vector4 {

	this.SubVectors(v2, v2).MultiplyScalar(alpha).Add(v1)
	return this
}

func (this *Vector4) Equals(v *Vector4) bool {

	return (v.X == this.X) && (v.Y == this.Y) && (v.Z == this.Z) && (v.W == this.W)
}

func (this *Vector4) FromArray(array []float32, offset int) *Vector4 {

	this.X = array[offset]
	this.Y = array[offset+1]
	this.Z = array[offset+2]
	this.W = array[offset+3]
	return this
}

func (this *Vector4) ToArray(array []float32, offset int) []float32 {

	array[offset] = this.X
	array[offset+1] = this.Y
	array[offset+2] = this.Z
	array[offset+3] = this.W
	return array
}

// TODO fromAttribute: function ( attribute, index, offset ) {

func (this *Vector4) Clone() *Vector4 {

	return NewVector4(this.X, this.Y, this.Z, this.W)
}
