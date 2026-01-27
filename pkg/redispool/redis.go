package redispool

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"go.ads.coffee/platform/pkg/circuitbreaker"
	"go.ads.coffee/platform/pkg/health"
)

const (
	writeStatsPeriod = 10 * time.Second
)

type CallCallback func(
	ctx context.Context,
	clu *redis.ClusterClient,
	keyFormatter KeyFormatter,
) error

type Redis struct {
	cfg    *Config
	logger *zap.Logger
	rd     *redis.ClusterClient

	metrics MetricsProvider
	circuit circuitbreaker.Circuit
	tracer  Tracer
}

func NewRedis(
	logger *zap.Logger,
	cfg *Config,
	circuit circuitbreaker.Circuit,
	metrics MetricsProvider,
	tracer Tracer,
) (*Redis, error) {
	redisClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    cfg.ClusterAddrs,
		Username: cfg.Username,
		Password: cfg.Password,

		MaxRetries:      cfg.MaxRetries,
		MinRetryBackoff: cfg.MinRetryBackoff,
		MaxRetryBackoff: cfg.MaxRetryBackoff,

		MaxRedirects: cfg.MaxRedirects,

		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,

		PoolSize:        cfg.PoolSize,
		PoolTimeout:     cfg.PoolTimeout,
		MinIdleConns:    cfg.MinIdleConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		MaxActiveConns:  cfg.MaxActiveConns,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
		ConnMaxLifetime: cfg.ConnMaxLifetime,

		ReadOnly: true,

		RouteRandomly:  cfg.RouteRandomly,
		RouteByLatency: cfg.RouteByLatency,
	})

	redisClient.AddHook(newMetricsHook(metrics)) // metrics hook

	r := &Redis{
		cfg:     cfg,
		logger:  logger,
		rd:      redisClient,
		metrics: metrics,
		circuit: circuit,

		tracer: tracer,
	}

	go r.reportStats()

	return r, nil
}

func (r *Redis) FormatKey(key string) string {
	if r.cfg.KeyPrefix == "" {
		return key
	}

	return fmt.Sprintf("%s:%s", r.cfg.KeyPrefix, key)
}

func (r *Redis) Call(
	ctx context.Context, name string, callback CallCallback,
) (err error) {
	if !r.cfg.Enabled {
		return nil
	}

	return r.circuit.Run(ctx, func(ctx context.Context) error {
		return r.call(ctx, name, callback)
	})
}

func (r *Redis) call(
	ctx context.Context, name string, callback CallCallback,
) (err error) {
	spanName := fmt.Sprintf("call_%s", name)
	ctx, span := r.tracer.StartSpan(ctx, spanName, trace.WithAttributes(
		attribute.String("redis.query.name", name),
	))
	defer span.End()

	defer func(ts time.Time) {
		r.metrics.Duration(
			name,
			err,
			time.Since(ts).Seconds(),
		)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
		}
	}(time.Now())

	err = callback(ctx, r.rd, r)
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Close() error {
	return r.rd.Close()
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.rd.Ping(ctx).Err()
}

func (r *Redis) HealthComponent() *health.Component {
	return &health.Component{
		Kind: health.ComponentKindLocal,
		Name: "redis",
		CheckFunc: func(ctx context.Context) error {
			return r.Ping(ctx)
		},
	}
}

func (r *Redis) reportStats() {
	for range time.NewTicker(writeStatsPeriod).C {
		stat := r.rd.PoolStats()
		r.metrics.ConnectionsCall("hits", float64(stat.Hits))
		r.metrics.ConnectionsCall("misses", float64(stat.Misses))
		r.metrics.ConnectionsCall("timeouts", float64(stat.Timeouts))

		r.metrics.Connections("idle", float64(stat.IdleConns))
		r.metrics.Connections("stale", float64(stat.StaleConns))
		r.metrics.Connections("total", float64(stat.TotalConns))
	}
}
