package tracker

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/go-chi/chi/v5"
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
	raw, err := base64.URLEncoding.DecodeString(chi.URLParam(state.Request, "data"))
	if err != nil {
		return false
	}

	info := ads.TrackerInfo{} // TODO: total use ads.Event?
	if err := json.Unmarshal(raw, &info); err != nil {
		return false
	}

	switch info.Action {
	case ads.ActionImpression:
		s.analytics.LogImpression(ctx, info)
	case ads.ActionClick:
		s.analytics.LogClick(ctx, info)
	default:
		return false
	}

	return true
}
