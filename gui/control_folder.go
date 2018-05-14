// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"fmt"
)

// ControlFolder represents a folder with controls.
type ControlFolder struct {
	Folder                      // Embedded folder
	tree   Tree                 // control tree
	styles *ControlFolderStyles // Pointer to styles
}

// ControlFolderStyles contains the styling for the valid GUI states of the components of a ControlFolder.
type ControlFolderStyles struct {
	Folder *FolderStyles
	Tree   *TreeStyles
}

// ControlFolderGroup represents a group of controls in the control folder.
type ControlFolderGroup struct {
	control *ControlFolder
	node    *TreeNode
}

// NewControlFolder creates and returns a pointer to a new control folder widget
// with the specified text and initial width
func NewControlFolder(text string, width float32) *ControlFolder {

	f := new(ControlFolder)
	f.Initialize(text, width)
	return f
}

// Initialize initializes the control folder with the specified text and initial width
// It is normally used when the control folder is embedded in another object
func (f *ControlFolder) Initialize(text string, width float32) {

	f.styles = &StyleDefault().ControlFolder
	f.tree.Initialize(width, width)
	f.tree.SetStyles(f.styles.Tree)
	f.tree.SetAutoHeight(600)
	f.tree.SetAutoWidth(400)

	f.Folder.Initialize(text, width, &f.tree)
	f.Folder.SetStyles(f.styles.Folder)
	f.Folder.SetAlignRight(false)
}

// Clear clears the control folder's tree
func (f *ControlFolder) Clear() {

	f.tree.Clear()
}

// RemoveAt removes the IPanel at the specified position from the control folder's tree
func (f *ControlFolder) RemoveAt(pos int) IPanel {

	return f.tree.RemoveAt(pos)
}

// AddPanel adds an IPanel to the control folder's tree
func (f *ControlFolder) AddPanel(pan IPanel) {

	f.tree.Add(pan)
}

// AddCheckBox adds a checkbox to the control folder's tree
func (f *ControlFolder) AddCheckBox(text string) *CheckRadio {

	cb := NewCheckBox(text)
	f.tree.Add(cb)
	return cb
}

// AddSlider adds a slider to the control folder's tree
func (f *ControlFolder) AddSlider(text string, sf, v float32) *Slider {

	cont, slider := f.newSlider(text, sf, v)
	f.tree.Add(cont)
	return slider
}

// AddGroup adds a group to the control folder
func (f *ControlFolder) AddGroup(text string) *ControlFolderGroup {

	g := new(ControlFolderGroup)
	g.control = f
	g.node = f.tree.AddNode(text)
	return g
}

// SetStyles set the folder styles overriding the default style
func (f *ControlFolder) SetStyles(fs *ControlFolderStyles) {

	f.styles = fs

	f.Folder.styles = fs.Folder
	f.tree.styles = fs.Tree

	f.tree.update()
	f.Folder.update()

}

// AddCheckBox adds a checkbox to the control folder group
func (g *ControlFolderGroup) AddCheckBox(text string) *CheckRadio {

	cb := NewCheckBox(text)
	g.node.Add(cb)
	return cb
}

// AddSlider adds a slider to the control folder group
func (g *ControlFolderGroup) AddSlider(text string, sf, v float32) *Slider {

	cont, slider := g.control.newSlider(text, sf, v)
	g.node.Add(cont)
	return slider
}

// AddPanel adds a panel to the control folder group
func (g *ControlFolderGroup) AddPanel(pan IPanel) {

	g.node.Add(pan)
}

func (f *ControlFolder) newSlider(text string, sf, value float32) (IPanel, *Slider) {

	// Creates container panel for the label and slider
	cont := NewPanel(200, 32)
	hbox := NewHBoxLayout()
	hbox.spacing = 4
	cont.SetLayout(hbox)

	// Adds label
	l := NewImageLabel(text)
	l.SetLayoutParams(&HBoxLayoutParams{AlignV: AlignCenter})
	cont.Add(l)

	// Adds slider
	s := NewHSlider(100, l.Height())
	s.SetScaleFactor(sf)
	s.SetScaleFactor(sf)
	s.SetValue(value)
	s.SetText(fmt.Sprintf("%1.1f", value))
	s.Subscribe(OnChange, func(evname string, ev interface{}) {
		s.SetText(fmt.Sprintf("%1.1f", s.Value()))
	})
	s.SetLayoutParams(&HBoxLayoutParams{AlignV: AlignCenter, Expand: 1})
	cont.Add(s)

	return cont, s
}
