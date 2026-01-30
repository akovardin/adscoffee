package plugins

import (
	"context"
)

type Format interface {
	Name() string
	Render(ctx context.Context, state *State) (any, error)
}
