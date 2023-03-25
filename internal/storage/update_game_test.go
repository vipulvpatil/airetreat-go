package storage

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/model"
)

func Test_Game_UpdateGameState(t *testing.T) {
	state := "PLAYERS_JOINED"
	currentTurnIndex := int64(4)
	turnOrder := []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"}
	stateHandled := true
	stateHandledAt := time.Now()
	lastQuestion := "what is the question?"
	lastQuestionTargetBotId := "bot_id2"
	stateTotalTime := int64(60)
	tests := []struct {
		name  string
		input struct {
			gameId     string
			updateOpts GameUpdateOptions
		}
		dbUpdateCheck   func(*sql.DB) bool
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		errorExpected   bool
		errorString     string
	}{
		{
			name: "updates game with all given details",
			input: struct {
				gameId     string
				updateOpts GameUpdateOptions
			}{
				gameId: "game_id1",
				updateOpts: GameUpdateOptions{
					State:                   &state,
					CurrentTurnIndex:        &currentTurnIndex,
					TurnOrder:               turnOrder,
					StateHandled:            &stateHandled,
					StateHandledAt:          &stateHandledAt,
					LastQuestion:            &lastQuestion,
					LastQuestionTargetBotId: &lastQuestionTargetBotId,
					StateTotalTime:          &stateTotalTime,
				},
			},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					scanState                   string
					scanCurrentTurnIndex        int64
					scanTurnOrder               []string
					scanStateHandled            bool
					scanStateHandledAt          sql.NullTime
					scanLastQuestion            sql.NullString
					scanLastQuestionTargetBotId sql.NullString
					scanStateTotalTime          int64
					updatedAt                   time.Time
				)
				row := db.QueryRow(
					`SELECT g.state, g.current_turn_index, g.turn_order, g.state_handled, g.state_handled_at,
					g.last_question, g.last_question_target_bot_id, g.state_total_time, g.updated_at
					FROM public."games" AS g
					WHERE g.id = 'game_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&scanState, &scanCurrentTurnIndex, pq.Array(&scanTurnOrder), &scanStateHandled, &scanStateHandledAt, &scanLastQuestion, &scanLastQuestionTargetBotId, &scanStateTotalTime, &updatedAt)
				assert.NoError(t, err)
				assert.Equal(t, state, scanState)
				assert.Equal(t, currentTurnIndex, scanCurrentTurnIndex)
				assert.Equal(t, turnOrder, scanTurnOrder)
				assert.Equal(t, stateHandled, scanStateHandled)
				assert.True(t, scanStateHandledAt.Valid)
				model.AssertTimeAlmostEqual(t, scanStateHandledAt.Time, stateHandledAt, 1*time.Second)
				assert.True(t, scanLastQuestion.Valid)
				assert.Equal(t, lastQuestion, scanLastQuestion.String)
				assert.True(t, scanLastQuestionTargetBotId.Valid)
				assert.Equal(t, lastQuestionTargetBotId, scanLastQuestionTargetBotId.String)
				assert.Equal(t, stateTotalTime, scanStateTotalTime)
				model.AssertTimeAlmostEqual(t, time.Now(), updatedAt, 1*time.Second)
				return true
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" (
						"id", "state", "current_turn_index", "turn_order", "state_handled", "created_at", "updated_at"
					)
					VALUES (
						'game_id1', 'STARTED', 0, Array['bot_id1'], false, $1, $2
					)`,
					Args: []any{
						time.Now().Add(-1 * time.Hour),
						time.Now().Add(-1 * time.Hour),
					},
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id2', 'bot2', 'AI', 'game_id1'
					)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'game_id1'`},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "updates game partially with given details",
			input: struct {
				gameId     string
				updateOpts GameUpdateOptions
			}{
				gameId: "game_id1",
				updateOpts: GameUpdateOptions{
					State:                   &state,
					StateHandled:            &stateHandled,
					StateHandledAt:          &stateHandledAt,
					LastQuestion:            &lastQuestion,
					LastQuestionTargetBotId: &lastQuestionTargetBotId,
					StateTotalTime:          &stateTotalTime,
				},
			},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					scanState                   string
					scanCurrentTurnIndex        int64
					scanTurnOrder               []string
					scanStateHandled            bool
					scanStateHandledAt          sql.NullTime
					scanLastQuestion            sql.NullString
					scanLastQuestionTargetBotId sql.NullString
					scanStateTotalTime          int64
					updatedAt                   time.Time
				)
				unchangedCurrentTurnIndex := int64(0)
				unchangedTurnOrder := []string{"bot_id1"}
				row := db.QueryRow(
					`SELECT g.state, g.current_turn_index, g.turn_order, g.state_handled, g.state_handled_at,
					g.last_question, g.last_question_target_bot_id, g.state_total_time, g.updated_at
					FROM public."games" AS g
					WHERE g.id = 'game_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&scanState, &scanCurrentTurnIndex, pq.Array(&scanTurnOrder), &scanStateHandled, &scanStateHandledAt, &scanLastQuestion, &scanLastQuestionTargetBotId, &scanStateTotalTime, &updatedAt)
				assert.NoError(t, err)
				assert.Equal(t, state, scanState)
				assert.Equal(t, unchangedCurrentTurnIndex, scanCurrentTurnIndex)
				assert.Equal(t, unchangedTurnOrder, scanTurnOrder)
				assert.Equal(t, stateHandled, scanStateHandled)
				assert.True(t, scanStateHandledAt.Valid)
				model.AssertTimeAlmostEqual(t, scanStateHandledAt.Time, stateHandledAt, 1*time.Second)
				assert.True(t, scanLastQuestion.Valid)
				assert.Equal(t, lastQuestion, scanLastQuestion.String)
				assert.True(t, scanLastQuestionTargetBotId.Valid)
				assert.Equal(t, lastQuestionTargetBotId, scanLastQuestionTargetBotId.String)
				assert.Equal(t, stateTotalTime, scanStateTotalTime)
				model.AssertTimeAlmostEqual(t, time.Now(), updatedAt, 1*time.Second)

				return true
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" (
						"id", "state", "current_turn_index", "turn_order", "state_handled", "created_at", "updated_at"
					)
					VALUES (
						'game_id1', 'STARTED', 0, Array['bot_id1'], false, $1, $2
					)`,
					Args: []any{
						time.Now().Add(-1 * time.Hour),
						time.Now().Add(-1 * time.Hour),
					},
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id2', 'bot2', 'AI', 'game_id1'
					)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'game_id1'`},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "errors if no update options provided",
			input: struct {
				gameId     string
				updateOpts GameUpdateOptions
			}{
				gameId:     "game_id1",
				updateOpts: GameUpdateOptions{},
			},
			dbUpdateCheck:   nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "no update options provided",
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
			err := s.UpdateGameState(tt.input.gameId, tt.input.updateOpts)
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

