package plugins

import "go.ads.coffee/platform/server/internal/domain/ads"

type Targeting interface {
	Name() string
	Filter(candidates []*ads.Banner, state *State) []*ads.Banner
}
