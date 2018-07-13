// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/g3n/engine/gui/assets/icon"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

const (
	// OnTableClick is the event generated when the table is right or left clicked
	// Parameter is TableClickEvent
	OnTableClick = "onTableClick"
	// OnTableRowCount is the event generated when the table row count changes (no parameters)
	OnTableRowCount = "onTableRowCount"
)

// TableSortType is the type used to specify the sort method for a table column
type TableSortType int

// The various sorting types
const (
	TableSortNone TableSortType = iota
	TableSortString
	TableSortNumber
)

// TableSelType is the type used to specify the table row selection
type TableSelType int

const (
	// TableSelSingleRow is the single row selection mode (default)
	TableSelSingleRow TableSelType = iota
	// TableSelMultiRow is the multiple row selection mode
	TableSelMultiRow
)

const (
	tableSortedNoneIcon = icon.SwapVert
	tableSortedAscIcon  = icon.ArrowDownward
	tableSortedDescIcon = icon.ArrowUpward
	tableSortedNone     = 0
	tableSortedAsc      = 1
	tableSortedDesc     = 2
	tableResizerPix     = 4
	tableColMinWidth    = 16
	tableErrInvRow      = "Invalid row index"
	tableErrInvCol      = "Invalid column id"
)

//
// Table implements a panel which can contains child panels
// organized in rows and columns.
//
type Table struct {
	Panel                       // Embedded panel
	styles         *TableStyles // pointer to current styles
	header         tableHeader  // table headers
	rows           []*tableRow  // array of table rows
	rowCursor      int          // index of row cursor
	firstRow       int          // index of the first visible row
	lastRow        int          // index of the last visible row
	vscroll        *ScrollBar   // vertical scroll bar
	statusPanel    Panel        // optional bottom status panel
	statusLabel    *Label       // status label
	scrollBarEvent bool         // do not update the scrollbar value in recalc() if true
	resizerPanel   Panel        // resizer panel
	resizeCol      int          // column being resized
	resizerX       float32      // initial resizer x coordinate
	resizing       bool         // dragging the column resizer
	selType        TableSelType // table selection type
}

// TableColumn describes a table column
type TableColumn struct {
	Id         string          // Column id used to reference the column. Must be unique
	Header     string          // Column name shown in the table header
	Width      float32         // Initial column width in pixels
	Minwidth   float32         // Minimum width in pixels for this column
	Hidden     bool            // Hidden flag
	Align      Align           // Cell content alignment: AlignLeft|AlignCenter|AlignRight
	Format     string          // Format string for formatting the columns' cells
	FormatFunc TableFormatFunc // Format function (overrides Format string)
	Expand     float32         // Column width expansion factor (0 for no expansion)
	Sort       TableSortType   // Column sort type
	Resize     bool            // Allow column to be resized by user
}

// TableCell describes a table cell.
// It is used as a parameter for formatting function
type TableCell struct {
	Tab   *Table      // Pointer to table
	Row   int         // Row index
	Col   string      // Column id
	Value interface{} // Cell value
}

// TableFormatFunc is the type for formatting functions
type TableFormatFunc func(cell TableCell) string

// TableHeaderStyle describes the style of the table header
type TableHeaderStyle BasicStyle

// TableRowStyle describes the style of the table row
type TableRowStyle BasicStyle

// TableStatusStyle describes the style of the table status line panel
type TableStatusStyle BasicStyle

// TableResizerStyle describes the style of the table resizer panel
type TableResizerStyle struct {
	Width       float32
	Border      RectBounds
	BorderColor math32.Color4
	BgColor     math32.Color4
}

// TableStyles describes all styles of the table header and rows
type TableStyles struct {
	Header    TableHeaderStyle
	RowEven   TableRowStyle
	RowOdd    TableRowStyle
	RowCursor TableRowStyle
	RowSel    TableRowStyle
	Status    TableStatusStyle
	Resizer   TableResizerStyle
}

// TableClickEvent describes a mouse click event over a table
// It contains the original mouse event plus additional information
type TableClickEvent struct {
	window.MouseEvent         // Embedded window mouse event
	X                 float32 // Table content area X coordinate
	Y                 float32 // Table content area Y coordinate
	Header            bool    // True if header was clicked
	Row               int     // Index of table row (may be -1)
	Col               string  // Id of table column (may be empty)
	ColOrder          int     // Current column exhibition order
}

// tableHeader is panel which contains the individual header panels for each column
type tableHeader struct {
	Panel                              // embedded panel
	cmap    map[string]*tableColHeader // maps column id with its panel/descriptor
	cols    []*tableColHeader          // array of individual column headers/descriptors
	lastPan Panel                      // last header panel not associated with a user column
}

// tableColHeader is panel for a column header
type tableColHeader struct {
	Panel                      // header panel
	label      *Label          // header label
	ricon      *Label          // header right icon (sort direction)
	id         string          // column id
	width      float32         // initial column width
	minWidth   float32         // minimum width
	format     string          // column format string
	formatFunc TableFormatFunc // column format function
	align      Align           // column alignment
	expand     float32         // column expand factor
	sort       TableSortType   // column sort type
	resize     bool            // column can be resized by user
	order      int             // row columns order
	sorted     int             // current sorted status
	xl         float32         // left border coordinate in pixels
	xr         float32         // right border coordinate in pixels
}

// tableRow is panel which contains an entire table row of cells
type tableRow struct {
	Panel                 // embedded panel
	selected bool         // row selected flag
	cells    []*tableCell // array of row cells
}

