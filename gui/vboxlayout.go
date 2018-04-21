// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

// VBoxLayout implements a panel layout which arranges the panel children vertically.
// The children can be separated by a space in pixels set by SetSpacing().
// The whole group of children can be aligned vertically by SetAlignV() which can
// accept the following types of alignment:
//
// 	AlignTop: Try to align the group of children to the top if the panel height is
// 	greater the the sum of the children heights + spacing.
//
// 	AlignBottom: Try to align the group of children to the bottoom if the panel height is
// 	greater the the sum of the children heights + spacing.
//
// 	AlignCenter: Try to align the group of children in the center if the panel height is
// 	greater the the sum of the children heights + spacing.
//
// 	AlignHeight: Try to align the individual children vertically with the same same space between each other.
// 	Each individual child can be aligned horizontally by SetLayoutParameters()
//
// If the layout method SetAutoHeight(true) is called, the panel minimum content height will be the
// sum of its children's heights plus the spacing.
//
// If the layout method SetAutoWidth(true) is called, the panel minimum content width will be the
// width of the widest child.
type VBoxLayout struct {
	pan        IPanel
	spacing    float32
	alignV     Align
	autoHeight bool
	autoWidth  bool
}

// VBoxLayoutParams specify the horizontal alignment of each individual child.
type VBoxLayoutParams struct {
	Expand float32 // item expand vertically factor (0 - no expand)
	AlignH Align   // item horizontal alignment
}

// NewVBoxLayout creates and returns a pointer to a new horizontal box layout
func NewVBoxLayout() *VBoxLayout {

	bl := new(VBoxLayout)
	bl.spacing = 0
	bl.alignV = AlignTop
	return bl
}

// SetSpacing sets the horizontal spacing between the items in pixels
// and updates the layout if possible
func (bl *VBoxLayout) SetSpacing(spacing float32) {

	bl.spacing = spacing
	bl.Recalc(bl.pan)
}

// SetAlignV sets the vertical alignment of the whole group of items
// inside the parent panel and updates the layout if possible.
// This only has any effect if there are no expanded items.
func (bl *VBoxLayout) SetAlignV(align Align) {

	bl.alignV = align
	bl.Recalc(bl.pan)
}

// SetAutoHeight sets if the panel minimum height should be the height of
// the largest of its children's height.
func (bl *VBoxLayout) SetAutoHeight(state bool) {

	bl.autoHeight = state
	bl.Recalc(bl.pan)
}

// SetAutoWidth sets if the panel minimum width should be sum of its
// children's width plus the spacing
func (bl *VBoxLayout) SetAutoWidth(state bool) {

	bl.autoWidth = state
	bl.Recalc(bl.pan)
}

