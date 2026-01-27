package kafkapool

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
)

var (
	// ErrPoolNotFound happens when try to get non existing pool.
	ErrPoolNotFound = errors.New("pool not found")
)

type Pool struct {
	cfgs  map[string]*Config
	pools map[string]*Kafka

	mtrx *Metrics

	logger *zap.Logger

	tracer    Tracer
	telemetry Telemetry
}

func NewPool(
	logger *zap.Logger,
	cfgs map[string]*Config,
	mtrx *Metrics,
	tracer Tracer,
	telemetry Telemetry,
) (*Pool, error) {
	p := &Pool{
		cfgs:  cfgs,
		pools: make(map[string]*Kafka, len(cfgs)),

		logger: logger,

		mtrx:      mtrx,
		tracer:    tracer,
		telemetry: telemetry,
	}

	for name, cfg := range p.cfgs {
		mtr, err := NewNamedMetrics(name, p.telemetry)
		if err != nil {
			return nil, fmt.Errorf("error on create metrics: %w", err)
		}

		kf, err := NewKafka(
			logger.Named(name),
			cfg,
			mtr,
			tracer,
		)
		if err != nil {
			return nil, fmt.Errorf("create kafka cli %s: %w", name, err)
		}
		p.pools[name] = kf
	}

	return p, nil
}

func (p *Pool) GetPool(name string) (*Kafka, error) {
	if kf, ok := p.pools[name]; ok {
		return kf, nil
	}

	return nil, ErrPoolNotFound
}
