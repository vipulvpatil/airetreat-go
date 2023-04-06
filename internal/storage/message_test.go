package storage

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

func Test_CreateMessage(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			sourceBotId string
			targetBotId string
			text        string
			messageType string
		}
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		idGenerator     utilities.CuidGenerator
		dbUpdateCheck   func(*sql.DB) bool
		errorExpected   bool
		errorString     string
	}{
		{
			name: "errors when sourceBotId is blank",
			input: struct {
				sourceBotId string
				targetBotId string
				text        string
				messageType string
			}{
				"",
				"",
				"this is a message",
				"question",
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			idGenerator:     &utilities.IdGeneratorMockConstant{Id: "message_id1"},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					id          string
					sourceBotId string
					targetBotId string
					text        string
					createdAt   time.Time
				)
				err := db.QueryRow(
					`SELECT "id", "source_bot_id", "target_bot_id", "text", "created_at"
						FROM public."messages" WHERE "id" = 'message_id1'`,
				).Scan(&id, &sourceBotId, &targetBotId, &text, &createdAt)
				assert.EqualError(t, err, "sql: no rows in result set")
				return true
			},
			errorExpected: true,
			errorString:   "sourceBotId cannot be blank",
		},
		{
			name: "errors when targetBotId is blank",
			input: struct {
				sourceBotId string
				targetBotId string
				text        string
				messageType string
			}{
				"bot_id1",
				"",
				"this is a message",
				"question",
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			idGenerator:     &utilities.IdGeneratorMockConstant{Id: "message_id1"},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					id          string
					sourceBotId string
					targetBotId string
					text        string
					createdAt   time.Time
				)
				err := db.QueryRow(
					`SELECT "id", "source_bot_id", "target_bot_id", "text", "created_at"
						FROM public."messages" WHERE "id" = 'message_id1'`,
				).Scan(&id, &sourceBotId, &targetBotId, &text, &createdAt)
				assert.EqualError(t, err, "sql: no rows in result set")
				return true
			},
			errorExpected: true,
			errorString:   "targetBotId cannot be blank",
		},
		{
			name: "errors when text is blank",
			input: struct {
				sourceBotId string
				targetBotId string
				text        string
				messageType string
			}{
				"bot_id1",
				"bot_id2",
				"",
				"question",
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			idGenerator:     &utilities.IdGeneratorMockConstant{Id: "message_id1"},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					id          string
					sourceBotId string
					targetBotId string
					text        string
					createdAt   time.Time
				)
				err := db.QueryRow(
					`SELECT "id", "source_bot_id", "target_bot_id", "text", "created_at"
						FROM public."messages" WHERE "id" = 'message_id1'`,
				).Scan(&id, &sourceBotId, &targetBotId, &text, &createdAt)
				assert.EqualError(t, err, "sql: no rows in result set")
				return true
			},
			errorExpected: true,
			errorString:   "text cannot be blank",
		},
		{
			name: "errors when type is invalid",
			input: struct {
				sourceBotId string
				targetBotId string
				text        string
				messageType string
			}{
				"bot_id1",
				"bot_id2",
				"message",
				"some_message",
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			idGenerator:     &utilities.IdGeneratorMockConstant{Id: "message_id1"},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					id          string
					sourceBotId string
					targetBotId string
					text        string
					createdAt   time.Time
				)
				err := db.QueryRow(
					`SELECT "id", "source_bot_id", "target_bot_id", "text", "created_at"
						FROM public."messages" WHERE "id" = 'message_id1'`,
				).Scan(&id, &sourceBotId, &targetBotId, &text, &createdAt)
				assert.EqualError(t, err, "sql: no rows in result set")
				return true
			},
			errorExpected: true,
			errorString:   "invalid messageType",
		},
		{
			name: "errors when source bot is same as target bot for question",
			input: struct {
				sourceBotId string
				targetBotId string
				text        string
				messageType string
			}{
				"bot_id1",
				"bot_id1",
				"some question",
				"question",
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			idGenerator:     &utilities.IdGeneratorMockConstant{Id: "message_id1"},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					id          string
					sourceBotId string
					targetBotId string
					text        string
					createdAt   time.Time
				)
				err := db.QueryRow(
					`SELECT "id", "source_bot_id", "target_bot_id", "text", "created_at"
						FROM public."messages" WHERE "id" = 'message_id1'`,
				).Scan(&id, &sourceBotId, &targetBotId, &text, &createdAt)
				assert.EqualError(t, err, "sql: no rows in result set")
				return true
			},
			errorExpected: true,
			errorString:   "question source and target bot cannot be same. bot_id1 bot_id1",
		},
		{
			name: "errors when source bot is different than target bot for answer",
			input: struct {
				sourceBotId string
				targetBotId string
				text        string
				messageType string
			}{
				"bot_id1",
				"bot_id2",
				"some answer",
				"answer",
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			idGenerator:     &utilities.IdGeneratorMockConstant{Id: "message_id1"},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					id          string
					sourceBotId string
					targetBotId string
					text        string
					createdAt   time.Time
				)
				err := db.QueryRow(
					`SELECT "id", "source_bot_id", "target_bot_id", "text", "created_at"
						FROM public."messages" WHERE "id" = 'message_id1'`,
				).Scan(&id, &sourceBotId, &targetBotId, &text, &createdAt)
				assert.EqualError(t, err, "sql: no rows in result set")
				return true
			},
			errorExpected: true,
			errorString:   "answer source and target bot should be same. bot_id1 bot_id2",
		},
		{
			name: "creates message successfully",
			input: struct {
				sourceBotId string
				targetBotId string
				text        string
				messageType string
			}{
				"bot_id1",
				"bot_id2",
				"this is a message",
				"question",
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
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'game_id1'`},
			},
			idGenerator: &utilities.IdGeneratorMockConstant{Id: "message_id1"},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					id          string
					sourceBotId string
					targetBotId string
					text        string
					createdAt   time.Time
				)
				err := db.QueryRow(
					`SELECT "id", "source_bot_id", "target_bot_id", "text", "created_at"
						FROM public."messages" WHERE "id" = 'message_id1'`,
				).Scan(&id, &sourceBotId, &targetBotId, &text, &createdAt)
				assert.NoError(t, err)
				assert.Equal(t, "message_id1", id)
				assert.Equal(t, "bot_id1", sourceBotId)
				assert.Equal(t, "bot_id2", targetBotId)
				assert.Equal(t, "this is a message", text)
				model.AssertTimeAlmostEqual(t, createdAt, time.Now(), 5*time.Second, "createdAt is not within expected range")
				return true
			},
			errorExpected: false,
			errorString:   "",
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
			err := s.CreateMessage(tt.input.sourceBotId, tt.input.targetBotId, tt.input.text, tt.input.messageType)
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
func Test_CreateMessageUsingTransaction(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			sourceBotId string
			targetBotId string
			text        string
			messageType string
		}
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		idGenerator     utilities.CuidGenerator
		dbUpdateCheck   func(*sql.DB) bool
		errorExpected   bool
		errorString     string
	}{
		{
			name: "creates message successfully",
			input: struct {
				sourceBotId string
				targetBotId string
				text        string
				messageType string
			}{
				"bot_id1",
				"bot_id2",
				"this is a message",
				"question",
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
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."games" WHERE id = 'game_id1'`},
			},
			idGenerator: &utilities.IdGeneratorMockConstant{Id: "message_id1"},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					id          string
					sourceBotId string
					targetBotId string
					text        string
					createdAt   time.Time
					messageType string
				)
				err := db.QueryRow(
					`SELECT "id", "source_bot_id", "target_bot_id", "text", "created_at", "type"
						FROM public."messages" WHERE "id" = 'message_id1'`,
				).Scan(&id, &sourceBotId, &targetBotId, &text, &createdAt, &messageType)
				assert.NoError(t, err)
				assert.Equal(t, "message_id1", id)
				assert.Equal(t, "bot_id1", sourceBotId)
				assert.Equal(t, "bot_id2", targetBotId)
				assert.Equal(t, "this is a message", text)
				assert.Equal(t, "question", messageType)
				model.AssertTimeAlmostEqual(t, createdAt, time.Now(), 5*time.Second, "createdAt is not within expected range")
				return true
			},
			errorExpected: false,
			errorString:   "",
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

			tx, err := s.BeginTransaction()
			assert.NoError(t, err)
			err = s.CreateMessageUsingTransaction(tt.input.sourceBotId, tt.input.targetBotId, tt.input.text, tt.input.messageType, tx)
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
