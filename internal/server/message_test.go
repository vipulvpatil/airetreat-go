package server

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
	pb "github.com/vipulvpatil/airetreat-go/protos"
)

func Test_SendMessage(t *testing.T) {
	tests := []struct {
		name               string
		input              *pb.SendMessageRequest
		output             *pb.SendMessageResponse
		transactionMock    *storage.DatabaseTransactionMock
		gameAccessorMock   storage.GameAccessor
		messageCreatorMock storage.MessageCreator
		txShouldCommit     bool
		errorExpected      bool
		errorString        string
	}{
		{
			name: "test runs successfully for a game waiting for human question",
			input: &pb.SendMessageRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
				BotId:    "bot_id2",
				Text:     "question message",
			},
			output:             &pb.SendMessageResponse{},
			transactionMock:    &storage.DatabaseTransactionMock{},
			messageCreatorMock: &storage.MessageCreatorMockSuccess{},
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
					expectedState := "WAITING_FOR_AI_ANSWER"
					expectedStateHandled := false
					expectedLastQuestion := "question message"
					expectedLastQuestionTargetBotId := "bot_id2"

					assert.Equal(t, storage.GameUpdateOptions{
						State:                   &expectedState,
						StateHandled:            &expectedStateHandled,
						LastQuestion:            &expectedLastQuestion,
						LastQuestionTargetBotId: &expectedLastQuestionTargetBotId,
					}, updateOpts, "game state should be updated with correct update options")
					return nil
				},
			},
			txShouldCommit: true,
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "test runs successfully for a game waiting for human answer",
			input: &pb.SendMessageRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
				BotId:    "bot_id1",
				Text:     "answer message",
			},
			output:             &pb.SendMessageResponse{},
			transactionMock:    &storage.DatabaseTransactionMock{},
			messageCreatorMock: &storage.MessageCreatorMockSuccess{},
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
				UpdateGameStateUsingTransactionInternal: func(gameId string, updateOpts storage.GameUpdateOptions, transaction storage.DatabaseTransaction) error {
					expectedState := "WAITING_FOR_AI_QUESTION"
					expectedStateHandled := false
					expectedCurrentTurnIndex := int64(3)

					assert.Equal(t, storage.GameUpdateOptions{
						State:            &expectedState,
						StateHandled:     &expectedStateHandled,
						CurrentTurnIndex: &expectedCurrentTurnIndex,
					}, updateOpts, "game state should be updated with correct update options")
					return nil
				},
			},
			txShouldCommit: true,
			errorExpected:  false,
			errorString:    "",
		},
		{
			name:               "errors if unable to get transaction",
			input:              &pb.SendMessageRequest{},
			output:             nil,
			transactionMock:    nil,
			messageCreatorMock: nil,
			gameAccessorMock:   nil,
			txShouldCommit:     false,
			errorExpected:      true,
			errorString:        "unable to begin a db transaction",
		},
		{
			name: "errors if cannot get game",
			input: &pb.SendMessageRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
				BotId:    "bot_id1",
				Text:     "answer message",
			},
			output:             &pb.SendMessageResponse{},
			transactionMock:    &storage.DatabaseTransactionMock{},
			messageCreatorMock: nil,
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
			input: &pb.SendMessageRequest{
				GameId:   "game_id1",
				PlayerId: "player_id2",
				BotId:    "bot_id1",
				Text:     "answer message",
			},
			output:             &pb.SendMessageResponse{},
			transactionMock:    &storage.DatabaseTransactionMock{},
			messageCreatorMock: nil,
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
			name: "errors if unable to get game update after incoming message",
			input: &pb.SendMessageRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
				BotId:    "bot_id1",
				Text:     "answer message",
			},
			output:             &pb.SendMessageResponse{},
			transactionMock:    &storage.DatabaseTransactionMock{},
			messageCreatorMock: nil,
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
							State:                   "WAITING_FOR_AI_ANSWER",
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
			errorString:    "expecting AI message but did not receive one",
		},
		{
			name: "errors if unable to update game state",
			input: &pb.SendMessageRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
				BotId:    "bot_id1",
				Text:     "answer message",
			},
			output:             &pb.SendMessageResponse{},
			transactionMock:    &storage.DatabaseTransactionMock{},
			messageCreatorMock: &storage.MessageCreatorMockSuccess{},
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
				UpdateGameStateUsingTransactionInternal: func(gameId string, updateOpts storage.GameUpdateOptions, transaction storage.DatabaseTransaction) error {
					return errors.New("could not update game")
				},
			},
			txShouldCommit: false,
			errorExpected:  true,
			errorString:    "could not update game",
		},
		{
			name: "errors and rollsback if unable to create message",
			input: &pb.SendMessageRequest{
				GameId:   "game_id1",
				PlayerId: "player_id1",
				BotId:    "bot_id1",
				Text:     "answer message",
			},
			output:             &pb.SendMessageResponse{},
			transactionMock:    &storage.DatabaseTransactionMock{},
			messageCreatorMock: &storage.MessageCreatorMockFailure{},
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
				UpdateGameStateUsingTransactionInternal: func(gameId string, updateOpts storage.GameUpdateOptions, transaction storage.DatabaseTransaction) error {
					expectedState := "WAITING_FOR_AI_QUESTION"
					expectedStateHandled := false
					expectedCurrentTurnIndex := int64(3)

					assert.Equal(t, storage.GameUpdateOptions{
						State:            &expectedState,
						StateHandled:     &expectedStateHandled,
						CurrentTurnIndex: &expectedCurrentTurnIndex,
					}, updateOpts, "game state should be updated with correct update options")
					return nil
				},
			},
			txShouldCommit: false,
			errorExpected:  true,
			errorString:    "unable to create message",
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
					storage.WithMessageCreatorMock(tt.messageCreatorMock),
				),
			})

			response, err := server.SendMessage(
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
