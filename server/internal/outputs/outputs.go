package outputs

import "go.ads.coffee/platform/server/internal/domain/plugins"

type Outputs = plugins.Registry[plugins.Output]

func New(list []plugins.Provider[plugins.Output]) *Outputs {
	return plugins.NewRegistry[plugins.Output]("outputs", list)
}
