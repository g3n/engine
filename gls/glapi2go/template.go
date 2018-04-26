package main

import (
	"bytes"
	"go/format"
	"os"
	"text/template"
)

// GLHeader is the definition of an OpenGL header file, with functions and constant definitions.
type GLHeader struct {
	Defines []GLDefine
	Funcs   []GLFunc
}

// GLDefine is the definition of an OpenGL constant.
type GLDefine struct {
	Name  string
	Value string
}

// GLFunc is the definition of an OpenGL function.
type GLFunc struct {
	Ptype    string    // type of function pointer (ex: PFNCULLFACEPROC)
	Spacer   string    // spacer string for formatting
	Pname    string    // pointer name (ex: pglCullFace)
	Rtype    string    // return type (ex: void)
	Fname    string    // name of function (ex: glCullFace)
	FnameGo  string    // name of function without the "gl" prefix
	CParams  string    // list of comma C parameter types and names ex:"GLenum mode, GLint x"
	Args     string    // list of comma separated argument names ex:"x, y, z"
	GoParams string    // list of comma separated Go parameters ex:"x float32, y float32"
	Params   []GLParam // array of function parameters
}

// GLParam is the definition of an argument to an OpenGL function (GLFunc).
type GLParam struct {
	Qualif string // optional parameter qualifier (ex: const)
	CType  string // parameter C type
	Arg    string // parameter name without pointer operator
	Name   string // parameter name with possible pointer operator
}

// genFile generates file from the specified template.
func genFile(templText string, td *GLHeader, fout string, gosrc bool) error {

	// Parses the template
	tmpl := template.New("tmpl")
	tmpl, err := tmpl.Parse(templText)
	if err != nil {
		return err
	}

	// Expands template to buffer
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, &td)
	if err != nil {
		return err
	}
	text := buf.Bytes()

	// Creates output file
	f, err := os.Create(fout)
	if err != nil {
		return nil
	}

	// If requested, formats generated text as Go source
	if gosrc {
		src, err := format.Source(text)
		if err != nil {
			return err
		}
		text = src
	}

	// Writes source to output file
	_, err = f.Write(text)
	if err != nil {
		return err
	}
	return f.Close()
}

