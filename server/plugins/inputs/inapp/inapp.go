package inapp

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
	"inputs.inapp",

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

type Inapp struct {
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
) *Inapp {
	return &Inapp{
		logger:     logger,
		analytics:  analytics,
		placements: placements,
		units:      units,
	}
}

func (w *Inapp) Name() string {
	return "inputs.inapp"
}

func (w *Inapp) Copy(cfg map[string]any) plugins.Input {
	return &Inapp{
		logger:     w.logger,
		analytics:  w.analytics,
		placements: w.placements,
		units:      w.units,
	}
}

func (w *Inapp) Do(ctx context.Context, state *plugins.State) bool {
	state.User = &plugins.User{}
	state.Device = &plugins.Device{}

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

	if err := w.analytics.LogRequest(ctx, state); err != nil {
		w.logger.Error("error on log request", zap.Error(err))
	}

	return true
}