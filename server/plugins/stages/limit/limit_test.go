package limit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type testLimit struct {
	name string
	drop func(*ads.Banner) bool
}

func (t *testLimit) Name() string {
	return t.name
}

func (t *testLimit) Apply(candidates []*ads.Banner, state *plugins.State) []*ads.Banner {
	kept := make([]*ads.Banner, 0, len(candidates))

	for _, candidate := range candidates {
		if !t.drop(candidate) {
			kept = append(kept, candidate)
		}
	}

	return kept
}

func newStage(t *testing.T) *Limit {
	t.Helper()

	v, err := New().Factory(nil)
	require.NoError(t, err)

	return v.(*Limit)
}

func TestLimit_Do(t *testing.T) {
	stage := newStage(t)
	stage.Limits([]plugins.Limit{
		&testLimit{name: "budget", drop: func(b *ads.Banner) bool { return b.ID == 2 }},
		&testLimit{name: "capping", drop: func(b *ads.Banner) bool { return b.ID == 3 }},
	})

	state := &plugins.State{
		Candidates: []*ads.Banner{{ID: 1}, {ID: 2}, {ID: 3}},
	}

	err := stage.Do(context.Background(), state)
	assert.NoError(t, err)
	assert.Len(t, state.Candidates, 1)
	assert.Equal(t, uint(1), state.Candidates[0].ID)
}

func TestLimit_New(t *testing.T) {
	v, err := New().Factory(map[string]any{})
	require.NoError(t, err)

	stage, ok := v.(*Limit)
	require.True(t, ok)
	assert.Equal(t, "stages.limit", stage.Name())
}
