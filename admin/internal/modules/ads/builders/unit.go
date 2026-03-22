package builders

import (
	"fmt"
	"strings"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/presets/gorm2op"
	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	"github.com/qor5/x/v3/ui/vuetifyx"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

type Unit struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewUnit(logger *zap.Logger, db *gorm.DB) *Unit {
	return &Unit{
		logger: logger,
		db:     db,
	}
}

func (n *Unit) Configure(b *presets.Builder) *presets.ModelBuilder {
	mn := b.Model(&models.Unit{}).
		MenuIcon("mdi-lan").
		// Label("Рекламодатели").
		RightDrawerWidth("1000")

	mnl := mn.Listing("ID", "Name", "Price", "Placement", "Network", "Active").
		SearchFunc(func(ctx *web.EventContext, params *presets.SearchParams) (result *presets.SearchResult, err error) {
			// по умоланию архивные сущности не показываются
			// только если явно выбрать их в фильтре
			exist := false
			for _, v := range params.SQLConditions {
				if strings.Contains(v.Query, "archived_at is not null") {
					exist = true
					break
				}

				if strings.Contains(v.Query, "(archived_at is not null or archived_at is null)") {
					exist = true
					break
				}
			}

			if !exist {
				qdb := n.db.Where("archived_at is null")
				return gorm2op.DataOperator(qdb).Search(ctx, params)
			} else {
				qdb := n.db.Where("")
				return gorm2op.DataOperator(qdb).Search(ctx, params)
			}
		}).
		SearchColumns("Name").
		// SelectableColumns(true).
		OrderableFields([]*presets.OrderableField{
			{
				FieldName: "ID",
				DBColumn:  "id",
			},
			{
				FieldName: "Name",
				DBColumn:  "name",
			},
			{
				FieldName: "Active",
				DBColumn:  "active",
			},
		})

	mnl.FilterDataFunc(func(ctx *web.EventContext) vuetifyx.FilterData {
		return []*vuetifyx.FilterItem{
			{
				Key:          "archived",
				Label:        "Архив",
				ItemType:     vuetifyx.ItemTypeSelect,
				SQLCondition: "archived_at is null",
				Options: []*vuetifyx.SelectItem{
					{

						Text:         "В архиве",
						Value:        "is_archived",
						SQLCondition: "archived_at is not null",
					},
					{
						Text:         "Все",
						Value:        "all",
						SQLCondition: "(archived_at is not null or archived_at is null)",
					},
				},
			},
			{
				Key:      "active",
				Label:    "Активность",
				ItemType: vuetifyx.ItemTypeSelect,
				Options: []*vuetifyx.SelectItem{
					{

						Text:         "Включен",
						Value:        "is_active",
						SQLCondition: "active = true",
					},
					{
						Text:         "Выключен",
						Value:        "not_active",
						SQLCondition: "active = false",
					},
				},
			},
		}
	})

	mnl.Field("Placement").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Unit)
		var placement models.Placement
		if c.PlacementID == 0 {
			return h.Td()
		}

		n.db.First(&placement, "id = ?", c.PlacementID)

		return h.Td().Text(placement.Name)
	})

	mnl.Field("Network").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Unit)
		var network models.Network
		if c.NetworkID == 0 {
			return h.Td()
		}

		n.db.First(&network, "id = ?", c.NetworkID)

		return h.Td().Text(network.Title)
	})

	mnl.Field("Active").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Unit)

		color := "red"
		text := "выключен"
		if c.Active {
			text = "включен"
			color = "green"
		}

		return h.Td().Children(h.Span(text).Style("color:" + color))
	})

	mne := mn.Editing(
		&presets.FieldsSection{
			// Title: "Info",
			Rows: [][]string{
				{"Name"},
				{"Price"},
				{"PlacementID"},
				{"NetworkID"},
				{"Data"},
			},
		},
	).ValidateFunc(func(obj interface{}, ctx *web.EventContext) (err web.ValidationErrors) {
		u := obj.(*models.Unit)

		if u.Name == "" {
			err.FieldError("Name", "Name is required")
		}
		return
	})

	mne.Field("PlacementID").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Unit)

		var items []models.Placement
		n.db.Find(&items)

		sel := v.VSelect().
			Variant("outlined").Density("compact").
			Label("Плейсмент").
			Items(items).
			ItemTitle("Name").
			ItemValue("ID").
			Attr(web.VField("PlacementID", c.PlacementID)...)

		return h.Div(
			sel,
		)
	})

	mne.Field("NetworkID").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Unit)

		var items []models.Network
		n.db.Find(&items)

		sel := v.VSelect().
			Variant("outlined").Density("compact").
			Label("Рекламная сеть").
			Items(items).
			ItemTitle("Title").
			ItemValue("ID").
			Attr(web.VField("NetworkID", c.NetworkID)...)

		return h.Div(
			sel,
		)
	})

	mne.Field("Data").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return v.VTextarea().
			Label(field.Label).
			Attr(web.VField(field.FormKey, fmt.Sprint(reflectutils.MustGet(obj, field.Name)))...).
			Disabled(field.Disabled).
			ErrorMessages(field.Errors...)
	})

	return mn
}
