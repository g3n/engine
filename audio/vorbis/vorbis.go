// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package vorbis implements the Go bindings of a subset (only one function) of the functions of the libvorbis library
// See API reference at: https://xiph.org/vorbis/doc/libvorbis/reference.html
package vorbis

// #cgo darwin   CFLAGS:  -DGO_DARWIN  -I/usr/include/vorbis -I/usr/local/include/vorbis
// #cgo freebsd  CFLAGS:  -DGO_FREEBSD -I/usr/local/include/vorbis
// #cgo linux    CFLAGS:  -DGO_LINUX   -I/usr/include/vorbis
// #cgo windows  CFLAGS:  -DGO_WINDOWS -I${SRCDIR}/../windows/libvorbis-1.3.5/include/vorbis -I${SRCDIR}/../windows/libogg-1.3.3/include
// #cgo darwin   LDFLAGS: -L/usr/lib -L/usr/local/lib -lvorbis
// #cgo freebsd  LDFLAGS: -L/usr/local/lib -lvorbis
// #cgo linux    LDFLAGS: -lvorbis
// #cgo windows  LDFLAGS: -L${SRCDIR}/../windows/bin -llibvorbis
// #include "codec.h"
import "C"

// VersionString returns a string giving version information for libvorbis
func VersionString() string {

	cstr := C.vorbis_version_string()
	return C.GoString(cstr)
}
