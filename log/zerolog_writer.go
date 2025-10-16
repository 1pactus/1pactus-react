package log

import (
	"io"

	"github.com/rs/zerolog"
)

func newZeroWriter(console, info, warn, err io.Writer) zerolog.LevelWriter {
	return &zeroWriter{
		console: console,
		info:    info,
		warn:    warn,
		err:     err,
	}
}

type zeroWriter struct {
	console io.Writer
	info    io.Writer
	warn    io.Writer
	err     io.Writer
}

func (zw zeroWriter) Write(p []byte) (n int, err error) {
	if zw.console != nil {
		zw.console.Write(p)
	}
	return zw.info.Write(p)
}

func (zw zeroWriter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	if zw.console != nil {
		zw.console.Write(p)
	}
	switch l {
	case zerolog.TraceLevel:
		fallthrough
	case zerolog.DebugLevel:
		fallthrough
	case zerolog.InfoLevel:
		return zw.info.Write(p)
	case zerolog.WarnLevel:
		return zw.warn.Write(p)
	case zerolog.ErrorLevel:
		fallthrough
	case zerolog.FatalLevel:
		fallthrough
	case zerolog.PanicLevel:
		return zw.err.Write(p)
	default:
		return zw.info.Write(p)
	}
}
