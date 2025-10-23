package service

import (
	"github.com/1pactus/1pactus-react/app/onepacd/service/chainextract"
	"github.com/1pactus/1pactus-react/app/onepacd/service/gather"
	"github.com/1pactus/1pactus-react/app/onepacd/service/webapi"
)

type Config struct {
	Gather       *gather.Config       `mapstructure:"gather"`
	WebApi       *webapi.Config       `mapstructure:"webapi"`
	ChainExtract *chainextract.Config `mapstructure:"chainextract"`
}

func NewDefaultServiceConfig() *Config {
	return &Config{
		Gather:       gather.NewDefaultConfig(),
		WebApi:       webapi.NewDefaultConfig(),
		ChainExtract: chainextract.NewDefaultConfig(),
	}
}
