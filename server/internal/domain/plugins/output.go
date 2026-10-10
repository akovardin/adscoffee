package plugins

import "context"

type Output interface {
	Name() string
	Do(ctx context.Context, state *State) error
}
