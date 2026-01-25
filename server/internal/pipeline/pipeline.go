package pipeline

import (
	"context"

	"go.ads.coffee/server/domain"
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

func (p *Pipeline) Process(
	ctx context.Context,
	state *domain.State,
) bool {

	p.input.Process(ctx, state)

	for _, stage := range p.stages {
		stage.Process(ctx, state)
	}

	p.output.Process(ctx, state)

	return true
}
