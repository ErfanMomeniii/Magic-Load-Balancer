package repository

import (
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/config"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/db"
)

type ServerWorkingDetail interface {
	AddWorkingTime(timeSecond int64) error
	GetWorkingTime() (int64, error)
	SetWorkingTime(timeSecond int64) error
}

type ServerWorkingTimeHandler struct {
	DB     db.DB
	Server config.Server
}

func (serverHandler *ServerWorkingTimeHandler) AddWorkingTime(timeSecond int64) error {
	time, _ := serverHandler.GetWorkingTime()
	time += timeSecond

	return serverHandler.SetWorkingTime(time)
}

func (serverHandler *ServerWorkingTimeHandler) GetWorkingTime() (int64, error) {
	result, err := serverHandler.DB.Get(serverHandler.Server.IP)

	return result.(int64), err
}

func (serverHandler *ServerWorkingTimeHandler) SetWorkingTime(timeSecond int64) error {
	return serverHandler.DB.Set(serverHandler.Server.IP, timeSecond)
}
