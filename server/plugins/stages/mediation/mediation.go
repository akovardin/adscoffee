package mediation

import (
	"context"
	"fmt"

	"github.com/mroth/weightedrand/v2"
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"stages.mediation",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Stage)),
			fx.ResultTags(`group:"stages"`),
		),
	),
)

type Mediation struct{}

func New() *Mediation {
	return &Mediation{}
}

func (t *Mediation) Name() string {
	return "stages.mediation"
}

func (t *Mediation) Copy(cfg map[string]any) plugins.Stage {
	return &Mediation{}
}

func (t *Mediation) Do(ctx context.Context, state *plugins.State) error {
	winners := state.Winners

	for _, u := range state.Placement.Units {
		winners = append(winners, ads.Banner{
			ID:      u.ID,
			Title:   u.Title,
			Price:   u.Price,
			Type:    ads.CreativeTypeMediator,
			Network: u.Network,
			Format:  u.Format,
		})
	}

	if len(winners) == 1 {
		state.Winners = winners

		return nil
	}

	winner, _, err := t.rotate(winners)
	if err != nil {
		return fmt.Errorf("error on rotate: %w", err)
	}

	state.Winners = []ads.Banner{winner}

	return nil
}

func (t *Mediation) rotate(candidates []ads.Banner) (ads.Banner, bool, error) {
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
