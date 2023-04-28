package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/vipulvpatil/airetreat-go/internal/config"
	"github.com/vipulvpatil/airetreat-go/internal/utilities/logger"
)

// This function will make a connection to the database only once.
func InitDb(cfg *config.Config) (*sql.DB, error) {
	var err error

	connStr := cfg.DbUrl
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	// this will be printed in the terminal, confirming the connection to the database
	logger.LogMessageln("The database is connected")
	return db, nil
}
