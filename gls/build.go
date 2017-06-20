// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package gls

// Generation of API files: glapi.c, glapi.h, consts.go
//go:generate glapi2go glcorearb.h

// // Platform build flags
// #cgo freebsd CFLAGS: -DGL_GLEXT_PROTOTYPES
// #cgo freebsd LDFLAGS: -ldl -lGL
//
// #cgo linux CFLAGS: -DGL_GLEXT_PROTOTYPES
// #cgo linux LDFLAGS: -ldl -lGL
//
// #cgo windows CFLAGS: -DGL_GEXT_PROTOTYPES
// #cgo windows LDFLAGS: -lopengl32
import "C"
