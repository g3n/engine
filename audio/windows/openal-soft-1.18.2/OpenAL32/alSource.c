/**
 * OpenAL cross platform audio library
 * Copyright (C) 1999-2007 by authors.
 * This library is free software; you can redistribute it and/or
 *  modify it under the terms of the GNU Library General Public
 *  License as published by the Free Software Foundation; either
 *  version 2 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 *  Library General Public License for more details.
 *
 * You should have received a copy of the GNU Library General Public
 *  License along with this library; if not, write to the
 *  Free Software Foundation, Inc.,
 *  51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 * Or go to http://www.gnu.org/copyleft/lgpl.html
 */

#include "config.h"

#include <stdlib.h>
#include <limits.h>
#include <math.h>
#include <float.h>

#include "AL/al.h"
#include "AL/alc.h"
#include "alMain.h"
#include "alError.h"
#include "alSource.h"
#include "alBuffer.h"
#include "alThunk.h"
#include "alAuxEffectSlot.h"

#include "backends/base.h"

#include "threads.h"
#include "almalloc.h"


extern inline void LockSourcesRead(ALCcontext *context);
extern inline void UnlockSourcesRead(ALCcontext *context);
extern inline void LockSourcesWrite(ALCcontext *context);
extern inline void UnlockSourcesWrite(ALCcontext *context);
extern inline struct ALsource *LookupSource(ALCcontext *context, ALuint id);
extern inline struct ALsource *RemoveSource(ALCcontext *context, ALuint id);

static void InitSourceParams(ALsource *Source, ALsizei num_sends);
static void DeinitSource(ALsource *source, ALsizei num_sends);
static void UpdateSourceProps(ALsource *source, ALvoice *voice, ALsizei num_sends);
static ALint64 GetSourceSampleOffset(ALsource *Source, ALCcontext *context, ALuint64 *clocktime);
static ALdouble GetSourceSecOffset(ALsource *Source, ALCcontext *context, ALuint64 *clocktime);
static ALdouble GetSourceOffset(ALsource *Source, ALenum name, ALCcontext *context);
static ALboolean GetSampleOffset(ALsource *Source, ALuint *offset, ALsizei *frac);
static ALboolean ApplyOffset(ALsource *Source, ALvoice *voice);

typedef enum SourceProp {
    srcPitch = AL_PITCH,
    srcGain = AL_GAIN,
    srcMinGain = AL_MIN_GAIN,
    srcMaxGain = AL_MAX_GAIN,
    srcMaxDistance = AL_MAX_DISTANCE,
    srcRolloffFactor = AL_ROLLOFF_FACTOR,
    srcDopplerFactor = AL_DOPPLER_FACTOR,
    srcConeOuterGain = AL_CONE_OUTER_GAIN,
    srcSecOffset = AL_SEC_OFFSET,
    srcSampleOffset = AL_SAMPLE_OFFSET,
    srcByteOffset = AL_BYTE_OFFSET,
    srcConeInnerAngle = AL_CONE_INNER_ANGLE,
    srcConeOuterAngle = AL_CONE_OUTER_ANGLE,
    srcRefDistance = AL_REFERENCE_DISTANCE,

    srcPosition = AL_POSITION,
    srcVelocity = AL_VELOCITY,
    srcDirection = AL_DIRECTION,

    srcSourceRelative = AL_SOURCE_RELATIVE,
    srcLooping = AL_LOOPING,
    srcBuffer = AL_BUFFER,
    srcSourceState = AL_SOURCE_STATE,
    srcBuffersQueued = AL_BUFFERS_QUEUED,
    srcBuffersProcessed = AL_BUFFERS_PROCESSED,
    srcSourceType = AL_SOURCE_TYPE,

    /* ALC_EXT_EFX */
    srcConeOuterGainHF = AL_CONE_OUTER_GAINHF,
    srcAirAbsorptionFactor = AL_AIR_ABSORPTION_FACTOR,
    srcRoomRolloffFactor =  AL_ROOM_ROLLOFF_FACTOR,
    srcDirectFilterGainHFAuto = AL_DIRECT_FILTER_GAINHF_AUTO,
    srcAuxSendFilterGainAuto = AL_AUXILIARY_SEND_FILTER_GAIN_AUTO,
    srcAuxSendFilterGainHFAuto = AL_AUXILIARY_SEND_FILTER_GAINHF_AUTO,
    srcDirectFilter = AL_DIRECT_FILTER,
    srcAuxSendFilter = AL_AUXILIARY_SEND_FILTER,

    /* AL_SOFT_direct_channels */
    srcDirectChannelsSOFT = AL_DIRECT_CHANNELS_SOFT,

    /* AL_EXT_source_distance_model */
    srcDistanceModel = AL_DISTANCE_MODEL,

    srcByteLengthSOFT = AL_BYTE_LENGTH_SOFT,
    srcSampleLengthSOFT = AL_SAMPLE_LENGTH_SOFT,
    srcSecLengthSOFT = AL_SEC_LENGTH_SOFT,

    /* AL_SOFT_source_latency */
    srcSampleOffsetLatencySOFT = AL_SAMPLE_OFFSET_LATENCY_SOFT,
    srcSecOffsetLatencySOFT = AL_SEC_OFFSET_LATENCY_SOFT,

    /* AL_EXT_STEREO_ANGLES */
    srcAngles = AL_STEREO_ANGLES,

    /* AL_EXT_SOURCE_RADIUS */
    srcRadius = AL_SOURCE_RADIUS,

    /* AL_EXT_BFORMAT */
    srcOrientation = AL_ORIENTATION,

    /* AL_SOFT_source_resampler */
    srcResampler = AL_SOURCE_RESAMPLER_SOFT,

    /* AL_SOFT_source_spatialize */
    srcSpatialize = AL_SOURCE_SPATIALIZE_SOFT,
} SourceProp;

static ALboolean SetSourcefv(ALsource *Source, ALCcontext *Context, SourceProp prop, const ALfloat *values);
static ALboolean SetSourceiv(ALsource *Source, ALCcontext *Context, SourceProp prop, const ALint *values);
static ALboolean SetSourcei64v(ALsource *Source, ALCcontext *Context, SourceProp prop, const ALint64SOFT *values);

static ALboolean GetSourcedv(ALsource *Source, ALCcontext *Context, SourceProp prop, ALdouble *values);
static ALboolean GetSourceiv(ALsource *Source, ALCcontext *Context, SourceProp prop, ALint *values);
static ALboolean GetSourcei64v(ALsource *Source, ALCcontext *Context, SourceProp prop, ALint64 *values);

static inline ALvoice *GetSourceVoice(const ALsource *source, const ALCcontext *context)
{
    ALvoice **voice = context->Voices;
    ALvoice **voice_end = voice + context->VoiceCount;
    while(voice != voice_end)
    {
        if(ATOMIC_LOAD(&(*voice)->Source, almemory_order_acquire) == source)
            return *voice;
        ++voice;
    }
    return NULL;
}

/**
 * Returns if the last known state for the source was playing or paused. Does
 * not sync with the mixer voice.
 */
static inline bool IsPlayingOrPaused(ALsource *source)
{
    ALenum state = ATOMIC_LOAD(&source->state, almemory_order_acquire);
    return state == AL_PLAYING || state == AL_PAUSED;
}

/**
 * Returns an updated source state using the matching voice's status (or lack
 * thereof).
 */
static inline ALenum GetSourceState(ALsource *source, ALvoice *voice)
{
    if(!voice)
    {
        ALenum state = AL_PLAYING;
        if(ATOMIC_COMPARE_EXCHANGE_STRONG(&source->state, &state, AL_STOPPED,
                                          almemory_order_acq_rel, almemory_order_acquire))
            return AL_STOPPED;
        return state;
    }
    return ATOMIC_LOAD(&source->state, almemory_order_acquire);
}

/**
 * Returns if the source should specify an update, given the context's
 * deferring state and the source's last known state.
 */
static inline bool SourceShouldUpdate(ALsource *source, ALCcontext *context)
{
    return !ATOMIC_LOAD(&context->DeferUpdates, almemory_order_acquire) &&
           IsPlayingOrPaused(source);
}

static ALint FloatValsByProp(ALenum prop)
{
    if(prop != (ALenum)((SourceProp)prop))
        return 0;
    switch((SourceProp)prop)
    {
        case AL_PITCH:
        case AL_GAIN:
        case AL_MIN_GAIN:
        case AL_MAX_GAIN:
        case AL_MAX_DISTANCE:
        case AL_ROLLOFF_FACTOR:
        case AL_DOPPLER_FACTOR:
        case AL_CONE_OUTER_GAIN:
        case AL_SEC_OFFSET:
        case AL_SAMPLE_OFFSET:
        case AL_BYTE_OFFSET:
        case AL_CONE_INNER_ANGLE:
        case AL_CONE_OUTER_ANGLE:
        case AL_REFERENCE_DISTANCE:
        case AL_CONE_OUTER_GAINHF:
        case AL_AIR_ABSORPTION_FACTOR:
        case AL_ROOM_ROLLOFF_FACTOR:
        case AL_DIRECT_FILTER_GAINHF_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAIN_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAINHF_AUTO:
        case AL_DIRECT_CHANNELS_SOFT:
        case AL_DISTANCE_MODEL:
        case AL_SOURCE_RELATIVE:
        case AL_LOOPING:
        case AL_SOURCE_STATE:
        case AL_BUFFERS_QUEUED:
        case AL_BUFFERS_PROCESSED:
        case AL_SOURCE_TYPE:
        case AL_BYTE_LENGTH_SOFT:
        case AL_SAMPLE_LENGTH_SOFT:
        case AL_SEC_LENGTH_SOFT:
        case AL_SOURCE_RADIUS:
        case AL_SOURCE_RESAMPLER_SOFT:
        case AL_SOURCE_SPATIALIZE_SOFT:
            return 1;

        case AL_STEREO_ANGLES:
            return 2;

        case AL_POSITION:
        case AL_VELOCITY:
        case AL_DIRECTION:
            return 3;

        case AL_ORIENTATION:
            return 6;

        case AL_SEC_OFFSET_LATENCY_SOFT:
            break; /* Double only */

        case AL_BUFFER:
        case AL_DIRECT_FILTER:
        case AL_AUXILIARY_SEND_FILTER:
            break; /* i/i64 only */
        case AL_SAMPLE_OFFSET_LATENCY_SOFT:
            break; /* i64 only */
    }
    return 0;
}
static ALint DoubleValsByProp(ALenum prop)
{
    if(prop != (ALenum)((SourceProp)prop))
        return 0;
    switch((SourceProp)prop)
    {
        case AL_PITCH:
        case AL_GAIN:
        case AL_MIN_GAIN:
        case AL_MAX_GAIN:
        case AL_MAX_DISTANCE:
        case AL_ROLLOFF_FACTOR:
        case AL_DOPPLER_FACTOR:
        case AL_CONE_OUTER_GAIN:
        case AL_SEC_OFFSET:
        case AL_SAMPLE_OFFSET:
        case AL_BYTE_OFFSET:
        case AL_CONE_INNER_ANGLE:
        case AL_CONE_OUTER_ANGLE:
        case AL_REFERENCE_DISTANCE:
        case AL_CONE_OUTER_GAINHF:
        case AL_AIR_ABSORPTION_FACTOR:
        case AL_ROOM_ROLLOFF_FACTOR:
        case AL_DIRECT_FILTER_GAINHF_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAIN_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAINHF_AUTO:
        case AL_DIRECT_CHANNELS_SOFT:
        case AL_DISTANCE_MODEL:
        case AL_SOURCE_RELATIVE:
        case AL_LOOPING:
        case AL_SOURCE_STATE:
        case AL_BUFFERS_QUEUED:
        case AL_BUFFERS_PROCESSED:
        case AL_SOURCE_TYPE:
        case AL_BYTE_LENGTH_SOFT:
        case AL_SAMPLE_LENGTH_SOFT:
        case AL_SEC_LENGTH_SOFT:
        case AL_SOURCE_RADIUS:
        case AL_SOURCE_RESAMPLER_SOFT:
        case AL_SOURCE_SPATIALIZE_SOFT:
            return 1;

        case AL_SEC_OFFSET_LATENCY_SOFT:
        case AL_STEREO_ANGLES:
            return 2;

        case AL_POSITION:
        case AL_VELOCITY:
        case AL_DIRECTION:
            return 3;

        case AL_ORIENTATION:
            return 6;

        case AL_BUFFER:
        case AL_DIRECT_FILTER:
        case AL_AUXILIARY_SEND_FILTER:
            break; /* i/i64 only */
        case AL_SAMPLE_OFFSET_LATENCY_SOFT:
            break; /* i64 only */
    }
    return 0;
}

