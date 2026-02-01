package static

import (
	"context"
	"fmt"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type Static struct {
	formats map[string]plugins.Format
}

func New(formats []plugins.Format) *Static {
	ff := map[string]plugins.Format{}
	for _, v := range formats {
		ff[v.Name()] = v
	}

	return &Static{
		formats: ff,
	}
}

func (w *Static) Name() string {
	return "outputs.static"
}

func (w *Static) Copy(cfg map[string]any) plugins.Output {
	ff := map[string]plugins.Format{}
	for k, f := range w.formats {
		ff[k] = f.Copy(cfg)
	}

	return &Static{
		formats: ff,
	}
}

func (w *Static) Do(ctx context.Context, state *plugins.State) error {
	data, err := w.formats["banner"].Render(ctx, state)

	if err != nil {
		return fmt.Errorf("error on render format: %w", err)
	}

	_, err = state.Response.Write([]byte(data.(string)))

	return err
}
