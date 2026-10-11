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
			fx.ResultTags(`group:"targetings"`),
		),
	),
)

// Placement фильтрует кандидатов по таргетингу на конкретный плейсмент.
// Таргетинг хранится в кампании и денормализуется в ads.Banner.Targeting.
type Placement struct{}

func New() plugins.Provider[plugins.Targeting] {
	return plugins.Provider[plugins.Targeting]{
		Name: "targetings.placement",
		Factory: func(cfg map[string]any) (plugins.Targeting, error) {
			return &Placement{}, nil
		},
	}
}

func (p *Placement) Name() string {
	return "targetings.placement"
}

func (p *Placement) Filter(candidates []*ads.Banner, state *plugins.State) []*ads.Banner {
	id := strconv.FormatUint(uint64(state.Placement.ID), 10)
	if id == "0" {
		return candidates
	}

	// in-place компакция: пишем выживших в тот же буфер, без аллокаций
	kept := candidates[:0]

	for _, candidate := range candidates {
		if candidate.Targeting.Placement.Validate([]string{id}) {
			kept = append(kept, candidate)
		}
	}

	return kept
}