static ALint IntValsByProp(ALenum prop)
{
    if(prop != (ALenum)((SourceProp)prop))
        return 0;
    switch((SourceProp)prop)
    {
        case AL_PITCH:
        case AL_GAIN:
        case AL_MIN_GAIN:
        case AL_MAX_GAIN:
        case AL_MAX_DISTANCE:
        case AL_ROLLOFF_FACTOR:
        case AL_DOPPLER_FACTOR:
        case AL_CONE_OUTER_GAIN:
        case AL_SEC_OFFSET:
        case AL_SAMPLE_OFFSET:
        case AL_BYTE_OFFSET:
        case AL_CONE_INNER_ANGLE:
        case AL_CONE_OUTER_ANGLE:
        case AL_REFERENCE_DISTANCE:
        case AL_CONE_OUTER_GAINHF:
        case AL_AIR_ABSORPTION_FACTOR:
        case AL_ROOM_ROLLOFF_FACTOR:
        case AL_DIRECT_FILTER_GAINHF_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAIN_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAINHF_AUTO:
        case AL_DIRECT_CHANNELS_SOFT:
        case AL_DISTANCE_MODEL:
        case AL_SOURCE_RELATIVE:
        case AL_LOOPING:
        case AL_BUFFER:
        case AL_SOURCE_STATE:
        case AL_BUFFERS_QUEUED:
        case AL_BUFFERS_PROCESSED:
        case AL_SOURCE_TYPE:
        case AL_DIRECT_FILTER:
        case AL_BYTE_LENGTH_SOFT:
        case AL_SAMPLE_LENGTH_SOFT:
        case AL_SEC_LENGTH_SOFT:
        case AL_SOURCE_RADIUS:
        case AL_SOURCE_RESAMPLER_SOFT:
        case AL_SOURCE_SPATIALIZE_SOFT:
            return 1;

        case AL_POSITION:
        case AL_VELOCITY:
        case AL_DIRECTION:
        case AL_AUXILIARY_SEND_FILTER:
            return 3;

        case AL_ORIENTATION:
            return 6;

        case AL_SAMPLE_OFFSET_LATENCY_SOFT:
            break; /* i64 only */
        case AL_SEC_OFFSET_LATENCY_SOFT:
            break; /* Double only */
        case AL_STEREO_ANGLES:
            break; /* Float/double only */
    }
    return 0;
}
static ALint Int64ValsByProp(ALenum prop)
{
    if(prop != (ALenum)((SourceProp)prop))
        return 0;
    switch((SourceProp)prop)
    {
        case AL_PITCH:
        case AL_GAIN:
        case AL_MIN_GAIN:
        case AL_MAX_GAIN:
        case AL_MAX_DISTANCE:
        case AL_ROLLOFF_FACTOR:
        case AL_DOPPLER_FACTOR:
        case AL_CONE_OUTER_GAIN:
        case AL_SEC_OFFSET:
        case AL_SAMPLE_OFFSET:
        case AL_BYTE_OFFSET:
        case AL_CONE_INNER_ANGLE:
        case AL_CONE_OUTER_ANGLE:
        case AL_REFERENCE_DISTANCE:
        case AL_CONE_OUTER_GAINHF:
        case AL_AIR_ABSORPTION_FACTOR:
        case AL_ROOM_ROLLOFF_FACTOR:
        case AL_DIRECT_FILTER_GAINHF_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAIN_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAINHF_AUTO:
        case AL_DIRECT_CHANNELS_SOFT:
        case AL_DISTANCE_MODEL:
        case AL_SOURCE_RELATIVE:
        case AL_LOOPING:
        case AL_BUFFER:
        case AL_SOURCE_STATE:
        case AL_BUFFERS_QUEUED:
        case AL_BUFFERS_PROCESSED:
        case AL_SOURCE_TYPE:
        case AL_DIRECT_FILTER:
        case AL_BYTE_LENGTH_SOFT:
        case AL_SAMPLE_LENGTH_SOFT:
        case AL_SEC_LENGTH_SOFT:
        case AL_SOURCE_RADIUS:
        case AL_SOURCE_RESAMPLER_SOFT:
        case AL_SOURCE_SPATIALIZE_SOFT:
            return 1;

        case AL_SAMPLE_OFFSET_LATENCY_SOFT:
            return 2;

        case AL_POSITION:
        case AL_VELOCITY:
        case AL_DIRECTION:
        case AL_AUXILIARY_SEND_FILTER:
            return 3;

        case AL_ORIENTATION:
            return 6;

        case AL_SEC_OFFSET_LATENCY_SOFT:
            break; /* Double only */
        case AL_STEREO_ANGLES:
            break; /* Float/double only */
    }
    return 0;
}


#define CHECKVAL(x) do {                                                      \
    if(!(x))                                                                  \
        SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_VALUE, AL_FALSE);      \
} while(0)

#define DO_UPDATEPROPS() do {                                                 \
    ALvoice *voice;                                                           \
    if(SourceShouldUpdate(Source, Context) &&                                 \
       (voice=GetSourceVoice(Source, Context)) != NULL)                       \
        UpdateSourceProps(Source, voice, device->NumAuxSends);                \
    else                                                                      \
        ATOMIC_FLAG_CLEAR(&Source->PropsClean, almemory_order_release);       \
} while(0)

static ALboolean SetSourcefv(ALsource *Source, ALCcontext *Context, SourceProp prop, const ALfloat *values)
{
    ALCdevice *device = Context->Device;
    ALint ival;

    switch(prop)
    {
        case AL_BYTE_LENGTH_SOFT:
        case AL_SAMPLE_LENGTH_SOFT:
        case AL_SEC_LENGTH_SOFT:
        case AL_SEC_OFFSET_LATENCY_SOFT:
            /* Query only */
            SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_OPERATION, AL_FALSE);

        case AL_PITCH:
            CHECKVAL(*values >= 0.0f);

            Source->Pitch = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_CONE_INNER_ANGLE:
            CHECKVAL(*values >= 0.0f && *values <= 360.0f);

            Source->InnerAngle = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_CONE_OUTER_ANGLE:
            CHECKVAL(*values >= 0.0f && *values <= 360.0f);

            Source->OuterAngle = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_GAIN:
            CHECKVAL(*values >= 0.0f);

            Source->Gain = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_MAX_DISTANCE:
            CHECKVAL(*values >= 0.0f);

            Source->MaxDistance = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_ROLLOFF_FACTOR:
            CHECKVAL(*values >= 0.0f);

            Source->RolloffFactor = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_REFERENCE_DISTANCE:
            CHECKVAL(*values >= 0.0f);

            Source->RefDistance = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_MIN_GAIN:
            CHECKVAL(*values >= 0.0f);

            Source->MinGain = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_MAX_GAIN:
            CHECKVAL(*values >= 0.0f);

            Source->MaxGain = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_CONE_OUTER_GAIN:
            CHECKVAL(*values >= 0.0f && *values <= 1.0f);

            Source->OuterGain = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_CONE_OUTER_GAINHF:
            CHECKVAL(*values >= 0.0f && *values <= 1.0f);

            Source->OuterGainHF = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_AIR_ABSORPTION_FACTOR:
            CHECKVAL(*values >= 0.0f && *values <= 10.0f);

            Source->AirAbsorptionFactor = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_ROOM_ROLLOFF_FACTOR:
            CHECKVAL(*values >= 0.0f && *values <= 10.0f);

            Source->RoomRolloffFactor = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_DOPPLER_FACTOR:
            CHECKVAL(*values >= 0.0f && *values <= 1.0f);

            Source->DopplerFactor = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_SEC_OFFSET:
        case AL_SAMPLE_OFFSET:
        case AL_BYTE_OFFSET:
            CHECKVAL(*values >= 0.0f);

            Source->OffsetType = prop;
            Source->Offset = *values;

            if(IsPlayingOrPaused(Source))
            {
                ALvoice *voice;

                ALCdevice_Lock(Context->Device);
                /* Double-check that the source is still playing while we have
                 * the lock.
                 */
                voice = GetSourceVoice(Source, Context);
                if(voice)
                {
                    WriteLock(&Source->queue_lock);
                    if(ApplyOffset(Source, voice) == AL_FALSE)
                    {
                        WriteUnlock(&Source->queue_lock);
                        ALCdevice_Unlock(Context->Device);
                        SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_VALUE, AL_FALSE);
                    }
                    WriteUnlock(&Source->queue_lock);
                }
                ALCdevice_Unlock(Context->Device);
            }
            return AL_TRUE;

        case AL_SOURCE_RADIUS:
            CHECKVAL(*values >= 0.0f && isfinite(*values));

            Source->Radius = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_STEREO_ANGLES:
            CHECKVAL(isfinite(values[0]) && isfinite(values[1]));

            Source->StereoPan[0] = values[0];
            Source->StereoPan[1] = values[1];
            DO_UPDATEPROPS();
            return AL_TRUE;


        case AL_POSITION:
            CHECKVAL(isfinite(values[0]) && isfinite(values[1]) && isfinite(values[2]));

            Source->Position[0] = values[0];
            Source->Position[1] = values[1];
            Source->Position[2] = values[2];
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_VELOCITY:
            CHECKVAL(isfinite(values[0]) && isfinite(values[1]) && isfinite(values[2]));

            Source->Velocity[0] = values[0];
            Source->Velocity[1] = values[1];
            Source->Velocity[2] = values[2];
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_DIRECTION:
            CHECKVAL(isfinite(values[0]) && isfinite(values[1]) && isfinite(values[2]));

            Source->Direction[0] = values[0];
            Source->Direction[1] = values[1];
            Source->Direction[2] = values[2];
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_ORIENTATION:
            CHECKVAL(isfinite(values[0]) && isfinite(values[1]) && isfinite(values[2]) &&
                     isfinite(values[3]) && isfinite(values[4]) && isfinite(values[5]));

            Source->Orientation[0][0] = values[0];
            Source->Orientation[0][1] = values[1];
            Source->Orientation[0][2] = values[2];
            Source->Orientation[1][0] = values[3];
            Source->Orientation[1][1] = values[4];
            Source->Orientation[1][2] = values[5];
            DO_UPDATEPROPS();
            return AL_TRUE;


        case AL_SOURCE_RELATIVE:
        case AL_LOOPING:
        case AL_SOURCE_STATE:
        case AL_SOURCE_TYPE:
        case AL_DISTANCE_MODEL:
        case AL_DIRECT_FILTER_GAINHF_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAIN_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAINHF_AUTO:
        case AL_DIRECT_CHANNELS_SOFT:
        case AL_SOURCE_RESAMPLER_SOFT:
        case AL_SOURCE_SPATIALIZE_SOFT:
            ival = (ALint)values[0];
            return SetSourceiv(Source, Context, prop, &ival);

        case AL_BUFFERS_QUEUED:
        case AL_BUFFERS_PROCESSED:
            ival = (ALint)((ALuint)values[0]);
            return SetSourceiv(Source, Context, prop, &ival);

        case AL_BUFFER:
        case AL_DIRECT_FILTER:
        case AL_AUXILIARY_SEND_FILTER:
        case AL_SAMPLE_OFFSET_LATENCY_SOFT:
            break;
    }

    ERR("Unexpected property: 0x%04x\n", prop);
    SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_ENUM, AL_FALSE);
}

