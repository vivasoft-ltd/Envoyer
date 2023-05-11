package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapWriter is a struct with zap.Logger and it implements the Writer interface
type ZapWriter struct {
	logger *zap.Logger
}

// NewZapWriter creates a new ZapWriter that implements the Writer interface
func NewZapWriter() (ZapWriter, error) {
	logConfig := zap.Config{
		OutputPaths: []string{"stdout"},
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error
	var log *zap.Logger
	if log, err = logConfig.Build(); err != nil {
		return ZapWriter{}, err
	}

	return ZapWriter{
		logger: log,
	}, nil
}

// Write logs the message and optional data to the zap logger
func (zw ZapWriter) Write(level Level, message string, data ...ExtraData) {
	defer zw.logger.Sync()
	var fields []zapcore.Field
	for _, contextData := range data {
		fields = append(fields, zap.Any(contextData.Key, contextData.Value))
	}

	switch level {
	case InfoLevel:
		zw.logger.Info(message, fields...)
	case ErrorLevel:
		zw.logger.Error(message, fields...)
	case FatalLevel:
		zw.logger.Fatal(message, fields...)
	default:
		zw.logger.Debug(message, fields...)
	}
}
