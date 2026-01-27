package circuitbreaker

import (
	"github.com/cep21/circuit/v4"
	"github.com/cep21/circuit/v4/closers/hystrix"
	"go.uber.org/zap"
)

type Pool struct {
	logger  *zap.Logger
	manager *circuit.Manager
}

func NewPool(logger *zap.Logger, metrics *Metrics, configs map[string]*Config) (*Pool, error) {
	hystrixConf := hystrix.Factory{
		CreateConfigureCloser: []func(circuitName string) hystrix.ConfigureCloser{
			func(circuitName string) hystrix.ConfigureCloser {
				c, ok := configs[circuitName]
				if !ok {
					return hystrix.ConfigureCloser{}
				}

				return hystrix.ConfigureCloser{
					SleepWindow:                  c.Hystrix.SleepWindow,
					HalfOpenAttempts:             c.Hystrix.HalfOpenAttempts,
					RequiredConcurrentSuccessful: c.Hystrix.RequiredConcurrentSuccessful,
				}
			},
		},

		CreateConfigureOpener: []func(circuitName string) hystrix.ConfigureOpener{
			func(circuitName string) hystrix.ConfigureOpener {
				c, ok := configs[circuitName]
				if !ok {
					return hystrix.ConfigureOpener{}
				}

				return hystrix.ConfigureOpener{
					ErrorThresholdPercentage: c.Hystrix.ErrorThresholdPercentage,
					RequestVolumeThreshold:   c.Hystrix.RequestVolumeThreshold,
					RollingDuration:          c.Hystrix.RollingDuration,
					NumBuckets:               c.Hystrix.NumBuckets,
				}
			},
		},
	}

	manager := &circuit.Manager{
		DefaultCircuitProperties: []circuit.CommandPropertiesConstructor{
			hystrixConf.Configure,
			func(circuitName string) circuit.Config {
				c, ok := configs[circuitName]
				if !ok {
					return circuit.Config{}
				}

				return c.ToCircuitBreakerConfig(logger.Named(circuitName), circuitName, metrics)
			},
		},
	}

	for name := range configs {
		if _, err := manager.CreateCircuit(name); err != nil {
			return nil, err
		}
	}

	return &Pool{
		logger:  logger,
		manager: manager,
	}, nil
}

func (p *Pool) Get(name string) Circuit {
	if p == nil {
		return &noopCircuit{}
	}

	c := p.manager.GetCircuit(name)
	if c == nil {
		return &noopCircuit{}
	}

	return c
}
