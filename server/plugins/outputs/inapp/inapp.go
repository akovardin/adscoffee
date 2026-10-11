package inapp

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"go.ads.coffee/platform/server/internal/analytics"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type Analytics interface {
	LogResponse(ctx context.Context, w *ads.Banner, state *plugins.State) error
}

type Config struct {
	Base   string `json:"base"`
	Format string `json:"format"`
}

type Inapp struct {
	format    string
	formats   map[string]plugins.Format
	analytics Analytics
	logger    *zap.Logger
}

func New(
	ff []plugins.Provider[plugins.Format],
	logger *zap.Logger,
	analytics *analytics.Analytics,
) plugins.Provider[plugins.Output] {
	return plugins.Provider[plugins.Output]{
		Name: "outputs.inapp",
		Factory: func(cfg map[string]any) (plugins.Output, error) {
			c, err := plugins.Decode[Config](cfg)
			if err != nil {
				return nil, fmt.Errorf("outputs.inapp config: %w", err)
			}

			formats := make(map[string]plugins.Format, len(ff))
			for _, p := range ff {
				f, err := p.Factory(cfg)
				if err != nil {
					return nil, fmt.Errorf("outputs.inapp format %q: %w", p.Name, err)
				}

				formats[p.Name] = f
			}

			return &Inapp{
				analytics: analytics,
				formats:   formats,
				logger:    logger,
				format:    cmp.Or(c.Format, "native"),
			}, nil
		},
	}
}

func (w *Inapp) Name() string {
	return "outputs.inapp"
}

//nolint:errcheck
func (w *Inapp) Do(ctx context.Context, state *plugins.State) error {
	format, ok := w.formats[w.format]
	if !ok {
		return fmt.Errorf("format %s not found", w.format)
	}

	result, err := format.Render(ctx, state)
	if err != nil {
		return fmt.Errorf("error on render format: %w", err)
	}

	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("error on render format: %w", err)
	}

	_, err = state.Response.Write(data)

	if len(state.Winners) > 0 {
		if err := w.analytics.LogResponse(ctx, state.Winners[0], state); err != nil {
			w.logger.Error("error on log response", zap.Error(err))
		}
	}

	return err
}
