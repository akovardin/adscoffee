package builders

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/presets/gorm2op"
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

const (
	copyCapmaignEvent      = "copyCampaign"
	archiveCampaignEvent   = "archiveCampaign"
	unarchiveCampaignEvent = "unarchiveCampaign"
)

func (m *Campaign) Configure(b *presets.Builder) {
	mc := b.Model(&models.Campaign{}).
		MenuIcon("mdi-bullseye-arrow").
		// Label("Кампании").
		RightDrawerWidth("1000")

	mcl := mc.Listing("ID", "Title", "Bundle", "Start", "End", "Advertiser", "Active").
		SearchFunc(func(ctx *web.EventContext, params *presets.SearchParams) (result *presets.SearchResult, err error) {
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
				qdb := m.db.Where("archived_at is null")
				return gorm2op.DataOperator(qdb).Search(ctx, params)
			} else {
				qdb := m.db.Where("")
				return gorm2op.DataOperator(qdb).Search(ctx, params)
			}
		}).
		SearchColumns("Title").
		SelectableColumns(true).
		OrderableFields([]*presets.OrderableField{
			{
				FieldName: "ID",
				DBColumn:  "id",
			},
			{
				FieldName: "Title",
				DBColumn:  "title",
			},
			{
				FieldName: "Start",
				DBColumn:  "start",
			},
			{
				FieldName: "End",
				DBColumn:  "end",
			},
		})

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
				Key:          "advertiser",
				Label:        "Рекламодатель",
				ItemType:     vuetifyx.ItemTypeSelect,
				SQLCondition: `advertiser_id %s ?`,
				Options:      options,
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
			{
				Key:          "created",
				Label:        "Создан",
				ItemType:     vuetifyx.ItemTypeDate,
				SQLCondition: `created_at %s ?`,
			},
		}
	})

	mcl.Field("Advertiser").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Campaign)
		var adv models.Advertiser
		if c.AdvertiserID == 0 {
			return h.Td()
		}

		m.db.First(&adv, "id = ?", c.AdvertiserID)

		return h.Td().Children(
			h.A().
				Text(adv.Title).
				Attr("onclick", "event.stopPropagation();").
				Href(fmt.Sprintf("/admin/campaigns?f_advertiser=%d", c.AdvertiserID)),
		)
	})

	mcl.Field("Title").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Campaign)

		style := ""
		text := ""
		if c.ArchivedAt != nil {
			style = "color:#bb0"
			text = " - архив"
		}

		return h.Td().Children(
			h.A().
				Text(c.Title+text).
				Style(style).
				Attr("onclick", "event.stopPropagation();").
				Href(fmt.Sprintf("/admin/bgroups?f_campaign=%d", c.ID)),
		)
	})

	mcl.Field("Active").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Campaign)

		color := "red"
		text := "выключен"
		if c.Active {
			text = "включен"
			color = "green"
		}

		return h.Td().Children(h.Span(text).Style("color:" + color))
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

	mcn.RowMenuItem("Archive").
		ComponentFunc(func(obj interface{}, id string, ctx *web.EventContext) h.HTMLComponent {
			item := obj.(*models.Campaign)
			if item.ArchivedAt == nil {
				return v.VListItem(
					web.Slot(
						v.VIcon("mdi-archive-arrow-down"), // Используем иконку копирования
					).Name("prepend"),
					v.VListItemTitle(
						h.Text("Архивировать"),
					),
				).Attr("@click",
					web.Plaid().EventFunc(archiveCampaignEvent).Query("id", id).Go(),
				)
			} else {
				return v.VListItem(
					web.Slot(
						v.VIcon("mdi-archive-arrow-up"), // Используем иконку копирования
					).Name("prepend"),
					v.VListItemTitle(
						h.Text("Разархивировать"),
					),
				).Attr("@click",
					web.Plaid().EventFunc(unarchiveCampaignEvent).Query("id", id).Go(),
				)
			}
		})

	// Регистрируем обработчик события копирования
	mc.RegisterEventFunc(copyCapmaignEvent, m.copyCampaign)
	mc.RegisterEventFunc(archiveCampaignEvent, m.archiveCampaign)
	mc.RegisterEventFunc(unarchiveCampaignEvent, m.unarchiveCampaign)

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

	targeting := components.NewTargeting(m.logger)
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
	nc, err := original.Copy(m.db, original.AdvertiserID)
	if err != nil {
		return r, fmt.Errorf("error on copy cmpaign: %w", err)
	}

	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Campaign{}),
		presets.PayloadModelsUpdated{Ids: []string{id, strconv.Itoa(int(nc.ID))}},
	)

	return r, nil
}

func (m *Campaign) archiveCampaign(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Campaign
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find campaign: %w", err)
	}

	now := time.Now()
	if err := original.Archive(m.db, &now); err != nil {
		return r, fmt.Errorf("failed to archive campaign: %w", err)
	}
	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Campaign{}),
		presets.PayloadModelsUpdated{Ids: []string{id}},
	)

	return r, nil
}

func (m *Campaign) unarchiveCampaign(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Campaign
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find advertiser: %w", err)
	}

	if err := original.Archive(m.db, nil); err != nil {
		return r, fmt.Errorf("failed to unarchive campaign: %w", err)
	}

	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Campaign{}),
		presets.PayloadModelsUpdated{Ids: []string{id}},
	)

	return r, nil
}
