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
		return nil, err
	}

	sourceBot := game.BotWithPlayerId(req.GetPlayerId())
	if sourceBot == nil {
		return nil, errors.New("incorrect game")
	}

	if !sourceBot.CanGetHelp() {
		return nil, errors.New("no more help possible")
	}

	currentTurnBot := game.GetBotThatGameIsWaitingOn()

	if sourceBot != currentTurnBot {
		return nil, errors.New("please wait for your turn")
	}

	responseText := "unable to help"

	tx, err := s.storage.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = s.storage.UpdateBotDecrementHelpCountUsingTransaction(sourceBot.Id(), tx)
	if err != nil {
		return nil, err
	}

	if game.IsInStateWaitingForAiQuestion() {
		aiBot := aibot.NewAiQuestionGenerator(
			aibot.AiBotOptions{
				BotId:        sourceBot.Id(),
				Game:         game,
				OpenAiClient: s.openAiClient,
			},
		)
		responseText = aiBot.GetNextQuestion()
	} else if game.IsInStateWaitingForAiAnswer() {
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
