/*
 * An example showing how to play a stream sync'd to video, using ffmpeg.
 *
 * Requires C++11.
 */

#include <condition_variable>
#include <functional>
#include <algorithm>
#include <iostream>
#include <iomanip>
#include <cstring>
#include <limits>
#include <thread>
#include <chrono>
#include <atomic>
#include <mutex>
#include <deque>
#include <array>

extern "C" {
#include "libavcodec/avcodec.h"
#include "libavformat/avformat.h"
#include "libavformat/avio.h"
#include "libavutil/time.h"
#include "libavutil/pixfmt.h"
#include "libavutil/avstring.h"
#include "libavutil/channel_layout.h"
#include "libswscale/swscale.h"
#include "libswresample/swresample.h"
}

#include "SDL.h"

#include "AL/alc.h"
#include "AL/al.h"
#include "AL/alext.h"

namespace
{

static const std::string AppName("alffplay");

static bool do_direct_out = false;
static bool has_latency_check = false;
static LPALGETSOURCEDVSOFT alGetSourcedvSOFT;

#define AUDIO_BUFFER_TIME 100 /* In milliseconds, per-buffer */
#define AUDIO_BUFFER_QUEUE_SIZE 8 /* Number of buffers to queue */
#define MAX_QUEUE_SIZE (15 * 1024 * 1024) /* Bytes of compressed data to keep queued */
#define AV_SYNC_THRESHOLD 0.01
#define AV_NOSYNC_THRESHOLD 10.0
#define SAMPLE_CORRECTION_MAX_DIFF 0.05
#define AUDIO_DIFF_AVG_NB 20
#define VIDEO_PICTURE_QUEUE_SIZE 16

enum {
    FF_UPDATE_EVENT = SDL_USEREVENT,
    FF_REFRESH_EVENT,
    FF_MOVIE_DONE_EVENT
};

enum {
    AV_SYNC_AUDIO_MASTER,
    AV_SYNC_VIDEO_MASTER,
    AV_SYNC_EXTERNAL_MASTER,

    DEFAULT_AV_SYNC_TYPE = AV_SYNC_EXTERNAL_MASTER
};


struct PacketQueue {
    std::deque<AVPacket> mPackets;
    std::atomic<int> mTotalSize;
    std::atomic<bool> mFinished;
    std::mutex mMutex;
    std::condition_variable mCond;

    PacketQueue() : mTotalSize(0), mFinished(false)
    { }
    ~PacketQueue()
    { clear(); }

    int put(const AVPacket *pkt);
    int peek(AVPacket *pkt, std::atomic<bool> &quit_var);
    void pop();

    void clear();
    void finish();
};


struct MovieState;

struct AudioState {
    MovieState *mMovie;

    AVStream *mStream;
    AVCodecContext *mCodecCtx;

    PacketQueue mQueue;

    /* Used for clock difference average computation */
    struct {
        std::atomic<int> Clocks; /* In microseconds */
        double Accum;
        double AvgCoeff;
        double Threshold;
        int AvgCount;
    } mDiff;

    /* Time (in seconds) of the next sample to be buffered */
    double mCurrentPts;

    /* Decompressed sample frame, and swresample context for conversion */
    AVFrame           *mDecodedFrame;
    struct SwrContext *mSwresCtx;

    /* Conversion format, for what gets fed to Alure */
    int                 mDstChanLayout;
    enum AVSampleFormat mDstSampleFmt;

    /* Storage of converted samples */
    uint8_t *mSamples;
    int mSamplesLen; /* In samples */
    int mSamplesPos;
    int mSamplesMax;

    /* OpenAL format */
    ALenum mFormat;
    ALsizei mFrameSize;

    std::recursive_mutex mSrcMutex;
    ALuint mSource;
    ALuint mBuffers[AUDIO_BUFFER_QUEUE_SIZE];
    ALsizei mBufferIdx;

    AudioState(MovieState *movie)
      : mMovie(movie), mStream(nullptr), mCodecCtx(nullptr)
      , mDiff{{0}, 0.0, 0.0, 0.0, 0}, mCurrentPts(0.0), mDecodedFrame(nullptr)
      , mSwresCtx(nullptr), mDstChanLayout(0), mDstSampleFmt(AV_SAMPLE_FMT_NONE)
      , mSamples(nullptr), mSamplesLen(0), mSamplesPos(0), mSamplesMax(0)
      , mFormat(AL_NONE), mFrameSize(0), mSource(0), mBufferIdx(0)
    {
        for(auto &buf : mBuffers)
            buf = 0;
    }
    ~AudioState()
    {
        if(mSource)
            alDeleteSources(1, &mSource);
        alDeleteBuffers(AUDIO_BUFFER_QUEUE_SIZE, mBuffers);

        av_frame_free(&mDecodedFrame);
        swr_free(&mSwresCtx);

        av_freep(&mSamples);

        avcodec_free_context(&mCodecCtx);
    }

    double getClock();

    int getSync();
    int decodeFrame();
    int readAudio(uint8_t *samples, int length);

    int handler();
};

struct VideoState {
    MovieState *mMovie;

    AVStream *mStream;
    AVCodecContext *mCodecCtx;

    PacketQueue mQueue;

    double  mClock;
    double  mFrameTimer;
    double  mFrameLastPts;
    double  mFrameLastDelay;
    double  mCurrentPts;
    /* time (av_gettime) at which we updated mCurrentPts - used to have running video pts */
    int64_t mCurrentPtsTime;

    /* Decompressed video frame, and swscale context for conversion */
    AVFrame           *mDecodedFrame;
    struct SwsContext *mSwscaleCtx;

    struct Picture {
        SDL_Texture *mImage;
        int mWidth, mHeight; /* Logical image size (actual size may be larger) */
        std::atomic<bool> mUpdated;
        double mPts;

        Picture()
          : mImage(nullptr), mWidth(0), mHeight(0), mUpdated(false), mPts(0.0)
        { }
        ~Picture()
        {
            if(mImage)
                SDL_DestroyTexture(mImage);
            mImage = nullptr;
        }
    };
    std::array<Picture,VIDEO_PICTURE_QUEUE_SIZE> mPictQ;
    size_t mPictQSize, mPictQRead, mPictQWrite;
    std::mutex mPictQMutex;
    std::condition_variable mPictQCond;
    bool mFirstUpdate;
    std::atomic<bool> mEOS;
    std::atomic<bool> mFinalUpdate;

