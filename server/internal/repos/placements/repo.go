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
	rows := []ads.Placement{}

	// Плейсменты выключенных сайтов не отдаём: запрос по такому плейсменту
	// завершится 200 OK с пустым ответом (SDK вызовет noad).
	err := b.db.Model(ads.Placement{}).
		Select("placements.*").
		Joins("join sites on sites.id = placements.site_id").
		Where("placements.deleted_at is null and placements.active = true").
		Where("sites.deleted_at is null and sites.active = true").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}
