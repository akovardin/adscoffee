package apps

import (
	"go.uber.org/fx"

	"go.ads.coffee/server/domain"
)

var Module = fx.Module(
	"targetings.apps",
	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Targeting)),
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

func (a *Apps) Copy(cfg map[string]any) domain.Targeting {
	return &Apps{}
}

func (a *Apps) Filter() {}
