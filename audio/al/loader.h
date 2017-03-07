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


int al_load();

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

