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
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/modules/ads/components"
	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

type Advertiser struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewAdvertiser(logger *zap.Logger, db *gorm.DB) *Advertiser {
	return &Advertiser{
		logger: logger,
		db:     db,
	}
}

const (
	copyAdvertiserEvent      = "copyAdvertiser"
	archiveAdvertiserEvent   = "archiveAdvertiser"
	unarchiveAdvertiserEvent = "unarchiveAdvertiser"
)

func (m *Advertiser) Configure(b *presets.Builder) {
	ma := b.Model(&models.Advertiser{}).
		MenuIcon("mdi-account-group").
		// Label("Рекламодатели").
		RightDrawerWidth("1000")

	mal := ma.Listing("ID", "Title", "Start", "End", "Active").
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

	mal.FilterDataFunc(func(ctx *web.EventContext) vuetifyx.FilterData {
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
			{
				Key:          "created",
				Label:        "Создан",
				ItemType:     vuetifyx.ItemTypeDate,
				SQLCondition: `created_at %s ?`,
			},
		}
	})

	mal.Field("Title").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Advertiser)

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
				Href(fmt.Sprintf("/admin/campaigns?f_advertiser=%d", c.ID)),
		)
	})

	mal.Field("Active").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Advertiser)

		color := "red"
		text := "выключен"
		if c.Active {
			text = "включен"
			color = "green"
		}

		return h.Td().Children(h.Span(text).Style("color:" + color))
	})

	man := mal.RowMenu()

	// Добавляем обработчик копирования
	man.RowMenuItem("Copy").
		ComponentFunc(func(obj interface{}, id string, ctx *web.EventContext) h.HTMLComponent {
			return v.VListItem(
				web.Slot(
					v.VIcon("mdi-content-copy"), // Используем иконку копирования
				).Name("prepend"),
				v.VListItemTitle(
					h.Text("Копировать"),
				),
			).Attr("@click",
				web.Plaid().EventFunc(copyAdvertiserEvent).Query("id", id).Go(),
			)
		})

	man.RowMenuItem("Archive").
		ComponentFunc(func(obj interface{}, id string, ctx *web.EventContext) h.HTMLComponent {
			item := obj.(*models.Advertiser)
			if item.ArchivedAt == nil {
				return v.VListItem(
					web.Slot(
						v.VIcon("mdi-archive-arrow-down"), // Используем иконку копирования
					).Name("prepend"),
					v.VListItemTitle(
						h.Text("Архивировать"),
					),
				).Attr("@click",
					web.Plaid().EventFunc(archiveAdvertiserEvent).Query("id", id).Go(),
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
					web.Plaid().EventFunc(unarchiveAdvertiserEvent).Query("id", id).Go(),
				)
			}
		})

	// Регистрируем обработчик события копирования
	ma.RegisterEventFunc(copyAdvertiserEvent, m.copyAdvertiser)
	ma.RegisterEventFunc(archiveAdvertiserEvent, m.archiveAdvertiser)
	ma.RegisterEventFunc(unarchiveAdvertiserEvent, m.unarchiveAdvertiser)

	mae := ma.Editing(
		&presets.FieldsSection{
			Title: "Info",
			Rows: [][]string{
				{"Title"},
				{"Info"},
				{"OrdContract"},
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
	mae.Field("Info").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return v.VTextarea().
			Label(field.Label).
			Attr(web.VField(field.FormKey, fmt.Sprint(reflectutils.MustGet(obj, field.Name)))...).
			Disabled(field.Disabled).
			ErrorMessages(field.Errors...)
	})
	timetable := components.NewTimetable(m.logger)
	mae.Field("Timetable").
		ComponentFunc(timetable.Component).
		SetterFunc(timetable.Setter)

	targeting := components.NewTargeting(m.logger)
	mae.Field("Targeting").
		ComponentFunc(targeting.Component).
		SetterFunc(targeting.Setter)

	budget := components.NewBudget(m.logger)
	mae.Field("Budget").
		ComponentFunc(budget.Component).
		SetterFunc(budget.Setter)

	capping := components.NewCapping(m.logger)
	mae.Field("Capping").
		ComponentFunc(capping.Component).
		SetterFunc(capping.Setter)
}

func (m *Advertiser) copyAdvertiser(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Advertiser
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find advertiser: %w", err)
	}

	// Создаем копию
	copyAdvertiser := models.Advertiser{
		Title:     original.Title + " (Копия)",
		Info:      original.Info,
		Start:     original.Start,
		End:       original.End,
		Timetable: original.Timetable,
		Targeting: original.Targeting,
		Budget:    original.Budget,
		Capping:   original.Capping,
		Active:    false,

		OrdContract: original.OrdContract,
	}

	// Сохраняем копию в базу данных
	if err := m.db.Create(&copyAdvertiser).Error; err != nil {
		return r, fmt.Errorf("failed to create copy: %w", err)
	}

	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Advertiser{}),
		presets.PayloadModelsUpdated{Ids: []string{id, strconv.Itoa(int(copyAdvertiser.ID))}},
	)

	return r, nil
}

func (m *Advertiser) archiveAdvertiser(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Advertiser
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find campaign: %w", err)
	}

	now := time.Now()
	if err := original.Archive(m.db, &now); err != nil {
		return r, fmt.Errorf("failed to archive campaign: %w", err)
	}
	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Advertiser{}),
		presets.PayloadModelsUpdated{Ids: []string{id}},
	)

	return r, nil
}

func (m *Advertiser) unarchiveAdvertiser(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Advertiser
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find advertiser: %w", err)
	}

	if err := original.Archive(m.db, nil); err != nil {
		return r, fmt.Errorf("failed to unarchive advertiser: %w", err)
	}

	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Advertiser{}),
		presets.PayloadModelsUpdated{Ids: []string{id}},
	)

	return r, nil
}
