package logger

import (
	"errors"
	"fmt"
	"sync"
)

// ErrInvalidLevelName is returned when an invalid level name is provided.
var ErrInvalidLevelName = errors.New("invalid log level name")

// Level defines a type to represent the severity level for a log entry.
type Level int8

// Initialize constants which represent a specific severity level. We use the iota
// keyword as a shortcut to assign successive integer values to the constants.
const (
	DebugLevel Level = iota // Has the value 0.
	InfoLevel
	ErrorLevel
	FatalLevel
	OffLevel
)

// Return a human-friendly string for the severity level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "Debug"
	case InfoLevel:
		return "Info"
	case ErrorLevel:
		return "Error"
	case FatalLevel:
		return "Fatal"
	case OffLevel:
		return "Off"
	default:
		return ""
	}
}

// LevelFromName returns the matched Level. If no match found the return InfoLevel as default Level
// and ErrInvalidLevelName error.
func LevelFromName(name string) (Level, error) {
	switch name {
	case "Debug":
		return DebugLevel, nil
	case "Info":
		return InfoLevel, nil
	case "Error":
		return ErrorLevel, nil
	case "Fatal":
		return FatalLevel, nil
	case "Off":
		return OffLevel, nil
	default:
		return InfoLevel, ErrInvalidLevelName
	}
}

// ExtraData is used to log additional information as a Key Value pair.
type ExtraData struct {
	Key   string
	Value interface{}
}

// Writer interface that can log messages
type Writer interface {
	// Write logs with log level, message and optional additional data as ExtraData
	Write(level Level, message string, data ...ExtraData)
}

// Logger can write logs at or above a minimum severity level
type Logger interface {
	// UpdateMinLevel updates the minimum level of the logger
	UpdateMinLevel(minLevel Level)
	// MinLevel returns the minimum level of the logger
	MinLevel() Level
	// Debug logs a message at DebugLevel
	Debug(message string, data ...ExtraData)
	// Info logs a message at InfoLevel.
	Info(message string, data ...ExtraData)
	// Error logs a message at ErrorLevel.
	Error(message string, data ...ExtraData)
	// Fatal logs a message at FatalLevel.
	Fatal(message string, data ...ExtraData)
}

// logger is a struct with a writer and a minimum severity level for logs
type logger struct {
	minLevel Level
	writer   Writer
}

// singleton instance of logger
var singletonLogger *logger
var lock sync.Mutex

// GetLogger return the logger instance. It creates new default logger(zap) with 'Info' level
// if logger is not created before.
func GetLogger() Logger {
	if singletonLogger == nil {
		fmt.Println("Logger not found. Creating default logger...")
		writer, err := NewZapWriter()
		if err != nil {
			fmt.Println("Default log writer(zap logger) not found")
		}
		NewLogger("Info", writer)
		fmt.Println("Default logger (zap) created with 'info' level")
	}
	return singletonLogger
}

// NewLogger Return a new logger instance which writes log entries at or above a minimum severity
// level to a specific output destination.
// If minLevel doesn't match with valid log level then logger will be configured with InfoLevel
func NewLogger(minLevel string, writer Writer) Logger {
	if singletonLogger == nil {
		lock.Lock()
		defer lock.Unlock()
		if singletonLogger == nil {
			logLevel, _ := LevelFromName(minLevel)
			singletonLogger = &logger{
				minLevel: logLevel,
				writer:   writer,
			}
		}
	} else {
		singletonLogger.Error("Logger is already created")
	}
	return singletonLogger
}

// Extra Add extra key-value pair to a logger's context
func Extra(key string, val interface{}) ExtraData {
	return ExtraData{
		Key:   key,
		Value: val,
	}
}

// UpdateMinLevel updates the minimum level of the logger
func (l *logger) UpdateMinLevel(minLevel Level) {
	l.minLevel = minLevel
}

// MinLevel returns the minimum level of the logger
func (l *logger) MinLevel() Level {
	return l.minLevel
}

// Debug logs a message at DebugLevel.
func (l *logger) Debug(message string, data ...ExtraData) {
	l.log(DebugLevel, message, data...)
}

// Info logs a message at InfoLevel.
func (l *logger) Info(message string, data ...ExtraData) {
	l.log(InfoLevel, message, data...)
}

// Error logs a message at ErrorLevel.
func (l *logger) Error(message string, data ...ExtraData) {
	l.log(ErrorLevel, message, data...)
}

// Fatal logs a message at FatalLevel.
func (l *logger) Fatal(message string, data ...ExtraData) {
	l.log(FatalLevel, message, data...)
}

// log check the current logger level and write in the writer.
func (l *logger) log(level Level, message string, data ...ExtraData) {
	if level < l.minLevel {
		return
	}
	l.writer.Write(level, message, data...)
}
