package onepacd

import (
	"github.com/frimin/1pactus-react/app/onepacd/service"
	"github.com/frimin/1pactus-react/config"
)

type Config struct {
	*config.ConfigBase `mapstructure:",squash"`
	Service            *service.Config `mapstructure:"service"`
}

var conf = Config{
	ConfigBase: config.NewDefaultConfigBaseConfig(),
	Service:    service.NewDefaultServiceConfig(),
}

func LoadConfig(app string, files []string, cliOverrides []string) (err error) {
	err = config.LoadConfig(app, files, cliOverrides, &conf)
	return
}
