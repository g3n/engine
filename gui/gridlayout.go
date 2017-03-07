// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type GridLayout struct {
	columns int
}

type GridLayoutParams struct {
	Row     int   // grid layout row number from 0
	Col     int   // grid layout column number from 0
	ColSpan int   // number of additional columns to ocuppy to the right
	AlignH  Align // vertical alignment
	AlignV  Align // horizontal alignment
}

// NewGridLayout creates and returns a pointer of a new grid layout
func NewGridLayout() *GridLayout {

	g := new(GridLayout)
	return g
}

func (g *GridLayout) Recalc(ipan IPanel) {

	type element struct {
		panel  *Panel
		params *GridLayoutParams
	}
	rows := 0
	cols := 0
	items := []element{}

	pan := ipan.GetPanel()
	for _, obj := range pan.Children() {
		// Get child panel
		child := obj.(IPanel).GetPanel()
		// Ignore if not visible
		if !child.Visible() {
			continue
		}
		// Ignore if no layout params
		if child.layoutParams == nil {
			continue
		}
		// Checks layout params
		params, ok := child.layoutParams.(*GridLayoutParams)
		if !ok {
			panic("layoutParams is not GridLayoutParams")
		}
		if params.Row >= rows {
			rows = params.Row + 1
		}
		if params.Col >= cols {
			cols = params.Col + 1
		}
		items = append(items, element{child, params})
	}
	// Check limits
	if rows > 100 {
		panic("Element row outsize limits")
	}
	if cols > 100 {
		panic("Element column outsize limits")
	}

	// Determine row and column maximum sizes
	colSizes := make([]int, cols)
	rowSizes := make([]int, rows)
	for _, el := range items {
		width := el.panel.Width()
		height := el.panel.Height()
		if int(width) > colSizes[el.params.Col] {
			colSizes[el.params.Col] = int(width)
		}
		if int(height) > rowSizes[el.params.Row] {
			rowSizes[el.params.Row] = int(height)
		}
	}

	// Determine row and column starting positions
	colStart := make([]int, cols)
	rowStart := make([]int, rows)
	for i := 1; i < len(colSizes); i++ {
		colStart[i] = colStart[i-1] + colSizes[i-1]
	}
	for i := 1; i < len(rowSizes); i++ {
		rowStart[i] = rowStart[i-1] + rowSizes[i-1]
	}

	// Position the elements
	for _, el := range items {
		row := el.params.Row
		col := el.params.Col
		cellHeight := rowSizes[row]
		// Current cell width
		cellWidth := 0
		for c := 0; c <= el.params.ColSpan; c++ {
			pos := col + c
			if pos >= len(colSizes) {
				break
			}
			cellWidth += colSizes[pos]
		}
		rstart := float32(rowStart[row])
		cstart := float32(colStart[col])
		// Horizontal alignment
		var dx float32 = 0
		switch el.params.AlignH {
		case AlignNone:
		case AlignLeft:
		case AlignRight:
			dx = float32(cellWidth) - el.panel.width
		case AlignCenter:
			dx = (float32(cellWidth) - el.panel.width) / 2
		default:
			panic("Invalid horizontal alignment")
		}
		// Vertical alignment
		var dy float32 = 0
		switch el.params.AlignV {
		case AlignNone:
		case AlignTop:
		case AlignBottom:
			dy = float32(cellHeight) - el.panel.height
		case AlignCenter:
			dy = (float32(cellHeight) - el.panel.height) / 2
		default:
			panic("Invalid vertical alignment")
		}
		el.panel.SetPosition(cstart+dx, rstart+dy)
	}
}
