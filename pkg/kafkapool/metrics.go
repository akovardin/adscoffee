package kafkapool

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"go.ads.coffee/platform/pkg/telemetry"
)

type Telemetry interface {
	Register(collectors ...prometheus.Collector) error
}

type Metrics struct {
	producerTotal       *prometheus.CounterVec
	consumerHandle      *prometheus.HistogramVec
	consumerGroupHandle *prometheus.HistogramVec
}

func NewMetrics(mr Telemetry) (*Metrics, error) {
	producerTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kafka",
			Subsystem: "producer",
			Name:      "send_count_total",
			Help:      "Produced events count",
		},
		[]string{"topic", telemetry.ErrLabel},
	)

	consumerHandle := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "kafka",
			Subsystem: "consumer",
			Name:      "duration_seconds",
			Help:      "Time elapsed to consume single messsage",
			Buckets:   telemetry.DefaultHistogramBuckets,
		},
		[]string{"topic", telemetry.ErrLabel},
	)

	consumerGroupHandle := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "kafka",
			Subsystem: "consumer_group",
			Name:      "duration_seconds",
			Help:      "Time elapsed to consume single messsage",
			Buckets:   telemetry.DefaultHistogramBuckets,
		},
		[]string{"group", "topic", telemetry.ErrLabel},
	)

	err := mr.Register(producerTotal, consumerHandle, consumerGroupHandle)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	return &Metrics{
		producerTotal:       producerTotal,
		consumerHandle:      consumerHandle,
		consumerGroupHandle: consumerGroupHandle,
	}, nil
}

func (m *Metrics) ProducerTotal(topic string, err error) {
	m.producerTotal.WithLabelValues(topic, telemetry.ErrLabelValue(err)).Inc()
}

func (m *Metrics) ConsumerHandleDuration(topic string, err error, seconds float64) {
	m.consumerHandle.WithLabelValues(topic, telemetry.ErrLabelValue(err)).Observe(seconds)
}

func (m *Metrics) ConsumerGroupHandleDuration(group, topic string, err error, seconds float64) {
	m.consumerGroupHandle.WithLabelValues(group, topic, telemetry.ErrLabelValue(err)).Observe(seconds)
}

type NamedMetrics struct {
	name                string
	producerTotal       *prometheus.CounterVec
	consumerHandle      *prometheus.HistogramVec
	consumerGroupHandle *prometheus.HistogramVec
}

func NewNamedMetrics(name string, mr Telemetry) (*NamedMetrics, error) {
	producerTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "kafka_pool",
			Subsystem: "producer",
			Name:      "send_count_total",
			Help:      "Produced events count",
		},
		[]string{"name", "topic", telemetry.ErrLabel},
	)

	consumerHandle := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "kafka_pool",
			Subsystem: "consumer",
			Name:      "duration_seconds",
			Help:      "Time elapsed to consume single messsage",
			Buckets:   telemetry.DefaultHistogramBuckets,
		},
		[]string{"name", "topic", telemetry.ErrLabel},
	)

	consumerGroupHandle := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "kafka_pool",
			Subsystem: "consumer_group",
			Name:      "duration_seconds",
			Help:      "Time elapsed to consume single messsage",
			Buckets:   telemetry.DefaultHistogramBuckets,
		},
		[]string{"name", "group", "topic", telemetry.ErrLabel},
	)

	err := mr.Register(producerTotal, consumerHandle, consumerGroupHandle)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	return &NamedMetrics{
		name:                name,
		producerTotal:       producerTotal,
		consumerHandle:      consumerHandle,
		consumerGroupHandle: consumerGroupHandle,
	}, nil
}

func (m *NamedMetrics) ProducerTotal(topic string, err error) {
	m.producerTotal.WithLabelValues(m.name, topic, telemetry.ErrLabelValue(err)).Inc()
}

func (m *NamedMetrics) ConsumerHandleDuration(topic string, err error, seconds float64) {
	m.consumerHandle.WithLabelValues(m.name, topic, telemetry.ErrLabelValue(err)).Observe(seconds)
}

func (m *NamedMetrics) ConsumerGroupHandleDuration(
	group, topic string,
	err error,
	seconds float64,
) {
	m.consumerGroupHandle.
		WithLabelValues(
			m.name,
			group,
			topic,
			telemetry.ErrLabelValue(err),
		).Observe(seconds)
}
