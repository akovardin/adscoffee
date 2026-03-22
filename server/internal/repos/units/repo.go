package units

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

func (r *Repo) All(ctx context.Context) ([]ads.Unit, error) {
	rows := []ads.Unit{}

	err := r.db.Model(ads.Unit{}).
		Preload("Network").
		Where("deleted_at is null and active = true and archived_at is NULL").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}
