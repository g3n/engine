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
#include <math.h>

#include "AL/al.h"
#include "AL/alc.h"
#include "alMain.h"
#include "alAuxEffectSlot.h"
#include "alThunk.h"
#include "alError.h"
#include "alListener.h"
#include "alSource.h"

#include "almalloc.h"


extern inline void LockEffectSlotsRead(ALCcontext *context);
extern inline void UnlockEffectSlotsRead(ALCcontext *context);
extern inline void LockEffectSlotsWrite(ALCcontext *context);
extern inline void UnlockEffectSlotsWrite(ALCcontext *context);
extern inline struct ALeffectslot *LookupEffectSlot(ALCcontext *context, ALuint id);
extern inline struct ALeffectslot *RemoveEffectSlot(ALCcontext *context, ALuint id);

static UIntMap EffectStateFactoryMap;
static inline ALeffectStateFactory *getFactoryByType(ALenum type)
{
    ALeffectStateFactory* (*getFactory)(void) = LookupUIntMapKey(&EffectStateFactoryMap, type);
    if(getFactory != NULL)
        return getFactory();
    return NULL;
}

static void ALeffectState_IncRef(ALeffectState *state);
static void ALeffectState_DecRef(ALeffectState *state);

#define DO_UPDATEPROPS() do {                                                 \
    if(!ATOMIC_LOAD(&context->DeferUpdates, almemory_order_acquire))          \
        UpdateEffectSlotProps(slot);                                          \
    else                                                                      \
        ATOMIC_FLAG_CLEAR(&slot->PropsClean, almemory_order_release);         \
} while(0)


AL_API ALvoid AL_APIENTRY alGenAuxiliaryEffectSlots(ALsizei n, ALuint *effectslots)
{
    ALCcontext *context;
    ALeffectslot **tmpslots = NULL;
    ALsizei cur;
    ALenum err;

    context = GetContextRef();
    if(!context) return;

    if(!(n >= 0))
        SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);
    tmpslots = al_malloc(DEF_ALIGN, sizeof(ALeffectslot*)*n);

    LockEffectSlotsWrite(context);
    for(cur = 0;cur < n;cur++)
    {
        ALeffectslot *slot = al_calloc(16, sizeof(ALeffectslot));
        err = AL_OUT_OF_MEMORY;
        if(!slot || (err=InitEffectSlot(slot)) != AL_NO_ERROR)
        {
            al_free(slot);
            UnlockEffectSlotsWrite(context);

            alDeleteAuxiliaryEffectSlots(cur, effectslots);
            SET_ERROR_AND_GOTO(context, err, done);
        }

        err = NewThunkEntry(&slot->id);
        if(err == AL_NO_ERROR)
            err = InsertUIntMapEntryNoLock(&context->EffectSlotMap, slot->id, slot);
        if(err != AL_NO_ERROR)
        {
            FreeThunkEntry(slot->id);
            ALeffectState_DecRef(slot->Effect.State);
            if(slot->Params.EffectState)
                ALeffectState_DecRef(slot->Params.EffectState);
            al_free(slot);
            UnlockEffectSlotsWrite(context);

            alDeleteAuxiliaryEffectSlots(cur, effectslots);
            SET_ERROR_AND_GOTO(context, err, done);
        }

        aluInitEffectPanning(slot);

        tmpslots[cur] = slot;
        effectslots[cur] = slot->id;
    }
    if(n > 0)
    {
        struct ALeffectslotArray *curarray = ATOMIC_LOAD(&context->ActiveAuxSlots, almemory_order_acquire);
        struct ALeffectslotArray *newarray = NULL;
        ALsizei newcount = curarray->count + n;
        ALCdevice *device;

        newarray = al_calloc(DEF_ALIGN, FAM_SIZE(struct ALeffectslotArray, slot, newcount));
        newarray->count = newcount;
        memcpy(newarray->slot, tmpslots, sizeof(ALeffectslot*)*n);
        if(curarray)
            memcpy(newarray->slot+n, curarray->slot, sizeof(ALeffectslot*)*curarray->count);

        newarray = ATOMIC_EXCHANGE_PTR(&context->ActiveAuxSlots, newarray,
                                       almemory_order_acq_rel);
        device = context->Device;
        while((ATOMIC_LOAD(&device->MixCount, almemory_order_acquire)&1))
            althrd_yield();
        al_free(newarray);
    }
    UnlockEffectSlotsWrite(context);

