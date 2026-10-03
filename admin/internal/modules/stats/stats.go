package stats

import (
	"fmt"
	"strings"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	h "github.com/theplant/htmlgo"

	"go.ads.coffee/platform/admin/internal/config"
)

const defaultGrafanaURL = "http://localhost:3001"

const grafanaDashboardPath = "/d/caffeine-analytics/caffeine-analytics?orgId=1&kiosk=tv&autofitpanels=1&theme=light"

type Stats struct {
	grafanaURL string
}

type Dashboard struct{}

func New(grafana config.Grafana) *Stats {
	url := strings.TrimRight(grafana.URL, "/")
	if url == "" {
		url = defaultGrafanaURL
	}

	return &Stats{grafanaURL: url}
}

func (s *Stats) Configure(pb *presets.Builder) {
	b := pb.Model(&Dashboard{}).
		MenuIcon("mdi-view-dashboard").
		URIName("dashboard")

	lb := b.Listing()

	lb.PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		r.Body = h.Div(
			h.RawHTML(fmt.Sprintf(`<iframe
				src="%s%s"
				style="width: 100%%; height: calc(100vh - 64px); border: none;"
				allow="fullscreen"
			></iframe>`, s.grafanaURL, grafanaDashboardPath)),
		)
		r.PageTitle = "Аналитика"
		return
	})
}
