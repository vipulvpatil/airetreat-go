package storage

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type MessageCreator interface {
	CreateMessage(sourceBotId, targetBotId, text string) error
	CreateMessageUsingTransaction(sourceBotId, targetBotId, text string, transaction DatabaseTransaction) error
}

func (s *Storage) CreateMessage(sourceBotId, targetBotId, text string) error {
	id := s.IdGenerator.Generate()
	return createMessageUsingCustomDbHandler(s.db, id, sourceBotId, targetBotId, text)
}

func (s *Storage) CreateMessageUsingTransaction(sourceBotId, targetBotId, text string, transaction DatabaseTransaction) error {
	id := s.IdGenerator.Generate()
	return createMessageUsingCustomDbHandler(transaction, id, sourceBotId, targetBotId, text)
}

func createMessageUsingCustomDbHandler(customDb customDbHandler, id, sourceBotId, targetBotId, text string) error {
	if utilities.IsBlank(sourceBotId) {
		return errors.New("sourceBotId cannot be blank")
	}

	if utilities.IsBlank(targetBotId) {
		return errors.New("targetBotId cannot be blank")
	}

	if utilities.IsBlank(text) {
		return errors.New("text cannot be blank")
	}

	result, err := customDb.Exec(
		`INSERT INTO public."messages" ("id", "source_bot_id", "target_bot_id", "text") VALUES ($1, $2, $3, $4)`, id, sourceBotId, targetBotId, text,
	)
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while inserting message: %s %s %s", sourceBotId, targetBotId, text))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while inserting message and changing db: %s %s %s", sourceBotId, targetBotId, text))
	}

	if rowsAffected != 1 {
		return utilities.NewBadError(fmt.Sprintf("Very few or too many rows were affected when inserting message in db. This is highly unexpected. rowsAffected: %d", rowsAffected))
	}

	return nil
}
