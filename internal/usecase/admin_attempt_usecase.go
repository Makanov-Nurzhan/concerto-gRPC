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
}

func NewAdminAttemptsUseCase(
	txManager tx.Manager,
	sessionRepo domain.SessionRepository,
	attemptsRepo domain.AttemptsRepository,
	exportQueueRepo domain.ExportQueueRepository,
) domain.AdminAttemptsUseCase {
	return &adminAttemptUseCase{
		txManager:       txManager,
		sessionRepo:     sessionRepo,
		attemptsRepo:    attemptsRepo,
		exportQueueRepo: exportQueueRepo,
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
