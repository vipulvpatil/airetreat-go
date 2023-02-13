package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/vipulvpatil/airetreat-go/internal/config"
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
	fmt.Println("The database is connected")

	// TODO: Delete once database connectivity can be verified by other means.
	rows, err := db.Query(`SELECT id, email FROM public."User"`)
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()
	for rows.Next() {
		var id string
		var email string

		rows.Scan(&id, &email)

		fmt.Println(id, email)
	}
	fmt.Println("finished printing result")

	return db, nil
}
