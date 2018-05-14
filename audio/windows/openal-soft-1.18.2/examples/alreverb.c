/*
 * OpenAL Reverb Example
 *
 * Copyright (c) 2012 by Chris Robinson <chris.kcat@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

/* This file contains an example for applying reverb to a sound. */

#include <stdio.h>
#include <assert.h>

#include <SDL_sound.h>

#include "AL/al.h"
#include "AL/alc.h"
#include "AL/alext.h"
#include "AL/efx-presets.h"

#include "common/alhelpers.h"


/* Effect object functions */
static LPALGENEFFECTS alGenEffects;
static LPALDELETEEFFECTS alDeleteEffects;
static LPALISEFFECT alIsEffect;
static LPALEFFECTI alEffecti;
static LPALEFFECTIV alEffectiv;
static LPALEFFECTF alEffectf;
static LPALEFFECTFV alEffectfv;
static LPALGETEFFECTI alGetEffecti;
static LPALGETEFFECTIV alGetEffectiv;
static LPALGETEFFECTF alGetEffectf;
static LPALGETEFFECTFV alGetEffectfv;

/* Auxiliary Effect Slot object functions */
static LPALGENAUXILIARYEFFECTSLOTS alGenAuxiliaryEffectSlots;
static LPALDELETEAUXILIARYEFFECTSLOTS alDeleteAuxiliaryEffectSlots;
static LPALISAUXILIARYEFFECTSLOT alIsAuxiliaryEffectSlot;
static LPALAUXILIARYEFFECTSLOTI alAuxiliaryEffectSloti;
static LPALAUXILIARYEFFECTSLOTIV alAuxiliaryEffectSlotiv;
static LPALAUXILIARYEFFECTSLOTF alAuxiliaryEffectSlotf;
static LPALAUXILIARYEFFECTSLOTFV alAuxiliaryEffectSlotfv;
static LPALGETAUXILIARYEFFECTSLOTI alGetAuxiliaryEffectSloti;
static LPALGETAUXILIARYEFFECTSLOTIV alGetAuxiliaryEffectSlotiv;
static LPALGETAUXILIARYEFFECTSLOTF alGetAuxiliaryEffectSlotf;
static LPALGETAUXILIARYEFFECTSLOTFV alGetAuxiliaryEffectSlotfv;


/* LoadEffect loads the given reverb properties into a new OpenAL effect
 * object, and returns the new effect ID. */
