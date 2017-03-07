// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package shader

import ()

type ProgramInfo struct {
	Vertex string // Vertex shader name
	Frag   string // Fragment shader name
}

var chunks = map[string]string{}
var shaders = map[string]string{}
var programs = map[string]ProgramInfo{}

func Chunks() map[string]string {

	return chunks
}

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

func AddProgram(name, vertexName, fragName string) {

	programs[name] = ProgramInfo{vertexName, fragName}
}