static ALboolean SetSourceiv(ALsource *Source, ALCcontext *Context, SourceProp prop, const ALint *values)
{
    ALCdevice *device = Context->Device;
    ALbuffer  *buffer = NULL;
    ALfilter  *filter = NULL;
    ALeffectslot *slot = NULL;
    ALbufferlistitem *oldlist;
    ALfloat fvals[6];

    switch(prop)
    {
        case AL_SOURCE_STATE:
        case AL_SOURCE_TYPE:
        case AL_BUFFERS_QUEUED:
        case AL_BUFFERS_PROCESSED:
        case AL_BYTE_LENGTH_SOFT:
        case AL_SAMPLE_LENGTH_SOFT:
        case AL_SEC_LENGTH_SOFT:
            /* Query only */
            SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_OPERATION, AL_FALSE);

        case AL_SOURCE_RELATIVE:
            CHECKVAL(*values == AL_FALSE || *values == AL_TRUE);

            Source->HeadRelative = (ALboolean)*values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_LOOPING:
            CHECKVAL(*values == AL_FALSE || *values == AL_TRUE);

            WriteLock(&Source->queue_lock);
            Source->Looping = (ALboolean)*values;
            if(IsPlayingOrPaused(Source))
            {
                ALvoice *voice = GetSourceVoice(Source, Context);
                if(voice)
                {
                    if(Source->Looping)
                        ATOMIC_STORE(&voice->loop_buffer, Source->queue, almemory_order_release);
                    else
                        ATOMIC_STORE(&voice->loop_buffer, NULL, almemory_order_release);

                    /* If the source is playing, wait for the current mix to finish
                     * to ensure it isn't currently looping back or reaching the
                     * end.
                     */
                    while((ATOMIC_LOAD(&device->MixCount, almemory_order_acquire)&1))
                        althrd_yield();
                }
            }
            WriteUnlock(&Source->queue_lock);
            return AL_TRUE;

        case AL_BUFFER:
            LockBuffersRead(device);
            if(!(*values == 0 || (buffer=LookupBuffer(device, *values)) != NULL))
            {
                UnlockBuffersRead(device);
                SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_VALUE, AL_FALSE);
            }

            WriteLock(&Source->queue_lock);
            {
                ALenum state = GetSourceState(Source, GetSourceVoice(Source, Context));
                if(state == AL_PLAYING || state == AL_PAUSED)
                {
                    WriteUnlock(&Source->queue_lock);
                    UnlockBuffersRead(device);
                    SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_OPERATION, AL_FALSE);
                }
            }

            oldlist = Source->queue;
            if(buffer != NULL)
            {
                /* Add the selected buffer to a one-item queue */
                ALbufferlistitem *newlist = al_calloc(DEF_ALIGN, sizeof(ALbufferlistitem));
                newlist->buffer = buffer;
                ATOMIC_INIT(&newlist->next, NULL);
                IncrementRef(&buffer->ref);

                /* Source is now Static */
                Source->SourceType = AL_STATIC;
                Source->queue = newlist;
            }
            else
            {
                /* Source is now Undetermined */
                Source->SourceType = AL_UNDETERMINED;
                Source->queue = NULL;
            }
            WriteUnlock(&Source->queue_lock);
            UnlockBuffersRead(device);

            /* Delete all elements in the previous queue */
            while(oldlist != NULL)
            {
                ALbufferlistitem *temp = oldlist;
                oldlist = ATOMIC_LOAD(&temp->next, almemory_order_relaxed);

                if(temp->buffer)
                    DecrementRef(&temp->buffer->ref);
                al_free(temp);
            }
            return AL_TRUE;

        case AL_SEC_OFFSET:
        case AL_SAMPLE_OFFSET:
        case AL_BYTE_OFFSET:
            CHECKVAL(*values >= 0);

            Source->OffsetType = prop;
            Source->Offset = *values;

            if(IsPlayingOrPaused(Source))
            {
                ALvoice *voice;

                ALCdevice_Lock(Context->Device);
                voice = GetSourceVoice(Source, Context);
                if(voice)
                {
                    WriteLock(&Source->queue_lock);
                    if(ApplyOffset(Source, voice) == AL_FALSE)
                    {
                        WriteUnlock(&Source->queue_lock);
                        ALCdevice_Unlock(Context->Device);
                        SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_VALUE, AL_FALSE);
                    }
                    WriteUnlock(&Source->queue_lock);
                }
                ALCdevice_Unlock(Context->Device);
            }
            return AL_TRUE;

        case AL_DIRECT_FILTER:
            LockFiltersRead(device);
            if(!(*values == 0 || (filter=LookupFilter(device, *values)) != NULL))
            {
                UnlockFiltersRead(device);
                SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_VALUE, AL_FALSE);
            }

            if(!filter)
            {
                Source->Direct.Gain = 1.0f;
                Source->Direct.GainHF = 1.0f;
                Source->Direct.HFReference = LOWPASSFREQREF;
                Source->Direct.GainLF = 1.0f;
                Source->Direct.LFReference = HIGHPASSFREQREF;
            }
            else
            {
                Source->Direct.Gain = filter->Gain;
                Source->Direct.GainHF = filter->GainHF;
                Source->Direct.HFReference = filter->HFReference;
                Source->Direct.GainLF = filter->GainLF;
                Source->Direct.LFReference = filter->LFReference;
            }
            UnlockFiltersRead(device);
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_DIRECT_FILTER_GAINHF_AUTO:
            CHECKVAL(*values == AL_FALSE || *values == AL_TRUE);

            Source->DryGainHFAuto = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_AUXILIARY_SEND_FILTER_GAIN_AUTO:
            CHECKVAL(*values == AL_FALSE || *values == AL_TRUE);

            Source->WetGainAuto = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_AUXILIARY_SEND_FILTER_GAINHF_AUTO:
            CHECKVAL(*values == AL_FALSE || *values == AL_TRUE);

            Source->WetGainHFAuto = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_DIRECT_CHANNELS_SOFT:
            CHECKVAL(*values == AL_FALSE || *values == AL_TRUE);

            Source->DirectChannels = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_DISTANCE_MODEL:
            CHECKVAL(*values == AL_NONE ||
                     *values == AL_INVERSE_DISTANCE ||
                     *values == AL_INVERSE_DISTANCE_CLAMPED ||
                     *values == AL_LINEAR_DISTANCE ||
                     *values == AL_LINEAR_DISTANCE_CLAMPED ||
                     *values == AL_EXPONENT_DISTANCE ||
                     *values == AL_EXPONENT_DISTANCE_CLAMPED);

            Source->DistanceModel = *values;
            if(Context->SourceDistanceModel)
                DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_SOURCE_RESAMPLER_SOFT:
            CHECKVAL(*values >= 0 && *values <= ResamplerMax);

            Source->Resampler = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;

        case AL_SOURCE_SPATIALIZE_SOFT:
            CHECKVAL(*values >= AL_FALSE && *values <= AL_AUTO_SOFT);

            Source->Spatialize = *values;
            DO_UPDATEPROPS();
            return AL_TRUE;


        case AL_AUXILIARY_SEND_FILTER:
            LockEffectSlotsRead(Context);
            LockFiltersRead(device);
            if(!((ALuint)values[1] < (ALuint)device->NumAuxSends &&
                 (values[0] == 0 || (slot=LookupEffectSlot(Context, values[0])) != NULL) &&
                 (values[2] == 0 || (filter=LookupFilter(device, values[2])) != NULL)))
            {
                UnlockFiltersRead(device);
                UnlockEffectSlotsRead(Context);
                SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_VALUE, AL_FALSE);
            }

            if(!filter)
            {
                /* Disable filter */
                Source->Send[values[1]].Gain = 1.0f;
                Source->Send[values[1]].GainHF = 1.0f;
                Source->Send[values[1]].HFReference = LOWPASSFREQREF;
                Source->Send[values[1]].GainLF = 1.0f;
                Source->Send[values[1]].LFReference = HIGHPASSFREQREF;
            }
            else
            {
                Source->Send[values[1]].Gain = filter->Gain;
                Source->Send[values[1]].GainHF = filter->GainHF;
                Source->Send[values[1]].HFReference = filter->HFReference;
                Source->Send[values[1]].GainLF = filter->GainLF;
                Source->Send[values[1]].LFReference = filter->LFReference;
            }
            UnlockFiltersRead(device);

            if(slot != Source->Send[values[1]].Slot && IsPlayingOrPaused(Source))
            {
                ALvoice *voice;
                /* Add refcount on the new slot, and release the previous slot */
                if(slot) IncrementRef(&slot->ref);
                if(Source->Send[values[1]].Slot)
                    DecrementRef(&Source->Send[values[1]].Slot->ref);
                Source->Send[values[1]].Slot = slot;

                /* We must force an update if the auxiliary slot changed on an
                 * active source, in case the slot is about to be deleted.
                 */
                if((voice=GetSourceVoice(Source, Context)) != NULL)
                    UpdateSourceProps(Source, voice, device->NumAuxSends);
                else
                    ATOMIC_FLAG_CLEAR(&Source->PropsClean, almemory_order_release);
            }
            else
            {
                if(slot) IncrementRef(&slot->ref);
                if(Source->Send[values[1]].Slot)
                    DecrementRef(&Source->Send[values[1]].Slot->ref);
                Source->Send[values[1]].Slot = slot;
                DO_UPDATEPROPS();
            }
            UnlockEffectSlotsRead(Context);

            return AL_TRUE;


        /* 1x float */
        case AL_CONE_INNER_ANGLE:
        case AL_CONE_OUTER_ANGLE:
        case AL_PITCH:
        case AL_GAIN:
        case AL_MIN_GAIN:
        case AL_MAX_GAIN:
        case AL_REFERENCE_DISTANCE:
        case AL_ROLLOFF_FACTOR:
        case AL_CONE_OUTER_GAIN:
        case AL_MAX_DISTANCE:
        case AL_DOPPLER_FACTOR:
        case AL_CONE_OUTER_GAINHF:
        case AL_AIR_ABSORPTION_FACTOR:
        case AL_ROOM_ROLLOFF_FACTOR:
        case AL_SOURCE_RADIUS:
            fvals[0] = (ALfloat)*values;
            return SetSourcefv(Source, Context, (int)prop, fvals);

        /* 3x float */
        case AL_POSITION:
        case AL_VELOCITY:
        case AL_DIRECTION:
            fvals[0] = (ALfloat)values[0];
            fvals[1] = (ALfloat)values[1];
            fvals[2] = (ALfloat)values[2];
            return SetSourcefv(Source, Context, (int)prop, fvals);

        /* 6x float */
        case AL_ORIENTATION:
            fvals[0] = (ALfloat)values[0];
            fvals[1] = (ALfloat)values[1];
            fvals[2] = (ALfloat)values[2];
            fvals[3] = (ALfloat)values[3];
            fvals[4] = (ALfloat)values[4];
            fvals[5] = (ALfloat)values[5];
            return SetSourcefv(Source, Context, (int)prop, fvals);

        case AL_SAMPLE_OFFSET_LATENCY_SOFT:
        case AL_SEC_OFFSET_LATENCY_SOFT:
        case AL_STEREO_ANGLES:
            break;
    }

    ERR("Unexpected property: 0x%04x\n", prop);
    SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_ENUM, AL_FALSE);
}

static ALboolean SetSourcei64v(ALsource *Source, ALCcontext *Context, SourceProp prop, const ALint64SOFT *values)
{
    ALfloat fvals[6];
    ALint   ivals[3];

    switch(prop)
    {
        case AL_SOURCE_TYPE:
        case AL_BUFFERS_QUEUED:
        case AL_BUFFERS_PROCESSED:
        case AL_SOURCE_STATE:
        case AL_SAMPLE_OFFSET_LATENCY_SOFT:
        case AL_BYTE_LENGTH_SOFT:
        case AL_SAMPLE_LENGTH_SOFT:
        case AL_SEC_LENGTH_SOFT:
            /* Query only */
            SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_OPERATION, AL_FALSE);


        /* 1x int */
        case AL_SOURCE_RELATIVE:
        case AL_LOOPING:
        case AL_SEC_OFFSET:
        case AL_SAMPLE_OFFSET:
        case AL_BYTE_OFFSET:
        case AL_DIRECT_FILTER_GAINHF_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAIN_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAINHF_AUTO:
        case AL_DIRECT_CHANNELS_SOFT:
        case AL_DISTANCE_MODEL:
        case AL_SOURCE_RESAMPLER_SOFT:
        case AL_SOURCE_SPATIALIZE_SOFT:
            CHECKVAL(*values <= INT_MAX && *values >= INT_MIN);

            ivals[0] = (ALint)*values;
            return SetSourceiv(Source, Context, (int)prop, ivals);

        /* 1x uint */
        case AL_BUFFER:
        case AL_DIRECT_FILTER:
            CHECKVAL(*values <= UINT_MAX && *values >= 0);

            ivals[0] = (ALuint)*values;
            return SetSourceiv(Source, Context, (int)prop, ivals);

        /* 3x uint */
        case AL_AUXILIARY_SEND_FILTER:
            CHECKVAL(values[0] <= UINT_MAX && values[0] >= 0 &&
                     values[1] <= UINT_MAX && values[1] >= 0 &&
                     values[2] <= UINT_MAX && values[2] >= 0);

            ivals[0] = (ALuint)values[0];
            ivals[1] = (ALuint)values[1];
            ivals[2] = (ALuint)values[2];
            return SetSourceiv(Source, Context, (int)prop, ivals);

        /* 1x float */
        case AL_CONE_INNER_ANGLE:
        case AL_CONE_OUTER_ANGLE:
        case AL_PITCH:
        case AL_GAIN:
        case AL_MIN_GAIN:
        case AL_MAX_GAIN:
        case AL_REFERENCE_DISTANCE:
        case AL_ROLLOFF_FACTOR:
        case AL_CONE_OUTER_GAIN:
        case AL_MAX_DISTANCE:
        case AL_DOPPLER_FACTOR:
        case AL_CONE_OUTER_GAINHF:
        case AL_AIR_ABSORPTION_FACTOR:
        case AL_ROOM_ROLLOFF_FACTOR:
        case AL_SOURCE_RADIUS:
            fvals[0] = (ALfloat)*values;
            return SetSourcefv(Source, Context, (int)prop, fvals);

        /* 3x float */
        case AL_POSITION:
        case AL_VELOCITY:
        case AL_DIRECTION:
            fvals[0] = (ALfloat)values[0];
            fvals[1] = (ALfloat)values[1];
            fvals[2] = (ALfloat)values[2];
            return SetSourcefv(Source, Context, (int)prop, fvals);

        /* 6x float */
        case AL_ORIENTATION:
            fvals[0] = (ALfloat)values[0];
            fvals[1] = (ALfloat)values[1];
            fvals[2] = (ALfloat)values[2];
            fvals[3] = (ALfloat)values[3];
            fvals[4] = (ALfloat)values[4];
            fvals[5] = (ALfloat)values[5];
            return SetSourcefv(Source, Context, (int)prop, fvals);

        case AL_SEC_OFFSET_LATENCY_SOFT:
        case AL_STEREO_ANGLES:
            break;
    }

    ERR("Unexpected property: 0x%04x\n", prop);
    SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_ENUM, AL_FALSE);
}

#undef CHECKVAL


