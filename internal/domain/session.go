package domain

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ProductData struct {
	ProductVariant  int32
	ProductLanguage string
}

type Attempts struct {
	TestTakerID uint64
	Attempts    int32
	Used        int32
	Refund      int32
}

type SessionData struct {
	ID               int32     `gorm:"column:id"`
	SessionStartTime time.Time `gorm:"column:startedTime"`
}

type ExportQueueData struct {
	ID          int32
	TestTakerID uint64
	TestAttempt int32
	CreatedAt   time.Time
}

type SessionStatus struct {
	TestTakerID       uint64
	HasActiveSession  bool
	CanUpdateAttempts bool
	SessionID         int32
	SessionStartTime  time.Time
}

type AdminUpdateAttemptsRequest struct {
	OperationID      string
	TestTakerID      uint64
	CurrentAttempts  int32
	CurrentUsed      int32
	AttemptsToRefund int32
	ProductData      ProductData
}

type AdminAddAttemptsRequest struct {
	OperationID     string
	TestTakerID     uint64
	CurrentAttempts int32
	CurrentUsed     int32
	AttemptsToAdd   int32
	ProductData     ProductData
}

type AdminUpdateAttemptsResponse struct {
	Success       bool
	ErrorCode     string
	ErrorMessage  string
	TestTakerID   uint64
	AttemptsTotal int32
	AttemptsUsed  int32
	Refund        int32
}

type SessionRepository interface {
	HasActiveSession(ctx context.Context, db *gorm.DB, testTakerID uint64) (SessionData, bool, error)
}

type AttemptsRepository interface {
	GetAttemptsForUpdate(ctx context.Context, db *gorm.DB, testTakerID uint64, data ProductData) (*Attempts, error)
	UpdateAttemptsRefund(ctx context.Context, db *gorm.DB, testTakerID uint64, refund int32, data ProductData) error
	AddAttempts(ctx context.Context, db *gorm.DB, testTakerID uint64, attempts int32, data ProductData) error
}

type ExportQueueRepository interface {
	HasFinishFirstDay(ctx context.Context, db *gorm.DB, testTakerID uint64, attempt int32, product ProductData) (ExportQueueData, bool, error)
}

type AdminAttemptsUseCase interface {
	GetSessionStatus(ctx context.Context, testTakerID uint64) (*SessionStatus, error)
	AdminUpdateAttempts(ctx context.Context, req AdminUpdateAttemptsRequest) (*AdminUpdateAttemptsResponse, error)
	AdminAddAttempts(ctx context.Context, req AdminAddAttemptsRequest) (*AdminUpdateAttemptsResponse, error)
}
