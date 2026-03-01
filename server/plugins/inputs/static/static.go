package static

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"go.ads.coffee/platform/server/internal/analytics"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/internal/repos/banners"
	"go.ads.coffee/platform/server/internal/repos/placements"
	"go.ads.coffee/platform/server/internal/repos/units"
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
	LogClick(ctx context.Context, data ads.TrackerInfo) error
}

type Cache interface {
	One(ctx context.Context, id uint) (ads.Banner, bool)
}

type Session interface {
	LoadWithExpire(r *http.Request) (sessions.Session, bool)
}

type Static struct {
	logger     *zap.Logger
	cache      Cache
	sessions   Session
	analytics  Analytics
	placements *placements.Cache
	units      *units.Cache
}

func New(
	logger *zap.Logger,
	cache *banners.Cache,
	sessions *sessions.Sessions,
	analytics *analytics.Analytics,
	placements *placements.Cache,
	units *units.Cache,
) *Static {
	return &Static{
		logger:     logger,
		cache:      cache,
		sessions:   sessions,
		analytics:  analytics,
		placements: placements,
		units:      units,
	}
}

func (s *Static) Name() string {
	return "inputs.static"
}

func (s *Static) Copy(cfg map[string]any) plugins.Input {
	return &Static{
		cache:      s.cache,
		logger:     s.logger,
		sessions:   s.sessions,
		analytics:  s.analytics,
		placements: s.placements,
		units:      s.units,
	}
}

func (s *Static) Do(ctx context.Context, state *plugins.State) bool {
	action := chi.URLParam(state.Request, "action")
	state.WithValue(actionKey, action)

	// нужно получить данные пользователя из запроса

	state.User = &plugins.User{}
	state.Device = &plugins.Device{}

	// проверить наличие placement

	id, _ := strconv.Atoi(chi.URLParam(state.Request, "placement"))

	placement, exit := s.placements.One(ctx, uint(id))
	if !exit {
		return false
	}

	state.Placement = placement

	units, exit := s.units.FindByPlacement(ctx, placement.ID)
	if exit {
		state.Units = units
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

		id, _ := strconv.Atoi(session.Value)
		banner, ok := s.cache.One(ctx, uint(id))
		if !ok {
			s.logger.Warn("error on load banner from cache")

			state.Response.WriteHeader(http.StatusNotFound)

			return false
		}

		if err := s.analytics.LogClick(ctx, ads.TrackerInfo{}); err != nil {
			s.logger.Error("error on log click", zap.Error(err))
		}

		http.Redirect(state.Response, state.Request, banner.Target, http.StatusSeeOther)

		return false
	}

	return true
}
