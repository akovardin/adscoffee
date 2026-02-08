package builders

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/qor5/admin/v3/media"
	"github.com/qor5/admin/v3/media/base"
	"github.com/qor5/admin/v3/media/media_library"
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

type Banner struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewBanner(logger *zap.Logger, db *gorm.DB) *Banner {
	return &Banner{
		logger: logger,
		db:     db,
	}
}

const (
	copyBannerEvent      = "copyBanner"
	archiveBannerEvent   = "archiveBanner"
	unarchiveBannerEvent = "unarchiveBanner"
)

func (m *Banner) Configure(b *presets.Builder) {
	mb := b.Model(&models.Banner{}).
		MenuIcon("mdi-image").
		// Label("Креативы").
		RightDrawerWidth("1000")

	mbl := mb.Listing("ID", "Title", "Icon", "Price", "Bgroup", "Active").
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
				qdb := m.db.Where("archived_at is null")
				return gorm2op.DataOperator(qdb).Search(ctx, params)
			} else {
				qdb := m.db.Where("")
				return gorm2op.DataOperator(qdb).Search(ctx, params)
			}
		}).
		SearchColumns("Title").
		// SelectableColumns(true).
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

	mbl.FilterDataFunc(func(ctx *web.EventContext) vuetifyx.FilterData {
		// msgr := i18n.MustGetModuleMessages(ctx.R, presets.ModelsI18nModuleKey, Messages_en_US).(*Messages)
		var companyOptions []*vuetifyx.SelectItem
		err := m.db.Model(&models.Bgroup{}).Select("title as text, id as value").Scan(&companyOptions).Error
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
				Key:          "group",
				Label:        "Группа",
				ItemType:     vuetifyx.ItemTypeSelect,
				SQLCondition: `bgroup_id %s ?`,
				Options:      companyOptions,
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

	mbl.Field("Bgroup").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Banner)
		var group models.Bgroup
		if c.BgroupID == 0 {
			return h.Td()
		}

		m.db.First(&group, "id = ?", c.BgroupID)

		return h.Td().Text(group.Title)
	})

	mbl.Field("Price").Label("CPM")

	mbl.Field("Active").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Banner)

		color := "red"
		text := "выключен"
		if c.Active {
			text = "включен"
			color = "green"
		}

		return h.Td().Children(h.Span(text).Style("color:" + color))
	})

	mbn := mbl.RowMenu()

	// Добавляем обработчик копирования
	mbn.RowMenuItem("Copy").
		ComponentFunc(func(obj interface{}, id string, ctx *web.EventContext) h.HTMLComponent {
			return v.VListItem(
				web.Slot(
					v.VIcon("mdi-content-copy"), // Используем иконку копирования
				).Name("prepend"),
				v.VListItemTitle(
					h.Text("Копировать"),
				),
			).Attr("@click",
				web.Plaid().EventFunc(copyBannerEvent).Query("id", id).Go(),
			)
		})

	mbn.RowMenuItem("Archive").
		ComponentFunc(func(obj interface{}, id string, ctx *web.EventContext) h.HTMLComponent {
			banner := obj.(*models.Banner)
			if banner.ArchivedAt == nil {
				return v.VListItem(
					web.Slot(
						v.VIcon("mdi-archive-arrow-down"), // Используем иконку копирования
					).Name("prepend"),
					v.VListItemTitle(
						h.Text("Архивировать"),
					),
				).Attr("@click",
					web.Plaid().EventFunc(archiveBannerEvent).Query("id", id).Go(),
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
					web.Plaid().EventFunc(unarchiveBannerEvent).Query("id", id).Go(),
				)
			}
		})

	// Регистрируем обработчик события копирования
	mb.RegisterEventFunc(copyBannerEvent, m.copyBanner)
	mb.RegisterEventFunc(archiveBannerEvent, m.archiveBanner)
	mb.RegisterEventFunc(unarchiveBannerEvent, m.unarchiveBanner)

	mbe := mb.Editing(
		&presets.FieldsSection{
			Title: "Info",
			Rows: [][]string{
				{"BgroupID"},
				{"Title"},
				{"Label"},
				{"Description"},
				{"Image", "Icon"},
				{"Active"},
			},
		},
		&presets.FieldsSection{
			Title: "Price",
			Rows: [][]string{
				{"Price"},
				{"ExpectedWinRate"},
			},
		},
		&presets.FieldsSection{
			Title: "Marker",
			Rows: [][]string{
				{"Erid"},
			},
		},
		&presets.FieldsSection{
			Title: "Timetable",
			Rows: [][]string{
				{"Timetable"},
			},
		},
		&presets.FieldsSection{
			Title: "Targeting",
			Rows: [][]string{
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
		&presets.FieldsSection{
			Title: "Tracking",
			Rows: [][]string{
				{"Clicktracker"},
				{"Imptracker"},
				{"Target"},
				{"Macros"},
			},
		},
	)

	mbe.Field("Price").Label("CPM")

	mbe.Field("Macros").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return h.Div(
			h.P(h.Text("Доступные макросы:")),
			h.P(h.Text("{gaid} {adid} {click_id} {ssp} {banner_id} {group_id} {campaign_id} {advertiser_id}")),
		)
	})

	mbe.Field("Image").
		WithContextValue(
			media.MediaBoxConfig,
			&media_library.MediaBoxConfig{
				AllowType: media_library.ALLOW_TYPE_IMAGE,
				Sizes: map[string]*base.Size{
					"image": {
						Width:  640,
						Height: 360,
					},
				},
			})

	mbe.Field("Icon").
		WithContextValue(
			media.MediaBoxConfig,
			&media_library.MediaBoxConfig{
				AllowType: media_library.ALLOW_TYPE_IMAGE,
				Sizes: map[string]*base.Size{
					"image": {
						Width:  64,
						Height: 64,
					},
				},
			})

	mbe.Field("BgroupID").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Banner)

		var comps []models.Bgroup
		m.db.Find(&comps)

		sel := v.VSelect().
			Label("Группа").
			Items(comps).
			ItemTitle("Title").
			ItemValue("ID").
			Attr(web.VField("BgroupID", c.BgroupID)...)

		return h.Div(
			sel,
		)
	})
	mbe.Field("Description").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return v.VTextarea().
			Label(field.Label).
			Attr(web.VField(field.FormKey, fmt.Sprint(reflectutils.MustGet(obj, field.Name)))...).
			Disabled(field.Disabled).
			ErrorMessages(field.Errors...)
	})

	mbe.Field("Erid").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return v.VTextField().
			Label(field.Label).
			Attr(web.VField(field.FormKey, fmt.Sprint(reflectutils.MustGet(obj, field.Name)))...).
			Disabled(field.Disabled).
			ErrorMessages(field.Errors...)
	})

	timetable := components.NewTimetable(m.logger)
	mbe.Field("Timetable").
		ComponentFunc(timetable.Component).
		SetterFunc(timetable.Setter)

	targeting := components.NewTargeting(m.logger, m.db)
	mbe.Field("Targeting").
		ComponentFunc(targeting.Component).
		SetterFunc(targeting.Setter)

	budget := components.NewBudget(m.logger)
	mbe.Field("Budget").
		ComponentFunc(budget.Component).
		SetterFunc(budget.Setter)

	capping := components.NewCapping(m.logger)
	mbe.Field("Capping").
		ComponentFunc(capping.Component).
		SetterFunc(capping.Setter)
}

