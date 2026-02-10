package web

import (
	"context"
	"encoding/json"
	"fmt"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type Web struct {
	format  string
	formats map[string]plugins.Format
}

func New(ff []plugins.Format) *Web {
	formats := map[string]plugins.Format{}

	for _, f := range ff {
		formats[f.Name()] = f
	}

	return &Web{
		formats: formats,
	}
}

func (w *Web) Name() string {
	return "outputs.web"
}

func (w *Web) Copy(cfg map[string]any) plugins.Output {
	format := "native" // default format
	if cfg != nil {
		format = cfg["format"].(string)
	}

	dest := make(map[string]plugins.Format, len(w.formats))
	for _, f := range w.formats {
		dest[f.Name()] = f.Copy(cfg)
	}

	return &Web{
		format:  format,
		formats: dest,
	}
}

//nolint:errcheck
func (w *Web) Do(ctx context.Context, state *plugins.State) error {
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

	return err
}