static ALboolean GetSourcedv(ALsource *Source, ALCcontext *Context, SourceProp prop, ALdouble *values)
{
    ALCdevice *device = Context->Device;
    ALbufferlistitem *BufferList;
    ClockLatency clocktime;
    ALuint64 srcclock;
    ALint ivals[3];
    ALboolean err;

    switch(prop)
    {
        case AL_GAIN:
            *values = Source->Gain;
            return AL_TRUE;

        case AL_PITCH:
            *values = Source->Pitch;
            return AL_TRUE;

        case AL_MAX_DISTANCE:
            *values = Source->MaxDistance;
            return AL_TRUE;

        case AL_ROLLOFF_FACTOR:
            *values = Source->RolloffFactor;
            return AL_TRUE;

        case AL_REFERENCE_DISTANCE:
            *values = Source->RefDistance;
            return AL_TRUE;

        case AL_CONE_INNER_ANGLE:
            *values = Source->InnerAngle;
            return AL_TRUE;

        case AL_CONE_OUTER_ANGLE:
            *values = Source->OuterAngle;
            return AL_TRUE;

        case AL_MIN_GAIN:
            *values = Source->MinGain;
            return AL_TRUE;

        case AL_MAX_GAIN:
            *values = Source->MaxGain;
            return AL_TRUE;

        case AL_CONE_OUTER_GAIN:
            *values = Source->OuterGain;
            return AL_TRUE;

        case AL_SEC_OFFSET:
        case AL_SAMPLE_OFFSET:
        case AL_BYTE_OFFSET:
            *values = GetSourceOffset(Source, prop, Context);
            return AL_TRUE;

        case AL_CONE_OUTER_GAINHF:
            *values = Source->OuterGainHF;
            return AL_TRUE;

        case AL_AIR_ABSORPTION_FACTOR:
            *values = Source->AirAbsorptionFactor;
            return AL_TRUE;

        case AL_ROOM_ROLLOFF_FACTOR:
            *values = Source->RoomRolloffFactor;
            return AL_TRUE;

        case AL_DOPPLER_FACTOR:
            *values = Source->DopplerFactor;
            return AL_TRUE;

        case AL_SEC_LENGTH_SOFT:
            ReadLock(&Source->queue_lock);
            if(!(BufferList=Source->queue))
                *values = 0;
            else
            {
                ALint length = 0;
                ALsizei freq = 1;
                do {
                    ALbuffer *buffer = BufferList->buffer;
                    if(buffer && buffer->SampleLen > 0)
                    {
                        freq = buffer->Frequency;
                        length += buffer->SampleLen;
                    }
                    BufferList = ATOMIC_LOAD(&BufferList->next, almemory_order_relaxed);
                } while(BufferList != NULL);
                *values = (ALdouble)length / (ALdouble)freq;
            }
            ReadUnlock(&Source->queue_lock);
            return AL_TRUE;

        case AL_SOURCE_RADIUS:
            *values = Source->Radius;
            return AL_TRUE;

        case AL_STEREO_ANGLES:
            values[0] = Source->StereoPan[0];
            values[1] = Source->StereoPan[1];
            return AL_TRUE;

        case AL_SEC_OFFSET_LATENCY_SOFT:
            /* Get the source offset with the clock time first. Then get the
             * clock time with the device latency. Order is important.
             */
            values[0] = GetSourceSecOffset(Source, Context, &srcclock);
            clocktime = V0(device->Backend,getClockLatency)();
            if(srcclock == (ALuint64)clocktime.ClockTime)
                values[1] = (ALdouble)clocktime.Latency / 1000000000.0;
            else
            {
                /* If the clock time incremented, reduce the latency by that
                 * much since it's that much closer to the source offset it got
                 * earlier.
                 */
                ALuint64 diff = clocktime.ClockTime - srcclock;
                values[1] = (ALdouble)(clocktime.Latency - minu64(clocktime.Latency, diff)) /
                            1000000000.0;
            }
            return AL_TRUE;

        case AL_POSITION:
            values[0] = Source->Position[0];
            values[1] = Source->Position[1];
            values[2] = Source->Position[2];
            return AL_TRUE;

        case AL_VELOCITY:
            values[0] = Source->Velocity[0];
            values[1] = Source->Velocity[1];
            values[2] = Source->Velocity[2];
            return AL_TRUE;

        case AL_DIRECTION:
            values[0] = Source->Direction[0];
            values[1] = Source->Direction[1];
            values[2] = Source->Direction[2];
            return AL_TRUE;

        case AL_ORIENTATION:
            values[0] = Source->Orientation[0][0];
            values[1] = Source->Orientation[0][1];
            values[2] = Source->Orientation[0][2];
            values[3] = Source->Orientation[1][0];
            values[4] = Source->Orientation[1][1];
            values[5] = Source->Orientation[1][2];
            return AL_TRUE;

        /* 1x int */
        case AL_SOURCE_RELATIVE:
        case AL_LOOPING:
        case AL_SOURCE_STATE:
        case AL_BUFFERS_QUEUED:
        case AL_BUFFERS_PROCESSED:
        case AL_SOURCE_TYPE:
        case AL_DIRECT_FILTER_GAINHF_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAIN_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAINHF_AUTO:
        case AL_DIRECT_CHANNELS_SOFT:
        case AL_BYTE_LENGTH_SOFT:
        case AL_SAMPLE_LENGTH_SOFT:
        case AL_DISTANCE_MODEL:
        case AL_SOURCE_RESAMPLER_SOFT:
        case AL_SOURCE_SPATIALIZE_SOFT:
            if((err=GetSourceiv(Source, Context, (int)prop, ivals)) != AL_FALSE)
                *values = (ALdouble)ivals[0];
            return err;

        case AL_BUFFER:
        case AL_DIRECT_FILTER:
        case AL_AUXILIARY_SEND_FILTER:
        case AL_SAMPLE_OFFSET_LATENCY_SOFT:
            break;
    }

    ERR("Unexpected property: 0x%04x\n", prop);
    SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_ENUM, AL_FALSE);
}

static ALboolean GetSourceiv(ALsource *Source, ALCcontext *Context, SourceProp prop, ALint *values)
{
    ALbufferlistitem *BufferList;
    ALdouble dvals[6];
    ALboolean err;

    switch(prop)
    {
        case AL_SOURCE_RELATIVE:
            *values = Source->HeadRelative;
            return AL_TRUE;

        case AL_LOOPING:
            *values = Source->Looping;
            return AL_TRUE;

        case AL_BUFFER:
            ReadLock(&Source->queue_lock);
            BufferList = (Source->SourceType == AL_STATIC) ? Source->queue : NULL;
            *values = (BufferList && BufferList->buffer) ? BufferList->buffer->id : 0;
            ReadUnlock(&Source->queue_lock);
            return AL_TRUE;

        case AL_SOURCE_STATE:
            *values = GetSourceState(Source, GetSourceVoice(Source, Context));
            return AL_TRUE;

        case AL_BYTE_LENGTH_SOFT:
            ReadLock(&Source->queue_lock);
            if(!(BufferList=Source->queue))
                *values = 0;
            else
            {
                ALint length = 0;
                do {
                    ALbuffer *buffer = BufferList->buffer;
                    if(buffer && buffer->SampleLen > 0)
                    {
                        ALuint byte_align, sample_align;
                        if(buffer->OriginalType == UserFmtIMA4)
                        {
                            ALsizei align = (buffer->OriginalAlign-1)/2 + 4;
                            byte_align = align * ChannelsFromFmt(buffer->FmtChannels);
                            sample_align = buffer->OriginalAlign;
                        }
                        else if(buffer->OriginalType == UserFmtMSADPCM)
                        {
                            ALsizei align = (buffer->OriginalAlign-2)/2 + 7;
                            byte_align = align * ChannelsFromFmt(buffer->FmtChannels);
                            sample_align = buffer->OriginalAlign;
                        }
                        else
                        {
                            ALsizei align = buffer->OriginalAlign;
                            byte_align = align * ChannelsFromFmt(buffer->FmtChannels);
                            sample_align = buffer->OriginalAlign;
                        }

                        length += buffer->SampleLen / sample_align * byte_align;
                    }
                    BufferList = ATOMIC_LOAD(&BufferList->next, almemory_order_relaxed);
                } while(BufferList != NULL);
                *values = length;
            }
            ReadUnlock(&Source->queue_lock);
            return AL_TRUE;

        case AL_SAMPLE_LENGTH_SOFT:
            ReadLock(&Source->queue_lock);
            if(!(BufferList=Source->queue))
                *values = 0;
            else
            {
                ALint length = 0;
                do {
                    ALbuffer *buffer = BufferList->buffer;
                    if(buffer) length += buffer->SampleLen;
                    BufferList = ATOMIC_LOAD(&BufferList->next, almemory_order_relaxed);
                } while(BufferList != NULL);
                *values = length;
            }
            ReadUnlock(&Source->queue_lock);
            return AL_TRUE;

        case AL_BUFFERS_QUEUED:
            ReadLock(&Source->queue_lock);
            if(!(BufferList=Source->queue))
                *values = 0;
            else
            {
                ALsizei count = 0;
                do {
                    ++count;
                    BufferList = ATOMIC_LOAD(&BufferList->next, almemory_order_relaxed);
                } while(BufferList != NULL);
                *values = count;
            }
            ReadUnlock(&Source->queue_lock);
            return AL_TRUE;

        case AL_BUFFERS_PROCESSED:
            ReadLock(&Source->queue_lock);
            if(Source->Looping || Source->SourceType != AL_STREAMING)
            {
                /* Buffers on a looping source are in a perpetual state of
                 * PENDING, so don't report any as PROCESSED */
                *values = 0;
            }
            else
            {
                const ALbufferlistitem *BufferList = Source->queue;
                const ALbufferlistitem *Current = NULL;
                ALsizei played = 0;
                ALvoice *voice;

                if((voice=GetSourceVoice(Source, Context)) != NULL)
                    Current = ATOMIC_LOAD_SEQ(&voice->current_buffer);
                else if(ATOMIC_LOAD_SEQ(&Source->state) == AL_INITIAL)
                    Current = BufferList;

                while(BufferList && BufferList != Current)
                {
                    played++;
                    BufferList = ATOMIC_LOAD(&CONST_CAST(ALbufferlistitem*,BufferList)->next,
                                             almemory_order_relaxed);
                }
                *values = played;
            }
            ReadUnlock(&Source->queue_lock);
            return AL_TRUE;

        case AL_SOURCE_TYPE:
            *values = Source->SourceType;
            return AL_TRUE;

        case AL_DIRECT_FILTER_GAINHF_AUTO:
            *values = Source->DryGainHFAuto;
            return AL_TRUE;

        case AL_AUXILIARY_SEND_FILTER_GAIN_AUTO:
            *values = Source->WetGainAuto;
            return AL_TRUE;

        case AL_AUXILIARY_SEND_FILTER_GAINHF_AUTO:
            *values = Source->WetGainHFAuto;
            return AL_TRUE;

        case AL_DIRECT_CHANNELS_SOFT:
            *values = Source->DirectChannels;
            return AL_TRUE;

        case AL_DISTANCE_MODEL:
            *values = Source->DistanceModel;
            return AL_TRUE;

        case AL_SOURCE_RESAMPLER_SOFT:
            *values = Source->Resampler;
            return AL_TRUE;

        case AL_SOURCE_SPATIALIZE_SOFT:
            *values = Source->Spatialize;
            return AL_TRUE;

        /* 1x float/double */
        case AL_CONE_INNER_ANGLE:
        case AL_CONE_OUTER_ANGLE:
        case AL_PITCH:
        case AL_GAIN:
        case AL_MIN_GAIN:
        case AL_MAX_GAIN:
        case AL_REFERENCE_DISTANCE:
        case AL_ROLLOFF_FACTOR:
        case AL_CONE_OUTER_GAIN:
        case AL_MAX_DISTANCE:
        case AL_SEC_OFFSET:
        case AL_SAMPLE_OFFSET:
        case AL_BYTE_OFFSET:
        case AL_DOPPLER_FACTOR:
        case AL_AIR_ABSORPTION_FACTOR:
        case AL_ROOM_ROLLOFF_FACTOR:
        case AL_CONE_OUTER_GAINHF:
        case AL_SEC_LENGTH_SOFT:
        case AL_SOURCE_RADIUS:
            if((err=GetSourcedv(Source, Context, prop, dvals)) != AL_FALSE)
                *values = (ALint)dvals[0];
            return err;

        /* 3x float/double */
        case AL_POSITION:
        case AL_VELOCITY:
        case AL_DIRECTION:
            if((err=GetSourcedv(Source, Context, prop, dvals)) != AL_FALSE)
            {
                values[0] = (ALint)dvals[0];
                values[1] = (ALint)dvals[1];
                values[2] = (ALint)dvals[2];
            }
            return err;

        /* 6x float/double */
        case AL_ORIENTATION:
            if((err=GetSourcedv(Source, Context, prop, dvals)) != AL_FALSE)
            {
                values[0] = (ALint)dvals[0];
                values[1] = (ALint)dvals[1];
                values[2] = (ALint)dvals[2];
                values[3] = (ALint)dvals[3];
                values[4] = (ALint)dvals[4];
                values[5] = (ALint)dvals[5];
            }
            return err;

        case AL_SAMPLE_OFFSET_LATENCY_SOFT:
            break; /* i64 only */
        case AL_SEC_OFFSET_LATENCY_SOFT:
            break; /* Double only */
        case AL_STEREO_ANGLES:
            break; /* Float/double only */

        case AL_DIRECT_FILTER:
        case AL_AUXILIARY_SEND_FILTER:
            break; /* ??? */
    }

    ERR("Unexpected property: 0x%04x\n", prop);
    SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_ENUM, AL_FALSE);
}

static ALboolean GetSourcei64v(ALsource *Source, ALCcontext *Context, SourceProp prop, ALint64 *values)
{
    ALCdevice *device = Context->Device;
    ClockLatency clocktime;
    ALuint64 srcclock;
    ALdouble dvals[6];
    ALint ivals[3];
    ALboolean err;

    switch(prop)
    {
        case AL_SAMPLE_OFFSET_LATENCY_SOFT:
            /* Get the source offset with the clock time first. Then get the
             * clock time with the device latency. Order is important.
             */
            values[0] = GetSourceSampleOffset(Source, Context, &srcclock);
            clocktime = V0(device->Backend,getClockLatency)();
            if(srcclock == (ALuint64)clocktime.ClockTime)
                values[1] = clocktime.Latency;
            else
            {
                /* If the clock time incremented, reduce the latency by that
                 * much since it's that much closer to the source offset it got
                 * earlier.
                 */
                ALuint64 diff = clocktime.ClockTime - srcclock;
                values[1] = clocktime.Latency - minu64(clocktime.Latency, diff);
            }
            return AL_TRUE;

        /* 1x float/double */
        case AL_CONE_INNER_ANGLE:
        case AL_CONE_OUTER_ANGLE:
        case AL_PITCH:
        case AL_GAIN:
        case AL_MIN_GAIN:
        case AL_MAX_GAIN:
        case AL_REFERENCE_DISTANCE:
        case AL_ROLLOFF_FACTOR:
        case AL_CONE_OUTER_GAIN:
        case AL_MAX_DISTANCE:
        case AL_SEC_OFFSET:
        case AL_SAMPLE_OFFSET:
        case AL_BYTE_OFFSET:
        case AL_DOPPLER_FACTOR:
        case AL_AIR_ABSORPTION_FACTOR:
        case AL_ROOM_ROLLOFF_FACTOR:
        case AL_CONE_OUTER_GAINHF:
        case AL_SEC_LENGTH_SOFT:
        case AL_SOURCE_RADIUS:
            if((err=GetSourcedv(Source, Context, prop, dvals)) != AL_FALSE)
                *values = (ALint64)dvals[0];
            return err;

        /* 3x float/double */
        case AL_POSITION:
        case AL_VELOCITY:
        case AL_DIRECTION:
            if((err=GetSourcedv(Source, Context, prop, dvals)) != AL_FALSE)
            {
                values[0] = (ALint64)dvals[0];
                values[1] = (ALint64)dvals[1];
                values[2] = (ALint64)dvals[2];
            }
            return err;

        /* 6x float/double */
        case AL_ORIENTATION:
            if((err=GetSourcedv(Source, Context, prop, dvals)) != AL_FALSE)
            {
                values[0] = (ALint64)dvals[0];
                values[1] = (ALint64)dvals[1];
                values[2] = (ALint64)dvals[2];
                values[3] = (ALint64)dvals[3];
                values[4] = (ALint64)dvals[4];
                values[5] = (ALint64)dvals[5];
            }
            return err;

        /* 1x int */
        case AL_SOURCE_RELATIVE:
        case AL_LOOPING:
        case AL_SOURCE_STATE:
        case AL_BUFFERS_QUEUED:
        case AL_BUFFERS_PROCESSED:
        case AL_BYTE_LENGTH_SOFT:
        case AL_SAMPLE_LENGTH_SOFT:
        case AL_SOURCE_TYPE:
        case AL_DIRECT_FILTER_GAINHF_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAIN_AUTO:
        case AL_AUXILIARY_SEND_FILTER_GAINHF_AUTO:
        case AL_DIRECT_CHANNELS_SOFT:
        case AL_DISTANCE_MODEL:
        case AL_SOURCE_RESAMPLER_SOFT:
        case AL_SOURCE_SPATIALIZE_SOFT:
            if((err=GetSourceiv(Source, Context, prop, ivals)) != AL_FALSE)
                *values = ivals[0];
            return err;

        /* 1x uint */
        case AL_BUFFER:
        case AL_DIRECT_FILTER:
            if((err=GetSourceiv(Source, Context, prop, ivals)) != AL_FALSE)
                *values = (ALuint)ivals[0];
            return err;

        /* 3x uint */
        case AL_AUXILIARY_SEND_FILTER:
            if((err=GetSourceiv(Source, Context, prop, ivals)) != AL_FALSE)
            {
                values[0] = (ALuint)ivals[0];
                values[1] = (ALuint)ivals[1];
                values[2] = (ALuint)ivals[2];
            }
            return err;

        case AL_SEC_OFFSET_LATENCY_SOFT:
            break; /* Double only */
        case AL_STEREO_ANGLES:
            break; /* Float/double only */
    }

    ERR("Unexpected property: 0x%04x\n", prop);
    SET_ERROR_AND_RETURN_VALUE(Context, AL_INVALID_ENUM, AL_FALSE);
}


