// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math32

import "errors"

// Matrix3 is 3x3 matrix organized internally as column matrix
type Matrix3 [9]float32

// NewMatrix3 creates and returns a pointer to a new Matrix3
// initialized as the identity matrix.
func NewMatrix3() *Matrix3 {

	var m Matrix3
	m.Identity()
	return &m
}

// Set sets all the elements of the matrix row by row starting at row1, column1,
// row1, column2, row1, column3 and so forth.
// Returns the pointer to this updated Matrix.
func (m *Matrix3) Set(n11, n12, n13, n21, n22, n23, n31, n32, n33 float32) *Matrix3 {

	m[0] = n11
	m[3] = n12
	m[6] = n13
	m[1] = n21
	m[4] = n22
	m[7] = n23
	m[2] = n31
	m[5] = n32
	m[8] = n33
	return m
}

// SetFromMatrix4 sets the matrix elements based on a Matrix4.
func (m *Matrix3) SetFromMatrix4(src *Matrix4) *Matrix3 {

	m.Set(
		src[0], src[4], src[8],
		src[1], src[5], src[9],
		src[2], src[6], src[10],
	)
	return m
}

// Identity sets this matrix as the identity matrix.
// Returns the pointer to this updated matrix.
func (m *Matrix3) Identity() *Matrix3 {

	m.Set(
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	)
	return m
}

// Zero sets this matrix as the zero matrix.
// Returns the pointer to this updated matrix.
func (m *Matrix3) Zero() *Matrix3 {

	m.Set(
		0, 0, 0,
		0, 0, 0,
		0, 0, 0,
	)
	return m
}

// Copy copies src matrix into this one.
// Returns the pointer to this updated matrix.
func (m *Matrix3) Copy(src *Matrix3) *Matrix3 {

	*m = *src
	return m
}

// MakeRotationFromQuaternion sets this matrix as a rotation matrix from the specified quaternion.
// Returns pointer to this updated matrix.
func (m *Matrix3) MakeRotationFromQuaternion(q *Quaternion) *Matrix3 {

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
	m[3] = xy - wz
	m[6] = xz + wy

	m[1] = xy + wz
	m[4] = 1 - (xx + zz)
	m[7] = yz - wx

	m[2] = xz - wy
	m[5] = yz + wx
	m[8] = 1 - (xx + yy)

	return m

}

// ApplyToVector3Array multiplies length vectors in the array starting at offset by this matrix.
// Returns pointer to the updated array.
// This matrix is unchanged.
func (m *Matrix3) ApplyToVector3Array(array []float32, offset int, length int) []float32 {

	var v1 Vector3
	j := offset
	for i := 0; i < length; i += 3 {
		v1.X = array[j]
		v1.Y = array[j+1]
		v1.Z = array[j+2]
		v1.ApplyMatrix3(m)
		array[j] = v1.X
		array[j+1] = v1.Y
		array[j+2] = v1.Z
	}
	return array
}

// Multiply multiply this matrix by the other matrix
// Returns pointer to this updated matrix.
func (m *Matrix3) Multiply(other *Matrix3) *Matrix3 {

	return m.MultiplyMatrices(m, other)
}

// MultiplyMatrices multiply matrix a by b storing the result in this matrix.
// Returns pointer to this updated matrix.
func (m *Matrix3) MultiplyMatrices(a, b *Matrix3) *Matrix3 {

	a11 := a[0]
	a12 := a[3]
	a13 := a[6]
	a21 := a[1]
	a22 := a[4]
	a23 := a[7]
	a31 := a[2]
	a32 := a[5]
	a33 := a[8]

	b11 := b[0]
	b12 := b[3]
	b13 := b[6]
	b21 := b[1]
	b22 := b[4]
	b23 := b[7]
	b31 := b[2]
	b32 := b[5]
	b33 := b[8]

	m[0] = a11*b11 + a12*b21 + a13*b31
	m[3] = a11*b12 + a12*b22 + a13*b32
	m[6] = a11*b13 + a12*b23 + a13*b33

	m[1] = a21*b11 + a22*b21 + a23*b31
	m[4] = a21*b12 + a22*b22 + a23*b32
	m[7] = a21*b13 + a22*b23 + a23*b33

	m[2] = a31*b11 + a32*b21 + a33*b31
	m[5] = a31*b12 + a32*b22 + a33*b32
	m[8] = a31*b13 + a32*b23 + a33*b33

	return m
}

