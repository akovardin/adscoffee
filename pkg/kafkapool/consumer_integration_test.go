package kafkapool_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.ads.coffee/platform/pkg/kafkapool"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type simpleTaskPayload struct {
	Name string `json:"name"`
}

var (
	defaultKafkaSeeds = "kafka:9092"
	kafkaSeeds        = os.Getenv("KAFKA_SEEDS")
)

type counter struct {
	mu      sync.RWMutex
	counter int64
}

func (c *counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.counter = 0
}

func (c *counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.counter++
}

func (c *counter) Counter() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.counter
}

type telemetry struct {
}

func (t *telemetry) Register(collectors ...prometheus.Collector) error {
	return nil
}

type tracer struct {
}

func (t *tracer) StartSpan(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return ctx, trace.SpanFromContext(ctx)
}

func (t *tracer) TracerProvider() trace.TracerProvider {
	return noop.NewTracerProvider()
}

func TestIntegration_Consumer(t *testing.T) {
	const (
		taskNum             int64 = 30
		producerMaxDuration       = 5 * time.Second
		consumerMaxDuration       = 5 * time.Second
		topic                     = "test_topic"
	)

	c := &counter{}
	topicHandler := func(_ context.Context, p kafkapool.Payload) error {
		c.Inc()

		return nil
	}

	logger := zaptest.NewLogger(t, zaptest.Level(zap.DebugLevel))
	seeds := kafkaSeeds
	if seeds == "" {
		seeds = defaultKafkaSeeds
	}
	cfg := kafkapool.Config{
		Prefix:                 "inter_test",
		Seeds:                  strings.Split(seeds, ","),
		AllowAutoTopicCreation: true,
	}

	mtrx, _ := kafkapool.NewMetrics(&telemetry{})

	kk, err := kafkapool.NewKafka(logger, &cfg, mtrx, &tracer{})
	require.NoError(t, err)
	require.NotNil(t, kk)

	producer, err := kafkapool.NewProducer(kk)
	require.NoError(t, err)
	require.NotNil(t, producer)

	for i := 1; i <= int(taskNum); i++ {
		taskName := fmt.Sprintf("Task_%d", i)
		payload, err := json.Marshal(simpleTaskPayload{
			Name: taskName,
		})
		require.NoError(t, err)
		err = producer.Send(context.Background(), topic, payload)
		require.NoError(t, err)
	}

	t.Run("consuming_from_start", func(t *testing.T) {
		consumer, err := kafkapool.NewConsumer(kk)
		require.NoError(t, err)
		require.NotNil(t, consumer)

		consumer.Subscribe(topic, topicHandler)

		consumer.Consume()

		assert.Eventuallyf(t, func() bool {
			return taskNum == c.Counter()
		}, 5*time.Second, 50*time.Millisecond,
			"tasks: %d counter: %d",
			taskNum, c.Counter(),
		)

		require.Equal(t, taskNum, c.Counter())
		consumer.Close()
	})

	t.Run("reconsume_from_start", func(t *testing.T) {
		c.Reset()

		consumer, err := kafkapool.NewConsumer(kk)
		require.NoError(t, err)
		require.NotNil(t, consumer)

		consumer.Subscribe(topic, topicHandler)

		consumer.Consume()

		assert.Eventuallyf(t, func() bool {
			return taskNum == c.Counter()
		}, 5*time.Second, 50*time.Millisecond,
			"tasks: %d counter: %d",
			taskNum, c.Counter(),
		)

		require.Equal(t, taskNum, c.Counter())
		consumer.Close()
	})

	t.Run("consume_only_new", func(t *testing.T) {
		n := time.Now()
		for i := 1; i <= int(taskNum); i++ {
			taskName := fmt.Sprintf("Task_%d_new", i)
			payload, err := json.Marshal(simpleTaskPayload{
				Name: taskName,
			})
			require.NoError(t, err)
			err = producer.Send(context.Background(), topic, payload)
			require.NoError(t, err)
		}

		c.Reset()

		consumer, err := kafkapool.NewConsumer(kk, kafkapool.WithConsumeResetOffsetConsumerOpt(kgo.NewOffset().AfterMilli(n.UnixMilli())))
		require.NoError(t, err)
		require.NotNil(t, consumer)

		consumer.Subscribe(topic, topicHandler)

		consumer.Consume()

		assert.Eventuallyf(t, func() bool {
			return taskNum == c.Counter()
		}, 5*time.Second, 50*time.Millisecond,
			"tasks: %d counter: %d",
			taskNum, c.Counter(),
		)

		require.Equal(t, taskNum, c.Counter())

		consumer.Close()
	})

	t.Run("double_consumers_new", func(t *testing.T) {
		n := time.Now()
		for i := 1; i <= int(taskNum); i++ {
			taskName := fmt.Sprintf("Task_%d_new_double", i)
			payload, err := json.Marshal(simpleTaskPayload{
				Name: taskName,
			})
			assert.NoError(t, err)
			err = producer.Send(context.Background(), topic, payload)
			assert.NoError(t, err)
		}

		c.Reset()

		consumer1, err := kafkapool.NewConsumer(kk, kafkapool.WithConsumeResetOffsetConsumerOpt(kgo.NewOffset().AfterMilli(n.UnixMilli())))
		require.NoError(t, err)
		require.NotNil(t, consumer1)

		consumer2, err := kafkapool.NewConsumer(kk, kafkapool.WithConsumeResetOffsetConsumerOpt(kgo.NewOffset().AfterMilli(n.UnixMilli())))
		require.NoError(t, err)
		require.NotNil(t, consumer1)

		consumer1.Subscribe(topic, topicHandler)
		consumer2.Subscribe(topic, topicHandler)

		consumer1.Consume()
		consumer2.Consume()

		assert.Eventuallyf(t, func() bool {
			return taskNum*2 == c.Counter()
		}, 5*time.Second, 50*time.Millisecond,
			"tasks: %d counter: %d",
			taskNum*2, c.Counter(),
		)

		require.Equal(t, taskNum*2, c.Counter())

		consumer1.Close()
		consumer2.Close()
	})

	t.Run("consumer_error", func(t *testing.T) {
		n := time.Now()
		err = producer.Send(context.Background(), topic, []byte("error_payload"))
		require.NoError(t, err)

		c.Reset()

		consumer, err := kafkapool.NewConsumer(kk,
			kafkapool.WithConsumeResetOffsetConsumerOpt(kgo.NewOffset().AfterMilli(n.UnixMilli())),
			kafkapool.WithCustomErrorHandlerConsumerOpt(func(_ error) {
				c.Inc()
			}),
		)
		require.NoError(t, err)
		require.NotNil(t, consumer)

		consumer.Subscribe(topic, func(_ context.Context, payload kafkapool.Payload) error {
			return fmt.Errorf("err")
		})

		consumer.Consume()

		assert.Eventuallyf(t, func() bool {
			return 1 == c.Counter()
		}, 5*time.Second, 50*time.Millisecond,
			"errs: %d counter: %d",
			1, c.Counter(),
		)

		require.Equal(t, int64(1), c.Counter())

		consumer.Close()
	})

	producer.Close(context.Background())

}
