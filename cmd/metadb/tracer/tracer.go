package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serviceName = "metadb"

type Flush func() error

func Init(otlpEndpoint string) (trace.Tracer, Flush, error) {
	ctx := context.Background()
	var tracerProvider trace.TracerProvider
	var conn *grpc.ClientConn
	var err error

	if len(otlpEndpoint) > 0 {
		conn, err = grpc.Dial(otlpEndpoint,
			// Note the use of insecure transport here. TLS is recommended in production.
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if err != nil {
			return nil, flush(tracerProvider, conn), err
		}
		exp, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
		if err != nil {
			return nil, flush(tracerProvider, conn), err
		}

		r, err := resource.Merge(
			resource.Default(),
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			),
		)
		if err != nil {
			return nil, flush(tracerProvider, conn), err
		}

		tracerProvider = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
			sdktrace.WithResource(r),
			sdktrace.WithSampler(
				sdktrace.TraceIDRatioBased(1),
			),
		)
	} else {
		tracerProvider = noop.NewTracerProvider()
	}

	otel.SetTracerProvider(tracerProvider)

	tracer := tracerProvider.Tracer("metadb_start")

	return tracer, flush(tracerProvider, conn), nil
}

func flush(tracerProvider trace.TracerProvider, grpcConn *grpc.ClientConn) Flush {
	return func() error {
		if tp, ok := tracerProvider.(*sdktrace.TracerProvider); ok {
			if err := tp.Shutdown(context.Background()); err != nil {
				return err
			}
			return grpcConn.Close()
		}
		return nil
	}
}