done:
    al_free(tmpslots);
    ALCcontext_DecRef(context);
}

AL_API ALvoid AL_APIENTRY alDeleteAuxiliaryEffectSlots(ALsizei n, const ALuint *effectslots)
{
    ALCcontext *context;
    ALeffectslot *slot;
    ALsizei i;

    context = GetContextRef();
    if(!context) return;

    LockEffectSlotsWrite(context);
    if(!(n >= 0))
        SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);
    for(i = 0;i < n;i++)
    {
        if((slot=LookupEffectSlot(context, effectslots[i])) == NULL)
            SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);
        if(ReadRef(&slot->ref) != 0)
            SET_ERROR_AND_GOTO(context, AL_INVALID_OPERATION, done);
    }

    // All effectslots are valid
    if(n > 0)
    {
        struct ALeffectslotArray *curarray = ATOMIC_LOAD(&context->ActiveAuxSlots, almemory_order_acquire);
        struct ALeffectslotArray *newarray = NULL;
        ALsizei newcount = curarray->count - n;
        ALCdevice *device;
        ALsizei j, k;

        assert(newcount >= 0);
        newarray = al_calloc(DEF_ALIGN, FAM_SIZE(struct ALeffectslotArray, slot, newcount));
        newarray->count = newcount;
        for(i = j = 0;i < newarray->count;)
        {
            slot = curarray->slot[j++];
            for(k = 0;k < n;k++)
            {
                if(slot->id == effectslots[k])
                    break;
            }
            if(k == n)
                newarray->slot[i++] = slot;
        }

        newarray = ATOMIC_EXCHANGE_PTR(&context->ActiveAuxSlots, newarray,
                                       almemory_order_acq_rel);
        device = context->Device;
        while((ATOMIC_LOAD(&device->MixCount, almemory_order_acquire)&1))
            althrd_yield();
        al_free(newarray);
    }

    for(i = 0;i < n;i++)
    {
        if((slot=RemoveEffectSlot(context, effectslots[i])) == NULL)
            continue;
        FreeThunkEntry(slot->id);

        DeinitEffectSlot(slot);

        memset(slot, 0, sizeof(*slot));
        al_free(slot);
    }

done:
    UnlockEffectSlotsWrite(context);
    ALCcontext_DecRef(context);
}

AL_API ALboolean AL_APIENTRY alIsAuxiliaryEffectSlot(ALuint effectslot)
{
    ALCcontext *context;
    ALboolean  ret;

    context = GetContextRef();
    if(!context) return AL_FALSE;

    LockEffectSlotsRead(context);
    ret = (LookupEffectSlot(context, effectslot) ? AL_TRUE : AL_FALSE);
    UnlockEffectSlotsRead(context);

    ALCcontext_DecRef(context);

    return ret;
}

AL_API ALvoid AL_APIENTRY alAuxiliaryEffectSloti(ALuint effectslot, ALenum param, ALint value)
{
    ALCdevice *device;
    ALCcontext *context;
    ALeffectslot *slot;
    ALeffect *effect = NULL;
    ALenum err;

    context = GetContextRef();
    if(!context) return;

    WriteLock(&context->PropLock);
    LockEffectSlotsRead(context);
    if((slot=LookupEffectSlot(context, effectslot)) == NULL)
        SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);
    switch(param)
    {
    case AL_EFFECTSLOT_EFFECT:
        device = context->Device;

        LockEffectsRead(device);
        effect = (value ? LookupEffect(device, value) : NULL);
        if(!(value == 0 || effect != NULL))
        {
            UnlockEffectsRead(device);
            SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);
        }
        err = InitializeEffect(device, slot, effect);
        UnlockEffectsRead(device);

        if(err != AL_NO_ERROR)
            SET_ERROR_AND_GOTO(context, err, done);
        break;

    case AL_EFFECTSLOT_AUXILIARY_SEND_AUTO:
        if(!(value == AL_TRUE || value == AL_FALSE))
            SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);
        slot->AuxSendAuto = value;
        break;

    default:
        SET_ERROR_AND_GOTO(context, AL_INVALID_ENUM, done);
    }
    DO_UPDATEPROPS();

