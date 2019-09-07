package gls

import (
	"math"
	"unsafe"
)

// Stats contains counters of WebGL resources being used as well
// the cumulative numbers of some WebGL calls for performance evaluation.
type Stats struct {
	Shaders    int    // Current number of shader programs
	Vaos       int    // Number of Vertex Array Objects
	Buffers    int    // Number of Buffer Objects
	Textures   int    // Number of Textures
	Caphits    uint64 // Cumulative number of hits for Enable/Disable
	UnilocHits uint64 // Cumulative number of uniform location cache hits
	UnilocMiss uint64 // Cumulative number of uniform location cache misses
	Unisets    uint64 // Cumulative number of uniform sets
	Drawcalls  uint64 // Cumulative number of draw calls
}

const (
	capUndef    = 0
	capDisabled = 1
	capEnabled  = 2
	uintUndef   = math.MaxUint32
	intFalse    = 0
	intTrue     = 1
)

const (
	FloatSize = int32(unsafe.Sizeof(float32(0)))
)
