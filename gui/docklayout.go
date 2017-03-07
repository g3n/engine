// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

type DockLayout struct {
}

type DockLayoutParams struct {
	Edge int
}

const (
	DockTop = iota + 1
	DockRight
	DockBottom
	DockLeft
	DockCenter
)

func NewDockLayout() *DockLayout {

	return new(DockLayout)
}

func (dl *DockLayout) Recalc(ipan IPanel) {

	pan := ipan.GetPanel()
	width := pan.Width()
	topY := float32(0)
	bottomY := pan.Height()
	leftX := float32(0)
	rightX := width

	// Top and bottom first
	for _, iobj := range pan.Children() {
		child := iobj.(IPanel).GetPanel()
		if child.layoutParams == nil {
			continue
		}
		params := child.layoutParams.(*DockLayoutParams)
		if params.Edge == DockTop {
			child.SetPosition(0, topY)
			topY += child.Height()
			child.SetWidth(width)
			continue
		}
		if params.Edge == DockBottom {
			child.SetPosition(0, bottomY-child.Height())
			bottomY -= child.Height()
			child.SetWidth(width)
			continue
		}
	}
	// Left and right
	for _, iobj := range pan.Children() {
		child := iobj.(IPanel).GetPanel()
		if child.layoutParams == nil {
			continue
		}
		params := child.layoutParams.(*DockLayoutParams)
		if params.Edge == DockLeft {
			child.SetPosition(leftX, topY)
			leftX += child.Width()
			child.SetHeight(bottomY - topY)
			continue
		}
		if params.Edge == DockRight {
			child.SetPosition(rightX-child.Width(), topY)
			rightX -= child.Width()
			child.SetHeight(bottomY - topY)
			continue
		}
	}
	// Center (only the first found)
	for _, iobj := range pan.Children() {
		child := iobj.(IPanel).GetPanel()
		if child.layoutParams == nil {
			continue
		}
		params := child.layoutParams.(*DockLayoutParams)
		if params.Edge == DockCenter {
			child.SetPosition(leftX, topY)
			child.SetSize(rightX-leftX, bottomY-topY)
			break
		}
	}
}
