package math32

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatrix4_MultiplyVector4(t *testing.T) {
	tests := []struct {
		matrix *Matrix4
		vector *Vector4
		expected *Vector4
	}{
		{
			vector:NewVector4(0,0,0,0),
			matrix:NewMatrix4().Set(0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0),
			expected:NewVector4(0,0,0,0),
		},
		{
			vector:NewVector4(1,1,1,1),
			matrix:NewMatrix4().Set(1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1),
			expected:NewVector4(4,4,4,4),
		},
		{
			vector:NewVector4(1,1,1,1),
			matrix:NewMatrix4().Set(1,1,1,1,2,2,2,2,3,3,3,3,4,4,4, 4),
			expected:NewVector4(4,8,12,16),
		},
		{
			vector:NewVector4(1,2,3,4),
			matrix:NewMatrix4().Set(1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1),
			expected:NewVector4(10,10,10,10),
		},
		{
			vector:NewVector4(1,2,3,4),
			matrix:NewMatrix4().Set(1,2,3,4,1,2,3,4,1,2,3,4,1,2,3,4),
			expected:NewVector4(30,30,30,30),
		},
		{
			vector:NewVector4(1,1,1,1),
			matrix:NewMatrix4().Set(2,2,2,2,1,1,1,1,1,1,1,1,1,1,1,1),
			expected:NewVector4(8,4,4,4),
		},
		{
			vector:NewVector4(2,1,1,1),
			matrix:NewMatrix4().Set(1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1),
			expected:NewVector4(5,5,5,5),
		},
		{
			vector:NewVector4(1,1,1,1),
			matrix:NewMatrix4().Set(1,1,1,1,0,0,0,0,0,0, 0,0,0,0,0,0),
			expected:NewVector4(4,0,0,0),
		},
		{
			vector:NewVector4(1,1,1,1),
			matrix:NewMatrix4().Set(1,0,0,0,1,0,0,0,1,0,0,0,1,0,0,0),
			expected:NewVector4(1,1,1,1),
		},
	}

	for i, test := range tests {
		actual := test.matrix.MultiplyVector4(test.vector)
		assert.Equalf(t, test.expected, actual, "Failed test %v", i)
	}
}