// tableCell is a panel which contains one cell (a label)
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
	t.styles = &StyleDefault().Table
	t.rowCursor = -1

	// Initialize table header
	t.header.Initialize(0, 0)
	t.header.cmap = make(map[string]*tableColHeader)
	t.header.cols = make([]*tableColHeader, 0)

	// Create column header panels
	for ci := 0; ci < len(cols); ci++ {
		cdesc := cols[ci]
		// Column id must not be empty
		if cdesc.Id == "" {
			return nil, fmt.Errorf("Column with empty id")
		}
		// Column id must be unique
		if t.header.cmap[cdesc.Id] != nil {
			return nil, fmt.Errorf("Column with duplicate id")
		}
		// Creates a column header
		c := new(tableColHeader)
		c.Initialize(0, 0)
		t.applyHeaderStyle(&c.Panel, false)
		c.label = NewLabel(cdesc.Header)
		c.Add(c.label)
		c.id = cdesc.Id
		c.minWidth = cdesc.Minwidth
		if c.minWidth < tableColMinWidth {
			c.minWidth = tableColMinWidth
		}
		c.width = cdesc.Width
		if c.width < c.minWidth {
			c.width = c.minWidth
		}
		c.align = cdesc.Align
		c.format = cdesc.Format
		c.formatFunc = cdesc.FormatFunc
		c.expand = cdesc.Expand
		c.sort = cdesc.Sort
		c.resize = cdesc.Resize
		// Adds optional sort icon
		if c.sort != TableSortNone {
			c.ricon = NewIcon(string(tableSortedNoneIcon))
			c.Add(c.ricon)
			c.ricon.Subscribe(OnMouseDown, func(evname string, ev interface{}) {
				t.onRicon(evname, c)
			})
		}
		// Sets default format and order
		if c.format == "" {
			c.format = "%v"
		}
		c.order = ci
		c.SetVisible(!cdesc.Hidden)
		t.header.cmap[c.id] = c
		// Sets column header width and height
		width := cdesc.Width
		if width < c.label.Width()+c.MinWidth() {
			width = c.label.Width() + c.MinWidth()
		}
		c.SetContentSize(width, c.label.Height())
		// Adds the column header to the header panel
		t.header.cols = append(t.header.cols, c)
		t.header.Panel.Add(c)
	}
	// Creates last header
	t.header.lastPan.Initialize(0, 0)
	t.applyHeaderStyle(&t.header.lastPan, true)
	t.header.Panel.Add(&t.header.lastPan)

	// Add header panel to the table panel
	t.Panel.Add(&t.header)

	// Creates resizer panel
	t.resizerPanel.Initialize(t.styles.Resizer.Width, 0)
	t.resizerPanel.SetVisible(false)
	t.applyResizerStyle()
	t.Panel.Add(&t.resizerPanel)

	// Creates status panel
	t.statusPanel.Initialize(0, 0)
	t.statusPanel.SetVisible(false)
	t.statusLabel = NewLabel("")
	t.applyStatusStyle()
	t.statusPanel.Add(t.statusLabel)
	t.Panel.Add(&t.statusPanel)

	// Subscribe to events
	t.Panel.Subscribe(OnCursorEnter, t.onCursor)
	t.Panel.Subscribe(OnCursorLeave, t.onCursor)
	t.Panel.Subscribe(OnCursor, t.onCursorPos)
	t.Panel.Subscribe(OnScroll, t.onScroll)
	t.Panel.Subscribe(OnMouseUp, t.onMouse)
	t.Panel.Subscribe(OnMouseDown, t.onMouse)
	t.Panel.Subscribe(OnKeyDown, t.onKey)
	t.Panel.Subscribe(OnKeyRepeat, t.onKey)
	t.Panel.Subscribe(OnResize, t.onResize)
	t.recalc()
	return t, nil
}

// SetStyles set this table styles overriding the default
func (t *Table) SetStyles(ts *TableStyles) {

	t.styles = ts
	t.recalc()
}

// SetSelectionType sets this table selection type
// Possible values are: TableSelSingleRow|TableSelMultiRow
func (t *Table) SetSelectionType(sel TableSelType) {

	t.selType = sel
}

// ShowHeader shows or hides the table header
func (t *Table) ShowHeader(show bool) {

	if t.header.Visible() == show {
		return
	}
	t.header.SetVisible(show)
	t.recalc()
}

// ShowColumn sets the visibility of the column with the specified id
// If the column id does not exit the function panics.
func (t *Table) ShowColumn(col string, show bool) {

	c := t.header.cmap[col]
	if c == nil {
		panic(tableErrInvCol)
	}
	if c.Visible() == show {
		return
	}
	c.SetVisible(show)
	t.recalc()
}

// ShowAllColumns shows all the table columns
func (t *Table) ShowAllColumns() {

	recalc := false
	for ci := 0; ci < len(t.header.cols); ci++ {
		c := t.header.cols[ci]
		if !c.Visible() {
			c.SetVisible(true)
			recalc = true
		}
	}
	if !recalc {
		return
	}
	t.recalc()
}

// RowCount returns the current number of rows in the table
func (t *Table) RowCount() int {

	return len(t.rows)
}

// SetRows clears all current rows of the table and
// sets new rows from the specifying parameter.
// Each row is a map keyed by the colum id.
// The map value currently can be a string or any number type
// If a row column is not found it is ignored
func (t *Table) SetRows(values []map[string]interface{}) {

	// Add missing rows
	if len(values) > len(t.rows) {
		count := len(values) - len(t.rows)
		for row := 0; row < count; row++ {
			t.insertRow(len(t.rows), nil)
		}
		// Remove remaining rows
	} else if len(values) < len(t.rows) {
		for row := len(values); row < len(t.rows); row++ {
			t.removeRow(row)
		}
	}

	// Set rows values
	for row := 0; row < len(values); row++ {
		t.setRow(row, values[row])
	}
	t.firstRow = 0
	t.rowCursor = -1
	t.recalc()
}

// SetRow sets the value of all the cells of the specified row from
// the specified map indexed by column id.
func (t *Table) SetRow(row int, values map[string]interface{}) {

	if row < 0 || row >= len(t.rows) {
		panic(tableErrInvRow)
	}
	t.setRow(row, values)
	t.recalc()
}

