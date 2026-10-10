package plugins

import "go.ads.coffee/platform/server/internal/domain/ads"

type Limit interface {
	Name() string
	Apply(candidates []ads.Banner, state *State) []ads.Banner
}
