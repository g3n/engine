// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gls

import (
	"fmt"
)

// Uniform represents an OpenGL uniform.
type Uniform struct {
	name      string // base name
	nameIdx   string // cached indexed name
	handle    uint32 // program handle
	location  int32  // last cached location
	lastIndex int32  // last index
}

// Init initializes this uniform location cache and sets its name.
func (u *Uniform) Init(name string) {

	u.name = name
	u.handle = 0     // invalid program handle
	u.location = -1  // invalid location
	u.lastIndex = -1 // invalid index
}

// Name returns the uniform name.
func (u *Uniform) Name() string {

	return u.name
}

// Location returns the location of this uniform for the current shader program.
// The returned location can be -1 if not found.
func (u *Uniform) Location(gs *GLS) int32 {

	handle := gs.prog.Handle()
	if handle != u.handle {
		u.location = gs.prog.GetUniformLocation(u.name)
		u.handle = handle
	}
	return u.location
}

// LocationIdx returns the location of this indexed uniform for the current shader program.
// The returned location can be -1 if not found.
func (u *Uniform) LocationIdx(gs *GLS, idx int32) int32 {

	if idx != u.lastIndex {
		u.nameIdx = fmt.Sprintf("%s[%d]", u.name, idx)
		u.lastIndex = idx
		u.handle = 0
	}
	handle := gs.prog.Handle()
	if handle != u.handle {
		u.location = gs.prog.GetUniformLocation(u.nameIdx)
		u.handle = handle
	}
	return u.location
}
