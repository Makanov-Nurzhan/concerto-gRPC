package app

import (
	adminv1 "github.com/Makanov-Nurzhan/concerto-gRPC/api/proto"
	deliverygrpc "github.com/Makanov-Nurzhan/concerto-gRPC/internal/delivery/grpc"
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/domain"
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/infra/tx"
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/repository"
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/usecase"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type App struct {
	TxManager tx.Manager

	SessionRepo     domain.SessionRepository
	AttemptsRepo    domain.AttemptsRepository
	ExportQueueRepo domain.ExportQueueRepository

	adminUC domain.AdminAttemptsUseCase

	GRPCServer *grpc.Server
}

func NewApp(db *gorm.DB) (*App, error) {
	a := &App{}

	txManager := tx.New(db)

	a.SessionRepo = repository.NewSessionRepository(db)
	a.AttemptsRepo = repository.NewAttemptsRepository(db)
	a.ExportQueueRepo = repository.NewExportQueueRepository(db)

	a.adminUC = usecase.NewAdminAttemptsUseCase(txManager, a.SessionRepo, a.AttemptsRepo, a.ExportQueueRepo)

	a.GRPCServer = grpc.NewServer()

	handler := deliverygrpc.NewServer(a.adminUC)

	adminv1.RegisterConcertoAdminServiceServer(a.GRPCServer, handler)

	return a, nil
}
