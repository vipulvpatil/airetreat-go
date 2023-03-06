package storage

import (
	"database/sql"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Game_DeleteGame(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		dbUpdateCheck   func(*sql.DB) bool
		setupSqlStmts   []string
		cleanupSqlStmts []string
		errorExpected   bool
		errorString     string
	}{
		{
			name:          "errors when game_id is blank",
			input:         "",
			dbUpdateCheck: nil,
			setupSqlStmts: []string{},
			errorExpected: true,
			errorString:   "gameId cannot be blank",
		},
		{
			name:          "errors when deleting a game that is not in db",
			input:         "game_id1",
			dbUpdateCheck: nil,
			setupSqlStmts: nil,
			errorExpected: true,
			errorString:   "THIS IS BAD: Very few or too many rows were affected when deleting game in db. This is highly unexpected. rowsAffected: 0",
		},
		{
			name:  "successfully deletes a game",
			input: "game_id1",
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					gameId string
					state  string
				)
				row := db.QueryRow(
					`SELECT g.id, g.state
					FROM public."games" AS g
					WHERE g.id = 'game_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&gameId, &state)
				assert.EqualError(t, err, "sql: no rows in result set")

				return true
			},
			setupSqlStmts: []string{
				`INSERT INTO public."games" (
					"id", "state", "current_turn_index", "turn_order", "state_handled"
				)
				VALUES (
					'game_id1', 'PLAYERS_JOINED', 0, Array['b','p1','b','p2'], false
				)`,
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id1', 'bot1', 'AI', 'game_id1'
				)`,
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id2', 'bot2', 'AI', 'game_id1'
				)`,
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
			err := s.DeleteGame(tt.input)
			if !tt.errorExpected {
				assert.NoError(t, err)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
			if tt.dbUpdateCheck != nil {
				assert.True(t, tt.dbUpdateCheck(s.db))
			}
		})
	}
}
