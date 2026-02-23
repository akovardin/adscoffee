package placements

import (
	"context"
	"sync"
	"time"

	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.uber.org/zap"
)

type Placements interface {
	All(ctx context.Context) ([]ads.Placement, error)
}

type Cache struct {
	logger *zap.Logger
	repo   Placements

	lock          sync.RWMutex
	placements    []ads.Placement
	placementById map[uint]ads.Placement
}

func NewCache(logger *zap.Logger, repo *Repo) *Cache {
	return &Cache{
		logger:        logger,
		repo:          repo,
		placementById: map[uint]ads.Placement{},
	}
}

func (c *Cache) All(ctx context.Context) []ads.Placement {
	c.lock.RLock()
	defer c.lock.RUnlock()

	placements := make([]ads.Placement, len(c.placements))
	copy(placements, c.placements)

	return placements
}

func (c *Cache) One(ctx context.Context, id uint) (ads.Placement, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	b, ok := c.placementById[id]

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
	placements, err := c.repo.All(context.Background())
	if err != nil {
		c.logger.Error("error on get placements from repo", zap.Error(err))

		return
	}

	c.lock.Lock()
	c.placements = placements
	for _, banner := range placements {
		c.placementById[banner.ID] = banner
	}
	c.lock.Unlock()
}
