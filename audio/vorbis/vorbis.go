// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 Package vorbis implements the Go bindings of a subset (only one function) of the functions of the libvorbis library
 It also implements a loader so the library can be dynamically loaded
 See API reference at: https://xiph.org/vorbis/doc/libvorbis/reference.html
*/
package vorbis

// #include "loader.h"
import "C"

import (
	"fmt"
)

// Load tries to load dinamically libvorbis share library/dll
func Load() error {

	// Loads libvorbis
	cres := C.vorbis_load()
	if cres != 0 {
		return fmt.Errorf("Error loading libvorbis shared library/dll")
	}
	return nil
}

// VersionString returns a string giving version information for libvorbis
func VersionString() string {

	cstr := C.vorbis_version_string()
	return C.GoString(cstr)
}
