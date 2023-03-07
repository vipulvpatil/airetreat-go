package server

import (
	"context"

	"github.com/pkg/errors"
	pb "github.com/vipulvpatil/airetreat-go/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AiRetreatGoService) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	gameId, err := s.storage.CreateGame()
	if err != nil {
		return nil, err
	}
	return &pb.CreateGameResponse{GameId: gameId}, nil
}

func (s *AiRetreatGoService) JoinGame(ctx context.Context, req *pb.JoinGameRequest) (*pb.JoinGameResponse, error) {
	err := s.storage.JoinGame(req.GetGameId(), req.GetPlayerId())
	if err != nil {
		return nil, err
	}
	return &pb.JoinGameResponse{}, nil
}

func (s *AiRetreatGoService) GetGameForPlayer(ctx context.Context, req *pb.GetGameForPlayerRequest) (*pb.GetGameForPlayerResponse, error) {
	game, err := s.storage.GetGame(req.GetGameId())
	if err != nil {
		return nil, err
	}

	gameView := game.GameViewForPlayer(req.GetPlayerId())
	if gameView == nil {
		return nil, errors.Errorf("Unable to get game %s for player %s", req.GetGameId(), req.GetPlayerId())
	}

	var stateStartedAt *timestamppb.Timestamp
	if gameView.StateStartedAt != nil {
		stateStartedAt = timestamppb.New(*gameView.StateStartedAt)
	}

	bots := []*pb.Bot{}
	for _, bot := range gameView.Bots {
		messages := []*pb.BotMessage{}
		for _, message := range bot.Messages {
			messages = append(messages, &pb.BotMessage{
				Text: message,
			})
		}
		bots = append(bots, &pb.Bot{
			Id:          bot.Id,
			Name:        bot.Name,
			BotMessages: messages,
		})
	}

	return &pb.GetGameForPlayerResponse{
		State:          gameView.State.String(),
		DisplayMessage: gameView.DisplayMessage,
		StateStartedAt: stateStartedAt,
		StateTotalTime: gameView.StateTotalTime,
		LastQuestion:   gameView.LastQuestion,
		MyBotId:        gameView.MyBotId,
		Bots:           bots,
	}, nil
}

func (s *AiRetreatGoService) GetGamesForPlayer(ctx context.Context, req *pb.GetGamesForPlayerRequest) (*pb.GetGamesForPlayerResponse, error) {
	gameIds, err := s.storage.GetGames(req.GetPlayerId())
	if err != nil {
		return nil, err
	}

	return &pb.GetGamesForPlayerResponse{GameIds: gameIds}, nil
}
