package migrator

import (
	"log"

	"github.com/sirupsen/logrus"
)

type LogLevel int

const (
	Debug LogLevel = iota + 1
	Info
	Warn
	Error
)

// Logger defines the logging interface for migration operations
//
// This interface provides leveled logging capabilities with printf-style formatting.
// Implementations should handle log output to various destinations (console, file, etc.).
//
// Log levels:
//
//	Debugf: Detailed diagnostic information (typically for development)
//	Infof:  General operational messages about migration progress
//	Warnf:  Potentially problematic situations that aren't errors
//	Errorf: Critical issues that require attention
//
// Implementations should ensure:
//   - Thread safety for concurrent operations
//   - Proper handling of formatting verbs (%s, %d, etc.)
//   - Reasonable performance for frequent logging
//
// Example implementation:
//
//	type ConsoleLogger struct{}
//	func (l *ConsoleLogger) Debugf(format string, v ...interface{}) {
//	    fmt.Printf("[DEBUG] "+format, v...)
//	}
//	// ... implement other methods similarly
type Logger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

type DefaultLogger struct {
	level LogLevel
}

func NewDefaultLogger(level LogLevel) Logger {
	return &DefaultLogger{level: level}
}

func (l *DefaultLogger) Debugf(format string, v ...interface{}) {
	if l.level <= Debug {
		log.Printf("[DEBUG] "+format, v...)
	}
}

func (l *DefaultLogger) Infof(format string, v ...interface{}) {
	if l.level <= Info {
		log.Printf("[INFO] "+format, v...)
	}
}

func (l *DefaultLogger) Warnf(format string, v ...interface{}) {
	if l.level <= Warn {
		log.Printf("[WARN] "+format, v...)
	}
}

func (l *DefaultLogger) Errorf(format string, v ...interface{}) {
	if l.level <= Error {
		log.Printf("[ERROR] "+format, v...)
	}
}

type Klog struct {
	klog *logrus.Entry
}

func NewKlogLogger(log *logrus.Entry) Logger {
	return &Klog{klog: log}
}

func (l *Klog) Debugf(format string, v ...interface{}) {
	l.klog.Debugf(format, v...)
}

func (l *Klog) Infof(format string, v ...interface{}) {
	l.klog.Infof(format, v...)
}

func (l *Klog) Warnf(format string, v ...interface{}) {
	l.klog.Warnf(format, v...)
}

func (l *Klog) Errorf(format string, v ...interface{}) {
	l.klog.Errorf(format, v...)
}
