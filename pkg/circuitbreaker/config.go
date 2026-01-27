package circuitbreaker

import (
	"context"
	"errors"
	"time"

	"github.com/cep21/circuit/v4"
	"go.uber.org/zap"
)

const (
	defaultTimeout               = 1 * time.Second
	defaultMaxConcurrentRequests = -1 // unlimited

	defaultHystrixSleepWindow                  = 5 * time.Second
	defaultHystrixHalfOpenAttempts             = 1
	defaultHystrixRequiredConcurrentSuccessful = 1
	defaultHystrixRequestVolumeThreshold       = 20
	defaultHystrixErrorThresholdPercentage     = 50
	defaultHystrixNumBuckets                   = 10
	defaultHystrixRollingDuration              = 10 * time.Second
)

type HystrixConfig struct {
	// Closer

	SleepWindow                  time.Duration `yaml:"sleep_window"`
	HalfOpenAttempts             int64         `yaml:"half_open_attempts"`
	RequiredConcurrentSuccessful int64         `yaml:"required_concurrent_successful"`

	// Opener

	ErrorThresholdPercentage int64         `yaml:"error_threshold_percentage"`
	RequestVolumeThreshold   int64         `yaml:"request_volume_threshold"`
	RollingDuration          time.Duration `yaml:"rolling_duration"`
	NumBuckets               int           `yaml:"num_buckets"`
}

type Config struct {
	Enabled                  bool          `yaml:"enabled"`
	Timeout                  time.Duration `yaml:"timeout"`
	MaxConcurrentRequests    int64         `yaml:"max_concurrent_requests"`
	ContextDeadlineIsAnError bool          `yaml:"context_deadline_is_an_error"`

	Hystrix HystrixConfig `yaml:"hystrix"`
}

func NewConfig() *Config {
	return &Config{
		Enabled:                  true,
		Timeout:                  defaultTimeout,
		MaxConcurrentRequests:    defaultMaxConcurrentRequests,
		ContextDeadlineIsAnError: false,
		Hystrix: HystrixConfig{
			SleepWindow:                  defaultHystrixSleepWindow,
			HalfOpenAttempts:             defaultHystrixHalfOpenAttempts,
			RequiredConcurrentSuccessful: defaultHystrixRequiredConcurrentSuccessful,
			ErrorThresholdPercentage:     defaultHystrixErrorThresholdPercentage,
			RequestVolumeThreshold:       defaultHystrixRequestVolumeThreshold,
			RollingDuration:              defaultHystrixRollingDuration,
			NumBuckets:                   defaultHystrixNumBuckets,
		},
	}
}

func (c *Config) ToCircuitBreakerConfig(
	logger *zap.Logger,
	name string,
	metrics *Metrics,
) circuit.Config {

	var (
		ignoreInterrupts bool
		isErrInterrupt   func(err error) bool
	)

	if c.ContextDeadlineIsAnError {
		ignoreInterrupts = true
		isErrInterrupt = func(err error) bool {
			return !errors.Is(err, context.DeadlineExceeded)
		}
	}

	return circuit.Config{
		General: circuit.GeneralConfig{
			Disabled: !c.Enabled,
			GoLostErrors: func(err error, pan interface{}) {
				switch {
				case err != nil:
					logger.Error("lost error", zap.Error(err))
				case pan != nil:
					logger.Error("lost panic", zap.Any("panic", pan))
				}
			},
		},
		Execution: circuit.ExecutionConfig{
			Timeout:               c.Timeout,
			MaxConcurrentRequests: c.MaxConcurrentRequests,
			IgnoreInterrupts:      ignoreInterrupts,
			IsErrInterrupt:        isErrInterrupt,
		},
		Fallback: circuit.FallbackConfig{
			Disabled:              false,
			MaxConcurrentRequests: c.MaxConcurrentRequests,
		},
		Metrics: circuit.MetricsCollectors{
			Run: []circuit.RunMetrics{
				NewRunMetrics(metrics, name),
			},
			Fallback: []circuit.FallbackMetrics{
				NewFallbackMetrics(metrics, name),
			},
			Circuit: []circuit.Metrics{
				NewCircuitMetrics(metrics, name),
			},
		},
	}
}
