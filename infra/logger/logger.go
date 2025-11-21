package logger

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Infof(msg string, keysAndValues ...interface{})
	Errorf(msg string, keysAndValues ...interface{})
}
