// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collision

// Matrix is a triangular collision matrix indicating which pairs of bodies are colliding.
type Matrix [][]bool

// NewMatrix creates and returns a pointer to a new collision Matrix.
func NewMatrix() Matrix {

	return make([][]bool, 0)
}

// Set sets whether i and j are colliding.
func (m *Matrix) Set(i, j int, val bool) {

	var s, l int
	if i < j {
		s = i
		l = j
	} else {
		s = j
		l = i
	}
	diff := s + 1 - len(*m)
	if diff > 0 {
		for i := 0; i < diff; i++ {
			*m = append(*m, make([]bool, 0))
		}
	}
	for idx := range *m {
		diff = l + 1 - len((*m)[idx]) - idx
		if diff > 0 {
			for i := 0; i < diff; i++ {
				(*m)[idx] = append((*m)[idx], false)
			}
		}
	}
	(*m)[s][l-s] = val
}

// Get returns whether i and j are colliding.
func (m *Matrix) Get(i, j int) bool {

	var s, l int
	if i < j {
		s = i
		l = j
	} else {
		s = j
		l = i
	}
	return (*m)[s][l-s]
}

// Reset clears all values.
func (m *Matrix) Reset() {

	*m = make([][]bool, 0)
}
