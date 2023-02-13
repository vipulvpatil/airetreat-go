package server

import (
	"context"
	"fmt"

	"github.com/vipulvpatil/airetreat-go/internal/storage"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

type AiRetreatGoService struct {
	pb.UnsafeAiRetreatGoServer
	storage storage.StorageAccessor
}

type ServerDependencies struct {
	Storage storage.StorageAccessor
}

func NewServer(deps ServerDependencies) (*AiRetreatGoService, error) {
	return &AiRetreatGoService{
		storage: deps.Storage,
	}, nil
}

func (s *AiRetreatGoService) Test(ctx context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	test := req.Test
	response := fmt.Sprintf("success: %s", test)
	return &pb.TestResponse{
		Test: response,
	}, nil
}
