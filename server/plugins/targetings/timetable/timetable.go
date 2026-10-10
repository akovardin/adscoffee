package timetable

import (
	"time"

	"go.uber.org/fx"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

var Module = fx.Module(
	"targetings.timetable",
	fx.Provide(
		fx.Annotate(
			New,
			fx.ResultTags(`group:"targetings"`),
		),
	),
)

// Timetable фильтрует кандидатов по расписанию активности (день недели + час)
// относительно времени сервера.
type Timetable struct {
	now func() time.Time
}

func New() plugins.Provider[plugins.Targeting] {
	return plugins.Provider[plugins.Targeting]{
		Name: "targetings.timetable",
		Factory: func(cfg map[string]any) (plugins.Targeting, error) {
			return &Timetable{now: time.Now}, nil
		},
	}
}

func (t *Timetable) Name() string {
	return "targetings.timetable"
}

func (t *Timetable) Filter(candidates []ads.Banner, state *plugins.State) []ads.Banner {
	now := t.now()

	// приводим Go weekday (Вс=0) к конвенции админки (Пн=0..Вс=6)
	day := (int(now.Weekday()) + 6) % 7
	hour := now.Hour()

	filtered := make([]ads.Banner, 0, len(candidates))

	for _, candidate := range candidates {
		if candidate.Timetable.Validate(day, hour) {
			filtered = append(filtered, candidate)
		}
	}

	return filtered
}
