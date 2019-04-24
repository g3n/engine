// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gls

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Program represents an OpenGL program.
// It must have Vertex and Fragment shaders.
// It can also have a Geometry shader.
type Program struct {
	gs         *GLS             // Reference to OpenGL state
	ShowSource bool             // Show source code in error messages
	handle     uint32           // OpenGL program handle
	shaders    []shaderInfo     // List of shaders for this program
	uniforms   map[string]int32 // List of uniforms
}

// shaderInfo contains OpenGL-related shader information.
type shaderInfo struct {
	stype  uint32 // OpenGL shader type (VERTEX_SHADER, FRAGMENT_SHADER, or GEOMETRY_SHADER)
	source string // Shader source code
	handle uint32 // OpenGL shader handle
}

// Map from shader types to names.
var shaderNames = map[uint32]string{
	VERTEX_SHADER:   "Vertex Shader",
	FRAGMENT_SHADER: "Fragment Shader",
	GEOMETRY_SHADER: "Geometry Shader",
}

// NewProgram creates and returns a new empty shader program object.
// Use this type methods to add shaders and build the final program.
func (gs *GLS) NewProgram() *Program {

	prog := new(Program)
	prog.gs = gs

	prog.shaders = make([]shaderInfo, 0)
	prog.uniforms = make(map[string]int32)
	prog.ShowSource = true
	return prog
}

// Handle returns the OpenGL handle of this program.
func (prog *Program) Handle() uint32 {

	return prog.handle
}

// AddShader adds a shader to this program.
// This must be done before the program is built.
func (prog *Program) AddShader(stype uint32, source string) {

	// Check if program already built
	if prog.handle != 0 {
		log.Fatal("Program already built")
	}
	prog.shaders = append(prog.shaders, shaderInfo{stype, source, 0})
}

// DeleteShaders deletes all of this program's shaders from OpenGL.
func (prog *Program) DeleteShaders() {

	for _, shaderInfo := range prog.shaders {
		if shaderInfo.handle != 0 {
			prog.gs.DeleteShader(shaderInfo.handle)
			shaderInfo.handle = 0
		}
	}
}

// Build builds the program, compiling and linking the previously supplied shaders.
func (prog *Program) Build() error {

	// Check if program already built
	if prog.handle != 0 {
		return fmt.Errorf("program already built")
	}

	// Check if shaders were provided
	if len(prog.shaders) == 0 {
		return fmt.Errorf("no shaders supplied")
	}

	// Create program
	prog.handle = prog.gs.CreateProgram()
	if prog.handle == 0 {
		return fmt.Errorf("error creating program")
	}

	// Clean unused GL allocated resources
	defer prog.DeleteShaders()

	// Compile and attach shaders
	for _, sinfo := range prog.shaders {
		shader, err := prog.CompileShader(sinfo.stype, sinfo.source)
		if err != nil {
			prog.gs.DeleteProgram(prog.handle)
			prog.handle = 0
			msg := fmt.Sprintf("error compiling %s: %s", shaderNames[sinfo.stype], err)
			if prog.ShowSource {
				msg += FormatSource(sinfo.source)
			}
			return errors.New(msg)
		}
		sinfo.handle = shader
		prog.gs.AttachShader(prog.handle, shader)
	}

	// Link program and check for errors
	prog.gs.LinkProgram(prog.handle)
	var status int32
	prog.gs.GetProgramiv(prog.handle, LINK_STATUS, &status)
	if status == FALSE {
		log := prog.gs.GetProgramInfoLog(prog.handle)
		prog.handle = 0
		return fmt.Errorf("error linking program: %v", log)
	}

	return nil
}

// GetAttribLocation returns the location of the specified attribute
// in this program. This location is internally cached.
func (prog *Program) GetAttribLocation(name string) int32 {

	return prog.gs.GetAttribLocation(prog.handle, name)
}

// GetUniformLocation returns the location of the specified uniform in this program.
// This location is internally cached.
func (prog *Program) GetUniformLocation(name string) int32 {

	// Try to get from the cache
	loc, ok := prog.uniforms[name]
	if ok {
		prog.gs.stats.UnilocHits++
		return loc
	}

	// Get location from OpenGL
	loc = prog.gs.GetUniformLocation(prog.handle, name)
	prog.gs.stats.UnilocMiss++

	// Cache result
	prog.uniforms[name] = loc
	if loc < 0 {
		log.Warn("Program.GetUniformLocation(%s): NOT FOUND", name)
	}

	return loc
}

// CompileShader creates and compiles an OpenGL shader of the specified type, with
// the specified source code, and returns a non-zero value by which it can be referenced.
func (prog *Program) CompileShader(stype uint32, source string) (uint32, error) {

	// Create shader object
	shader := prog.gs.CreateShader(stype)
	if shader == 0 {
		return 0, fmt.Errorf("error creating shader")
	}

	// Set shader source and compile it
	prog.gs.ShaderSource(shader, source)
	prog.gs.CompileShader(shader)

	// Get the shader compiler log
	slog := prog.gs.GetShaderInfoLog(shader)

	// Get the shader compile status
	var status int32
	prog.gs.GetShaderiv(shader, COMPILE_STATUS, &status)
	if status == FALSE {
		return shader, fmt.Errorf("%s", slog)
	}

	// If the shader compiled OK but the log has data,
	// log this data instead of returning error
	if len(slog) > 2 {
		log.Warn("%s", slog)
	}

	return shader, nil
}

// FormatSource returns the supplied program source code with
// line numbers prepended.
func FormatSource(source string) string {

	// Reads all lines from the source string
	lines := make([]string, 0)
	buf := bytes.NewBuffer([]byte(source))
	for {
		line, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		lines = append(lines, string(line[:len(line)-1]))
	}
	// Adds a final line terminator
	lines = append(lines, "\n")

	// Prepends the line number for each line
	ndigits := len(strconv.Itoa(len(lines)))
	format := "%0" + strconv.Itoa(ndigits) + "d:%s"
	formatted := make([]string, 0)
	for pos, l := range lines {
		fline := fmt.Sprintf(format, pos+1, l)
		formatted = append(formatted, fline)
	}

	return strings.Join(formatted, "\n")
}
