package http

import (
	"context"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/app"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/config"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/handler"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type server struct {
	e *echo.Echo
}

type LoadBalancerServer interface {
	Serve()
	AddEndpointFromConfig()
}

func NewServer() *server {
	e := echo.New()

	e.HideBanner = true
	e.Server.ReadTimeout = config.C.HTTPServer.ReadTimeout
	e.Server.WriteTimeout = config.C.HTTPServer.WriteTimeout
	e.Server.ReadHeaderTimeout = config.C.HTTPServer.ReadHeaderTimeout
	e.Server.IdleTimeout = config.C.HTTPServer.IdleTimeout

	return &server{
		e: e,
	}
}

func (s *server) Serve() {
	s.AddEndpointFromConfig()

	go func() {
		if err := s.e.Start(config.C.HTTPServer.Listen); err != nil && err != http.ErrServerClosed {
			s.e.Logger.Fatalf("shutting down the server (%v). err: %v", config.C.HTTPServer.Listen, err)
		}
	}()

	go func() {
		<-app.A.Ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := s.e.Shutdown(ctx); err != nil {
			s.e.Logger.Fatal(err)
		}
	}()
}

func (s *server) AddEndpointFromConfig() {
	for _, endpoint := range config.C.Endpoints {
		s.e.Add(endpoint.HttpMethod, endpoint.URL, func(ctx echo.Context) error {
			return handler.SendToServers(ctx, endpoint)
		})
	}
}
