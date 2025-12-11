package main

import (
	"context"
	"fmt"
	adminv1 "github.com/Makanov-Nurzhan/concerto-gRPC/api/gen/adminv1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("failed to create gRPC client", err)
	}

	client := adminv1.NewConcertoAdminServiceClient(conn)

	resp, err := client.GetSessionStatus(ctx, &adminv1.GetSessionStatusRequest{
		TestTakerId: 1,
	})
	if err != nil {
		log.Fatal("failed to get session status", err)
	}
	fmt.Printf("Response: %d:%t:%t:%d:%s\n",
		resp.TestTakerId,
		resp.HasActiveSession,
		resp.CanUpdateAttempts,
		resp.SessionId,
		resp.SessionStartDate)

	fmt.Println("▶ Calling AdminUpdateAttempts...")
	opID := uuid.NewString()
	fmt.Println("OperationID:", opID)
	updateReq := &adminv1.AdminUpdateAttemptsRequest{
		OperationId:      opID,
		TestTakerId:      1,
		CurrentAttempts:  4,
		CurrentUsed:      2,
		AttemptsToRefund: 1,
		ProductVariant:   1,
		ProductLanguage:  "ru",
	}

	updateResp, err := client.AdminUpdateAttempts(ctx, updateReq)
	if err != nil {
		log.Fatal("failed to update attempts: ", err)
	}

	if !updateResp.Success {
		fmt.Printf("AdminUpdateAttempts → ERROR [%s]: %s\n",
			updateResp.ErrorCode, updateResp.ErrorMessage)
	} else {
		fmt.Printf("AdminUpdateAttempts → OK: AttemptsTotal=%d, Used=%d, Refund=%d\n",
			updateResp.AttemptsTotal, updateResp.AttemptsUsed, updateResp.Refund)
	}

}
