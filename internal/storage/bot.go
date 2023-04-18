package storage

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type BotAccessor interface {
	UpdateBotWithPlayerIdUsingTransaction(botId, playerId string, transaction DatabaseTransaction) error
	UpdateBotDecrementHelpCountUsingTransaction(botId string, transaction DatabaseTransaction) error
}

func (s *Storage) UpdateBotWithPlayerIdUsingTransaction(botId, playerId string, transaction DatabaseTransaction) error {
	return connectPlayerToBot(transaction, playerId, botId)
}

func (s *Storage) UpdateBotDecrementHelpCountUsingTransaction(botId string, transaction DatabaseTransaction) error {
	return decrementHelpCount(transaction, botId)
}

func connectPlayerToBot(customDb customDbHandler, playerId, botId string) error {
	if utilities.IsBlank(botId) {
		return errors.New("botId cannot be blank")
	}

	if utilities.IsBlank(playerId) {
		return errors.New("playerId cannot be blank")
	}

	result, err := customDb.Exec(
		`UPDATE public."bots" SET "player_id" = $1, "type" = 'HUMAN', "help_count" = 3 WHERE id = $2`, playerId, botId,
	)
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while connecting player to bot: %s %s", playerId, botId))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while checking affected row after connecting player to bot: %s %s", playerId, botId))
	}

	if rowsAffected != 1 {
		return utilities.NewBadError("No rows were affected when player was connected to Bot. This is highly unexpected.")
	}

	return nil
}

func decrementHelpCount(customDb customDbHandler, botId string) error {
	if utilities.IsBlank(botId) {
		return errors.New("botId cannot be blank")
	}

	var helpCount int
	row := customDb.QueryRow(`SELECT b.help_count FROM public."bots" AS b WHERE b.id = $1 FOR UPDATE OF b`, botId)
	err := row.Scan(&helpCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.Errorf("getting help_count for %s: no such bot", botId)
		}
		return errors.Errorf("getting help_count for %s: %v", botId, err)
	}

	if helpCount <= 0 {
		return errors.Errorf("help_count should not be updated below 0 for %s", botId)
	}

	result, err := customDb.Exec(`UPDATE public."bots" SET "help_count" = $2 WHERE id = $1`, botId, helpCount-1)
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while decrementing bot help count: %s", botId))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utilities.WrapBadError(err, fmt.Sprintf("dbError while checking affected row while decrementing bot help count: %s", botId))
	}

	if rowsAffected != 1 {
		return utilities.NewBadError("No rows were affected while decrementing bot help count. This is highly unexpected.")
	}

	return nil
}
