package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"go.ads.coffee/platform/pkg/telemetry"
)

const (
	KindInput    = "input"
	KindPipeline = "pipeline"
)

// requests — общий счётчик входящих запросов. RPS/RPM считается в Prometheus как
// rate(server_requests_total[...]) с разбивкой по kind (input|pipeline) и name.
var requests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "total",
		Help:      "Total number of incoming requests (use rate() for RPS/RPM).",
	},
	[]string{"kind", "name"},
)

// errors — счётчик ошибок обработки запросов (на уровне пайплайна).
var errors = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "errors_total",
		Help:      "Total number of failed requests.",
	},
	[]string{"kind", "name"},
)

// duration — время обработки запроса (на уровне пайплайна).
var duration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "server",
		Subsystem: "requests",
		Name:      "duration_seconds",
		Help:      "Request processing duration in seconds.",
		Buckets:   telemetry.DefaultHistogramBuckets,
	},
	[]string{"kind", "name"},
)

// IncInput увеличивает счётчик входящих запросов для input-плагина.
func IncInput(name string) {
	requests.WithLabelValues(KindInput, name).Inc()
}

// IncPipeline увеличивает счётчик входящих запросов для pipeline-роута.
func IncPipeline(name string) {
	requests.WithLabelValues(KindPipeline, name).Inc()
}

// IncError увеличивает счётчик ошибок.
func IncError(kind, name string) {
	errors.WithLabelValues(kind, name).Inc()
}

// ObserveDuration фиксирует время обработки запроса.
func ObserveDuration(kind, name string, d time.Duration) {
	duration.WithLabelValues(kind, name).Observe(d.Seconds())
}
