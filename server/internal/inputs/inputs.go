package inputs

import "go.ads.coffee/platform/server/internal/domain/plugins"

type Inputs = plugins.Registry[plugins.Input]

func New(list []plugins.Provider[plugins.Input]) *Inputs {
	return plugins.NewRegistry[plugins.Input]("inputs", list)
}
