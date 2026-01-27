package redispool

import (
	"context"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	metricsHook struct {
		metrics MetricsProvider
	}
)

func newMetricsHook(metrics MetricsProvider) *metricsHook {
	return &metricsHook{
		metrics: metrics,
	}
}

func (h *metricsHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		h.metrics.ConnectionCreate(addr)

		return next(ctx, network, addr)
	}
}

func (h *metricsHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		defer func(t time.Time) {
			h.metrics.SingleCommands(cmd.Name(), time.Since(t).Seconds(), cmd.Err())
		}(time.Now())

		err := next(ctx, cmd)
		if err != nil {
			return err
		}

		return nil
	}
}

func (h *metricsHook) ProcessPipelineHook(
	next redis.ProcessPipelineHook,
) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) (err error) {
		defer func(t time.Time) {
			h.metrics.SingleCommands("pipeline", time.Since(t).Seconds(), err)
			h.metrics.PipelinedCommands(err)
		}(time.Now())

		err = next(ctx, cmds)
		if err != nil {
			return err
		}

		return nil
	}
}