done:
    UnlockEffectSlotsRead(context);
    WriteUnlock(&context->PropLock);
    ALCcontext_DecRef(context);
}

AL_API ALvoid AL_APIENTRY alAuxiliaryEffectSlotiv(ALuint effectslot, ALenum param, const ALint *values)
{
    ALCcontext *context;

    switch(param)
    {
    case AL_EFFECTSLOT_EFFECT:
    case AL_EFFECTSLOT_AUXILIARY_SEND_AUTO:
        alAuxiliaryEffectSloti(effectslot, param, values[0]);
        return;
    }

    context = GetContextRef();
    if(!context) return;

    LockEffectSlotsRead(context);
    if(LookupEffectSlot(context, effectslot) == NULL)
        SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);
    switch(param)
    {
    default:
        SET_ERROR_AND_GOTO(context, AL_INVALID_ENUM, done);
    }

done:
    UnlockEffectSlotsRead(context);
    ALCcontext_DecRef(context);
}

AL_API ALvoid AL_APIENTRY alAuxiliaryEffectSlotf(ALuint effectslot, ALenum param, ALfloat value)
{
    ALCcontext *context;
    ALeffectslot *slot;

    context = GetContextRef();
    if(!context) return;

    WriteLock(&context->PropLock);
    LockEffectSlotsRead(context);
    if((slot=LookupEffectSlot(context, effectslot)) == NULL)
        SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);
    switch(param)
    {
    case AL_EFFECTSLOT_GAIN:
        if(!(value >= 0.0f && value <= 1.0f))
            SET_ERROR_AND_GOTO(context, AL_INVALID_VALUE, done);
        slot->Gain = value;
        break;

    default:
        SET_ERROR_AND_GOTO(context, AL_INVALID_ENUM, done);
    }
    DO_UPDATEPROPS();

done:
    UnlockEffectSlotsRead(context);
    WriteUnlock(&context->PropLock);
    ALCcontext_DecRef(context);
}

AL_API ALvoid AL_APIENTRY alAuxiliaryEffectSlotfv(ALuint effectslot, ALenum param, const ALfloat *values)
{
    ALCcontext *context;

    switch(param)
    {
    case AL_EFFECTSLOT_GAIN:
        alAuxiliaryEffectSlotf(effectslot, param, values[0]);
        return;
    }

    context = GetContextRef();
    if(!context) return;

    LockEffectSlotsRead(context);
    if(LookupEffectSlot(context, effectslot) == NULL)
        SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);
    switch(param)
    {
    default:
        SET_ERROR_AND_GOTO(context, AL_INVALID_ENUM, done);
    }

done:
    UnlockEffectSlotsRead(context);
    ALCcontext_DecRef(context);
}

AL_API ALvoid AL_APIENTRY alGetAuxiliaryEffectSloti(ALuint effectslot, ALenum param, ALint *value)
{
    ALCcontext *context;
    ALeffectslot *slot;

    context = GetContextRef();
    if(!context) return;

    LockEffectSlotsRead(context);
    if((slot=LookupEffectSlot(context, effectslot)) == NULL)
        SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);
    switch(param)
    {
    case AL_EFFECTSLOT_AUXILIARY_SEND_AUTO:
        *value = slot->AuxSendAuto;
        break;

    default:
        SET_ERROR_AND_GOTO(context, AL_INVALID_ENUM, done);
    }

done:
    UnlockEffectSlotsRead(context);
    ALCcontext_DecRef(context);
}

