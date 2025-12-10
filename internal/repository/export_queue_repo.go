package repository

import (
	"context"
	"errors"
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/domain"
	"gorm.io/gorm"
)

type exportQueueRepo struct {
	db *gorm.DB
}

func NewExportQueueRepository(db *gorm.DB) domain.ExportQueueRepository {
	return &exportQueueRepo{
		db: db,
	}
}

func (e exportQueueRepo) HasFinishFirstDay(ctx context.Context, db *gorm.DB, testTakerID uint64, attempt int32, product domain.ProductData) (domain.ExportQueueData, bool, error) {
	if db == nil {
		db = e.db
	}
	var data domain.ExportQueueData

	err := db.WithContext(ctx).
		Table("online_ko_export_queue").
		Select("id, user_id as test_taker_id, test_attempt, created_at").
		Where("user_id = ? AND test_attempt = ? AND variant = ? AND lang = ?",
			testTakerID, attempt, product.ProductVariant, product.ProductLanguage).
		First(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ExportQueueData{}, false, nil
		}
		return domain.ExportQueueData{}, false, err
	}

	return data, true, nil
}
