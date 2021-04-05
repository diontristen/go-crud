package util

import (
	"log"
)

// USAGE (Error): Print plain error without format
func Error(args ...interface{}) {
	log.Println(args...)
}

// USAGE (Error): Print error with format
func Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// USAGE (Info): Print info without format
func Info(args ...interface{}) {
	log.Println(args...)
}

// USAGE (Info): Print info with format
func Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// VARIABLE: Used for checking if there is a debug flag.
var DebugFlag bool

// USAGE (DEBUG) : Checks if there is debug flag and print the logs
func Debug(args ...interface{}) {
	if !DebugFlag {
		return
	}
	log.Println(args...)
}

// USAGE (DEBUG) : Checks if there is debug flag and print the logs with format
func Debugf(format string, args ...interface{}) {
	if !DebugFlag {
		return
	}
	log.Printf(format, args...)
}
