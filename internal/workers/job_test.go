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
			name: "updates game successfully to WAITING_FOR_BOT_QUESTION",
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
					assert.Equal(t, "WAITING_FOR_BOT_QUESTION", *opts.State)
					assert.Equal(t, int64(0), *opts.CurrentTurnIndex)
					assert.Equal(t, []string{"bot_id3", "bot_id4", "bot_id2", "bot_id1", "bot_id5"}, opts.TurnOrder)
					return nil
				},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "updates game successfully to WAITING_FOR_PLAYER_QUESTION",
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
					assert.Equal(t, "WAITING_FOR_PLAYER_QUESTION", *opts.State)
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
