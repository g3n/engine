// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Levels to filter log output
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Flags used to format the log date/time
const (
	// Show date
	FDATE = 1 << iota
	// Show hour, minutes and seconds
	FTIME
	// Show milliseconds after FTIME
	FMILIS
	// Show microseconds after FTIME
	FMICROS
	// Show nanoseconfs after TIME
	FNANOS
)

// List of level names
var levelNames = [...]string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

// Default logger and global mutex
var Default *Logger = nil
var rootLoggers = []*Logger{}
var mutex sync.Mutex

// Interface for all logger writers
type LoggerWriter interface {
	Write(*Event)
	Close()
	Sync()
}

// Logger Object state structure
type Logger struct {
	name     string
	prefix   string
	enabled  bool
	level    int
	format   int
	outputs  []LoggerWriter
	parent   *Logger
	children []*Logger
}

// Logger event passed from the logger to its writers.
type Event struct {
	time    time.Time
	level   int
	usermsg string
	fmsg    string
}

// creates the default logger
func init() {
	Default = New("G3N", nil)
	Default.SetFormat(FTIME | FMICROS)
	Default.AddWriter(NewConsole(false))
}

// New() creates and returns a new logger with the specified name.
// If a parent logger is specified, the created logger inherits the
// parent's configuration.
func New(name string, parent *Logger) *Logger {

	self := new(Logger)
	self.name = name
	self.prefix = name
	self.enabled = true
	self.level = ERROR
	self.format = FDATE | FTIME | FMICROS
	self.outputs = make([]LoggerWriter, 0)
	self.children = make([]*Logger, 0)
	self.parent = parent
	if parent != nil {
		self.prefix = parent.prefix + "/" + name
		self.enabled = parent.enabled
		self.level = parent.level
		self.format = parent.format
		parent.children = append(parent.children, self)
	} else {
		rootLoggers = append(rootLoggers, self)
	}
	return self
}

// SetLevel set the current level of this logger
// Only log messages with levels with the same or higher
// priorities than the current level will be emitted.
func (self *Logger) SetLevel(level int) {

	if level < DEBUG || level > FATAL {
		return
	}
	self.level = level
}

// SetLevelByName sets the current level of this logger by level name:
// debug|info|warn|error|fatal (case ignored.)
// Only log messages with levels with the same or higher
// priorities than the current level will be emitted.
func (self *Logger) SetLevelByName(lname string) error {
	var level int

	lname = strings.ToUpper(lname)
	for level = 0; level < len(levelNames); level++ {
		if lname == levelNames[level] {
			self.level = level
			return nil
		}
	}
	return fmt.Errorf("Invalid log level name: %s", lname)
}

// SetFormat sets the logger date/time message format
func (self *Logger) SetFormat(format int) {

	self.format = format
}

// AddWriter adds a writer to the current outputs of this logger.
func (self *Logger) AddWriter(writer LoggerWriter) {

	self.outputs = append(self.outputs, writer)
}

// RemoveWriter removes the specified writer from  the current outputs of this logger.
func (self *Logger) RemoveWriter(writer LoggerWriter) {

	for pos, w := range self.outputs {
		if w != writer {
			continue
		}
		self.outputs = append(self.outputs[:pos], self.outputs[pos+1:]...)
	}
}

// EnableChild enables or disables this logger child logger with
// the specified name.
func (self *Logger) EnableChild(name string, state bool) {

	for _, c := range self.children {
		if c.name == name {
			c.enabled = state
		}
	}
}

// Debug emits a DEBUG level log message
func (self *Logger) Debug(format string, v ...interface{}) {

	self.Log(DEBUG, format, v...)
}

// Info emits an INFO level log message
func (self *Logger) Info(format string, v ...interface{}) {

	self.Log(INFO, format, v...)
}

// Warn emits a WARN level log message
func (self *Logger) Warn(format string, v ...interface{}) {

	self.Log(WARN, format, v...)
}

// Error emits an ERROR level log message
func (self *Logger) Error(format string, v ...interface{}) {

	self.Log(ERROR, format, v...)
}

// Fatal emits a FATAL level log message
func (self *Logger) Fatal(format string, v ...interface{}) {

	self.Log(FATAL, format, v...)
}

// Logs emits a log message with the specified level
func (self *Logger) Log(level int, format string, v ...interface{}) {

	// Ignores message if logger not enabled or with level bellow the current one.
	if !self.enabled || level < self.level {
		return
	}

	// Formats date
	now := time.Now().UTC()
	year, month, day := now.Date()
	hour, min, sec := now.Clock()
	fdate := []string{}

	if self.format&FDATE != 0 {
		fdate = append(fdate, fmt.Sprintf("%04d/%02d/%02d", year, month, day))
	}
	if self.format&FTIME != 0 {
		if len(fdate) > 0 {
			fdate = append(fdate, "-")
		}
		fdate = append(fdate, fmt.Sprintf("%02d:%02d:%02d", hour, min, sec))
		var sdecs string
		if self.format&FMILIS != 0 {
			sdecs = fmt.Sprintf(".%.03d", now.Nanosecond()/1000000)
		} else if self.format&FMICROS != 0 {
			sdecs = fmt.Sprintf(".%.06d", now.Nanosecond()/1000)
		} else if self.format&FNANOS != 0 {
			sdecs = fmt.Sprintf(".%.09d", now.Nanosecond())
		}
		fdate = append(fdate, sdecs)
	}

	// Formats message
	usermsg := fmt.Sprintf(format, v...)
	prefix := self.prefix
	msg := fmt.Sprintf("%s:%s:%s:%s\n", strings.Join(fdate, ""), levelNames[level][:1], prefix, usermsg)

	// Log event
	var event = Event{
		time:    now,
		level:   level,
		usermsg: usermsg,
		fmsg:    msg,
	}

	// Writes message to this logger and its ancestors.
	mutex.Lock()
	defer mutex.Unlock()
	self.writeAll(&event)

	// Close all logger writers
	if level == FATAL {
		for _, w := range self.outputs {
			w.Close()
		}
		panic("LOG FATAL")
	}
}

// write message to this logger output and of all of its ancestors.
func (self *Logger) writeAll(event *Event) {

	for _, w := range self.outputs {
		w.Write(event)
		w.Sync()
	}
	if self.parent != nil {
		self.parent.writeAll(event)
	}
}

//
// Functions for the Default Logger
//

func Log(level int, format string, v ...interface{}) {

	Default.Log(level, format, v...)
}

func SetLevel(level int) {

	Default.SetLevel(level)
}

func SetLevelByName(lname string) {

	Default.SetLevelByName(lname)
}

func SetFormat(format int) {

	Default.SetFormat(format)
}

func AddWriter(writer LoggerWriter) {

	Default.AddWriter(writer)
}

func Debug(format string, v ...interface{}) {

	Default.Debug(format, v...)
}

func Info(format string, v ...interface{}) {

	Default.Info(format, v...)
}

func Warn(format string, v ...interface{}) {

	Default.Warn(format, v...)
}

func Error(format string, v ...interface{}) {

	Default.Error(format, v...)
}

func Fatal(format string, v ...interface{}) {

	Default.Fatal(format, v...)
}

// Find finds a logger with the specified path.
func Find(path string) *Logger {

	parts := strings.Split(strings.ToUpper(path), "/")
	level := 0
	var find func([]*Logger) *Logger

	find = func(logs []*Logger) *Logger {

		for _, l := range logs {
			if l.name != parts[level] {
				continue
			}
			if level == len(parts)-1 {
				return l
			}
			level++
			return find(l.children)
		}
		return nil
	}
	return find(rootLoggers)
}
