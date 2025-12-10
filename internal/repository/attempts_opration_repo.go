package repository

import (
	"context"
	"errors"
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/domain"
	"gorm.io/gorm"
)

type attemptsOperationRepo struct {
	db *gorm.DB
}

func NewAttemptsOperationRepository(db *gorm.DB) domain.AttemptsOperationRepository {
	return &attemptsOperationRepo{db: db}
}
func (a attemptsOperationRepo) GetByOperationID(ctx context.Context, db *gorm.DB, opID string) (*domain.AttemptsOperation, bool, error) {
	if db == nil {
		db = a.db
	}

	var row struct {
		OperationID    string
		TestTakerID    uint64 `gorm:"column:test_taker_id"`
		Variant        int32
		Lang           string
		AttemptsRefund int32 `gorm:"column:attempts_refund"`
		Status         string
	}

	err := db.WithContext(ctx).
		Table("online_ko_admin_attempts_operations").
		Where("operation_id = ?", opID).
		First(&row).Error
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, false, nil
		}
		return nil, false, err
	}

	op := &domain.AttemptsOperation{
		OperationID:      row.OperationID,
		TestTakerID:      row.TestTakerID,
		Variant:          row.Variant,
		Lang:             row.Lang,
		AttemptsToRefund: row.AttemptsRefund,
		Status:           domain.OperationStatus(row.Status),
	}

	return op, true, nil
}

func (a attemptsOperationRepo) CreatePending(ctx context.Context, db *gorm.DB, op *domain.AttemptsOperation) error {
	if db == nil {
		db = a.db
	}

	return db.WithContext(ctx).
		Exec(`
            INSERT INTO online_ko_admin_attempts_operations
                (operation_id, test_taker_id, variant, lang, attempts_refund, status)
            VALUES (?, ?, ?, ?, ?, 'pending')
        `,
			op.OperationID,
			op.TestTakerID,
			op.Variant,
			op.Lang,
			op.AttemptsToRefund,
		).Error
}

func (a attemptsOperationRepo) MarkSuccess(ctx context.Context, db *gorm.DB, opID string) error {
	if db == nil {
		db = a.db
	}

	return db.WithContext(ctx).
		Exec(`
            UPDATE online_ko_admin_attempts_operations
            SET status = 'success'
            WHERE operation_id = ?
        `, opID).Error
}
