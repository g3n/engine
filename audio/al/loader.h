#ifndef LOADER_H
#define LOADER_H

#ifdef _WIN32
#include <stdlib.h>
#include <stdio.h>
#include "AL/al.h"
#include "AL/alc.h"
#include "AL/efx.h"
#elif defined(__APPLE__) || defined(__APPLE_CC__)
#include <stdlib.h>
#include "AL/al.h"
#include "AL/alc.h"
#include "AL/efx.h"
#else
#include <stdlib.h>
#include <stdio.h>
#include "AL/al.h"
#include "AL/alc.h"
#include "AL/efx.h"
#endif

// Function declarations
int al_load();

ALCcontext* _alcCreateContext(ALCdevice *device, const ALCint* attrlist);
ALCboolean _alcMakeContextCurrent(ALCcontext *context);
void _alcProcessContext(ALCcontext *context);
void _alcSuspendContext(ALCcontext *context);
void _alcDestroyContext(ALCcontext *context);
ALCcontext* _alcGetCurrentContext(void);
ALCdevice* _alcGetContextsDevice(ALCcontext *context);
ALCdevice* _alcOpenDevice(const ALCchar *devicename);
ALCboolean _alcCloseDevice(ALCdevice *device);
ALCenum _alcGetError(ALCdevice *device);
ALCboolean _alcIsExtensionPresent(ALCdevice *device, const ALCchar *extname);
void* _alcGetProcAddress(ALCdevice *device, const ALCchar *funcname);
ALCenum _alcGetEnumValue(ALCdevice *device, const ALCchar *enumname);
const ALCchar* _alcGetString(ALCdevice *device, ALCenum param);
void _alcGetIntegerv(ALCdevice *device, ALCenum param, ALCsizei size, ALCint *values);
ALCdevice* _alcCaptureOpenDevice(const ALCchar *devicename, ALCuint frequency, ALCenum format, ALCsizei buffersize);
ALCboolean _alcCaptureCloseDevice(ALCdevice *device);
void _alcCaptureStart(ALCdevice *device);
void _alcCaptureStop(ALCdevice *device);
void _alcCaptureSamples(ALCdevice *device, ALCvoid *buffer, ALCsizei samples);

void _alEnable(ALenum capability);
void _alDisable(ALenum capability);
ALboolean _alIsEnabled(ALenum capability);
const ALchar* _alGetString(ALenum param);
void _alGetBooleanv(ALenum param, ALboolean *values);
void _alGetIntegerv(ALenum param, ALint *values);
void _alGetFloatv(ALenum param, ALfloat *values);
void _alGetDoublev(ALenum param, ALdouble *values);
ALboolean _alGetBoolean(ALenum param);
ALint _alGetInteger(ALenum param);
ALfloat _alGetFloat(ALenum param);
ALdouble _alGetDouble(ALenum param);
ALenum _alGetError(void);
ALboolean _alIsExtensionPresent(const ALchar *extname);
void* _alGetProcAddress(const ALchar *fname);
ALenum _alGetEnumValue(const ALchar *ename);
void _alListenerf(ALenum param, ALfloat value);
void _alListener3f(ALenum param, ALfloat value1, ALfloat value2, ALfloat value3);
void _alListenerfv(ALenum param, const ALfloat *values);
void _alListeneri(ALenum param, ALint value);
void _alListener3i(ALenum param, ALint value1, ALint value2, ALint value3);
void _alListeneriv(ALenum param, const ALint *values);
void  _alGetListenerf(ALenum param, ALfloat *value);
void  _alGetListener3f(ALenum param, ALfloat *value1, ALfloat *value2, ALfloat *value3);
void _alGetListenerfv(ALenum param, ALfloat *values);
void _alGetListeneri(ALenum param, ALint *value);
void _alGetListener3i(ALenum param, ALint *value1, ALint *value2, ALint *value3);
void _alGetListeneriv(ALenum param, ALint *values);
void _alGenSources(ALsizei n, ALuint *sources);
void  _alDeleteSources(ALsizei n, const ALuint *sources);
ALboolean _alIsSource(ALuint source);
void _alSourcef(ALuint source, ALenum param, ALfloat value);
void _alSource3f(ALuint source, ALenum param, ALfloat value1, ALfloat value2, ALfloat value3);
void _alSourcefv(ALuint source, ALenum param, const ALfloat *values);
void _alSourcei(ALuint source, ALenum param, ALint value);
void _alSource3i(ALuint source, ALenum param, ALint value1, ALint value2, ALint value3);
void _alSourceiv(ALuint source, ALenum param, const ALint *values);
void _alGetSourcef(ALuint source, ALenum param, ALfloat *value);
void _alGetSource3f(ALuint source, ALenum param, ALfloat *value1, ALfloat *value2, ALfloat *value3);
void _alGetSourcefv(ALuint source, ALenum param, ALfloat *values);
void _alGetSourcei(ALuint source, ALenum param, ALint *value);
void _alGetSource3i(ALuint source, ALenum param, ALint *value1, ALint *value2, ALint *value3);
void _alGetSourceiv(ALuint source, ALenum param, ALint *values);
void _alSourcePlayv(ALsizei n, const ALuint *sources);
void _alSourceStopv(ALsizei n, const ALuint *sources);
void _alSourceRewindv(ALsizei n, const ALuint *sources);
void _alSourcePausev(ALsizei n, const ALuint *sources);
void _alSourcePlay(ALuint source);
void _alSourceStop(ALuint source);
void _alSourceRewind(ALuint source);
void _alSourcePause(ALuint source);
void _alSourceQueueBuffers(ALuint source, ALsizei nb, const ALuint *buffers);
void _alSourceUnqueueBuffers(ALuint source, ALsizei nb, ALuint *buffers);
void _alGenBuffers(ALsizei n, ALuint *buffers);
void _alDeleteBuffers(ALsizei n, const ALuint *buffers);
ALboolean _alIsBuffer(ALuint buffer);
void _alBufferData(ALuint buffer, ALenum format, const ALvoid *data, ALsizei size, ALsizei freq);
void _alBufferf(ALuint buffer, ALenum param, ALfloat value);
void _alBuffer3f(ALuint buffer, ALenum param, ALfloat value1, ALfloat value2, ALfloat value3);
void _alBufferfv(ALuint buffer, ALenum param, const ALfloat *values);
void _alBufferi(ALuint buffer, ALenum param, ALint value);
void _alBuffer3i(ALuint buffer, ALenum param, ALint value1, ALint value2, ALint value3);
void _alBufferiv(ALuint buffer, ALenum param, const ALint *values);
void _alGetBufferf(ALuint buffer, ALenum param, ALfloat *value);
void _alGetBuffer3f(ALuint buffer, ALenum param, ALfloat *value1, ALfloat *value2, ALfloat *value3);
void _alGetBufferfv(ALuint buffer, ALenum param, ALfloat *values);
void _alGetBufferi(ALuint buffer, ALenum param, ALint *value);
void _alGetBuffer3i(ALuint buffer, ALenum param, ALint *value1, ALint *value2, ALint *value3);
void _alGetBufferiv(ALuint buffer, ALenum param, ALint *values);

