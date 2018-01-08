/**
 * OpenAL cross platform audio library
 * Copyright (C) 2013 by Mike Gorchak
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

#include <math.h>
#include <stdlib.h>

#include "alMain.h"
#include "alFilter.h"
#include "alAuxEffectSlot.h"
#include "alError.h"
#include "alu.h"


enum FlangerWaveForm {
    FWF_Triangle = AL_FLANGER_WAVEFORM_TRIANGLE,
    FWF_Sinusoid = AL_FLANGER_WAVEFORM_SINUSOID
};

typedef struct ALflangerState {
    DERIVE_FROM_TYPE(ALeffectState);

    ALfloat *SampleBuffer[2];
    ALsizei BufferLength;
    ALsizei offset;
    ALsizei lfo_range;
    ALfloat lfo_scale;
    ALint lfo_disp;

    /* Gains for left and right sides */
    ALfloat Gain[2][MAX_OUTPUT_CHANNELS];

    /* effect parameters */
    enum FlangerWaveForm waveform;
    ALint delay;
    ALfloat depth;
    ALfloat feedback;
} ALflangerState;

static ALvoid ALflangerState_Destruct(ALflangerState *state);
static ALboolean ALflangerState_deviceUpdate(ALflangerState *state, ALCdevice *Device);
static ALvoid ALflangerState_update(ALflangerState *state, const ALCdevice *Device, const ALeffectslot *Slot, const ALeffectProps *props);
static ALvoid ALflangerState_process(ALflangerState *state, ALsizei SamplesToDo, const ALfloat (*restrict SamplesIn)[BUFFERSIZE], ALfloat (*restrict SamplesOut)[BUFFERSIZE], ALsizei NumChannels);
DECLARE_DEFAULT_ALLOCATORS(ALflangerState)

DEFINE_ALEFFECTSTATE_VTABLE(ALflangerState);


static void ALflangerState_Construct(ALflangerState *state)
{
    ALeffectState_Construct(STATIC_CAST(ALeffectState, state));
    SET_VTABLE2(ALflangerState, ALeffectState, state);

    state->BufferLength = 0;
    state->SampleBuffer[0] = NULL;
    state->SampleBuffer[1] = NULL;
    state->offset = 0;
    state->lfo_range = 1;
    state->waveform = FWF_Triangle;
}

static ALvoid ALflangerState_Destruct(ALflangerState *state)
{
    al_free(state->SampleBuffer[0]);
    state->SampleBuffer[0] = NULL;
    state->SampleBuffer[1] = NULL;

    ALeffectState_Destruct(STATIC_CAST(ALeffectState,state));
}

static ALboolean ALflangerState_deviceUpdate(ALflangerState *state, ALCdevice *Device)
{
    ALsizei maxlen;
    ALsizei it;

    maxlen = fastf2i(AL_FLANGER_MAX_DELAY * 2.0f * Device->Frequency) + 1;
    maxlen = NextPowerOf2(maxlen);

    if(maxlen != state->BufferLength)
    {
        void *temp = al_calloc(16, maxlen * sizeof(ALfloat) * 2);
        if(!temp) return AL_FALSE;

        al_free(state->SampleBuffer[0]);
        state->SampleBuffer[0] = temp;
        state->SampleBuffer[1] = state->SampleBuffer[0] + maxlen;

        state->BufferLength = maxlen;
    }

    for(it = 0;it < state->BufferLength;it++)
    {
        state->SampleBuffer[0][it] = 0.0f;
        state->SampleBuffer[1][it] = 0.0f;
    }

    return AL_TRUE;
}

static ALvoid ALflangerState_update(ALflangerState *state, const ALCdevice *Device, const ALeffectslot *Slot, const ALeffectProps *props)
{
    ALfloat frequency = (ALfloat)Device->Frequency;
    ALfloat coeffs[MAX_AMBI_COEFFS];
    ALfloat rate;
    ALint phase;

    switch(props->Flanger.Waveform)
    {
        case AL_FLANGER_WAVEFORM_TRIANGLE:
            state->waveform = FWF_Triangle;
            break;
        case AL_FLANGER_WAVEFORM_SINUSOID:
            state->waveform = FWF_Sinusoid;
            break;
    }
    state->feedback = props->Flanger.Feedback;
    state->delay = fastf2i(props->Flanger.Delay * frequency);
    /* The LFO depth is scaled to be relative to the sample delay. */
    state->depth = props->Flanger.Depth * state->delay;

    /* Gains for left and right sides */
    CalcAngleCoeffs(-F_PI_2, 0.0f, 0.0f, coeffs);
    ComputePanningGains(Device->Dry, coeffs, Slot->Params.Gain, state->Gain[0]);
    CalcAngleCoeffs( F_PI_2, 0.0f, 0.0f, coeffs);
    ComputePanningGains(Device->Dry, coeffs, Slot->Params.Gain, state->Gain[1]);

    phase = props->Flanger.Phase;
    rate = props->Flanger.Rate;
    if(!(rate > 0.0f))
    {
        state->lfo_scale = 0.0f;
        state->lfo_range = 1;
        state->lfo_disp = 0;
    }
    else
    {
        /* Calculate LFO coefficient */
        state->lfo_range = fastf2i(frequency/rate + 0.5f);
        switch(state->waveform)
        {
            case FWF_Triangle:
                state->lfo_scale = 4.0f / state->lfo_range;
                break;
            case FWF_Sinusoid:
                state->lfo_scale = F_TAU / state->lfo_range;
                break;
        }

        /* Calculate lfo phase displacement */
        if(phase >= 0)
            state->lfo_disp = fastf2i(state->lfo_range * (phase/360.0f));
        else
            state->lfo_disp = fastf2i(state->lfo_range * ((360+phase)/360.0f));
    }
}

