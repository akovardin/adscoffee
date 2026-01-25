package main

import (
	"context"
	"net/http"

	"go.uber.org/config"
	"go.uber.org/fx"

	cfg "go.ads.coffee/server/config"
	"go.ads.coffee/server/domain"
	"go.ads.coffee/server/pipeline"
	"go.ads.coffee/server/plugins"
)

func main() {
	fx.New(
		fx.Provide(
			func() (cfg.Config, error) {

				base := config.File("config.yaml")

				provider, err := config.NewYAML(base)
				if err != nil {
					return cfg.Config{}, err
				}

				c := cfg.Config{}

				if err := provider.Get("").Populate(&c); err != nil {
					return cfg.Config{}, err
				}

				return c, nil
			},
		),

		plugins.Module,

		fx.Invoke(
			start,
		),
	).Run()
}

func start(manager *pipeline.Manager) {
	r, _ := http.NewRequest(http.MethodGet, "", nil)
	w := &Response{}

	state := domain.State{
		Request:  r,
		Response: w,
	}

	manager.Process(context.Background(), state)
}

type Response struct {
}

func (r *Response) Write(b []byte) (int, error) {
	return 0, nil
}

func (r *Response) Header() http.Header {
	return http.Header{}
}

func (r *Response) WriteHeader(statusCode int) {

}
