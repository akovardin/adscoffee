package pipeline

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/internal/inputs"
	"go.ads.coffee/platform/server/internal/outputs"
	"go.ads.coffee/platform/server/internal/stages"
	"go.ads.coffee/platform/server/internal/targetings"
)

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
		router.Mount(p.Route(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			state := &plugins.State{
				Request:  r,
				Response: w,
			}

			if err := p.Do(ctx, state); err != nil {
				w.WriteHeader(http.StatusNotFound)
			}
		}))
	}
}
