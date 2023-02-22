package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	pb "github.com/vipulvpatil/airetreat-go/protos"
	"google.golang.org/grpc/metadata"
)

func contextWithPrefilledRequestingUser() context.Context {
	return metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(
			requestingUserIdCtxKey, "internalUserId1",
			requestingUserEmailCtxKey, "test@example.com",
		),
	)
}

func Test_Test(t *testing.T) {
	t.Run("test runs successfully", func(t *testing.T) {
		server, _ := NewServer(ServerDependencies{})

		response, err := server.Test(
			contextWithPrefilledRequestingUser(),
			&pb.TestRequest{Test: "test_string"},
		)
		assert.NoError(t, err)
		assert.EqualValues(t, response, &pb.TestResponse{Test: "success: test_string"})
	})
}

func Test_GetPlayerId(t *testing.T) {
	tests := []struct {
		name              string
		output            *pb.GetPlayerIdResponse
		playerCreatorMock storage.PlayerCreator
		errorExpected     bool
		errorString       string
	}{
		{
			name:              "test runs successfully",
			output:            &pb.GetPlayerIdResponse{PlayerId: "player_id1"},
			playerCreatorMock: &storage.PlayerCreatorMockSuccess{PlayerId: "player_id1"},
			errorExpected:     false,
			errorString:       "",
		},
		{
			name:              "errors if player creation fails",
			output:            nil,
			playerCreatorMock: &storage.PlayerCreatorMockFailure{},
			errorExpected:     true,
			errorString:       "unable to create player",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithPlayerCreatorMock(
						tt.playerCreatorMock,
					),
				),
			})

			response, err := server.GetPlayerId(
				context.Background(),
				&pb.GetPlayerIdRequest{},
			)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.EqualValues(t, response, tt.output)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}

func Test_CreateGame(t *testing.T) {
	tests := []struct {
		name            string
		output          *pb.CreateGameResponse
		gameCreatorMock storage.GameAccessor
		errorExpected   bool
		errorString     string
	}{
		{
			name:            "test runs successfully",
			output:          &pb.CreateGameResponse{GameId: "game_id1"},
			gameCreatorMock: &storage.GameCreatorMockSuccess{GameId: "game_id1"},
			errorExpected:   false,
			errorString:     "",
		},
		{
			name:            "errors if game creation fails",
			output:          nil,
			gameCreatorMock: &storage.GameCreatorMockFailure{},
			errorExpected:   true,
			errorString:     "unable to create game",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithGameAccessorMock(
						tt.gameCreatorMock,
					),
				),
			})

			response, err := server.CreateGame(
				context.Background(),
				&pb.CreateGameRequest{},
			)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.EqualValues(t, response, tt.output)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}

func Test_JoinGame(t *testing.T) {
	tests := []struct {
		name           string
		input          *pb.JoinGameRequest
		output         *pb.JoinGameResponse
		gameJoinerMock storage.GameAccessor
		errorExpected  bool
		errorString    string
	}{
		{
			name: "test runs successfully",
			input: &pb.JoinGameRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
			},
			output:         &pb.JoinGameResponse{},
			gameJoinerMock: &storage.GameJoinerMockSuccess{},
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "errors if joining game fails",
			input: &pb.JoinGameRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
			},
			output:         nil,
			gameJoinerMock: &storage.GameJoinerMockFailure{},
			errorExpected:  true,
			errorString:    "unable to join game",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithGameAccessorMock(
						tt.gameJoinerMock,
					),
				),
			})

			response, err := server.JoinGame(
				context.Background(),
				tt.input,
			)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.EqualValues(t, response, tt.output)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}

func Test_SendMessage(t *testing.T) {
	tests := []struct {
		name               string
		input              *pb.SendMessageRequest
		output             *pb.SendMessageResponse
		messageCreatorMock storage.MessageCreator
		errorExpected      bool
		errorString        string
	}{
		{
			name: "test runs successfully",
			input: &pb.SendMessageRequest{
				BotId: "bot_id1",
				Text:  "message",
			},
			output:             &pb.SendMessageResponse{},
			messageCreatorMock: &storage.MessageCreatorMockSuccess{},
			errorExpected:      false,
			errorString:        "",
		},
		{
			name: "errors if message creation fails",
			input: &pb.SendMessageRequest{
				BotId: "bot_id1",
				Text:  "message",
			},
			output:             nil,
			messageCreatorMock: &storage.MessageCreatorMockFailure{},
			errorExpected:      true,
			errorString:        "unable to create message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithMessageCreatorMock(
						tt.messageCreatorMock,
					),
				),
			})

			response, err := server.SendMessage(
				context.Background(),
				tt.input,
			)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.EqualValues(t, response, tt.output)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}

func Test_GetGame(t *testing.T) {
	player, _ := model.NewPlayer(
		model.PlayerOptions{
			Id: "player_id1",
		},
	)
	bots := []*model.Bot{}
	for i := 0; i < 5; i++ {
		botOpts := model.BotOptions{
			Id:        fmt.Sprintf("bot_id%d", i+1),
			Name:      fmt.Sprintf("bot%d", i+1),
			TypeOfBot: "AI",
		}
		switch botOpts.Id {
		case "bot_id1":
			botOpts.Messages = []string{
				"Q1: what is your name?",
				"A1: My name is Antony Gonsalvez",
				"Q2: Where is the gold?",
				"A2: what gold!",
			}
		case "bot_id2":
			botOpts.Messages = []string{
				"Q1: What is your name?",
				"A1: Bot 2 Dot 2",
			}
		}
		bot, _ := model.NewBot(botOpts)
		bots = append(bots, bot)
	}
	bots[4].ConnectPlayer(player)
	game, _ := model.NewGame(
		model.GameOptions{
			Id:               "game_id1",
			State:            "STARTED",
			CurrentTurnIndex: 0,
			TurnOrder:        []string{"b", "p1", "b", "p2"},
			StateHandled:     false,
			StateTotalTime:   0,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			Bots:             bots,
		},
	)
	tests := []struct {
		name           string
		input          *pb.GetGameRequest
		output         *pb.GetGameResponse
		gameGetterMock storage.GameAccessor
		errorExpected  bool
		errorString    string
	}{
		{
			name: "test runs successfully",
			input: &pb.GetGameRequest{
				GameId: "game_id1",
			},
			output: &pb.GetGameResponse{},
			gameGetterMock: &storage.GameGetterMockSuccess{
				Game: game,
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "errors if get game fails",
			input: &pb.GetGameRequest{
				GameId: "game_id1",
			},
			output:         nil,
			gameGetterMock: &storage.GameGetterMockFailure{},
			errorExpected:  true,
			errorString:    "unable to get game",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithGameAccessorMock(
						tt.gameGetterMock,
					),
				),
			})

			response, err := server.GetGame(
				context.Background(),
				tt.input,
			)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.EqualValues(t, response, tt.output)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}
