package kafkapool

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"

	"go.ads.coffee/platform/pkg/health"
)

type ProducerOpt func(*Producer)

type Producer struct {
	kf *Kafka
	cl *kgo.Client
}

func NewProducer(kf *Kafka) (*Producer, error) {
	p := &Producer{
		kf: kf,
	}

	cl, err := kgo.NewClient(append(kf.opts, kf.producerOpts()...)...)
	if err != nil {
		return nil, err
	}

	p.cl = cl

	return p, nil

}

func (p *Producer) Send(ctx context.Context, topic string, payload Payload) error {
	if !p.kf.Enabled() {
		return nil
	}

	topic = p.kf.formatWithPrefix(topic)

	ctx, span := p.kf.tracer.StartSpan(ctx, "producer_send")
	defer span.End()

	err := p.cl.ProduceSync(
		ctx,
		convertPayloadToRecord(topic, nil, payload),
	).FirstErr()

	p.kf.metrics.ProducerTotal(topic, err)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return fmt.Errorf("produce: %w", err)
	}

	return nil
}

func (p *Producer) SendWithKey(
	ctx context.Context,
	topic string,
	key Key,
	payload Payload,
) error {
	if !p.kf.Enabled() {
		return nil
	}

	topic = p.kf.formatWithPrefix(topic)

	ctx, span := p.kf.tracer.StartSpan(ctx, "producer_send")
	defer span.End()

	err := p.cl.ProduceSync(
		ctx,
		convertPayloadToRecord(topic, key, payload),
	).FirstErr()

	p.kf.metrics.ProducerTotal(topic, err)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return fmt.Errorf("produce: %w", err)
	}

	return nil
}

func (p *Producer) SendBatch(ctx context.Context, topic string, payloads ...Payload) []error {
	if !p.kf.Enabled() {
		return nil
	}

	topic = p.kf.formatWithPrefix(topic)

	ctx, span := p.kf.tracer.StartSpan(ctx, "producer_send_batch")
	defer span.End()

	results := p.cl.ProduceSync(ctx, convertPayloadsToRecords(topic, payloads...)...)
	errors := make([]error, 0, len(payloads))
	for i := range results {
		r := &results[i]

		p.kf.metrics.ProducerTotal(topic, r.Err)

		if r.Err != nil {
			span.RecordError(r.Err)
			errors = append(errors, fmt.Errorf("produce error: %w", r.Err))

			continue
		}
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func (p *Producer) SendAsync(ctx context.Context, topic string, payload Payload) {
	if !p.kf.Enabled() {
		return
	}

	topic = p.kf.formatWithPrefix(topic)
	ctx, span := p.kf.tracer.StartSpan(ctx, "producer_send_async")
	defer span.End()

	p.cl.Produce(ctx,
		convertPayloadToRecord(topic, nil, payload),
		func(_ *kgo.Record, err error) {
			if err != nil {
				span.RecordError(err)
				p.kf.logger.Error("Kafka: async producer error", zap.Error(err), zap.String("topic", topic))
			}
			p.kf.metrics.ProducerTotal(topic, err)
		})
}

func (p *Producer) Ping(ctx context.Context) error {
	if err := p.cl.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping Kafka: %w", err)
	}

	return nil
}

func (p *Producer) Close(_ context.Context) error {
	p.cl.Close()

	return nil
}

func (p *Producer) HealthComponent() *health.Component {
	return &health.Component{
		Kind: health.ComponentKindLocal,
		Name: "kafka_producer",
		CheckFunc: func(ctx context.Context) error {
			return p.Ping(ctx)
		},
	}
}
