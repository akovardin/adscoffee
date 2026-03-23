package builders

import (
	"fmt"
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

	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

type Placement struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewPlacement(logger *zap.Logger, db *gorm.DB) *Placement {
	return &Placement{
		logger: logger,
		db:     db,
	}
}

const (
	archivePlacementEvent   = "archivePlacement"
	unarchivePlacementEvent = "unarchivePlacement"
)

func (m *Placement) Configure(b *presets.Builder) {
	mp := b.Model(&models.Placement{}).
		MenuIcon("mdi-lan").
		// Label("Рекламодатели").
		RightDrawerWidth("1000")

	mpl := mp.Listing("ID", "Title", "Active").
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
				FieldName: "Active",
				DBColumn:  "active",
			},
		})

	mp.Editing(
		&presets.FieldsSection{
			// Title: "Info",
			Rows: [][]string{
				{"Title"},
				{"Active"},
			},
		},
	).ValidateFunc(func(obj interface{}, ctx *web.EventContext) (err web.ValidationErrors) {
		u := obj.(*models.Placement)

		if u.Title == "" {
			err.FieldError("Name", "Name is required")
		}
		return
	})

	mpl.FilterDataFunc(func(ctx *web.EventContext) vuetifyx.FilterData {
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

	mpl.Field("Active").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Placement)

		color := "red"
		text := "выключен"
		if c.Active {
			text = "включен"
			color = "green"
		}

		return h.Td().Children(h.Span(text).Style("color:" + color))
	})

	mpn := mpl.RowMenu()

	mpn.RowMenuItem("Archive").
		ComponentFunc(func(obj interface{}, id string, ctx *web.EventContext) h.HTMLComponent {
			item := obj.(*models.Placement)
			if item.ArchivedAt == nil {
				return v.VListItem(
					web.Slot(
						v.VIcon("mdi-archive-arrow-down"), // Используем иконку копирования
					).Name("prepend"),
					v.VListItemTitle(
						h.Text("Архивировать"),
					),
				).Attr("@click",
					web.Plaid().EventFunc(archivePlacementEvent).Query("id", id).Go(),
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
					web.Plaid().EventFunc(unarchivePlacementEvent).Query("id", id).Go(),
				)
			}
		})

	// Регистрируем обработчик события копирования
	mp.RegisterEventFunc(archivePlacementEvent, m.archive)
	mp.RegisterEventFunc(unarchivePlacementEvent, m.unarchive)
}

func (m *Placement) archive(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Placement
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find placement: %w", err)
	}

	now := time.Now()
	if err := original.Archive(m.db, &now); err != nil {
		return r, fmt.Errorf("failed to archive placement: %w", err)
	}

	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Placement{}),
		presets.PayloadModelsUpdated{Ids: []string{id}},
	)

	return r, nil
}

func (m *Placement) unarchive(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Placement
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find placement: %w", err)
	}

	if err := original.Archive(m.db, nil); err != nil {
		return r, fmt.Errorf("failed to unarchive placement: %w", err)
	}

	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Placement{}),
		presets.PayloadModelsUpdated{Ids: []string{id}},
	)

	return r, nil
}
