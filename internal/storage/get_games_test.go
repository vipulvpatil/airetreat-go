package storage

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Game_GetGames(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		output          []string
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		errorExpected   bool
		errorString     string
	}{
		{
			name:            "errors when playerId is blank",
			input:           "",
			output:          []string{},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "cannot GetGames for a blank playerId",
		},
		{
			name:   "returns a list of gameIds for the player",
			input:  "player_id1",
			output: []string{"game_id5", "game_id2", "game_id1"},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
				VALUES ('game_id1', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
				VALUES ('bot_id1', 'bot1', 'AI', 'game_id1')`,
				},
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
				VALUES ('game_id2', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
				VALUES ('bot_id2', 'bot2', 'AI', 'game_id2')`,
				},
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
					VALUES ('game_id3', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id3', 'bot3', 'AI', 'game_id3')`,
				},
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
					VALUES ('game_id4', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id4', 'bot4', 'AI', 'game_id4')`,
				},
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
					VALUES ('game_id5', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id5', 'bot5', 'AI', 'game_id5')`,
				},
				{Query: `INSERT INTO public."players" ("id") VALUES ('player_id1')`},
				{Query: `INSERT INTO public."players" ("id") VALUES ('player_id2')`},
				{Query: `UPDATE public."bots" SET "player_id" = 'player_id1', "type" = 'HUMAN' WHERE id = 'bot_id1'`},
				{Query: `UPDATE public."bots" SET "player_id" = 'player_id1', "type" = 'HUMAN' WHERE id = 'bot_id2'`},
				{Query: `UPDATE public."bots" SET "player_id" = 'player_id2', "type" = 'HUMAN' WHERE id = 'bot_id3'`},
				{Query: `UPDATE public."bots" SET "player_id" = 'player_id2', "type" = 'HUMAN' WHERE id = 'bot_id4'`},
				{Query: `UPDATE public."bots" SET "player_id" = 'player_id1', "type" = 'HUMAN' WHERE id = 'bot_id5'`},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'game_id1'`},
				{Query: `DELETE FROM public."games" WHERE id = 'game_id2'`},
				{Query: `DELETE FROM public."games" WHERE id = 'game_id3'`},
				{Query: `DELETE FROM public."games" WHERE id = 'game_id4'`},
				{Query: `DELETE FROM public."games" WHERE id = 'game_id5'`},
				{Query: `DELETE FROM public."players" WHERE id = 'player_id1'`},
				{Query: `DELETE FROM public."players" WHERE id = 'player_id2'`},
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

func Test_Game_GetOldGames(t *testing.T) {
	tests := []struct {
		name            string
		input           time.Duration
		output          []string
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		errorExpected   bool
		errorString     string
	}{
		{
			name:            "errors when timeDuration is invalid",
			input:           1 * time.Second,
			output:          []string{},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "invalid game expiry duration. Max acceptable time is -5 minutes.",
		},
		{
			name:   "returns a list of gameIds that are old",
			input:  -1 * time.Hour,
			output: []string{"game_id5", "game_id2", "game_id1"},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled", "created_at")
					VALUES ('game_id1', 'STARTED', 0, Array['b','p1','b','p2'], false, $1)`,
					Args: []any{time.Now().Add(-2 * time.Hour)},
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id1', 'bot1', 'AI', 'game_id1')`,
				},
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled", "created_at")
					VALUES ('game_id2', 'STARTED', 0, Array['b','p1','b','p2'], false, $1)`,
					Args: []any{time.Now().Add(-2 * time.Hour)},
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id2', 'bot2', 'AI', 'game_id2')`,
				},
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
					VALUES ('game_id3', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id3', 'bot3', 'AI', 'game_id3')`,
				},
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled")
					VALUES ('game_id4', 'STARTED', 0, Array['b','p1','b','p2'], false)`,
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id4', 'bot4', 'AI', 'game_id4')`,
				},
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled", "created_at")
					VALUES ('game_id5', 'STARTED', 0, Array['b','p1','b','p2'], false, $1)`,
					Args: []any{time.Now().Add(-2 * time.Hour)},
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id5', 'bot5', 'AI', 'game_id5')`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'game_id1'`},
				{Query: `DELETE FROM public."games" WHERE id = 'game_id2'`},
				{Query: `DELETE FROM public."games" WHERE id = 'game_id3'`},
				{Query: `DELETE FROM public."games" WHERE id = 'game_id4'`},
				{Query: `DELETE FROM public."games" WHERE id = 'game_id5'`},
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
			games, err := s.GetOldGames(tt.input)
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
func Test_Game_GetPublicJoinableGames(t *testing.T) {
	tests := []struct {
		name            string
		output          []string
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		errorExpected   bool
		errorString     string
	}{
		{
			name:   "returns a list of gameIds that are old",
			output: []string{"game_id5", "game_id1"},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled", "created_at", "public")
					VALUES ('game_id1', 'STARTED', 0, Array['b','p1','b','p2'], false, $1, true)`,
					Args: []any{time.Now().Add(-15 * time.Minute)},
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id1', 'bot1', 'HUMAN', 'game_id1')`,
				},
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled", "created_at", "public")
					VALUES ('game_id2', 'STARTED', 0, Array['b','p1','b','p2'], false, $1, true)`,
					Args: []any{time.Now().Add(-35 * time.Minute)},
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id2', 'bot2', 'HUMAN', 'game_id2')`,
				},
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled", "created_at", "public")
					VALUES ('game_id3', 'STARTED', 0, Array['b','p1','b','p2'], false, $1, false)`,
					Args: []any{time.Now().Add(-5 * time.Minute)},
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id3', 'bot3', 'HUMAN', 'game_id3')`,
				},
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled", "created_at", "public")
					VALUES ('game_id4', 'STARTED', 0, Array['b','p1','b','p2'], false, $1, true)`,
					Args: []any{time.Now().Add(-5 * time.Minute)},
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id4', 'bot4', 'AI', 'game_id4')`,
				},
				{
					Query: `INSERT INTO public."games" ("id", "state", "current_turn_index", "turn_order", "state_handled", "created_at", "public")
					VALUES ('game_id5', 'STARTED', 0, Array['b','p1','b','p2'], false, $1, true)`,
					Args: []any{time.Now().Add(-5 * time.Minute)},
				},
				{
					Query: `INSERT INTO public."bots" ("id", "name", "type", "game_id")
					VALUES ('bot_id5', 'bot5', 'HUMAN', 'game_id5')`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'game_id1'`},
				{Query: `DELETE FROM public."games" WHERE id = 'game_id2'`},
				{Query: `DELETE FROM public."games" WHERE id = 'game_id3'`},
				{Query: `DELETE FROM public."games" WHERE id = 'game_id4'`},
				{Query: `DELETE FROM public."games" WHERE id = 'game_id5'`},
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
			games, err := s.GetPublicJoinableGames()
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
