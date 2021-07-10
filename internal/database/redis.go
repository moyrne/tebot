package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var Redis *redis.Client

func ConnectRedis() error {
	Redis = redis.NewClient(&redis.Options{
		Addr: viper.GetString("Redis.Host"),
	})

	return Redis.Ping(context.Background()).Err()
}
