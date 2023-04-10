package server

import (
	"context"
	"errors"

	"github.com/vipulvpatil/airetreat-go/internal/storage"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

func (s *AiRetreatGoService) Tag(ctx context.Context, req *pb.TagRequest) (*pb.TagResponse, error) {
	tx, err := s.storage.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	game, err := s.storage.GetGameUsingTransaction(req.GetGameId(), tx)
	if err != nil {
		return nil, err
	}

	sourceBot := game.BotWithPlayerId(req.GetPlayerId())
	if sourceBot == nil {
		return nil, errors.New("incorrect game")
	}

	gameUpdate, err := game.GetGameUpdateAfterTag(sourceBot.Id(), req.GetBotId())
	if err != nil {
		return nil, err
	}

	newGameState := gameUpdate.State.String()
	updateOptions := storage.GameUpdateOptions{
		State: &newGameState,
	}

	err = s.storage.UpdateGameStateUsingTransaction(req.GetGameId(), updateOptions, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	return &pb.TagResponse{}, err
}
