package placement

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

func TestPlacement_Name(t *testing.T) {
	assert.Equal(t, "targetings.placement", New().Name())
}

func TestPlacement_Copy(t *testing.T) {
	p := New()
	copied := p.Copy(map[string]any{})
	assert.IsType(t, &Placement{}, copied)
}

func TestPlacement_Filter(t *testing.T) {
	p := New()

	state := &plugins.State{
		Placement: ads.Placement{ID: 5},
	}

	candidates := []ads.Banner{
		{ID: 1, Targeting: ads.Targeting{Placement: ads.ExcludeInclude{IncludeOr: []string{"5"}}}},
		{ID: 2, Targeting: ads.Targeting{Placement: ads.ExcludeInclude{IncludeOr: []string{"6"}}}},
		{ID: 3, Targeting: ads.Targeting{}},
		{ID: 4, Targeting: ads.Targeting{Placement: ads.ExcludeInclude{ExcludeOr: []string{"5"}}}},
	}

	filtered := p.Filter(candidates, state)

	assert.Len(t, filtered, 2)
	assert.Equal(t, uint(1), filtered[0].ID)
	assert.Equal(t, uint(3), filtered[1].ID)
}

func TestPlacement_FilterWithoutPlacement(t *testing.T) {
	p := New()

	candidates := []ads.Banner{
		{ID: 1, Targeting: ads.Targeting{Placement: ads.ExcludeInclude{IncludeOr: []string{"5"}}}},
	}

	filtered := p.Filter(candidates, &plugins.State{})

	assert.Len(t, filtered, 1)
}
