
#include "config.h"

#include <stddef.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#include "AL/alc.h"
#include "router.h"
#include "almalloc.h"


#define COUNTOF(x)  (sizeof(x)/sizeof(x[0]))

#define DECL(x) { #x, (ALCvoid*)(x) }
static const struct {
    const ALCchar *funcName;
    ALCvoid *address;
} alcFunctions[] = {
    DECL(alcCreateContext),
    DECL(alcMakeContextCurrent),
    DECL(alcProcessContext),
    DECL(alcSuspendContext),
    DECL(alcDestroyContext),
    DECL(alcGetCurrentContext),
    DECL(alcGetContextsDevice),
    DECL(alcOpenDevice),
    DECL(alcCloseDevice),
    DECL(alcGetError),
    DECL(alcIsExtensionPresent),
    DECL(alcGetProcAddress),
    DECL(alcGetEnumValue),
    DECL(alcGetString),
    DECL(alcGetIntegerv),
    DECL(alcCaptureOpenDevice),
    DECL(alcCaptureCloseDevice),
    DECL(alcCaptureStart),
    DECL(alcCaptureStop),
    DECL(alcCaptureSamples),

    DECL(alcSetThreadContext),
    DECL(alcGetThreadContext),

    DECL(alEnable),
    DECL(alDisable),
    DECL(alIsEnabled),
    DECL(alGetString),
    DECL(alGetBooleanv),
    DECL(alGetIntegerv),
    DECL(alGetFloatv),
    DECL(alGetDoublev),
    DECL(alGetBoolean),
    DECL(alGetInteger),
    DECL(alGetFloat),
    DECL(alGetDouble),
    DECL(alGetError),
    DECL(alIsExtensionPresent),
    DECL(alGetProcAddress),
    DECL(alGetEnumValue),
    DECL(alListenerf),
    DECL(alListener3f),
    DECL(alListenerfv),
    DECL(alListeneri),
    DECL(alListener3i),
    DECL(alListeneriv),
    DECL(alGetListenerf),
    DECL(alGetListener3f),
    DECL(alGetListenerfv),
    DECL(alGetListeneri),
    DECL(alGetListener3i),
    DECL(alGetListeneriv),
    DECL(alGenSources),
    DECL(alDeleteSources),
    DECL(alIsSource),
    DECL(alSourcef),
    DECL(alSource3f),
    DECL(alSourcefv),
    DECL(alSourcei),
    DECL(alSource3i),
    DECL(alSourceiv),
    DECL(alGetSourcef),
    DECL(alGetSource3f),
    DECL(alGetSourcefv),
    DECL(alGetSourcei),
    DECL(alGetSource3i),
    DECL(alGetSourceiv),
    DECL(alSourcePlayv),
    DECL(alSourceStopv),
    DECL(alSourceRewindv),
    DECL(alSourcePausev),
    DECL(alSourcePlay),
    DECL(alSourceStop),
    DECL(alSourceRewind),
    DECL(alSourcePause),
    DECL(alSourceQueueBuffers),
    DECL(alSourceUnqueueBuffers),
    DECL(alGenBuffers),
    DECL(alDeleteBuffers),
    DECL(alIsBuffer),
    DECL(alBufferData),
    DECL(alBufferf),
    DECL(alBuffer3f),
    DECL(alBufferfv),
    DECL(alBufferi),
    DECL(alBuffer3i),
    DECL(alBufferiv),
    DECL(alGetBufferf),
    DECL(alGetBuffer3f),
    DECL(alGetBufferfv),
    DECL(alGetBufferi),
    DECL(alGetBuffer3i),
    DECL(alGetBufferiv),
    DECL(alDopplerFactor),
    DECL(alDopplerVelocity),
    DECL(alSpeedOfSound),
    DECL(alDistanceModel),
};
#undef DECL

