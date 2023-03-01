package storage

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Game_GetGames(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		output          []string
		setupSqlStmts   []string
		cleanupSqlStmts []string
		errorExpected   bool
		errorString     string
	}{
		{
			name:            "errors when playerId is blank",
			input:           "",
			output:          []string{},
			setupSqlStmts:   []string{},
			cleanupSqlStmts: []string{},
			errorExpected:   true,
			errorString:     "cannot GetGames for a blank playerId",
		},
		{
			name:   "returns a list of gameIds for the player",
			input:  "player_id1",
			output: []string{"game_id5", "game_id2", "game_id1"},
			setupSqlStmts: []string{
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
				VALUES ('game_id1', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
				VALUES ('bot_id1', 'bot1', 'AI', 'game_id1')`,
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
				VALUES ('game_id2', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
				VALUES ('bot_id2', 'bot2', 'AI', 'game_id2')`,
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
					VALUES ('game_id3', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id3', 'bot3', 'AI', 'game_id3')`,
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
					VALUES ('game_id4', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id4', 'bot4', 'AI', 'game_id4')`,
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
					VALUES ('game_id5', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id5', 'bot5', 'AI', 'game_id5')`,
				`INSERT INTO public."players" ("id") VALUES ('player_id1')`,
				`INSERT INTO public."players" ("id") VALUES ('player_id2')`,
				`UPDATE public."bots" SET "player_id" = 'player_id1', "type" = 'HUMAN' WHERE id = 'bot_id1'`,
				`UPDATE public."bots" SET "player_id" = 'player_id1', "type" = 'HUMAN' WHERE id = 'bot_id2'`,
				`UPDATE public."bots" SET "player_id" = 'player_id2', "type" = 'HUMAN' WHERE id = 'bot_id3'`,
				`UPDATE public."bots" SET "player_id" = 'player_id2', "type" = 'HUMAN' WHERE id = 'bot_id4'`,
				`UPDATE public."bots" SET "player_id" = 'player_id1', "type" = 'HUMAN' WHERE id = 'bot_id5'`,
			},
			cleanupSqlStmts: []string{
				`DELETE FROM public."games" WHERE id = 'game_id1'`,
				`DELETE FROM public."games" WHERE id = 'game_id2'`,
				`DELETE FROM public."games" WHERE id = 'game_id3'`,
				`DELETE FROM public."games" WHERE id = 'game_id4'`,
				`DELETE FROM public."games" WHERE id = 'game_id5'`,
				`DELETE FROM public."players" WHERE id = 'player_id1'`,
				`DELETE FROM public."players" WHERE id = 'player_id2'`,
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
			games, err := s.GetGames(tt.input)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.Equal(t, tt.output, games)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}
