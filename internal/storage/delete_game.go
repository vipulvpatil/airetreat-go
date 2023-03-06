package storage

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

func (s *Storage) DeleteGame(gameId string) error {
	if utilities.IsBlank(gameId) {
		return errors.New("gameId cannot be blank")
	}

	result, err := s.db.Exec(`DELETE FROM public."games" WHERE id = $1`, gameId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, "dbError while deleting game and changing db")
	}

	if rowsAffected != 1 {
		return utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when deleting game in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
	}
	return nil
}
