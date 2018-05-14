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

#include <math.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>
#include <assert.h>

#include "alMain.h"
#include "alSource.h"
#include "alBuffer.h"
#include "alListener.h"
#include "alAuxEffectSlot.h"
#include "alu.h"
#include "bs2b.h"
#include "hrtf.h"
#include "uhjfilter.h"
#include "bformatdec.h"
#include "static_assert.h"

#include "mixer_defs.h"

#include "backends/base.h"


struct ChanMap {
    enum Channel channel;
    ALfloat angle;
    ALfloat elevation;
};

/* Cone scalar */
ALfloat ConeScale = 1.0f;

/* Localized Z scalar for mono sources */
ALfloat ZScale = 1.0f;

extern inline ALfloat minf(ALfloat a, ALfloat b);
extern inline ALfloat maxf(ALfloat a, ALfloat b);
extern inline ALfloat clampf(ALfloat val, ALfloat min, ALfloat max);

extern inline ALdouble mind(ALdouble a, ALdouble b);
extern inline ALdouble maxd(ALdouble a, ALdouble b);
extern inline ALdouble clampd(ALdouble val, ALdouble min, ALdouble max);

extern inline ALuint minu(ALuint a, ALuint b);
extern inline ALuint maxu(ALuint a, ALuint b);
extern inline ALuint clampu(ALuint val, ALuint min, ALuint max);

extern inline ALint mini(ALint a, ALint b);
extern inline ALint maxi(ALint a, ALint b);
extern inline ALint clampi(ALint val, ALint min, ALint max);

extern inline ALint64 mini64(ALint64 a, ALint64 b);
extern inline ALint64 maxi64(ALint64 a, ALint64 b);
extern inline ALint64 clampi64(ALint64 val, ALint64 min, ALint64 max);

extern inline ALuint64 minu64(ALuint64 a, ALuint64 b);
extern inline ALuint64 maxu64(ALuint64 a, ALuint64 b);
extern inline ALuint64 clampu64(ALuint64 val, ALuint64 min, ALuint64 max);

extern inline ALfloat lerp(ALfloat val1, ALfloat val2, ALfloat mu);
extern inline ALfloat resample_fir4(ALfloat val0, ALfloat val1, ALfloat val2, ALfloat val3, ALsizei frac);

extern inline void aluVectorSet(aluVector *restrict vector, ALfloat x, ALfloat y, ALfloat z, ALfloat w);

extern inline void aluMatrixfSetRow(aluMatrixf *matrix, ALuint row,
                                    ALfloat m0, ALfloat m1, ALfloat m2, ALfloat m3);
extern inline void aluMatrixfSet(aluMatrixf *matrix,
                                 ALfloat m00, ALfloat m01, ALfloat m02, ALfloat m03,
                                 ALfloat m10, ALfloat m11, ALfloat m12, ALfloat m13,
                                 ALfloat m20, ALfloat m21, ALfloat m22, ALfloat m23,
                                 ALfloat m30, ALfloat m31, ALfloat m32, ALfloat m33);

const aluMatrixf IdentityMatrixf = {{
    { 1.0f, 0.0f, 0.0f, 0.0f },
    { 0.0f, 1.0f, 0.0f, 0.0f },
    { 0.0f, 0.0f, 1.0f, 0.0f },
    { 0.0f, 0.0f, 0.0f, 1.0f },
}};


void DeinitVoice(ALvoice *voice)
{
    struct ALvoiceProps *props;
    size_t count = 0;

    props = ATOMIC_EXCHANGE_PTR_SEQ(&voice->Update, NULL);
    if(props) al_free(props);

    props = ATOMIC_EXCHANGE_PTR(&voice->FreeList, NULL, almemory_order_relaxed);
    while(props)
    {
        struct ALvoiceProps *next;
        next = ATOMIC_LOAD(&props->next, almemory_order_relaxed);
        al_free(props);
        props = next;
        ++count;
    }
    /* This is excessively spammy if it traces every voice destruction, so just
     * warn if it was unexpectedly large.
     */
    if(count > 3)
        WARN("Freed "SZFMT" voice property objects\n", count);
}


static inline HrtfDirectMixerFunc SelectHrtfMixer(void)
{
#ifdef HAVE_NEON
    if((CPUCapFlags&CPU_CAP_NEON))
        return MixDirectHrtf_Neon;
#endif
#ifdef HAVE_SSE
    if((CPUCapFlags&CPU_CAP_SSE))
        return MixDirectHrtf_SSE;
#endif

    return MixDirectHrtf_C;
}


/* Prior to VS2013, MSVC lacks the round() family of functions. */
#if defined(_MSC_VER) && _MSC_VER < 1800
static float roundf(float val)
{
    if(val < 0.0f)
        return ceilf(val-0.5f);
    return floorf(val+0.5f);
}
#endif

/* This RNG method was created based on the math found in opusdec. It's quick,
 * and starting with a seed value of 22222, is suitable for generating
 * whitenoise.
 */
static inline ALuint dither_rng(ALuint *seed)
{
    *seed = (*seed * 96314165) + 907633515;
    return *seed;
}


static inline void aluCrossproduct(const ALfloat *inVector1, const ALfloat *inVector2, ALfloat *outVector)
{
    outVector[0] = inVector1[1]*inVector2[2] - inVector1[2]*inVector2[1];
    outVector[1] = inVector1[2]*inVector2[0] - inVector1[0]*inVector2[2];
    outVector[2] = inVector1[0]*inVector2[1] - inVector1[1]*inVector2[0];
}

static inline ALfloat aluDotproduct(const aluVector *vec1, const aluVector *vec2)
{
    return vec1->v[0]*vec2->v[0] + vec1->v[1]*vec2->v[1] + vec1->v[2]*vec2->v[2];
}

static ALfloat aluNormalize(ALfloat *vec)
{
    ALfloat length = sqrtf(vec[0]*vec[0] + vec[1]*vec[1] + vec[2]*vec[2]);
    if(length > 0.0f)
    {
        ALfloat inv_length = 1.0f/length;
        vec[0] *= inv_length;
        vec[1] *= inv_length;
        vec[2] *= inv_length;
    }
    return length;
}

static void aluMatrixfFloat3(ALfloat *vec, ALfloat w, const aluMatrixf *mtx)
{
    ALfloat v[4] = { vec[0], vec[1], vec[2], w };

    vec[0] = v[0]*mtx->m[0][0] + v[1]*mtx->m[1][0] + v[2]*mtx->m[2][0] + v[3]*mtx->m[3][0];
    vec[1] = v[0]*mtx->m[0][1] + v[1]*mtx->m[1][1] + v[2]*mtx->m[2][1] + v[3]*mtx->m[3][1];
    vec[2] = v[0]*mtx->m[0][2] + v[1]*mtx->m[1][2] + v[2]*mtx->m[2][2] + v[3]*mtx->m[3][2];
}

static aluVector aluMatrixfVector(const aluMatrixf *mtx, const aluVector *vec)
{
    aluVector v;
    v.v[0] = vec->v[0]*mtx->m[0][0] + vec->v[1]*mtx->m[1][0] + vec->v[2]*mtx->m[2][0] + vec->v[3]*mtx->m[3][0];
    v.v[1] = vec->v[0]*mtx->m[0][1] + vec->v[1]*mtx->m[1][1] + vec->v[2]*mtx->m[2][1] + vec->v[3]*mtx->m[3][1];
    v.v[2] = vec->v[0]*mtx->m[0][2] + vec->v[1]*mtx->m[1][2] + vec->v[2]*mtx->m[2][2] + vec->v[3]*mtx->m[3][2];
    v.v[3] = vec->v[0]*mtx->m[0][3] + vec->v[1]*mtx->m[1][3] + vec->v[2]*mtx->m[2][3] + vec->v[3]*mtx->m[3][3];
    return v;
}


/* Prepares the interpolator for a given rate (determined by increment).  A
 * result of AL_FALSE indicates that the filter output will completely cut
 * the input signal.
 *
 * With a bit of work, and a trade of memory for CPU cost, this could be
 * modified for use with an interpolated increment for buttery-smooth pitch
 * changes.
 */
