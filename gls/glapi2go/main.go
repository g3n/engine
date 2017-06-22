package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// Current Version
const (
	PROGNAME = "glapi2go"
	VMAJOR   = 0
	VMINOR   = 1
)

// Command line options
var (
	oGLVersion = flag.String("glversion", "GL_VERSION_3_3", "OpenGL version to use")
)

const (
	fileGLAPIC   = "glapi.c"
	fileGLAPIH   = "glapi.h"
	fileGLPARAMH = "glparam.h"
	fileCONSTS   = "consts.go"
)

// Maps OpenGL types to Go
var mapCType2Go = map[string]string{
	"GLenum":     "uint",
	"GLfloat":    "float32",
	"GLchar":     "byte",
	"GLbyte":     "byte",
	"GLboolean":  "bool",
	"GLshort":    "int16",
	"GLushort":   "uint16",
	"GLint":      "int",
	"GLint64":    "int64",
	"GLsizei":    "int",
	"GLbitfield": "uint",
	"GLdouble":   "float64",
	"GLuint":     "uint",
	"GLuint64":   "uint64",
	"GLubyte":    "byte",
	"GLintptr":   "uintptr",
	"GLsizeiptr": "uintptr",
	"GLsync":     "unsafe.Pointer",
}

func main() {

	// Parse command line parameters
	flag.Usage = usage
	flag.Parse()

	// Checks for input header file
	if len(flag.Args()) == 0 {
		usage()
		return
	}
	fname := flag.Args()[0]

	// Open input header file
	fin, err := os.Open(fname)
	if err != nil {
		abort(err)
	}

	// Parses the header and builds GLHeader struct
	// with all the information necessary to expand all templates.
	var glh GLHeader
	err = parser(fin, &glh)
	if err != nil {
		abort(err)
	}

	// Generates glapi.c
	err = genFile(templGLAPIC, &glh, fileGLAPIC, false)
	if err != nil {
		abort(err)
	}

	// Generates glapi.h
	err = genFile(templGLAPIH, &glh, fileGLAPIH, false)
	if err != nil {
		abort(err)
	}

	// Generates consts.go
	err = genFile(templCONSTS, &glh, fileCONSTS, true)
	if err != nil {
		abort(err)
	}
}

// parser parses the header file and builds the Template structure
func parser(fheader io.Reader, h *GLHeader) error {

	// Regex to parser #endif line to detect end of definitions for
	// specific OpenGL version: ex:"#endif /* GL_VERSION_3_3 */"
	rexEndif := regexp.MustCompile(`#endif\s+/\*\s+(\w+)\s+\*/`)

	// Regex to parse define line, capturing name (1) and value (2)
	rexDefine := regexp.MustCompile(`#define\s+(\w+)\s+(\w+)`)

	// Regex to parse function definition line,
	// capturing return value (1), function name (2) and parameters (3)
	rexApi := regexp.MustCompile(`GLAPI\s+(.*)APIENTRY\s+(\w+)\s+\((.*)\)`)

	h.Defines = make([]GLDefine, 0)
	h.Funcs = make([]GLFunc, 0)
	bufin := bufio.NewReader(fheader)
	maxLength := 0
	for {
		// Reads next line and abort on error (not EOF)
		line, err := bufin.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		// Checks for "#endif" identifying end of definitions for specified
		// OpenGL version
		res := rexEndif.FindStringSubmatch(line)
		if len(res) > 0 {
			if res[1] == *oGLVersion {
				break
			}
		}

		// Checks for "#define" of GL constants
		res = rexDefine.FindStringSubmatch(line)
		if len(res) >= 3 {
			dname := res[1]
			if strings.HasPrefix(dname, "GL_") {
				h.Defines = append(h.Defines, GLDefine{
					Name:  gldef2go(res[1]),
					Value: glval2go(res[2]),
				})
			}
		}

		// Checks for function declaration
		res = rexApi.FindStringSubmatch(line)
		if len(res) >= 2 {
			var f GLFunc
			f.Rtype = strings.Trim(res[1], " ")
			f.Ptype = "PFN" + strings.ToUpper(res[2]) + "PROC"
			f.Fname = res[2]
			f.FnameGo = glfname2go(res[2])
			f.Pname = "p" + f.Fname
			f.CParams = res[3]
			err := parseParams(res[3], &f)
			if err != nil {
				return err
			}
			h.Funcs = append(h.Funcs, f)
			if len(f.Ptype) > maxLength {
				maxLength = len(f.Ptype)
			}
		}
		// If EOF ends of parsing.
		if err == io.EOF {
			break
		}
	}
	// Sets spacer string
	for i := 0; i < len(h.Funcs); i++ {
		h.Funcs[i].Spacer = strings.Repeat(" ", maxLength-len(h.Funcs[i].Ptype)+1)
	}

	return nil
}