AL_API ALvoid AL_APIENTRY alGenSources(ALsizei n, ALuint *sources)
{
    ALCdevice *device;
    ALCcontext *context;
    ALsizei cur = 0;
    ALenum err;

    context = GetContextRef();
    if(!context) return;

    if(!(n >= 0))
        SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);
    device = context->Device;
    for(cur = 0;cur < n;cur++)
    {
        ALsource *source = al_calloc(16, sizeof(ALsource));
        if(!source)
        {
            alDeleteSources(cur, sources);
            SET_ERROR_AND_GOTO(context, AL_OUT_OF_MEMORY, done);
        }
        InitSourceParams(source, device->NumAuxSends);

        err = NewThunkEntry(&source->id);
        if(err == AL_NO_ERROR)
            err = InsertUIntMapEntry(&context->SourceMap, source->id, source);
        if(err != AL_NO_ERROR)
        {
            FreeThunkEntry(source->id);
            memset(source, 0, sizeof(ALsource));
            al_free(source);

            alDeleteSources(cur, sources);
            SET_ERROR_AND_GOTO(context, err, done);
        }

        sources[cur] = source->id;
    }

done:
    ALCcontext_DecRef(context);
}


AL_API ALvoid AL_APIENTRY alDeleteSources(ALsizei n, const ALuint *sources)
{
    ALCdevice *device;
    ALCcontext *context;
    ALsource *Source;
    ALsizei i;

    context = GetContextRef();
    if(!context) return;

    LockSourcesWrite(context);
    if(!(n >= 0))
        SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);

    /* Check that all Sources are valid */
    for(i = 0;i < n;i++)
    {
        if(LookupSource(context, sources[i]) == NULL)
            SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);
    }
    device = context->Device;
    for(i = 0;i < n;i++)
    {
        ALvoice *voice;

        if((Source=RemoveSource(context, sources[i])) == NULL)
            continue;
        FreeThunkEntry(Source->id);

        ALCdevice_Lock(device);
        if((voice=GetSourceVoice(Source, context)) != NULL)
        {
            ATOMIC_STORE(&voice->Source, NULL, almemory_order_relaxed);
            ATOMIC_STORE(&voice->Playing, false, almemory_order_release);
        }
        ALCdevice_Unlock(device);

        DeinitSource(Source, device->NumAuxSends);

        memset(Source, 0, sizeof(*Source));
        al_free(Source);
    }

done:
    UnlockSourcesWrite(context);
    ALCcontext_DecRef(context);
}


AL_API ALboolean AL_APIENTRY alIsSource(ALuint source)
{
    ALCcontext *context;
    ALboolean ret;

    context = GetContextRef();
    if(!context) return AL_FALSE;

    LockSourcesRead(context);
    ret = (LookupSource(context, source) ? AL_TRUE : AL_FALSE);
    UnlockSourcesRead(context);

    ALCcontext_DecRef(context);

    return ret;
}


AL_API ALvoid AL_APIENTRY alSourcef(ALuint source, ALenum param, ALfloat value)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    WriteLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!(FloatValsByProp(param) == 1))
        alSetError(Context, AL_INVALID_ENUM);
    else
        SetSourcefv(Source, Context, param, &value);
    UnlockSourcesRead(Context);
    WriteUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}

AL_API ALvoid AL_APIENTRY alSource3f(ALuint source, ALenum param, ALfloat value1, ALfloat value2, ALfloat value3)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    WriteLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!(FloatValsByProp(param) == 3))
        alSetError(Context, AL_INVALID_ENUM);
    else
    {
        ALfloat fvals[3] = { value1, value2, value3 };
        SetSourcefv(Source, Context, param, fvals);
    }
    UnlockSourcesRead(Context);
    WriteUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}

AL_API ALvoid AL_APIENTRY alSourcefv(ALuint source, ALenum param, const ALfloat *values)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    WriteLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!values)
        alSetError(Context, AL_INVALID_VALUE);
    else if(!(FloatValsByProp(param) > 0))
        alSetError(Context, AL_INVALID_ENUM);
    else
        SetSourcefv(Source, Context, param, values);
    UnlockSourcesRead(Context);
    WriteUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}


AL_API ALvoid AL_APIENTRY alSourcedSOFT(ALuint source, ALenum param, ALdouble value)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    WriteLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!(DoubleValsByProp(param) == 1))
        alSetError(Context, AL_INVALID_ENUM);
    else
    {
        ALfloat fval = (ALfloat)value;
        SetSourcefv(Source, Context, param, &fval);
    }
    UnlockSourcesRead(Context);
    WriteUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}

AL_API ALvoid AL_APIENTRY alSource3dSOFT(ALuint source, ALenum param, ALdouble value1, ALdouble value2, ALdouble value3)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    WriteLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!(DoubleValsByProp(param) == 3))
        alSetError(Context, AL_INVALID_ENUM);
    else
    {
        ALfloat fvals[3] = { (ALfloat)value1, (ALfloat)value2, (ALfloat)value3 };
        SetSourcefv(Source, Context, param, fvals);
    }
    UnlockSourcesRead(Context);
    WriteUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}

AL_API ALvoid AL_APIENTRY alSourcedvSOFT(ALuint source, ALenum param, const ALdouble *values)
{
    ALCcontext *Context;
    ALsource   *Source;
    ALint      count;

    Context = GetContextRef();
    if(!Context) return;

    WriteLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!values)
        alSetError(Context, AL_INVALID_VALUE);
    else if(!((count=DoubleValsByProp(param)) > 0 && count <= 6))
        alSetError(Context, AL_INVALID_ENUM);
    else
    {
        ALfloat fvals[6];
        ALint i;

        for(i = 0;i < count;i++)
            fvals[i] = (ALfloat)values[i];
        SetSourcefv(Source, Context, param, fvals);
    }
    UnlockSourcesRead(Context);
    WriteUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}


AL_API ALvoid AL_APIENTRY alSourcei(ALuint source, ALenum param, ALint value)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    WriteLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!(IntValsByProp(param) == 1))
        alSetError(Context, AL_INVALID_ENUM);
    else
        SetSourceiv(Source, Context, param, &value);
    UnlockSourcesRead(Context);
    WriteUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}

AL_API void AL_APIENTRY alSource3i(ALuint source, ALenum param, ALint value1, ALint value2, ALint value3)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    WriteLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!(IntValsByProp(param) == 3))
        alSetError(Context, AL_INVALID_ENUM);
    else
    {
        ALint ivals[3] = { value1, value2, value3 };
        SetSourceiv(Source, Context, param, ivals);
    }
    UnlockSourcesRead(Context);
    WriteUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}

AL_API void AL_APIENTRY alSourceiv(ALuint source, ALenum param, const ALint *values)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    WriteLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!values)
        alSetError(Context, AL_INVALID_VALUE);
    else if(!(IntValsByProp(param) > 0))
        alSetError(Context, AL_INVALID_ENUM);
    else
        SetSourceiv(Source, Context, param, values);
    UnlockSourcesRead(Context);
    WriteUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}


AL_API ALvoid AL_APIENTRY alSourcei64SOFT(ALuint source, ALenum param, ALint64SOFT value)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    WriteLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!(Int64ValsByProp(param) == 1))
        alSetError(Context, AL_INVALID_ENUM);
    else
        SetSourcei64v(Source, Context, param, &value);
    UnlockSourcesRead(Context);
    WriteUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}

AL_API void AL_APIENTRY alSource3i64SOFT(ALuint source, ALenum param, ALint64SOFT value1, ALint64SOFT value2, ALint64SOFT value3)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    WriteLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!(Int64ValsByProp(param) == 3))
        alSetError(Context, AL_INVALID_ENUM);
    else
    {
        ALint64SOFT i64vals[3] = { value1, value2, value3 };
        SetSourcei64v(Source, Context, param, i64vals);
    }
    UnlockSourcesRead(Context);
    WriteUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}

AL_API void AL_APIENTRY alSourcei64vSOFT(ALuint source, ALenum param, const ALint64SOFT *values)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    WriteLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!values)
        alSetError(Context, AL_INVALID_VALUE);
    else if(!(Int64ValsByProp(param) > 0))
        alSetError(Context, AL_INVALID_ENUM);
    else
        SetSourcei64v(Source, Context, param, values);
    UnlockSourcesRead(Context);
    WriteUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}


AL_API ALvoid AL_APIENTRY alGetSourcef(ALuint source, ALenum param, ALfloat *value)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    ReadLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!value)
        alSetError(Context, AL_INVALID_VALUE);
    else if(!(FloatValsByProp(param) == 1))
        alSetError(Context, AL_INVALID_ENUM);
    else
    {
        ALdouble dval;
        if(GetSourcedv(Source, Context, param, &dval))
            *value = (ALfloat)dval;
    }
    UnlockSourcesRead(Context);
    ReadUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}


AL_API ALvoid AL_APIENTRY alGetSource3f(ALuint source, ALenum param, ALfloat *value1, ALfloat *value2, ALfloat *value3)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    ReadLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!(value1 && value2 && value3))
        alSetError(Context, AL_INVALID_VALUE);
    else if(!(FloatValsByProp(param) == 3))
        alSetError(Context, AL_INVALID_ENUM);
    else
    {
        ALdouble dvals[3];
        if(GetSourcedv(Source, Context, param, dvals))
        {
            *value1 = (ALfloat)dvals[0];
            *value2 = (ALfloat)dvals[1];
            *value3 = (ALfloat)dvals[2];
        }
    }
    UnlockSourcesRead(Context);
    ReadUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}


AL_API ALvoid AL_APIENTRY alGetSourcefv(ALuint source, ALenum param, ALfloat *values)
{
    ALCcontext *Context;
    ALsource   *Source;
    ALint      count;

    Context = GetContextRef();
    if(!Context) return;

    ReadLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!values)
        alSetError(Context, AL_INVALID_VALUE);
    else if(!((count=FloatValsByProp(param)) > 0 && count <= 6))
        alSetError(Context, AL_INVALID_ENUM);
    else
    {
        ALdouble dvals[6];
        if(GetSourcedv(Source, Context, param, dvals))
        {
            ALint i;
            for(i = 0;i < count;i++)
                values[i] = (ALfloat)dvals[i];
        }
    }
    UnlockSourcesRead(Context);
    ReadUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}


AL_API void AL_APIENTRY alGetSourcedSOFT(ALuint source, ALenum param, ALdouble *value)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    ReadLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!value)
        alSetError(Context, AL_INVALID_VALUE);
    else if(!(DoubleValsByProp(param) == 1))
        alSetError(Context, AL_INVALID_ENUM);
    else
        GetSourcedv(Source, Context, param, value);
    UnlockSourcesRead(Context);
    ReadUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}

AL_API void AL_APIENTRY alGetSource3dSOFT(ALuint source, ALenum param, ALdouble *value1, ALdouble *value2, ALdouble *value3)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    ReadLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!(value1 && value2 && value3))
        alSetError(Context, AL_INVALID_VALUE);
    else if(!(DoubleValsByProp(param) == 3))
        alSetError(Context, AL_INVALID_ENUM);
    else
    {
        ALdouble dvals[3];
        if(GetSourcedv(Source, Context, param, dvals))
        {
            *value1 = dvals[0];
            *value2 = dvals[1];
            *value3 = dvals[2];
        }
    }
    UnlockSourcesRead(Context);
    ReadUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}

AL_API void AL_APIENTRY alGetSourcedvSOFT(ALuint source, ALenum param, ALdouble *values)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    ReadLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!values)
        alSetError(Context, AL_INVALID_VALUE);
    else if(!(DoubleValsByProp(param) > 0))
        alSetError(Context, AL_INVALID_ENUM);
    else
        GetSourcedv(Source, Context, param, values);
    UnlockSourcesRead(Context);
    ReadUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}


AL_API ALvoid AL_APIENTRY alGetSourcei(ALuint source, ALenum param, ALint *value)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    ReadLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!value)
        alSetError(Context, AL_INVALID_VALUE);
    else if(!(IntValsByProp(param) == 1))
        alSetError(Context, AL_INVALID_ENUM);
    else
        GetSourceiv(Source, Context, param, value);
    UnlockSourcesRead(Context);
    ReadUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}


