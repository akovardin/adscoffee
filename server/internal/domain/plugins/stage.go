package plugins

import "context"

type Stage interface {
	Name() string
	Copy(cfg map[string]any) Stage
	Do(ctx context.Context, state *State) error
}

type WithTargetings interface {
	Targetings(tt []Targeting)
}
