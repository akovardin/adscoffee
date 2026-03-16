package tracker

import (
	"context"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/analytics"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"inputs.tracker",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Input)),
			fx.ResultTags(`group:"inputs"`),
		),
	),
)

type Analytics interface {
	LogImpression(ctx context.Context, data ads.TrackerInfo) error
	LogClick(ctx context.Context, data ads.TrackerInfo) error
}

type Tracker struct {
	analytics Analytics
}

func New(analytics *analytics.Analytics) *Tracker {
	return &Tracker{
		analytics: analytics,
	}
}

func (s *Tracker) Name() string {
	return "inputs.tracker"
}

func (s *Tracker) Copy(cfg map[string]any) plugins.Input {
	return &Tracker{
		analytics: s.analytics,
	}
}

func (s *Tracker) Do(ctx context.Context, state *plugins.State) bool {
	// нужно получить данные пользователя из запроса

	// decode from url base 64
	// save impression/click
	// save unic
	// log to analytics

	// s.analytics.LogImpression()

	state.User = &plugins.User{}
	state.Device = &plugins.Device{}

	return true
}
