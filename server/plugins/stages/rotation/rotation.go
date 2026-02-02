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
			fx.As(new(plugins.Stage)),
			fx.ResultTags(`group:"stages"`),
		),
	),
)

type Rotattion struct{}

func New() *Rotattion {
	return &Rotattion{}
}

func (r *Rotattion) Name() string {
	return "stages.rotation"
}

func (r *Rotattion) Copy(cfg map[string]any) plugins.Stage {
	return &Rotattion{}
}

func (r *Rotattion) Do(ctx context.Context, state *plugins.State) error {
	winners, ok, err := r.rotate(state.Candidates)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	state.Winners = []ads.Banner{winners}

	return nil
}

func (r *Rotattion) rotate(candidates []ads.Banner) (ads.Banner, bool, error) {
	choices := []weightedrand.Choice[ads.Banner, int]{}
	for _, candidate := range candidates {
		choices = append(choices, weightedrand.NewChoice(candidate, candidate.Price))
	}

	if len(choices) == 0 {
		return ads.Banner{}, false, nil
	}

	chooser, err := weightedrand.NewChooser(choices...)
	if err != nil {
		return ads.Banner{}, false, fmt.Errorf("error on chooser: %w", err)
	}

	return chooser.Pick(), true, nil
}
