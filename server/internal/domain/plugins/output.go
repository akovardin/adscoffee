package plugins

import "context"

type Output interface {
	Name() string
	Copy(cfg map[string]any) Output
	Formats(ff []Format)
	Do(ctx context.Context, state *State)
}
