package targetings

import "go.ads.coffee/platform/server/internal/domain/plugins"

type Targetings = plugins.Registry[plugins.Targeting]

func New(list []plugins.Provider[plugins.Targeting]) *Targetings {
	return plugins.NewRegistry[plugins.Targeting]("targetings", list)
}
