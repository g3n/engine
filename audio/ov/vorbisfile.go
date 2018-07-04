// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ov implements the Go bindings of a subset of the functions of the Ogg Vorbis File C library.
// The libvorbisfile C API reference is at: https://xiph.org/vorbis/doc/vorbisfile/reference.html
package ov

// #cgo darwin   CFLAGS:  -DGO_DARWIN  -I/usr/include/vorbis -I/usr/local/include/vorbis
// #cgo freebsd  CFLAGS:  -DGO_FREEBSD -I/usr/include/vorbis -I/usr/local/include/vorbis
// #cgo linux    CFLAGS:  -DGO_LINUX   -I/usr/include/vorbis
// #cgo windows  CFLAGS:  -DGO_WINDOWS -I${SRCDIR}/../windows/libvorbis-1.3.5/include/vorbis -I${SRCDIR}/../windows/libogg-1.3.3/include
// #cgo darwin   LDFLAGS: -L/usr/lib -L/usr/local/lib -lvorbisfile
// #cgo freebsd  LDFLAGS: -L/usr/lib -L/usr/local/lib -lvorbisfile
// #cgo linux    LDFLAGS: -lvorbisfile
// #cgo windows  LDFLAGS: -L${SRCDIR}/../windows/bin -llibvorbisfile
// #include <stdlib.h>
// #include "vorbisfile.h"
import "C"

import (
	"fmt"
	"unsafe"
)

// File type encapsulates a pointer to C allocated OggVorbis_File structure
type File struct {
	vf *C.OggVorbis_File
}

type VorbisInfo struct {
	Version        int
	Channels       int
	Rate           int
	BitrateUpper   int
	BitrateNominal int
	BitrateLower   int
	BitrateWindow  int
}

const (
	Eread      = C.OV_EREAD
	Efault     = C.OV_EFAULT
	Eimpl      = C.OV_EIMPL
	Einval     = C.OV_EINVAL
	EnotVorbis = C.OV_ENOTVORBIS
	EbadHeader = C.OV_EBADHEADER
	Eversion   = C.OV_EVERSION
	EnotAudio  = C.OV_ENOTAUDIO
	EbadPacket = C.OV_EBADPACKET
	EbadLink   = C.OV_EBADLINK
	EnoSeek    = C.OV_ENOSEEK
)

// Maps ogg vorbis error codes to string
var errCodes = map[C.int]string{
	C.OV_EREAD:      "Eread",
	C.OV_EFAULT:     "Efault",
	C.OV_EIMPL:      "Eimpl",
	C.OV_EINVAL:     "Einval",
	C.OV_ENOTVORBIS: "EnotVorbis",
	C.OV_EVERSION:   "Eversion",
	C.OV_ENOTAUDIO:  "EnotAudio",
	C.OV_EBADPACKET: "EbadPacket",
	C.OV_EBADLINK:   "EbadLink",
	C.OV_ENOSEEK:    "EnoSeek",
}

// Fopen opens an ogg vorbis file for decoding
// Returns an opaque pointer to the internal decode structure and an error
func Fopen(path string) (*File, error) {

	// Allocates pointer to vorbisfile structure using C memory
	var f File
	f.vf = (*C.OggVorbis_File)(C.malloc(C.size_t(unsafe.Sizeof(C.OggVorbis_File{}))))

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	cerr := C.ov_fopen(cpath, f.vf)
	if cerr == 0 {
		return &f, nil
	}
	return nil, fmt.Errorf("Error:%s from Fopen", errCodes[cerr])
}

// Clear clears the decoded buffers and closes the file
func Clear(f *File) error {

	cerr := C.ov_clear(f.vf)
	if cerr == 0 {
		C.free(unsafe.Pointer(f.vf))
		f.vf = nil
		return nil
	}
	return fmt.Errorf("Error:%s from Clear", errCodes[cerr])
}

// Read decodes next data from the file updating the specified buffer contents and
// returns the number of bytes read, the number of current logical bitstream and an error
func Read(f *File, buffer unsafe.Pointer, length int, bigendianp bool, word int, sgned bool) (int, int, error) {

	var cbigendianp C.int = 0
	var csgned C.int = 0
	var bitstream C.int

	if bigendianp {
		cbigendianp = 1
	}
	if sgned {
		csgned = 1
	}
	cres := C.ov_read(f.vf, (*C.char)(buffer), C.int(length), cbigendianp, C.int(word), csgned, &bitstream)
	if cres < 0 {
		return 0, 0, fmt.Errorf("Error:%s from Read()", errCodes[C.int(cres)])
	}
	return int(cres), int(bitstream), nil
}

// Info updates the specified VorbisInfo structure with contains basic
// information about the audio in a vorbis stream
func Info(f *File, link int, info *VorbisInfo) error {

	vi := C.ov_info(f.vf, C.int(link))
	if vi == nil {
		return fmt.Errorf("Error returned from 'ov_info'")
	}
	info.Version = int(vi.version)
	info.Channels = int(vi.channels)
	info.Rate = int(vi.rate)
	info.BitrateUpper = int(vi.bitrate_upper)
	info.BitrateNominal = int(vi.bitrate_nominal)
	info.BitrateLower = int(vi.bitrate_lower)
	info.BitrateWindow = int(vi.bitrate_window)
	return nil
}

// Seekable returns indication whether or not the bitstream is seekable
func Seekable(f *File) bool {

	cres := C.ov_seekable(f.vf)
	if cres == 0 {
		return false
	}
	return true
}

// Seek seeks to the offset specified (in number pcm samples) within the physical bitstream.
// This function only works for seekable streams.
// Updates everything needed within the decoder, so you can immediately call Read()
// and get data from the newly seeked to position.
func PcmSeek(f *File, pos int64) error {

	cres := C.ov_pcm_seek(f.vf, C.ogg_int64_t(pos))
	if cres == 0 {
		return nil
	}
	return fmt.Errorf("Error:%s from 'ov_pcm_seek()'", errCodes[C.int(cres)])
}

// PcmTotal returns the total number of pcm samples of the physical bitstream or a specified logical bit stream.
// To retrieve the total pcm samples for the entire physical bitstream, the 'link' parameter should be set to -1
func PcmTotal(f *File, i int) (int64, error) {

	cres := C.ov_pcm_total(f.vf, C.int(i))
	if cres < 0 {
		return 0, fmt.Errorf("Error:%s from 'ov_pcm_total()'", errCodes[C.int(cres)])
	}
	return int64(cres), nil
}

// TimeTotal returns the total time in seconds of the physical bitstream or a specified logical bitstream
// To retrieve the time total for the entire physical bitstream, 'i' should be set to -1.
func TimeTotal(f *File, i int) (float64, error) {

	cres := C.ov_time_total(f.vf, C.int(i))
	if cres < 0 {
		return 0, fmt.Errorf("Error:%s from 'ov_time_total()'", errCodes[C.int(cres)])
	}
	return float64(cres), nil
}

// TimeTell returns the current decoding offset in seconds.
func TimeTell(f *File) (float64, error) {

	cres := C.ov_time_tell(f.vf)
	if cres < 0 {
		return 0, fmt.Errorf("Error:%s from 'ov_time_total()'", errCodes[C.int(cres)])
	}
	return float64(cres), nil
}
