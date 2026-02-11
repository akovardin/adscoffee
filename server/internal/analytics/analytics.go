package analytics

import (
	"context"
	"fmt"

	"go.ads.coffee/platform/pkg/kafkapool"
	"go.ads.coffee/platform/pkg/telemetry"
	"go.ads.coffee/platform/server/internal/domain/ads"
)

type Analytics struct {
	tel      *telemetry.Telemetry
	producer *kafkapool.Producer
}

func New(pool *kafkapool.Pool, tel *telemetry.Telemetry) (*Analytics, error) {
	if err := tel.Register(actions, money); err != nil {
		return nil, err
	}

	kfk, err := pool.GetPool("main")
	if err != nil {
		return nil, fmt.Errorf("kafka pool error: %w", err)
	}

	producer, err := kafkapool.NewProducer(kfk)
	if err != nil {
		return nil, fmt.Errorf("kafka producer error: %w", err)
	}

	return &Analytics{
		tel:      tel,
		producer: producer,
	}, nil
}

func (r *Analytics) Log(ctx context.Context, name string, event ads.Event) error {
	actions.WithLabelValues(
		name,
		event.Network,
	).Inc()

	data, err := event.JSON()
	if err != nil {
		return err
	}
	// write to kafka
	r.producer.SendAsync(
		context.WithoutCancel(ctx), // нужно проверить на отмену по таймауту
		name,
		data,
	)

	return nil
}

func (a *Analytics) LogImpression(ctx context.Context) {

}

func (a *Analytics) LogClick(ctx context.Context) {

}