#define DECL(x) { #x, (x) }
static const struct {
    const ALCchar *enumName;
    ALCenum value;
} alcEnumerations[] = {
    DECL(ALC_INVALID),
    DECL(ALC_FALSE),
    DECL(ALC_TRUE),

    DECL(ALC_MAJOR_VERSION),
    DECL(ALC_MINOR_VERSION),
    DECL(ALC_ATTRIBUTES_SIZE),
    DECL(ALC_ALL_ATTRIBUTES),
    DECL(ALC_DEFAULT_DEVICE_SPECIFIER),
    DECL(ALC_DEVICE_SPECIFIER),
    DECL(ALC_ALL_DEVICES_SPECIFIER),
    DECL(ALC_DEFAULT_ALL_DEVICES_SPECIFIER),
    DECL(ALC_EXTENSIONS),
    DECL(ALC_FREQUENCY),
    DECL(ALC_REFRESH),
    DECL(ALC_SYNC),
    DECL(ALC_MONO_SOURCES),
    DECL(ALC_STEREO_SOURCES),
    DECL(ALC_CAPTURE_DEVICE_SPECIFIER),
    DECL(ALC_CAPTURE_DEFAULT_DEVICE_SPECIFIER),
    DECL(ALC_CAPTURE_SAMPLES),

    DECL(ALC_NO_ERROR),
    DECL(ALC_INVALID_DEVICE),
    DECL(ALC_INVALID_CONTEXT),
    DECL(ALC_INVALID_ENUM),
    DECL(ALC_INVALID_VALUE),
    DECL(ALC_OUT_OF_MEMORY),

    DECL(AL_INVALID),
    DECL(AL_NONE),
    DECL(AL_FALSE),
    DECL(AL_TRUE),

    DECL(AL_SOURCE_RELATIVE),
    DECL(AL_CONE_INNER_ANGLE),
    DECL(AL_CONE_OUTER_ANGLE),
    DECL(AL_PITCH),
    DECL(AL_POSITION),
    DECL(AL_DIRECTION),
    DECL(AL_VELOCITY),
    DECL(AL_LOOPING),
    DECL(AL_BUFFER),
    DECL(AL_GAIN),
    DECL(AL_MIN_GAIN),
    DECL(AL_MAX_GAIN),
    DECL(AL_ORIENTATION),
    DECL(AL_REFERENCE_DISTANCE),
    DECL(AL_ROLLOFF_FACTOR),
    DECL(AL_CONE_OUTER_GAIN),
    DECL(AL_MAX_DISTANCE),
    DECL(AL_SEC_OFFSET),
    DECL(AL_SAMPLE_OFFSET),
    DECL(AL_BYTE_OFFSET),
    DECL(AL_SOURCE_TYPE),
    DECL(AL_STATIC),
    DECL(AL_STREAMING),
    DECL(AL_UNDETERMINED),

    DECL(AL_SOURCE_STATE),
    DECL(AL_INITIAL),
    DECL(AL_PLAYING),
    DECL(AL_PAUSED),
    DECL(AL_STOPPED),

    DECL(AL_BUFFERS_QUEUED),
    DECL(AL_BUFFERS_PROCESSED),

    DECL(AL_FORMAT_MONO8),
    DECL(AL_FORMAT_MONO16),
    DECL(AL_FORMAT_STEREO8),
    DECL(AL_FORMAT_STEREO16),

    DECL(AL_FREQUENCY),
    DECL(AL_BITS),
    DECL(AL_CHANNELS),
    DECL(AL_SIZE),

    DECL(AL_UNUSED),
    DECL(AL_PENDING),
    DECL(AL_PROCESSED),

    DECL(AL_NO_ERROR),
    DECL(AL_INVALID_NAME),
    DECL(AL_INVALID_ENUM),
    DECL(AL_INVALID_VALUE),
    DECL(AL_INVALID_OPERATION),
    DECL(AL_OUT_OF_MEMORY),

    DECL(AL_VENDOR),
    DECL(AL_VERSION),
    DECL(AL_RENDERER),
    DECL(AL_EXTENSIONS),

    DECL(AL_DOPPLER_FACTOR),
    DECL(AL_DOPPLER_VELOCITY),
    DECL(AL_DISTANCE_MODEL),
    DECL(AL_SPEED_OF_SOUND),

    DECL(AL_INVERSE_DISTANCE),
    DECL(AL_INVERSE_DISTANCE_CLAMPED),
    DECL(AL_LINEAR_DISTANCE),
    DECL(AL_LINEAR_DISTANCE_CLAMPED),
    DECL(AL_EXPONENT_DISTANCE),
    DECL(AL_EXPONENT_DISTANCE_CLAMPED),
};
#undef DECL

