package static

import (
	"context"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
	"go.ads.coffee/platform/server/internal/repos/banners"
	"go.ads.coffee/platform/server/internal/sessions"
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

type Static struct {
	logger   *zap.Logger
	cache    *banners.Cache
	sessions *sessions.Sessions
}

func New(
	logger *zap.Logger,
	cache *banners.Cache,
	sessions *sessions.Sessions,
) *Static {
	return &Static{
		logger:   logger,
		cache:    cache,
		sessions: sessions,
	}
}

func (s *Static) Name() string {
	return "inputs.static"
}

func (s *Static) Copy(cfg map[string]any) plugins.Input {
	return &Static{
		cache:    s.cache,
		logger:   s.logger,
		sessions: s.sessions,
	}
}

func (s *Static) Do(ctx context.Context, state *plugins.State) bool {
	action := chi.URLParam(state.Request, "action")
	state.WithValue("action", action)

	if session, ok := s.sessions.LoadWithExpire(state.Request); ok {
		// banner := session.Value

		banner, ok := s.cache.One(ctx, session.Value)
		if !ok {
			s.logger.Warn("error on load banner from cache")

			return false
		}

		switch action {
		case "img":
			//render img

			state.Response.Write([]byte(banner.Image.Full("")))

		case "click":
			// redirect to url

			state.Response.Write([]byte(banner.Target))
		}

		return false // дальше не идем
	}

	// нужно получить данные пользователя из запроса

	state.User = &plugins.User{}
	state.Device = &plugins.Device{}

	// проверить наличие placement

	placement := chi.URLParam(state.Request, "placement")

	state.Placement = &plugins.Placement{
		ID: placement,
		Units: []ads.Unit{
			{
				ID:      "yandex-1",
				Network: "yandex",
				Price:   10,
				Format:  "banner",
			},
		},
	}

	return true
}
