#ifndef ROUTER_ROUTER_H
#define ROUTER_ROUTER_H

#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <winnt.h>

#include <stdio.h>

#include "AL/alc.h"
#include "AL/al.h"
#include "AL/alext.h"
#include "atomic.h"
#include "rwlock.h"
#include "threads.h"


#ifndef UNUSED
#if defined(__cplusplus)
#define UNUSED(x)
#elif defined(__GNUC__)
#define UNUSED(x) UNUSED_##x __attribute__((unused))
#elif defined(__LCLINT__)
#define UNUSED(x) /*@unused@*/ x
#else
#define UNUSED(x) x
#endif
#endif

#define MAKE_ALC_VER(major, minor) (((major)<<8) | (minor))

typedef struct DriverIface {
    WCHAR Name[32];
    HMODULE Module;
    int ALCVer;

    LPALCCREATECONTEXT alcCreateContext;
    LPALCMAKECONTEXTCURRENT alcMakeContextCurrent;
    LPALCPROCESSCONTEXT alcProcessContext;
    LPALCSUSPENDCONTEXT alcSuspendContext;
    LPALCDESTROYCONTEXT alcDestroyContext;
    LPALCGETCURRENTCONTEXT alcGetCurrentContext;
    LPALCGETCONTEXTSDEVICE alcGetContextsDevice;
    LPALCOPENDEVICE alcOpenDevice;
    LPALCCLOSEDEVICE alcCloseDevice;
    LPALCGETERROR alcGetError;
    LPALCISEXTENSIONPRESENT alcIsExtensionPresent;
    LPALCGETPROCADDRESS alcGetProcAddress;
    LPALCGETENUMVALUE alcGetEnumValue;
    LPALCGETSTRING alcGetString;
    LPALCGETINTEGERV alcGetIntegerv;
    LPALCCAPTUREOPENDEVICE alcCaptureOpenDevice;
    LPALCCAPTURECLOSEDEVICE alcCaptureCloseDevice;
    LPALCCAPTURESTART alcCaptureStart;
    LPALCCAPTURESTOP alcCaptureStop;
    LPALCCAPTURESAMPLES alcCaptureSamples;

    PFNALCSETTHREADCONTEXTPROC alcSetThreadContext;
    PFNALCGETTHREADCONTEXTPROC alcGetThreadContext;

    LPALENABLE alEnable;
    LPALDISABLE alDisable;
    LPALISENABLED alIsEnabled;
    LPALGETSTRING alGetString;
    LPALGETBOOLEANV alGetBooleanv;
    LPALGETINTEGERV alGetIntegerv;
    LPALGETFLOATV alGetFloatv;
    LPALGETDOUBLEV alGetDoublev;
    LPALGETBOOLEAN alGetBoolean;
    LPALGETINTEGER alGetInteger;
    LPALGETFLOAT alGetFloat;
    LPALGETDOUBLE alGetDouble;
    LPALGETERROR alGetError;
    LPALISEXTENSIONPRESENT alIsExtensionPresent;
    LPALGETPROCADDRESS alGetProcAddress;
    LPALGETENUMVALUE alGetEnumValue;
    LPALLISTENERF alListenerf;
    LPALLISTENER3F alListener3f;
    LPALLISTENERFV alListenerfv;
    LPALLISTENERI alListeneri;
    LPALLISTENER3I alListener3i;
    LPALLISTENERIV alListeneriv;
    LPALGETLISTENERF alGetListenerf;
    LPALGETLISTENER3F alGetListener3f;
    LPALGETLISTENERFV alGetListenerfv;
    LPALGETLISTENERI alGetListeneri;
    LPALGETLISTENER3I alGetListener3i;
    LPALGETLISTENERIV alGetListeneriv;
    LPALGENSOURCES alGenSources;
    LPALDELETESOURCES alDeleteSources;
    LPALISSOURCE alIsSource;
    LPALSOURCEF alSourcef;
    LPALSOURCE3F alSource3f;
    LPALSOURCEFV alSourcefv;
    LPALSOURCEI alSourcei;
    LPALSOURCE3I alSource3i;
    LPALSOURCEIV alSourceiv;
    LPALGETSOURCEF alGetSourcef;
    LPALGETSOURCE3F alGetSource3f;
    LPALGETSOURCEFV alGetSourcefv;
    LPALGETSOURCEI alGetSourcei;
    LPALGETSOURCE3I alGetSource3i;
    LPALGETSOURCEIV alGetSourceiv;
    LPALSOURCEPLAYV alSourcePlayv;
    LPALSOURCESTOPV alSourceStopv;
    LPALSOURCEREWINDV alSourceRewindv;
    LPALSOURCEPAUSEV alSourcePausev;
    LPALSOURCEPLAY alSourcePlay;
    LPALSOURCESTOP alSourceStop;
    LPALSOURCEREWIND alSourceRewind;
    LPALSOURCEPAUSE alSourcePause;
    LPALSOURCEQUEUEBUFFERS alSourceQueueBuffers;
    LPALSOURCEUNQUEUEBUFFERS alSourceUnqueueBuffers;
    LPALGENBUFFERS alGenBuffers;
    LPALDELETEBUFFERS alDeleteBuffers;
    LPALISBUFFER alIsBuffer;
    LPALBUFFERF alBufferf;
    LPALBUFFER3F alBuffer3f;
    LPALBUFFERFV alBufferfv;
    LPALBUFFERI alBufferi;
    LPALBUFFER3I alBuffer3i;
    LPALBUFFERIV alBufferiv;
    LPALGETBUFFERF alGetBufferf;
    LPALGETBUFFER3F alGetBuffer3f;
    LPALGETBUFFERFV alGetBufferfv;
    LPALGETBUFFERI alGetBufferi;
    LPALGETBUFFER3I alGetBuffer3i;
    LPALGETBUFFERIV alGetBufferiv;
    LPALBUFFERDATA alBufferData;
    LPALDOPPLERFACTOR alDopplerFactor;
    LPALDOPPLERVELOCITY alDopplerVelocity;
    LPALSPEEDOFSOUND alSpeedOfSound;
    LPALDISTANCEMODEL alDistanceModel;
} DriverIface;

