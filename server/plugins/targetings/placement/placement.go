package placement

import (
	"strconv"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"targetings.placement",
	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Targeting)),
			fx.ResultTags(`group:"targetings"`),
		),
	),
)

// Placement фильтрует кандидатов по таргетингу на конкретный плейсмент.
// Таргетинг хранится в кампании и денормализуется в ads.Banner.Targeting.
type Placement struct{}

func New() *Placement {
	return &Placement{}
}

func (p *Placement) Name() string {
	return "targetings.placement"
}

func (p *Placement) Copy(cfg map[string]any) plugins.Targeting {
	return &Placement{}
}

func (p *Placement) Filter(candidates []ads.Banner, state *plugins.State) []ads.Banner {
	id := strconv.FormatUint(uint64(state.Placement.ID), 10)
	if id == "0" {
		return candidates
	}

	filtered := make([]ads.Banner, 0, len(candidates))

	for _, candidate := range candidates {
		if candidate.Targeting.Placement.Validate([]string{id}) {
			filtered = append(filtered, candidate)
		}
	}

	return filtered
}
