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

type PostgresConfig struct {
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Database    string `mapstructure:"database"`
	Healthcheck int    `mapstructure:"healthcheck"`

	MaxOpenConns    int `mapstructure:"max_open_conns"`
	MaxIdleConns    int `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"`
}

func NewDefaultPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Username:        "root",
		Password:        "",
		Host:            "localhost",
		Port:            5432,
		Database:        "db",
		Healthcheck:     30,
		MaxOpenConns:    32,
		MaxIdleConns:    16,
		ConnMaxLifetime: 300,
	}
}

type KafkaConfig struct {
	Enable                   bool       `mapstructure:"enable"`
	Brokers                  []string   `mapstructure:"brokers"`
	ConsumerGroup            string     `mapstructure:"consumer_group"`
	AutoCreateTopics         bool       `mapstructure:"auto_create_topics"`
	DefaultPartitions        int        `mapstructure:"default_partitions"`
	DefaultReplicationFactor int        `mapstructure:"default_replication_factor"`
	Healthcheck              int        `mapstructure:"healthcheck"`
	SASL                     SASLConfig `mapstructure:"sasl"`
}

type SASLConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Mechanism string `mapstructure:"mechanism"`
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
}

func NewDefaultKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Enable:                   false,
		Brokers:                  []string{"localhost:9092"},
		ConsumerGroup:            "default-group",
		AutoCreateTopics:         true,
		DefaultPartitions:        1,
		DefaultReplicationFactor: 1,
		Healthcheck:              30,
		SASL: SASLConfig{
			Enabled:   false,
			Mechanism: "PLAIN",
			Username:  "",
			Password:  "",
		},
	}
}
