package inapp

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"go.ads.coffee/platform/server/internal/analytics"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type Analytics interface {
	LogResponse(ctx context.Context, w ads.Banner, state *plugins.State) error
}

type Inapp struct {
	format    string
	formats   map[string]plugins.Format
	analytics Analytics
	logger    *zap.Logger
}

func New(ff []plugins.Format, logger *zap.Logger, analytics *analytics.Analytics) *Inapp {
	formats := map[string]plugins.Format{}

	for _, f := range ff {
		formats[f.Name()] = f
	}

	return &Inapp{
		analytics: analytics,
		formats:   formats,
		logger:    logger,
	}
}

func (w *Inapp) Name() string {
	return "outputs.inapp"
}

func (w *Inapp) Copy(cfg map[string]any) plugins.Output {
	format := "native"
	if cfg != nil {
		format = cfg["format"].(string)
	}

	dest := make(map[string]plugins.Format, len(w.formats))
	for _, f := range w.formats {
		dest[f.Name()] = f.Copy(cfg)
	}

	return &Inapp{
		analytics: w.analytics,
		format:    format,
		formats:   dest,
		logger:    w.logger,
	}
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