type Format struct {
	Title string
	Value string
}

func (m *Banner) copyBanner(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Banner
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find banner: %w", err)
	}

	// Создаем копию
	nb, err := original.Copy(m.db, original.BgroupID)
	if err != nil {
		return r, fmt.Errorf("error on copy banner: %w", err)
	}

	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Banner{}),
		presets.PayloadModelsUpdated{Ids: []string{id, strconv.Itoa(int(nb.ID))}},
	)

	return r, nil
}

func (m *Banner) archiveBanner(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Banner
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find banner: %w", err)
	}

	now := time.Now()
	original.ArchivedAt = &now

	m.db.Save(original)

	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Banner{}),
		presets.PayloadModelsUpdated{Ids: []string{id}},
	)

	return r, nil
}

func (m *Banner) unarchiveBanner(ctx *web.EventContext) (r web.EventResponse, err error) {
	id := ctx.R.FormValue("id")
	if id == "" {
		return r, fmt.Errorf("id is required")
	}

	// Находим оригинальную запись
	var original models.Banner
	if err := m.db.First(&original, id).Error; err != nil {
		return r, fmt.Errorf("failed to find banner: %w", err)
	}

	original.ArchivedAt = nil

	m.db.Save(original)

	// Обновляем список
	r.Emit(
		presets.NotifModelsUpdated(&models.Banner{}),
		presets.PayloadModelsUpdated{Ids: []string{id}},
	)

	return r, nil
}
