package web

import (
	"context"
	"encoding/json"
	"fmt"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type Web struct {
	formats []plugins.Format
}

func New(formats []plugins.Format) *Web {
	return &Web{
		formats: formats,
	}
}

func (w *Web) Name() string {
	return "outputs.web"
}

func (w *Web) Copy(cfg map[string]any) plugins.Output {
	return &Web{
		formats: w.formats,
	}
}

//nolint:errcheck
func (w *Web) Do(ctx context.Context, state *plugins.State) error {
	result := map[string]any{}

	for _, f := range w.formats {
		val, err := f.Render(ctx, state)
		if err != nil {
			return fmt.Errorf("error on render format: %w", err)
		}

		result[f.Name()] = val
	}

	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("error on render format: %w", err)
	}

	_, err = state.Response.Write(data)

	return err
}
