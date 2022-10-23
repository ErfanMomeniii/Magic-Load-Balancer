package handler

import (
	"fmt"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/app"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/config"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type LoadBalancerActivity interface {
	SendToServers(ctx echo.Context, service config.Service, path string) error
	FindSuitableServer(service config.Service) (config.Server, error)
}

func SendToServers(ctx echo.Context, endpoint config.Endpoint) error {
	if config.C.Tracing.Enabled {
		_, span := app.Tracer.Start(ctx.Request().Context(), "SendTOServers")
		defer span.End()
	}

	server, err := FindSuitableServer(endpoint.Service)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"result": "we have problem on choosing server for this request",
		})
	}

	url := fmt.Sprintf("%s/%s", server.IP, endpoint.URL)

	req, err := http.NewRequest(endpoint.HttpMethod, url, ctx.Request().Body)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"result": "we have problem on sending your request",
		})
	}

	req.Header.Set("Content-Type", ctx.Request().Header.Get("Content-Type"))

	client := &http.Client{}

	resp, err := client.Do(req)

	if err == nil {
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		return ctx.JSON(resp.StatusCode, body)
	}
	return nil
}

func FindSuitableServer(service config.Service) (config.Server, error) {
	switch config.C.Algorithm.Name {
	case "random":
		index := rand.Int63n(int64(len(service.Servers)))
		return service.Servers[index], nil
	case "round-robin":
	case "magic":
	}
	return config.Server{}, nil
}