//
// Template for glapi C file
//
const templGLAPIC = `
// This file was generated automatically by "glapi2go" and contains functions to
// open the platform's OpenGL dll/shared library and to load all OpenGL function
// pointers for an specified OpenGL version described by the header file "glcorearb.h",
// from "https://www.khronos.org/registry/OpenGL/api/GL/glcorearb.h".
//
// As Go cgo cannot call directly to C pointers it also creates C function wrappers
// for all loaded OpenGL pointers.
// The code was heavily based on "https://github.com/skaslev/gl3w"

#include <stdlib.h>
#include <stdio.h>
#include "glapi.h"

//
// OpenGL function loader for Windows
//
#ifdef _WIN32
#define WIN32_LEAN_AND_MEAN 1
#include <windows.h>
#undef near
#undef far

static HMODULE libgl;

// open_libgl opens the OpenGL dll for Windows
static int open_libgl(void) {

	libgl = LoadLibraryA("opengl32.dll");
	if (libgl == NULL) {
		return -1;
	}
	return 0;
}

// close_libgl closes the OpenGL dll object for Windows
static void close_libgl(void) {

	FreeLibrary(libgl);
}

// get_proc gets the pointer for an OpenGL function for Windows
static void* get_proc(const char *proc) {

	void* res;
	res = (void*)wglGetProcAddress(proc);
	if (!res) {
		res = (void*)GetProcAddress(libgl, proc);
	}
	return res;
}

//
// OpenGL function loader for Mac OS
//
#elif defined(__APPLE__)
#include <dlfcn.h>

static void *libgl;

static int open_libgl(void) {

	libgl = dlopen("/System/Library/Frameworks/OpenGL.framework/OpenGL", RTLD_LAZY | RTLD_GLOBAL);
	if (!libgl) {
		return -1;
	}
	return 0;
}

static void close_libgl(void) {

	dlclose(libgl);
}

static void* get_proc(const char *proc) {

	void* res;
	*(void **)(&res) = dlsym(libgl, proc);
	return res;
}

//
// OpenGL function loader for Linux, Unix*
//
#else
#include <dlfcn.h>
#include <GL/glx.h>

static void *libgl;
static PFNGLXGETPROCADDRESSPROC glx_get_proc_address;

// open_libgl opens the OpenGL shared object for Linux/Freebsd
static int open_libgl(void) {

	libgl = dlopen("libGL.so.1", RTLD_LAZY | RTLD_GLOBAL);
	if (libgl == NULL) {
		return -1;
	}
	*(void **)(&glx_get_proc_address) = dlsym(libgl, "glXGetProcAddressARB");
	if (glx_get_proc_address == NULL) {
		return -1;
	}
	return 0;
}

// close_libgl closes the OpenGL shared object for Linux/Freebsd
static void close_libgl(void) {

	dlclose(libgl);
}

// get_proc gets the pointer for an OpenGL function for Linux/Freebsd
static void* get_proc(const char *proc) {

	void* res;
	res = glx_get_proc_address((const GLubyte *)proc);
	if (!res) {
		*(void **)(&res) = dlsym(libgl, proc);
	}
	return res;
}
#endif

// Internal global flag to check error from OpenGL functions
static int checkError = 1;

// Declaration of internal function for loading OpenGL function pointers
static void load_procs();

//
// glapiLoad() tries to load functions addresses from the OpenGL library
//
int glapiLoad(void) {

	int res = open_libgl();
	if (res != 0) {
		return res;
	}
	load_procs();
	close_libgl();
	return 0;
}

//
// glapiCheckError sets the state of the internal flag which determines
// if error checking must be done for OpenGL calls
//
void glapiCheckError(int check) {

	checkError = check;
}

// Internal function to abort process when error
static void panic(GLenum err, const char* fname) {

		const char *msg;
		switch(err) {
    	case GL_NO_ERROR:
    		msg = "No error";
    		break;
    	case GL_INVALID_ENUM:
    		msg = "An unacceptable value is specified for an enumerated argument";
    		break;
    	case GL_INVALID_VALUE:
    		msg = "A numeric argument is out of range";
    		break;
    	case GL_INVALID_OPERATION:
    		msg = "The specified operation is not allowed in the current state";
    		break;
    	case GL_INVALID_FRAMEBUFFER_OPERATION:
    		msg = "The framebuffer object is not complete";
    		break;
    	case GL_OUT_OF_MEMORY:
    		msg = "There is not enough memory left to execute the command";
    		break;
    	case GL_STACK_UNDERFLOW:
    		msg = "An attempt has been made to perform an operation that would cause an internal stack to underflow";
    		break;
    	case GL_STACK_OVERFLOW:
    		msg = "An attempt has been made to perform an operation that would cause an internal stack to overflow";
    		break;
    	default:
    		msg = "Unexpected error";
    		break;
    }
    printf("\nGLAPI Error: %s (%d) calling: %s\n", msg, err, fname);
    exit(1);
}


//
// Definitions of function pointers variables
//
{{- range .Funcs}}
static {{.Ptype}} {{.Spacer}} {{.Pname}};
{{- end}}

//
// load_procs loads all gl functions addresses into the pointers
//
static void load_procs() {

	{{range .Funcs}}
	{{- .Pname}} = ({{.Ptype}})get_proc("{{.Fname}}"); 
	{{end}}
}

//
// Definitions of C wrapper functions for all OpenGL loaded pointers
// which call the pointer and optionally cals glGetError() to check
// for OpenGL errors.
//
{{range .Funcs}}
{{.Rtype}} {{.Fname}} ({{.CParams}}) {

	{{if ne .Rtype "void"}}
		{{- .Rtype}} res = {{.Pname}}({{.Args}});
	{{else}}
		{{- .Pname}}({{.Args}});
	{{end -}}
	if (checkError) {
		GLenum err = pglGetError();
		if (err != GL_NO_ERROR) {
			panic(err, "{{.Fname}}");
		}
	}
	{{if ne .Rtype "void" -}}
		return res;
	{{- end}}
}
{{end}}

`

//
// Template for glapi.h file
//
const templGLAPIH = `
// This file was generated automatically by "glapi2go" and contains declarations
// of public functions from "glapli.c".

#ifndef _glapi_h_
#define _glapi_h_

#include "glcorearb.h"

// Loads the OpenGL function pointers
int glapiLoad(void);

// Set the internal flag to enable/disable OpenGL error checking
void glapiCheckError(int check);

#endif
`

//
// Template for glparam.h file
//
const templGLPARAMH = `
#ifndef _glparam_h_
#define _glparam_h_

#include "glcorearb.h"

//
// Definition of structures for passing parameters to queued OpenGL functions
//
{{range .Funcs}}
typedef struct {
{{- range .Params}}
	{{.CType}} {{.Name -}};
{{- end}}
} Param{{.Fname}};
{{end}}

#endif
`

//
// Template for consts.go file
//
const templCONSTS = `
// This file was generated automatically by "glapi2go" and contains all
// OpenGL constants specified by "#define GL_*" directives contained in
// "glcorearb.h" for an specific OpenGL version converted to Go constants.

package gls

const (
	{{range .Defines}}
	{{.Name}} = {{.Value -}}
	{{end}}
)
`
