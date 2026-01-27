package banners

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"go.ads.coffee/platform/server/internal/domain/ads"
)

type Cache struct {
	logger *zap.Logger
	repo   *Repo

	lock        sync.RWMutex
	banners     []ads.Banner
	bannersById map[string]ads.Banner
}

func NewCache(logger *zap.Logger, repo *Repo) *Cache {
	return &Cache{
		logger:      logger,
		repo:        repo,
		bannersById: map[string]ads.Banner{},
	}
}

func (c *Cache) All(ctx context.Context) []ads.Banner {
	c.lock.RLock()
	defer c.lock.RUnlock()

	banners := make([]ads.Banner, len(c.banners))
	copy(banners, c.banners)

	return banners
}

func (c *Cache) One(ctx context.Context, id string) (ads.Banner, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	b, ok := c.bannersById[id]

	return b, ok
}

// Start reload banners cache.
func (c *Cache) Start(ctx context.Context) {
	c.reload()

	ticker := time.NewTicker(time.Minute * 1)

	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		default:
			c.reload()
		}
	}
}

func (c *Cache) reload() {
	banners, err := c.repo.All(context.Background())
	if err != nil {
		c.logger.Error("error on get banners from repo", zap.Error(err))

		return
	}

	c.lock.Lock()
	c.banners = banners
	for _, banner := range banners {
		c.bannersById[banner.ID] = banner
	}
	c.lock.Unlock()
}
