package redis

import (
	"github.com/minhhoanq/ecommerce/catalog_service/configs"
	"github.com/minhhoanq/ecommerce/common/logger"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	cfg configs.Config
	l   logger.Interface
}

func NewRedis(cfg configs.Config, l logger.Interface) *Redis {
	return &Redis{
		cfg: cfg,
		l:   l,
	}
}

func (r *Redis) Connect() *redis.Client {
	// Connect to Redis
	opts := &redis.Options{
		Addr: r.cfg.RedisAddress,
	}
	client := redis.NewClient(opts)

	r.l.Info("connecting to redis", logger.String("address", r.cfg.RedisAddress))

	return client
}