AL_API ALvoid AL_APIENTRY alGetAuxiliaryEffectSlotiv(ALuint effectslot, ALenum param, ALint *values)
{
    ALCcontext *context;

    switch(param)
    {
    case AL_EFFECTSLOT_EFFECT:
    case AL_EFFECTSLOT_AUXILIARY_SEND_AUTO:
        alGetAuxiliaryEffectSloti(effectslot, param, values);
        return;
    }

    context = GetContextRef();
    if(!context) return;

    LockEffectSlotsRead(context);
    if(LookupEffectSlot(context, effectslot) == NULL)
        SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);
    switch(param)
    {
    default:
        SET_ERROR_AND_GOTO(context, AL_INVALID_ENUM, done);
    }

done:
    UnlockEffectSlotsRead(context);
    ALCcontext_DecRef(context);
}

AL_API ALvoid AL_APIENTRY alGetAuxiliaryEffectSlotf(ALuint effectslot, ALenum param, ALfloat *value)
{
    ALCcontext *context;
    ALeffectslot *slot;

    context = GetContextRef();
    if(!context) return;

    LockEffectSlotsRead(context);
    if((slot=LookupEffectSlot(context, effectslot)) == NULL)
        SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);
    switch(param)
    {
    case AL_EFFECTSLOT_GAIN:
        *value = slot->Gain;
        break;

    default:
        SET_ERROR_AND_GOTO(context, AL_INVALID_ENUM, done);
    }

done:
    UnlockEffectSlotsRead(context);
    ALCcontext_DecRef(context);
}

AL_API ALvoid AL_APIENTRY alGetAuxiliaryEffectSlotfv(ALuint effectslot, ALenum param, ALfloat *values)
{
    ALCcontext *context;

    switch(param)
    {
    case AL_EFFECTSLOT_GAIN:
        alGetAuxiliaryEffectSlotf(effectslot, param, values);
        return;
    }

    context = GetContextRef();
    if(!context) return;

    LockEffectSlotsRead(context);
    if(LookupEffectSlot(context, effectslot) == NULL)
        SET_ERROR_AND_GOTO(context, AL_INVALID_NAME, done);
    switch(param)
    {
    default:
        SET_ERROR_AND_GOTO(context, AL_INVALID_ENUM, done);
    }

done:
    UnlockEffectSlotsRead(context);
    ALCcontext_DecRef(context);
}


void InitEffectFactoryMap(void)
{
    InitUIntMap(&EffectStateFactoryMap, INT_MAX);

    InsertUIntMapEntry(&EffectStateFactoryMap, AL_EFFECT_NULL, ALnullStateFactory_getFactory);
    InsertUIntMapEntry(&EffectStateFactoryMap, AL_EFFECT_EAXREVERB, ALreverbStateFactory_getFactory);
    InsertUIntMapEntry(&EffectStateFactoryMap, AL_EFFECT_REVERB, ALreverbStateFactory_getFactory);
    InsertUIntMapEntry(&EffectStateFactoryMap, AL_EFFECT_CHORUS, ALchorusStateFactory_getFactory);
    InsertUIntMapEntry(&EffectStateFactoryMap, AL_EFFECT_COMPRESSOR, ALcompressorStateFactory_getFactory);
    InsertUIntMapEntry(&EffectStateFactoryMap, AL_EFFECT_DISTORTION, ALdistortionStateFactory_getFactory);
    InsertUIntMapEntry(&EffectStateFactoryMap, AL_EFFECT_ECHO, ALechoStateFactory_getFactory);
    InsertUIntMapEntry(&EffectStateFactoryMap, AL_EFFECT_EQUALIZER, ALequalizerStateFactory_getFactory);
    InsertUIntMapEntry(&EffectStateFactoryMap, AL_EFFECT_FLANGER, ALflangerStateFactory_getFactory);
    InsertUIntMapEntry(&EffectStateFactoryMap, AL_EFFECT_RING_MODULATOR, ALmodulatorStateFactory_getFactory);
    InsertUIntMapEntry(&EffectStateFactoryMap, AL_EFFECT_DEDICATED_DIALOGUE, ALdedicatedStateFactory_getFactory);
    InsertUIntMapEntry(&EffectStateFactoryMap, AL_EFFECT_DEDICATED_LOW_FREQUENCY_EFFECT, ALdedicatedStateFactory_getFactory);
}

