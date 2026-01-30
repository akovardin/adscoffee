package plugins

import "context"

type Output interface {
	Name() string
	Copy(cfg map[string]any) Output
	Do(ctx context.Context, state *State) error
}
