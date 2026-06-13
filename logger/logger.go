package logger

import (
	"io"
	"log"
	"sync/atomic"
)

type LogLevel int

const (
	DEBUG   LogLevel = iota // 0
	INFO                    // 1
	WARNING                 // 2
	ERROR                   // 3
	NONE                    // 4
)

func (l LogLevel) String() string {
	s := ""
	switch l {
	case DEBUG:
		s = "Debug"
	case INFO:
		s = "Info"
	case WARNING:
		s = "Warning"
	case ERROR:
		s = "Error"
	case NONE:
		s = "None"
	}
	return s
}

type Logger struct {
	l *log.Logger

	level int32
}

var defaultLogger = &Logger{
	l:     log.Default(),
	level: int32(DEBUG),
}

func SetOutput(w io.Writer) {
	defaultLogger.SetOutput(w)
}

func SetLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

func Fatal(v ...any) {
	defaultLogger.Fatal(v...)
}
func Error(v ...any) {
	defaultLogger.Error(v...)
}
func Warning(v ...any) {
	defaultLogger.Warning(v...)
}
func Info(v ...any) {
	defaultLogger.Info(v...)
}
func Debug(v ...any) {
	defaultLogger.Debug(v...)
}

func Fatalf(format string, v ...any) {
	defaultLogger.Fatalf(format, v...)
}
func Errorf(format string, v ...any) {
	defaultLogger.Errorf(format, v...)
}
func Warningf(format string, v ...any) {
	defaultLogger.Warningf(format, v...)
}
func Infof(format string, v ...any) {
	defaultLogger.Infof(format, v...)
}
func Debugf(format string, v ...any) {
	defaultLogger.Debugf(format, v...)
}

func New(out io.Writer, prefix string, flag int, level LogLevel) *Logger {
	return &Logger{
		l:     log.New(out, prefix, flag),
		level: int32(level),
	}
}
func (l *Logger) isEligible(op LogLevel) bool {
	level := atomic.LoadInt32(&l.level)
	return int32(op) >= level
}

func (l *Logger) SetOutput(w io.Writer) {
	l.l.SetOutput(w)
}

func (l *Logger) SetLevel(level LogLevel) {
	atomic.StoreInt32(&l.level, int32(level))
}

func (l *Logger) Fatal(v ...any) {
	args := make([]any, len(v)+1)
	args[0] = "[FATAL] "
	copy(args[1:], v)

	l.l.Fatal(args...)
}
func (l *Logger) Error(v ...any) {
	if l.isEligible(ERROR) {
		args := make([]any, len(v)+1)
		args[0] = "[ERROR] "
		copy(args[1:], v)
		l.l.Print(args...)
	}
}
func (l *Logger) Warning(v ...any) {
	if l.isEligible(WARNING) {
		args := make([]any, len(v)+1)
		args[0] = "[WARNING] "
		copy(args[1:], v)
		l.l.Print(args...)
	}
}
func (l *Logger) Info(v ...any) {
	if l.isEligible(INFO) {
		args := make([]any, len(v)+1)
		args[0] = "[INFO] "
		copy(args[1:], v)
		l.l.Print(args...)
	}
}
func (l *Logger) Debug(v ...any) {
	if l.isEligible(DEBUG) {
		args := make([]any, len(v)+1)
		args[0] = "[DEBUG] "
		copy(args[1:], v)
		l.l.Print(args...)
	}
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.l.Fatalf("[FATAL] "+format, v...)
}
func (l *Logger) Errorf(format string, v ...any) {
	if l.isEligible(ERROR) {
		l.l.Printf("[ERROR] "+format, v...)
	}
}
func (l *Logger) Warningf(format string, v ...any) {
	if l.isEligible(WARNING) {
		l.l.Printf("[WARNING] "+format, v...)
	}
}
func (l *Logger) Infof(format string, v ...any) {
	if l.isEligible(INFO) {
		l.l.Printf("[INFO] "+format, v...)
	}
}
func (l *Logger) Debugf(format string, v ...any) {
	if l.isEligible(DEBUG) {
		l.l.Printf("[DEBUG] "+format, v...)
	}
}
