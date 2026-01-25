package pipeline

import (
	"context"

	"go.ads.coffee/server/config"
	"go.ads.coffee/server/domain"
	"go.ads.coffee/server/plugins/inputs"
	"go.ads.coffee/server/plugins/outputs"
	"go.ads.coffee/server/plugins/stages"
	"go.ads.coffee/server/plugins/targetings"
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
) *Manager {

	// собираем пайплайн по спецификации в конфиге

	m := &Manager{}
	for _, c := range cfg.Pipelines {
		tt := []domain.Targeting{}

		for _, t := range c.Targetings {
			tt = append(tt, targetings.Get(t.Name, t.Config))
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

		m.pipelines = append(m.pipelines, NewPipeline(
			inputs.Get(c.Input.Name, c.Input.Config),
			outputs.Get(c.Output.Name, c.Output.Config),
			ss,
		))
	}

	return m
}

func (m *Manager) Process(
	ctx context.Context,
	state domain.State,
) {
	for _, p := range m.pipelines {
		if p.Process(ctx, state) {
			return
		}
	}
}
