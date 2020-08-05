// +build !go1.13

package wasm

import (
	"syscall/js"
)

func SliceToTypedArray(s interface{}) (val js.Value, free func()) {
	dataTA := js.TypedArrayOf(s)
	return dataTA.Value, func(){dataTA.Release()}
}

func Equal(a, b js.Value) bool {
	return a == b
}