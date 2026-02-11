package web

import (
	"context"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/analytics"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"inputs.web",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Input)),
			fx.ResultTags(`group:"inputs"`),
		),
	),
)

type Analytics interface {
	LogRequest(ctx context.Context, state *plugins.State) error
}

type Web struct {
	analytics Analytics
}

func New(analytics *analytics.Analytics) *Web {
	return &Web{
		analytics: analytics,
	}
}

func (s *Web) Name() string {
	return "inputs.web"
}

func (s *Web) Copy(cfg map[string]any) plugins.Input {
	return &Web{
		analytics: s.analytics,
	}
}

func (s *Web) Do(ctx context.Context, state *plugins.State) bool {
	// нужно получить данные пользователя из запроса

	state.User = &plugins.User{}
	state.Device = &plugins.Device{}

	// проверить наличие placement

	placement := chi.URLParam(state.Request, "placement")

	state.Placement = &plugins.Placement{
		ID: placement,
		Units: []ads.Unit{
			{
				ID:      "yandex-1",
				Network: "yandex",
				Price:   10,
				Format:  "banner",
			},
		},
	}

	// check error
	_ = s.analytics.LogRequest(ctx, state)

	return true
}
