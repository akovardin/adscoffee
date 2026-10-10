package targeting

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

func (t *testTargeting) Filter(candidates []ads.Banner, state *plugins.State) []ads.Banner {
	filtered := make([]ads.Banner, 0, len(candidates))

	for _, candidate := range candidates {
		if t.keep(candidate) {
			filtered = append(filtered, candidate)
		}
	}

	return filtered
}

func newStage(t *testing.T) *Targeting {
	t.Helper()

	v, err := New().Factory(nil)
	require.NoError(t, err)

	return v.(*Targeting)
}

func TestTargeting_Do(t *testing.T) {
	stage := newStage(t)
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

func TestTargeting_New(t *testing.T) {
	v, err := New().Factory(map[string]any{})
	require.NoError(t, err)

	stage, ok := v.(*Targeting)
	require.True(t, ok)
	assert.Equal(t, "stages.targeting", stage.Name())
}