// SetCell sets the value of the cell specified by its row and column id
// The function panics if the passed row or column id is invalid
func (t *Table) SetCell(row int, colid string, value interface{}) {

	if row < 0 || row >= len(t.rows) {
		panic(tableErrInvRow)
	}
	if t.header.cmap[colid] == nil {
		panic(tableErrInvCol)
	}
	t.setCell(row, colid, value)
	t.recalc()
}

// SetColFormat sets the formatting string (Printf) for the specified column
// Update must be called to update the table.
func (t *Table) SetColFormat(id, format string) {

	c := t.header.cmap[id]
	if c == nil {
		panic(tableErrInvCol)
	}
	c.format = format
}

// SetColOrder sets the exhibition order of the specified column.
// The previous column which has the specified order will have
// the original column order.
func (t *Table) SetColOrder(colid string, order int) {

	// Checks column id
	c := t.header.cmap[colid]
	if c == nil {
		panic(tableErrInvCol)
	}
	// Checks exhibition order
	if order < 0 || order > len(t.header.cols) {
		panic(tableErrInvRow)
	}
	// Find the exhibition order for the specified column
	for ci := 0; ci < len(t.header.cols); ci++ {
		if t.header.cols[ci] == c {
			// If the order of the specified column is the same, nothing to do
			if ci == order {
				return
			}
			// Swap column orders
			prev := t.header.cols[order]
			t.header.cols[order] = c
			t.header.cols[ci] = prev
			break
		}
	}

	// Recalculates the header and all rows
	t.recalc()
}

// EnableColResize enable or disables if the specified column can be resized by the
// user using the mouse.
func (t *Table) EnableColResize(colid string, enable bool) {

	// Checks column id
	c := t.header.cmap[colid]
	if c == nil {
		panic(tableErrInvCol)
	}
	c.resize = enable
}

// SetColWidth sets the specified column width and may
// change the widths of the columns to the right
func (t *Table) SetColWidth(colid string, width float32) {

	// Checks column id
	c := t.header.cmap[colid]
	if c == nil {
		panic(tableErrInvCol)
	}
	t.setColWidth(c, width)
}

// SetColExpand sets the column expand factor.
// When the table width is increased the columns widths are
// increased proportionally to their expand factor.
// A column with expand factor = 0 is not increased.
func (t *Table) SetColExpand(colid string, expand float32) {

	// Checks column id
	c := t.header.cmap[colid]
	if c == nil {
		panic(tableErrInvCol)
	}
	if expand < 0 {
		c.expand = 0
	} else {
		c.expand = expand
	}
	t.recalc()
}

// AddRow adds a new row at the end of the table with the specified values
func (t *Table) AddRow(values map[string]interface{}) {

	t.InsertRow(len(t.rows), values)
}

// InsertRow inserts the specified values in a new row at the specified index
func (t *Table) InsertRow(row int, values map[string]interface{}) {

	// Checks row index
	if row < 0 || row > len(t.rows) {
		panic(tableErrInvRow)
	}
	t.insertRow(row, values)
	t.recalc()
	t.Dispatch(OnTableRowCount, nil)
}

// RemoveRow removes from the specified row from the table
func (t *Table) RemoveRow(row int) {

	// Checks row index
	if row < 0 || row >= len(t.rows) {
		panic(tableErrInvRow)
	}
	t.removeRow(row)
	maxFirst := t.calcMaxFirst()
	if t.firstRow > maxFirst {
		t.firstRow = maxFirst
	}
	t.recalc()
	t.Dispatch(OnTableRowCount, nil)
}

// Clear removes all rows from the table
func (t *Table) Clear() {

	for ri := 0; ri < len(t.rows); ri++ {
		trow := t.rows[ri]
		t.Panel.Remove(trow)
		trow.DisposeChildren(true)
		trow.Dispose()
	}
	t.rows = nil
	t.firstRow = 0
	t.rowCursor = -1
	t.recalc()
	t.Dispatch(OnTableRowCount, nil)
}

// SelectedRows returns a slice with the indexes of the currently selected rows
// If no row are selected returns an empty slice
func (t *Table) SelectedRows() []int {

	res := make([]int, 0)
	if t.rowCursor >= 0 {
		res = append(res, t.rowCursor)
	}
	for ri := 0; ri < len(t.rows); ri++ {
		if t.rows[ri].selected && ri != t.rowCursor {
			res = append(res, ri)
		}
	}
	return res
}

// ShowStatus sets the visibility of the status lines at the bottom of the table
func (t *Table) ShowStatus(show bool) {

	if t.statusPanel.Visible() == show {
		return
	}
	t.statusPanel.SetVisible(show)
	t.recalcStatus()
	t.recalc()
}

// SetStatusText sets the text of status line at the bottom of the table
// It does not change its current visibility
func (t *Table) SetStatusText(text string) {

	t.statusLabel.SetText(text)
}

// Rows returns a slice of maps with the contents of the table rows
// specified by the rows first and last index.
// To get all the table rows, use Rows(0, -1)
func (t *Table) Rows(fi, li int) []map[string]interface{} {

	if fi < 0 || fi >= len(t.header.cols) {
		panic(tableErrInvRow)
	}
	if li < 0 {
		li = len(t.rows) - 1
	} else if li < 0 || li >= len(t.rows) {
		panic(tableErrInvRow)
	}
	if li < fi {
		panic("Last index less than first index")
	}
	res := make([]map[string]interface{}, li-li+1)
	for ri := fi; ri <= li; ri++ {
		trow := t.rows[ri]
		rmap := make(map[string]interface{})
		for ci := 0; ci < len(t.header.cols); ci++ {
			c := t.header.cols[ci]
			rmap[c.id] = trow.cells[c.order].value
		}
		res = append(res, rmap)
	}
	return res
}

// Row returns a map with the current contents of the specified row index
func (t *Table) Row(ri int) map[string]interface{} {

	if ri < 0 || ri > len(t.header.cols) {
		panic(tableErrInvRow)
	}
	res := make(map[string]interface{})
	trow := t.rows[ri]
	for ci := 0; ci < len(t.header.cols); ci++ {
		c := t.header.cols[ci]
		res[c.id] = trow.cells[c.order].value
	}
	return res
}

