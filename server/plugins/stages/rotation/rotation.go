package rotation

import (
	"context"
	"fmt"

	"github.com/mroth/weightedrand/v2"
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"stages.rotation",

	fx.Provide(
		fx.Annotate(
			New,
			fx.ResultTags(`group:"stages"`),
		),
	),
)

type Rotattion struct{}

func New() plugins.Provider[plugins.Stage] {
	return plugins.Provider[plugins.Stage]{
		Name: "stages.rotation",
		Factory: func(cfg map[string]any) (plugins.Stage, error) {
			return &Rotattion{}, nil
		},
	}
}

func (r *Rotattion) Name() string {
	return "stages.rotation"
}

func (r *Rotattion) Do(ctx context.Context, state *plugins.State) error {
	winner, ok, err := r.rotate(state.Candidates)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	state.Winners = []*ads.Banner{winner}

	return nil
}

func (r *Rotattion) rotate(candidates []*ads.Banner) (*ads.Banner, bool, error) {
	choices := []weightedrand.Choice[*ads.Banner, int]{}
	for _, candidate := range candidates {
		choices = append(choices, weightedrand.NewChoice(candidate, candidate.Price))
	}

	if len(choices) == 0 {
		return nil, false, nil
	}

	chooser, err := weightedrand.NewChooser(choices...)
	if err != nil {
		return nil, false, fmt.Errorf("error on chooser: %w", err)
	}

	return chooser.Pick(), true, nil
}
