package builders

import (
	"fmt"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	v "github.com/qor5/x/v3/ui/vuetify"
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

const copyAdvertiserEvent = "copyAdvertiser"

func (m *Advertiser) Configure(b *presets.Builder) {
	ma := b.Model(&models.Advertiser{}).
		MenuIcon("mdi-account-group").
		// Label("Рекламодатели").
		RightDrawerWidth("1000")

	mal := ma.Listing("ID", "Title", "Start", "End", "Active")

	mal.Field("Title").ComponentFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		c := obj.(*models.Advertiser)
		return h.Td().Children(
			h.A().
				Text(c.Title).
				Attr("onclick", "event.stopPropagation();").
				Href(fmt.Sprintf("/admin/campaigns?f_advertiser=%d", c.ID)),
		)
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

	// Регистрируем обработчик события копирования
	ma.RegisterEventFunc(copyAdvertiserEvent, m.copyAdvertiser)

	mae := ma.Editing(
		&presets.FieldsSection{
			Title: "Info",
			Rows: [][]string{
				{"Title"},
				{"Info"},
				{"OrdContract"},
				{"Start", "End"},
				{"Active"},
				{"OrdEnable"},
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

	targeting := components.NewTargeting(m.logger, m.db)
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
		OrdEnable:   original.OrdEnable,
	}

	// Сохраняем копию в базу данных
	if err := m.db.Create(&copyAdvertiser).Error; err != nil {
		return r, fmt.Errorf("failed to create copy: %w", err)
	}

	// Обновляем список
	r.Reload = true

	return r, nil
}