// Cell returns the current content of the specified cell
func (t *Table) Cell(col string, ri int) interface{} {

	c := t.header.cmap[col]
	if c == nil {
		panic(tableErrInvCol)
	}
	if ri < 0 || ri >= len(t.rows) {
		panic(tableErrInvRow)
	}
	trow := t.rows[ri]
	return trow.cells[c.order].value
}

// SortColumn sorts the specified column interpreting its values as strings or numbers
// and sorting in ascending or descending order.
// This sorting is independent of the sort configuration of column set when the table was created
func (t *Table) SortColumn(col string, asString bool, asc bool) {

	c := t.header.cmap[col]
	if c == nil {
		panic(tableErrInvCol)
	}
	if len(t.rows) < 2 {
		return
	}
	if asString {
		ts := tableSortString{rows: t.rows, col: c.order, asc: asc, format: c.format}
		sort.Sort(ts)
	} else {
		ts := tableSortNumber{rows: t.rows, col: c.order, asc: asc}
		sort.Sort(ts)
	}
	t.recalc()
}

// setRow sets the value of all the cells of the specified row from
// the specified map indexed by column id.
func (t *Table) setRow(row int, values map[string]interface{}) {

	for ci := 0; ci < len(t.header.cols); ci++ {
		c := t.header.cols[ci]
		cv, ok := values[c.id]
		if !ok {
			continue
		}
		t.setCell(row, c.id, cv)
	}
}

// setCell sets the value of the cell specified by its row and column id
func (t *Table) setCell(row int, colid string, value interface{}) {

	c := t.header.cmap[colid]
	if c == nil {
		return
	}
	cell := t.rows[row].cells[c.order]
	cell.label.SetText(fmt.Sprintf(c.format, value))
	cell.value = value
}

// insertRow is the internal version of InsertRow which does not call recalc()
func (t *Table) insertRow(row int, values map[string]interface{}) {

	// Creates tableRow panel
	trow := new(tableRow)
	trow.Initialize(0, 0)
	trow.cells = make([]*tableCell, 0)
	for ci := 0; ci < len(t.header.cols); ci++ {
		// Creates tableRow cell panel
		cell := new(tableCell)
		cell.Initialize(0, 0)
		cell.label.initialize("", StyleDefault().Font)
		cell.Add(&cell.label)
		trow.cells = append(trow.cells, cell)
		trow.Panel.Add(cell)
	}
	t.Panel.Add(trow)

	// Inserts tableRow in the table rows at the specified index
	t.rows = append(t.rows, nil)
	copy(t.rows[row+1:], t.rows[row:])
	t.rows[row] = trow
	t.updateRowStyle(row)

	// Sets the new row values from the specified map
	if values != nil {
		t.SetRow(row, values)
	}
	t.recalcRow(row)
}

// ScrollDown scrolls the table the specified number of rows down if possible
func (t *Table) scrollDown(n int) {

	// Calculates number of rows to scroll down
	maxFirst := t.calcMaxFirst()
	maxScroll := maxFirst - t.firstRow
	if maxScroll <= 0 {
		return
	}
	if n > maxScroll {
		n = maxScroll
	}

	t.firstRow += n
	if t.rowCursor < t.firstRow {
		t.rowCursor = t.firstRow
		t.Dispatch(OnChange, nil)
	}
	t.recalc()
	return
}

// ScrollUp scrolls the table the specified number of rows up if possible
func (t *Table) scrollUp(n int) {

	// Calculates number of rows to scroll up
	if t.firstRow == 0 {
		return
	}
	if n > t.firstRow {
		n = t.firstRow
	}
	t.firstRow -= n
	lastRow := t.lastRow - n
	if t.rowCursor > lastRow {
		t.rowCursor = lastRow
		t.Dispatch(OnChange, nil)
	}
	t.recalc()
}

// removeRow removes from the table the row specified its index
func (t *Table) removeRow(row int) {

	// Get row to be removed
	trow := t.rows[row]

	// Remove row from table children
	t.Panel.Remove(trow)

	// Remove row from rows array
	copy(t.rows[row:], t.rows[row+1:])
	t.rows[len(t.rows)-1] = nil
	t.rows = t.rows[:len(t.rows)-1]

	// Dispose row resources
	trow.DisposeChildren(true)
	trow.Dispose()
}

// onCursor process subscribed cursor events
func (t *Table) onCursor(evname string, ev interface{}) {

	switch evname {
	case OnCursorEnter:
		t.root.SetScrollFocus(t)
	case OnCursorLeave:
		t.root.SetScrollFocus(nil)
	}
	t.root.StopPropagation(Stop3D)
}

// onCursorPos process subscribed cursor position events
func (t *Table) onCursorPos(evname string, ev interface{}) {

	// Convert mouse window coordinates to table content coordinates
	cev := ev.(*window.CursorEvent)
	cx, _ := t.ContentCoords(cev.Xpos, cev.Ypos)

	// If user is dragging the resizer, updates its position
	if t.resizing {
		t.resizerPanel.SetPosition(cx, 0)
		return
	}

	// Checks if the mouse cursor is near the border of a resizable column
	found := false
	for ci := 0; ci < len(t.header.cols); ci++ {
		c := t.header.cols[ci]
		dx := math32.Abs(cx - c.xr)
		if dx < tableResizerPix {
			if c.resize {
				found = true
				t.resizeCol = ci
				t.resizerX = c.xr
				t.root.SetCursorHResize()
			}
			break
		}
	}
	// If column not found but previously was near a resizable column,
	// resets the the window cursor.
	if !found && t.resizeCol >= 0 {
		t.root.SetCursorNormal()
		t.resizeCol = -1
	}
	t.root.StopPropagation(Stop3D)
}

