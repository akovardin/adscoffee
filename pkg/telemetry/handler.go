package telemetry

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (t *Telemetry) Handler() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}
