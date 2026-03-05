package web

import (
	"context"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"go.ads.coffee/platform/server/internal/analytics"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/internal/repos/placements"
	"go.ads.coffee/platform/server/internal/repos/units"
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

type Placements interface {
	One(ctx context.Context, id uint) (ads.Placement, bool)
}

type Units interface {
	FindByPlacement(ctx context.Context, id uint) ([]ads.Unit, bool)
}

type Web struct {
	logger     *zap.Logger
	analytics  Analytics
	placements Placements
	units      Units
}

func New(
	logger *zap.Logger,
	analytics *analytics.Analytics,
	placements *placements.Cache,
	units *units.Cache,
) *Web {
	return &Web{
		logger:     logger,
		analytics:  analytics,
		placements: placements,
		units:      units,
	}
}

func (w *Web) Name() string {
	return "inputs.web"
}

func (w *Web) Copy(cfg map[string]any) plugins.Input {
	return &Web{
		logger:     w.logger,
		analytics:  w.analytics,
		placements: w.placements,
		units:      w.units,
	}
}

func (w *Web) Do(ctx context.Context, state *plugins.State) bool {
	// нужно получить данные пользователя из запроса

	state.User = &plugins.User{}
	state.Device = &plugins.Device{}

	// проверить наличие placement
	id, _ := strconv.Atoi(chi.URLParam(state.Request, "placement"))

	placement, exit := w.placements.One(ctx, uint(id))
	if !exit {
		return false
	}

	state.Placement = placement

	units, exit := w.units.FindByPlacement(ctx, placement.ID)
	if exit {
		state.Units = units
	}

	// check error
	if err := w.analytics.LogRequest(ctx, state); err != nil {
		w.logger.Error("error on log request", zap.Error(err))
	}

	return true
}
