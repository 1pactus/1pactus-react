package log

import "github.com/rs/zerolog"

type ILogger interface {
	Fatal(format ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})

	WithFields(fields map[string]interface{}) ILogger
	WithField(key string, value interface{}) ILogger
	WithError(err error) ILogger

	WithKv(key, value string) ILogger

	GetInternalLogger() zerolog.Logger
	Skip(upper ...int) ILogger

	Write(p []byte) (n int, err error)
}
