// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"os"
)

// Ansi terminal color codes
const (
	csi      = "\x1B["
	black    = "30m"
	red      = "31m"
	green    = "32m"
	yellow   = "33m"
	blue     = "34m"
	magenta  = "35m"
	cyan     = "36m"
	white    = "37m"
	bblack   = "30m"
	bred     = "31;1m"
	bgreen   = "32;1m"
	byellow  = "33;1m"
	bblue    = "34;1m"
	bmagenta = "35;1m"
	bcyan    = "36;1m"
	bwhite   = "37;1m"
)

// Maps log level to color sequence
var colorMap = map[int]string{
	DEBUG: white,
	INFO:  green,
	WARN:  byellow,
	ERROR: bred,
	FATAL: bmagenta,
}

// Console is a console writer used for logging.
type Console struct {
	writer *os.File
	color  bool
}

// NewConsole creates and returns a new logger Console writer
// If color is true, this writer uses Ansi codes to write
// log messages in color accordingly to its level.
func NewConsole(color bool) *Console {

	return &Console{os.Stdout, color}
}

// Write writes the provided logger event to the console.
func (w *Console) Write(event *Event) {

	if w.color {
		w.writer.Write([]byte(csi))
		w.writer.Write([]byte(colorMap[event.level]))
	}
	w.writer.Write([]byte(event.fmsg))
	if w.color {
		w.writer.Write([]byte(csi))
		w.writer.Write([]byte(white))
	}
}

func (w *Console) Close() {

}

func (w *Console) Sync() {

}
