// infrastructure/zaplogger/zap_logger.go
package logger

import "go.uber.org/zap"

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger() *ZapLogger {
	l, _ := zap.NewProduction()
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