extern DriverIface *DriverList;
extern int DriverListSize;

extern altss_t ThreadCtxDriver;
extern ATOMIC(DriverIface*) CurrentCtxDriver;


typedef struct PtrIntMap {
    ALvoid **keys;
    /* Shares memory with keys. */
    ALint *values;

    ALsizei size;
    ALsizei capacity;
    RWLock lock;
} PtrIntMap;
#define PTRINTMAP_STATIC_INITIALIZE { NULL, NULL, 0, 0, RWLOCK_STATIC_INITIALIZE }

void InitPtrIntMap(PtrIntMap *map);
void ResetPtrIntMap(PtrIntMap *map);
ALenum InsertPtrIntMapEntry(PtrIntMap *map, ALvoid *key, ALint value);
ALint RemovePtrIntMapKey(PtrIntMap *map, ALvoid *key);
ALint LookupPtrIntMapKey(PtrIntMap *map, ALvoid *key);


void InitALC(void);
void ReleaseALC(void);


enum LogLevel {
    LogLevel_None  = 0,
    LogLevel_Error = 1,
    LogLevel_Warn  = 2,
    LogLevel_Trace = 3,
};
extern enum LogLevel LogLevel;
extern FILE *LogFile;

#define TRACE(...) do {                                   \
    if(LogLevel >= LogLevel_Trace)                        \
    {                                                     \
        fprintf(LogFile, "AL Router (II): " __VA_ARGS__); \
        fflush(LogFile);                                  \
    }                                                     \
} while(0)
#define WARN(...) do {                                    \
    if(LogLevel >= LogLevel_Warn)                         \
    {                                                     \
        fprintf(LogFile, "AL Router (WW): " __VA_ARGS__); \
        fflush(LogFile);                                  \
    }                                                     \
} while(0)
#define ERR(...) do {                                     \
    if(LogLevel >= LogLevel_Error)                        \
    {                                                     \
        fprintf(LogFile, "AL Router (EE): " __VA_ARGS__); \
        fflush(LogFile);                                  \
    }                                                     \
} while(0)

#endif /* ROUTER_ROUTER_H */
