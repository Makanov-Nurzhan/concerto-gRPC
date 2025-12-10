package main

import (
	DI "github.com/Makanov-Nurzhan/concerto-gRPC/internal/app"
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/config"
	"github.com/Makanov-Nurzhan/concerto-gRPC/internal/infra/db"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	cfg := config.Load()
	dbConn := db.InitDB(cfg)

	app, err := DI.NewApp(dbConn)
	if err != nil {
		log.Fatal("failed to initialize application", err)
	}

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatal("failed to listen for connections", err)
	}

	reflection.Register(app.GRPCServer)

	log.Printf("gRPC server listening on port %s", cfg.GRPCPort)
	if err := app.GRPCServer.Serve(lis); err != nil {
		log.Fatal("failed to serve gRPC", err)
	}
}
