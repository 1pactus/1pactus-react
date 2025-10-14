package config

type MongoConfig struct {
	Uri         string `mapstructure:"uri"`
	Database    string `mapstructure:"database"`
	Password    string `mapstructure:"password"`
	Healthcheck int    `mapstructure:"healthcheck"`
}

func NewDefaultMongoConfig() *MongoConfig {
	return &MongoConfig{
		Uri:         "mongodb://localhost:27017",
		Database:    "db",
		Password:    "",
		Healthcheck: 30,
	}
}

type RedisConfig struct {
	Addr         string   `mapstructure:"addr"`
	ClusterAddrs []string `mapstructure:"cluster_addrs"`
	Password     string   `mapstructure:"password"`
	Db           int      `mapstructure:"db"`
	Healthcheck  int      `mapstructure:"healthcheck"`
}

func NewDefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:        "localhost:6379",
		Password:    "",
		Healthcheck: 30,
		Db:          0,
	}
}
