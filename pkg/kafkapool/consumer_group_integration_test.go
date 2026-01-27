//go:build integration
// +build integration

package kafkapool_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.ads.coffee/platform/pkg/kafkapool"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestIntegration_ConsumerGroup(t *testing.T) {
	const (
		taskNum             int64 = 30
		producerMaxDuration       = 5 * time.Second
		consumerMaxDuration       = 5 * time.Second
	)

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

	t.Run("consuming", func(t *testing.T) {
		topic := "test_group_topic"
		c := &counter{}
		topicHandler := func(_ context.Context, p kafkapool.Payload) error {
			c.Inc()

			return nil
		}

		for i := 1; i <= int(taskNum); i++ {
			taskName := fmt.Sprintf("Task_%d", i)
			payload, err := json.Marshal(simpleTaskPayload{
				Name: taskName,
			})

			require.NoError(t, err)
			err = producer.Send(context.Background(), topic, payload)
			require.NoError(t, err)
		}

		consumer, err := kafkapool.NewConsumerGroup(kk, "test_group_1")
		require.NoError(t, err)
		require.NotNil(t, consumer)

		consumer.Subscribe(topic, topicHandler)

		consumer.Consume(10)

		assert.Eventuallyf(t, func() bool {
			return taskNum == c.Counter()
		}, 10*time.Second, 50*time.Millisecond,
			"tasks: %d counter: %d",
			taskNum, c.Counter(),
		)

		require.Equal(t, taskNum, c.Counter())
		consumer.Close()
	})

	t.Run("consuming_two_consumers", func(t *testing.T) {
		topic := "test_group_topic_2"

		c1 := &counter{}
		topicHandler1 := func(_ context.Context, p kafkapool.Payload) error {
			c1.Inc()

			return nil
		}
		c2 := &counter{}
		topicHandler2 := func(_ context.Context, p kafkapool.Payload) error {
			c2.Inc()

			return nil
		}

		consumer1, err := kafkapool.NewConsumerGroup(kk, "test_group_2")
		require.NoError(t, err)
		require.NotNil(t, consumer1)
		consumer2, err := kafkapool.NewConsumerGroup(kk, "test_group_2")
		require.NoError(t, err)
		require.NotNil(t, consumer2)

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			consumer1.Subscribe(topic, topicHandler1)

			consumer1.Consume(int(taskNum))
			wg.Done()
		}()

		go func() {
			consumer2.Subscribe(topic, topicHandler2)

			consumer2.Consume(int(taskNum))
			wg.Done()
		}()

		wg.Wait()

		for i := 1; i <= int(taskNum); i++ {
			taskName := fmt.Sprintf("Task_%d", i)
			payload, err := json.Marshal(simpleTaskPayload{
				Name: taskName,
			})
			require.NoError(t, err)
			producer.Send(context.Background(), topic, payload)
		}

		assert.Eventuallyf(t, func() bool {
			return taskNum == c1.Counter()+c2.Counter()
		}, 15*time.Second, 50*time.Millisecond,
			"tasks: %d counter: %d",
			taskNum, c1.Counter()+c2.Counter(),
		)

		require.Equal(t, taskNum, c1.Counter()+c2.Counter())
		logger.Info("counters", zap.Int64("counter1", c1.Counter()), zap.Int64("counter2", c2.Counter()))
		consumer1.Close()
		consumer2.Close()
	})

	producer.Close(context.Background())
}