ALboolean BsincPrepare(const ALuint increment, BsincState *state)
{
    static const ALfloat scaleBase = 1.510578918e-01f, scaleRange = 1.177936623e+00f;
    static const ALuint m[BSINC_SCALE_COUNT] = { 24, 24, 24, 24, 24, 24, 24, 20, 20, 20, 16, 16, 16, 12, 12, 12 };
    static const ALuint to[4][BSINC_SCALE_COUNT] =
    {
        { 0, 24, 408, 792, 1176, 1560, 1944, 2328, 2648, 2968, 3288, 3544, 3800, 4056, 4248, 4440 },
        { 4632, 5016, 5400, 5784, 6168, 6552, 6936, 7320, 7640, 7960, 8280, 8536, 8792, 9048, 9240, 0 },
        { 0, 9432, 9816, 10200, 10584, 10968, 11352, 11736, 12056, 12376, 12696, 12952, 13208, 13464, 13656, 13848 },
        { 14040, 14424, 14808, 15192, 15576, 15960, 16344, 16728, 17048, 17368, 17688, 17944, 18200, 18456, 18648, 0 }
    };
    static const ALuint tm[2][BSINC_SCALE_COUNT] =
    {
        { 0, 24, 24, 24, 24, 24, 24, 20, 20, 20, 16, 16, 16, 12, 12, 12 },
        { 24, 24, 24, 24, 24, 24, 24, 20, 20, 20, 16, 16, 16, 12, 12, 0 }
    };
    ALfloat sf;
    ALsizei si, pi;
    ALboolean uncut = AL_TRUE;

    if(increment > FRACTIONONE)
    {
        sf = (ALfloat)FRACTIONONE / increment;
        if(sf < scaleBase)
        {
            /* Signal has been completely cut.  The return result can be used
             * to skip the filter (and output zeros) as an optimization.
             */
            sf = 0.0f;
            si = 0;
            uncut = AL_FALSE;
        }
        else
        {
            sf = (BSINC_SCALE_COUNT - 1) * (sf - scaleBase) * scaleRange;
            si = fastf2i(sf);
            /* The interpolation factor is fit to this diagonally-symmetric
             * curve to reduce the transition ripple caused by interpolating
             * different scales of the sinc function.
             */
            sf = 1.0f - cosf(asinf(sf - si));
        }
    }
    else
    {
        sf = 0.0f;
        si = BSINC_SCALE_COUNT - 1;
    }

    state->sf = sf;
    state->m = m[si];
    state->l = -(ALint)((m[si] / 2) - 1);
    /* The CPU cost of this table re-mapping could be traded for the memory
     * cost of a complete table map (1024 elements large).
     */
    for(pi = 0;pi < BSINC_PHASE_COUNT;pi++)
    {
        state->coeffs[pi].filter  = &bsincTab[to[0][si] + tm[0][si]*pi];
        state->coeffs[pi].scDelta = &bsincTab[to[1][si] + tm[1][si]*pi];
        state->coeffs[pi].phDelta = &bsincTab[to[2][si] + tm[0][si]*pi];
        state->coeffs[pi].spDelta = &bsincTab[to[3][si] + tm[1][si]*pi];
    }
    return uncut;
}


static ALboolean CalcListenerParams(ALCcontext *Context)
{
    ALlistener *Listener = Context->Listener;
    ALfloat N[3], V[3], U[3], P[3];
    struct ALlistenerProps *props;
    aluVector vel;

    props = ATOMIC_EXCHANGE_PTR(&Listener->Update, NULL, almemory_order_acq_rel);
    if(!props) return AL_FALSE;

    /* AT then UP */
    N[0] = props->Forward[0];
    N[1] = props->Forward[1];
    N[2] = props->Forward[2];
    aluNormalize(N);
    V[0] = props->Up[0];
    V[1] = props->Up[1];
    V[2] = props->Up[2];
    aluNormalize(V);
    /* Build and normalize right-vector */
    aluCrossproduct(N, V, U);
    aluNormalize(U);

    aluMatrixfSet(&Listener->Params.Matrix,
        U[0], V[0], -N[0], 0.0,
        U[1], V[1], -N[1], 0.0,
        U[2], V[2], -N[2], 0.0,
         0.0,  0.0,   0.0, 1.0
    );

    P[0] = props->Position[0];
    P[1] = props->Position[1];
    P[2] = props->Position[2];
    aluMatrixfFloat3(P, 1.0, &Listener->Params.Matrix);
    aluMatrixfSetRow(&Listener->Params.Matrix, 3, -P[0], -P[1], -P[2], 1.0f);

    aluVectorSet(&vel, props->Velocity[0], props->Velocity[1], props->Velocity[2], 0.0f);
    Listener->Params.Velocity = aluMatrixfVector(&Listener->Params.Matrix, &vel);

    Listener->Params.Gain = props->Gain * Context->GainBoost;
    Listener->Params.MetersPerUnit = props->MetersPerUnit;

    Listener->Params.DopplerFactor = props->DopplerFactor;
    Listener->Params.SpeedOfSound = props->SpeedOfSound * props->DopplerVelocity;

    Listener->Params.SourceDistanceModel = props->SourceDistanceModel;
    Listener->Params.DistanceModel = props->DistanceModel;

    ATOMIC_REPLACE_HEAD(struct ALlistenerProps*, &Listener->FreeList, props);
    return AL_TRUE;
}

static ALboolean CalcEffectSlotParams(ALeffectslot *slot, ALCdevice *device)
{
    struct ALeffectslotProps *props;
    ALeffectState *state;

    props = ATOMIC_EXCHANGE_PTR(&slot->Update, NULL, almemory_order_acq_rel);
    if(!props) return AL_FALSE;

    slot->Params.Gain = props->Gain;
    slot->Params.AuxSendAuto = props->AuxSendAuto;
    slot->Params.EffectType = props->Type;
    if(IsReverbEffect(slot->Params.EffectType))
    {
        slot->Params.RoomRolloff = props->Props.Reverb.RoomRolloffFactor;
        slot->Params.DecayTime = props->Props.Reverb.DecayTime;
        slot->Params.DecayHFRatio = props->Props.Reverb.DecayHFRatio;
        slot->Params.DecayHFLimit = props->Props.Reverb.DecayHFLimit;
        slot->Params.AirAbsorptionGainHF = props->Props.Reverb.AirAbsorptionGainHF;
    }
    else
    {
        slot->Params.RoomRolloff = 0.0f;
        slot->Params.DecayTime = 0.0f;
        slot->Params.DecayHFRatio = 0.0f;
        slot->Params.DecayHFLimit = AL_FALSE;
        slot->Params.AirAbsorptionGainHF = 1.0f;
    }

    /* Swap effect states. No need to play with the ref counts since they keep
     * the same number of refs.
     */
    state = props->State;
    props->State = slot->Params.EffectState;
    slot->Params.EffectState = state;

    V(state,update)(device, slot, &props->Props);

    ATOMIC_REPLACE_HEAD(struct ALeffectslotProps*, &slot->FreeList, props);
    return AL_TRUE;
}


static const struct ChanMap MonoMap[1] = {
    { FrontCenter, 0.0f, 0.0f }
}, RearMap[2] = {
    { BackLeft,  DEG2RAD(-150.0f), DEG2RAD(0.0f) },
    { BackRight, DEG2RAD( 150.0f), DEG2RAD(0.0f) }
}, QuadMap[4] = {
    { FrontLeft,  DEG2RAD( -45.0f), DEG2RAD(0.0f) },
    { FrontRight, DEG2RAD(  45.0f), DEG2RAD(0.0f) },
    { BackLeft,   DEG2RAD(-135.0f), DEG2RAD(0.0f) },
    { BackRight,  DEG2RAD( 135.0f), DEG2RAD(0.0f) }
}, X51Map[6] = {
    { FrontLeft,   DEG2RAD( -30.0f), DEG2RAD(0.0f) },
    { FrontRight,  DEG2RAD(  30.0f), DEG2RAD(0.0f) },
    { FrontCenter, DEG2RAD(   0.0f), DEG2RAD(0.0f) },
    { LFE, 0.0f, 0.0f },
    { SideLeft,    DEG2RAD(-110.0f), DEG2RAD(0.0f) },
    { SideRight,   DEG2RAD( 110.0f), DEG2RAD(0.0f) }
}, X61Map[7] = {
    { FrontLeft,    DEG2RAD(-30.0f), DEG2RAD(0.0f) },
    { FrontRight,   DEG2RAD( 30.0f), DEG2RAD(0.0f) },
    { FrontCenter,  DEG2RAD(  0.0f), DEG2RAD(0.0f) },
    { LFE, 0.0f, 0.0f },
    { BackCenter,   DEG2RAD(180.0f), DEG2RAD(0.0f) },
    { SideLeft,     DEG2RAD(-90.0f), DEG2RAD(0.0f) },
    { SideRight,    DEG2RAD( 90.0f), DEG2RAD(0.0f) }
}, X71Map[8] = {
    { FrontLeft,   DEG2RAD( -30.0f), DEG2RAD(0.0f) },
    { FrontRight,  DEG2RAD(  30.0f), DEG2RAD(0.0f) },
    { FrontCenter, DEG2RAD(   0.0f), DEG2RAD(0.0f) },
    { LFE, 0.0f, 0.0f },
    { BackLeft,    DEG2RAD(-150.0f), DEG2RAD(0.0f) },
    { BackRight,   DEG2RAD( 150.0f), DEG2RAD(0.0f) },
    { SideLeft,    DEG2RAD( -90.0f), DEG2RAD(0.0f) },
    { SideRight,   DEG2RAD(  90.0f), DEG2RAD(0.0f) }
};