static void GetTriangleDelays(ALint *restrict delays, ALsizei offset, const ALsizei lfo_range,
                              const ALfloat lfo_scale, const ALfloat depth, const ALsizei delay,
                              const ALsizei todo)
{
    ALsizei i;
    for(i = 0;i < todo;i++)
    {
        delays[i] = fastf2i((1.0f - fabsf(2.0f - lfo_scale*offset)) * depth) + delay;
        offset = (offset+1)%lfo_range;
    }
}

static void GetSinusoidDelays(ALint *restrict delays, ALsizei offset, const ALsizei lfo_range,
                              const ALfloat lfo_scale, const ALfloat depth, const ALsizei delay,
                              const ALsizei todo)
{
    ALsizei i;
    for(i = 0;i < todo;i++)
    {
        delays[i] = fastf2i(sinf(lfo_scale*offset) * depth) + delay;
        offset = (offset+1)%lfo_range;
    }
}

static ALvoid ALflangerState_process(ALflangerState *state, ALsizei SamplesToDo, const ALfloat (*restrict SamplesIn)[BUFFERSIZE], ALfloat (*restrict SamplesOut)[BUFFERSIZE], ALsizei NumChannels)
{
    ALfloat *restrict leftbuf = state->SampleBuffer[0];
    ALfloat *restrict rightbuf = state->SampleBuffer[1];
    const ALsizei bufmask = state->BufferLength-1;
    const ALfloat feedback = state->feedback;
    ALsizei offset = state->offset;
    ALsizei i, c;
    ALsizei base;

    for(base = 0;base < SamplesToDo;)
    {
        const ALsizei todo = mini(128, SamplesToDo-base);
        ALfloat temps[128][2];
        ALint moddelays[2][128];

        switch(state->waveform)
        {
            case FWF_Triangle:
                GetTriangleDelays(moddelays[0], offset%state->lfo_range, state->lfo_range,
                                  state->lfo_scale, state->depth, state->delay, todo);
                GetTriangleDelays(moddelays[1], (offset+state->lfo_disp)%state->lfo_range,
                                  state->lfo_range, state->lfo_scale, state->depth, state->delay,
                                  todo);
                break;
            case FWF_Sinusoid:
                GetSinusoidDelays(moddelays[0], offset%state->lfo_range, state->lfo_range,
                                  state->lfo_scale, state->depth, state->delay, todo);
                GetSinusoidDelays(moddelays[1], (offset+state->lfo_disp)%state->lfo_range,
                                  state->lfo_range, state->lfo_scale, state->depth, state->delay,
                                  todo);
                break;
        }

        for(i = 0;i < todo;i++)
        {
            leftbuf[offset&bufmask] = SamplesIn[0][base+i];
            temps[i][0] = leftbuf[(offset-moddelays[0][i])&bufmask] * feedback;
            leftbuf[offset&bufmask] += temps[i][0];

            rightbuf[offset&bufmask] = SamplesIn[0][base+i];
            temps[i][1] = rightbuf[(offset-moddelays[1][i])&bufmask] * feedback;
            rightbuf[offset&bufmask] += temps[i][1];

            offset++;
        }

        for(c = 0;c < NumChannels;c++)
        {
            ALfloat gain = state->Gain[0][c];
            if(fabsf(gain) > GAIN_SILENCE_THRESHOLD)
            {
                for(i = 0;i < todo;i++)
                    SamplesOut[c][i+base] += temps[i][0] * gain;
            }

            gain = state->Gain[1][c];
            if(fabsf(gain) > GAIN_SILENCE_THRESHOLD)
            {
                for(i = 0;i < todo;i++)
                    SamplesOut[c][i+base] += temps[i][1] * gain;
            }
        }

        base += todo;
    }

    state->offset = offset;
}


typedef struct ALflangerStateFactory {
    DERIVE_FROM_TYPE(ALeffectStateFactory);
} ALflangerStateFactory;

ALeffectState *ALflangerStateFactory_create(ALflangerStateFactory *UNUSED(factory))
{
    ALflangerState *state;

    NEW_OBJ0(state, ALflangerState)();
    if(!state) return NULL;

    return STATIC_CAST(ALeffectState, state);
}

