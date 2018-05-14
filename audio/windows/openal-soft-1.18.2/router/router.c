
#include "config.h"

#include "router.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "AL/alc.h"
#include "AL/al.h"
#include "almalloc.h"

#include "version.h"


DriverIface *DriverList = NULL;
int DriverListSize = 0;
static int DriverListSizeMax = 0;

altss_t ThreadCtxDriver;

enum LogLevel LogLevel = LogLevel_Error;
FILE *LogFile;

static void LoadDriverList(void);


BOOL APIENTRY DllMain(HINSTANCE UNUSED(module), DWORD reason, void* UNUSED(reserved))
{
    const char *str;
    int i;

    switch(reason)
    {
        case DLL_PROCESS_ATTACH:
            LogFile = stderr;
            str = getenv("ALROUTER_LOGFILE");
            if(str && *str != '\0')
            {
                FILE *f = fopen(str, "w");
                if(f == NULL)
                    ERR("Could not open log file: %s\n", str);
                else
                    LogFile = f;
            }
            str = getenv("ALROUTER_LOGLEVEL");
            if(str && *str != '\0')
            {
                char *end = NULL;
                long l = strtol(str, &end, 0);
                if(!end || *end != '\0')
                    ERR("Invalid log level value: %s\n", str);
                else if(l < LogLevel_None || l > LogLevel_Trace)
                    ERR("Log level out of range: %s\n", str);
                else
                    LogLevel = l;
            }
            TRACE("Initializing router v0.1-%s %s\n", ALSOFT_GIT_COMMIT_HASH, ALSOFT_GIT_BRANCH);
            LoadDriverList();

            altss_create(&ThreadCtxDriver, NULL);
            InitALC();
            break;

        case DLL_THREAD_ATTACH:
        case DLL_THREAD_DETACH:
            break;

        case DLL_PROCESS_DETACH:
            ReleaseALC();
            altss_delete(ThreadCtxDriver);

            for(i = 0;i < DriverListSize;i++)
            {
                if(DriverList[i].Module)
                    FreeLibrary(DriverList[i].Module);
            }
            al_free(DriverList);
            DriverList = NULL;
            DriverListSize = 0;
            DriverListSizeMax = 0;

            if(LogFile && LogFile != stderr)
                fclose(LogFile);
            LogFile = NULL;
            break;
    }
    return TRUE;
}


#ifdef __GNUC__
#define CAST_FUNC(x) (__typeof(x))
#else
#define CAST_FUNC(x) (void*)
#endif