static const ALCchar alcNoError[] = "No Error";
static const ALCchar alcErrInvalidDevice[] = "Invalid Device";
static const ALCchar alcErrInvalidContext[] = "Invalid Context";
static const ALCchar alcErrInvalidEnum[] = "Invalid Enum";
static const ALCchar alcErrInvalidValue[] = "Invalid Value";
static const ALCchar alcErrOutOfMemory[] = "Out of Memory";
static const ALCchar alcExtensionList[] =
    "ALC_ENUMERATE_ALL_EXT ALC_ENUMERATION_EXT ALC_EXT_CAPTURE "
    "ALC_EXT_thread_local_context";

static const ALCint alcMajorVersion = 1;
static const ALCint alcMinorVersion = 1;


static almtx_t EnumerationLock;
static almtx_t ContextSwitchLock;

static ATOMIC(ALCenum) LastError = ATOMIC_INIT_STATIC(ALC_NO_ERROR);
static PtrIntMap DeviceIfaceMap = PTRINTMAP_STATIC_INITIALIZE;
static PtrIntMap ContextIfaceMap = PTRINTMAP_STATIC_INITIALIZE;


typedef struct EnumeratedList {
    ALCchar *Names;
    ALCchar *NamesEnd;
    ALCint *Indicies;
    ALCsizei IndexSize;
} EnumeratedList;
static EnumeratedList DevicesList = { NULL, NULL, NULL, 0 };
static EnumeratedList AllDevicesList = { NULL, NULL, NULL, 0 };
static EnumeratedList CaptureDevicesList = { NULL, NULL, NULL, 0 };

static void ClearDeviceList(EnumeratedList *list)
{
    al_free(list->Names);
    list->Names = NULL;
    list->NamesEnd = NULL;

    al_free(list->Indicies);
    list->Indicies = NULL;
    list->IndexSize = 0;
}

static void AppendDeviceList(EnumeratedList *list, const ALCchar *names, ALint idx)
{
    const ALCchar *name_end = names;
    ALCsizei count = 0;
    ALCchar *new_list;
    ALCint *new_indicies;
    size_t len;
    ALCsizei i;

    if(!name_end)
        return;
    while(*name_end)
    {
        TRACE("Enumerated \"%s\", driver %d\n", name_end, idx);
        count++;
        name_end += strlen(name_end)+1;
    }
    if(names == name_end)
        return;

    len = (list->NamesEnd - list->Names) + (name_end - names);
    new_list = al_calloc(DEF_ALIGN, len + 1);
    memcpy(new_list,  list->Names, list->NamesEnd - list->Names);
    memcpy(new_list + (list->NamesEnd - list->Names), names, name_end - names);
    al_free(list->Names);
    list->Names = new_list;
    list->NamesEnd = list->Names + len;

    new_indicies = al_calloc(16, sizeof(ALCint)*(list->IndexSize + count));
    for(i = 0;i < list->IndexSize;i++)
        new_indicies[i] = list->Indicies[i];
    for(i = 0;i < count;i++)
        new_indicies[list->IndexSize+i] = idx;
    al_free(list->Indicies);
    list->Indicies = new_indicies;
    list->IndexSize += count;
}

static ALint GetDriverIndexForName(const EnumeratedList *list, const ALCchar *name)
{
    const ALCchar *devnames = list->Names;
    const ALCint *index = list->Indicies;

    while(devnames && *devnames)
    {
        if(strcmp(name, devnames) == 0)
            return *index;
        devnames += strlen(devnames)+1;
        index++;
    }
    return -1;
}

void InitALC(void)
{
    almtx_init(&EnumerationLock, almtx_recursive);
    almtx_init(&ContextSwitchLock, almtx_plain);
}