DEFINE_ALEFFECTSTATEFACTORY_VTABLE(ALflangerStateFactory);

ALeffectStateFactory *ALflangerStateFactory_getFactory(void)
{
    static ALflangerStateFactory FlangerFactory = { { GET_VTABLE2(ALflangerStateFactory, ALeffectStateFactory) } };

    return STATIC_CAST(ALeffectStateFactory, &FlangerFactory);
}


void ALflanger_setParami(ALeffect *effect, ALCcontext *context, ALenum param, ALint val)
{
    ALeffectProps *props = &effect->Props;
    switch(param)
    {
        case AL_FLANGER_WAVEFORM:
            if(!(val >= AL_FLANGER_MIN_WAVEFORM && val <= AL_FLANGER_MAX_WAVEFORM))
                SET_ERROR_AND_RETURN(context, AL_INVALID_VALUE);
            props->Flanger.Waveform = val;
            break;

        case AL_FLANGER_PHASE:
            if(!(val >= AL_FLANGER_MIN_PHASE && val <= AL_FLANGER_MAX_PHASE))
                SET_ERROR_AND_RETURN(context, AL_INVALID_VALUE);
            props->Flanger.Phase = val;
            break;

        default:
            SET_ERROR_AND_RETURN(context, AL_INVALID_ENUM);
    }
}
void ALflanger_setParamiv(ALeffect *effect, ALCcontext *context, ALenum param, const ALint *vals)
{
    ALflanger_setParami(effect, context, param, vals[0]);
}
void ALflanger_setParamf(ALeffect *effect, ALCcontext *context, ALenum param, ALfloat val)
{
    ALeffectProps *props = &effect->Props;
    switch(param)
    {
        case AL_FLANGER_RATE:
            if(!(val >= AL_FLANGER_MIN_RATE && val <= AL_FLANGER_MAX_RATE))
                SET_ERROR_AND_RETURN(context, AL_INVALID_VALUE);
            props->Flanger.Rate = val;
            break;

        case AL_FLANGER_DEPTH:
            if(!(val >= AL_FLANGER_MIN_DEPTH && val <= AL_FLANGER_MAX_DEPTH))
                SET_ERROR_AND_RETURN(context, AL_INVALID_VALUE);
            props->Flanger.Depth = val;
            break;

        case AL_FLANGER_FEEDBACK:
            if(!(val >= AL_FLANGER_MIN_FEEDBACK && val <= AL_FLANGER_MAX_FEEDBACK))
                SET_ERROR_AND_RETURN(context, AL_INVALID_VALUE);
            props->Flanger.Feedback = val;
            break;

        case AL_FLANGER_DELAY:
            if(!(val >= AL_FLANGER_MIN_DELAY && val <= AL_FLANGER_MAX_DELAY))
                SET_ERROR_AND_RETURN(context, AL_INVALID_VALUE);
            props->Flanger.Delay = val;
            break;

        default:
            SET_ERROR_AND_RETURN(context, AL_INVALID_ENUM);
    }
}
void ALflanger_setParamfv(ALeffect *effect, ALCcontext *context, ALenum param, const ALfloat *vals)
{
    ALflanger_setParamf(effect, context, param, vals[0]);
}

void ALflanger_getParami(const ALeffect *effect, ALCcontext *context, ALenum param, ALint *val)
{
    const ALeffectProps *props = &effect->Props;
    switch(param)
    {
        case AL_FLANGER_WAVEFORM:
            *val = props->Flanger.Waveform;
            break;

        case AL_FLANGER_PHASE:
            *val = props->Flanger.Phase;
            break;

        default:
            SET_ERROR_AND_RETURN(context, AL_INVALID_ENUM);
    }
}
void ALflanger_getParamiv(const ALeffect *effect, ALCcontext *context, ALenum param, ALint *vals)
{
    ALflanger_getParami(effect, context, param, vals);
}
void ALflanger_getParamf(const ALeffect *effect, ALCcontext *context, ALenum param, ALfloat *val)
{
    const ALeffectProps *props = &effect->Props;
    switch(param)
    {
        case AL_FLANGER_RATE:
            *val = props->Flanger.Rate;
            break;

        case AL_FLANGER_DEPTH:
            *val = props->Flanger.Depth;
            break;

        case AL_FLANGER_FEEDBACK:
            *val = props->Flanger.Feedback;
            break;

        case AL_FLANGER_DELAY:
            *val = props->Flanger.Delay;
            break;

        default:
            SET_ERROR_AND_RETURN(context, AL_INVALID_ENUM);
    }
}
void ALflanger_getParamfv(const ALeffect *effect, ALCcontext *context, ALenum param, ALfloat *vals)
{
    ALflanger_getParamf(effect, context, param, vals);
}

DEFINE_ALEFFECT_VTABLE(ALflanger);
