package repository

import (
	"context"
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/domain"
	"gorm.io/gorm"
)

type attemptsRepo struct {
	db *gorm.DB
}

func NewAttemptsRepository(db *gorm.DB) domain.AttemptsRepository {
	return &attemptsRepo{
		db: db,
	}
}

func (a *attemptsRepo) GetAttemptsForUpdate(ctx context.Context, db *gorm.DB, testTakerID uint64, data domain.ProductData) (*domain.Attempts, error) {
	if db == nil {
		db = a.db
	}
	var row domain.Attempts

	if err := db.WithContext(ctx).
		Raw(`
			SELECT 
			    user_id AS test_taker_id,
			    attempts,
			    used,
			    refund
			FROM online_ko_variants_access
			WHERE user_id = ?
			  AND variant = ?
			  AND lang = ?
			FOR UPDATE`,
			testTakerID, data.ProductVariant, data.ProductLanguage,
		).Scan(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (a *attemptsRepo) UpdateAttemptsRefund(ctx context.Context, db *gorm.DB, testTakerID uint64, refund int32, data domain.ProductData) error {
	if db == nil {
		db = a.db
	}
	return db.WithContext(ctx).
		Exec(`
			UPDATE online_ko_variants_access
			SET refund = refund + ?
			Where user_id = ?
			AND variant = ?
			AND lang = ?
			`, refund, testTakerID, data.ProductVariant, data.ProductLanguage).Error
}

func (a *attemptsRepo) AddAttempts(ctx context.Context, db *gorm.DB, testTakerID uint64, attempts int32, data domain.ProductData) error {
	if db == nil {
		db = a.db
	}
	return db.WithContext(ctx).
		Exec(`
			UPDATE online_ko_variants_access
			SET attempts = attempts + ?
			Where user_id = ?
			AND variant = ?
			AND lang = ?
			`, attempts, testTakerID, data.ProductVariant, data.ProductLanguage).Error
}
