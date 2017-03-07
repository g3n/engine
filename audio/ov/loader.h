#ifndef VBF_LOADER_H
#define VBF_LOADER_H

#include "vorbis/vorbisfile.h"

#if defined(_WIN32)
 #define VBF_APIENTRY __cdecl
#else
 #define VBF_APIENTRY
#endif


// API function pointers type definitions
typedef int (VBF_APIENTRY *LPOVCLEAR)(OggVorbis_File *vf);
typedef int (VBF_APIENTRY *LPOVFOPEN)(const char *path,OggVorbis_File *vf);
typedef int (VBF_APIENTRY *LPOVOPEN)(FILE *f,OggVorbis_File *vf,const char *initial,long ibytes);
typedef int (VBF_APIENTRY *LPOVOPENCALLBACKS)(void *datasource, OggVorbis_File *vf, const char *initial, long ibytes, ov_callbacks callbacks);

typedef int (VBF_APIENTRY *LPOVTEST)(FILE *f,OggVorbis_File *vf,const char *initial,long ibytes);
typedef int (VBF_APIENTRY *LPOVTESTCALLBACKS)(void *datasource, OggVorbis_File *vf, const char *initial, long ibytes, ov_callbacks callbacks);
typedef int (VBF_APIENTRY *LPOVTESTOPEN)(OggVorbis_File *vf);

typedef long (VBF_APIENTRY *LPOVBITRATE)(OggVorbis_File *vf,int i);
typedef long (VBF_APIENTRY *LPOVBITRATEINSTANT)(OggVorbis_File *vf);
typedef int (VBF_APIENTRY *LPOVSTREAMS)(OggVorbis_File *vf);
typedef int (VBF_APIENTRY *LPOVSEEKABLE)(OggVorbis_File *vf);
typedef int (VBF_APIENTRY *LPOVSERIALNUMBER)(OggVorbis_File *vf,int i);

typedef ogg_int64_t (VBF_APIENTRY *LPOVRAWTOTAL)(OggVorbis_File *vf,int i);
typedef ogg_int64_t (VBF_APIENTRY *LPOVPCMTOTAL)(OggVorbis_File *vf,int i);
typedef double (VBF_APIENTRY *LPOVTIMETOTAL)(OggVorbis_File *vf,int i);

typedef int (VBF_APIENTRY *LPOVRAWSEEK)(OggVorbis_File *vf,ogg_int64_t pos);
typedef int (VBF_APIENTRY *LPOVPCMSEEK)(OggVorbis_File *vf,ogg_int64_t pos);
typedef int (VBF_APIENTRY *LPOVPCMSEEKPAGE)(OggVorbis_File *vf,ogg_int64_t pos);
typedef int (VBF_APIENTRY *LPOVTIMESEEK)(OggVorbis_File *vf,double pos);
typedef int (VBF_APIENTRY *LPOVTIMESEEKPAGE)(OggVorbis_File *vf,double pos);

typedef int (VBF_APIENTRY *LPOVRAWSEEKLAP)(OggVorbis_File *vf,ogg_int64_t pos);
typedef int (VBF_APIENTRY *LPOVPCMSEEKLAP)(OggVorbis_File *vf,ogg_int64_t pos);
typedef int (VBF_APIENTRY *LPOVPCMSEEKPAGELAP)(OggVorbis_File *vf,ogg_int64_t pos);
typedef int (VBF_APIENTRY *LPOVTIMESEEKLAP)(OggVorbis_File *vf,double pos);
typedef int (VBF_APIENTRY *LPOVTIMESEEKPAGELAP)(OggVorbis_File *vf,double pos);

typedef ogg_int64_t (VBF_APIENTRY *LPOVRAWTELL)(OggVorbis_File *vf);
typedef ogg_int64_t (VBF_APIENTRY *LPOVPCMTELL)(OggVorbis_File *vf);
typedef double (VBF_APIENTRY *LPOVTIMETELL)(OggVorbis_File *vf);

typedef vorbis_info* (VBF_APIENTRY *LPOVINFO)(OggVorbis_File *vf,int link);
typedef vorbis_comment* (VBF_APIENTRY *LPOVCOMMENT)(OggVorbis_File *vf,int link);

typedef long (VBF_APIENTRY *LPOVREADFLOAT)(OggVorbis_File *vf,float ***pcm_channels,int samples, int *bitstream);
typedef long (VBF_APIENTRY *LPOVREADFILTER)(OggVorbis_File *vf,char *buffer,int length, int bigendianp,int word,int sgned,int *bitstream, void (*filter)(float **pcm,long channels,long samples,void *filter_param),void *filter_param); 
typedef long (VBF_APIENTRY *LPOVREAD)(OggVorbis_File *vf,char *buffer,int length, int bigendianp,int word,int sgned,int *bitstream);
typedef int (VBF_APIENTRY *LPOVCROSSLAP)(OggVorbis_File *vf1,OggVorbis_File *vf2);
typedef int (VBF_APIENTRY *LPOVHALFRATE)(OggVorbis_File *vf,int flag);
typedef int (VBF_APIENTRY *LPOVHALFRATEP)(OggVorbis_File *vf);


int vorbisfile_load();


extern LPOVCLEAR           p_ov_clear;
extern LPOVFOPEN           p_ov_fopen;
extern LPOVOPEN            p_ov_open;
extern LPOVOPENCALLBACKS   p_ov_open_callbacks;
extern LPOVTEST            p_ov_test;
extern LPOVTESTCALLBACKS   p_ov_test_callbacks;
extern LPOVTESTOPEN        p_ov_test_open;
extern LPOVBITRATE         p_ov_bitrate;
extern LPOVBITRATEINSTANT  p_ov_bitrate_instant;
extern LPOVSTREAMS         p_ov_streams;
extern LPOVSEEKABLE        p_ov_seekable;
extern LPOVSERIALNUMBER    p_ov_serialnumber;
extern LPOVRAWTOTAL        p_ov_raw_total;
extern LPOVPCMTOTAL        p_ov_pcm_total;
extern LPOVTIMETOTAL       p_ov_time_total;
extern LPOVRAWSEEK         p_ov_raw_seek;
extern LPOVPCMSEEK         p_ov_pcm_seek;
extern LPOVPCMSEEKPAGE     p_ov_pcm_seek_page;
extern LPOVTIMESEEK        p_ov_time_seek;
extern LPOVTIMESEEKPAGE    p_ov_time_seek_page;
extern LPOVRAWSEEKLAP      p_ov_raw_seek_lap;
extern LPOVPCMSEEKLAP      p_ov_pcm_seek_lap;
extern LPOVPCMSEEKPAGELAP  p_ov_pcm_seek_page_lap;
extern LPOVTIMESEEKLAP     p_ov_time_seek_lap;
extern LPOVTIMESEEKPAGELAP p_ov_time_seek_page_lap;
extern LPOVRAWTELL         p_ov_raw_tell;
extern LPOVPCMTELL         p_ov_pcm_tell;
extern LPOVTIMETELL        p_ov_time_tell;
extern LPOVINFO            p_ov_info;
extern LPOVCOMMENT         p_ov_comment;
extern LPOVREADFLOAT       p_ov_read_float;
extern LPOVREADFILTER      p_ov_read_filter;
extern LPOVREAD            p_ov_read;
extern LPOVCROSSLAP        p_ov_crosslap;
extern LPOVHALFRATE        p_ov_halfrate;
extern LPOVHALFRATEP       p_ov_halfrate_p;



#endif

