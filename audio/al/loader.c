#include "loader.h"

typedef void (*alProc)(void);

//
// Windows --------------------------------------------------------------------
//
#ifdef _WIN32
#define WIN32_LEAN_AND_MEAN 1
#include <windows.h>

static HMODULE libal;

static int open_libal(void) {

	libal = LoadLibraryA("OpenAL32.dll");
    if (libal == NULL) {
        return -1;
    }
    return 0;
}

static void close_libal(void) {
	FreeLibrary(libal);
}

static alProc get_proc(const char *proc) {
    return (alProc) GetProcAddress(libal, proc);
}
//
// Mac --------------------------------------------------------------------
//
#elif defined(__APPLE__) || defined(__APPLE_CC__)
#include <Carbon/Carbon.h>

CFBundleRef bundle;
CFURLRef bundleURL;

static void open_libal(void) {
	bundleURL = CFURLCreateWithFileSystemPath(kCFAllocatorDefault,
		CFSTR("/System/Library/Frameworks/OpenAL.framework"),
		kCFURLPOSIXPathStyle, true);
	bundle = CFBundleCreate(kCFAllocatorDefault, bundleURL);
	assert(bundle != NULL);
}

static void close_libal(void) {
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

static void *libal;

static char* lib_names[] = {
    "libopenal.so",
    "libopenal.so.1",
    NULL
};

static int open_libal(void) {

    int i = 0;
    while (lib_names[i] != NULL) {
	    libal = dlopen(lib_names[i], RTLD_LAZY | RTLD_GLOBAL);
        if (libal != NULL) {
            dlerror(); // clear errors
            return 0;
        }
        i++;
    }
    return -1;
}

static void close_libal(void) {
	dlclose(libal);
}

static alProc get_proc(const char *proc) {
    return dlsym(libal, proc);
}
#endif

// Prototypes of local functions
static void load_procs(void);
static void load_efx_procs(void);


// Pointers to functions loaded from shared library
LPALENABLE                  palEnable;
LPALDISABLE                 palDisable;
LPALISENABLED               palIsEnabled;
LPALGETSTRING               palGetString;
LPALGETBOOLEANV             palGetBooleanv;
LPALGETINTEGERV             palGetIntegerv;
LPALGETFLOATV               palGetFloatv;
LPALGETDOUBLEV              palGetDoublev;
LPALGETBOOLEAN              palGetBoolean;
LPALGETINTEGER              palGetInteger;
LPALGETFLOAT                palGetFloat;
LPALGETDOUBLE               palGetDouble;
LPALGETERROR                palGetError;
LPALISEXTENSIONPRESENT      palIsExtensionPresent;
LPALGETPROCADDRESS          palGetProcAddress;
LPALGETENUMVALUE            palGetEnumValue;
LPALLISTENERF               palListenerf;
LPALLISTENER3F              palListener3f;
LPALLISTENERFV              palListenerfv;
LPALLISTENERI               palListeneri;
LPALLISTENER3I              palListener3i;
LPALLISTENERIV              palListeneriv;
LPALGETLISTENERF            palGetListenerf;
LPALGETLISTENER3F           palGetListener3f;
LPALGETLISTENERFV           palGetListenerfv;
LPALGETLISTENERI            palGetListeneri;
LPALGETLISTENER3I           palGetListener3i;
LPALGETLISTENERIV           palGetListeneriv;
LPALGENSOURCES              palGenSources;
LPALDELETESOURCES           palDeleteSources;
LPALISSOURCE                palIsSource;
LPALSOURCEF                 palSourcef;
LPALSOURCE3F                palSource3f;
LPALSOURCEFV                palSourcefv;
LPALSOURCEI                 palSourcei;
LPALSOURCE3I                palSource3i;
LPALSOURCEIV                palSourceiv;
LPALGETSOURCEF              palGetSourcef;
LPALGETSOURCE3F             palGetSource3f;
LPALGETSOURCEFV             palGetSourcefv;
LPALGETSOURCEI              palGetSourcei;
LPALGETSOURCE3I             palGetSource3i;
LPALGETSOURCEIV             palGetSourceiv;
LPALSOURCEPLAYV             palSourcePlayv;
LPALSOURCESTOPV             palSourceStopv;
LPALSOURCEREWINDV           palSourceRewindv;
LPALSOURCEPAUSEV            palSourcePausev;
LPALSOURCEPLAY              palSourcePlay;
LPALSOURCESTOP              palSourceStop;
LPALSOURCEREWIND            palSourceRewind;
LPALSOURCEPAUSE             palSourcePause;
LPALSOURCEQUEUEBUFFERS      palSourceQueueBuffers;
LPALSOURCEUNQUEUEBUFFERS    palSourceUnqueueBuffers;
LPALGENBUFFERS              palGenBuffers;
LPALDELETEBUFFERS           palDeleteBuffers;
LPALISBUFFER                palIsBuffer;
LPALBUFFERDATA              palBufferData;
LPALBUFFERF                 palBufferf;
LPALBUFFER3F                palBuffer3f;
LPALBUFFERFV                palBufferfv;
LPALBUFFERI                 palBufferi;
LPALBUFFER3I                palBuffer3i;
LPALBUFFERIV                palBufferiv;
LPALGETBUFFERF              palGetBufferf;
LPALGETBUFFER3F             palGetBuffer3f;
LPALGETBUFFERFV             palGetBufferfv;
LPALGETBUFFERI              palGetBufferi;
LPALGETBUFFER3I             palGetBuffer3i;
LPALGETBUFFERIV             palGetBufferiv;
LPALDOPPLERFACTOR           palDopplerFactor;
LPALDOPPLERVELOCITY         palDopplerVelocity;
LPALSPEEDOFSOUND            palSpeedOfSound;
LPALDISTANCEMODEL           palDistanceModel;

LPALCCREATECONTEXT          palcCreateContext;
LPALCMAKECONTEXTCURRENT     palcMakeContextCurrent;     
LPALCPROCESSCONTEXT         palcProcessContext;
LPALCSUSPENDCONTEXT         palcSuspendContext;
LPALCDESTROYCONTEXT         palcDestroyContext;
LPALCGETCURRENTCONTEXT      palcGetCurrentContext;
LPALCGETCONTEXTSDEVICE      palcGetContextsDevice;
LPALCOPENDEVICE             palcOpenDevice;
LPALCCLOSEDEVICE            palcCloseDevice;
LPALCGETERROR               palcGetError;
LPALCISEXTENSIONPRESENT     palcIsExtensionPresent;
LPALCGETPROCADDRESS         palcGetProcAddress;
LPALCGETENUMVALUE           palcGetEnumValue;
LPALCGETSTRING              palcGetString;
LPALCGETINTEGERV            palcGetIntegerv;
LPALCCAPTUREOPENDEVICE      palcCaptureOpenDevice;
LPALCCAPTURECLOSEDEVICE     palcCaptureCloseDevice;
LPALCCAPTURESTART           palcCaptureStart;
LPALCCAPTURESTOP            palcCaptureStop;
LPALCCAPTURESAMPLES         palcCaptureSamples;

// Pointers to EFX extension functions
LPALGENEFFECTS                   palGenEffects;
LPALDELETEEFFECTS                palDeleteEffects;
LPALISEFFECT                     palIsEffect;
LPALEFFECTI                      palEffecti;
LPALEFFECTIV                     palEffectiv;
LPALEFFECTF                      palEffectf;
LPALEFFECTFV                     palEffectfv;
LPALGETEFFECTI                   palGetEffecti;
LPALGETEFFECTIV                  palGetEffectiv;
LPALGETEFFECTF                   palGetEffectf;
LPALGETEFFECTFV                  palGetEffectfv;

LPALGENFILTERS                   palGenFilters;
LPALDELETEFILTERS                palDeleteFilters;
LPALISFILTER                     palIsFilter;
LPALFILTERI                      palFilteri;
LPALFILTERIV                     palFilteriv;
LPALFILTERF                      palFilterf;
LPALFILTERFV                     palFilterfv;
LPALGETFILTERI                   palGetFilteri;
LPALGETFILTERIV                  palGetFilteriv;
LPALGETFILTERF                   palGetFilterf;
LPALGETFILTERFV                  palGetFilterfv;

LPALGENAUXILIARYEFFECTSLOTS      palGenAuxiliaryEffectsSlots;
LPALDELETEAUXILIARYEFFECTSLOTS   palDeleteAuxiliaryEffectsSlots;
LPALISAUXILIARYEFFECTSLOT        palIsAuxiliaryEffectSlot;
LPALAUXILIARYEFFECTSLOTI         palAuxiliaryEffectSloti;
LPALAUXILIARYEFFECTSLOTIV        palAuxiliaryEffectSlotiv;
LPALAUXILIARYEFFECTSLOTF         palAuxiliaryEffectSlotf;
LPALAUXILIARYEFFECTSLOTFV        palAuxiliaryEffectSlotfv;
LPALGETAUXILIARYEFFECTSLOTI      palGetAuxiliaryEffectSloti;
LPALGETAUXILIARYEFFECTSLOTIV     palGetAuxiliaryEffectSlotif;
LPALGETAUXILIARYEFFECTSLOTF      palGetAuxiliaryEffectSlotf;
LPALGETAUXILIARYEFFECTSLOTFV     palGetAuxiliaryEffectSlotfv;


int al_load() {

    int res = open_libal();
    if (res) {
        return res;
    }
    load_procs();
    load_efx_procs();
    return 0;
}

static void load_procs(void) {
    palEnable               = (LPALENABLE)get_proc("alEnable");
    palDisable              = (LPALDISABLE)get_proc("alDisable");
    palIsEnabled            = (LPALISENABLED)get_proc("alIsEnabled");
    palGetString            = (LPALGETSTRING)get_proc("alGetString");
    palGetBooleanv          = (LPALGETBOOLEANV)get_proc("alGetBooleanv");
    palGetIntegerv          = (LPALGETINTEGERV)get_proc("alGetIntegerv");
    palGetFloatv            = (LPALGETFLOATV)get_proc("alGetFloatv");
    palGetDoublev           = (LPALGETDOUBLEV)get_proc("alGetDoublev");
    palGetBoolean           = (LPALGETBOOLEAN)get_proc("alGetBoolean");
    palGetInteger           = (LPALGETINTEGER)get_proc("alGetInteger");
    palGetFloat             = (LPALGETFLOAT)get_proc("alGetFloat");
    palGetDouble            = (LPALGETDOUBLE)get_proc("alGetDouble");
    palGetError             = (LPALGETERROR)get_proc("alGetError");
    palIsExtensionPresent   = (LPALISEXTENSIONPRESENT)get_proc("alIsExtensionPresent");
    palGetProcAddress       = (LPALGETPROCADDRESS)get_proc("alGetProcAddress");
    palGetEnumValue         = (LPALGETENUMVALUE)get_proc("alGetEnumValue");
    palListenerf            = (LPALLISTENERF)get_proc("alListeners");
    palListener3f           = (LPALLISTENER3F)get_proc("alListener3f");
    palListenerfv           = (LPALLISTENERFV)get_proc("alListenerfv");
    palListeneri            = (LPALLISTENERI)get_proc("alListeneri");
    palListener3i           = (LPALLISTENER3I)get_proc("alListener3i");
    palListeneriv           = (LPALLISTENERIV)get_proc("alListeneriv");
    palGetListenerf         = (LPALGETLISTENERF)get_proc("alGetListenerf");
    palGetListener3f        = (LPALGETLISTENER3F)get_proc("alGetListener3f");
    palGetListenerfv        = (LPALGETLISTENERFV)get_proc("alGetListenerfv");
    palGetListeneri         = (LPALGETLISTENERI)get_proc("alGetListeneri");
    palGetListener3i        = (LPALGETLISTENER3I)get_proc("alGetListener3i");
    palGetListeneriv        = (LPALGETLISTENERIV)get_proc("alGetListeneriv");
    palGenSources           = (LPALGENSOURCES)get_proc("alGenSources");
    palDeleteSources        = (LPALDELETESOURCES)get_proc("alDeleteSources");
    palIsSource             = (LPALISSOURCE)get_proc("alIsSource");
    palSourcef              = (LPALSOURCEF)get_proc("alSourcef");
    palSource3f             = (LPALSOURCE3F)get_proc("alSource3f");
    palSourcefv             = (LPALSOURCEFV)get_proc("alSourcefv");
    palSourcei              = (LPALSOURCEI)get_proc("alSourcei");
    palSource3i             = (LPALSOURCE3I)get_proc("alSource3i");
    palSourceiv             = (LPALSOURCEIV)get_proc(" alSourceiv");
    palGetSourcef           = (LPALGETSOURCEF)get_proc("alGetSourcef");
    palGetSource3f          = (LPALGETSOURCE3F)get_proc("alGetSource3f");
    palGetSourcefv          = (LPALGETSOURCEFV)get_proc("alGetSourcefv");
    palGetSourcei           = (LPALGETSOURCEI)get_proc("alGetSourcei");
    palGetSource3i          = (LPALGETSOURCE3I)get_proc("alGetSource3i");
    palGetSourceiv          = (LPALGETSOURCEIV)get_proc("alGetSourceiv");
    palSourcePlayv          = (LPALSOURCEPLAYV)get_proc("alSourcePlayv");
    palSourceStopv          = (LPALSOURCESTOPV)get_proc("alSourceStopv");
    palSourceRewindv        = (LPALSOURCEREWINDV)get_proc("alSourceRewindv");
    palSourcePausev         = (LPALSOURCEPAUSEV)get_proc("alSourcePausev");
    palSourcePlay           = (LPALSOURCEPLAY)get_proc("alSourcePlay");
    palSourceStop           = (LPALSOURCESTOP)get_proc("alSourceStop");
    palSourceRewind         = (LPALSOURCEREWIND)get_proc("alSourceRewind");
    palSourcePause          = (LPALSOURCEPAUSE)get_proc("alSourcePause");
    palSourceQueueBuffers   = (LPALSOURCEQUEUEBUFFERS)get_proc("alSourceQueueBuffers");
    palSourceUnqueueBuffers = (LPALSOURCEUNQUEUEBUFFERS)get_proc("alSourceUnqueueBuffers");
    palGenBuffers           = (LPALGENBUFFERS)get_proc("alGenBuffers");
    palDeleteBuffers        = (LPALDELETEBUFFERS)get_proc("alDeleteBuffers");
    palIsBuffer             = (LPALISBUFFER)get_proc("alIsBuffer");
    palBufferData           = (LPALBUFFERDATA)get_proc("alBufferData");
    palBufferf              = (LPALBUFFERF)get_proc("alBufferf");
    palBuffer3f             = (LPALBUFFER3F)get_proc("alBuffer3f");
    palBufferfv             = (LPALBUFFERFV)get_proc("alBufferfv");
    palBufferi              = (LPALBUFFERI)get_proc("alBufferi");
    palBuffer3i             = (LPALBUFFER3I)get_proc("alBuffer3i");
    palBufferiv             = (LPALBUFFERIV)get_proc("alBufferiv");
    palGetBufferf           = (LPALGETBUFFERF)get_proc("alGetBufferf");
    palGetBuffer3f          = (LPALGETBUFFER3F)get_proc("alGetBuffer3f");
    palGetBufferfv          = (LPALGETBUFFERFV)get_proc("alGetBufferfv");
    palGetBufferi           = (LPALGETBUFFERI)get_proc("alGetBufferi");
    palGetBuffer3i          = (LPALGETBUFFER3I)get_proc("alGetBuffer3i");
    palGetBufferiv          = (LPALGETBUFFERIV)get_proc("alGetBufferiv");
    palDopplerFactor        = (LPALDOPPLERFACTOR)get_proc("alDopplerFactor");
    palDopplerVelocity      = (LPALDOPPLERVELOCITY)get_proc("alDopplerVelocity");
    palSpeedOfSound         = (LPALSPEEDOFSOUND)get_proc("alSpeedOfSound");
    palDistanceModel        = (LPALDISTANCEMODEL)get_proc("alDistanceModel");

    palcCreateContext       = (LPALCCREATECONTEXT)get_proc("alcCreateContext");
    palcMakeContextCurrent  = (LPALCMAKECONTEXTCURRENT)get_proc("alcMakeContextCurrent");     
    palcProcessContext      = (LPALCPROCESSCONTEXT)get_proc("alcProcessContext");
    palcSuspendContext      = (LPALCSUSPENDCONTEXT)get_proc("alcSuspendContext");
    palcDestroyContext      = (LPALCDESTROYCONTEXT)get_proc("alcDestroyContext");
    palcGetCurrentContext   = (LPALCGETCURRENTCONTEXT)get_proc("alcGetCurrentContext");
    palcGetContextsDevice   = (LPALCGETCONTEXTSDEVICE)get_proc("alcGetContextsDevice");
    palcOpenDevice          = (LPALCOPENDEVICE)get_proc("alcOpenDevice");
    palcCloseDevice         = (LPALCCLOSEDEVICE)get_proc("alcCloseDevice");
    palcGetError            = (LPALCGETERROR)get_proc("alcGetError");
    palcIsExtensionPresent  = (LPALCISEXTENSIONPRESENT)get_proc("alcIsExtensionPresent");
    palcGetProcAddress      = (LPALCGETPROCADDRESS)get_proc("alcGetProcAddress");
    palcGetEnumValue        = (LPALCGETENUMVALUE)get_proc("alcGetEnumValue");
    palcGetString           = (LPALCGETSTRING)get_proc("alcGetString");
    palcGetIntegerv         = (LPALCGETINTEGERV)get_proc("alcGetIntegerv");
    palcCaptureOpenDevice   = (LPALCCAPTUREOPENDEVICE)get_proc("alcCaptureOpenDevice");
    palcCaptureCloseDevice  = (LPALCCAPTURECLOSEDEVICE)get_proc("alcCaptureCloseDevice");
    palcCaptureStart        = (LPALCCAPTURESTART)get_proc("alcCaptureStart");
    palcCaptureStop         = (LPALCCAPTURESTOP)get_proc("alcCaptureStop");
    palcCaptureSamples      = (LPALCCAPTURESAMPLES)get_proc("alcCaptureSamples");
}

static void load_efx_procs(void) {

    palGenEffects       = palGetProcAddress("alGenEffects");
    palDeleteEffects    = palGetProcAddress("alDeleteEffects");
    palIsEffect         = palGetProcAddress("alIsEffect");
    palEffecti          = palGetProcAddress("alEffecti");
    palEffectiv         = palGetProcAddress("alEffectiv");
    palEffectf          = palGetProcAddress("alEffectf");
    palEffectfv         = palGetProcAddress("alEffectfv");
    palGetEffecti       = palGetProcAddress("alGetEffectiv");
    palGetEffectiv      = palGetProcAddress("alGetEffectiv");
    palGetEffectf       = palGetProcAddress("alGetEffectf");
    palGetEffectfv      = palGetProcAddress("alGetEffectfv");

    palGenFilters       = palGetProcAddress("alGenFilters");
    palDeleteFilters    = palGetProcAddress("alDeleteFilters");
    palIsFilter         = palGetProcAddress("alIsFilter");
    palFilteri          = palGetProcAddress("alFilteri");
    palFilteriv         = palGetProcAddress("alFilteriv");
    palFilterf          = palGetProcAddress("alFilterf");
    palFilterfv         = palGetProcAddress("alFilterfv");
    palGetFilteri       = palGetProcAddress("GetFilteri");
    palGetFilteriv      = palGetProcAddress("GetFilteriv");
    palGetFilterf       = palGetProcAddress("GetFilterf");
    palGetFilterfv      = palGetProcAddress("GetFilterfv");

    palGenAuxiliaryEffectsSlots     = palGetProcAddress("alGenAuxiliaryEffectSlots");
    palDeleteAuxiliaryEffectsSlots  = palGetProcAddress("alDeleteAuxiliaryEffectsSlots");
    palIsAuxiliaryEffectSlot        = palGetProcAddress("alIsAuxiliaryEffectSlot");
    palAuxiliaryEffectSloti         = palGetProcAddress("alAuxiliaryEffectSloti");
    palAuxiliaryEffectSlotiv        = palGetProcAddress("alAuxiliaryEffectSlotiv");
    palAuxiliaryEffectSlotf         = palGetProcAddress("alAuxiliaryEffectSlotf");
    palAuxiliaryEffectSlotfv        = palGetProcAddress("alAuxiliaryEffectSlotfv");
    palGetAuxiliaryEffectSloti      = palGetProcAddress("alGetAuxiliaryEffectSloti");
    palGetAuxiliaryEffectSlotif     = palGetProcAddress("alGetAuxiliaryEffectSlotif");
    palGetAuxiliaryEffectSlotf      = palGetProcAddress("alGetAuxiliaryEffectSlotf");
    palGetAuxiliaryEffectSlotfv     = palGetProcAddress("alGetAuxiliaryEffectSlotfv");
}

//
// Go code cannot call C function pointers directly
// The following C functions call the corresponding function pointers and can be
// called by Go code.
//

//
// alc.h
//

ALC_API ALCcontext* ALC_APIENTRY alcCreateContext(ALCdevice *device, const ALCint* attrlist) {
    return palcCreateContext(device, attrlist);
}

ALC_API ALCboolean  ALC_APIENTRY alcMakeContextCurrent(ALCcontext *context) {
    return palcMakeContextCurrent(context);
}

ALC_API void ALC_APIENTRY alcProcessContext(ALCcontext *context) {
        palcProcessContext(context);
}

ALC_API void ALC_APIENTRY alcSuspendContext(ALCcontext *context) {
    palcSuspendContext(context);
}

ALC_API void ALC_APIENTRY alcDestroyContext(ALCcontext *context) {
    palcDestroyContext(context);
}

ALC_API ALCcontext* ALC_APIENTRY alcGetCurrentContext(void) {
    return palcGetCurrentContext();
}

ALC_API ALCdevice* ALC_APIENTRY alcGetContextsDevice(ALCcontext *context) {
    return palcGetContextsDevice(context);
}

ALC_API ALCdevice* ALC_APIENTRY alcOpenDevice(const ALCchar *devicename) {
    return palcOpenDevice(devicename);
}

ALC_API ALCboolean ALC_APIENTRY alcCloseDevice(ALCdevice *device) {
    return palcCloseDevice(device);
}

ALC_API ALCenum ALC_APIENTRY alcGetError(ALCdevice *device) {
    return palcGetError(device);
}

ALC_API ALCboolean ALC_APIENTRY alcIsExtensionPresent(ALCdevice *device, const ALCchar *extname) {
    return palcIsExtensionPresent(device, extname);
}

ALC_API void* ALC_APIENTRY alcGetProcAddress(ALCdevice *device, const ALCchar *funcname) {
    return palcGetProcAddress(device, funcname);
}

ALC_API ALCenum ALC_APIENTRY alcGetEnumValue(ALCdevice *device, const ALCchar *enumname) {
    return palcGetEnumValue(device, enumname);
}

ALC_API const ALCchar* ALC_APIENTRY alcGetString(ALCdevice *device, ALCenum param) {
    return palcGetString(device, param);
}

ALC_API void ALC_APIENTRY alcGetIntegerv(ALCdevice *device, ALCenum param, ALCsizei size, ALCint *values) {
    return palcGetIntegerv(device, param, size, values);
}

ALC_API ALCdevice* ALC_APIENTRY alcCaptureOpenDevice(const ALCchar *devicename, ALCuint frequency, ALCenum format, ALCsizei buffersize) {
    return palcCaptureOpenDevice(devicename, frequency, format, buffersize);
}

ALC_API ALCboolean ALC_APIENTRY alcCaptureCloseDevice(ALCdevice *device) {
    return palcCaptureCloseDevice(device);
}

ALC_API void ALC_APIENTRY alcCaptureStart(ALCdevice *device) {
    palcCaptureStart(device);
}

ALC_API void ALC_APIENTRY alcCaptureStop(ALCdevice *device) {
    palcCaptureStop(device);
}

ALC_API void ALC_APIENTRY alcCaptureSamples(ALCdevice *device, ALCvoid *buffer, ALCsizei samples) {
    palcCaptureSamples(device, buffer, samples);
}

//
// al.h
//

AL_API void AL_APIENTRY alEnable(ALenum capability) {
    palEnable(capability);
}

AL_API void AL_APIENTRY alDisable(ALenum capability) {
    palDisable(capability);
}

AL_API ALboolean AL_APIENTRY alIsEnabled(ALenum capability) {
    return palIsEnabled(capability);
}

AL_API const ALchar* AL_APIENTRY alGetString(ALenum param) {
    return palGetString(param);
}

AL_API void AL_APIENTRY alGetBooleanv(ALenum param, ALboolean *values) {
    palGetBooleanv(param, values);
}

AL_API void AL_APIENTRY alGetIntegerv(ALenum param, ALint *values) {
    palGetIntegerv(param, values);
}

AL_API void AL_APIENTRY alGetFloatv(ALenum param, ALfloat *values) {
    palGetFloatv(param, values);
}

AL_API void AL_APIENTRY alGetDoublev(ALenum param, ALdouble *values) {
    palGetDoublev(param, values);
}

AL_API ALboolean AL_APIENTRY alGetBoolean(ALenum param) {
    return palGetBoolean(param);
}

AL_API ALint AL_APIENTRY alGetInteger(ALenum param) {
    return palGetInteger(param);
}

AL_API ALfloat AL_APIENTRY alGetFloat(ALenum param) {
    return palGetFloat(param);
}

AL_API ALdouble AL_APIENTRY alGetDouble(ALenum param) {
    return palGetDouble(param);
}

AL_API ALenum AL_APIENTRY alGetError(void) {
    return palGetError();
}

AL_API ALboolean AL_APIENTRY alIsExtensionPresent(const ALchar *extname) {
    return palIsExtensionPresent(extname);
}

AL_API void* AL_APIENTRY alGetProcAddress(const ALchar *fname) {
    return palGetProcAddress(fname);
}

AL_API ALenum AL_APIENTRY alGetEnumValue(const ALchar *ename) {
    return palGetEnumValue(ename);
}

AL_API void AL_APIENTRY alListenerf(ALenum param, ALfloat value) {
    palListenerf(param, value);
}

AL_API void AL_APIENTRY alListener3f(ALenum param, ALfloat value1, ALfloat value2, ALfloat value3) {
    palListener3f(param, value1, value2, value3);
}

AL_API void AL_APIENTRY alListenerfv(ALenum param, const ALfloat *values) {
    palListenerfv(param, values);
}

AL_API void AL_APIENTRY alListeneri(ALenum param, ALint value) {
    palListeneri(param, value);
}

AL_API void AL_APIENTRY alListener3i(ALenum param, ALint value1, ALint value2, ALint value3) {
    palListener3i(param, value1, value2, value3);
}

AL_API void AL_APIENTRY alListeneriv(ALenum param, const ALint *values) {
    palListeneriv(param, values);
}

AL_API void AL_APIENTRY alGetListenerf(ALenum param, ALfloat *value) {
    palGetListenerf(param, value);
}

AL_API void AL_APIENTRY alGetListener3f(ALenum param, ALfloat *value1, ALfloat *value2, ALfloat *value3) {
    palGetListener3f(param, value1, value2, value3);
}

AL_API void AL_APIENTRY alGetListenerfv(ALenum param, ALfloat *values) {
    palGetListenerfv(param, values);
}

AL_API void AL_APIENTRY alGetListeneri(ALenum param, ALint *value) {
    palGetListeneri(param, value);
}

AL_API void AL_APIENTRY alGetListener3i(ALenum param, ALint *value1, ALint *value2, ALint *value3) {
    palGetListener3i(param, value1, value2, value3);
}

AL_API void AL_APIENTRY alGetListeneriv(ALenum param, ALint *values) {
    palGetListeneriv(param, values);
}

AL_API void AL_APIENTRY alGenSources(ALsizei n, ALuint *sources) {
    palGenSources(n, sources);
}

AL_API void AL_APIENTRY alDeleteSources(ALsizei n, const ALuint *sources) {
    palDeleteSources(n, sources);
}

AL_API ALboolean AL_APIENTRY alIsSource(ALuint source) {
    return palIsSource(source);
}

AL_API void AL_APIENTRY alSourcef(ALuint source, ALenum param, ALfloat value) {
    palSourcef(source, param, value);
}

AL_API void AL_APIENTRY alSource3f(ALuint source, ALenum param, ALfloat value1, ALfloat value2, ALfloat value3) {
    palSource3f(source, param, value1, value2, value3);
}

AL_API void AL_APIENTRY alSourcefv(ALuint source, ALenum param, const ALfloat *values) {
    palSourcefv(source, param, values);
}

AL_API void AL_APIENTRY alSourcei(ALuint source, ALenum param, ALint value) {
    palSourcei(source, param, value);
}

AL_API void AL_APIENTRY alSource3i(ALuint source, ALenum param, ALint value1, ALint value2, ALint value3) {
    palSource3i(source, param, value1, value2, value3);
}

AL_API void AL_APIENTRY alSourceiv(ALuint source, ALenum param, const ALint *values) {
    palSourceiv(source, param, values);
}

AL_API void AL_APIENTRY alGetSourcef(ALuint source, ALenum param, ALfloat *value) {
    palGetSourcef(source, param, value);
}

AL_API void AL_APIENTRY alGetSource3f(ALuint source, ALenum param, ALfloat *value1, ALfloat *value2, ALfloat *value3) {
    palGetSource3f(source, param, value1, value2, value3);
}

AL_API void AL_APIENTRY alGetSourcefv(ALuint source, ALenum param, ALfloat *values) {
    palGetSourcefv(source, param, values);
}

AL_API void AL_APIENTRY alGetSourcei(ALuint source, ALenum param, ALint *value) {
    palGetSourcei(source, param, value);
}

AL_API void AL_APIENTRY alGetSource3i(ALuint source, ALenum param, ALint *value1, ALint *value2, ALint *value3) {
    palGetSource3i(source, param, value1, value2, value3);
}

AL_API void AL_APIENTRY alGetSourceiv(ALuint source, ALenum param, ALint *values) {
    palGetSourceiv(source, param, values);
}

AL_API void AL_APIENTRY alSourcePlayv(ALsizei n, const ALuint *sources) {
    palSourcePlayv(n, sources);
}

AL_API void AL_APIENTRY alSourceStopv(ALsizei n, const ALuint *sources) {
    palSourceStopv(n, sources);
}

AL_API void AL_APIENTRY alSourceRewindv(ALsizei n, const ALuint *sources) {
    palSourceRewindv(n, sources);
}

AL_API void AL_APIENTRY alSourcePausev(ALsizei n, const ALuint *sources) {
    palSourcePausev(n, sources);
}

AL_API void AL_APIENTRY alSourcePlay(ALuint source) {
    palSourcePlay(source);
}

AL_API void AL_APIENTRY alSourceStop(ALuint source) {
    palSourceStop(source);
}

AL_API void AL_APIENTRY alSourceRewind(ALuint source) {
    palSourceRewind(source);
}

AL_API void AL_APIENTRY alSourcePause(ALuint source) {
    palSourcePause(source);
}

AL_API void AL_APIENTRY alSourceQueueBuffers(ALuint source, ALsizei nb, const ALuint *buffers) {
    palSourceQueueBuffers(source, nb, buffers);
}

AL_API void AL_APIENTRY alSourceUnqueueBuffers(ALuint source, ALsizei nb, ALuint *buffers) {
    palSourceUnqueueBuffers(source, nb, buffers);
}

AL_API void AL_APIENTRY alGenBuffers(ALsizei n, ALuint *buffers) {
    palGenBuffers(n, buffers);
}

AL_API void AL_APIENTRY alDeleteBuffers(ALsizei n, const ALuint *buffers) {
    palDeleteBuffers(n, buffers);
}

AL_API ALboolean AL_APIENTRY alIsBuffer(ALuint buffer) {
    return palIsBuffer(buffer);
}

AL_API void AL_APIENTRY alBufferData(ALuint buffer, ALenum format, const ALvoid *data, ALsizei size, ALsizei freq) {
    palBufferData(buffer, format, data, size, freq);
}

AL_API void AL_APIENTRY alBufferf(ALuint buffer, ALenum param, ALfloat value) {
    palBufferf(buffer, param, value);
}

AL_API void AL_APIENTRY alBuffer3f(ALuint buffer, ALenum param, ALfloat value1, ALfloat value2, ALfloat value3) {
    palBuffer3f(buffer, param, value1, value2, value3);
}

AL_API void AL_APIENTRY alBufferfv(ALuint buffer, ALenum param, const ALfloat *values) {
    palBufferfv(buffer, param, values);
}

AL_API void AL_APIENTRY alBufferi(ALuint buffer, ALenum param, ALint value) {
    palBufferi(buffer, param, value);
}

AL_API void AL_APIENTRY alBuffer3i(ALuint buffer, ALenum param, ALint value1, ALint value2, ALint value3) {
    palBuffer3i(buffer, param, value1, value2, value3);
}

AL_API void AL_APIENTRY alBufferiv(ALuint buffer, ALenum param, const ALint *values) {
    palBufferiv(buffer, param, values);
}

AL_API void AL_APIENTRY alGetBufferf(ALuint buffer, ALenum param, ALfloat *value) {
    palGetBufferf(buffer, param, value);
}

AL_API void AL_APIENTRY alGetBuffer3f(ALuint buffer, ALenum param, ALfloat *value1, ALfloat *value2, ALfloat *value3) {
    palGetBuffer3f(buffer, param, value1, value2, value3);
}

AL_API void AL_APIENTRY alGetBufferfv(ALuint buffer, ALenum param, ALfloat *values) {
    palGetBufferfv(buffer, param, values);
}

AL_API void AL_APIENTRY alGetBufferi(ALuint buffer, ALenum param, ALint *value) {
    palGetBufferi(buffer, param, value);
}

AL_API void AL_APIENTRY alGetBuffer3i(ALuint buffer, ALenum param, ALint *value1, ALint *value2, ALint *value3) {
    palGetBuffer3i(buffer, param, value1, value2, value3);
}

AL_API void AL_APIENTRY alGetBufferiv(ALuint buffer, ALenum param, ALint *values) {
    palGetBufferiv(buffer, param, values);
}

