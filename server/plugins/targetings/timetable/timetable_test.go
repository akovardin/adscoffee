package timetable

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

func TestTimetable_Name(t *testing.T) {
	assert.Equal(t, "targetings.timetable", New().Name())
}

func TestTimetable_Copy(t *testing.T) {
	now := func() time.Time { return time.Date(2026, time.October, 5, 10, 0, 0, 0, time.UTC) }
	tt := &Timetable{now: now}

	copied, ok := tt.Copy(map[string]any{}).(*Timetable)
	assert.True(t, ok)
	assert.NotNil(t, copied.now)
}

func TestTimetable_Filter(t *testing.T) {
	// 5 октября 2026 — понедельник, 10:00
	now := func() time.Time { return time.Date(2026, time.October, 5, 10, 0, 0, 0, time.UTC) }

	tt := &Timetable{now: now}

	candidates := []ads.Banner{
		{ID: 1, Timetable: ads.Timetable{0: {10: true}}},  // пн 10:00 — проходит
		{ID: 2, Timetable: ads.Timetable{0: {10: false}}}, // пн 10:00 выключен — отсев
		{ID: 3, Timetable: ads.Timetable{1: {10: true}}},  // вт — отсев
		{ID: 4, Timetable: ads.Timetable{}},               // расписание не задано — проходит
	}

	filtered := tt.Filter(candidates, &plugins.State{})

	assert.Len(t, filtered, 2)
	assert.Equal(t, uint(1), filtered[0].ID)
	assert.Equal(t, uint(4), filtered[1].ID)
}

func TestTimetable_FilterSunday(t *testing.T) {
	// 4 октября 2026 — воскресенье, должно мапиться в admin-индекс 6
	now := func() time.Time { return time.Date(2026, time.October, 4, 12, 0, 0, 0, time.UTC) }

	tt := &Timetable{now: now}

	candidates := []ads.Banner{
		{ID: 1, Timetable: ads.Timetable{6: {12: true}}},
	}

	filtered := tt.Filter(candidates, &plugins.State{})

	assert.Len(t, filtered, 1)
	assert.Equal(t, uint(1), filtered[0].ID)
}
