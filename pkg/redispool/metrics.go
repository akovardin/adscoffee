package redispool

import (
	"github.com/prometheus/client_golang/prometheus"

	"go.ads.coffee/platform/pkg/telemetry"
)

type Telemetry interface {
	Register(collectors ...prometheus.Collector) error
}

type Metrics struct {
	duration        *prometheus.HistogramVec
	connections     *prometheus.GaugeVec
	connectionsCall *prometheus.GaugeVec

	poolConnCreatedTotal *prometheus.CounterVec
	singleCommands       *prometheus.HistogramVec
	pipelinedCommands    *prometheus.CounterVec
}

type NamedMetrics struct {
	name string

	mtrx *Metrics
}

func NewMetrics(tel Telemetry) (*Metrics, error) {
	duration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "redis_pool_action_duration_seconds",
			Help:    "redis pool action request latencies in seconds.",
			Buckets: telemetry.DefaultHistogramBuckets,
		},
		[]string{"name", "action", telemetry.ErrLabel},
	)

	connections := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redis_pool_connections",
			Help: "redis pool connections",
		},
		[]string{"name", "status"},
	)

	connectionsCall := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redis_pool_connections_calls",
			Help: "Number of redis pool connections calls.",
		},
		[]string{"name", "status"},
	)

	poolConnCreatedTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_pool_connection_created_total",
			Help: "Total number of created connections in pool.",
		},
		[]string{"name", "addr"},
	)

	singleCommands := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "redis_single_commands",
		Help:    "Histogram of single Redis commands",
		Buckets: telemetry.DefaultHistogramBuckets,
	}, []string{"name", "command", telemetry.ErrLabel})

	pipelinedCommands := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "redis_pipelined_commands_total",
		Help: "Number of pipelined Redis commands",
	}, []string{"name", telemetry.ErrLabel})

	if err := tel.Register(
		duration,
		connections,
		connectionsCall,
		poolConnCreatedTotal,
		singleCommands,
		pipelinedCommands,
	); err != nil {
		return nil, err
	}

	return &Metrics{
		duration:             duration,
		connections:          connections,
		connectionsCall:      connectionsCall,
		poolConnCreatedTotal: poolConnCreatedTotal,
		singleCommands:       singleCommands,
		pipelinedCommands:    pipelinedCommands,
	}, nil
}

func (m *Metrics) Duration(name, action string, err error, d float64) {
	m.duration.WithLabelValues(
		name,
		action,
		telemetry.ErrLabelValue(err),
	).Observe(d)
}

func (m *Metrics) Connections(name, status string, val float64) {
	m.connections.WithLabelValues(name, status).Set(val)
}

func (m *Metrics) ConnectionsCall(name, status string, val float64) {
	m.connectionsCall.WithLabelValues(name, status).Set(val)
}

func (m *Metrics) ConnectionCreate(name, addr string) {
	m.poolConnCreatedTotal.WithLabelValues(name, addr).Inc()
}

func (m *Metrics) SingleCommands(name, cmdName string, d float64, err error) {
	m.singleCommands.WithLabelValues(
		name,
		cmdName,
		telemetry.ErrLabelValue(err),
	).Observe(d)
}

func (m *Metrics) PipelinedCommands(name string, err error) {
	m.pipelinedCommands.WithLabelValues(
		name,
		telemetry.ErrLabelValue(err),
	).Inc()
}

func NewNamedMetrics(name string, mtrx *Metrics) *NamedMetrics {
	return &NamedMetrics{
		name: name,
		mtrx: mtrx,
	}
}

func (m *NamedMetrics) Duration(action string, err error, d float64) {
	m.mtrx.Duration(m.name, action, err, d)
}

func (m *NamedMetrics) Connections(status string, val float64) {
	m.mtrx.Connections(m.name, status, val)
}

func (m *NamedMetrics) ConnectionsCall(status string, val float64) {
	m.mtrx.ConnectionsCall(m.name, status, val)
}

func (m *NamedMetrics) ConnectionCreate(addr string) {
	m.mtrx.ConnectionCreate(m.name, addr)
}

func (m *NamedMetrics) SingleCommands(cmdName string, d float64, err error) {
	m.mtrx.SingleCommands(m.name, cmdName, d, err)
}

func (m *NamedMetrics) PipelinedCommands(err error) {
	m.mtrx.PipelinedCommands(m.name, err)
}
