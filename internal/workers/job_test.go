package workers

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/gocraft/work"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/storage"
)

func Test_startGameOncePlayersHaveJoined(t *testing.T) {
	tests := []struct {
		name             string
		input            map[string]interface{}
		gameAccessorMock storage.GameAccessor
		errorExpected    bool
		errorString      string
	}{
		{
			name: "updates game successfully to WAITING_FOR_AI_QUESTION",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(string) (*model.Game, error) {
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
					return model.NewGame(
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
				},
				UpdateGameStateInternal: func(gameId string, opts storage.GameUpdateOptions) error {
					assert.Equal(t, "game_id1", gameId)
					assert.Equal(t, "WAITING_FOR_AI_QUESTION", *opts.State)
					assert.Equal(t, int64(0), *opts.CurrentTurnIndex)
					assert.Equal(t, []string{"bot_id3", "bot_id4", "bot_id2", "bot_id1", "bot_id5"}, opts.TurnOrder)
					return nil
				},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "updates game successfully to WAITING_FOR_HUMAN_QUESTION",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(string) (*model.Game, error) {
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
					bots[2].ConnectPlayer(player)
					return model.NewGame(
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
				},
				UpdateGameStateInternal: func(gameId string, opts storage.GameUpdateOptions) error {
					assert.Equal(t, "game_id1", gameId)
					assert.Equal(t, "WAITING_FOR_HUMAN_QUESTION", *opts.State)
					assert.Equal(t, int64(0), *opts.CurrentTurnIndex)
					assert.Equal(t, []string{"bot_id3", "bot_id4", "bot_id2", "bot_id1", "bot_id5"}, opts.TurnOrder)
					return nil
				},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "errors if game is not in db",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(string) (*model.Game, error) {

					return nil, errors.New("game not in db")
				},
			},
			errorExpected: true,
			errorString:   "game not in db",
		},
		{
			name: "errors if game is in wrong state",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(string) (*model.Game, error) {
					bot, _ := model.NewBot(
						model.BotOptions{
							Id:        "bot_id1",
							Name:      "bot1",
							TypeOfBot: "AI",
						},
					)
					return model.NewGame(
						model.GameOptions{
							Id:               "game_id1",
							State:            "STARTED",
							CurrentTurnIndex: 0,
							TurnOrder:        []string{"b", "p1", "b", "p2"},
							StateHandled:     false,
							StateTotalTime:   0,
							CreatedAt:        time.Now(),
							UpdatedAt:        time.Now(),
							Bots:             []*model.Bot{bot},
						},
					)
				},
			},
			errorExpected: true,
			errorString:   "game should be in PlayersJoined state: game_id1",
		},
		{
			name: "errors if game is has already been handled",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(string) (*model.Game, error) {
					bot, _ := model.NewBot(
						model.BotOptions{
							Id:        "bot_id1",
							Name:      "bot1",
							TypeOfBot: "AI",
						},
					)
					return model.NewGame(
						model.GameOptions{
							Id:               "game_id1",
							State:            "STARTED",
							CurrentTurnIndex: 0,
							TurnOrder:        []string{"b", "p1", "b", "p2"},
							StateHandled:     true,
							StateTotalTime:   0,
							CreatedAt:        time.Now(),
							UpdatedAt:        time.Now(),
							Bots:             []*model.Bot{bot},
						},
					)
				},
			},
			errorExpected: true,
			errorString:   "game has already been handled: game_id1",
		},
		{
			name: "errors if gameId is blank",
			input: map[string]interface{}{
				"gameId": "",
			},
			gameAccessorMock: nil,
			errorExpected:    true,
			errorString:      "gameId is required",
		},
	}

	for _, tt := range tests {
		workerStorage = storage.NewStorageAccessorMock(
			storage.WithGameAccessorMock(tt.gameAccessorMock),
		)

		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(0)
			jc := jobContext{}
			err := jc.startGameOncePlayersHaveJoined(&work.Job{
				Args: tt.input,
			})
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_askQuestionOnBehalfOfBot(t *testing.T) {
	tests := []struct {
		name               string
		input              map[string]interface{}
		transactionMock    *storage.DatabaseTransactionMock
		messageCreatorMock storage.MessageCreator
		gameAccessorMock   storage.GameAccessor
		txShouldCommit     bool
		errorExpected      bool
		errorString        string
	}{
		{
			name: "updates game successfully",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
			transactionMock:    &storage.DatabaseTransactionMock{},
			messageCreatorMock: &storage.MessageCreatorMockSuccess{},
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) (*model.Game, error) {
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
					return model.NewGame(
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
				},
				UpdateGameStateUsingTransactionInternal: func(gameId string, updateOpts storage.GameUpdateOptions, transaction storage.DatabaseTransaction) error {
					expectedState := "WAITING_FOR_AI_ANSWER"
					expectedStateHandled := false
					expectedLastQuestion := "Some question from AI"
					expectedLastQuestionTargetBotId := "bot_id4"

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
			name: "errors if gameId not provided",
			input: map[string]interface{}{
				"gameId": "",
			},
			transactionMock:    nil,
			messageCreatorMock: nil,
			gameAccessorMock:   nil,
			txShouldCommit:     false,
			errorExpected:      true,
			errorString:        "gameId is required",
		},
		{
			name: "errors if unable to get transaction",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
			transactionMock:    nil,
			messageCreatorMock: nil,
			gameAccessorMock:   nil,
			txShouldCommit:     false,
			errorExpected:      true,
			errorString:        "unable to begin a db transaction",
		},
		{
			name: "errors if cannot get game",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
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
			name: "errors if game has already been handled",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
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
							StateHandled:            true,
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
			errorString:    "game has already been handled: game_id1",
		},
		{
			name: "errors if game not in correct state",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
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
			errorString:    "game should be in WaitingForAiQuestion state: game_id1",
		},
		{
			name: "errors if cannot determine bot to ask a question",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
			transactionMock:    &storage.DatabaseTransactionMock{},
			messageCreatorMock: nil,
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) (*model.Game, error) {
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
					game, _ := model.NewGame(
						model.GameOptions{
							Id:                      "game_id1",
							State:                   "WAITING_FOR_AI_QUESTION",
							CurrentTurnIndex:        0,
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
			errorString:    "cannot get random bot from an empty list",
		},
		{
			name: "errors if unable to update game state",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
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
							State:                   "WAITING_FOR_AI_QUESTION",
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
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
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
							State:                   "WAITING_FOR_AI_QUESTION",
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
					expectedState := "WAITING_FOR_AI_ANSWER"
					expectedStateHandled := false
					expectedLastQuestion := "Some question from AI"
					expectedLastQuestionTargetBotId := "bot_id4"

					assert.Equal(t, storage.GameUpdateOptions{
						State:                   &expectedState,
						StateHandled:            &expectedStateHandled,
						LastQuestion:            &expectedLastQuestion,
						LastQuestionTargetBotId: &expectedLastQuestionTargetBotId,
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
		workerStorage = storage.NewStorageAccessorMock(
			storage.WithDatabaseTransactionProviderMock(&storage.DatabaseTransactionProviderMock{
				Transaction: tt.transactionMock,
			}),
			storage.WithGameAccessorMock(tt.gameAccessorMock),
			storage.WithMessageCreatorMock(tt.messageCreatorMock),
		)

		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(0)
			jc := jobContext{}
			err := jc.askQuestionOnBehalfOfBot(&work.Job{
				Args: tt.input,
			})
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
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

func Test_answerQuestionOnBehalfOfBot(t *testing.T) {
	tests := []struct {
		name               string
		input              map[string]interface{}
		transactionMock    *storage.DatabaseTransactionMock
		messageCreatorMock storage.MessageCreator
		gameAccessorMock   storage.GameAccessor
		txShouldCommit     bool
		errorExpected      bool
		errorString        string
	}{
		{
			name: "updates game successfully",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
			transactionMock:    &storage.DatabaseTransactionMock{},
			messageCreatorMock: &storage.MessageCreatorMockSuccess{},
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameUsingTransactionInternal: func(gameId string, transaction storage.DatabaseTransaction) (*model.Game, error) {
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
					return model.NewGame(
						model.GameOptions{
							Id:                      "game_id1",
							State:                   "WAITING_FOR_AI_ANSWER",
							CurrentTurnIndex:        0,
							TurnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
							StateHandled:            false,
							StateTotalTime:          0,
							LastQuestion:            "Here is a question?",
							LastQuestionTargetBotId: "bot_id4",
							CreatedAt:               time.Now(),
							UpdatedAt:               time.Now(),
							Bots:                    bots,
						},
					)
				},
				UpdateGameStateUsingTransactionInternal: func(gameId string, updateOpts storage.GameUpdateOptions, transaction storage.DatabaseTransaction) error {
					expectedState := "WAITING_FOR_AI_QUESTION"
					expectedStateHandled := false
					expectedCurrentTurnIndex := int64(1)

					assert.Equal(t, storage.GameUpdateOptions{
						State:            &expectedState,
						CurrentTurnIndex: &expectedCurrentTurnIndex,
						StateHandled:     &expectedStateHandled,
					}, updateOpts, "game state should be updated with correct update options")
					return nil
				},
			},
			txShouldCommit: true,
			errorExpected:  false,
			errorString:    "",
		},
		{
			name: "errors if gameId not provided",
			input: map[string]interface{}{
				"gameId": "",
			},
			transactionMock:    nil,
			messageCreatorMock: nil,
			gameAccessorMock:   nil,
			txShouldCommit:     false,
			errorExpected:      true,
			errorString:        "gameId is required",
		},
		{
			name: "errors if unable to get transaction",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
			transactionMock:    nil,
			messageCreatorMock: nil,
			gameAccessorMock:   nil,
			txShouldCommit:     false,
			errorExpected:      true,
			errorString:        "unable to begin a db transaction",
		},
		{
			name: "errors if cannot get game",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
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
			name: "errors if game has already been handled",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
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
							StateHandled:            true,
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
			errorString:    "game has already been handled: game_id1",
		},
		{
			name: "errors if game not in correct state",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
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
			errorString:    "game should be in WaitingForAiAnswer state: game_id1",
		},
		{
			name: "errors if unable to update game state",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
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
							State:                   "WAITING_FOR_AI_ANSWER",
							CurrentTurnIndex:        2,
							TurnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
							StateHandled:            false,
							StateTotalTime:          0,
							LastQuestion:            "what is the answer?",
							LastQuestionTargetBotId: "bot_id3",
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
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
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
							State:                   "WAITING_FOR_AI_ANSWER",
							CurrentTurnIndex:        2,
							TurnOrder:               []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"},
							StateHandled:            false,
							StateTotalTime:          0,
							LastQuestion:            "what is the answer?",
							LastQuestionTargetBotId: "bot_id3",
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
						CurrentTurnIndex: &expectedCurrentTurnIndex,
						StateHandled:     &expectedStateHandled,
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
		workerStorage = storage.NewStorageAccessorMock(
			storage.WithDatabaseTransactionProviderMock(&storage.DatabaseTransactionProviderMock{
				Transaction: tt.transactionMock,
			}),
			storage.WithGameAccessorMock(tt.gameAccessorMock),
			storage.WithMessageCreatorMock(tt.messageCreatorMock),
		)

		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(0)
			jc := jobContext{}
			err := jc.answerQuestionOnBehalfOfBot(&work.Job{
				Args: tt.input,
			})
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
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

func Test_deleteExpiredGames(t *testing.T) {
	tests := []struct {
		name             string
		input            map[string]interface{}
		gameAccessorMock storage.GameAccessor
		errorExpected    bool
		errorString      string
	}{
		{
			name: "deletes game if present and is not recently updated",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(string) (*model.Game, error) {
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
					return model.NewGame(
						model.GameOptions{
							Id:               "game_id1",
							State:            "PLAYERS_JOINED",
							CurrentTurnIndex: 0,
							TurnOrder:        []string{"b", "p1", "b", "p2"},
							StateHandled:     false,
							StateTotalTime:   0,
							CreatedAt:        time.Now(),
							UpdatedAt:        time.Now().Add(-10 * time.Hour),
							Bots:             bots,
						},
					)
				},
				DeleteGameInternal: func(gameId string) error {
					assert.Equal(t, "game_id1", gameId)
					return nil
				},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "continues silently if game is present and is recently updated",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(string) (*model.Game, error) {
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
					return model.NewGame(
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
				},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "errors if game is not in db",
			input: map[string]interface{}{
				"gameId": "game_id1",
			},
			gameAccessorMock: &storage.GameAccessorConfigurableMock{
				GetGameInternal: func(string) (*model.Game, error) {

					return nil, errors.New("game not in db")
				},
			},
			errorExpected: true,
			errorString:   "game not in db",
		},
		{
			name: "errors if gameId is blank",
			input: map[string]interface{}{
				"gameId": "",
			},
			gameAccessorMock: nil,
			errorExpected:    true,
			errorString:      "gameId is required",
		},
	}

	for _, tt := range tests {
		workerStorage = storage.NewStorageAccessorMock(
			storage.WithGameAccessorMock(tt.gameAccessorMock),
		)

		t.Run(tt.name, func(t *testing.T) {
			rand.Seed(0)
			jc := jobContext{}
			err := jc.deleteExpiredGames(&work.Job{
				Args: tt.input,
			})
			if tt.errorExpected {
				assert.EqualError(t, err, tt.errorString)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