// onMouseEvent process subscribed mouse events
func (t *Table) onMouse(evname string, ev interface{}) {

	e := ev.(*window.MouseEvent)
	t.root.SetKeyFocus(t)
	switch evname {
	case OnMouseDown:
		// If over a resizable column border, shows the resizer panel
		if t.resizeCol >= 0 && e.Button == window.MouseButtonLeft {
			t.resizing = true
			height := t.ContentHeight()
			if t.statusPanel.Visible() {
				height -= t.statusPanel.Height()
			}
			px := t.resizerX - t.resizerPanel.Width()/2
			t.resizerPanel.SetPositionX(px)
			t.resizerPanel.SetHeight(height)
			t.resizerPanel.SetVisible(true)
			t.SetTopChild(&t.resizerPanel)
			return
		}
		// Find click position
		var tce TableClickEvent
		tce.MouseEvent = *e
		t.findClick(&tce)
		// If row is clicked, selects it
		if tce.Row >= 0 && e.Button == window.MouseButtonLeft {
			t.rowCursor = tce.Row
			if t.selType == TableSelMultiRow && e.Mods == window.ModControl {
				t.toggleRowSel(t.rowCursor)
			}
			t.recalc()
			t.Dispatch(OnChange, nil)
		}
		// Creates and dispatch TableClickEvent for user's context menu
		t.Dispatch(OnTableClick, tce)
	case OnMouseUp:
		// If user was resizing a column, hides the resizer and
		// sets the new column width if possible
		if t.resizing {
			t.resizing = false
			t.resizerPanel.SetVisible(false)
			t.root.SetCursorNormal()
			// Calculates the new column width
			cx, _ := t.ContentCoords(e.Xpos, e.Ypos)
			c := t.header.cols[t.resizeCol]
			width := cx - c.xl
			t.setColWidth(c, width)
		}
	default:
		return
	}
	t.root.StopPropagation(StopAll)
}

// onKeyEvent receives subscribed key events for this table
func (t *Table) onKey(evname string, ev interface{}) {

	kev := ev.(*window.KeyEvent)
	if kev.Keycode == window.KeyUp && kev.Mods == 0 {
		t.selPrev()
	} else if kev.Keycode == window.KeyDown && kev.Mods == 0 {
		t.selNext()
	} else if kev.Keycode == window.KeyPageUp && kev.Mods == 0 {
		t.prevPage()
	} else if kev.Keycode == window.KeyPageDown && kev.Mods == 0 {
		t.nextPage()
	} else if kev.Keycode == window.KeyPageUp && kev.Mods == window.ModControl {
		t.firstPage()
	} else if kev.Keycode == window.KeyPageDown && kev.Mods == window.ModControl {
		t.lastPage()
	} else if kev.Keycode == window.KeyEnter && kev.Mods == window.ModControl {
		if t.selType == TableSelMultiRow {
			t.toggleRowSel(t.rowCursor)
		}
	}
}

// onResize receives subscribed resize events for this table
func (t *Table) onResize(evname string, ev interface{}) {

	t.recalc()
	t.recalcStatus()
}

// onScroll receives subscribed scroll events for this table
func (t *Table) onScroll(evname string, ev interface{}) {

	sev := ev.(*window.ScrollEvent)
	if sev.Yoffset > 0 {
		t.scrollUp(1)
	} else if sev.Yoffset < 0 {
		t.scrollDown(1)
	}
	t.root.StopPropagation(Stop3D)
}

// onRicon receives subscribed events for column header right icon
func (t *Table) onRicon(evname string, c *tableColHeader) {

	icon := tableSortedNoneIcon
	var asc bool
	if c.sorted == tableSortedNone || c.sorted == tableSortedDesc {
		c.sorted = tableSortedAsc
		icon = tableSortedAscIcon
		asc = false
	} else {
		c.sorted = tableSortedDesc
		icon = tableSortedDescIcon
		asc = true
	}

	var asString bool
	if c.sort == TableSortString {
		asString = true
	} else {
		asString = false
	}
	t.SortColumn(c.id, asString, asc)
	c.ricon.SetText(string(icon))
}

// findClick finds where in the table the specified mouse click event
// occurred updating the specified TableClickEvent with the click coordinates.
func (t *Table) findClick(ev *TableClickEvent) {

	x, y := t.ContentCoords(ev.Xpos, ev.Ypos)
	ev.X = x
	ev.Y = y
	ev.Row = -1
	// Find column id
	colx := float32(0)
	for ci := 0; ci < len(t.header.cols); ci++ {
		c := t.header.cols[ci]
		if !c.Visible() {
			continue
		}
		colx += t.header.cols[ci].Width()
		if x < colx {
			ev.Col = c.id
			ev.ColOrder = ci
			break
		}
	}
	// If column not found the user clicked at the right of rows
	if ev.Col == "" {
		return
	}
	// Checks if is in header
	if t.header.Visible() && y < t.header.Height() {
		ev.Header = true
	}

	// Find row clicked
	rowy := float32(0)
	if t.header.Visible() {
		rowy = t.header.Height()
	}
	theight := t.ContentHeight()
	for ri := t.firstRow; ri < len(t.rows); ri++ {
		trow := t.rows[ri]
		rowy += trow.height
		if rowy > theight {
			break
		}
		if y < rowy {
			ev.Row = ri
			break
		}
	}
}

// selNext selects the next row if possible
func (t *Table) selNext() {

	// If selected row is last, nothing to do
	if t.rowCursor == len(t.rows)-1 {
		return
	}
	// If no selected row, selects first visible row
	if t.rowCursor < 0 {
		t.rowCursor = t.firstRow
		t.recalc()
		t.Dispatch(OnChange, nil)
		return
	}
	// Selects next row
	t.rowCursor++
	t.Dispatch(OnChange, nil)

	// Scroll down if necessary
	if t.rowCursor > t.lastRow {
		t.scrollDown(1)
	} else {
		t.recalc()
	}
}