static ALuint LoadEffect(const EFXEAXREVERBPROPERTIES *reverb)
{
    ALuint effect = 0;
    ALenum err;

    /* Create the effect object and check if we can do EAX reverb. */
    alGenEffects(1, &effect);
    if(alGetEnumValue("AL_EFFECT_EAXREVERB") != 0)
    {
        printf("Using EAX Reverb\n");

        /* EAX Reverb is available. Set the EAX effect type then load the
         * reverb properties. */
        alEffecti(effect, AL_EFFECT_TYPE, AL_EFFECT_EAXREVERB);

        alEffectf(effect, AL_EAXREVERB_DENSITY, reverb->flDensity);
        alEffectf(effect, AL_EAXREVERB_DIFFUSION, reverb->flDiffusion);
        alEffectf(effect, AL_EAXREVERB_GAIN, reverb->flGain);
        alEffectf(effect, AL_EAXREVERB_GAINHF, reverb->flGainHF);
        alEffectf(effect, AL_EAXREVERB_GAINLF, reverb->flGainLF);
        alEffectf(effect, AL_EAXREVERB_DECAY_TIME, reverb->flDecayTime);
        alEffectf(effect, AL_EAXREVERB_DECAY_HFRATIO, reverb->flDecayHFRatio);
        alEffectf(effect, AL_EAXREVERB_DECAY_LFRATIO, reverb->flDecayLFRatio);
        alEffectf(effect, AL_EAXREVERB_REFLECTIONS_GAIN, reverb->flReflectionsGain);
        alEffectf(effect, AL_EAXREVERB_REFLECTIONS_DELAY, reverb->flReflectionsDelay);
        alEffectfv(effect, AL_EAXREVERB_REFLECTIONS_PAN, reverb->flReflectionsPan);
        alEffectf(effect, AL_EAXREVERB_LATE_REVERB_GAIN, reverb->flLateReverbGain);
        alEffectf(effect, AL_EAXREVERB_LATE_REVERB_DELAY, reverb->flLateReverbDelay);
        alEffectfv(effect, AL_EAXREVERB_LATE_REVERB_PAN, reverb->flLateReverbPan);
        alEffectf(effect, AL_EAXREVERB_ECHO_TIME, reverb->flEchoTime);
        alEffectf(effect, AL_EAXREVERB_ECHO_DEPTH, reverb->flEchoDepth);
        alEffectf(effect, AL_EAXREVERB_MODULATION_TIME, reverb->flModulationTime);
        alEffectf(effect, AL_EAXREVERB_MODULATION_DEPTH, reverb->flModulationDepth);
        alEffectf(effect, AL_EAXREVERB_AIR_ABSORPTION_GAINHF, reverb->flAirAbsorptionGainHF);
        alEffectf(effect, AL_EAXREVERB_HFREFERENCE, reverb->flHFReference);
        alEffectf(effect, AL_EAXREVERB_LFREFERENCE, reverb->flLFReference);
        alEffectf(effect, AL_EAXREVERB_ROOM_ROLLOFF_FACTOR, reverb->flRoomRolloffFactor);
        alEffecti(effect, AL_EAXREVERB_DECAY_HFLIMIT, reverb->iDecayHFLimit);
    }
    else
    {
        printf("Using Standard Reverb\n");

        /* No EAX Reverb. Set the standard reverb effect type then load the
         * available reverb properties. */
        alEffecti(effect, AL_EFFECT_TYPE, AL_EFFECT_REVERB);

        alEffectf(effect, AL_REVERB_DENSITY, reverb->flDensity);
        alEffectf(effect, AL_REVERB_DIFFUSION, reverb->flDiffusion);
        alEffectf(effect, AL_REVERB_GAIN, reverb->flGain);
        alEffectf(effect, AL_REVERB_GAINHF, reverb->flGainHF);
        alEffectf(effect, AL_REVERB_DECAY_TIME, reverb->flDecayTime);
        alEffectf(effect, AL_REVERB_DECAY_HFRATIO, reverb->flDecayHFRatio);
        alEffectf(effect, AL_REVERB_REFLECTIONS_GAIN, reverb->flReflectionsGain);
        alEffectf(effect, AL_REVERB_REFLECTIONS_DELAY, reverb->flReflectionsDelay);
        alEffectf(effect, AL_REVERB_LATE_REVERB_GAIN, reverb->flLateReverbGain);
        alEffectf(effect, AL_REVERB_LATE_REVERB_DELAY, reverb->flLateReverbDelay);
        alEffectf(effect, AL_REVERB_AIR_ABSORPTION_GAINHF, reverb->flAirAbsorptionGainHF);
        alEffectf(effect, AL_REVERB_ROOM_ROLLOFF_FACTOR, reverb->flRoomRolloffFactor);
        alEffecti(effect, AL_REVERB_DECAY_HFLIMIT, reverb->iDecayHFLimit);
    }

    /* Check if an error occured, and clean up if so. */
    err = alGetError();
    if(err != AL_NO_ERROR)
    {
        fprintf(stderr, "OpenAL error: %s\n", alGetString(err));
        if(alIsEffect(effect))
            alDeleteEffects(1, &effect);
        return 0;
    }

    return effect;
}


/* LoadBuffer loads the named audio file into an OpenAL buffer object, and
 * returns the new buffer ID.
 */
static ALuint LoadSound(const char *filename)
{
    Sound_Sample *sample;
    ALenum err, format;
    ALuint buffer;
    Uint32 slen;

    /* Open the audio file */
    sample = Sound_NewSampleFromFile(filename, NULL, 65536);
    if(!sample)
    {
        fprintf(stderr, "Could not open audio in %s\n", filename);
        return 0;
    }

    /* Get the sound format, and figure out the OpenAL format */
    if(sample->actual.channels == 1)
    {
        if(sample->actual.format == AUDIO_U8)
            format = AL_FORMAT_MONO8;
        else if(sample->actual.format == AUDIO_S16SYS)
            format = AL_FORMAT_MONO16;
        else
        {
            fprintf(stderr, "Unsupported sample format: 0x%04x\n", sample->actual.format);
            Sound_FreeSample(sample);
            return 0;
        }
    }
    else if(sample->actual.channels == 2)
    {
        if(sample->actual.format == AUDIO_U8)
            format = AL_FORMAT_STEREO8;
        else if(sample->actual.format == AUDIO_S16SYS)
            format = AL_FORMAT_STEREO16;
        else
        {
            fprintf(stderr, "Unsupported sample format: 0x%04x\n", sample->actual.format);
            Sound_FreeSample(sample);
            return 0;
        }
    }
    else
    {
        fprintf(stderr, "Unsupported channel count: %d\n", sample->actual.channels);
        Sound_FreeSample(sample);
        return 0;
    }

    /* Decode the whole audio stream to a buffer. */
    slen = Sound_DecodeAll(sample);
    if(!sample->buffer || slen == 0)
    {
        fprintf(stderr, "Failed to read audio from %s\n", filename);
        Sound_FreeSample(sample);
        return 0;
    }

    /* Buffer the audio data into a new buffer object, then free the data and
     * close the file. */
    buffer = 0;
    alGenBuffers(1, &buffer);
    alBufferData(buffer, format, sample->buffer, slen, sample->actual.rate);
    Sound_FreeSample(sample);

    /* Check if an error occured, and clean up if so. */
    err = alGetError();
    if(err != AL_NO_ERROR)
    {
        fprintf(stderr, "OpenAL Error: %s\n", alGetString(err));
        if(buffer && alIsBuffer(buffer))
            alDeleteBuffers(1, &buffer);
        return 0;
    }

    return buffer;
}


