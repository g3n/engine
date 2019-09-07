// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

import "errors"

// Matrix4 is 4x4 matrix organized internally as column matrix.
type Matrix4 [16]float32

// NewMatrix4 creates and returns a pointer to a new Matrix4
// initialized as the identity matrix.
func NewMatrix4() *Matrix4 {

	var mat Matrix4
	mat.Identity()
	return &mat
}

// Set sets all the elements of this matrix row by row starting at row1, column1,
// row1, column2, row1, column3 and so forth.
// Returns pointer to this updated Matrix.
func (m *Matrix4) Set(n11, n12, n13, n14, n21, n22, n23, n24, n31, n32, n33, n34, n41, n42, n43, n44 float32) *Matrix4 {

	m[0] = n11
	m[4] = n12
	m[8] = n13
	m[12] = n14
	m[1] = n21
	m[5] = n22
	m[9] = n23
	m[13] = n24
	m[2] = n31
	m[6] = n32
	m[10] = n33
	m[14] = n34
	m[3] = n41
	m[7] = n42
	m[11] = n43
	m[15] = n44
	return m
}

// Identity sets this matrix as the identity matrix.
// Returns pointer to this updated matrix.
func (m *Matrix4) Identity() *Matrix4 {

	m.Set(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	return m
}

// Zero sets this matrix as the zero matrix.
// Returns the pointer to this updated matrix.
func (m *Matrix4) Zero() *Matrix4 {

	m.Set(
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
		0, 0, 0, 0,
	)
	return m
}

// Copy copies src matrix into this one.
// Returns pointer to this updated matrix.
func (m *Matrix4) Copy(src *Matrix4) *Matrix4 {

	*m = *src
	return m
}

// CopyPosition copies the position elements of the src matrix into this one.
// Returns pointer to this updated matrix.
func (m *Matrix4) CopyPosition(src *Matrix4) *Matrix4 {

	m[12] = src[12]
	m[13] = src[13]
	m[14] = src[14]
	return m
}

// ExtractBasis updates the specified vectors with the basis vectors of this matrix.
// Returns pointer to this unchanged matrix.
func (m *Matrix4) ExtractBasis(xAxis, yAxis, zAxis *Vector3) *Matrix4 {

	xAxis.Set(m[0], m[1], m[2])
	yAxis.Set(m[4], m[5], m[6])
	zAxis.Set(m[8], m[9], m[10])
	return m
}

// MakeBasis sets this matrix basis vectors from the specified vectors.
// Returns pointer to this updated matrix.
func (m *Matrix4) MakeBasis(xAxis, yAxis, zAxis *Vector3) *Matrix4 {

	m.Set(
		xAxis.X, yAxis.X, zAxis.X, 0,
		xAxis.Y, yAxis.Y, zAxis.Y, 0,
		xAxis.Z, yAxis.Z, zAxis.Z, 0,
		0, 0, 0, 1,
	)
	return m
}

// ExtractRotation set this matrix as rotation matrix from the src transformation matrix.
// Returns pointer to this updated matrix.
func (m *Matrix4) ExtractRotation(src *Matrix4) *Matrix4 {

	var v1 Vector3

	scaleX := 1 / v1.Set(src[0], src[1], src[2]).Length()
	scaleY := 1 / v1.Set(src[4], src[5], src[6]).Length()
	scaleZ := 1 / v1.Set(src[8], src[9], src[10]).Length()

	m[0] = src[0] * scaleX
	m[1] = src[1] * scaleX
	m[2] = src[2] * scaleX

	m[4] = src[4] * scaleY
	m[5] = src[5] * scaleY
	m[6] = src[6] * scaleY

	m[8] = src[8] * scaleZ
	m[9] = src[9] * scaleZ
	m[10] = src[10] * scaleZ
	return m
}

// MakeRotationFromEuler set this a matrix as a rotation matrix from the specified euler angles.
// Returns pointer to this updated matrix.
func (m *Matrix4) MakeRotationFromEuler(euler *Vector3) *Matrix4 {

	x := euler.X
	y := euler.Y
	z := euler.Z
	a := Cos(x)
	b := Sin(x)
	c := Cos(y)
	d := Sin(y)
	e := Cos(z)
	f := Sin(z)

	ae := a * e
	af := a * f
	be := b * e
	bf := b * f
	m[0] = c * e
	m[4] = -c * f
	m[8] = d
	m[1] = af + be*d
	m[5] = ae - bf*d
	m[9] = -b * c
	m[2] = bf - ae*d
	m[6] = be + af*d
	m[10] = a * c

	// bottom row
	m[3] = 0
	m[7] = 0
	m[11] = 0

	// last column
	m[12] = 0
	m[13] = 0
	m[14] = 0
	m[15] = 1
	return m
}

// MakeRotationFromQuaternion sets this matrix as a rotation matrix from the specified quaternion.
// Returns pointer to this updated matrix.
func (m *Matrix4) MakeRotationFromQuaternion(q *Quaternion) *Matrix4 {

	x := q.X
	y := q.Y
	z := q.Z
	w := q.W

	x2 := x + x
	y2 := y + y
	z2 := z + z

	xx := x * x2
	xy := x * y2
	xz := x * z2
	yy := y * y2
	yz := y * z2
	zz := z * z2
	wx := w * x2
	wy := w * y2
	wz := w * z2

	m[0] = 1 - (yy + zz)
	m[4] = xy - wz
	m[8] = xz + wy

	m[1] = xy + wz
	m[5] = 1 - (xx + zz)
	m[9] = yz - wx

	m[2] = xz - wy
	m[6] = yz + wx
	m[10] = 1 - (xx + yy)

	// bottom row
	m[3] = 0
	m[7] = 0
	m[11] = 0

	// last column
	m[12] = 0
	m[13] = 0
	m[14] = 0
	m[15] = 1
	return m
}

// LookAt sets this matrix as view transform matrix with origin at eye,
// looking at target and using the up vector.
// Returns pointer to this updated matrix.
func (m *Matrix4) LookAt(eye, target, up *Vector3) *Matrix4 {

	var x, y, z Vector3

	z.SubVectors(eye, target)
	if z.LengthSq() == 0 {
		// Eye and target are in the same position
		z.Z = 1
	}
	z.Normalize()

	x.CrossVectors(up, &z)
	if x.LengthSq() == 0 {
		// Up and Z are parallel
		if Abs(up.Z) == 1 {
			z.X += 0.0001
		} else {
			z.Z += 0.0001
		}
		z.Normalize()
		x.CrossVectors(up, &z)
	}
	x.Normalize()

	y.CrossVectors(&z, &x)

	m[0] = x.X
	m[1] = x.Y
	m[2] = x.Z

	m[4] = y.X
	m[5] = y.Y
	m[6] = y.Z

	m[8] = z.X
	m[9] = z.Y
	m[10] = z.Z

	return m
}

// Multiply multiply this matrix by the other matrix
// Returns pointer to this updated matrix.
func (m *Matrix4) Multiply(other *Matrix4) *Matrix4 {

	return m.MultiplyMatrices(m, other)
}

// MultiplyMatrices computes a * b, storing the result in this matrix.
// Returns pointer to this updated matrix.
func (m *Matrix4) MultiplyMatrices(a, b *Matrix4) *Matrix4 {

	a11 := a[0]
	a12 := a[4]
	a13 := a[8]
	a14 := a[12]
	a21 := a[1]
	a22 := a[5]
	a23 := a[9]
	a24 := a[13]
	a31 := a[2]
	a32 := a[6]
	a33 := a[10]
	a34 := a[14]
	a41 := a[3]
	a42 := a[7]
	a43 := a[11]
	a44 := a[15]

	b11 := b[0]
	b12 := b[4]
	b13 := b[8]
	b14 := b[12]
	b21 := b[1]
	b22 := b[5]
	b23 := b[9]
	b24 := b[13]
	b31 := b[2]
	b32 := b[6]
	b33 := b[10]
	b34 := b[14]
	b41 := b[3]
	b42 := b[7]
	b43 := b[11]
	b44 := b[15]

	m[0] = a11*b11 + a12*b21 + a13*b31 + a14*b41
	m[4] = a11*b12 + a12*b22 + a13*b32 + a14*b42
	m[8] = a11*b13 + a12*b23 + a13*b33 + a14*b43
	m[12] = a11*b14 + a12*b24 + a13*b34 + a14*b44

	m[1] = a21*b11 + a22*b21 + a23*b31 + a24*b41
	m[5] = a21*b12 + a22*b22 + a23*b32 + a24*b42
	m[9] = a21*b13 + a22*b23 + a23*b33 + a24*b43
	m[13] = a21*b14 + a22*b24 + a23*b34 + a24*b44

	m[2] = a31*b11 + a32*b21 + a33*b31 + a34*b41
	m[6] = a31*b12 + a32*b22 + a33*b32 + a34*b42
	m[10] = a31*b13 + a32*b23 + a33*b33 + a34*b43
	m[14] = a31*b14 + a32*b24 + a33*b34 + a34*b44

	m[3] = a41*b11 + a42*b21 + a43*b31 + a44*b41
	m[7] = a41*b12 + a42*b22 + a43*b32 + a44*b42
	m[11] = a41*b13 + a42*b23 + a43*b33 + a44*b43
	m[15] = a41*b14 + a42*b24 + a43*b34 + a44*b44

	return m
}

// MultiplyScalar multiplies each element of this matrix by the specified scalar.
// Returns pointer to this updated matrix.
func (m *Matrix4) MultiplyScalar(s float32) *Matrix4 {

	m[0] *= s
	m[4] *= s
	m[8] *= s
	m[12] *= s
	m[1] *= s
	m[5] *= s
	m[9] *= s
	m[13] *= s
	m[2] *= s
	m[6] *= s
	m[10] *= s
	m[14] *= s
	m[3] *= s
	m[7] *= s
	m[11] *= s
	m[15] *= s
	return m
}

// ApplyToVector3Array multiplies length vectors in the array starting at offset by this matrix.
// Returns pointer to the updated array.
// This matrix is unchanged.
func (m *Matrix4) ApplyToVector3Array(array []float32, offset int, length int) []float32 {

	var v1 Vector3
	j := offset
	for i := 0; i < length; i += 3 {
		v1.X = array[j]
		v1.Y = array[j+1]
		v1.Z = array[j+2]

		v1.ApplyMatrix4(m)

		array[j] = v1.X
		array[j+1] = v1.Y
		array[j+2] = v1.Z
		j += 3
	}
	return array
}

// Determinant calculates and returns the determinat of this matrix.
func (m *Matrix4) Determinant() float32 {

	n11 := m[0]
	n12 := m[4]
	n13 := m[8]
	n14 := m[12]
	n21 := m[1]
	n22 := m[5]
	n23 := m[9]
	n24 := m[13]
	n31 := m[2]
	n32 := m[6]
	n33 := m[10]
	n34 := m[14]
	n41 := m[3]
	n42 := m[7]
	n43 := m[11]
	n44 := m[15]

	return n41*(+n14*n23*n32-n13*n24*n32-n14*n22*n33+n12*n24*n33+n13*n22*n34-n12*n23*n34) +
		n42*(+n11*n23*n34-n11*n24*n33+n14*n21*n33-n13*n21*n34+n13*n24*n31-n14*n23*n31) +
		n43*(+n11*n24*n32-n11*n22*n34-n14*n21*n32+n12*n21*n34+n14*n22*n31-n12*n24*n31) +
		n44*(-n13*n22*n31-n11*n23*n32+n11*n22*n33+n13*n21*n32-n12*n21*n33+n12*n23*n31)

}

// Transpose transposes this matrix.
// Returns pointer to this updated matrix.
func (m *Matrix4) Transpose() *Matrix4 {

	m[1], m[4] = m[4], m[1]
	m[2], m[8] = m[8], m[2]
	m[6], m[9] = m[9], m[6]
	m[3], m[12] = m[12], m[3]
	m[7], m[13] = m[13], m[7]
	m[11], m[14] = m[14], m[11]
	return m
}

// SetPosition sets this transformation matrix position fields from the specified vector v.
// Returns pointer to this updated matrix.
func (m *Matrix4) SetPosition(v *Vector3) *Matrix4 {

	m[12] = v.X
	m[13] = v.Y
	m[14] = v.Z
	return m
}

// GetInverse sets this matrix to the inverse of the src matrix.
// If the src matrix cannot be inverted returns error and
// sets this matrix to the identity matrix.
func (m *Matrix4) GetInverse(src *Matrix4) error {

	n11 := src[0]
	n12 := src[4]
	n13 := src[8]
	n14 := src[12]
	n21 := src[1]
	n22 := src[5]
	n23 := src[9]
	n24 := src[13]
	n31 := src[2]
	n32 := src[6]
	n33 := src[10]
	n34 := src[14]
	n41 := src[3]
	n42 := src[7]
	n43 := src[11]
	n44 := src[15]

	t11 := n23*n34*n42 - n24*n33*n42 + n24*n32*n43 - n22*n34*n43 - n23*n32*n44 + n22*n33*n44
	t12 := n14*n33*n42 - n13*n34*n42 - n14*n32*n43 + n12*n34*n43 + n13*n32*n44 - n12*n33*n44
	t13 := n13*n24*n42 - n14*n23*n42 + n14*n22*n43 - n12*n24*n43 - n13*n22*n44 + n12*n23*n44
	t14 := n14*n23*n32 - n13*n24*n32 - n14*n22*n33 + n12*n24*n33 + n13*n22*n34 - n12*n23*n34

	det := n11*t11 + n21*t12 + n31*t13 + n41*t14

	if det == 0 {
		m.Identity()
		return errors.New("cannot invert matrix")
	}

	m[0] = t11
	m[1] = n24*n33*n41 - n23*n34*n41 - n24*n31*n43 + n21*n34*n43 + n23*n31*n44 - n21*n33*n44
	m[2] = n22*n34*n41 - n24*n32*n41 + n24*n31*n42 - n21*n34*n42 - n22*n31*n44 + n21*n32*n44
	m[3] = n23*n32*n41 - n22*n33*n41 - n23*n31*n42 + n21*n33*n42 + n22*n31*n43 - n21*n32*n43
	m[4] = t12
	m[5] = n13*n34*n41 - n14*n33*n41 + n14*n31*n43 - n11*n34*n43 - n13*n31*n44 + n11*n33*n44
	m[6] = n14*n32*n41 - n12*n34*n41 - n14*n31*n42 + n11*n34*n42 + n12*n31*n44 - n11*n32*n44
	m[7] = n12*n33*n41 - n13*n32*n41 + n13*n31*n42 - n11*n33*n42 - n12*n31*n43 + n11*n32*n43
	m[8] = t13
	m[9] = n14*n23*n41 - n13*n24*n41 - n14*n21*n43 + n11*n24*n43 + n13*n21*n44 - n11*n23*n44
	m[10] = n12*n24*n41 - n14*n22*n41 + n14*n21*n42 - n11*n24*n42 - n12*n21*n44 + n11*n22*n44
	m[11] = n13*n22*n41 - n12*n23*n41 - n13*n21*n42 + n11*n23*n42 + n12*n21*n43 - n11*n22*n43
	m[12] = t14
	m[13] = n13*n24*n31 - n14*n23*n31 + n14*n21*n33 - n11*n24*n33 - n13*n21*n34 + n11*n23*n34
	m[14] = n14*n22*n31 - n12*n24*n31 - n14*n21*n32 + n11*n24*n32 + n12*n21*n34 - n11*n22*n34
	m[15] = n12*n23*n31 - n13*n22*n31 + n13*n21*n32 - n11*n23*n32 - n12*n21*n33 + n11*n22*n33

	m.MultiplyScalar(1.0 / det)

	return nil
}

// Scale multiply the first column of this matrix by the vector X component,
// the second column by the vector Y component and the third column by
// the vector Z component. The matrix fourth column is unchanged.
// Returns pointer to this updated matrix.
func (m *Matrix4) Scale(v *Vector3) *Matrix4 {

	m[0] *= v.X
	m[4] *= v.Y
	m[8] *= v.Z
	m[1] *= v.X
	m[5] *= v.Y
	m[9] *= v.Z
	m[2] *= v.X
	m[6] *= v.Y
	m[10] *= v.Z
	m[3] *= v.X
	m[7] *= v.Y
	m[11] *= v.Z
	return m
}

// GetMaxScaleOnAxis returns the maximum scale value of the 3 axes.
func (m *Matrix4) GetMaxScaleOnAxis() float32 {

	scaleXSq := m[0]*m[0] + m[1]*m[1] + m[2]*m[2]
	scaleYSq := m[4]*m[4] + m[5]*m[5] + m[6]*m[6]
	scaleZSq := m[8]*m[8] + m[9]*m[9] + m[10]*m[10]
	return Sqrt(Max(scaleXSq, Max(scaleYSq, scaleZSq)))
}

// MakeTranslation sets this matrix to a translation matrix from the specified x, y and z values.
// Returns pointer to this updated matrix.
func (m *Matrix4) MakeTranslation(x, y, z float32) *Matrix4 {

	m.Set(
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	)
	return m
}

// MakeRotationX sets this matrix to a rotation matrix of angle theta around the X axis.
// Returns pointer to this updated matrix.
func (m *Matrix4) MakeRotationX(theta float32) *Matrix4 {

	c := Cos(theta)
	s := Sin(theta)

	m.Set(
		1, 0, 0, 0,
		0, c, -s, 0,
		0, s, c, 0,
		0, 0, 0, 1,
	)
	return m
}

// MakeRotationY sets this matrix to a rotation matrix of angle theta around the Y axis.
// Returns pointer to this updated matrix.
func (m *Matrix4) MakeRotationY(theta float32) *Matrix4 {

	c := Cos(theta)
	s := Sin(theta)
	m.Set(
		c, 0, s, 0,
		0, 1, 0, 0,
		-s, 0, c, 0,
		0, 0, 0, 1,
	)
	return m
}

// MakeRotationZ sets this matrix to a rotation matrix of angle theta around the Z axis.
// Returns pointer to this updated matrix.
func (m *Matrix4) MakeRotationZ(theta float32) *Matrix4 {

	c := Cos(theta)
	s := Sin(theta)
	m.Set(
		c, -s, 0, 0,
		s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	return m
}

// MakeRotationAxis sets this matrix to a rotation matrix of the specified angle around the specified axis.
// Returns pointer to this updated matrix.
func (m *Matrix4) MakeRotationAxis(axis *Vector3, angle float32) *Matrix4 {

	c := Cos(angle)
	s := Sin(angle)
	t := 1 - c
	x := axis.X
	y := axis.Y
	z := axis.Z
	tx := t * x
	ty := t * y
	m.Set(
		tx*x+c, tx*y-s*z, tx*z+s*y, 0,
		tx*y+s*z, ty*y+c, ty*z-s*x, 0,
		tx*z-s*y, ty*z+s*x, t*z*z+c, 0,
		0, 0, 0, 1,
	)
	return m
}

// MakeScale sets this matrix to a scale transformation matrix using the specified x, y and z values.
// Returns pointer to this updated matrix.
func (m *Matrix4) MakeScale(x, y, z float32) *Matrix4 {

	m.Set(
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, 1,
	)
	return m
}

// Compose sets this matrix to a transformation matrix for the specified position,
// rotation specified by the quaternion and scale.
// Returns pointer to this updated matrix.
func (m *Matrix4) Compose(position *Vector3, quaternion *Quaternion, scale *Vector3) *Matrix4 {

	m.MakeRotationFromQuaternion(quaternion)
	m.Scale(scale)
	m.SetPosition(position)
	return m
}

// Decompose updates the position vector, quaternion and scale from this transformation matrix.
// Returns pointer to this unchanged matrix.
func (m *Matrix4) Decompose(position *Vector3, quaternion *Quaternion, scale *Vector3) *Matrix4 {

	var vector Vector3
	var matrix = *m

	position.X = m[12]
	position.Y = m[13]
	position.Z = m[14]

	scale.X = vector.Set(m[0], m[1], m[2]).Length()
	scale.Y = vector.Set(m[4], m[5], m[6]).Length()
	scale.Z = vector.Set(m[8], m[9], m[10]).Length()

	// If determinant is negative, we need to invert one scale
	det := m.Determinant()
	if det < 0 {
		scale.X = -scale.X
	}

	// Scale the rotation part
	invSX := 1 / scale.X
	invSY := 1 / scale.Y
	invSZ := 1 / scale.Z

	matrix[0] *= invSX
	matrix[1] *= invSX
	matrix[2] *= invSX

	matrix[4] *= invSY
	matrix[5] *= invSY
	matrix[6] *= invSY

	matrix[8] *= invSZ
	matrix[9] *= invSZ
	matrix[10] *= invSZ

	quaternion.SetFromRotationMatrix(&matrix)

	return m
}

// MakeFrustum sets this matrix to a projection frustum matrix bounded by the specified planes.
// Returns pointer to this updated matrix.
func (m *Matrix4) MakeFrustum(left, right, bottom, top, near, far float32) *Matrix4 {

	m[0] = 2 * near / (right - left)
	m[1] = 0
	m[2] = 0
	m[3] = 0
	m[4] = 0
	m[5] = 2 * near / (top - bottom)
	m[6] = 0
	m[7] = 0
	m[8] = (right + left) / (right - left)
	m[9] = (top + bottom) / (top - bottom)
	m[10] = -(far + near) / (far - near)
	m[11] = -1
	m[12] = 0
	m[13] = 0
	m[14] = -(2 * far * near) / (far - near)
	m[15] = 0
	return m
}

// MakePerspective sets this matrix to a perspective projection matrix
// with the specified field of view in degrees,
// aspect ratio (width/height) and near and far planes.
// Returns pointer to this updated matrix.
func (m *Matrix4) MakePerspective(fov, aspect, near, far float32) *Matrix4 {

	ymax := near * Tan(DegToRad(fov*0.5))
	ymin := -ymax
	xmin := ymin * aspect
	xmax := ymax * aspect
	return m.MakeFrustum(xmin, xmax, ymin, ymax, near, far)
}

// MakeOrthographic sets this matrix to an orthographic projection matrix
// bounded by the specified planes.
// Returns pointer to this updated matrix.
func (m *Matrix4) MakeOrthographic(left, right, top, bottom, near, far float32) *Matrix4 {

	w := right - left
	h := top - bottom
	p := far - near

	x := (right + left) / w
	y := (top + bottom) / h
	z := (far + near) / p

	m[0] = 2 / w
	m[4] = 0
	m[8] = 0
	m[12] = -x
	m[1] = 0
	m[5] = 2 / h
	m[9] = 0
	m[13] = -y
	m[2] = 0
	m[6] = 0
	m[10] = -2 / p
	m[14] = -z
	m[3] = 0
	m[7] = 0
	m[11] = 0
	m[15] = 1
	return m
}

// FromArray set this matrix elements from the array starting at offset.
// Returns pointer to this updated matrix.
func (m *Matrix4) FromArray(array []float32, offset int) *Matrix4 {

	copy(m[:], array[offset:offset+16])
	return m
}

// ToArray copies this matrix elements to array starting at offset.
// Returns pointer to the updated array.
func (m *Matrix4) ToArray(array []float32, offset int) []float32 {

	copy(array[offset:], m[:])
	return array
}

// Clone creates and returns a pointer to a copy of this matrix.
func (m *Matrix4) Clone() *Matrix4 {

	var cloned Matrix4
	cloned = *m
	return &cloned
}