func Test_Game_UpdateGameStateUsingTrasaction(t *testing.T) {
	state := "PLAYERS_JOINED"
	currentTurnIndex := int64(4)
	turnOrder := []string{"bot_id1", "bot_id2", "bot_id3", "bot_id4", "bot_id5"}
	stateHandled := true
	stateHandledAt := time.Now()
	lastQuestion := "what is the question?"
	lastQuestionTargetBotId := "bot_id2"
	stateTotalTime := int64(60)
	tests := []struct {
		name  string
		input struct {
			gameId     string
			updateOpts GameUpdateOptions
		}
		dbUpdateCheck   func(*sql.DB) bool
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		errorExpected   bool
		errorString     string
	}{
		{
			name: "updates game with all given details",
			input: struct {
				gameId     string
				updateOpts GameUpdateOptions
			}{
				gameId: "game_id1",
				updateOpts: GameUpdateOptions{
					State:                   &state,
					CurrentTurnIndex:        &currentTurnIndex,
					TurnOrder:               turnOrder,
					StateHandled:            &stateHandled,
					StateHandledAt:          &stateHandledAt,
					LastQuestion:            &lastQuestion,
					LastQuestionTargetBotId: &lastQuestionTargetBotId,
					StateTotalTime:          &stateTotalTime,
				},
			},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					scanState                   string
					scanCurrentTurnIndex        int64
					scanTurnOrder               []string
					scanStateHandled            bool
					scanStateHandledAt          sql.NullTime
					scanLastQuestion            sql.NullString
					scanLastQuestionTargetBotId sql.NullString
					scanStateTotalTime          int64
					updatedAt                   time.Time
				)
				row := db.QueryRow(
					`SELECT g.state, g.current_turn_index, g.turn_order, g.state_handled, g.state_handled_at,
					g.last_question, g.last_question_target_bot_id, g.state_total_time, g.updated_at
					FROM public."games" AS g
					WHERE g.id = 'game_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&scanState, &scanCurrentTurnIndex, pq.Array(&scanTurnOrder), &scanStateHandled, &scanStateHandledAt, &scanLastQuestion, &scanLastQuestionTargetBotId, &scanStateTotalTime, &updatedAt)
				assert.NoError(t, err)
				assert.Equal(t, state, scanState)
				assert.Equal(t, currentTurnIndex, scanCurrentTurnIndex)
				assert.Equal(t, turnOrder, scanTurnOrder)
				assert.Equal(t, stateHandled, scanStateHandled)
				assert.True(t, scanStateHandledAt.Valid)
				model.AssertTimeAlmostEqual(t, scanStateHandledAt.Time, stateHandledAt, 1*time.Second)
				assert.True(t, scanLastQuestion.Valid)
				assert.Equal(t, lastQuestion, scanLastQuestion.String)
				assert.True(t, scanLastQuestionTargetBotId.Valid)
				assert.Equal(t, lastQuestionTargetBotId, scanLastQuestionTargetBotId.String)
				assert.Equal(t, stateTotalTime, scanStateTotalTime)
				model.AssertTimeAlmostEqual(t, time.Now(), updatedAt, 1*time.Second)
				return true
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" (
						"id", "state", "current_turn_index", "turn_order", "state_handled", "created_at", "updated_at"
					)
					VALUES (
						'game_id1', 'STARTED', 0, Array['bot_id1'], false, $1, $2
					)`,
					Args: []any{
						time.Now().Add(-1 * time.Hour),
						time.Now().Add(-1 * time.Hour),
					},
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id2', 'bot2', 'AI', 'game_id1'
					)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'game_id1'`},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name: "updates game partially with given details",
			input: struct {
				gameId     string
				updateOpts GameUpdateOptions
			}{
				gameId: "game_id1",
				updateOpts: GameUpdateOptions{
					State:                   &state,
					StateHandled:            &stateHandled,
					StateHandledAt:          &stateHandledAt,
					LastQuestion:            &lastQuestion,
					LastQuestionTargetBotId: &lastQuestionTargetBotId,
					StateTotalTime:          &stateTotalTime,
				},
			},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					scanState                   string
					scanCurrentTurnIndex        int64
					scanTurnOrder               []string
					scanStateHandled            bool
					scanStateHandledAt          sql.NullTime
					scanLastQuestion            sql.NullString
					scanLastQuestionTargetBotId sql.NullString
					scanStateTotalTime          int64
					updatedAt                   time.Time
				)
				unchangedCurrentTurnIndex := int64(0)
				unchangedTurnOrder := []string{"bot_id1"}
				row := db.QueryRow(
					`SELECT g.state, g.current_turn_index, g.turn_order, g.state_handled, g.state_handled_at,
					g.last_question, g.last_question_target_bot_id, g.state_total_time, g.updated_at
					FROM public."games" AS g
					WHERE g.id = 'game_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&scanState, &scanCurrentTurnIndex, pq.Array(&scanTurnOrder), &scanStateHandled, &scanStateHandledAt, &scanLastQuestion, &scanLastQuestionTargetBotId, &scanStateTotalTime, &updatedAt)
				assert.NoError(t, err)
				assert.Equal(t, state, scanState)
				assert.Equal(t, unchangedCurrentTurnIndex, scanCurrentTurnIndex)
				assert.Equal(t, unchangedTurnOrder, scanTurnOrder)
				assert.Equal(t, stateHandled, scanStateHandled)
				assert.True(t, scanStateHandledAt.Valid)
				model.AssertTimeAlmostEqual(t, scanStateHandledAt.Time, stateHandledAt, 1*time.Second)
				assert.True(t, scanLastQuestion.Valid)
				assert.Equal(t, lastQuestion, scanLastQuestion.String)
				assert.True(t, scanLastQuestionTargetBotId.Valid)
				assert.Equal(t, lastQuestionTargetBotId, scanLastQuestionTargetBotId.String)
				assert.Equal(t, stateTotalTime, scanStateTotalTime)
				model.AssertTimeAlmostEqual(t, time.Now(), updatedAt, 1*time.Second)

				return true
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" (
						"id", "state", "current_turn_index", "turn_order", "state_handled", "created_at", "updated_at"
					)
					VALUES (
						'game_id1', 'STARTED', 0, Array['bot_id1'], false, $1, $2
					)`,
					Args: []any{
						time.Now().Add(-1 * time.Hour),
						time.Now().Add(-1 * time.Hour),
					},
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id2', 'bot2', 'AI', 'game_id1'
					)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'game_id1'`},
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
			tx, err := s.BeginTransaction()
			assert.NoError(t, err)
			err = s.UpdateGameStateUsingTransaction(tt.input.gameId, tt.input.updateOpts, tx)
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