static void CalcPanningAndFilters(ALvoice *voice, const ALfloat Distance, const ALfloat *Dir,
                                  const ALfloat Spread, const ALfloat DryGain,
                                  const ALfloat DryGainHF, const ALfloat DryGainLF,
                                  const ALfloat *WetGain, const ALfloat *WetGainLF,
                                  const ALfloat *WetGainHF, ALeffectslot **SendSlots,
                                  const ALbuffer *Buffer, const struct ALvoiceProps *props,
                                  const ALlistener *Listener, const ALCdevice *Device)
{
    struct ChanMap StereoMap[2] = {
        { FrontLeft,  DEG2RAD(-30.0f), DEG2RAD(0.0f) },
        { FrontRight, DEG2RAD( 30.0f), DEG2RAD(0.0f) }
    };
    bool DirectChannels = props->DirectChannels;
    const ALsizei NumSends = Device->NumAuxSends;
    const ALuint Frequency = Device->Frequency;
    const struct ChanMap *chans = NULL;
    ALsizei num_channels = 0;
    bool isbformat = false;
    ALfloat downmix_gain = 1.0f;
    ALsizei c, i, j;

    switch(Buffer->FmtChannels)
    {
    case FmtMono:
        chans = MonoMap;
        num_channels = 1;
        /* Mono buffers are never played direct. */
        DirectChannels = false;
        break;

    case FmtStereo:
        /* Convert counter-clockwise to clockwise. */
        StereoMap[0].angle = -props->StereoPan[0];
        StereoMap[1].angle = -props->StereoPan[1];

        chans = StereoMap;
        num_channels = 2;
        downmix_gain = 1.0f / 2.0f;
        break;

    case FmtRear:
        chans = RearMap;
        num_channels = 2;
        downmix_gain = 1.0f / 2.0f;
        break;

    case FmtQuad:
        chans = QuadMap;
        num_channels = 4;
        downmix_gain = 1.0f / 4.0f;
        break;

    case FmtX51:
        chans = X51Map;
        num_channels = 6;
        /* NOTE: Excludes LFE. */
        downmix_gain = 1.0f / 5.0f;
        break;

    case FmtX61:
        chans = X61Map;
        num_channels = 7;
        /* NOTE: Excludes LFE. */
        downmix_gain = 1.0f / 6.0f;
        break;

    case FmtX71:
        chans = X71Map;
        num_channels = 8;
        /* NOTE: Excludes LFE. */
        downmix_gain = 1.0f / 7.0f;
        break;

    case FmtBFormat2D:
        num_channels = 3;
        isbformat = true;
        DirectChannels = false;
        break;

    case FmtBFormat3D:
        num_channels = 4;
        isbformat = true;
        DirectChannels = false;
        break;
    }

    voice->Flags &= ~(VOICE_HAS_HRTF | VOICE_HAS_NFC);
    if(isbformat)
    {
        /* Special handling for B-Format sources. */

        if(Distance > FLT_EPSILON)
        {
            /* Panning a B-Format sound toward some direction is easy. Just pan
             * the first (W) channel as a normal mono sound and silence the
             * others.
             */
            ALfloat coeffs[MAX_AMBI_COEFFS];

            if(Device->AvgSpeakerDist > 0.0f && Listener->Params.MetersPerUnit > 0.0f)
            {
                ALfloat mdist = Distance * Listener->Params.MetersPerUnit;
                ALfloat w0 = SPEEDOFSOUNDMETRESPERSEC /
                             (mdist * (ALfloat)Device->Frequency);
                ALfloat w1 = SPEEDOFSOUNDMETRESPERSEC /
                             (Device->AvgSpeakerDist * (ALfloat)Device->Frequency);
                /* Clamp w0 for really close distances, to prevent excessive
                 * bass.
                 */
                w0 = minf(w0, w1*4.0f);

                /* Only need to adjust the first channel of a B-Format source. */
                NfcFilterAdjust1(&voice->Direct.Params[0].NFCtrlFilter[0], w0);
                NfcFilterAdjust2(&voice->Direct.Params[0].NFCtrlFilter[1], w0);
                NfcFilterAdjust3(&voice->Direct.Params[0].NFCtrlFilter[2], w0);

                for(i = 0;i < MAX_AMBI_ORDER+1;i++)
                    voice->Direct.ChannelsPerOrder[i] = Device->Dry.NumChannelsPerOrder[i];
                voice->Flags |= VOICE_HAS_NFC;
            }

            if(Device->Render_Mode == StereoPair)
            {
                ALfloat ev = asinf(Dir[1]);
                ALfloat az = atan2f(Dir[0], -Dir[2]);
                CalcAnglePairwiseCoeffs(az, ev, Spread, coeffs);
            }
            else
                CalcDirectionCoeffs(Dir, Spread, coeffs);

            /* NOTE: W needs to be scaled by sqrt(2) due to FuMa normalization. */
            ComputePanningGains(Device->Dry, coeffs, DryGain*1.414213562f,
                                voice->Direct.Params[0].Gains.Target);
            for(c = 1;c < num_channels;c++)
            {
                for(j = 0;j < MAX_OUTPUT_CHANNELS;j++)
                    voice->Direct.Params[c].Gains.Target[j] = 0.0f;
            }

            for(i = 0;i < NumSends;i++)
            {
                const ALeffectslot *Slot = SendSlots[i];
                if(Slot)
                    ComputePanningGainsBF(Slot->ChanMap, Slot->NumChannels,
                        coeffs, WetGain[i]*1.414213562f, voice->Send[i].Params[0].Gains.Target
                    );
                else
                    for(j = 0;j < MAX_EFFECT_CHANNELS;j++)
                        voice->Send[i].Params[0].Gains.Target[j] = 0.0f;
                for(c = 1;c < num_channels;c++)
                {
                    for(j = 0;j < MAX_EFFECT_CHANNELS;j++)
                        voice->Send[i].Params[c].Gains.Target[j] = 0.0f;
                }
            }
        }
        else
        {
            /* Local B-Format sources have their XYZ channels rotated according
             * to the orientation.
             */
            ALfloat N[3], V[3], U[3];
            aluMatrixf matrix;
            ALfloat scale;

            if(Device->AvgSpeakerDist > 0.0f)
            {
                /* NOTE: The NFCtrlFilters were created with a w0 of 0, which
                 * is what we want for FOA input. The first channel may have
                 * been previously re-adjusted if panned, so reset it.
                 */
                NfcFilterAdjust1(&voice->Direct.Params[0].NFCtrlFilter[0], 0.0f);
                NfcFilterAdjust2(&voice->Direct.Params[0].NFCtrlFilter[1], 0.0f);
                NfcFilterAdjust3(&voice->Direct.Params[0].NFCtrlFilter[2], 0.0f);

                voice->Direct.ChannelsPerOrder[0] = 1;
                voice->Direct.ChannelsPerOrder[1] = mini(voice->Direct.Channels-1, 3);
                for(i = 2;i < MAX_AMBI_ORDER+1;i++)
                    voice->Direct.ChannelsPerOrder[2] = 0;
                voice->Flags |= VOICE_HAS_NFC;
            }

            /* AT then UP */
            N[0] = props->Orientation[0][0];
            N[1] = props->Orientation[0][1];
            N[2] = props->Orientation[0][2];
            aluNormalize(N);
            V[0] = props->Orientation[1][0];
            V[1] = props->Orientation[1][1];
            V[2] = props->Orientation[1][2];
            aluNormalize(V);
            if(!props->HeadRelative)
            {
                const aluMatrixf *lmatrix = &Listener->Params.Matrix;
                aluMatrixfFloat3(N, 0.0f, lmatrix);
                aluMatrixfFloat3(V, 0.0f, lmatrix);
            }
            /* Build and normalize right-vector */
            aluCrossproduct(N, V, U);
            aluNormalize(U);

            /* Build a rotate + conversion matrix (FuMa -> ACN+N3D). */
            scale = 1.732050808f;
            aluMatrixfSet(&matrix,
                1.414213562f,        0.0f,        0.0f,        0.0f,
                        0.0f, -N[0]*scale,  N[1]*scale, -N[2]*scale,
                        0.0f,  U[0]*scale, -U[1]*scale,  U[2]*scale,
                        0.0f, -V[0]*scale,  V[1]*scale, -V[2]*scale
            );

            voice->Direct.Buffer = Device->FOAOut.Buffer;
            voice->Direct.Channels = Device->FOAOut.NumChannels;
            for(c = 0;c < num_channels;c++)
                ComputeFirstOrderGains(Device->FOAOut, matrix.m[c], DryGain,
                                       voice->Direct.Params[c].Gains.Target);
            for(i = 0;i < NumSends;i++)
            {
                const ALeffectslot *Slot = SendSlots[i];
                if(Slot)
                {
                    for(c = 0;c < num_channels;c++)
                        ComputeFirstOrderGainsBF(Slot->ChanMap, Slot->NumChannels,
                            matrix.m[c], WetGain[i], voice->Send[i].Params[c].Gains.Target
                        );
                }
                else
                {
                    for(c = 0;c < num_channels;c++)
                        for(j = 0;j < MAX_EFFECT_CHANNELS;j++)
                            voice->Send[i].Params[c].Gains.Target[j] = 0.0f;
                }
            }
        }
    }
    else if(DirectChannels)
    {
        /* Direct source channels always play local. Skip the virtual channels
         * and write inputs to the matching real outputs.
         */
        voice->Direct.Buffer = Device->RealOut.Buffer;
        voice->Direct.Channels = Device->RealOut.NumChannels;

        for(c = 0;c < num_channels;c++)
        {
            int idx;
            for(j = 0;j < MAX_OUTPUT_CHANNELS;j++)
                voice->Direct.Params[c].Gains.Target[j] = 0.0f;
            if((idx=GetChannelIdxByName(Device->RealOut, chans[c].channel)) != -1)
                voice->Direct.Params[c].Gains.Target[idx] = DryGain;
        }

        /* Auxiliary sends still use normal channel panning since they mix to
         * B-Format, which can't channel-match.
         */
        for(c = 0;c < num_channels;c++)
        {
            ALfloat coeffs[MAX_AMBI_COEFFS];
            CalcAngleCoeffs(chans[c].angle, chans[c].elevation, 0.0f, coeffs);

            for(i = 0;i < NumSends;i++)
            {
                const ALeffectslot *Slot = SendSlots[i];
                if(Slot)
                    ComputePanningGainsBF(Slot->ChanMap, Slot->NumChannels,
                        coeffs, WetGain[i], voice->Send[i].Params[c].Gains.Target
                    );
                else
                    for(j = 0;j < MAX_EFFECT_CHANNELS;j++)
                        voice->Send[i].Params[c].Gains.Target[j] = 0.0f;
            }
        }
    }
    else if(Device->Render_Mode == HrtfRender)
    {
        /* Full HRTF rendering. Skip the virtual channels and render to the
         * real outputs.
         */
        voice->Direct.Buffer = Device->RealOut.Buffer;
        voice->Direct.Channels = Device->RealOut.NumChannels;

        if(Distance > FLT_EPSILON)
        {
            ALfloat coeffs[MAX_AMBI_COEFFS];
            ALfloat ev, az;

            ev = asinf(Dir[1]);
            az = atan2f(Dir[0], -Dir[2]);

            /* Get the HRIR coefficients and delays just once, for the given
             * source direction.
             */
            GetHrtfCoeffs(Device->HrtfHandle, ev, az, Spread,
                          voice->Direct.Params[0].Hrtf.Target.Coeffs,
                          voice->Direct.Params[0].Hrtf.Target.Delay);
            voice->Direct.Params[0].Hrtf.Target.Gain = DryGain * downmix_gain;

            /* Remaining channels use the same results as the first. */
            for(c = 1;c < num_channels;c++)
            {
                /* Skip LFE */
                if(chans[c].channel == LFE)
                    memset(&voice->Direct.Params[c].Hrtf.Target, 0,
                           sizeof(voice->Direct.Params[c].Hrtf.Target));
                else
                    voice->Direct.Params[c].Hrtf.Target = voice->Direct.Params[0].Hrtf.Target;
            }

            /* Calculate the directional coefficients once, which apply to all
             * input channels of the source sends.
             */
            CalcDirectionCoeffs(Dir, Spread, coeffs);

            for(i = 0;i < NumSends;i++)
            {
                const ALeffectslot *Slot = SendSlots[i];
                if(Slot)
                    for(c = 0;c < num_channels;c++)
                    {
                        /* Skip LFE */
                        if(chans[c].channel == LFE)
                            for(j = 0;j < MAX_EFFECT_CHANNELS;j++)
                                voice->Send[i].Params[c].Gains.Target[j] = 0.0f;
                        else
                            ComputePanningGainsBF(Slot->ChanMap,
                                Slot->NumChannels, coeffs, WetGain[i] * downmix_gain,
                                voice->Send[i].Params[c].Gains.Target
                            );
                    }
                else
                    for(j = 0;j < MAX_EFFECT_CHANNELS;j++)
                        voice->Send[i].Params[c].Gains.Target[j] = 0.0f;
            }
        }
        else
        {
            /* Local sources on HRTF play with each channel panned to its
             * relative location around the listener, providing "virtual
             * speaker" responses.
             */
            for(c = 0;c < num_channels;c++)
            {
                ALfloat coeffs[MAX_AMBI_COEFFS];

                if(chans[c].channel == LFE)
                {
                    /* Skip LFE */
                    memset(&voice->Direct.Params[c].Hrtf.Target, 0,
                           sizeof(voice->Direct.Params[c].Hrtf.Target));
                    for(i = 0;i < NumSends;i++)
                    {
                        for(j = 0;j < MAX_EFFECT_CHANNELS;j++)
                            voice->Send[i].Params[c].Gains.Target[j] = 0.0f;
                    }
                    continue;
                }

                /* Get the HRIR coefficients and delays for this channel
                 * position.
                 */
                GetHrtfCoeffs(Device->HrtfHandle,
                    chans[c].elevation, chans[c].angle, Spread,
                    voice->Direct.Params[c].Hrtf.Target.Coeffs,
                    voice->Direct.Params[c].Hrtf.Target.Delay
                );
                voice->Direct.Params[c].Hrtf.Target.Gain = DryGain;

                /* Normal panning for auxiliary sends. */
                CalcAngleCoeffs(chans[c].angle, chans[c].elevation, Spread, coeffs);

                for(i = 0;i < NumSends;i++)
                {
                    const ALeffectslot *Slot = SendSlots[i];
                    if(Slot)
                        ComputePanningGainsBF(Slot->ChanMap, Slot->NumChannels,
                            coeffs, WetGain[i], voice->Send[i].Params[c].Gains.Target
                        );
                    else
                        for(j = 0;j < MAX_EFFECT_CHANNELS;j++)
                            voice->Send[i].Params[c].Gains.Target[j] = 0.0f;
                }
            }
        }

        voice->Flags |= VOICE_HAS_HRTF;
    }
    else
    {
        /* Non-HRTF rendering. Use normal panning to the output. */

        if(Distance > FLT_EPSILON)
        {
            ALfloat coeffs[MAX_AMBI_COEFFS];
            ALfloat w0 = 0.0f;

            /* Calculate NFC filter coefficient if needed. */
            if(Device->AvgSpeakerDist > 0.0f && Listener->Params.MetersPerUnit > 0.0f)
            {
                ALfloat mdist = Distance * Listener->Params.MetersPerUnit;
                ALfloat w1 = SPEEDOFSOUNDMETRESPERSEC /
                             (Device->AvgSpeakerDist * (ALfloat)Device->Frequency);
                w0 = SPEEDOFSOUNDMETRESPERSEC /
                     (mdist * (ALfloat)Device->Frequency);
                /* Clamp w0 for really close distances, to prevent excessive
                 * bass.
                 */
                w0 = minf(w0, w1*4.0f);

                for(i = 0;i < MAX_AMBI_ORDER+1;i++)
                    voice->Direct.ChannelsPerOrder[i] = Device->Dry.NumChannelsPerOrder[i];
                voice->Flags |= VOICE_HAS_NFC;
            }

            /* Calculate the directional coefficients once, which apply to all
             * input channels.
             */
            if(Device->Render_Mode == StereoPair)
            {
                ALfloat ev = asinf(Dir[1]);
                ALfloat az = atan2f(Dir[0], -Dir[2]);
                CalcAnglePairwiseCoeffs(az, ev, Spread, coeffs);
            }
            else
                CalcDirectionCoeffs(Dir, Spread, coeffs);

            for(c = 0;c < num_channels;c++)
            {
                /* Adjust NFC filters if needed. */
                if((voice->Flags&VOICE_HAS_NFC))
                {
                    NfcFilterAdjust1(&voice->Direct.Params[c].NFCtrlFilter[0], w0);
                    NfcFilterAdjust2(&voice->Direct.Params[c].NFCtrlFilter[1], w0);
                    NfcFilterAdjust3(&voice->Direct.Params[c].NFCtrlFilter[2], w0);
                }

                /* Special-case LFE */
                if(chans[c].channel == LFE)
                {
                    for(j = 0;j < MAX_OUTPUT_CHANNELS;j++)
                        voice->Direct.Params[c].Gains.Target[j] = 0.0f;
                    if(Device->Dry.Buffer == Device->RealOut.Buffer)
                    {
                        int idx = GetChannelIdxByName(Device->RealOut, chans[c].channel);
                        if(idx != -1) voice->Direct.Params[c].Gains.Target[idx] = DryGain;
                    }
                    continue;
                }

                ComputePanningGains(Device->Dry,
                    coeffs, DryGain * downmix_gain, voice->Direct.Params[c].Gains.Target
                );
            }

            for(i = 0;i < NumSends;i++)
            {
                const ALeffectslot *Slot = SendSlots[i];
                if(Slot)
                    for(c = 0;c < num_channels;c++)
                    {
                        /* Skip LFE */
                        if(chans[c].channel == LFE)
                            for(j = 0;j < MAX_EFFECT_CHANNELS;j++)
                                voice->Send[i].Params[c].Gains.Target[j] = 0.0f;
                        else
                            ComputePanningGainsBF(Slot->ChanMap,
                                Slot->NumChannels, coeffs, WetGain[i] * downmix_gain,
                                voice->Send[i].Params[c].Gains.Target
                            );
                    }
                else
                    for(c = 0;c < num_channels;c++)
                    {
                        for(j = 0;j < MAX_EFFECT_CHANNELS;j++)
                            voice->Send[i].Params[c].Gains.Target[j] = 0.0f;
                    }
            }
        }
        else
        {
            ALfloat w0 = 0.0f;

            if(Device->AvgSpeakerDist > 0.0f)
            {
                /* If the source distance is 0, set w0 to w1 to act as a pass-
                 * through. We still want to pass the signal through the
                 * filters so they keep an appropriate history, in case the
                 * source moves away from the listener.
                 */
                w0 = SPEEDOFSOUNDMETRESPERSEC /
                     (Device->AvgSpeakerDist * (ALfloat)Device->Frequency);

                for(i = 0;i < MAX_AMBI_ORDER+1;i++)
                    voice->Direct.ChannelsPerOrder[i] = Device->Dry.NumChannelsPerOrder[i];
                voice->Flags |= VOICE_HAS_NFC;
            }

            for(c = 0;c < num_channels;c++)
            {
                ALfloat coeffs[MAX_AMBI_COEFFS];

                if((voice->Flags&VOICE_HAS_NFC))
                {
                    NfcFilterAdjust1(&voice->Direct.Params[c].NFCtrlFilter[0], w0);
                    NfcFilterAdjust2(&voice->Direct.Params[c].NFCtrlFilter[1], w0);
                    NfcFilterAdjust3(&voice->Direct.Params[c].NFCtrlFilter[2], w0);
                }

                /* Special-case LFE */
                if(chans[c].channel == LFE)
                {
                    for(j = 0;j < MAX_OUTPUT_CHANNELS;j++)
                        voice->Direct.Params[c].Gains.Target[j] = 0.0f;
                    if(Device->Dry.Buffer == Device->RealOut.Buffer)
                    {
                        int idx = GetChannelIdxByName(Device->RealOut, chans[c].channel);
                        if(idx != -1) voice->Direct.Params[c].Gains.Target[idx] = DryGain;
                    }

                    for(i = 0;i < NumSends;i++)
                    {
                        for(j = 0;j < MAX_EFFECT_CHANNELS;j++)
                            voice->Send[i].Params[c].Gains.Target[j] = 0.0f;
                    }
                    continue;
                }

                if(Device->Render_Mode == StereoPair)
                    CalcAnglePairwiseCoeffs(chans[c].angle, chans[c].elevation, Spread, coeffs);
                else
                    CalcAngleCoeffs(chans[c].angle, chans[c].elevation, Spread, coeffs);
                ComputePanningGains(Device->Dry,
                    coeffs, DryGain, voice->Direct.Params[c].Gains.Target
                );

                for(i = 0;i < NumSends;i++)
                {
                    const ALeffectslot *Slot = SendSlots[i];
                    if(Slot)
                        ComputePanningGainsBF(Slot->ChanMap, Slot->NumChannels,
                            coeffs, WetGain[i], voice->Send[i].Params[c].Gains.Target
                        );
                    else
                        for(j = 0;j < MAX_EFFECT_CHANNELS;j++)
                            voice->Send[i].Params[c].Gains.Target[j] = 0.0f;
                }
            }
        }
    }

    {
        ALfloat hfScale = props->Direct.HFReference / Frequency;
        ALfloat lfScale = props->Direct.LFReference / Frequency;
        ALfloat gainHF = maxf(DryGainHF, 0.001f); /* Limit -60dB */
        ALfloat gainLF = maxf(DryGainLF, 0.001f);

        voice->Direct.FilterType = AF_None;
        if(gainHF != 1.0f) voice->Direct.FilterType |= AF_LowPass;
        if(gainLF != 1.0f) voice->Direct.FilterType |= AF_HighPass;
        ALfilterState_setParams(
            &voice->Direct.Params[0].LowPass, ALfilterType_HighShelf,
            gainHF, hfScale, calc_rcpQ_from_slope(gainHF, 1.0f)
        );
        ALfilterState_setParams(
            &voice->Direct.Params[0].HighPass, ALfilterType_LowShelf,
            gainLF, lfScale, calc_rcpQ_from_slope(gainLF, 1.0f)
        );
        for(c = 1;c < num_channels;c++)
        {
            ALfilterState_copyParams(&voice->Direct.Params[c].LowPass,
                                     &voice->Direct.Params[0].LowPass);
            ALfilterState_copyParams(&voice->Direct.Params[c].HighPass,
                                     &voice->Direct.Params[0].HighPass);
        }
    }
    for(i = 0;i < NumSends;i++)
    {
        ALfloat hfScale = props->Send[i].HFReference / Frequency;
        ALfloat lfScale = props->Send[i].LFReference / Frequency;
        ALfloat gainHF = maxf(WetGainHF[i], 0.001f);
        ALfloat gainLF = maxf(WetGainLF[i], 0.001f);

        voice->Send[i].FilterType = AF_None;
        if(gainHF != 1.0f) voice->Send[i].FilterType |= AF_LowPass;
        if(gainLF != 1.0f) voice->Send[i].FilterType |= AF_HighPass;
        ALfilterState_setParams(
            &voice->Send[i].Params[0].LowPass, ALfilterType_HighShelf,
            gainHF, hfScale, calc_rcpQ_from_slope(gainHF, 1.0f)
        );
        ALfilterState_setParams(
            &voice->Send[i].Params[0].HighPass, ALfilterType_LowShelf,
            gainLF, lfScale, calc_rcpQ_from_slope(gainLF, 1.0f)
        );
        for(c = 1;c < num_channels;c++)
        {
            ALfilterState_copyParams(&voice->Send[i].Params[c].LowPass,
                                     &voice->Send[i].Params[0].LowPass);
            ALfilterState_copyParams(&voice->Send[i].Params[c].HighPass,
                                     &voice->Send[i].Params[0].HighPass);
        }
    }
}