// Recalc recalculates and sets the position and sizes of all children
func (bl *VBoxLayout) Recalc(ipan IPanel) {

	// Saves the received panel
	bl.pan = ipan
	if bl.pan == nil {
		return
	}
	parent := ipan.GetPanel()
	if len(parent.Children()) == 0 {
		return
	}

	// If autoHeight is set, get the sum of heights of this panel's children plus the spacings.
	// If the panel content height is less than this height, set its content height to this value.
	if bl.autoHeight {
		var totalHeight float32
		for _, ichild := range parent.Children() {
			child := ichild.(IPanel).GetPanel()
			if !child.Visible() || !child.Bounded() {
				continue
			}
			totalHeight += child.Height()
		}
		// Adds spacing
		totalHeight += bl.spacing * float32(len(parent.Children())-1)
		if parent.ContentHeight() < totalHeight {
			parent.setContentSize(parent.ContentWidth(), totalHeight, false)
		}
	}

	// If autoWidth is set, get the maximum width of all the panel's children
	// and if the panel content width is less than this maximum, set its content width to this value.
	if bl.autoWidth {
		var maxWidth float32
		for _, ichild := range parent.Children() {
			child := ichild.(IPanel).GetPanel()
			if !child.Visible() || !child.Bounded() {
				continue
			}
			if child.Width() > maxWidth {
				maxWidth = child.Width()
			}
		}
		if parent.ContentWidth() < maxWidth {
			parent.setContentSize(maxWidth, parent.ContentHeight(), false)
		}
	}

	// Calculates the total height, expanded height, fixed height and
	// the sum of the expand factor for all items.
	var theight float32
	//var eheight float32
	var fheight float32
	var texpand float32
	ecount := 0
	paramsDef := VBoxLayoutParams{Expand: 0, AlignH: AlignLeft}
	for pos, obj := range parent.Children() {
		pan := obj.(IPanel).GetPanel()
		if !pan.Visible() || !pan.Bounded() {
			continue
		}
		// Get item layout parameters or use default
		params := paramsDef
		if pan.layoutParams != nil {
			params = *pan.layoutParams.(*VBoxLayoutParams)
		}
		// Calculate total height
		theight += pan.Height()
		if pos > 0 {
			theight += bl.spacing
		}
		// Calculate height of expanded items
		if params.Expand > 0 {
			texpand += params.Expand
			//eheight += pan.Height()
			//if pos > 0 {
			//	eheight += bl.spacing
			//}
			ecount++
			// Calculate width of fixed items
		} else {
			fheight += pan.Height()
			if pos > 0 {
				fheight += bl.spacing
			}
		}
	}

	// If there is at least on expanded item, all free space will be occupied
	spaceMiddle := bl.spacing
	var posY float32
	if texpand > 0 {
		// If there is free space, distribute space between expanded items
		totalSpace := parent.ContentHeight() - theight
		if totalSpace > 0 {
			for _, obj := range parent.Children() {
				pan := obj.(IPanel).GetPanel()
				if !pan.Visible() || !pan.Bounded() {
					continue
				}
				// Get item layout parameters or use default
				params := paramsDef
				if pan.layoutParams != nil {
					params = *pan.layoutParams.(*VBoxLayoutParams)
				}
				if params.Expand > 0 {
					iheight := totalSpace * params.Expand / texpand
					pan.SetHeight(pan.Height() + iheight)
				}
			}
			// No free space: distribute expanded items heights
		} else {
			for _, obj := range parent.Children() {
				pan := obj.(IPanel).GetPanel()
				if !pan.Visible() || !pan.Bounded() {
					continue
				}
				// Get item layout parameters or use default
				params := paramsDef
				if pan.layoutParams != nil {
					params = *pan.layoutParams.(*VBoxLayoutParams)
				}
				if params.Expand > 0 {
					spacing := bl.spacing * float32(ecount-1)
					iheight := (parent.ContentHeight() - spacing - fheight - bl.spacing) * params.Expand / texpand
					pan.SetHeight(iheight)
				}
			}
		}
		// No expanded items: checks block vertical alignment
	} else {
		// Calculates initial y position which depends
		// on the current horizontal alignment.
		switch bl.alignV {
		case AlignNone, AlignTop:
			posY = 0
		case AlignCenter:
			posY = (parent.ContentHeight() - theight) / 2
		case AlignBottom:
			posY = parent.ContentHeight() - theight
		case AlignHeight:
			space := parent.ContentHeight() - theight + bl.spacing*float32(len(parent.Children())-1)
			if space < 0 {
				space = bl.spacing * float32(len(parent.Children())-1)
			}
			spaceMiddle = space / float32(len(parent.Children())+1)
			posY = spaceMiddle
		default:
			log.Fatal("VBoxLayout: invalid global vertical alignment")
		}
	}

	// Calculates the X position of each item considering
	// it horizontal alignment
	var posX float32
	width := parent.ContentWidth()
	for pos, obj := range parent.Children() {
		pan := obj.(IPanel).GetPanel()
		if !pan.Visible() || !pan.Bounded() {
			continue
		}
		// Get item layout parameters or use default
		params := paramsDef
		if pan.layoutParams != nil {
			params = *pan.layoutParams.(*VBoxLayoutParams)
		}
		cwidth := pan.Width()
		switch params.AlignH {
		case AlignNone, AlignLeft:
			posX = 0
		case AlignCenter:
			posX = (width - cwidth) / 2
		case AlignRight:
			posX = width - cwidth
		case AlignWidth:
			posX = 0
			pan.SetWidth(width)
		default:
			log.Fatal("VBoxLayout: invalid item horizontal alignment")
		}
		// Sets the child position
		pan.SetPosition(posX, posY)
		// Calculates next position
		posY += pan.Height()
		if pos < len(parent.Children())-1 {
			posY += spaceMiddle
		}
	}
}
