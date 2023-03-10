package storage

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/airetreat-go/internal/config"
)

var testDb *sql.DB

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	cfg, _ := config.NewConfigFromEnvVars()

	testDb, err = openTestDatabaseConnection(cfg)
	if err != nil {
		return -1, fmt.Errorf("could not setup test database: %w", err)
	}
	resetTestDatabase(testDb)
	err = populateTableSchemaIntoTestDatabase(testDb)
	if err != nil {
		return -1, fmt.Errorf("could not populate test database: %w", err)
	}

	defer func() {
		resetTestDatabase(testDb)
		closeTestDatabaseConnection(testDb)
	}()

	return m.Run(), nil
}

func openTestDatabaseConnection(cfg *config.Config) (*sql.DB, error) {
	connStr := cfg.TestDbUrl
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func populateTableSchemaIntoTestDatabase(db *sql.DB) error {
	schemaCreationSql, err := os.ReadFile("./database_schema_test.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(schemaCreationSql))

	if err != nil {
		return err
	}

	fmt.Println("DB schema is setup for testing")
	return nil
}

func resetTestDatabase(db *sql.DB) {
	deleteAllTablesSql := "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO public;"
	_, err := db.Exec(string(deleteAllTablesSql))
	if err != nil {
		fmt.Println("Unable to clean up test db")
	}
}

func closeTestDatabaseConnection(db *sql.DB) {
	db.Close()
}

func Test_InitDb(t *testing.T) {
	t.Run("Test that DB connectivity works", func(t *testing.T) {
		cfg, _ := config.NewConfigFromEnvVars()
		storage, err := InitDb(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, storage)
	})
}

func runSqlOnDb(t *testing.T, db *sql.DB, sqlStmts []string) {
	for _, sqlStmts := range sqlStmts {
		_, err := db.Exec(sqlStmts)
		assert.NoError(t, err)
	}
}
