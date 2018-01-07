// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vorbis implements the Go bindings of a subset (only one function) of the functions of the libvorbis library
// See API reference at: https://xiph.org/vorbis/doc/libvorbis/reference.html
package vorbis

// #cgo darwin   CFLAGS:  -DGO_DARWIN  -I/usr/include/vorbis
// #cgo linux    CFLAGS:  -DGO_LINUX   -I/usr/include/vorbis
// #cgo windows  CFLAGS:  -DGO_WINDOWS -I/usr/include/vorbis
// #cgo darwin   LDFLAGS: -lvorbis
// #cgo linux    LDFLAGS: -lvorbis
// #cgo windows  LDFLAGS: -lvorbis
// #include "codec.h"
import "C"

// VersionString returns a string giving version information for libvorbis
func VersionString() string {

	cstr := C.vorbis_version_string()
	return C.GoString(cstr)
}
