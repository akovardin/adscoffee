package pipeline

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/internal/inputs"
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
) *Manager {
	m := &Manager{}
	for _, c := range pipelines {
		tt := []plugins.Targeting{}

		for _, t := range c.Targetings {
			tt = append(tt, targetings.Get(t.Name, t.Config))
		}

		ss := []plugins.Stage{}

		for _, s := range c.Stages {
			v := stages.Get(s.Name, s.Config)
			switch s := v.(type) {
			case plugins.WithTargetings:
				s.Targetings(tt)
			default:
				ss = append(ss, v)
			}
		}

		m.pipelines = append(m.pipelines, NewPipeline(
			c.Name,
			c.Route,
			inputs.Get(c.Input.Name, c.Input.Config),
			outputs.Get(c.Output.Name, c.Output.Config),
			ss,
		))
	}

	return m
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
