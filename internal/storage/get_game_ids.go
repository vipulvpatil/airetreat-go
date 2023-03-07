package storage

import (
	"errors"

	"github.com/vipulvpatil/airetreat-go/internal/model"
	"github.com/vipulvpatil/airetreat-go/internal/utilities"
)

func (s *Storage) GetUnhandledGameIdsForState(gameStateString string) ([]string, error) {
	gameState := model.GameState(gameStateString)
	if !gameState.Valid() {
		return nil, errors.New("invalid game state")
	}

	rows, err := s.db.Query(
		`SELECT id
		FROM public."games"
		WHERE state = $1
		AND state_handled = false
		ORDER BY created_at DESC, id DESC
		`, gameState.String(),
	)
	if err != nil {
		return nil, utilities.WrapBadError(err, "error getting unhandled games")
	}
	defer rows.Close()

	gameIds := []string{}

	for rows.Next() {
		var gameId string
		err := rows.Scan(
			&gameId,
		)

		if err != nil {
			return nil, utilities.WrapBadError(err, "failed while scanning rows")
		}
		gameIds = append(gameIds, gameId)
	}

	err = rows.Err()
	if err != nil {
		return nil, utilities.WrapBadError(err, "failed to correctly go through bot rows")
	}
	return gameIds, nil
}