static void CalcNonAttnSourceParams(ALvoice *voice, const struct ALvoiceProps *props, const ALbuffer *ALBuffer, const ALCcontext *ALContext)
{
    static const ALfloat dir[3] = { 0.0f, 0.0f, -1.0f };
    const ALCdevice *Device = ALContext->Device;
    const ALlistener *Listener = ALContext->Listener;
    ALfloat DryGain, DryGainHF, DryGainLF;
    ALfloat WetGain[MAX_SENDS];
    ALfloat WetGainHF[MAX_SENDS];
    ALfloat WetGainLF[MAX_SENDS];
    ALeffectslot *SendSlots[MAX_SENDS];
    ALfloat Pitch;
    ALsizei i;

    voice->Direct.Buffer = Device->Dry.Buffer;
    voice->Direct.Channels = Device->Dry.NumChannels;
    for(i = 0;i < Device->NumAuxSends;i++)
    {
        SendSlots[i] = props->Send[i].Slot;
        if(!SendSlots[i] && i == 0)
            SendSlots[i] = ALContext->DefaultSlot;
        if(!SendSlots[i] || SendSlots[i]->Params.EffectType == AL_EFFECT_NULL)
        {
            SendSlots[i] = NULL;
            voice->Send[i].Buffer = NULL;
            voice->Send[i].Channels = 0;
        }
        else
        {
            voice->Send[i].Buffer = SendSlots[i]->WetBuffer;
            voice->Send[i].Channels = SendSlots[i]->NumChannels;
        }
    }

    /* Calculate the stepping value */
    Pitch = (ALfloat)ALBuffer->Frequency/(ALfloat)Device->Frequency * props->Pitch;
    if(Pitch > (ALfloat)MAX_PITCH)
        voice->Step = MAX_PITCH<<FRACTIONBITS;
    else
        voice->Step = maxi(fastf2i(Pitch*FRACTIONONE + 0.5f), 1);
    BsincPrepare(voice->Step, &voice->ResampleState.bsinc);
    voice->Resampler = SelectResampler(props->Resampler);

    /* Calculate gains */
    DryGain  = clampf(props->Gain, props->MinGain, props->MaxGain);
    DryGain *= props->Direct.Gain * Listener->Params.Gain;
    DryGain  = minf(DryGain, GAIN_MIX_MAX);
    DryGainHF = props->Direct.GainHF;
    DryGainLF = props->Direct.GainLF;
    for(i = 0;i < Device->NumAuxSends;i++)
    {
        WetGain[i]  = clampf(props->Gain, props->MinGain, props->MaxGain);
        WetGain[i] *= props->Send[i].Gain * Listener->Params.Gain;
        WetGain[i]  = minf(WetGain[i], GAIN_MIX_MAX);
        WetGainHF[i] = props->Send[i].GainHF;
        WetGainLF[i] = props->Send[i].GainLF;
    }

    CalcPanningAndFilters(voice, 0.0f, dir, 0.0f, DryGain, DryGainHF, DryGainLF, WetGain,
                          WetGainLF, WetGainHF, SendSlots, ALBuffer, props, Listener, Device);
}

