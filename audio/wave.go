// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package audio

import (
	"fmt"
	"github.com/g3n/engine/audio/al"
	"os"
)

// WaveSpecs describes the characteristics of the audio encoded in a wave file.
type WaveSpecs struct {
	Format     int     // OpenAl Format
	Type       int     // Type field from wave header
	Channels   int     // Number of channels
	SampleRate int     // Sample rate in hz
	BitsSample int     // Number of bits per sample (8 or 16)
	DataSize   int     // Total data size in bytes
	BytesSec   int     // Bytes per second
	TotalTime  float64 // Total time in seconds
}

const (
	waveHeaderSize = 44
	fileMark       = "RIFF"
	fileHead       = "WAVE"
)

// WaveCheck checks if the specified filepath corresponds to a an audio wave file.
// If the file is a valid wave file, return a pointer to WaveSpec structure
// with information about the encoded audio data.
func WaveCheck(filepath string) (*WaveSpecs, error) {

	// Open file
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Reads header
	header := make([]uint8, waveHeaderSize)
	n, err := f.Read(header)
	if err != nil {
		return nil, err
	}
	if n < waveHeaderSize {
		return nil, fmt.Errorf("File size less than header")
	}
	// Checks file marks
	if string(header[0:4]) != fileMark {
		return nil, fmt.Errorf("'RIFF' mark not found")
	}
	if string(header[8:12]) != fileHead {
		return nil, fmt.Errorf("'WAVE' mark not found")
	}

	// Decodes header fields
	var ws WaveSpecs
	ws.Format = -1
	ws.Type = int(header[20]) + int(header[21])<<8
	ws.Channels = int(header[22]) + int(header[23])<<8
	ws.SampleRate = int(header[24]) + int(header[25])<<8 + int(header[26])<<16 + int(header[27])<<24
	ws.BitsSample = int(header[34]) + int(header[35])<<8
	ws.DataSize = int(header[40]) + int(header[41])<<8 + int(header[42])<<16 + int(header[43])<<24

	// Sets OpenAL format field if possible
	if ws.Channels == 1 {
		if ws.BitsSample == 8 {
			ws.Format = al.FormatMono8
		} else if ws.BitsSample == 16 {
			ws.Format = al.FormatMono16
		}
	} else if ws.Channels == 2 {
		if ws.BitsSample == 8 {
			ws.Format = al.FormatStereo8
		} else if ws.BitsSample == 16 {
			ws.Format = al.FormatStereo16
		}
	}

	// Calculates bytes/sec and total time
	var bytesChannel int
	if ws.BitsSample == 8 {
		bytesChannel = 1
	} else {
		bytesChannel = 2
	}
	ws.BytesSec = ws.SampleRate * ws.Channels * bytesChannel
	ws.TotalTime = float64(ws.DataSize) / float64(ws.BytesSec)
	return &ws, nil
}
