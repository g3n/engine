//
// Dynamically loads the vorbis file shared library / dll 
//
#include "loader.h"


typedef void (*alProc)(void);

//
// Windows --------------------------------------------------------------------
//
#ifdef _WIN32
#define WIN32_LEAN_AND_MEAN 1
#include <windows.h>

static HMODULE libvbf;

static int open_libvbf(void) {

	libvbf = LoadLibraryA("libvorbisfile.dll");
    if (libvbf == NULL) {
        return -1;
    }
    return 0;
}

static void close_libvbf(void) {
	FreeLibrary(libvbf);
}

static alProc get_proc(const char *proc) {
    return (alProc) GetProcAddress(libvbf, proc);
}
//
// Mac --------------------------------------------------------------------
//
#elif defined(__APPLE__) || defined(__APPLE_CC__)
#include <Carbon/Carbon.h>

CFBundleRef bundle;
CFURLRef bundleURL;

static int open_libvbf(void) {
	bundleURL = CFURLCreateWithFileSystemPath(kCFAllocatorDefault,
		CFSTR("/System/Library/Frameworks/OpenAL.framework"),
		kCFURLPOSIXPathStyle, true);
	bundle = CFBundleCreate(kCFAllocatorDefault, bundleURL);
	if (bundle == NULL) {
		return -1;
	}
	return 0;
}

static void close_libvbf(void) {
	CFRelease(bundle);
	CFRelease(bundleURL);
}

static alProc get_proc(const char *proc) {
	GL3WglProc res;
	CFStringRef procname = CFStringCreateWithCString(kCFAllocatorDefault, proc,
		kCFStringEncodingASCII);
	res = (GL3WglProc) CFBundleGetFunctionPointerForName(bundle, procname);
	CFRelease(procname);
	return res;
}
//
// Linux --------------------------------------------------------------------
//
#else
#include <dlfcn.h>

static void *libvbf;

static char* lib_names[] = {
    "libvorbisfile.so.3",
    "libvorbisfile.so",
    NULL
};

static int open_libvbf(void) {

    int i = 0;
    while (lib_names[i] != NULL) {
	    libvbf = dlopen(lib_names[i], RTLD_LAZY | RTLD_GLOBAL);
        if (libvbf != NULL) {
            dlerror(); // clear errors
            return 0;
        }
        i++;
    }
    return -1;
}

static void close_libvbf(void) {
	dlclose(libvbf);
}

static alProc get_proc(const char *proc) {
    return dlsym(libvbf, proc);
}
#endif

// Prototypes of local functions
static void load_procs(void);


// Pointers to functions loaded from shared library
LPOVCLEAR           p_ov_clear;
LPOVFOPEN           p_ov_fopen;
LPOVOPEN            p_ov_open;
LPOVOPENCALLBACKS   p_ov_open_callbacks;
LPOVTEST            p_ov_test;
LPOVTESTCALLBACKS   p_ov_test_callbacks;
LPOVTESTOPEN        p_ov_test_open;
LPOVBITRATE         p_ov_bitrate;
LPOVBITRATEINSTANT  p_ov_bitrate_instant;
LPOVSTREAMS         p_ov_streams;
LPOVSEEKABLE        p_ov_seekable;
LPOVSERIALNUMBER    p_ov_serialnumber;
LPOVRAWTOTAL        p_ov_raw_total;
LPOVPCMTOTAL        p_ov_pcm_total;
LPOVTIMETOTAL       p_ov_time_total;
LPOVRAWSEEK         p_ov_raw_seek;
LPOVPCMSEEK         p_ov_pcm_seek;
LPOVPCMSEEKPAGE     p_ov_pcm_seek_page;
LPOVTIMESEEK        p_ov_time_seek;
LPOVTIMESEEKPAGE    p_ov_time_seek_page;
LPOVRAWSEEKLAP      p_ov_raw_seek_lap;
LPOVPCMSEEKLAP      p_ov_pcm_seek_lap;
LPOVPCMSEEKPAGELAP  p_ov_pcm_seek_page_lap;
LPOVTIMESEEKLAP     p_ov_time_seek_lap;
LPOVTIMESEEKPAGELAP p_ov_time_seek_page_lap;
LPOVRAWTELL         p_ov_raw_tell;
LPOVPCMTELL         p_ov_pcm_tell;
LPOVTIMETELL        p_ov_time_tell;
LPOVINFO            p_ov_info;
LPOVCOMMENT         p_ov_comment;
LPOVREADFLOAT       p_ov_read_float;
LPOVREADFILTER      p_ov_read_filter;
LPOVREAD            p_ov_read;
LPOVCROSSLAP        p_ov_crosslap;
LPOVHALFRATE        p_ov_halfrate;
LPOVHALFRATEP       p_ov_halfrate_p;


