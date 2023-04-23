package storage

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

func Test_GetPlayerUsingTransaction(t *testing.T) {
	userId := "user_id1"
	playerWithoutUser, _ := model.NewPlayer(model.PlayerOptions{Id: "player_id1"})
	playerWithUser, _ := model.NewPlayer(model.PlayerOptions{Id: "player_id1", UserId: &userId})
	tests := []struct {
		name            string
		input           string
		output          *model.Player
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		errorExpected   bool
		errorString     string
	}{
		{
			name:            "errors if player id is nil",
			input:           "",
			output:          nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "playerId cannot be blank",
		},
		{
			name:            "errors if player not in db",
			input:           "player_id1",
			output:          nil,
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			errorExpected:   true,
			errorString:     "getting player for player_id1: no such player",
		},
		{
			name:   "gets player successfully without user id",
			input:  "player_id1",
			output: playerWithoutUser,
			setupSqlStmts: []TestSqlStmts{
				{Query: `INSERT INTO public."players" ("id", "user_id") VALUES ('player_id1', NULL)`},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."players" WHERE id = 'player_id1'`},
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name:   "gets player with user id successfully",
			input:  "player_id1",
			output: playerWithUser,
			setupSqlStmts: []TestSqlStmts{
				{Query: `INSERT INTO public."users" ("id") VALUES ('user_id1')`},
				{Query: `INSERT INTO public."players" ("id", "user_id") VALUES ('player_id1', 'user_id1')`},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."users" WHERE id = 'user_id1'`},
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

			tx, err := s.BeginTransaction()
			assert.NoError(t, err)
			playerId, err := s.GetPlayerUsingTransaction(tt.input, tx)
			tx.Commit()

			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.Equal(t, tt.output, playerId)
			} else {
				assert.NotEmpty(t, tt.errorString)
				assert.EqualError(t, err, tt.errorString)
			}
		})
	}
}

func Test_CreatePlayer(t *testing.T) {
	userId := "user_id1"
	tests := []struct {
		name            string
		input           *string
		output          string
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		idGenerator     utilities.CuidGenerator
		dbUpdateCheck   func(*sql.DB) bool
		errorExpected   bool
		errorString     string
	}{
		{
			name:          "creates player successfully without user id",
			input:         nil,
			output:        "player_id1",
			setupSqlStmts: nil,
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."players" WHERE id = 'player_id1'`},
			},
			idGenerator: &utilities.IdGeneratorMockConstant{Id: "player_id1"},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					id string
				)
				err := db.QueryRow(
					`SELECT "id" FROM public."players" WHERE "id" = 'player_id1'`,
				).Scan(&id)
				assert.NoError(t, err)
				assert.Equal(t, "player_id1", id)
				return true
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name:   "creates player successfully with user id",
			input:  &userId,
			output: "player_id1",
			setupSqlStmts: []TestSqlStmts{
				{Query: `INSERT INTO public."users" ("id") VALUES ('user_id1')`},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."users" WHERE id = 'user_id1'`},
			},
			idGenerator: &utilities.IdGeneratorMockConstant{Id: "player_id1"},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					id     string
					userId string
				)
				err := db.QueryRow(
					`SELECT "id", "user_id" FROM public."players" WHERE "id" = 'player_id1'`,
				).Scan(&id, &userId)
				assert.NoError(t, err)
				assert.Equal(t, "player_id1", id)
				assert.Equal(t, "user_id1", userId)
				return true
			},
			errorExpected: false,
			errorString:   "",
		},
		{
			name:   "errors and does not update anything, if Player ID already exists in DB",
			output: "",
			setupSqlStmts: []TestSqlStmts{
				{Query: `INSERT INTO public."players" ("id") VALUES ('id1')`},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."players" WHERE id = 'player_id1'`},
			},
			idGenerator: &utilities.IdGeneratorMockConstant{Id: "id1"},
			dbUpdateCheck: func(db *sql.DB) bool {
				var id string
				err := db.QueryRow(
					`SELECT id FROM public."players" WHERE "id" = 'id1'`,
				).Scan(&id)
				assert.NoError(t, err)
				assert.Equal(t, "id1", id)
				return true
			},
			errorExpected: true,
			errorString:   "pq: duplicate key value violates unique constraint \"players_pkey\"",
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

			playerId, err := s.CreatePlayer(tt.input)
			if !tt.errorExpected {
				assert.NoError(t, err)
				assert.Equal(t, tt.output, playerId)
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

func Test_UpdatePlayerWithUserIdUsingTransaction(t *testing.T) {
	tests := []struct {
		name  string
		input struct {
			playerId string
			userId   string
		}
		setupSqlStmts   []TestSqlStmts
		cleanupSqlStmts []TestSqlStmts
		dbUpdateCheck   func(*sql.DB) bool
		errorExpected   bool
		errorString     string
	}{
		{
			name: "errors if userId is blank",
			input: struct {
				playerId string
				userId   string
			}{
				playerId: "player_id1",
				userId:   "",
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   true,
			errorString:     "userId cannot be blank",
		},
		{
			name: "errors if playerId is blank",
			input: struct {
				playerId string
				userId   string
			}{
				playerId: "",
				userId:   "",
			},
			setupSqlStmts:   nil,
			cleanupSqlStmts: nil,
			dbUpdateCheck:   nil,
			errorExpected:   true,
			errorString:     "playerId cannot be blank",
		},
		{
			name: "errors if db update errors",
			input: struct {
				playerId string
				userId   string
			}{
				playerId: "player_id1",
				userId:   "user_id1",
			},
			setupSqlStmts: []TestSqlStmts{
				{Query: `INSERT INTO public."players" ("id") VALUES ('player_id1')`},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."players" WHERE id = 'player_id1'`},
			},
			dbUpdateCheck: nil,
			errorExpected: true,
			errorString:   "THIS IS BAD: dbError while attempting player update: pq: insert or update on table \"players\" violates foreign key constraint \"players_user_id_fkey\"",
		},
		{
			name: "updates player successfully with user id",
			input: struct {
				playerId string
				userId   string
			}{
				playerId: "player_id1",
				userId:   "user_id1",
			},
			setupSqlStmts: []TestSqlStmts{
				{Query: `INSERT INTO public."users" ("id") VALUES ('user_id1')`},
				{Query: `INSERT INTO public."players" ("id") VALUES ('player_id1')`},
			},
			cleanupSqlStmts: []TestSqlStmts{
				{Query: `DELETE FROM public."users" WHERE id = 'user_id1'`},
			},
			dbUpdateCheck: func(db *sql.DB) bool {
				var (
					userId string
				)
				err := db.QueryRow(
					`SELECT "user_id" FROM public."players" WHERE "id" = 'player_id1'`,
				).Scan(&userId)
				assert.NoError(t, err)
				assert.Equal(t, "user_id1", userId)
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
					Db: testDb,
				},
			)

			runSqlOnDb(t, s.db, tt.setupSqlStmts)
			defer runSqlOnDb(t, s.db, tt.cleanupSqlStmts)

			tx, err := s.BeginTransaction()
			assert.NoError(t, err)
			err = s.UpdatePlayerWithUserIdUsingTransaction(tt.input.playerId, tt.input.userId, tx)
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
