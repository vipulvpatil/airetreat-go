package storage

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

func Test_CreateGame(t *testing.T) {
	tests := []struct {
		name            string
		setupSqlStmts   []string
		cleanupSqlStmts []string
		idGenerator     utilities.CuidGenerator
		dbUpdateCheck   func(*sql.DB) bool
		errorExpected   bool
		errorString     string
	}{
		{
			name:          "creates game successfully",
			setupSqlStmts: []string{},
			cleanupSqlStmts: []string{
				`DELETE FROM public."games" WHERE id = 'game_id1'`,
			},
			idGenerator: &utilities.IdGeneratorMockSeries{Series: []string{"game_id1", "bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"}},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					id                      string
					state                   string
					currentTurnIndex        int
					turnOrder               []string
					stateHandled            bool
					stateHandledAt          pq.NullTime
					stateTotalTime          int
					lastQuestion            sql.NullString
					lastQuestionTargetBotId sql.NullString
					createdAt               pq.NullTime
					updatedAt               pq.NullTime
				)
				err := db.QueryRow(
					`SELECT "id", "state", "current_turn_index", "turn_order", "state_handled", "state_handled_at", "state_total_time", "last_question", "last_question_target_bot_id", "created_at", "updated_at"
					FROM public."games" WHERE "id" = 'game_id1'`,
				).Scan(&id, &state, &currentTurnIndex, pq.Array(&turnOrder), &stateHandled, &stateHandledAt, &stateTotalTime, &lastQuestion, &lastQuestionTargetBotId, &createdAt, &updatedAt)
				assert.NoError(t, err)
				assert.Equal(t, "game_id1", id)
				assert.Equal(t, "STARTED", state)
				assert.Equal(t, 0, currentTurnIndex)
				assert.EqualValues(t, []string{"b", "p1", "b", "p2"}, turnOrder)
				assert.False(t, stateHandled)
				assert.False(t, stateHandledAt.Valid)
				assert.Equal(t, 0, stateTotalTime)
				assert.False(t, lastQuestion.Valid)
				assert.False(t, lastQuestionTargetBotId.Valid)
				assert.True(t, createdAt.Valid)
				assert.True(t, updatedAt.Valid)

				rows, err := db.Query(
					`SELECT
					"id", "name", "type", "player_id", "question_count", "created_at"
					FROM public."bots" WHERE "game_id" = 'game_id1'`,
				)
				defer rows.Close()
				assert.NoError(t, err)

				botIndex := 0
				expectedBotNames := []string{"Electronic Device-209", "ED-I", "B.O.B.Z", "T-800X", "GLaDOODLES"}

				for rows.Next() {
					var (
						id            string
						name          string
						typeOfBot     string
						playerId      sql.NullString
						questionCount int
						createdAt     pq.NullTime
					)
					err = rows.Scan(&id, &name, &typeOfBot, &playerId, &questionCount, &createdAt)
					assert.NoError(t, err)
					assert.Equal(t, fmt.Sprintf("bot_id%d", botIndex+1), id)
					assert.Equal(t, expectedBotNames[botIndex], name)
					assert.Equal(t, "AI", typeOfBot)
					assert.False(t, playerId.Valid)
					assert.Equal(t, 0, questionCount)
					assert.NotNil(t, createdAt)
					botIndex++
				}
				assert.Equal(t, 5, botIndex)
				return true
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "errors and does not update anything, if Game ID already exists in DB",
			setupSqlStmts: []string{
				`INSERT INTO public."games" (
					"id", "state", "current_turn_index", "turn_order", "state_handled"
				)
				VALUES (
					'id1', 'STARTED', 0, ARRAY['b', 'p1', 'b', 'p2'], false
				)
				`,
			},
			cleanupSqlStmts: []string{
				`DELETE FROM public."games" WHERE id = 'id1'`,
			},
			idGenerator: &utilities.IdGeneratorMockConstant{Id: "id1"},
			dbUpdateCheck: func(db *sql.DB) bool {
				var id string
				err := db.QueryRow(
					`SELECT id FROM public."games" WHERE "id" = 'id1'`,
				).Scan(&id)
				assert.NoError(t, err)
				assert.Equal(t, "id1", id)

				rows, err := db.Query(
					`SELECT
					"id", "name", "type", "player_id", "question_count", "created_at"
					FROM public."bots" WHERE "game_id" = 'id1'`,
				)
				defer rows.Close()
				assert.NoError(t, err)
				assert.False(t, rows.Next())
				return true
			},
			errorExpected: true,
			errorString:   "pq: duplicate key value violates unique constraint \"games_pkey\"",
		},
		{
			name:            "errors and does not update anything, if Bot ID already exists in DB",
			setupSqlStmts:   []string{},
			cleanupSqlStmts: []string{},
			idGenerator:     &utilities.IdGeneratorMockConstant{Id: "id1"},
			dbUpdateCheck: func(db *sql.DB) bool {
				var id string
				err := db.QueryRow(
					`SELECT id FROM public."games" WHERE "id" = 'id1'`,
				).Scan(&id)
				assert.Error(t, err, "some error")

				rows, err := db.Query(
					`SELECT
					"id", "name", "type", "player_id", "question_count", "created_at"
					FROM public."bots" WHERE "game_id" = 'id1'`,
				)
				defer rows.Close()
				assert.NoError(t, err)
				assert.False(t, rows.Next())
				return true
			},
			errorExpected: true,
			errorString:   "pq: duplicate key value violates unique constraint \"bots_pkey\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := NewDbStorage(
				StorageOptions{
					Db:          testDb,
					IdGenerator: tt.idGenerator,
				},
			)

			runSqlOnDb(t, s.db, tt.setupSqlStmts)
			defer runSqlOnDb(t, s.db, tt.cleanupSqlStmts)

			rand.Seed(0)
			err := s.CreateGame()
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
