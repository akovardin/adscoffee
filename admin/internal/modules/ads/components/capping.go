package components

import (
	"errors"
	"strconv"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
	"go.uber.org/zap"

	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

type Capping struct {
	logger *zap.Logger
}

func NewCapping(logger *zap.Logger) *Capping {
	return &Capping{
		logger: logger,
	}
}

func (c *Capping) Component(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	data, ok := field.Value(obj).(string)
	if !ok {
		c.logger.Error("capping field value is not string", zap.String("field", field.Name), zap.Any("value", field.Value(obj)))
	}

	capping, err := models.NewCapping(data)
	if err != nil {
		c.logger.Error("error unmarshal budget", zap.Error(err), zap.String("field", field.Name), zap.Any("value", field.Value(obj)))
	}

	components := []h.HTMLComponent{
		v.VRow(
			[]h.HTMLComponent{
				v.VCol([]h.HTMLComponent{
					h.Div(
						h.Label("Показы").Class("text-subtitle-2"),
					).Style("padding-bottom: 12px;"),
					v.VTextField().
						Hint("1000").
						Variant("outlined").Density("compact").
						Attr(web.VField("Capping.Count", capping.Count)...),
				}...),

				v.VCol([]h.HTMLComponent{
					h.Div(
						h.Label("Период (часы)").Class("text-subtitle-2"),
					).Style("padding-bottom: 12px;"),
					v.VTextField().
						Hint("1000").
						Variant("outlined").Density("compact").
						Attr(web.VField("Capping.Period", capping.Period)...),
				}...),
			}...,
		),
	}

	return h.Div(components...).Class("capping-field")
}

func (c *Capping) Setter(obj any, field *presets.FieldContext, ctx *web.EventContext) error {
	data, ok := field.Value(obj).(string)
	if !ok {
		return errors.New("capping field value is not string")
	}

	capping, err := models.NewCapping(data)
	if err != nil {
		return err
	}

	capping.Count, err = strconv.Atoi(ctx.R.FormValue("Capping.Count"))
	if err != nil {
		return err
	}

	capping.Period, err = strconv.Atoi(ctx.R.FormValue("Capping.Period"))
	if err != nil {
		return err
	}

	return reflectutils.Set(obj, field.Name, capping.String())
}
