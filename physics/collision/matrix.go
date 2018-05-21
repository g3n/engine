// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package collision implements collision related algorithms and data structures.
package collision

// Matrix is a triangular collision matrix indicating which pairs of bodies are colliding.
type Matrix struct {
	col [][]bool
}

// NewMatrix creates and returns a pointer to a new collision Matrix.
func NewMatrix() *Matrix {

	m := new(Matrix)
	m.col = make([][]bool, 0)
	return m
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
	diff := s + 1 - len(m.col)
	if diff > 0 {
		for i := 0; i < diff; i++ {
			m.col = append(m.col, make([]bool,0))
		}
	}
	for idx := range m.col {
		diff = l + 1 - len(m.col[idx]) - idx
		if diff > 0 {
			for i := 0; i < diff; i++ {
				m.col[idx] = append(m.col[idx], false)
			}
		}
	}
	m.col[s][l-s] = val
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
	return m.col[s][l-s]
}