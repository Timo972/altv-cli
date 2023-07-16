package logging

import (
	"io"
	"log"
	"os"
)

var (
	InfoLogger   *log.Logger
	WarnLogger   *log.Logger
	ErrLogger    *log.Logger
	DebugLogger  *log.Logger
	defaultFlags = log.Ldate | log.Ltime
	debugFlags   = log.Ldate | log.Ltime | log.Lshortfile
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", defaultFlags)
	WarnLogger = log.New(os.Stdout, "WARN: ", defaultFlags)
	ErrLogger = log.New(os.Stderr, "ERROR: ", defaultFlags)
	DebugLogger = log.New(io.Discard, "DEBUG: ", debugFlags)
}

func SetDebug(debug bool) {
	if debug {
		DebugLogger.SetOutput(os.Stdout)
		InfoLogger.SetFlags(debugFlags)
		WarnLogger.SetFlags(debugFlags)
		ErrLogger.SetFlags(debugFlags)
	} else {
		DebugLogger.SetOutput(io.Discard)
		InfoLogger.SetFlags(defaultFlags)
		WarnLogger.SetFlags(defaultFlags)
		ErrLogger.SetFlags(defaultFlags)
	}
}

func Disable() {
	InfoLogger.SetOutput(io.Discard)
	WarnLogger.SetOutput(io.Discard)
	ErrLogger.SetOutput(io.Discard)
	DebugLogger.SetOutput(io.Discard)
}
