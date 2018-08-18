// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package animation
package animation

// Animation is a keyframe animation, containing channels.
// Each channel animates a specific property of an object.
// Animations can span multiple objects and properties.
type Animation struct {
	name     string     // Animation name
	loop     bool       // Whether the animation loops
	paused   bool       // Whether the animation is paused
	start    float32    // Initial time offset value
	time     float32    // Total running time
	minTime  float32    // Minimum time value across all channels
	maxTime  float32    // Maximum time value across all channels
	speed    float32    // Animation speed multiplier
	channels []IChannel // List of channels
}

// NewAnimation creates and returns a pointer to a new Animation object.
func NewAnimation() *Animation {

	anim := new(Animation)
	anim.speed = 1
	return anim
}

// SetName sets the animation name.
func (anim *Animation) SetName(name string) {

	anim.name = name
}

// Name returns the animation name.
func (anim *Animation) Name() string {

	return anim.name
}

// SetSpeed sets the animation speed.
func (anim *Animation) SetSpeed(speed float32) {

	anim.speed = speed
}

// Speed returns the animation speed.
func (anim *Animation) Speed() float32 {

	return anim.speed
}

// Reset resets the animation to the beginning.
func (anim *Animation) Reset() {

	anim.time = anim.start

	// Update all channels
	for i := range anim.channels {
		ch := anim.channels[i]
		ch.Update(anim.start)
	}
}

// SetPaused sets whether the animation is paused.
func (anim *Animation) SetPaused(state bool) {

	anim.paused = state
}

// Paused returns whether the animation is paused.
func (anim *Animation) Paused() bool {

	return anim.paused
}

// SetLoop sets whether the animation is looping.
func (anim *Animation) SetLoop(state bool) {

	anim.loop = state
}

// Loop returns whether the animation is looping.
func (anim *Animation) Loop() bool {

	return anim.loop
}

// SetStart sets the initial time offset value.
func (anim *Animation) SetStart(v float32) {

	anim.start = v
}

// Update interpolates and updates the target values for each channel.
// If the animation is paused, returns false. If the animation is not paused,
// returns true if the input value is inside the key frames ranges or false otherwise.
func (anim *Animation) Update(delta float32) {

	// Check if paused
	if anim.paused {
		return
	}

	// Check if input is less than minimum
	anim.time = anim.time + delta * anim.speed
	if anim.time < anim.minTime {
		return
	}

	// Check if input is greater than maximum
	if anim.time > anim.maxTime {
		if anim.loop {
			anim.time = anim.time - anim.maxTime
		} else {
			anim.time = anim.maxTime - 0.000001
			anim.SetPaused(true)
		}
	}

	// Update all channels
	for i := range anim.channels {
		ch := anim.channels[i]
		ch.Update(anim.time)
	}
}

// AddChannel adds a channel to the animation.
func (anim *Animation) AddChannel(ch IChannel) {

	// TODO (maybe) prevent user from adding two channels of the same type that share target ?

	// Add the channel
	anim.channels = append(anim.channels, ch)

	// Update maxTime and minTime values
	kf := ch.Keyframes()
	firstTime := kf[0]
	if anim.minTime > firstTime {
		anim.minTime = firstTime
	}
	lastTime := kf[len(kf)-1]
	if anim.maxTime < lastTime {
		anim.maxTime = lastTime
	}
}
