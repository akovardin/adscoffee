package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.uber.org/zap"

	otelBridge "go.opentelemetry.io/otel/bridge/opentracing"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type Telemetry struct {
	cfg    Config
	logger *zap.Logger

	registry       prometheus.Registerer
	tracer         trace.TracerProvider
	bridgeTracer   *otelBridge.BridgeTracer
	tracerStopFunc func(ctx context.Context) error
}

func New(
	logger *zap.Logger,
	config Config,
	pr prometheus.Registerer,
) (*Telemetry, error) {
	t := &Telemetry{
		cfg:          config,
		logger:       logger,
		registry:     pr,
		tracer:       noop.NewTracerProvider(), // default empty tracer
		bridgeTracer: otelBridge.NewBridgeTracer(),
	}

	{ // tracer
		if config.Jaeger.Enabled {
			var err error
			t.tracer, t.bridgeTracer, t.tracerStopFunc, err = tracerProvider(config)
			if err != nil {
				return nil, fmt.Errorf("tracer provider: %w", err)
			}
		}
	}

	return t, nil
}

func (t *Telemetry) Registry() prometheus.Registerer {
	return t.registry
}

func (t *Telemetry) Register(collectors ...prometheus.Collector) error {
	for _, collector := range collectors {
		if err := t.registry.Register(collector); err != nil {
			return fmt.Errorf("failed to register collector: %w", err)
		}
	}

	return nil
}

func (t *Telemetry) Stop(ctx context.Context) error {
	if t.tracerStopFunc != nil {
		err := t.tracerStopFunc(ctx)
		if err != nil {
			return fmt.Errorf("stop func: %w", err)
		}
	}

	return nil
}

func tracerProvider(config Config) (
	trace.TracerProvider,
	*otelBridge.BridgeTracer,
	func(ctx context.Context) error,
	error,
) {
	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName(config.ServiceName),
			semconv.ServiceVersion(config.Version),
			semconv.ServiceInstanceID(config.Hostname),
			semconv.HostName(config.Hostname),
		),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// nolint:staticcheck // TODO change later
	conn, err := grpc.DialContext(ctx, config.Jaeger.Endpoint,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // nolint:staticcheck // TODO change later
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithGRPCConn(conn),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(config.Jaeger.SamplingRatio/100.0)),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	bridgeTracer, wrapperTracerProvider := otelBridge.NewTracerPair(tracerProvider.Tracer(""))

	// set the Tracer Provider and the W3C Trace Context propagator as globals
	otel.SetTracerProvider(wrapperTracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return wrapperTracerProvider, bridgeTracer, tracerProvider.Shutdown, nil
}

func (t *Telemetry) StartSpan(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return t.tracer.Tracer(t.cfg.ServiceName).Start(ctx, spanName, opts...)
}

func (t *Telemetry) TracerProvider() trace.TracerProvider {
	return t.tracer
}

func (t *Telemetry) OpenTracer() *otelBridge.BridgeTracer {
	return t.bridgeTracer
}
