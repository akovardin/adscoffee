package stats

import (
	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	h "github.com/theplant/htmlgo"
)

type Stats struct{}

type Dashboard struct{}

func New() *Stats {
	return &Stats{}
}

func (s *Stats) Configure(pb *presets.Builder) {
	b := pb.Model(&Dashboard{}).
		MenuIcon("mdi-view-dashboard").
		URIName("dashboard")

	lb := b.Listing()

	lb.PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		r.Body = h.Div(
			h.RawHTML(`<iframe
				src="http://localhost:3001/d/caffeine-analytics/caffeine-analytics?orgId=1&kiosk=tv&autofitpanels=1&theme=light"
				style="width: 100%; height: calc(100vh - 64px); border: none;"
				allow="fullscreen"
			></iframe>`),
		)
		r.PageTitle = "Аналитика"
		return
	})
}