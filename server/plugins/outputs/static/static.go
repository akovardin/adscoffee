package static

import (
	"context"
	"fmt"
	"net/http"

	"go.ads.coffee/platform/server/internal/analytics"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/internal/sessions"
	"go.ads.coffee/platform/server/plugins/outputs/static/formats"
	"go.uber.org/fx"
)

const (
	baseUrlKey = "base"
	actionImg  = "img"
	actionKey  = "action"
)

var Module = fx.Module(
	"outputs.static",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Output)),
			fx.ResultTags(`group:"outputs"`),
		),
		formats.NewBanner,
	),
)

type Analytics interface {
	LogImpression(ctx context.Context, data ads.TrackerInfo) error
}

type Banner interface {
	Banner(ctx context.Context, base string, banner ads.Banner, w http.ResponseWriter) error
}

type Session interface {
	Start(r *http.Request, value string) error
}

type Static struct {
	base      string
	sessions  Session
	analytics Analytics
	format    Banner
}

func New(
	format *formats.Banner,
	sessions *sessions.Sessions,
	analytics *analytics.Analytics,
) *Static {
	return &Static{
		format:    format,
		sessions:  sessions,
		analytics: analytics,
	}
}

func (w *Static) Name() string {
	return "outputs.static"
}

func (w *Static) Copy(cfg map[string]any) plugins.Output {
	base := ""

	if cfg != nil {
		base = cfg[baseUrlKey].(string)
	}

	return &Static{
		base:      base,
		format:    w.format,
		sessions:  w.sessions,
		analytics: w.analytics,
	}
}

func (w *Static) Do(ctx context.Context, state *plugins.State) error {
	action := state.Value(actionKey).(string)

	// сюда мы попадаем только для экшена img
	if action != actionImg {
		state.Response.WriteHeader(http.StatusNotFound)

		return nil
	}

	if len(state.Winners) == 0 {
		state.Response.WriteHeader(http.StatusNotFound)

		return nil
	}

	banner := state.Winners[0]

	if err := w.sessions.Start(state.Request, banner.ID); err != nil {
		return fmt.Errorf("error on start session: %w", err)
	}

	// check error
	_ = w.analytics.LogImpression(ctx, ads.TrackerInfo{})

	return w.format.Banner(ctx, w.base, banner, state.Response)
}
