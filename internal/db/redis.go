package db

import (
	"context"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/config"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/log"
	"github.com/go-redis/redis/v8"
	"math"
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
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

func (rc *RedisConnection) Set(key string, value interface{}) error {
	err := rc.Client.Set(context.Background(), key, value, 365*1000*time.Hour).Err()

	return err
}

func (rc *RedisConnection) Get(key string) (interface{}, error) {
	val, err := rc.Client.Get(context.Background(), key).Result()

	return val, err
}

func (rc *RedisConnection) ChooseServerRoundly(service config.Service) config.Server {
	index, err := rc.Get(service.Name)

	if err == redis.Nil {
		index = 0
	}

	err = rc.Set(service.Name, (index.(int)+1)%len(service.Servers))
	if err != nil {
		log.Logger.Error(err.Error())
	}

	return service.Servers[index.(int)]
}

func (rc *RedisConnection) ChooseServerMagically(service config.Service) config.Server {
	index := 0
	minTime := math.MaxInt
	for i, server := range service.Servers {
		serverWorkingTime, err := rc.Get(server.IP)

		if err == redis.Nil {
			serverWorkingTime = 0
			_ = rc.Set(server.IP, 0)
		}

		if serverWorkingTime.(int) < minTime {
			minTime = serverWorkingTime.(int)
			index = i
		}
	}

	return service.Servers[index]
}
