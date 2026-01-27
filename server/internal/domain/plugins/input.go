package plugins

import "context"

type Input interface {
	Name() string
	Copy(cfg map[string]any) Input
	Do(ctx context.Context, state *State) bool
}
