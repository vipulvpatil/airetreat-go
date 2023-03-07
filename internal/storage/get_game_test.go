package storage

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/model"
)

// GetGameUsingTransaction

func Test_GetGame(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		outputFunc      func() *model.Game
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		errorExpected   bool
		errorString     string
	}{
		{
			name:  "errors when gameId is blank",
			input: "",
			outputFunc: func() *model.Game {
				return nil
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "cannot getGame for a blank gameId",
		},
		{
			name:  "errors nicely if no game found",
			input: "game_id1",
			outputFunc: func() *model.Game {
				return nil
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "game not found: game_id1",
		},
		{
			name:  "bad error when found game with bad rows",
			input: "game_id1",
			outputFunc: func() *model.Game {
				return nil
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" (
						"id", "state", "current_turn_index", "turn_order", "state_handled"
					)
					VALUES (
						'game_id1', 'STARTED', 0, Array['b','p1','b','p2'], 'false'
					)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{
					Query: `DELETE FROM public."games" WHERE id = 'game_id1'`,
				},
			},
			errorExpected: true,
			errorString:   "THIS IS BAD: failed while scanning rows: sql: Scan error on column index 9, name \"id\": converting NULL to string is unsupported",
		},
		{
			name:  "error when found bot with bad data",
			input: "game_id1",
			outputFunc: func() *model.Game {
				return nil
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" (
						"id", "state", "current_turn_index", "turn_order", "state_handled"
					)
					VALUES (
						'game_id1', 'STARTED', 0, Array['b','p1','b','p2'], false
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id1', 'bot1', 'WHAT', 'game_id1'
					)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{
					Query: `DELETE FROM public."games" WHERE id = 'game_id1'`,
				},
			},
			errorExpected: true,
			errorString:   "THIS IS BAD: failed to create bot: cannot create bot with an invalid botType",
		},
		{
			name:  "error when found game with bad data",
			input: "game_id1",
			outputFunc: func() *model.Game {
				return nil
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" (
						"id", "state", "current_turn_index", "turn_order", "state_handled"
					)
					VALUES (
						'game_id1', 'NOTSTARTED', 0, Array['b','p1','b','p2'], false
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id1', 'bot1', 'AI', 'game_id1'
					)`,
				},
				{
					Query: `INSERT INTO public."players" ("id") VALUES ('player_id1')`,
				},
				{
					Query: `UPDATE public."bots" SET
					"player_id" = 'player_id1',
					"type" = 'HUMAN'
					WHERE id = 'bot_id1'`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'game_id1'`},
				{Query: `DELETE FROM public."players" WHERE id = 'player_id1'`},
			},
			errorExpected: true,
			errorString:   "THIS IS BAD: failed to create game: cannot create game with an invalid state",
		},
		{
			name:  "gets a game successfully",
			input: "game_id1",
			outputFunc: func() *model.Game {
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
				return game
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" (
						"id", "state", "current_turn_index", "turn_order", "state_handled", "state_handled_at"
					)
					VALUES (
						'game_id1', 'STARTED', 0, Array['b','p1','b','p2'], false, current_timestamp
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id1', 'bot1', 'AI', 'game_id1'
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id2', 'bot2', 'AI', 'game_id1'
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id3', 'bot3', 'AI', 'game_id1'
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id4', 'bot4', 'AI', 'game_id1'
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id5', 'bot5', 'AI', 'game_id1'
					)`,
				},
				{Query: `INSERT INTO public."players" ("id") VALUES ('player_id1')`},
				{
					Query: `UPDATE public."bots" SET
					"player_id" = 'player_id1',
					"type" = 'HUMAN'
					WHERE id = 'bot_id5'`,
				},
				{Query: `INSERT INTO public."messages" ("id", "bot_id", "text") VALUES ('message_id1', 'bot_id1', 'Q1: what is your name?')`},
				{Query: `INSERT INTO public."messages" ("id", "bot_id", "text") VALUES ('message_id2', 'bot_id1', 'A1: My name is Antony Gonsalvez')`},
				{Query: `INSERT INTO public."messages" ("id", "bot_id", "text") VALUES ('message_id3', 'bot_id2', 'Q1: What is your name?')`},
				{Query: `INSERT INTO public."messages" ("id", "bot_id", "text") VALUES ('message_id4', 'bot_id2', 'A1: Bot 2 Dot 2')`},
				{Query: `INSERT INTO public."messages" ("id", "bot_id", "text") VALUES ('message_id5', 'bot_id1', 'Q2: Where is the gold?')`},
				{Query: `INSERT INTO public."messages" ("id", "bot_id", "text") VALUES ('message_id6', 'bot_id1', 'A2: what gold!')`},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'game_id1'`},
				{Query: `DELETE FROM public."players" WHERE id = 'player_id1'`},
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewDbStorage(
				StorageOptions{
					Db: testDb,
				},
			)

			runSqlOnDb(t, s.db, tt.setupSqlStmts)
			defer runSqlOnDb(t, s.db, tt.cleanupSqlStmts)

			rand.Seed(0)
			result, err := s.GetGame(tt.input)
			if !tt.errorExpected {
				assert.NoError(t, err)
				output := tt.outputFunc()
				model.AssertEqualGame(t, output, result)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}

