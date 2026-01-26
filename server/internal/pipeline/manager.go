package pipeline

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.ads.coffee/platform/server/config"
	"go.ads.coffee/platform/server/domain"
	"go.ads.coffee/platform/server/internal/formats"
	"go.ads.coffee/platform/server/internal/inputs"
	"go.ads.coffee/platform/server/internal/outputs"
	"go.ads.coffee/platform/server/internal/stages"
	"go.ads.coffee/platform/server/internal/targetings"
)

type Manager struct {
	pipelines []*Pipeline
}

func NewManager(
	cfg config.Config,
	inputs *inputs.Inputs,
	outputs *outputs.Outputs,
	stages *stages.Stages,
	targetings *targetings.Targetings,
	formats *formats.Formats,
) *Manager {
	m := &Manager{}
	for _, c := range cfg.Pipelines {
		tt := []domain.Targeting{}

		for _, t := range c.Targetings {
			tt = append(tt, targetings.Get(t.Name, t.Config))
		}

		ff := []domain.Format{}

		for _, f := range c.Formats {
			ff = append(ff, formats.Get(f.Name, f.Config))
		}

		ss := []domain.Stage{}

		for _, s := range c.Stages {
			v := stages.Get(s.Name, s.Config)
			switch s := v.(type) {
			case domain.WithTargetings:
				s.Targetings(tt)
			default:
				ss = append(ss, v)
			}
		}

		out := outputs.Get(c.Output.Name, c.Output.Config)
		out.Formats(ff)

		m.pipelines = append(m.pipelines, NewPipeline(
			c.Name,
			c.Route,
			inputs.Get(c.Input.Name, c.Input.Config),
			out,
			ss,
			ff,
		))
	}

	return m
}

func (m *Manager) Mount(router *chi.Mux) {
	for _, p := range m.pipelines {
		router.Mount(p.Route(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			state := &domain.State{
				Request:  r,
				Response: w,
			}

			p.Do(ctx, state)
		}))
	}
}
