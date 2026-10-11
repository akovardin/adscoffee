package banners

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"go.ads.coffee/platform/server/internal/domain/ads"
)

type Banners interface {
	All(ctx context.Context) ([]ads.Banner, error)
}

type Cache struct {
	logger *zap.Logger
	repo   Banners

	lock        sync.RWMutex
	banners     []*ads.Banner
	bannersById map[uint]*ads.Banner
}

func NewCache(logger *zap.Logger, repo *Repo) *Cache {
	return &Cache{
		logger:      logger,
		repo:        repo,
		bannersById: map[uint]*ads.Banner{},
	}
}

// All возвращает общий слайс указателей на баннеры. Контракт: результат
// read-only — вызывающая сторона обязана скопировать его (n*8 байт) перед
// любой модификацией. Это позволяет не копировать сами структуры на каждый
// запрос.
func (c *Cache) All(ctx context.Context) []*ads.Banner {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.banners
}

func (c *Cache) One(ctx context.Context, id uint) (*ads.Banner, bool) {
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

	ptrs := make([]*ads.Banner, len(banners))
	byID := make(map[uint]*ads.Banner, len(banners))

	for i := range banners {
		banner := &banners[i]
		ptrs[i] = banner
		byID[banner.ID] = banner
	}

	c.lock.Lock()
	c.banners = ptrs
	c.bannersById = byID
	c.lock.Unlock()
}
