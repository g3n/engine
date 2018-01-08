/*
 * OpenAL Source Latency Example
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

/* This file contains an example for checking the latency of a sound. */

#include <stdio.h>
#include <assert.h>

#include <SDL_sound.h>

#include "AL/al.h"
#include "AL/alc.h"
#include "AL/alext.h"

#include "common/alhelpers.h"


static LPALSOURCEDSOFT alSourcedSOFT;
static LPALSOURCE3DSOFT alSource3dSOFT;
static LPALSOURCEDVSOFT alSourcedvSOFT;
static LPALGETSOURCEDSOFT alGetSourcedSOFT;
static LPALGETSOURCE3DSOFT alGetSource3dSOFT;
static LPALGETSOURCEDVSOFT alGetSourcedvSOFT;
static LPALSOURCEI64SOFT alSourcei64SOFT;
static LPALSOURCE3I64SOFT alSource3i64SOFT;
static LPALSOURCEI64VSOFT alSourcei64vSOFT;
static LPALGETSOURCEI64SOFT alGetSourcei64SOFT;
static LPALGETSOURCE3I64SOFT alGetSource3i64SOFT;
static LPALGETSOURCEI64VSOFT alGetSourcei64vSOFT;

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
    ALuint source, buffer;
    ALdouble offsets[2];
    ALenum state;

    /* Print out usage if no arguments were specified */
    if(argc < 2)
    {
        fprintf(stderr, "Usage: %s [-device <name>] <filename>\n", argv[0]);
        return 1;
    }

    /* Initialize OpenAL, and check for source_latency support. */
    argv++; argc--;
    if(InitAL(&argv, &argc) != 0)
        return 1;

    if(!alIsExtensionPresent("AL_SOFT_source_latency"))
    {
        fprintf(stderr, "Error: AL_SOFT_source_latency not supported\n");
        CloseAL();
        return 1;
    }

    /* Define a macro to help load the function pointers. */
#define LOAD_PROC(x)  ((x) = alGetProcAddress(#x))
    LOAD_PROC(alSourcedSOFT);
    LOAD_PROC(alSource3dSOFT);
    LOAD_PROC(alSourcedvSOFT);
    LOAD_PROC(alGetSourcedSOFT);
    LOAD_PROC(alGetSource3dSOFT);
    LOAD_PROC(alGetSourcedvSOFT);
    LOAD_PROC(alSourcei64SOFT);
    LOAD_PROC(alSource3i64SOFT);
    LOAD_PROC(alSourcei64vSOFT);
    LOAD_PROC(alGetSourcei64SOFT);
    LOAD_PROC(alGetSource3i64SOFT);
    LOAD_PROC(alGetSourcei64vSOFT);
#undef LOAD_PROC

    /* Initialize SDL_sound. */
    Sound_Init();

    /* Load the sound into a buffer. */
    buffer = LoadSound(argv[0]);
    if(!buffer)
    {
        Sound_Quit();
        CloseAL();
        return 1;
    }

    /* Create the source to play the sound with. */
    source = 0;
    alGenSources(1, &source);
    alSourcei(source, AL_BUFFER, buffer);
    assert(alGetError()==AL_NO_ERROR && "Failed to setup sound source");

    /* Play the sound until it finishes. */
    alSourcePlay(source);
    do {
        al_nssleep(10000000);
        alGetSourcei(source, AL_SOURCE_STATE, &state);

        /* Get the source offset and latency. AL_SEC_OFFSET_LATENCY_SOFT will
         * place the offset (in seconds) in offsets[0], and the time until that
         * offset will be heard (in seconds) in offsets[1]. */
        alGetSourcedvSOFT(source, AL_SEC_OFFSET_LATENCY_SOFT, offsets);
        printf("\rOffset: %f - Latency:%3u ms  ", offsets[0], (ALuint)(offsets[1]*1000));
        fflush(stdout);
    } while(alGetError() == AL_NO_ERROR && state == AL_PLAYING);
    printf("\n");

    /* All done. Delete resources, and close down SDL_sound and OpenAL. */
    alDeleteSources(1, &source);
    alDeleteBuffers(1, &buffer);

    Sound_Quit();
    CloseAL();

    return 0;
}
