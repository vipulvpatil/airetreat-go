package server

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/config"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	pb "github.com/vipulvpatil/airetreat-go/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AiRetreatGoService struct {
	pb.UnsafeAiRetreatGoServer
	storage storage.StorageAccessor
	config  *config.Config
}

type ServerDependencies struct {
	Storage storage.StorageAccessor
	Config  *config.Config
}

func NewServer(deps ServerDependencies) (*AiRetreatGoService, error) {
	return &AiRetreatGoService{
		storage: deps.Storage,
		config:  deps.Config,
	}, nil
}

func (s *AiRetreatGoService) Test(ctx context.Context, req *pb.TestRequest) (*pb.TestResponse, error) {
	test := req.Test
	response := fmt.Sprintf("success: %s", test)
	return &pb.TestResponse{
		Test: response,
	}, nil
}

func (s *AiRetreatGoService) GetPlayerId(ctx context.Context, req *pb.GetPlayerIdRequest) (*pb.GetPlayerIdResponse, error) {
	playerId, err := s.storage.CreatePlayer()
	if err != nil {
		return nil, err
	}
	return &pb.GetPlayerIdResponse{PlayerId: playerId}, nil
}

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

func (s *AiRetreatGoService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	err := s.storage.CreateMessage(req.GetBotId(), req.GetText())
	if err != nil {
		return nil, err
	}
	return &pb.SendMessageResponse{}, nil
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