func Test_GetGameUsingTransaction(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		outputFunc      func() *model.Game
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		errorExpected   bool
		errorString     string
	}{
		{
			name:  "gets a game successfully",
			input: "game_id1",
			outputFunc: func() *model.Game {
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
				return game
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" (
						"id", "state", "current_turn_index", "turn_order", "state_handled", "state_handled_at"
					)
					VALUES (
						'game_id1', 'STARTED', 0, Array['b','p1','b','p2'], false, current_timestamp
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id1', 'bot1', 'AI', 'game_id1'
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id2', 'bot2', 'AI', 'game_id1'
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id3', 'bot3', 'AI', 'game_id1'
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id4', 'bot4', 'AI', 'game_id1'
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id5', 'bot5', 'AI', 'game_id1'
					)`,
				},
				{Query: `INSERT INTO public."players" ("id") VALUES ('player_id1')`},
				{
					Query: `UPDATE public."bots" SET
					"player_id" = 'player_id1',
					"type" = 'HUMAN'
					WHERE id = 'bot_id5'`,
				},
				{Query: `INSERT INTO public."messages" ("id", "bot_id", "text") VALUES ('message_id1', 'bot_id1', 'Q1: what is your name?')`},
				{Query: `INSERT INTO public."messages" ("id", "bot_id", "text") VALUES ('message_id2', 'bot_id1', 'A1: My name is Antony Gonsalvez')`},
				{Query: `INSERT INTO public."messages" ("id", "bot_id", "text") VALUES ('message_id3', 'bot_id2', 'Q1: What is your name?')`},
				{Query: `INSERT INTO public."messages" ("id", "bot_id", "text") VALUES ('message_id4', 'bot_id2', 'A1: Bot 2 Dot 2')`},
				{Query: `INSERT INTO public."messages" ("id", "bot_id", "text") VALUES ('message_id5', 'bot_id1', 'Q2: Where is the gold?')`},
				{Query: `INSERT INTO public."messages" ("id", "bot_id", "text") VALUES ('message_id6', 'bot_id1', 'A2: what gold!')`},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'game_id1'`},
				{Query: `DELETE FROM public."players" WHERE id = 'player_id1'`},
			},
			errorExpected: false,
			errorString:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewDbStorage(
				StorageOptions{
					Db: testDb,
				},
			)

			runSqlOnDb(t, s.db, tt.setupSqlStmts)
			defer runSqlOnDb(t, s.db, tt.cleanupSqlStmts)

			tx, err := s.GetTx()
			assert.NoError(t, err)
			result, err := s.GetGameUsingTransaction(tt.input, tx)
			tx.Commit()
			if !tt.errorExpected {
				assert.NoError(t, err)
				output := tt.outputFunc()
				model.AssertEqualGame(t, output, result)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}
