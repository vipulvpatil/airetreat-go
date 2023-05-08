package server

import (
	"context"
	"math/rand"

	"github.com/pkg/errors"
	pb "github.com/vipulvpatil/airetreat-go/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *AiRetreatGoService) CreateGame(ctx context.Context, req *pb.CreateGameRequest) (*pb.CreateGameResponse, error) {
	gameId, err := s.storage.CreateGame(req.GetPublic())
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}
	return &pb.CreateGameResponse{GameId: gameId}, nil
}

func (s *AiRetreatGoService) JoinGame(ctx context.Context, req *pb.JoinGameRequest) (*pb.JoinGameResponse, error) {
	tx, err := s.storage.BeginTransaction()
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}
	defer tx.Rollback()

	game, err := s.storage.GetGameUsingTransaction(req.GetGameId(), tx)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	if game.HasPlayer(req.GetPlayerId()) {
		return &pb.JoinGameResponse{}, nil
	}

	if !game.HasJustStarted() {
		s.logger.LogError(err)
		return nil, errors.New("cannot join this game")
	}

	aiBot, err := game.GetOneRandomAiBot()
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	err = s.storage.UpdateBotWithPlayerIdUsingTransaction(aiBot.Id(), req.GetPlayerId(), tx)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	err = s.storage.UpdateGameStateIfEnoughPlayersHaveJoinedUsingTransaction(req.GetGameId(), tx)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	err = tx.Commit()
	return &pb.JoinGameResponse{}, err
}

func (s *AiRetreatGoService) GetGameForPlayer(ctx context.Context, req *pb.GetGameForPlayerRequest) (*pb.GetGameForPlayerResponse, error) {
	game, err := s.storage.GetGame(req.GetGameId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	gameView := game.GameViewForPlayer(req.GetPlayerId())
	if gameView == nil {
		return nil, status.Errorf(codes.NotFound, "unable to get game %s for player %s", req.GetGameId(), req.GetPlayerId())
	}

	var stateStartedAt *timestamppb.Timestamp
	if gameView.StateStartedAt != nil {
		stateStartedAt = timestamppb.New(*gameView.StateStartedAt)
	}

	bots := []*pb.Bot{}
	for _, bot := range gameView.Bots {
		bots = append(bots, &pb.Bot{
			Id:   bot.Id,
			Name: bot.Name,
		})
	}

	gameMessages := []*pb.GameMessage{}
	for _, detailedMessage := range gameView.DetailedMessages {
		gameMessages = append(gameMessages, &pb.GameMessage{
			SourceBotId: detailedMessage.SourceBotId,
			TargetBotId: detailedMessage.TargetBotId,
			Text:        detailedMessage.Text,
			Type:        detailedMessage.MessageType,
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
		Messages:       gameMessages,
		WinningBotId:   gameView.WinningBotId,
		MyHelpCount:    gameView.MyHelpCount,
	}, nil
}

func (s *AiRetreatGoService) GetGamesForPlayer(ctx context.Context, req *pb.GetGamesForPlayerRequest) (*pb.GetGamesForPlayerResponse, error) {
	gameIds, err := s.storage.GetGames(req.GetPlayerId())
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	return &pb.GetGamesForPlayerResponse{GameIds: gameIds}, nil
}

func (s *AiRetreatGoService) AutoJoinGame(ctx context.Context, req *pb.AutoJoinGameRequest) (*pb.AutoJoinGameResponse, error) {
	gameIds, err := s.storage.GetAutoJoinableGames()
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	randomlySelectedGameId, err := getRandomGameId(gameIds)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	tx, err := s.storage.BeginTransaction()
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}
	defer tx.Rollback()

	game, err := s.storage.GetGameUsingTransaction(randomlySelectedGameId, tx)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	if game.HasPlayer(req.GetPlayerId()) {
		return &pb.AutoJoinGameResponse{}, nil
	}

	if !game.HasJustStarted() {
		s.logger.LogError(err)
		return nil, errors.New("cannot join this game")
	}

	aiBot, err := game.GetOneRandomAiBot()
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	err = s.storage.UpdateBotWithPlayerIdUsingTransaction(aiBot.Id(), req.GetPlayerId(), tx)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	err = s.storage.UpdateGameStateIfEnoughPlayersHaveJoinedUsingTransaction(randomlySelectedGameId, tx)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}

	err = tx.Commit()
	return &pb.AutoJoinGameResponse{
		GameId: randomlySelectedGameId,
	}, err
}

func getRandomGameId(gameIds []string) (string, error) {
	if len(gameIds) == 0 {
		return "", errors.New("no auto joinable games")
	}

	rand.Shuffle(len(gameIds), func(i, j int) {
		gameIds[i], gameIds[j] = gameIds[j], gameIds[i]
	})
	return gameIds[0], nil
}
