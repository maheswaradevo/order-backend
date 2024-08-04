package config

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func NewRedis(config *Config, log *zap.Logger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Address,
		Password: config.RedisConfig.Password,
		DB:       config.RedisConfig.DB,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("could not connect to redis: ", zap.Error(err))
	}

	return rdb
}