// selPrev selects the previous row if possible
func (t *Table) selPrev() {

	// If selected row is first, nothing to do
	sel := t.rowCursor
	if sel == 0 {
		return
	}
	// If no selected row, selects last visible row
	if sel < 0 {
		t.rowCursor = t.lastRow
		t.recalc()
		t.Dispatch(OnChange, nil)
		return
	}
	// Selects previous row and selects previous
	prev := sel - 1
	t.rowCursor = prev

	// Scroll up if necessary
	if prev < t.firstRow && t.firstRow > 0 {
		t.scrollUp(1)
	} else {
		t.recalc()
	}
	t.Dispatch(OnChange, nil)
}

// nextPage shows the next page of rows and selects its first row
func (t *Table) nextPage() {

	if len(t.rows) == 0 {
		return
	}
	if t.lastRow == len(t.rows)-1 {
		t.rowCursor = t.lastRow
		t.recalc()
		t.Dispatch(OnChange, nil)
		return
	}
	plen := t.lastRow - t.firstRow
	if plen <= 0 {
		return
	}
	t.scrollDown(plen)
}

// prevPage shows the previous page of rows and selects its last row
func (t *Table) prevPage() {

	if t.firstRow == 0 {
		t.rowCursor = 0
		t.recalc()
		t.Dispatch(OnChange, nil)
		return
	}
	plen := t.lastRow - t.firstRow
	if plen <= 0 {
		return
	}
	t.scrollUp(plen)
}

// firstPage shows the first page of rows and selects the first row
func (t *Table) firstPage() {

	if len(t.rows) == 0 {
		return
	}
	t.firstRow = 0
	t.rowCursor = 0
	t.recalc()
	t.Dispatch(OnChange, nil)
}

// lastPage shows the last page of rows and selects the last row
func (t *Table) lastPage() {

	if len(t.rows) == 0 {
		return
	}
	maxFirst := t.calcMaxFirst()
	t.firstRow = maxFirst
	t.rowCursor = len(t.rows) - 1
	t.recalc()
	t.Dispatch(OnChange, nil)
}

// selectRow selects the specified row.
// Should be used only when multi row selection is enabled
func (t *Table) selectRow(ri int) {

	trow := t.rows[ri]
	trow.selected = true
	t.Dispatch(OnChange, nil)
}

// toggleRowSel toogles the specified row selection state
// Should be used only when multi row selection is enabled
func (t *Table) toggleRowSel(ri int) {

	trow := t.rows[ri]
	trow.selected = !trow.selected
	t.Dispatch(OnChange, nil)
}

// setColWidth sets the width of the specified column
func (t *Table) setColWidth(c *tableColHeader, width float32) {

	// Sets the column width
	if width < c.minWidth {
		width = c.minWidth
	}
	if c.Width() == width {
		return
	}
	dw := width - c.Width()
	c.SetWidth(width)

	// Find the column index and if any column has expand != 0
	hasExpand := false
	ci := -1
	for i := 0; i < len(t.header.cols); i++ {
		current := t.header.cols[i]
		if current == c {
			ci = i
		}
		if current.expand > 0 && current.Visible() {
			hasExpand = true
		}
	}
	if ci >= len(t.header.cols) {
		panic("Internal: column not found")
	}
	// If no column is expandable, nothing more todo
	if !hasExpand {
		t.recalc()
		return
	}
	// Calculates the width of the columns at the right
	rwidth := float32(0)
	for i := ci + 1; i < len(t.header.cols); i++ {
		c := t.header.cols[i]
		if !c.Visible() {
			continue
		}
		rwidth += c.Width()
	}
	// Distributes the delta to the columns at the right
	for i := ci + 1; i < len(t.header.cols); i++ {
		c := t.header.cols[i]
		if !c.Visible() {
			continue
		}
		cdelta := -dw * (c.Width() / rwidth)
		newWidth := c.Width() + cdelta
		if newWidth < c.minWidth {
			newWidth = c.minWidth
		}
		c.SetWidth(newWidth)
	}
	t.recalc()
}

// recalcHeader recalculates and sets the position and size of the header panels
func (t *Table) recalcHeader() {

	// Calculates total width, height, expansion and available width space
	hwidth := float32(0)
	height := float32(0)
	wspace := float32(0)
	totalExpand := float32(0)
	for ci := 0; ci < len(t.header.cols); ci++ {
		// If column is invisible, ignore
		c := t.header.cols[ci]
		if !c.Visible() {
			continue
		}
		if c.Height() > height {
			height = c.Height()
		}
		if c.expand > 0 {
			totalExpand += c.expand
		}
		hwidth += c.Width()
	}
	// Total table width
	twidth := t.ContentWidth()
	if t.vscroll != nil && t.vscroll.Visible() {
		twidth -= t.vscroll.Width()
	}
	// Available space for columns: may be negative
	wspace = twidth - hwidth

	// If no expandable column, keeps the columns widths
	if totalExpand == 0 {
	} else if wspace >= 0 {
		for ci := 0; ci < len(t.header.cols); ci++ {
			// If column is invisible, ignore
			c := t.header.cols[ci]
			if !c.Visible() {
				continue
			}
			// There is space available and if column is expandable,
			// expands it proportionally to the other expandable columns
			factor := c.expand / totalExpand
			w := factor * wspace
			c.SetWidth(c.Width() + w)
		}
	} else {
		acols := make([]*tableColHeader, 0)
		awidth := float32(0)
		widthAvail := twidth
		// Sets the widths of the columns
		for ci := 0; ci < len(t.header.cols); ci++ {
			// If column is invisible, ignore
			c := t.header.cols[ci]
			if !c.Visible() {
				continue
			}
			// The table was reduced so shrinks this column proportionally to its current width
			factor := c.Width() / hwidth
			newWidth := factor * twidth
			if newWidth < c.minWidth {
				newWidth = c.minWidth
				c.SetWidth(newWidth)
				widthAvail -= c.minWidth
			} else {
				acols = append(acols, c)
				awidth += c.Width()
			}
		}
		for ci := 0; ci < len(acols); ci++ {
			c := acols[ci]
			factor := c.Width() / awidth
			newWidth := factor * widthAvail
			c.SetWidth(newWidth)
		}
	}

	// Sets the header panel and its internal panels positions
	posx := float32(0)
	for ci := 0; ci < len(t.header.cols); ci++ {
		// If column is invisible, ignore
		c := t.header.cols[ci]
		if !c.Visible() {
			continue
		}
		// Sets the right icon position inside the column header panel
		if c.ricon != nil {
			ix := c.ContentWidth() - c.ricon.Width()
			if ix < 0 {
				ix = 0
			}
			c.ricon.SetPosition(ix, 0)
		}
		// Sets the column header panel position
		c.SetPosition(posx, 0)
		c.SetVisible(true)
		c.xl = posx
		posx += c.Width()
		c.xr = posx
	}

	// Last header
	w := t.ContentWidth() - posx
	if w > 0 {
		t.header.lastPan.SetVisible(true)
		t.header.lastPan.SetSize(w, height)
		t.header.lastPan.SetPosition(posx, 0)
	} else {
		t.header.lastPan.SetVisible(false)
	}

	// Header container
	t.header.SetWidth(t.ContentWidth())
	t.header.SetContentHeight(height)
}