AL_API void AL_APIENTRY alGetSource3i(ALuint source, ALenum param, ALint *value1, ALint *value2, ALint *value3)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    ReadLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!(value1 && value2 && value3))
        alSetError(Context, AL_INVALID_VALUE);
    else if(!(IntValsByProp(param) == 3))
        alSetError(Context, AL_INVALID_ENUM);
    else
    {
        ALint ivals[3];
        if(GetSourceiv(Source, Context, param, ivals))
        {
            *value1 = ivals[0];
            *value2 = ivals[1];
            *value3 = ivals[2];
        }
    }
    UnlockSourcesRead(Context);
    ReadUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}


AL_API void AL_APIENTRY alGetSourceiv(ALuint source, ALenum param, ALint *values)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    ReadLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!values)
        alSetError(Context, AL_INVALID_VALUE);
    else if(!(IntValsByProp(param) > 0))
        alSetError(Context, AL_INVALID_ENUM);
    else
        GetSourceiv(Source, Context, param, values);
    UnlockSourcesRead(Context);
    ReadUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}


AL_API void AL_APIENTRY alGetSourcei64SOFT(ALuint source, ALenum param, ALint64SOFT *value)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    ReadLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!value)
        alSetError(Context, AL_INVALID_VALUE);
    else if(!(Int64ValsByProp(param) == 1))
        alSetError(Context, AL_INVALID_ENUM);
    else
        GetSourcei64v(Source, Context, param, value);
    UnlockSourcesRead(Context);
    ReadUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}

AL_API void AL_APIENTRY alGetSource3i64SOFT(ALuint source, ALenum param, ALint64SOFT *value1, ALint64SOFT *value2, ALint64SOFT *value3)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    ReadLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!(value1 && value2 && value3))
        alSetError(Context, AL_INVALID_VALUE);
    else if(!(Int64ValsByProp(param) == 3))
        alSetError(Context, AL_INVALID_ENUM);
    else
    {
        ALint64 i64vals[3];
        if(GetSourcei64v(Source, Context, param, i64vals))
        {
            *value1 = i64vals[0];
            *value2 = i64vals[1];
            *value3 = i64vals[2];
        }
    }
    UnlockSourcesRead(Context);
    ReadUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}

AL_API void AL_APIENTRY alGetSourcei64vSOFT(ALuint source, ALenum param, ALint64SOFT *values)
{
    ALCcontext *Context;
    ALsource   *Source;

    Context = GetContextRef();
    if(!Context) return;

    ReadLock(&Context->PropLock);
    LockSourcesRead(Context);
    if((Source=LookupSource(Context, source)) == NULL)
        alSetError(Context, AL_INVALID_NAME);
    else if(!values)
        alSetError(Context, AL_INVALID_VALUE);
    else if(!(Int64ValsByProp(param) > 0))
        alSetError(Context, AL_INVALID_ENUM);
    else
        GetSourcei64v(Source, Context, param, values);
    UnlockSourcesRead(Context);
    ReadUnlock(&Context->PropLock);

    ALCcontext_DecRef(Context);
}


AL_API ALvoid AL_APIENTRY alSourcePlay(ALuint source)
{
    alSourcePlayv(1, &source);
}
AL_API ALvoid AL_APIENTRY alSourcePlayv(ALsizei n, const ALuint *sources)
{
    ALCcontext *context;
    ALCdevice *device;
    ALsource *source;
    ALvoice *voice;
    ALsizei i, j;

    context = GetContextRef();
    if(!context) return;

    LockSourcesRead(context);
    if(!(n >= 0))
        SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);
    for(i = 0;i < n;i++)
    {
        if(!LookupSource(context, sources[i]))
            SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);
    }

    device = context->Device;
    ALCdevice_Lock(device);
    /* If the device is disconnected, go right to stopped. */
    if(!device->Connected)
    {
        for(i = 0;i < n;i++)
        {
            source = LookupSource(context, sources[i]);
            ATOMIC_STORE(&source->state, AL_STOPPED, almemory_order_relaxed);
        }
        ALCdevice_Unlock(device);
        goto done;
    }

    while(n > context->MaxVoices-context->VoiceCount)
    {
        ALsizei newcount = context->MaxVoices << 1;
        if(context->MaxVoices >= newcount)
        {
            ALCdevice_Unlock(device);
            SET_ERROR_AND_GOTO(context, AL_OUT_OF_MEMORY, done);
        }
        AllocateVoices(context, newcount, device->NumAuxSends);
    }

    for(i = 0;i < n;i++)
    {
        ALbufferlistitem *BufferList;
        ALbuffer *buffer = NULL;
        bool start_fading = false;
        ALsizei s;

        source = LookupSource(context, sources[i]);
        WriteLock(&source->queue_lock);
        /* Check that there is a queue containing at least one valid, non zero
         * length Buffer.
         */
        BufferList = source->queue;
        while(BufferList)
        {
            if((buffer=BufferList->buffer) != NULL && buffer->SampleLen > 0)
                break;
            BufferList = ATOMIC_LOAD(&BufferList->next, almemory_order_relaxed);
        }

        /* If there's nothing to play, go right to stopped. */
        if(!BufferList)
        {
            /* NOTE: A source without any playable buffers should not have an
             * ALvoice since it shouldn't be in a playing or paused state. So
             * there's no need to look up its voice and clear the source.
             */
            ATOMIC_STORE(&source->state, AL_STOPPED, almemory_order_relaxed);
            source->OffsetType = AL_NONE;
            source->Offset = 0.0;
            goto finish_play;
        }

        voice = GetSourceVoice(source, context);
        switch(GetSourceState(source, voice))
        {
            case AL_PLAYING:
                assert(voice != NULL);
                /* A source that's already playing is restarted from the beginning. */
                ATOMIC_STORE(&voice->current_buffer, BufferList, almemory_order_relaxed);
                ATOMIC_STORE(&voice->position, 0, almemory_order_relaxed);
                ATOMIC_STORE(&voice->position_fraction, 0, almemory_order_release);
                goto finish_play;

            case AL_PAUSED:
                assert(voice != NULL);
                /* A source that's paused simply resumes. */
                ATOMIC_STORE(&voice->Playing, true, almemory_order_release);
                ATOMIC_STORE(&source->state, AL_PLAYING, almemory_order_release);
                goto finish_play;

            default:
                break;
        }

        /* Make sure this source isn't already active, and if not, look for an
         * unused voice to put it in.
         */
        assert(voice == NULL);
        for(j = 0;j < context->VoiceCount;j++)
        {
            if(ATOMIC_LOAD(&context->Voices[j]->Source, almemory_order_acquire) == NULL)
            {
                voice = context->Voices[j];
                break;
            }
        }
        if(voice == NULL)
            voice = context->Voices[context->VoiceCount++];
        ATOMIC_STORE(&voice->Playing, false, almemory_order_release);

        ATOMIC_FLAG_TEST_AND_SET(&source->PropsClean, almemory_order_acquire);
        UpdateSourceProps(source, voice, device->NumAuxSends);

        /* A source that's not playing or paused has any offset applied when it
         * starts playing.
         */
        if(source->Looping)
            ATOMIC_STORE(&voice->loop_buffer, source->queue, almemory_order_relaxed);
        else
            ATOMIC_STORE(&voice->loop_buffer, NULL, almemory_order_relaxed);
        ATOMIC_STORE(&voice->current_buffer, BufferList, almemory_order_relaxed);
        ATOMIC_STORE(&voice->position, 0, almemory_order_relaxed);
        ATOMIC_STORE(&voice->position_fraction, 0, almemory_order_relaxed);
        if(source->OffsetType != AL_NONE)
        {
            ApplyOffset(source, voice);
            start_fading = ATOMIC_LOAD(&voice->position, almemory_order_relaxed) != 0 ||
                ATOMIC_LOAD(&voice->position_fraction, almemory_order_relaxed) != 0 ||
                ATOMIC_LOAD(&voice->current_buffer, almemory_order_relaxed) != BufferList;
        }

        voice->NumChannels = ChannelsFromFmt(buffer->FmtChannels);
        voice->SampleSize  = BytesFromFmt(buffer->FmtType);

        /* Clear previous samples. */
        memset(voice->PrevSamples, 0, sizeof(voice->PrevSamples));

        /* Clear the stepping value so the mixer knows not to mix this until
         * the update gets applied.
         */
        voice->Step = 0;

        voice->Flags = start_fading ? VOICE_IS_FADING : 0;
        memset(voice->Direct.Params, 0, sizeof(voice->Direct.Params[0])*voice->NumChannels);
        for(s = 0;s < device->NumAuxSends;s++)
            memset(voice->Send[s].Params, 0, sizeof(voice->Send[s].Params[0])*voice->NumChannels);
        if(device->AvgSpeakerDist > 0.0f)
        {
            ALfloat w1 = SPEEDOFSOUNDMETRESPERSEC /
                        (device->AvgSpeakerDist * device->Frequency);
            for(j = 0;j < voice->NumChannels;j++)
            {
                NfcFilterCreate1(&voice->Direct.Params[j].NFCtrlFilter[0], 0.0f, w1);
                NfcFilterCreate2(&voice->Direct.Params[j].NFCtrlFilter[1], 0.0f, w1);
                NfcFilterCreate3(&voice->Direct.Params[j].NFCtrlFilter[2], 0.0f, w1);
            }
        }

        ATOMIC_STORE(&voice->Source, source, almemory_order_relaxed);
        ATOMIC_STORE(&voice->Playing, true, almemory_order_release);
        ATOMIC_STORE(&source->state, AL_PLAYING, almemory_order_release);
    finish_play:
        WriteUnlock(&source->queue_lock);
    }
    ALCdevice_Unlock(device);

done:
    UnlockSourcesRead(context);
    ALCcontext_DecRef(context);
}

AL_API ALvoid AL_APIENTRY alSourcePause(ALuint source)
{
    alSourcePausev(1, &source);
}
AL_API ALvoid AL_APIENTRY alSourcePausev(ALsizei n, const ALuint *sources)
{
    ALCcontext *context;
    ALCdevice *device;
    ALsource *source;
    ALvoice *voice;
    ALsizei i;

    context = GetContextRef();
    if(!context) return;

    LockSourcesRead(context);
    if(!(n >= 0))
        SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);
    for(i = 0;i < n;i++)
    {
        if(!LookupSource(context, sources[i]))
            SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);
    }

    device = context->Device;
    ALCdevice_Lock(device);
    for(i = 0;i < n;i++)
    {
        source = LookupSource(context, sources[i]);
        WriteLock(&source->queue_lock);
        if((voice=GetSourceVoice(source, context)) != NULL)
        {
            ATOMIC_STORE(&voice->Playing, false, almemory_order_release);
            while((ATOMIC_LOAD(&device->MixCount, almemory_order_acquire)&1))
                althrd_yield();
        }
        if(GetSourceState(source, voice) == AL_PLAYING)
            ATOMIC_STORE(&source->state, AL_PAUSED, almemory_order_release);
        WriteUnlock(&source->queue_lock);
    }
    ALCdevice_Unlock(device);

done:
    UnlockSourcesRead(context);
    ALCcontext_DecRef(context);
}

AL_API ALvoid AL_APIENTRY alSourceStop(ALuint source)
{
    alSourceStopv(1, &source);
}
AL_API ALvoid AL_APIENTRY alSourceStopv(ALsizei n, const ALuint *sources)
{
    ALCcontext *context;
    ALCdevice *device;
    ALsource *source;
    ALvoice *voice;
    ALsizei i;

    context = GetContextRef();
    if(!context) return;

    LockSourcesRead(context);
    if(!(n >= 0))
        SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);
    for(i = 0;i < n;i++)
    {
        if(!LookupSource(context, sources[i]))
            SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);
    }

    device = context->Device;
    ALCdevice_Lock(device);
    for(i = 0;i < n;i++)
    {
        source = LookupSource(context, sources[i]);
        WriteLock(&source->queue_lock);
        if((voice=GetSourceVoice(source, context)) != NULL)
        {
            ATOMIC_STORE(&voice->Source, NULL, almemory_order_relaxed);
            ATOMIC_STORE(&voice->Playing, false, almemory_order_release);
            while((ATOMIC_LOAD(&device->MixCount, almemory_order_acquire)&1))
                althrd_yield();
        }
        if(ATOMIC_LOAD(&source->state, almemory_order_acquire) != AL_INITIAL)
            ATOMIC_STORE(&source->state, AL_STOPPED, almemory_order_relaxed);
        source->OffsetType = AL_NONE;
        source->Offset = 0.0;
        WriteUnlock(&source->queue_lock);
    }
    ALCdevice_Unlock(device);

done:
    UnlockSourcesRead(context);
    ALCcontext_DecRef(context);
}

AL_API ALvoid AL_APIENTRY alSourceRewind(ALuint source)
{
    alSourceRewindv(1, &source);
}
AL_API ALvoid AL_APIENTRY alSourceRewindv(ALsizei n, const ALuint *sources)
{
    ALCcontext *context;
    ALCdevice *device;
    ALsource *source;
    ALvoice *voice;
    ALsizei i;

    context = GetContextRef();
    if(!context) return;

    LockSourcesRead(context);
    if(!(n >= 0))
        SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);
    for(i = 0;i < n;i++)
    {
        if(!LookupSource(context, sources[i]))
            SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);
    }

    device = context->Device;
    ALCdevice_Lock(device);
    for(i = 0;i < n;i++)
    {
        source = LookupSource(context, sources[i]);
        WriteLock(&source->queue_lock);
        if((voice=GetSourceVoice(source, context)) != NULL)
        {
            ATOMIC_STORE(&voice->Source, NULL, almemory_order_relaxed);
            ATOMIC_STORE(&voice->Playing, false, almemory_order_release);
            while((ATOMIC_LOAD(&device->MixCount, almemory_order_acquire)&1))
                althrd_yield();
        }
        if(ATOMIC_LOAD(&source->state, almemory_order_acquire) != AL_INITIAL)
            ATOMIC_STORE(&source->state, AL_INITIAL, almemory_order_relaxed);
        source->OffsetType = AL_NONE;
        source->Offset = 0.0;
        WriteUnlock(&source->queue_lock);
    }
    ALCdevice_Unlock(device);

