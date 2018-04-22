// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gui

var defaultStyle *Style

// init sets the default style
func init() {

	defaultStyle = NewDarkStyle()
}

// StyleDefault returns a pointer to the current default style
func StyleDefault() *Style {

	return defaultStyle
}

// SetStyleDefault sets the default style
func SetStyleDefault(s *Style) {

	defaultStyle = s
}
