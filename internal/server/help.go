package server

import (
	"context"
	"errors"

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

	// TODO: Check if bot is allowed to ask for help. limit it to 3.

	currentTurnBot := game.GetBotThatGameIsWaitingOn()

	if sourceBot != currentTurnBot {
		return nil, errors.New("please wait for your turn")
	}

	responseText := "unable to help"

	// TODO: Uncomment once implemented.
	// if game.IsInStateWaitingForAiQuestion() {
	// 	aiBot := aibot.NewAiQuestionGenerator(
	// 		aibot.AiBotOptions{
	// 			BotId:        sourceBot.Id(),
	// 			Game:         game,
	// 			OpenAiClient: nil, // TODO: create/provide appropriate openAIClient
	// 		},
	// 	)
	// 	responseText = aiBot.GetNextQuestion()
	// } else if game.IsInStateWaitingForAiAnswer() {
	// 	aiBot := aibot.NewAiAnswerGenerator(
	// 		aibot.AiBotOptions{
	// 			BotId:        sourceBot.Id(),
	// 			Game:         game,
	// 			OpenAiClient: nil, // TODO: create/provide appropriate openAIClient
	// 		},
	// 	)
	// 	responseText = aiBot.GetNextAnswer()
	// }

	return &pb.HelpResponse{Text: responseText}, err
}
