package server

import (
	"context"
	"errors"

	"github.com/vipulvpatil/airetreat-go/internal/storage"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

func (s *AiRetreatGoService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
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

	gameUpdate, err := game.GetGameUpdateAfterIncomingMessage(sourceBot.Id(), req.GetBotId(), req.GetText())
	if err != nil {
		return nil, err
	}

	newGameState := gameUpdate.State.String()
	updateOptions := storage.GameUpdateOptions{
		State:                   &newGameState,
		CurrentTurnIndex:        gameUpdate.CurrentTurnIndex,
		StateHandled:            gameUpdate.StateHandled,
		LastQuestion:            gameUpdate.LastQuestion,
		LastQuestionTargetBotId: gameUpdate.LastQuestionTargetBotId,
	}

	err = s.storage.UpdateGameStateUsingTransaction(req.GetGameId(), updateOptions, tx)
	if err != nil {
		return nil, err
	}

	err = s.storage.CreateMessageUsingTransaction(sourceBot.Id(), req.GetBotId(), req.GetText(), tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	return &pb.SendMessageResponse{}, err
}