static void AddModule(HMODULE module, const WCHAR *name)
{
    DriverIface newdrv;
    int err = 0;
    int i;

    for(i = 0;i < DriverListSize;i++)
    {
        if(DriverList[i].Module == module)
        {
            TRACE("Skipping already-loaded module %p\n", module);
            FreeLibrary(module);
            return;
        }
        if(wcscmp(DriverList[i].Name, name) == 0)
        {
            TRACE("Skipping similarly-named module %ls\n", name);
            FreeLibrary(module);
            return;
        }
    }

    if(DriverListSize == DriverListSizeMax)
    {
        int newmax = DriverListSizeMax ? DriverListSizeMax<<1 : 4;
        void *newlist = al_calloc(DEF_ALIGN, sizeof(DriverList[0])*newmax);
        if(!newlist) return;

        memcpy(newlist, DriverList, DriverListSize*sizeof(DriverList[0]));
        al_free(DriverList);
        DriverList = newlist;
        DriverListSizeMax = newmax;
    }

    memset(&newdrv, 0, sizeof(newdrv));
    /* Load required functions. */
#define LOAD_PROC(x) do {                                                     \
    newdrv.x = CAST_FUNC(newdrv.x) GetProcAddress(module, #x);                \
    if(!newdrv.x)                                                             \
    {                                                                         \
        ERR("Failed to find entry point for %s in %ls\n", #x, name);          \
        err = 1;                                                              \
    }                                                                         \
} while(0)
    LOAD_PROC(alcCreateContext);
    LOAD_PROC(alcMakeContextCurrent);
    LOAD_PROC(alcProcessContext);
    LOAD_PROC(alcSuspendContext);
    LOAD_PROC(alcDestroyContext);
    LOAD_PROC(alcGetCurrentContext);
    LOAD_PROC(alcGetContextsDevice);
    LOAD_PROC(alcOpenDevice);
    LOAD_PROC(alcCloseDevice);
    LOAD_PROC(alcGetError);
    LOAD_PROC(alcIsExtensionPresent);
    LOAD_PROC(alcGetProcAddress);
    LOAD_PROC(alcGetEnumValue);
    LOAD_PROC(alcGetString);
    LOAD_PROC(alcGetIntegerv);
    LOAD_PROC(alcCaptureOpenDevice);
    LOAD_PROC(alcCaptureCloseDevice);
    LOAD_PROC(alcCaptureStart);
    LOAD_PROC(alcCaptureStop);
    LOAD_PROC(alcCaptureSamples);

    LOAD_PROC(alEnable);
    LOAD_PROC(alDisable);
    LOAD_PROC(alIsEnabled);
    LOAD_PROC(alGetString);
    LOAD_PROC(alGetBooleanv);
    LOAD_PROC(alGetIntegerv);
    LOAD_PROC(alGetFloatv);
    LOAD_PROC(alGetDoublev);
    LOAD_PROC(alGetBoolean);
    LOAD_PROC(alGetInteger);
    LOAD_PROC(alGetFloat);
    LOAD_PROC(alGetDouble);
    LOAD_PROC(alGetError);
    LOAD_PROC(alIsExtensionPresent);
    LOAD_PROC(alGetProcAddress);
    LOAD_PROC(alGetEnumValue);
    LOAD_PROC(alListenerf);
    LOAD_PROC(alListener3f);
    LOAD_PROC(alListenerfv);
    LOAD_PROC(alListeneri);
    LOAD_PROC(alListener3i);
    LOAD_PROC(alListeneriv);
    LOAD_PROC(alGetListenerf);
    LOAD_PROC(alGetListener3f);
    LOAD_PROC(alGetListenerfv);
    LOAD_PROC(alGetListeneri);
    LOAD_PROC(alGetListener3i);
    LOAD_PROC(alGetListeneriv);
    LOAD_PROC(alGenSources);
    LOAD_PROC(alDeleteSources);
    LOAD_PROC(alIsSource);
    LOAD_PROC(alSourcef);
    LOAD_PROC(alSource3f);
    LOAD_PROC(alSourcefv);
    LOAD_PROC(alSourcei);
    LOAD_PROC(alSource3i);
    LOAD_PROC(alSourceiv);
    LOAD_PROC(alGetSourcef);
    LOAD_PROC(alGetSource3f);
    LOAD_PROC(alGetSourcefv);
    LOAD_PROC(alGetSourcei);
    LOAD_PROC(alGetSource3i);
    LOAD_PROC(alGetSourceiv);
    LOAD_PROC(alSourcePlayv);
    LOAD_PROC(alSourceStopv);
    LOAD_PROC(alSourceRewindv);
    LOAD_PROC(alSourcePausev);
    LOAD_PROC(alSourcePlay);
    LOAD_PROC(alSourceStop);
    LOAD_PROC(alSourceRewind);
    LOAD_PROC(alSourcePause);
    LOAD_PROC(alSourceQueueBuffers);
    LOAD_PROC(alSourceUnqueueBuffers);
    LOAD_PROC(alGenBuffers);
    LOAD_PROC(alDeleteBuffers);
    LOAD_PROC(alIsBuffer);
    LOAD_PROC(alBufferf);
    LOAD_PROC(alBuffer3f);
    LOAD_PROC(alBufferfv);
    LOAD_PROC(alBufferi);
    LOAD_PROC(alBuffer3i);
    LOAD_PROC(alBufferiv);
    LOAD_PROC(alGetBufferf);
    LOAD_PROC(alGetBuffer3f);
    LOAD_PROC(alGetBufferfv);
    LOAD_PROC(alGetBufferi);
    LOAD_PROC(alGetBuffer3i);
    LOAD_PROC(alGetBufferiv);
    LOAD_PROC(alBufferData);
    LOAD_PROC(alDopplerFactor);
    LOAD_PROC(alDopplerVelocity);
    LOAD_PROC(alSpeedOfSound);
    LOAD_PROC(alDistanceModel);
    if(!err)
    {
        ALCint alc_ver[2] = { 0, 0 };
        wcsncpy(newdrv.Name, name, 32);
        newdrv.Module = module;
        newdrv.alcGetIntegerv(NULL, ALC_MAJOR_VERSION, 1, &alc_ver[0]);
        newdrv.alcGetIntegerv(NULL, ALC_MINOR_VERSION, 1, &alc_ver[1]);
        if(newdrv.alcGetError(NULL) == ALC_NO_ERROR)
            newdrv.ALCVer = MAKE_ALC_VER(alc_ver[0], alc_ver[1]);
        else
            newdrv.ALCVer = MAKE_ALC_VER(1, 0);

#undef LOAD_PROC
#define LOAD_PROC(x) do {                                                     \
    newdrv.x = CAST_FUNC(newdrv.x) newdrv.alcGetProcAddress(NULL, #x);        \
    if(!newdrv.x)                                                             \
    {                                                                         \
        ERR("Failed to find entry point for %s in %ls\n", #x, name);          \
        err = 1;                                                              \
    }                                                                         \
} while(0)
        if(newdrv.alcIsExtensionPresent(NULL, "ALC_EXT_thread_local_context"))
        {
            LOAD_PROC(alcSetThreadContext);
            LOAD_PROC(alcGetThreadContext);
        }
    }

    if(!err)
    {
        TRACE("Loaded module %p, %ls, ALC %d.%d\n", module, name,
              newdrv.ALCVer>>8, newdrv.ALCVer&255);
        DriverList[DriverListSize++] = newdrv;
    }
#undef LOAD_PROC
}

static void SearchDrivers(WCHAR *path)
{
    WCHAR srchPath[MAX_PATH+1] = L"";
    WIN32_FIND_DATAW fdata;
    HANDLE srchHdl;

    TRACE("Searching for drivers in %ls...\n", path);
    wcsncpy(srchPath, path, MAX_PATH);
    wcsncat(srchPath, L"\\*oal.dll", MAX_PATH - lstrlenW(srchPath));
    srchHdl = FindFirstFileW(srchPath, &fdata);
    if(srchHdl != INVALID_HANDLE_VALUE)
    {
        do {
            HMODULE mod;

            wcsncpy(srchPath, path, MAX_PATH);
            wcsncat(srchPath, L"\\", MAX_PATH - lstrlenW(srchPath));
            wcsncat(srchPath, fdata.cFileName, MAX_PATH - lstrlenW(srchPath));
            TRACE("Found %ls\n", srchPath);

            mod = LoadLibraryW(srchPath);
            if(!mod)
                WARN("Could not load %ls\n", srchPath);
            else
                AddModule(mod, fdata.cFileName);
        } while(FindNextFileW(srchHdl, &fdata));
        FindClose(srchHdl);
    }
}

static WCHAR *strrchrW(WCHAR *str, WCHAR ch)
{
    WCHAR *res = NULL;
    while(str && *str != '\0')
    {
        if(*str == ch)
            res = str;
        ++str;
    }
    return res;
}

static int GetLoadedModuleDirectory(const WCHAR *name, WCHAR *moddir, DWORD length)
{
    HMODULE module = NULL;
    WCHAR *sep0, *sep1;

    if(name)
    {
        module = GetModuleHandleW(name);
        if(!module) return 0;
    }

    if(GetModuleFileNameW(module, moddir, length) == 0)
        return 0;

    sep0 = strrchrW(moddir, '/');
    if(sep0) sep1 = strrchrW(sep0+1, '\\');
    else sep1 = strrchrW(moddir, '\\');

    if(sep1) *sep1 = '\0';
    else if(sep0) *sep0 = '\0';
    else *moddir = '\0';

    return 1;
}

void LoadDriverList(void)
{
    WCHAR path[MAX_PATH+1] = L"";
    int len;

    if(GetLoadedModuleDirectory(L"OpenAL32.dll", path, MAX_PATH))
        SearchDrivers(path);

    GetCurrentDirectoryW(MAX_PATH, path);
    len = lstrlenW(path);
    if(len > 0 && (path[len-1] == '\\' || path[len-1] == '/'))
        path[len-1] = '\0';
    SearchDrivers(path);

    if(GetLoadedModuleDirectory(NULL, path, MAX_PATH))
        SearchDrivers(path);

    GetSystemDirectoryW(path, MAX_PATH);
    len = lstrlenW(path);
    if(len > 0 && (path[len-1] == '\\' || path[len-1] == '/'))
        path[len-1] = '\0';
    SearchDrivers(path);
}


void InitPtrIntMap(PtrIntMap *map)
{
    map->keys = NULL;
    map->values = NULL;
    map->size = 0;
    map->capacity = 0;
    RWLockInit(&map->lock);
}

void ResetPtrIntMap(PtrIntMap *map)
{
    WriteLock(&map->lock);
    al_free(map->keys);
    map->keys = NULL;
    map->values = NULL;
    map->size = 0;
    map->capacity = 0;
    WriteUnlock(&map->lock);
}

ALenum InsertPtrIntMapEntry(PtrIntMap *map, ALvoid *key, ALint value)
{
    ALsizei pos = 0;

    WriteLock(&map->lock);
    if(map->size > 0)
    {
        ALsizei count = map->size;
        do {
            ALsizei step = count>>1;
            ALsizei i = pos+step;
            if(!(map->keys[i] < key))
                count = step;
            else
            {
                pos = i+1;
                count -= step+1;
            }
        } while(count > 0);
    }

    if(pos == map->size || map->keys[pos] != key)
    {
        if(map->size == map->capacity)
        {
            ALvoid **keys = NULL;
            ALint *values;
            ALsizei newcap;

            newcap = (map->capacity ? (map->capacity<<1) : 4);
            if(newcap > map->capacity)
                keys = al_calloc(16, (sizeof(map->keys[0])+sizeof(map->values[0]))*newcap);
            if(!keys)
            {
                WriteUnlock(&map->lock);
                return AL_OUT_OF_MEMORY;
            }
            values = (ALint*)&keys[newcap];

            if(map->keys)
            {
                memcpy(keys, map->keys, map->size*sizeof(map->keys[0]));
                memcpy(values, map->values, map->size*sizeof(map->values[0]));
            }
            al_free(map->keys);
            map->keys = keys;
            map->values = values;
            map->capacity = newcap;
        }

        if(pos < map->size)
        {
            memmove(&map->keys[pos+1], &map->keys[pos],
                    (map->size-pos)*sizeof(map->keys[0]));
            memmove(&map->values[pos+1], &map->values[pos],
                    (map->size-pos)*sizeof(map->values[0]));
        }
        map->size++;
    }
    map->keys[pos] = key;
    map->values[pos] = value;
    WriteUnlock(&map->lock);

    return AL_NO_ERROR;
}

ALint RemovePtrIntMapKey(PtrIntMap *map, ALvoid *key)
{
    ALint ret = -1;
    WriteLock(&map->lock);
    if(map->size > 0)
    {
        ALsizei pos = 0;
        ALsizei count = map->size;
        do {
            ALsizei step = count>>1;
            ALsizei i = pos+step;
            if(!(map->keys[i] < key))
                count = step;
            else
            {
                pos = i+1;
                count -= step+1;
            }
        } while(count > 0);
        if(pos < map->size && map->keys[pos] == key)
        {
            ret = map->values[pos];
            if(pos < map->size-1)
            {
                memmove(&map->keys[pos], &map->keys[pos+1],
                        (map->size-1-pos)*sizeof(map->keys[0]));
                memmove(&map->values[pos], &map->values[pos+1],
                        (map->size-1-pos)*sizeof(map->values[0]));
            }
            map->size--;
        }
    }
    WriteUnlock(&map->lock);
    return ret;
}

ALint LookupPtrIntMapKey(PtrIntMap *map, ALvoid *key)
{
    ALint ret = -1;
    ReadLock(&map->lock);
    if(map->size > 0)
    {
        ALsizei pos = 0;
        ALsizei count = map->size;
        do {
            ALsizei step = count>>1;
            ALsizei i = pos+step;
            if(!(map->keys[i] < key))
                count = step;
            else
            {
                pos = i+1;
                count -= step+1;
            }
        } while(count > 0);
        if(pos < map->size && map->keys[pos] == key)
            ret = map->values[pos];
    }
    ReadUnlock(&map->lock);
    return ret;
}
