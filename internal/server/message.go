package server

import (
	"context"

	pb "github.com/vipulvpatil/airetreat-go/protos"
)

func (s *AiRetreatGoService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	err := s.storage.CreateMessage(req.GetBotId(), req.GetText())
	if err != nil {
		return nil, err
	}
	return &pb.SendMessageResponse{}, nil
}
