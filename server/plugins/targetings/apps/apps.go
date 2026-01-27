package apps

import (
	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"targetings.apps",
	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(plugins.Targeting)),
			fx.ResultTags(`group:"targetings"`),
		),
	),
)

type Apps struct{}

func New() *Apps {
	return &Apps{}
}

func (a *Apps) Name() string {
	return "targetings.apps"
}

func (a *Apps) Copy(cfg map[string]any) plugins.Targeting {
	return &Apps{}
}

func (a *Apps) Filter() {}