// recalcStatus recalculates and sets the position and size of the status panel and its label
func (t *Table) recalcStatus() {

	if !t.statusPanel.Visible() {
		return
	}
	t.statusPanel.SetContentHeight(t.statusLabel.Height())
	py := t.ContentHeight() - t.statusPanel.Height()
	t.statusPanel.SetPosition(0, py)
	t.statusPanel.SetWidth(t.ContentWidth())
}

// recalc calculates the visibility, positions and sizes of all row cells.
// should be called in the following situations:
// - the table is resized
// - row is added, inserted or removed
// - column alignment and expansion changed
// - column visibility is changed
// - horizontal or vertical scroll position changed
func (t *Table) recalc() {

	// Get available row height for rows
	starty, theight := t.rowsHeight()

	// Determines if it is necessary to show the scrollbar or not.
	scroll := false
	py := starty
	for ri := 0; ri < len(t.rows); ri++ {
		trow := t.rows[ri]
		py += trow.height
		if py > starty+theight {
			scroll = true
			break
		}
	}
	t.setVScrollBar(scroll)
	// Recalculates the header
	t.recalcHeader()

	// Sets the position and sizes of all cells of the visible rows
	py = starty
	for ri := 0; ri < len(t.rows); ri++ {
		trow := t.rows[ri]
		// If row is before first row or its y coordinate is greater the table height,
		// sets it invisible
		if ri < t.firstRow || py > starty+theight {
			trow.SetVisible(false)
			continue
		}
		t.recalcRow(ri)
		// Set row y position and visible
		trow.SetPosition(0, py)
		trow.SetVisible(true)
		t.updateRowStyle(ri)
		// Set the last completely visible row index
		if py+trow.Height() <= starty+theight {
			t.lastRow = ri
		}
		//log.Error("ri:%v py:%v theight:%v", ri, py, theight)
		py += trow.height
	}
	// Status panel must be on top of all the row panels
	t.SetTopChild(&t.statusPanel)
}

// recalcRow recalculates the positions and sizes of all cells of the specified row
// Should be called when the row is created and column visibility or order is changed.
func (t *Table) recalcRow(ri int) {

	trow := t.rows[ri]
	// Calculates and sets row height
	maxheight := float32(0)
	for ci := 0; ci < len(t.header.cols); ci++ {
		// If column is hidden, ignore
		c := t.header.cols[ci]
		if !c.Visible() {
			continue
		}
		cell := trow.cells[c.order]
		cellHeight := cell.MinHeight() + cell.label.Height()
		if cellHeight > maxheight {
			maxheight = cellHeight
		}
	}
	trow.SetContentHeight(maxheight)

	// Sets row cells sizes and positions and sets row width
	px := float32(0)
	for ci := 0; ci < len(t.header.cols); ci++ {
		// If column is hidden, ignore
		c := t.header.cols[ci]
		cell := trow.cells[c.order]
		if !c.Visible() {
			cell.SetVisible(false)
			continue
		}
		// Sets cell position and size
		cell.SetPosition(px, 0)
		cell.SetVisible(true)
		cell.SetSize(c.Width(), trow.ContentHeight())
		// Checks for format function
		if c.formatFunc != nil {
			text := c.formatFunc(TableCell{t, ri, c.id, cell.value})
			cell.label.SetText(text)
		}
		// Sets the cell label alignment inside the cell
		ccw := cell.ContentWidth()
		lw := cell.label.Width()
		space := ccw - lw
		lx := float32(0)
		switch c.align {
		case AlignLeft:
		case AlignRight:
			if space > 0 {
				lx = ccw - lw
			}
		case AlignCenter:
			if space > 0 {
				lx = space / 2
			}
		}
		cell.label.SetPosition(lx, 0)
		px += c.Width()
	}
	trow.SetContentWidth(px)
}

// rowsHeight returns the available start y coordinate and height in the table for rows,
// considering the visibility of the header and status panels.
func (t *Table) rowsHeight() (float32, float32) {

	start := float32(0)
	height := t.ContentHeight()
	if t.header.Visible() {
		height -= t.header.Height()
		start += t.header.Height()
	}
	if t.statusPanel.Visible() {
		height -= t.statusPanel.Height()
	}
	if height < 0 {
		return 0, 0
	}
	return start, height
}

