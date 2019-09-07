package util

import (
	"time"
)

// FrameRater implements a frame rate controller
type FrameRater struct {
	targetFPS      uint          // desired number of frames per second
	targetDuration time.Duration // calculated desired duration of frame
	frameStart     time.Time     // start time of last frame
	frameTimes     time.Duration // accumulated frame times for potential FPS calculation
	frameCount     uint          // accumulated number of frames for FPS calculation
	lastUpdate     time.Time     // time of last FPS calculation update
	timer          *time.Timer   // timer for sleeping during frame
}

// NewFrameRater returns a frame rate controller object for the specified
// number of target frames per second
func NewFrameRater(targetFPS uint) *FrameRater {

	f := new(FrameRater)
	f.targetDuration = time.Second / time.Duration(targetFPS)
	f.frameTimes = 0
	f.frameCount = 0
	f.lastUpdate = time.Now()
	f.timer = time.NewTimer(0)
	<-f.timer.C
	return f
}

// Start should be called at the start of the frame
func (f *FrameRater) Start() {

	f.frameStart = time.Now()
}

// Wait should be called at the end of the frame
// If necessary it will sleep to achieve the desired frame rate
func (f *FrameRater) Wait() {

	// Calculates the time duration of this frame
	elapsed := time.Now().Sub(f.frameStart)
	// Accumulates this frame time for potential FPS calculation
	f.frameCount++
	f.frameTimes += elapsed
	// If this frame duration is less than the target duration, sleeps
	// during the difference
	diff := f.targetDuration - elapsed
	if diff > 0 {
		f.timer.Reset(diff)
		<-f.timer.C
	}
}

// FPS calculates and returns the current measured FPS and the maximum
// potential FPS after the specified time interval has elapsed.
// It returns an indication if the results are valid
func (f *FrameRater) FPS(t time.Duration) (float64, float64, bool) {

	// If the time from the last update has not passed, nothing to do
	elapsed := time.Now().Sub(f.lastUpdate)
	if elapsed < t {
		return 0, 0, false
	}

	// Calculates the measured average frame rate
	fps := float64(f.frameCount) / elapsed.Seconds()
	// Calculates the average duration of a frame and the potential FPS
	frameDur := f.frameTimes.Seconds() / float64(f.frameCount)
	pfps := 1.0 / frameDur
	// Resets the frame counter and times
	f.frameCount = 0
	f.frameTimes = 0
	f.lastUpdate = time.Now()
	return fps, pfps, true
}
