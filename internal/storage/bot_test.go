package storage

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UpdateBotWithPlayerIdUsingTransaction(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			botId    string
			playerId string
		}
		dbUpdateCheck   func(*sql.DB) bool
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		errorExpected   bool
		errorString     string
	}{
		{
			name: "errors if botId is blank",
			input: struct {
				botId    string
				playerId string
			}{
				botId:    "",
				playerId: "",
			},
			dbUpdateCheck:   nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "botId cannot be blank",
		},
		{
			name: "errors if playerId is blank",
			input: struct {
				botId    string
				playerId string
			}{
				botId:    "bot_id1",
				playerId: "",
			},
			dbUpdateCheck:   nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "playerId cannot be blank",
		},
		{
			name: "errors if player id is not in db",
			input: struct {
				botId    string
				playerId string
			}{
				botId:    "bot_id1",
				playerId: "player_id1",
			},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					playerId  sql.NullString
					typeOfBot string
				)
				row := db.QueryRow(
					`SELECT player_id, type
				FROM public."bots"
				WHERE id = 'bot_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&playerId, &typeOfBot)
				assert.NoError(t, err)
				assert.False(t, playerId.Valid)
				assert.Equal(t, "AI", typeOfBot)

				return true
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" (
						"id", "state", "current_turn_index", "turn_order", "state_handled"
					)
					VALUES (
						'bot_id1', 'STARTED', 0, Array['b','p1','b','p2'], false
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id1', 'bot1', 'AI', 'bot_id1'
					)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'bot_id1'`},
			},
			errorExpected: true,
			errorString:   "THIS IS BAD: dbError while connecting player to bot: player_id1 bot_id1: pq: insert or update on table \"bots\" violates foreign key constraint \"bots_player_id_fkey\"",
		},
		{
			name: "bot updates successfully with the playerId",
			input: struct {
				botId    string
				playerId string
			}{
				botId:    "bot_id1",
				playerId: "player_id1",
			},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					playerId  string
					typeOfBot string
				)
				row := db.QueryRow(
					`SELECT player_id, type
				FROM public."bots"
				WHERE id = 'bot_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&playerId, &typeOfBot)
				assert.NoError(t, err)
				assert.Equal(t, "player_id1", playerId)
				assert.Equal(t, "HUMAN", typeOfBot)

				return true
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" (
						"id", "state", "current_turn_index", "turn_order", "state_handled"
					)
					VALUES (
						'bot_id1', 'STARTED', 0, Array['b','p1','b','p2'], false
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id1', 'bot1', 'AI', 'bot_id1'
					)`,
				},
				{Query: `INSERT INTO public."players" ("id") VALUES ('player_id1')`},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'bot_id1'`},
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
			err = s.UpdateBotWithPlayerIdUsingTransaction(tt.input.botId, tt.input.playerId, tx)
			tx.Commit()
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
