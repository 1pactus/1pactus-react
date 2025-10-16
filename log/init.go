package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var defaultLog = &ZeroLog{zerolog.New(os.Stdout).With().Timestamp().Logger()}
var recordLog = &ZeroLog{zerolog.New(os.Stdout).With().Timestamp().Logger()}

type Options struct {
	App        string `yaml:"app"`
	RunMode    string `yaml:"run_mode"`
	OutPath    string `yaml:"out_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
	ConsoleLog bool   `yaml:"console_log"`
}

func NewDefaultOptions() *Options {
	return &Options{
		App:        "",
		RunMode:    "debug",
		OutPath:    "./log",
		MaxSize:    1024,
		MaxBackups: 30,
		MaxAge:     14,
		Compress:   false,
		ConsoleLog: true,
	}
}

func Setup(opt Options) {
	defaultLog = &ZeroLog{NewLog(opt)}
}

func SetupRecord(opt Options) {
	recordLog = &ZeroLog{NewLog(opt)}
}

func NewLog(opt Options) zerolog.Logger {
	if !strings.HasSuffix(opt.OutPath, "/") {
		opt.OutPath += "/"
	}
	var infoWriter = &lumberjack.Logger{
		Filename:   opt.OutPath + fmt.Sprintf("%s_info.log", opt.App),
		MaxSize:    opt.MaxSize,
		MaxBackups: opt.MaxBackups,
		MaxAge:     opt.MaxAge,
		Compress:   opt.Compress,
	}
	var warnWriter = &lumberjack.Logger{
		Filename:   opt.OutPath + fmt.Sprintf("%s_warn.log", opt.App),
		MaxSize:    opt.MaxSize,
		MaxBackups: opt.MaxBackups,
		MaxAge:     opt.MaxAge,
		Compress:   opt.Compress,
	}
	var errorWriter = &lumberjack.Logger{
		Filename:   opt.OutPath + fmt.Sprintf("%s_error.log", opt.App),
		MaxSize:    opt.MaxSize,
		MaxBackups: opt.MaxBackups,
		MaxAge:     opt.MaxAge,
		Compress:   opt.Compress,
	}
	// set log level
	if opt.RunMode == "debug" {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	// 使用 zeroWriter 将不同级别的日志写入不同的输出
	var consoleWriter io.Writer
	if opt.ConsoleLog {
		consoleWriter = zerolog.NewConsoleWriter()
	}
	log.Logger = zerolog.New(newZeroWriter(consoleWriter, infoWriter, warnWriter, errorWriter)).
		With().Timestamp().Str("app", opt.App).Logger()
	return log.Logger
}
