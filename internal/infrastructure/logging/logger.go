package logging

import (
	"log"
	"os"
	"sync"
)

// Logger defines the interface for logging operations.
type Logger interface {
	Info(format string, v ...interface{})
	Error(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

// logger implements the Logger interface.
type logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

var (
	once     sync.Once
	instance *logger
)

// New creates and returns a new Logger instance.
func New() Logger {
	once.Do(func() {
		instance = &logger{
			infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
			errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
			debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile), // Can be changed to os.Stderr or disabled
		}
	})
	return instance
}

// Info logs messages at the info level.
func (l *logger) Info(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

// Error logs messages at the error level.
func (l *logger) Error(format string, v ...interface{}) {
	l.errorLogger.Printf(format, v...)
}

// Debug logs messages at the debug level.
func (l *logger) Debug(format string, v ...interface{}) {
	l.debugLogger.Printf(format, v...)
}
