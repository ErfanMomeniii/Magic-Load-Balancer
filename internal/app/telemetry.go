package app

import (
	"context"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/config"
	"github.com/ErfanMomeniii/Magic-Load-Balancer/internal/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.uber.org/zap"
	"os"
	"time"
)

func WithTelemetry() func() {
	cfg := config.C.Tracing
	if !cfg.Enabled {
		return func() {}
	}

	exp, err := jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost(cfg.AgentHost),
		jaeger.WithAgentPort(cfg.AgentPort)))
	if err != nil {
		log.Logger.Fatal("error in setting jaeger endpoint", zap.Error(err))
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Logger.Fatal("error in getting hostname", zap.Error(err))
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SamplerRatio)),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceInstanceIDKey.String(hostname),
			semconv.ServiceNameKey.String(Name),
			semconv.ServiceVersionKey.String(GitTag),
		)),
	)

	deferFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
		defer cancel()

		err := traceProvider.Shutdown(ctx)
		if err != nil {
			log.Logger.Fatal("error in shutting down telemetry", zap.Error(err))
		}
	}

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	Tracer = otel.Tracer(Name)

	log.Logger.Info("tracing set up successfully")

	return deferFunc
}
