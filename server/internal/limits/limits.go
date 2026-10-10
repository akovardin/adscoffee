package limits

import "go.ads.coffee/platform/server/internal/domain/plugins"

type Limits = plugins.Registry[plugins.Limit]

func New(list []plugins.Provider[plugins.Limit]) *Limits {
	return plugins.NewRegistry[plugins.Limit]("limits", list)
}
