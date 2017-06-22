// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gls

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/g3n/engine/math32"
	"io"
	"strconv"
	"strings"
)

// Shader Program Object
type Program struct {
	// Shows source code in error messages
	ShowSource bool
	gs         *GLS
	handle     uint32
	shaders    []shaderInfo
	uniforms   map[string]int32
	Specs      interface{}
}

type shaderInfo struct {
	stype   uint32
	source  string
	defines map[string]interface{}
	handle  uint32
}

// Map shader types to names
var shaderNames = map[uint32]string{
	VERTEX_SHADER:   "Vertex Shader",
	FRAGMENT_SHADER: "Fragment Shader",
}

// NewProgram creates a new empty shader program object.
// Use this type methods to add shaders and build the final program.
func (gs *GLS) NewProgram() *Program {

	prog := new(Program)
	prog.gs = gs

	prog.shaders = make([]shaderInfo, 0)
	prog.uniforms = make(map[string]int32)
	prog.ShowSource = true
	return prog
}

// AddShaders adds a shader to this program.
// This must be done before the program is built.
func (prog *Program) AddShader(stype uint32, source string, defines map[string]interface{}) {

	if prog.handle != 0 {
		log.Fatal("Program already built")
	}
	prog.shaders = append(prog.shaders, shaderInfo{stype, source, defines, 0})
}

// Build builds the program compiling and linking the previously supplied shaders.
func (prog *Program) Build() error {

	if prog.handle != 0 {
		return fmt.Errorf("Program already built")
	}

	// Checks if shaders were provided
	if len(prog.shaders) == 0 {
		return fmt.Errorf("No shaders supplied")
	}

	// Create program
	prog.handle = prog.gs.CreateProgram()
	if prog.handle == 0 {
		return fmt.Errorf("Error creating program")
	}

	// Clean unused GL allocated resources
	defer func() {
		for _, sinfo := range prog.shaders {
			if sinfo.handle != 0 {
				prog.gs.DeleteShader(sinfo.handle)
				sinfo.handle = 0
			}
		}
	}()

	// Compiles and attach each shader
	for _, sinfo := range prog.shaders {
		// Creates string with defines from specified parameters
		deflines := make([]string, 0)
		if sinfo.defines != nil {
			for pname, pval := range sinfo.defines {
				line := "#define " + pname + " "
				switch val := pval.(type) {
				case bool:
					if val {
						deflines = append(deflines, line)
					}
				case float32:
					line += strconv.FormatFloat(float64(val), 'f', -1, 32)
					deflines = append(deflines, line)
				default:
					panic("Parameter type not supported")
				}
			}
		}
		deftext := strings.Join(deflines, "\n")
		// Compile shader
		shader, err := prog.CompileShader(sinfo.stype, sinfo.source+deftext)
		if err != nil {
			prog.gs.DeleteProgram(prog.handle)
			prog.handle = 0
			msg := fmt.Sprintf("Error compiling %s: %s", shaderNames[sinfo.stype], err)
			if prog.ShowSource {
				source := FormatSource(sinfo.source + deftext)
				msg += source
			}
			return errors.New(msg)
		}
		sinfo.handle = shader
		prog.gs.AttachShader(prog.handle, shader)
	}

	// Link program and checks for errors
	prog.gs.LinkProgram(prog.handle)
	var status int32
	prog.gs.GetProgramiv(prog.handle, LINK_STATUS, &status)
	if status == FALSE {
		log := prog.gs.GetProgramInfoLog(prog.handle)
		prog.handle = 0
		return fmt.Errorf("Error linking program: %v", log)
	}

	return nil
}

// Handle returns the handle of this program
func (prog *Program) Handle() uint32 {

	return prog.handle
}

// GetAttributeLocation returns the location of the specified attribute
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
		return loc
	}
	// Get location from GL
	loc = prog.gs.GetUniformLocation(prog.handle, name)
	// Cache result
	prog.uniforms[name] = loc
	if loc < 0 {
		log.Warn("GetUniformLocation(%s) NOT FOUND", name)
	}
	return loc
}

// SetUniformInt sets this program uniform variable specified by
// its location to the the value of the specified int
func (prog *Program) SetUniformInt(loc int32, v int) {

	prog.gs.Uniform1i(loc, int32(v))
}

// SetUniformFloat sets this program uniform variable specified by
// its location to the the value of the specified float
func (prog *Program) SetUniformFloat(loc int32, v float32) {

	prog.gs.Uniform1f(loc, v)
}

