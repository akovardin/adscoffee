package stats

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

type Stats struct {
	db       *gorm.DB
	query    *Query
	template *template.Template
}

func New(db *gorm.DB, query *Query) *Stats {
	tm := template.Must(template.New("").Parse(page))

	return &Stats{
		db:       db,
		query:    query,
		template: tm,
	}
}

type Dashboard struct{}

type DateRow struct {
	ID   int
	Date string
}

type Option struct {
	ID   string
	Name string
}

func (s *Stats) Configure(pb *presets.Builder) {
	b := pb.Model(&Dashboard{}).
		MenuIcon("mdi-view-dashboard").
		URIName("dashboard")

	lb := b.Listing()

	lb.PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		// filters
		startedAt := started(ctx)
		endedAt := ended(ctx, startedAt)

		metrics := parse(ctx, "metrics")
		if metrics == nil {
			metrics = []string{MetricImpressions, MetricClicks}
		}

		grouped := parse(ctx, "grouped")
		if grouped == nil {
			grouped = []string{GroupNetwork}
		}

		banenrs := parse(ctx, "banenrs")
		groups := parse(ctx, "groups")
		campaigns := parse(ctx, "campaigns")
		advertisers := parse(ctx, "advertisers")

		bundles := parse(ctx, "bundles")
		networks := parse(ctx, "networks")

		// load stats
		data, err := s.query.Select(context.Background(), Condition{
			From:    startedAt,
			To:      endedAt,
			Metrics: metrics,
			Filters: []Filter{
				{
					Field: FilterAdvertiserId,
					Value: advertisers,
				},
				{
					Field: FilterCampaignId,
					Value: campaigns,
				},
				{
					Field: FilterGroupId,
					Value: groups,
				},
				{
					Field: FilterBannerId,
					Value: banenrs,
				},
				{
					Field: FilterBundle,
					Value: bundles,
				},
				{
					Field: FilterNetwork,
					Value: networks,
				},
			},
			Groups: grouped,
		})

		if err != nil {
			return r, err
		}

		notes := make([]*DateRow, 0, len(data.Labels))

		for i := len(data.Labels); i > 0; i-- {
			notes = append(notes, &DateRow{
				ID:   i - 1,
				Date: data.Labels[i-1],
			})
		}

		// prepare table
		dt := vx.DataTable(notes).WithoutHeader(false)

		bundleOptions, _ := s.query.bundles(context.Background())
		networkOptions, _ := s.query.networks(context.Background())

		dt.Hover(true)

		dt.Column("Date").Title("Дата и час").CellComponentFunc(func(obj interface{}, fieldName string, ctx *web.EventContext) h.HTMLComponent {
			n := obj.(*DateRow)
			return h.Td(h.Text(n.Date)).Style("width: 136px; min-width: 136px;")
		})

		for key, dataset := range data.Datasets {
			dt.Column(key).Title(key).CellComponentFunc(func(obj interface{}, fieldName string, ctx *web.EventContext) h.HTMLComponent {
				n := obj.(*DateRow)
				return h.Td(h.Text(fmt.Sprintf("%.2f", dataset[n.Date]))).Style("min-width: 100px")
			})
		}

		// graph
		var script bytes.Buffer
		if err := s.template.Execute(&script, data); err != nil {
			return r, err
		}

		// render page
		body := v.VContainer(
			h.Div(
				h.RawHTML(`<div style="width: 100%; height: 400px;">
					<canvas id="stats-chart" style="width: 100%; height: 400px;"></canvas>
				</div>`),
				h.Script("").Attr("src", "https://cdn.jsdelivr.net/npm/chart.js"),
				h.Script(script.String()),
			).Style("margin-bottom: 30px"),

			h.Div(
				h.H3("Срезы"),
			).Style("margin-bottom: 20px; margin-top: 30px"),

			v.VRow(
				h.Div(
					vx.VXSelect().
						Items(s.metrics()).
						ItemTitle("Name").
						ItemValue("ID").
						Multiple(true).
						Clearable(true).
						Chips(true).
						Attr(web.VField("Metrics", metrics)...).
						Label("Метрики"),
				).Id("Metrics").Style("margin-right: 20px; padding-left: 12px; width: 250px"),

				h.Div(
					vx.VXSelect().
						Items(s.grouped()).
						ItemTitle("Name").
						ItemValue("ID").
						Multiple(true).
						Clearable(true).
						Chips(true).
						Attr(web.VField("Grouped", grouped)...).
						Label("Группировки"),
				).Id("Grouped").Style("margin-right: 20px; padding-left: 12px; width: 250px"),

				h.Div(
					vx.VXDatepicker().Type("datetimepicker").
						Format("YYYY-MM-DD HH:mm").
						Clearable(true).
						Id("StartedAt").
						Attr(web.VField("StartedAt", startedAt)...).
						Label("Начало").
						Width(240),
				).Style("margin-right: 34px; padding-left: 12px;"),

				vx.VXDatepicker().Type("datetimepicker").
					Format("YYYY-MM-DD HH:mm").
					Clearable(true).
					Id("EndedAt").
					Attr(web.VField("EndedAt", endedAt)...).
					Label("Конец").
					Width(240),
			),

			h.Div(
				h.H3("Фильтры"),
			).Style("margin-bottom: 20px; margin-top: 30px"),

			v.VRow(
				h.Div(
					vx.VXSelect().
						Items(s.banners()).
						ItemTitle("Name").
						ItemValue("ID").
						Multiple(true).
						Clearable(true).
						Chips(true).
						Attr(web.VField("BannerId", banenrs)...).
						Label("Баннер"),
				).Id("BannerId").Style("margin-right: 20px; padding-left: 12px; width: 250px"),

				h.Div(
					vx.VXSelect().
						Items(s.groups()).
						ItemTitle("Name").
						ItemValue("ID").
						Multiple(true).
						Clearable(true).
						Chips(true).
						Attr(web.VField("GroupId", groups)...).
						Label("Группа"),
				).Id("GroupId").Style("margin-right: 20px; padding-left: 12px; width: 250px"),

				h.Div(
					vx.VXSelect().
						Items(s.campaigns()).
						ItemTitle("Name").
						ItemValue("ID").
						Multiple(true).
						Clearable(true).
						Chips(true).
						Attr(web.VField("CampaignId", campaigns)...).
						Label("Кампания"),
				).Id("CampaignId").Style("margin-right: 20px; padding-left: 12px; width: 250px"),

				h.Div(
					vx.VXSelect().
						Items(s.advertisers()).
						ItemTitle("Name").
						ItemValue("ID").
						Multiple(true).
						Clearable(true).
						Chips(true).
						Attr(web.VField("AdvertiserId", advertisers)...).
						Label("Рекламодатель"),
				).Id("AdvertiserId").Style("margin-right: 20px; padding-left: 12px; width: 250px"),
			),

			v.VRow(
				h.Div(
					vx.VXSelect().
						Items(bundleOptions).
						ItemTitle("Name").
						ItemValue("ID").
						Multiple(true).
						Clearable(true).
						Chips(true).
						Attr(web.VField("Bundle", bundles)...).
						Label("Бандл"),
				).Id("Bundle").Style("margin-right: 20px; padding-left: 12px; width: 250px"),

				h.Div(
					vx.VXSelect().
						Items(networkOptions).
						ItemTitle("Name").
						ItemValue("ID").
						Multiple(true).
						Clearable(true).
						Chips(true).
						Attr(web.VField("Network", networks)...).
						Label("Сеть"),
				).Id("Network").Style("margin-right: 20px; padding-left: 12px; width: 250px"),
			),

			h.Div(
				v.VBtn("Обновить").Attr("onclick", `(function() {
						// filtesr
						const startedAt = document.querySelector('#StartedAt input[type="text"]').value;
						const endedAt = document.querySelector('#EndedAt input[type="text"]').value;
						const banenrs = document.querySelector('#BannerId input[type="text"]').value;
						const groups = document.querySelector('#GroupId input[type="text"]').value;
						const campaigns = document.querySelector('#CampaignId input[type="text"]').value;
						const advertisers = document.querySelector('#AdvertiserId input[type="text"]').value;
						const bundles = document.querySelector('#Bundle input[type="text"]').value;
						const networks = document.querySelector('#Network input[type="text"]').value;
						
						// grouped
						const grouped = document.querySelector('#Grouped input[type="text"]').value;

						// metrics
						const metrics = document.querySelector('#Metrics input[type="text"]').value;

						const params = new URLSearchParams();

						if (startedAt != "") {
							params.append("started_at", startedAt);
						}	
					
						if (endedAt != "") {
							params.append("ended_at", endedAt);
						}

						if (banenrs != "") {
							params.append("banners", banenrs);
						}

						if (groups != "") {
							params.append("groups", groups);
						}

						if (campaigns != "") {
							params.append("campaigns", campaigns);
						}

						if (advertisers != "") {
							params.append("advertisers", advertisers);
						}

						if (grouped != "") {
							params.append("grouped", grouped);
						}

						if (metrics != "") {
							params.append("metrics", metrics);
						}

						if (bundles != "") {
							params.append("bundles", bundles);
						}

						if (networks != "") {
							params.append("networks", networks);
						}


						window.location.href = '/admin/dashboard?' + params.toString();
					})()`).Variant(v.VariantFlat).Class("bg-primary"),
			).Style("margin-top: 20px; margin-bottom: 30px"),

			h.Div(
				dt,
			).Style("margin-top: 20px"),
		)

		r.Body = body
		r.PageTitle = "Отчеты"

		return
	})
}

