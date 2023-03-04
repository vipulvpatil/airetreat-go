package storage

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Game_GetUnhandledGameIdsForState(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		output          []string
		setupSqlStmts   []string
		cleanupSqlStmts []string
	}{
		{
			name:            "returns nil when incorrect game state is passed in",
			input:           "",
			output:          nil,
			setupSqlStmts:   []string{},
			cleanupSqlStmts: []string{},
		},
		{
			name:   "returns a list of unhandled game Ids matching state",
			input:  "PLAYERS_JOINED",
			output: []string{"game_id3", "game_id1"},
			setupSqlStmts: []string{
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
				VALUES ('game_id1', 'PLAYERS_JOINED', 0, Array['b','p1','b','p2'], false)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
				VALUES ('bot_id1', 'bot1', 'AI', 'game_id1')`,
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
				VALUES ('game_id2', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
				VALUES ('bot_id2', 'bot2', 'AI', 'game_id2')`,
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
					VALUES ('game_id3', 'PLAYERS_JOINED', 0, Array['b','p1','b','p2'], false)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id3', 'bot3', 'AI', 'game_id3')`,
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
					VALUES ('game_id4', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id4', 'bot4', 'AI', 'game_id4')`,
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
					VALUES ('game_id5', 'PLAYERS_JOINED', 0, Array['b','p1','b','p2'], true)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id5', 'bot5', 'AI', 'game_id5')`,
			},
			cleanupSqlStmts: []string{
				`DELETE FROM public."games" WHERE id = 'game_id1'`,
				`DELETE FROM public."games" WHERE id = 'game_id2'`,
				`DELETE FROM public."games" WHERE id = 'game_id3'`,
				`DELETE FROM public."games" WHERE id = 'game_id4'`,
				`DELETE FROM public."games" WHERE id = 'game_id5'`,
			},
		},
		{
			name:   "returns an empty list if no games match",
			input:  "PLAYERS_JOINED",
			output: []string{},
			setupSqlStmts: []string{
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
				VALUES ('game_id1', 'PLAYERS_JOINED', 0, Array['b','p1','b','p2'], true)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
				VALUES ('bot_id1', 'bot1', 'AI', 'game_id1')`,
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
				VALUES ('game_id2', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
				VALUES ('bot_id2', 'bot2', 'AI', 'game_id2')`,
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
					VALUES ('game_id3', 'PLAYERS_JOINED', 0, Array['b','p1','b','p2'], true)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id3', 'bot3', 'AI', 'game_id3')`,
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
					VALUES ('game_id4', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id4', 'bot4', 'AI', 'game_id4')`,
				`INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
					VALUES ('game_id5', 'PLAYERS_JOINED', 0, Array['b','p1','b','p2'], true)`,
				`INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id5', 'bot5', 'AI', 'game_id5')`,
			},
			cleanupSqlStmts: []string{
				`DELETE FROM public."games" WHERE id = 'game_id1'`,
				`DELETE FROM public."games" WHERE id = 'game_id2'`,
				`DELETE FROM public."games" WHERE id = 'game_id3'`,
				`DELETE FROM public."games" WHERE id = 'game_id4'`,
				`DELETE FROM public."games" WHERE id = 'game_id5'`,
			},
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
			gameIds := s.GetUnhandledGameIdsForState(tt.input)
			assert.Equal(t, tt.output, gameIds)
		})
	}
}
