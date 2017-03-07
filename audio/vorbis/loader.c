//
// Dynamically loads the vorbis shared library / dll
// Currently only get the pointer to the function to get the library version
//
#include "loader.h"


typedef void (*alProc)(void);

//
// Windows --------------------------------------------------------------------
//
#ifdef _WIN32
#define WIN32_LEAN_AND_MEAN 1
#include <windows.h>

static HMODULE libvb;

static int open_libvb(void) {

	libvb = LoadLibraryA("libvorbis.dll");
    if (libvb == NULL) {
        return -1;
    }
    return 0;
}

static void close_libvb(void) {
	FreeLibrary(libvb);
}

static alProc get_proc(const char *proc) {
    return (alProc) GetProcAddress(libvb, proc);
}
//
// Mac --------------------------------------------------------------------
//
#elif defined(__APPLE__) || defined(__APPLE_CC__)



//
// Linux --------------------------------------------------------------------
//
#else
#include <dlfcn.h>

static void *libvb;

static char* lib_names[] = {
    "libvorbis.so",
    "libvorbis.so.0",
    NULL
};

static int open_libvb(void) {

    int i = 0;
    while (lib_names[i] != NULL) {
	    libvb = dlopen(lib_names[i], RTLD_LAZY | RTLD_GLOBAL);
        if (libvb != NULL) {
            dlerror(); // clear errors
            return 0;
        }
        i++;
    }
    return -1;
}

static void close_libvb(void) {
	dlclose(libvb);
}

static alProc get_proc(const char *proc) {
    return dlsym(libvb, proc);
}
#endif

// Prototypes of local functions
static void load_procs(void);


// Pointers to functions loaded from shared library
LPVORBISVERSIONSTRING   p_vorbis_version_string;


// Load functions from shared library
int vorbis_load() {

    int res = open_libvb();
    if (res) {
        return res;
    }
    load_procs();
    return 0;
}

static void load_procs(void) {
    p_vorbis_version_string = (LPVORBISVERSIONSTRING)get_proc("vorbis_version_string");
}

//
// Go code cannot directly call the vorbis file function pointers loaded dynamically
// The following C functions call the corresponding function pointers and can be
// called by Go code.
//
const char *vorbis_version_string(void) {

    return p_vorbis_version_string();
}

