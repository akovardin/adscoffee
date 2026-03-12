package units

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"go.ads.coffee/platform/server/internal/domain/ads"
)

type Placements interface {
	All(ctx context.Context) ([]ads.Unit, error)
}

type Cache struct {
	logger *zap.Logger
	repo   Placements

	lock             sync.RWMutex
	units            []ads.Unit
	unitsByPlacement map[uint][]ads.Unit
}

func NewCache(logger *zap.Logger, repo *Repo) *Cache {
	return &Cache{
		logger:           logger,
		repo:             repo,
		unitsByPlacement: map[uint][]ads.Unit{},
	}
}

func (c *Cache) All(ctx context.Context) []ads.Unit {
	c.lock.RLock()
	defer c.lock.RUnlock()

	units := make([]ads.Unit, len(c.units))
	copy(units, c.units)

	return units
}

func (c *Cache) FindByPlacement(ctx context.Context, placement uint) ([]ads.Unit, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	b, ok := c.unitsByPlacement[placement]

	units := make([]ads.Unit, len(b))
	copy(units, b)

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
	units, err := c.repo.All(context.Background())
	if err != nil {
		c.logger.Error("error on get placements from repo", zap.Error(err))

		return
	}

	c.lock.Lock()
	c.units = units
	c.unitsByPlacement = map[uint][]ads.Unit{}
	for _, unit := range units {
		c.unitsByPlacement[unit.PlacementID] = append(c.unitsByPlacement[unit.PlacementID], unit)
	}
	c.lock.Unlock()
}