void ReleaseALC(void)
{
    ClearDeviceList(&DevicesList);
    ClearDeviceList(&AllDevicesList);
    ClearDeviceList(&CaptureDevicesList);

    ResetPtrIntMap(&ContextIfaceMap);
    ResetPtrIntMap(&DeviceIfaceMap);

    almtx_destroy(&ContextSwitchLock);
    almtx_destroy(&EnumerationLock);
}


ALC_API ALCdevice* ALC_APIENTRY alcOpenDevice(const ALCchar *devicename)
{
    ALCdevice *device = NULL;
    ALint idx;

    /* Prior to the enumeration extension, apps would hardcode these names as a
     * quality hint for the wrapper driver. Ignore them since there's no sane
     * way to map them.
     */
    if(devicename && (devicename[0] == '\0' ||
                      strcmp(devicename, "DirectSound3D") == 0 ||
                      strcmp(devicename, "DirectSound") == 0 ||
                      strcmp(devicename, "MMSYSTEM") == 0))
        devicename = NULL;
    if(devicename)
    {
        almtx_lock(&EnumerationLock);
        if(!DevicesList.Names)
            (void)alcGetString(NULL, ALC_DEVICE_SPECIFIER);
        idx = GetDriverIndexForName(&DevicesList, devicename);
        if(idx < 0)
        {
            if(!AllDevicesList.Names)
                (void)alcGetString(NULL, ALC_ALL_DEVICES_SPECIFIER);
            idx = GetDriverIndexForName(&AllDevicesList, devicename);
        }
        almtx_unlock(&EnumerationLock);
        if(idx < 0)
        {
            ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_VALUE);
            TRACE("Failed to find driver for name \"%s\"\n", devicename);
            return NULL;
        }
        TRACE("Found driver %d for name \"%s\"\n", idx, devicename);
        device = DriverList[idx].alcOpenDevice(devicename);
    }
    else
    {
        int i;
        for(i = 0;i < DriverListSize;i++)
        {
            if(DriverList[i].ALCVer >= MAKE_ALC_VER(1, 1) ||
               DriverList[i].alcIsExtensionPresent(NULL, "ALC_ENUMERATION_EXT"))
            {
                idx = i;
                TRACE("Using default device from driver %d\n", idx);
                device = DriverList[idx].alcOpenDevice(NULL);
                break;
            }
        }
    }

    if(device)
    {
        if(InsertPtrIntMapEntry(&DeviceIfaceMap, device, idx) != ALC_NO_ERROR)
        {
            DriverList[idx].alcCloseDevice(device);
            device = NULL;
        }
    }

    return device;
}

ALC_API ALCboolean ALC_APIENTRY alcCloseDevice(ALCdevice *device)
{
    ALint idx;

    if(!device || (idx=LookupPtrIntMapKey(&DeviceIfaceMap, device)) < 0)
    {
        ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_DEVICE);
        return ALC_FALSE;
    }
    if(!DriverList[idx].alcCloseDevice(device))
        return ALC_FALSE;
    RemovePtrIntMapKey(&DeviceIfaceMap, device);
    return ALC_TRUE;
}


ALC_API ALCcontext* ALC_APIENTRY alcCreateContext(ALCdevice *device, const ALCint *attrlist)
{
    ALCcontext *context;
    ALint idx;

    if(!device || (idx=LookupPtrIntMapKey(&DeviceIfaceMap, device)) < 0)
    {
        ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_DEVICE);
        return ALC_FALSE;
    }
    context = DriverList[idx].alcCreateContext(device, attrlist);
    if(context)
    {
        if(InsertPtrIntMapEntry(&ContextIfaceMap, context, idx) != ALC_NO_ERROR)
        {
            DriverList[idx].alcDestroyContext(context);
            context = NULL;
        }
    }

    return context;
}

