// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

// HBoxLayout implements a panel layout which arranges the panel children horizontally.
// The children can be separated by a space in pixels set by SetSpacing().
// The whole group of children can be aligned horizontally by SetAlignH() which can
// accept the following types of alignment:
//
// 	AlignLeft: Try to align the group of children to the left if the panel width is
// 	greater the the sum of the children widths + spacing.
//
// 	AlignRight: Try to align the group of children to the right if the panel width is
// 	greater the the sum of the children widths + spacing.
//
// 	AlignCenter: Try to align the group of children in the center if the panel width is
// 	greater the the sum of the children widths + spacing.
//
// 	AlignWidth - Try to align the individual children with the same same space between each other.
// 	Each individual child can be aligned vertically by SetLayoutParameters()
//
// If the layout method SetAutoHeight(true) is called, the panel minimum content height will be the
// height of the child with the largest height.
//
// If the layout method SetAutoWidth(true) is called, the panel minimum content width will be the
// sum of its children's widths plus the spacing.
type HBoxLayout struct {
	pan        IPanel
	spacing    float32
	alignH     Align
	autoHeight bool
	minHeight  bool
}

// HBoxLayoutParams specify the vertical alignment of each individual child.
type HBoxLayoutParams struct {
	Expand float32 // item expand horizontally factor (0 - no expand)
	AlignV Align   // item vertical alignment
}

// NewHBoxLayout creates and returns a pointer to a new horizontal box layout
func NewHBoxLayout() *HBoxLayout {

	bl := new(HBoxLayout)
	bl.spacing = 0
	bl.alignH = AlignLeft
	return bl
}

// SetSpacing sets the horizontal spacing between the items in pixels
// and updates the layout if possible
func (bl *HBoxLayout) SetSpacing(spacing float32) {

	bl.spacing = spacing
	bl.Recalc(bl.pan)
}

// SetAlignH sets the horizontal alignment of the whole group of items
// inside the parent panel and updates the layout.
// This only has any effect if there are no expanded items.
func (bl *HBoxLayout) SetAlignH(align Align) {

	bl.alignH = align
	bl.Recalc(bl.pan)
}

// SetAutoHeight sets if the panel minimum height should be the height of
// the largest of its children's height.
func (bl *HBoxLayout) SetAutoHeight(state bool) {

	bl.autoHeight = state
	bl.Recalc(bl.pan)
}

// SetAutoWidth sets if the panel minimum width should be sum of its
// children's width plus the spacing
func (bl *HBoxLayout) SetAutoWidth(state bool) {

	bl.minHeight = state
	bl.Recalc(bl.pan)
}

