// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package al implements the Go bindings of a subset of the functions of the OpenAL C library.
// The OpenAL documentation can be accessed at https://openal.org/documentation/
package al

/*
#cgo darwin   CFLAGS:  -DGO_DARWIN  -I/usr/local/opt/openal-soft/include/AL -I/usr/include/AL
#cgo freebsd  CFLAGS:  -DGO_FREEBSD -I/usr/local/include/AL
#cgo linux    CFLAGS:  -DGO_LINUX   -I/usr/include/AL
#cgo windows  CFLAGS:  -DGO_WINDOWS -I${SRCDIR}/../windows/openal-soft-1.18.2/include/AL
#cgo darwin   LDFLAGS: -L/usr/local/opt/openal-soft/lib -lopenal
#cgo freebsd  LDFLAGS: -L/usr/local/lib -lopenal
#cgo linux    LDFLAGS: -lopenal
#cgo windows  LDFLAGS: -L${SRCDIR}/../windows/bin -lOpenAL32

#ifdef GO_DARWIN
#include <stdlib.h>
#include "al.h"
#include "alc.h"
#include "efx.h"
#endif

#ifdef GO_FREEBSD
#include <stdlib.h>
#include "al.h"
#include "alc.h"
#include "efx.h"
#endif

#ifdef GO_LINUX
#include <stdlib.h>
#include "al.h"
#include "alc.h"
#include "efx.h"
#endif

#ifdef GO_WINDOWS
#include <stdlib.h>
#include "al.h"
#include "alc.h"
#include "efx.h"
#endif
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// AL constants
const (
	None                    = C.AL_NONE
	False                   = C.AL_FALSE
	True                    = C.AL_TRUE
	SourceRelative          = C.AL_SOURCE_RELATIVE
	ConeInnerAngle          = C.AL_CONE_INNER_ANGLE
	ConeOuterAngle          = C.AL_CONE_OUTER_ANGLE
	Pitch                   = C.AL_PITCH
	Position                = C.AL_POSITION
	Direction               = C.AL_DIRECTION
	Velocity                = C.AL_VELOCITY
	Looping                 = C.AL_LOOPING
	Buffer                  = C.AL_BUFFER
	Gain                    = C.AL_GAIN
	MinGain                 = C.AL_MIN_GAIN
	MaxGain                 = C.AL_MAX_GAIN
	Orientation             = C.AL_ORIENTATION
	SourceState             = C.AL_SOURCE_STATE
	Initial                 = C.AL_INITIAL
	Playing                 = C.AL_PLAYING
	Paused                  = C.AL_PAUSED
	Stopped                 = C.AL_STOPPED
	BuffersQueued           = C.AL_BUFFERS_QUEUED
	BuffersProcessed        = C.AL_BUFFERS_PROCESSED
	ReferenceDistance       = C.AL_REFERENCE_DISTANCE
	RolloffFactor           = C.AL_ROLLOFF_FACTOR
	ConeOuterGain           = C.AL_CONE_OUTER_GAIN
	MaxDistance             = C.AL_MAX_DISTANCE
	SecOffset               = C.AL_SEC_OFFSET
	SampleOffset            = C.AL_SAMPLE_OFFSET
	ByteOffset              = C.AL_BYTE_OFFSET
	SourceType              = C.AL_SOURCE_TYPE
	Static                  = C.AL_STATIC
	Streaming               = C.AL_STREAMING
	Undetermined            = C.AL_UNDETERMINED
	FormatMono8             = C.AL_FORMAT_MONO8
	FormatMono16            = C.AL_FORMAT_MONO16
	FormatStereo8           = C.AL_FORMAT_STEREO8
	FormatStereo16          = C.AL_FORMAT_STEREO16
	Frequency               = C.AL_FREQUENCY
	Bits                    = C.AL_BITS
	Channels                = C.AL_CHANNELS
	Size                    = C.AL_SIZE
	Unused                  = C.AL_UNUSED
	Pending                 = C.AL_PENDING
	Processed               = C.AL_PROCESSED
	NoError                 = C.AL_NO_ERROR
	InvalidName             = C.AL_INVALID_NAME
	InvalidEnum             = C.AL_INVALID_ENUM
	InvalidValue            = C.AL_INVALID_VALUE
	InvalidOperation        = C.AL_INVALID_OPERATION
	OutOfMemory             = C.AL_OUT_OF_MEMORY
	Vendor                  = C.AL_VENDOR
	Version                 = C.AL_VERSION
	Renderer                = C.AL_RENDERER
	Extensions              = C.AL_EXTENSIONS
	DopplerFactor           = C.AL_DOPPLER_FACTOR
	DopplerVelocity         = C.AL_DOPPLER_VELOCITY
	SpeedOfSound            = C.AL_SPEED_OF_SOUND
	DistanceModel           = C.AL_DISTANCE_MODEL
	InverseDistance         = C.AL_INVERSE_DISTANCE
	InverseDistanceClamped  = C.AL_INVERSE_DISTANCE_CLAMPED
	LinearDistance          = C.AL_LINEAR_DISTANCE
	LinearDistanceClamped   = C.AL_LINEAR_DISTANCE_CLAMPED
	ExponentDistance        = C.AL_EXPONENT_DISTANCE
	ExponentDistanceClamped = C.AL_EXPONENT_DISTANCE_CLAMPED
)

// ALC constants
const (
	AttributesSize                = C.ALC_ATTRIBUTES_SIZE
	AllAttributes                 = C.ALC_ALL_ATTRIBUTES
	DefaultDeviceSpecifier        = C.ALC_DEFAULT_DEVICE_SPECIFIER
	DeviceSpecifier               = C.ALC_DEVICE_SPECIFIER
	CtxExtensions                 = C.ALC_EXTENSIONS
	ExtCapture                    = C.ALC_EXT_CAPTURE
	CaptureDeviceSpecifier        = C.ALC_CAPTURE_DEVICE_SPECIFIER
	CaptureDefaultDeviceSpecifier = C.ALC_CAPTURE_DEFAULT_DEVICE_SPECIFIER
	CtxCaptureSamples             = C.ALC_CAPTURE_SAMPLES
	EnumerateAllExt               = C.ALC_ENUMERATE_ALL_EXT
	DefaultAllDevicesSpecifier    = C.ALC_DEFAULT_ALL_DEVICES_SPECIFIER
	AllDevicesSpecifier           = C.ALC_ALL_DEVICES_SPECIFIER
)

// AL EFX extension constants
const (
	EFX_MAJOR_VERSION                       = C.ALC_EFX_MAJOR_VERSION
	EFX_MINOR_VERSION                       = C.ALC_EFX_MINOR_VERSION
	MAX_AUXILIARY_SENDS                     = C.ALC_MAX_AUXILIARY_SENDS
	METERS_PER_UNIT                         = C.AL_METERS_PER_UNIT
	AL_DIRECT_FILTER                        = C.AL_DIRECT_FILTER
	AL_AUXILIARY_SEND_FILTER                = C.AL_AUXILIARY_SEND_FILTER
	AL_AIR_ABSORPTION_FACTOR                = C.AL_AIR_ABSORPTION_FACTOR
	AL_ROOM_ROLLOFF_FACTOR                  = C.AL_ROOM_ROLLOFF_FACTOR
	AL_CONE_OUTER_GAINHF                    = C.AL_CONE_OUTER_GAINHF
	AL_DIRECT_FILTER_GAINHF_AUTO            = C.AL_DIRECT_FILTER_GAINHF_AUTO
	AL_AUXILIARY_SEND_FILTER_GAIN_AUTO      = C.AL_AUXILIARY_SEND_FILTER_GAIN_AUTO
	AL_AUXILIARY_SEND_FILTER_GAINHF_AUTO    = C.AL_AUXILIARY_SEND_FILTER_GAINHF_AUTO
	AL_REVERB_DENSITY                       = C.AL_REVERB_DENSITY
	AL_REVERB_DIFFUSION                     = C.AL_REVERB_DIFFUSION
	AL_REVERB_GAIN                          = C.AL_REVERB_GAIN
	AL_REVERB_GAINHF                        = C.AL_REVERB_GAINHF
	AL_REVERB_DECAY_TIME                    = C.AL_REVERB_DECAY_TIME
	AL_REVERB_DECAY_HFRATIO                 = C.AL_REVERB_DECAY_HFRATIO
	AL_REVERB_REFLECTIONS_GAIN              = C.AL_REVERB_REFLECTIONS_GAIN
	AL_REVERB_REFLECTIONS_DELAY             = C.AL_REVERB_REFLECTIONS_DELAY
	AL_REVERB_LATE_REVERB_GAIN              = C.AL_REVERB_LATE_REVERB_GAIN
	AL_REVERB_LATE_REVERB_DELAY             = C.AL_REVERB_LATE_REVERB_DELAY
	AL_REVERB_AIR_ABSORPTION_GAINHF         = C.AL_REVERB_AIR_ABSORPTION_GAINHF
	AL_REVERB_ROOM_ROLLOFF_FACTOR           = C.AL_REVERB_ROOM_ROLLOFF_FACTOR
	AL_REVERB_DECAY_HFLIMIT                 = C.AL_REVERB_DECAY_HFLIMIT
	AL_EAXREVERB_DENSITY                    = C.AL_EAXREVERB_DENSITY
	AL_EAXREVERB_DIFFUSION                  = C.AL_EAXREVERB_DIFFUSION
	AL_EAXREVERB_GAIN                       = C.AL_EAXREVERB_GAIN
	AL_EAXREVERB_GAINHF                     = C.AL_EAXREVERB_GAINHF
	AL_EAXREVERB_GAINLF                     = C.AL_EAXREVERB_GAINLF
	AL_EAXREVERB_DECAY_TIME                 = C.AL_EAXREVERB_DECAY_TIME
	AL_EAXREVERB_DECAY_HFRATIO              = C.AL_EAXREVERB_DECAY_HFRATIO
	AL_EAXREVERB_DECAY_LFRATIO              = C.AL_EAXREVERB_DECAY_LFRATIO
	AL_EAXREVERB_REFLECTIONS_GAIN           = C.AL_EAXREVERB_REFLECTIONS_GAIN
	AL_EAXREVERB_REFLECTIONS_DELAY          = C.AL_EAXREVERB_REFLECTIONS_DELAY
	AL_EAXREVERB_REFLECTIONS_PAN            = C.AL_EAXREVERB_REFLECTIONS_PAN
	AL_EAXREVERB_LATE_REVERB_GAIN           = C.AL_EAXREVERB_LATE_REVERB_GAIN
	AL_EAXREVERB_LATE_REVERB_DELAY          = C.AL_EAXREVERB_LATE_REVERB_DELAY
	AL_EAXREVERB_LATE_REVERB_PAN            = C.AL_EAXREVERB_LATE_REVERB_PAN
	AL_EAXREVERB_ECHO_TIME                  = C.AL_EAXREVERB_ECHO_TIME
	AL_EAXREVERB_ECHO_DEPTH                 = C.AL_EAXREVERB_ECHO_DEPTH
	AL_EAXREVERB_MODULATION_TIME            = C.AL_EAXREVERB_MODULATION_TIME
	AL_EAXREVERB_MODULATION_DEPTH           = C.AL_EAXREVERB_MODULATION_DEPTH
	AL_EAXREVERB_AIR_ABSORPTION_GAINHF      = C.AL_EAXREVERB_AIR_ABSORPTION_GAINHF
	AL_EAXREVERB_HFREFERENCE                = C.AL_EAXREVERB_HFREFERENCE
	AL_EAXREVERB_LFREFERENCE                = C.AL_EAXREVERB_LFREFERENCE
	AL_EAXREVERB_ROOM_ROLLOFF_FACTOR        = C.AL_EAXREVERB_ROOM_ROLLOFF_FACTOR
	AL_EAXREVERB_DECAY_HFLIMIT              = C.AL_EAXREVERB_DECAY_HFLIMIT
	AL_CHORUS_WAVEFORM                      = C.AL_CHORUS_WAVEFORM
	AL_CHORUS_PHASE                         = C.AL_CHORUS_PHASE
	AL_CHORUS_RATE                          = C.AL_CHORUS_RATE
	AL_CHORUS_DEPTH                         = C.AL_CHORUS_DEPTH
	AL_CHORUS_FEEDBACK                      = C.AL_CHORUS_FEEDBACK
	AL_CHORUS_DELAY                         = C.AL_CHORUS_DELAY
	AL_DISTORTION_EDGE                      = C.AL_DISTORTION_EDGE
	AL_DISTORTION_GAIN                      = C.AL_DISTORTION_GAIN
	AL_DISTORTION_LOWPASS_CUTOFF            = C.AL_DISTORTION_LOWPASS_CUTOFF
	AL_DISTORTION_EQCENTER                  = C.AL_DISTORTION_EQCENTER
	AL_DISTORTION_EQBANDWIDTH               = C.AL_DISTORTION_EQBANDWIDTH
	AL_ECHO_DELAY                           = C.AL_ECHO_DELAY
	AL_ECHO_LRDELAY                         = C.AL_ECHO_LRDELAY
	AL_ECHO_DAMPING                         = C.AL_ECHO_DAMPING
	AL_ECHO_FEEDBACK                        = C.AL_ECHO_FEEDBACK
	AL_ECHO_SPREAD                          = C.AL_ECHO_SPREAD
	AL_FLANGER_WAVEFORM                     = C.AL_FLANGER_WAVEFORM
	AL_FLANGER_PHASE                        = C.AL_FLANGER_PHASE
	AL_FLANGER_RATE                         = C.AL_FLANGER_RATE
	AL_FLANGER_DEPTH                        = C.AL_FLANGER_DEPTH
	AL_FLANGER_FEEDBACK                     = C.AL_FLANGER_FEEDBACK
	AL_FLANGER_DELAY                        = C.AL_FLANGER_DELAY
	AL_FREQUENCY_SHIFTER_FREQUENCY          = C.AL_FREQUENCY_SHIFTER_FREQUENCY
	AL_FREQUENCY_SHIFTER_LEFT_DIRECTION     = C.AL_FREQUENCY_SHIFTER_LEFT_DIRECTION
	AL_FREQUENCY_SHIFTER_RIGHT_DIRECTION    = C.AL_FREQUENCY_SHIFTER_RIGHT_DIRECTION
	AL_VOCAL_MORPHER_PHONEMEA               = C.AL_VOCAL_MORPHER_PHONEMEA
	AL_VOCAL_MORPHER_PHONEMEA_COARSE_TUNING = C.AL_VOCAL_MORPHER_PHONEMEA_COARSE_TUNING
	AL_VOCAL_MORPHER_PHONEMEB               = C.AL_VOCAL_MORPHER_PHONEMEB
	AL_VOCAL_MORPHER_PHONEMEB_COARSE_TUNING = C.AL_VOCAL_MORPHER_PHONEMEB_COARSE_TUNING
	AL_VOCAL_MORPHER_WAVEFORM               = C.AL_VOCAL_MORPHER_WAVEFORM
	AL_VOCAL_MORPHER_RATE                   = C.AL_VOCAL_MORPHER_RATE
	AL_PITCH_SHIFTER_COARSE_TUNE            = C.AL_PITCH_SHIFTER_COARSE_TUNE
	AL_PITCH_SHIFTER_FINE_TUNE              = C.AL_PITCH_SHIFTER_FINE_TUNE
	AL_RING_MODULATOR_FREQUENCY             = C.AL_RING_MODULATOR_FREQUENCY
	AL_RING_MODULATOR_HIGHPASS_CUTOFF       = C.AL_RING_MODULATOR_HIGHPASS_CUTOFF
	AL_RING_MODULATOR_WAVEFORM              = C.AL_RING_MODULATOR_WAVEFORM
	AL_AUTOWAH_ATTACK_TIME                  = C.AL_AUTOWAH_ATTACK_TIME
	AL_AUTOWAH_RELEASE_TIME                 = C.AL_AUTOWAH_RELEASE_TIME
	AL_AUTOWAH_RESONANCE                    = C.AL_AUTOWAH_RESONANCE
	AL_AUTOWAH_PEAK_GAIN                    = C.AL_AUTOWAH_PEAK_GAIN
	AL_COMPRESSOR_ONOFF                     = C.AL_COMPRESSOR_ONOFF
	AL_EQUALIZER_LOW_GAIN                   = C.AL_EQUALIZER_LOW_GAIN
	AL_EQUALIZER_LOW_CUTOFF                 = C.AL_EQUALIZER_LOW_CUTOFF
	AL_EQUALIZER_MID1_GAIN                  = C.AL_EQUALIZER_MID1_GAIN
	AL_EQUALIZER_MID1_CENTER                = C.AL_EQUALIZER_MID1_CENTER
	AL_EQUALIZER_MID1_WIDTH                 = C.AL_EQUALIZER_MID1_WIDTH
	AL_EQUALIZER_MID2_GAIN                  = C.AL_EQUALIZER_MID2_GAIN
	AL_EQUALIZER_MID2_CENTER                = C.AL_EQUALIZER_MID2_CENTER
	AL_EQUALIZER_MID2_WIDTH                 = C.AL_EQUALIZER_MID2_WIDTH
	AL_EQUALIZER_HIGH_GAIN                  = C.AL_EQUALIZER_HIGH_GAIN
	AL_EQUALIZER_HIGH_CUTOFF                = C.AL_EQUALIZER_HIGH_CUTOFF
	AL_EFFECT_FIRST_PARAMETER               = C.AL_EFFECT_FIRST_PARAMETER
	AL_EFFECT_LAST_PARAMETER                = C.AL_EFFECT_LAST_PARAMETER
	AL_EFFECT_TYPE                          = C.AL_EFFECT_TYPE
	AL_EFFECT_NULL                          = C.AL_EFFECT_NULL
	AL_EFFECT_REVERB                        = C.AL_EFFECT_REVERB
	AL_EFFECT_CHORUS                        = C.AL_EFFECT_CHORUS
	AL_EFFECT_DISTORTION                    = C.AL_EFFECT_DISTORTION
	AL_EFFECT_ECHO                          = C.AL_EFFECT_ECHO
	AL_EFFECT_FLANGER                       = C.AL_EFFECT_FLANGER
	AL_EFFECT_FREQUENCY_SHIFTER             = C.AL_EFFECT_FREQUENCY_SHIFTER
	AL_EFFECT_VOCAL_MORPHER                 = C.AL_EFFECT_VOCAL_MORPHER
	AL_EFFECT_PITCH_SHIFTER                 = C.AL_EFFECT_PITCH_SHIFTER
	AL_EFFECT_RING_MODULATOR                = C.AL_EFFECT_RING_MODULATOR
	AL_EFFECT_AUTOWAH                       = C.AL_EFFECT_AUTOWAH
	AL_EFFECT_COMPRESSOR                    = C.AL_EFFECT_COMPRESSOR
	AL_EFFECT_EQUALIZER                     = C.AL_EFFECT_EQUALIZER
	AL_EFFECT_EAXREVERB                     = C.AL_EFFECT_EAXREVERB
	AL_EFFECTSLOT_EFFECT                    = C.AL_EFFECTSLOT_EFFECT
	AL_EFFECTSLOT_GAIN                      = C.AL_EFFECTSLOT_GAIN
	AL_EFFECTSLOT_AUXILIARY_SEND_AUTO       = C.AL_EFFECTSLOT_AUXILIARY_SEND_AUTO
	AL_EFFECTSLOT_NULL                      = C.AL_EFFECTSLOT_NULL
	AL_LOWPASS_GAIN                         = C.AL_LOWPASS_GAIN
	AL_LOWPASS_GAINHF                       = C.AL_LOWPASS_GAINHF
	AL_HIGHPASS_GAIN                        = C.AL_HIGHPASS_GAIN
	AL_HIGHPASS_GAINLF                      = C.AL_HIGHPASS_GAINLF
	AL_BANDPASS_GAIN                        = C.AL_BANDPASS_GAIN
	AL_BANDPASS_GAINLF                      = C.AL_BANDPASS_GAINLF
	AL_BANDPASS_GAINHF                      = C.AL_BANDPASS_GAINHF
	AL_FILTER_FIRST_PARAMETER               = C.AL_FILTER_FIRST_PARAMETER
	AL_FILTER_LAST_PARAMETER                = C.AL_FILTER_LAST_PARAMETER
	AL_FILTER_TYPE                          = C.AL_FILTER_TYPE
	AL_FILTER_NULL                          = C.AL_FILTER_NULL
	AL_FILTER_LOWPASS                       = C.AL_FILTER_LOWPASS
	AL_FILTER_HIGHPASS                      = C.AL_FILTER_HIGHPASS
	AL_FILTER_BANDPASS                      = C.AL_FILTER_BANDPASS
)

var errCodes = map[uint]string{
	C.AL_INVALID_NAME:      "AL_INVALID_NAME",
	C.AL_INVALID_ENUM:      "AL_INVALID_ENUM",
	C.AL_INVALID_VALUE:     "AL_INVALID_VALUE",
	C.AL_INVALID_OPERATION: "AL_INVALID_OPERATION",
	C.AL_OUT_OF_MEMORY:     "AL_OUT_OF_MEMORY",
}

type Device struct {
	cdev *C.ALCdevice
}

type Context struct {
	cctx *C.ALCcontext
}

// Statistics
type Stats struct {
	Sources  int   // Current number of sources
	Buffers  int   // Current number of buffers
	CgoCalls int64 // Accumulated cgo calls
	Callocs  int   // Current number of C allocations
}

// Maps C pointer to device to Go pointer to Device
var mapDevice = map[*C.ALCdevice]*Device{}

// Global statistics structure
var stats Stats

// GetStats returns copy of the statistics structure
func GetStats() Stats {

	return stats
}

func checkCtxError(dev *Device) {

	err := CtxGetError(dev)
	if err != nil {
		panic(err)
	}
}

func CreateContext(dev *Device, attrlist []int) (*Context, error) {

	var plist unsafe.Pointer
	if len(attrlist) != 0 {
		plist = (unsafe.Pointer)(&attrlist[0])
	}
	ctx := C.alcCreateContext(dev.cdev, (*C.ALCint)(plist))
	if ctx != nil {
		return &Context{ctx}, nil
	}
	return nil, fmt.Errorf("%s", errCodes[uint(C.alcGetError(dev.cdev))])
}

func MakeContextCurrent(ctx *Context) error {

	cres := C.alcMakeContextCurrent(ctx.cctx)
	if cres == C.ALC_TRUE {
		return nil
	}
	return fmt.Errorf("%s", errCodes[uint(C.alGetError())])
}

func ProcessContext(ctx *Context) {

	C.alcProcessContext(ctx.cctx)
}

func SuspendContext(ctx *Context) {

	C.alcSuspendContext(ctx.cctx)
}

func DestroyContext(ctx *Context) {

	C.alcDestroyContext(ctx.cctx)
}

func GetContextsDevice(ctx *Context) *Device {

	cdev := C.alcGetContextsDevice(ctx.cctx)
	if cdev == nil {
		return nil
	}
	return mapDevice[cdev]
}

func OpenDevice(name string) (*Device, error) {

	cstr := (*C.ALCchar)(C.CString(name))
	defer C.free(unsafe.Pointer(cstr))
	cdev := C.alcOpenDevice(cstr)
	if cdev != nil {
		dev := &Device{cdev}
		mapDevice[cdev] = dev
		return dev, nil
	}
	return nil, fmt.Errorf("%s", errCodes[uint(C.alGetError())])
}

func CloseDevice(dev *Device) error {

	cres := C.alcCloseDevice(dev.cdev)
	if cres == C.ALC_TRUE {
		delete(mapDevice, dev.cdev)
		return nil
	}
	return fmt.Errorf("%s", errCodes[uint(C.alGetError())])
}

func CtxGetError(dev *Device) error {

	cerr := C.alcGetError(dev.cdev)
	if cerr == C.AL_NONE {
		return nil
	}
	return fmt.Errorf("%s", errCodes[uint(cerr)])
}

func CtxIsExtensionPresent(dev *Device, extname string) bool {

	cname := (*C.ALCchar)(C.CString(extname))
	defer C.free(unsafe.Pointer(cname))
	cres := C.alcIsExtensionPresent(dev.cdev, cname)
	if cres == C.AL_TRUE {
		return true
	}
	return false
}

func CtxGetEnumValue(dev *Device, enumName string) uint32 {

	cname := (*C.ALCchar)(C.CString(enumName))
	defer C.free(unsafe.Pointer(cname))
	cres := C.alcGetEnumValue(dev.cdev, cname)
	return uint32(cres)
}

func CtxGetString(dev *Device, param uint) string {

	var cdev *C.ALCdevice = nil
	if dev != nil {
		cdev = dev.cdev
	}
	cstr := C.alcGetString(cdev, C.ALCenum(param))
	return C.GoString((*C.char)(cstr))
}

func CtxGetIntegerv(dev *Device, param uint32, values []int32) {

	C.alcGetIntegerv(dev.cdev, C.ALCenum(param), C.ALCsizei(len(values)), (*C.ALCint)(unsafe.Pointer(&values[0])))
}

func CaptureOpenDevice(devname string, frequency uint32, format uint32, buffersize uint32) (*Device, error) {

	cstr := (*C.ALCchar)(C.CString(devname))
	defer C.free(unsafe.Pointer(cstr))
	cdev := C.alcCaptureOpenDevice(cstr, C.ALCuint(frequency), C.ALCenum(format), C.ALCsizei(buffersize))
	if cdev != nil {
		dev := &Device{cdev}
		mapDevice[cdev] = dev
		return dev, nil
	}
	return nil, fmt.Errorf("%s", errCodes[uint(C.alGetError())])
}

func CaptureCloseDevice(dev *Device) error {

	cres := C.alcCaptureCloseDevice(dev.cdev)
	if cres == C.AL_TRUE {
		return nil
	}
	return fmt.Errorf("%s", errCodes[uint(C.alGetError())])
}

func CaptureStart(dev *Device) {

	C.alcCaptureStart(dev.cdev)
	checkCtxError(dev)
}

func CaptureStop(dev *Device) {

	C.alcCaptureStop(dev.cdev)
	checkCtxError(dev)
}

func CaptureSamples(dev *Device, buffer []byte, nsamples uint) {

	C.alcCaptureSamples(dev.cdev, unsafe.Pointer(&buffer[0]), C.ALCsizei(nsamples))
	checkCtxError(dev)
}

func Enable(capability uint) {

	C.alEnable(C.ALenum(capability))
}

func Disable(capability uint) {

	C.alDisable(C.ALenum(capability))
}

func IsEnabled(capability uint) bool {

	cres := C.alIsEnabled(C.ALenum(capability))
	if cres == C.AL_TRUE {
		return true
	}
	return false
}

func GetString(param uint32) string {

	cstr := C.alGetString(C.ALenum(param))
	return C.GoString((*C.char)(cstr))
}

func GetBooleanv(param uint32, values []bool) {

	cvals := make([]C.ALboolean, len(values))
	C.alGetBooleanv(C.ALenum(param), &cvals[0])
	for i := 0; i < len(cvals); i++ {
		if cvals[i] == C.AL_TRUE {
			values[i] = true
		} else {
			values[i] = false
		}
	}
}

func GetIntegerv(param uint32, values []int32) {

	C.alGetIntegerv(C.ALenum(param), (*C.ALint)(unsafe.Pointer(&values[0])))
}

func GetFloatv(param uint32, values []float32) {

	C.alGetFloatv(C.ALenum(param), (*C.ALfloat)(unsafe.Pointer(&values[0])))
}

func GetDoublev(param uint32, values []float64) {

	C.alGetDoublev(C.ALenum(param), (*C.ALdouble)(unsafe.Pointer(&values[0])))
}

func GetBoolean(param uint32) bool {

	cres := C.alGetBoolean(C.ALenum(param))
	if cres == C.AL_TRUE {
		return true
	}
	return false
}

func GetInteger(param uint32) int32 {

	cres := C.alGetInteger(C.ALenum(param))
	return int32(cres)
}

func GetFloat(param uint32) float32 {

	cres := C.alGetFloat(C.ALenum(param))
	return float32(cres)
}

func GetDouble(param uint32) float64 {

	cres := C.alGetDouble(C.ALenum(param))
	return float64(cres)
}

func GetError() error {

	cerr := C.alGetError()
	if cerr == C.AL_NONE {
		return nil
	}
	return fmt.Errorf("%s", errCodes[uint(cerr)])
}

func IsExtensionPresent(extName string) bool {

	cstr := (*C.ALchar)(C.CString(extName))
	defer C.free(unsafe.Pointer(cstr))
	cres := C.alIsExtensionPresent(cstr)
	if cres == 0 {
		return false
	}
	return true
}

func GetEnumValue(enam string) uint32 {

	cenam := (*C.ALchar)(C.CString(enam))
	defer C.free(unsafe.Pointer(cenam))
	cres := C.alGetEnumValue(cenam)
	return uint32(cres)
}

func Listenerf(param uint32, value float32) {

	C.alListenerf(C.ALenum(param), C.ALfloat(value))
}

func Listener3f(param uint32, value1, value2, value3 float32) {

	C.alListener3f(C.ALenum(param), C.ALfloat(value1), C.ALfloat(value2), C.ALfloat(value3))
}

func Listenerfv(param uint32, values []float32) {

	C.alListenerfv(C.ALenum(param), (*C.ALfloat)(unsafe.Pointer(&values[0])))
}

func Listeneri(param uint32, value int32) {

	C.alListeneri(C.ALenum(param), C.ALint(value))
}

func Listener3i(param uint32, value1, value2, value3 int32) {

	C.alListener3i(C.ALenum(param), C.ALint(value1), C.ALint(value2), C.ALint(value3))
}

func Listeneriv(param uint32, values []int32) {

	C.alListeneriv(C.ALenum(param), (*C.ALint)(unsafe.Pointer(&values[0])))
}

func GetListenerf(param uint32) float32 {

	var cval C.ALfloat
	C.alGetListenerf(C.ALenum(param), &cval)
	return float32(cval)
}

func GetListener3f(param uint32) (float32, float32, float32) {

	var cval1 C.ALfloat
	var cval2 C.ALfloat
	var cval3 C.ALfloat
	C.alGetListener3f(C.ALenum(param), &cval1, &cval2, &cval3)
	return float32(cval1), float32(cval2), float32(cval3)
}

func GetListenerfv(param uint32, values []uint32) {

	C.alGetListenerfv(C.ALenum(param), (*C.ALfloat)(unsafe.Pointer(&values[0])))
}

func GetListeneri(param uint32) int32 {

	var cval C.ALint
	C.alGetListeneri(C.ALenum(param), &cval)
	return int32(cval)
}

func GetListener3i(param uint32) (int32, int32, int32) {

	var cval1 C.ALint
	var cval2 C.ALint
	var cval3 C.ALint
	C.alGetListener3i(C.ALenum(param), &cval1, &cval2, &cval3)
	return int32(cval1), int32(cval2), int32(cval3)
}

func GetListeneriv(param uint32, values []int32) {

	if len(values) < 3 {
		panic("Slice length less than minimum")
	}
	C.alGetListeneriv(C.ALenum(param), (*C.ALint)(unsafe.Pointer(&values[0])))
}

func GenSource() uint32 {

	var csource C.ALuint
	C.alGenSources(1, &csource)
	stats.Sources++
	return uint32(csource)
}

func GenSources(sources []uint32) {

	C.alGenSources(C.ALsizei(len(sources)), (*C.ALuint)(unsafe.Pointer(&sources[0])))
	stats.Sources += len(sources)
}

func DeleteSource(source uint32) {

	C.alDeleteSources(1, (*C.ALuint)(unsafe.Pointer(&source)))
	stats.Sources--
}

func DeleteSources(sources []uint32) {

	C.alDeleteSources(C.ALsizei(len(sources)), (*C.ALuint)(unsafe.Pointer(&sources[0])))
	stats.Sources -= len(sources)
}

func IsSource(source uint32) bool {

	cres := C.alIsSource(C.ALuint(source))
	if cres == C.AL_TRUE {
		return true
	}
	return false
}

func Sourcef(source uint32, param uint32, value float32) {

	C.alSourcef(C.ALuint(source), C.ALenum(param), C.ALfloat(value))
}

func Source3f(source uint32, param uint32, value1, value2, value3 float32) {

	C.alSource3f(C.ALuint(source), C.ALenum(param), C.ALfloat(value1), C.ALfloat(value2), C.ALfloat(value3))
}

func Sourcefv(source uint32, param uint32, values []float32) {

	if len(values) < 3 {
		panic("Slice length less than minimum")
	}
	C.alSourcefv(C.ALuint(source), C.ALenum(param), (*C.ALfloat)(unsafe.Pointer(&values[0])))
}

func Sourcei(source uint32, param uint32, value int32) {

	C.alSourcei(C.ALuint(source), C.ALenum(param), C.ALint(value))
}

func Source3i(source uint32, param uint32, value1, value2, value3 int32) {

	C.alSource3i(C.ALuint(source), C.ALenum(param), C.ALint(value1), C.ALint(value2), C.ALint(value3))
}

func Sourceiv(source uint32, param uint32, values []int32) {

	if len(values) < 3 {
		panic("Slice length less than minimum")
	}
	C.alSourceiv(C.ALuint(source), C.ALenum(param), (*C.ALint)(unsafe.Pointer(&values[0])))
}

func GetSourcef(source uint32, param uint32) float32 {

	var value C.ALfloat
	C.alGetSourcef(C.ALuint(source), C.ALenum(param), &value)
	return float32(value)
}

func GetSource3f(source uint32, param uint32) (float32, float32, float32) {

	var cval1 C.ALfloat
	var cval2 C.ALfloat
	var cval3 C.ALfloat
	C.alGetSource3f(C.ALuint(source), C.ALenum(param), &cval1, &cval2, &cval3)
	return float32(cval1), float32(cval2), float32(cval3)
}

func GetSourcefv(source uint32, param uint32, values []float32) {

	if len(values) < 3 {
		panic("Slice length less than minimum")
	}
	C.alGetSourcefv(C.ALuint(source), C.ALenum(param), (*C.ALfloat)(unsafe.Pointer(&values[0])))
}

func GetSourcei(source uint32, param uint32) int32 {

	var value C.ALint
	C.alGetSourcei(C.ALuint(source), C.ALenum(param), &value)
	return int32(value)
}

func GetSource3i(source uint32, param uint32) (int32, int32, int32) {

	var cval1 C.ALint
	var cval2 C.ALint
	var cval3 C.ALint
	C.alGetSource3i(C.ALuint(source), C.ALenum(param), &cval1, &cval2, &cval3)
	return int32(cval1), int32(cval2), int32(cval3)
}

func GetSourceiv(source uint32, param uint32, values []int32) {

	if len(values) < 3 {
		panic("Slice length less than minimum")
	}
	C.alGetSourceiv(C.ALuint(source), C.ALenum(param), (*C.ALint)(unsafe.Pointer(&values[0])))
}

func SourcePlayv(sources []uint32) {

	C.alSourcePlayv(C.ALsizei(len(sources)), (*C.ALuint)(unsafe.Pointer(&sources[0])))
}

func SourceStopv(sources []uint32) {

	C.alSourceStopv(C.ALsizei(len(sources)), (*C.ALuint)(unsafe.Pointer(&sources[0])))
}

func SourceRewindv(sources []uint32) {

	C.alSourceRewindv(C.ALsizei(len(sources)), (*C.ALuint)(unsafe.Pointer(&sources[0])))
}

func SourcePausev(sources []uint32) {

	C.alSourcePausev(C.ALsizei(len(sources)), (*C.ALuint)(unsafe.Pointer(&sources[0])))
}

func SourcePlay(source uint32) {

	C.alSourcePlay(C.ALuint(source))
}

func SourceStop(source uint32) {

	C.alSourceStop(C.ALuint(source))
}

func SourceRewind(source uint32) {

	C.alSourceRewind(C.ALuint(source))
}

func SourcePause(source uint32) {

	C.alSourcePause(C.ALuint(source))
}

func SourceQueueBuffers(source uint32, buffers ...uint32) {

	C.alSourceQueueBuffers(C.ALuint(source), C.ALsizei(len(buffers)), (*C.ALuint)(unsafe.Pointer(&buffers[0])))
}

func SourceUnqueueBuffers(source uint32, n uint32, buffers []uint32) {

	removed := make([]C.ALuint, n)
	C.alSourceUnqueueBuffers(C.ALuint(source), C.ALsizei(n), &removed[0])
}

func GenBuffers(n uint32) []uint32 {

	buffers := make([]uint32, n)
	C.alGenBuffers(C.ALsizei(len(buffers)), (*C.ALuint)(unsafe.Pointer(&buffers[0])))
	return buffers
}

func DeleteBuffers(buffers []uint32) {

	C.alDeleteBuffers(C.ALsizei(len(buffers)), (*C.ALuint)(unsafe.Pointer(&buffers[0])))
}

func IsBuffer(buffer uint32) bool {

	cres := C.alIsBuffer(C.ALuint(buffer))
	if cres == C.AL_TRUE {
		return true
	}
	return false
}

func BufferData(buffer uint32, format uint32, data unsafe.Pointer, size uint32, freq uint32) {

	C.alBufferData(C.ALuint(buffer), C.ALenum(format), data, C.ALsizei(size), C.ALsizei(freq))
}

func Bufferf(buffer uint32, param uint32, value float32) {

	C.alBufferf(C.ALuint(buffer), C.ALenum(param), C.ALfloat(value))
}

func Buffer3f(buffer uint32, param uint32, value1, value2, value3 float32) {

	C.alBuffer3f(C.ALuint(buffer), C.ALenum(param), C.ALfloat(value1), C.ALfloat(value2), C.ALfloat(value3))
}

func Bufferfv(buffer uint32, param uint32, values []float32) {

	C.alBufferfv(C.ALuint(buffer), C.ALenum(param), (*C.ALfloat)(unsafe.Pointer(&values[0])))
}

func Bufferi(buffer uint32, param uint32, value int32) {

	C.alBufferi(C.ALuint(buffer), C.ALenum(param), C.ALint(value))
}

func Buffer3i(buffer uint32, param uint32, value1, value2, value3 int32) {

	C.alBuffer3i(C.ALuint(buffer), C.ALenum(param), C.ALint(value1), C.ALint(value2), C.ALint(value3))
}

func Bufferiv(buffer uint32, param uint32, values []int32) {

	C.alBufferiv(C.ALuint(buffer), C.ALenum(param), (*C.ALint)(unsafe.Pointer(&values[0])))
}

func GetBufferf(buffer uint32, param uint32) float32 {

	var value C.ALfloat
	C.alGetBufferf(C.ALuint(buffer), C.ALenum(param), &value)
	return float32(value)
}

func GetBuffer3f(buffer uint32, param uint32) (v1 float32, v2 float32, v3 float32) {

	var value1, value2, value3 C.ALfloat
	C.alGetBuffer3f(C.ALuint(buffer), C.ALenum(param), &value1, &value2, &value3)
	return float32(value1), float32(value2), float32(value3)
}

func GetBufferfv(buffer uint32, param uint32, values []float32) {

	C.alGetBufferfv(C.ALuint(buffer), C.ALenum(param), (*C.ALfloat)(unsafe.Pointer(&values[0])))
}

func GetBufferi(buffer uint32, param uint32) int32 {

	var value C.ALint
	C.alGetBufferi(C.ALuint(buffer), C.ALenum(param), &value)
	return int32(value)
}

func GetBuffer3i(buffer uint32, param uint32) (int32, int32, int32) {

	var value1, value2, value3 C.ALint
	C.alGetBuffer3i(C.ALuint(buffer), C.ALenum(param), &value1, &value2, &value3)
	return int32(value1), int32(value2), int32(value3)
}

func GetBufferiv(buffer uint32, param uint32, values []int32) {

	C.alGetBufferiv(C.ALuint(buffer), C.ALenum(param), (*C.ALint)(unsafe.Pointer(&values[0])))
}
