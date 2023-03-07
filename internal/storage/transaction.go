package storage

import "database/sql"

type TxProvider interface {
	GetTx() (*sql.Tx, error)
}

func (s *Storage) GetTx() (*sql.Tx, error) {
	return s.db.Begin()
}
