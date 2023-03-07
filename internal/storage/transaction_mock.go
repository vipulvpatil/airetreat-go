package storage

import (
	"errors"
)

type databaseTransactionMock struct {
	customDbHandler
}

func (s databaseTransactionMock) Commit() error {
	return nil
}

func (s databaseTransactionMock) Rollback() error {
	return nil
}

type DatabaseTransactionProviderMockSuccess struct{}

func (s *DatabaseTransactionProviderMockSuccess) BeginTransaction() (DatabaseTransaction, error) {
	return &databaseTransactionMock{}, nil
}

type DatabaseTransactionProviderMockFailure struct{}

func (s *DatabaseTransactionProviderMockFailure) BeginTransaction() (DatabaseTransaction, error) {
	return nil, errors.New("unable to begin a db transaction")
}
