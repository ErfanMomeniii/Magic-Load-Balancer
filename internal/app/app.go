package app

import (
	"context"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/log"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

const Name = "Load-Balancer"

var (
	GitTag string
	Tracer trace.Tracer
)

type application struct {
	Ctx        context.Context
	cancelFunc context.CancelFunc
}

var (
	A *application
)

func init() {
	A = &application{}
}

func WithGracefulShutdown() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	A.Ctx, A.cancelFunc = context.WithCancel(context.Background())

	go func() {
		sig := <-c
		log.Logger.Info("system call", zap.Any("signal", sig))
		A.cancelFunc()
	}()
}

func Wait() {
	defer A.cancelFunc()
	<-A.Ctx.Done()
}
