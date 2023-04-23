package server

import (
	"context"

	"github.com/pkg/errors"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

func (s *AiRetreatGoService) GetPlayerId(ctx context.Context, req *pb.GetPlayerIdRequest) (*pb.GetPlayerIdResponse, error) {
	playerId, err := s.storage.CreatePlayer(nil)
	if err != nil {
		return nil, err
	}
	return &pb.GetPlayerIdResponse{PlayerId: playerId}, nil
}

func (s *AiRetreatGoService) RegisterPlayerId(ctx context.Context, req *pb.RegisterPlayerIdRequest) (*pb.RegisterPlayerIdResponse, error) {
	user, err := getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	userId := user.GetId()
	playerId := req.GetPlayerId()

	tx, err := s.storage.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	player, err := s.storage.GetPlayerUsingTransaction(playerId, tx)
	if err != nil {
		return nil, err
	}

	if player.UserId() != nil {
		return nil, errors.New("player is already registered")
	}

	err = s.storage.UpdatePlayerWithUserIdUsingTransaction(playerId, userId, tx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	return &pb.RegisterPlayerIdResponse{ConfirmedPlayerId: playerId}, err
}
