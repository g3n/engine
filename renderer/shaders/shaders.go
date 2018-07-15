// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package shaders contains the several shaders used by the engine
package shaders

// Generates shaders sources from this directory and include directory *.glsl files
//go:generate g3nshaders -in=. -out=sources.go -pkg=shaders -v

// ProgramInfo contains information for a registered shader program
type ProgramInfo struct {
	Vertex   string // Vertex shader name
	Fragment string // Fragment shader name
	Geometry string // Geometry shader name (optional)
}

// AddInclude adds a chunk of shader code to the default shaders registry
// which can be included in a shader using the "#include <name>" directive
func AddInclude(name string, source string) {

	if len(name) == 0 || len(source) == 0 {
		panic("Invalid include name and/or source")
	}
	includeMap[name] = source
}

// AddShader add a shader to default shaders registry.
// The specified name can be used when adding programs to the registry
func AddShader(name string, source string) {

	if len(name) == 0 || len(source) == 0 {
		panic("Invalid shader name and/or source")
	}
	shaderMap[name] = source
}

// AddProgram adds a shader program to the default registry of programs.
// Currently up to 3 shaders: vertex, fragment and geometry (optional) can be specified.
func AddProgram(name string, vertex string, frag string, others ...string) {

	if len(name) == 0 || len(vertex) == 0 || len(frag) == 0 {
		panic("Program and/or shader name empty")
	}
	if shaderMap[vertex] == "" {
		panic("Invalid vertex shader name")
	}
	if shaderMap[frag] == "" {
		panic("Invalid vertex shader name")
	}
	var geom = ""
	if len(others) > 0 {
		geom = others[0]
		if shaderMap[geom] == "" {
			panic("Invalid geometry shader name")
		}
	}
	programMap[name] = ProgramInfo{
		Vertex:   vertex,
		Fragment: frag,
		Geometry: geom,
	}
}

// Includes returns list with the names of all include chunks currently in the default shaders registry.
func Includes() []string {

	list := make([]string, 0)
	for name := range includeMap {
		list = append(list, name)
	}
	return list
}

// IncludeSource returns the source code of the specified shader include chunk.
// If the name is not found an empty string is returned.
func IncludeSource(name string) string {

	return includeMap[name]
}

// Shaders returns list with the names of all shaders currently in the default shaders registry.
func Shaders() []string {

	list := make([]string, 0)
	for name := range shaderMap {
		list = append(list, name)
	}
	return list
}

// ShaderSource returns the source code of the specified shader in the default shaders registry.
// If the name is not found an empty string is returned
func ShaderSource(name string) string {

	return shaderMap[name]
}

// Programs returns list with the names of all programs currently in the default shaders registry.
func Programs() []string {

	list := make([]string, 0)
	for name := range programMap {
		list = append(list, name)
	}
	return list
}

// GetProgramInfo returns ProgramInfo struct for the specified program name
// in the default shaders registry
func GetProgramInfo(name string) ProgramInfo {

	return programMap[name]
}
