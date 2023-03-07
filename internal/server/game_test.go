package server

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	pb "github.com/vipulvpatil/airetreat-go/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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
				assert.EqualValues(t, tt.output, response)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}

func Test_JoinGame(t *testing.T) {
	tests := []struct {
		name             string
		input            *pb.JoinGameRequest
		output           *pb.JoinGameResponse
		txProviderMock   storage.DatabaseTransactionProvider
		gameAccessorMock storage.GameAccessor
		botAccessorMock  storage.BotAccessor
		txShouldCommit   bool
		errorExpected    bool
		errorString      string
	}{
		{
			name:             "errors if unable to get transaction",
			input:            &pb.JoinGameRequest{},
			output:           nil,
			txProviderMock:   &storage.DatabaseTransactionProviderMockFailure{},
			txShouldCommit:   false,
			gameAccessorMock: nil,
			botAccessorMock:  nil,
			errorExpected:    true,
			errorString:      "unable to begin a db transaction",
		},
		{
			name: "errors if unable to get game",
			input: &pb.JoinGameRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
			},
			output:         nil,
			txProviderMock: &storage.DatabaseTransactionProviderMockSuccess{},
			txShouldCommit: false,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) (*model.Game, error) {
					return nil, errors.New("unable to get game")
				},
			},
			botAccessorMock: nil,
			errorExpected:   true,
			errorString:     "unable to get game",
		},
		{
			name: "errors if game is not waiting for more people to join",
			input: &pb.JoinGameRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
			},
			output:         nil,
			txProviderMock: &storage.DatabaseTransactionProviderMockSuccess{},
			txShouldCommit: false,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) (*model.Game, error) {
					bots := []*model.Bot{}
					for i := 0; i < 5; i++ {
						botOpts := model.BotOptions{
							Id:        fmt.Sprintf("bot_id%d", i+1),
							Name:      fmt.Sprintf("bot%d", i+1),
							TypeOfBot: "AI",
						}
						bot, _ := model.NewBot(botOpts)
						bots = append(bots, bot)
					}
					game, _ := model.NewGame(
						model.GameOptions{
							Id:               "game_id1",
							State:            "PLAYERS_JOINED",
							CurrentTurnIndex: 0,
							TurnOrder:        []string{"b", "p1", "b", "p2"},
							StateHandled:     false,
							StateTotalTime:   0,
							CreatedAt:        time.Now(),
							UpdatedAt:        time.Now(),
							Bots:             bots,
						},
					)

					return game, nil
				},
			},
			botAccessorMock: nil,
			errorExpected:   true,
			errorString:     "cannot join this game",
		},
		{
			name: "does not error if player is already in game",
			input: &pb.JoinGameRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
			},
			output:         &pb.JoinGameResponse{},
			txProviderMock: &storage.DatabaseTransactionProviderMockSuccess{},
			txShouldCommit: false,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) (*model.Game, error) {
					player1, _ := model.NewPlayer(
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
						bot, _ := model.NewBot(botOpts)
						bots = append(bots, bot)
					}
					bots[1].ConnectPlayer(player1)
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

					return game, nil
				},
			},
			botAccessorMock: nil,
			errorExpected:   false,
			errorString:     "",
		},
		{
			name: "errors if no ai bots in game to convert to human bot",
			input: &pb.JoinGameRequest{
				GameId:   "game_id1",
				PlayerId: "player_id2",
			},
			output:         nil,
			txProviderMock: &storage.DatabaseTransactionProviderMockSuccess{},
			txShouldCommit: false,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) (*model.Game, error) {
					player1, _ := model.NewPlayer(
						model.PlayerOptions{
							Id: "player_id1",
						},
					)
					bots := []*model.Bot{}
					for i := 0; i < 1; i++ {
						botOpts := model.BotOptions{
							Id:        fmt.Sprintf("bot_id%d", i+1),
							Name:      fmt.Sprintf("bot%d", i+1),
							TypeOfBot: "AI",
						}
						bot, _ := model.NewBot(botOpts)
						bots = append(bots, bot)
					}
					bots[0].ConnectPlayer(player1)
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

					return game, nil
				},
			},
			botAccessorMock: nil,
			errorExpected:   true,
			errorString:     "no AI bots in the game",
		},
		{
			name: "errors if unable to update bot",
			input: &pb.JoinGameRequest{
				GameId:   "game_id1",
				PlayerId: "player_id2",
			},
			output:         nil,
			txProviderMock: &storage.DatabaseTransactionProviderMockSuccess{},
			txShouldCommit: false,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) (*model.Game, error) {
					player1, _ := model.NewPlayer(
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
						bot, _ := model.NewBot(botOpts)
						bots = append(bots, bot)
					}
					bots[0].ConnectPlayer(player1)
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
					return game, nil
				},
			},
			botAccessorMock: &storage.BotAccessorMockFailure{},
			errorExpected:   true,
			errorString:     "unable to update bot",
		},
		{
			name: "errors if unable to update game state",
			input: &pb.JoinGameRequest{
				GameId:   "game_id1",
				PlayerId: "player_id2",
			},
			output:         nil,
			txProviderMock: &storage.DatabaseTransactionProviderMockSuccess{},
			txShouldCommit: false,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) (*model.Game, error) {
					player1, _ := model.NewPlayer(
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
						bot, _ := model.NewBot(botOpts)
						bots = append(bots, bot)
					}
					bots[0].ConnectPlayer(player1)
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
					return game, nil
				},
				UpdateGameStateIfEnoughPlayersHaveJoinedUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) error {
					return errors.New("unable to update game")
				},
			},
			botAccessorMock: &storage.BotAccessorMockSuccess{},
			errorExpected:   true,
			errorString:     "unable to update game",
		},
		{
			name: "success",
			input: &pb.JoinGameRequest{
				GameId:   "game_id1",
				PlayerId: "player_id2",
			},
			output:         &pb.JoinGameResponse{},
			txProviderMock: &storage.DatabaseTransactionProviderMockSuccess{},
			txShouldCommit: false,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) (*model.Game, error) {
					player1, _ := model.NewPlayer(
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
						bot, _ := model.NewBot(botOpts)
						bots = append(bots, bot)
					}
					bots[0].ConnectPlayer(player1)
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
					return game, nil
				},
				UpdateGameStateIfEnoughPlayersHaveJoinedUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) error {
					return nil
				},
			},
			botAccessorMock: &storage.BotAccessorMockSuccess{},
			errorExpected:   false,
			errorString:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithDatabaseTransactionProviderMock(tt.txProviderMock),
					storage.WithGameAccessorMock(tt.gameAccessorMock),
					storage.WithBotAccessorMock(tt.botAccessorMock),
				),
			})

			rand.Seed(0)
			response, err := server.JoinGame(
				context.Background(),
				tt.input,
			)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.output, response)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)

			}
			if tt.txShouldCommit {
				// TODO: Commit was called
			} else {
				// TODO: Check rollback was called
			}
		})
	}
}

