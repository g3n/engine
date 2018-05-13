// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package audio

import (
	"fmt"
	"github.com/g3n/engine/audio/al"
	"github.com/g3n/engine/audio/ov"
	"io"
	"os"
	"unsafe"
)

// AudioInfo represents the information associated to an audio file
type AudioInfo struct {
	Format     int     // OpenAl Format
	Channels   int     // Number of channels
	SampleRate int     // Sample rate in hz
	BitsSample int     // Number of bits per sample (8 or 16)
	DataSize   int     // Total data size in bytes
	BytesSec   int     // Bytes per second
	TotalTime  float64 // Total time in seconds
}

// AudioFile represents an audio file
type AudioFile struct {
	wavef   *os.File  // Pointer to wave file opened filed (nil for vorbis)
	vorbisf *ov.File  // Pointer to vorbis file structure (nil for wave)
	info    AudioInfo // Audio information structure
	looping bool      // Looping flag
}

// NewAudioFile creates and returns a pointer to a new audio file object and an error
func NewAudioFile(filename string) (*AudioFile, error) {

	// Checks if file exists
	_, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	af := new(AudioFile)

	// Try to open as a wave file
	if af.openWave(filename) == nil {
		return af, nil
	}

	// Try to open as an ogg vorbis file
	if af.openVorbis(filename) == nil {
		return af, nil
	}

	return nil, fmt.Errorf("Unsuported file type")
}

// Close closes the audiofile
func (af *AudioFile) Close() error {

	if af.wavef != nil {
		return af.wavef.Close()
	}
	return ov.Clear(af.vorbisf)
}

// Read reads decoded data from the audio file
func (af *AudioFile) Read(pdata unsafe.Pointer, nbytes int) (int, error) {

	// Slice to access buffer
	bs := (*[1 << 30]byte)(pdata)[0:nbytes:nbytes]

	// Reads wave file directly
	if af.wavef != nil {
		n, err := af.wavef.Read(bs)
		if err != nil {
			return 0, err
		}
		if !af.looping {
			return n, nil
		}
		if n == nbytes {
			return n, nil
		}
		// EOF reached. Position file at the beginning
		_, err = af.wavef.Seek(int64(waveHeaderSize), 0)
		if err != nil {
			return 0, nil
		}
		// Reads next data into the remaining buffer space
		n2, err := af.wavef.Read(bs[n:])
		if err != nil {
			return 0, err
		}
		return n + n2, err
	}

	// Decodes Ogg vorbis
	decoded := 0
	for decoded < nbytes {
		n, _, err := ov.Read(af.vorbisf, unsafe.Pointer(&bs[decoded]), nbytes-decoded, false, 2, true)
		// Error
		if err != nil {
			return 0, err
		}
		// EOF
		if n == 0 {
			if !af.looping {
				break
			}
			// Position file at the beginning
			err = ov.PcmSeek(af.vorbisf, 0)
			if err != nil {
				return 0, err
			}
		}
		decoded += n
	}
	if nbytes > 0 && decoded == 0 {
		return 0, io.EOF
	}
	return decoded, nil
}

// Seek sets the file reading position relative to the origin
func (af *AudioFile) Seek(pos uint) error {

	if af.wavef != nil {
		_, err := af.wavef.Seek(int64(waveHeaderSize+pos), 0)
		return err
	}
	return ov.PcmSeek(af.vorbisf, int64(pos))
}

// Info returns the audio info structure for this audio file
func (af *AudioFile) Info() AudioInfo {

	return af.info
}

// CurrentTime returns the current time in seconds for the current file read position
func (af *AudioFile) CurrentTime() float64 {

	if af.vorbisf != nil {
		pos, _ := ov.TimeTell(af.vorbisf)
		return pos
	}
	pos, err := af.wavef.Seek(0, 1)
	if err != nil {
		return 0
	}
	return float64(pos) / float64(af.info.BytesSec)
}

