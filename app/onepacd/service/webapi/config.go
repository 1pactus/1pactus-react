package webapi

type Config struct {
	HttpListen string `mapstructure:"http_listen"`
}

func NewDefaultConfig() *Config {
	return &Config{
		HttpListen: ":8080",
	}
}
