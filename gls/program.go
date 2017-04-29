// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gls

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/g3n/engine/math32"
	"github.com/go-gl/gl/v3.3-core/gl"
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
	gl.VERTEX_SHADER:   "Vertex Shader",
	gl.FRAGMENT_SHADER: "Fragment Shader",
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
	prog.handle = gl.CreateProgram()
	if prog.handle == 0 {
		return fmt.Errorf("Error creating program")
	}

	// Clean unused GL allocated resources
	defer func() {
		for _, sinfo := range prog.shaders {
			if sinfo.handle != 0 {
				gl.DeleteShader(sinfo.handle)
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
		shader, err := CompileShader(sinfo.stype, sinfo.source+deftext)
		if err != nil {
			gl.DeleteProgram(prog.handle)
			prog.handle = 0
			msg := fmt.Sprintf("Error compiling %s: %s", shaderNames[sinfo.stype], err)
			if prog.ShowSource {
				source := FormatSource(sinfo.source + deftext)
				msg += source
			}
			return errors.New(msg)
		}
		sinfo.handle = shader
		gl.AttachShader(prog.handle, shader)
	}

	// Link program and checks for errors
	gl.LinkProgram(prog.handle)
	var status int32
	gl.GetProgramiv(prog.handle, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(prog.handle, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(prog.handle, logLength, nil, gl.Str(log))
		prog.handle = 0
		return fmt.Errorf("Error linking program: %v", log)
	}

	return nil
}

// Handle returns the handle of this program
func (prog *Program) Handle() uint32 {

	return prog.handle
}

// GetActiveUniformBlockSize returns the minimum number of bytes
// to contain the data for the uniform block specified by its index.
func (prog *Program) GetActiveUniformBlockSize(ubindex uint32) int32 {

	var uboSize int32
	gl.GetActiveUniformBlockiv(prog.handle, ubindex, gl.UNIFORM_BLOCK_DATA_SIZE, &uboSize)
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("GetUniformBlockSize(%v) error: %d", ubindex, ecode)
		}
	}
	return uboSize
}

// GetActiveUniformsiv returns information about the specified uniforms
// specified by its indices
func (prog *Program) GetActiveUniformsiv(indices []uint32, pname uint32) []int32 {

	data := make([]int32, len(indices))
	gl.GetActiveUniformsiv(prog.handle, int32(len(indices)), &indices[0], pname, &data[0])
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("GetActiveUniformsiv() error: %d", ecode)
		}
	}
	return data
}

// GetAttributeLocation returns the location of the specified attribute
// in this program. This location is internally cached.
func (prog *Program) GetAttribLocation(name string) int32 {

	loc := gl.GetAttribLocation(prog.handle, gl.Str(name+"\x00"))
	prog.gs.checkError("GetAttribLocation")
	return loc
}

// GetUniformBlockIndex returns the index of the named uniform block.
// If the supplied name is not valid, the function returns gl.INVALID_INDEX
func (prog *Program) GetUniformBlockIndex(name string) uint32 {

	index := gl.GetUniformBlockIndex(prog.handle, gl.Str(name+"\x00"))
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("GetUniformBlockIndex(%s) error", name)
		}
	}
	return index
}

// GetUniformIndices returns the indices for each specified named
// uniform. If an specified name is not valid the corresponding
// index value will be gl.INVALID_INDEX
func (prog *Program) GetUniformIndices(names []string) []uint32 {

	// Add C terminators to uniform names
	for _, s := range names {
		s += "\x00"
	}
	unames, freefunc := gl.Strs(names...)

	indices := make([]uint32, len(names))
	gl.GetUniformIndices(prog.handle, int32(len(names)), unames, &indices[0])
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("GetUniformIndices() error: %d", ecode)
		}
	}
	freefunc()
	return indices
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
	loc = gl.GetUniformLocation(prog.handle, gl.Str(name+"\x00"))
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("GetUniformLocation(%s) error: %d", name, ecode)
		}
	}
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

	gl.Uniform1i(loc, int32(v))
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformInt() error: %d", ecode)
		}
	}
}

// SetUniformFloat sets this program uniform variable specified by
// its location to the the value of the specified float
func (prog *Program) SetUniformFloat(loc int32, v float32) {

	gl.Uniform1f(loc, v)
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformFloat() error: %d", ecode)
		}
	}
}

// SetUniformVector2 sets this program uniform variable specified by
// its location to the the value of the specified Vector2
func (prog *Program) SetUniformVector2(loc int32, v *math32.Vector2) {

	gl.Uniform2f(loc, v.X, v.Y)
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformVector2() error: %d", ecode)
		}
	}
}

// SetUniformVector3 sets this program uniform variable specified by
// its location to the the value of the specified Vector3
func (prog *Program) SetUniformVector3(loc int32, v *math32.Vector3) {

	gl.Uniform3f(loc, v.X, v.Y, v.Z)
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformVector3() error: %d", ecode)
		}
	}
}

// SetUniformVector4 sets this program uniform variable specified by
// its location to the the value of the specified Vector4
func (prog *Program) SetUniformVector4(loc int32, v *math32.Vector4) {

	gl.Uniform4f(loc, v.X, v.Y, v.Z, v.W)
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformVector4() error: %d", ecode)
		}
	}
}