// Looping returns the current looping state of this audio file
func (af *AudioFile) Looping() bool {

	return af.looping
}

// SetLooping sets the looping state of this audio file
func (af *AudioFile) SetLooping(looping bool) {

	af.looping = looping
}

// openWave tries to open the specified file as a wave file
// and if succesfull, sets the file pointer positioned after the header.
func (af *AudioFile) openWave(filename string) error {

	// Open file
	osf, err := os.Open(filename)
	if err != nil {
		return err
	}

	// Reads header
	header := make([]uint8, waveHeaderSize)
	n, err := osf.Read(header)
	if err != nil {
		osf.Close()
		return err
	}
	if n < waveHeaderSize {
		osf.Close()
		return fmt.Errorf("File size less than header")
	}
	// Checks file marks
	if string(header[0:4]) != fileMark {
		osf.Close()
		return fmt.Errorf("'RIFF' mark not found")
	}
	if string(header[8:12]) != fileHead {
		osf.Close()
		return fmt.Errorf("'WAVE' mark not found")
	}

	// Decodes header fields
	af.info.Format = -1
	af.info.Channels = int(header[22]) + int(header[23])<<8
	af.info.SampleRate = int(header[24]) + int(header[25])<<8 + int(header[26])<<16 + int(header[27])<<24
	af.info.BitsSample = int(header[34]) + int(header[35])<<8
	af.info.DataSize = int(header[40]) + int(header[41])<<8 + int(header[42])<<16 + int(header[43])<<24

	// Sets OpenAL format field if possible
	if af.info.Channels == 1 {
		if af.info.BitsSample == 8 {
			af.info.Format = al.FormatMono8
		} else if af.info.BitsSample == 16 {
			af.info.Format = al.FormatMono16
		}
	} else if af.info.Channels == 2 {
		if af.info.BitsSample == 8 {
			af.info.Format = al.FormatStereo8
		} else if af.info.BitsSample == 16 {
			af.info.Format = al.FormatStereo16
		}
	}
	if af.info.Format == -1 {
		osf.Close()
		return fmt.Errorf("Unsupported OpenAL format")
	}

	// Calculates bytes/sec and total time
	var bytesChannel int
	if af.info.BitsSample == 8 {
		bytesChannel = 1
	} else {
		bytesChannel = 2
	}
	af.info.BytesSec = af.info.SampleRate * af.info.Channels * bytesChannel
	af.info.TotalTime = float64(af.info.DataSize) / float64(af.info.BytesSec)

	// Seeks after the header
	_, err = osf.Seek(waveHeaderSize, 0)
	if err != nil {
		osf.Close()
		return err
	}

	af.wavef = osf
	return nil
}

// openVorbis tries to open the specified file as an ogg vorbis file
// and if succesfull, sets up the player for playing this file
func (af *AudioFile) openVorbis(filename string) error {

	// Try to open file as ogg vorbis
	vf, err := ov.Fopen(filename)
	if err != nil {
		return err
	}

	// Get info for opened vorbis file
	var info ov.VorbisInfo
	err = ov.Info(vf, -1, &info)
	if err != nil {
		return err
	}
	if info.Channels == 1 {
		af.info.Format = al.FormatMono16
	} else if info.Channels == 2 {
		af.info.Format = al.FormatStereo16
	} else {
		return fmt.Errorf("Unsupported number of channels")
	}
	totalSamples, err := ov.PcmTotal(vf, -1)
	if err != nil {
		ov.Clear(vf)
		return nil
	}
	timeTotal, err := ov.TimeTotal(vf, -1)
	if err != nil {
		ov.Clear(vf)
		return nil
	}

	af.vorbisf = vf
	af.info.SampleRate = info.Rate
	af.info.BitsSample = 16
	af.info.Channels = info.Channels
	af.info.DataSize = int(totalSamples) * info.Channels * 2
	af.info.TotalTime = timeTotal
	return nil
}