done:
    UnlockSourcesRead(context);
    ALCcontext_DecRef(context);
}


AL_API ALvoid AL_APIENTRY alSourceQueueBuffers(ALuint src, ALsizei nb, const ALuint *buffers)
{
    ALCdevice *device;
    ALCcontext *context;
    ALsource *source;
    ALsizei i;
    ALbufferlistitem *BufferListStart;
    ALbufferlistitem *BufferList;
    ALbuffer *BufferFmt = NULL;

    if(nb == 0)
        return;

    context = GetContextRef();
    if(!context) return;

    device = context->Device;

    LockSourcesRead(context);
    if(!(nb >= 0))
        SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);
    if((source=LookupSource(context, src)) == NULL)
        SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);

    WriteLock(&source->queue_lock);
    if(source->SourceType == AL_STATIC)
    {
        WriteUnlock(&source->queue_lock);
        /* Can't queue on a Static Source */
        SET_ERROR_AND_GOTO(context, AL_INVALID_OPERATION, done);
    }

    /* Check for a valid Buffer, for its frequency and format */
    BufferList = source->queue;
    while(BufferList)
    {
        if(BufferList->buffer)
        {
            BufferFmt = BufferList->buffer;
            break;
        }
        BufferList = ATOMIC_LOAD(&BufferList->next, almemory_order_relaxed);
    }

    LockBuffersRead(device);
    BufferListStart = NULL;
    BufferList = NULL;
    for(i = 0;i < nb;i++)
    {
        ALbuffer *buffer = NULL;
        if(buffers[i] && (buffer=LookupBuffer(device, buffers[i])) == NULL)
        {
            WriteUnlock(&source->queue_lock);
            SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, buffer_error);
        }

        if(!BufferListStart)
        {
            BufferListStart = al_calloc(DEF_ALIGN, sizeof(ALbufferlistitem));
            BufferList = BufferListStart;
        }
        else
        {
            ALbufferlistitem *item = al_calloc(DEF_ALIGN, sizeof(ALbufferlistitem));
            ATOMIC_STORE(&BufferList->next, item, almemory_order_relaxed);
            BufferList = item;
        }
        BufferList->buffer = buffer;
        ATOMIC_INIT(&BufferList->next, NULL);
        if(!buffer) continue;

        /* Hold a read lock on each buffer being queued while checking all
         * provided buffers. This is done so other threads don't see an extra
         * reference on some buffers if this operation ends up failing. */
        ReadLock(&buffer->lock);
        IncrementRef(&buffer->ref);

        if(BufferFmt == NULL)
            BufferFmt = buffer;
        else if(BufferFmt->Frequency != buffer->Frequency ||
                BufferFmt->OriginalChannels != buffer->OriginalChannels ||
                BufferFmt->OriginalType != buffer->OriginalType)
        {
            WriteUnlock(&source->queue_lock);
            SET_ERROR_AND_GOTO(context, AL_INVALID_OPERATION, buffer_error);

        buffer_error:
            /* A buffer failed (invalid ID or format), so unlock and release
             * each buffer we had. */
            while(BufferListStart)
            {
                ALbufferlistitem *next = ATOMIC_LOAD(&BufferListStart->next,
                                                     almemory_order_relaxed);
                if((buffer=BufferListStart->buffer) != NULL)
                {
                    DecrementRef(&buffer->ref);
                    ReadUnlock(&buffer->lock);
                }
                al_free(BufferListStart);
                BufferListStart = next;
            }
            UnlockBuffersRead(device);
            goto done;
        }
    }
    /* All buffers good, unlock them now. */
    BufferList = BufferListStart;
    while(BufferList != NULL)
    {
        ALbuffer *buffer = BufferList->buffer;
        if(buffer) ReadUnlock(&buffer->lock);
        BufferList = ATOMIC_LOAD(&BufferList->next, almemory_order_relaxed);
    }
    UnlockBuffersRead(device);

    /* Source is now streaming */
    source->SourceType = AL_STREAMING;

    if(!(BufferList=source->queue))
        source->queue = BufferListStart;
    else
    {
        ALbufferlistitem *next;
        while((next=ATOMIC_LOAD(&BufferList->next, almemory_order_relaxed)) != NULL)
            BufferList = next;
        ATOMIC_STORE(&BufferList->next, BufferListStart, almemory_order_release);
    }
    WriteUnlock(&source->queue_lock);

done:
    UnlockSourcesRead(context);
    ALCcontext_DecRef(context);
}

AL_API ALvoid AL_APIENTRY alSourceUnqueueBuffers(ALuint src, ALsizei nb, ALuint *buffers)
{
    ALCcontext *context;
    ALsource *source;
    ALbufferlistitem *OldHead;
    ALbufferlistitem *OldTail;
    ALbufferlistitem *Current;
    ALvoice *voice;
    ALsizei i = 0;

    context = GetContextRef();
    if(!context) return;

    LockSourcesRead(context);
    if(!(nb >= 0))
        SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);

    if((source=LookupSource(context, src)) == NULL)
        SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);

    /* Nothing to unqueue. */
    if(nb == 0) goto done;

    WriteLock(&source->queue_lock);
    if(source->Looping || source->SourceType != AL_STREAMING)
    {
        WriteUnlock(&source->queue_lock);
        /* Trying to unqueue buffers on a looping or non-streaming source. */
        SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);
    }

    /* Find the new buffer queue head */
    OldTail = source->queue;
    Current = NULL;
    if((voice=GetSourceVoice(source, context)) != NULL)
        Current = ATOMIC_LOAD_SEQ(&voice->current_buffer);
    else if(ATOMIC_LOAD_SEQ(&source->state) == AL_INITIAL)
        Current = OldTail;
    if(OldTail != Current)
    {
        for(i = 1;i < nb;i++)
        {
            ALbufferlistitem *next = ATOMIC_LOAD(&OldTail->next, almemory_order_relaxed);
            if(!next || next == Current) break;
            OldTail = next;
        }
    }
    if(i != nb)
    {
        WriteUnlock(&source->queue_lock);
        /* Trying to unqueue pending buffers. */
        SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);
    }

    /* Swap it, and cut the new head from the old. */
    OldHead = source->queue;
    source->queue = ATOMIC_EXCHANGE_PTR(&OldTail->next, NULL, almemory_order_acq_rel);
    WriteUnlock(&source->queue_lock);

    while(OldHead != NULL)
    {
        ALbufferlistitem *next = ATOMIC_LOAD(&OldHead->next, almemory_order_relaxed);
        ALbuffer *buffer = OldHead->buffer;

        if(!buffer)
            *(buffers++) = 0;
        else
        {
            *(buffers++) = buffer->id;
            DecrementRef(&buffer->ref);
        }

        al_free(OldHead);
        OldHead = next;
    }

done:
    UnlockSourcesRead(context);
    ALCcontext_DecRef(context);
}


static void InitSourceParams(ALsource *Source, ALsizei num_sends)
{
    ALsizei i;

    RWLockInit(&Source->queue_lock);

    Source->InnerAngle = 360.0f;
    Source->OuterAngle = 360.0f;
    Source->Pitch = 1.0f;
    Source->Position[0] = 0.0f;
    Source->Position[1] = 0.0f;
    Source->Position[2] = 0.0f;
    Source->Velocity[0] = 0.0f;
    Source->Velocity[1] = 0.0f;
    Source->Velocity[2] = 0.0f;
    Source->Direction[0] = 0.0f;
    Source->Direction[1] = 0.0f;
    Source->Direction[2] = 0.0f;
    Source->Orientation[0][0] =  0.0f;
    Source->Orientation[0][1] =  0.0f;
    Source->Orientation[0][2] = -1.0f;
    Source->Orientation[1][0] =  0.0f;
    Source->Orientation[1][1] =  1.0f;
    Source->Orientation[1][2] =  0.0f;
    Source->RefDistance = 1.0f;
    Source->MaxDistance = FLT_MAX;
    Source->RolloffFactor = 1.0f;
    Source->Gain = 1.0f;
    Source->MinGain = 0.0f;
    Source->MaxGain = 1.0f;
    Source->OuterGain = 0.0f;
    Source->OuterGainHF = 1.0f;

    Source->DryGainHFAuto = AL_TRUE;
    Source->WetGainAuto = AL_TRUE;
    Source->WetGainHFAuto = AL_TRUE;
    Source->AirAbsorptionFactor = 0.0f;
    Source->RoomRolloffFactor = 0.0f;
    Source->DopplerFactor = 1.0f;
    Source->HeadRelative = AL_FALSE;
    Source->Looping = AL_FALSE;
    Source->DistanceModel = DefaultDistanceModel;
    Source->Resampler = ResamplerDefault;
    Source->DirectChannels = AL_FALSE;
    Source->Spatialize = SpatializeAuto;

    Source->StereoPan[0] = DEG2RAD( 30.0f);
    Source->StereoPan[1] = DEG2RAD(-30.0f);

    Source->Radius = 0.0f;

    Source->Direct.Gain = 1.0f;
    Source->Direct.GainHF = 1.0f;
    Source->Direct.HFReference = LOWPASSFREQREF;
    Source->Direct.GainLF = 1.0f;
    Source->Direct.LFReference = HIGHPASSFREQREF;
    Source->Send = al_calloc(16, num_sends*sizeof(Source->Send[0]));
    for(i = 0;i < num_sends;i++)
    {
        Source->Send[i].Slot = NULL;
        Source->Send[i].Gain = 1.0f;
        Source->Send[i].GainHF = 1.0f;
        Source->Send[i].HFReference = LOWPASSFREQREF;
        Source->Send[i].GainLF = 1.0f;
        Source->Send[i].LFReference = HIGHPASSFREQREF;
    }

    Source->Offset = 0.0;
    Source->OffsetType = AL_NONE;
    Source->SourceType = AL_UNDETERMINED;
    ATOMIC_INIT(&Source->state, AL_INITIAL);

    Source->queue = NULL;

    /* No way to do an 'init' here, so just test+set with relaxed ordering and
     * ignore the test.
     */
    ATOMIC_FLAG_TEST_AND_SET(&Source->PropsClean, almemory_order_relaxed);
}

static void DeinitSource(ALsource *source, ALsizei num_sends)
{
    ALbufferlistitem *BufferList;
    ALsizei i;

    BufferList = source->queue;
    while(BufferList != NULL)
    {
        ALbufferlistitem *next = ATOMIC_LOAD(&BufferList->next, almemory_order_relaxed);
        if(BufferList->buffer != NULL)
            DecrementRef(&BufferList->buffer->ref);
        al_free(BufferList);
        BufferList = next;
    }
    source->queue = NULL;

    if(source->Send)
    {
        for(i = 0;i < num_sends;i++)
        {
            if(source->Send[i].Slot)
                DecrementRef(&source->Send[i].Slot->ref);
            source->Send[i].Slot = NULL;
        }
        al_free(source->Send);
        source->Send = NULL;
    }
}

static void UpdateSourceProps(ALsource *source, ALvoice *voice, ALsizei num_sends)
{
    struct ALvoiceProps *props;
    ALsizei i;

    /* Get an unused property container, or allocate a new one as needed. */
    props = ATOMIC_LOAD(&voice->FreeList, almemory_order_acquire);
    if(!props)
        props = al_calloc(16, FAM_SIZE(struct ALvoiceProps, Send, num_sends));
    else
    {
        struct ALvoiceProps *next;
        do {
            next = ATOMIC_LOAD(&props->next, almemory_order_relaxed);
        } while(ATOMIC_COMPARE_EXCHANGE_PTR_WEAK(&voice->FreeList, &props, next,
                almemory_order_acq_rel, almemory_order_acquire) == 0);
    }

    /* Copy in current property values. */
    props->Pitch = source->Pitch;
    props->Gain = source->Gain;
    props->OuterGain = source->OuterGain;
    props->MinGain = source->MinGain;
    props->MaxGain = source->MaxGain;
    props->InnerAngle = source->InnerAngle;
    props->OuterAngle = source->OuterAngle;
    props->RefDistance = source->RefDistance;
    props->MaxDistance = source->MaxDistance;
    props->RolloffFactor = source->RolloffFactor;
    for(i = 0;i < 3;i++)
        props->Position[i] = source->Position[i];
    for(i = 0;i < 3;i++)
        props->Velocity[i] = source->Velocity[i];
    for(i = 0;i < 3;i++)
        props->Direction[i] = source->Direction[i];
    for(i = 0;i < 2;i++)
    {
        ALsizei j;
        for(j = 0;j < 3;j++)
            props->Orientation[i][j] = source->Orientation[i][j];
    }
    props->HeadRelative = source->HeadRelative;
    props->DistanceModel = source->DistanceModel;
    props->Resampler = source->Resampler;
    props->DirectChannels = source->DirectChannels;
    props->SpatializeMode = source->Spatialize;

    props->DryGainHFAuto = source->DryGainHFAuto;
    props->WetGainAuto = source->WetGainAuto;
    props->WetGainHFAuto = source->WetGainHFAuto;
    props->OuterGainHF = source->OuterGainHF;

    props->AirAbsorptionFactor = source->AirAbsorptionFactor;
    props->RoomRolloffFactor = source->RoomRolloffFactor;
    props->DopplerFactor = source->DopplerFactor;

    props->StereoPan[0] = source->StereoPan[0];
    props->StereoPan[1] = source->StereoPan[1];

    props->Radius = source->Radius;

    props->Direct.Gain = source->Direct.Gain;
    props->Direct.GainHF = source->Direct.GainHF;
    props->Direct.HFReference = source->Direct.HFReference;
    props->Direct.GainLF = source->Direct.GainLF;
    props->Direct.LFReference = source->Direct.LFReference;

    for(i = 0;i < num_sends;i++)
    {
        props->Send[i].Slot = source->Send[i].Slot;
        props->Send[i].Gain = source->Send[i].Gain;
        props->Send[i].GainHF = source->Send[i].GainHF;
        props->Send[i].HFReference = source->Send[i].HFReference;
        props->Send[i].GainLF = source->Send[i].GainLF;
        props->Send[i].LFReference = source->Send[i].LFReference;
    }

    /* Set the new container for updating internal parameters. */
    props = ATOMIC_EXCHANGE_PTR(&voice->Update, props, almemory_order_acq_rel);
    if(props)
    {
        /* If there was an unused update container, put it back in the
         * freelist.
         */
        ATOMIC_REPLACE_HEAD(struct ALvoiceProps*, &voice->FreeList, props);
    }
}