// SetUniformVector2 sets this program uniform variable specified by
// its location to the the value of the specified Vector2
func (prog *Program) SetUniformVector2(loc int32, v *math32.Vector2) {

	prog.gs.Uniform2f(loc, v.X, v.Y)
}

// SetUniformVector3 sets this program uniform variable specified by
// its location to the the value of the specified Vector3
func (prog *Program) SetUniformVector3(loc int32, v *math32.Vector3) {

	prog.gs.Uniform3f(loc, v.X, v.Y, v.Z)
}

// SetUniformVector4 sets this program uniform variable specified by
// its location to the the value of the specified Vector4
func (prog *Program) SetUniformVector4(loc int32, v *math32.Vector4) {

	prog.gs.Uniform4f(loc, v.X, v.Y, v.Z, v.W)
}

// SetUniformMatrix3 sets this program uniform variable specified by
// its location with the values from the specified Matrix3.
func (prog *Program) SetUniformMatrix3(loc int32, m *math32.Matrix3) {

	prog.gs.UniformMatrix3fv(loc, 1, false, &m[0])
}

// SetUniformMatrix4 sets this program uniform variable specified by
// its location with the values from the specified Matrix4.
func (prog *Program) SetUniformMatrix4(loc int32, m *math32.Matrix4) {

	prog.gs.UniformMatrix4fv(loc, 1, false, &m[0])
}

// SetUniformIntByName sets this program uniform variable specified by
// its name to the value of the specified int.
// The specified name location is cached internally.
func (prog *Program) SetUniformIntByName(name string, v int) {

	prog.gs.Uniform1i(prog.GetUniformLocation(name), int32(v))
}

// SetUniformFloatByName sets this program uniform variable specified by
// its name to the value of the specified float32.
// The specified name location is cached internally.
func (prog *Program) SetUniformFloatByName(name string, v float32) {

	prog.gs.Uniform1f(prog.GetUniformLocation(name), v)
}

// SetUniformVector2ByName sets this program uniform variable specified by
// its name to the values from the specified Vector2.
// The specified name location is cached internally.
func (prog *Program) SetUniformVector2ByName(name string, v *math32.Vector2) {

	prog.gs.Uniform2f(prog.GetUniformLocation(name), v.X, v.Y)
}

// SetUniformVector3ByName sets this program uniform variable specified by
// its name to the values from the specified Vector3.
// The specified name location is cached internally.
func (prog *Program) SetUniformVector3ByName(name string, v *math32.Vector3) {

	prog.gs.Uniform3f(prog.GetUniformLocation(name), v.X, v.Y, v.Z)
}

// SetUniformVector4ByName sets this program uniform variable specified by
// its name to the values from the specified Vector4.
// The specified name location is cached internally.
func (prog *Program) SetUniformVector4ByName(name string, v *math32.Vector4) {

	prog.gs.Uniform4f(prog.GetUniformLocation(name), v.X, v.Y, v.Z, v.W)
}

// SetUniformMatrix3ByName sets this program uniform variable specified by
// its name with the values from the specified Matrix3.
// The specified name location is cached internally.
func (prog *Program) SetUniformMatrix3ByName(name string, m *math32.Matrix3) {

	prog.gs.UniformMatrix3fv(prog.GetUniformLocation(name), 1, false, &m[0])
}

// SetUniformMatrix4ByName sets this program uniform variable specified by
// its name with the values from the specified Matrix4.
// The location of the name is cached internally.
func (prog *Program) SetUniformMatrix4ByName(name string, m *math32.Matrix4) {

	prog.gs.UniformMatrix4fv(prog.GetUniformLocation(name), 1, false, &m[0])
}

// SetUniformColorByName set this program uniform variable specified by
// its name to the values from the specified Color
// The specified name location is cached internally.
func (prog *Program) SetUniformColorByName(name string, c *math32.Color) {

	prog.gs.Uniform3f(prog.GetUniformLocation(name), c.R, c.G, c.B)
}

// SetUniformColor4ByName set this program uniform variable specified by
// its name to the values from the specified Color4
// The specified name location is cached internally.
func (prog *Program) SetUniformColor4ByName(name string, c *math32.Color4) {

	prog.gs.Uniform4f(prog.GetUniformLocation(name), c.R, c.G, c.B, c.A)
}

// CompileShader creates and compiles a shader of the specified type and with
// the specified source code and returns a non-zero value by which
// it can be referenced.
func (prog *Program) CompileShader(stype uint32, source string) (uint32, error) {

	// Creates shader object
	shader := prog.gs.CreateShader(stype)
	if shader == 0 {
		return 0, fmt.Errorf("Error creating shader")
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
	// logs this data instead of returning error
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
