#ifndef VB_LOADER_H
#define VB_LOADER_H

#include "vorbis/vorbisenc.h"

#if defined(_WIN32)
 #define VB_APIENTRY __cdecl
#else
 #define VB_APIENTRY
#endif


// API function pointers type definitions
typedef const char* (VB_APIENTRY *LPVORBISVERSIONSTRING)(void);


// Loader
int vorbis_load();


// Pointers to functions
extern LPVORBISVERSIONSTRING   p_vorbis_version_string;



#endif


