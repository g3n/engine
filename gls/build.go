// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gls

// Generation of API files: glapi.c, glapi.h, consts.go
//go:generate glapi2go -glversion GL_VERSION_3_3 glcorearb.h

// // Platform build flags
// #cgo freebsd CFLAGS:  -DGL_GLEXT_PROTOTYPES
// #cgo freebsd LDFLAGS: 
//
// #cgo linux   CFLAGS:  -DGL_GLEXT_PROTOTYPES
// #cgo linux   LDFLAGS: -ldl
//
// #cgo windows CFLAGS:  -DGL_GLEXT_PROTOTYPES
// #cgo windows LDFLAGS: -lopengl32
//
// #cgo darwin  CFLAGS:  -DGL_GLEXT_PROTOTYPES
// #cgo darwin  LDFLAGS: -framework OpenGL
import "C"
