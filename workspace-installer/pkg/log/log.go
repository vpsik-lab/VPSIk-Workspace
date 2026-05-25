package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

var levelNames = map[Level]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
}

type Logger struct {
	level Level
	out   io.Writer
	err   io.Writer
}

var Default = New(LevelInfo, os.Stdout, os.Stderr)

func New(level Level, out, err io.Writer) *Logger {
	return &Logger{level: level, out: out, err: err}
}

func SetLevel(level Level) {
	Default.level = level
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05")
	w := l.out
	if level >= LevelWarn {
		w = l.err
	}
	fmt.Fprintf(w, "%s [%s] %s\n", timestamp, levelNames[level], msg)
}

func Debug(format string, args ...interface{}) { Default.log(LevelDebug, format, args...) }
func Info(format string, args ...interface{})  { Default.log(LevelInfo, format, args...) }
func Warn(format string, args ...interface{})  { Default.log(LevelWarn, format, args...) }
func Error(format string, args ...interface{}) { Default.log(LevelError, format, args...) }

func (l *Logger) Debug(format string, args ...interface{}) { l.log(LevelDebug, format, args...) }
func (l *Logger) Info(format string, args ...interface{})  { l.log(LevelInfo, format, args...) }
func (l *Logger) Warn(format string, args ...interface{})  { l.log(LevelWarn, format, args...) }
func (l *Logger) Error(format string, args ...interface{}) { l.log(LevelError, format, args...) }