func Test_GetGameForPlayer(t *testing.T) {
	player1, _ := model.NewPlayer(
		model.PlayerOptions{
			Id: "player_id1",
		},
	)
	player2, _ := model.NewPlayer(
		model.PlayerOptions{
			Id: "player_id2",
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
	bots[4].ConnectPlayer(player1)
	bots[3].ConnectPlayer(player2)
	stateHandledTime := time.Now()
	game, _ := model.NewGame(
		model.GameOptions{
			Id:               "game_id1",
			State:            "WAITING_FOR_PLAYER_QUESTION",
			CurrentTurnIndex: 1,
			TurnOrder:        []string{"bot_id4", "bot_id5", "bot_id3", "bot_id2"},
			StateHandled:     false,
			StateHandledAt:   &stateHandledTime,
			StateTotalTime:   30,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			Bots:             bots,
		},
	)
	tests := []struct {
		name           string
		input          *pb.GetGameForPlayerRequest
		output         *pb.GetGameForPlayerResponse
		gameGetterMock storage.GameAccessor
		errorExpected  bool
		errorString    string
	}{
		{
			name: "errors if get game fails",
			input: &pb.GetGameForPlayerRequest{
				GameId: "game_id1",
			},
			output:         nil,
			gameGetterMock: &storage.GameGetterMockFailure{},
			errorExpected:  true,
			errorString:    "unable to get game",
		},
		{
			name: "errors if player is not in the game",
			input: &pb.GetGameForPlayerRequest{
				GameId:   "game_id1",
				PlayerId: "player_id3",
			},
			output:         nil,
			gameGetterMock: &storage.GameGetterMockSuccess{Game: game},
			errorExpected:  true,
			errorString:    "Unable to get game game_id1 for player player_id3",
		},
		{
			name: "returns correct game view for player 1",
			input: &pb.GetGameForPlayerRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
			},
			output: &pb.GetGameForPlayerResponse{
				State:          "WAITING_ON_YOU_TO_ASK_A_QUESTION",
				DisplayMessage: "Please type a question. OR Click suggest for help!",
				StateStartedAt: timestamppb.New(stateHandledTime),
				StateTotalTime: 30,
				LastQuestion:   "no question",
				MyBotId:        "bot_id5",
				Bots: []*pb.Bot{
					{
						Id:   "bot_id1",
						Name: "bot1",
						BotMessages: []*pb.BotMessage{
							{Text: "Q1: what is your name?"},
							{Text: "A1: My name is Antony Gonsalvez"},
							{Text: "Q2: Where is the gold?"},
							{Text: "A2: what gold!"},
						},
					},
					{
						Id:   "bot_id2",
						Name: "bot2",
						BotMessages: []*pb.BotMessage{
							{Text: "Q1: What is your name?"},
							{Text: "A1: Bot 2 Dot 2"},
						},
					},
					{
						Id:          "bot_id3",
						Name:        "bot3",
						BotMessages: []*pb.BotMessage{},
					},
					{
						Id:          "bot_id4",
						Name:        "bot4",
						BotMessages: []*pb.BotMessage{},
					},
					{
						Id:          "bot_id5",
						Name:        "bot5",
						BotMessages: []*pb.BotMessage{},
					},
				},
			},
			gameGetterMock: &storage.GameGetterMockSuccess{Game: game},
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "returns correct game view for player 2",
			input: &pb.GetGameForPlayerRequest{
				GameId:   "game_id1",
				PlayerId: "player_id2",
			},
			output: &pb.GetGameForPlayerResponse{
				State:          "WAITING_ON_BOT_TO_ASK_A_QUESTION",
				DisplayMessage: "Please wait as someone is asking a question",
				StateStartedAt: timestamppb.New(stateHandledTime),
				StateTotalTime: 30,
				LastQuestion:   "no question",
				MyBotId:        "bot_id4",
				Bots: []*pb.Bot{
					{
						Id:   "bot_id1",
						Name: "bot1",
						BotMessages: []*pb.BotMessage{
							{Text: "Q1: what is your name?"},
							{Text: "A1: My name is Antony Gonsalvez"},
							{Text: "Q2: Where is the gold?"},
							{Text: "A2: what gold!"},
						},
					},
					{
						Id:   "bot_id2",
						Name: "bot2",
						BotMessages: []*pb.BotMessage{
							{Text: "Q1: What is your name?"},
							{Text: "A1: Bot 2 Dot 2"},
						},
					},
					{
						Id:          "bot_id3",
						Name:        "bot3",
						BotMessages: []*pb.BotMessage{},
					},
					{
						Id:          "bot_id4",
						Name:        "bot4",
						BotMessages: []*pb.BotMessage{},
					},
					{
						Id:          "bot_id5",
						Name:        "bot5",
						BotMessages: []*pb.BotMessage{},
					},
				},
			},
			gameGetterMock: &storage.GameGetterMockSuccess{Game: game},
			errorExpected:  false,
			errorString:    "",
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

			response, err := server.GetGameForPlayer(
				context.Background(),
				tt.input,
			)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.output, response)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}

func Test_GetGamesForPlayer(t *testing.T) {
	tests := []struct {
		name             string
		input            *pb.GetGamesForPlayerRequest
		output           *pb.GetGamesForPlayerResponse
		gameAccessorMock storage.GameAccessor
		errorExpected    bool
		errorString      string
	}{
		{
			name:             "test runs successfully",
			input:            &pb.GetGamesForPlayerRequest{PlayerId: "player_id1"},
			output:           &pb.GetGamesForPlayerResponse{GameIds: []string{"game_id1", "game_id2"}},
			gameAccessorMock: &storage.GamesGetterMockSuccess{GameIds: []string{"game_id1", "game_id2"}},
			errorExpected:    false,
			errorString:      "",
		},
		{
			name:             "errors if getting games fails",
			input:            &pb.GetGamesForPlayerRequest{PlayerId: "player_id1"},
			output:           nil,
			gameAccessorMock: &storage.GamesGetterMockFailure{},
			errorExpected:    true,
			errorString:      "unable to get games",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := NewServer(ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithGameAccessorMock(
						tt.gameAccessorMock,
					),
				),
			})

			response, err := server.GetGamesForPlayer(
				context.Background(),
				tt.input,
			)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.output, response)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}