func Test_Game_UpdateGameStateIfEnoughPlayersHaveJoinedUsingTransaction(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		dbUpdateCheck   func(*sql.DB) bool
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		errorExpected   bool
		errorString     string
	}{
		{
			name:  "updates game if enough players have joined",
			input: "game_id1",
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					scanState string
					updatedAt time.Time
				)
				row := db.QueryRow(
					`SELECT g.state, g.updated_at
					FROM public."games" AS g
					WHERE g.id = 'game_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&scanState, &updatedAt)
				assert.NoError(t, err)
				assert.Equal(t, "PLAYERS_JOINED", scanState)
				model.AssertTimeAlmostEqual(t, time.Now(), updatedAt, 1*time.Second)
				return true
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" (
						"id", "state", "current_turn_index", "turn_order", "state_handled", "updated_at"
					)
					VALUES (
						'game_id1', 'STARTED', 0, Array['bot_id1'], false, $1
					)`,
					Args: []any{
						time.Now().Add(-1 * time.Hour),
					},
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id1', 'bot1', 'HUMAN', 'game_id1'
					)`,
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id2', 'bot2', 'HUMAN', 'game_id1'
					)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'game_id1'`},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name:  "does not update game if enough players have not joined",
			input: "game_id1",
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					scanState string
					updatedAt time.Time
				)
				row := db.QueryRow(
					`SELECT g.state, g.updated_at
					FROM public."games" AS g
					WHERE g.id = 'game_id1'`,
				)
				assert.NoError(t, row.Err())
				err := row.Scan(&scanState, &updatedAt)
				assert.NoError(t, err)
				assert.Equal(t, "STARTED", scanState)
				model.AssertTimeAlmostEqual(t, time.Now().Add(-1*time.Hour), updatedAt, 1*time.Second)
				return true
			},
			setupSqlStmts: []TestSqlStmts{
				{
					Query: `INSERT INTO public."games" (
						"id", "state", "current_turn_index", "turn_order", "state_handled", "updated_at"
					)
					VALUES (
						'game_id1', 'STARTED', 0, Array['bot_id1'], false, $1
					)`,
					Args: []any{
						time.Now().Add(-1 * time.Hour),
					},
				},
				{
					Query: `INSERT INTO public."bots" (
						"id", "name", "type", "game_id"
					)
					VALUES (
						'bot_id1', 'bot1', 'HUMAN', 'game_id1'
					)`,
				},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'game_id1'`},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name:            "errors if gameId is blank",
			input:           "",
			dbUpdateCheck:   nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "gameId cannot be blank",
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

			tx, err := s.BeginTransaction()
			assert.NoError(t, err)

			err = s.UpdateGameStateIfEnoughPlayersHaveJoinedUsingTransaction(tt.input, tx)
			err = tx.Commit()
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
