package pipeline

import (
	"context"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type Pipeline struct {
	name   string
	route  string
	input  plugins.Input
	output plugins.Output
	stages []plugins.Stage
}

func NewPipeline(
	name string,
	route string,
	input plugins.Input,
	output plugins.Output,
	stages []plugins.Stage,
) *Pipeline {
	return &Pipeline{
		name:   name,
		route:  route,
		input:  input,
		output: output,
		stages: stages,
	}
}

func (p *Pipeline) Name() string {
	return p.name
}

func (p *Pipeline) Route() string {
	return p.route
}

func (p *Pipeline) Do(
	ctx context.Context,
	state *plugins.State,
) error {
	if ok := p.input.Do(ctx, state); !ok {
		return nil
	}

	for _, stage := range p.stages {
		if err := stage.Do(ctx, state); err != nil {
			return err
		}
	}

	return p.output.Do(ctx, state)
}
