package storage

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type MessageCreator interface {
	CreateMessage(botId, text string) error
}

func (s *Storage) CreateMessage(botId, text string) error {
	if utilities.IsBlank(botId) {
		return errors.New("botId cannot be blank")
	}

	if utilities.IsBlank(text) {
		return errors.New("text cannot be blank")
	}

	id := s.IdGenerator.Generate()

	result, err := s.db.Exec(
		`INSERT INTO public."messages" ("id", "bot_id", "text") VALUES ($1, $2, $3)`, id, botId, text,
	)
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while inserting message: %s %s", botId, text))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while inserting message and changing db: %s %s", botId, text))
	}

	if rowsAffected != 1 {
		return utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when inserting message in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
	}

	return nil
}
