package kafkapool

import (
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
)

type SASLMechanism string

const (
	SASLMechanismNone        SASLMechanism = ""
	SASLMechanismPlain       SASLMechanism = "PLAIN"
	SASLMechanismScramSHA256 SASLMechanism = "SCRAM-SHA-256"
	SASLMechanismScramSHA512 SASLMechanism = "SCRAM-SHA-512"
)

type Config struct {
	Enabled bool     `yaml:"enabled"`
	Prefix  string   `yaml:"prefix"`
	Seeds   []string `yaml:"seeds"`

	SASLMechanism SASLMechanism `yaml:"sasl_mechanism"`
	Username      string        `yaml:"username"`
	Password      string        `yaml:"password"`
	Certificate   string        `yaml:"certificate"`
	TLS           bool          `yaml:"tls"`

	AllowAutoTopicCreation bool `yaml:"allow_auto_topic_creation"`

	Producer ConfigProducer `yaml:"producer"`
	Consumer ConfigConsumer `yaml:"consumer"`
}

type ConfigProducer struct {
	DisableIdempotentWrite bool          `yaml:"disable_idempotent_write"`
	RequestTimeout         time.Duration `yaml:"request_timeout"`
	RecordRetries          int           `yaml:"record_retries"`
}

type ConfigConsumer struct {
	ConfigBackoff ConfigBackoff `yaml:"backoff"` //nolint:tagliatelle
}

type ConfigBackoff struct {
	Enabled bool `yaml:"enabled"`

	MaxRetries *uint64 `yaml:"max_retries"`

	InitialInterval     *time.Duration `yaml:"initial_interval"`
	MaxInterval         *time.Duration `yaml:"max_interval"`
	MaxElapsedTime      *time.Duration `yaml:"max_elapsed_time"`
	Multiplier          *float64       `yaml:"multiplier"`
	RandomizationFactor *float64       `yaml:"randomization_factor"`
}

func (c *Config) Prepare() *Config {
	if len(c.Seeds) == 1 {
		c.Seeds = strings.Split(c.Seeds[0], ",")
	}

	return c
}

func (c ConfigBackoff) GetBackOff() backoff.BackOff {
	if !c.Enabled {
		return &backoff.StopBackOff{}
	}

	var opts []backoff.ExponentialBackOffOpts

	if c.InitialInterval != nil {
		opts = append(opts, backoff.WithInitialInterval(*c.InitialInterval))
	}

	if c.MaxInterval != nil {
		opts = append(opts, backoff.WithMaxInterval(*c.MaxInterval))
	}

	if c.MaxElapsedTime != nil {
		opts = append(opts, backoff.WithMaxElapsedTime(*c.MaxElapsedTime))
	}

	if c.Multiplier != nil {
		opts = append(opts, backoff.WithMultiplier(*c.Multiplier))
	}

	if c.RandomizationFactor != nil {
		opts = append(opts, backoff.WithRandomizationFactor(*c.RandomizationFactor))
	}

	var bo backoff.BackOff = backoff.NewExponentialBackOff(opts...)

	if c.MaxRetries != nil {
		bo = backoff.WithMaxRetries(bo, *c.MaxRetries)
	}

	return bo
}
