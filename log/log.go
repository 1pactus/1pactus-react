package log

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

func Raw() *ZeroLog {
	return defaultLog
}

func WithCaller() ILogger {
	return defaultLog.WithCaller()
}

func WithKv(key, value string) ILogger {
	return defaultLog.WithKv(key, value)
}

func Skip(upper ...int) ILogger {
	var skip = 1
	if len(upper) > 0 {
		skip = upper[0] + 1
	}
	// 在这里获取上一级函数的行号
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		index := strings.LastIndex(file, "/")
		if index == -1 {
			return defaultLog.WithKv("file", fmt.Sprintf("%s:%d", file, line))
		} else {
			return defaultLog.WithKv("file", fmt.Sprintf("%s:%d", file[index+1:], line))
		}
	}
	return defaultLog.WithKv("file", "no found file")
}

func WithFields(fields map[string]interface{}) ILogger {
	return defaultLog.WithFields(fields)
}

func Debugf(format string, args ...interface{}) {
	defaultLog.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	defaultLog.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	defaultLog.Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	defaultLog.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLog.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	defaultLog.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	defaultLog.Panicf(format, args...)
}

func Debug(args ...interface{}) {
	defaultLog.Debug(args...)
}

func Info(args ...interface{}) {
	defaultLog.Info(args...)
}

func Warn(args ...interface{}) {
	defaultLog.Warn(args...)
}

func Warning(args ...interface{}) {
	defaultLog.Warning(args...)
}

func Error(args ...interface{}) {
	defaultLog.Error(args...)
}

func Fatal(args ...interface{}) {
	defaultLog.Fatal(args...)
}

func Panic(args ...interface{}) {
	defaultLog.Panic(args...)
}

type RecordObject interface {
	MarshalRecordObject(e *zerolog.Event) string
}

func RecordInfo(recordObj RecordObject) {
	event := recordLog.logger.Info()
	name := recordObj.MarshalRecordObject(event) // 调用我们写的方法，手动展开字段
	event.Msg(name)
}
