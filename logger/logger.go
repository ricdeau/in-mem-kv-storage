package logger

import (
	"fmt"
	"log"
)

const (
	info  = "[INFO]"
	err   = "[ERROR]"
	fatal = "[FATAL]"
)

func Infof(msg string, args ...interface{}) {
	logLevel(info, msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	logLevel(err, msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	log.Fatalf("%s %s", fatal, fmt.Sprintf(msg, args...))
}

func logLevel(level string, msg string, args ...interface{}) {
	log.Printf("%s %s", level, fmt.Sprintf(msg, args...))
}
