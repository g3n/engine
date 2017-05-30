// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import ()

//
// Table implements a panel which can contains row and columns of child panels
//
type Table struct {
	Panel // Embedded panel
}

// NewTable creates and returns a pointer to a new Table with the
// specified initial width and height
func NewTable(width, height float32) *Table {

	t := new(Table)
	t.Panel.Initialize(width, height)

	return t
}
