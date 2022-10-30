package db

import "github.com/ErfanMomeniii/Magic-Load-Balancer/internal/repository"

type DB interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

func Start() {
	//we can add config for selection from dbs but for now I don't want another db,
	//so that we add choose it directly(without switch case)
	rc := &RedisConnection{}

	rc.Client = *NewRedisServer()

	repository.SelectHandler = &repository.ServerSelectionHandler{
		DB: rc,
	}
}
