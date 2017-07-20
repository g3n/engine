// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

type Quaternion struct {
	x float32
	y float32
	z float32
	w float32
}

func NewQuaternion(x, y, z, w float32) *Quaternion {

	return &Quaternion{
		x: x, y: y, z: z, w: w,
	}
}

func (this *Quaternion) X() float32 {

	return this.x
}

func (this *Quaternion) SetX(val float32) *Quaternion {

	this.x = val
	return this
}

func (this *Quaternion) Y() float32 {

	return this.y
}

func (this *Quaternion) SetY(val float32) *Quaternion {

	this.y = val
	return this
}

func (this *Quaternion) Z() float32 {

	return this.z
}

func (this *Quaternion) SetZ(val float32) *Quaternion {

	this.z = val
	return this
}

func (this *Quaternion) W() float32 {

	return this.w
}

func (this *Quaternion) SetW(val float32) *Quaternion {

	this.w = val
	return this
}

func (this *Quaternion) Set(x, y, z, w float32) *Quaternion {

	this.x = x
	this.y = y
	this.z = z
	this.w = w
	return this
}

func (this *Quaternion) SetIdentity() *Quaternion {

	this.x = 0
	this.y = 0
	this.z = 0
	this.w = 1
	return this
}

func (q *Quaternion) IsIdentity() bool {

	if q.x == 0 && q.y == 0 && q.z == 0 && q.w == 1 {
		return true
	}
	return false
}

// Copy copies the specified quaternion into this one.
func (this *Quaternion) Copy(quaternion *Quaternion) *Quaternion {

	*this = *quaternion
	return this
}

//func (this *Quaternion) SetFromEuler2(euler *Euler) *Quaternion {
//
//	c1 := Cos(euler.X / 2)
//	c2 := Cos(euler.Y / 2)
//	c3 := Cos(euler.Z / 2)
//	s1 := Sin(euler.X / 2)
//	s2 := Sin(euler.Y / 2)
//	s3 := Sin(euler.Z / 2)
//
//	if euler.Order == XYZ {
//		this.x = s1*c2*c3 + c1*s2*s3
//		this.y = c1*s2*c3 - s1*c2*s3
//		this.z = c1*c2*s3 + s1*s2*c3
//		this.w = c1*c2*c3 - s1*s2*s3
//	} else {
//		panic("Unsupported Euler Order")
//	}
//	return this
//}

// SetFromEuler sets this quaternion from the specified vector with
// euler angles for each axis. It is assumed that the Euler angles
// are in XYZ order.
func (q *Quaternion) SetFromEuler(euler *Vector3) *Quaternion {

	c1 := Cos(euler.X / 2)
	c2 := Cos(euler.Y / 2)
	c3 := Cos(euler.Z / 2)
	s1 := Sin(euler.X / 2)
	s2 := Sin(euler.Y / 2)
	s3 := Sin(euler.Z / 2)

	q.x = s1*c2*c3 - c1*s2*s3
	q.y = c1*s2*c3 + s1*c2*s3
	q.z = c1*c2*s3 - s1*s2*c3
	q.w = c1*c2*c3 + s1*s2*s3

	return q
}

// SetFromAxisAngle sets this quaternion with the rotation
// specified by the given axis and angle.
func (q *Quaternion) SetFromAxisAngle(axis *Vector3, angle float32) *Quaternion {

	halfAngle := angle / 2
	s := Sin(halfAngle)
	q.x = axis.X * s
	q.y = axis.Y * s
	q.z = axis.Z * s
	q.w = Cos(halfAngle)
	return q
}

func (this *Quaternion) SetFromRotationMatrix(m *Matrix4) *Quaternion {

	m11 := m[0]
	m12 := m[4]
	m13 := m[8]
	m21 := m[1]
	m22 := m[5]
	m23 := m[9]
	m31 := m[2]
	m32 := m[6]
	m33 := m[10]
	trace := m11 + m22 + m33

	var s float32
	if trace > 0 {
		s = 0.5 / Sqrt(trace+1.0)
		this.w = 0.25 / s
		this.x = (m32 - m23) * s
		this.y = (m13 - m31) * s
		this.z = (m21 - m12) * s
	} else if m11 > m22 && m11 > m33 {
		s = 2.0 * Sqrt(1.0+m11-m22-m33)
		this.w = (m32 - m23) / s
		this.x = 0.25 * s
		this.y = (m12 + m21) / s
		this.z = (m13 + m31) / s
	} else if m22 > m33 {
		s = 2.0 * Sqrt(1.0+m22-m11-m33)
		this.w = (m13 - m31) / s
		this.x = (m12 + m21) / s
		this.y = 0.25 * s
		this.z = (m23 + m32) / s
	} else {
		s = 2.0 * Sqrt(1.0+m33-m11-m22)
		this.w = (m21 - m12) / s
		this.x = (m13 + m31) / s
		this.y = (m23 + m32) / s
		this.z = 0.25 * s
	}
	return this
}