// Function pointers declarations
extern LPALENABLE                  palEnable;
extern LPALDISABLE                 palDisable;
extern LPALISENABLED               palIsEnabled;
extern LPALGETSTRING               palGetString;
extern LPALGETBOOLEANV             palGetBooleanv;
extern LPALGETINTEGERV             palGetIntegerv;
extern LPALGETFLOATV               palGetFloatv;
extern LPALGETDOUBLEV              palGetDoublev;
extern LPALGETBOOLEAN              palGetBoolean;
extern LPALGETINTEGER              palGetInteger;
extern LPALGETFLOAT                palGetFloat;
extern LPALGETDOUBLE               palGetDouble;
extern LPALGETERROR                palGetError;
extern LPALISEXTENSIONPRESENT      palIsExtensionPresent;
extern LPALGETPROCADDRESS          palGetProcAddress;
extern LPALGETENUMVALUE            palGetEnumValue;
extern LPALLISTENERF               palListenerf;
extern LPALLISTENER3F              palListener3f;
extern LPALLISTENERFV              palListenerfv;
extern LPALLISTENERI               palListeneri;
extern LPALLISTENER3I              palListener3i;
extern LPALLISTENERIV              palListeneriv;
extern LPALGETLISTENERF            palGetListenerf;
extern LPALGETLISTENER3F           palGetListener3f;
extern LPALGETLISTENERFV           palGetListenerfv;
extern LPALGETLISTENERI            palGetListeneri;
extern LPALGETLISTENER3I           palGetListener3i;
extern LPALGETLISTENERIV           palGetListeneriv;
extern LPALGENSOURCES              palGenSources;
extern LPALDELETESOURCES           palDeleteSources;
extern LPALISSOURCE                palIsSource;
extern LPALSOURCEF                 palSourcef;
extern LPALSOURCE3F                palSource3f;
extern LPALSOURCEFV                palSourcefv;
extern LPALSOURCEI                 palSourcei;
extern LPALSOURCE3I                palSource3i;
extern LPALSOURCEIV                palSourceiv;
extern LPALGETSOURCEF              palGetSourcef;
extern LPALGETSOURCE3F             palGetSource3f;
extern LPALGETSOURCEFV             palGetSourcefv;
extern LPALGETSOURCEI              palGetSourcei;
extern LPALGETSOURCE3I             palGetSource3i;
extern LPALGETSOURCEIV             palGetSourceiv;
extern LPALSOURCEPLAYV             palSourcePlayv;
extern LPALSOURCESTOPV             palSourceStopv;
extern LPALSOURCEREWINDV           palSourceRewindv;
extern LPALSOURCEPAUSEV            palSourcePausev;
extern LPALSOURCEPLAY              palSourcePlay;
extern LPALSOURCESTOP              palSourceStop;
extern LPALSOURCEREWIND            palSourceRewind;
extern LPALSOURCEPAUSE             palSourcePause;
extern LPALSOURCEQUEUEBUFFERS      palSourceQueueBuffers;
extern LPALSOURCEUNQUEUEBUFFERS    palSourceUnqueueBuffers;
extern LPALGENBUFFERS              palGenBuffers;
extern LPALDELETEBUFFERS           palDeleteBuffers;
extern LPALISBUFFER                palIsBuffer;
extern LPALBUFFERDATA              palBufferData;
extern LPALBUFFERF                 palBufferf;
extern LPALBUFFER3F                palBuffer3f;
extern LPALBUFFERFV                palBufferfv;
extern LPALBUFFERI                 palBufferi;
extern LPALBUFFER3I                palBuffer3i;
extern LPALBUFFERIV                palBufferiv;
extern LPALGETBUFFERF              palGetBufferf;
extern LPALGETBUFFER3F             palGetBuffer3f;
extern LPALGETBUFFERFV             palGetBufferfv;
extern LPALGETBUFFERI              palGetBufferi;
extern LPALGETBUFFER3I             palGetBuffer3i;
extern LPALGETBUFFERIV             palGetBufferiv;
extern LPALDOPPLERFACTOR           palDopplerFactor;
extern LPALDOPPLERVELOCITY         palDopplerVelocity;
extern LPALSPEEDOFSOUND            palSpeedOfSound;
extern LPALDISTANCEMODEL           palDistanceModel;