static void CalcAttnSourceParams(ALvoice *voice, const struct ALvoiceProps *props, const ALbuffer *ALBuffer, const ALCcontext *ALContext)
{
    const ALCdevice *Device = ALContext->Device;
    const ALlistener *Listener = ALContext->Listener;
    const ALsizei NumSends = Device->NumAuxSends;
    aluVector Position, Velocity, Direction, SourceToListener;
    ALfloat Distance, ClampedDist, DopplerFactor;
    ALeffectslot *SendSlots[MAX_SENDS];
    ALfloat RoomRolloff[MAX_SENDS];
    ALfloat DecayDistance[MAX_SENDS];
    ALfloat DecayHFDistance[MAX_SENDS];
    ALfloat DryGain, DryGainHF, DryGainLF;
    ALfloat WetGain[MAX_SENDS];
    ALfloat WetGainHF[MAX_SENDS];
    ALfloat WetGainLF[MAX_SENDS];
    bool directional;
    ALfloat dir[3];
    ALfloat spread;
    ALfloat Pitch;
    ALint i;

    /* Set mixing buffers and get send parameters. */
    voice->Direct.Buffer = Device->Dry.Buffer;
    voice->Direct.Channels = Device->Dry.NumChannels;
    for(i = 0;i < NumSends;i++)
    {
        SendSlots[i] = props->Send[i].Slot;
        if(!SendSlots[i] && i == 0)
            SendSlots[i] = ALContext->DefaultSlot;
        if(!SendSlots[i] || SendSlots[i]->Params.EffectType == AL_EFFECT_NULL)
        {
            SendSlots[i] = NULL;
            RoomRolloff[i] = 0.0f;
            DecayDistance[i] = 0.0f;
            DecayHFDistance[i] = 0.0f;
        }
        else if(SendSlots[i]->Params.AuxSendAuto)
        {
            RoomRolloff[i] = SendSlots[i]->Params.RoomRolloff + props->RoomRolloffFactor;
            DecayDistance[i] = SendSlots[i]->Params.DecayTime * SPEEDOFSOUNDMETRESPERSEC;
            DecayHFDistance[i] = DecayDistance[i] * SendSlots[i]->Params.DecayHFRatio;
            if(SendSlots[i]->Params.DecayHFLimit)
            {
                ALfloat airAbsorption = SendSlots[i]->Params.AirAbsorptionGainHF;
                if(airAbsorption < 1.0f)
                {
                    ALfloat limitRatio = log10f(REVERB_DECAY_GAIN) / log10f(airAbsorption);
                    DecayHFDistance[i] = minf(limitRatio, DecayHFDistance[i]);
                }
            }
        }
        else
        {
            /* If the slot's auxiliary send auto is off, the data sent to the
             * effect slot is the same as the dry path, sans filter effects */
            RoomRolloff[i] = props->RolloffFactor;
            DecayDistance[i] = 0.0f;
            DecayHFDistance[i] = 0.0f;
        }

        if(!SendSlots[i])
        {
            voice->Send[i].Buffer = NULL;
            voice->Send[i].Channels = 0;
        }
        else
        {
            voice->Send[i].Buffer = SendSlots[i]->WetBuffer;
            voice->Send[i].Channels = SendSlots[i]->NumChannels;
        }
    }

    /* Transform source to listener space (convert to head relative) */
    aluVectorSet(&Position, props->Position[0], props->Position[1], props->Position[2], 1.0f);
    aluVectorSet(&Direction, props->Direction[0], props->Direction[1], props->Direction[2], 0.0f);
    aluVectorSet(&Velocity, props->Velocity[0], props->Velocity[1], props->Velocity[2], 0.0f);
    if(props->HeadRelative == AL_FALSE)
    {
        const aluMatrixf *Matrix = &Listener->Params.Matrix;
        /* Transform source vectors */
        Position = aluMatrixfVector(Matrix, &Position);
        Velocity = aluMatrixfVector(Matrix, &Velocity);
        Direction = aluMatrixfVector(Matrix, &Direction);
    }
    else
    {
        const aluVector *lvelocity = &Listener->Params.Velocity;
        /* Offset the source velocity to be relative of the listener velocity */
        Velocity.v[0] += lvelocity->v[0];
        Velocity.v[1] += lvelocity->v[1];
        Velocity.v[2] += lvelocity->v[2];
    }

    directional = aluNormalize(Direction.v) > FLT_EPSILON;
    SourceToListener.v[0] = -Position.v[0];
    SourceToListener.v[1] = -Position.v[1];
    SourceToListener.v[2] = -Position.v[2];
    SourceToListener.v[3] = 0.0f;
    Distance = aluNormalize(SourceToListener.v);

    /* Initial source gain */
    DryGain = props->Gain;
    DryGainHF = 1.0f;
    DryGainLF = 1.0f;
    for(i = 0;i < NumSends;i++)
    {
        WetGain[i] = props->Gain;
        WetGainHF[i] = 1.0f;
        WetGainLF[i] = 1.0f;
    }

    /* Calculate distance attenuation */
    ClampedDist = Distance;

    switch(Listener->Params.SourceDistanceModel ?
           props->DistanceModel : Listener->Params.DistanceModel)
    {
        case InverseDistanceClamped:
            ClampedDist = clampf(ClampedDist, props->RefDistance, props->MaxDistance);
            if(props->MaxDistance < props->RefDistance)
                break;
            /*fall-through*/
        case InverseDistance:
            if(!(props->RefDistance > 0.0f))
                ClampedDist = props->RefDistance;
            else
            {
                ALfloat dist = lerp(props->RefDistance, ClampedDist, props->RolloffFactor);
                if(dist > 0.0f) DryGain *= props->RefDistance / dist;
                for(i = 0;i < NumSends;i++)
                {
                    dist = lerp(props->RefDistance, ClampedDist, RoomRolloff[i]);
                    if(dist > 0.0f) WetGain[i] *= props->RefDistance / dist;
                }
            }
            break;

        case LinearDistanceClamped:
            ClampedDist = clampf(ClampedDist, props->RefDistance, props->MaxDistance);
            if(props->MaxDistance < props->RefDistance)
                break;
            /*fall-through*/
        case LinearDistance:
            if(!(props->MaxDistance != props->RefDistance))
                ClampedDist = props->RefDistance;
            else
            {
                ALfloat attn = props->RolloffFactor * (ClampedDist-props->RefDistance) /
                               (props->MaxDistance-props->RefDistance);
                DryGain *= maxf(1.0f - attn, 0.0f);
                for(i = 0;i < NumSends;i++)
                {
                    attn = RoomRolloff[i] * (ClampedDist-props->RefDistance) /
                           (props->MaxDistance-props->RefDistance);
                    WetGain[i] *= maxf(1.0f - attn, 0.0f);
                }
            }
            break;

        case ExponentDistanceClamped:
            ClampedDist = clampf(ClampedDist, props->RefDistance, props->MaxDistance);
            if(props->MaxDistance < props->RefDistance)
                break;
            /*fall-through*/
        case ExponentDistance:
            if(!(ClampedDist > 0.0f && props->RefDistance > 0.0f))
                ClampedDist = props->RefDistance;
            else
            {
                DryGain *= powf(ClampedDist/props->RefDistance, -props->RolloffFactor);
                for(i = 0;i < NumSends;i++)
                    WetGain[i] *= powf(ClampedDist/props->RefDistance, -RoomRolloff[i]);
            }
            break;

        case DisableDistance:
            ClampedDist = props->RefDistance;
            break;
    }

    /* Distance-based air absorption */
    if(ClampedDist > props->RefDistance && props->RolloffFactor > 0.0f)
    {
        ALfloat meters_base = (ClampedDist-props->RefDistance) * props->RolloffFactor *
                              Listener->Params.MetersPerUnit;
        if(props->AirAbsorptionFactor > 0.0f)
        {
            ALfloat hfattn = powf(AIRABSORBGAINHF, meters_base * props->AirAbsorptionFactor);
            DryGainHF *= hfattn;
            for(i = 0;i < NumSends;i++)
                WetGainHF[i] *= hfattn;
        }

        if(props->WetGainAuto)
        {
            /* Apply a decay-time transformation to the wet path, based on the
             * source distance in meters. The initial decay of the reverb
             * effect is calculated and applied to the wet path.
             */
            for(i = 0;i < NumSends;i++)
            {
                ALfloat gain;

                if(!(DecayDistance[i] > 0.0f))
                    continue;

                gain = powf(REVERB_DECAY_GAIN, meters_base/DecayDistance[i]);
                WetGain[i] *= gain;
                /* Yes, the wet path's air absorption is applied with
                 * WetGainAuto on, rather than WetGainHFAuto.
                 */
                if(gain > 0.0f)
                {
                    ALfloat gainhf = powf(REVERB_DECAY_GAIN, meters_base/DecayHFDistance[i]);
                    WetGainHF[i] *= minf(gainhf / gain, 1.0f);
                }
            }
        }
    }

    /* Calculate directional soundcones */
    if(directional && props->InnerAngle < 360.0f)
    {
        ALfloat ConeVolume;
        ALfloat ConeHF;
        ALfloat Angle;

        Angle = acosf(aluDotproduct(&Direction, &SourceToListener));
        Angle = RAD2DEG(Angle * ConeScale * 2.0f);
        if(!(Angle > props->InnerAngle))
        {
            ConeVolume = 1.0f;
            ConeHF = 1.0f;
        }
        else if(Angle < props->OuterAngle)
        {
            ALfloat scale = (            Angle-props->InnerAngle) /
                            (props->OuterAngle-props->InnerAngle);
            ConeVolume = lerp(1.0f, props->OuterGain, scale);
            ConeHF = lerp(1.0f, props->OuterGainHF, scale);
        }
        else
        {
            ConeVolume = props->OuterGain;
            ConeHF = props->OuterGainHF;
        }

        DryGain *= ConeVolume;
        if(props->DryGainHFAuto)
            DryGainHF *= ConeHF;
        if(props->WetGainAuto)
        {
            for(i = 0;i < NumSends;i++)
                WetGain[i] *= ConeVolume;
        }
        if(props->WetGainHFAuto)
        {
            for(i = 0;i < NumSends;i++)
                WetGainHF[i] *= ConeHF;
        }
    }

    /* Apply gain and frequency filters */
    DryGain  = clampf(DryGain, props->MinGain, props->MaxGain);
    DryGain  = minf(DryGain*props->Direct.Gain*Listener->Params.Gain, GAIN_MIX_MAX);
    DryGainHF *= props->Direct.GainHF;
    DryGainLF *= props->Direct.GainLF;
    for(i = 0;i < NumSends;i++)
    {
        WetGain[i]  = clampf(WetGain[i], props->MinGain, props->MaxGain);
        WetGain[i]  = minf(WetGain[i]*props->Send[i].Gain*Listener->Params.Gain, GAIN_MIX_MAX);
        WetGainHF[i] *= props->Send[i].GainHF;
        WetGainLF[i] *= props->Send[i].GainLF;
    }


    /* Initial source pitch */
    Pitch = props->Pitch;

    /* Calculate velocity-based doppler effect */
    DopplerFactor = props->DopplerFactor * Listener->Params.DopplerFactor;
    if(DopplerFactor > 0.0f)
    {
        const aluVector *lvelocity = &Listener->Params.Velocity;
        const ALfloat SpeedOfSound = Listener->Params.SpeedOfSound;
        ALfloat vss, vls;

        vss = aluDotproduct(&Velocity, &SourceToListener) * DopplerFactor;
        vls = aluDotproduct(lvelocity, &SourceToListener) * DopplerFactor;

        if(!(vls < SpeedOfSound))
        {
            /* Listener moving away from the source at the speed of sound.
             * Sound waves can't catch it.
             */
            Pitch = 0.0f;
        }
        else if(!(vss < SpeedOfSound))
        {
            /* Source moving toward the listener at the speed of sound. Sound
             * waves bunch up to extreme frequencies.
             */
            Pitch = HUGE_VALF;
        }
        else
        {
            /* Source and listener movement is nominal. Calculate the proper
             * doppler shift.
             */
            Pitch *= (SpeedOfSound-vls) / (SpeedOfSound-vss);
        }
    }

    /* Adjust pitch based on the buffer and output frequencies, and calculate
     * fixed-point stepping value.
     */
    Pitch *= (ALfloat)ALBuffer->Frequency/(ALfloat)Device->Frequency;
    if(Pitch > (ALfloat)MAX_PITCH)
        voice->Step = MAX_PITCH<<FRACTIONBITS;
    else
        voice->Step = maxi(fastf2i(Pitch*FRACTIONONE + 0.5f), 1);
    BsincPrepare(voice->Step, &voice->ResampleState.bsinc);
    voice->Resampler = SelectResampler(props->Resampler);

    if(Distance > FLT_EPSILON)
    {
        dir[0] = -SourceToListener.v[0];
        /* Clamp Y, in case rounding errors caused it to end up outside of
         * -1...+1.
         */
        dir[1] = clampf(-SourceToListener.v[1], -1.0f, 1.0f);
        dir[2] = -SourceToListener.v[2] * ZScale;
    }
    else
    {
        dir[0] =  0.0f;
        dir[1] =  0.0f;
        dir[2] = -1.0f;
    }
    if(props->Radius > Distance)
        spread = F_TAU - Distance/props->Radius*F_PI;
    else if(Distance > FLT_EPSILON)
        spread = asinf(props->Radius / Distance) * 2.0f;
    else
        spread = 0.0f;

    CalcPanningAndFilters(voice, Distance, dir, spread, DryGain, DryGainHF, DryGainLF, WetGain,
                          WetGainLF, WetGainHF, SendSlots, ALBuffer, props, Listener, Device);
}