// parseParams receives a string with the declaration of the parameters of a C function
// and parses it into an array of GLParam types with are then saved in the specified
// GLfunc object.
func parseParams(gparams string, f *GLFunc) error {

	params := strings.Split(gparams, ",")
	res := make([]GLParam, 0)
	args := make([]string, 0)
	goParams := make([]string, 0)
	for _, tn := range params {
		parts := strings.Split(strings.TrimSpace(tn), " ")
		var qualif string
		var name string
		var ctype string
		switch len(parts) {
		case 1:
			ctype = parts[0]
			if ctype != "void" {
				panic("Should be void but is:" + ctype)
			}
			continue
		case 2:
			ctype = parts[0]
			name = parts[1]
		case 3:
			qualif = parts[0]
			ctype = parts[1]
			name = parts[2]
		default:
			return fmt.Errorf("Invalid parameter:[%s]", tn)
		}
		arg := getArgName(name)
		args = append(args, arg)
		res = append(res, GLParam{Qualif: qualif, CType: ctype, Arg: arg, Name: name})
		// Go parameter
		goarg, gotype := gltypearg2go(ctype, name)
		goParams = append(goParams, goarg+" "+gotype)
	}
	f.Args = strings.Join(args, ", ")
	f.Params = res
	f.GoParams = strings.Join(goParams, ", ")
	return nil
}

// getArgName remove qualifiers and array brackets from the argument
// returning only the argument name. Ex: *const*indices -> indices
func getArgName(arg string) string {

	if strings.HasPrefix(arg, "*const*") {
		return strings.TrimPrefix(arg, "*const*")
	}
	if strings.HasPrefix(arg, "**") {
		return strings.TrimPrefix(arg, "**")
	}
	if strings.HasPrefix(arg, "*") {
		return strings.TrimPrefix(arg, "*")
	}
	// Checks for array index: [?]
	aidx := strings.Index(arg, "[")
	if aidx > 0 {
		return arg[:aidx]
	}
	return arg
}

// glfname2go converts the name of an OpenGL C function to Go
func glfname2go(glfname string) string {

	if strings.HasPrefix(glfname, "gl") {
		return strings.TrimPrefix(glfname, "gl")
	}
	return glfname
}

// gldef2go converts a name such as GL_LINE_LOOP to LINE_LOOP
func gldef2go(gldef string) string {

	return strings.TrimPrefix(gldef, "GL_")
}

// glval2go converts a C OpenGL value to a Go value
func glval2go(glval string) string {

	val := glval
	if strings.HasSuffix(val, "u") {
		val = strings.TrimSuffix(val, "u")
	}
	if strings.HasSuffix(val, "ull") {
		val = strings.TrimSuffix(val, "ull")
	}
	return val
}

// gltypearg2go converts a C OpenGL function type/argument to a Go argument/type
// GLfloat *param -> param *float32
// GLuint type    -> ptype uint
// void *pixels   -> pixels unsafe.Pointer
// void **params  -> params *unsafe.Pointer
func gltypearg2go(gltype, glarg string) (goarg string, gotype string) {

	// Replace parameter names using Go keywords
	gokeys := []string{"type", "func"}
	for _, k := range gokeys {
		if strings.HasSuffix(glarg, k) {
			glarg = strings.TrimSuffix(glarg, k) + "p" + k
			break
		}
	}

	if gltype == "void" {
		gotype = "unsafe.Pointer"
		if strings.HasPrefix(glarg, "**") {
			goarg = strings.TrimPrefix(glarg, "**")
			gotype = "*" + gotype
			return goarg, gotype
		}
		if strings.HasPrefix(glarg, "*") {
			goarg = strings.TrimPrefix(glarg, "*")
			return goarg, gotype
		}
		return "???", "???"
	}

	goarg = glarg
	gotype = mapCType2Go[gltype]
	if strings.HasPrefix(glarg, "*") {
		gotype = "*" + gotype
		goarg = strings.TrimPrefix(goarg, "*")
	}
	return goarg, gotype
}

//
// Shows application usage
//
func usage() {

	fmt.Fprintf(os.Stderr, "%s v%d.%d\n", PROGNAME, VMAJOR, VMINOR)
	fmt.Fprintf(os.Stderr, "usage:%s [options] <glheader>\n", strings.ToLower(PROGNAME))
	flag.PrintDefaults()
	os.Exit(2)
}

func abort(err error) {

	fmt.Fprintf(os.Stderr, "%s\n", err)
	os.Exit(1)
}
