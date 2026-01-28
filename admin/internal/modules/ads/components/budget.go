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

type Budget struct {
	logger *zap.Logger
}

func NewBudget(logger *zap.Logger) *Budget {
	return &Budget{
		logger: logger,
	}
}

func (b *Budget) Component(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	data, ok := field.Value(obj).(string)
	if !ok {
		b.logger.Error("budget field value is not string", zap.String("field", field.Name), zap.Any("value", field.Value(obj)))
	}

	budget, err := models.NewBudget(data)
	if err != nil {
		b.logger.Error("error unmarshal budget", zap.Error(err), zap.String("field", field.Name), zap.Any("value", field.Value(obj)))
	}

	components := []h.HTMLComponent{
		v.VCol([]h.HTMLComponent{
			v.VRow(
				[]h.HTMLComponent{
					h.Label("Показы").Style("width: 120px; margin-left: 12px; margin-top: 7px;"),
					v.VTextField().
						Label("Суточный").
						Hint("1000").
						Variant("outlined").Density("compact").
						Attr(web.VField("Budget.Impressions.Daily", budget.Impressions.Daily)...),
					v.VTextField().
						Label("Общий").
						Hint("1000").
						Variant("outlined").Density("compact").
						Attr(web.VField("Budget.Impressions.Total", budget.Impressions.Total)...),
					v.VCheckbox().
						Label("Равномерный").
						Density("compact").
						Attr(web.VField("Budget.Impressions.Uniform", budget.Impressions.Uniform)...),
				}...,
			),
			v.VRow(
				[]h.HTMLComponent{
					h.Label("Клики").Style("width: 120px; margin-left: 12px; margin-top: 7px;"),
					v.VTextField().
						Hint("1000").
						Variant("outlined").Density("compact").
						Attr(web.VField("Budget.Clicks.Daily", budget.Clicks.Daily)...),
					v.VTextField().
						Hint("1000").
						Variant("outlined").Density("compact").
						Attr(web.VField("Budget.Clicks.Total", budget.Clicks.Total)...),
					v.VCheckbox().
						Label("Равномерный").
						Density("compact").
						Attr(web.VField("Budget.Clicks.Uniform", budget.Clicks.Uniform)...),
				}...,
			),
			v.VRow(
				[]h.HTMLComponent{
					h.Label("Конверсии").Style("width: 120px; margin-left: 12px; margin-top: 7px;"),
					v.VTextField().
						Hint("1000").
						Variant("outlined").Density("compact").
						Attr(web.VField("Budget.Conversions.Daily", budget.Conversions.Daily)...),
					v.VTextField().
						Hint("1000").
						Variant("outlined").Density("compact").
						Attr(web.VField("Budget.Conversions.Total", budget.Conversions.Total)...),
					v.VCheckbox().
						Label("Равномерный").
						Density("compact").
						Attr(web.VField("Budget.Conversions.Uniform", budget.Conversions.Uniform)...),
				}...,
			),
			v.VRow(
				[]h.HTMLComponent{
					h.Label("Деньги").Style("width: 120px; margin-left: 12px; margin-top: 7px;"),
					v.VTextField().
						Hint("1000").
						Variant("outlined").Density("compact").
						Attr(web.VField("Budget.Money.Daily", budget.Money.Daily)...),
					v.VTextField().
						Hint("1000").
						Variant("outlined").Density("compact").
						Attr(web.VField("Budget.Money.Total", budget.Money.Total)...),
					v.VCheckbox().
						Label("Равномерный").
						Density("compact").
						Attr(web.VField("Budget.Money.Uniform", budget.Money.Uniform)...),
				}...,
			),
		}...,
		),
	}

	return h.Div(components...).Class("budget-field")
}

func (b *Budget) Setter(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) error {
	data, ok := field.Value(obj).(string)
	if !ok {
		return errors.New("budget field value is not string")
	}

	budget, err := models.NewBudget(data)
	if err != nil {
		return err
	}

	budget.Impressions.Daily, err = strconv.Atoi(ctx.R.FormValue("Budget.Impressions.Daily"))
	if err != nil {
		return err
	}

	budget.Impressions.Total, err = strconv.Atoi(ctx.R.FormValue("Budget.Impressions.Total"))
	if err != nil {
		return err
	}

	budget.Impressions.Uniform, err = strconv.ParseBool(ctx.R.FormValue("Budget.Impressions.Uniform"))
	if err != nil {
		return err
	}

	budget.Clicks.Daily, err = strconv.Atoi(ctx.R.FormValue("Budget.Clicks.Daily"))
	if err != nil {
		return err
	}

	budget.Clicks.Total, err = strconv.Atoi(ctx.R.FormValue("Budget.Clicks.Total"))
	if err != nil {
		return err
	}

	budget.Clicks.Uniform, err = strconv.ParseBool(ctx.R.FormValue("Budget.Clicks.Uniform"))
	if err != nil {
		return err
	}

	budget.Conversions.Daily, err = strconv.Atoi(ctx.R.FormValue("Budget.Conversions.Daily"))
	if err != nil {
		return err
	}

	budget.Conversions.Total, err = strconv.Atoi(ctx.R.FormValue("Budget.Conversions.Total"))
	if err != nil {
		return err
	}

	budget.Conversions.Uniform, err = strconv.ParseBool(ctx.R.FormValue("Budget.Conversions.Uniform"))
	if err != nil {
		return err
	}

	budget.Money.Daily, err = strconv.Atoi(ctx.R.FormValue("Budget.Money.Daily"))
	if err != nil {
		return err
	}

	budget.Money.Total, err = strconv.Atoi(ctx.R.FormValue("Budget.Money.Total"))
	if err != nil {
		return err
	}

	budget.Money.Uniform, err = strconv.ParseBool(ctx.R.FormValue("Budget.Money.Uniform"))
	if err != nil {
		return err
	}

	return reflectutils.Set(obj, field.Name, budget.String())
}