// Recalc recalculates and sets the position and sizes of all children
func (bl *HBoxLayout) Recalc(ipan IPanel) {

	// Saves the received panel
	bl.pan = ipan
	if bl.pan == nil {
		return
	}
	parent := ipan.GetPanel()
	if len(parent.Children()) == 0 {
		return
	}

	// If autoHeight is set, get the maximum height of all the panel's children
	// and if the panel content height is less than this maximum, set its content height to this value.
	if bl.autoHeight {
		var maxHeight float32
		for _, ichild := range parent.Children() {
			child := ichild.(IPanel).GetPanel()
			if !child.Visible() {
				continue
			}
			if child.Height() > maxHeight {
				maxHeight = child.Height()
			}
		}
		if parent.ContentHeight() < maxHeight {
			parent.setContentSize(parent.ContentWidth(), maxHeight, false)
		}
	}

	// If minHeight is set, get the sum of widths of this panel's children plus the spacings.
	// If the panel content width is less than this width, set its content width to this value.
	if bl.minHeight {
		var totalWidth float32
		for _, ichild := range parent.Children() {
			child := ichild.(IPanel).GetPanel()
			if !child.Visible() {
				continue
			}
			totalWidth += child.Width()
		}
		// Adds spacing
		totalWidth += bl.spacing * float32(len(parent.Children())-1)
		if parent.ContentWidth() < totalWidth {
			parent.setContentSize(totalWidth, parent.ContentHeight(), false)
		}
	}

	// Calculates the total width, expanded width, fixed width and
	// the sum of the expand factor for all items.
	var twidth float32
	//var ewidth float32
	var fwidth float32
	var texpand float32
	ecount := 0
	paramsDef := HBoxLayoutParams{Expand: 0, AlignV: AlignTop}
	for pos, obj := range parent.Children() {
		pan := obj.(IPanel).GetPanel()
		if !pan.Visible() {
			continue
		}
		// Get item layout parameters or use default
		params := paramsDef
		if pan.layoutParams != nil {
			params = *pan.layoutParams.(*HBoxLayoutParams)
		}
		// Calculate total width
		twidth += pan.Width()
		if pos > 0 {
			twidth += bl.spacing
		}
		// Calculate width of expanded items
		if params.Expand > 0 {
			texpand += params.Expand
			//ewidth += pan.Width()
			//if pos > 0 {
			//	ewidth += bl.spacing
			//}
			ecount++
			// Calculate width of fixed items
		} else {
			fwidth += pan.Width()
			if pos > 0 {
				fwidth += bl.spacing
			}
		}
	}

	// If there is at least on expanded item, all free space will be occupied
	spaceMiddle := bl.spacing
	var posX float32
	if texpand > 0 {
		// If there is free space, distribute space between expanded items
		totalSpace := parent.ContentWidth() - twidth
		if totalSpace > 0 {
			for _, obj := range parent.Children() {
				pan := obj.(IPanel).GetPanel()
				if !pan.Visible() {
					continue
				}
				// Get item layout parameters or use default
				params := paramsDef
				if pan.layoutParams != nil {
					params = *pan.layoutParams.(*HBoxLayoutParams)
				}
				if params.Expand > 0 {
					iwidth := totalSpace * params.Expand / texpand
					pan.SetWidth(pan.Width() + iwidth)
				}
			}
			// No free space: distribute expanded items widths
		} else {
			for _, obj := range parent.Children() {
				pan := obj.(IPanel).GetPanel()
				if !pan.Visible() {
					continue
				}
				// Get item layout parameters or use default
				params := paramsDef
				if pan.layoutParams != nil {
					params = *pan.layoutParams.(*HBoxLayoutParams)
				}
				if params.Expand > 0 {
					spacing := bl.spacing * (float32(ecount) - 1)
					iwidth := (parent.ContentWidth() - spacing - fwidth - bl.spacing) * params.Expand / texpand
					pan.SetWidth(iwidth)
				}
			}
		}
		// No expanded items: checks block horizontal alignment
	} else {
		// Calculates initial x position which depends
		// on the current horizontal alignment.
		switch bl.alignH {
		case AlignLeft:
			posX = 0
		case AlignCenter:
			posX = (parent.ContentWidth() - twidth) / 2
		case AlignRight:
			posX = parent.ContentWidth() - twidth
		case AlignWidth:
			space := parent.ContentWidth() - twidth + bl.spacing*float32(len(parent.Children())-1)
			if space < 0 {
				space = bl.spacing * float32(len(parent.Children())-1)
			}
			spaceMiddle = space / float32(len(parent.Children())+1)
			posX = spaceMiddle
		default:
			log.Fatal("HBoxLayout: invalid global horizontal alignment")
		}
	}

	// Calculates the Y position of each item considering its vertical alignment
	var posY float32
	height := parent.ContentHeight()
	for pos, obj := range parent.Children() {
		pan := obj.(IPanel).GetPanel()
		if !pan.Visible() {
			continue
		}
		// Get item layout parameters or use default
		params := paramsDef
		if pan.layoutParams != nil {
			params = *pan.layoutParams.(*HBoxLayoutParams)
		}
		cheight := pan.Height()
		switch params.AlignV {
		case AlignNone, AlignTop:
			posY = 0
		case AlignCenter:
			posY = (height - cheight) / 2
		case AlignBottom:
			posY = height - cheight
		case AlignHeight:
			posY = 0
			pan.SetHeight(height)
		default:
			log.Fatal("HBoxLayout: invalid item vertical alignment")
		}
		// Sets the child position
		pan.SetPosition(posX, posY)
		// Calculates next position
		posX += pan.Width()
		if pos < len(parent.Children())-1 {
			posX += spaceMiddle
		}
	}
}
