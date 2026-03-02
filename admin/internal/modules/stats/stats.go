package stats

import (
	"fmt"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
	"gorm.io/gorm"
)

type Stats struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Stats {
	return &Stats{
		db: db,
	}
}

type Dashboard struct{}

type Note struct {
	ID          int
	Date        string
	Requests    int
	Impressions int
	Clicks      int
}

func (s *Stats) Configure(pb *presets.Builder) {
	b := pb.Model(&Dashboard{}).Label("Dashboard").URIName("dashboard")

	notes := []*Note{
		{
			ID:          1,
			Date:        "10.02",
			Requests:    12,
			Impressions: 12,
		},
		{
			ID:          2,
			Date:        "11.02",
			Requests:    19,
			Impressions: 10,
			Clicks:      3,
		},
		{
			ID:          3,
			Date:        "12.02",
			Requests:    3,
			Impressions: 2,
			Clicks:      0,
		},
	}
	dt := vx.DataTable(notes).WithoutHeader(false)

	dt.Column("Date").Title("Дата").CellComponentFunc(func(obj interface{}, fieldName string, ctx *web.EventContext) h.HTMLComponent {
		n := obj.(*Note)
		return h.Td(h.Text(n.Date))
	})
	dt.Column("Value").Title("Запросы").CellComponentFunc(func(obj interface{}, fieldName string, ctx *web.EventContext) h.HTMLComponent {
		n := obj.(*Note)
		return h.Td(h.Text(fmt.Sprintf("%d", n.Requests)))
	})
	dt.Column("Value").Title("Показы").CellComponentFunc(func(obj interface{}, fieldName string, ctx *web.EventContext) h.HTMLComponent {
		n := obj.(*Note)
		return h.Td(h.Text(fmt.Sprintf("%d", n.Impressions)))
	})
	dt.Column("Value").Title("Клики").CellComponentFunc(func(obj interface{}, fieldName string, ctx *web.EventContext) h.HTMLComponent {
		n := obj.(*Note)
		return h.Td(h.Text(fmt.Sprintf("%d", n.Clicks)))
	})

	lb := b.Listing()

	lb.PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {

		body := vuetify.VContainer(
			h.RawHTML(
				`<div style="width: 100%; height: 400px;">
					<canvas id="myChart" style="width: 100%; height: 400px;"></canvas>
				</div>`,
			),
			h.Script("").Attr("src", "https://cdn.jsdelivr.net/npm/chart.js"),
			h.Script(`
			document.addEventListener('DOMContentLoaded', function() {
    const ctx = document.getElementById('myChart').getContext('2d');
    
    new Chart(ctx, {
      type: 'line',
      data: {
        labels: ['10.02', '11.02', '12.02', '13.02', '14.02', '15.02'],
        datasets: [{
          label: '# of Votes',
          data: [12, 19, 3, 5, 2, 3],
          borderWidth: 1
        }]
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
  });`),
			h.Div(
				dt,
			).Style("height:400px"),
		)

		r.Body = body
		r.PageTitle = "Dashboard"

		return
	})
}

func errorBody(msg string) h.HTMLComponent {
	return vuetify.VContainer(
		h.P().Text(msg),
	)
}