// Load functions from shared library
int vorbisfile_load() {

    int res = open_libvbf();
    if (res) {
        return res;
    }
    load_procs();
    return 0;
}

// Loads function addresses and store in the pointers
static void load_procs(void) {
    p_ov_clear              = (LPOVCLEAR)get_proc("ov_clear");
    p_ov_fopen              = (LPOVFOPEN)get_proc("ov_fopen");
    p_ov_open               = (LPOVOPEN)get_proc("ov_open");
    p_ov_open_callbacks     = (LPOVOPENCALLBACKS)get_proc("ov_open_callbacks");
    p_ov_test               = (LPOVTEST)get_proc("ov_test");
    p_ov_test_callbacks     = (LPOVTESTCALLBACKS)get_proc("ov_test_callbacks");
    p_ov_test_open          = (LPOVTESTOPEN)get_proc("ov_test_open");
    p_ov_bitrate            = (LPOVBITRATE)get_proc("ov_bitrate");
    p_ov_bitrate_instant    = (LPOVBITRATEINSTANT)get_proc("ov_bitrate_instant");
    p_ov_streams            = (LPOVSTREAMS)get_proc("ov_streams");
    p_ov_seekable           = (LPOVSEEKABLE)get_proc("ov_seekable");
    p_ov_serialnumber       = (LPOVSERIALNUMBER)get_proc("ov_serialnumber");
    p_ov_raw_total          = (LPOVRAWTOTAL)get_proc("ov_raw_total");
    p_ov_pcm_total          = (LPOVPCMTOTAL)get_proc("ov_pcm_total");
    p_ov_time_total         = (LPOVTIMETOTAL)get_proc("ov_time_total");
    p_ov_raw_seek           = (LPOVRAWSEEK)get_proc("ov_raw_seek");
    p_ov_pcm_seek           = (LPOVPCMSEEK)get_proc("ov_pcm_seek");
    p_ov_pcm_seek_page      = (LPOVPCMSEEKPAGE)get_proc("ov_pcm_seek_page");
    p_ov_time_seek          = (LPOVTIMESEEK)get_proc("ov_time_seek");
    p_ov_time_seek_page     = (LPOVTIMESEEKPAGE)get_proc("ov_time_seek_page");
    p_ov_raw_seek_lap       = (LPOVRAWSEEKLAP)get_proc("ov_raw_seek_lap");
    p_ov_pcm_seek_lap       = (LPOVPCMSEEKLAP)get_proc("ov_pcm_seek_lap");
    p_ov_pcm_seek_page_lap  = (LPOVPCMSEEKPAGELAP)get_proc("ov_pcm_seek_page_lap");
    p_ov_time_seek_lap      = (LPOVTIMESEEKLAP)get_proc("ov_time_seek_lap");
    p_ov_time_seek_page_lap = (LPOVTIMESEEKPAGELAP)get_proc("ov_time_seek_page_lap");
    p_ov_raw_tell           = (LPOVRAWTELL)get_proc("ov_raw_tell");
    p_ov_pcm_tell           = (LPOVPCMTELL)get_proc("ov_pcm_tell");
    p_ov_time_tell          = (LPOVTIMETELL)get_proc("ov_time_tell");
    p_ov_info               = (LPOVINFO)get_proc("ov_info");
    p_ov_comment            = (LPOVCOMMENT)get_proc("ov_comment");
    p_ov_read_float         = (LPOVREADFLOAT)get_proc("ov_read_float");
    p_ov_read_filter        = (LPOVREADFILTER)get_proc("ov_read_filter");
    p_ov_read               = (LPOVREAD)get_proc("ov_read");
    p_ov_crosslap           = (LPOVCROSSLAP)get_proc("ov_crosslap");
    p_ov_halfrate           = (LPOVHALFRATE)get_proc("ov_halfrate");
    p_ov_halfrate_p         = (LPOVHALFRATEP)get_proc("ov_halfrate_p");
}

//
// Go code cannot directly call the vorbis file function pointers loaded dynamically
// The following C functions call the corresponding function pointers and can be
// called by Go code.
//

int ov_clear(OggVorbis_File *vf) {
    return p_ov_clear(vf);
}

int ov_fopen(const char *path,OggVorbis_File *vf) {
    return p_ov_fopen(path, vf);
}

int ov_open(FILE *f,OggVorbis_File *vf,const char *initial,long ibytes) {
    return ov_open(f, vf, initial, ibytes);
}

int ov_open_callbacks(void *datasource, OggVorbis_File *vf, const char *initial, long ibytes, ov_callbacks callbacks) {
    return p_ov_open_callbacks(datasource, vf, initial, ibytes, callbacks);
}

