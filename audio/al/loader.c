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

static int open_libal(void) {
	bundleURL = CFURLCreateWithFileSystemPath(kCFAllocatorDefault,
		CFSTR("/System/Library/Frameworks/OpenAL.framework"),
		kCFURLPOSIXPathStyle, true);
	bundle = CFBundleCreate(kCFAllocatorDefault, bundleURL);
	if (bundle == NULL) {
		return -1;
	}
	return 0;
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
ALCcontext* _alcCreateContext(ALCdevice *device, const ALCint* attrlist) {
    return palcCreateContext(device, attrlist);
}

ALCboolean _alcMakeContextCurrent(ALCcontext *context) {
    return palcMakeContextCurrent(context);
}

void _alcProcessContext(ALCcontext *context) {
    palcProcessContext(context);
}

void _alcSuspendContext(ALCcontext *context) {
    palcSuspendContext(context);
}

void _alcDestroyContext(ALCcontext *context) {
    palcDestroyContext(context);
}

ALCcontext* _alcGetCurrentContext(void) {
    return palcGetCurrentContext();
}

ALCdevice* _alcGetContextsDevice(ALCcontext *context) {
    return palcGetContextsDevice(context);
}

ALCdevice* _alcOpenDevice(const ALCchar *devicename) {
    return palcOpenDevice(devicename);
}

ALCboolean _alcCloseDevice(ALCdevice *device) {
    return palcCloseDevice(device);
}

ALCenum _alcGetError(ALCdevice *device) {
    return palcGetError(device);
}

ALCboolean _alcIsExtensionPresent(ALCdevice *device, const ALCchar *extname) {
    return palcIsExtensionPresent(device, extname);
}

void* _alcGetProcAddress(ALCdevice *device, const ALCchar *funcname) {
    return palcGetProcAddress(device, funcname);
}

ALCenum _alcGetEnumValue(ALCdevice *device, const ALCchar *enumname) {
    return palcGetEnumValue(device, enumname);
}

const ALCchar* _alcGetString(ALCdevice *device, ALCenum param) {
    return palcGetString(device, param);
}

void _alcGetIntegerv(ALCdevice *device, ALCenum param, ALCsizei size, ALCint *values) {
    return palcGetIntegerv(device, param, size, values);
}

ALCdevice* _alcCaptureOpenDevice(const ALCchar *devicename, ALCuint frequency, ALCenum format, ALCsizei buffersize) {
    return palcCaptureOpenDevice(devicename, frequency, format, buffersize);
}

ALCboolean _alcCaptureCloseDevice(ALCdevice *device) {
    return palcCaptureCloseDevice(device);
}

void _alcCaptureStart(ALCdevice *device) {
    palcCaptureStart(device);
}

void _alcCaptureStop(ALCdevice *device) {
    palcCaptureStop(device);
}

void _alcCaptureSamples(ALCdevice *device, ALCvoid *buffer, ALCsizei samples) {
    palcCaptureSamples(device, buffer, samples);
}

//
// al.h
//
void _alEnable(ALenum capability) {
    palEnable(capability);
}

void _alDisable(ALenum capability) {
    palDisable(capability);
}

ALboolean _alIsEnabled(ALenum capability) {
    return palIsEnabled(capability);
}

const ALchar* _alGetString(ALenum param) {
    return palGetString(param);
}

void _alGetBooleanv(ALenum param, ALboolean *values) {
    palGetBooleanv(param, values);
}

void _alGetIntegerv(ALenum param, ALint *values) {
    palGetIntegerv(param, values);
}

void _alGetFloatv(ALenum param, ALfloat *values) {
    palGetFloatv(param, values);
}

void _alGetDoublev(ALenum param, ALdouble *values) {
    palGetDoublev(param, values);
}

ALboolean _alGetBoolean(ALenum param) {
    return palGetBoolean(param);
}

ALint _alGetInteger(ALenum param) {
    return palGetInteger(param);
}

ALfloat _alGetFloat(ALenum param) {
    return palGetFloat(param);
}

ALdouble _alGetDouble(ALenum param) {
    return palGetDouble(param);
}

ALenum _alGetError(void) {
    return palGetError();
}

ALboolean _alIsExtensionPresent(const ALchar *extname) {
    return palIsExtensionPresent(extname);
}

void* _alGetProcAddress(const ALchar *fname) {
    return palGetProcAddress(fname);
}

ALenum _alGetEnumValue(const ALchar *ename) {
    return palGetEnumValue(ename);
}

void _alListenerf(ALenum param, ALfloat value) {
    palListenerf(param, value);
}

void _alListener3f(ALenum param, ALfloat value1, ALfloat value2, ALfloat value3) {
    palListener3f(param, value1, value2, value3);
}

void _alListenerfv(ALenum param, const ALfloat *values) {
    palListenerfv(param, values);
}

void _alListeneri(ALenum param, ALint value) {
    palListeneri(param, value);
}

void _alListener3i(ALenum param, ALint value1, ALint value2, ALint value3) {
    palListener3i(param, value1, value2, value3);
}

void _alListeneriv(ALenum param, const ALint *values) {
    palListeneriv(param, values);
}

void  _alGetListenerf(ALenum param, ALfloat *value) {
    palGetListenerf(param, value);
}

void  _alGetListener3f(ALenum param, ALfloat *value1, ALfloat *value2, ALfloat *value3) {
    palGetListener3f(param, value1, value2, value3);
}

void _alGetListenerfv(ALenum param, ALfloat *values) {
    palGetListenerfv(param, values);
}

void _alGetListeneri(ALenum param, ALint *value) {
    palGetListeneri(param, value);
}

void _alGetListener3i(ALenum param, ALint *value1, ALint *value2, ALint *value3) {
    palGetListener3i(param, value1, value2, value3);
}

void _alGetListeneriv(ALenum param, ALint *values) {
    palGetListeneriv(param, values);
}

void _alGenSources(ALsizei n, ALuint *sources) {
    palGenSources(n, sources);
}

void  _alDeleteSources(ALsizei n, const ALuint *sources) {
    palDeleteSources(n, sources);
}

ALboolean _alIsSource(ALuint source) {
    return palIsSource(source);
}

void _alSourcef(ALuint source, ALenum param, ALfloat value) {
    palSourcef(source, param, value);
}

void _alSource3f(ALuint source, ALenum param, ALfloat value1, ALfloat value2, ALfloat value3) {
    palSource3f(source, param, value1, value2, value3);
}

void _alSourcefv(ALuint source, ALenum param, const ALfloat *values) {
    palSourcefv(source, param, values);
}

void _alSourcei(ALuint source, ALenum param, ALint value) {
    palSourcei(source, param, value);
}

void _alSource3i(ALuint source, ALenum param, ALint value1, ALint value2, ALint value3) {
    palSource3i(source, param, value1, value2, value3);
}

void _alSourceiv(ALuint source, ALenum param, const ALint *values) {
    palSourceiv(source, param, values);
}

void _alGetSourcef(ALuint source, ALenum param, ALfloat *value) {
    palGetSourcef(source, param, value);
}

void _alGetSource3f(ALuint source, ALenum param, ALfloat *value1, ALfloat *value2, ALfloat *value3) {
    palGetSource3f(source, param, value1, value2, value3);
}

void _alGetSourcefv(ALuint source, ALenum param, ALfloat *values) {
    palGetSourcefv(source, param, values);
}

void _alGetSourcei(ALuint source, ALenum param, ALint *value) {
    palGetSourcei(source, param, value);
}

void _alGetSource3i(ALuint source, ALenum param, ALint *value1, ALint *value2, ALint *value3) {
    palGetSource3i(source, param, value1, value2, value3);
}

void _alGetSourceiv(ALuint source, ALenum param, ALint *values) {
    palGetSourceiv(source, param, values);
}

void _alSourcePlayv(ALsizei n, const ALuint *sources) {
    palSourcePlayv(n, sources);
}

void _alSourceStopv(ALsizei n, const ALuint *sources) {
    palSourceStopv(n, sources);
}

void _alSourceRewindv(ALsizei n, const ALuint *sources) {
    palSourceRewindv(n, sources);
}

void _alSourcePausev(ALsizei n, const ALuint *sources) {
    palSourcePausev(n, sources);
}

void _alSourcePlay(ALuint source) {
    palSourcePlay(source);
}

void _alSourceStop(ALuint source) {
    palSourceStop(source);
}

void _alSourceRewind(ALuint source) {
    palSourceRewind(source);
}

void _alSourcePause(ALuint source) {
    palSourcePause(source);
}

void _alSourceQueueBuffers(ALuint source, ALsizei nb, const ALuint *buffers) {
    palSourceQueueBuffers(source, nb, buffers);
}

void _alSourceUnqueueBuffers(ALuint source, ALsizei nb, ALuint *buffers) {
    palSourceUnqueueBuffers(source, nb, buffers);
}

void _alGenBuffers(ALsizei n, ALuint *buffers) {
    palGenBuffers(n, buffers);
}

void _alDeleteBuffers(ALsizei n, const ALuint *buffers) {
    palDeleteBuffers(n, buffers);
}

ALboolean _alIsBuffer(ALuint buffer) {
    return palIsBuffer(buffer);
}

void _alBufferData(ALuint buffer, ALenum format, const ALvoid *data, ALsizei size, ALsizei freq) {
    palBufferData(buffer, format, data, size, freq);
}

void _alBufferf(ALuint buffer, ALenum param, ALfloat value) {
    palBufferf(buffer, param, value);
}

void _alBuffer3f(ALuint buffer, ALenum param, ALfloat value1, ALfloat value2, ALfloat value3) {
    palBuffer3f(buffer, param, value1, value2, value3);
}

void _alBufferfv(ALuint buffer, ALenum param, const ALfloat *values) {
    palBufferfv(buffer, param, values);
}

void _alBufferi(ALuint buffer, ALenum param, ALint value) {
    palBufferi(buffer, param, value);
}

void _alBuffer3i(ALuint buffer, ALenum param, ALint value1, ALint value2, ALint value3) {
    palBuffer3i(buffer, param, value1, value2, value3);
}

void _alBufferiv(ALuint buffer, ALenum param, const ALint *values) {
    palBufferiv(buffer, param, values);
}

void _alGetBufferf(ALuint buffer, ALenum param, ALfloat *value) {
    palGetBufferf(buffer, param, value);
}

void _alGetBuffer3f(ALuint buffer, ALenum param, ALfloat *value1, ALfloat *value2, ALfloat *value3) {
    palGetBuffer3f(buffer, param, value1, value2, value3);
}

void _alGetBufferfv(ALuint buffer, ALenum param, ALfloat *values) {
    palGetBufferfv(buffer, param, values);
}

void _alGetBufferi(ALuint buffer, ALenum param, ALint *value) {
    palGetBufferi(buffer, param, value);
}

void _alGetBuffer3i(ALuint buffer, ALenum param, ALint *value1, ALint *value2, ALint *value3) {
    palGetBuffer3i(buffer, param, value1, value2, value3);
}

void _alGetBufferiv(ALuint buffer, ALenum param, ALint *values) {
    palGetBufferiv(buffer, param, values);
}

