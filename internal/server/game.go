package server

import (
	"context"

	"github.com/pkg/errors"
	pb "github.com/vipulvpatil/airetreat-go/protos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	tx, err := s.storage.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	game, err := s.storage.GetGameUsingTransaction(req.GetGameId(), tx)
	if err != nil {
		return nil, err
	}

	if game.HasPlayer(req.GetPlayerId()) {
		return &pb.JoinGameResponse{}, nil
	}

	if !game.HasJustStarted() {
		return nil, errors.New("cannot join this game")
	}

	aiBot, err := game.GetOneRandomAiBot()
	if err != nil {
		return nil, err
	}

	err = s.storage.UpdateBotWithPlayerIdUsingTransaction(aiBot.Id(), req.GetPlayerId(), tx)
	if err != nil {
		return nil, err
	}

	err = s.storage.UpdateGameStateIfEnoughPlayersHaveJoinedUsingTransaction(req.GetGameId(), tx)
	if err != nil {
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
		// messages := []*pb.BotMessage{}
		// for _, message := range bot.Messages {
		// 	messages = append(messages, &pb.BotMessage{
		// 		Text: message,
		// 	})
		// }
		bots = append(bots, &pb.Bot{
			Id:   bot.Id,
			Name: bot.Name,
		})
	}

	conversation := []*pb.ConversationElement{}
	for _, conversationElement := range gameView.Conversation {
		conversation = append(conversation, &pb.ConversationElement{
			IsQuestion: conversationElement.IsQuestion,
			BotId:      conversationElement.BotId,
			Text:       conversationElement.Text,
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
		Conversation:   conversation,
	}, nil
}

func (s *AiRetreatGoService) GetGamesForPlayer(ctx context.Context, req *pb.GetGamesForPlayerRequest) (*pb.GetGamesForPlayerResponse, error) {
	gameIds, err := s.storage.GetGames(req.GetPlayerId())
	if err != nil {
		return nil, err
	}

	return &pb.GetGamesForPlayerResponse{GameIds: gameIds}, nil
}
