package usecase

import (
	"context"
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/domain"
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/infra/tx"
	"gorm.io/gorm"
)

type adminAttemptUseCase struct {
	txManager       tx.Manager
	sessionRepo     domain.SessionRepository
	attemptsRepo    domain.AttemptsRepository
	exportQueueRepo domain.ExportQueueRepository
	opRepo          domain.AttemptsOperationRepository
}

func NewAdminAttemptsUseCase(
	txManager tx.Manager,
	sessionRepo domain.SessionRepository,
	attemptsRepo domain.AttemptsRepository,
	exportQueueRepo domain.ExportQueueRepository,
	opRepo domain.AttemptsOperationRepository,
) domain.AdminAttemptsUseCase {
	return &adminAttemptUseCase{
		txManager:       txManager,
		sessionRepo:     sessionRepo,
		attemptsRepo:    attemptsRepo,
		exportQueueRepo: exportQueueRepo,
		opRepo:          opRepo,
	}
}

func (a *adminAttemptUseCase) GetSessionStatus(ctx context.Context, testTakerID uint64) (*domain.SessionStatus, error) {
	session, hasSession, err := a.sessionRepo.HasActiveSession(ctx, nil, testTakerID)
	if err != nil {
		return nil, err
	}

	status := &domain.SessionStatus{
		TestTakerID:       testTakerID,
		HasActiveSession:  hasSession,
		CanUpdateAttempts: !hasSession,
	}
	if hasSession {
		status.SessionID = session.ID
		status.SessionStartTime = session.SessionStartTime
	}
	return status, nil
}

func (a *adminAttemptUseCase) AdminUpdateAttempts(ctx context.Context, req domain.AdminUpdateAttemptsRequest) (*domain.AdminUpdateAttemptsResponse, error) {
	var result *domain.AdminUpdateAttemptsResponse

	err := a.txManager.Do(ctx, func(txDB *gorm.DB) error {
		if req.OperationID != "" {
			op, exists, err := a.opRepo.GetByOperationID(ctx, txDB, req.OperationID)
			if err != nil {
				return err
			}
			if exists {
				switch op.Status {
				case domain.OperationSuccess:

					return domain.ErrOperationAlreadyProcessed
				case domain.OperationPending:
					return domain.ErrOperationInProgress
				}
			}
			newOp := &domain.AttemptsOperation{
				OperationID:      req.OperationID,
				TestTakerID:      req.TestTakerID,
				Variant:          req.ProductData.ProductVariant,
				Lang:             req.ProductData.ProductLanguage,
				AttemptsToRefund: req.AttemptsToRefund,
				Status:           domain.OperationPending,
			}
			if err := a.opRepo.CreatePending(ctx, txDB, newOp); err != nil {
				op, exists, err2 := a.opRepo.GetByOperationID(ctx, txDB, req.OperationID)
				if err2 != nil {
					return err2
				}
				if exists {
					if op.Status == domain.OperationSuccess {
						return domain.ErrOperationAlreadyProcessed
					}
					return domain.ErrOperationInProgress
				}
				return err
			}
		}
		if req.AttemptsToRefund <= 0 {
			return domain.ErrInvalidRefund
		}
		_, hasSession, err := a.sessionRepo.HasActiveSession(ctx, txDB, req.TestTakerID)
		if err != nil {
			return err
		}
		if hasSession {
			return domain.ErrHasActiveSession
		}
		attempts, err := a.attemptsRepo.GetAttemptsForUpdate(ctx, txDB, req.TestTakerID, req.ProductData)
		if err != nil {
			return err
		}
		if attempts.Attempts != req.CurrentAttempts || attempts.Used != req.CurrentUsed {
			return domain.ErrInvalidAttemptCount
		}
		availableAttempts := attempts.Attempts - attempts.Used - attempts.Refund
		if availableAttempts <= 0 {
			return domain.ErrInvalidTotalAttempts
		}

		if req.AttemptsToRefund > availableAttempts {
			return domain.ErrInvalidRefund
		}
		attemptsNumbers := missingAttempts(attempts.Attempts, attempts.Used, req.AttemptsToRefund)
		usedForFirstDay := make([]int32, 0, len(attemptsNumbers))
		for _, att := range attemptsNumbers {
			_, has, err := a.exportQueueRepo.HasFinishFirstDay(ctx, txDB, req.TestTakerID, att, req.ProductData)
			if err != nil {
				return err
			}
			if has {
				usedForFirstDay = append(usedForFirstDay, att)
			}
		}
		if len(usedForFirstDay) > 0 {
			return domain.ErrFirstDayAttempts
		}
		if err := a.attemptsRepo.UpdateAttemptsRefund(ctx, txDB, req.TestTakerID, req.AttemptsToRefund, req.ProductData); err != nil {
			return domain.ErrFailedToUpdate
		}

		if req.OperationID != "" {
			if err := a.opRepo.MarkSuccess(ctx, txDB, req.OperationID); err != nil {
				return err
			}
		}

		result = &domain.AdminUpdateAttemptsResponse{
			Success:       true,
			ErrorCode:     "",
			ErrorMessage:  "Success",
			TestTakerID:   req.TestTakerID,
			AttemptsTotal: attempts.Attempts,
			AttemptsUsed:  attempts.Used,
			Refund:        req.AttemptsToRefund + attempts.Refund,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func missingAttempts(maxAttempts, usedAttempts, refundAttempts int32) []int32 {
	if usedAttempts >= maxAttempts {
		return []int32{}
	}
	available := make([]int32, 0, maxAttempts-usedAttempts)
	for i := usedAttempts + 1; i <= maxAttempts; i++ {
		available = append(available, i)
	}

	if refundAttempts >= int32(len(available)) {
		return available
	}

	start := int32(len(available)) - refundAttempts
	return available[start:]

}
