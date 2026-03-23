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
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.ads.coffee/platform/admin/internal/modules/ads/models"
)

type Network struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewNetwork(logger *zap.Logger, db *gorm.DB) *Network {
	return &Network{
		logger: logger,
		db:     db,
	}
}

const (
	archiveNetworkEvent   = "archiveNetwork"
	unarchiveNetworkEvent = "unarchiveNetwork"
)

func (n *Network) Configure(b *presets.Builder) *presets.ModelBuilder {
	mn := b.Model(&models.Network{}).
		MenuIcon("mdi-lan").
		// Label("Рекламодатели").
		RightDrawerWidth("1000")

	mnl := mn.Listing("ID", "Title", "Name", "Active").
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
				qdb := n.db.Where("archived_at is null")
				return gorm2op.DataOperator(qdb).Search(ctx, params)
			} else {
				qdb := n.db.Where("")
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

	mne := mn.Editing(
		&presets.FieldsSection{
			// Title: "Info",
			Rows: [][]string{
				{"Title"},
				{"Name"},
				{"Data"},
				{"Active"},
			},
		},
	).ValidateFunc(func(obj interface{}, ctx *web.EventContext) (err web.ValidationErrors) {
		u := obj.(*models.Network)
		if u.Title == "" {
			err.FieldError("Title", "Title is required")
		}
		if u.Name == "" {
			err.FieldError("Name", "Name is required")
		}
		return
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
			{
				Key:          "created",
				Label:        "Создан",
				ItemType:     vuetifyx.ItemTypeDate,
				SQLCondition: `created_at %s ?`,
			},
		}
	})

	mnl.Field("Active").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Network)

		color := "red"
		text := "выключен"
		if c.Active {
			text = "включен"
			color = "green"
		}

		return h.Td().Children(h.Span(text).Style("color:" + color))
	})

	mne.Field("Data").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return v.VTextarea().
			Label(field.Label).
			Attr(web.VField(field.FormKey, fmt.Sprint(reflectutils.MustGet(obj, field.Name)))...).
			Disabled(field.Disabled).
			ErrorMessages(field.Errors...)
	})

	mnn := mnl.RowMenu()

	mnn.RowMenuItem("Archive").
		ComponentFunc(func(obj interface{}, id string, ctx *web.EventContext) h.HTMLComponent {
			item := obj.(*models.Network)
			if item.ArchivedAt == nil {
				return v.VListItem(
					web.Slot(
						v.VIcon("mdi-archive-arrow-down"), // Используем иконку копирования
					).Name("prepend"),
					v.VListItemTitle(
						h.Text("Архивировать"),
					),
				).Attr("@click",
					web.Plaid().EventFunc(archiveNetworkEvent).Query("id", id).Go(),
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
					web.Plaid().EventFunc(unarchiveNetworkEvent).Query("id", id).Go(),
				)
			}
		})

	// Регистрируем обработчик события копирования
	mn.RegisterEventFunc(archiveNetworkEvent, n.archive)
	mn.RegisterEventFunc(unarchiveNetworkEvent, n.unarchive)

	return mn
}

func (n *Network) archive(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Network
	if err := n.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find network: %w", err)
	}

	now := time.Now()
	if err := original.Archive(n.db, &now); err != nil {
		return r, fmt.Errorf("failed to archive network: %w", err)
	}

	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Network{}),
		presets.PayloadModelsUpdated{Ids: []string{id}},
	)

	return r, nil
}

func (n *Network) unarchive(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Network
	if err := n.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find network: %w", err)
	}

	if err := original.Archive(n.db, nil); err != nil {
		return r, fmt.Errorf("failed to unarchive network: %w", err)
	}

	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Network{}),
		presets.PayloadModelsUpdated{Ids: []string{id}},
	)

	return r, nil
}