int ov_test(FILE *f,OggVorbis_File *vf,const char *initial,long ibytes) {
    return p_ov_test(f, vf, initial, ibytes);
}

int ov_test_callbacks(void *datasource, OggVorbis_File *vf, const char *initial, long ibytes, ov_callbacks callbacks) {
    return p_ov_test_callbacks(datasource, vf, initial, ibytes, callbacks);
}

int ov_test_open(OggVorbis_File *vf) {
    return p_ov_test_open(vf);
}

long ov_bitrate(OggVorbis_File *vf,int i) {
    return p_ov_bitrate(vf, i);
}

long ov_bitrate_instant(OggVorbis_File *vf) {
    return p_ov_bitrate_instant(vf);
}

long ov_streams(OggVorbis_File *vf) {
    return p_ov_streams(vf);
}

long ov_seekable(OggVorbis_File *vf) {
    return p_ov_seekable(vf);
}

long ov_serialnumber(OggVorbis_File *vf,int i) {
    return p_ov_serialnumber(vf, i);
}

ogg_int64_t ov_raw_total(OggVorbis_File *vf,int i) {
    return p_ov_raw_total(vf, i);
}

ogg_int64_t ov_pcm_total(OggVorbis_File *vf,int i) {
    return p_ov_pcm_total(vf, i);
}

double ov_time_total(OggVorbis_File *vf,int i) {
    return p_ov_time_total(vf, i);
}

int ov_raw_seek(OggVorbis_File *vf,ogg_int64_t pos) {
    return p_ov_raw_seek(vf, pos);
}

int ov_pcm_seek(OggVorbis_File *vf,ogg_int64_t pos) {
    return p_ov_pcm_seek(vf, pos);
}

int ov_pcm_seek_page(OggVorbis_File *vf,ogg_int64_t pos) {
    return p_ov_pcm_seek_page(vf, pos);
}

int ov_time_seek(OggVorbis_File *vf,double pos) {
    return p_ov_time_seek(vf, pos);
}

int ov_time_seek_page(OggVorbis_File *vf,double pos) {
    return p_ov_time_seek(vf, pos);
}

int ov_raw_seek_lap(OggVorbis_File *vf,ogg_int64_t pos) {
    return p_ov_raw_seek_lap(vf, pos);
}

int ov_pcm_seek_lap(OggVorbis_File *vf,ogg_int64_t pos) {
    return p_ov_pcm_seek(vf, pos);
}

int ov_pcm_seek_page_lap(OggVorbis_File *vf,ogg_int64_t pos) {
    return p_ov_pcm_seek_page_lap(vf, pos);
}

int ov_time_seek_lap(OggVorbis_File *vf,double pos) {
    return p_ov_time_seek_lap(vf, pos);
}

int ov_time_seek_page_lap(OggVorbis_File *vf,double pos) {
    return p_ov_time_seek_page_lap(vf, pos);
}

ogg_int64_t ov_raw_tell(OggVorbis_File *vf) {
    return p_ov_raw_tell(vf);
}

ogg_int64_t ov_pcm_tell(OggVorbis_File *vf) {
    return p_ov_pcm_tell(vf);
}

double ov_time_tell(OggVorbis_File *vf) {
    return p_ov_time_tell(vf);
}

vorbis_info *ov_info(OggVorbis_File *vf,int link) {
    return p_ov_info(vf, link);
}

vorbis_comment *ov_comment(OggVorbis_File *vf,int link) {
    return p_ov_comment(vf, link);
}

long ov_read_float(OggVorbis_File *vf,float ***pcm_channels,int samples, int *bitstream) {
    return p_ov_read_float(vf, pcm_channels, samples, bitstream);
}

long ov_read_filter(OggVorbis_File *vf,char *buffer,int length, int bigendianp,int word,int sgned,int *bitstream,
                          void (*filter)(float **pcm,long channels,long samples,void *filter_param),void *filter_param) {
    return p_ov_read_filter(vf, buffer, length, bigendianp, word, sgned, bitstream, filter, filter_param);
}

long ov_read(OggVorbis_File *vf,char *buffer,int length, int bigendianp,int word,int sgned,int *bitstream) {
    return p_ov_read(vf, buffer, length, bigendianp, word, sgned, bitstream);
}

int ov_crosslap(OggVorbis_File *vf1,OggVorbis_File *vf2) {
    return p_ov_crosslap(vf1, vf2);
}

int ov_halfrate(OggVorbis_File *vf,int flag) {
    return p_ov_halfrate(vf, flag);
}

int ov_halfrate_p(OggVorbis_File *vf) {
    return p_ov_halfrate_p(vf);
}


