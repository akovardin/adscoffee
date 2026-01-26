package pipeline

import (
	"context"

	"go.ads.coffee/platform/server/domain"
)

type Pipeline struct {
	name    string
	route   string
	input   domain.Input
	output  domain.Output
	stages  []domain.Stage
	formats []domain.Format
}

func NewPipeline(
	name string,
	route string,
	input domain.Input,
	output domain.Output,
	stages []domain.Stage,
	formats []domain.Format,
) *Pipeline {
	return &Pipeline{
		name:    name,
		route:   route,
		input:   input,
		output:  output,
		stages:  stages,
		formats: formats,
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
	state *domain.State,
) bool {
	if ok := p.input.Do(ctx, state); !ok {
		return false
	}

	for _, stage := range p.stages {
		stage.Do(ctx, state)
	}

	p.output.Do(ctx, state)

	return true
}
