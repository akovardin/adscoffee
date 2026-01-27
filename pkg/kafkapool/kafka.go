package kafkapool

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"github.com/twmb/franz-go/plugin/kotel"
	"github.com/twmb/franz-go/plugin/kzap"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Key []byte
type Payload []byte

type ConsumerCallback func(ctx context.Context, payload Payload) error

type Tracer interface {
	StartSpan(
		ctx context.Context,
		spanName string,
		opts ...trace.SpanStartOption,
	) (context.Context, trace.Span)

	TracerProvider() trace.TracerProvider
}

type MetricsProvider interface {
	ProducerTotal(topic string, err error)
	ConsumerHandleDuration(topic string, err error, seconds float64)
	ConsumerGroupHandleDuration(group, topic string, err error, seconds float64)
}

// Kafka contains all things for default connection, metrics, and logger.
type Kafka struct {
	logger *zap.Logger
	cfg    *Config

	metrics MetricsProvider
	tracer  Tracer
	kTracer *kotel.Tracer

	opts []kgo.Opt
}

// nolint:cyclop
func NewKafka(
	logger *zap.Logger,
	cfg *Config,
	metrics MetricsProvider,
	tracer Tracer,
) (*Kafka, error) {
	if cfg == nil {
		return nil, fmt.Errorf("empty config")
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Seeds...),
		kgo.WithLogger(kzap.New(logger, kzap.Level(kgo.LogLevelWarn))),
	}

	logger.Info("kafka mechanism config", zap.String("mechanism", string(cfg.SASLMechanism)))

	switch cfg.SASLMechanism {
	case SASLMechanismNone:
		// do nothing
	case SASLMechanismPlain:
		opts = append(opts, kgo.SASL(plain.Auth{
			User: cfg.Username,
			Pass: cfg.Password,
		}.AsMechanism()))
	case SASLMechanismScramSHA256:
		opts = append(opts, kgo.SASL(scram.Auth{
			User: cfg.Username,
			Pass: cfg.Password,
		}.AsSha256Mechanism()))
	case SASLMechanismScramSHA512:
		opts = append(opts, kgo.SASL(scram.Auth{
			User: cfg.Username,
			Pass: cfg.Password,
		}.AsSha512Mechanism()))
	default:
		return nil, fmt.Errorf("unknown sasl mechanism: %s", cfg.SASLMechanism)
	}

	if cfg.TLS {
		opts = append(opts, kgo.DialTLS())
	}
	if cfg.Certificate != "" {
		caCert, err := os.ReadFile(cfg.Certificate)
		if err != nil {
			return nil, fmt.Errorf("read certificate: %w", err)
		}

		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(caCert)
		if !ok {
			return nil, fmt.Errorf("append certs from pem")
		}

		opts = append(opts, kgo.DialTLSConfig(
			//nolint: gosec
			&tls.Config{
				MinVersion: tls.VersionTLS10,
				RootCAs:    caCertPool,
			}))
	}

	if cfg.AllowAutoTopicCreation {
		opts = append(opts, kgo.AllowAutoTopicCreation())
	}

	kTracer := kotel.NewTracer(
		kotel.TracerProvider(tracer.TracerProvider()),
		kotel.TracerPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})),
	)
	kotelOps := []kotel.Opt{
		kotel.WithTracer(kTracer),
	}
	kotelService := kotel.NewKotel(kotelOps...)
	opts = append(opts, kgo.WithHooks(kotelService.Hooks()...))

	return &Kafka{
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
		kTracer: kTracer,
		cfg:     cfg,
		opts:    opts,
	}, nil
}

func (k *Kafka) Enabled() bool {
	return k.cfg.Enabled
}

func (k *Kafka) formatWithPrefix(str string) string {
	if k.cfg.Prefix != "" {
		return fmt.Sprintf("%s_%s", k.cfg.Prefix, str)
	}

	return str
}

func (k *Kafka) producerOpts() []kgo.Opt {
	res := make([]kgo.Opt, 0)
	pCfg := k.cfg.Producer

	if pCfg.DisableIdempotentWrite {
		res = append(res, kgo.DisableIdempotentWrite())
	}
	if pCfg.RequestTimeout > 0 {
		res = append(res, kgo.ProduceRequestTimeout(pCfg.RequestTimeout))
	}
	if pCfg.RecordRetries > 0 {
		res = append(res, kgo.RecordRetries(pCfg.RecordRetries))
	}

	return res

}
func convertPayloadToRecord(
	topic string,
	key Key,
	payload Payload,
) *kgo.Record {
	return &kgo.Record{
		Topic: topic,
		Key:   key,
		Value: payload,
	}
}

func convertPayloadsToRecords(
	topic string,
	payload ...Payload,
) []*kgo.Record {
	res := make([]*kgo.Record, 0, len(payload))
	for _, p := range payload {
		res = append(res, convertPayloadToRecord(topic, nil, p))
	}

	return res
}

func convertRecordToPayload(r *kgo.Record) Payload {
	return r.Value
}
