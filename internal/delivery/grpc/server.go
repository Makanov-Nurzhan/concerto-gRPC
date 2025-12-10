package grpc

import (
	"context"
	adminv1 "github.com/Makanov-Nurzhan/concerto-gRPC/api/proto"
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	adminv1.UnimplementedConcertoAdminServiceServer
	adminUC domain.AdminAttemptsUseCase
}

func NewServer(adminUC domain.AdminAttemptsUseCase) *Server {
	return &Server{adminUC: adminUC}
}

func (s *Server) GetSessionStatus(ctx context.Context, req *adminv1.GetSessionStatusRequest) (*adminv1.GetSessionStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is null")
	}
	if req.TestTakerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "test_taker_id is required")
	}
	st, err := s.adminUC.GetSessionStatus(ctx, req.TestTakerId)
	if err != nil {
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
	return resp, nil
}

func (s *Server) AdminUpdateAttempts(ctx context.Context, req *adminv1.AdminUpdateAttemptsRequest) (*adminv1.AdminUpdateAttemptsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	if req.TestTakerId == 0 {
		return nil, status.Error(codes.InvalidArgument, "test_taker_id is required")
	}

	input := domain.AdminUpdateAttemptsRequest{
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
	switch err {
	case domain.ErrInvalidRefund:
		return "INVALID_REFUND", err.Error()
	case domain.ErrHasActiveSession:
		return "HAS_ACTIVE_SESSION", err.Error()
	case domain.ErrInvalidAttemptCount:
		return "INVALID_ATTEMPT_COUNT", err.Error()
	case domain.ErrInvalidTotalAttempts:
		return "INVALID_TOTAL_ATTEMPTS", err.Error()
	case domain.ErrFirstDayAttempts:
		return "FIRST_DAY_ATTEMPTS", err.Error()
	case domain.ErrFailedToUpdate:
		return "FAILED_TO_UPDATE", err.Error()
	default:
		return "INTERNAL_ERROR", err.Error()
	}
}
