package health

import (
	"context"

	pb "github.com/vipulvpatil/airetreat-go/protos"
)

type AiRetreatGoHealthService struct {
	pb.UnsafeAiRetreatGoHealthServer
}

func (s *AiRetreatGoHealthService) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	return &pb.CheckResponse{}, nil
}