// MultiplyScalar multiplies each of this matrix's components by the specified scalar.
// Returns pointer to this updated matrix.
func (m *Matrix3) MultiplyScalar(s float32) *Matrix3 {

	m[0] *= s
	m[3] *= s
	m[6] *= s
	m[1] *= s
	m[4] *= s
	m[7] *= s
	m[2] *= s
	m[5] *= s
	m[8] *= s
	return m
}

// ScaleColumns multiplies the matrix columns by the vector components.
// This can be used when multiplying this matrix by a diagonal matrix if we store the diagonal components as a vector.
// Returns pointer to this updated matrix.
func (m *Matrix3) ScaleColumns(v *Vector3) *Matrix3 {

	m[0] *= v.X
	m[1] *= v.X
	m[2] *= v.X
	m[3] *= v.Y
	m[4] *= v.Y
	m[5] *= v.Y
	m[6] *= v.Z
	m[7] *= v.Z
	m[8] *= v.Z
	return m
}

// Determinant calculates and returns the determinant of this matrix.
func (m *Matrix3) Determinant() float32 {

	return m[0]*m[4]*m[8] -
		m[0]*m[5]*m[7] -
		m[1]*m[3]*m[8] +
		m[1]*m[5]*m[6] +
		m[2]*m[3]*m[7] -
		m[2]*m[4]*m[6]
}

// GetInverse sets this matrix to the inverse of the src matrix.
// If the src matrix cannot be inverted returns error and
// sets this matrix to the identity matrix.
func (m *Matrix3) GetInverse(src *Matrix3) error {

	n11 := src[0]
	n21 := src[1]
	n31 := src[2]
	n12 := src[3]
	n22 := src[4]
	n32 := src[5]
	n13 := src[6]
	n23 := src[7]
	n33 := src[8]

	t11 := n33*n22 - n32*n23
	t12 := n32*n13 - n33*n12
	t13 := n23*n12 - n22*n13

	det := n11*t11 + n21*t12 + n31*t13

	// no inverse
	if det == 0 {
		m.Identity()
		return errors.New("cannot invert matrix")
	}

	detInv := 1 / det

	m[0] = t11 * detInv
	m[1] = (n31*n23 - n33*n21) * detInv
	m[2] = (n32*n21 - n31*n22) * detInv
	m[3] = t12 * detInv
	m[4] = (n33*n11 - n31*n13) * detInv
	m[5] = (n31*n12 - n32*n11) * detInv
	m[6] = t13 * detInv
	m[7] = (n21*n13 - n23*n11) * detInv
	m[8] = (n22*n11 - n21*n12) * detInv

	return nil
}

// Transpose transposes this matrix.
// Returns pointer to this updated matrix.
func (m *Matrix3) Transpose() *Matrix3 {

	m[1], m[3] = m[3], m[1]
	m[2], m[6] = m[6], m[2]
	m[5], m[7] = m[7], m[5]
	return m
}

// GetNormalMatrix set this matrix to the matrix to transform the normal vectors
// from the src matrix to transform the vertices.
// If the src matrix cannot be inverted returns error.
func (m *Matrix3) GetNormalMatrix(src *Matrix4) error {

	m.SetFromMatrix4(src)
	err := m.GetInverse(m)
	m.Transpose()
	return err
}

// FromArray set this matrix array starting at offset.
// Returns pointer to this updated matrix.
func (m *Matrix3) FromArray(array []float32, offset int) *Matrix3 {

	copy(m[:], array[offset:offset+9])
	return m
}

// ToArray copies this matrix to array starting at offset.
// Returns pointer to the updated array.
func (m *Matrix3) ToArray(array []float32, offset int) []float32 {

	copy(array[offset:], m[:])
	return array
}

// Clone creates and returns a pointer to a copy of this matrix.
func (m *Matrix3) Clone() *Matrix3 {

	var cloned Matrix3
	cloned = *m
	return &cloned
}
