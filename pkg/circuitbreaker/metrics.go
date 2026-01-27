package circuitbreaker

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.ads.coffee/platform/pkg/telemetry"
)

type MetricsType string

const (
	MetricsTypeSuccess                   MetricsType = "success"
	MetricsTypeErrFailure                MetricsType = "failure"
	MetricsTypeErrTimeout                MetricsType = "timeout"
	MetricsTypeErrBadRequest             MetricsType = "bad_request"
	MetricsTypeErrInterrupt              MetricsType = "interrupt"
	MetricsTypeErrConcurrencyLimitReject MetricsType = "concurrency_limit_reject"
	MetricsTypeErrShortCircuit           MetricsType = "short_circuit"
	MetricsTypeOpened                    MetricsType = "opened"
	MetricsTypeClosed                    MetricsType = "closed"
)

type Telemetry interface {
	Register(collectors ...prometheus.Collector) error
}

type Metrics struct {
	total    *prometheus.CounterVec
	duration *prometheus.HistogramVec
	state    *prometheus.CounterVec
}

func NewMetrics(tel Telemetry) (*Metrics, error) {
	const subsystem = "circuit_breaker"

	labelNames := []string{"type", "name", "fallback"}
	stateLabelNames := []string{"type", "name"}

	total := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      "calls_total",
			Help:      "total number of circuit breaker calls.",
		},
		labelNames,
	)

	duration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: subsystem,
			Name:      "call_duration_seconds",
			Help:      "circuit breaker call latencies in seconds.",
			Buckets:   telemetry.DefaultHistogramBuckets,
		},
		labelNames,
	)

	state := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      "state_total",
			Help:      "circuit breaker state",
		},
		stateLabelNames,
	)

	if err := tel.Register(
		total,
		duration,
		state,
	); err != nil {
		return nil, err
	}

	return &Metrics{
		total:    total,
		duration: duration,
		state:    state,
	}, nil
}

func (m *Metrics) Success(
	_ context.Context,
	_ time.Time,
	duration time.Duration,
	fallback bool,
	name string,
) {
	var vals []string

	if fallback {
		vals = []string{string(MetricsTypeSuccess), name, "true"}
	} else {
		vals = []string{string(MetricsTypeSuccess), name, "false"}
	}

	m.total.WithLabelValues(vals...).Inc()
	m.duration.WithLabelValues(vals...).Observe(duration.Seconds())
}

func (m *Metrics) ErrFailure(
	_ context.Context,
	_ time.Time,
	duration time.Duration,
	fallback bool,
	name string,
) {
	var vals []string

	if fallback {
		vals = []string{string(MetricsTypeErrFailure), name, "true"}
	} else {
		vals = []string{string(MetricsTypeErrFailure), name, "false"}
	}

	m.total.WithLabelValues(vals...).Inc()
	m.duration.WithLabelValues(vals...).Observe(duration.Seconds())
}

func (m *Metrics) ErrTimeout(
	_ context.Context,
	_ time.Time,
	duration time.Duration,
	name string,
) {
	vals := []string{string(MetricsTypeErrTimeout), name, "false"}

	m.total.WithLabelValues(vals...).Inc()
	m.duration.WithLabelValues(vals...).Observe(duration.Seconds())
}

func (m *Metrics) ErrBadRequest(
	_ context.Context,
	_ time.Time,
	duration time.Duration,
	name string,
) {
	vals := []string{string(MetricsTypeErrBadRequest), name, "false"}

	m.total.WithLabelValues(vals...).Inc()
	m.duration.WithLabelValues(vals...).Observe(duration.Seconds())
}

func (m *Metrics) ErrInterrupt(
	_ context.Context,
	_ time.Time,
	duration time.Duration,
	name string,
) {
	vals := []string{string(MetricsTypeErrInterrupt), name, "false"}

	m.total.WithLabelValues(vals...).Inc()
	m.duration.WithLabelValues(vals...).Observe(duration.Seconds())
}

