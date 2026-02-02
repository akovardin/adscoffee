package static

import (
	"context"
	"fmt"

	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/internal/sessions"
)

type Static struct {
	formats  map[string]plugins.Format
	sessions *sessions.Sessions
}

func New(
	formats []plugins.Format,
	sessions *sessions.Sessions,
) *Static {
	ff := map[string]plugins.Format{}
	for _, v := range formats {
		ff[v.Name()] = v
	}

	return &Static{
		formats:  ff,
		sessions: sessions,
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
		formats:  ff,
		sessions: w.sessions,
	}
}

func (w *Static) Do(ctx context.Context, state *plugins.State) error {
	action := state.Value("action").(string)
	banner := state.Winners[0]

	if err := w.sessions.Start(state.Request, banner.ID); err != nil {
		return fmt.Errorf("error on start session: %w", err)
	}

	switch action {
	case "img":
		//render img

		state.Response.Write([]byte(banner.Image.Full("")))

	case "click":
		// redirect to url

		state.Response.Write([]byte(banner.Target))
	}

	// data, err := w.formats["banner"].Render(ctx, state)

	// if err != nil {
	// 	return fmt.Errorf("error on render format: %w", err)
	// }

	// _, err = state.Response.Write([]byte(data.(string)))

	return nil
}
