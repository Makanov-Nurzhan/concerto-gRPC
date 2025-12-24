package domain

import (
	"context"
	"gorm.io/gorm"
)

type OperationStatus string

const (
	OperationPending OperationStatus = "pending"
	OperationSuccess OperationStatus = "success"
)
const (
	OperationRefund string = "refund"
	OperationAdded  string = "added"
)

type AttemptsOperation struct {
	OperationID   string
	TestTakerID   uint64
	Variant       int32
	Lang          string
	Attempts      int32
	OperationType string
	Status        OperationStatus
}

type AttemptsOperationRepository interface {
	GetByOperationID(ctx context.Context, db *gorm.DB, opID string) (*AttemptsOperation, bool, error)
	CreatePending(ctx context.Context, db *gorm.DB, op *AttemptsOperation) error
	MarkSuccess(ctx context.Context, db *gorm.DB, opID string) error
}
