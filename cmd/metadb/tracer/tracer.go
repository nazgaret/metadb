package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type Flush func() error

func Init(url string) (trace.Tracer, Flush, error) {
	ctx := context.Background()
	var tracerProvider trace.TracerProvider

	if len(url) > 0 {
		exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(url))
		if err != nil {
			return nil, flush(tracerProvider), err
		}

		tracerProvider = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
		)
	} else {
		tracerProvider = noop.NewTracerProvider()
	}

	otel.SetTracerProvider(tracerProvider)

	tracer := tracerProvider.Tracer("metadb_start")

	return tracer, flush(tracerProvider), nil
}

func flush(tracerProvider trace.TracerProvider) Flush {
	return func() error {
		if tp, ok := tracerProvider.(*sdktrace.TracerProvider); ok {
			if err := tp.Shutdown(context.Background()); err != nil {
				return err
			}
		}
		return nil
	}
}