func (m *Metrics) ErrConcurrencyLimitReject(
	_ context.Context,
	_ time.Time,
	fallback bool,
	name string,
) {
	var vals []string

	if fallback {
		vals = []string{string(MetricsTypeErrConcurrencyLimitReject), name, "true"}
	} else {
		vals = []string{string(MetricsTypeErrConcurrencyLimitReject), name, "false"}
	}

	m.total.WithLabelValues(vals...).Inc()
}

func (m *Metrics) ErrShortCircuit(
	_ context.Context,
	_ time.Time,
	name string,
) {
	vals := []string{string(MetricsTypeErrShortCircuit), name, "false"}

	m.total.WithLabelValues(vals...).Inc()
}

func (m *Metrics) Opened(
	_ context.Context,
	_ time.Time,
	name string,
) {
	m.state.WithLabelValues(string(MetricsTypeOpened), name).Add(1)
}

func (m *Metrics) Closed(
	_ context.Context,
	_ time.Time,
	name string,
) {
	m.state.WithLabelValues(string(MetricsTypeClosed), name).Add(1)
}

type RunMetrics struct {
	metrics *Metrics

	name string
}

func NewRunMetrics(metrics *Metrics, name string) RunMetrics {
	return RunMetrics{
		metrics: metrics,
		name:    name,
	}
}

func (m RunMetrics) Success(
	ctx context.Context,
	now time.Time,
	duration time.Duration,
) {
	m.metrics.Success(ctx, now, duration, false, m.name)
}

func (m RunMetrics) ErrFailure(
	ctx context.Context,
	now time.Time,
	duration time.Duration,
) {
	m.metrics.ErrFailure(ctx, now, duration, false, m.name)
}

func (m RunMetrics) ErrTimeout(
	ctx context.Context,
	now time.Time,
	duration time.Duration,
) {
	m.metrics.ErrTimeout(ctx, now, duration, m.name)
}

func (m RunMetrics) ErrBadRequest(
	ctx context.Context,
	now time.Time,
	duration time.Duration,
) {
	m.metrics.ErrBadRequest(ctx, now, duration, m.name)
}

func (m RunMetrics) ErrInterrupt(
	ctx context.Context,
	now time.Time,
	duration time.Duration,
) {
	m.metrics.ErrInterrupt(ctx, now, duration, m.name)
}

func (m RunMetrics) ErrConcurrencyLimitReject(
	ctx context.Context,
	now time.Time,
) {
	m.metrics.ErrConcurrencyLimitReject(ctx, now, false, m.name)
}

func (m RunMetrics) ErrShortCircuit(
	ctx context.Context,
	now time.Time,
) {
	m.metrics.ErrShortCircuit(ctx, now, m.name)
}

type FallbackMetrics struct {
	metrics *Metrics

	name string
}

func NewFallbackMetrics(metrics *Metrics, name string) FallbackMetrics {
	return FallbackMetrics{
		metrics: metrics,
		name:    name,
	}
}

func (m FallbackMetrics) Success(
	ctx context.Context,
	now time.Time,
	duration time.Duration,
) {
	m.metrics.Success(ctx, now, duration, true, m.name)
}

func (m FallbackMetrics) ErrFailure(
	ctx context.Context,
	now time.Time,
	duration time.Duration,
) {
	m.metrics.ErrFailure(ctx, now, duration, true, m.name)
}

func (m FallbackMetrics) ErrConcurrencyLimitReject(
	ctx context.Context,
	now time.Time,
) {
	m.metrics.ErrConcurrencyLimitReject(ctx, now, true, m.name)
}

type CircuitMetrics struct {
	metrics *Metrics

	name string
}

func NewCircuitMetrics(metrics *Metrics, name string) CircuitMetrics {
	return CircuitMetrics{
		metrics: metrics,
		name:    name,
	}
}

func (m CircuitMetrics) Opened(
	ctx context.Context,
	now time.Time,
) {
	m.metrics.Opened(ctx, now, m.name)
}

func (m CircuitMetrics) Closed(
	ctx context.Context,
	now time.Time,
) {
	m.metrics.Closed(ctx, now, m.name)
}
