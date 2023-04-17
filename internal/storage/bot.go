package storage

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

type BotAccessor interface {
	UpdateBotWithPlayerIdUsingTransaction(botId, playerId string, transaction DatabaseTransaction) error
}

func (s *Storage) UpdateBotWithPlayerIdUsingTransaction(botId, playerId string, transaction DatabaseTransaction) error {
	return connectPlayerToBot(transaction, playerId, botId)
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