void DeinitEffectFactoryMap(void)
{
    ResetUIntMap(&EffectStateFactoryMap);
}


ALenum InitializeEffect(ALCdevice *Device, ALeffectslot *EffectSlot, ALeffect *effect)
{
    ALenum newtype = (effect ? effect->type : AL_EFFECT_NULL);
    struct ALeffectslotProps *props;
    ALeffectState *State;

    if(newtype != EffectSlot->Effect.Type)
    {
        ALeffectStateFactory *factory;

        factory = getFactoryByType(newtype);
        if(!factory)
        {
            ERR("Failed to find factory for effect type 0x%04x\n", newtype);
            return AL_INVALID_ENUM;
        }
        State = V0(factory,create)();
        if(!State) return AL_OUT_OF_MEMORY;

        START_MIXER_MODE();
        almtx_lock(&Device->BackendLock);
        State->OutBuffer = Device->Dry.Buffer;
        State->OutChannels = Device->Dry.NumChannels;
        if(V(State,deviceUpdate)(Device) == AL_FALSE)
        {
            almtx_unlock(&Device->BackendLock);
            LEAVE_MIXER_MODE();
            ALeffectState_DecRef(State);
            return AL_OUT_OF_MEMORY;
        }
        almtx_unlock(&Device->BackendLock);
        END_MIXER_MODE();

        if(!effect)
        {
            EffectSlot->Effect.Type = AL_EFFECT_NULL;
            memset(&EffectSlot->Effect.Props, 0, sizeof(EffectSlot->Effect.Props));
        }
        else
        {
            EffectSlot->Effect.Type = effect->type;
            EffectSlot->Effect.Props = effect->Props;
        }

        ALeffectState_DecRef(EffectSlot->Effect.State);
        EffectSlot->Effect.State = State;
    }
    else if(effect)
        EffectSlot->Effect.Props = effect->Props;

    /* Remove state references from old effect slot property updates. */
    props = ATOMIC_LOAD_SEQ(&EffectSlot->FreeList);
    while(props)
    {
        if(props->State)
            ALeffectState_DecRef(props->State);
        props->State = NULL;
        props = ATOMIC_LOAD(&props->next, almemory_order_relaxed);
    }

    return AL_NO_ERROR;
}


static void ALeffectState_IncRef(ALeffectState *state)
{
    uint ref;
    ref = IncrementRef(&state->Ref);
    TRACEREF("%p increasing refcount to %u\n", state, ref);
}

static void ALeffectState_DecRef(ALeffectState *state)
{
    uint ref;
    ref = DecrementRef(&state->Ref);
    TRACEREF("%p decreasing refcount to %u\n", state, ref);
    if(ref == 0) DELETE_OBJ(state);
}


void ALeffectState_Construct(ALeffectState *state)
{
    InitRef(&state->Ref, 1);

    state->OutBuffer = NULL;
    state->OutChannels = 0;
}

void ALeffectState_Destruct(ALeffectState *UNUSED(state))
{
}


ALenum InitEffectSlot(ALeffectslot *slot)
{
    ALeffectStateFactory *factory;

    slot->Effect.Type = AL_EFFECT_NULL;

    factory = getFactoryByType(AL_EFFECT_NULL);
    if(!(slot->Effect.State=V0(factory,create)()))
        return AL_OUT_OF_MEMORY;

    slot->Gain = 1.0;
    slot->AuxSendAuto = AL_TRUE;
    ATOMIC_FLAG_TEST_AND_SET(&slot->PropsClean, almemory_order_relaxed);
    InitRef(&slot->ref, 0);

    ATOMIC_INIT(&slot->Update, NULL);
    ATOMIC_INIT(&slot->FreeList, NULL);

    slot->Params.Gain = 1.0f;
    slot->Params.AuxSendAuto = AL_TRUE;
    ALeffectState_IncRef(slot->Effect.State);
    slot->Params.EffectState = slot->Effect.State;
    slot->Params.RoomRolloff = 0.0f;
    slot->Params.DecayTime = 0.0f;
    slot->Params.DecayHFRatio = 0.0f;
    slot->Params.DecayHFLimit = AL_FALSE;
    slot->Params.AirAbsorptionGainHF = 1.0f;

    return AL_NO_ERROR;
}

