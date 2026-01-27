package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"

	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"go.ads.coffee/platform/pkg/telemetry"
)

type Telemetry interface {
	Register(collectors ...prometheus.Collector) error
}

func Metrics(tel Telemetry) (func(next http.Handler) http.Handler, error) {
	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "http",
			Subsystem: "requests",
			Name:      "total",
			Help:      "Total number of HTTP requests.",
		},
		[]string{"method", "pattern", "status"},
	)

	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "http",
			Subsystem: "requests",
			Name:      "duration",
			Help:      "HTTP request latencies in seconds.",
			Buckets:   telemetry.DefaultHistogramBuckets,
		},
		[]string{"method", "pattern"},
	)

	if err := tel.Register(requestsTotal, requestDuration); err != nil {
		return nil, err
	}

	metricsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			chictx := chi.RouteContext(r.Context())
			ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

			ts := time.Now()

			next.ServeHTTP(ww, r)

			cost := time.Since(ts)

			requestsTotal.WithLabelValues(
				r.Method,
				chictx.RoutePattern(),
				fmt.Sprint(ww.Status()),
			).Inc()

			requestDuration.WithLabelValues(
				r.Method,
				chictx.RoutePattern(),
			).Observe(cost.Seconds())
		})
	}

	return metricsHandler, nil
}
