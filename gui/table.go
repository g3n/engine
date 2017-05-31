// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"fmt"
	"github.com/g3n/engine/math32"
)

//
// Table implements a panel which can contains child panels
// organized in rows and columns.
//
type Table struct {
	Panel                                // Embedded panel
	styles       *TableStyles            // pointer to current styles
	cols         []TableColumn           // array of columns descriptors
	colmap       map[string]*TableColumn // maps column id to column descriptor
	firstRow     int                     // index of the first row of data to show
	rows         []tableRow              // array of table rows
	headerHeight float32
}

// TableColumn describes a table column
type TableColumn struct {
	Id     string  // column id used to reference the column
	Name   string  // column name shown in the header
	Width  float32 // column preferable width in pixels
	Hidden bool    // hidden flag
	Format string  // format string for numbers and strings
	Align  int     // cell align
	order  int     // show order
	header *Panel  // header panel
	label  *Label  // header label
}

// TableHeaderStyle describes the style of the table header
type TableHeaderStyle struct {
	Border      BorderSizes
	Paddings    BorderSizes
	BorderColor math32.Color4
	BgColor     math32.Color
	FgColor     math32.Color
}

// TableRowStyle describes the style of the table row
type TableRowStyle struct {
	Border      BorderSizes
	Paddings    BorderSizes
	BorderColor math32.Color4
	BgColor     math32.Color
	FgColor     math32.Color
}

// TableRowStyles describes all styles for the table row
type TableRowStyles struct {
	Normal   TableRowStyle
	Selected TableRowStyle
}

// TableStyles describes all styles of the table header and rows
type TableStyles struct {
	Header *TableHeaderStyle
	Row    *TableRowStyles
}

type tableRow struct {
	height float32      // row height
	cells  []*tableCell // array of row cells
}

type tableCell struct {
	Panel             // embedded panel
	label Label       // cell label
	value interface{} // cell current value
}

// NewTable creates and returns a pointer to a new Table with the
// specified width, height and columns
func NewTable(width, height float32, cols []TableColumn) (*Table, error) {

	t := new(Table)
	t.Panel.Initialize(width, height)
	t.styles = &StyleDefault.Table

	// Checks columns descriptors
	t.colmap = make(map[string]*TableColumn)
	t.cols = make([]TableColumn, len(cols))
	copy(t.cols, cols)
	for i := 0; i < len(t.cols); i++ {
		c := &t.cols[i]
		if c.Format == "" {
			c.Format = "%v"
		}
		c.order = i
		if c.Id == "" {
			return nil, fmt.Errorf("Column with empty id")
		}
		if t.colmap[c.Id] != nil {
			return nil, fmt.Errorf("Column with duplicate id")
		}
		t.colmap[c.Id] = c
	}

	// Create header panels
	for i := 0; i < len(t.cols); i++ {
		c := &t.cols[i]
		c.header = NewPanel(0, 0)
		t.applyHeaderStyle(c.header)
		c.label = NewLabel(c.Name)
		c.header.Add(c.label)
		width := c.Width
		if width < c.label.Width()+c.header.MinWidth() {
			width = c.label.Width() + c.header.MinWidth()
		}
		c.header.SetContentSize(width, c.label.Height())
		t.headerHeight = c.header.Height()
		t.Panel.Add(c.header)
	}
	t.recalcHeader()

	return t, nil
}

// SetRows clears all current rows of the table and
// sets new rows from the specifying parameter.
// Each row is a map keyed by the colum id.
// The map value currently can be a string or any number type
// If a row column is not found it is ignored
func (t *Table) SetRows(rows []map[string]interface{}) {

	// Create rows if necessary
	if len(rows) > len(t.rows) {
		count := len(rows) - len(t.rows)
		for ri := 0; ri < count; ri++ {
			var trow tableRow
			trow.cells = make([]*tableCell, 0)
			for coli := 0; coli < len(t.cols); coli++ {
				cell := new(tableCell)
				cell.Initialize(10, 10)
				cell.label.initialize("", StyleDefault.Font)
				cell.Add(&cell.label)
				cell.SetBorders(1, 1, 1, 1)
				trow.cells = append(trow.cells, cell)
				t.Panel.Add(cell)
			}
			t.rows = append(t.rows, trow)
		}
	}

	for ri := 0; ri < len(rows); ri++ {
		t.SetRow(ri, rows[ri])
	}
	t.firstRow = 0
	t.recalc()
}

func (t *Table) SetRow(row int, values map[string]interface{}) {

	if row < 0 || row >= len(t.rows) {
		panic("Invalid row index")
	}

	for ci := 0; ci < len(t.cols); ci++ {
		c := t.cols[ci]
		cv := values[c.Id]
		if cv == nil {
			continue
		}
		t.SetCell(row, c.Id, values[c.Id])
	}
}

func (t *Table) SetCell(row int, colid string, value interface{}) {

	if row < 0 || row >= len(t.rows) {
		panic("Invalid row index")
	}
	c := t.colmap[colid]
	if c == nil {
		return
	}
	cell := t.rows[row].cells[c.order]
	cell.label.SetText(fmt.Sprintf(c.Format, value))
}

// SetColFormat sets the formatting string (Printf) for the specified column
// Update must be called to update the table.
func (t *Table) SetColFormat(id, format string) error {

	c := t.colmap[id]
	if c == nil {
		return fmt.Errorf("No column with id:%s", id)
	}
	c.Format = format
	return nil
}

// recalcHeader recalculates and sets the position and size of the header panels
func (t *Table) recalcHeader() {

	posx := float32(0)
	for i := 0; i < len(t.cols); i++ {
		c := t.cols[i]
		if c.Hidden {
			continue
		}
		c.header.SetPosition(posx, 0)
		posx += c.header.Width()
	}
}

// recalc calculates the positions...
func (t *Table) recalc() {

	// Assumes that the TableColum array is sorted in show order
	py := t.headerHeight
	for ri := t.firstRow; ri < len(t.rows); ri++ {
		row := t.rows[ri]
		px := float32(0)
		// Get maximum height for row
		for ci := 0; ci < len(t.cols); ci++ {
			// If column is hidden, ignore
			c := t.cols[ci]
			if c.Hidden {
				continue
			}
			cell := row.cells[c.order]
			cellHeight := cell.MinHeight() + cell.label.Height()
			if cellHeight > row.height {
				row.height = cellHeight
			}
		}
		// Sets position and size of row cells
		for ci := 0; ci < len(t.cols); ci++ {
			// If column is hidden, ignore
			c := t.cols[ci]
			if c.Hidden {
				continue
			}
			// Sets cell position and size
			cell := row.cells[c.order]
			cell.SetPosition(px, py)
			cell.SetSize(c.header.Width(), row.height)
			log.Error("Cell(%v,%v)(%p) size:%v/%v pos:%v/%v", ri, c.Id, &cell, cell.Width(), cell.Height(), cell.Position().X, cell.Position().Y)
			px += cell.Width()
		}
		py += row.height
	}
}

// colAtOrder returns the TableColumn at the specified view order.
// Returns nil if not found
func (t *Table) colAtOrder(order int) *TableColumn {

	for ci := 0; ci < len(t.cols); ci++ {
		if t.cols[ci].order == order {
			return &t.cols[ci]
		}
	}
	return nil
}

// applyStyle applies the specified menu body style
func (t *Table) applyHeaderStyle(hp *Panel) {

	s := t.styles.Header
	hp.SetBordersFrom(&s.Border)
	hp.SetBordersColor4(&s.BorderColor)
	hp.SetPaddingsFrom(&s.Paddings)
	hp.SetColor(&s.BgColor)
}
