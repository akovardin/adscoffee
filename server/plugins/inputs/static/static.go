package static

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"go.ads.coffee/platform/server/internal/analytics"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/internal/repos/banners"
	"go.ads.coffee/platform/server/internal/sessions"
)

const (
	actionClick = "click"
	actionKey   = "action"
)

var Module = fx.Module(
	"inputs.static",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Input)),
			fx.ResultTags(`group:"inputs"`),
		),
	),
)

type Analytics interface {
	LogClick(ctx context.Context)
}

type Static struct {
	logger    *zap.Logger
	cache     *banners.Cache
	sessions  *sessions.Sessions
	analytics Analytics
}

func New(
	logger *zap.Logger,
	cache *banners.Cache,
	sessions *sessions.Sessions,
	analytics *analytics.Analytics,
) *Static {
	return &Static{
		logger:    logger,
		cache:     cache,
		sessions:  sessions,
		analytics: analytics,
	}
}

func (s *Static) Name() string {
	return "inputs.static"
}

func (s *Static) Copy(cfg map[string]any) plugins.Input {
	return &Static{
		cache:     s.cache,
		logger:    s.logger,
		sessions:  s.sessions,
		analytics: s.analytics,
	}
}

func (s *Static) Do(ctx context.Context, state *plugins.State) bool {
	action := chi.URLParam(state.Request, "action")
	state.WithValue(actionKey, action)

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

	// проверяем есть ли в сессии баннер для экшена click
	// если баннер в сессии, то редиректим на трекер url
	if action == actionClick {
		session, ok := s.sessions.LoadWithExpire(state.Request)
		if !ok {
			s.logger.Warn("error on load banner from cache")

			state.Response.WriteHeader(http.StatusNotFound)

			return false
		}

		banner, ok := s.cache.One(ctx, session.Value)
		if !ok {
			s.logger.Warn("error on load banner from cache")

			state.Response.WriteHeader(http.StatusNotFound)

			return false
		}

		s.analytics.LogClick(ctx)

		http.Redirect(state.Response, state.Request, banner.Target, http.StatusSeeOther)

		return false
	}

	// эешен img будет обрабатываться в outputs.static

	return true
}
