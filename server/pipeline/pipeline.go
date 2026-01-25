package pipeline

import (
	"context"

	"go.ads.coffee/server/domain"
)

type Pipeline struct {
	input  domain.Input
	output domain.Output
	stages []domain.Stage
}

// в конструкторе все возможные плагины
func NewPipeline(
	input domain.Input,
	output domain.Output,
	stages []domain.Stage,
) *Pipeline {
	return &Pipeline{
		input:  input,
		output: output,
		stages: stages,
	}
}

func (p *Pipeline) Process(
	ctx context.Context,
	state domain.State,
) bool {

	// правильно вызываем input -> stages -> output

	return true
}
