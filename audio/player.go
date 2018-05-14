// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package audio

// #include <stdlib.h>
import "C"

import (
	"github.com/g3n/engine/audio/al"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/math32"
	"io"
	"time"
	"unsafe"
)

const (
	playerBufferCount = 2
	playerBufferSize  = 32 * 1024
)

// Player is a 3D (spatial) audio file player
// It embeds a core.Node so it can be inserted as a child in any other 3D object.
type Player struct {
	core.Node                // Embedded node
	af        *AudioFile     // Pointer to media audio file
	buffers   []uint32       // OpenAL buffer names
	source    uint32         // OpenAL source name
	nextBuf   int            // Index of next buffer to fill
	pdata     unsafe.Pointer // Pointer to C allocated storage
	disposed  bool           // Disposed flag
	gchan     chan (string)  // Channel for informing of goroutine end
}

// NewPlayer creates and returns a pointer to a new audio player object
// which will play the audio encoded in the specified file.
// Currently it supports wave and Ogg Vorbis formats.
func NewPlayer(filename string) (*Player, error) {

	// Try to open audio file
	af, err := NewAudioFile(filename)
	if err != nil {
		return nil, err
	}

	// Creates player
	p := new(Player)
	p.Node.Init()
	p.af = af

	// Generate buffers names
	p.buffers = al.GenBuffers(playerBufferCount)

	// Generate source name
	p.source = al.GenSource()

	// Allocates C memory buffer
	p.pdata = C.malloc(playerBufferSize)

	// Initialize channel for communication with internal goroutine
	p.gchan = make(chan string, 1)
	return p, nil
}

// Dispose disposes of this player resources
func (p *Player) Dispose() {

	p.Stop()

	// Close file
	p.af.Close()

	// Release OpenAL resources
	al.DeleteSource(p.source)
	al.DeleteBuffers(p.buffers)

	// Release C memory
	C.free(p.pdata)
	p.pdata = nil
	p.disposed = true
}

// State returns the current state of this player
func (p *Player) State() int {

	return int(al.GetSourcei(p.source, al.SourceState))
}

// Play starts playing this player
func (p *Player) Play() error {

	state := p.State()
	// Already playing, nothing to do
	if state == al.Playing {
		return nil
	}

	// If paused, goroutine should be running, just starts playing
	if state == al.Paused {
		al.SourcePlay(p.source)
		return nil
	}

	// Inactive or Stopped state
	if state == al.Initial || state == al.Stopped {

		// Sets file pointer to the beginning
		err := p.af.Seek(0)
		if err != nil {
			return err
		}

		// Fill buffers with decoded data
		for i := 0; i < playerBufferCount; i++ {
			err = p.fillBuffer(p.buffers[i])
			if err != nil {
				if err != io.EOF {
					return err
				}
				break
			}
		}
		p.nextBuf = 0

		// Clear previous goroutine response channel
		select {
		case _ = <-p.gchan:
		default:
		}
		// Starts playing and starts goroutine to fill buffers
		al.SourcePlay(p.source)
		go p.run()
		return nil
	}
	return nil
}

// Pause sets the player in the pause state
func (p *Player) Pause() {

	if p.State() == al.Paused {
		return
	}
	al.SourcePause(p.source)
}

// Stop stops the player
func (p *Player) Stop() {

	state := p.State()
	if state == al.Stopped || state == al.Initial {
		return
	}
	al.SourceStop(p.source)
	// Waits for goroutine to finish
	<-p.gchan
}

// CurrentTime returns the current time in seconds spent in the stream
func (p *Player) CurrentTime() float64 {

	return p.af.CurrentTime()
}

// TotalTime returns the total time in seconds to play this stream
func (p *Player) TotalTime() float64 {

	return p.af.info.TotalTime
}

// Gain returns the current gain (volume) of this player
func (p *Player) Gain() float32 {

	return al.GetSourcef(p.source, al.Gain)
}

// SetGain sets the gain (volume) of this player
func (p *Player) SetGain(gain float32) {

	al.Sourcef(p.source, al.Gain, gain)
}

// MinGain returns the current minimum gain of this player
func (p *Player) MinGain() float32 {

	return al.GetSourcef(p.source, al.MinGain)
}

// SetMinGain sets the minimum gain (volume) of this player
func (p *Player) SetMinGain(gain float32) {

	al.Sourcef(p.source, al.MinGain, gain)
}

// MaxGain returns the current maximum gain of this player
func (p *Player) MaxGain() float32 {

	return al.GetSourcef(p.source, al.MaxGain)
}

// SetMaxGain sets the maximum gain (volume) of this player
func (p *Player) SetMaxGain(gain float32) {

	al.Sourcef(p.source, al.MaxGain, gain)
}

// Pitch returns the current pitch factor of this player
func (p *Player) Pitch() float32 {

	return al.GetSourcef(p.source, al.Pitch)
}