// SetUniformMatrix3 sets this program uniform variable specified by
// its location with the values from the specified Matrix3.
func (prog *Program) SetUniformMatrix3(loc int32, m *math32.Matrix3) {

	gl.UniformMatrix3fv(loc, 1, false, &m[0])
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformMatrix3() error: %d", ecode)
		}
	}
}

// SetUniformMatrix4 sets this program uniform variable specified by
// its location with the values from the specified Matrix4.
func (prog *Program) SetUniformMatrix4(loc int32, m *math32.Matrix4) {

	gl.UniformMatrix4fv(loc, 1, false, &m[0])
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformMatrix4() error: %d", ecode)
		}
	}
}

// SetUniformIntByName sets this program uniform variable specified by
// its name to the value of the specified int.
// The specified name location is cached internally.
func (prog *Program) SetUniformIntByName(name string, v int) {

	gl.Uniform1i(prog.GetUniformLocation(name), int32(v))
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("GetUniformIntByName(%s) error: %d", name, ecode)
		}
	}
}

// SetUniformFloatByName sets this program uniform variable specified by
// its name to the value of the specified float32.
// The specified name location is cached internally.
func (prog *Program) SetUniformFloatByName(name string, v float32) {

	gl.Uniform1f(prog.GetUniformLocation(name), v)
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformFloatByName(%s) error: %d", name, ecode)
		}
	}
}

// SetUniformVector2ByName sets this program uniform variable specified by
// its name to the values from the specified Vector2.
// The specified name location is cached internally.
func (prog *Program) SetUniformVector2ByName(name string, v *math32.Vector2) {

	gl.Uniform2f(prog.GetUniformLocation(name), v.X, v.Y)
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformVector2ByName(%s) error: %d", name, ecode)
		}
	}
}

// SetUniformVector3ByName sets this program uniform variable specified by
// its name to the values from the specified Vector3.
// The specified name location is cached internally.
func (prog *Program) SetUniformVector3ByName(name string, v *math32.Vector3) {

	gl.Uniform3f(prog.GetUniformLocation(name), v.X, v.Y, v.Z)
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformVector3ByName(%s) error: %d", name, ecode)
		}
	}
}

// SetUniformVector4ByName sets this program uniform variable specified by
// its name to the values from the specified Vector4.
// The specified name location is cached internally.
func (prog *Program) SetUniformVector4ByName(name string, v *math32.Vector4) {

	gl.Uniform4f(prog.GetUniformLocation(name), v.X, v.Y, v.Z, v.W)
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformVector4ByName(%s) error: %d", name, ecode)
		}
	}
}

// SetUniformMatrix3ByName sets this program uniform variable specified by
// its name with the values from the specified Matrix3.
// The specified name location is cached internally.
func (prog *Program) SetUniformMatrix3ByName(name string, m *math32.Matrix3) {

	gl.UniformMatrix3fv(prog.GetUniformLocation(name), 1, false, &m[0])
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformMatrix3ByName(%s) error: %d", name, ecode)
		}
	}
}

// SetUniformMatrix4ByName sets this program uniform variable specified by
// its name with the values from the specified Matrix4.
// The location of the name is cached internally.
func (prog *Program) SetUniformMatrix4ByName(name string, m *math32.Matrix4) {

	gl.UniformMatrix4fv(prog.GetUniformLocation(name), 1, false, &m[0])
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformMatrix4ByName(%s) error: %d", name, ecode)
		}
	}
}

// SetUniformColorByName set this program uniform variable specified by
// its name to the values from the specified Color
// The specified name location is cached internally.
func (prog *Program) SetUniformColorByName(name string, c *math32.Color) {

	gl.Uniform3f(prog.GetUniformLocation(name), c.R, c.G, c.B)
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformColorByName(%s) error: %d", name, ecode)
		}
	}
}

// SetUniformColor4ByName set this program uniform variable specified by
// its name to the values from the specified Color4
// The specified name location is cached internally.
func (prog *Program) SetUniformColor4ByName(name string, c *math32.Color4) {

	gl.Uniform4f(prog.GetUniformLocation(name), c.R, c.G, c.B, c.A)
	if prog.gs.CheckErrors() {
		ecode := gl.GetError()
		if ecode != 0 {
			log.Fatal("SetUniformColor4ByName(%s) error: %d", name, ecode)
		}
	}
}

// CompileShader creates and compiles a shader of the specified type and with
// the specified source code and returns a non-zero value by which
// it can be referenced.
func CompileShader(stype uint32, source string) (uint32, error) {

	shader := gl.CreateShader(stype)
	if shader == 0 {
		return 0, fmt.Errorf("Error creating shader")
	}

	// Allocates C string to store the source
	csource, freeSource := gl.Strs(source + "\x00")
	defer freeSource()

	// Set shader source and compile it
	gl.ShaderSource(shader, 1, csource, nil)
	gl.CompileShader(shader)

	// Get the shader compiler log
	var logLength int32
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
	slog := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(slog))

	// Get the shader compile status
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
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
