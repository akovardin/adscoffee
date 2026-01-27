package redispool_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"go.ads.coffee/platform/pkg/circuitbreaker"
	"go.ads.coffee/platform/pkg/redispool"
)

var (
	defaultRedisSeeds = "127.0.0.1:7001,127.0.0.1:7002,127.0.0.1:7003,127.0.0.1:7004,127.0.0.1:7005,127.0.0.1:7006"
	redisSeeds        = os.Getenv("REDIS_SEEDS")
)

type tracer struct {
}

func (t *tracer) StartSpan(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return ctx, trace.SpanFromContext(ctx)
}

type telemetry struct {
}

func (t *telemetry) Register(collectors ...prometheus.Collector) error {
	return nil
}

func TestIntegration_Redis(t *testing.T) {
	logger := zaptest.NewLogger(t, zaptest.Level(zap.DebugLevel))

	if redisSeeds == "" {
		redisSeeds = defaultRedisSeeds
	}

	mtrx, _ := redispool.NewMetrics(&telemetry{})
	nm := redispool.NewNamedMetrics("test", mtrx)

	cfg := &redispool.Config{
		KeyPrefix:      "test",
		ClusterAddrs:   strings.Split(redisSeeds, ","),
		PoolSize:       4,
		MaxRetries:     3,
		MaxRedirects:   8,
		RouteByLatency: true,
	}

	circuit := &circuitMock{}
	circuit.On("Run", mock.Anything, mock.Anything)

	rds, err := redispool.NewRedis(logger, cfg, circuit, nm, &tracer{})
	require.NoError(t, err)
	require.NotNil(t, rds)

	defer rds.Close()

	t.Run("primary_call", func(t *testing.T) {
		testKey := "primary"
		value := "value"
		err := rds.Call(
			context.Background(),
			"set_primary",
			func(ctx context.Context, clu *redis.ClusterClient, kf redispool.KeyFormatter) error {
				return clu.Set(ctx, kf.FormatKey(testKey), value, 0).Err()
			},
		)

		require.NoError(t, err)
	})

	t.Run("primary_secondary_call", func(t *testing.T) {
		testKey := "primary_secondary"
		value := "value"
		err := rds.Call(
			context.Background(),
			"set_primary",
			func(ctx context.Context, clu *redis.ClusterClient, kf redispool.KeyFormatter) error {
				return clu.Set(ctx, kf.FormatKey(testKey), value, 0).Err()
			},
		)
		require.NoError(t, err)

		var resValue string
		err = rds.Call(
			context.Background(),
			"get_secondary_primary",
			func(ctx context.Context, clu *redis.ClusterClient, kf redispool.KeyFormatter) error {
				return clu.Get(ctx, kf.FormatKey(testKey)).Scan(&resValue)
			},
		)

		require.NoError(t, err)
		require.Equal(t, value, resValue)
	})
}

func TestIntegration_RedisWithCircuitBreaker(t *testing.T) {
	logger := zaptest.NewLogger(t, zaptest.Level(zap.DebugLevel))

	mtrx, _ := redispool.NewMetrics(&telemetry{})
	nm := redispool.NewNamedMetrics("test", mtrx)

	cfg := &redispool.Config{
		KeyPrefix:      "test",
		ClusterAddrs:   []string{"bad_addr:7001"},
		PoolSize:       4,
		MaxRetries:     3,
		MaxRedirects:   8,
		RouteByLatency: true,
	}

	cbMetrix, err := circuitbreaker.NewMetrics(&telemetry{})
	require.NoError(t, err)

	circuitName := redispool.CircuitName("test")

	cbPool, err := circuitbreaker.NewPool(
		logger,
		cbMetrix,
		map[string]*circuitbreaker.Config{
			circuitName: {
				Enabled:               true,
				Timeout:               100 * time.Millisecond,
				MaxConcurrentRequests: 1,
				Hystrix: circuitbreaker.HystrixConfig{
					RequestVolumeThreshold: 2,
				},
			},
		},
	)

	require.NoError(t, err)

	rds, err := redispool.NewRedis(logger, cfg, cbPool.Get(circuitName), nm, &tracer{})
	require.NoError(t, err)
	require.NotNil(t, rds)

	defer rds.Close()

	t.Run("primary_call_cb", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			testKey := "primary"
			value := "value"
			err := rds.Call(
				context.Background(),
				"set_primary",
				func(ctx context.Context, clu *redis.ClusterClient, kf redispool.KeyFormatter) error {
					return clu.Set(ctx, kf.FormatKey(testKey), value, 0).Err()
				},
			)

			if i > 1 {
				var cbErr circuitbreaker.Error

				require.ErrorAs(t, err, &cbErr)
			} else {
				require.Error(t, err)
			}
		}
	})
}

type circuitMock struct {
	mock.Mock
}

func (c *circuitMock) Go(ctx context.Context, runFunc func(context.Context) error, fallbackFunc func(context.Context, error) error) error {
	c.Called(ctx, runFunc, fallbackFunc)

	return runFunc(ctx)
}

func (c *circuitMock) Run(ctx context.Context, runFunc func(context.Context) error) error {
	c.Called(ctx, runFunc)

	return runFunc(ctx)
}

func (c *circuitMock) Execute(ctx context.Context, runFunc func(context.Context) error, fallbackFunc func(context.Context, error) error) error {
	c.Called(ctx, runFunc, fallbackFunc)

	return runFunc(ctx)
}