// SetPitch sets the pitch factor of this player
func (p *Player) SetPitch(pitch float32) {

	al.Sourcef(p.source, al.Pitch, pitch)
}

// Looping returns the current looping state of this player
func (p *Player) Looping() bool {

	return p.af.Looping()
}

// SetLooping sets the looping state of this player
func (p *Player) SetLooping(looping bool) {

	p.af.SetLooping(looping)
}

// InnerCone returns the inner cone angle in degrees
func (p *Player) InnerCone() float32 {

	return al.GetSourcef(p.source, al.ConeInnerAngle)
}

// SetInnerCone sets the inner cone angle in degrees
func (p *Player) SetInnerCone(inner float32) {

	al.Sourcef(p.source, al.ConeInnerAngle, inner)
}

// OuterCone returns the outer cone angle in degrees
func (p *Player) OuterCone() float32 {

	return al.GetSourcef(p.source, al.ConeOuterAngle)
}

// SetOuterCone sets the outer cone angle in degrees
func (p *Player) SetOuterCone(outer float32) {

	al.Sourcef(p.source, al.ConeOuterAngle, outer)
}

// SetVelocity sets the velocity of this player
// It is used to calculate Doppler effects
func (p *Player) SetVelocity(vx, vy, vz float32) {

	al.Source3f(p.source, al.Velocity, vx, vy, vz)
}

// SetVelocityVec sets the velocity of this player from the specified vector
// It is used to calculate Doppler effects
func (p Player) SetVelocityVec(v *math32.Vector3) {

	al.Source3f(p.source, al.Velocity, v.X, v.Y, v.Z)
}

// Velocity returns this player velocity
func (p *Player) Velocity() (float32, float32, float32) {

	return al.GetSource3f(p.source, al.Velocity)
}

// VelocityVec returns this player velocity vector
func (p *Player) VelocityVec() math32.Vector3 {

	vx, vy, vz := al.GetSource3f(p.source, al.Velocity)
	return math32.Vector3{vx, vy, vz}
}

// SetRolloffFactor sets this player rolloff factor user to calculate
// the gain attenuation by distance
func (p *Player) SetRolloffFactor(rfactor float32) {

	al.Sourcef(p.source, al.RolloffFactor, rfactor)
}

// Render satisfies the INode interface.
// It is called by renderer at every frame and is used to
// update the audio source position and direction
func (p *Player) Render(gl *gls.GLS) {

	// Sets the player source world position
	var wpos math32.Vector3
	p.WorldPosition(&wpos)
	al.Source3f(p.source, al.Position, wpos.X, wpos.Y, wpos.Z)

	// Sets the player source world direction
	var wdir math32.Vector3
	p.WorldDirection(&wdir)
	al.Source3f(p.source, al.Direction, wdir.X, wdir.Y, wdir.Z)
}

// Goroutine to fill PCM buffers with decoded data for OpenAL
func (p *Player) run() {

	for {
		// Get current state of player source
		state := al.GetSourcei(p.source, al.SourceState)
		processed := al.GetSourcei(p.source, al.BuffersProcessed)
		queued := al.GetSourcei(p.source, al.BuffersQueued)
		//log.Debug("state:%x processed:%v queued:%v", state, processed, queued)

		// If stopped, unqueues all buffer before exiting
		if state == al.Stopped {
			if queued == 0 {
				break
			}
			// Unqueue buffers
			if processed > 0 {
				al.SourceUnqueueBuffers(p.source, uint32(processed), nil)
			}
			continue
		}

		// If no buffers processed, sleeps and try again
		if processed == 0 {
			time.Sleep(20 * time.Millisecond)
			continue
		}

		// Remove processed buffers from the queue
		al.SourceUnqueueBuffers(p.source, uint32(processed), nil)
		// Fill and enqueue buffers with new data
		for i := 0; i < int(processed); i++ {
			err := p.fillBuffer(p.buffers[p.nextBuf])
			if err != nil {
				break
			}
			p.nextBuf = (p.nextBuf + 1) % playerBufferCount
		}
	}
	// Sends indication of goroutine end
	p.gchan <- "end"
}

// fillBuffer fills the specified OpenAL buffer with next decoded data
// and queues the buffer to this player source
func (p *Player) fillBuffer(buf uint32) error {

	// Reads next decoded data
	n, err := p.af.Read(p.pdata, playerBufferSize)
	if err != nil {
		return err
	}
	// Sends data to buffer
	//log.Debug("BufferData:%v format:%x n:%v rate:%v", buf, p.af.info.Format, n, p.af.info.SampleRate)
	al.BufferData(buf, uint32(p.af.info.Format), p.pdata, uint32(n), uint32(p.af.info.SampleRate))
	al.SourceQueueBuffers(p.source, buf)
	return nil
}
