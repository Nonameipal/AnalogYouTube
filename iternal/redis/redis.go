package redis

import (
	"fmt"

	"github.com/Nonameipal/P2P/internal/configs"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedisConnection() error {
	cfg := configs.AppSettings

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisParams.Host, cfg.RedisParams.Port),
		Password: cfg.RedisParams.Password,
		DB:       cfg.RedisParams.DB,
	})

	return nil
}

func GetRedisClient() *redis.Client {
	return rdb
}

func CloseRedisConnection() error {
	if rdb != nil {
		return rdb.Close()
	}
	return nil
}