ALC_API ALCboolean ALC_APIENTRY alcMakeContextCurrent(ALCcontext *context)
{
    ALint idx = -1;

    almtx_lock(&ContextSwitchLock);
    if(context)
    {
        idx = LookupPtrIntMapKey(&ContextIfaceMap, context);
        if(idx < 0)
        {
            ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_CONTEXT);
            almtx_unlock(&ContextSwitchLock);
            return ALC_FALSE;
        }
        if(!DriverList[idx].alcMakeContextCurrent(context))
        {
            almtx_unlock(&ContextSwitchLock);
            return ALC_FALSE;
        }
    }

    /* Unset the context from the old driver if it's different from the new
     * current one.
     */
    if(idx < 0)
    {
        DriverIface *oldiface = altss_get(ThreadCtxDriver);
        if(oldiface) oldiface->alcSetThreadContext(NULL);
        oldiface = ATOMIC_EXCHANGE_PTR_SEQ(&CurrentCtxDriver, NULL);
        if(oldiface) oldiface->alcMakeContextCurrent(NULL);
    }
    else
    {
        DriverIface *oldiface = altss_get(ThreadCtxDriver);
        if(oldiface && oldiface != &DriverList[idx])
            oldiface->alcSetThreadContext(NULL);
        oldiface = ATOMIC_EXCHANGE_PTR_SEQ(&CurrentCtxDriver, &DriverList[idx]);
        if(oldiface && oldiface != &DriverList[idx])
            oldiface->alcMakeContextCurrent(NULL);
    }
    almtx_unlock(&ContextSwitchLock);
    altss_set(ThreadCtxDriver, NULL);

    return ALC_TRUE;
}

ALC_API void ALC_APIENTRY alcProcessContext(ALCcontext *context)
{
    if(context)
    {
        ALint idx = LookupPtrIntMapKey(&ContextIfaceMap, context);
        if(idx >= 0)
            return DriverList[idx].alcProcessContext(context);
    }
    ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_CONTEXT);
}

ALC_API void ALC_APIENTRY alcSuspendContext(ALCcontext *context)
{
    if(context)
    {
        ALint idx = LookupPtrIntMapKey(&ContextIfaceMap, context);
        if(idx >= 0)
            return DriverList[idx].alcSuspendContext(context);
    }
    ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_CONTEXT);
}

ALC_API void ALC_APIENTRY alcDestroyContext(ALCcontext *context)
{
    ALint idx;

    if(!context || (idx=LookupPtrIntMapKey(&ContextIfaceMap, context)) < 0)
    {
        ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_CONTEXT);
        return;
    }

    DriverList[idx].alcDestroyContext(context);
    RemovePtrIntMapKey(&ContextIfaceMap, context);
}

ALC_API ALCcontext* ALC_APIENTRY alcGetCurrentContext(void)
{
    DriverIface *iface = altss_get(ThreadCtxDriver);
    if(!iface) iface = ATOMIC_LOAD_SEQ(&CurrentCtxDriver);
    return iface ? iface->alcGetCurrentContext() : NULL;
}

ALC_API ALCdevice* ALC_APIENTRY alcGetContextsDevice(ALCcontext *context)
{
    if(context)
    {
        ALint idx = LookupPtrIntMapKey(&ContextIfaceMap, context);
        if(idx >= 0)
            return DriverList[idx].alcGetContextsDevice(context);
    }
    ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_CONTEXT);
    return NULL;
}


ALC_API ALCenum ALC_APIENTRY alcGetError(ALCdevice *device)
{
    if(device)
    {
        ALint idx = LookupPtrIntMapKey(&DeviceIfaceMap, device);
        if(idx < 0) return ALC_INVALID_DEVICE;
        return DriverList[idx].alcGetError(device);
    }
    return ATOMIC_EXCHANGE_SEQ(&LastError, ALC_NO_ERROR);
}

ALC_API ALCboolean ALC_APIENTRY alcIsExtensionPresent(ALCdevice *device, const ALCchar *extname)
{
    const char *ptr;
    size_t len;

    if(device)
    {
        ALint idx = LookupPtrIntMapKey(&DeviceIfaceMap, device);
        if(idx < 0)
        {
            ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_DEVICE);
            return ALC_FALSE;
        }
        return DriverList[idx].alcIsExtensionPresent(device, extname);
    }

    len = strlen(extname);
    ptr = alcExtensionList;
    while(ptr && *ptr)
    {
        if(strncasecmp(ptr, extname, len) == 0 && (ptr[len] == '\0' || isspace(ptr[len])))
            return ALC_TRUE;
        if((ptr=strchr(ptr, ' ')) != NULL)
        {
            do {
                ++ptr;
            } while(isspace(*ptr));
        }
    }
    return ALC_FALSE;
}