void UpdateAllSourceProps(ALCcontext *context)
{
    ALsizei num_sends = context->Device->NumAuxSends;
    ALsizei pos;

    for(pos = 0;pos < context->VoiceCount;pos++)
    {
        ALvoice *voice = context->Voices[pos];
        ALsource *source = ATOMIC_LOAD(&voice->Source, almemory_order_acquire);
        if(source && !ATOMIC_FLAG_TEST_AND_SET(&source->PropsClean, almemory_order_acq_rel))
            UpdateSourceProps(source, voice, num_sends);
    }
}


/* GetSourceSampleOffset
 *
 * Gets the current read offset for the given Source, in 32.32 fixed-point
 * samples. The offset is relative to the start of the queue (not the start of
 * the current buffer).
 */
static ALint64 GetSourceSampleOffset(ALsource *Source, ALCcontext *context, ALuint64 *clocktime)
{
    ALCdevice *device = context->Device;
    const ALbufferlistitem *Current;
    ALuint64 readPos;
    ALuint refcount;
    ALvoice *voice;

    ReadLock(&Source->queue_lock);
    do {
        Current = NULL;
        readPos = 0;
        while(((refcount=ATOMIC_LOAD(&device->MixCount, almemory_order_acquire))&1))
            althrd_yield();
        *clocktime = GetDeviceClockTime(device);

        voice = GetSourceVoice(Source, context);
        if(voice)
        {
            Current = ATOMIC_LOAD(&voice->current_buffer, almemory_order_relaxed);

            readPos  = (ALuint64)ATOMIC_LOAD(&voice->position, almemory_order_relaxed) << 32;
            readPos |= (ALuint64)ATOMIC_LOAD(&voice->position_fraction, almemory_order_relaxed) <<
                       (32-FRACTIONBITS);
        }
        ATOMIC_THREAD_FENCE(almemory_order_acquire);
    } while(refcount != ATOMIC_LOAD(&device->MixCount, almemory_order_relaxed));

    if(voice)
    {
        const ALbufferlistitem *BufferList = Source->queue;
        while(BufferList && BufferList != Current)
        {
            if(BufferList->buffer)
                readPos += (ALuint64)BufferList->buffer->SampleLen << 32;
            BufferList = ATOMIC_LOAD(&CONST_CAST(ALbufferlistitem*,BufferList)->next,
                                     almemory_order_relaxed);
        }
        readPos = minu64(readPos, U64(0x7fffffffffffffff));
    }

    ReadUnlock(&Source->queue_lock);
    return (ALint64)readPos;
}

/* GetSourceSecOffset
 *
 * Gets the current read offset for the given Source, in seconds. The offset is
 * relative to the start of the queue (not the start of the current buffer).
 */
static ALdouble GetSourceSecOffset(ALsource *Source, ALCcontext *context, ALuint64 *clocktime)
{
    ALCdevice *device = context->Device;
    const ALbufferlistitem *Current;
    ALuint64 readPos;
    ALuint refcount;
    ALdouble offset;
    ALvoice *voice;

    ReadLock(&Source->queue_lock);
    do {
        Current = NULL;
        readPos = 0;
        while(((refcount=ATOMIC_LOAD(&device->MixCount, almemory_order_acquire))&1))
            althrd_yield();
        *clocktime = GetDeviceClockTime(device);

        voice = GetSourceVoice(Source, context);
        if(voice)
        {
            Current = ATOMIC_LOAD(&voice->current_buffer, almemory_order_relaxed);

            readPos  = (ALuint64)ATOMIC_LOAD(&voice->position, almemory_order_relaxed) <<
                       FRACTIONBITS;
            readPos |= ATOMIC_LOAD(&voice->position_fraction, almemory_order_relaxed);
        }
        ATOMIC_THREAD_FENCE(almemory_order_acquire);
    } while(refcount != ATOMIC_LOAD(&device->MixCount, almemory_order_relaxed));

    offset = 0.0;
    if(voice)
    {
        const ALbufferlistitem *BufferList = Source->queue;
        const ALbuffer *BufferFmt = NULL;
        while(BufferList && BufferList != Current)
        {
            const ALbuffer *buffer = BufferList->buffer;
            if(buffer != NULL)
            {
                if(!BufferFmt) BufferFmt = buffer;
                readPos += (ALuint64)buffer->SampleLen << FRACTIONBITS;
            }
            BufferList = ATOMIC_LOAD(&CONST_CAST(ALbufferlistitem*,BufferList)->next,
                                     almemory_order_relaxed);
        }

        while(BufferList && !BufferFmt)
        {
            BufferFmt = BufferList->buffer;
            BufferList = ATOMIC_LOAD(&CONST_CAST(ALbufferlistitem*,BufferList)->next,
                                     almemory_order_relaxed);
        }
        assert(BufferFmt != NULL);

        offset = (ALdouble)readPos / (ALdouble)FRACTIONONE /
                 (ALdouble)BufferFmt->Frequency;
    }

    ReadUnlock(&Source->queue_lock);
    return offset;
}

/* GetSourceOffset
 *
 * Gets the current read offset for the given Source, in the appropriate format
 * (Bytes, Samples or Seconds). The offset is relative to the start of the
 * queue (not the start of the current buffer).
 */
static ALdouble GetSourceOffset(ALsource *Source, ALenum name, ALCcontext *context)
{
    ALCdevice *device = context->Device;
    const ALbufferlistitem *Current;
    ALuint readPos;
    ALsizei readPosFrac;
    ALuint refcount;
    ALdouble offset;
    ALvoice *voice;

    ReadLock(&Source->queue_lock);
    do {
        Current = NULL;
        readPos = readPosFrac = 0;
        while(((refcount=ATOMIC_LOAD(&device->MixCount, almemory_order_acquire))&1))
            althrd_yield();
        voice = GetSourceVoice(Source, context);
        if(voice)
        {
            Current = ATOMIC_LOAD(&voice->current_buffer, almemory_order_relaxed);

            readPos = ATOMIC_LOAD(&voice->position, almemory_order_relaxed);
            readPosFrac = ATOMIC_LOAD(&voice->position_fraction, almemory_order_relaxed);
        }
        ATOMIC_THREAD_FENCE(almemory_order_acquire);
    } while(refcount != ATOMIC_LOAD(&device->MixCount, almemory_order_relaxed));

    offset = 0.0;
    if(voice)
    {
        const ALbufferlistitem *BufferList = Source->queue;
        const ALbuffer *BufferFmt = NULL;
        ALboolean readFin = AL_FALSE;
        ALuint totalBufferLen = 0;

        while(BufferList != NULL)
        {
            const ALbuffer *buffer;
            readFin = readFin || (BufferList == Current);
            if((buffer=BufferList->buffer) != NULL)
            {
                if(!BufferFmt) BufferFmt = buffer;
                totalBufferLen += buffer->SampleLen;
                if(!readFin) readPos += buffer->SampleLen;
            }
            BufferList = ATOMIC_LOAD(&CONST_CAST(ALbufferlistitem*,BufferList)->next,
                                     almemory_order_relaxed);
        }
        assert(BufferFmt != NULL);

        if(Source->Looping)
            readPos %= totalBufferLen;
        else
        {
            /* Wrap back to 0 */
            if(readPos >= totalBufferLen)
                readPos = readPosFrac = 0;
        }

        offset = 0.0;
        switch(name)
        {
            case AL_SEC_OFFSET:
                offset = (readPos + (ALdouble)readPosFrac/FRACTIONONE) / BufferFmt->Frequency;
                break;

            case AL_SAMPLE_OFFSET:
                offset = readPos + (ALdouble)readPosFrac/FRACTIONONE;
                break;

            case AL_BYTE_OFFSET:
                if(BufferFmt->OriginalType == UserFmtIMA4)
                {
                    ALsizei align = (BufferFmt->OriginalAlign-1)/2 + 4;
                    ALuint BlockSize = align * ChannelsFromFmt(BufferFmt->FmtChannels);
                    ALuint FrameBlockSize = BufferFmt->OriginalAlign;

                    /* Round down to nearest ADPCM block */
                    offset = (ALdouble)(readPos / FrameBlockSize * BlockSize);
                }
                else if(BufferFmt->OriginalType == UserFmtMSADPCM)
                {
                    ALsizei align = (BufferFmt->OriginalAlign-2)/2 + 7;
                    ALuint BlockSize = align * ChannelsFromFmt(BufferFmt->FmtChannels);
                    ALuint FrameBlockSize = BufferFmt->OriginalAlign;

                    /* Round down to nearest ADPCM block */
                    offset = (ALdouble)(readPos / FrameBlockSize * BlockSize);
                }
                else
                {
                    ALuint FrameSize = FrameSizeFromUserFmt(BufferFmt->OriginalChannels,
                                                            BufferFmt->OriginalType);
                    offset = (ALdouble)(readPos * FrameSize);
                }
                break;
        }
    }

    ReadUnlock(&Source->queue_lock);
    return offset;
}


/* ApplyOffset
 *
 * Apply the stored playback offset to the Source. This function will update
 * the number of buffers "played" given the stored offset.
 */
static ALboolean ApplyOffset(ALsource *Source, ALvoice *voice)
{
    ALbufferlistitem *BufferList;
    const ALbuffer *Buffer;
    ALuint bufferLen, totalBufferLen;
    ALuint offset = 0;
    ALsizei frac = 0;

    /* Get sample frame offset */
    if(!GetSampleOffset(Source, &offset, &frac))
        return AL_FALSE;

    totalBufferLen = 0;
    BufferList = Source->queue;
    while(BufferList && totalBufferLen <= offset)
    {
        Buffer = BufferList->buffer;
        bufferLen = Buffer ? Buffer->SampleLen : 0;

        if(bufferLen > offset-totalBufferLen)
        {
            /* Offset is in this buffer */
            ATOMIC_STORE(&voice->position, offset - totalBufferLen, almemory_order_relaxed);
            ATOMIC_STORE(&voice->position_fraction, frac, almemory_order_relaxed);
            ATOMIC_STORE(&voice->current_buffer, BufferList, almemory_order_release);
            return AL_TRUE;
        }

        totalBufferLen += bufferLen;

        BufferList = ATOMIC_LOAD(&BufferList->next, almemory_order_relaxed);
    }

    /* Offset is out of range of the queue */
    return AL_FALSE;
}


/* GetSampleOffset
 *
 * Retrieves the sample offset into the Source's queue (from the Sample, Byte
 * or Second offset supplied by the application). This takes into account the
 * fact that the buffer format may have been modifed since.
 */
static ALboolean GetSampleOffset(ALsource *Source, ALuint *offset, ALsizei *frac)
{
    const ALbuffer *BufferFmt = NULL;
    const ALbufferlistitem *BufferList;
    ALdouble dbloff, dblfrac;

    /* Find the first valid Buffer in the Queue */
    BufferList = Source->queue;
    while(BufferList)
    {
        if((BufferFmt=BufferList->buffer) != NULL)
            break;
        BufferList = ATOMIC_LOAD(&CONST_CAST(ALbufferlistitem*,BufferList)->next,
                                 almemory_order_relaxed);
    }
    if(!BufferFmt)
    {
        Source->OffsetType = AL_NONE;
        Source->Offset = 0.0;
        return AL_FALSE;
    }

    switch(Source->OffsetType)
    {
    case AL_BYTE_OFFSET:
        /* Determine the ByteOffset (and ensure it is block aligned) */
        *offset = (ALuint)Source->Offset;
        if(BufferFmt->OriginalType == UserFmtIMA4)
        {
            ALsizei align = (BufferFmt->OriginalAlign-1)/2 + 4;
            *offset /= align * ChannelsFromUserFmt(BufferFmt->OriginalChannels);
            *offset *= BufferFmt->OriginalAlign;
        }
        else if(BufferFmt->OriginalType == UserFmtMSADPCM)
        {
            ALsizei align = (BufferFmt->OriginalAlign-2)/2 + 7;
            *offset /= align * ChannelsFromUserFmt(BufferFmt->OriginalChannels);
            *offset *= BufferFmt->OriginalAlign;
        }
        else
            *offset /= FrameSizeFromUserFmt(BufferFmt->OriginalChannels,
                                            BufferFmt->OriginalType);
        *frac = 0;
        break;

    case AL_SAMPLE_OFFSET:
        dblfrac = modf(Source->Offset, &dbloff);
        *offset = (ALuint)mind(dbloff, UINT_MAX);
        *frac = (ALsizei)mind(dblfrac*FRACTIONONE, FRACTIONONE-1.0);
        break;

    case AL_SEC_OFFSET:
        dblfrac = modf(Source->Offset*BufferFmt->Frequency, &dbloff);
        *offset = (ALuint)mind(dbloff, UINT_MAX);
        *frac = (ALsizei)mind(dblfrac*FRACTIONONE, FRACTIONONE-1.0);
        break;
    }
    Source->OffsetType = AL_NONE;
    Source->Offset = 0.0;

    return AL_TRUE;
}


/* ReleaseALSources
 *
 * Destroys all sources in the source map.
 */
ALvoid ReleaseALSources(ALCcontext *Context)
{
    ALCdevice *device = Context->Device;
    ALsizei pos;
    for(pos = 0;pos < Context->SourceMap.size;pos++)
    {
        ALsource *temp = Context->SourceMap.values[pos];
        Context->SourceMap.values[pos] = NULL;

        DeinitSource(temp, device->NumAuxSends);

        FreeThunkEntry(temp->id);
        memset(temp, 0, sizeof(*temp));
        al_free(temp);
    }
}
