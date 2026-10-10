package plugins

import "context"

type Input interface {
	Name() string
	Do(ctx context.Context, state *State) bool
}
