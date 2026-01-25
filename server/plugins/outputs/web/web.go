package web

import (
	"context"

	"go.ads.coffee/server/domain"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"outputs.web",

	fx.Provide(
		fx.Annotate(
			New,
			fx.As(new(domain.Output)),
			fx.ResultTags(`group:"outputs"`),
		),
	),
)

type Web struct {
	formtats []domain.Format
}

func New() *Web {
	input := &Web{}

	return input
}

func (w *Web) Name() string {
	return "outputs.web"
}

func (w *Web) Copy(cfg map[string]any) domain.Output {
	return &Web{
		formtats: w.formtats,
	}
}

func (w *Web) Formats(ff []domain.Format) {
	w.formtats = ff
}

func (w *Web) Process(ctx context.Context, state *domain.State) {

	// обработка разных типов запросов тоже
	// может быть вынесена в пллагины

	// а тут я могу использовать разные форматы в зависимости
	// от конфига и запроса

	state.Response.Write([]byte(`:)`))
}
