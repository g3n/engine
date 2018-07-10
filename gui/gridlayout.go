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
	pan     IPanel    // parent panel
	columns []colInfo // columns alignment info
	alignh  Align     // global cell horizontal alignment
	alignv  Align     // global cell vertical alignment
	expandh bool      // expand horizontally flag
	expandv bool      // expand vertically flag
}

// GridLayoutParams describes layout parameter for an specific child
type GridLayoutParams struct {
	ColSpan int   // Number of additional columns to ocuppy to the right
	AlignH  Align // Vertical alignment
	AlignV  Align // Horizontal alignment
}

// colInfo keeps information about each grid column
type colInfo struct {
	alignh *Align // optional column horizontal alignment
	alignv *Align // optional column vertical alignment
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
	g.Recalc(g.pan)
}

// SetAlignH sets the horizontal alignment for all the grid cells
// The alignment of an individual cell can be set by settings its layout parameters.
func (g *GridLayout) SetAlignH(align Align) {

	g.alignh = align
	g.Recalc(g.pan)
}

// SetExpandH sets it the columns should expand horizontally if possible
func (g *GridLayout) SetExpandH(expand bool) {

	g.expandh = expand
	g.Recalc(g.pan)
}

// SetExpandV sets it the rowss should expand vertically if possible
func (g *GridLayout) SetExpandV(expand bool) {

	g.expandv = expand
	g.Recalc(g.pan)
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
	g.Recalc(g.pan)
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
	g.Recalc(g.pan)
}

// Recalc sets the position and sizes of all of the panel's children.
// It is normally called by the parent panel when its size changes or
// a child is added or removed.
func (g *GridLayout) Recalc(ipan IPanel) {

	type cell struct {
		panel     *Panel           // pointer to cell panel
		params    GridLayoutParams // copy of params or default
		paramsDef bool             // true if parameters are default
	}

	type row struct {
		cells  []*cell // array of row cells
		height float32 // row height
	}

	// Saves the received panel
	g.pan = ipan
	if g.pan == nil {
		return
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
		var paramsDef bool
		if ip != nil {
			params, ok = child.layoutParams.(*GridLayoutParams)
			if !ok {
				panic("layoutParams is not GridLayoutParams")
			}
			paramsDef = false
		} else {
			params = &GridLayoutParams{}
			paramsDef = true

		}
		// If first column, creates row and appends to rows
		if icol == 0 {
			var r row
			r.cells = make([]*cell, len(g.columns))
			rows = append(rows, r)
		}
		// Set current child panel to current cells
		rows[irow].cells[icol] = &cell{child, *params, paramsDef}
		// Updates next cell column and row
		icol += 1 + params.ColSpan
		if icol >= len(g.columns) {
			irow++
			icol = 0
		}
	}

	// Sets the height of each row to the height of the heightest child
	// Sets the width of each column to the width of the widest column
	colWidths := make([]float32, len(g.columns))
	spanWidths := make([]float32, len(g.columns))
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
			// If this cell span columns compare with other span cell widths
			if cell.params.ColSpan > 0 {
				if cell.panel.Width() > spanWidths[ci] {
					spanWidths[ci] = cell.panel.Width()
				}
			} else {
				if cell.panel.Width() > colWidths[ci] {
					colWidths[ci] = cell.panel.Width()
				}
			}
		}
	}
	// The final width for each column is the maximum no span column width
	// but if it is zero, is the maximum span column width
	for i := 0; i < len(colWidths); i++ {
		if colWidths[i] == 0 {
			colWidths[i] = spanWidths[i]
		}
	}

	// If expand horizontally set, distribute available space between all columns
	if g.expandh {
		var twidth float32
		for i := 0; i < len(colWidths); i++ {
			twidth += colWidths[i]
		}
		space := pan.ContentWidth() - twidth
		if space > 0 {
			colspace := space / float32(len(colWidths))
			for i := 0; i < len(colWidths); i++ {
				colWidths[i] += colspace
			}
		}
	}

	// If expand vertically set, distribute available space between all rows
	if g.expandv {
		// Calculates the sum of all row heights
		var theight float32
		for _, r := range rows {
			theight += r.height
		}
		// If space available distribute between all rows
		space := pan.ContentHeight() - theight
		if space > 0 {
			rowspace := space / float32(len(rows))
			for i := 0; i < len(rows); i++ {
				rows[i].height += rowspace
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
			colspan := 0
			// Default grid cell alignment
			alignv := g.alignv
			alignh := g.alignh
			// If column has alignment, use them
			if g.columns[ci].alignv != nil {
				alignv = *g.columns[ci].alignv
			}
			if g.columns[ci].alignh != nil {
				alignh = *g.columns[ci].alignh
			}
			// If cell has layout parameters, use them
			if !cell.paramsDef {
				alignh = cell.params.AlignH
				alignv = cell.params.AlignV
				colspan = cell.params.ColSpan
			}
			// Calculates the available width for the cell considering colspan
			var cellWidth float32
			for i := ci; i < ci+colspan+1; i++ {
				if i >= len(colWidths) {
					break
				}
				cellWidth += colWidths[i]
			}
			// Determines child panel horizontal position
			px := cellx
			switch alignh {
			case AlignNone:
			case AlignLeft:
			case AlignRight:
				space := cellWidth - cell.panel.Width()
				if space > 0 {
					px += space
				}
			case AlignCenter:
				space := (cellWidth - cell.panel.Width()) / 2
				if space > 0 {
					px += space
				}
			default:
				panic("Invalid horizontal alignment")
			}
			// Determines child panel vertical position
			py := celly
			switch alignv {
			case AlignNone:
			case AlignTop:
			case AlignBottom:
				space := r.height - cell.panel.Height()
				if space > 0 {
					py += space
				}
			case AlignCenter:
				space := (r.height - cell.panel.Height()) / 2
				if space > 0 {
					py += space
				}
			default:
				panic("Invalid vertical alignment")
			}
			// Sets child panel position
			cell.panel.SetPosition(px, py)
			// Advances to next row cell considering colspan
			cellx += cellWidth
		}
		celly += r.height
	}
}