static void CalcSourceParams(ALvoice *voice, ALCcontext *context, ALboolean force)
{
    ALbufferlistitem *BufferListItem;
    struct ALvoiceProps *props;

    props = ATOMIC_EXCHANGE_PTR(&voice->Update, NULL, almemory_order_acq_rel);
    if(!props && !force) return;

    if(props)
    {
        memcpy(voice->Props, props,
            FAM_SIZE(struct ALvoiceProps, Send, context->Device->NumAuxSends)
        );

        ATOMIC_REPLACE_HEAD(struct ALvoiceProps*, &voice->FreeList, props);
    }
    props = voice->Props;

    BufferListItem = ATOMIC_LOAD(&voice->current_buffer, almemory_order_relaxed);
    while(BufferListItem != NULL)
    {
        const ALbuffer *buffer;
        if((buffer=BufferListItem->buffer) != NULL)
        {
            if(props->SpatializeMode == SpatializeOn ||
               (props->SpatializeMode == SpatializeAuto && buffer->FmtChannels == FmtMono))
                CalcAttnSourceParams(voice, props, buffer, context);
            else
                CalcNonAttnSourceParams(voice, props, buffer, context);
            break;
        }
        BufferListItem = ATOMIC_LOAD(&BufferListItem->next, almemory_order_acquire);
    }
}


