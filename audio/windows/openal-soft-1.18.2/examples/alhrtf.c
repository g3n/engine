/*
 * OpenAL HRTF Example
 *
 * Copyright (c) 2015 by Chris Robinson <chris.kcat@gmail.com>
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

/* This file contains an example for selecting an HRTF. */

#include <stdio.h>
#include <assert.h>
#include <math.h>

#include <SDL_sound.h>

#include "AL/al.h"
#include "AL/alc.h"
#include "AL/alext.h"

#include "common/alhelpers.h"


#ifndef M_PI
#define M_PI                         (3.14159265358979323846)
#endif

static LPALCGETSTRINGISOFT alcGetStringiSOFT;
static LPALCRESETDEVICESOFT alcResetDeviceSOFT;

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
    ALCdevice *device;
    ALboolean has_angle_ext;
    ALuint source, buffer;
    const char *soundname;
    const char *hrtfname;
    ALCint hrtf_state;
    ALCint num_hrtf;
    ALdouble angle;
    ALenum state;

    /* Print out usage if no arguments were specified */
    if(argc < 2)
    {
        fprintf(stderr, "Usage: %s [-device <name>] [-hrtf <name>] <soundfile>\n", argv[0]);
        return 1;
    }

    /* Initialize OpenAL, and check for HRTF support. */
    argv++; argc--;
    if(InitAL(&argv, &argc) != 0)
        return 1;

    device = alcGetContextsDevice(alcGetCurrentContext());
    if(!alcIsExtensionPresent(device, "ALC_SOFT_HRTF"))
    {
        fprintf(stderr, "Error: ALC_SOFT_HRTF not supported\n");
        CloseAL();
        return 1;
    }

    /* Define a macro to help load the function pointers. */
#define LOAD_PROC(d, x)  ((x) = alcGetProcAddress((d), #x))
    LOAD_PROC(device, alcGetStringiSOFT);
    LOAD_PROC(device, alcResetDeviceSOFT);
#undef LOAD_PROC

    /* Check for the AL_EXT_STEREO_ANGLES extension to be able to also rotate
     * stereo sources.
     */
    has_angle_ext = alIsExtensionPresent("AL_EXT_STEREO_ANGLES");
    printf("AL_EXT_STEREO_ANGLES%s found\n", has_angle_ext?"":" not");

    /* Check for user-preferred HRTF */
    if(strcmp(argv[0], "-hrtf") == 0)
    {
        hrtfname = argv[1];
        soundname = argv[2];
    }
    else
    {
        hrtfname = NULL;
        soundname = argv[0];
    }

    /* Enumerate available HRTFs, and reset the device using one. */
    alcGetIntegerv(device, ALC_NUM_HRTF_SPECIFIERS_SOFT, 1, &num_hrtf);
    if(!num_hrtf)
        printf("No HRTFs found\n");
    else
    {
        ALCint attr[5];
        ALCint index = -1;
        ALCint i;

        printf("Available HRTFs:\n");
        for(i = 0;i < num_hrtf;i++)
        {
            const ALCchar *name = alcGetStringiSOFT(device, ALC_HRTF_SPECIFIER_SOFT, i);
            printf("    %d: %s\n", i, name);

            /* Check if this is the HRTF the user requested. */
            if(hrtfname && strcmp(name, hrtfname) == 0)
                index = i;
        }

        i = 0;
        attr[i++] = ALC_HRTF_SOFT;
        attr[i++] = ALC_TRUE;
        if(index == -1)
        {
            if(hrtfname)
                printf("HRTF \"%s\" not found\n", hrtfname);
            printf("Using default HRTF...\n");
        }
        else
        {
            printf("Selecting HRTF %d...\n", index);
            attr[i++] = ALC_HRTF_ID_SOFT;
            attr[i++] = index;
        }
        attr[i] = 0;

        if(!alcResetDeviceSOFT(device, attr))
            printf("Failed to reset device: %s\n", alcGetString(device, alcGetError(device)));
    }

    /* Check if HRTF is enabled, and show which is being used. */
    alcGetIntegerv(device, ALC_HRTF_SOFT, 1, &hrtf_state);
    if(!hrtf_state)
        printf("HRTF not enabled!\n");
    else
    {
        const ALchar *name = alcGetString(device, ALC_HRTF_SPECIFIER_SOFT);
        printf("HRTF enabled, using %s\n", name);
    }
    fflush(stdout);

    /* Initialize SDL_sound. */
    Sound_Init();

    /* Load the sound into a buffer. */
    buffer = LoadSound(soundname);
    if(!buffer)
    {
        Sound_Quit();
        CloseAL();
        return 1;
    }

    /* Create the source to play the sound with. */
    source = 0;
    alGenSources(1, &source);
    alSourcei(source, AL_SOURCE_RELATIVE, AL_TRUE);
    alSource3f(source, AL_POSITION, 0.0f, 0.0f, -1.0f);
    alSourcei(source, AL_BUFFER, buffer);
    assert(alGetError()==AL_NO_ERROR && "Failed to setup sound source");

    /* Play the sound until it finishes. */
    angle = 0.0;
    alSourcePlay(source);
    do {
        al_nssleep(10000000);

        /* Rotate the source around the listener by about 1/4 cycle per second,
         * and keep it within -pi...+pi.
         */
        angle += 0.01 * M_PI * 0.5;
        if(angle > M_PI)
            angle -= M_PI*2.0;

        /* This only rotates mono sounds. */
        alSource3f(source, AL_POSITION, (ALfloat)sin(angle), 0.0f, -(ALfloat)cos(angle));

        if(has_angle_ext)
        {
            /* This rotates stereo sounds with the AL_EXT_STEREO_ANGLES
             * extension. Angles are specified counter-clockwise in radians.
             */
            ALfloat angles[2] = { (ALfloat)(M_PI/6.0 - angle), (ALfloat)(-M_PI/6.0 - angle) };
            alSourcefv(source, AL_STEREO_ANGLES, angles);
        }

        alGetSourcei(source, AL_SOURCE_STATE, &state);
    } while(alGetError() == AL_NO_ERROR && state == AL_PLAYING);

    /* All done. Delete resources, and close down SDL_sound and OpenAL. */
    alDeleteSources(1, &source);
    alDeleteBuffers(1, &buffer);

    Sound_Quit();
    CloseAL();

    return 0;
}
