package builders

import (
	"fmt"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
	"github.com/qor5/x/v3/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/modules/ads/components"
	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

type Campaign struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewCampaign(logger *zap.Logger, db *gorm.DB) *Campaign {
	return &Campaign{
		logger: logger,
		db:     db,
	}
}

const copyCapmaignEvent = "copyCampaign"

func (m *Campaign) Configure(b *presets.Builder) {
	mc := b.Model(&models.Campaign{}).
		MenuIcon("mdi-bullseye-arrow").
		// Label("Кампании").
		RightDrawerWidth("1000")

	mcl := mc.Listing("ID", "Title", "Bundle", "Start", "End", "Advertiser", "Active")

	mcl.FilterDataFunc(func(ctx *web.EventContext) vuetifyx.FilterData {
		// msgr := i18n.MustGetModuleMessages(ctx.R, presets.ModelsI18nModuleKey, Messages_en_US).(*Messages)
		var options []*vuetifyx.SelectItem
		err := m.db.Model(&models.Advertiser{}).Select("title as text, id as value").Scan(&options).Error
		if err != nil {
			m.logger.Error("erro on load advertisers", zap.Error(err))

			return nil
		}

		return []*vuetifyx.FilterItem{
			{
				Key:          "advertiser",
				Label:        "Рекламодатель",
				ItemType:     vuetifyx.ItemTypeSelect,
				SQLCondition: `advertiser_id %s ?`,
				Options:      options,
			},
		}
	})

	mcl.Field("Advertiser").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Campaign)
		var comp models.Advertiser
		if c.AdvertiserID == 0 {
			return h.Td()
		}

		m.db.First(&comp, "id = ?", c.AdvertiserID)

		return h.Td().Text(comp.Title)
	})

	mcl.Field("Title").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Campaign)
		return h.Td().Children(
			h.A().
				Text(c.Title).
				Attr("onclick", "event.stopPropagation();").
				Href(fmt.Sprintf("/admin/bgroups?f_campaign=%d", c.ID)),
		)
	})

	mcn := mcl.RowMenu()

	// Добавляем обработчик копирования
	mcn.RowMenuItem("Copy").
		ComponentFunc(func(obj interface{}, id string, ctx *web.EventContext) h.HTMLComponent {
			return v.VListItem(
				web.Slot(
					v.VIcon("mdi-content-copy"), // Используем иконку копирования
				).Name("prepend"),
				v.VListItemTitle(
					h.Text("Копировать"),
				),
			).Attr("@click",
				web.Plaid().EventFunc(copyCapmaignEvent).Query("id", id).Go(),
			)
		})

	// Регистрируем обработчик события копирования
	mc.RegisterEventFunc(copyCapmaignEvent, m.copyCampaign)

	mce := mc.Editing(
		&presets.FieldsSection{
			Title: "Info",
			Rows: [][]string{
				{"AdvertiserID"},
				{"Title", "Bundle"},
				{"Start", "End"},
				{"Active"},
			},
		},
		&presets.FieldsSection{
			Title: "Targeting",
			Rows: [][]string{
				{"Timetable"},
				{"Targeting"},
			},
		},
		&presets.FieldsSection{
			Title: "Budget",
			Rows: [][]string{
				{"Budget"},
			},
		},
		&presets.FieldsSection{
			Title: "Capping",
			Rows: [][]string{
				{"Capping"},
			},
		},
	)
	mce.Field("AdvertiserID").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Campaign)

		var comps []models.Advertiser
		m.db.Find(&comps)

		sel := v.VSelect().
			Label("Рекламодатель").
			Items(comps).
			ItemTitle("Title").ItemValue("ID").
			Attr(web.VField("AdvertiserID", c.AdvertiserID)...)

		return h.Div(
			sel,
		)
	})

	timetable := components.NewTimetable(m.logger)
	mce.Field("Timetable").
		ComponentFunc(timetable.Component).
		SetterFunc(timetable.Setter)

	targeting := components.NewTargeting(m.logger, m.db)
	mce.Field("Targeting").
		ComponentFunc(targeting.Component).
		SetterFunc(targeting.Setter)

	budget := components.NewBudget(m.logger)
	mce.Field("Budget").
		ComponentFunc(budget.Component).
		SetterFunc(budget.Setter)

	capping := components.NewCapping(m.logger)
	mce.Field("Capping").
		ComponentFunc(capping.Component).
		SetterFunc(capping.Setter)
}

func (m *Campaign) copyCampaign(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Campaign
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find camapign: %w", err)
	}

	// Создаем копию
	if _, err := original.Copy(m.db, original.AdvertiserID); err != nil {
		return r, fmt.Errorf("error on copy banner: %w", err)
	}

	// Обновляем список
	r.Reload = true

	return r, nil
}
