package stages

import "go.ads.coffee/platform/server/internal/domain/plugins"

type Stages = plugins.Registry[plugins.Stage]

func New(list []plugins.Provider[plugins.Stage]) *Stages {
	return plugins.NewRegistry[plugins.Stage]("stages", list)
}
