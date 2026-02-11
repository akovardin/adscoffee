package analytics

import "github.com/prometheus/client_golang/prometheus"

var actions = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "analytics",
		Subsystem: "actions",
		Name:      "total",
		Help:      "Total number of actions.",
	},
	[]string{"action"},
)

var money = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "analytics",
		Subsystem: "actions",
		Name:      "price",
		Help:      "Total money.",
	},
	[]string{"action"},
)
