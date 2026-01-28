package components

import (
	"fmt"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
	"go.ads.coffee/platform/admin/internal/modules/ads/models"
	"go.uber.org/zap"
)

type Timetable struct {
	logger *zap.Logger
}

func NewTimetable(logger *zap.Logger) *Timetable {
	return &Timetable{
		logger: logger,
	}
}

func (t *Timetable) Component(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	data, ok := field.Value(obj).(string)
	if !ok {
		t.logger.Error("timetable field value is not string", zap.String("field", field.Name), zap.Any("value", field.Value(obj)))
	}

	timetable, err := models.NewTimetable(data)
	if err != nil {
		t.logger.Error("error unmarshal timetable", zap.Error(err), zap.String("field", field.Name), zap.Any("value", field.Value(obj)))
	}

	days := []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}

	components := []h.HTMLComponent{
		h.H3("Расписание активности").Style("margin-bottom: 8px;"),
	}

	titles := []h.HTMLComponent{}
	for hour := 0; hour < 24; hour++ {
		titles = append(titles, h.Span(fmt.Sprintf("%d", hour)).Style("font-size: 8px; width: 27px; display: inline-block; text-align: center;"))
	}

	components = append(components, h.Div(
		titles...,
	).Style("margin-left: 46px;"))

	for day, name := range days {
		dayComponents := []h.HTMLComponent{
			h.Strong(name).Style("display: inline-block; width: 40px;"),
		}

		for hour := 0; hour < 24; hour++ {
			key := fmt.Sprintf("Timetable[%d][%d]", day, hour)
			val, ok := timetable[day][hour]
			if !ok {
				val = true
			}

			dayComponents = append(dayComponents,
				h.Div(
					h.Input("").Type("checkbox").Name(key).Attr(web.VField(key, val)...),
				).Style("display: inline-block; margin: 2px; margin-right: 8px; padding: 2px 2px; cursor: pointer;"),
			)
		}

		components = append(components,
			v.VCol(dayComponents...),
		)
	}

	return h.Div(components...).Class("timetable-field")
}

func (t *Timetable) Setter(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) error {
	timetable := models.Timetable{}

	for day := 0; day < 7; day++ {
		for hour := 0; hour < 24; hour++ {
			if ctx.R.FormValue(fmt.Sprintf("Timetable[%d][%d]", day, hour)) == "true" {
				timetable.Set(day, hour, true)
			} else {
				timetable.Set(day, hour, false)
			}
		}
	}

	return reflectutils.Set(obj, field.Name, timetable.String())
}
