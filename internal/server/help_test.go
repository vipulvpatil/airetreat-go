package server

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/clients/openai"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

func Test_Help(t *testing.T) {
	tests := []struct {
		name             string
		input            *pb.HelpRequest
		output           *pb.HelpResponse
		transactionMock  *storage.DatabaseTransactionMock
		gameAccessorMock storage.GameAccessor
		botAccessorMock  storage.BotAccessor
		openAiResponse   string
		txShouldCommit   bool
		errorExpected    bool
		errorString      string
	}{
		{
			name: "errors if unable to get game",
			input: &pb.HelpRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
			},
			output:          nil,
			transactionMock: nil,
			openAiResponse:  "",
			txShouldCommit:  false,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(gameId string) (*model.Game, error) {
					return nil, errors.New("unable to get game")
				},
			},
			botAccessorMock: nil,
			errorExpected:   true,
			errorString:     "unable to get game",
		},
		{
			name: "errors if player not in game",
			input: &pb.HelpRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
			},
			output:          nil,
			transactionMock: nil,
			openAiResponse:  "",
			txShouldCommit:  false,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(gameId string) (*model.Game, error) {
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
							State:            "WAITING_FOR_AI_QUESTION",
							CurrentTurnIndex: 0,
							TurnOrder:        []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
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
			errorString:     "incorrect game",
		},
		{
			name: "errors if bot has used all help",
			input: &pb.HelpRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
			},
			output:          nil,
			transactionMock: nil,
			openAiResponse:  "",
			txShouldCommit:  false,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(gameId string) (*model.Game, error) {
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
							HelpCount: 0,
						}
						bot, _ := model.NewBot(botOpts)
						bots = append(bots, bot)
					}
					bots[1].ConnectPlayer(player1)
					game, _ := model.NewGame(
						model.GameOptions{
							Id:               "game_id1",
							State:            "WAITING_FOR_AI_QUESTION",
							CurrentTurnIndex: 0,
							TurnOrder:        []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
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
			errorString:     "no more help possible",
		},
		{
			name: "errors if not this bot's turn",
			input: &pb.HelpRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
			},
			output:          nil,
			transactionMock: nil,
			openAiResponse:  "",
			txShouldCommit:  false,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(gameId string) (*model.Game, error) {
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
							HelpCount: 3,
						}
						bot, _ := model.NewBot(botOpts)
						bots = append(bots, bot)
					}
					bots[1].ConnectPlayer(player1)
					game, _ := model.NewGame(
						model.GameOptions{
							Id:               "game_id1",
							State:            "WAITING_FOR_AI_QUESTION",
							CurrentTurnIndex: 3,
							TurnOrder:        []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
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
			errorString:     "please wait for your turn",
		},
		{
			name: "errors if unable to get transaction",
			input: &pb.HelpRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
			},
			output:          nil,
			transactionMock: nil,
			openAiResponse:  "",
			txShouldCommit:  false,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(gameId string) (*model.Game, error) {
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
							HelpCount: 3,
						}
						bot, _ := model.NewBot(botOpts)
						bots = append(bots, bot)
					}
					bots[1].ConnectPlayer(player1)
					game, _ := model.NewGame(
						model.GameOptions{
							Id:               "game_id1",
							State:            "WAITING_FOR_AI_QUESTION",
							CurrentTurnIndex: 1,
							TurnOrder:        []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
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
			errorString:     "unable to begin a db transaction",
		},
		{
			name: "errors if unable to update bot",
			input: &pb.HelpRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
			},
			output:          nil,
			transactionMock: &storage.DatabaseTransactionMock{},
			openAiResponse:  "",
			txShouldCommit:  false,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(gameId string) (*model.Game, error) {
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
							HelpCount: 3,
						}
						bot, _ := model.NewBot(botOpts)
						bots = append(bots, bot)
					}
					bots[1].ConnectPlayer(player1)
					game, _ := model.NewGame(
						model.GameOptions{
							Id:               "game_id1",
							State:            "WAITING_FOR_AI_QUESTION",
							CurrentTurnIndex: 1,
							TurnOrder:        []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
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
			name: "success when waiting for human question",
			input: &pb.HelpRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
			},
			output:          &pb.HelpResponse{Text: "sample response"},
			transactionMock: &storage.DatabaseTransactionMock{},
			openAiResponse:  "sample response",
			txShouldCommit:  true,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(gameId string) (*model.Game, error) {
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
							HelpCount: 3,
						}
						bot, _ := model.NewBot(botOpts)
						bots = append(bots, bot)
					}
					bots[1].ConnectPlayer(player1)
					game, _ := model.NewGame(
						model.GameOptions{
							Id:               "game_id1",
							State:            "WAITING_FOR_HUMAN_QUESTION",
							CurrentTurnIndex: 1,
							TurnOrder:        []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
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
			botAccessorMock: &storage.BotAccessorMockSuccess{},
			errorExpected:   false,
			errorString:     "",
		},
		{
			name: "success when waiting for human answer",
			input: &pb.HelpRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
			},
			output:          &pb.HelpResponse{Text: "sample response"},
			transactionMock: &storage.DatabaseTransactionMock{},
			openAiResponse:  "sample response",
			txShouldCommit:  true,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(gameId string) (*model.Game, error) {
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
							HelpCount: 3,
						}
						bot, _ := model.NewBot(botOpts)
						bots = append(bots, bot)
					}
					bots[1].ConnectPlayer(player1)
					game, _ := model.NewGame(
						model.GameOptions{
							Id:                      "game_id1",
							State:                   "WAITING_FOR_HUMAN_ANSWER",
							CurrentTurnIndex:        1,
							TurnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
							StateHandled:            false,
							StateTotalTime:          0,
							LastQuestionTargetBotId: "bot_id2",
							CreatedAt:               time.Now(),
							UpdatedAt:               time.Now(),
							Bots:                    bots,
						},
					)
					return game, nil
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
					storage.WithDatabaseTransactionProviderMock(&storage.DatabaseTransactionProviderMock{
						Transaction: tt.transactionMock,
					}),
					storage.WithGameAccessorMock(tt.gameAccessorMock),
					storage.WithBotAccessorMock(tt.botAccessorMock),
				),
				OpenAiClient: &openai.MockClientSuccess{Text: tt.openAiResponse},
				Logger:       &utilities.NullLogger{},
			})

			rand.Seed(0)
			response, err := server.Help(
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
			if tt.transactionMock != nil {
				if tt.txShouldCommit {
					assert.True(t, tt.transactionMock.Committed, "transaction should have committed")
				} else {
					assert.True(t, tt.transactionMock.Rolledback, "transaction should have rolledback")
					assert.False(t, tt.transactionMock.Committed, "transaction should not have committed")
				}
			}
		})
	}
}