ALC_API void* ALC_APIENTRY alcGetProcAddress(ALCdevice *device, const ALCchar *funcname)
{
    size_t i;

    if(device)
    {
        ALint idx = LookupPtrIntMapKey(&DeviceIfaceMap, device);
        if(idx < 0)
        {
            ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_DEVICE);
            return NULL;
        }
        return DriverList[idx].alcGetProcAddress(device, funcname);
    }

    for(i = 0;i < COUNTOF(alcFunctions);i++)
    {
        if(strcmp(funcname, alcFunctions[i].funcName) == 0)
            return alcFunctions[i].address;
    }
    return NULL;
}

ALC_API ALCenum ALC_APIENTRY alcGetEnumValue(ALCdevice *device, const ALCchar *enumname)
{
    size_t i;

    if(device)
    {
        ALint idx = LookupPtrIntMapKey(&DeviceIfaceMap, device);
        if(idx < 0)
        {
            ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_DEVICE);
            return 0;
        }
        return DriverList[idx].alcGetEnumValue(device, enumname);
    }

    for(i = 0;i < COUNTOF(alcEnumerations);i++)
    {
        if(strcmp(enumname, alcEnumerations[i].enumName) == 0)
            return alcEnumerations[i].value;
    }
    return 0;
}

ALC_API const ALCchar* ALC_APIENTRY alcGetString(ALCdevice *device, ALCenum param)
{
    ALsizei i = 0;

    if(device)
    {
        ALint idx = LookupPtrIntMapKey(&DeviceIfaceMap, device);
        if(idx < 0)
        {
            ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_DEVICE);
            return NULL;
        }
        return DriverList[idx].alcGetString(device, param);
    }

    switch(param)
    {
    case ALC_NO_ERROR:
        return alcNoError;
    case ALC_INVALID_ENUM:
        return alcErrInvalidEnum;
    case ALC_INVALID_VALUE:
        return alcErrInvalidValue;
    case ALC_INVALID_DEVICE:
        return alcErrInvalidDevice;
    case ALC_INVALID_CONTEXT:
        return alcErrInvalidContext;
    case ALC_OUT_OF_MEMORY:
        return alcErrOutOfMemory;
    case ALC_EXTENSIONS:
        return alcExtensionList;

    case ALC_DEVICE_SPECIFIER:
        almtx_lock(&EnumerationLock);
        ClearDeviceList(&DevicesList);
        for(i = 0;i < DriverListSize;i++)
        {
            /* Only enumerate names from drivers that support it. */
            if(DriverList[i].ALCVer >= MAKE_ALC_VER(1, 1) ||
               DriverList[i].alcIsExtensionPresent(NULL, "ALC_ENUMERATION_EXT"))
                AppendDeviceList(&DevicesList,
                    DriverList[i].alcGetString(NULL, ALC_DEVICE_SPECIFIER), i
                );
        }
        almtx_unlock(&EnumerationLock);
        return DevicesList.Names;

    case ALC_ALL_DEVICES_SPECIFIER:
        almtx_lock(&EnumerationLock);
        ClearDeviceList(&AllDevicesList);
        for(i = 0;i < DriverListSize;i++)
        {
            /* If the driver doesn't support ALC_ENUMERATE_ALL_EXT, substitute
             * standard enumeration.
             */
            if(DriverList[i].alcIsExtensionPresent(NULL, "ALC_ENUMERATE_ALL_EXT"))
                AppendDeviceList(&AllDevicesList,
                    DriverList[i].alcGetString(NULL, ALC_ALL_DEVICES_SPECIFIER), i
                );
            else if(DriverList[i].ALCVer >= MAKE_ALC_VER(1, 1) ||
                    DriverList[i].alcIsExtensionPresent(NULL, "ALC_ENUMERATION_EXT"))
                AppendDeviceList(&AllDevicesList,
                    DriverList[i].alcGetString(NULL, ALC_DEVICE_SPECIFIER), i
                );
        }
        almtx_unlock(&EnumerationLock);
        return AllDevicesList.Names;

    case ALC_CAPTURE_DEVICE_SPECIFIER:
        almtx_lock(&EnumerationLock);
        ClearDeviceList(&CaptureDevicesList);
        for(i = 0;i < DriverListSize;i++)
        {
            if(DriverList[i].ALCVer >= MAKE_ALC_VER(1, 1) ||
               DriverList[i].alcIsExtensionPresent(NULL, "ALC_EXT_CAPTURE"))
                AppendDeviceList(&CaptureDevicesList,
                    DriverList[i].alcGetString(NULL, ALC_CAPTURE_DEVICE_SPECIFIER), i
                );
        }
        almtx_unlock(&EnumerationLock);
        return CaptureDevicesList.Names;

    case ALC_DEFAULT_DEVICE_SPECIFIER:
        for(i = 0;i < DriverListSize;i++)
        {
            if(DriverList[i].ALCVer >= MAKE_ALC_VER(1, 1) ||
               DriverList[i].alcIsExtensionPresent(NULL, "ALC_ENUMERATION_EXT"))
                return DriverList[i].alcGetString(NULL, ALC_DEFAULT_DEVICE_SPECIFIER);
        }
        return "";

    case ALC_DEFAULT_ALL_DEVICES_SPECIFIER:
        for(i = 0;i < DriverListSize;i++)
        {
            if(DriverList[i].alcIsExtensionPresent(NULL, "ALC_ENUMERATE_ALL_EXT"))
                return DriverList[i].alcGetString(NULL, ALC_DEFAULT_ALL_DEVICES_SPECIFIER);
        }
        return "";

    case ALC_CAPTURE_DEFAULT_DEVICE_SPECIFIER:
        for(i = 0;i < DriverListSize;i++)
        {
            if(DriverList[i].ALCVer >= MAKE_ALC_VER(1, 1) ||
               DriverList[i].alcIsExtensionPresent(NULL, "ALC_EXT_CAPTURE"))
                return DriverList[i].alcGetString(NULL, ALC_CAPTURE_DEFAULT_DEVICE_SPECIFIER);
        }
        return "";

    default:
        ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_ENUM);
        break;
    }
    return NULL;
}

