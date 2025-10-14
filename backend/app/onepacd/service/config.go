package service

import "github.com/frimin/1pactus-react/backend/app/onepacd/service/gather"

type Config struct {
	Gather *gather.Config `mapstructure:"gather"`
}

func NewDefaultServiceConfig() *Config {
	return &Config{
		Gather: gather.NewDefaultConfig(),
	}
}
