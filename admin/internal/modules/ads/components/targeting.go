package components

import (
	"errors"
	"strings"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/form"
	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

type Targeting struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewTargeting(logger *zap.Logger, db *gorm.DB) *Targeting {
	return &Targeting{
		logger: logger,
		db:     db,
	}
}

func (t *Targeting) Component(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	// include or
	// exclude or
	data, ok := field.Value(obj).(string)
	if !ok {
		t.logger.Error("targeting field value is not string", zap.String("field", field.Name), zap.Any("value", field.Value(obj)))
	}

	targeting, err := models.NewTargeting(data)
	if err != nil {
		t.logger.Error("error unmarshal targeting", zap.Error(err), zap.String("field", field.Name), zap.Any("value", field.Value(obj)))
	}

	border := "border: 1px solid #ddd; border-radius: 4px; margin-bottom: 10px;"
	header := "display: inline; margin: 0;"
	summary := "cursor: pointer; padding: 12px;"

	networks := []models.Network{}
	if err := t.db.Model(&models.Network{}).Find(&networks).Error; err != nil {
		t.logger.Error("erro on load networks", zap.Error(err))
	}

	audiences := []models.Audience{}
	if err := t.db.Model(&models.Audience{}).Find(&audiences).Error; err != nil {
		t.logger.Error("erro on load audiences", zap.Error(err))
	}

	components := []h.HTMLComponent{
		v.VCol([]h.HTMLComponent{
			h.Details(
				h.Summary(
					h.H3("SSP").Style(header),
				).Style(summary),

				h.Div(
					h.Div(
						h.Label("Включить").Class("v-label theme--dark"),

						v.VSelect().
							Multiple(true).
							Items(networks).
							ItemTitle("Title").ItemValue("Name").
							Attr(web.VField(field.Name+".Network.IncludeOr", targeting.Network.IncludeOr)...).
							ErrorMessages(field.Errors...),
					).Class("mb-4"),
					h.Div(
						h.Label("Исключить").Class("v-label theme--dark"),

						v.VSelect().
							Multiple(true).
							Items(networks).
							ItemTitle("Title").ItemValue("Name").
							Attr(web.VField(field.Name+".Network.ExcludeOr", targeting.Network.ExcludeOr)...).
							ErrorMessages(field.Errors...),
					),
				).Style("padding: 16px;"),
			).Style(border),

			h.Details(
				h.Summary(
					h.H3("Бандлы").Style(header),
				).Style(summary),

				h.Div(
					h.Div([]h.HTMLComponent{
						h.Label("Включить").Class("v-label theme--dark"),
						v.VTextarea().
							Hint("com.example ru.rustore").
							Attr(web.VField("Targeting.Bundle.IncludeOr",
								strings.Join(targeting.Bundle.IncludeOr, " "))...).
							Disabled(false).
							ErrorMessages(field.Errors...),
					}...),
					h.Div([]h.HTMLComponent{
						h.Label("Исключить").Class("v-label theme--dark"),
						v.VTextarea().
							Hint("com.example ru.rustore").
							Attr(web.VField("Targeting.Bundle.ExcludeOr",
								strings.Join(targeting.Bundle.ExcludeOr, " "))...).
							Disabled(false).
							ErrorMessages(field.Errors...),
					}...),
				).Style("padding: 16px;"),
			).Style(border),

			h.Details(
				h.Summary(
					h.H3("Bapp").Style(header),
				).Style(summary),

				h.Div(
					h.Div([]h.HTMLComponent{
						h.Label("Включить").Class("v-label theme--dark"),
						v.VTextarea().
							Hint("com.example ru.rustore").
							Attr(web.VField("Targeting.Bapp.IncludeOr",
								strings.Join(targeting.Bapp.IncludeOr, " "))...).
							Disabled(false).
							ErrorMessages(field.Errors...),
					}...),
					h.Div([]h.HTMLComponent{
						h.Label("Исключить").Class("v-label theme--dark"),
						v.VTextarea().
							Hint("com.example ru.rustore").
							Attr(web.VField("Targeting.Bapp.ExcludeOr",
								strings.Join(targeting.Bapp.ExcludeOr, " "))...).
							Disabled(false).
							ErrorMessages(field.Errors...),
					}...),
				).Style("padding: 16px;"),
			).Style(border),

			h.Details(
				h.Summary(
					h.H3("Страны").Style(header),
				).Style(summary),

				h.Div(
					h.Div([]h.HTMLComponent{
						h.Label("Включить").Class("v-label theme--dark"),
						v.VTextarea().
							Hint("RU US").
							Attr(web.VField("Targeting.Country.IncludeOr",
								strings.Join(targeting.Country.IncludeOr, " "))...).
							Disabled(false).
							ErrorMessages(field.Errors...),
					}...),
					h.Div([]h.HTMLComponent{
						h.Label("Исключить").Class("v-label theme--dark"),
						v.VTextarea().
							Hint("RU US").
							Attr(web.VField("Targeting.Country.ExcludeOr",
								strings.Join(targeting.Country.ExcludeOr, " "))...).
							Disabled(false).
							ErrorMessages(field.Errors...),
					}...),
				).Style("padding: 16px;"),
			).Style(border),

			h.Details(
				h.Summary(
					h.H3("Регионы").Style(header),
				).Style(summary),

				h.Div(
					h.Div([]h.HTMLComponent{
						h.Label("Включить").Class("v-label theme--dark"),
						v.VTextarea().
							Hint("SPE MOW").
							Attr(web.VField("Targeting.Region.IncludeOr",
								strings.Join(targeting.Region.IncludeOr, " "))...).
							Disabled(false).
							ErrorMessages(field.Errors...),
					}...),
					h.Div([]h.HTMLComponent{
						h.Label("Исключить").Class("v-label theme--dark"),
						v.VTextarea().
							Hint("SPE MOW").
							Attr(web.VField("Targeting.Region.ExcludeOr",
								strings.Join(targeting.Region.ExcludeOr, " "))...).
							Disabled(false).
							ErrorMessages(field.Errors...),
					}...),
				).Style("padding: 16px;"),
			).Style(border),

			h.Details(
				h.Summary(
					h.H3("Города").Style(header),
				).Style(summary),

				h.Div(
					h.Div([]h.HTMLComponent{
						h.Label("Включить").Class("v-label theme--dark"),
						v.VTextarea().
							Hint("KUF OMS").
							Attr(web.VField("Targeting.City.IncludeOr",
								strings.Join(targeting.City.IncludeOr, " "))...).
							Disabled(false).
							ErrorMessages(field.Errors...),
					}...),
					h.Div([]h.HTMLComponent{
						h.Label("Исключить").Class("v-label theme--dark"),
						v.VTextarea().
							Hint("KUF OMS").
							Attr(web.VField("Targeting.City.ExcludeOr",
								strings.Join(targeting.City.ExcludeOr, " "))...).
							Disabled(false).
							ErrorMessages(field.Errors...),
					}...),
				).Style("padding: 16px;"),
			).Style(border),

			h.Details(
				h.Summary(
					h.H3("IP").Style(header),
				).Style(summary),

				h.Div(
					h.Div([]h.HTMLComponent{
						h.Label("Включить").Class("v-label theme--dark"),
						v.VTextarea().
							Hint("188.170.172.0/22 188.170.192.0/22").
							Attr(web.VField("Targeting.IP.Include",
								strings.Join(targeting.IP.Include, " "))...).
							Disabled(false).
							ErrorMessages(field.Errors...),
					}...),
					h.Div([]h.HTMLComponent{
						h.Label("Исключить").Class("v-label theme--dark"),
						v.VTextarea().
							Hint("188.170.172.0/22 188.170.192.0/22").
							Attr(web.VField("Targeting.IP.Exclude",
								strings.Join(targeting.IP.Exclude, " "))...).
							Disabled(false).
							ErrorMessages(field.Errors...),
					}...),
				).Style("padding: 16px;"),
			).Style(border),

			h.Details(
				h.Summary(
					h.H3("Аудитории").Style(header),
				).Style(summary),

				h.Div(
					h.Div([]h.HTMLComponent{
						h.Label("Включить").Class("v-label theme--dark"),

						v.VSelect().
							Multiple(true).
							Items(audiences).
							ItemTitle("Title").ItemValue("Name").
							Attr(web.VField(field.Name+".Audience.IncludeOr", targeting.Audience.IncludeOr)...).
							ErrorMessages(field.Errors...),
					}...),
					h.Div([]h.HTMLComponent{
						h.Label("Исключить").Class("v-label theme--dark"),

						v.VSelect().
							Multiple(true).
							Items(audiences).
							ItemTitle("Title").ItemValue("Name").
							Attr(web.VField(field.Name+".Audience.ExcludeOr", targeting.Audience.ExcludeOr)...).
							ErrorMessages(field.Errors...),
					}...),
				).Style("padding: 16px;"),
			).Style(border),
		}...),
	}

	return h.Div(components...).Class("targeting-field")
}

func (t *Targeting) Setter(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) error {
	data, ok := field.Value(obj).(string)
	if !ok {
		return errors.New("budget field value is not string")
	}

	targeting, err := models.NewTargeting(data)
	if err != nil {
		return err
	}

	if ctx.R.Form.Has("Targeting.Bundle.IncludeOr") {
		targeting.Bundle.IncludeOr = strings.Fields(ctx.R.FormValue("Targeting.Bundle.IncludeOr"))
	}
	if ctx.R.Form.Has("Targeting.Bundle.ExcludeOr") {
		targeting.Bundle.ExcludeOr = strings.Fields(ctx.R.FormValue("Targeting.Bundle.ExcludeOr"))
	}

	if ctx.R.Form.Has("Targeting.Bapp.IncludeOr") {
		targeting.Bapp.IncludeOr = strings.Fields(ctx.R.FormValue("Targeting.Bapp.IncludeOr"))
	}
	if ctx.R.Form.Has("Targeting.Bapp.ExcludeOr") {
		targeting.Bapp.ExcludeOr = strings.Fields(ctx.R.FormValue("Targeting.Bapp.ExcludeOr"))
	}

	if ctx.R.Form.Has("Targeting.Country.IncludeOr") {
		targeting.Country.IncludeOr = strings.Fields(ctx.R.FormValue("Targeting.Country.IncludeOr"))
	}
	if ctx.R.Form.Has("Targeting.Region.IncludeOr") {
		targeting.Region.IncludeOr = strings.Fields(ctx.R.FormValue("Targeting.Region.IncludeOr"))
	}
	if ctx.R.Form.Has("Targeting.City.IncludeOr") {
		targeting.City.IncludeOr = strings.Fields(ctx.R.FormValue("Targeting.City.IncludeOr"))
	}

	if ctx.R.Form.Has("Targeting.Country.ExcludeOr") {
		targeting.Country.ExcludeOr = strings.Fields(ctx.R.FormValue("Targeting.Country.ExcludeOr"))
	}
	if ctx.R.Form.Has("Targeting.Region.ExcludeOr") {
		targeting.Region.ExcludeOr = strings.Fields(ctx.R.FormValue("Targeting.Region.ExcludeOr"))
	}
	if ctx.R.Form.Has("Targeting.City.ExcludeOr") {
		targeting.City.ExcludeOr = strings.Fields(ctx.R.FormValue("Targeting.City.ExcludeOr"))
	}

	if ctx.R.Form.Has("Targeting.IP.Include") {
		targeting.IP.Include = strings.Fields(ctx.R.FormValue("Targeting.IP.Include"))
	}
	if ctx.R.Form.Has("Targeting.IP.Exclude") {
		targeting.IP.Exclude = strings.Fields(ctx.R.FormValue("Targeting.IP.Exclude"))
	}

	if ctx.R.Form.Has("Targeting.Network.IncludeOr") {
		values := form.Values(ctx.R, "Targeting.Network.IncludeOr")
		targeting.Network.IncludeOr = values
	}
	if ctx.R.Form.Has("Targeting.Network.ExcludeOr") {
		values := form.Values(ctx.R, "Targeting.Network.ExcludeOr")
		targeting.Network.ExcludeOr = values
	}

	if ctx.R.Form.Has("Targeting.Audience.IncludeOr") {
		values := form.Values(ctx.R, "Targeting.Audience.IncludeOr")
		targeting.Audience.IncludeOr = values
	}
	if ctx.R.Form.Has("Targeting.Audience.ExcludeOr") {
		values := form.Values(ctx.R, "Targeting.Audience.ExcludeOr")
		targeting.Audience.ExcludeOr = values
	}

	return reflectutils.Set(obj, field.Name, targeting.String())
}
