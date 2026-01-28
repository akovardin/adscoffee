package builders

import (
	"fmt"

	"github.com/qor5/admin/v3/media"
	"github.com/qor5/admin/v3/media/base"
	"github.com/qor5/admin/v3/media/media_library"
	"github.com/qor5/admin/v3/presets"
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

const copyBannerEvent = "copyBanner"

func (m *Banner) Configure(b *presets.Builder) {
	mb := b.Model(&models.Banner{}).
		MenuIcon("mdi-image").
		// Label("Креативы").
		RightDrawerWidth("1000")

	mbl := mb.Listing("ID", "Title", "Icon", "Price", "Bgroup", "Active")

	mbl.FilterDataFunc(func(ctx *web.EventContext) vuetifyx.FilterData {
		// msgr := i18n.MustGetModuleMessages(ctx.R, presets.ModelsI18nModuleKey, Messages_en_US).(*Messages)
		var companyOptions []*vuetifyx.SelectItem
		err := m.db.Model(&models.Bgroup{}).Select("title as text, id as value").Scan(&companyOptions).Error
		if err != nil {
			panic(err)
		}

		return []*vuetifyx.FilterItem{
			{
				Key:          "group",
				Label:        "Группа",
				ItemType:     vuetifyx.ItemTypeSelect,
				SQLCondition: `bgroup_id %s ?`,
				Options:      companyOptions,
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

	// Регистрируем обработчик события копирования
	mb.RegisterEventFunc(copyBannerEvent, m.copyBanner)

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

	mb.EventsHub.RegisterEventFunc("erid", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		banner := models.Banner{}

		id := ctx.R.URL.Query().Get("id")
		targeting := ctx.R.FormValue("Ord.Targeting")
		category := ctx.R.FormValue("Ord.Category")
		format := ctx.R.FormValue("Ord.Format")
		kktu := ctx.R.FormValue("Ord.Kktu")

		m.db.Model(models.Banner{}).Where("id = ?", id).Preload("Bgroup.Campaign.Advertiser").Preload("Campaign").First(&banner)

		if !banner.Bgroup.Campaign.Advertiser.OrdEnable {
			return web.EventResponse{
				Reload: true,
			}, nil
		}

		banner.OrdCategory = category
		banner.OrdTargeting = targeting
		banner.OrdFormat = format
		banner.OrdKktu = kktu

		targets := []string{}
		if banner.Clicktracker != "" {
			targets = append(targets, banner.Clicktracker)
		}
		if banner.Target != "" {
			targets = append(targets, banner.Target)
		}

		if err := m.db.Save(&banner).Error; err != nil {
			m.logger.Error("error on save banner", zap.Error(err))

			return web.EventResponse{
				Reload: true,
			}, fmt.Errorf("error on save banner: %w", err)
		}

		return web.EventResponse{
			Reload: true,
		}, nil
	})

	mbe.AppendTabsPanelFunc(func(obj interface{}, ctx *web.EventContext) (tab h.HTMLComponent, content h.HTMLComponent) {
		c := obj.(*models.Banner)

		tab = v.VTab(h.Text("ОРД")).Value("2")
		if c.ID == 0 {
			content = v.VTabsWindowItem(
				h.Text("Сначала нужно создать баннер"),
			).Value("2").Class("pa-4")

			return
		}

		banner := models.Banner{}
		m.db.Model(banner).Where("id = ?", c.ID).Preload("Bgroup.Campaign.Advertiser").Preload("Campaign").First(&banner)

		if !banner.Bgroup.Campaign.Advertiser.OrdEnable {
			content = v.VTabsWindowItem(
				h.Text("На уровне рекламодателя выключена работа с ОРД"),
			).Value("2").Class("pa-4")

			return
		}

		formats := []Format{
			{
				Value: "banner",
				Title: "Баннер",
			},
			{
				Value: "text_block",
				Title: "Текстовый блок",
			},
			{
				Value: "text_graphic_block",
				Title: "Текстовый-графический блок",
			},
			{
				Value: "banner_html5",
				Title: "HTML5-баннер",
			},
		}

		content = v.VTabsWindowItem(
			v.VCard(
				h.Div(v.VCol([]h.HTMLComponent{
					v.VRow(
						[]h.HTMLComponent{
							v.VTextField().
								Label("Категория").
								Hint("Реклама мобильных приложений").
								Variant("outlined").Density("compact").
								Attr(web.VField("Ord.Category", c.OrdCategory)...),
						}...,
					),
					v.VRow(
						[]h.HTMLComponent{
							v.VTextField().
								Label("Таргетинг").
								Hint("Жители России").
								Variant("outlined").Density("compact").
								Attr(web.VField("Ord.Targeting", c.OrdTargeting)...),
						}...,
					),
					v.VRow(
						[]h.HTMLComponent{
							v.VTextField().
								Label("ККТУ").
								Hint("30.15.1"). // "30.15.1"
								Variant("outlined").Density("compact").
								Attr(web.VField("Ord.Kktu", c.OrdKktu)...),
						}...,
					),
					v.VRow(
						[]h.HTMLComponent{
							v.VSelect().
								Label("Формат креатива").
								Items(formats).
								ItemTitle("Title").ItemValue("Value").
								Variant("outlined").Density("compact").
								Attr(web.VField("Ord.Format", c.OrdFormat)...),
						}...,
					),

					v.VCardActions(
						v.VSpacer(),
						v.VBtn("Сгенерировать").Attr("@click",
							web.Plaid().EventFunc("erid").Query("id", c.ID).Go(),
						).Variant(v.VariantFlat).Class("bg-primary"),
					),
				}...),
				),
			).Variant(v.VariantFlat).Class("mx-0 mt-1 px-4 pb-0 pt-4 section-body"),
		).Value("2").Class("pa-4")

		return
	})

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
	if _, err := original.Copy(m.db, original.BgroupID); err != nil {
		return r, fmt.Errorf("error on copy banner: %w", err)
	}

	// Обновляем список
	r.Reload = true

	return r, nil
}