int main(int argc, char **argv)
{
    EFXEAXREVERBPROPERTIES reverb = EFX_REVERB_PRESET_GENERIC;
    ALuint source, buffer, effect, slot;
    ALenum state;

    /* Print out usage if no arguments were specified */
    if(argc < 2)
    {
        fprintf(stderr, "Usage: %s [-device <name] <filename>\n", argv[0]);
        return 1;
    }

    /* Initialize OpenAL, and check for EFX support. */
    argv++; argc--;
    if(InitAL(&argv, &argc) != 0)
        return 1;

    if(!alcIsExtensionPresent(alcGetContextsDevice(alcGetCurrentContext()), "ALC_EXT_EFX"))
    {
        fprintf(stderr, "Error: EFX not supported\n");
        CloseAL();
        return 1;
    }

    /* Define a macro to help load the function pointers. */
#define LOAD_PROC(x)  ((x) = alGetProcAddress(#x))
    LOAD_PROC(alGenEffects);
    LOAD_PROC(alDeleteEffects);
    LOAD_PROC(alIsEffect);
    LOAD_PROC(alEffecti);
    LOAD_PROC(alEffectiv);
    LOAD_PROC(alEffectf);
    LOAD_PROC(alEffectfv);
    LOAD_PROC(alGetEffecti);
    LOAD_PROC(alGetEffectiv);
    LOAD_PROC(alGetEffectf);
    LOAD_PROC(alGetEffectfv);

    LOAD_PROC(alGenAuxiliaryEffectSlots);
    LOAD_PROC(alDeleteAuxiliaryEffectSlots);
    LOAD_PROC(alIsAuxiliaryEffectSlot);
    LOAD_PROC(alAuxiliaryEffectSloti);
    LOAD_PROC(alAuxiliaryEffectSlotiv);
    LOAD_PROC(alAuxiliaryEffectSlotf);
    LOAD_PROC(alAuxiliaryEffectSlotfv);
    LOAD_PROC(alGetAuxiliaryEffectSloti);
    LOAD_PROC(alGetAuxiliaryEffectSlotiv);
    LOAD_PROC(alGetAuxiliaryEffectSlotf);
    LOAD_PROC(alGetAuxiliaryEffectSlotfv);
#undef LOAD_PROC

    /* Initialize SDL_sound. */
    Sound_Init();

    /* Load the sound into a buffer. */
    buffer = LoadSound(argv[0]);
    if(!buffer)
    {
        CloseAL();
        Sound_Quit();
        return 1;
    }

    /* Load the reverb into an effect. */
    effect = LoadEffect(&reverb);
    if(!effect)
    {
        alDeleteBuffers(1, &buffer);
        Sound_Quit();
        CloseAL();
        return 1;
    }

    /* Create the effect slot object. This is what "plays" an effect on sources
     * that connect to it. */
    slot = 0;
    alGenAuxiliaryEffectSlots(1, &slot);

    /* Tell the effect slot to use the loaded effect object. Note that the this
     * effectively copies the effect properties. You can modify or delete the
     * effect object afterward without affecting the effect slot.
     */
    alAuxiliaryEffectSloti(slot, AL_EFFECTSLOT_EFFECT, effect);
    assert(alGetError()==AL_NO_ERROR && "Failed to set effect slot");

    /* Create the source to play the sound with. */
    source = 0;
    alGenSources(1, &source);
    alSourcei(source, AL_BUFFER, buffer);

    /* Connect the source to the effect slot. This tells the source to use the
     * effect slot 'slot', on send #0 with the AL_FILTER_NULL filter object.
     */
    alSource3i(source, AL_AUXILIARY_SEND_FILTER, slot, 0, AL_FILTER_NULL);
    assert(alGetError()==AL_NO_ERROR && "Failed to setup sound source");

    /* Play the sound until it finishes. */
    alSourcePlay(source);
    do {
        al_nssleep(10000000);
        alGetSourcei(source, AL_SOURCE_STATE, &state);
    } while(alGetError() == AL_NO_ERROR && state == AL_PLAYING);

    /* All done. Delete resources, and close down SDL_sound and OpenAL. */
    alDeleteSources(1, &source);
    alDeleteAuxiliaryEffectSlots(1, &slot);
    alDeleteEffects(1, &effect);
    alDeleteBuffers(1, &buffer);

    Sound_Quit();
    CloseAL();

    return 0;
}
