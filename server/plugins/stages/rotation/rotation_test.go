package rotation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.ads.coffee/platform/server/internal/domain/ads"
)

func TestRotattion_rotate(t *testing.T) {
	r := &Rotattion{}

	t.Run("empty candidates", func(t *testing.T) {
		banner, ok, err := r.rotate([]ads.Banner{})
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

		banner, ok, err := r.rotate(candidates)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, candidates[0], banner)
	})

	t.Run("multiple candidates", func(t *testing.T) {
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
			banner, ok, err := r.rotate(candidates)
			assert.NoError(t, err)
			assert.True(t, ok)
			results[banner.ID]++
		}

		assert.Len(t, results, 3)

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
			banner, ok, err := r.rotate(candidates)
			assert.NoError(t, err)
			assert.True(t, ok)
			if banner.ID == "1" {
				zeroCount++
			}
		}

		assert.Less(t, zeroCount, 50)
	})
}