// setVScrollBar sets the visibility state of the vertical scrollbar
func (t *Table) setVScrollBar(state bool) {

	// Visible
	if state {
		var scrollWidth float32 = 20
		// Creates scroll bar if necessary
		if t.vscroll == nil {
			t.vscroll = NewVScrollBar(0, 0)
			t.vscroll.SetBorders(0, 0, 0, 1)
			t.vscroll.Subscribe(OnChange, t.onVScrollBar)
			t.Panel.Add(t.vscroll)
		}
		// Sets the scroll bar size and positions
		py, height := t.rowsHeight()
		t.vscroll.SetSize(scrollWidth, height)
		t.vscroll.SetPositionX(t.ContentWidth() - scrollWidth)
		t.vscroll.SetPositionY(py)
		t.vscroll.recalc()
		t.vscroll.SetVisible(true)
		if !t.scrollBarEvent {
			maxFirst := t.calcMaxFirst()
			t.vscroll.SetValue(float32(t.firstRow) / float32(maxFirst))
		} else {
			t.scrollBarEvent = false
		}
		// scroll bar must be on top of all table rows
		t.SetTopChild(t.vscroll)
		// Not visible
	} else {
		if t.vscroll != nil {
			t.vscroll.SetVisible(false)
		}
	}
}

// onVScrollBar is called when a vertical scroll bar event is received
func (t *Table) onVScrollBar(evname string, ev interface{}) {

	// Calculates the new first visible line
	pos := t.vscroll.Value()
	maxFirst := t.calcMaxFirst()
	first := int(math.Floor((float64(maxFirst) * pos) + 0.5))

	// Sets the new selected row
	sel := t.rowCursor
	selChange := false
	if sel < first {
		t.rowCursor = first
		selChange = true
	} else {
		lines := first - t.firstRow
		lastRow := t.lastRow + lines
		if sel > lastRow {
			t.rowCursor = lastRow
			selChange = true
		}
	}
	t.scrollBarEvent = true
	t.firstRow = first
	t.recalc()
	if selChange {
		t.Dispatch(OnChange, nil)
	}
}

// calcMaxFirst calculates the maximum index of the first visible row
// such as the remaining rows fits completely inside the table
// It is used when scrolling the table vertically
func (t *Table) calcMaxFirst() int {

	_, total := t.rowsHeight()
	ri := len(t.rows) - 1
	if ri < 0 {
		return 0
	}
	height := float32(0)
	for {
		trow := t.rows[ri]
		height += trow.height
		if height > total {
			break
		}
		ri--
		if ri < 0 {
			break
		}
	}
	return ri + 1
}

// updateRowStyle applies the correct style for the specified row
func (t *Table) updateRowStyle(ri int) {

	row := t.rows[ri]
	var trs TableRowStyle
	if ri == t.rowCursor {
		trs = t.styles.RowCursor
	} else if row.selected {
		trs = t.styles.RowSel
	} else {
		if ri%2 == 0 {
			trs = t.styles.RowEven
		} else {
			trs = t.styles.RowOdd
		}
	}
	t.applyRowStyle(row, &trs)
}

// applyHeaderStyle applies style to the specified table header
// the last header panel does not the right border.
func (t *Table) applyHeaderStyle(h *Panel, last bool) {

	styleCopy := t.styles.Header.PanelStyle
	if last {
		styleCopy.Border.Right = 0
	}
	h.ApplyStyle(&styleCopy)
}

// applyRowStyle applies the specified style to all cells for the specified table row
func (t *Table) applyRowStyle(trow *tableRow, trs *TableRowStyle) {

	for i := 0; i < len(trow.cells); i++ {
		cell := trow.cells[i]
		cell.ApplyStyle(&trs.PanelStyle)
	}
}

// applyStatusStyle applies the status style
func (t *Table) applyStatusStyle() {

	s := t.styles.Status
	t.statusPanel.ApplyStyle(&s.PanelStyle)
}

// applyResizerStyle applies the status style
func (t *Table) applyResizerStyle() {

	s := t.styles.Resizer
	t.resizerPanel.SetBordersFrom(&s.Border)
	t.resizerPanel.SetBordersColor4(&s.BorderColor)
	t.resizerPanel.SetColor4(&s.BgColor)
}

// tableSortString is an internal type implementing the sort.Interface
// and is used to sort a table column interpreting its values as strings
type tableSortString struct {
	rows   []*tableRow
	col    int
	asc    bool
	format string
}

func (ts tableSortString) Len() int      { return len(ts.rows) }
func (ts tableSortString) Swap(i, j int) { ts.rows[i], ts.rows[j] = ts.rows[j], ts.rows[i] }
func (ts tableSortString) Less(i, j int) bool {

	vi := ts.rows[i].cells[ts.col].value
	vj := ts.rows[j].cells[ts.col].value
	si := fmt.Sprintf(ts.format, vi)
	sj := fmt.Sprintf(ts.format, vj)
	if ts.asc {
		return si < sj
	}
	return sj < si
}

// tableSortNumber is an internal type implementing the sort.Interface
// and is used to sort a table column interpreting its values as numbers
type tableSortNumber struct {
	rows []*tableRow
	col  int
	asc  bool
}

func (ts tableSortNumber) Len() int      { return len(ts.rows) }
func (ts tableSortNumber) Swap(i, j int) { ts.rows[i], ts.rows[j] = ts.rows[j], ts.rows[i] }
func (ts tableSortNumber) Less(i, j int) bool {

	vi := ts.rows[i].cells[ts.col].value
	vj := ts.rows[j].cells[ts.col].value
	ni := cv2f64(vi)
	nj := cv2f64(vj)
	if ts.asc {
		return ni < nj
	}
	return nj < ni
}

// Try to convert an interface value to a float64 number
func cv2f64(v interface{}) float64 {

	if v == nil {
		return 0
	}
	switch n := v.(type) {
	case uint8:
		return float64(n)
	case uint16:
		return float64(n)
	case uint32:
		return float64(n)
	case uint64:
		return float64(n)
	case uint:
		return float64(n)
	case int8:
		return float64(n)
	case int16:
		return float64(n)
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	case int:
		return float64(n)
	case string:
		sv, err := strconv.ParseFloat(n, 64)
		if err == nil {
			return sv
		}
		return 0
	default:
		return 0
	}
}
