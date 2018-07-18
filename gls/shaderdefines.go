// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gls

// ShaderDefines is a store of shader defines ("#define <key> <value>").
type ShaderDefines map[string]string

// NewShaderDefines creates and returns a pointer to a ShaderDefines object.
func NewShaderDefines() *ShaderDefines {

	sd := ShaderDefines(make(map[string]string))
	return &sd
}

// Set sets a shader define with the specified value.
func (sd *ShaderDefines) Set(name, value string) {

	(*sd)[name] = value
}

// Unset removes the specified name from the shader defines.
func (sd *ShaderDefines) Unset(name string) {

	delete(*sd, name)
}

// Add adds to this ShaderDefines all the key-value pairs in the specified ShaderDefines.
func (sd *ShaderDefines) Add(other *ShaderDefines) {

	for k, v := range map[string]string(*other){
		(*sd)[k] = v
	}
}

// Equals compares two ShaderDefines and return true if they contain the same key-value pairs.
func (sd *ShaderDefines) Equals(other *ShaderDefines) bool {

	if sd == nil && other == nil {
		return true
	}
	if sd != nil && other != nil {
		if len(*sd) != len(*other) {
			return false
		}
		for k := range map[string]string(*sd) {
			v1, ok1 := (*sd)[k]
			v2, ok2 := (*other)[k]
			if v1 != v2 || ok1 != ok2 {
				return false
			}
		}
		return true
	}
	// One is nil and the other is not nil
	return false
}
