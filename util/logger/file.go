// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"os"
)

// File is a file writer used for logging.
type File struct {
	writer *os.File
}

// NewFile creates and returns a pointer to a new File object along with any error that occurred.
func NewFile(filename string) (*File, error) {

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &File{file}, nil
}

// Write writes the provided logger event to the file.
func (f *File) Write(event *Event) {

	f.writer.Write([]byte(event.fmsg))
}

// Close closes the file.
func (f *File) Close() {

	f.writer.Close()
	f.writer = nil
}

// Sync commits the current contents of the file to stable storage.
func (f *File) Sync() {

	f.writer.Sync()
}
