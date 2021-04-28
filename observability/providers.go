package observability

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	zipkinExporter "go.opentelemetry.io/otel/exporters/trace/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

// ProvidersConfig will
type ProvidersConfig struct {
	ServiceName string
	ZipkinURL   string
	JaegerURL   string
}

// Providers will
type Providers struct {
	Stdout *tracesdk.TracerProvider
	Zipkin *tracesdk.TracerProvider
	Jaeger *tracesdk.TracerProvider
}

// InitTracerProviders will return a struct with both providers: jaeger and stdout
func (c ProvidersConfig) InitTracerProviders() (p Providers, err error) {
	zp, err := ZipkinTracerProvider(c.ServiceName, c.ZipkinURL)
	if err != nil {
		return p, err
	}

	jp, err := JaegerTracerProvider(c.ServiceName, c.JaegerURL)
	if err != nil {
		return p, err
	}

	sp, err := StdoutTracerProvider(c.ServiceName)
	if err != nil {
		return p, err
	}

	p = Providers{
		Zipkin: zp,
		Jaeger: jp,
		Stdout: sp,
	}

	return p, nil
}

// JaegerTracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func JaegerTracerProvider(service string, url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter

	exp, err := jaeger.NewRawExporter(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(url),
		),
	)

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
		)),
	)

	if err != nil {
		return nil, err
	}
	return tp, nil
}

// ZipkinTracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Zipkin exporter
func ZipkinTracerProvider(service string, url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := zipkinExporter.NewRawExporter(url)

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
		)),
	)

	if err != nil {
		return nil, err
	}
	return tp, nil
}

// StdoutTracerProvider returns an OpenTelemetry TracerProvider configured to use
// the stdout exporter
func StdoutTracerProvider(service string) (*tracesdk.TracerProvider, error) {
	exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithSyncer(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", environment),
		)),
	)

	return tp, nil
}