ALC_API void ALC_APIENTRY alcGetIntegerv(ALCdevice *device, ALCenum param, ALCsizei size, ALCint *values)
{
    if(device)
    {
        ALint idx = LookupPtrIntMapKey(&DeviceIfaceMap, device);
        if(idx < 0)
        {
            ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_DEVICE);
            return;
        }
        return DriverList[idx].alcGetIntegerv(device, param, size, values);
    }

    if(size <= 0 || values == NULL)
    {
        ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_VALUE);
        return;
    }

    switch(param)
    {
        case ALC_MAJOR_VERSION:
            if(size >= 1)
            {
                values[0] = alcMajorVersion;
                return;
            }
            /*fall-through*/
        case ALC_MINOR_VERSION:
            if(size >= 1)
            {
                values[0] = alcMinorVersion;
                return;
            }
            ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_VALUE);
            return;

        case ALC_ATTRIBUTES_SIZE:
        case ALC_ALL_ATTRIBUTES:
        case ALC_FREQUENCY:
        case ALC_REFRESH:
        case ALC_SYNC:
        case ALC_MONO_SOURCES:
        case ALC_STEREO_SOURCES:
        case ALC_CAPTURE_SAMPLES:
            ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_DEVICE);
            return;

        default:
            ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_ENUM);
            return;
    }
}


