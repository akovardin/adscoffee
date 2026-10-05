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

type Site struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewSite(logger *zap.Logger, db *gorm.DB) *Site {
	return &Site{
		logger: logger,
		db:     db,
	}
}

const (
	archiveSiteEvent   = "archiveSite"
	unarchiveSiteEvent = "unarchiveSite"
)

func (m *Site) Configure(b *presets.Builder) {
	mp := b.Model(&models.Site{}).
		URIName("sites").
		MenuIcon("mdi-web").
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
		SearchColumns("Title")

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
		}
	})

	mpl.Field("Title").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Site)

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
				Href(fmt.Sprintf("/placements?f_site=%d", c.ID)),
		)
	})

	mpl.Field("Active").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Site)

		color := "red"
		text := "выключен"
		if c.Active {
			text = "включен"
			color = "green"
		}

		return h.Td().Children(h.Span(text).Style("color:" + color))
	})

	mp.Editing(
		&presets.FieldsSection{
			Rows: [][]string{
				{"Title"},
				{"Active"},
			},
		},
	).ValidateFunc(func(obj interface{}, ctx *web.EventContext) (err web.ValidationErrors) {
		u := obj.(*models.Site)

		if u.Title == "" {
			err.FieldError("Name", "Name is required")
		}
		return
	})

	mpn := mpl.RowMenu()

	mpn.RowMenuItem("Archive").
		ComponentFunc(func(obj interface{}, id string, ctx *web.EventContext) h.HTMLComponent {
			item := obj.(*models.Site)
			if item.ArchivedAt == nil {
				return v.VListItem(
					web.Slot(
						v.VIcon("mdi-archive-arrow-down"),
					).Name("prepend"),
					v.VListItemTitle(
						h.Text("Архивировать"),
					),
				).Attr("@click",
					web.Plaid().EventFunc(archiveSiteEvent).Query("id", id).Go(),
				)
			} else {
				return v.VListItem(
					web.Slot(
						v.VIcon("mdi-archive-arrow-up"),
					).Name("prepend"),
					v.VListItemTitle(
						h.Text("Разархивировать"),
					),
				).Attr("@click",
					web.Plaid().EventFunc(unarchiveSiteEvent).Query("id", id).Go(),
				)
			}
		})

	mp.RegisterEventFunc(archiveSiteEvent, m.archive)
	mp.RegisterEventFunc(unarchiveSiteEvent, m.unarchive)
}

func (m *Site) archive(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	var original models.Site
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find site: %w", err)
	}

	now := time.Now()
	if err := original.Archive(m.db, &now); err != nil {
		return r, fmt.Errorf("failed to archive site: %w", err)
	}

	r.Emit(
		presets.NotifModelsUpdated(&models.Site{}),
		presets.PayloadModelsUpdated{Ids: []string{id}},
	)

	return r, nil
}

func (m *Site) unarchive(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	var original models.Site
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find site: %w", err)
	}

	if err := original.Archive(m.db, nil); err != nil {
		return r, fmt.Errorf("failed to unarchive site: %w", err)
	}

	r.Emit(
		presets.NotifModelsUpdated(&models.Site{}),
		presets.PayloadModelsUpdated{Ids: []string{id}},
	)

	return r, nil
}