    VideoState(MovieState *movie)
      : mMovie(movie), mStream(nullptr), mCodecCtx(nullptr), mClock(0.0)
      , mFrameTimer(0.0), mFrameLastPts(0.0), mFrameLastDelay(0.0)
      , mCurrentPts(0.0), mCurrentPtsTime(0), mDecodedFrame(nullptr)
      , mSwscaleCtx(nullptr), mPictQSize(0), mPictQRead(0), mPictQWrite(0)
      , mFirstUpdate(true), mEOS(false), mFinalUpdate(false)
    { }
    ~VideoState()
    {
        sws_freeContext(mSwscaleCtx);
        mSwscaleCtx = nullptr;
        av_frame_free(&mDecodedFrame);
        avcodec_free_context(&mCodecCtx);
    }

    double getClock();

    static Uint32 SDLCALL sdl_refresh_timer_cb(Uint32 interval, void *opaque);
    void schedRefresh(int delay);
    void display(SDL_Window *screen, SDL_Renderer *renderer);
    void refreshTimer(SDL_Window *screen, SDL_Renderer *renderer);
    void updatePicture(SDL_Window *screen, SDL_Renderer *renderer);
    int queuePicture(double pts);
    double synchronize(double pts);
    int handler();
};

struct MovieState {
    AVFormatContext *mFormatCtx;
    int mVideoStream, mAudioStream;

    int mAVSyncType;

    int64_t mExternalClockBase;

    std::atomic<bool> mQuit;

    AudioState mAudio;
    VideoState mVideo;

    std::thread mParseThread;
    std::thread mAudioThread;
    std::thread mVideoThread;

    std::string mFilename;

    MovieState(std::string fname)
      : mFormatCtx(nullptr), mVideoStream(0), mAudioStream(0)
      , mAVSyncType(DEFAULT_AV_SYNC_TYPE), mExternalClockBase(0), mQuit(false)
      , mAudio(this), mVideo(this), mFilename(std::move(fname))
    { }
    ~MovieState()
    {
        mQuit = true;
        if(mParseThread.joinable())
            mParseThread.join();
        avformat_close_input(&mFormatCtx);
    }

    static int decode_interrupt_cb(void *ctx);
    bool prepare();
    void setTitle(SDL_Window *window);

    double getClock();

    double getMasterClock();

