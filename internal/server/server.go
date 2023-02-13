package server

import (
	"context"
	"fmt"

	"github.com/vipulvpatil/airetreat-go/internal/clients/instagram"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	"github.com/vipulvpatil/airetreat-go/internal/workers"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

type AiRetreatGoService struct {
	pb.UnsafeAiRetreatGoServer
	JobStarter      workers.JobStarter
	storage         storage.StorageAccessor
	instagramClient instagram.InstagramClient
}

type ServerDependencies struct {
	JobStarter      workers.JobStarter
	Storage         storage.StorageAccessor
	InstagramClient instagram.InstagramClient
}

func NewServer(deps ServerDependencies) (*AiRetreatGoService, error) {
	return &AiRetreatGoService{
		JobStarter:      deps.JobStarter,
		storage:         deps.Storage,
		instagramClient: deps.InstagramClient,
	}, nil
}

func (s *AiRetreatGoService) Test(ctx context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	test := req.Test
	response := fmt.Sprintf("success: %s", test)
	return &pb.TestResponse{
		Test: response,
	}, nil
}
