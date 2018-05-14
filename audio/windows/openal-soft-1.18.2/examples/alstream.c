/*
 * OpenAL Audio Stream Example
 *
 * Copyright (c) 2011 by Chris Robinson <chris.kcat@gmail.com>
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

/* This file contains a relatively simple streaming audio player. */

#include <string.h>
#include <stdlib.h>
#include <stdio.h>
#include <signal.h>
#include <assert.h>

#include <SDL_sound.h>

#include "AL/al.h"
#include "AL/alc.h"
#include "AL/alext.h"

#include "common/alhelpers.h"


#ifndef SDL_AUDIO_MASK_BITSIZE
#define SDL_AUDIO_MASK_BITSIZE (0xFF)
#endif
#ifndef SDL_AUDIO_BITSIZE
#define SDL_AUDIO_BITSIZE(x) (x & SDL_AUDIO_MASK_BITSIZE)
#endif

/* Define the number of buffers and buffer size (in milliseconds) to use. 4
 * buffers with 200ms each gives a nice per-chunk size, and lets the queue last
 * for almost one second. */
#define NUM_BUFFERS 4
#define BUFFER_TIME_MS 200

typedef struct StreamPlayer {
    /* These are the buffers and source to play out through OpenAL with */
    ALuint buffers[NUM_BUFFERS];
    ALuint source;

    /* Handle for the audio file */
    Sound_Sample *sample;

    /* The format of the output stream */
    ALenum format;
    ALsizei srate;
} StreamPlayer;

static StreamPlayer *NewPlayer(void);
static void DeletePlayer(StreamPlayer *player);
static int OpenPlayerFile(StreamPlayer *player, const char *filename);
static void ClosePlayerFile(StreamPlayer *player);
static int StartPlayer(StreamPlayer *player);
static int UpdatePlayer(StreamPlayer *player);

/* Creates a new player object, and allocates the needed OpenAL source and
 * buffer objects. Error checking is simplified for the purposes of this
 * example, and will cause an abort if needed. */
static StreamPlayer *NewPlayer(void)
{
    StreamPlayer *player;

    player = calloc(1, sizeof(*player));
    assert(player != NULL);

    /* Generate the buffers and source */
    alGenBuffers(NUM_BUFFERS, player->buffers);
    assert(alGetError() == AL_NO_ERROR && "Could not create buffers");

    alGenSources(1, &player->source);
    assert(alGetError() == AL_NO_ERROR && "Could not create source");

    /* Set parameters so mono sources play out the front-center speaker and
     * won't distance attenuate. */
    alSource3i(player->source, AL_POSITION, 0, 0, -1);
    alSourcei(player->source, AL_SOURCE_RELATIVE, AL_TRUE);
    alSourcei(player->source, AL_ROLLOFF_FACTOR, 0);
    assert(alGetError() == AL_NO_ERROR && "Could not set source parameters");

    return player;
}

/* Destroys a player object, deleting the source and buffers. No error handling
 * since these calls shouldn't fail with a properly-made player object. */
static void DeletePlayer(StreamPlayer *player)
{
    ClosePlayerFile(player);

    alDeleteSources(1, &player->source);
    alDeleteBuffers(NUM_BUFFERS, player->buffers);
    if(alGetError() != AL_NO_ERROR)
        fprintf(stderr, "Failed to delete object IDs\n");

    memset(player, 0, sizeof(*player));
    free(player);
}


/* Opens the first audio stream of the named file. If a file is already open,
 * it will be closed first. */
static int OpenPlayerFile(StreamPlayer *player, const char *filename)
{
    Uint32 frame_size;

    ClosePlayerFile(player);

    /* Open the file and get the first stream from it */
    player->sample = Sound_NewSampleFromFile(filename, NULL, 0);
    if(!player->sample)
    {
        fprintf(stderr, "Could not open audio in %s\n", filename);
        goto error;
    }

    /* Get the stream format, and figure out the OpenAL format */
    if(player->sample->actual.channels == 1)
    {
        if(player->sample->actual.format == AUDIO_U8)
            player->format = AL_FORMAT_MONO8;
        else if(player->sample->actual.format == AUDIO_S16SYS)
            player->format = AL_FORMAT_MONO16;
        else
        {
            fprintf(stderr, "Unsupported sample format: 0x%04x\n", player->sample->actual.format);
            goto error;
        }
    }
    else if(player->sample->actual.channels == 2)
    {
        if(player->sample->actual.format == AUDIO_U8)
            player->format = AL_FORMAT_STEREO8;
        else if(player->sample->actual.format == AUDIO_S16SYS)
            player->format = AL_FORMAT_STEREO16;
        else
        {
            fprintf(stderr, "Unsupported sample format: 0x%04x\n", player->sample->actual.format);
            goto error;
        }
    }
    else
    {
        fprintf(stderr, "Unsupported channel count: %d\n", player->sample->actual.channels);
        goto error;
    }
    player->srate = player->sample->actual.rate;

    frame_size = player->sample->actual.channels *
                 SDL_AUDIO_BITSIZE(player->sample->actual.format) / 8;

    /* Set the buffer size, given the desired millisecond length. */
    Sound_SetBufferSize(player->sample, (Uint32)((Uint64)player->srate*BUFFER_TIME_MS/1000) *
                                        frame_size);

    return 1;

error:
    if(player->sample)
        Sound_FreeSample(player->sample);
    player->sample = NULL;

    return 0;
}