    int streamComponentOpen(int stream_index);
    int parse_handler();
};


int PacketQueue::put(const AVPacket *pkt)
{
    std::unique_lock<std::mutex> lock(mMutex);
    mPackets.push_back(AVPacket{});
    if(av_packet_ref(&mPackets.back(), pkt) != 0)
    {
        mPackets.pop_back();
        return -1;
    }
    mTotalSize += mPackets.back().size;
    lock.unlock();

    mCond.notify_one();
    return 0;
}

int PacketQueue::peek(AVPacket *pkt, std::atomic<bool> &quit_var)
{
    std::unique_lock<std::mutex> lock(mMutex);
    while(!quit_var.load())
    {
        if(!mPackets.empty())
        {
            if(av_packet_ref(pkt, &mPackets.front()) != 0)
                return -1;
            return 1;
        }

        if(mFinished.load())
            return 0;
        mCond.wait(lock);
    }
    return -1;
}

void PacketQueue::pop()
{
    std::unique_lock<std::mutex> lock(mMutex);
    AVPacket *pkt = &mPackets.front();
    mTotalSize -= pkt->size;
    av_packet_unref(pkt);
    mPackets.pop_front();
}

void PacketQueue::clear()
{
    std::unique_lock<std::mutex> lock(mMutex);
    std::for_each(mPackets.begin(), mPackets.end(),
        [](AVPacket &pkt) { av_packet_unref(&pkt); }
    );
    mPackets.clear();
    mTotalSize = 0;
}
void PacketQueue::finish()
{
    std::unique_lock<std::mutex> lock(mMutex);
    mFinished = true;
    lock.unlock();
    mCond.notify_all();
}


double AudioState::getClock()
{
    double pts;

    std::unique_lock<std::recursive_mutex> lock(mSrcMutex);
    /* The audio clock is the timestamp of the sample currently being heard.
     * It's based on 4 components:
     * 1 - The timestamp of the next sample to buffer (state->current_pts)
     * 2 - The length of the source's buffer queue
     * 3 - The offset OpenAL is currently at in the source (the first value
     *     from AL_SEC_OFFSET_LATENCY_SOFT)
     * 4 - The latency between OpenAL and the DAC (the second value from
     *     AL_SEC_OFFSET_LATENCY_SOFT)
     *
     * Subtracting the length of the source queue from the next sample's
     * timestamp gives the timestamp of the sample at start of the source
     * queue. Adding the source offset to that results in the timestamp for
     * OpenAL's current position, and subtracting the source latency from that
     * gives the timestamp of the sample currently at the DAC.
     */
    pts = mCurrentPts;
    if(mSource)
    {
        ALdouble offset[2];
        ALint queue_size;
        ALint status;

        /* NOTE: The source state must be checked last, in case an underrun
         * occurs and the source stops between retrieving the offset+latency
         * and getting the state. */
        if(has_latency_check)
        {
            alGetSourcedvSOFT(mSource, AL_SEC_OFFSET_LATENCY_SOFT, offset);
            alGetSourcei(mSource, AL_BUFFERS_QUEUED, &queue_size);
        }
        else
        {
            ALint ioffset;
            alGetSourcei(mSource, AL_SAMPLE_OFFSET, &ioffset);
            alGetSourcei(mSource, AL_BUFFERS_QUEUED, &queue_size);
            offset[0] = (double)ioffset / (double)mCodecCtx->sample_rate;
            offset[1] = 0.0f;
        }
        alGetSourcei(mSource, AL_SOURCE_STATE, &status);

        /* If the source is AL_STOPPED, then there was an underrun and all
         * buffers are processed, so ignore the source queue. The audio thread
         * will put the source into an AL_INITIAL state and clear the queue
         * when it starts recovery. */
        if(status != AL_STOPPED)
            pts -= queue_size*((double)AUDIO_BUFFER_TIME/1000.0) - offset[0];
        if(status == AL_PLAYING)
            pts -= offset[1];
    }
    lock.unlock();

    return std::max(pts, 0.0);
}

int AudioState::getSync()
{
    double diff, avg_diff, ref_clock;

    if(mMovie->mAVSyncType == AV_SYNC_AUDIO_MASTER)
        return 0;

    ref_clock = mMovie->getMasterClock();
    diff = ref_clock - getClock();

    if(!(fabs(diff) < AV_NOSYNC_THRESHOLD))
    {
        /* Difference is TOO big; reset diff stuff */
        mDiff.Accum = 0.0;
        return 0;
    }

    /* Accumulate the diffs */
    mDiff.Accum = mDiff.Accum*mDiff.AvgCoeff + diff;
    avg_diff = mDiff.Accum*(1.0 - mDiff.AvgCoeff);
    if(fabs(avg_diff) < mDiff.Threshold)
        return 0;

    /* Constrain the per-update difference to avoid exceedingly large skips */
    if(!(diff <= SAMPLE_CORRECTION_MAX_DIFF))
        diff = SAMPLE_CORRECTION_MAX_DIFF;
    else if(!(diff >= -SAMPLE_CORRECTION_MAX_DIFF))
        diff = -SAMPLE_CORRECTION_MAX_DIFF;
    return (int)(diff*mCodecCtx->sample_rate);
}

int AudioState::decodeFrame()
{
    while(!mMovie->mQuit.load())
    {
        while(!mMovie->mQuit.load())
        {
            /* Get the next packet */
            AVPacket pkt{};
            if(mQueue.peek(&pkt, mMovie->mQuit) <= 0)
                return -1;

            int ret = avcodec_send_packet(mCodecCtx, &pkt);
            if(ret != AVERROR(EAGAIN))
            {
                if(ret < 0)
                    std::cerr<< "Failed to send encoded packet: 0x"<<std::hex<<ret<<std::dec <<std::endl;
                mQueue.pop();
            }
            av_packet_unref(&pkt);
            if(ret == 0 || ret == AVERROR(EAGAIN))
                break;
        }

        int ret = avcodec_receive_frame(mCodecCtx, mDecodedFrame);
        if(ret == AVERROR(EAGAIN))
            continue;
        if(ret == AVERROR_EOF || ret < 0)
        {
            std::cerr<< "Failed to decode frame: "<<ret <<std::endl;
            return 0;
        }

        if(mDecodedFrame->nb_samples <= 0)
        {
            av_frame_unref(mDecodedFrame);
            continue;
        }

        /* If provided, update w/ pts */
        int64_t pts = av_frame_get_best_effort_timestamp(mDecodedFrame);
        if(pts != AV_NOPTS_VALUE)
            mCurrentPts = av_q2d(mStream->time_base)*pts;

        if(mDecodedFrame->nb_samples > mSamplesMax)
        {
            av_freep(&mSamples);
            av_samples_alloc(
                &mSamples, nullptr, mCodecCtx->channels,
                mDecodedFrame->nb_samples, mDstSampleFmt, 0
            );
            mSamplesMax = mDecodedFrame->nb_samples;
        }
        /* Return the amount of sample frames converted */
        int data_size = swr_convert(mSwresCtx, &mSamples, mDecodedFrame->nb_samples,
            (const uint8_t**)mDecodedFrame->data, mDecodedFrame->nb_samples
        );

        av_frame_unref(mDecodedFrame);
        return data_size;
    }

    return 0;
}

/* Duplicates the sample at in to out, count times. The frame size is a
 * multiple of the template type size.
 */
template<typename T>
static void sample_dup(uint8_t *out, const uint8_t *in, int count, int frame_size)
{
    const T *sample = reinterpret_cast<const T*>(in);
    T *dst = reinterpret_cast<T*>(out);
    if(frame_size == sizeof(T))
        std::fill_n(dst, count, *sample);
    else
    {
        /* NOTE: frame_size is a multiple of sizeof(T). */
        int type_mult = frame_size / sizeof(T);
        int i = 0;
        std::generate_n(dst, count*type_mult,
            [sample,type_mult,&i]() -> T
            {
                T ret = sample[i];
                i = (i+1)%type_mult;
                return ret;
            }
        );
    }
}


int AudioState::readAudio(uint8_t *samples, int length)
{
    int sample_skip = getSync();
    int audio_size = 0;

    /* Read the next chunk of data, refill the buffer, and queue it
     * on the source */
    length /= mFrameSize;
    while(audio_size < length)
    {
        if(mSamplesLen <= 0 || mSamplesPos >= mSamplesLen)
        {
            int frame_len = decodeFrame();
            if(frame_len <= 0) break;

            mSamplesLen = frame_len;
            mSamplesPos = std::min(mSamplesLen, sample_skip);
            sample_skip -= mSamplesPos;

            mCurrentPts += (double)mSamplesPos / (double)mCodecCtx->sample_rate;
            continue;
        }

        int rem = length - audio_size;
        if(mSamplesPos >= 0)
        {
            int len = mSamplesLen - mSamplesPos;
            if(rem > len) rem = len;
            memcpy(samples, mSamples + mSamplesPos*mFrameSize, rem*mFrameSize);
        }
        else
        {
            rem = std::min(rem, -mSamplesPos);

            /* Add samples by copying the first sample */
            if((mFrameSize&7) == 0)
                sample_dup<uint64_t>(samples, mSamples, rem, mFrameSize);
            else if((mFrameSize&3) == 0)
                sample_dup<uint32_t>(samples, mSamples, rem, mFrameSize);
            else if((mFrameSize&1) == 0)
                sample_dup<uint16_t>(samples, mSamples, rem, mFrameSize);
            else
                sample_dup<uint8_t>(samples, mSamples, rem, mFrameSize);
        }

        mSamplesPos += rem;
        mCurrentPts += (double)rem / mCodecCtx->sample_rate;
        samples += rem*mFrameSize;
        audio_size += rem;
    }

    if(audio_size < length && audio_size > 0)
    {
        int rem = length - audio_size;
        std::fill_n(samples, rem*mFrameSize,
                    (mDstSampleFmt == AV_SAMPLE_FMT_U8) ? 0x80 : 0x00);
        mCurrentPts += (double)rem / mCodecCtx->sample_rate;
        audio_size += rem;
    }

    return audio_size * mFrameSize;
}


int AudioState::handler()
{
    std::unique_lock<std::recursive_mutex> lock(mSrcMutex);
    ALenum fmt;

    /* Find a suitable format for Alure. */
    mDstChanLayout = 0;
    if(mCodecCtx->sample_fmt == AV_SAMPLE_FMT_U8 || mCodecCtx->sample_fmt == AV_SAMPLE_FMT_U8P)
    {
        mDstSampleFmt = AV_SAMPLE_FMT_U8;
        mFrameSize = 1;
        if(mCodecCtx->channel_layout == AV_CH_LAYOUT_7POINT1 &&
           alIsExtensionPresent("AL_EXT_MCFORMATS") &&
           (fmt=alGetEnumValue("AL_FORMAT_71CHN8")) != AL_NONE && fmt != -1)
        {
            mDstChanLayout = mCodecCtx->channel_layout;
            mFrameSize *= 8;
            mFormat = fmt;
        }
        if((mCodecCtx->channel_layout == AV_CH_LAYOUT_5POINT1 ||
            mCodecCtx->channel_layout == AV_CH_LAYOUT_5POINT1_BACK) &&
           alIsExtensionPresent("AL_EXT_MCFORMATS") &&
           (fmt=alGetEnumValue("AL_FORMAT_51CHN8")) != AL_NONE && fmt != -1)
        {
            mDstChanLayout = mCodecCtx->channel_layout;
            mFrameSize *= 6;
            mFormat = fmt;
        }
        if(mCodecCtx->channel_layout == AV_CH_LAYOUT_MONO)
        {
            mDstChanLayout = mCodecCtx->channel_layout;
            mFrameSize *= 1;
            mFormat = AL_FORMAT_MONO8;
        }
        if(!mDstChanLayout)
        {
            mDstChanLayout = AV_CH_LAYOUT_STEREO;
            mFrameSize *= 2;
            mFormat = AL_FORMAT_STEREO8;
        }
    }
    if((mCodecCtx->sample_fmt == AV_SAMPLE_FMT_FLT || mCodecCtx->sample_fmt == AV_SAMPLE_FMT_FLTP) &&
       alIsExtensionPresent("AL_EXT_FLOAT32"))
    {
        mDstSampleFmt = AV_SAMPLE_FMT_FLT;
        mFrameSize = 4;
        if(mCodecCtx->channel_layout == AV_CH_LAYOUT_7POINT1 &&
           alIsExtensionPresent("AL_EXT_MCFORMATS") &&
           (fmt=alGetEnumValue("AL_FORMAT_71CHN32")) != AL_NONE && fmt != -1)
        {
            mDstChanLayout = mCodecCtx->channel_layout;
            mFrameSize *= 8;
            mFormat = fmt;
        }
        if((mCodecCtx->channel_layout == AV_CH_LAYOUT_5POINT1 ||
            mCodecCtx->channel_layout == AV_CH_LAYOUT_5POINT1_BACK) &&
           alIsExtensionPresent("AL_EXT_MCFORMATS") &&
           (fmt=alGetEnumValue("AL_FORMAT_51CHN32")) != AL_NONE && fmt != -1)
        {
            mDstChanLayout = mCodecCtx->channel_layout;
            mFrameSize *= 6;
            mFormat = fmt;
        }
        if(mCodecCtx->channel_layout == AV_CH_LAYOUT_MONO)
        {
            mDstChanLayout = mCodecCtx->channel_layout;
            mFrameSize *= 1;
            mFormat = AL_FORMAT_MONO_FLOAT32;
        }
        if(!mDstChanLayout)
        {
            mDstChanLayout = AV_CH_LAYOUT_STEREO;
            mFrameSize *= 2;
            mFormat = AL_FORMAT_STEREO_FLOAT32;
        }
    }
    if(!mDstChanLayout)
    {
        mDstSampleFmt = AV_SAMPLE_FMT_S16;
        mFrameSize = 2;
        if(mCodecCtx->channel_layout == AV_CH_LAYOUT_7POINT1 &&
           alIsExtensionPresent("AL_EXT_MCFORMATS") &&
           (fmt=alGetEnumValue("AL_FORMAT_71CHN16")) != AL_NONE && fmt != -1)
        {
            mDstChanLayout = mCodecCtx->channel_layout;
            mFrameSize *= 8;
            mFormat = fmt;
        }
        if((mCodecCtx->channel_layout == AV_CH_LAYOUT_5POINT1 ||
            mCodecCtx->channel_layout == AV_CH_LAYOUT_5POINT1_BACK) &&
           alIsExtensionPresent("AL_EXT_MCFORMATS") &&
           (fmt=alGetEnumValue("AL_FORMAT_51CHN16")) != AL_NONE && fmt != -1)
        {
            mDstChanLayout = mCodecCtx->channel_layout;
            mFrameSize *= 6;
            mFormat = fmt;
        }
        if(mCodecCtx->channel_layout == AV_CH_LAYOUT_MONO)
        {
            mDstChanLayout = mCodecCtx->channel_layout;
            mFrameSize *= 1;
            mFormat = AL_FORMAT_MONO16;
        }
        if(!mDstChanLayout)
        {
            mDstChanLayout = AV_CH_LAYOUT_STEREO;
            mFrameSize *= 2;
            mFormat = AL_FORMAT_STEREO16;
        }
    }
    ALsizei buffer_len = mCodecCtx->sample_rate * AUDIO_BUFFER_TIME / 1000 *
            mFrameSize;
    void *samples = av_malloc(buffer_len);

    mSamples = NULL;
    mSamplesMax = 0;
    mSamplesPos = 0;
    mSamplesLen = 0;

    if(!(mDecodedFrame=av_frame_alloc()))
    {
        std::cerr<< "Failed to allocate audio frame" <<std::endl;
        goto finish;
    }

    mSwresCtx = swr_alloc_set_opts(nullptr,
        mDstChanLayout, mDstSampleFmt, mCodecCtx->sample_rate,
        mCodecCtx->channel_layout ? mCodecCtx->channel_layout :
            (uint64_t)av_get_default_channel_layout(mCodecCtx->channels),
        mCodecCtx->sample_fmt, mCodecCtx->sample_rate,
        0, nullptr
    );
    if(!mSwresCtx || swr_init(mSwresCtx) != 0)
    {
        std::cerr<< "Failed to initialize audio converter" <<std::endl;
        goto finish;
    }

    alGenBuffers(AUDIO_BUFFER_QUEUE_SIZE, mBuffers);
    alGenSources(1, &mSource);

    if(do_direct_out)
    {
        if(!alIsExtensionPresent("AL_SOFT_direct_channels"))
            std::cerr<< "AL_SOFT_direct_channels not supported for direct output" <<std::endl;
        else
        {
            alSourcei(mSource, AL_DIRECT_CHANNELS_SOFT, AL_TRUE);
            std::cout<< "Direct out enabled" <<std::endl;
        }
    }

    while(alGetError() == AL_NO_ERROR && !mMovie->mQuit.load())
    {
        /* First remove any processed buffers. */
        ALint processed;
        alGetSourcei(mSource, AL_BUFFERS_PROCESSED, &processed);
        if(processed > 0)
        {
            std::array<ALuint,AUDIO_BUFFER_QUEUE_SIZE> tmp;
            alSourceUnqueueBuffers(mSource, processed, tmp.data());
        }

        /* Refill the buffer queue. */
        ALint queued;
        alGetSourcei(mSource, AL_BUFFERS_QUEUED, &queued);
        while(queued < AUDIO_BUFFER_QUEUE_SIZE)
        {
            int audio_size;

            /* Read the next chunk of data, fill the buffer, and queue it on
             * the source */
            audio_size = readAudio(reinterpret_cast<uint8_t*>(samples), buffer_len);
            if(audio_size <= 0) break;

            ALuint bufid = mBuffers[mBufferIdx++];
            mBufferIdx %= AUDIO_BUFFER_QUEUE_SIZE;

            alBufferData(bufid, mFormat, samples, audio_size, mCodecCtx->sample_rate);
            alSourceQueueBuffers(mSource, 1, &bufid);
            queued++;
        }
        if(queued == 0)
            break;

        /* Check that the source is playing. */
        ALint state;
        alGetSourcei(mSource, AL_SOURCE_STATE, &state);
        if(state == AL_STOPPED)
        {
            /* AL_STOPPED means there was an underrun. Rewind the source to get
             * it back into an AL_INITIAL state.
             */
            alSourceRewind(mSource);
            continue;
        }

        lock.unlock();

        /* (re)start the source if needed, and wait for a buffer to finish */
        if(state != AL_PLAYING && state != AL_PAUSED)
            alSourcePlay(mSource);
        SDL_Delay(AUDIO_BUFFER_TIME / 3);

        lock.lock();
    }

finish:
    alSourceRewind(mSource);
    alSourcei(mSource, AL_BUFFER, 0);

    av_frame_free(&mDecodedFrame);
    swr_free(&mSwresCtx);

    av_freep(&mSamples);

    return 0;
}


double VideoState::getClock()
{
    double delta = (av_gettime() - mCurrentPtsTime) / 1000000.0;
    return mCurrentPts + delta;
}

Uint32 SDLCALL VideoState::sdl_refresh_timer_cb(Uint32 /*interval*/, void *opaque)
{
    SDL_Event evt{};
    evt.user.type = FF_REFRESH_EVENT;
    evt.user.data1 = opaque;
    SDL_PushEvent(&evt);
    return 0; /* 0 means stop timer */
}

/* Schedules an FF_REFRESH_EVENT event to occur in 'delay' ms. */
void VideoState::schedRefresh(int delay)
{
    SDL_AddTimer(delay, sdl_refresh_timer_cb, this);
}

/* Called by VideoState::refreshTimer to display the next video frame. */
void VideoState::display(SDL_Window *screen, SDL_Renderer *renderer)
{
    Picture *vp = &mPictQ[mPictQRead];

    if(!vp->mImage)
        return;

    float aspect_ratio;
    int win_w, win_h;
    int w, h, x, y;

    if(mCodecCtx->sample_aspect_ratio.num == 0)
        aspect_ratio = 0.0f;
    else
    {
        aspect_ratio = av_q2d(mCodecCtx->sample_aspect_ratio) * mCodecCtx->width /
                       mCodecCtx->height;
    }
    if(aspect_ratio <= 0.0f)
        aspect_ratio = (float)mCodecCtx->width / (float)mCodecCtx->height;

    SDL_GetWindowSize(screen, &win_w, &win_h);
    h = win_h;
    w = ((int)rint(h * aspect_ratio) + 3) & ~3;
    if(w > win_w)
    {
        w = win_w;
        h = ((int)rint(w / aspect_ratio) + 3) & ~3;
    }
    x = (win_w - w) / 2;
    y = (win_h - h) / 2;

    SDL_Rect src_rect{ 0, 0, vp->mWidth, vp->mHeight };
    SDL_Rect dst_rect{ x, y, w, h };
    SDL_RenderCopy(renderer, vp->mImage, &src_rect, &dst_rect);
    SDL_RenderPresent(renderer);
}

/* FF_REFRESH_EVENT handler called on the main thread where the SDL_Renderer
 * was created. It handles the display of the next decoded video frame (if not
 * falling behind), and sets up the timer for the following video frame.
 */
void VideoState::refreshTimer(SDL_Window *screen, SDL_Renderer *renderer)
{
    if(!mStream)
    {
        if(mEOS)
        {
            mFinalUpdate = true;
            std::unique_lock<std::mutex>(mPictQMutex).unlock();
            mPictQCond.notify_all();
            return;
        }
        schedRefresh(100);
        return;
    }

    std::unique_lock<std::mutex> lock(mPictQMutex);
retry:
    if(mPictQSize == 0)
    {
        if(mEOS)
            mFinalUpdate = true;
        else
            schedRefresh(1);
        lock.unlock();
        mPictQCond.notify_all();
        return;
    }

    Picture *vp = &mPictQ[mPictQRead];
    mCurrentPts = vp->mPts;
    mCurrentPtsTime = av_gettime();

    /* Get delay using the frame pts and the pts from last frame. */
    double delay = vp->mPts - mFrameLastPts;
    if(delay <= 0 || delay >= 1.0)
    {
        /* If incorrect delay, use previous one. */
        delay = mFrameLastDelay;
    }
    /* Save for next frame. */
    mFrameLastDelay = delay;
    mFrameLastPts = vp->mPts;

    /* Update delay to sync to clock if not master source. */
    if(mMovie->mAVSyncType != AV_SYNC_VIDEO_MASTER)
    {
        double ref_clock = mMovie->getMasterClock();
        double diff = vp->mPts - ref_clock;

        /* Skip or repeat the frame. Take delay into account. */
        double sync_threshold = std::min(delay, AV_SYNC_THRESHOLD);
        if(fabs(diff) < AV_NOSYNC_THRESHOLD)
        {
            if(diff <= -sync_threshold)
                delay = 0;
            else if(diff >= sync_threshold)
                delay *= 2.0;
        }
    }

    mFrameTimer += delay;
    /* Compute the REAL delay. */
    double actual_delay = mFrameTimer - (av_gettime() / 1000000.0);
    if(!(actual_delay >= 0.010))
    {
        /* We don't have time to handle this picture, just skip to the next one. */
        mPictQRead = (mPictQRead+1)%mPictQ.size();
        mPictQSize--;
        goto retry;
    }
    schedRefresh((int)(actual_delay*1000.0 + 0.5));

    /* Show the picture! */
    display(screen, renderer);

    /* Update queue for next picture. */
    mPictQRead = (mPictQRead+1)%mPictQ.size();
    mPictQSize--;
    lock.unlock();
    mPictQCond.notify_all();
}

/* FF_UPDATE_EVENT handler, updates the picture's texture. It's called on the
 * main thread where the renderer was created.
 */
void VideoState::updatePicture(SDL_Window *screen, SDL_Renderer *renderer)
{
    Picture *vp = &mPictQ[mPictQWrite];
    bool fmt_updated = false;

    /* allocate or resize the buffer! */
    if(!vp->mImage || vp->mWidth != mCodecCtx->width || vp->mHeight != mCodecCtx->height)
    {
        fmt_updated = true;
        if(vp->mImage)
            SDL_DestroyTexture(vp->mImage);
        vp->mImage = SDL_CreateTexture(
            renderer, SDL_PIXELFORMAT_IYUV, SDL_TEXTUREACCESS_STREAMING,
            mCodecCtx->coded_width, mCodecCtx->coded_height
        );
        if(!vp->mImage)
            std::cerr<< "Failed to create YV12 texture!" <<std::endl;
        vp->mWidth = mCodecCtx->width;
        vp->mHeight = mCodecCtx->height;

        if(mFirstUpdate && vp->mWidth > 0 && vp->mHeight > 0)
        {
            /* For the first update, set the window size to the video size. */
            mFirstUpdate = false;

            int w = vp->mWidth;
            int h = vp->mHeight;
            if(mCodecCtx->sample_aspect_ratio.den != 0)
            {
                double aspect_ratio = av_q2d(mCodecCtx->sample_aspect_ratio);
                if(aspect_ratio >= 1.0)
                    w = (int)(w*aspect_ratio + 0.5);
                else if(aspect_ratio > 0.0)
                    h = (int)(h/aspect_ratio + 0.5);
            }
            SDL_SetWindowSize(screen, w, h);
        }
    }

    if(vp->mImage)
    {
        AVFrame *frame = mDecodedFrame;
        void *pixels = nullptr;
        int pitch = 0;

        if(mCodecCtx->pix_fmt == AV_PIX_FMT_YUV420P)
            SDL_UpdateYUVTexture(vp->mImage, nullptr,
                frame->data[0], frame->linesize[0],
                frame->data[1], frame->linesize[1],
                frame->data[2], frame->linesize[2]
            );
        else if(SDL_LockTexture(vp->mImage, nullptr, &pixels, &pitch) != 0)
            std::cerr<< "Failed to lock texture" <<std::endl;
        else
        {
            // Convert the image into YUV format that SDL uses
            int coded_w = mCodecCtx->coded_width;
            int coded_h = mCodecCtx->coded_height;
            int w = mCodecCtx->width;
            int h = mCodecCtx->height;
            if(!mSwscaleCtx || fmt_updated)
            {
                sws_freeContext(mSwscaleCtx);
                mSwscaleCtx = sws_getContext(
                    w, h, mCodecCtx->pix_fmt,
                    w, h, AV_PIX_FMT_YUV420P, 0,
                    nullptr, nullptr, nullptr
                );
            }

            /* point pict at the queue */
            uint8_t *pict_data[3];
            pict_data[0] = reinterpret_cast<uint8_t*>(pixels);
            pict_data[1] = pict_data[0] + coded_w*coded_h;
            pict_data[2] = pict_data[1] + coded_w*coded_h/4;

            int pict_linesize[3];
            pict_linesize[0] = pitch;
            pict_linesize[1] = pitch / 2;
            pict_linesize[2] = pitch / 2;

            sws_scale(mSwscaleCtx, (const uint8_t**)frame->data,
                      frame->linesize, 0, h, pict_data, pict_linesize);
            SDL_UnlockTexture(vp->mImage);
        }
    }

    std::unique_lock<std::mutex> lock(mPictQMutex);
    vp->mUpdated = true;
    lock.unlock();
    mPictQCond.notify_one();
}

int VideoState::queuePicture(double pts)
{
    /* Wait until we have space for a new pic */
    std::unique_lock<std::mutex> lock(mPictQMutex);
    while(mPictQSize >= mPictQ.size() && !mMovie->mQuit.load())
        mPictQCond.wait(lock);
    lock.unlock();

    if(mMovie->mQuit.load())
        return -1;

    Picture *vp = &mPictQ[mPictQWrite];

    /* We have to create/update the picture in the main thread  */
    vp->mUpdated = false;
    SDL_Event evt{};
    evt.user.type = FF_UPDATE_EVENT;
    evt.user.data1 = this;
    SDL_PushEvent(&evt);

    /* Wait until the picture is updated. */
    lock.lock();
    while(!vp->mUpdated && !mMovie->mQuit.load())
        mPictQCond.wait(lock);
    if(mMovie->mQuit.load())
        return -1;
    vp->mPts = pts;

    mPictQWrite = (mPictQWrite+1)%mPictQ.size();
    mPictQSize++;
    lock.unlock();

    return 0;
}

double VideoState::synchronize(double pts)
{
    double frame_delay;

    if(pts == 0.0) /* if we aren't given a pts, set it to the clock */
        pts = mClock;
    else /* if we have pts, set video clock to it */
        mClock = pts;

    /* update the video clock */
    frame_delay = av_q2d(mCodecCtx->time_base);
    /* if we are repeating a frame, adjust clock accordingly */
    frame_delay += mDecodedFrame->repeat_pict * (frame_delay * 0.5);
    mClock += frame_delay;
    return pts;
}

int VideoState::handler()
{
    mDecodedFrame = av_frame_alloc();
    while(!mMovie->mQuit)
    {
        while(!mMovie->mQuit)
        {
            AVPacket packet{};
            if(mQueue.peek(&packet, mMovie->mQuit) <= 0)
                goto finish;

            int ret = avcodec_send_packet(mCodecCtx, &packet);
            if(ret != AVERROR(EAGAIN))
            {
                if(ret < 0)
                    std::cerr<< "Failed to send encoded packet: 0x"<<std::hex<<ret<<std::dec <<std::endl;
                mQueue.pop();
            }
            av_packet_unref(&packet);
            if(ret == 0 || ret == AVERROR(EAGAIN))
                break;
        }

        /* Decode video frame */
        int ret = avcodec_receive_frame(mCodecCtx, mDecodedFrame);
        if(ret == AVERROR(EAGAIN))
            continue;
        if(ret < 0)
        {
            std::cerr<< "Failed to decode frame: "<<ret <<std::endl;
            break;
        }

        double pts = synchronize(
            av_q2d(mStream->time_base) * av_frame_get_best_effort_timestamp(mDecodedFrame)
        );
        if(queuePicture(pts) < 0)
            break;
        av_frame_unref(mDecodedFrame);
    }
finish:
    mEOS = true;
    av_frame_free(&mDecodedFrame);

    std::unique_lock<std::mutex> lock(mPictQMutex);
    if(mMovie->mQuit)
    {
        mPictQRead = 0;
        mPictQWrite = 0;
        mPictQSize = 0;
    }
    while(!mFinalUpdate)
        mPictQCond.wait(lock);

    return 0;
}


int MovieState::decode_interrupt_cb(void *ctx)
{
    return reinterpret_cast<MovieState*>(ctx)->mQuit;
}

bool MovieState::prepare()
{
    mFormatCtx = avformat_alloc_context();
    mFormatCtx->interrupt_callback.callback = decode_interrupt_cb;
    mFormatCtx->interrupt_callback.opaque = this;
    if(avio_open2(&mFormatCtx->pb, mFilename.c_str(), AVIO_FLAG_READ,
                  &mFormatCtx->interrupt_callback, nullptr))
    {
        std::cerr<< "Failed to open "<<mFilename <<std::endl;
        return false;
    }

    /* Open movie file */
    if(avformat_open_input(&mFormatCtx, mFilename.c_str(), nullptr, nullptr) != 0)
    {
        std::cerr<< "Failed to open "<<mFilename <<std::endl;
        return false;
    }

    /* Retrieve stream information */
    if(avformat_find_stream_info(mFormatCtx, nullptr) < 0)
    {
        std::cerr<< mFilename<<": failed to find stream info" <<std::endl;
        return false;
    }

    mVideo.schedRefresh(40);

    mParseThread = std::thread(std::mem_fn(&MovieState::parse_handler), this);
    return true;
}

void MovieState::setTitle(SDL_Window *window)
{
    auto pos1 = mFilename.rfind('/');
    auto pos2 = mFilename.rfind('\\');
    auto fpos = ((pos1 == std::string::npos) ? pos2 :
                 (pos2 == std::string::npos) ? pos1 :
                 std::max(pos1, pos2)) + 1;
    SDL_SetWindowTitle(window, (mFilename.substr(fpos)+" - "+AppName).c_str());
}

double MovieState::getClock()
{
    return (av_gettime()-mExternalClockBase) / 1000000.0;
}

double MovieState::getMasterClock()
{
    if(mAVSyncType == AV_SYNC_VIDEO_MASTER)
        return mVideo.getClock();
    if(mAVSyncType == AV_SYNC_AUDIO_MASTER)
        return mAudio.getClock();
    return getClock();
}

int MovieState::streamComponentOpen(int stream_index)
{
    if(stream_index < 0 || (unsigned int)stream_index >= mFormatCtx->nb_streams)
        return -1;

    /* Get a pointer to the codec context for the stream, and open the
     * associated codec.
     */
    AVCodecContext *avctx = avcodec_alloc_context3(nullptr);
    if(!avctx) return -1;

    if(avcodec_parameters_to_context(avctx, mFormatCtx->streams[stream_index]->codecpar))
    {
        avcodec_free_context(&avctx);
        return -1;
    }

    AVCodec *codec = avcodec_find_decoder(avctx->codec_id);
    if(!codec || avcodec_open2(avctx, codec, nullptr) < 0)
    {
        std::cerr<< "Unsupported codec: "<<avcodec_get_name(avctx->codec_id)
                 << " (0x"<<std::hex<<avctx->codec_id<<std::dec<<")" <<std::endl;
        avcodec_free_context(&avctx);
        return -1;
    }

    /* Initialize and start the media type handler */
    switch(avctx->codec_type)
    {
        case AVMEDIA_TYPE_AUDIO:
            mAudioStream = stream_index;
            mAudio.mStream = mFormatCtx->streams[stream_index];
            mAudio.mCodecCtx = avctx;

            /* Averaging filter for audio sync */
            mAudio.mDiff.AvgCoeff = exp(log(0.01) / AUDIO_DIFF_AVG_NB);
            /* Correct audio only if larger error than this */
            mAudio.mDiff.Threshold = 0.050/* 50 ms */;

            mAudioThread = std::thread(std::mem_fn(&AudioState::handler), &mAudio);
            break;

        case AVMEDIA_TYPE_VIDEO:
            mVideoStream = stream_index;
            mVideo.mStream = mFormatCtx->streams[stream_index];
            mVideo.mCodecCtx = avctx;

            mVideo.mCurrentPtsTime = av_gettime();
            mVideo.mFrameTimer = (double)mVideo.mCurrentPtsTime / 1000000.0;
            mVideo.mFrameLastDelay = 40e-3;

            mVideoThread = std::thread(std::mem_fn(&VideoState::handler), &mVideo);
            break;

        default:
            avcodec_free_context(&avctx);
            break;
    }

    return 0;
}

int MovieState::parse_handler()
{
    int video_index = -1;
    int audio_index = -1;

    mVideoStream = -1;
    mAudioStream = -1;

    /* Dump information about file onto standard error */
    av_dump_format(mFormatCtx, 0, mFilename.c_str(), 0);

    /* Find the first video and audio streams */
    for(unsigned int i = 0;i < mFormatCtx->nb_streams;i++)
    {
        if(mFormatCtx->streams[i]->codecpar->codec_type == AVMEDIA_TYPE_VIDEO && video_index < 0)
            video_index = i;
        else if(mFormatCtx->streams[i]->codecpar->codec_type == AVMEDIA_TYPE_AUDIO && audio_index < 0)
            audio_index = i;
    }
    /* Start the external clock in 50ms, to give the audio and video
     * components time to start without needing to skip ahead.
     */
    mExternalClockBase = av_gettime() + 50000;
    if(audio_index >= 0)
        streamComponentOpen(audio_index);
    if(video_index >= 0)
        streamComponentOpen(video_index);

    if(mVideoStream < 0 && mAudioStream < 0)
    {
        std::cerr<< mFilename<<": could not open codecs" <<std::endl;
        mQuit = true;
    }

    /* Main packet handling loop */
    while(!mQuit.load())
    {
        if(mAudio.mQueue.mTotalSize + mVideo.mQueue.mTotalSize >= MAX_QUEUE_SIZE)
        {
            std::this_thread::sleep_for(std::chrono::milliseconds(10));
            continue;
        }

        AVPacket packet;
        if(av_read_frame(mFormatCtx, &packet) < 0)
            break;

        /* Copy the packet in the queue it's meant for. */
        if(packet.stream_index == mVideoStream)
            mVideo.mQueue.put(&packet);
        else if(packet.stream_index == mAudioStream)
            mAudio.mQueue.put(&packet);
        av_packet_unref(&packet);
    }
    mVideo.mQueue.finish();
    mAudio.mQueue.finish();

    /* all done - wait for it */
    if(mVideoThread.joinable())
        mVideoThread.join();
    if(mAudioThread.joinable())
        mAudioThread.join();

    mVideo.mEOS = true;
    std::unique_lock<std::mutex> lock(mVideo.mPictQMutex);
    while(!mVideo.mFinalUpdate)
        mVideo.mPictQCond.wait(lock);
    lock.unlock();

    SDL_Event evt{};
    evt.user.type = FF_MOVIE_DONE_EVENT;
    SDL_PushEvent(&evt);

    return 0;
}

} // namespace


