package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Level int

const (
	LevelNone Level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

var (
	level       = LevelInfo
	errorLogger = log.New(os.Stderr, "", log.LstdFlags)
	stdLogger   = log.New(os.Stdout, "", log.LstdFlags)
)

func ParseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "none", "off", "silent":
		return LevelNone
	case "error", "err":
		return LevelError
	case "warn", "warning":
		return LevelWarn
	case "info":
		return LevelInfo
	case "debug":
		return LevelDebug
	default:
		return LevelInfo
	}
}

func (l Level) String() string {
	switch l {
	case LevelNone:
		return "none"
	case LevelError:
		return "error"
	case LevelWarn:
		return "warn"
	case LevelInfo:
		return "info"
	case LevelDebug:
		return "debug"
	default:
		return "unknown"
	}
}

func SetLevel(l Level) {
	level = l
	if l == LevelNone {
		stdLogger.SetOutput(io.Discard)
		errorLogger.SetOutput(io.Discard)
	} else {
		stdLogger.SetOutput(os.Stdout)
		errorLogger.SetOutput(os.Stderr)
	}
}

func Debug(format string, v ...interface{}) {
	if level >= LevelDebug {
		stdLogger.Printf("[DEBUG] "+format, v...)
	}
}

func Info(format string, v ...interface{}) {
	if level >= LevelInfo {
		stdLogger.Printf(format, v...)
	}
}

func Warn(format string, v ...interface{}) {
	if level >= LevelWarn {
		stdLogger.Printf("[WARN] "+format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if level >= LevelError {
		errorLogger.Printf("[ERROR] "+format, v...)
	}
}

func Fatal(format string, v ...interface{}) {
	log.Fatalf("[FATAL] "+format, v...)
}

func Startup(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}

func Result(format string, v ...interface{}) {
	if level >= LevelInfo {
		stdLogger.Printf(format, v...)
	}
}
