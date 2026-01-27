package placements

import (
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"go.ads.coffee/platform/server/internal/domain/ads"
)

type Repo struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewRepo(logger *zap.Logger, db *gorm.DB) *Repo {
	r := &Repo{
		logger: logger,
		db:     db,
	}

	return r
}

func (b *Repo) All(ctx context.Context) ([]ads.Placement, error) {
	return nil, nil
}
