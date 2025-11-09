package service

import (
	"github.com/1pactus/1pactus-react/app/onepacd/service/chainextract"
	"github.com/1pactus/1pactus-react/app/onepacd/service/chainscan"
	"github.com/1pactus/1pactus-react/app/onepacd/service/webapi"
)

type Config struct {
	Chainscan    *chainscan.Config    `mapstructure:"chainscan"`
	WebApi       *webapi.Config       `mapstructure:"webapi"`
	ChainExtract *chainextract.Config `mapstructure:"chainextract"`
}

func NewDefaultServiceConfig() *Config {
	return &Config{
		Chainscan:    chainscan.NewDefaultConfig(),
		WebApi:       webapi.NewDefaultConfig(),
		ChainExtract: chainextract.NewDefaultConfig(),
	}
}