int main(int argc, char *argv[])
{
    std::unique_ptr<MovieState> movState;

    if(argc < 2)
    {
        std::cerr<< "Usage: "<<argv[0]<<" [-device <device name>] [-direct] <files...>" <<std::endl;
        return 1;
    }
    /* Register all formats and codecs */
    av_register_all();
    /* Initialize networking protocols */
    avformat_network_init();

    if(SDL_Init(SDL_INIT_VIDEO | SDL_INIT_TIMER))
    {
        std::cerr<< "Could not initialize SDL - <<"<<SDL_GetError() <<std::endl;
        return 1;
    }

    /* Make a window to put our video */
    SDL_Window *screen = SDL_CreateWindow(AppName.c_str(), 0, 0, 640, 480, SDL_WINDOW_RESIZABLE);
    if(!screen)
    {
        std::cerr<< "SDL: could not set video mode - exiting" <<std::endl;
        return 1;
    }
    /* Make a renderer to handle the texture image surface and rendering. */
    SDL_Renderer *renderer = SDL_CreateRenderer(screen, -1, SDL_RENDERER_ACCELERATED);
    if(renderer)
    {
        SDL_RendererInfo rinf{};
        bool ok = false;

        /* Make sure the renderer supports IYUV textures. If not, fallback to a
         * software renderer. */
        if(SDL_GetRendererInfo(renderer, &rinf) == 0)
        {
            for(Uint32 i = 0;!ok && i < rinf.num_texture_formats;i++)
                ok = (rinf.texture_formats[i] == SDL_PIXELFORMAT_IYUV);
        }
        if(!ok)
        {
            std::cerr<< "IYUV pixelformat textures not supported on renderer "<<rinf.name <<std::endl;
            SDL_DestroyRenderer(renderer);
            renderer = nullptr;
        }
    }
    if(!renderer)
        renderer = SDL_CreateRenderer(screen, -1, SDL_RENDERER_SOFTWARE);
    if(!renderer)
    {
        std::cerr<< "SDL: could not create renderer - exiting" <<std::endl;
        return 1;
    }
    SDL_SetRenderDrawColor(renderer, 0, 0, 0, 255);
    SDL_RenderFillRect(renderer, nullptr);
    SDL_RenderPresent(renderer);

    /* Open an audio device */
    int fileidx = 1;
    ALCdevice *device = [argc,argv,&fileidx]() -> ALCdevice*
    {
        ALCdevice *dev = NULL;
        if(argc > 3 && strcmp(argv[1], "-device") == 0)
        {
            fileidx = 3;
            dev = alcOpenDevice(argv[2]);
            if(dev) return dev;
            std::cerr<< "Failed to open \""<<argv[2]<<"\" - trying default" <<std::endl;
        }
        return alcOpenDevice(nullptr);
    }();
    ALCcontext *context = alcCreateContext(device, nullptr);
    if(!context || alcMakeContextCurrent(context) == ALC_FALSE)
    {
        std::cerr<< "Failed to set up audio device" <<std::endl;
        if(context)
            alcDestroyContext(context);
        return 1;
    }

    const ALCchar *name = nullptr;
    if(alcIsExtensionPresent(device, "ALC_ENUMERATE_ALL_EXT"))
        name = alcGetString(device, ALC_ALL_DEVICES_SPECIFIER);
    if(!name || alcGetError(device) != AL_NO_ERROR)
        name = alcGetString(device, ALC_DEVICE_SPECIFIER);
    std::cout<< "Opened \""<<name<<"\"" <<std::endl;

    if(fileidx < argc && strcmp(argv[fileidx], "-direct") == 0)
    {
        ++fileidx;
        do_direct_out = true;
    }

    while(fileidx < argc && !movState)
    {
        movState = std::unique_ptr<MovieState>(new MovieState(argv[fileidx++]));
        if(!movState->prepare()) movState = nullptr;
    }
    if(!movState)
    {
        std::cerr<< "Could not start a video" <<std::endl;
        return 1;
    }
    movState->setTitle(screen);

    /* Default to going to the next movie at the end of one. */
    enum class EomAction {
        Next, Quit
    } eom_action = EomAction::Next;
    SDL_Event event;
    while(SDL_WaitEvent(&event) == 1)
    {
        switch(event.type)
        {
            case SDL_KEYDOWN:
                switch(event.key.keysym.sym)
                {
                    case SDLK_ESCAPE:
                        movState->mQuit = true;
                        eom_action = EomAction::Quit;
                        break;

                    case SDLK_n:
                        movState->mQuit = true;
                        eom_action = EomAction::Next;
                        break;

                    default:
                        break;
                }
                break;

            case SDL_WINDOWEVENT:
                switch(event.window.event)
                {
                    case SDL_WINDOWEVENT_RESIZED:
                        SDL_SetRenderDrawColor(renderer, 0, 0, 0, 255);
                        SDL_RenderFillRect(renderer, nullptr);
                        break;

                    default:
                        break;
                }
                break;

            case SDL_QUIT:
                movState->mQuit = true;
                eom_action = EomAction::Quit;
                break;

            case FF_UPDATE_EVENT:
                reinterpret_cast<VideoState*>(event.user.data1)->updatePicture(
                    screen, renderer
                );
                break;

            case FF_REFRESH_EVENT:
                reinterpret_cast<VideoState*>(event.user.data1)->refreshTimer(
                    screen, renderer
                );
                break;

            case FF_MOVIE_DONE_EVENT:
                if(eom_action != EomAction::Quit)
                {
                    movState = nullptr;
                    while(fileidx < argc && !movState)
                    {
                        movState = std::unique_ptr<MovieState>(new MovieState(argv[fileidx++]));
                        if(!movState->prepare()) movState = nullptr;
                    }
                    if(movState)
                    {
                        movState->setTitle(screen);
                        break;
                    }
                }

                /* Nothing more to play. Shut everything down and quit. */
                movState = nullptr;

                alcMakeContextCurrent(nullptr);
                alcDestroyContext(context);
                alcCloseDevice(device);

                SDL_DestroyRenderer(renderer);
                renderer = nullptr;
                SDL_DestroyWindow(screen);
                screen = nullptr;

                SDL_Quit();
                exit(0);

            default:
                break;
        }
    }

    std::cerr<< "SDL_WaitEvent error - "<<SDL_GetError() <<std::endl;
    return 1;
}
