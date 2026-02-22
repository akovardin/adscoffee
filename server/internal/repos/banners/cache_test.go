package banners

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"

	"go.ads.coffee/platform/server/internal/domain/ads"
)

// MockBanners is a mock implementation of the Banners interface
type MockBanners struct {
	mock.Mock
}

func (m *MockBanners) All(ctx context.Context) ([]ads.Banner, error) {
	args := m.Called(ctx)
	return args.Get(0).([]ads.Banner), args.Error(1)
}

// TestCacheAll tests the All method of Cache
func TestCacheAll(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Create cache instance directly without NewCache
	cache := &Cache{
		logger:      logger,
		banners:     []ads.Banner{},
		bannersById: map[uint]ads.Banner{},
	}

	// Add some test data
	banners := []ads.Banner{
		{ID: 1, Title: "Banner 1"},
		{ID: 2, Title: "Banner 2"},
	}

	cache.lock.Lock()
	cache.banners = banners
	cache.lock.Unlock()

	// Test the All method
	result := cache.All(context.Background())

	// Check that the result is correct
	assert.Len(t, result, 2)
	assert.Equal(t, "Banner 1", result[0].Title)
	assert.Equal(t, "Banner 2", result[1].Title)

}

// TestCacheOne tests the One method of Cache
func TestCacheOne(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Create cache instance directly without NewCache
	cache := &Cache{
		logger:      logger,
		banners:     []ads.Banner{},
		bannersById: map[uint]ads.Banner{},
	}

	// Add some test data
	banner := ads.Banner{ID: 1, Title: "Test Banner"}
	cache.lock.Lock()
	cache.bannersById[1] = banner
	cache.lock.Unlock()

	// Test the One method with existing banner
	result, ok := cache.One(context.Background(), 1)

	// Check that the result is correct
	assert.True(t, ok)
	assert.Equal(t, "Test Banner", result.Title)

	// Test the One method with non-existing banner
	_, ok = cache.One(context.Background(), 2)
	assert.False(t, ok)
}

// TestCacheReload tests the reload method of Cache
func TestCacheReload(t *testing.T) {
	logger := zaptest.NewLogger(t)
	repo := &MockBanners{}

	// Create cache instance directly without NewCache
	cache := &Cache{
		logger:      logger,
		repo:        repo,
		banners:     []ads.Banner{},
		bannersById: map[uint]ads.Banner{},
	}

	// Set up mock expectations
	banners := []ads.Banner{
		{ID: 1, Title: "Banner 1"},
		{ID: 2, Title: "Banner 2"},
	}
	repo.On("All", mock.Anything).Return(banners, nil)

	// Call the reload method
	cache.reload()

	// Check that the cache was updated correctly
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	assert.Len(t, cache.banners, 2)
	assert.Equal(t, "Banner 1", cache.banners[0].Title)
	assert.Equal(t, "Banner 2", cache.banners[1].Title)

	assert.Len(t, cache.bannersById, 2)
	assert.Equal(t, "Banner 1", cache.bannersById[1].Title)
	assert.Equal(t, "Banner 2", cache.bannersById[2].Title)

	// Verify mock expectations
	repo.AssertExpectations(t)
}

// TestCacheStart tests the Start method of Cache
func TestCacheStart(t *testing.T) {
	logger := zaptest.NewLogger(t)
	repo := &MockBanners{}

	// Create cache instance directly without NewCache
	cache := &Cache{
		logger:      logger,
		repo:        repo,
		banners:     []ads.Banner{},
		bannersById: map[uint]ads.Banner{},
	}

	// Set up mock expectations
	banners := []ads.Banner{
		{ID: 1, Title: "Banner 1"},
	}
	repo.On("All", mock.Anything).Return(banners, nil)

	// Create a context with timeout for testing
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Start the cache in a separate goroutine
	go cache.Start(ctx)

	// Give some time for the reload to happen
	time.Sleep(50 * time.Millisecond)

	// Check that the cache was updated correctly
	cache.lock.RLock()
	defer cache.lock.RUnlock()

	assert.Len(t, cache.banners, 1)
	assert.Equal(t, "Banner 1", cache.banners[0].Title)

	// Verify mock expectations
	repo.AssertExpectations(t)
}
