package repository

import (
	"context"
	"errors"
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/domain"
	"gorm.io/gorm"
)

type sessionRepo struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) domain.SessionRepository {
	return &sessionRepo{db: db}
}

func (s sessionRepo) HasActiveSession(ctx context.Context, db *gorm.DB, testTakerID uint64) (domain.SessionData, bool, error) {
	if db == nil {
		db = s.db
	}

	var sessionData domain.SessionData
	err := db.WithContext(ctx).
		Raw(`SELECT id, startedTime FROM online_ko_sessions WHERE user_id = ? LIMIT 1`, testTakerID).
		Scan(&sessionData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.SessionData{}, false, nil
		}
		return domain.SessionData{}, false, err
	}

	if sessionData.ID == 0 {
		return domain.SessionData{}, false, nil
	}

	return sessionData, true, nil
}
