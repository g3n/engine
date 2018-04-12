// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

// FillLayout is the simple layout where the assigned panel "fills" its parent in the specified dimension(s)
type FillLayout struct {
	width  bool
	height bool
}

// NewFillLayout creates and returns a pointer of a new fill layout
func NewFillLayout(width, height bool) *FillLayout {

	f := new(FillLayout)
	f.width = width
	f.height = height
	return f
}

// Recalc is called by the panel which has this layout
func (f *FillLayout) Recalc(ipan IPanel) {

	parent := ipan.GetPanel()
	children := parent.Children()
	if len(children) == 0 {
		return
	}
	child := children[0].(IPanel).GetPanel()

	if f.width {
		child.SetWidth(parent.ContentWidth())
	}
	if f.height {
		child.SetHeight(parent.ContentHeight())
	}
}
