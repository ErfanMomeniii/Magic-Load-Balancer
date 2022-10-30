package db

import (
	"context"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/config"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisConnection struct {
	Client redis.Client
}

func NewRedisServer() *redis.Client {
	if config.C.Redis.Enable {
		rdb := redis.NewClient(&redis.Options{
			Addr:     config.C.Redis.Address,
			Password: config.C.Redis.Password,
			DB:       config.C.Redis.DB,
		})

		return rdb
	}

	return nil
}

type DB interface {
	Set(key string, value string) error
	Get(key string) (string, error)
}

func (rc *RedisConnection) Set(key string, value string) error {
	err := rc.Client.Set(context.Background(), key, value, 365*1000*time.Hour).Err()

	return err
}

func (rc *RedisConnection) Get(key string) (string, error) {
	val, err := rc.Client.Get(context.Background(), key).Result()

	return val, err
}