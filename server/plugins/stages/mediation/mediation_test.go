package mediation

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

func TestNew(t *testing.T) {
	mediation := New()
	assert.NotNil(t, mediation)
}

func TestMediation_Name(t *testing.T) {
	mediation := New()
	name := mediation.Name()
	assert.Equal(t, "stages.mediation", name)
}

func TestMediation_Copy(t *testing.T) {
	mediation := New()
	cfgMap := map[string]any{"key": "value"}
	copied := mediation.Copy(cfgMap)
	assert.NotNil(t, copied)
	assert.IsType(t, &Mediation{}, copied)
}

func TestMediation_Do(t *testing.T) {
	mediation := New()

	t.Run("single unit", func(t *testing.T) {
		ctx := context.Background()
		state := &plugins.State{
			Request: &http.Request{},
			Placement: &plugins.Placement{
				ID: "test-placement",
				Units: []ads.Unit{
					{
						ID:      "1",
						Title:   "Test Unit",
						Price:   100,
						Network: "test-network",
					},
				},
			},
			Winners: []ads.Banner{},
		}

		err := mediation.Do(ctx, state)
		assert.NoError(t, err)
		assert.Len(t, state.Winners, 1)
		assert.Equal(t, "1", state.Winners[0].ID)
		assert.Equal(t, "Test Unit", state.Winners[0].Title)
		assert.Equal(t, 100, state.Winners[0].Price)
		assert.Equal(t, ads.CreativeTypeMediator, state.Winners[0].Type)
		assert.Equal(t, "test-network", state.Winners[0].Network)
	})

	t.Run("multiple units", func(t *testing.T) {
		ctx := context.Background()
		state := &plugins.State{
			Request: &http.Request{},
			Placement: &plugins.Placement{
				ID: "test-placement",
				Units: []ads.Unit{
					{
						ID:      "1",
						Title:   "Unit 1",
						Price:   100,
						Network: "network-1",
					},
					{
						ID:      "2",
						Title:   "Unit 2",
						Price:   200,
						Network: "network-2",
					},
				},
			},
			Winners: []ads.Banner{},
		}

		err := mediation.Do(ctx, state)
		assert.NoError(t, err)
		assert.Len(t, state.Winners, 1)
		assert.Contains(t, []string{"1", "2"}, state.Winners[0].ID)
	})

	t.Run("no units", func(t *testing.T) {
		ctx := context.Background()
		state := &plugins.State{
			Request:   &http.Request{},
			Placement: &plugins.Placement{ID: "test-placement", Units: []ads.Unit{}},
			Winners:   []ads.Banner{},
		}

		err := mediation.Do(ctx, state)
		assert.NoError(t, err)
		// When there are no units and no existing winners,
		// the rotate function returns ok=false, but the Do method
		// still sets state.Winners to a slice with an empty banner
		assert.Len(t, state.Winners, 1)
		assert.Equal(t, "", state.Winners[0].ID)
	})

	t.Run("existing winners", func(t *testing.T) {
		ctx := context.Background()
		state := &plugins.State{
			Request: &http.Request{},
			Placement: &plugins.Placement{
				ID: "test-placement",
				Units: []ads.Unit{
					{
						ID:      "1",
						Title:   "Test Unit",
						Price:   100,
						Network: "test-network",
					},
				},
			},
			Winners: []ads.Banner{
				{
					ID:    "existing",
					Title: "Existing Banner",
					Price: 50,
				},
			},
		}

		err := mediation.Do(ctx, state)
		assert.NoError(t, err)
		assert.Len(t, state.Winners, 1)
		// Should select from both existing winners and placement units
		assert.Contains(t, []string{"existing", "1"}, state.Winners[0].ID)
	})

	t.Run("no existing winners and no units", func(t *testing.T) {
		ctx := context.Background()
		state := &plugins.State{
			Request:   &http.Request{},
			Placement: &plugins.Placement{ID: "test-placement", Units: []ads.Unit{}},
			Winners:   []ads.Banner{},
		}

		err := mediation.Do(ctx, state)
		assert.NoError(t, err)
		// When there are no units and no existing winners,
		// the rotate function returns ok=false, but the Do method
		// still sets state.Winners to a slice with an empty banner
		assert.Len(t, state.Winners, 1)
		assert.Equal(t, "", state.Winners[0].ID)
	})
}

func TestMediation_rotate(t *testing.T) {
	m := &Mediation{}

	t.Run("empty candidates", func(t *testing.T) {
		banner, ok, err := m.rotate([]ads.Banner{})
		assert.NoError(t, err)
		assert.False(t, ok)
		assert.Equal(t, ads.Banner{}, banner)
	})

	t.Run("single candidate", func(t *testing.T) {
		candidates := []ads.Banner{
			{
				ID:    "1",
				Title: "Test Banner",
				Price: 100,
			},
		}

		banner, ok, err := m.rotate(candidates)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, candidates[0], banner)
	})

	t.Run("multiple candidates with different prices", func(t *testing.T) {
		candidates := []ads.Banner{
			{
				ID:    "1",
				Title: "Banner 1",
				Price: 100,
			},
			{
				ID:    "2",
				Title: "Banner 2",
				Price: 200,
			},
			{
				ID:    "3",
				Title: "Banner 3",
				Price: 300,
			},
		}

		results := make(map[string]int)
		for i := 0; i < 1000; i++ {
			banner, ok, err := m.rotate(candidates)
			assert.NoError(t, err)
			assert.True(t, ok)
			results[banner.ID]++
		}

		assert.Len(t, results, 3)

		// Higher price banners should be selected more often
		assert.Greater(t, results["3"], results["2"])
		assert.Greater(t, results["2"], results["1"])
	})

	t.Run("candidates with zero price", func(t *testing.T) {
		candidates := []ads.Banner{
			{
				ID:    "1",
				Title: "Banner 1",
				Price: 0,
			},
			{
				ID:    "2",
				Title: "Banner 2",
				Price: 100,
			},
		}

		zeroCount := 0
		for i := 0; i < 1000; i++ {
			banner, ok, err := m.rotate(candidates)
			assert.NoError(t, err)
			assert.True(t, ok)
			if banner.ID == "1" {
				zeroCount++
			}
		}

		// Banner with zero price should be selected rarely
		assert.Less(t, zeroCount, 50)
	})

	t.Run("all candidates with zero price", func(t *testing.T) {
		candidates := []ads.Banner{
			{
				ID:    "1",
				Title: "Banner 1",
				Price: 0,
			},
			{
				ID:    "2",
				Title: "Banner 2",
				Price: 0,
			},
		}

		// When all weights are zero, the chooser will return an error
		_, ok, err := m.rotate(candidates)
		assert.Error(t, err)
		assert.False(t, ok)
	})

	t.Run("candidates with equal prices", func(t *testing.T) {
		candidates := []ads.Banner{
			{
				ID:    "1",
				Title: "Banner 1",
				Price: 100,
			},
			{
				ID:    "2",
				Title: "Banner 2",
				Price: 100,
			},
			{
				ID:    "3",
				Title: "Banner 3",
				Price: 100,
			},
		}

		results := make(map[string]int)
		for i := 0; i < 1000; i++ {
			banner, ok, err := m.rotate(candidates)
			assert.NoError(t, err)
			assert.True(t, ok)
			results[banner.ID]++
		}

		// All banners should be selected approximately equally
		assert.InDelta(t, results["1"], results["2"], 100)
		assert.InDelta(t, results["2"], results["3"], 100)
		assert.InDelta(t, results["1"], results["3"], 100)
	})
}
