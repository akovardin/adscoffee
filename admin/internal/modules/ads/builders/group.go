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

type Group struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewGroup(logger *zap.Logger, db *gorm.DB) *Group {
	return &Group{
		logger: logger,
		db:     db,
	}
}

// Константа для имени события копирования
const (
	copyGroupEvent      = "copyGroup"
	archiveGroupEvent   = "archiveGroup"
	unarchiveGroupEvent = "unarchiveGroup"
)

func (m *Group) Configure(b *presets.Builder) {
	mg := b.Model(&models.Bgroup{}).
		MenuIcon("mdi-lightbulb-group").
		// Label("Группы").
		RightDrawerWidth("1000")

	mgl := mg.Listing("ID", "Title", "Price", "Start", "End", "Campaign", "Active").
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

	mgl.FilterDataFunc(func(ctx *web.EventContext) vuetifyx.FilterData {
		// msgr := i18n.MustGetModuleMessages(ctx.R, presets.ModelsI18nModuleKey, Messages_en_US).(*Messages)
		var options []*vuetifyx.SelectItem
		err := m.db.Model(&models.Campaign{}).Select("title as text, id as value").Scan(&options).Error
		if err != nil {
			panic(err)
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
				Key:          "campaign",
				Label:        "Кампания",
				ItemType:     vuetifyx.ItemTypeSelect,
				SQLCondition: `campaign_id %s ?`,
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

	mgl.Field("Campaign").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Bgroup)
		var comp models.Campaign
		if c.CampaignID == 0 {
			return h.Td()
		}

		m.db.First(&comp, "id = ?", c.CampaignID)

		return h.Td().Children(
			h.A().
				Text(comp.Title).
				Attr("onclick", "event.stopPropagation();").
				Href(fmt.Sprintf("/admin/bgroups?f_campaign=%d", c.CampaignID)),
		)
	})

	mgl.Field("Title").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Bgroup)

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
				Href(fmt.Sprintf("/admin/banners?f_group=%d", c.ID)),
		)
	})

	mgl.Field("Active").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Bgroup)

		color := "red"
		text := "выключен"
		if c.Active {
			text = "включен"
			color = "green"
		}

		return h.Td().Children(h.Span(text).Style("color:" + color))
	})

	mgl.Field("Price").Label("CPM")

	mgn := mgl.RowMenu()

	// Добавляем обработчик копирования
	mgn.RowMenuItem("Copy").
		ComponentFunc(func(obj interface{}, id string, ctx *web.EventContext) h.HTMLComponent {
			return v.VListItem(
				web.Slot(
					v.VIcon("mdi-content-copy"), // Используем иконку копирования
				).Name("prepend"),
				v.VListItemTitle(
					h.Text("Копировать"),
				),
			).Attr("@click",
				web.Plaid().EventFunc(copyGroupEvent).Query("id", id).Go(),
			)
		})

	mgn.RowMenuItem("Archive").
		ComponentFunc(func(obj interface{}, id string, ctx *web.EventContext) h.HTMLComponent {
			item := obj.(*models.Bgroup)
			if item.ArchivedAt == nil {
				return v.VListItem(
					web.Slot(
						v.VIcon("mdi-archive-arrow-down"), // Используем иконку копирования
					).Name("prepend"),
					v.VListItemTitle(
						h.Text("Архивировать"),
					),
				).Attr("@click",
					web.Plaid().EventFunc(archiveGroupEvent).Query("id", id).Go(),
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
					web.Plaid().EventFunc(unarchiveGroupEvent).Query("id", id).Go(),
				)
			}
		})

	// Регистрируем обработчик события копирования
	mg.RegisterEventFunc(copyGroupEvent, m.copyGroup)
	mg.RegisterEventFunc(archiveGroupEvent, m.archiveGroup)
	mg.RegisterEventFunc(unarchiveGroupEvent, m.unarchiveGroup)

	mge := mg.Editing(
		&presets.FieldsSection{
			Title: "Info",
			Rows: [][]string{
				{"CampaignID"},
				{"Title", "Price"},
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

	mge.Field("Price").Label("CPM")

	mge.Field("CampaignID").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Bgroup)

		var comps []models.Campaign
		m.db.Find(&comps)

		sel := v.VSelect().
			Label("Кампания").
			Items(comps).
			ItemTitle("Title").ItemValue("ID").
			Attr(web.VField("CampaignID", c.CampaignID)...)

		return h.Div(
			sel,
		)
	})
	timetable := components.NewTimetable(m.logger)
	mge.Field("Timetable").
		ComponentFunc(timetable.Component).
		SetterFunc(timetable.Setter)

	targeting := components.NewTargeting(m.logger, m.db)
	mge.Field("Targeting").
		ComponentFunc(targeting.Component).
		SetterFunc(targeting.Setter)

	budget := components.NewBudget(m.logger)
	mge.Field("Budget").
		ComponentFunc(budget.Component).
		SetterFunc(budget.Setter)

	capping := components.NewCapping(m.logger)
	mge.Field("Capping").
		ComponentFunc(capping.Component).
		SetterFunc(capping.Setter)
}

// Функция для копирования группы
func (m *Group) copyGroup(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Bgroup
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find group: %w", err)
	}

	// Создаем копию
	ng, err := original.Copy(m.db, original.CampaignID)
	if err != nil {
		return r, fmt.Errorf("failed on copy group: %w", err)
	}

	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Bgroup{}),
		presets.PayloadModelsUpdated{Ids: []string{id, strconv.Itoa(int(ng.ID))}},
	)

	return r, nil
}

func (m *Group) archiveGroup(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Bgroup
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find group: %w", err)
	}

	now := time.Now()
	if err := original.Archive(m.db, &now); err != nil {
		return r, fmt.Errorf("failed to archive group: %w", err)
	}

	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Bgroup{}),
		presets.PayloadModelsUpdated{Ids: []string{id}},
	)

	return r, nil
}

func (m *Group) unarchiveGroup(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Bgroup
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find group: %w", err)
	}

	if err := original.Archive(m.db, nil); err != nil {
		return r, fmt.Errorf("failed to unarchive group: %w", err)
	}

	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Bgroup{}),
		presets.PayloadModelsUpdated{Ids: []string{id}},
	)

	return r, nil
}
