package repository

import (
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/config"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/db"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/log"
	"github.com/go-redis/redis/v8"
	"math"
)

var SelectHandler *ServerSelectionHandler

type ServerSelection interface {
	SelectServerRoundly(service config.Service) config.Server
	SelectServerMagically(service config.Service) config.Server
}

type ServerSelectionHandler struct {
	DB db.DB
}

func (handler *ServerSelectionHandler) SelectServerRoundly(service config.Service) config.Server {
	index, err := handler.DB.Get(service.Name)

	if err == redis.Nil {
		index = 0
	}

	err = handler.DB.Set(service.Name, (index.(int)+1)%len(service.Servers))
	if err != nil {
		log.Logger.Error(err.Error())
	}

	return service.Servers[index.(int)]
}

func (handler *ServerSelectionHandler) SelectServerMagically(service config.Service) config.Server {
	index := 0
	minTime := math.MaxInt
	for i, server := range service.Servers {
		serverWorkingTime, err := handler.DB.Get(server.IP)

		if err == redis.Nil {
			serverWorkingTime = 0
			_ = handler.DB.Set(server.IP, 0)
		}

		if serverWorkingTime.(int) < minTime {
			minTime = serverWorkingTime.(int)
			index = i
		}
	}

	return service.Servers[index]
}
