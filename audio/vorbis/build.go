package vorbis

// #cgo darwin   CFLAGS:  -DGO_DARWIN
// #cgo linux    CFLAGS:  -DGO_LINUX   -I../include
// #cgo windows  CFLAGS:  -DGO_WINDOWS -I../include
// #cgo darwin   LDFLAGS:
// #cgo linux    LDFLAGS: -ldl
// #cgo windows  LDFLAGS:
import "C"