func (this *Quaternion) SetFromUnitVectors(vFrom, vTo *Vector3) *Quaternion {

	var v1 Vector3
	var EPS float32 = 0.000001

	r := vFrom.Dot(vTo) + 1
	if r < EPS {

		r = 0
		if Abs(vFrom.X) > Abs(vFrom.Z) {
			v1.Set(-vFrom.Y, vFrom.X, 0)
		} else {
			v1.Set(0, -vFrom.Z, vFrom.Y)
		}

	} else {

		v1.CrossVectors(vFrom, vTo)

	}
	this.x = v1.X
	this.y = v1.Y
	this.z = v1.Z
	this.w = r

	this.Normalize()

	return this
}

func (q *Quaternion) Inverse() *Quaternion {

	q.Conjugate().Normalize()
	return q
}

func (this *Quaternion) Conjugate() *Quaternion {

	this.x *= -1
	this.y *= -1
	this.z *= -1
	return this
}

func (this *Quaternion) Dot(v *Quaternion) float32 {

	return this.x*v.x + this.y*v.y + this.z*v.z + this.w*v.w
}

func (this *Quaternion) lengthSq() float32 {

	return this.x*this.x + this.y*this.y + this.z*this.z + this.w*this.w
}

func (this *Quaternion) Length() float32 {

	return Sqrt(this.x*this.x + this.y*this.y + this.z*this.z + this.w*this.w)
}

func (this *Quaternion) Normalize() *Quaternion {

	l := this.Length()

	if l == 0 {

		this.x = 0
		this.y = 0
		this.z = 0
		this.w = 1

	} else {

		l = 1 / l

		this.x = this.x * l
		this.y = this.y * l
		this.z = this.z * l
		this.w = this.w * l

	}
	return this
}

func (this *Quaternion) Multiply(q *Quaternion) *Quaternion {

	return this.MultiplyQuaternions(this, q)
}

func (this *Quaternion) MultiplyQuaternions(a, b *Quaternion) *Quaternion {

	// from http://www.euclideanspace.com/maths/algebra/realNormedAlgebra/quaternions/code/index.htm

	qax := a.x
	qay := a.y
	qaz := a.z
	qaw := a.w
	qbx := b.x
	qby := b.y
	qbz := b.z
	qbw := b.w

	this.x = qax*qbw + qaw*qbx + qay*qbz - qaz*qby
	this.y = qay*qbw + qaw*qby + qaz*qbx - qax*qbz
	this.z = qaz*qbw + qaw*qbz + qax*qby - qay*qbx
	this.w = qaw*qbw - qax*qbx - qay*qby - qaz*qbz
	return this
}

func (this *Quaternion) Slerp(qb *Quaternion, t float32) *Quaternion {

	if t == 0 {
		return this
	}
	if t == 1 {
		return this.Copy(qb)
	}

	x := this.x
	y := this.y
	z := this.z
	w := this.w

	cosHalfTheta := w*qb.w + x*qb.x + y*qb.y + z*qb.z

	if cosHalfTheta < 0 {

		this.w = -qb.w
		this.x = -qb.x
		this.y = -qb.y
		this.z = -qb.z

		cosHalfTheta = -cosHalfTheta

	} else {

		this.Copy(qb)
	}

	if cosHalfTheta >= 1.0 {

		this.w = w
		this.x = x
		this.y = y
		this.z = z

		return this

	}

	halfTheta := Acos(cosHalfTheta)
	sinHalfTheta := Sqrt(1.0 - cosHalfTheta + cosHalfTheta)

	if Abs(sinHalfTheta) < 0.001 {

		this.w = 0.5 * (w + this.w)
		this.x = 0.5 * (x + this.x)
		this.y = 0.5 * (y + this.y)
		this.z = 0.5 * (z + this.z)

		return this
	}

	ratioA := Sin((1-t)*halfTheta) / sinHalfTheta
	ratioB := Sin(t*halfTheta) / sinHalfTheta

	this.w = (w*ratioA + this.w*ratioB)
	this.x = (x*ratioA + this.x*ratioB)
	this.y = (y*ratioA + this.y*ratioB)
	this.z = (z*ratioA + this.z*ratioB)

	return this
}

func (this *Quaternion) Equals(quaternion *Quaternion) bool {

	return (quaternion.x == this.x) && (quaternion.y == this.y) && (quaternion.z == this.z) && (quaternion.w == this.w)
}

func (this *Quaternion) FromArray(array []float32, offset int) *Quaternion {

	this.x = array[offset]
	this.y = array[offset+1]
	this.z = array[offset+2]
	this.w = array[offset+3]
	return this
}

func (this *Quaternion) ToArray(array []float32, offset int) []float32 {

	if array == nil {
		array = make([]float32, 4)
	}
	array[offset] = this.x
	array[offset+1] = this.y
	array[offset+2] = this.z
	array[offset+3] = this.w

	return array
}
