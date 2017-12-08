// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package gui

// GridLayout is a panel layout which arranges its children in a rectangular grid.
// It is necessary to set the number of columns of the grid when the layout is created.
// The panel's child elements are positioned in the grid cells accordingly to the
// order they were added to the panel.
// The height of each row is determined by the height of the heightest child in the row.
// The width of each column is determined by the width of the widest child in the column
type GridLayout struct {
	columns []colInfo
	alignh  Align   // global cell horizontal alignment
	alignv  Align   // global cell vertical alignment
	spaceh  float32 // space between rows
	specev  float32 // space between columns
}

// GridLayoutParams describes layout parameter for an specific child
type GridLayoutParams struct {
	ColSpan int   // Number of additional columns to ocuppy to the right
	AlignH  Align // Vertical alignment
	AlignV  Align // Horizontal alignment
}

// colInfo keeps information about each grid column
type colInfo struct {
	width  float32 // column width
	alignh *Align  // optional column horizontal alignment
	alignv *Align  // optional column vertical alignment
}

// NewGridLayout creates and returns a pointer of a new grid layout
func NewGridLayout(ncols int) *GridLayout {

	if ncols <= 0 {
		panic("Invalid number of columns")
	}
	gl := new(GridLayout)
	gl.columns = make([]colInfo, ncols)
	return gl
}

// SetAlignV sets the vertical alignment for all the grid cells
// The alignment of an individual cell can be set by settings its layout parameters.
func (g *GridLayout) SetAlignV(align Align) {

	g.alignv = align
}

// SetAlignH sets the horizontal alignment for all the grid cells
// The alignment of an individual cell can be set by settings its layout parameters.
func (g *GridLayout) SetAlignH(align Align) {

	g.alignh = align
}

// SetColAlignV sets the vertical alignment for all the cells of
// the specified column. The function panics if the supplied column is invalid
func (g *GridLayout) SetColAlignV(col int, align Align) {

	if col < 0 || col >= len(g.columns) {
		panic("Invalid column")
	}
	if g.columns[col].alignv == nil {
		g.columns[col].alignv = new(Align)
	}
	*g.columns[col].alignv = align
}

// SetColAlignH sets the horizontal alignment for all the cells of
// the specified column. The function panics if the supplied column is invalid.
func (g *GridLayout) SetColAlignH(col int, align Align) {

	if col < 0 || col >= len(g.columns) {
		panic("Invalid column")
	}
	if g.columns[col].alignh == nil {
		g.columns[col].alignh = new(Align)
	}
	*g.columns[col].alignh = align
}

// Recalc sets the position and sizes of all of the panel's children.
// It is normally called by the parent panel when its size changes or
// a child is added or removed.
func (g *GridLayout) Recalc(ipan IPanel) {

	type cell struct {
		panel  *Panel
		params *GridLayoutParams
	}

	type row struct {
		cells  []*cell
		height float32
	}

	// Builds array of child rows
	pan := ipan.GetPanel()
	var irow int
	var icol int
	rows := []row{}
	for _, node := range pan.Children() {
		// Ignore invisible child
		child := node.(IPanel).GetPanel()
		if !child.Visible() {
			continue
		}
		// Checks child layout params, if supplied
		ip := child.layoutParams
		var params *GridLayoutParams
		var ok bool
		if ip != nil {
			params, ok = child.layoutParams.(*GridLayoutParams)
			if !ok {
				panic("layoutParams is not GridLayoutParams")
			}
		}
		// If first column, creates row and appends to rows
		if icol == 0 {
			var r row
			r.cells = make([]*cell, len(g.columns))
			rows = append(rows, r)
		}
		// Set current child panel to current cells
		rows[irow].cells[icol] = &cell{child, params}
		// Checks child panel colspan layout params
		coljump := 1
		if params != nil && params.ColSpan > 0 {
			coljump += params.ColSpan
		}
		// Updates next cell column and row
		icol += coljump
		if icol >= len(g.columns) {
			irow++
			icol = 0
		}
	}

	// Resets columns widths
	for i := 0; i < len(g.columns); i++ {
		g.columns[i].width = 0
	}

	// Sets the height of each row to the height of the heightest child
	// Sets the width of each column to the width of the widest column
	for i := 0; i < len(rows); i++ {
		r := &rows[i]
		r.height = 0
		for ci, cell := range r.cells {
			if cell == nil {
				continue
			}
			if cell.panel.Height() > r.height {
				r.height = cell.panel.Height()
			}
			if cell.panel.Width() > g.columns[ci].width {
				g.columns[ci].width = cell.panel.Width()
			}
		}
	}

	// Position each child panel in the parent panel
	var celly float32
	for _, r := range rows {
		var cellx float32
		for ci, cell := range r.cells {
			if cell == nil {
				continue
			}
			colWidth := g.columns[ci].width
			colspan := 0
			// Default grid cell alignment
			alignv := g.alignv
			alignh := g.alignh
			// If column has aligment:
			if g.columns[ci].alignv != nil {
				alignv = *g.columns[ci].alignv
			}
			if g.columns[ci].alignh != nil {
				alignh = *g.columns[ci].alignh
			}
			// If cell has layout parameters:
			if cell.params != nil {
				alignv = cell.params.AlignV
				alignv = cell.params.AlignH
				colspan = cell.params.ColSpan
			}
			// Determines child panel horizontal position
			px := cellx
			switch alignh {
			case AlignNone:
			case AlignLeft:
			case AlignRight:
				px += float32(colWidth) - cell.panel.Width()
			case AlignCenter:
				px += (float32(colWidth) - cell.panel.Width()) / 2
			default:
				panic("Invalid horizontal alignment")
			}
			// Determines child panel vertical position
			py := celly
			switch alignv {
			case AlignNone:
			case AlignTop:
			case AlignBottom:
				py += r.height - cell.panel.Height()
			case AlignCenter:
				py += (r.height - cell.panel.Height()) / 2
			default:
				panic("Invalid vertical alignment")
			}
			// Sets child panel position
			cell.panel.SetPosition(px, py)
			// Advances to next row cell considering colspan
			for i := ci; i < ci+colspan+1; i++ {
				if i >= len(g.columns) {
					break
				}
				cellx += g.columns[i].width
			}
		}
		celly += r.height
	}
}
