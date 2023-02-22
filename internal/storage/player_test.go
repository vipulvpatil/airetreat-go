package storage

import (
	"database/sql"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

func Test_CreatePlayer(t *testing.T) {
	tests := []struct {
		name            string
		output          string
		setupSqlStmts   []string
		cleanupSqlStmts []string
		idGenerator     utilities.CuidGenerator
		dbUpdateCheck   func(*sql.DB) bool
		errorExpected   bool
		errorString     string
	}{
		{
			name:          "creates player successfully",
			output:        "player_id1",
			setupSqlStmts: []string{},
			cleanupSqlStmts: []string{
				`DELETE FROM public."players" WHERE id = 'player_id1'`,
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
			name:   "errors and does not update anything, if Player ID already exists in DB",
			output: "",
			setupSqlStmts: []string{
				`INSERT INTO public."players" ("id") VALUES ('id1')`,
			},
			cleanupSqlStmts: []string{
				`DELETE FROM public."games" WHERE id = 'id1'`,
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

			rand.Seed(0)
			playerId, err := s.CreatePlayer()
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
