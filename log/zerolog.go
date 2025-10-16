package log

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

func Wrapper(logger zerolog.Logger) ILogger {
	return &ZeroLog{logger}
}

type ZeroLog struct {
	logger zerolog.Logger
}

func (z *ZeroLog) Skip(upper ...int) ILogger {
	var skip = 1
	if len(upper) > 0 {
		skip = upper[0] + 1
	}
	// 在这里获取上一级函数的行号
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		index := strings.LastIndex(file, "/")
		if index == -1 {
			return z.WithKv("file", fmt.Sprintf("%s:%d", file, line))
		} else {
			return z.WithKv("file", fmt.Sprintf("%s:%d", file[index+1:], line))
		}
	}
	return z.WithKv("file", "no found file")
}

func (z *ZeroLog) WithCaller() ILogger {
	return Wrapper(z.logger.With().Caller().Logger())
}

func (z *ZeroLog) WithKv(key, value string) ILogger {
	return Wrapper(z.logger.With().Str(key, value).Logger())
}

func (z *ZeroLog) WithKvs(key string, vals []string) ILogger {
	return Wrapper(z.logger.With().Strs(key, vals).Logger())
}

func (z *ZeroLog) WithJson(key string, b []byte) ILogger {
	return Wrapper(z.logger.With().RawJSON(key, b).Logger())
}

func (z *ZeroLog) WithError(err error) ILogger {
	return Wrapper(z.logger.With().Err(err).Logger())
}

func (z *ZeroLog) WithErrors(key string, errs []error) ILogger {
	return Wrapper(z.logger.With().Errs(key, errs).Logger())
}

func (z *ZeroLog) WithField(key string, value interface{}) ILogger {
	return z.WithFields(map[string]interface{}{key: value})
}

func (z *ZeroLog) WithFields(fields map[string]interface{}) ILogger {
	return Wrapper(z.logger.With().Fields(fields).Logger())
}

func (z *ZeroLog) Debugf(format string, args ...interface{}) {
	if zerolog.GlobalLevel() > zerolog.DebugLevel {
		return
	}
	z.logger.Debug().Msgf(format, args...)
}

func (z *ZeroLog) Infof(format string, args ...interface{}) {
	if zerolog.GlobalLevel() > zerolog.InfoLevel {
		return
	}
	z.logger.Info().Msgf(format, args...)
}

func (z *ZeroLog) Warnf(format string, args ...interface{}) {
	if zerolog.GlobalLevel() > zerolog.WarnLevel {
		return
	}
	z.logger.Warn().Msgf(format, args...)
}

func (z *ZeroLog) Warningf(format string, args ...interface{}) {
	if zerolog.GlobalLevel() > zerolog.WarnLevel {
		return
	}
	z.logger.Warn().Msgf(format, args...)
}

func (z *ZeroLog) Errorf(format string, args ...interface{}) {
	if zerolog.GlobalLevel() > zerolog.ErrorLevel {
		return
	}
	z.logger.Error().Msgf(format, args...)
}

func (z *ZeroLog) Fatalf(format string, args ...interface{}) {
	if zerolog.GlobalLevel() > zerolog.FatalLevel {
		return
	}
	z.logger.Panic().Msgf(format, args...)
}

func (z *ZeroLog) Panicf(format string, args ...interface{}) {
	if zerolog.GlobalLevel() > zerolog.PanicLevel {
		return
	}
	z.logger.Panic().Msgf(format, args...)
}

func (z *ZeroLog) Debug(args ...interface{}) {
	if zerolog.GlobalLevel() > zerolog.DebugLevel {
		return
	}
	z.logger.Debug().Msg(fmt.Sprint(args...))
}

func (z *ZeroLog) Info(args ...interface{}) {
	if zerolog.GlobalLevel() > zerolog.InfoLevel {
		return
	}
	z.logger.Info().Msg(fmt.Sprint(args...))
}

func (z *ZeroLog) Warn(args ...interface{}) {
	if zerolog.GlobalLevel() > zerolog.WarnLevel {
		return
	}
	z.logger.Warn().Msg(fmt.Sprint(args...))
}

func (z *ZeroLog) Warning(args ...interface{}) {
	if zerolog.GlobalLevel() > zerolog.WarnLevel {
		return
	}
	z.logger.Warn().Msg(fmt.Sprint(args...))
}

func (z *ZeroLog) Error(args ...interface{}) {
	if zerolog.GlobalLevel() > zerolog.ErrorLevel {
		return
	}
	z.logger.Error().Msg(fmt.Sprint(args...))
}

func (z *ZeroLog) Fatal(args ...interface{}) {
	if zerolog.GlobalLevel() > zerolog.FatalLevel {
		return
	}
	z.logger.Panic().Msg(fmt.Sprint(args...))
}

func (z *ZeroLog) Panic(args ...interface{}) {
	if zerolog.GlobalLevel() > zerolog.PanicLevel {
		return
	}
	z.logger.Panic().Msg(fmt.Sprint(args...))
}

func (z *ZeroLog) Debugln(args ...interface{}) {
	z.Debug(args...)
}

func (z *ZeroLog) Infoln(args ...interface{}) {
	z.Info(args...)
}

func (z *ZeroLog) Warnln(args ...interface{}) {
	z.Warn(args...)
}

func (z *ZeroLog) Errorln(args ...interface{}) {
	z.Error(args...)
}

func (z *ZeroLog) Fatalln(args ...interface{}) {
	z.Panic(args...)
}

func (z *ZeroLog) Panicln(args ...interface{}) {
	z.Panic(args...)
}

func (z *ZeroLog) GetInternalLogger() zerolog.Logger {
	return z.logger
}

func (w *ZeroLog) Write(p []byte) (n int, err error) {
	message := string(p)
	if len(message) > 0 && message[len(message)-1] == '\n' {
		message = message[:len(message)-1]
	}

	w.logger.Info().Msg(message)
	return len(p), nil
}