void DeinitEffectSlot(ALeffectslot *slot)
{
    struct ALeffectslotProps *props;
    size_t count = 0;

    props = ATOMIC_LOAD_SEQ(&slot->Update);
    if(props)
    {
        if(props->State) ALeffectState_DecRef(props->State);
        TRACE("Freed unapplied AuxiliaryEffectSlot update %p\n", props);
        al_free(props);
    }
    props = ATOMIC_LOAD(&slot->FreeList, almemory_order_relaxed);
    while(props)
    {
        struct ALeffectslotProps *next = ATOMIC_LOAD(&props->next, almemory_order_relaxed);
        if(props->State) ALeffectState_DecRef(props->State);
        al_free(props);
        props = next;
        ++count;
    }
    TRACE("Freed "SZFMT" AuxiliaryEffectSlot property object%s\n", count, (count==1)?"":"s");

    ALeffectState_DecRef(slot->Effect.State);
    if(slot->Params.EffectState)
        ALeffectState_DecRef(slot->Params.EffectState);
}

void UpdateEffectSlotProps(ALeffectslot *slot)
{
    struct ALeffectslotProps *props;
    ALeffectState *oldstate;

    /* Get an unused property container, or allocate a new one as needed. */
    props = ATOMIC_LOAD(&slot->FreeList, almemory_order_relaxed);
    if(!props)
        props = al_calloc(16, sizeof(*props));
    else
    {
        struct ALeffectslotProps *next;
        do {
            next = ATOMIC_LOAD(&props->next, almemory_order_relaxed);
        } while(ATOMIC_COMPARE_EXCHANGE_PTR_WEAK(&slot->FreeList, &props, next,
                almemory_order_seq_cst, almemory_order_acquire) == 0);
    }

    /* Copy in current property values. */
    props->Gain = slot->Gain;
    props->AuxSendAuto = slot->AuxSendAuto;

    props->Type = slot->Effect.Type;
    props->Props = slot->Effect.Props;
    /* Swap out any stale effect state object there may be in the container, to
     * delete it.
     */
    ALeffectState_IncRef(slot->Effect.State);
    oldstate = props->State;
    props->State = slot->Effect.State;

    /* Set the new container for updating internal parameters. */
    props = ATOMIC_EXCHANGE_PTR(&slot->Update, props, almemory_order_acq_rel);
    if(props)
    {
        /* If there was an unused update container, put it back in the
         * freelist.
         */
        ATOMIC_REPLACE_HEAD(struct ALeffectslotProps*, &slot->FreeList, props);
    }

    if(oldstate)
        ALeffectState_DecRef(oldstate);
}

void UpdateAllEffectSlotProps(ALCcontext *context)
{
    struct ALeffectslotArray *auxslots;
    ALsizei i;

    LockEffectSlotsRead(context);
    auxslots = ATOMIC_LOAD(&context->ActiveAuxSlots, almemory_order_acquire);
    for(i = 0;i < auxslots->count;i++)
    {
        ALeffectslot *slot = auxslots->slot[i];
        if(!ATOMIC_FLAG_TEST_AND_SET(&slot->PropsClean, almemory_order_acq_rel))
            UpdateEffectSlotProps(slot);
    }
    UnlockEffectSlotsRead(context);
}

ALvoid ReleaseALAuxiliaryEffectSlots(ALCcontext *Context)
{
    ALsizei pos;
    for(pos = 0;pos < Context->EffectSlotMap.size;pos++)
    {
        ALeffectslot *temp = Context->EffectSlotMap.values[pos];
        Context->EffectSlotMap.values[pos] = NULL;

        DeinitEffectSlot(temp);

        FreeThunkEntry(temp->id);
        memset(temp, 0, sizeof(ALeffectslot));
        al_free(temp);
    }
}