extern LPALCCREATECONTEXT          palcCreateContext;
extern LPALCMAKECONTEXTCURRENT     palcMakeContextCurrent;     
extern LPALCPROCESSCONTEXT         palcProcessContext;
extern LPALCSUSPENDCONTEXT         palcSuspendContext;
extern LPALCDESTROYCONTEXT         palcDestroyContext;
extern LPALCGETCURRENTCONTEXT      palcGetCurrentContext;
extern LPALCGETCONTEXTSDEVICE      palcGetContextsDevice;
extern LPALCOPENDEVICE             palcOpenDevice;
extern LPALCCLOSEDEVICE            palcCloseDevice;
extern LPALCGETERROR               palcGetError;
extern LPALCISEXTENSIONPRESENT     palcIsExtensionPresent;
extern LPALCGETPROCADDRESS         palcGetProcAddress;
extern LPALCGETENUMVALUE           palcGetEnumValue;
extern LPALCGETSTRING              palcGetString;
extern LPALCGETINTEGERV            palcGetIntegerv;
extern LPALCCAPTUREOPENDEVICE      palcCaptureOpenDevice;
extern LPALCCAPTURECLOSEDEVICE     palcCaptureCloseDevice;
extern LPALCCAPTURESTART           palcCaptureStart;
extern LPALCCAPTURESTOP            palcCaptureStop;
extern LPALCCAPTURESAMPLES         palcCaptureSamples;

// EFX extension
extern LPALGENEFFECTS              palGenEffects;
extern LPALDELETEEFFECTS           palDeleteEffects;
extern LPALISEFFECT                palIsEffect;
extern LPALEFFECTI                 palEffecti;
extern LPALEFFECTIV                palEffectiv;
extern LPALEFFECTF                 palEffectf;
extern LPALEFFECTFV                palEffectfv;
extern LPALGETEFFECTI              palGetEffecti;
extern LPALGETEFFECTIV             palGetEffectiv;
extern LPALGETEFFECTF              palGetEffectf;
extern LPALGETEFFECTFV             palGetEffectfv;

extern LPALGENFILTERS              palGenFilters;
extern LPALDELETEFILTERS           palDeleteFilters;
extern LPALISFILTER                palIsFilter;
extern LPALFILTERI                 palFilteri;
extern LPALFILTERIV                palFilteriv;
extern LPALFILTERF                 palFilterf;
extern LPALFILTERFV                palFilterfv;
extern LPALGETFILTERI              palGetFilteri;
extern LPALGETFILTERIV             palGetFilteriv;
extern LPALGETFILTERF              palGetFilterf;
extern LPALGETFILTERFV             palGetFilterfv;

extern LPALGENAUXILIARYEFFECTSLOTS      palGenAuxiliaryEffectSlos;
extern LPALDELETEAUXILIARYEFFECTSLOTS   palDeleteAuxiliaryEffectSlots;
extern LPALISAUXILIARYEFFECTSLOT        palIsAuxiliaryEffectSlot;
extern LPALAUXILIARYEFFECTSLOTI         palAuxiliaryEffectSloti;
extern LPALAUXILIARYEFFECTSLOTIV        palAuxiliaryEffectSlotiv;
extern LPALAUXILIARYEFFECTSLOTF         palAuxiliaryEffectSlotf;
extern LPALAUXILIARYEFFECTSLOTFV        palAuxiliaryEffectSlotfv;
extern LPALGETAUXILIARYEFFECTSLOTI      palGetAuxiliaryEffectSloti;
extern LPALGETAUXILIARYEFFECTSLOTIV     palGetAuxiliaryEffectSlotif;
extern LPALGETAUXILIARYEFFECTSLOTF      palGetAuxiliaryEffectSlotf;
extern LPALGETAUXILIARYEFFECTSLOTFV     palGetAuxiliaryEffectSlotfv;

#endif