ALC_API ALCdevice* ALC_APIENTRY alcCaptureOpenDevice(const ALCchar *devicename, ALCuint frequency, ALCenum format, ALCsizei buffersize)
{
    ALCdevice *device = NULL;
    ALint idx;

    if(devicename && devicename[0] == '\0')
        devicename = NULL;
    if(devicename)
    {
        almtx_lock(&EnumerationLock);
        if(!CaptureDevicesList.Names)
            (void)alcGetString(NULL, ALC_CAPTURE_DEVICE_SPECIFIER);
        idx = GetDriverIndexForName(&CaptureDevicesList, devicename);
        almtx_unlock(&EnumerationLock);
        if(idx < 0)
        {
            ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_VALUE);
            TRACE("Failed to find driver for name \"%s\"\n", devicename);
            return NULL;
        }
        TRACE("Found driver %d for name \"%s\"\n", idx, devicename);
        device = DriverList[idx].alcCaptureOpenDevice(
            devicename, frequency, format, buffersize
        );
    }
    else
    {
        int i;
        for(i = 0;i < DriverListSize;i++)
        {
            if(DriverList[i].ALCVer >= MAKE_ALC_VER(1, 1) ||
               DriverList[i].alcIsExtensionPresent(NULL, "ALC_EXT_CAPTURE"))
            {
                idx = i;
                TRACE("Using default capture device from driver %d\n", idx);
                device = DriverList[idx].alcCaptureOpenDevice(
                    NULL, frequency, format, buffersize
                );
                break;
            }
        }
    }

    if(device)
    {
        if(InsertPtrIntMapEntry(&DeviceIfaceMap, device, idx) != ALC_NO_ERROR)
        {
            DriverList[idx].alcCaptureCloseDevice(device);
            device = NULL;
        }
    }

    return device;
}

ALC_API ALCboolean ALC_APIENTRY alcCaptureCloseDevice(ALCdevice *device)
{
    ALint idx;

    if(!device || (idx=LookupPtrIntMapKey(&DeviceIfaceMap, device)) < 0)
    {
        ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_DEVICE);
        return ALC_FALSE;
    }
    if(!DriverList[idx].alcCaptureCloseDevice(device))
        return ALC_FALSE;
    RemovePtrIntMapKey(&DeviceIfaceMap, device);
    return ALC_TRUE;
}

ALC_API void ALC_APIENTRY alcCaptureStart(ALCdevice *device)
{
    if(device)
    {
        ALint idx = LookupPtrIntMapKey(&DeviceIfaceMap, device);
        if(idx >= 0)
            return DriverList[idx].alcCaptureStart(device);
    }
    ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_DEVICE);
}

ALC_API void ALC_APIENTRY alcCaptureStop(ALCdevice *device)
{
    if(device)
    {
        ALint idx = LookupPtrIntMapKey(&DeviceIfaceMap, device);
        if(idx >= 0)
            return DriverList[idx].alcCaptureStop(device);
    }
    ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_DEVICE);
}

ALC_API void ALC_APIENTRY alcCaptureSamples(ALCdevice *device, ALCvoid *buffer, ALCsizei samples)
{
    if(device)
    {
        ALint idx = LookupPtrIntMapKey(&DeviceIfaceMap, device);
        if(idx >= 0)
            return DriverList[idx].alcCaptureSamples(device, buffer, samples);
    }
    ATOMIC_STORE_SEQ(&LastError, ALC_INVALID_DEVICE);
}


ALC_API ALCboolean ALC_APIENTRY alcSetThreadContext(ALCcontext *context)
{
    ALCenum err = ALC_INVALID_CONTEXT;
    ALint idx;

    if(!context)
    {
        DriverIface *oldiface = altss_get(ThreadCtxDriver);
        if(oldiface && !oldiface->alcSetThreadContext(NULL))
            return ALC_FALSE;
        altss_set(ThreadCtxDriver, NULL);
        return ALC_TRUE;
    }

    idx = LookupPtrIntMapKey(&ContextIfaceMap, context);
    if(idx >= 0)
    {
        if(DriverList[idx].alcSetThreadContext(context))
        {
            DriverIface *oldiface = altss_get(ThreadCtxDriver);
            if(oldiface != &DriverList[idx])
            {
                altss_set(ThreadCtxDriver, &DriverList[idx]);
                if(oldiface) oldiface->alcSetThreadContext(NULL);
            }
            return ALC_TRUE;
        }
        err = DriverList[idx].alcGetError(NULL);
    }
    ATOMIC_STORE_SEQ(&LastError, err);
    return ALC_FALSE;
}

ALC_API ALCcontext* ALC_APIENTRY alcGetThreadContext(void)
{
    DriverIface *iface = altss_get(ThreadCtxDriver);
    if(iface) return iface->alcGetThreadContext();
    return NULL;
}
