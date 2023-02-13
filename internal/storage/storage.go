package storage

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type StorageAccessor interface {
	UserRetriever
}

type Storage struct {
	db          *sql.DB
	IdGenerator utilities.CuidGenerator
}

type StorageOptions struct {
	Db          *sql.DB
	IdGenerator utilities.CuidGenerator
}

func NewDbStorage(opts StorageOptions) (*Storage, error) {
	if opts.Db == nil {
		return nil, errors.New("Needs a backing database")
	}

	if opts.IdGenerator == nil {
		opts.IdGenerator = &utilities.RandomIdGenerator{}
	}

	return &Storage{
		db:          opts.Db,
		IdGenerator: opts.IdGenerator,
	}, nil
}
