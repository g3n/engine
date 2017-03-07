package vorbis

// #cgo darwin   CFLAGS:  -DGO_DARWIN
// #cgo linux    CFLAGS:  -DGO_LINUX   -I.
// #cgo windows  CFLAGS:  -DGO_WINDOWS -I.
// #cgo darwin   LDFLAGS:
// #cgo linux    LDFLAGS: -ldl
// #cgo windows  LDFLAGS:
import "C"
