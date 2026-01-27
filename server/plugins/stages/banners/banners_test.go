//nolint:errcheck
package banners

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/domain/plugins"
)

type mockBannersCache struct {
	banners []ads.Banner
}

func (m *mockBannersCache) All(ctx context.Context) []ads.Banner {
	return m.banners
}

func TestNew(t *testing.T) {
	mockCache := &mockBannersCache{}

	banners := &Banners{
		cache: mockCache,
	}

	assert.NotNil(t, banners)
	assert.Equal(t, mockCache, banners.cache)
}

func TestBanners_Name(t *testing.T) {
	mockCache := &mockBannersCache{}
	banners := &Banners{
		cache: mockCache,
	}

	name := banners.Name()

	assert.Equal(t, "stages.banners", name)
}

func TestBanners_Copy(t *testing.T) {
	mockCache := &mockBannersCache{}
	banners := &Banners{
		cache: mockCache,
	}

	cfgMap := map[string]any{"key": "value"}
	copied := banners.Copy(cfgMap)

	assert.NotNil(t, copied)
	assert.IsType(t, &Banners{}, copied)
}

func TestBanners_Do(t *testing.T) {
	testBanners := []ads.Banner{
		{
			ID:     "1",
			Title:  "Test Banner 1",
			Price:  100,
			Active: true,
		},
		{
			ID:     "2",
			Title:  "Test Banner 2",
			Price:  200,
			Active: false,
		},
	}

	mockCache := &mockBannersCache{
		banners: testBanners,
	}
	banners := &Banners{
		cache: mockCache,
	}

	ctx := context.Background()
	state := &plugins.State{
		Request:    &http.Request{},
		Response:   nil,
		User:       nil,
		Device:     nil,
		Candidates: []ads.Banner{},
		Winners:    []ads.Banner{},
	}

	banners.Do(ctx, state)

	assert.Equal(t, testBanners, state.Candidates)
	assert.Len(t, state.Candidates, 2)
}

func TestBanners_DoWithEmptyCache(t *testing.T) {
	mockCache := &mockBannersCache{
		banners: []ads.Banner{},
	}
	banners := &Banners{
		cache: mockCache,
	}

	ctx := context.Background()
	state := &plugins.State{
		Request:    &http.Request{},
		Response:   nil,
		User:       nil,
		Device:     nil,
		Candidates: []ads.Banner{},
		Winners:    []ads.Banner{},
	}

	banners.Do(ctx, state)

	assert.Empty(t, state.Candidates)
	assert.Len(t, state.Candidates, 0)
}
