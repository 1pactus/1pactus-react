package chainextract

type Config struct {
	GrpcServers []string `mapstructure:"grpc_servers"`
}

func NewDefaultConfig() *Config {
	return &Config{
		GrpcServers: []string{"127.0.0.1:50051"},
	}
}
