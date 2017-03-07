// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"os"
)

type File struct {
	writer *os.File
}

func NewFile(filename string) (*File, error) {

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &File{file}, nil
}

func (f *File) Write(event *Event) {

	f.writer.Write([]byte(event.fmsg))
}

func (f *File) Close() {

	f.writer.Close()
	f.writer = nil
}

func (f *File) Sync() {

	f.writer.Sync()
}
