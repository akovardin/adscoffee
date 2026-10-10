package pipeline

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/internal/inputs"
	"go.ads.coffee/platform/server/internal/limits"
	"go.ads.coffee/platform/server/internal/metrics"
	"go.ads.coffee/platform/server/internal/outputs"
	"go.ads.coffee/platform/server/internal/stages"
	"go.ads.coffee/platform/server/internal/targetings"
)

var tracer = otel.Tracer("server")

type Manager struct {
	pipelines []*Pipeline
}

func NewManager(
	pipelines []Config,
	inputs *inputs.Inputs,
	outputs *outputs.Outputs,
	stages *stages.Stages,
	targetings *targetings.Targetings,
	limits *limits.Limits,
) (*Manager, error) {
	m := &Manager{}
	for _, c := range pipelines {
		tt := []plugins.Targeting{}

		for _, t := range c.Targetings {
			v, err := targetings.Get(t.Name, t.Config)
			if err != nil {
				return nil, fmt.Errorf("pipeline %q targeting: %w", c.Name, err)
			}
			tt = append(tt, v)
		}

		ll := []plugins.Limit{}

		for _, l := range c.Limits {
			v, err := limits.Get(l.Name, l.Config)
			if err != nil {
				return nil, fmt.Errorf("pipeline %q limit: %w", c.Name, err)
			}
			ll = append(ll, v)
		}

		in, err := inputs.Get(c.Input.Name, c.Input.Config)
		if err != nil {
			return nil, fmt.Errorf("pipeline %q input: %w", c.Name, err)
		}

		out, err := outputs.Get(c.Output.Name, c.Output.Config)
		if err != nil {
			return nil, fmt.Errorf("pipeline %q output: %w", c.Name, err)
		}

		ss := []plugins.Stage{}

		for _, s := range c.Stages {
			v, err := stages.Get(s.Name, s.Config)
			if err != nil {
				return nil, fmt.Errorf("pipeline %q stage: %w", c.Name, err)
			}

			if w, ok := v.(plugins.WithTargetings); ok {
				w.Targetings(tt)
			}
			if w, ok := v.(plugins.WithLimits); ok {
				w.Limits(ll)
			}

			ss = append(ss, v)
		}

		m.pipelines = append(m.pipelines, NewPipeline(
			c.Name,
			c.Route,
			in,
			out,
			ss,
		))
	}

	return m, nil
}

func (m *Manager) Mount(router *chi.Mux) {
	for _, p := range m.pipelines {
		pipeline := p

		router.Mount(pipeline.Route(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// извлекаем контекст трассировки из заголовков и стартуем спан пайплайна
			ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			ctx, span := tracer.Start(ctx, "pipeline "+pipeline.Name())
			defer span.End()

			// счётчик RPS на уровне пайплайна (server_requests_total{kind="pipeline"})
			metrics.IncPipeline(pipeline.Name())

			state := &plugins.State{
				RequestID: uuid.NewString(),
				ClickID:   uuid.NewString(),
				Request:   r.WithContext(ctx),
				Response:  w,
			}

			err := pipeline.Do(ctx, state)

			// время ответа сервиса
			metrics.ObserveDuration(metrics.KindPipeline, pipeline.Name(), time.Since(start))

			if err != nil {
				// счётчик ошибок на уровне пайплайна
				metrics.IncError(metrics.KindPipeline, pipeline.Name())
				span.RecordError(err)
				w.WriteHeader(http.StatusNotFound)
			}
		}))
	}
}
