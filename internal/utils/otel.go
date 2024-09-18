package utils

import (
	"context"
	"os"
	"time"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	otelTrace "go.opentelemetry.io/otel/trace"
)

const OTEL_COLLECTOR_ENDPOINT = "localhost:4317"
const SpanTracer = "order-service-tracer"

func InitOtelInstrumentation() {
	ctx := context.Background()

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("order-service"),
		semconv.DeploymentEnvironment(os.Getenv("ENV")),
	)

	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure(), otlpmetricgrpc.WithEndpoint(OTEL_COLLECTOR_ENDPOINT))
	if err != nil {
		panic(err)
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter, metric.WithInterval(time.Second*10))),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(OTEL_COLLECTOR_ENDPOINT))
	if err != nil {
		panic(err)
	}

	traceID := xray.NewIDGenerator()
	tracerProvider := trace.NewTracerProvider(
		trace.WithIDGenerator(traceID),
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)
}

func Trace(context context.Context, name string) (context.Context, otelTrace.Span) {
	tracer := otel.Tracer(SpanTracer)
	return tracer.Start(context, name)
}