const page = `
document.addEventListener('DOMContentLoaded', function() {
	const ctx = document.getElementById('stats-chart').getContext('2d');

	new Chart(ctx, {
		type: 'line',
		data: {
			labels: [
				{{range $index, $elem := .Labels}}
					{{if $index}},{{end}}
					'{{$elem}}'
				{{end}}
			],
			datasets: [
				{{$labels := .Labels}}

				{{range $name, $dataset := .Datasets}}
				{
					label: '{{$name}}',
					data: [
					{{range $index, $elem := $labels}}
						{{if $index}},{{end}}
						{{ index $dataset $elem }}
					{{end}}
					],
					borderWidth: 1
				},
				{{end}}
			]
		},
		options: {
			responsive: true,
			maintainAspectRatio: false,
			scales: {
				y: {
					beginAtZero: true
				}
			}
		}
	});
});
`

func (s *Stats) banners() []Option {
	modles := []models.Banner{}

	s.db.Model(&models.Banner{}).Find(&modles)

	items := make([]Option, 0, len(modles))
	for _, model := range modles {
		items = append(items, Option{
			ID:   fmt.Sprintf("%d", model.ID),
			Name: model.Title,
		})
	}

	return items
}

func (s *Stats) groups() []Option {
	modles := []models.Bgroup{}

	s.db.Model(&models.Bgroup{}).Find(&modles)

	items := make([]Option, 0, len(modles))
	for _, model := range modles {
		items = append(items, Option{
			ID:   fmt.Sprintf("%d", model.ID),
			Name: model.Title,
		})
	}

	return items
}

