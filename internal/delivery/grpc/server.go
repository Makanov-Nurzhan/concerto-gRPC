package grpc

import (
	"context"
	"errors"
	adminv1 "github.com/Makanov-Nurzhan/concerto-gRPC/api/gen/adminv1"
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type Server struct {
	adminv1.UnimplementedConcertoAdminServiceServer
	adminUC domain.AdminAttemptsUseCase
}

func NewServer(adminUC domain.AdminAttemptsUseCase) *Server {
	return &Server{adminUC: adminUC}
}

func (s *Server) GetSessionStatus(ctx context.Context, req *adminv1.GetSessionStatusRequest) (*adminv1.GetSessionStatusResponse, error) {
	logger := slog.Default().With("method", "GetSessionStatus")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is null")
	}
	if req.TestTakerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "test_taker_id is required")
	}
	logger = logger.With("test_taker_id", req.TestTakerId)
	st, err := s.adminUC.GetSessionStatus(ctx, req.TestTakerId)
	if err != nil {
		logger.ErrorContext(ctx, "get session status failed", "error", err)
		return nil, err
	}

	resp := &adminv1.GetSessionStatusResponse{
		TestTakerId:       st.TestTakerID,
		HasActiveSession:  st.HasActiveSession,
		CanUpdateAttempts: st.CanUpdateAttempts,
		SessionId:         st.SessionID,
	}
	if st.HasActiveSession && !st.SessionStartTime.IsZero() {
		resp.SessionStartDate = st.SessionStartTime.Format("2006-01-02 15:04:05")
	} else {
		resp.SessionStartDate = ""
	}
	logger.InfoContext(ctx, "get session status success",
		"has_active_session", resp.HasActiveSession,
		"can_update_attempts", resp.CanUpdateAttempts,
		"session_id", resp.SessionId,
		"session_start_date", resp.SessionStartDate,
	)
	return resp, nil
}

func (s *Server) AdminUpdateAttempts(ctx context.Context, req *adminv1.AdminUpdateAttemptsRequest) (*adminv1.AdminUpdateAttemptsResponse, error) {
	logger := slog.Default().With("method", "AdminUpdateAttempts")
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	if req.TestTakerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "test_taker_id is required")
	}
	if req.OperationId == "" {
		return nil, status.Error(codes.InvalidArgument, "operation_id is required")
	}
	if req.ProductLanguage == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid arguments")
	}

	logger = logger.With(
		"operation_id", req.OperationId,
		"test_taker_id", req.TestTakerId,
		"product_variant", req.ProductVariant,
		"product_language", req.ProductLanguage,
		"attempts_to_refund", req.AttemptsToRefund,
		"current_attempts", req.CurrentAttempts,
		"current_used", req.CurrentUsed,
	)

	input := domain.AdminUpdateAttemptsRequest{
		OperationID:      req.OperationId,
		TestTakerID:      req.TestTakerId,
		AttemptsToRefund: req.AttemptsToRefund,
		CurrentAttempts:  req.CurrentAttempts,
		CurrentUsed:      req.CurrentUsed,
		ProductData: domain.ProductData{
			ProductVariant:  req.ProductVariant,
			ProductLanguage: req.ProductLanguage,
		},
	}

	resp, err := s.adminUC.AdminUpdateAttempts(ctx, input)
	if err != nil {
		code, msg := mapDomainError(err)
		logger.WarnContext(ctx, "admin update attempts failed", "code", code, "error", err)
		return &adminv1.AdminUpdateAttemptsResponse{
			Success:       false,
			ErrorCode:     code,
			ErrorMessage:  msg,
			TestTakerId:   req.TestTakerId,
			AttemptsTotal: 0,
			AttemptsUsed:  0,
			Refund:        0,
		}, nil
	}

	logger.InfoContext(ctx, "admin update attempts success",
		"attempts_total", resp.AttemptsTotal,
		"attempts_used", resp.AttemptsUsed,
		"refund", resp.Refund,
	)
	return &adminv1.AdminUpdateAttemptsResponse{
		Success:       resp.Success,
		ErrorCode:     resp.ErrorCode,
		ErrorMessage:  resp.ErrorMessage,
		TestTakerId:   resp.TestTakerID,
		AttemptsTotal: resp.AttemptsTotal,
		AttemptsUsed:  resp.AttemptsUsed,
		Refund:        resp.Refund,
	}, nil
}

func mapDomainError(err error) (string, string) {
	switch {
	case errors.Is(err, domain.ErrOperationAlreadyProcessed):
		return "ALREADY_PROCESSED", err.Error()
	case errors.Is(err, domain.ErrOperationInProgress):
		return "OPERATION_IN_PROGRESS", err.Error()
	case errors.Is(err, domain.ErrInvalidRefund):
		return "INVALID_REFUND", err.Error()
	case errors.Is(err, domain.ErrHasActiveSession):
		return "HAS_ACTIVE_SESSION", err.Error()
	case errors.Is(err, domain.ErrInvalidAttemptCount):
		return "INVALID_ATTEMPT_COUNT", err.Error()
	case errors.Is(err, domain.ErrInvalidTotalAttempts):
		return "INVALID_TOTAL_ATTEMPTS", err.Error()
	case errors.Is(err, domain.ErrFirstDayAttempts):
		return "FIRST_DAY_ATTEMPTS", err.Error()
	case errors.Is(err, domain.ErrFailedToUpdate):
		return "FAILED_TO_UPDATE", err.Error()
	default:
		return "INTERNAL_ERROR", err.Error()
	}
}
