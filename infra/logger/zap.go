// infrastructure/zaplogger/zap_logger.go
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger(level zapcore.Level, logFilePath string, errorFilePath string) *ZapLogger {
	cfg := zap.Config{
		Encoding:         "json", // or "console"
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{logFilePath},   // where logs go
		ErrorOutputPaths: []string{errorFilePath}, // where errors go
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	l, err := cfg.Build()
	if err != nil {
		panic(err) // or handle error properly
	}

	return &ZapLogger{logger: l.Sugar()}
}

func (z *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	z.logger.Infow(msg, keysAndValues...)
}

func (z *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	z.logger.Errorw(msg, keysAndValues...)
}

func (z *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	z.logger.Debugw(msg, keysAndValues...)
}

func (z *ZapLogger) Infof(format string, args ...interface{}) {
	z.logger.Infof(format, args...)
}

func (z *ZapLogger) Errorf(format string, args ...interface{}) {
	z.logger.Errorf(format, args...)
}

func (z *ZapLogger) Debugf(format string, args ...interface{}) {
	z.logger.Debugf(format, args...)
}
