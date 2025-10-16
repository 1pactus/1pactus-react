package service

import (
	"github.com/frimin/1pactus-react/app/onepacd/service/gather"
	"github.com/frimin/1pactus-react/app/onepacd/service/webapi"
)

type Config struct {
	Gather *gather.Config `mapstructure:"gather"`
	WebApi *webapi.Config `mapstructure:"webapi"`
}

func NewDefaultServiceConfig() *Config {
	return &Config{
		Gather: gather.NewDefaultConfig(),
		WebApi: webapi.NewDefaultConfig(),
	}
}
