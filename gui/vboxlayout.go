// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type VBoxLayout struct {
	pan     IPanel
	spacing float32 // vertical spacing between the children in pixels.
	alignV  Align   // vertical alignment of the whole block of children
}

// Parameters for individual children
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

// SetAlignH sets the horizontal alignment of the whole group of items
// inside the parent panel and updates the layout if possible.
// This only has any effect if there are no expanded items.
func (bl *VBoxLayout) SetAlignV(align Align) {

	bl.alignV = align
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

	// Calculates the total height, expanded height, fixed height and
	// the sum of the expand factor for all items.
	var theight float32 = 0
	var eheight float32 = 0
	var fheight float32 = 0
	var texpand float32 = 0
	ecount := 0
	paramsDef := VBoxLayoutParams{Expand: 0, AlignH: AlignLeft}
	for pos, obj := range parent.Children() {
		pan := obj.(IPanel).GetPanel()
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
			eheight += pan.Height()
			if pos > 0 {
				eheight += bl.spacing
			}
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
	var posY float32 = 0
	if texpand > 0 {
		// If there is free space, distribute space between expanded items
		totalSpace := parent.ContentHeight() - theight
		if totalSpace > 0 {
			for _, obj := range parent.Children() {
				pan := obj.(IPanel).GetPanel()
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
		case AlignTop:
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
		// Get item layout parameters or use default
		params := paramsDef
		if pan.layoutParams != nil {
			params = *pan.layoutParams.(*VBoxLayoutParams)
		}
		cwidth := pan.Width()
		switch params.AlignH {
		case AlignLeft:
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
