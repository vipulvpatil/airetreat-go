package server

import (
	"context"

	pb "github.com/vipulvpatil/airetreat-go/protos"
)

func (s *AiRetreatGoService) GetPlayerId(ctx context.Context, req *pb.GetPlayerIdRequest) (*pb.GetPlayerIdResponse, error) {
	playerId, err := s.storage.CreatePlayer()
	if err != nil {
		return nil, err
	}
	return &pb.GetPlayerIdResponse{PlayerId: playerId}, nil
}
