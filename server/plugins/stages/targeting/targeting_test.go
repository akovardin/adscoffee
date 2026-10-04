package targeting

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type testTargeting struct {
	name string
	keep func(ads.Banner) bool
}

func (t *testTargeting) Name() string {
	return t.name
}

func (t *testTargeting) Copy(cfg map[string]any) plugins.Targeting {
	return t
}

func (t *testTargeting) Filter(candidates []ads.Banner, state *plugins.State) []ads.Banner {
	filtered := make([]ads.Banner, 0, len(candidates))

	for _, candidate := range candidates {
		if t.keep(candidate) {
			filtered = append(filtered, candidate)
		}
	}

	return filtered
}

func TestTargeting_Do(t *testing.T) {
	stage := New()
	stage.Targetings([]plugins.Targeting{
		&testTargeting{name: "first", keep: func(b ads.Banner) bool { return b.ID != 2 }},
		&testTargeting{name: "second", keep: func(b ads.Banner) bool { return b.ID != 3 }},
	})

	state := &plugins.State{
		Candidates: []ads.Banner{{ID: 1}, {ID: 2}, {ID: 3}},
	}

	err := stage.Do(context.Background(), state)
	assert.NoError(t, err)
	assert.Len(t, state.Candidates, 1)
	assert.Equal(t, uint(1), state.Candidates[0].ID)
}

func TestTargeting_Copy(t *testing.T) {
	stage := New()
	stage.Targetings([]plugins.Targeting{&testTargeting{name: "first", keep: func(ads.Banner) bool { return true }}})

	copied, ok := stage.Copy(map[string]any{}).(*Targeting)
	assert.True(t, ok)
	assert.Len(t, copied.targetings, 1)
}
