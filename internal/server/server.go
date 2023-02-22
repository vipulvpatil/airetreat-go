package server

import (
	"context"
	"fmt"

	"github.com/vipulvpatil/airetreat-go/internal/config"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	pb "github.com/vipulvpatil/airetreat-go/protos"
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
	err := s.storage.JoinGame(req.GameId, req.PlayerId)
	if err != nil {
		return nil, err
	}
	return &pb.JoinGameResponse{}, nil
}
func (s *AiRetreatGoService) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	return nil, nil
}
