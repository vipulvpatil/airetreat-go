package server

import (
	"context"
	"errors"
	"strings"

	"github.com/vipulvpatil/airetreat-go/internal/storage"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

func (s *AiRetreatGoService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	tx, err := s.storage.BeginTransaction()
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}
	defer tx.Rollback()

	messageText := req.GetText()

	err = validateMessageText(messageText)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	game, err := s.storage.GetGameUsingTransaction(req.GetGameId(), tx)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	sourceBot := game.BotWithPlayerId(req.GetPlayerId())
	if sourceBot == nil {
		err := errors.New("incorrect game")
		s.logger.LogError(err)
		return nil, err
	}

	gameUpdate, err := game.GetGameUpdateAfterIncomingMessage(sourceBot.Id(), req.GetBotId(), messageText)
	if err != nil {
		s.logger.LogError(err)
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
		s.logger.LogError(err)
		return nil, err
	}

	err = s.storage.CreateMessageUsingTransaction(sourceBot.Id(), req.GetBotId(), req.GetText(), req.GetType(), tx)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	err = tx.Commit()
	return &pb.SendMessageResponse{}, err
}

// TODO: Find a more appropriate place for this function
func validateMessageText(text string) error {
	if utilities.IsBlank(text) {
		return errors.New("message cannot be blank")
	}

	if len(strings.Trim(text, " ")) > 120 {
		return errors.New("message cannot be this long")
	}

	return nil
}
