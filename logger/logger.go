package logger

import (
	"fmt"
	"log"
)

// Log levels
const (
	info  = "[INFO]"
	err   = "[ERROR]"
	fatal = "[FATAL]"
)

// Infof prints to standard logger with INFO log level.
func Infof(msg string, args ...interface{}) {
	logLevel(info, msg, args...)
}

// Errorf prints to standard logger with ERROR log level.
func Errorf(msg string, args ...interface{}) {
	logLevel(err, msg, args...)
}

// Fatalf prints to standard logger with FATAL log level.
// Behaves as log.Fatalf.
func Fatalf(msg string, args ...interface{}) {
	log.Fatalf("%s %s", fatal, fmt.Sprintf(msg, args...))
}

func logLevel(level string, msg string, args ...interface{}) {
	log.Printf("%s %s", level, fmt.Sprintf(msg, args...))
}
