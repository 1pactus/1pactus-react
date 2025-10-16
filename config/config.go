package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/frimin/1pactus-react/log"

	"github.com/a8m/envsubst"
	"github.com/spf13/viper"
)

type AppConfig struct {
	RunMode string `mapstructure:"run_mode"`
}

func NewDefaultAppConfig() *AppConfig {
	return &AppConfig{
		RunMode: "debug",
	}
}

type ConfigBase struct {
	App   *AppConfig   `mapstructure:"app"`
	Log   *log.Options `mapstructure:"log"`
	Mongo *MongoConfig `mapstructure:"mongo"`
	Redis *RedisConfig `mapstructure:"redis"`
}

func NewDefaultConfigBaseConfig() *ConfigBase {
	return &ConfigBase{
		App:   NewDefaultAppConfig(),
		Log:   log.NewDefaultOptions(),
		Mongo: NewDefaultMongoConfig(),
		Redis: NewDefaultRedisConfig(),
	}
}

type logSetter interface {
	setLog(log *log.Options)
}

func (b *ConfigBase) setLog(log *log.Options) {
	b.Log = log
}

func expandEnv(configContent string) (string, error) {
	result, err := envsubst.String(configContent)
	if err != nil {
		panic(err)
	}
	return result, nil
}

func LoadConfig(app string, files []string, cliOverrides []string, config interface{}) error {
	v := viper.New()
	v.SetConfigType("yaml")

	for i, file := range files {
		absfile, err := filepath.Abs(file)
		if err != nil {
			return fmt.Errorf("error getting absolute path of config file: %w", err)
		}
		files[i] = absfile
	}

	for i, file := range files {
		rawConfig, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}

		configString, err := expandEnv(string(rawConfig))

		if err != nil {
			return fmt.Errorf("error expanding env: %w", err)
		}

		if i == 0 {
			if err := v.ReadConfig(bytes.NewBuffer([]byte(configString))); err != nil {
				return fmt.Errorf("error reading config file: %w", err)
			}
		} else {
			if err := v.MergeConfig(bytes.NewBuffer([]byte(configString))); err != nil {
				return fmt.Errorf("error reading config file: %w", err)
			}
		}
	}

	for _, override := range cliOverrides {
		parts := strings.SplitN(override, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			v.Set(key, value)
		} else {
			return fmt.Errorf("invalid override format '%s', skipping", override)
		}
	}

	if s, ok := config.(logSetter); ok {
		var opt log.Options

		opt.App = app
		opt.RunMode = v.GetString("app.run_mode")
		opt.OutPath = v.GetString("log.out_path")
		opt.MaxSize = v.GetInt("log.max_size")
		opt.MaxBackups = v.GetInt("log.max_backups")
		opt.MaxAge = v.GetInt("log.max_age")
		opt.Compress = v.GetBool("log.compress")
		opt.ConsoleLog = v.GetBool("log.console_log")

		s.setLog(&opt)
		log.Setup(opt)
	}

	if _, ok := config.(logSetter); ok {
		var opt log.Options

		opt.App = app
		opt.RunMode = v.GetString("app.run_mode")
		opt.OutPath = v.GetString("record_log.out_path")
		opt.MaxSize = v.GetInt("record_log.max_size")
		opt.MaxBackups = v.GetInt("record_log.max_backups")
		opt.MaxAge = v.GetInt("record_log.max_age")
		opt.Compress = v.GetBool("record_log.compress")
		opt.ConsoleLog = v.GetBool("record_log.console_log")

		log.SetupRecord(opt)
	}

	err := v.Unmarshal(config)

	if err != nil {
		return err
	}

	for _, file := range files {
		log.WithKv("file", file).Infof("load config")
	}

	for _, override := range cliOverrides {
		log.WithKv("override", override).Infof("load config")
	}

	return nil
}