func (s *Stats) campaigns() []Option {
	modles := []models.Campaign{}

	s.db.Model(&models.Campaign{}).Find(&modles)

	items := make([]Option, 0, len(modles))
	for _, group := range modles {
		items = append(items, Option{
			ID:   fmt.Sprintf("%d", group.ID),
			Name: group.Title,
		})
	}

	return items
}

func (s *Stats) advertisers() []Option {
	modles := []models.Advertiser{}

	s.db.Model(&models.Advertiser{}).Find(&modles)

	items := make([]Option, 0, len(modles))
	for _, group := range modles {
		items = append(items, Option{
			ID:   fmt.Sprintf("%d", group.ID),
			Name: group.Title,
		})
	}

	return items
}

func (s *Stats) metrics() []Option {
	return []Option{
		{
			ID:   MetricRequests,
			Name: "Запросы",
		},
		{
			ID:   MetricResponses,
			Name: "Респонсы",
		},
		{
			ID:   MetricWins,
			Name: "Победы",
		},
		{
			ID:   MetricImpressions,
			Name: "Показы",
		},
		{
			ID:   MetricClicks,
			Name: "Клики",
		},
		{
			ID:   MetricConversions,
			Name: "Конверсии",
		},

		{
			ID:   MetricPrice,
			Name: "Деньги",
		},
	}
}

func (s *Stats) grouped() []Option {
	return []Option{
		{
			ID:   GroupBanner,
			Name: "Баннеры",
		},
		{
			ID:   GroupGroup,
			Name: "Группы",
		},
		{
			ID:   GroupCampaign,
			Name: "Кампании",
		},
		{
			ID:   GroupAdvertiser,
			Name: "Рекламодатели",
		},
		{
			ID:   GroupNetwork,
			Name: "Сети",
		},
		{
			ID:   GroupBundle,
			Name: "Бандлы",
		},
	}
}

func parse(ctx *web.EventContext, key string) []string {
	v := ctx.R.FormValue(key)
	if v == "" {
		return nil
	}

	items := []string{}

	for _, c := range strings.Split(v, ",") {
		items = append(items, strings.TrimSpace(c))
	}

	return items
}

func started(ctx *web.EventContext) time.Time {
	in := ctx.R.FormValue("started_at")
	if in == "" {
		return time.Now().Add(-(time.Hour * 24))
	}

	t, err := time.Parse(DateHourMinute, in)
	if err != nil {
		return time.Now().Add(-(time.Hour * 24))
	}

	return t
}

func ended(ctx *web.EventContext, started time.Time) time.Time {
	var t time.Time

	in := ctx.R.FormValue("ended_at")
	if in == "" {
		t = time.Now()
	}

	t, err := time.Parse(DateHourMinute, in)
	if err != nil {
		t = time.Now()
	}

	if t.Before(started) {
		return started
	}

	return t
}