static void UpdateContextSources(ALCcontext *ctx, const struct ALeffectslotArray *slots)
{
    ALvoice **voice, **voice_end;
    ALsource *source;
    ALsizei i;

    IncrementRef(&ctx->UpdateCount);
    if(!ATOMIC_LOAD(&ctx->HoldUpdates, almemory_order_acquire))
    {
        ALboolean force = CalcListenerParams(ctx);
        for(i = 0;i < slots->count;i++)
            force |= CalcEffectSlotParams(slots->slot[i], ctx->Device);

        voice = ctx->Voices;
        voice_end = voice + ctx->VoiceCount;
        for(;voice != voice_end;++voice)
        {
            source = ATOMIC_LOAD(&(*voice)->Source, almemory_order_acquire);
            if(source) CalcSourceParams(*voice, ctx, force);
        }
    }
    IncrementRef(&ctx->UpdateCount);
}


static void ApplyDistanceComp(ALfloatBUFFERSIZE *restrict Samples, DistanceComp *distcomp,
                              ALfloat *restrict Values, ALsizei SamplesToDo, ALsizei numchans)
{
    ALsizei i, c;

    Values = ASSUME_ALIGNED(Values, 16);
    for(c = 0;c < numchans;c++)
    {
        ALfloat *restrict inout = ASSUME_ALIGNED(Samples[c], 16);
        const ALfloat gain = distcomp[c].Gain;
        const ALsizei base = distcomp[c].Length;
        ALfloat *restrict distbuf = ASSUME_ALIGNED(distcomp[c].Buffer, 16);

        if(base == 0)
        {
            if(gain < 1.0f)
            {
                for(i = 0;i < SamplesToDo;i++)
                    inout[i] *= gain;
            }
            continue;
        }

        if(SamplesToDo >= base)
        {
            for(i = 0;i < base;i++)
                Values[i] = distbuf[i];
            for(;i < SamplesToDo;i++)
                Values[i] = inout[i-base];
            memcpy(distbuf, &inout[SamplesToDo-base], base*sizeof(ALfloat));
        }
        else
        {
            for(i = 0;i < SamplesToDo;i++)
                Values[i] = distbuf[i];
            memmove(distbuf, distbuf+SamplesToDo, (base-SamplesToDo)*sizeof(ALfloat));
            memcpy(distbuf+base-SamplesToDo, inout, SamplesToDo*sizeof(ALfloat));
        }
        for(i = 0;i < SamplesToDo;i++)
            inout[i] = Values[i]*gain;
    }
}

static void ApplyDither(ALfloatBUFFERSIZE *restrict Samples, ALuint *dither_seed,
                        const ALfloat quant_scale, const ALsizei SamplesToDo,
                        const ALsizei numchans)
{
    const ALfloat invscale = 1.0f / quant_scale;
    ALuint seed = *dither_seed;
    ALsizei c, i;

    /* Dithering. Step 1, generate whitenoise (uniform distribution of random
     * values between -1 and +1). Step 2 is to add the noise to the samples,
     * before rounding and after scaling up to the desired quantization depth.
     */
    for(c = 0;c < numchans;c++)
    {
        ALfloat *restrict samples = Samples[c];
        for(i = 0;i < SamplesToDo;i++)
        {
            ALfloat val = samples[i] * quant_scale;
            ALuint rng0 = dither_rng(&seed);
            ALuint rng1 = dither_rng(&seed);
            val += (ALfloat)(rng0*(1.0/UINT_MAX) - rng1*(1.0/UINT_MAX));
            samples[i] = roundf(val) * invscale;
        }
    }
    *dither_seed = seed;
}


static inline ALfloat Conv_ALfloat(ALfloat val)
{ return val; }
static inline ALint Conv_ALint(ALfloat val)
{
    /* Floats only have a 24-bit mantissa, so [-16777216, +16777216] is the max
     * integer range normalized floats can be safely converted to (a bit of the
     * exponent helps out, effectively giving 25 bits).
     */
    return fastf2i(clampf(val*16777216.0f, -16777216.0f, 16777215.0f))<<7;
}
static inline ALshort Conv_ALshort(ALfloat val)
{ return fastf2i(clampf(val*32768.0f, -32768.0f, 32767.0f)); }
static inline ALbyte Conv_ALbyte(ALfloat val)
{ return fastf2i(clampf(val*128.0f, -128.0f, 127.0f)); }

/* Define unsigned output variations. */
#define DECL_TEMPLATE(T, func, O)                             \
static inline T Conv_##T(ALfloat val) { return func(val)+O; }

DECL_TEMPLATE(ALubyte, Conv_ALbyte, 128)
DECL_TEMPLATE(ALushort, Conv_ALshort, 32768)
DECL_TEMPLATE(ALuint, Conv_ALint, 2147483648u)

#undef DECL_TEMPLATE

#define DECL_TEMPLATE(T, A)                                                   \
static void Write##A(const ALfloatBUFFERSIZE *InBuffer, ALvoid *OutBuffer,    \
                     ALsizei Offset, ALsizei SamplesToDo, ALsizei numchans)   \
{                                                                             \
    ALsizei i, j;                                                             \
    for(j = 0;j < numchans;j++)                                               \
    {                                                                         \
        const ALfloat *restrict in = ASSUME_ALIGNED(InBuffer[j], 16);         \
        T *restrict out = (T*)OutBuffer + Offset*numchans + j;                \
                                                                              \
        for(i = 0;i < SamplesToDo;i++)                                        \
            out[i*numchans] = Conv_##T(in[i]);                                \
    }                                                                         \
}

