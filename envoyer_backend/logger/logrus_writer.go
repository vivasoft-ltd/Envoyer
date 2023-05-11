package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

// LogrusWriter is a struct with logrus.Logger and it implements the Writer interface
type LogrusWriter struct {
	logger *logrus.Logger
}

// NewLogrusWriter creates a new LogrusWriter that implements the Writer interface.
// Also Hooks can be added to the logger if necessary.
func NewLogrusWriter(hooks ...logrus.Hook) (LogrusWriter, error) {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	customFormatter.ForceColors = true

	log := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{}, //customFormatter,  //&logrus.TextFormatter{},
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.DebugLevel,
	}

	for _, hook := range hooks {
		log.Hooks.Add(hook)
	}

	return LogrusWriter{
		logger: log,
	}, nil
}

// Write logs the message and optional data to the logrus logger
func (lw LogrusWriter) Write(level Level, message string, data ...ExtraData) {
	fields := logrus.Fields{}
	for _, contextData := range data {
		fields[contextData.Key] = contextData.Value
	}

	switch level {
	case InfoLevel:
		lw.logger.WithFields(fields).Info(message)
	case ErrorLevel:
		lw.logger.WithFields(fields).Error(message)
	case FatalLevel:
		lw.logger.WithFields(fields).Fatal(message)
	default:
		lw.logger.WithFields(fields).Debug(message)
	}
}
