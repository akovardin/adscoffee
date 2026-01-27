package redispool

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"go.ads.coffee/platform/pkg/circuitbreaker"
	"go.ads.coffee/platform/pkg/health"
)

const (
	CircuitNamePrefix = "redis:"
)

var (
	// ErrPoolNotFound happens when try to get non existing pool.
	ErrPoolNotFound = errors.New("pool not found")
)

type Tracer interface {
	StartSpan(
		ctx context.Context,
		spanName string,
		opts ...trace.SpanStartOption,
	) (context.Context, trace.Span)
}

type MetricsProvider interface {
	Duration(action string, err error, d float64)
	Connections(status string, val float64)
	ConnectionsCall(status string, val float64)
	ConnectionCreate(addr string)
	SingleCommands(cmdName string, d float64, err error)
	PipelinedCommands(err error)
}

type Pool struct {
	cfgs  map[string]*Config
	pools map[string]*Redis
	mu    sync.RWMutex

	logger *zap.Logger

	mtrx   *Metrics
	tracer Tracer
}

func NewPool(
	logger *zap.Logger,
	cfgs map[string]*Config,
	mtrx *Metrics,
	tracer Tracer,
	pool *circuitbreaker.Pool,
) (*Pool, error) {
	p := &Pool{
		cfgs:  cfgs,
		pools: make(map[string]*Redis, len(cfgs)),

		logger: logger,

		mtrx:   mtrx,
		tracer: tracer,
	}

	for name, cfg := range p.cfgs {
		redis, err := NewRedis(
			logger,
			cfg,
			pool.Get(CircuitName(name)),
			NewNamedMetrics(name, mtrx),
			p.tracer,
		)

		if err != nil {
			return nil, err
		}

		p.pools[name] = redis
	}

	return p, nil
}

func (p *Pool) GetPool(name string) (*Redis, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if pool, ok := p.pools[name]; ok {
		return pool, nil
	}

	return nil, ErrPoolNotFound
}

func (p *Pool) Stop(_ context.Context) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, pool := range p.pools {
		err := pool.Close()
		if err != nil {
			p.logger.Error("stop redis", zap.Error(err))

			return err
		}
	}

	return nil
}

func (p *Pool) HealthComponent() *health.Component {
	return &health.Component{
		Kind: health.ComponentKindLocal,
		Name: "redis_pool",
		CheckFunc: func(ctx context.Context) error {
			p.mu.RLock()
			defer p.mu.RUnlock()

			for name, conn := range p.pools {
				err := conn.Ping(ctx)
				if err != nil {
					return fmt.Errorf("%s: %w", name, err)
				}
			}

			return nil
		},
	}
}

func CircuitName(name string) string {
	return fmt.Sprintf("%s%s", CircuitNamePrefix, name)
}
