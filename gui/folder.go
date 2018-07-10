// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

import (
	"github.com/g3n/engine/math32"
)

// Folder represents a folder GUI element.
type Folder struct {
	Panel               // Embedded panel
	label        Label  // Folder label
	icon         Label  // Folder icon
	contentPanel IPanel // Content panel
	styles       *FolderStyles
	cursorOver   bool
	alignRight   bool
}

// FolderStyle contains the styling of a Folder.
type FolderStyle struct {
	PanelStyle
	FgColor math32.Color4
	Icons   [2]string
}

// FolderStyles contains a FolderStyle for each valid GUI state.
type FolderStyles struct {
	Normal   FolderStyle
	Over     FolderStyle
	Focus    FolderStyle
	Disabled FolderStyle
}

// NewFolder creates and returns a pointer to a new folder widget
// with the specified text and initial width.
func NewFolder(text string, width float32, contentPanel IPanel) *Folder {

	f := new(Folder)
	f.Initialize(text, width, contentPanel)
	return f
}

// Initialize initializes the Folder with the specified text and initial width
// It is normally used when the folder is embedded in another object.
func (f *Folder) Initialize(text string, width float32, contentPanel IPanel) {

	f.Panel.Initialize(width, 0)
	f.styles = &StyleDefault().Folder

	// Initialize label
	f.label.initialize(text, StyleDefault().Font)
	f.Panel.Add(&f.label)

	// Create icon
	f.icon.initialize("", StyleDefault().FontIcon)
	f.icon.SetFontSize(StyleDefault().Label.PointSize * 1.3)
	f.Panel.Add(&f.icon)

	// Setup content panel
	f.contentPanel = contentPanel
	contentPanel.GetPanel().bounded = false
	contentPanel.GetPanel().SetVisible(false)
	f.Panel.Add(f.contentPanel)

	// Set event callbacks
	f.Panel.Subscribe(OnMouseDown, f.onMouse)
	f.Panel.Subscribe(OnCursorEnter, f.onCursor)
	f.Panel.Subscribe(OnCursorLeave, f.onCursor)

	f.alignRight = true
	f.update()
	f.recalc()
}

// SetStyles set the folder styles overriding the default style.
func (f *Folder) SetStyles(fs *FolderStyles) {

	f.styles = fs
	f.update()
}

// SetAlignRight sets the side of the alignment of the content panel
// in relation to the folder.
func (f *Folder) SetAlignRight(state bool) {

	f.alignRight = state
	f.recalc()
}

// TotalHeight returns this folder total height
// considering the contents panel, if visible.
func (f *Folder) TotalHeight() float32 {

	height := f.Height()
	if f.contentPanel.GetPanel().Visible() {
		height += f.contentPanel.GetPanel().Height()
	}
	return height
}

// onMouse receives mouse button events over the folder panel.
func (f *Folder) onMouse(evname string, ev interface{}) {

	switch evname {
	case OnMouseDown:
		cont := f.contentPanel.GetPanel()
		if !cont.Visible() {
			cont.SetVisible(true)
		} else {
			cont.SetVisible(false)
		}
		f.update()
		f.recalc()
	default:
		return
	}
}

// onCursor receives cursor events over the folder panel
func (f *Folder) onCursor(evname string, ev interface{}) {

	switch evname {
	case OnCursorEnter:
		f.cursorOver = true
		f.update()
	case OnCursorLeave:
		f.cursorOver = false
		f.update()
	default:
		return
	}
}

// update updates the folder visual state
func (f *Folder) update() {

	if f.cursorOver {
		f.applyStyle(&f.styles.Over)
		return
	}
	f.applyStyle(&f.styles.Normal)
}

// applyStyle applies the specified style
func (f *Folder) applyStyle(s *FolderStyle) {

	f.Panel.ApplyStyle(&s.PanelStyle)

	icode := 0
	if f.contentPanel.GetPanel().Visible() {
		icode = 1
	}
	f.icon.SetText(string(s.Icons[icode]))
	f.icon.SetColor4(&s.FgColor)
	f.label.SetBgColor4(&s.BgColor)
	f.label.SetColor4(&s.FgColor)
}

func (f *Folder) recalc() {

	// icon position
	f.icon.SetPosition(0, 0)

	// Label position and width
	f.label.SetPosition(f.icon.Width()+4, 0)
	f.Panel.SetContentHeight(f.label.Height())

	// Sets position of the base folder scroller panel
	cont := f.contentPanel.GetPanel()
	if f.alignRight {
		cont.SetPosition(0, f.Panel.Height())
	} else {
		dx := cont.Width() - f.Panel.Width()
		cont.SetPosition(-dx, f.Panel.Height())
	}
}