/* Closes the audio file stream */
static void ClosePlayerFile(StreamPlayer *player)
{
    if(player->sample)
        Sound_FreeSample(player->sample);
    player->sample = NULL;
}


/* Prebuffers some audio from the file, and starts playing the source */
static int StartPlayer(StreamPlayer *player)
{
    size_t i;

    /* Rewind the source position and clear the buffer queue */
    alSourceRewind(player->source);
    alSourcei(player->source, AL_BUFFER, 0);

    /* Fill the buffer queue */
    for(i = 0;i < NUM_BUFFERS;i++)
    {
        /* Get some data to give it to the buffer */
        Uint32 slen = Sound_Decode(player->sample);
        if(slen == 0) break;

        alBufferData(player->buffers[i], player->format,
                     player->sample->buffer, slen, player->srate);
    }
    if(alGetError() != AL_NO_ERROR)
    {
        fprintf(stderr, "Error buffering for playback\n");
        return 0;
    }

    /* Now queue and start playback! */
    alSourceQueueBuffers(player->source, i, player->buffers);
    alSourcePlay(player->source);
    if(alGetError() != AL_NO_ERROR)
    {
        fprintf(stderr, "Error starting playback\n");
        return 0;
    }

    return 1;
}

static int UpdatePlayer(StreamPlayer *player)
{
    ALint processed, state;

    /* Get relevant source info */
    alGetSourcei(player->source, AL_SOURCE_STATE, &state);
    alGetSourcei(player->source, AL_BUFFERS_PROCESSED, &processed);
    if(alGetError() != AL_NO_ERROR)
    {
        fprintf(stderr, "Error checking source state\n");
        return 0;
    }

    /* Unqueue and handle each processed buffer */
    while(processed > 0)
    {
        ALuint bufid;
        Uint32 slen;

        alSourceUnqueueBuffers(player->source, 1, &bufid);
        processed--;

        if((player->sample->flags&(SOUND_SAMPLEFLAG_EOF|SOUND_SAMPLEFLAG_ERROR)))
            continue;

        /* Read the next chunk of data, refill the buffer, and queue it
         * back on the source */
        slen = Sound_Decode(player->sample);
        if(slen > 0)
        {
            alBufferData(bufid, player->format, player->sample->buffer, slen,
                         player->srate);
            alSourceQueueBuffers(player->source, 1, &bufid);
        }
        if(alGetError() != AL_NO_ERROR)
        {
            fprintf(stderr, "Error buffering data\n");
            return 0;
        }
    }

    /* Make sure the source hasn't underrun */
    if(state != AL_PLAYING && state != AL_PAUSED)
    {
        ALint queued;

        /* If no buffers are queued, playback is finished */
        alGetSourcei(player->source, AL_BUFFERS_QUEUED, &queued);
        if(queued == 0)
            return 0;

        alSourcePlay(player->source);
        if(alGetError() != AL_NO_ERROR)
        {
            fprintf(stderr, "Error restarting playback\n");
            return 0;
        }
    }

    return 1;
}


int main(int argc, char **argv)
{
    StreamPlayer *player;
    int i;

    /* Print out usage if no arguments were specified */
    if(argc < 2)
    {
        fprintf(stderr, "Usage: %s [-device <name>] <filenames...>\n", argv[0]);
        return 1;
    }

    argv++; argc--;
    if(InitAL(&argv, &argc) != 0)
        return 1;

    Sound_Init();

    player = NewPlayer();

    /* Play each file listed on the command line */
    for(i = 0;i < argc;i++)
    {
        const char *namepart;

        if(!OpenPlayerFile(player, argv[i]))
            continue;

        /* Get the name portion, without the path, for display. */
        namepart = strrchr(argv[i], '/');
        if(namepart || (namepart=strrchr(argv[i], '\\')))
            namepart++;
        else
            namepart = argv[i];

        printf("Playing: %s (%s, %dhz)\n", namepart, FormatName(player->format),
               player->srate);
        fflush(stdout);

        if(!StartPlayer(player))
        {
            ClosePlayerFile(player);
            continue;
        }

        while(UpdatePlayer(player))
            al_nssleep(10000000);

        /* All done with this file. Close it and go to the next */
        ClosePlayerFile(player);
    }
    printf("Done.\n");

    /* All files done. Delete the player, and close down SDL_sound and OpenAL */
    DeletePlayer(player);
    player = NULL;

    Sound_Quit();
    CloseAL();

    return 0;
}
