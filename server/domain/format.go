package domain

import (
	"context"
)

type Format interface {
	Name() string
	Copy(cfg map[string]any) Format
	Render(ctx context.Context, state *State)
}
