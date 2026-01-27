package kafkapool

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/chapsuk/wait"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kotel"
	"go.uber.org/zap"

	"go.ads.coffee/platform/pkg/health"
)

// ConsumerOpt allows customize consumer.
type ConsumerOpt func(c *Consumer)

func WithConsumeResetOffsetConsumerOpt(offset kgo.Offset) ConsumerOpt {
	return func(c *Consumer) {
		c.opts = append(c.opts, kgo.ConsumeResetOffset(offset))
	}
}

func WithCustomErrorHandlerConsumerOpt(eh func(err error)) ConsumerOpt {
	return func(c *Consumer) {
		c.errHandler = eh
	}
}

type Consumer struct {
	kf     *Kafka
	logger *zap.Logger

	opts       []kgo.Opt
	errHandler func(error)

	client *kgo.Client

	consumers map[string]tconsumer
	wg        wait.Group
}

type tconsumer struct {
	mtrx          MetricsProvider
	kTracer       *kotel.Tracer
	logger        *zap.Logger
	ctx           context.Context
	contextCancel context.CancelFunc
	rec           chan *kgo.Record
	backoff       backoff.BackOff
}

func NewConsumer(kf *Kafka, opts ...ConsumerOpt) (*Consumer, error) {
	c := &Consumer{
		kf:        kf,
		logger:    kf.logger.Named("consumer"),
		opts:      make([]kgo.Opt, 0, len(opts)),
		consumers: make(map[string]tconsumer, 1),
	}
	c.errHandler = func(err error) {
		c.logger.Warn("Kafka: consumer error", zap.Error(err))
	}

	for _, opt := range opts {
		opt(c)
	}

	c.opts = append(c.opts, kf.opts...)

	cl, err := kgo.NewClient(c.opts...)
	if err != nil {
		return nil, err
	}

	c.client = cl

	return c, nil
}

func (c *Consumer) Subscribe(
	topic string,
	callback ConsumerCallback,
) {
	topic = c.kf.formatWithPrefix(topic)
	c.client.AddConsumeTopics(topic)

	ctx, cancel := context.WithCancel(context.Background())
	tc := tconsumer{
		mtrx:          c.kf.metrics,
		kTracer:       c.kf.kTracer,
		logger:        c.logger.Named("topic_consumer"),
		rec:           make(chan *kgo.Record),
		ctx:           ctx,
		contextCancel: cancel,
		backoff:       c.kf.cfg.Consumer.ConfigBackoff.GetBackOff(),
	}

	c.wg.Add(func() {
		tc.consume(topic, callback, c.errHandler)
	})

	c.consumers[topic] = tc
}

func (c *Consumer) Consume() {
	c.wg.Add(func() {
		for {
			fetches := c.client.PollFetches(context.Background())
			if fetches.IsClientClosed() {
				c.logger.Info("stop consumer")

				return
			}
			fetches.EachError(func(t string, p int32, err error) {
				c.errHandler(fmt.Errorf("topic: %s, partition: %d error: %w", t, p, err))
			})

			fetches.EachTopic(c.consumeTopic)
		}
	})
}

func (c *Consumer) Ping(ctx context.Context) error {
	if err := c.client.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping Kafka: %w", err)
	}

	return nil
}

func (c *Consumer) consumeTopic(t kgo.FetchTopic) {
	tconsumers, ok := c.consumers[t.Topic]
	if !ok {
		return
	}

	t.EachRecord(func(r *kgo.Record) {
		tconsumers.rec <- r
	})
}

func (c *Consumer) Close() {
	c.client.Close()
	for _, tc := range c.consumers {
		tc.stop()
	}

	c.wg.Wait()
}

func (c *Consumer) HealthComponent() *health.Component {
	return &health.Component{
		Kind: health.ComponentKindExternal,
		Name: "kafka_consumer",
		CheckFunc: func(ctx context.Context) error {
			return c.Ping(ctx)
		},
	}
}

func (tc *tconsumer) consume(
	topic string,
	callback ConsumerCallback,
	errHandler func(err error),
) {
	tc.logger.Info("start consume", zap.String("topic", topic))
	for {
		select {
		case <-tc.ctx.Done():
			tc.logger.Info("stop consuming", zap.String("topic", topic))

			return
		case rec := <-tc.rec:
			start := time.Now()
			ctx, span := tc.kTracer.WithProcessSpan(rec)

			tryNum := 0
			err := backoff.RetryNotify(
				func() error {
					tryNum++

					return callback(ctx, convertRecordToPayload(rec))

				},
				backoff.WithContext(tc.backoff, ctx),
				func(err error, duration time.Duration) {
					tc.logger.Warn(
						"consumer call callback, retrying",
						zap.String("topic", topic),
						zap.Int("try_num", tryNum),
						zap.Duration("retry_delay", duration),
						zap.Error(err),
					)
				},
			)

			tc.mtrx.ConsumerHandleDuration(topic, err, time.Since(start).Seconds())

			if err != nil {
				span.RecordError(err)
				errHandler(err)
			}

			span.End()

		}
	}
}

func (tc *tconsumer) stop() {
	tc.contextCancel()
}