DECL_TEMPLATE(ALfloat, F32)
DECL_TEMPLATE(ALuint, UI32)
DECL_TEMPLATE(ALint, I32)
DECL_TEMPLATE(ALushort, UI16)
DECL_TEMPLATE(ALshort, I16)
DECL_TEMPLATE(ALubyte, UI8)
DECL_TEMPLATE(ALbyte, I8)

#undef DECL_TEMPLATE


void aluMixData(ALCdevice *device, ALvoid *OutBuffer, ALsizei NumSamples)
{
    ALsizei SamplesToDo;
    ALsizei SamplesDone;
    ALCcontext *ctx;
    ALsizei i, c;

    START_MIXER_MODE();
    for(SamplesDone = 0;SamplesDone < NumSamples;)
    {
        SamplesToDo = mini(NumSamples-SamplesDone, BUFFERSIZE);
        for(c = 0;c < device->Dry.NumChannels;c++)
            memset(device->Dry.Buffer[c], 0, SamplesToDo*sizeof(ALfloat));
        if(device->Dry.Buffer != device->FOAOut.Buffer)
            for(c = 0;c < device->FOAOut.NumChannels;c++)
                memset(device->FOAOut.Buffer[c], 0, SamplesToDo*sizeof(ALfloat));
        if(device->Dry.Buffer != device->RealOut.Buffer)
            for(c = 0;c < device->RealOut.NumChannels;c++)
                memset(device->RealOut.Buffer[c], 0, SamplesToDo*sizeof(ALfloat));

        IncrementRef(&device->MixCount);

        ctx = ATOMIC_LOAD(&device->ContextList, almemory_order_acquire);
        while(ctx)
        {
            const struct ALeffectslotArray *auxslots;

            auxslots = ATOMIC_LOAD(&ctx->ActiveAuxSlots, almemory_order_acquire);
            UpdateContextSources(ctx, auxslots);

            for(i = 0;i < auxslots->count;i++)
            {
                ALeffectslot *slot = auxslots->slot[i];
                for(c = 0;c < slot->NumChannels;c++)
                    memset(slot->WetBuffer[c], 0, SamplesToDo*sizeof(ALfloat));
            }

            /* source processing */
            for(i = 0;i < ctx->VoiceCount;i++)
            {
                ALvoice *voice = ctx->Voices[i];
                ALsource *source = ATOMIC_LOAD(&voice->Source, almemory_order_acquire);
                if(source && ATOMIC_LOAD(&voice->Playing, almemory_order_relaxed) &&
                   voice->Step > 0)
                {
                    if(!MixSource(voice, source, device, SamplesToDo))
                    {
                        ATOMIC_STORE(&voice->Source, NULL, almemory_order_relaxed);
                        ATOMIC_STORE(&voice->Playing, false, almemory_order_release);
                    }
                }
            }

            /* effect slot processing */
            for(i = 0;i < auxslots->count;i++)
            {
                const ALeffectslot *slot = auxslots->slot[i];
                ALeffectState *state = slot->Params.EffectState;
                V(state,process)(SamplesToDo, slot->WetBuffer, state->OutBuffer,
                                 state->OutChannels);
            }

            ctx = ctx->next;
        }

        /* Increment the clock time. Every second's worth of samples is
         * converted and added to clock base so that large sample counts don't
         * overflow during conversion. This also guarantees an exact, stable
         * conversion. */
        device->SamplesDone += SamplesToDo;
        device->ClockBase += (device->SamplesDone/device->Frequency) * DEVICE_CLOCK_RES;
        device->SamplesDone %= device->Frequency;
        IncrementRef(&device->MixCount);

        if(device->HrtfHandle)
        {
            HrtfDirectMixerFunc HrtfMix;
            DirectHrtfState *state;
            int lidx, ridx;

            if(device->AmbiUp)
                ambiup_process(device->AmbiUp,
                    device->Dry.Buffer, device->Dry.NumChannels,
                    SAFE_CONST(ALfloatBUFFERSIZE*,device->FOAOut.Buffer), SamplesToDo
                );

            lidx = GetChannelIdxByName(device->RealOut, FrontLeft);
            ridx = GetChannelIdxByName(device->RealOut, FrontRight);
            assert(lidx != -1 && ridx != -1);

            HrtfMix = SelectHrtfMixer();
            state = device->Hrtf;
            for(c = 0;c < device->Dry.NumChannels;c++)
            {
                HrtfMix(device->RealOut.Buffer[lidx], device->RealOut.Buffer[ridx],
                    device->Dry.Buffer[c], state->Offset, state->IrSize,
                    SAFE_CONST(ALfloat2*,state->Chan[c].Coeffs),
                    state->Chan[c].Values, SamplesToDo
                );
            }
            state->Offset += SamplesToDo;
        }
        else if(device->AmbiDecoder)
        {
            if(device->Dry.Buffer != device->FOAOut.Buffer)
                bformatdec_upSample(device->AmbiDecoder,
                    device->Dry.Buffer, SAFE_CONST(ALfloatBUFFERSIZE*,device->FOAOut.Buffer),
                    device->FOAOut.NumChannels, SamplesToDo
                );
            bformatdec_process(device->AmbiDecoder,
                device->RealOut.Buffer, device->RealOut.NumChannels,
                SAFE_CONST(ALfloatBUFFERSIZE*,device->Dry.Buffer), SamplesToDo
            );
        }
        else if(device->AmbiUp)
        {
            ambiup_process(device->AmbiUp,
                device->RealOut.Buffer, device->RealOut.NumChannels,
                SAFE_CONST(ALfloatBUFFERSIZE*,device->FOAOut.Buffer), SamplesToDo
            );
        }
        else if(device->Uhj_Encoder)
        {
            int lidx = GetChannelIdxByName(device->RealOut, FrontLeft);
            int ridx = GetChannelIdxByName(device->RealOut, FrontRight);
            if(lidx != -1 && ridx != -1)
            {
                /* Encode to stereo-compatible 2-channel UHJ output. */
                EncodeUhj2(device->Uhj_Encoder,
                    device->RealOut.Buffer[lidx], device->RealOut.Buffer[ridx],
                    device->Dry.Buffer, SamplesToDo
                );
            }
        }
        else if(device->Bs2b)
        {
            int lidx = GetChannelIdxByName(device->RealOut, FrontLeft);
            int ridx = GetChannelIdxByName(device->RealOut, FrontRight);
            if(lidx != -1 && ridx != -1)
            {
                /* Apply binaural/crossfeed filter */
                bs2b_cross_feed(device->Bs2b, device->RealOut.Buffer[lidx],
                                device->RealOut.Buffer[ridx], SamplesToDo);
            }
        }

        if(OutBuffer)
        {
            ALfloat (*Buffer)[BUFFERSIZE] = device->RealOut.Buffer;
            ALsizei Channels = device->RealOut.NumChannels;

            /* Use NFCtrlData for temp value storage. */
            ApplyDistanceComp(Buffer, device->ChannelDelay, device->NFCtrlData,
                              SamplesToDo, Channels);

            if(device->Limiter)
                ApplyCompression(device->Limiter, Channels, SamplesToDo, Buffer);

            if(device->DitherDepth > 0.0f)
                ApplyDither(Buffer, &device->DitherSeed, device->DitherDepth, SamplesToDo,
                            Channels);

            switch(device->FmtType)
            {
                case DevFmtByte:
                    WriteI8(Buffer, OutBuffer, SamplesDone, SamplesToDo, Channels);
                    break;
                case DevFmtUByte:
                    WriteUI8(Buffer, OutBuffer, SamplesDone, SamplesToDo, Channels);
                    break;
                case DevFmtShort:
                    WriteI16(Buffer, OutBuffer, SamplesDone, SamplesToDo, Channels);
                    break;
                case DevFmtUShort:
                    WriteUI16(Buffer, OutBuffer, SamplesDone, SamplesToDo, Channels);
                    break;
                case DevFmtInt:
                    WriteI32(Buffer, OutBuffer, SamplesDone, SamplesToDo, Channels);
                    break;
                case DevFmtUInt:
                    WriteUI32(Buffer, OutBuffer, SamplesDone, SamplesToDo, Channels);
                    break;
                case DevFmtFloat:
                    WriteF32(Buffer, OutBuffer, SamplesDone, SamplesToDo, Channels);
                    break;
            }
        }

        SamplesDone += SamplesToDo;
    }
    END_MIXER_MODE();
}


void aluHandleDisconnect(ALCdevice *device)
{
    ALCcontext *ctx;

    device->Connected = ALC_FALSE;

    ctx = ATOMIC_LOAD_SEQ(&device->ContextList);
    while(ctx)
    {
        ALsizei i;
        for(i = 0;i < ctx->VoiceCount;i++)
        {
            ALvoice *voice = ctx->Voices[i];
            ALsource *source;

            source = ATOMIC_EXCHANGE_PTR(&voice->Source, NULL, almemory_order_acq_rel);
            ATOMIC_STORE(&voice->Playing, false, almemory_order_release);

            if(source)
            {
                ALenum playing = AL_PLAYING;
                (void)(ATOMIC_COMPARE_EXCHANGE_STRONG_SEQ(&source->state, &playing, AL_STOPPED));
            }
        }
        ctx->VoiceCount = 0;

        ctx = ctx->next;
    }
}
