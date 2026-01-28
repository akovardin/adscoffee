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
const copyGroupEvent = "copyGroup"

func (m *Group) Configure(b *presets.Builder) {
	mg := b.Model(&models.Bgroup{}).
		MenuIcon("mdi-lightbulb-group").
		// Label("Группы").
		RightDrawerWidth("1000")

	mgl := mg.Listing("ID", "Title", "Price", "Start", "End", "Campaign", "Active")

	mgl.FilterDataFunc(func(ctx *web.EventContext) vuetifyx.FilterData {
		// msgr := i18n.MustGetModuleMessages(ctx.R, presets.ModelsI18nModuleKey, Messages_en_US).(*Messages)
		var options []*vuetifyx.SelectItem
		err := m.db.Model(&models.Campaign{}).Select("title as text, id as value").Scan(&options).Error
		if err != nil {
			panic(err)
		}

		return []*vuetifyx.FilterItem{
			{
				Key:          "campaign",
				Label:        "Кампания",
				ItemType:     vuetifyx.ItemTypeSelect,
				SQLCondition: `campaign_id %s ?`,
				Options:      options,
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

		return h.Td().Text(comp.Title)
	})

	mgl.Field("Title").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Bgroup)
		return h.Td().Children(
			h.A().
				Text(c.Title).
				Attr("onclick", "event.stopPropagation();").
				Href(fmt.Sprintf("/admin/banners?f_group=%d", c.ID)),
		)
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

	// Регистрируем обработчик события копирования
	mg.RegisterEventFunc(copyGroupEvent, m.copyGroup)

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
	if _, err := original.Copy(m.db, original.CampaignID); err != nil {
		return r, fmt.Errorf("failed to create copy: %w", err)
	}

	// Обновляем список
	r.Reload = true

	return r, nil
}
