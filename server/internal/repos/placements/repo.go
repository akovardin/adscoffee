package placements

import (
	"context"
	"time"

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
	rows := []ads.Placement{}

	err := b.db.Model(ads.Placement{}).
		Where("deleted_at is not nil and start > ? and end < ? and active = true", time.Now()).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}
