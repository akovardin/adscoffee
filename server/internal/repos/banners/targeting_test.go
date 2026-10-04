package banners

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTargetingPlacementParseAndMerge(t *testing.T) {
	// банер (креатив) затаргетирован на плейсмент 2
	banner, err := newTargeting(`{"placement":{"include_or":["2"]}}`)
	require.NoError(t, err)

	// остальные уровни без таргетинга
	empty, err := newTargeting("")
	require.NoError(t, err)

	merged := empty.merge(empty).merge(empty).merge(banner)
	domain := merged.toDomain()

	assert.Equal(t, []string{"2"}, domain.Placement.IncludeOr)
	assert.True(t, domain.Placement.Validate([]string{"2"}))
	assert.False(t, domain.Placement.Validate([]string{"1"}))
}
