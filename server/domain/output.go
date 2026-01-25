package domain

import "context"

type Output interface {
	Name() string
	Copy(cfg map[string]any) Output
	Formats(ff []Format)
	Process(ctx context.Context, state *State)
}
