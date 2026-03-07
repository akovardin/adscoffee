package stats

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
	"gorm.io/gorm"
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

type FilterItem struct {
	ID   string
	Name string
}

func (s *Stats) Configure(pb *presets.Builder) {
	b := pb.Model(&Dashboard{}).Label("Dashboard").URIName("dashboard")

	lb := b.Listing()

	lb.PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		// load data from stats
		// prepare data for chart

		// filters
		startedAt := ctx.R.FormValue("StartedAt")
		_ = startedAt
		endedAt := ctx.R.FormValue("EndedAt")

		example := ctx.R.FormValue("example")
		_ = example

		metrics := ctx.R.FormValue("Metrics")
		_ = metrics

		filters := ctx.R.FormValue("Filters")
		_ = filters

		groups := ctx.R.FormValue("Groups")
		_ = groups

		// load stats
		data, err := s.query.Select(context.Background(), Condition{
			From:    time.Now().Add(-time.Hour * 24 * 6),
			To:      time.Now().Add(-time.Hour * 24 * 3),
			Metrics: []string{MetricImpressions, MetricClicks},
			Filters: []Filter{
				{
					Field: FilterCampaignId,
					Value: []any{2},
				},
				{
					Field: FilterGroupId,
					Value: []any{85, 89, 78},
				},
			},
			Groups: []string{"group_id", "network"},
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

		dt.Column("Date").Title("Дата и час").CellComponentFunc(func(obj interface{}, fieldName string, ctx *web.EventContext) h.HTMLComponent {
			n := obj.(*DateRow)
			return h.Td(h.Text(n.Date))
		})

		for key, dataset := range data.Datasets {
			dt.Column(key).Title(key).CellComponentFunc(func(obj interface{}, fieldName string, ctx *web.EventContext) h.HTMLComponent {
				n := obj.(*DateRow)
				return h.Td(h.Text(fmt.Sprintf("%.2f", dataset[n.Date])))
			})
		}

		// graph
		var script bytes.Buffer
		if err := s.template.Execute(&script, data); err != nil {
			return r, err
		}

		items := []FilterItem{
			{ID: "1", Name: "Item 1"},
			{ID: "2", Name: "Item 2"},
		}

		bannerId := "1"

		// render page
		body := v.VContainer(
			h.Div(
				h.H3("Период"),
			).Style("margin-bottom: 20px"),
			v.VRow(
				h.Div(
					vx.VXDatepicker().Type("datetimepicker").
						Format("YYYY-MM-DD HH:mm").
						Clearable(true).
						Id("StartedAt").
						Attr(web.VField("StartedAt", startedAt)...).
						Label("Начало"),
				).Style("margin-right: 30px; padding-left: 12px;"),
				vx.VXDatepicker().Type("datetimepicker").
					Format("YYYY-MM-DD HH:mm").
					Clearable(true).
					Id("EndedAt").
					Attr(web.VField("EndedAt", endedAt)...).
					Label("Конец"),
			),
			h.Div(
				h.H3("Фильтры"),
			).Style("margin-bottom: 20px"),
			v.VRow(
				h.Div(
					vx.VXSelect().
						Items(items).
						ItemTitle("Name").
						ItemValue("ID").
						Multiple(true).
						Attr(web.VField("BannerId", bannerId)...).
						Label("Баннер"),
				).Id("BannerId").Style("margin-right: 20px; padding-left: 12px; width: 200px"),
				h.Div(
					vx.VXSelect().
						Items(items).
						ItemTitle("Name").
						ItemValue("ID").
						Multiple(true).
						// Attr(web.VField("CampaignId", campaignId)...).
						Label("Группа"),
				).Style("margin-right: 20px"),
				h.Div(
					vx.VXSelect().
						Items(items).
						ItemTitle("Name").
						ItemValue("ID").
						// Attr(web.VField("CampaignId", campaignId)...).
						Label("Кампания"),
				).Style("margin-right: 20px"),
			),

			h.Div(
				h.H3("Метрики"),
			).Style("margin-bottom: 20px"),

			h.Div(
				h.H3("Групировки"),
			).Style("margin-bottom: 20px"),

			h.Div(
				v.VBtn("Обновить").Attr("onclick", `(function() {
						const startedAt = document.querySelector('#StartedAt input[type="text"]').value;
						const endedAt = document.querySelector('#EndedAt input[type="text"]').value;
						const banenrs = document.querySelector('#BannerId input[type="text"]').value;
				
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

						window.location.href = '/admin/dashboard?' + params.toString();
					})()`).Variant(v.VariantFlat).Class("bg-primary"),
			),

			h.RawHTML(
				`<div style="width: 100%; height: 400px;">
					<canvas id="myChart" style="width: 100%; height: 400px;"></canvas>
				</div>`,
			),
			h.Script("").Attr("src", "https://cdn.jsdelivr.net/npm/chart.js"),
			h.Script(script.String()),
			h.Div(
				dt,
			).Style("height:400px"),
		)

		r.Body = body
		r.PageTitle = "Dashboard"

		return
	})
}

const page = `
document.addEventListener('DOMContentLoaded', function() {
	const ctx = document.getElementById('myChart').getContext('2d');

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
