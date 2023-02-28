package storage

import (
	"database/sql"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Game_JoinGame(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			gameId   string
			playerId string
		}
		dbUpdateCheck   func(*sql.DB) bool
		setupSqlStmts   []string
		cleanupSqlStmts []string
		errorExpected   bool
		errorString     string
	}{
		{
			name: "errors when game_id is blank",
			input: struct {
				gameId   string
				playerId string
			}{
				gameId:   "",
				playerId: "player_id2",
			},
			dbUpdateCheck:   nil,
			setupSqlStmts:   []string{},
			cleanupSqlStmts: []string{},
			errorExpected:   true,
			errorString:     "gameId cannot be blank",
		},
		{
			name: "errors when player_id is blank",
			input: struct {
				gameId   string
				playerId string
			}{
				gameId:   "game_id1",
				playerId: "",
			},
			dbUpdateCheck:   nil,
			setupSqlStmts:   []string{},
			cleanupSqlStmts: []string{},
			errorExpected:   true,
			errorString:     "playerId cannot be blank",
		},
		{
			name: "errors when trying to join a game that has no AI bots",
			input: struct {
				gameId   string
				playerId string
			}{
				gameId:   "game_id1",
				playerId: "player_id2",
			},
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
				assert.NoError(t, err)
				assert.Equal(t, "game_id1", gameId)
				assert.Equal(t, "STARTED", state)

				return true
			},
			setupSqlStmts: []string{
				`INSERT INTO public."games" (
					"id", "state", "current_turn_index", "turn_order", "state_handled"
				)
				VALUES (
					'game_id1', 'STARTED', 0, Array['b','p1','b','p2'], false
				)`,
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id1', 'bot1', 'AI', 'game_id1'
				)`,
				`INSERT INTO public."players" ("id") VALUES ('player_id1')`,
				`UPDATE public."bots" SET
				"player_id" = 'player_id1',
				"type" = 'HUMAN'
				WHERE id = 'bot_id1'`,
				`INSERT INTO public."players" ("id") VALUES ('player_id2')`,
			},
			cleanupSqlStmts: []string{
				`DELETE FROM public."games" WHERE id = 'game_id1'`,
				`DELETE FROM public."players" WHERE id = 'player_id1'`,
				`DELETE FROM public."players" WHERE id = 'player_id2'`,
			},
			errorExpected: true,
			errorString:   "no AI bots in the game",
		},
		{
			name: "errors when trying to join a game that has already started (state: playersJoined)",
			input: struct {
				gameId   string
				playerId string
			}{
				gameId:   "game_id1",
				playerId: "player_id1",
			},
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
				assert.NoError(t, err)
				assert.Equal(t, "game_id1", gameId)
				assert.Equal(t, "PLAYERS_JOINED", state)

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
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id3', 'bot3', 'AI', 'game_id1'
				)`,
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id4', 'bot4', 'AI', 'game_id1'
				)`,
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id5', 'bot5', 'AI', 'game_id1'
				)`,
			},
			cleanupSqlStmts: []string{
				`DELETE FROM public."games" WHERE id = 'game_id1'`,
			},
			errorExpected: true,
			errorString:   "Cannot join this game: game_id1",
		},
		{
			name: "does not error or update anynthing if player is already in the game",
			input: struct {
				gameId   string
				playerId string
			}{
				gameId:   "game_id1",
				playerId: "player_id1",
			},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					gameId    string
					state     string
					botId     string
					playerId  string
					typeOfBot string
				)
				row := db.QueryRow(
					`SELECT g.id, g.state, b.id, b.player_id, b.type
					FROM public."games" AS g
					JOIN public."bots" AS b ON b.game_id = g.id
					JOIN public."players" AS p ON p.id = b.player_id
					WHERE g.id = 'game_id1' AND p.id = 'player_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&gameId, &state, &botId, &playerId, &typeOfBot)
				assert.NoError(t, err)
				assert.Equal(t, "game_id1", gameId)
				assert.Equal(t, "STARTED", state)
				assert.Equal(t, "bot_id5", botId)
				assert.Equal(t, "player_id1", playerId)
				assert.Equal(t, "HUMAN", typeOfBot)

				return true
			},
			setupSqlStmts: []string{
				`INSERT INTO public."games" (
					"id", "state", "current_turn_index", "turn_order", "state_handled"
				)
				VALUES (
					'game_id1', 'STARTED', 0, Array['b','p1','b','p2'], false
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
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id3', 'bot3', 'AI', 'game_id1'
				)`,
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id4', 'bot4', 'AI', 'game_id1'
				)`,
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id5', 'bot5', 'AI', 'game_id1'
				)`,
				`INSERT INTO public."players" ("id") VALUES ('player_id1')`,
				`UPDATE public."bots" SET
				"player_id" = 'player_id1',
				"type" = 'HUMAN'
				WHERE id = 'bot_id5'`,
			},
			cleanupSqlStmts: []string{
				`DELETE FROM public."games" WHERE id = 'game_id1'`,
				`DELETE FROM public."players" WHERE id = 'player_id1'`,
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "player joins a game successfully and but does not update state, if there is one human player in the game after joining",
			input: struct {
				gameId   string
				playerId string
			}{
				gameId:   "game_id1",
				playerId: "player_id1",
			},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					gameId    string
					state     string
					botId     string
					playerId  string
					typeOfBot string
				)
				row := db.QueryRow(
					`SELECT g.id, g.state, b.id, b.player_id, b.type
					FROM public."games" AS g
					JOIN public."bots" AS b ON b.game_id = g.id
					JOIN public."players" AS p ON p.id = b.player_id
					WHERE g.id = 'game_id1' AND p.id = 'player_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&gameId, &state, &botId, &playerId, &typeOfBot)
				assert.NoError(t, err)
				assert.Equal(t, "game_id1", gameId)
				assert.Equal(t, "STARTED", state)
				assert.Equal(t, "bot_id3", botId)
				assert.Equal(t, "player_id1", playerId)
				assert.Equal(t, "HUMAN", typeOfBot)

				return true
			},
			setupSqlStmts: []string{
				`INSERT INTO public."games" (
					"id", "state", "current_turn_index", "turn_order", "state_handled"
				)
				VALUES (
					'game_id1', 'STARTED', 0, Array['b','p1','b','p2'], false
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
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id3', 'bot3', 'AI', 'game_id1'
				)`,
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id4', 'bot4', 'AI', 'game_id1'
				)`,
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id5', 'bot5', 'AI', 'game_id1'
				)`,
				`INSERT INTO public."players" ("id") VALUES ('player_id1')`,
			},
			cleanupSqlStmts: []string{
				`DELETE FROM public."games" WHERE id = 'game_id1'`,
				`DELETE FROM public."players" WHERE id = 'player_id1'`,
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "player joins a game successfully and updates state, if there are two human players in the game after joining",
			input: struct {
				gameId   string
				playerId string
			}{
				gameId:   "game_id1",
				playerId: "player_id2",
			},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					gameId    string
					state     string
					botId     string
					playerId  string
					typeOfBot string
				)
				row := db.QueryRow(
					`SELECT g.id, g.state, b.id, b.player_id, b.type
					FROM public."games" AS g
					JOIN public."bots" AS b ON b.game_id = g.id
					JOIN public."players" AS p ON p.id = b.player_id
					WHERE g.id = 'game_id1' AND p.id = 'player_id2'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&gameId, &state, &botId, &playerId, &typeOfBot)
				assert.NoError(t, err)
				assert.Equal(t, "game_id1", gameId)
				assert.Equal(t, "PLAYERS_JOINED", state)
				assert.Equal(t, "bot_id3", botId)
				assert.Equal(t, "player_id2", playerId)
				assert.Equal(t, "HUMAN", typeOfBot)

				return true
			},
			setupSqlStmts: []string{
				`INSERT INTO public."games" (
					"id", "state", "current_turn_index", "turn_order", "state_handled"
				)
				VALUES (
					'game_id1', 'STARTED', 0, Array['b','p1','b','p2'], false
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
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id3', 'bot3', 'AI', 'game_id1'
				)`,
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id4', 'bot4', 'AI', 'game_id1'
				)`,
				`INSERT INTO public."bots" (
					"id", "name", "type", "game_id"
				)
				VALUES (
					'bot_id5', 'bot5', 'AI', 'game_id1'
				)`,
				`INSERT INTO public."players" ("id") VALUES ('player_id1')`,
				`UPDATE public."bots" SET
				"player_id" = 'player_id1',
				"type" = 'HUMAN'
				WHERE id = 'bot_id5'`,
				`INSERT INTO public."players" ("id") VALUES ('player_id2')`,
			},
			cleanupSqlStmts: []string{
				`DELETE FROM public."games" WHERE id = 'game_id1'`,
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
			err := s.JoinGame(tt.input.gameId, tt.input.playerId)
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
