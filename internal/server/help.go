package server

import (
	"context"
	"errors"

	aibot "github.com/vipulvpatil/airetreat-go/internal/services/ai-bot"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

func (s *AiRetreatGoService) Help(ctx context.Context, req *pb.HelpRequest) (*pb.HelpResponse, error) {
	game, err := s.storage.GetGame(req.GetGameId())
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

	if !sourceBot.CanGetHelp() {
		err := errors.New("no more help possible")
		s.logger.LogError(err)
		return nil, err
	}

	currentTurnBot := game.GetBotThatGameIsWaitingOn()

	if sourceBot != currentTurnBot {
		err := errors.New("please wait for your turn")
		s.logger.LogError(err)
		return nil, err
	}

	responseText := "unable to help"

	tx, err := s.storage.BeginTransaction()
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}
	defer tx.Rollback()

	err = s.storage.UpdateBotDecrementHelpCountUsingTransaction(sourceBot.Id(), tx)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	if game.IsInStateWaitingForHumanQuestion() {
		aiBot := aibot.NewAiQuestionGenerator(
			aibot.AiBotOptions{
				BotId:        sourceBot.Id(),
				Game:         game,
				OpenAiClient: s.openAiClient,
			},
		)
		responseText = aiBot.GetNextQuestion()
	} else if game.IsInStateWaitingForHumanAnswer() {
		aiBot := aibot.NewAiAnswerGenerator(
			aibot.AiBotOptions{
				BotId:        sourceBot.Id(),
				Game:         game,
				OpenAiClient: s.openAiClient,
			},
		)
		responseText = aiBot.GetNextAnswer()
	}

	err = tx.Commit()
	return &pb.HelpResponse{Text: responseText}, err
}
