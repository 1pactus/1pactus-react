package storedriver

import (
	"context"
	"fmt"

	"github.com/frimin/1pactus-react/backend/config"
	"github.com/frimin/1pactus-react/backend/log"

	"time"

	"github.com/redis/go-redis/v9"
)

type IRedisStore interface {
	Init(store Redis)
}

type Redis interface {
	GetClient() redis.UniversalClient
	GetTimeout() time.Duration
}

type redisImpl struct {
	conf    *config.RedisConfig
	client  redis.UniversalClient
	stores  []IRedisStore
	timeout time.Duration
	log     log.ILogger
}

func RedisStart(name string, conf *config.RedisConfig, stores []IRedisStore) error {
	if len(conf.ClusterAddrs) > 0 && len(conf.Addr) != 0 {
		return fmt.Errorf("redis config error: both ClusterAddrs and Addr are set, please use one of them")
	} else if len(conf.ClusterAddrs) == 0 && len(conf.Addr) == 0 {
		return fmt.Errorf("redis config error: neither ClusterAddrs nor Addr is set, please set one of them")
	}

	isCluserMode := len(conf.ClusterAddrs) > 0

	r := &redisImpl{
		conf:    conf,
		stores:  stores,
		timeout: time.Second * 10,
	}

	if isCluserMode {
		r.log = log.WithKv("module", "store").WithKv("redis_cluster", name)
	} else {
		r.log = log.WithKv("module", "store").WithKv("redis", name)
	}

	var err error

	maxRetry := 10

	for {
		if err = r.connect(); err == nil {
			r.initStores()

			r.log.Infof("redis connect and initialized success")

			go r.monitorConnection()
			return nil
		} else {
			maxRetry -= 1

			if maxRetry <= 0 {
				return fmt.Errorf("redis [%s] connect failed after retries: %v", name, err)
			}

			r.log.Errorf("redis [%s] connect failed: %v", name, err)
			time.Sleep(time.Second * 5)
		}
	}
}

func (r *redisImpl) connect() error {
	var rdb redis.UniversalClient

	if len(r.conf.ClusterAddrs) > 0 {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    r.conf.ClusterAddrs,
			Password: r.conf.Password,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:     r.conf.Addr,
			Password: r.conf.Password,
			DB:       0,
		})
	}

	r.client = rdb

	_, err := r.client.Ping(context.Background()).Result()
	if err != nil {
		return fmt.Errorf("redis connect error: %v", err)
	}
	return nil
}

func (r *redisImpl) initStores() {
	for _, store := range r.stores {
		store.Init(r)
	}
}

func (r *redisImpl) GetClient() redis.UniversalClient {
	return r.client
}

func (r *redisImpl) GetTimeout() time.Duration {
	return r.timeout
}

func (r *redisImpl) monitorConnection() {
	healthcheck := time.Duration(r.conf.Healthcheck)

	for {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

			defer cancel()

			if _, err := r.client.Ping(ctx).Result(); err != nil {
				r.log.Errorf("lost redis connection, retrying: %v", err)

				for {
					if err := r.connect(); err == nil {
						r.initStores()
						r.log.Info("successfully reconnected to redis")
						break
					}
					r.log.Errorf("error connecting to redis: %v", err)
					time.Sleep(5 * time.Second)
				}
			}

		}()

		time.Sleep(healthcheck * time.Second)
	}
}
