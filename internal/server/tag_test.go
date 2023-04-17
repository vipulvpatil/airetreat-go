package server

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

func Test_Tag(t *testing.T) {
	tests := []struct {
		name             string
		input            *pb.TagRequest
		output           *pb.TagResponse
		transactionMock  *storage.DatabaseTransactionMock
		gameAccessorMock storage.GameAccessor
		txShouldCommit   bool
		errorExpected    bool
		errorString      string
	}{
		{
			name: "test runs successfully if game is not finished",
			input: &pb.TagRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
				BotId:    "bot_id2",
			},
			output:          &pb.TagResponse{},
			transactionMock: &storage.DatabaseTransactionMock{},
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) (*model.Game, error) {
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
						bot, _ := model.NewBot(botOpts)
						bots = append(bots, bot)
					}
					bots[0].ConnectPlayer(player1)
					bots[1].ConnectPlayer(player2)
					game, _ := model.NewGame(
						model.GameOptions{
							Id:               "game_id1",
							State:            "WAITING_FOR_HUMAN_QUESTION",
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
				UpdateGameStateUsingTransactionInternal: func(gameId string, updateOpts storage.GameUpdateOptions, transaction storage.DatabaseTransaction) error {
					expectedState := "FINISHED"
					result := "bot1 tagged bot2 and won."
					winningBotId := "bot_id1"

					assert.Equal(t, storage.GameUpdateOptions{
						State:        &expectedState,
						Result:       &result,
						WinningBotId: &winningBotId,
					}, updateOpts, "game state should be updated with correct update options")
					return nil
				},
			},
			txShouldCommit: true,
			errorExpected:  false,
			errorString:    "",
		},

		{
			name:             "errors if unable to get transaction",
			input:            &pb.TagRequest{},
			output:           nil,
			transactionMock:  nil,
			gameAccessorMock: nil,
			txShouldCommit:   false,
			errorExpected:    true,
			errorString:      "unable to begin a db transaction",
		},
		{
			name: "errors if cannot get game",
			input: &pb.TagRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
				BotId:    "bot_id1",
			},
			output:          &pb.TagResponse{},
			transactionMock: &storage.DatabaseTransactionMock{},
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) (*model.Game, error) {
					return nil, errors.New("cannot get game")
				},
			},
			txShouldCommit: false,
			errorExpected:  true,
			errorString:    "cannot get game",
		},
		{
			name: "errors if player not in game",
			input: &pb.TagRequest{
				GameId:   "game_id1",
				PlayerId: "player_id2",
				BotId:    "bot_id1",
			},
			output:          &pb.TagResponse{},
			transactionMock: &storage.DatabaseTransactionMock{},
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
							Id:                      "game_id1",
							State:                   "WAITING_FOR_HUMAN_ANSWER",
							CurrentTurnIndex:        2,
							TurnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
							StateHandled:            false,
							StateTotalTime:          0,
							LastQuestion:            "what is the answer?",
							LastQuestionTargetBotId: "bot_id1",
							CreatedAt:               time.Now(),
							UpdatedAt:               time.Now(),
							Bots:                    bots,
						},
					)
					return game, nil
				},
			},
			txShouldCommit: false,
			errorExpected:  true,
			errorString:    "incorrect game",
		},
		{
			name: "errors if unable to get game update after tag",
			input: &pb.TagRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
				BotId:    "bot_id1",
			},
			output:          &pb.TagResponse{},
			transactionMock: &storage.DatabaseTransactionMock{},
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
							Id:                      "game_id1",
							State:                   "FINISHED",
							CurrentTurnIndex:        2,
							TurnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
							StateHandled:            false,
							StateTotalTime:          0,
							LastQuestion:            "what is the answer?",
							LastQuestionTargetBotId: "bot_id1",
							CreatedAt:               time.Now(),
							UpdatedAt:               time.Now(),
							Bots:                    bots,
						},
					)
					return game, nil
				},
			},
			txShouldCommit: false,
			errorExpected:  true,
			errorString:    "game has already finished",
		},
		{
			name: "errors if unable to update game state",
			input: &pb.TagRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
				BotId:    "bot_id3",
			},
			output:          &pb.TagResponse{},
			transactionMock: &storage.DatabaseTransactionMock{},
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) (*model.Game, error) {
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
						bot, _ := model.NewBot(botOpts)
						bots = append(bots, bot)
					}
					bots[0].ConnectPlayer(player1)
					bots[1].ConnectPlayer(player2)
					game, _ := model.NewGame(
						model.GameOptions{
							Id:                      "game_id1",
							State:                   "WAITING_FOR_HUMAN_ANSWER",
							CurrentTurnIndex:        2,
							TurnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
							StateHandled:            false,
							StateTotalTime:          0,
							LastQuestion:            "what is the answer?",
							LastQuestionTargetBotId: "bot_id1",
							CreatedAt:               time.Now(),
							UpdatedAt:               time.Now(),
							Bots:                    bots,
						},
					)
					return game, nil
				},
				UpdateGameStateUsingTransactionInternal: func(gameId string, updateOpts storage.GameUpdateOptions, transaction storage.DatabaseTransaction) error {
					return errors.New("could not update game")
				},
			},
			txShouldCommit: false,
			errorExpected:  true,
			errorString:    "could not update game",
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
				),
			})

			response, err := server.Tag(
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