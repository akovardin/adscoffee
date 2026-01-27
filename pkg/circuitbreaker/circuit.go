package circuitbreaker

import "context"

type Circuit interface {
	Go(
		ctx context.Context,
		runFunc func(context.Context) error,
		fallbackFunc func(context.Context, error) error,
	) error
	Run(ctx context.Context, runFunc func(context.Context) error) error
	Execute(
		ctx context.Context,
		runFunc func(context.Context) error,
		fallbackFunc func(context.Context, error) error,
	) error
}

type noopCircuit struct{}

func (c *noopCircuit) Go(
	ctx context.Context,
	runFunc func(context.Context) error,
	_ func(context.Context, error) error,
) error {
	return runFunc(ctx)
}

func (c *noopCircuit) Run(ctx context.Context, runFunc func(context.Context) error) error {
	return runFunc(ctx)
}

func (c *noopCircuit) Execute(
	ctx context.Context,
	runFunc func(context.Context) error,
	_ func(context.Context, error) error,
) error {
	return runFunc(ctx)
}
