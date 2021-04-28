package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	environment = "development"
	id          = 1
)

// InitGlobalTrace will initiate a global tracer based on provider
func InitGlobalTrace(tracerProvider *tracesdk.TracerProvider) {
	tp := tracerProvider

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

// Span will test a child span
func Span(ctx context.Context, component string, name string, spanTags []attribute.KeyValue) (childCtx context.Context, span trace.Span) {
	// Use the global TracerProvider.
	tr := otel.Tracer(component)

	childCtx, span = tr.Start(ctx, name)
	if len(spanTags) > 0 {
		for _, kvTag := range spanTags {
			fmt.Println(kvTag)
			span.SetAttributes(attribute.Key(kvTag.Key).String(kvTag.Value.AsString()))
		}
	}

	return childCtx, span
}
