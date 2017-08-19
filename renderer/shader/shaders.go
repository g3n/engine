// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

import (
	"fmt"

	"github.com/g3n/engine/gls"
)

// ProgramInfo contains information for a registered program name
type ProgramInfo struct {
	Vertex   string // Vertex shader name
	Frag     string // Fragment shader name
	Geometry string // Geometry shader name (maybe an empty string)
}

// Internal global maps of shader chunks, shader sources and programs
var chunks = map[string]string{}
var shaders = map[string]string{}
var programs = map[string]ProgramInfo{}

// Chunks returns a map with all registered shader chunks names
// associated with its glsl source code.
func Chunks() map[string]string {

	return chunks
}

// Shaders returns a map with all registered shader names
// associated with
func Shaders() map[string]string {

	return shaders
}

func Programs() map[string]ProgramInfo {

	return programs
}

func AddChunk(name, source string) {

	chunks[name] = source
}

func AddShader(name, source string) {

	shaders[name] = source
}

// AddProgram adds a program name to the global program registry
func AddProgram(name, vertexName, fragName string) {

	programs[name] = ProgramInfo{vertexName, fragName, ""}
}

// SetProgramShader sets the shader type and name for a previously
// specified program name.
// It panics if the specified program or shader name not found or
// if an invalid shader type was specified.
func SetProgramShader(pname string, stype int, sname string) {

	// Checks if program name is valid
	pinfo, ok := programs[pname]
	if !ok {
		panic(fmt.Sprintf("Program name:%s not found", pname))
	}

	// Checks if shader name is valid
	_, ok = shaders[sname]
	if !ok {
		panic(fmt.Sprintf("Shader name:%s not found", sname))
	}

	// Sets the program shader name for the specified type
	switch stype {
	case gls.VERTEX_SHADER:
		pinfo.Vertex = sname
	case gls.FRAGMENT_SHADER:
		pinfo.Frag = sname
	case gls.GEOMETRY_SHADER:
		pinfo.Geometry = sname
	default:
		panic("Invalid shader type")
	}
	programs[pname] = pinfo
}
