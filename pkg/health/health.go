package health

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

const (
	AppComponentName = "application"
)

type CheckCallback func(component *Component, err error, duration time.Duration)

type Health struct {
	logger *zap.Logger
	config Config

	isReady atomic.Value // global readiness flag, initial is false

	components         []*Component
	componentProviders []ComponentProvider

	m sync.RWMutex
}

func NewHealth(p Params) *Health {
	h := &Health{
		logger:     p.Logger,
		config:     p.Config,
		components: []*Component{},
	}

	h.isReady.Store(false)

	h.components = append(h.components, &Component{
		Kind: ComponentKindApp,
		Name: AppComponentName,
		CheckFunc: func(context.Context) error {
			// nolint:errcheck
			if isReady, _ := h.isReady.Load().(bool); !isReady {
				return ErrApplicationIsNotReady
			}

			return nil
		},
	})

	h.components = append(h.components, p.Components...)

	h.componentProviders = p.ComponentProviders

	return h
}

func (h *Health) Config() Config {
	return h.config
}

func (h *Health) SetReady(ready bool) {
	h.isReady.Store(ready)
}

func (h *Health) Iter(requestedKind ComponentKind, callback func(*Component)) {
	if callback == nil {
		return
	}

	h.m.RLock()
	defer h.m.RUnlock()

	for _, component := range h.components {
		if component.Kind&requestedKind == 0 {
			continue
		}

		callback(component)
	}

	for _, provider := range h.componentProviders {
		for _, component := range provider.HealthComponents() {
			if component.Kind&requestedKind == 0 {
				continue
			}

			callback(component)
		}
	}
}

func (h *Health) Check(ctx context.Context, requestedKind ComponentKind) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second) // 3 seconds for check
	defer cancel()

	h.m.Lock()
	defer h.m.Unlock()

	for _, component := range h.components {
		if component.Kind&requestedKind == 0 {
			continue
		}

		component.Check(ctx)

		if component.CheckErr != nil {
			h.logger.Error(
				"health check failed",
				zap.String("component", component.Name),
				zap.Error(component.CheckErr),
			)
		}
	}

	for _, provider := range h.componentProviders {
		for _, component := range provider.HealthComponents() {
			if component.Kind&requestedKind == 0 {
				continue
			}

			component.Check(ctx)

			if component.CheckErr != nil {
				h.logger.Error(
					"health check failed",
					zap.String("component", component.Name),
					zap.Error(component.CheckErr),
				)
			}
		}
	}
}
